#!/bin/sh

# Set credentials supplied as CircleCI environment variables
cat <<EOF
export NUODB_CP_USER="$DBAAS_TEST_USER"
export NUODB_CP_PASSWORD="$DBAAS_TEST_PASSWORD"
export NUODB_CP_URL_BASE="$DBAAS_API_ENDPOINT"

export CONTAINER_SCHEDULING_ENABLED="true"
export WEBHOOKS_ENABLED="true"
export ORGANIZATION_BOUND_USER="true"
export TESTARGS="-tags shared_env"
EOF
