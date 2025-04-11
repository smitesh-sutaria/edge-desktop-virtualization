kubectl delete  -f deploy/manifests/maverikflats-device-plugin.yaml
sleep 5
kubectl get po -A
rm device-plugin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o device-plugin cmd/main.go
docker build --no-cache -t localhost:5000/device-plugin:v17 .
docker push localhost:5000/device-plugin:v17
docker tag localhost:5000/device-plugin:v17 10.190.167.103:5000/device-plugin:v17
kubectl apply -f deploy/manifests/maverikflats-device-plugin.yaml
sleep 1
kubectl get po -A
