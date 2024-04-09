#!/bin/sh

cd "$(dirname "$0")"

# Collect K8s diagnostics
: ${TMP_DIR:="/tmp"}
: ${KUBECONFIG:="$TMP_DIR/kubeconfig.yaml"}
export KUBECONFIG
../k8s/logs.sh

# Collect minikube-specific logging
if [ -d "$OUTPUT_DIR" ]; then
    minikube logs --problems --file "$OUTPUT_DIR/minikube-logs.out"
else
    minikube logs --problems
fi
