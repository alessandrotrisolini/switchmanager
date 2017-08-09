#!/bin/bash

MANAGERDIR="manager-certs"

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
    find $MANAGERDIR -type f ! -name "*.cnf" -delete
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

_info "Generating manager private key"
openssl genrsa -aes256 -out $MANAGERDIR/manager.key.encrypted
check_exit "Error while generating manager private key"

_info "Generating manager certificate signing request"
openssl req -config $MANAGERDIR/manager.cnf \
            -key $MANAGERDIR/manager.key.encrypted \
            -new -sha256 -out $MANAGERDIR/manager.csr
check_exit "Error while generating manager certificate signing request"

_info "Generating manager certificate"
openssl ca -config $MANAGERDIR/manager.cnf \
           -md sha256 -notext \
           -in $MANAGERDIR/manager.csr \
           -out $MANAGERDIR/manager.pem
check_exit "Error while generating manager certificate"

_info "Generating manager plain private key"
openssl rsa -in $MANAGERDIR/manager.key.encrypted \
            -out $MANAGERDIR/manager.key
check_exit "Error while generating manager plain private key"

_info "Done"

