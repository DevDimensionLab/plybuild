#!/bin/sh
set -eu
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin"
export GOTOOLCHAIN=local GOFLAGS=-mod=readonly
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
test -x "$PLY_CANDIDATE_BINARY"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
export PLY_DELIVERY_TEST_BINARY="$PLY_CANDIDATE_BINARY"
# Real candidate CLI, Git and sealed records; modeled provider/Herdr only.
# Runtime approval and the later real reviewer test remain separate observations.
PLY_AL15_EVIDENCE_DIR=$(mktemp -d "$(dirname "$PLY_CANDIDATE_BINARY")/al17-cli.XXXXXX")
export PLY_AL15_EVIDENCE_DIR
export PLY_AL16_CLI_EVIDENCE="$PLY_AL15_EVIDENCE_DIR"
printf 'Raw CLI evidence (modeled provider/Herdr): %s\n' "$PLY_AL15_EVIDENCE_DIR"
log="$PLY_AL15_EVIDENCE_DIR/tests.stdout"
if go test -v -count=1 -timeout=180s ./internal/taskrun ./internal/workerskills -run '^(TestAL15GeneratedReviewTestMandate|TestAL15ReviewTestMandateThroughCandidateCLI|TestRV01ReviewerKeepsCandidateAndCannotPassADirtyWorktree|TestRV01ReviewerGoTestPreservesModuleMetadata|TestAL17ReviewerBundleScopesOwnerCallbacks|TestRoleBundlesContainExactSelectedBytesAndReferences|TestQualifiedCLIInstallationPublicConsumer|TestQualifiedCLIInstalledSmokePreservesNonzeroExit)$' >"$log" 2>"$PLY_AL15_EVIDENCE_DIR/tests.stderr"; then
 printf '0\n' >"$PLY_AL15_EVIDENCE_DIR/tests.exit"
else
 result=$?
 printf '%s\n' "$result" >"$PLY_AL15_EVIDENCE_DIR/tests.exit"
 cat "$log" "$PLY_AL15_EVIDENCE_DIR/tests.stderr"
 exit "$result"
fi
cat "$log" "$PLY_AL15_EVIDENCE_DIR/tests.stderr"
for name in TestAL15GeneratedReviewTestMandate TestAL15ReviewTestMandateThroughCandidateCLI TestRV01ReviewerKeepsCandidateAndCannotPassADirtyWorktree TestRV01ReviewerGoTestPreservesModuleMetadata TestAL17ReviewerBundleScopesOwnerCallbacks TestRoleBundlesContainExactSelectedBytesAndReferences TestQualifiedCLIInstallationPublicConsumer TestQualifiedCLIInstalledSmokePreservesNonzeroExit; do
 grep -q "^--- PASS: $name " "$log"
done
go vet ./internal/taskrun ./internal/workerskills ./cmd
git diff --check
test -z "$(git status --porcelain)"
