PROJECT_DIR := $(shell pwd)
BIN_DIR ?= $(PROJECT_DIR)/bin
export TEST_RESULTS ?= $(PROJECT_DIR)/test-results
export TMP_DIR ?= $(PROJECT_DIR)/tmp
export PATH := $(BIN_DIR):$(PATH)
export OUTPUT_DIR ?= $(TMP_DIR)/test-artifacts

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

TERRAFORM_VERSION ?= 1.7.3
TOFU_VERSION ?= 1.7.1
KUBECTL_VERSION ?= 1.28.3
KWOKCTL_VERSION ?= 0.5.1
HELM_VERSION ?= 3.14.3
MINIKUBE_VERSION ?= 1.32.0
NUODB_CP_VERSION ?= 2.6.0

GOTESTSUM := bin/gotestsum
TFPLUGINDOCS := bin/tfplugindocs
OAPI_CODEGEN := bin/oapi-codegen
TERRAFORM := bin/terraform
TOFU := bin/tofu
KUBECTL := bin/kubectl
KWOKCTL := bin/kwokctl
HELM := bin/helm
MINIKUBE := bin/minikube
GOLANGCI_LINT := bin/golangci-lint
NUODB_CP := bin/nuodb-cp

# For actual releases, GoReleaser uses the Git tag to obtain the version and
# not this variable, but this is used by the `make package` target which is
# invoked by the e2e app test. Scrape the value from the main.go file, which is
# also overridden by the Git tag, so that we do not specify the same value in
# multiple places.
PUBLISH_VERSION ?= $(shell sed -n 's|^\t*version string *= "\([^"]*\)" // {{version}}|\1|p' main.go)
PUBLISH_DIR ?= $(PROJECT_DIR)/dist

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

##@ Dependencies

$(TERRAFORM):
	mkdir -p bin tmp
	curl -L -s https://releases.hashicorp.com/terraform/$(TERRAFORM_VERSION)/terraform_$(TERRAFORM_VERSION)_$(OS)_$(ARCH).zip -o tmp/terraform.zip
	cd tmp && unzip terraform.zip
	mv tmp/terraform $(TERRAFORM)

$(TOFU):
	mkdir -p bin tmp
	curl -L -s https://github.com/opentofu/opentofu/releases/download/v$(TOFU_VERSION)/tofu_$(TOFU_VERSION)_$(OS)_$(ARCH).zip -o tmp/tofu.zip
	cd tmp && unzip tofu.zip tofu
	mv tmp/tofu $(TOFU)

$(KUBECTL):
	mkdir -p bin
	curl -L -s https://storage.googleapis.com/kubernetes-release/release/v$(KUBECTL_VERSION)/bin/$(OS)/$(ARCH)/kubectl -o $(KUBECTL)
	chmod +x $(KUBECTL)

$(KWOKCTL):
	mkdir -p bin
	curl -L -s https://github.com/kubernetes-sigs/kwok/releases/download/v$(KWOKCTL_VERSION)/kwokctl-$(OS)-$(ARCH) -o $(KWOKCTL)
	chmod +x $(KWOKCTL)

$(HELM):
	mkdir -p bin
	curl -L -s https://get.helm.sh/helm-v$(HELM_VERSION)-$(OS)-$(ARCH).tar.gz | tar -xz -C bin --strip-components=1 $(OS)-$(ARCH)/helm
	chmod +x $(HELM)

$(MINIKUBE):
	mkdir -p bin
	curl -L -s https://storage.googleapis.com/minikube/releases/v$(MINIKUBE_VERSION)/minikube-$(OS)-$(ARCH) -o $(MINIKUBE)
	chmod +x $(MINIKUBE)

$(NUODB_CP):
	mkdir -p bin
	curl -L -s https://github.com/nuodb/nuodb-cp-releases/releases/download/v$(NUODB_CP_VERSION)/nuodb-cp -o $(NUODB_CP)
	chmod +x $(NUODB_CP)

bin/%:
	$(MAKE) install-tools

.PHONY: install-tools
install-tools: ## Install tools declared as dependencies in tools.go
	@echo "Installing build tools declared in tools.go..."
	@go list -e -f '{{range .Imports}}{{.}} {{end}}' tools.go | GOBIN=$(BIN_DIR) xargs go install

##@ Development

.PHONY: check-no-changes
check-no-changes: ## Check that there are no uncommitted changes
	$(eval GIT_STATUS := $(shell git status --porcelain))
	@[ "$(GIT_STATUS)" = "" ] || ( echo "There are uncommitted changes:\n$(GIT_STATUS)"; exit 1; )

