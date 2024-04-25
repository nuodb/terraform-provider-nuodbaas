#!/bin/sh

set -e
cd "$(dirname "$0")"

# Get hostname for Nginx load balancer service, which could be IP or hostname
hostname="$(kubectl get service ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}')"
: ${hostname:="$(kubectl get service ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')"}

# Get port for Nginx load balancer service
port="$(kubectl get service ingress-nginx-controller -o jsonpath='{.spec.ports[?(@.appProtocol=="http")].port}')"

# Get password for system/admin user generate by `helm install`
export NUODB_CP_URL_BASE="http://$hostname:$port/nuodb-cp"
export NUODB_CP_USER=system/admin
export NUODB_CP_PASSWORD="$(kubectl get secret dbaas-user-system-admin -o jsonpath='{.data.password}' | base64 -d)"

# Use token authentication if server is configured to generate tokens
if [ -n "$(kubectl get secrets nuodb-cp-runtime-config -o jsonpath='{.data.secretPassword}' --ignore-not-found)" ]; then
    NUODB_CP_TOKEN="$(nuodb-cp httpclient POST login --jsonpath token --unquote)"
    if [ -n "$NUODB_CP_TOKEN" ]; then
        unset NUODB_CP_USER
        unset NUODB_CP_PASSWORD
    fi
fi

cat <<EOF
export NUODB_CP_URL_BASE="$NUODB_CP_URL_BASE"
export NUODB_CP_USER="$NUODB_CP_USER"
export NUODB_CP_PASSWORD="$NUODB_CP_PASSWORD"
export NUODB_CP_TOKEN="$NUODB_CP_TOKEN"

export PAUSE_OPERATOR_COMMAND="$(pwd)/pause-operator.sh"
export RESUME_OPERATOR_COMMAND="$(pwd)/resume-operator.sh"

export CONTAINER_SCHEDULING_ENABLED="true"
export WEBHOOKS_ENABLED="true"
EOF
