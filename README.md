# cncamp

## Image
We can pull the image through below command:
```shell
docker pull joey100/cncamp-httpserver:v2.0
```

## Deployment

Follow the steps below:
```
kubectl apply -f k8s/nginx-ingress-deployment.yaml   # we must deploy ingress nginx controller first
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
INGRESS_SVC_IP=`kubectl get svc ingress-nginx-controller -n ingress-nginx|awk 'END{print $3}'`
curl -H "Host: cncamp-httpserver.com" http://${INGRESS_SVC_IP}/foo
```
