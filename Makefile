apply:
	kubectl apply -f kubernetes/
	kubectl apply -f microservices/user-service/kubernetes
	kubectl apply -f microservices/middleware-service/kubernetes
	minikube addons enable ingress

