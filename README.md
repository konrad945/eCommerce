# eCommerce [WIP]

eCommerce will be a cloud-native microservice application. Let's start!

### Current functionality - already implemented
#### Services
* catalog - svc containing basic information about all products

#### DBs
* PostgreSQL (Relational DB...)

#### More...
* Logging (EFK stack)
* Tracing (OpenTelemetry)
* Vault (Secrets, etc...)

#### What's next to do...
* Setup db and svc to use vault, move svc to separate namespace 

## Prerequisites

For the full-blown experience <b>docker, kind</b> and <b> kubectl</b> should be installed.

### Kubectl setup
For linux systems with snap package manager use command for kubectl installation:
```shell
sudo snap install kubectl --classic
```

### Kind setup
Kind setup is as easy as downloading binary with go install command
```shell
go install sigs.k8s.io/kind@latest
```
If after installation kind command is not available, you should check if GOBIN is 
added to PATH. On ubuntu you can add those two lines to your .bashrc
```
GOBIN=$HOME/go/bin
export PATH=$PATH:$GOBIN
```

## Quick Start
To quickly bootstrap application execute ```make start-all ```

It will:
  * Create kind cluster
  * Create local docker registry
  * Install needed backend services (DB, Jaeger, EFK, (more will come)... )
  * Build, load and install apps (catalog, (more will come)... ) 

To delete k8s cluster and docker registry execute ```make clean-all```

As for now there are no external endpoints (ingress, api gateway) created thus for access to specific 
services port forwarding is required. ```kubectl port-forward svc/<svc-name> <local-port>:<svc-port>```
