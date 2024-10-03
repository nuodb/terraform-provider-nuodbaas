#!/bin/sh

set -e
cd "$(dirname "$0")"

: ${KUBERNETES_VERSION:="1.31.1"}
: ${TMP_DIR:="/tmp"}
: ${KUBECONFIG:="$TMP_DIR/kubeconfig.yaml"}
export KUBECONFIG

# Start minikube cluster
minikube start --vm-driver=docker --kubernetes-version=v"$KUBERNETES_VERSION"
minikube status

# Create tunnel to enable external connectivity
nohup minikube tunnel > "$OUTPUT_DIR/minikube-tunnel.out" 2>&1 &
echo "$!" > "$OUTPUT_DIR/minikube-tunnel.pid"

# Install Cert-manager, Nginx Ingress, and the Control Plane into K8s cluster
../k8s/setup.sh

# Override service tiers and Helm features
kubectl apply -f tiers.yaml
