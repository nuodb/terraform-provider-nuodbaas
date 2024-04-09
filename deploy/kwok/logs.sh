#!/bin/sh

collect_logs() {
    if [ -d "$OUTPUT_DIR" ]; then
        docker logs "$1" > "$OUTPUT_DIR/$1.log" 2>&1
    else
        docker logs "$1"
    fi
}

# Collect logging for Control Plane containers and volume provisioner
collect_logs nuodb-cp-rest
collect_logs nuodb-cp-operator
collect_logs noop-provisioner

# Collect K8s diagnostics
: ${TMP_DIR:="/tmp"}
: ${KUBECONFIG:="$TMP_DIR/kubeconfig.yaml"}
export KUBECONFIG
if [ -d "$OUTPUT_DIR" ]; then
    mkdir -p "$OUTPUT_DIR/cluster-info"
    kubectl cluster-info dump -A --output-directory="$OUTPUT_DIR/cluster-info"
else
    kubectl cluster-info dump -A
fi
