#!/bin/bash

if [[ ! $(command -v k3d) ]]; then
  echo "k3d is not installed. Please install k3d first."
  exit 1
fi

if [[ -z $CLUSTER_NAME ]]; then
  export CLUSTER_NAME=test
fi


if [[ $(k3d cluster list | grep ${CLUSTER_NAME}) ]]; then
  if [[ $1 == "reset" ]]; then
    k3d cluster delete ${CLUSTER_NAME}
    k3d cluster create ${CLUSTER_NAME}
  else
    kubectl config use-context k3d-${CLUSTER_NAME}
  fi
else
  k3d cluster create ${CLUSTER_NAME}
fi


helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

kubectl wait --for=condition=ready node --all --timeout=300s

if [[ -z $NAMESPACE ]]; then
  export NAMESPACE=test
fi

if [[ -z $GOOGLE_JSON_KEY ]]; then
  export GOOGLE_JSON_KEY=./secrets/credentials.json
fi

if [[ -z $GCS_BUCKET_NAME ]]; then
  export GCS_BUCKET_NAME=pixo-test-bucket
fi

if [[ -z $SA_NAME ]]; then
  export SA_NAME=test-sa
fi

kubectl create ns $NAMESPACE
kubectl config set-context --current --namespace=$NAMESPACE

kubectl create sa $SA_NAME -n $NAMESPACE

kubectl create clusterrolebinding test-admin-binding --clusterrole=cluster-admin --serviceaccount=$NAMESPACE:$SA_NAME
kubectl create clusterrolebinding test-deployment-binding --clusterrole=system:controller:deployment-controller --serviceaccount=$NAMESPACE:$SA_NAME

helm upgrade --install argo-events argo/argo-events \
  --namespace argo-events \
  --create-namespace

helm upgrade --install argo-workflows argo/argo-workflows \
  --namespace argo-workflows \
  --create-namespace \
  --version 0.22.15

kubectl wait --namespace argo-events --for=condition=ready pod --selector=app.kubernetes.io/instance=argo-events --timeout=300s
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
      bucket: $GCS_BUCKET_NAME
      path: artifacts
      serviceAccountKeySecret:
        name: google-credentials
        key: credentials
EOF

kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: $SA_NAME.service-account-token
  annotations:
    kubernetes.io/service-account.name: $SA_NAME
type: kubernetes.io/service-account-token
EOF

ARGO_WORKFLOWS_TOKEN="$(kubectl get secret $SA_NAME.service-account-token -o=jsonpath='{.data.token}' | base64 --decode)"
echo
echo "ARGO_WORKFLOWS_TOKEN:"
echo "Bearer $ARGO_WORKFLOWS_TOKEN"

echo
echo "curl http://localhost:2746/api/v1/workflows/argo -H \"Authorization: $ARGO_WORKFLOWS_TOKEN\""

echo
echo "kubectl port-forward svc/argo-workflows-server 2746:2746 -n argo-workflows"
