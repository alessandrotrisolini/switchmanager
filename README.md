# switchmanager
`switchmanager` is a set of tools that let a system administrator to manage a x64-based layer-2 switch in a distributed environment. It is composed by a manager daemon and by an agent daemon:
- `managerd`: is the daemon that runs into the management server. It is in charge of monitoring all the agent that are deployed in a network;
- `agentd`: is the daemon that runs into each x64-based layer-2 switch.

## Architecture

## Install

### Go
All the component in this repository are written in [Go](https://golang.org) (version >= 1.7), so a Go distribution must be installed on your system in order to build the source code.
Go distributions can be downloaded from the official [page](https://golang.org/dl/) or by using the package management system that comes with your Linux distribution. If you install Go by downloading it, there is an installation [guide](https://golang.org/doc/install) available.

#### Dependecies installation
Once Go has been installed, some dependencies are needed:
```sh
$ go get -u github.com/gorilla/mux
$ go get -u github.com/spf13/viper
```
### Open vSwitch
[Open vSwitch](http://openvswitch.org/) is the software switch that is used as core switching engine inside each switch machine. 
It can be installed from the main project GitHub repository by following the installation [guide](https://github.com/openvswitch/ovs/blob/master/Documentation/intro/install/general.rst) or by using the package management system.

### hostapd
**hostapd** has to be installed on each switch. Every physical port that is supposed to be part of the switch need an instance of **hostapd** to manage both the 802.1X authentication and the MACsec channel generation. 

**agentd** is in charge of run and manage all the life process of **hostapd** intances.

## Usage examples
In this section the usage of both **managerd** and **agentd** will be explained.
### managerd usage
On the manager machine we have to launch the **managerd** and a CLI will appear:
```sh
$ managerd -config /path/to/config
```
The configuration file is a YAML which structure is composed by the following fields:
```sh
manager_cert: "/path/to/manager/pem"
manager_key: "/path/to/manager/key"
ca_cert: "/path/to/ca/pem"
```

Now we can interact with the **managerd** CLI with several commands:
- `list` : lists all the registred agents;
- `run -hostname <agent.hostname>` : runs an instance of **hostapd** on a registered agent;
- `dump -hostname <agent.hostname>` : lists all the instances of **hostapd** of a registered agent;
- `kill -hostname <agent.hostname> -pid <pid>` : kills a specific instance of **hostapd**.

[![demo](https://asciinema.org/a/ydWjUTmYwOJVyk3yguT8fujCh.png)](https://asciinema.org/a/ydWjUTmYwOJVyk3yguT8fujCh?autoplay=1)

### agentd usage
