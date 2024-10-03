#!/bin/sh

if [ -d "$OUTPUT_DIR" ]; then
    mkdir -p "$OUTPUT_DIR/cluster-info"
    kubectl cluster-info dump -A --output-directory="$OUTPUT_DIR/cluster-info"

    # For any containers that have restarted, collect the exited container output
    (
        kubectl get pod -A -o jsonpath="{range .items[*]}{.metadata.namespace}{' '}{.metadata.name}{' '}{.status.containerStatuses[?(.restartCount > 0)].name}{'\n'}{end}"
    ) | while read ns pod containers; do
        for container in $containers; do
            kubectl logs -n "$ns" "$pod" -c "$container" --previous > "$OUTPUT_DIR/cluster-info/$ns/$pod/logs-$container-previous.txt"
        done
    done
else
    kubectl cluster-info dump -A
fi

