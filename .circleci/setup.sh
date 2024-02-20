#!/usr/bin/env bash

set -ex

# Create test results and output directories
mkdir -p "$TEST_RESULTS"
mkdir -p "$OUTPUT_DIR"

# Download kubectl
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v"${KUBERNETES_VERSION}"/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/

if [ "$E2E_TEST" = true ]; then
    # Download minikube
    curl -Lo minikube https://storage.googleapis.com/minikube/releases/v"${MINIKUBE_VERSION}"/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/

    # Start minikube
    minikube start --vm-driver=docker --kubernetes-version=v"${KUBERNETES_VERSION}"
    minikube status
    kubectl cluster-info

    nohup minikube tunnel > "${OUTPUT_DIR}/minikube_tunnel.log" 2>&1 &

    # Install helm
    curl -Lo /tmp/helm.tar.gz https://get.helm.sh/helm-"${HELM_VERSION}"-linux-amd64.tar.gz
    tar xzf /tmp/helm.tar.gz -C /tmp --strip-components=1 && chmod +x /tmp/helm && sudo mv /tmp/helm /usr/local/bin

    # Install Cert-manager, Nginx Ingress, and the Control Plane into K8s cluster
    make deploy-cp

    # Extract credentials and make them available as environment variables for subsequent steps in job
    make extract-creds >> "$BASH_ENV"
else
    # Download and start test helper for REST service that includes CRUD-only K8s and mock operators
    make deploy-test-helper

    # Allow kubectl to be used against cluster by subsequent commands if needed
    echo "export KUBECONFIG=\"$OUTPUT_DIR/kubeconfig.yml\"" >> "$BASH_ENV"
fi
