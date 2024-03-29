package agentserver

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	a "switchmanager/agentd/agent"
	cmn "switchmanager/common"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"
)

// AgentServer implements the server functionalities of the
// agent daemon
type AgentServer struct {
	agent  *a.Agent
	log    *l.Log
	router *mux.Router
	server *http.Server

	certPath string
	keyPath  string
}

// NewAgentServer initializes the agent server
func NewAgentServer(certPath string, keyPath string, caCertPath string) (*AgentServer, error) {
	agent := a.NewAgent()
	log := l.GetLogger()
	router := mux.NewRouter()
	server := &http.Server{Handler: router}

	err := cmn.SetupTLSServer(server, caCertPath)
	if err != nil {
		return nil, err
	}

	as := &AgentServer{
		agent:    agent,
		router:   router,
		server:   server,
		certPath: certPath,
		keyPath:  keyPath,
		log:      log,
	}

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	})
	router.Handle("/processes", doRun(as)).Methods("POST")
	router.Handle("/processes/{pid}", doKill(as)).Methods("DELETE")
	router.Handle("/processes", doDump(as)).Methods("GET")
	router.Handle("/alive", doAlive(as)).Methods("GET")

	ticker := time.NewTicker(time.Second).C
	go func() {
		for {
			<-ticker
			agent.RefreshProcesses()
		}
	}()

	return as, nil
}

// Start starts the agent server
func (as *AgentServer) Start(port string) {
	as.server.Addr = ":" + port
	err := as.server.ListenAndServeTLS(as.certPath, as.keyPath)
	if err != nil {
		as.log.Error(err)
		os.Exit(1)
	}
}

func doAlive(as *AgentServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func doRun(as *AgentServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var hostapdConfig dm.HostapdConfig
		_ = json.NewDecoder(req.Body).Decode(&hostapdConfig)
		as.log.Info(hostapdConfig)
		configFile := newHostapdConfigFile(as, hostapdConfig)
		cmd := exec.Command("hostapd", configFile, "-z", hostapdConfig.OpenvSwitch)
		// cmd := exec.Command("./foo")

		err := cmd.Start()
		if err != nil {
			as.log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			pid := cmd.Process.Pid
			as.log.Info("Process started - PID:", pid)

			state, err := cmn.GetProcessState(pid)
			if err != nil {
				as.log.Error(err)
			} else {
				p := &a.Process{
					Process: cmd.Process,
					State:   state,
				}
				as.agent.AddProcess(pid, p)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(dm.ProcessDescriptor{Pid: pid, State: state})
			}
		}
	})
}

func doKill(as *AgentServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		pid, err := strconv.ParseUint(vars["pid"], 10, 32)
		if err == nil {
			as.log.Info("Trying to kill process with PID:", pid)
			if as.agent.CheckPid(int(pid)) {
				err := as.agent.DeleteProcess(int(pid))
				if err != nil {
					as.log.Error("Cannot stop process with PID:", int(pid))
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					as.log.Info("Process killed!")
					w.WriteHeader(http.StatusOK)
				}
			} else {
				as.log.Error("Process with PID", pid, "does not exist")
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			as.log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func doDump(as *AgentServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(as.agent.DumpProcesses())
	})
}

func newHostapdConfigFile(as *AgentServer, hostapdConfig dm.HostapdConfig) string {
	var buffer bytes.Buffer

	buffer.WriteString("interface=" + hostapdConfig.Interface + "\n")
	buffer.WriteString("driver=macsec_linux\n")

	buffer.WriteString("ieee8021x=1\n")
	buffer.WriteString("eap_reauth_period=" + strconv.FormatUint(hostapdConfig.ReauthTimeout, 10) + "\n")
	buffer.WriteString("use_pae_group_addr=1\n")

	buffer.WriteString("own_ip_addr=10.1.1.4\n")
	buffer.WriteString("radius_client_addr=10.1.1.4\n")

	buffer.WriteString("auth_server_addr=" + hostapdConfig.RadiusAuthServer + "\n")
	buffer.WriteString("auth_server_port=1812\n")
	buffer.WriteString("auth_server_shared_secret=" + hostapdConfig.RadiusSecret + "\n")

	buffer.WriteString("acct_server_addr=" + hostapdConfig.RadiusAcctServer + "\n")
	buffer.WriteString("acct_server_port=1813\n")
	buffer.WriteString("acct_server_shared_secret=" + hostapdConfig.RadiusSecret + "\n")

	f, _ := ioutil.TempFile("", "/tmp")
	as.log.Info(f.Name())
	_ = ioutil.WriteFile(f.Name(), buffer.Bytes(), 0644)
	return f.Name()
}
