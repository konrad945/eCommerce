# eCommerce [WIP]

eCommerce will be a cloud-native microservice application. Let's start!

## Prerequisites

For the full-blown experience <b>kind, kubectl</b> and <b>helm</b> should be installed.

### Kubectl setup
For linux systems with snap package manager use command for kubectl installation:
```shell
sudo snap install kubectl --classic
```

### Helm setup
For linux systems with snap package manager use command for helm installation:
```shell
sudo snap install helm --classic
```
After successful helm installation let's add some public repositories
```shell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
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
To quickly bootstrap application execute ```make start-all ```. It will:
  * Create kind cluster
  * Create needed namespaces
  * Install needed backend services (DB, Jaeger, (more will come)... )
  * Build, load and install apps (catalog, (more will come)... ) 

As for now there are no external endpoints (ingress, api gateway) created thus for access to specific 
services port forwarding is required. ```kubectl port-forward svc/<svc-name> <local-port>:<svc-port>```
