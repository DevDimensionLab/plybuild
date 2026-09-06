GO ?= go
GOFMT ?= gofmt
BASH ?= /bin/bash
GORELEASER ?= goreleaser
GOLANGCI_LINT_VERSION := 2.12.2
GOLANGCI_LINT ?= golangci-lint
APIDIFF_VERSION := v0.0.0-20260709172345-9ea1abe57597
APIDIFF ?= apidiff
DOCKER ?= docker
PYTHON ?= python3
REPO_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SCRIPTS_DIR ?= $(REPO_ROOT)/scripts
QUALITY_AUDIT ?= $(REPO_ROOT)/.quality/tools/quality-audit.sh
QUALITY_BASELINE ?= $(REPO_ROOT)/.quality/baseline/scorecard.json
QUALITY_MANUAL_EVIDENCE ?=
QUALITY_OUTPUT_ROOT ?=
QUALITY_GOCACHE ?=
QUALITY_GOMODCACHE ?=
override PLY_QUALITY_MUTATIONS := cli-context config-cloud maven-sorting template file-shell spring http interactive-build

.DEFAULT_GOAL := all

.PHONY: acceptance acceptance-docker acceptance-snapshot all build compat-api compat-cli compatibility docker-build docker-run docker-publish format install lint preflight quality release release-brew run snapshot test test-agent-start test-cli-surface test-compatibility test-distribution test-install test-lint test-preflight test-quality test-toolchain upgrade

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

test: test-agent-start test-distribution test-lint test-preflight test-quality test-toolchain
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

test-quality:
	$(BASH) test/makefile_quality_test.sh

test-toolchain:
	$(BASH) test/toolchain_declarations_test.sh

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

acceptance-docker:
	$(BASH) "$(SCRIPTS_DIR)/accept-docker"

preflight: compatibility
	$(GO) build ./...
	$(GO) test ./... -count=1
	$(GO) vet ./...
	$(MAKE) --no-print-directory -C "$(REPO_ROOT)" lint
	$(BASH) "$(REPO_ROOT)/test/codex_dev_start_test.sh"
	$(BASH) "$(REPO_ROOT)/test/makefile_distribution_test.sh"
	$(BASH) "$(REPO_ROOT)/test/makefile_lint_test.sh"
	$(BASH) "$(REPO_ROOT)/test/makefile_install_test.sh"
	$(BASH) "$(REPO_ROOT)/test/toolchain_declarations_test.sh"
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

