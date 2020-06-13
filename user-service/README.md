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

// delete 
kubectl delete deployment -l app=devsmd-user-service
kubectl delete services -l app=devsmd-user-service

// restart 
kubectl rollout restart deploy devsmd-user-service

// Push to docker hub
docker build -t devsmd/user-service .
docker tag devsmd/user-service devsmd/user-service:1.3.2
docker push devsmd/user-service:1.3.2