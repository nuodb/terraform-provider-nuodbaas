HELM_JETSTACK_RELEASE ?= cert-manager
JETSTACK_CHARTS_VERSION ?= 1.13.3
JETSTACK_CHART := https://charts.jetstack.io/charts/cert-manager-v$(JETSTACK_CHARTS_VERSION).tgz

CP_CHARTS_VERSION ?= 2.3.1

HELM_CP_CRD_RELEASE ?= nuodb-cp-crd
CP_CRD_CHART ?= https://github.com/nuodb/nuodb-cp-releases/releases/download/v$(CP_CHARTS_VERSION)/nuodb-cp-crd-$(CP_CHARTS_VERSION).tgz

HELM_CP_OPERATOR_RELEASE ?= nuodb-cp-operator
CP_OPERATOR_CHART ?= https://github.com/nuodb/nuodb-cp-releases/releases/download/v$(CP_CHARTS_VERSION)/nuodb-cp-operator-$(CP_CHARTS_VERSION).tgz

HELM_CP_REST_RELEASE ?= nuodb-cp-rest
CP_REST_CHART ?= https://github.com/nuodb/nuodb-cp-releases/releases/download/v$(CP_CHARTS_VERSION)/nuodb-cp-rest-$(CP_CHARTS_VERSION).tgz

HELM_NGINX_RELEASE ?= ingress-nginx
NGINX_CHARTS_VERSION ?= 4.7.1
NGINX_CHART := https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-$(NGINX_CHARTS_VERSION)/ingress-nginx-$(NGINX_CHARTS_VERSION).tgz
NGINX_INGRESS_VERSION ?= 1.8.1

PROJECT_DIR := $(shell pwd)
BIN_DIR ?= 	$(PROJECT_DIR)/bin
TEST_RESULTS ?= $(BIN_DIR)/test-results

GOTESTSUM_VERSION ?= v1.11.0
GOTESTSUM_BIN := $(BIN_DIR)/gotestsum

PUBLISH_VERSION ?= 0.1.0
PUBLISH_DIR ?= $(BIN_DIR)/dist

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
deploy-cp: ## Deploy a local Control Plane to test with
	helm upgrade --install $(HELM_JETSTACK_RELEASE) $(JETSTACK_CHART) \
		--namespace cert-manager \
		--set installCRDs=true \
		--create-namespace

	helm upgrade --install $(HELM_NGINX_RELEASE) $(NGINX_CHART) \
		--set controller.image.tag="$(NGINX_INGRESS_VERSION)" \
		--set controller.ingressClassResource.default=true \
		--set controller.replicaCount=1 \
		--set controller.service.enablePorts.http=false \
		--set controller.service.nodePorts.https="30500" \
		--set controller.extraArgs.enable-ssl-passthrough=true

	helm upgrade --install $(HELM_CP_CRD_RELEASE) $(CP_CRD_CHART)

	@# Wait for all Control Plane dependencies to be ready
	kubectl -n $(HELM_JETSTACK_RELEASE) wait pod --all --for=condition=Ready
	kubectl -l app.kubernetes.io/instance="$(HELM_NGINX_RELEASE)" wait pod --for=condition=ready --timeout=120s

	helm upgrade --install $(HELM_CP_OPERATOR_RELEASE) $(CP_OPERATOR_CHART) \
		--set cpOperator.webhooks.enabled=true \
		--set cpOperator.samples.serviceTiers.enabled=false

	helm upgrade --install $(HELM_CP_REST_RELEASE) $(CP_REST_CHART) \
		--set cpRest.ingress.enabled=true \
		--set cpRest.ingress.className=nginx \
		--set cpRest.authentication.admin.create=true

	kubectl apply -f test/tiers.yaml

	kubectl wait pod -l app.kubernetes.io/instance="$(HELM_CP_OPERATOR_RELEASE)" --for=condition=ready --timeout=120s
	kubectl wait pod -l app.kubernetes.io/instance="$(HELM_CP_REST_RELEASE)" --for=condition=ready --timeout=120s

