# switchmanager
`switchmanager` is a set of tools that let a system administrator to manage a x86-based layer-2 switch in a distributed environment. It is composed by a manager daemon and by an agent daemon:
- `managerd`:
- `agent`:

## Install
All the component in this repository are written in [Go](https://golang.org) (version >= 1.7), so a Go distribution must be installed on your system in order to build the source code.
Go distributions can be downloaded from https://golang.org/dl/ or by using the package management system that comes with your Linux distribution. 

### Dependecies installation
Once Go has been installed, some dependencies are needed:
```sh
go get -u github.com/gorilla/mux
go get menteslibres.net/gosexy/yaml
```
