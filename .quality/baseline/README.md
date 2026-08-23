# Baseline measurement context

`raw-upstream-scorecard.md` is the untouched Markdown emitted by the pinned
upstream audit. It may contain upstream verdicts corrected in `scorecard.json`;
the structured scorecard is authoritative.

The source checkout was commit
`5635d50bd161a9a5aa81fc4332cc0c9d68885d08`, tree
`ad3b300e55a9ff08936c09fdbeb875a35bbcf376`. The audit overlaid
`.quality/baseline/inventory` as `.quality/inventory` before measuring because
that file is the declared population for the lift and was not present in the
source commit. The overlay
SHA-256 was
`4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d`.
The same metadata and the `8 / 8 / 4 / 5` declared-population denominators are
recorded structurally in `scorecard.json`. The authoritative local measurements
are `80 / 80` direct external call sites outside the five intended adapters,
`0 / 8` subjects with real production roots and their named executable harness,
and `2 / 33` skipped tests with no central leak-guard apparatus. All five
intended adapter directories were absent, so Q1.3 is a project finding and is
not a comparable ratchet precondition until those exact directories exist.

The structured baseline is also bound to the selected Go build context. It was
measured with `go version go1.26.2 darwin/arm64`, `GOOS=darwin`, `GOARCH=arm64`,
`CGO_ENABLED=0`, `GOENV=off`, `GOWORK=off`, `GOARM64=v8.0`, and empty
`GOFLAGS`/`GOEXPERIMENT`; `go list`
selected 86 production and 14 test files. All architecture-specific Go build
selectors and the path-and-content digests are recorded under `tool.go_build`.
The scorecard also pins the wrapper, parser, source scanner, Q0.6 contract,
vendored audit, and report-template checksums. Reproduction requires the same
instrument and toolchain configuration unless the baseline is deliberately
re-established.

Reproduce both the raw upstream Markdown and authoritative structured JSON from
the repository root with:

```sh
quality_root="$PWD/.quality"
baseline_repo="$(mktemp -d /tmp/ply-baseline.XXXXXX)"
audit_out="$(mktemp -d /tmp/ply-baseline-audit.XXXXXX)"
git clone --shared --no-checkout "$PWD" "$baseline_repo"
git -C "$baseline_repo" checkout --detach 5635d50bd161a9a5aa81fc4332cc0c9d68885d08
mkdir -p "$baseline_repo/.quality"
cp "$quality_root/baseline/inventory" "$baseline_repo/.quality/inventory"
export CGO_ENABLED=0 GOENV=off GOWORK=off

test "$(go version)" = "$(python3 -c \
  'import json,sys; print(json.load(open(sys.argv[1]))["tool"]["go_build"]["version"])' \
  "$quality_root/baseline/scorecard.json")"
for key in GOOS GOARCH CGO_ENABLED GOFLAGS \
  GO386 GOAMD64 GOARM GOARM64 GOMIPS GOMIPS64 GOPPC64 GORISCV64 GOWASM \
  GOEXPERIMENT; do
  expected="$(python3 -c \
    'import json,sys; print(json.load(open(sys.argv[1]))["tool"]["go_build"][sys.argv[2]])' \
    "$quality_root/baseline/scorecard.json" \
    "$(printf '%s' "$key" | tr '[:upper:]' '[:lower:]')")"
  test "$(go env "$key")" = "$expected"
done

if LC_ALL=C LANG=C bash "$quality_root/tools/vendor/quality-audit.sh" \
  "$baseline_repo" --out "$audit_out"; then
  upstream_rc=0
else
  upstream_rc=$?
fi
test "$upstream_rc" -eq 1

if LC_ALL=C LANG=C python3 "$quality_root/tools/scorecard.py" \
  --report "$audit_out/scorecard.md" \
  --output "$audit_out/scorecard.json" \
  --repo "$baseline_repo" \
  --commit 5635d50bd161a9a5aa81fc4332cc0c9d68885d08 \
  --upstream "$quality_root/tools/vendor/quality-audit.sh" \
  --upstream-exit 1 \
  --manual-evidence "$quality_root/baseline/manual-evidence.json" \
  --inventory-overlay-sha256 \
    4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d; then
  structured_rc=0
else
  structured_rc=$?
fi
test "$structured_rc" -eq 1

tail -n +2 "$audit_out/scorecard.md" > "$audit_out/reproduced-raw.body"
tail -n +2 "$quality_root/baseline/raw-upstream-scorecard.md" > "$audit_out/stored-raw.body"
cmp "$audit_out/reproduced-raw.body" "$audit_out/stored-raw.body"
cmp "$audit_out/scorecard.json" "$quality_root/baseline/scorecard.json"
```

Both audit commands are expected to exit `1`: that means measured project
findings, not a broken audit. Exit `2`, a checksum mismatch, a dirty checkout
beyond the declared inventory overlay, or either failed `cmp` invalidates the
reproduction. Raw-report line 1 is excluded from the comparison because upstream
embeds the checkout's absolute path there; every measured line remains covered.
