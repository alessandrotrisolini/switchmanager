#!/bin/bash

CADIR="ca"

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
    find $CADIR -type f ! -name "*.cnf" -delete    
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

touch $CADIR/index.txt
echo 1000 > $CADIR/serial

_info "Generating CA private key"
openssl genrsa -aes256 -out $CADIR/ca.key 4096
check_exit "Error while generating CA private key"

_info "Generating CA certificate"
openssl req -config $CADIR/ca.cnf \
            -key $CADIR/ca.key \
            -new -x509 -days 3650 -sha256 -extensions v3_ca \
            -out $CADIR/ca.pem
check_exit "Error while generating CA certificate"

_info "Done"

