kind-start:
	kind create cluster

kind-stop:
	kind delete cluster

k8s-namespace-create:
	kubectl create namespace backend
	kubectl create namespace apps

k8s-postgresql-install:
	helm repo update
	helm install postgre-db -f helm/postgresql/values.yaml bitnami/postgresql -n backend

k8s-jaeger-install:
	helm repo update
	helm install jaeger -f helm/jaeger/values.yaml jaegertracing/jaeger -n backend

k8s-images-load:
	kind load docker-image catalog:latest

service-catalog-build:
	docker build . -t catalog:latest

k8s-catalog-install:
	helm install catalog svc/catalog/helm/ -n backend

start-all: kind-start k8s-namespace-create k8s-postgresql-install k8s-jaeger-install service-catalog-build k8s-images-load k8s-catalog-install

run-tests:
	go test ./svc/...