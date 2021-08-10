.PHONY: test build init create delete stop

test:
	go test $(go list ./...) -v

build:
	docker build -t ghcr.io/$(REPO)/stripe-api:latest .
	docker push ghcr.io/$(REPO)/stripe-api:latest

init:
	minikube config set vm-driver hyperkit
	minikube start
	minikube addons enable ingress
	kubectl create secret docker-registry regcred --docker-server=https://ghcr.io/v2/ --docker-username=$(REPO) --docker-password=$(CR_PAT)
	echo "$(minikube ip) hello.world" | sudo tee -a /etc/hosts

create:
	kubectl apply -f ./kubernetes/secret.yml
	kubectl apply -f ./kubernetes/minikube-ingress.yml
	kubectl apply -f ./kubernetes/api-demo-deployment.yml
	kubectl apply -f ./kubernetes/api-demo-service.yml

delete:
	kubectl delete -f ./kubernetes/secret.yml
	kubectl delete -f ./kubernetes/minikube-ingress.yml
	kubectl delete -f ./kubernetes/api-demo-deployment.yml
	kubectl delete -f ./kubernetes/api-demo-service.yml

stop:
	minikube stop