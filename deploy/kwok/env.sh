#!/bin/sh

set -e
cd "$(dirname "$0")"

: ${TMP_DIR:="/tmp"}
cat <<EOF
export KUBECONFIG="${KUBECONFIG:-"$TMP_DIR/kubeconfig.yaml"}"

export NUODB_CP_URL_BASE="http://localhost:8080"
export NUODB_CP_USER="${NUODB_CP_USER:-"system/admin"}"
export NUODB_CP_PASSWORD="${NUODB_CP_PASSWORD:-"changeIt"}"

export PAUSE_OPERATOR_COMMAND="$(pwd)/pause-operator.sh"
export RESUME_OPERATOR_COMMAND="$(pwd)/resume-operator.sh"

export POD_SCHEDULING_ENABLED="false"
export WEBHOOKS_ENABLED="false"
EOF
