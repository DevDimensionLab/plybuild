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
measured with `go version go1.26.7 darwin/arm64`, `GOOS=darwin`, `GOARCH=arm64`,
`CGO_ENABLED=0`, `GOENV=off`, `GOWORK=off`, `GOARM64=v8.0`, and empty
`GOFLAGS`/`GOEXPERIMENT`; `go list`
selected 86 production and 14 test files. All architecture-specific Go build
selectors and the path-and-content digests are recorded under `tool.go_build`.
The scorecard also pins the wrapper, parser, source scanner, Q0.6 contract,
vendored audit, and report-template checksums. Reproduction requires the same
instrument and toolchain configuration unless the baseline is deliberately
re-established.

P7 deliberately re-established only the exact Go toolchain identity. The old
Go 1.26.2 run reproduced the stored schema-2 scorecard byte-for-byte at SHA-256
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`.
The Go 1.26.7 run produced the current authoritative scorecard at SHA-256
`9f044510ae95d1df85b6c0868323c11453f97015c16fc46721d70c9da0b3a463`.
The raw report body remained byte-identical, every denominator and criterion
object matched, and all 228 numeric debt leaves were preserved. Exact binary,
archive, image, pinned-tool, support-source, and comparison identities are in
`toolchain-migration.json`. Go 1.27 was not selected because pinned
golangci-lint 2.12.2 was built with Go 1.26.2 and its support policy does not
claim target Go versions newer than the Go line used to compile the linter.
The module retains its characterized `go 1.18` language-compatibility floor and
adds the distinct preferred-toolchain contract `toolchain go1.26.7`. This is
the Go toolchain selection model's intended separation: builds select the exact
maintained toolchain without silently enabling newer language, vet, or lint
semantics in this bounded move.

P1B deliberately migrated the structured schema without migrating debt.
`instrument-migration.json` records both complete instrument identities. The
old schema-1 parser is pinned by source commit `8d2cc113241ad58ed21efda5c5a219f2d9ef24b6`
and SHA-256 `277c2ed18ded263e883ef6ac1a94e10e4b721d16bdbf3f09299ee5aaa59fb071`;
its scorecard SHA-256 is
`d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`.
The schema-2 parser is pinned by source commit
`4887222d38f3dd45a5f61e4231f3b19c27583440` and SHA-256
`f56dc96885c0f3ab5e18bdb3ccbe155411efc0ce7bacfeb4fe701d8a23d0c31a`;
the pre-P7 schema-2 scorecard SHA-256 was
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`.
Both runs produce identical denominators and criterion objects, including all
228 numeric leaves under those objects, and retain Q3.9 PASS. Only schema,
instrument, and structured manual-evidence identity changed.

Reproduce the current raw upstream Markdown and authoritative schema-2 JSON
from the repository root with:

```sh
quality_root="$PWD/.quality"
execution_repo="$(mktemp -d /private/tmp/ply-baseline-execution.XXXXXX)"
structured_repo="$(mktemp -d /private/tmp/ply-baseline-structured.XXXXXX)"
reproduction_home="$(mktemp -d /private/tmp/ply-baseline-home.XXXXXX)"
audit_out="$(mktemp -d /private/tmp/ply-baseline-audit.XXXXXX)"
audit_gocache="$(mktemp -d /private/tmp/ply-baseline-gocache.XXXXXX)"
audit_gotmpdir="$(mktemp -d /private/tmp/ply-baseline-gotmp.XXXXXX)"
baseline_gomodcache="$(mktemp -d /private/tmp/ply-baseline-gomodcache.XXXXXX)"
for baseline_repo in "$execution_repo" "$structured_repo"; do
  git clone --shared --no-checkout "$PWD" "$baseline_repo"
  git -C "$baseline_repo" checkout --detach 5635d50bd161a9a5aa81fc4332cc0c9d68885d08
  mkdir -p "$baseline_repo/.quality"
  cp "$quality_root/baseline/inventory" "$baseline_repo/.quality/inventory"
done
cp -R "$quality_root/baseline/reproduction-home/." "$reproduction_home"
active_profile="$reproduction_home/.co-pilot/profiles/.active_profile"
test -d "$active_profile" && test ! -L "$active_profile"
test "$(shasum -a 256 "$active_profile/.fixture" | awk '{print $1}')" = \
  cb95f24c35d3987f8aba51231aade19580ffe9242324804fcad2ccff350d1c9a
export CGO_ENABLED=0 GOENV=off GOWORK=off GOTOOLCHAIN=local \
  GOCACHE="$audit_gocache" GOTMPDIR="$audit_gotmpdir" \
  GOMODCACHE="$baseline_gomodcache"

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

if LC_ALL=C LANG=C HOME="$reproduction_home" \
  bash "$quality_root/tools/vendor/quality-audit.sh" \
  "$execution_repo" --out "$audit_out"; then
  upstream_rc=0
else
  upstream_rc=$?
fi
test "$upstream_rc" -eq 1

if LC_ALL=C LANG=C HOME="$reproduction_home" \
  python3 "$quality_root/tools/scorecard.py" \
  --report "$audit_out/scorecard.md" \
  --output "$audit_out/scorecard.json" \
  --repo "$structured_repo" \
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

The directory-shaped legacy active-profile fixture is part of the pinned
reproduction environment recorded in `instrument-migration.json`; it reproduces
the stored historical test failure without changing the measured source. The
vendored audit runs in an execution replica because the exact historical tests
can create ignored files there. The structured parser reads an independently
verified pristine replica of the same commit and inventory overlay, so those
test effects cannot be mistaken for source present before measurement.

The audit meta-test's T15 control is the reproducible old/new migration recipe.
It extracts the old instrument and its schema-1 evidence from the recorded Git
commit. Before using the pinned reproduction environment, it runs old and new
vendored audits in separate clean-HOME replicas, records every changed path and
raw byte digest, requires their path, mode, and relocation-independent content
digests to match, and proves pre-existing source drift is rejected. Only the
two disposable checkout roots are normalized for that comparison; the raw
digests remain in the retained manifests. It then runs old and new instruments
over equivalent `5635d50`
execution and structured replicas under the exact build context above,
reproduces both scorecard hashes byte-for-byte, compares normalized raw output,
verifies both instrument identities, compares every denominator and criterion
object, counts and compares all 228 numeric debt leaves, and asserts Q3.9
remains PASS:

```sh
LC_ALL=C LANG=C bash .quality/tools/test-quality-audit.sh
```
