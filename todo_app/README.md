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