#!/bin/bash

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

if [[ -z $NAMESPACE ]]; then
  export NAMESPACE=test
fi

if [[ -z $SA_NAME ]]; then
  export SA_NAME=test-sa
fi

kubectl create ns $NAMESPACE
kubectl config set-context --current --namespace=$NAMESPACE

kubectl create sa $SA_NAME -n $NAMESPACE

kubectl create clusterrolebinding test-sa-admin-binding --clusterrole=cluster-admin --serviceaccount=$NAMESPACE:$SA_NAME
kubectl create clusterrolebinding test-sa-deployment-binding --clusterrole=system:controller:deployment-controller --serviceaccount=$NAMESPACE:$SA_NAME

helm install agones agones/agones --namespace agones-system --create-namespace --set 'gameservers.namespaces[0]'=$NAMESPACE
helm install argo-workflows argo/argo-workflows --namespace argo-workflows --create-namespace --version 0.22.15

for crd in gameservers.agones.dev gameserversets.agones.dev fleets.agones.dev fleetautoscalers.autoscaling.agones.dev; do
  kubectl wait --for=condition=established --timeout=300s crd/$crd
done

kubectl wait --namespace agones-system --for=condition=ready pod --selector=app=agones --timeout=300s
kubectl wait --namespace argo-workflows --for=condition=ready pod --selector=app.kubernetes.io/instance=argo-workflows --timeout=300s

kubectl create secret generic google-credentials --from-file=credentials=$GOOGLE_JSON_KEY -n $NAMESPACE

kubectl apply -n $NAMESPACE -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: artifact-repositories
data:
  gcs-artifact-repository: |
    gcs:
      bucket: pixo-test-bucket
      path: artifacts
      serviceAccountKeySecret:
        name: google-credentials
        key: credentials
EOF
