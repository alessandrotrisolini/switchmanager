# `agentd` code structure

The code structure of `agentd` is pretty straightforward. Here is reported the structure of this directory and its subfolders:

```
.
├── agent
│   └── agent.go
├── agentd.go
├── agentserver
│   └── agentserver.go
├── config
│   └── config.go
└── managerapi
    ├── apiwrapper.go
    └── managerapi.go
```

The entry point of `agentd` is `agentd.go`: it reads the configuration file, initialises both the client and the server of `agentd`, and registers `agentd` to the manager server.

All the communications are performed via REST API. `agentd` exposes a REST API to the manager and it is defined in `agentserver/agentserver.go`:

- `POST /processes`: instantiates a new process on the agent machine;
- `GET /processes`: dumps all the processes that are currently running on the agent machine;
- `DELETE /processes/{pid}`: kills a process given its PID;
- `GET /alive`: it is used by the manager to check if an `agentd` is alive.

All the actions that `agentd` performs on the system are defined in `agent/agent.go` (i.e. create, kill, and get the status of a process).
