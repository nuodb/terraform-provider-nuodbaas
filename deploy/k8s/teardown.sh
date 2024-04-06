#!/bin/sh

set -e
cd "$(dirname "$0")"
./check.sh

if [ "$IGNORE_NOT_FOUND" = true ]; then
    set +e
fi

echo "Deleting Control Plane resources..."
kubectl get databases.cp.nuodb.com -o name | xargs -r kubectl delete
kubectl get domains.cp.nuodb.com -o name | xargs -r kubectl delete
kubectl get servicetiers.cp.nuodb.com -o name | xargs -r kubectl delete
kubectl get helmfeatures.cp.nuodb.com -o name | xargs -r kubectl delete
kubectl get databasequotas.cp.nuodb.com -o name | xargs -r kubectl delete
kubectl get pvc -o name --selector=group=nuodb | xargs -r kubectl delete

echo "Uninstalling Control Plane and dependencies..."
helm uninstall nuodb-cp-rest
helm uninstall nuodb-cp-operator
helm uninstall nuodb-cp-crd
helm uninstall ingress-nginx
helm uninstall cert-manager -n cert-manager

echo "Deleting Cert Manager leases to accelerate restart..."
kubectl -n kube-system delete leases.coordination.k8s.io cert-manager-cainjector-leader-election
kubectl -n kube-system delete leases.coordination.k8s.io cert-manager-controller
