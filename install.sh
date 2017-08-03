#!/bin/bash

RC=0

# -----------------------------------------------------------------------------
# logging helpers
# -----------------------------------------------------------------------------

function _log {
	level=$1
	msg=$2

	case "$level" in
		info)
			tag="\e[1;36minfo\e[0m"
			;;
		err)
			tag="\e[1;31merr \e[0m"
			;;
		warn)
			tag="\e[1;33mwarn\e[0m"
			;;
		ok)
			tag="\e[1;32m ok \e[0m"
			;;
		fail)
			tag="\e[1;31mfail\e[0m"
			;;
		*)
			tag="    "
			;;
	esac
	echo -e "`date +%Y-%m-%dT%H:%M:%S` [$tag] $msg"
}

function _err {
	msg=$1
	_log "err" "$msg"
}

function _warn {
	msg=$1
	_log "warn" "$msg"
}

function _info {
	msg=$1
	_log "info" "$msg"
}

function _success {
	msg=$1
	_log "ok" "$msg"
}

function _fail {
	msg=$1
	_log "fail" "$msg"
}

# -----------------------------------------------------------------------------
# Checking dependencies
# -----------------------------------------------------------------------------
function check_dependencies {
	IFS=':' read -a PATHS <<< "$PATH"
	DEPS=("tar" "wget")
	for d in ${DEPS[@]}; do
		_info "Looking for '$d'"
		found=0
		for p in ${PATHS[@]}; do
			if [ -f "$p/$d" ]; then
				_success "Found '$d' in $p/$d"
				found=1
			fi
		done
		if [ $found -eq 0 ]; then
			_fail "Missing dependency: $d"
			case $d in
				"tar")
					_info "On Debian/Ubuntu: apt-get install tar"
					_info "On RedHat/CentOS/Fedora: yum install tar"
					;;
				"wget")
					_info "On Debian/Ubuntu: apt-get install wget"
					_info "On RedHat/CentOS/Fedora: yum install wget"
					;;
				*)
					_info "No suggestions available. Good luck... :-)"
					;;
			esac
			exit 1
		fi
	done
}

function check_bin {
	IFS=':' read -a PATHS <<< "$PATH"
	BIN=$1
	_info "Looking for '$BIN'"
	found=0
	for p in ${PATHS[@]}; do
		if [ -f "$p/$BIN" ]; then
			_success "Found '$BIN' in $p/$BIN"
			found=1
		fi
	done
	if [ $found -eq 0 ]; then
		_fail "Missing binary: $BIN"
		RC=1
	else
		RC=0
	fi
}


function check_go {
	check_bin go
}

function check_go_version {
	_info "Checking Go version..."
	v=`go version | sed -e "s/.*go\([0-9]*\.[0-9]*\).*/\1/g"`
	if [ "$v" < "1.7" ]; then
	        RC=1
	else
		RC=0
	fi
}

function install_go {
	_info "Installing Go..."
	if [ ! -d "$HOME/go" ]; then
        mkdir ~/go
    fi
	wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
	sudo tar -C /usr/local -xzf go1.8.3.linux-amd64.tar.gz
	grep -q "/usr/local/go/bin" ~/.profile || \
            echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile && \
            source ~/.profile
        grep -q "export GOPATH=~/go" ~/.profile || \
            echo "export GOPATH=~/go" >> ~/.profile && \
            source ~/.profile
	rm go1.8.3.linux-amd64.tar.gz
	check_go
}

# -----------------------------------------------------------------------------
# main
# -----------------------------------------------------------------------------

if [ "$0" == "$BASH_SOURCE" ]; then
    _fail "Please source the script -> source install.sh"
    exit 1
fi

check_dependencies
check_go

if [ $RC -eq 1 ]; then
	install_go
	if [ $RC -eq 0 ]; then
		_info "Go has been correctly installed"
	else
		_fail "Go has not been correctly installed"
		exit 1
	fi
else
	check_go_version
	if [ $RC -eq 1 ]; then
		_warn "Your Go version is older than 1.7"
		_warn "Uninstall it and re-run this script"
		exit 0
	else
		_info "Go is already installed on your system"
	fi
fi

_info "Installing gorilla/mux"
go get -v -u github.com/gorilla/mux
_info "Installing spf13/viper"
go get -v -u github.com/spf13/viper
_info "Installing mitchellh/cli"
go get -v -u github.com/mitchellh/cli
_info "Installing vishvananda/netlink"
go get -v -u github.com/vishvananda/netlink
_info "Installing vishvananda/netlink"
_info "Installing chzyer/readline"
go get -v -u github.com/chzyer/readline

 _info "Copying switchmanager folder to $GOPATH/src/switchmanager"
 if [ ! -d "$GOPATH/src/switchmanager" ]; then
     mkdir $GOPATH/src/switchmanager
 fi
 cp -r . $GOPATH/src/switchmanager

_info "Setup finished" 
