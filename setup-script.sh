#!/bin/bash

set -e

echo "Building Docker image for antaeus.."
docker build -t tinjis_pleo-antaeus:0.1 .
echo "Building Docker image for payment service.."
docker build -f Dockerfile-payment -t tinjis_pleo-payment:0.1 .

cd kubernetes
echo "Creating payment deployment.."
kubectl apply -f payment.yml
echo "Creating payment service.."
kubectl apply -f payment-svc.yml
echo "Creating antaeus deployment.."
kubectl apply -f antaeus.yml
echo "Creating antaeus service"
kubectl apply -f antaeus-svc.yml

sleep 10
kubectl get pods
antaeus=$(kubectl get pod -l run=antaeus -o go-template='{{range .items}}{{.metadata.name}}{{end}}')
payment=$(kubectl get pod -l run=payment -o go-template='{{range .items}}{{.metadata.name}}{{end}}')
NODEPORT=$(kubectl get service antaeus-service -o go-template='{{(index .spec.ports 0).nodePort}}')
echo "Nodeport for Antaeus service - $NODEPORT"
echo "Execute following command to export nodeport.." 
echo 
echo 'export NODEPORT=$(kubectl get service antaeus-service -o go-template='"'"'{{(index .spec.ports 0).nodePort}}'"'"')'
