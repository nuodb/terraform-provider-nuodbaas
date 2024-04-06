#!/bin/sh

set -e

# Scale Operator deployment back up to 1
kubectl scale deployment nuodb-cp-operator --replicas=1
