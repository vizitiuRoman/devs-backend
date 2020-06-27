// Apply folder

kubectl apply -f kubernetes/

// Apply single

kubectl create -f user-middleware-secret.yaml

kubectl get secrets

kubectl apply -f user-middleware-deployment.yaml

kubectl apply -f user-middleware.yaml

// delete 

kubectl delete deployment -l app=devsmd-middleware-service

kubectl delete services -l app=devsmd-middleware-service

kubectl delete deployment --all

kubectl delete secrets --all

kubectl delete services --all

// restart 

kubectl rollout restart deploy devsmd-middleware-service

// Push to docker hub

docker build -t devsmd/middleware-service .

docker tag devsmd/user-middleware devsmd/user-service:1.3.2

docker push devsmd/middleware-service:1.3.2