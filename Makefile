GO ?= go
BASH ?= /bin/bash

.DEFAULT_GOAL := all

.PHONY: all build docker-build docker-run docker-publish install run test test-agent-start test-install lint release release-brew upgrade

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

test: test-agent-start
	$(GO) test -v -cover ./...
	bash test/makefile_install_test.sh

test-agent-start:
	$(BASH) test/codex_dev_start_test.sh

test-install:
	bash test/makefile_install_test.sh

lint:
	gofmt -w pkg
	gofmt -w cmd

release:
	goreleaser release --clean

release-brew:
	goreleaser release --clean --skip=validate -f .goreleaser.brews.yml

upgrade:
	$(GO) get github.com/devdimensionlab/mvn-pom-mutator
	$(GO) get -u ./...
	$(GO) clean

all: build
