#!/usr/bin/env bash

set -ex

# Download kubectl
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v"${KUBERNETES_VERSION}"/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# download minikube
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v"${MINIKUBE_VERSION}"/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/

# start minikube
minikube start --vm-driver=docker --kubernetes-version=v"${KUBERNETES_VERSION}"
minikube status
kubectl cluster-info

nohup minikube tunnel > "${TEST_RESULTS}/minikube_tunnel.log" 2>&1 &

# install helm
wget https://get.helm.sh/helm-"${HELM_VERSION}"-linux-amd64.tar.gz -O /tmp/helm.tar.gz
tar xzf /tmp/helm.tar.gz -C /tmp --strip-components=1 && chmod +x /tmp/helm && sudo mv /tmp/helm /usr/local/bin
