#!/usr/bin/env bash
set -e

CLUSTER_NAME=aether-e2e

make test

# Spin up test cluster.
kind create cluster --config=test/kind-conf.yaml --name ${CLUSTER_NAME}
export KUBECONFIG="$(kind get kubeconfig-path --name="$CLUSTER_NAME")"

echo "Loading images..."
kind load docker-image compute-aether --name=${CLUSTER_NAME}
kind load docker-image aether-tests --name=${CLUSTER_NAME}

# Deploy Compute Aether.
kubectl apply -f ./deploy.yaml

kubectl apply -f ./test/container/deploy.yaml
