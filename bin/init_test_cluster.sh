#!/bin/bash

# Check for k3d
if [[ ! $(command -v k3d) ]]; then
  echo "k3d is not installed. Please install k3d first."
  exit 1
fi

if [[ $(k3d cluster list | grep test-cluster) ]]; then
  k3d cluster delete test-cluster
fi

k3d cluster create test-cluster

helm repo add agones https://agones.dev/chart/stable
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

kubectl wait --for=condition=ready node --all --timeout=300s

kubectl create ns dev-multiplayer
helm install agones agones/agones --namespace agones-system --create-namespace --set 'gameservers.namespaces[0]'=dev-multiplayer
helm install argo-workflows argo/argo-workflows --namespace argo-workflows --create-namespace

for crd in gameservers.agones.dev gameserversets.agones.dev fleets.agones.dev fleetautoscalers.autoscaling.agones.dev; do
  kubectl wait --for=condition=established --timeout=300s crd/$crd
done

kubectl wait --namespace agones-system --for=condition=ready pod --selector=app=agones --timeout=300s
kubectl wait --namespace argo-workflows --for=condition=ready pod --selector=app.kubernetes.io/instance=argo-workflows --timeout=300s
