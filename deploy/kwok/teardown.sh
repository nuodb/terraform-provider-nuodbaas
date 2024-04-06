#!/bin/sh

: ${TMP_DIR:="/tmp"}
: ${KWOK_WORKDIR:="$TMP_DIR/kwok"}
export KWOK_WORKDIR

docker rm -f nuodb-cp-rest nuodb-cp-operator noop-provisioner
kwokctl delete cluster
