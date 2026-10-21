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
evidence=$(mktemp -d "$(dirname "$PLY_CANDIDATE_BINARY")/al19-acceptance.XXXXXX")
printf 'Raw evidence (actual candidate CLI; modeled provider/Herdr; no Slack): %s\n' "$evidence"
capture() {
 name=$1
 shift
 if "$@" >"$evidence/$name.stdout" 2>"$evidence/$name.stderr"; then
  printf '0\n' >"$evidence/$name.exit"
 else
  result=$?
  printf '%s\n' "$result" >"$evidence/$name.exit"
  cat "$evidence/$name.stdout" "$evidence/$name.stderr"
  exit "$result"
 fi
 cat "$evidence/$name.stdout" "$evidence/$name.stderr"
}
# Receipt/idempotence coverage stops before unrelated local cleanup and writer inspection.
capture tests go test -v -count=1 -timeout=180s ./internal/taskrun ./cmd -run '^(TestDeliveryStartupCLI|TestDeliveryObservedStopRequiresOwnerResolutionBeforeProgress|TestDeliveryResolveStopContinuesNamedChangedActionWithoutPermissionClaim|TestDeliveryReceiptResolutionPreservesUncertainExternalEffects|TestDeliveryReceiptResolutionSelectsFrozenStartupEntrypoint|TestDeliveryBeginUnconfirmedAndWrongOwnerCannotAccept|TestDeliveryReceiptDoesNotRequireBlanketPermission|TestExecuteCheckInHasNoBlanketPermissionFlag|TestExecuteResolveStopNamesChangedActionWithoutPermissionClaim)$'
for name in TestDeliveryStartupCLI TestDeliveryObservedStopRequiresOwnerResolutionBeforeProgress TestDeliveryResolveStopContinuesNamedChangedActionWithoutPermissionClaim TestDeliveryReceiptResolutionPreservesUncertainExternalEffects TestDeliveryReceiptResolutionSelectsFrozenStartupEntrypoint TestDeliveryBeginUnconfirmedAndWrongOwnerCannotAccept TestDeliveryReceiptDoesNotRequireBlanketPermission TestExecuteCheckInHasNoBlanketPermissionFlag TestExecuteResolveStopNamesChangedActionWithoutPermissionClaim; do
 grep -q "^--- PASS: $name " "$evidence/tests.stdout"
done
capture vet go vet ./internal/taskrun ./cmd
capture diffcheck git diff --check
test -z "$(git status --porcelain)"
