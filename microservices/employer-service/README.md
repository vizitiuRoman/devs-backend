// Apply folder

kubectl apply -f kubernetes/

// Apply single

kubectl create -f employer-service-secret.yaml

kubectl get secrets

kubectl apply -f employer-service-deployment.yaml

kubectl apply -f employer-service.yaml

// delete 

kubectl delete deployment -l app=devsmd-employer-service

kubectl delete services -l app=devsmd-employer-service

kubectl delete deployment --all

kubectl delete secrets --all

kubectl delete services --all

// restart 

kubectl rollout restart deploy devsmd-employer-service

// Push to docker hub

docker build -t devsmd/employer-service .

docker tag devsmd/employer-service devsmd/employer-service:1.3.2

docker push devsmd/employer-service:1.3.2