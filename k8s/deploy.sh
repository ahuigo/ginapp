kubectl delete  -f k8s/ingress.yml 
kubectl delete  -f k8s/deployment.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yml

# kubectl get deploy
kubectl get service ginapp
kubectl logs -l app=ginapp  --all-containers=true -f 