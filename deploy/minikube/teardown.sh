#!/bin/sh

# Kill minikube tunnel process
if [ -f "$OUTPUT_DIR/minikube-tunnel.pid" ]; then
    kill "$(cat "$OUTPUT_DIR/minikube-tunnel.pid")"
    rm "$OUTPUT_DIR/minikube-tunnel.pid"
fi

# Stop minikube cluster
minikube stop
