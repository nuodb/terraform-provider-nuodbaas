#!/bin/sh

if [ -d "$OUTPUT_DIR" ]; then
    mkdir -p "$OUTPUT_DIR/cluster-info"
    kubectl cluster-info dump -A --output-directory="$OUTPUT_DIR/cluster-info"

    # For any containers that have restarted, collect the exited container output
    ns="$(kubectl get serviceaccounts -o jsonpath='{.items[0].metadata.namespace}')"
    (
        kubectl get pod -o jsonpath="{range .items[*]}{.metadata.name}{' '}{.status.containerStatuses[?(.restartCount > 0)].name}{'\n'}{end}"
    ) | while read pod containers; do
        for container in $containers; do
            kubectl logs "$pod" -c "$container" -p > "$OUTPUT_DIR/cluster-info/$ns/$pod/$container-previous.txt"
        done
    done
else
    kubectl cluster-info dump -A
fi

