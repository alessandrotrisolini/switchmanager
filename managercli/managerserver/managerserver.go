package managerserver

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	cmn "switchmanager/common"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"
	m "switchmanager/managercli/manager"
)

// Port where the Manager server exposes the service
const port string = ":5000"

//ManagerServer is the data type that models the server
type ManagerServer struct {
	manager  *m.Manager
	router   *mux.Router
	server   *http.Server
	certPath string
	keyPath  string
	log      *l.Log
}

func doRegister(ms *ManagerServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var conf dm.AgentConfig
		_ = json.NewDecoder(req.Body).Decode(&conf)
		ms.manager.RegisterAgent(conf)
		ms.log.Trace("Registered agent with config:", conf)
		w.WriteHeader(http.StatusOK)
	})
}

// NewManagerServer initializes the manager server
func NewManagerServer(manager *m.Manager, certPath, keyPath, caCertPath string) (*ManagerServer, error) {
	log := l.GetLogger()
	router := mux.NewRouter()
	server := &http.Server{Addr: port}
	server.Handler = router

	err := cmn.SetupTLSServer(server, caCertPath)
	if err != nil {
		return nil, err
	}

	ms := &ManagerServer{
		manager:  manager,
		router:   router,
		server:   server,
		certPath: certPath,
		keyPath:  keyPath,
		log:      log,
	}

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	})

	router.Handle("/agents", doRegister(ms)).Methods("POST")

	return ms, nil
}

// Start starts the server
func (ms *ManagerServer) Start() {
	ms.log.Trace("Starting manager server...")
	err := ms.server.ListenAndServeTLS(ms.certPath, ms.keyPath)
	if err != nil {
		ms.log.Error("Cannot start the server:", err)
		os.Exit(1)
	}
}

