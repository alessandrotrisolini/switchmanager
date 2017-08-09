# Manual installation
In this document it will be explained the manual installation process of all the building blocks needed by `managercli` and `agentd`.

## Go
All the component in this repository are written in [Go](https://golang.org) (version >= 1.7), so a Go distribution must be installed on your system in order to build the source code.
Go distributions can be downloaded from the official [page](https://golang.org/dl/) or by using the package management system that comes with your Linux distribution. If you install Go by downloading it, there is an installation [guide](https://golang.org/doc/install) available.

## Prepare Go environment
Before starting to build and install Go applications, `GOPATH` environment variable must be set. It contains the root path where all the Go toolchain refer to when building, running, or installing a Go application.

```sh
$ export GOPATH=$HOME/go    # we suppose to use $HOME/go as path, but it can be different
```
**Beware**: you have to set `GOPATH` every time you open a new shell. You might want to add the above command at the bottom of `$HOME/.profile` to automate the `GOPATH` assignment. 

## Dependecies installation
Once Go has been installed, some dependencies are needed:
```sh
$ go get -u github.com/gorilla/mux
$ go get -u github.com/spf13/viper
$ go get -u github.com/mitchellh/cli
$ go get -u github.com/vishvananda/netlink
$ go get -u github.com/chzyer/readline
```

