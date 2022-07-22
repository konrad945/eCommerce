cluster-deploy:
	terraform -chdir=terraform/create-cluster init
	terraform -chdir=terraform/create-cluster apply -auto-approve

cluster-destroy:
	terraform -chdir=terraform/create-cluster destroy -auto-approve

apps-deploy:
	terraform -chdir=terraform/setup-cluster init
	terraform -chdir=terraform/setup-cluster apply -auto-approve

apps-destroy:
	terraform -chdir=terraform/setup-cluster destroy -auto-approve

service-catalog-build:
	docker build . -t localhost:5001/catalog:latest

images-push:
	docker push localhost:5001/catalog:latest

start-all: service-catalog-build cluster-deploy images-push apps-deploy
clean-all: apps-destroy cluster-destroy

run-tests:
	go test ./svc/...