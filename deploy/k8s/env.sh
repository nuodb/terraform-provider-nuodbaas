#!/bin/sh

set -e
cd "$(dirname "$0")"

# Get hostname for Nginx load balancer service, which could be IP or hostname
hostname="$(kubectl get service ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}')"
: ${hostname:="$(kubectl get service ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')"}

# Get port for Nginx load balancer service
port="$(kubectl get service ingress-nginx-controller -o jsonpath='{.spec.ports[?(@.appProtocol=="http")].port}')"

# Get password for system/admin user generate by `helm install`
NUODB_CP_URL_BASE="http://$hostname:$port/nuodb-cp"
NUODB_CP_USER=system/admin
NUODB_CP_PASSWORD="$(kubectl get secret dbaas-user-system-admin -o jsonpath='{.data.password}' | base64 -d)"

cat <<EOF
export NUODB_CP_URL_BASE="$NUODB_CP_URL_BASE"
export NUODB_CP_USER="$NUODB_CP_USER"
export NUODB_CP_PASSWORD="$NUODB_CP_PASSWORD"

export PAUSE_OPERATOR_COMMAND="$(pwd)/pause-operator.sh"
export RESUME_OPERATOR_COMMAND="$(pwd)/resume-operator.sh"

export POD_SCHEDULING_ENABLED="true"
export WEBHOOKS_ENABLED="true"
EOF
