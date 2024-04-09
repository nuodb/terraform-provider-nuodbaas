#!/bin/sh

set -e

# Get the API server for the Kubernetes cluster
context="$(kubectl config current-context)"
cluster="$(kubectl config view -o jsonpath="{.contexts[?(@.name == \"$context\")].context.cluster}")"
server="$(kubectl config view -o jsonpath="{.clusters[?(@.name == \"$cluster\")].cluster.server}")"

# Check that API server does not have .com domain. This is to prevent from
# making changes to a shared cluster like EKS.
if echo "$server" | grep -qE '^https://[^/]*[.]com(:[0-9]+)?/?$'; then
    echo "ERROR: Cannot make changes in Kubernetes cluster with API server $server" >&2
    exit 1
fi
