# To Do Application

## How to run

### For exercise 1.2

Run this command:
```
kubectl create deployment todo-app --image=sericy/todo-app
```

### For exercise 1.4

Run this command:
```
kubectl apply -f manifests/deployment.yaml
```

### For exercise 1.6

Run this command:
```
kubectl apply -f manifests/service.yaml
```
Make sure port 30080 is connected to host port with k3d, and no other service is using port 30080

### For exercise 1.8

Run this command:
```
kubectl apply -f manifests/ingress.yaml
```
Make sure there is no other ingress listening on port 80.

### For exercise 1.12

Run this command:
```
kubectl apply -f manifests
```
Make sure there is no other ingress listening on port 80 and there is no other persistent volume on /tmp/kube path.