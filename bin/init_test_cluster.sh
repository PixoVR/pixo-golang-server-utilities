#!/bin/bash

if [[ ! $(command -v k3d) ]]; then
  echo "k3d is not installed. Please install k3d first."
  exit 1
fi

if [[ -z $CLUSTER_NAME ]]; then
  export CLUSTER_NAME=test
fi


if [[ $(k3d cluster list | grep ${CLUSTER_NAME}) ]]; then
  kubectl config use-context k3d-${CLUSTER_NAME}
else
  k3d cluster create ${CLUSTER_NAME}
fi


helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

kubectl wait --for=condition=ready node --all --timeout=300s

if [[ -z $NAMESPACE ]]; then
  export NAMESPACE=test
fi

if [[ -z $SA_NAME ]]; then
  export SA_NAME=test
fi

kubectl create ns $NAMESPACE
kubectl config set-context --current --namespace=$NAMESPACE

kubectl create sa $SA_NAME -n $NAMESPACE

kubectl create clusterrolebinding test-admin-binding --clusterrole=cluster-admin --serviceaccount=$NAMESPACE:$SA_NAME
kubectl create clusterrolebinding test-deployment-binding --clusterrole=system:controller:deployment-controller --serviceaccount=$NAMESPACE:$SA_NAME

helm upgrade --install argo-workflows argo/argo-workflows \
  --namespace argo-workflows \
  --create-namespace \
  --version 0.22.15

kubectl wait --namespace argo-workflows --for=condition=ready pod --selector=app.kubernetes.io/instance=argo-workflows --timeout=300s
