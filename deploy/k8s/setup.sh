#!/bin/sh

set -e
cd "$(dirname "$0")"
./check.sh

: ${NUODB_CP_VERSION:="2.5.0"}
: ${NUODB_CP_REPO:="https://nuodb.github.io/nuodb-cp-releases/charts"}

: ${CERT_MANAGER_VERSION:="1.13.3"}
: ${CERT_MANAGER_REPO:="https://charts.jetstack.io"}

: ${NGINX_INGRESS_VERSION:="4.7.1"}
: ${NGINX_INGRESS_REPO:="https://kubernetes.github.io/ingress-nginx"}

helm upgrade --install cert-manager cert-manager \
    --repo "$CERT_MANAGER_REPO" --version "$CERT_MANAGER_VERSION" \
    --namespace cert-manager \
    --set installCRDs=true \
    --create-namespace

helm upgrade --install ingress-nginx ingress-nginx \
    --repo "$NGINX_INGRESS_REPO" --version "$NGINX_INGRESS_VERSION" \
    --set controller.ingressClassResource.default=true \
    --set controller.replicaCount=1 \
    --set controller.service.ports.http=8080 \
    --set controller.service.ports.https=8443 \
    --set controller.extraArgs.enable-ssl-passthrough=true

echo "Waiting for Control Plane dependencies to become ready..."
kubectl -n cert-manager wait pod --all --for=condition=Ready
kubectl -l app.kubernetes.io/instance=ingress-nginx wait pod --for=condition=ready --timeout=120s

helm upgrade --install nuodb-cp-crd nuodb-cp-crd \
    --repo "$NUODB_CP_REPO" --version "$NUODB_CP_VERSION"

helm upgrade --install nuodb-cp-operator nuodb-cp-operator \
    --repo "$NUODB_CP_REPO" --version "$NUODB_CP_VERSION" \
    --set cpOperator.webhooks.enabled=true \
    --set cpOperator.samples.serviceTiers.enabled=false

helm upgrade --install nuodb-cp-rest nuodb-cp-rest \
    --repo "$NUODB_CP_REPO" --version "$NUODB_CP_VERSION" \
    --set cpRest.ingress.enabled=true \
    --set cpRest.ingress.className=nginx \
    --set cpRest.authentication.admin.create=true

echo "Waiting for Control Plane to become ready..."
kubectl wait pod -l app.kubernetes.io/instance=nuodb-cp-operator --for=condition=ready --timeout=120s
kubectl wait pod -l app.kubernetes.io/instance=nuodb-cp-rest --for=condition=ready --timeout=120s
