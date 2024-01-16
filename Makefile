HELM_HAPROXY_RELEASE ?= haproxy-ingress
HAPROXY_CHARTS_VERSION ?= 1.25.1
HAPROXY_CHART ?= https://github.com/haproxytech/helm-charts/releases/download/kubernetes-ingress-$(HAPROXY_CHARTS_VERSION)/kubernetes-ingress-$(HAPROXY_CHARTS_VERSION).tgz

HELM_CP_CRD_RELEASE ?= nuodb-cp-crd
HELM_CP_OPERATOR_RELEASE ?= nuodb-cp-operator
HELM_CP_REST_RELEASE ?= nuodb-cp-rest

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
	helm upgrade --install $(HELM_HAPROXY_RELEASE) $(HAPROXY_CHART) \
              --set controller.ingressClassResource.default=true \
              --set controller.service.type=LoadBalancer

	kubectl wait pod -l app.kubernetes.io/instance="$(HELM_HAPROXY_RELEASE)" --for=condition=ready --timeout=120s

	helm repo add nuodb-cp https://nuodb.github.io/nuodb-cp-releases/charts
	helm repo update

	helm upgrade --install $(HELM_CP_CRD_RELEASE) nuodb-cp/nuodb-cp-crd

	helm upgrade --install $(HELM_CP_OPERATOR_RELEASE) nuodb-cp/nuodb-cp-operator \
		--set cpOperator.webhooks.enabled=true \
		--set cpOperator.samples.serviceTiers.enabled=true

	helm upgrade --install $(HELM_CP_REST_RELEASE) nuodb-cp/nuodb-cp-rest \
		--set cpRest.ingress.enabled=true \
		--set cpRest.ingress.className=haproxy \
		--set cpRest.authentication.admin.create=true

	kubectl wait pod -l app.kubernetes.io/instance="$(HELM_CP_OPERATOR_RELEASE)" --for=condition=ready --timeout=120s
	kubectl wait pod -l app.kubernetes.io/instance="$(HELM_CP_REST_RELEASE)" --for=condition=ready --timeout=120s

.PHONY: undeploy-cp
undeploy-cp: ## Uninstall a local Control Plane previously installed by this script
	helm uninstall $(HELM_CP_REST_RELEASE) $(HELM_CP_OPERATOR_RELEASE) || $(IGNORE_NOT_FOUND)
	helm uninstall $(HELM_CP_CRD_RELEASE) $(HELM_HAPROXY_RELEASE) || $(IGNORE_NOT_FOUND)


.PHONY: discover-test
discover-test: ## Discover a local control plane and run tests against it
	@NUODB_CP_PASSWORD=$(shell kubectl get secret dbaas-user-system-admin -o jsonpath='{.data.password}' | base64 -d) \
		NUODB_CP_URL_BASE="http://localhost/nuodb-cp" \
		NUODB_CP_USER="admin" \
		NUODB_CP_ORGANIZATION="system" \
		$(MAKE) testacc


.PHONY: testacc
testacc: ## Run acceptance tests
	TF_ACC=1 go test ./plugin/... -v $(TESTARGS)
