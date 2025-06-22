# Log Output Application

## How to run

### For exercise 1.1

Run this command:
```
kubectl create deployment log-output --image=sericy/log_output
```

### For exercise 1.3

Run this command:
```
kubectl apply -f manifests/deployment.yaml
```

### For exercise 1.7

Run this command:
```
kubectl apply -f manifests
```
Make sure there is no other ingress listening on port 80.