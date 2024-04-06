#!/bin/sh

set -e

context="$(kubectl config current-context)"
cluster="$(kubectl config view -o jsonpath="{.contexts[?(@.name == \"$context\")].context.cluster}")"
server="$(kubectl config view -o jsonpath="{.clusters[?(@.name == \"$cluster\")].cluster.server}")"

if echo "$server" | grep -qE '^https://[^/]*[.]com(:[0-9]+)?/?$'; then
    echo "ERROR: Cannot make changes in Kubernetes cluster with API server $server" >&2
    exit 1
fi