.PHONY: generate
generate: $(TFPLUGINDOCS) $(OAPI_CODEGEN) $(TERRAFORM) ## Generate Golang client for the NuoDB REST API and Terraform provider documentation
	@if git tag --points-at HEAD | grep -q "^v"; then \
		echo "Updating openapi.yaml because commit has version tag..." ;\
		$(MAKE) update-spec ;\
	fi
	go mod tidy
	go generate

.PHONY: update-spec
update-spec: ## Update spec to released Control Plane version
	curl -s https://raw.githubusercontent.com/nuodb/nuodb-cp-releases/v$(NUODB_CP_VERSION)/openapi.yaml -o openapi.yaml

.PHONY: lint
lint: $(GOLANGCI_LINT) ## Run linters to check code quality and find for common errors
	$(GOLANGCI_LINT) run

##@ Testing

.PHONY: kwok-deps
kwok-deps: $(KWOKCTL) $(KUBECTL) $(HELM)

.PHONY: k8s-deps
k8s-deps: $(KUBECTL) $(HELM) $(NUODB_CP)

.PHONY: minikube-deps
minikube-deps: $(MINIKUBE) k8s-deps

.PHONY: external-deps

setup-%: %-deps
	mkdir -p $(OUTPUT_DIR) $(TMP_DIR)
	[ ! -x "./deploy/$*/setup.sh" ] || ./deploy/$*/setup.sh

env-%: %-deps
	@[ ! -x "./deploy/$*/env.sh" ] || ./deploy/$*/env.sh

logs-%: %-deps
	[ ! -x "./deploy/$*/logs.sh" ] || ./deploy/$*/logs.sh

teardown-%: %-deps
	[ ! -x "./deploy/$*/teardown.sh" ] || ./deploy/$*/teardown.sh

.PHONY: testacc
testacc: $(GOTESTSUM) $(TERRAFORM) ## Run acceptance tests
	@if [ "$$USE_TOFU" = true ]; then \
		$(MAKE) $(TOFU) ;\
	fi
	mkdir -p $(TEST_RESULTS)
	TF_ACC=1 $(GOTESTSUM) --junitfile $(TEST_RESULTS)/gotestsum-report.xml \
		   --format testname -- -v -count=1 -p 1 -timeout 30m \
		   -coverprofile $(TEST_RESULTS)/cover.out -coverpkg ./internal/... \
		   $(TESTARGS) ./...

.PHONY: coverage-report
coverage-report:
	go tool cover -html=$(TEST_RESULTS)/cover.out -o $(OUTPUT_DIR)/coverage.html
	go tool cover -func $(TEST_RESULTS)/cover.out -o $(OUTPUT_DIR)/coverage.txt

##@ Packaging

.PHONY: package
package: ## Generate the provider for this machines OS and Architecture
	PACKAGE_OS=$(OS) PACKAGE_ARCH=$(ARCH) $(MAKE) package-all

.PHONY: package-all
package-all: ## Generate the provider for every OS and Architecture
	rm -r $(PUBLISH_DIR) || $(IGNORE_NOT_FOUND)
	mkdir -p $(PUBLISH_DIR)
	$(eval PACKAGE_OS ?= darwin linux windows)
	$(eval PACKAGE_ARCH ?= amd64 arm64)
	$(foreach OS, $(PACKAGE_OS), \
		$(foreach ARCH, $(PACKAGE_ARCH), $(call package-os,$(OS),$(ARCH))))
	$(eval PUBLISH_MIRROR ?= $(PUBLISH_DIR)/pkg_mirror/registry.terraform.io/nuodb/nuodbaas)
	mkdir -p $(PUBLISH_MIRROR)
	cp $(PUBLISH_DIR)/*.zip $(PUBLISH_MIRROR)

# Build the release package for a given OS and Architecture
define package-os
$(eval PUBLISH_STAGING := $(PUBLISH_DIR)/staging_$(1)_$(2))
$(eval PLUGIN_PKG := $(PUBLISH_DIR)/terraform-provider-nuodbaas_$(PUBLISH_VERSION)_$(1)_$(2).zip)
mkdir -p $(PUBLISH_STAGING)
GOOS=$(1) GOARCH=$(2) go build -o $(PUBLISH_STAGING)/terraform-provider-nuodbaas_v$(PUBLISH_VERSION)
cd $(PUBLISH_STAGING) && zip $(PLUGIN_PKG) ./*
endef
