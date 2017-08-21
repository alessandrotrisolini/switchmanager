# `managercli` code structure

The code structure of `managercli` is:

```
.
├── agentapi
│   ├── agentapi.go
│   └── apiwrapper.go
├── cli
│   └── cli.go
├── config
│   └── config.go
├── manager
│   └── manager.go
├── managercli.go
└── managerserver
    └── managerserver.go
```
 
The entry point of `managercli` is `managercli.go`: it reads the configuration file, initialises both the client and the server of `managercli`, and starts the command line interface.

All the communications are performed via REST API. `managercli` exposes a REST API to the agents and it is defined in `managerserver/managerserver.go`:

- `POST /agents`: register a new `agetnd` to the manager server;

