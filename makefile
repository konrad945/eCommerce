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