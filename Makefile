GO ?= go
GOFMT ?= gofmt
BASH ?= /bin/bash
GORELEASER ?= goreleaser
GOLANGCI_LINT_VERSION := 2.12.2
APIDIFF_VERSION := v0.0.0-20260709172345-9ea1abe57597
APIDIFF ?= apidiff
REPO_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SCRIPTS_DIR ?= $(REPO_ROOT)/scripts

.DEFAULT_GOAL := all

.PHONY: acceptance acceptance-snapshot all build compat-api compat-cli compatibility docker-build docker-run docker-publish format install lint preflight release release-brew run snapshot test test-agent-start test-cli-surface test-compatibility test-distribution test-install test-lint test-preflight upgrade

build:
	$(GO) build -o ply ./cmd/ply

docker-build:
	docker build --tag ply:latest .

docker-run:
	docker run ply $(ARGS)

docker-publish:
	./docker-publish.sh

install:
	$(GO) install ./cmd/ply

run:
	$(GO) run ./cmd/ply

test: test-agent-start test-distribution test-lint test-preflight
	$(GO) test -v -cover ./...
	bash test/makefile_install_test.sh

test-agent-start:
	$(BASH) test/codex_dev_start_test.sh

test-distribution:
	$(BASH) test/makefile_distribution_test.sh

test-install:
	bash test/makefile_install_test.sh

test-lint:
	$(BASH) test/makefile_lint_test.sh

test-preflight:
	$(BASH) test/makefile_preflight_test.sh

lint:
	@lint_bin="$${GOLANGCI_LINT:-}"; \
	if [ -z "$$lint_bin" ]; then \
		lint_bin="$$(command -v golangci-lint 2>/dev/null || true)"; \
	fi; \
	if [ -z "$$lint_bin" ] || [ ! -f "$$lint_bin" ] || [ ! -x "$$lint_bin" ]; then \
		printf '%s\n' \
			'golangci-lint $(GOLANGCI_LINT_VERSION) is required; install the official binary with:' \
			'  curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b "$$(go env GOPATH)/bin" v$(GOLANGCI_LINT_VERSION)' >&2; \
		exit 1; \
	fi; \
	if ! actual_version="$$(cd "$(REPO_ROOT)" && \
		GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local "$$lint_bin" version --short 2>/dev/null)"; then \
		printf 'golangci-lint version check failed for %s\n' "$$lint_bin" >&2; \
		exit 1; \
	fi; \
	if [ "$$actual_version" != "$(GOLANGCI_LINT_VERSION)" ]; then \
		printf 'golangci-lint version mismatch: expected %s, got %s\n' \
			'$(GOLANGCI_LINT_VERSION)' "$$actual_version" >&2; \
		exit 1; \
	fi; \
	cd "$(REPO_ROOT)" && \
		GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local "$$lint_bin" config verify --config .golangci.yml && \
		GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local "$$lint_bin" run --config .golangci.yml --modules-download-mode=readonly ./...

format:
	cd "$(REPO_ROOT)" && find . -path './vendor' -prune -o -type f -name '*.go' -exec "$(GOFMT)" -w {} +

compat-api:
	APIDIFF="$(APIDIFF)" APIDIFF_VERSION="$(APIDIFF_VERSION)" \
		$(BASH) "$(REPO_ROOT)/scripts/check-api-compat.sh"

compat-cli:
	$(BASH) "$(REPO_ROOT)/scripts/check-cli-compat.sh"
	$(BASH) "$(REPO_ROOT)/test/cli_surface_contract_test.sh"

compatibility: compat-api compat-cli

test-compatibility:
	$(BASH) "$(REPO_ROOT)/scripts/test-check-api-compat.sh"
	$(BASH) "$(REPO_ROOT)/scripts/test-check-cli-compat.sh"

test-cli-surface:
	$(BASH) "$(REPO_ROOT)/test/cli_surface_contract_test.sh"

