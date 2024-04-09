#!/bin/sh

if [ -d "$OUTPUT_DIR" ]; then
    mkdir -p "$OUTPUT_DIR/cluster-info"
    kubectl cluster-info dump -A --output-directory="$OUTPUT_DIR/cluster-info"
else
    kubectl cluster-info dump -A
fi
