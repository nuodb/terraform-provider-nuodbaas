#!/bin/sh

set -e

# Scale Operator deployment to 0 and wait for pods to be deleted
kubectl scale deployment nuodb-cp-operator --replicas=0
kubectl wait pod --for=delete -l app=nuodb-cp-operator --timeout=30s