acceptance:
	$(BASH) "$(REPO_ROOT)/scripts/verify-install"
	$(BASH) "$(REPO_ROOT)/scripts/verify-status"
	$(BASH) "$(REPO_ROOT)/scripts/verify-upgrade"
	$(BASH) "$(REPO_ROOT)/scripts/verify-build"

acceptance-snapshot:
	PLY_SNAPSHOT_GORELEASER="$(GORELEASER)" \
		$(BASH) "$(SCRIPTS_DIR)/accept-snapshot"

preflight: compatibility
	$(GO) build ./...
	$(GO) test ./... -count=1
	$(GO) vet ./...
	$(MAKE) --no-print-directory -C "$(REPO_ROOT)" lint
	$(BASH) "$(REPO_ROOT)/test/codex_dev_start_test.sh"
	$(BASH) "$(REPO_ROOT)/test/makefile_distribution_test.sh"
	$(BASH) "$(REPO_ROOT)/test/makefile_lint_test.sh"
	$(BASH) "$(REPO_ROOT)/test/makefile_install_test.sh"
	@set -eu; \
	scripts_dir="$(SCRIPTS_DIR)"; \
	if [ ! -d "$$scripts_dir" ] || [ -L "$$scripts_dir" ]; then \
		printf '%s\n' 'preflight: script directory does not exist or is not a regular directory' >&2; \
		exit 1; \
	fi; \
	production_count=0; \
	for script in "$$scripts_dir"/*; do \
		[ -f "$$script" ] || continue; \
		name=$${script##*/}; \
		case $$name in test-*) continue ;; esac; \
		if [ -L "$$script" ]; then \
			printf 'preflight: production script is a symlink: %s\n' "$$name" >&2; \
			exit 1; \
		fi; \
		production_count=$$((production_count + 1)); \
		if [ ! -f "$$scripts_dir/test-$$name" ] || [ -L "$$scripts_dir/test-$$name" ]; then \
			printf 'preflight: missing meta-test for %s\n' "$$name" >&2; \
			exit 1; \
		fi; \
	done; \
	if [ "$$production_count" -eq 0 ]; then \
		printf '%s\n' 'preflight: production script population is empty' >&2; \
		exit 1; \
	fi; \
	test_count=0; \
	for test_script in "$$scripts_dir"/test-*; do \
		[ -f "$$test_script" ] || continue; \
		name=$${test_script##*/}; \
		production_name=$${name#test-}; \
		if [ -L "$$test_script" ]; then \
			printf 'preflight: meta-test is a symlink: %s\n' "$$name" >&2; \
			exit 1; \
		fi; \
		if [ ! -f "$$scripts_dir/$$production_name" ] || [ -L "$$scripts_dir/$$production_name" ]; then \
			printf 'preflight: meta-test has no production script: %s\n' "$$name" >&2; \
			exit 1; \
		fi; \
		test_count=$$((test_count + 1)); \
	done; \
	if [ "$$test_count" -ne "$$production_count" ]; then \
		printf 'preflight: script population mismatch: %s production, %s meta-tests\n' \
			"$$production_count" "$$test_count" >&2; \
		exit 1; \
	fi; \
	for test_script in "$$scripts_dir"/test-*; do \
		$(BASH) "$$test_script"; \
	done
	$(BASH) "$(REPO_ROOT)/.quality/tools/test-quality-audit.sh"

snapshot:
	GITHUB_TOKEN= GITLAB_TOKEN= GITEA_TOKEN= HOMEBREW_TAP_GITHUB_TOKEN= SNAPCRAFT_STORE_CREDENTIALS= \
		"$(GORELEASER)" release --snapshot --clean --skip=publish

release: snapshot

release-brew:
	@printf '%s\n' 'Homebrew distribution is inactive; release-brew is disabled.' >&2
	@false

upgrade:
	$(GO) get github.com/devdimensionlab/mvn-pom-mutator
	$(GO) get -u ./...
	$(GO) clean

all: build
