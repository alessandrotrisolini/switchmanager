# switchmanager
`switchmanager` is a set of tools that let a system administrator to manage a x86-based layer-2 switch in a distributed environment. It is composed by a manager daemon and by an agent daemon:
- `managerd`: is the daemon that runs into the management server. It is in charge of monitoring all the agent that are deployed in a network;
- `agent`: is the daemon that runs into each x86-based layer-2 switch.

## Architecture

## Install

### Go
All the component in this repository are written in [Go](https://golang.org) (version >= 1.7), so a Go distribution must be installed on your system in order to build the source code.
Go distributions can be downloaded from the official [page](https://golang.org/dl/) or by using the package management system that comes with your Linux distribution. If you install Go by downloading it, there is an installation [guide](https://golang.org/doc/install) available.

#### Dependecies installation
Once Go has been installed, some dependencies are needed:
```sh
go get -u github.com/gorilla/mux
go get menteslibres.net/gosexy/yaml
```
### Open vSwitch
[Open vSwitch](http://openvswitch.org/) is the software switch that is used as core switching engine inside each switch machine. 
It can be installed from the main project GitHub repository by following the installation [guide](https://github.com/openvswitch/ovs/blob/master/Documentation/intro/install/general.rst) or by using the package management system.
### wpa_supplicant/hostapd
TODO
