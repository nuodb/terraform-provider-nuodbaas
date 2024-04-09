#!/bin/sh

set -e
cd "$(dirname "$0")"
./check.sh

if [ "$IGNORE_NOT_FOUND" = true ]; then
    set +e
fi

echo "Deleting Control Plane resources..."
for crd in $(kubectl get crd -o name | sed -n 's|.*/\(.*\.cp\.nuodb\.com\)|\1|p'); do
    kubectl get "$crd" -o name | xargs -r kubectl delete
done

echo "Uninstalling Control Plane and dependencies..."
helm uninstall nuodb-cp-rest
helm uninstall nuodb-cp-operator
helm uninstall nuodb-cp-crd
helm uninstall ingress-nginx
helm uninstall cert-manager -n cert-manager

echo "Deleting Cert Manager leases to accelerate restart..."
kubectl -n kube-system delete leases.coordination.k8s.io cert-manager-cainjector-leader-election
kubectl -n kube-system delete leases.coordination.k8s.io cert-manager-controller
