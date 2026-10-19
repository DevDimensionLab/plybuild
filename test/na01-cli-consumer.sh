#!/bin/sh
# NA01 public CLI consumer check for an exact or ordinary installed Ply.
#
#   PLY_BIN=/absolute/ply sh test/na01-cli-consumer.sh
#
# Runs the existing native SL01 question/answer journey in a disposable
# fixture (real native worker, answer store, continuation and week controller;
# owned stand-in provider, no live provider, Slack, phone or passkey). At each
# stage, unanswered -> received -> continuing -> confirmed use -> commit ->
# review -> native verify -> integration -> completed (single-attempt
# journey), PLY_BIN's public reads must agree with the native projections:
#   workflow worker show, workflow week show,
#   workflow execute question, workflow execute answer status.
# Run from a source tree at the delivered commit. Read commands only.
set -eu

fail() { echo "na01-cli-consumer: FAIL: $*" >&2; exit 1; }

: "${PLY_BIN:?PLY_BIN must name the ply executable to check}"
case "$PLY_BIN" in /*) ;; *) fail "PLY_BIN must be absolute" ;; esac
[ -x "$PLY_BIN" ] || fail "PLY_BIN is not executable"
echo "na01-cli-consumer: PLY_BIN $PLY_BIN sha256 $(shasum -a 256 "$PLY_BIN" | cut -d' ' -f1)"

# Private physical TMPDIR outside any Ply workspace, also for Go temps.
tmp=$(mktemp -d /private/tmp/ply-na01-cli.XXXXXX) || fail "cannot create private TMPDIR"
chmod 700 "$tmp"
tmp=$(cd "$tmp" && pwd -P)
trap 'rm -rf "$tmp"' EXIT
trap 'exit 1' HUP INT TERM
TMPDIR="$tmp"
GOTMPDIR="$tmp/go-tmp"
GOCACHE=${GOCACHE:-"$tmp/go-build"}
GOFLAGS=-mod=readonly
GOPROXY=off
mkdir -p "$GOTMPDIR" "$GOCACHE"
export TMPDIR GOTMPDIR GOCACHE GOFLAGS GOPROXY PLY_BIN

log="$tmp/na01-cli-consumer.log"
if ! go test ./internal/taskrun -count=1 -timeout 600s -v -run '^(TestSL01NativeParentQuestionAnswerAndCompletedPauseResume|TestSL01NativeSingleAttemptZeroRetryAnswer)$' >"$log" 2>&1; then
	cat "$log" >&2
	fail "public CLI reads disagree with the native journey"
fi
grep -E '^(--- |ok)|na01_' "$log" || true
reads=$(grep -c 'na01 public CLI' "$log" || true)
[ "$reads" -ge 30 ] || fail "expected the journey's public CLI reads, saw $reads"
echo "na01-cli-consumer: PASS ($reads public CLI reads agreed)"
