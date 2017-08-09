#!/bin/bash

AGENTDIR="agent-certs"

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
# generic helpers
# -----------------------------------------------------------------------------
function cleanup {
    find $AGENTDIR -type f ! -name "*.cnf" -delete
}

function check_exit {
    if [ $? -ne 0 ]; then
        _fail "$1"
        {
        cleanup
        } > /dev/null 2>&1
        exit 1
    fi
}

# -----------------------------------------------------------------------------
# main
# -----------------------------------------------------------------------------

{
cleanup
} > /dev/null 2>&1

if [ ! -f ca/ca.pem ]; then
    _fail "Certificate authority has not been created"
    exit 1
fi

_info "Generating agent private key" 
openssl genrsa -aes256 -out $AGENTDIR/agent.key.encrypted
check_exit "Error while generating agent private key"

_info "Generating agent certificate signing request"
openssl req -config $AGENTDIR/agent.cnf \
            -key $AGENTDIR/agent.key.encrypted \
            -new -sha256 -out $AGENTDIR/agent.csr
check_exit "Error while generating agent certificate signing request"

_info "Generating agent certificate"
openssl ca -config $AGENTDIR/agent.cnf \
           -md sha256 -notext \
           -in $AGENTDIR/agent.csr \
           -out $AGENTDIR/agent.pem
check_exit "Error while generating agent certificate"

_info "Generating agent plain private key"
openssl rsa -in $AGENTDIR/agent.key.encrypted \
            -out $AGENTDIR/agent.key
check_exit "Error while generating agent plain private key"

_info "Done"

