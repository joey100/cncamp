# cncamp

## Image
We can pull the image through below command:
```shell
docker pull joey100/cncamp-httpserver:v2.0
docker pull joey100/cncamp-httpserver:v2.0-metrics
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

### Install loki and grafana

```
helm repo add grafana https://grafana.github.io/helm-charts
helm upgrade --install loki grafana/loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
```



### Change the grafana service to NodePort type and access it
```sh
kubectl edit svc loki-grafana -oyaml -n default
```

And change ClusterIP type to NodePort.

Login password is in secret `loki-grafana`

```sh
kubectl get secret loki-grafana -oyaml -n default
```

Find admin-password: `xxx`

```sh
echo 'xxx' | base64 -d
```

Then you will get grafana login password, the login username is 'admin' on default.

> Note: `xxx` is the value of key `admin-password` in your yaml.

### Change the grafana service to NodePort type and access it

Login password is in secret `loki-grafana`


## Metrics.go reference
https://github.com/kubernetes/autoscaler/blob/master/vertical-pod-autoscaler/pkg/utils/metrics/metrics.go




## Istio traffic management

### Install istio

```sh
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.13.2
cp bin/istioctl /usr/local/bin
istioctl install --set profile=demo -y
```

### Deploy httpsserver

```sh
kubectl create ns securesvc
kubectl label ns securesvc istio-injection=enabled
kubectl create -f k8s/deployment.yaml -n securesvc
kubectl create -f k8s/service.yaml -n securesvc
```

```sh
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt
kubectl create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt
kubectl apply -f istio-specs.yaml -n securesvc
```

### Access httpsserver

```sh
export INGRESS_IP=`kubectl get svc istio-ingressgateway -n istio-system|awk 'NR!=1{print $3}'`
# below should be OK
curl --resolve httpsserver.cncamp.io:443:$INGRESS_IP -H "user: jesse" https://httpsserver.cncamp.io/healthz -v -k
# below should not be OK
curl --resolve httpsserver.cncamp.io:443:$INGRESS_IP -H "user: jesse1" https://httpsserver.cncamp.io/healthz -v -k
```

### Tracing

```sh
kubectl apply -f k8s/jaeger.yaml
kubectl edit configmap istio -n istio-system
set tracing.sampling=100
```

#### Check tracing dashboard

```sh
istioctl dashboard jaeger --address 0.0.0.0
```