.PHONY: undeploy-cp
undeploy-cp: ## Uninstall a local Control Plane previously installed by this script
	# Clean up any left over DBaaS resources
	kubectl get databases.cp.nuodb.com -o name | xargs -r kubectl delete || $(IGNORE_NOT_FOUND)
	kubectl get domains.cp.nuodb.com -o name | xargs -r kubectl delete || $(IGNORE_NOT_FOUND)
	kubectl get servicetiers.cp.nuodb.com -o name | xargs -r kubectl delete || $(IGNORE_NOT_FOUND)
	kubectl get helmfeatures.cp.nuodb.com -o name | xargs -r kubectl delete || $(IGNORE_NOT_FOUND)
	kubectl get databasequotas.cp.nuodb.com -o name | xargs -r kubectl delete || $(IGNORE_NOT_FOUND)
	kubectl get pvc -o name --selector=group=nuodb | xargs -r kubectl delete || $(IGNORE_NOT_FOUND)

	# Uninstall DBaaS helm charts
	helm uninstall $(HELM_CP_REST_RELEASE) $(HELM_CP_OPERATOR_RELEASE) || $(IGNORE_NOT_FOUND)
	helm uninstall $(HELM_CP_CRD_RELEASE) $(HELM_NGINX_RELEASE) || $(IGNORE_NOT_FOUND)
	helm uninstall $(HELM_JETSTACK_RELEASE) --namespace cert-manager || $(IGNORE_NOT_FOUND)

	# Delete lease resources so that next time Cert Manager is deployed it does not have to wait for them to expire
	kubectl -n kube-system delete leases.coordination.k8s.io cert-manager-cainjector-leader-election --ignore-not-found=$(IGNORE_NOT_FOUND)
	kubectl -n kube-system delete leases.coordination.k8s.io cert-manager-controller --ignore-not-found=$(IGNORE_NOT_FOUND)

$(GOTESTSUM_BIN):
	$(call go-get-tool,gotest.tools/gotestsum@$(GOTESTSUM_VERSION))

.PHONY: discover-test
discover-test: ## Discover a local control plane and run tests against it
	$(eval HOST := $(or $(shell kubectl get service $(HELM_NGINX_RELEASE)-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}'), \
						$(shell kubectl get service $(HELM_NGINX_RELEASE)-controller -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')))
	$(eval PORT := $(shell kubectl get service $(HELM_NGINX_RELEASE)-controller -o jsonpath='{.spec.ports[?(@.appProtocol=="http")].port}'))

	@NUODB_CP_PASSWORD=$(shell kubectl get secret dbaas-user-system-admin -o jsonpath='{.data.password}' | base64 -d) \
		NUODB_CP_URL_BASE="http://$(HOST):$(PORT)/nuodb-cp" \
		NUODB_CP_USER="admin" \
		NUODB_CP_ORGANIZATION="system" \
		$(MAKE) testacc


.PHONY: testacc
testacc: $(GOTESTSUM_BIN) ## Run acceptance tests
	TF_ACC=1 $(GOTESTSUM_BIN) --junitfile $(TEST_RESULTS)/gotestsum-report.xml --format testname -- -v -timeout 30m $(TESTARGS) ./plugin/...

.PHONY: package
package:
	rm -r $(PUBLISH_DIR) || $(IGNORE_NOT_FOUND)
	mkdir -p $(PUBLISH_DIR)
	$(eval PACKAGE_OS ?= darwin linux windows)
	$(eval PACKAGE_ARCH ?= amd64 arm64)
	$(foreach OS, $(PACKAGE_OS), \
		$(foreach ARCH, $(PACKAGE_ARCH), $(call package-os,$(OS),$(ARCH))))
	$(eval PUBLISH_MIRROR ?= $(PUBLISH_DIR)/pkg_mirror/hashicorp.com/edu/nuodbaas)
	mkdir -p $(PUBLISH_MIRROR)
	cp $(PUBLISH_DIR)/*.zip $(PUBLISH_MIRROR)

# Build the release package for a given OS and architechture
define package-os
$(eval PUBLISH_STAGING := $(PUBLISH_DIR)/staging_$(1)_$(2))
$(eval PLUGIN_PKG := $(PUBLISH_DIR)/terraform-provider-nuodbaas_$(PUBLISH_VERSION)_$(1)_$(2).zip)
mkdir -p $(PUBLISH_STAGING)
GOOS=$(1) GOARCH=$(2) go build -C plugin -o $(PUBLISH_STAGING)/terraform-provider-nuodbaas_v$(PUBLISH_VERSION)
cd $(PUBLISH_STAGING) && zip $(PLUGIN_PKG) ./*
endef

# Use 'go install' to fetch tool, falling back to 'go get' if the tool is not
# made available locally. 'go install' should work in Golang versions >=1.16,
# while 'go get' should work in Golang <1.16.
define go-get-tool
@echo "Downloading $(1)"
export GO111MODULE=on GOBIN=$(BIN_DIR) && go install $(1) || go get $(1)
endef
