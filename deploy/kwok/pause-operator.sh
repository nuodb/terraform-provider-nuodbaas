#!/bin/sh

set -e

# Stop Operator container and wait for it to exit
docker stop nuodb-cp-operator
docker wait nuodb-cp-operator
