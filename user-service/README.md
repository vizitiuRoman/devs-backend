kubectl create -f pg-secret.yaml
kubectl get secrets
kubectl apply -f pg-db-pv.yaml
kubectl apply -f pg-db-pvc.yaml
kubectl apply -f pg-db-deployment.yaml
kubectl apply -f pg-db-service.yaml
kubectl apply -f user-service-deployment.yaml
kubectl apply -f user-service.yaml
kubectl apply -f redis-master-deployment.yaml
kubectl apply -f redis-master-service.yaml
kubectl apply -f redis-slave-deployment.yaml
kubectl apply -f redis-slave-service.yaml

docker build -t devsmd/user-service .
docker tag devsmd/user-service devsmd/user-service:1.3.2
docker push devsmd/user-service:1.3.2