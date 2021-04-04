# papertrader-api
back end for a paper paper trading app


## setup

### local development
_the api assumes a configured database, you can get one running using the directions in the **kubernetes** portion below_




### docker image
```bash
cd $GOPATH/src/github.com/papertrader-api/kube

#push the newest api image
docker build -t calebwilliamsdatastax/papertrader-api .
docker push calebwilliamsdatastax/papertrader-api
```

### kubernetes
```bash
cd $GOPATH/src/github.com/papertrader-api/kube

#create the cluster from our kind.yaml file
kind create cluster --config kind.yaml

#create our kubernetes namespace
kubectl create namespace papertrader

#set our context
kubectl config set-context --current --namespace papertrader

#cluster configuration & database
kubectl create -f configmap.yaml
kubectl create -f stargate-service.yaml
kubectl create -f stargate-deployment.yaml

#optional for api development
kubectl create -f papertrader-service.yaml
kubectl create -f papertrader-deployment.yaml

#check the stuff
kubectl get pods
```