quality:
	@set -eu; \
	fail() { printf 'quality: %s\n' "$$*" >&2; exit 1; }; \
	repo_root='$(REPO_ROOT)'; \
	require_external_tool() { \
		label=$$1; path=$$2; \
		[ -n "$$path" ] || fail "$$label path is required"; \
		case $$path in /*) ;; *) fail "$$label path must be absolute" ;; esac; \
		[ -f "$$path" ] && [ -x "$$path" ] && [ ! -L "$$path" ] || \
			fail "$$label path is not a regular executable"; \
		physical=$$(realpath "$$path") || fail "cannot resolve $$label path"; \
		case $$physical in "$$repo_root"|"$$repo_root"/*) fail "$$label path must be outside the repository" ;; esac; \
	}; \
	require_external_dir() { \
		label=$$1; path=$$2; \
		[ -n "$$path" ] || fail "$$label path is required"; \
		case $$path in /*) ;; *) fail "$$label path must be absolute" ;; esac; \
		[ -d "$$path" ] && [ ! -L "$$path" ] || fail "$$label path is not a regular directory"; \
		physical=$$(cd "$$path" && pwd -P) || fail "cannot resolve $$label path"; \
		case $$physical in "$$repo_root"|"$$repo_root"/*) fail "$$label path must be outside the repository" ;; esac; \
	}; \
	for tool in \
		'go:$(GO)' 'bash:$(BASH)' 'goreleaser:$(GORELEASER)' \
		'apidiff:$(APIDIFF)' 'golangci-lint:$(GOLANGCI_LINT)' \
		'docker:$(DOCKER)' 'python:$(PYTHON)'; do \
		require_external_tool "$${tool%%:*}" "$${tool#*:}"; \
	done; \
	require_external_dir GOCACHE '$(QUALITY_GOCACHE)'; \
	require_external_dir GOMODCACHE '$(QUALITY_GOMODCACHE)'; \
	evidence='$(QUALITY_MANUAL_EVIDENCE)'; \
	[ -n "$$evidence" ] || fail 'QUALITY_MANUAL_EVIDENCE is required'; \
	case $$evidence in /*) ;; *) fail 'manual evidence path must be absolute' ;; esac; \
	[ -f "$$evidence" ] && [ -s "$$evidence" ] && [ ! -L "$$evidence" ] || \
		fail 'manual evidence must be a non-empty regular file'; \
	evidence=$$(realpath "$$evidence") || fail 'cannot resolve manual evidence path'; \
	case $$evidence in "$$repo_root"|"$$repo_root"/*) fail 'manual evidence must be outside the repository' ;; esac; \
	output='$(QUALITY_OUTPUT_ROOT)'; \
	[ -n "$$output" ] || fail 'QUALITY_OUTPUT_ROOT is required'; \
	case $$output in /*) ;; *) fail 'quality output root must be absolute' ;; esac; \
	case $$output in *[![:graph:]]*) fail 'quality output root cannot contain whitespace' ;; esac; \
	[ ! -e "$$output" ] && [ ! -L "$$output" ] || fail 'quality output root is stale or pre-existing'; \
	parent=$$(dirname "$$output"); name=$$(basename "$$output"); \
	[ -d "$$parent" ] && [ ! -L "$$parent" ] || fail 'quality output parent is not a regular directory'; \
	parent=$$(cd "$$parent" && pwd -P) || fail 'cannot resolve quality output parent'; \
	output="$$parent/$$name"; \
	case $$output in "$$repo_root"|"$$repo_root"/*) fail 'quality output root must be outside the repository' ;; esac; \
	[ -f '$(QUALITY_AUDIT)' ] && [ ! -L '$(QUALITY_AUDIT)' ] || fail 'quality audit is not a regular file'; \
	[ -f '$(QUALITY_BASELINE)' ] && [ -s '$(QUALITY_BASELINE)' ] && [ ! -L '$(QUALITY_BASELINE)' ] || \
		fail 'quality baseline is not a non-empty regular file'; \
	mutation_count=0; seen=' '; \
	for name in $(PLY_QUALITY_MUTATIONS); do \
		case "$$seen" in *" $$name "*) fail "duplicate mutation stage: $$name" ;; esac; \
		seen="$$seen$$name "; mutation_count=$$((mutation_count + 1)); \
		for script in "$(SCRIPTS_DIR)/mutate-$$name" "$(SCRIPTS_DIR)/test-mutate-$$name"; do \
			[ -f "$$script" ] && [ -x "$$script" ] && [ ! -L "$$script" ] || \
				fail "required mutation stage is not a regular executable: $$script"; \
		done; \
	done; \
	[ "$$mutation_count" -eq 8 ] || fail 'mutation stage population must contain exactly 8 subjects'; \
	acceptance_count=0; \
	for name in install status upgrade build; do \
		for script in "$(SCRIPTS_DIR)/verify-$$name" "$(SCRIPTS_DIR)/test-verify-$$name"; do \
			[ -f "$$script" ] && [ -x "$$script" ] && [ ! -L "$$script" ] || \
				fail "required acceptance stage is not a regular executable: $$script"; \
		done; \
		acceptance_count=$$((acceptance_count + 1)); \
	done; \
	[ "$$acceptance_count" -eq 4 ] || fail 'acceptance stage population must contain exactly 4 features'; \
	for name in snapshot docker; do \
		for script in "$(SCRIPTS_DIR)/accept-$$name" "$(SCRIPTS_DIR)/test-accept-$$name"; do \
			[ -f "$$script" ] && [ -x "$$script" ] && [ ! -L "$$script" ] || \
				fail "required artifact stage is not a regular executable: $$script"; \
		done; \
	done; \
	mkdir "$$output" "$$output/logs" "$$output/tmp" \
		"$$output/golangci-lint-cache" || \
		fail 'cannot initialize fresh external quality output'; \
	: >"$$output/stage-ledger.txt"
	@set -eu; log='$(QUALITY_OUTPUT_ROOT)/logs/preflight.log'; \
	if ! TMPDIR='$(QUALITY_OUTPUT_ROOT)/tmp' GOCACHE='$(QUALITY_GOCACHE)' GOMODCACHE='$(QUALITY_GOMODCACHE)' \
		GOLANGCI_LINT_CACHE='$(QUALITY_OUTPUT_ROOT)/golangci-lint-cache' \
		API_COMPAT_REPORT_OUT='$(QUALITY_OUTPUT_ROOT)/api-compatibility.json' \
		CLI_COMPAT_REPORT_OUT='$(QUALITY_OUTPUT_ROOT)/cli-compatibility.json' \
		GO='$(GO)' BASH='$(BASH)' GORELEASER='$(GORELEASER)' APIDIFF='$(APIDIFF)' \
		GOLANGCI_LINT='$(GOLANGCI_LINT)' SCRIPTS_DIR='$(SCRIPTS_DIR)' \
		$(MAKE) --no-print-directory -C "$(REPO_ROOT)" MAKEOVERRIDES= preflight >"$$log" 2>&1; then \
		sed -n '1,240p' "$$log" >&2; exit 1; \
	fi; \
	printf 'preflight\n' >>'$(QUALITY_OUTPUT_ROOT)/stage-ledger.txt'
	@set -eu; \
	for name in $(PLY_QUALITY_MUTATIONS); do \
		log='$(QUALITY_OUTPUT_ROOT)/logs/mutation-meta-'$$name'.log'; \
		if ! TMPDIR='$(QUALITY_OUTPUT_ROOT)/tmp' GOCACHE='$(QUALITY_GOCACHE)' \
			GOMODCACHE='$(QUALITY_GOMODCACHE)' '$(BASH)' "$(SCRIPTS_DIR)/test-mutate-$$name" >"$$log" 2>&1; then \
			sed -n '1,240p' "$$log" >&2; exit 1; \
		fi; \
		[ "$$(grep -Fxc "test-mutate-$$name: PASS (T1-T10, declared=10 killed=10 survived=0 unusable=0)" "$$log")" -eq 1 ] || { \
			printf 'quality: mutation meta-stage has no exact non-empty receipt: %s\n' "$$name" >&2; exit 1; \
		}; \
		printf 'mutation-meta:%s\n' "$$name" >>'$(QUALITY_OUTPUT_ROOT)/stage-ledger.txt'; \
	done
	@set -eu; \
	for name in $(PLY_QUALITY_MUTATIONS); do \
		upper=$$(printf '%s' "$$name" | tr '[:lower:]-' '[:upper:]_'); \
		log='$(QUALITY_OUTPUT_ROOT)/logs/mutation-'$$name'.log'; \
		if ! env "MUTATION_$${upper}_KEEP_WORK=0" "MUTATION_$${upper}_WORK_ROOT=" \
			TMPDIR='$(QUALITY_OUTPUT_ROOT)/tmp' GOCACHE='$(QUALITY_GOCACHE)' \
			GOMODCACHE='$(QUALITY_GOMODCACHE)' '$(BASH)' "$(SCRIPTS_DIR)/mutate-$$name" >"$$log" 2>&1; then \
			sed -n '1,240p' "$$log" >&2; exit 1; \
		fi; \
		[ "$$(grep -Fxc 'declared=10 killed=10 survived=0 unusable=0' "$$log")" -eq 1 ] || { \
			printf 'quality: mutation stage has no exact non-empty receipt: %s\n' "$$name" >&2; exit 1; \
		}; \
		printf 'mutation:%s\n' "$$name" >>'$(QUALITY_OUTPUT_ROOT)/stage-ledger.txt'; \
	done
	@set -eu; log='$(QUALITY_OUTPUT_ROOT)/logs/acceptance-host.log'; \
	if ! TMPDIR='$(QUALITY_OUTPUT_ROOT)/tmp' GOCACHE='$(QUALITY_GOCACHE)' GOMODCACHE='$(QUALITY_GOMODCACHE)' \
		$(MAKE) --no-print-directory -C "$(REPO_ROOT)" acceptance GO='$(GO)' BASH='$(BASH)' \
			SCRIPTS_DIR='$(SCRIPTS_DIR)' >"$$log" 2>&1; then \
		sed -n '1,240p' "$$log" >&2; exit 1; \
	fi; \
	for name in install status upgrade build; do \
		[ "$$(grep -Fxc "verify-$$name: PASS" "$$log")" -eq 1 ] || { \
			printf 'quality: host acceptance has no exact receipt: %s\n' "$$name" >&2; exit 1; \
		}; \
	done; \
	printf 'acceptance:host\n' >>'$(QUALITY_OUTPUT_ROOT)/stage-ledger.txt'
	@set -eu; log='$(QUALITY_OUTPUT_ROOT)/logs/acceptance-snapshot.log'; \
	if ! TMPDIR='$(QUALITY_OUTPUT_ROOT)/tmp' GOCACHE='$(QUALITY_GOCACHE)' GOMODCACHE='$(QUALITY_GOMODCACHE)' \
		PLY_SNAPSHOT_EVIDENCE_ROOT='$(QUALITY_OUTPUT_ROOT)/snapshot-evidence' \
		$(MAKE) --no-print-directory -C "$(REPO_ROOT)" acceptance-snapshot GO='$(GO)' BASH='$(BASH)' \
			GORELEASER='$(GORELEASER)' SCRIPTS_DIR='$(SCRIPTS_DIR)' >"$$log" 2>&1; then \
		sed -n '1,240p' "$$log" >&2; exit 1; \
	fi; \
	[ -s '$(QUALITY_OUTPUT_ROOT)/snapshot-evidence/report.txt' ] && \
	[ "$$(grep -Fxc 'snapshot-acceptance: PASS' '$(QUALITY_OUTPUT_ROOT)/snapshot-evidence/report.txt')" -eq 1 ] || { \
		printf '%s\n' 'quality: snapshot acceptance produced no fresh exact receipt' >&2; exit 1; \
	}; \
	printf 'acceptance:snapshot\n' >>'$(QUALITY_OUTPUT_ROOT)/stage-ledger.txt'
	@set -eu; log='$(QUALITY_OUTPUT_ROOT)/logs/acceptance-docker.log'; \
	if ! TMPDIR='$(QUALITY_OUTPUT_ROOT)/tmp' GOCACHE='$(QUALITY_GOCACHE)' GOMODCACHE='$(QUALITY_GOMODCACHE)' \
		PLY_DOCKER_EVIDENCE_ROOT='$(QUALITY_OUTPUT_ROOT)/docker-evidence' PLY_DOCKER_BIN='$(DOCKER)' \
		$(MAKE) --no-print-directory -C "$(REPO_ROOT)" acceptance-docker GO='$(GO)' BASH='$(BASH)' \
			SCRIPTS_DIR='$(SCRIPTS_DIR)' >"$$log" 2>&1; then \
		sed -n '1,240p' "$$log" >&2; exit 1; \
	fi; \
	[ -s '$(QUALITY_OUTPUT_ROOT)/docker-evidence/report.txt' ] && \
	[ "$$(grep -Fxc 'docker-acceptance: PASS' '$(QUALITY_OUTPUT_ROOT)/docker-evidence/report.txt')" -eq 1 ] || { \
		printf '%s\n' 'quality: Docker acceptance produced no fresh exact receipt' >&2; exit 1; \
	}; \
	printf 'acceptance:docker\n' >>'$(QUALITY_OUTPUT_ROOT)/stage-ledger.txt'
	@set -eu; log='$(QUALITY_OUTPUT_ROOT)/logs/audit-q0-q2.log'; \
	if ! TMPDIR='$(QUALITY_OUTPUT_ROOT)/tmp' GOCACHE='$(QUALITY_GOCACHE)' GOMODCACHE='$(QUALITY_GOMODCACHE)' \
		'$(BASH)' '$(QUALITY_AUDIT)' '$(REPO_ROOT)' --baseline '$(QUALITY_BASELINE)' \
			--manual-evidence '$(QUALITY_MANUAL_EVIDENCE)' --out '$(QUALITY_OUTPUT_ROOT)/audit-q0-q2' \
			--only 'Q0.*,Q1.*,Q2.*' >"$$log" 2>&1; then \
		sed -n '1,240p' "$$log" >&2; exit 1; \
	fi; \
	printf 'audit:Q0.*,Q1.*,Q2.*\n' >>'$(QUALITY_OUTPUT_ROOT)/stage-ledger.txt'
	@'$(PYTHON)' -c 'import json,sys; d=json.load(open(sys.argv[1], encoding="utf-8")); expected={*("Q0.%d"%n for n in range(1,9)), *("Q1.%d"%n for n in range(1,10)), *("Q2.%d"%n for n in range(1,11))}; criteria=d.get("criteria",[]); actual={c.get("id") for c in criteria}; den=d.get("denominators",{}); ratchet=d.get("ratchet",{}); manual=d.get("manual_evidence",{}); tree=d.get("repository",{}).get("tree",{}); ok=(len(criteria)==27 and actual==expected and all(c.get("verdict")=="PASS" for c in criteria) and d.get("attained_level")=="L2" and den.get("declared_subjects")==8 and den.get("mutation_harnesses")==8 and den.get("declared_features")==4 and den.get("acceptance_scripts")==4 and den.get("test_functions",0)>0 and manual.get("status")=="valid" and manual.get("criterion_receipts")==6 and ratchet.get("held")==0 and ratchet.get("regressed")==0 and ratchet.get("current_not_comparable")==0 and tree.get("measurement_clean") is True and tree.get("dirty_paths")==[]); sys.exit(0 if ok else 1)' '$(QUALITY_OUTPUT_ROOT)/audit-q0-q2/scorecard.json' || { \
		printf '%s\n' 'quality: authoritative Q0-Q2 scorecard did not attain clean L2' >&2; exit 1; \
	}
	@set -eu; expected='$(QUALITY_OUTPUT_ROOT)/expected-stage-ledger.txt'; \
	{ \
		printf 'preflight\n'; \
		for name in $(PLY_QUALITY_MUTATIONS); do printf 'mutation-meta:%s\n' "$$name"; done; \
		for name in $(PLY_QUALITY_MUTATIONS); do printf 'mutation:%s\n' "$$name"; done; \
		printf 'acceptance:host\nacceptance:snapshot\nacceptance:docker\naudit:Q0.*,Q1.*,Q2.*\n'; \
	} >"$$expected"; \
	cmp -s "$$expected" '$(QUALITY_OUTPUT_ROOT)/stage-ledger.txt' || { \
		printf '%s\n' 'quality: required stage population is missing, duplicated, or out of order' >&2; exit 1; \
	}; \
	for name in tmp mutations golangci-lint-cache; do \
		path='$(QUALITY_OUTPUT_ROOT)/'$$name; \
		[ ! -L "$$path" ] || { printf 'quality: refusing transient symlink: %s\n' "$$path" >&2; exit 1; }; \
		if [ -d "$$path" ]; then \
			find "$$path" -type d -exec chmod u+rwx {} + 2>/dev/null || \
				{ printf 'quality: cannot make transient tree removable: %s\n' "$$path" >&2; exit 1; }; \
			rm -rf "$$path" || { printf 'quality: cannot remove transient tree: %s\n' "$$path" >&2; exit 1; }; \
		fi; \
	done; \
	printf 'quality: PASS (Q0-Q2 attained L2)\n'

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
