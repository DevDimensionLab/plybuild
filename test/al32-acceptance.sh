#!/bin/sh
set -eu
cached_go=/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin
export PATH="$cached_go:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin"
export GOTOOLCHAIN=local GOFLAGS=-mod=readonly LC_ALL=C
export GOMODCACHE=/Users/perottochristensen/go/pkg/mod
export GOCACHE=/Users/perottochristensen/Library/Caches/go-build
: "${PLY_CANDIDATE_BINARY:?exact private candidate binary required}"
: "${PLY_CANDIDATE_OID:?exact candidate OID required}"
: "${PLY_CANDIDATE_TREE:?exact candidate tree required}"
test -x "$PLY_CANDIDATE_BINARY"
test "$(git rev-parse HEAD)" = "$PLY_CANDIDATE_OID"
test "$(git rev-parse 'HEAD^{tree}')" = "$PLY_CANDIDATE_TREE"
test -z "$(git status --porcelain)"
"$PLY_CANDIDATE_BINARY" version --format json
help="$TMPDIR/al32-review-return-help.txt"
"$PLY_CANDIDATE_BINARY" workflow execute review-return --help >"$help"
cat "$help"
for phrase in 'summary: 1-2000 Unicode codepoints' 'findings[].summary: 1-600 Unicode codepoints' '256 KiB (262144 bytes)'; do
 grep -F -q "$phrase" "$help"
done
# Generated brief/template and native parser, with modeled fixture state only.
# Reuse AL13 long-prose/control coverage without a real reviewer or process probe.
log="$TMPDIR/al32-tests.log"
if go test -v -count=1 -timeout=90s ./internal/taskrun -run '^(TestAL32GeneratedReviewReturnLimits|TestAL13ReviewReturnProseParser)$' >"$log" 2>&1; then
 cat "$log"
else
 result=$?
 cat "$log"
 exit "$result"
fi
for name in TestAL32GeneratedReviewReturnLimits TestAL13ReviewReturnProseParser; do
 grep -q "^--- PASS: $name " "$log"
done
go test -v -count=1 ./cmd -run '^TestExecuteReviewIsOneCoordinatorOperation$'
git diff --check
test -z "$(git status --porcelain)"
