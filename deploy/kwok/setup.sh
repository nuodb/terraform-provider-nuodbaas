#!/bin/sh

set -e

: ${TMP_DIR:="/tmp"}
: ${KWOK_WORKDIR:="$TMP_DIR/kwok"}
export KWOK_WORKDIR

: ${KUBECONFIG:="$TMP_DIR/kubeconfig.yaml"}
export KUBECONFIG

# Special volume provisioner required for scheduling stateful workloads in
# KWOK. See https://github.com/adriansuarez/noop-provisioner.
: ${PROVISIONER_IMAGE:="ghcr.io/adriansuarez/noop-provisioner:latest"}

: ${NUODB_CP_REPO:="https://nuodb.github.io/nuodb-cp-releases/charts"}
: ${NUODB_CP_VERSION:="2.7.0"}
: ${NUODB_CP_IMAGE:="ghcr.io/nuodb/nuodb-cp-images:$NUODB_CP_VERSION"}
: ${NUODB_CP_USER:="system/admin"}
: ${NUODB_CP_PASSWORD:="changeIt"}

echo "Creating K8s cluster..."
kwokctl create cluster --wait 1m
kwokctl scale node --replicas 5

CLUSTER_DIR="$KWOK_WORKDIR/clusters/kwok"
DOCKER_NET="kwok-kwok"

# Make everything in kwok directory world-readable to avoid annoying Docker
# volume permissions issues. The NuoDB CP image has uid 1000 for the runtime
# user and would not be able access files owned by root on the host that are
# not group- or world-readable. In CircleCI, the root user invokes this script,
# and by default the K8s credentials are not world-readable.
chmod -R a+X,a+r "$CLUSTER_DIR"

# Start noop volume provisioner that immediately binds PVCs
docker run -d --name noop-provisioner \
    -v "$CLUSTER_DIR"/kubeconfig:/kubeconfig \
    -v "$CLUSTER_DIR"/pki:/etc/kubernetes/pki \
    --network "$DOCKER_NET" \
    "$PROVISIONER_IMAGE" -kubeconfig /kubeconfig

# Create storage class for noop provisioner
kubectl apply -f - <<EOF
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  name: noop
provisioner: nuodb.github.io/noop-provisioner
EOF

echo "Installing CRDs for DBaaS..."
if [ -n "$NUODB_CP_CRD_CHART" ]; then
    helm install nuodb-cp-crd "$NUODB_CP_CRD_CHART"
else
    helm install nuodb-cp-crd nuodb-cp-crd --repo "$NUODB_CP_REPO" --version "$NUODB_CP_VERSION"
fi

# Create service tiers. The Helm features do not matter, since they are not
# exposed in Terraform.
kubectl apply -f - <<EOF
apiVersion: cp.nuodb.com/v1beta1
kind: ServiceTier
metadata:
  name: n0.nano
spec:
  features: []
---
apiVersion: cp.nuodb.com/v1beta1
kind: ServiceTier
metadata:
  name: n0.small
spec:
  features: []
---
apiVersion: cp.nuodb.com/v1beta1
kind: ServiceTier
metadata:
  name: n1.small
spec:
  features: []
EOF

# Create service account needed for AP
kubectl create serviceaccount nuodb

# Start Operator with webhooks and backups disabled, since they cannot work in
# this environment
echo "Starting DBaaS Operator..."
docker run -d --name nuodb-cp-operator \
    -v "$CLUSTER_DIR"/kubeconfig:/home/nuodb/.kube/config \
    -v "$CLUSTER_DIR"/pki:/etc/kubernetes/pki \
    -e ENABLE_WEBHOOKS=false \
    --network "$DOCKER_NET" \
    "$NUODB_CP_IMAGE" \
    controller --feature-gates EmbeddedDatabaseBackupPlugin=false

echo "Starting DBaaS REST service..."
docker run -d --name nuodb-cp-rest -p 8080:8080 \
    -v "$CLUSTER_DIR"/kubeconfig:/home/nuodb/.kube/config \
    -v "$CLUSTER_DIR"/pki:/etc/kubernetes/pki \
    --network "$DOCKER_NET" \
    "$NUODB_CP_IMAGE" \
    nuodb-cp server start

echo "Waiting for REST service to become ready..."
check() {
    out="$(docker exec nuodb-cp-rest nuodb-cp httpclient GET /healthz 2>&1)"
}
i=0
n=15
while ! check; do
    if [ $i -ge $n ]; then
        echo "Readiness check failed after $n seconds:"
        echo "$out"
        exit 1
    fi
    echo "Retrying in 1 second..."
    sleep 1
    i=$((i + 1))
done

echo "Creating DBaaS system/admin user..."
docker exec nuodb-cp-rest nuodb-cp user create "$NUODB_CP_USER" --user-password "$NUODB_CP_PASSWORD" --allow 'all:*' --allow-cross-organization

cat <<EOF
Setup complete.

To inspect Kubernetes state:

  export KUBECONFIG="$KUBECONFIG"
  kubectl get domains,databases,pods

To inspect DBaaS state:

  export NUODB_CP_USER="$NUODB_CP_USER"
  export NUODB_CP_PASSWORD="$NUODB_CP_PASSWORD"
  export NUODB_CP_URL_BASE=http://localhost:8080
  alias nuodb-cp="docker exec nuodb-cp-rest nuodb-cp"
  nuodb-cp project list /
  nuodb-cp database list /

EOF
