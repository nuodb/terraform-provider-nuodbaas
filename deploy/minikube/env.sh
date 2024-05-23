#!/bin/sh

set -e
cd "$(dirname "$0")"

# Print KUBECONFIG used by minikube
: ${TMP_DIR:="/tmp"}
: ${KUBECONFIG:="$TMP_DIR/kubeconfig.yaml"}
export KUBECONFIG
cat <<EOF
export KUBECONFIG="$KUBECONFIG"
EOF

# Print common K8s environment variables
../k8s/env.sh
