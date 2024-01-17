HELM_HAPROXY_RELEASE ?= haproxy-ingress
HAPROXY_CHARTS_VERSION ?= 1.25.1
HAPROXY_CHART ?= https://github.com/haproxytech/helm-charts/releases/download/kubernetes-ingress-$(HAPROXY_CHARTS_VERSION)/kubernetes-ingress-$(HAPROXY_CHARTS_VERSION).tgz

HELM_JETSTACK_RELEASE ?= cert-manager
JETSTACK_CHARTS_VERSION ?= 1.13.3

HELM_CP_CRD_RELEASE ?= nuodb-cp-crd
HELM_CP_OPERATOR_RELEASE ?= nuodb-cp-operator
HELM_CP_REST_RELEASE ?= nuodb-cp-rest

HELM_NGINX_RELEASE=ingress-nginx
NGINX_CHART=https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-4.7.1/ingress-nginx-4.7.1.tgz
NGINX_INGRESS_VERSION=1.8.1

IGNORE_NOT_FOUND ?= true

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Development

.PHONY: deploy-cp
deploy-cp: ## Deploy a local Control Plane
	helm repo add jetstack https://charts.jetstack.io
	helm repo add nuodb-cp https://nuodb.github.io/nuodb-cp-releases/charts
	helm repo update

	helm upgrade --install $(HELM_JETSTACK_RELEASE) jetstack/cert-manager \
		--version $(JETSTACK_CHARTS_VERSION) \
		--namespace cert-manager \
		--set installCRDs=true \
		--create-namespace

	helm upgrade --install "$(HELM_NGINX_RELEASE)" "$(NGINX_CHART)" \
           --set controller.image.tag="$(NGINX_INGRESS_VERSION)" \
           --set controller.ingressClassResource.default=true \
           --set controller.replicaCount=1 \
           --set controller.service.enablePorts.http=false \
           --set controller.service.nodePorts.https="30500" \
           --set controller.extraArgs.enable-ssl-passthrough=true

	helm upgrade --install $(HELM_CP_CRD_RELEASE) nuodb-cp/nuodb-cp-crd

	@# Wait for all Control Plane dependencies to be ready
	kubectl -n $(HELM_JETSTACK_RELEASE) wait pod --all --for=condition=Ready
	kubectl -l app.kubernetes.io/instance="$(HELM_NGINX_RELEASE)" wait pod --for=condition=ready --timeout=120s

	helm upgrade --install $(HELM_CP_OPERATOR_RELEASE) nuodb-cp/nuodb-cp-operator \
		--set cpOperator.webhooks.enabled=true

	helm upgrade --install $(HELM_CP_REST_RELEASE) nuodb-cp/nuodb-cp-rest \
		--set cpRest.ingress.enabled=true \
		--set cpRest.ingress.className=nginx \
		--set cpRest.authentication.admin.create=true

	kubectl apply -f test/tiers.yaml

	kubectl wait pod -l app.kubernetes.io/instance="$(HELM_CP_OPERATOR_RELEASE)" --for=condition=ready --timeout=120s
	kubectl wait pod -l app.kubernetes.io/instance="$(HELM_CP_REST_RELEASE)" --for=condition=ready --timeout=120s

.PHONY: undeploy-cp
undeploy-cp: ## Uninstall a local Control Plane previously installed by this script
	kubectl delete --wait=false -f test/tiers.yaml || $(IGNORE_NOT_FOUND)
	helm uninstall $(HELM_CP_REST_RELEASE) $(HELM_CP_OPERATOR_RELEASE) || $(IGNORE_NOT_FOUND)
	helm uninstall $(HELM_CP_CRD_RELEASE) $(HELM_NGINX_RELEASE) || $(IGNORE_NOT_FOUND)
	helm uninstall $(HELM_JETSTACK_RELEASE) --namespace cert-manager || $(IGNORE_NOT_FOUND)


.PHONY: discover-test
discover-test: ## Discover a local control plane and run tests against it
	$(eval HOST := $(or $(shell kubectl get service ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}'), \
						$(shell kubectl get service ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')))
	$(eval PORT := $(shell kubectl get service ingress-nginx-controller -o jsonpath='{.spec.ports[?(@.appProtocol=="http")].port}'))

	@NUODB_CP_PASSWORD=$(shell kubectl get secret dbaas-user-system-admin -o jsonpath='{.data.password}' | base64 -d) \
		NUODB_CP_URL_BASE="http://$(HOST):$(PORT)/nuodb-cp" \
		NUODB_CP_USER="admin" \
		NUODB_CP_ORGANIZATION="system" \
		$(MAKE) testacc


.PHONY: testacc
testacc: ## Run acceptance tests
	TF_ACC=1 go test ./plugin/... -v $(TESTARGS) -timeout 30m
