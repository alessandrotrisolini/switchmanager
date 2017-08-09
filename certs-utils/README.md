In this directory you can find a set of script that are useful for building a simple Certificate Authority (CA) and    for creating new certificates.

# Directory structure
There are three directories:

- `ca` directory: it will contain all the files related to the CA;
- `manager-certs` directory: output directory when creating a manager certificate;
- `agent-certs` directory: output directory when creating an agent certificate.

Each one of these folders contain a `.cnf` file, which is fed to `openssl` as configuration file.

# Scripts
Three scripts are available as tools for creating new certificates:

- `generate-ca.sh` script: creates a new CA based on the content of `ca/ca.cnf` configuration file;
- `generate-manager-cert.sh` script: creates a new manager certificate based on the content of `manager-certs/manager.cnf` configuration file;
- `generate-agent-cert.sh` script: creates a new agent certificate based on the content of `agent-certs/agent.cnf` configuration file.

# Usage

Before start, make sure that all the scripts can be executed. If not:

```bash
$ chmod +x *.sh
```

**Beware**: all the scripts clear the content of the relative directory every time are they executed (i.e. `generate-ca.sh` clears `ca` directory content).

## Certificate Authority
First of all, the CA **must** be created:

```bash
$ ./generate-ca.sh
```

You have to follow the `openssl` prompt and set the private key password when requested. After the script finished its execution, the `ca` directory contains all the files needed by the CA to correctly sign certificates.

**Beware**: the file `ca.cnf` in `ca` directory contains the configuration of the CA certificate. If you want to modify some parameter you have to edit that file and launch the script again: it will erase all the content of `ca` directory and create a CA accondingly with `ca.cnf` content.

## Manager certificate creation
Before creating a new manager certificate, you might want to modify the `manager.cnf` file contained in `manager-certs` directory, especially the `commonName`:  

```
[...]

[ server ]
countryName            = IT
stateOrProvinceName    = Italy
localityName           = Turin
organizationName       = Torsec
commonName             = manager.torsec.it        <---  MODIFY ME, PLEASE

[ v3_server ]
basicConstraints       = CA:FALSE
keyUsage               = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage       = clientAuth, serverAuth
```

Now the private key and the certificate can be created:

```bash
$ ./generate-manager-cert.sh
```

Now you copy `manager.pem` and `manager.key` files to another path and use it within the `managercli` application.


## Agent certificate creation
As already noticed in the manager certificate creation section, `agent.cnf` file, contained in `agent-certs` directory, has to be modified before launching the script. Once it has been modified, launch:

```bash
$ ./generate-agent-cert.sh
```

Now you copy `agent.pem` and `agent.key` files to another path and use it within the `agentd` application.

