GO ?= go
GOFMT ?= gofmt
BASH ?= /bin/bash
GOLANGCI_LINT_VERSION := 2.12.2
REPO_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

.DEFAULT_GOAL := all

.PHONY: all build docker-build docker-run docker-publish format install lint release release-brew run test test-agent-start test-install test-lint upgrade

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

test: test-agent-start test-lint
	$(GO) test -v -cover ./...
	bash test/makefile_install_test.sh

test-agent-start:
	$(BASH) test/codex_dev_start_test.sh

test-install:
	bash test/makefile_install_test.sh

test-lint:
	$(BASH) test/makefile_lint_test.sh

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

release:
	goreleaser release --clean

release-brew:
	goreleaser release --clean --skip=validate -f .goreleaser.brews.yml

upgrade:
	$(GO) get github.com/devdimensionlab/mvn-pom-mutator
	$(GO) get -u ./...
	$(GO) clean

all: build
