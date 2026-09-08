# Agent Session: Evaluate Bgentry Speakeasy Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-08T023421+0200-evaluate-bgentry-speakeasy-dependency`
Created: `2026-09-08T02:34:21+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ad9291ae115c91e01fa403345b73a551c4d3e1f7321c91ed5f31a13e2beaa808`
Previous: [2026-09-08T005841+0200-evaluate-beorn7-perks-dependency.md](2026-09-08T005841+0200-evaluate-beorn7-perks-dependency.md)
Next: [2026-09-08T044755+0200-evaluate-bketelsen-crypt-dependency.md](2026-09-08T044755+0200-evaluate-bketelsen-crypt-dependency.md)
Outcome: Upgraded exact-path Speakeasy v0.1.0 to canonical stable latest
  v0.2.0 in dependency-only commit `41f9561`; its one-module closure preserves
  Go 1.18 and the complete quality and vulnerability-parity gates pass.

# Answer

Select `github.com/bgentry/speakeasy v0.2.0`. Exact Go 1.26.7
`go get github.com/bgentry/speakeasy@v0.2.0` produced dependency-only commit
`41f9561f6ea2f5b6395c5c4d9bcc56a54533133a`, exact parent
`1f55aaa31280a4610ed66c1c37ddda953bfa6a8e`, tree
`a50d4f256fe742e472f1f7e9cf0589a596e971b0`. Only `go.mod` and `go.sum`
changed, with three insertions and no deletion or production, toolchain,
quality, release, or packaging change.

## Canonical Identity And Release Qualification

The exact proxy lists only stable v0.1.0 and v0.2.0; no prerelease or
pseudo-version is listed. Exact `@latest`, `@v0`, `@master`, the default branch,
and tag v0.2.0 all resolve commit
`760eaf8b681647364e7a400b856e0921248728a5` at
2022-09-10T01:20:23Z. Go-import metadata maps the exact module path to
`https://github.com/bgentry/speakeasy.git`. GitHub reports that public source
enabled, unarchived, undisabled, non-fork, and on default branch `master`, with
one branch, two tags, and one Release. There is no retraction, module
deprecation, `/v2` module, or later default-branch commit.

Selected v0.1.0 is unsigned commit
`4aabc24848ce5fd31929f7d1e4ea74d3709c14cd`, tree
`4905bd85608c760d3ebf869da51b5740731c50b1`, at 2017-04-17T20:07:03Z. Signed annotated tag object
`12abe455afa69a18b1098acb18b1f004a4a6c33c` was created later, at
2017-06-15T22:05:56Z. Its tag signature verifies locally with fingerprint
`757FD463E177A2F1CD1C89038B6EDBF713E83E69`; the key is now expired, while
GitHub reports the signature valid. The target commit itself is unsigned.

Candidate v0.2.0 is a lightweight tag, so it has no independent tag object,
tagger time, or tag signature. Its commit has tree
`b2198079290577be7ff0fb5732f778664e28fd64`, parent
`a4057f540bab4628fbdc15a38dfbccdd34e6fa5d`, and a
cryptographically valid GitHub web-flow signature with fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; that key is now expired. The
sole non-draft, non-prerelease GitHub Release is v0.2.0. It was created at the
2022 commit time but published much later, at 2024-06-27T20:45:36Z. That
delayed GitHub Release publication is distinct from the lightweight tag and
commit time.

Both proxy ZIPs contain eight regular files and match their exact Git tag
archives byte-for-byte. V0.1.0 and v0.2.0 ZIP SHA-256 values are respectively
`d4bfd48b9bf68c87f92c94478ac910bcdab272e15eb909d58f1fb939233f75f0`
and `54db02088d979118dabe83139b80efa5d84d8bba03f3effe458be5c61f20f09f`.
Their normalized manifests are
`9ac15e9129b5abc956a5c93d9bcd86a6c13ea948ec63e09a221fa9e7eb7f4ebe`
and `ad860d846b60438708efba826f5b05649fa4aa956882e2a6402377136ad842d7`.
Sumdb verifies the v0.1.0 source/module pair
`h1:ByYyxL9InA1OWqxJqqp2A5pYHUrCiAL6K3J+LKSsQkY=` /
`h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=` and the v0.2.0 pair
`h1:tgObeVOf8WAvtuAX6DhJ4xks4CFNwPDZiqzGqIHE51E=` / the same module hash.

## Floor, Tests, And Consumer Behavior

Neither module declares a Go version or a requirement. The complete minimal
declared closure is nevertheless proved, not inferred: one Speakeasy module,
two packages, and only standard-library imports. Proxy and exact-Git forms of
both versions pass verify/list, count-1, two repeated count-10 runs, race, and
vet under exact Go 1.26.7 and exact Go 1.18.10, with source bytes unchanged.
The module has no `_test.go` files, so those runs prove complete package
buildability and absence of upstream failures rather than behavioral coverage.

Selected `github.com/mitchellh/cli v1.1.0` is the real graph consumer. Its
`BasicUi.AskSecret` calls the sole used symbol, `speakeasy.Ask("")`, only for
a terminal stdin. Focused consumer tests pass count-1/count-10/race and vet
against both versions under Go 1.18.10. A scratch-only real-PTY fixture waits
until the ECHO bit is disabled before sending a secret, then proves the secret
is absent from captured terminal output, the returned value is exact, and ECHO
is restored for both versions. An initial immediate-write probe exposed a
harness scheduling race and was superseded by that terminal-state-aware run;
it was not a dependency failure and touched no repository file.

## Exact Selection Closure

The exact candidate move changes only Speakeasy v0.1.0 to v0.2.0. Selected
modules remain 234, complete packages 429, loaded modules 41, and loaded
packages 197. Graph edges grow 3,581 to 3,582 solely for the new main-module
edge to v0.2.0; Mitchellh CLI v1.0.0 and v1.1.0 retain their historical
v0.1.0 requests. `go.sum` grows 1,045 to 1,047 lines only for the candidate
source and module hashes. No Speakeasy package is loaded by Ply, Ply contains
no source import, and `go mod why -m` says the main module does not need it.

The unapplied tidy projection grows 361 to 371 lines. Tidy removes the explicit
candidate pin and its pair, restores inherited v0.1.0 selection, and leaves all
historical tidy debt. Tidy was evidence only, never implementation. Relative
to accepted go-cmp commit `c314bcb`, accepted dependency metadata now adds
exactly 31 checksum lines.

## Repository, Quality, And Vulnerability Results

Repository module verification, build, complete count-1/count-10/race tests,
vet, Windows-amd64 build, and pinned golangci-lint 2.12.2 pass. Parent and
candidate root/status/upgrade/build help plus API and CLI reports are
byte-identical. API and CLI SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
Complete corrected preflight passes all 62 launcher controls, every Make and
production-script meta-contract, and all 15 audit controls. One nested launcher
signal-retention timing probe flaked once; an immediate standalone repeat and
the complete preflight repeat pass all 62 controls.

Exact `make quality` exits 0 with all 21 stages and all 27 Q0-Q2 rows passing
at L2, 80/80 mutations killed, valid six-receipt schema-2 evidence, and zero
held, regressed, not-comparable, or dirty counts. Its scorecard SHA-256 is
`48decac359a9ebab23e59c29680e63682ab6d1a65141d13911644143a8db2a01`.
Host, fresh snapshot, and fresh Docker acceptance pass; report SHA-256 values
are `ceed704ab74e0241549502b952bb3e52dab154a3b8501cfbc2460f28a58e7793`
and `11a03f099dcf20ff195e5f04044b1e780aca7532b6318fe90d5071e1cc5005b2`.
The separate full audit exits expected 1, never 2, only for queued Q3.1, Q3.3,
Q3.4, and Q3.7; scorecard SHA-256 is
`53f8d233daa28ad7c5777bc79c8fa46434ce9989fc0c15719ddd3614a2453c15`.
Empty-HOME count-2 passes.

Fresh primary data updated 2026-09-02T19:12:04Z contains 1,392 module records
and no exact Speakeasy record. Old/candidate Darwin reachable results are
identical at 20 IDs/22 traces, Darwin module results at 30 IDs, and Windows
reachable results at 20 IDs/22 traces. No trace names Speakeasy. The exact
reachable normalized result SHA-256 is
`7757df547ed97c0b709bfe8cbbbf807356331727c870bdeb7f0602b1f393ba79`.

The identical old/candidate reachable ID set is `GO-2023-1989`,
`GO-2023-1990`, `GO-2024-2937`, `GO-2024-3333`, `GO-2025-3595`,
`GO-2026-4440`, `GO-2026-4441`, `GO-2026-4815`, `GO-2026-4961`,
`GO-2026-5025`, `GO-2026-5027`, `GO-2026-5028`, `GO-2026-5029`,
`GO-2026-5030`, `GO-2026-5031`, `GO-2026-5032`, `GO-2026-5061`,
`GO-2026-5062`, `GO-2026-5066`, and `GO-2026-6222`. The 22 exact trace
groups are seven 15-frame `golang.org/x/image/tiff.Decode -> cmd.ExecuteE`
traces (1989, 1990, 2937, 4815, 5032, 5062, 5066), nine 13-frame
`golang.org/x/net/html.Parse -> cmd.ExecuteE` traces (3333, 3595, 4440,
4441, 5025, 5027, 5028, 5029, 5030), two
15-frame `golang.org/x/image/webp.Decode -> cmd.ExecuteE` traces (4961,
5061), and one each for 15-frame
`golang.org/x/image/bmp.Decode -> cmd.ExecuteE` (5031), 17-frame
`golang.org/x/image/vp8l.Decode -> cmd.ExecuteE` (6222), five-frame
`golang.org/x/image/tiff.(*buffer).ReadAt -> pkg/http.GetJsonWithAccessToken`
(4815), and four-frame `golang.org/x/image/webp.init -> cmd.init` (5061).
Darwin and
Windows reachable evidence have the same IDs and trace groups; only
platform-specific positions differ. The 30-ID Darwin module population is
also byte-identical old versus candidate; its normalized ID-list SHA-256 is
`be304d82c9db9981b650fd0effa3a214aed50b2ca1ba79b1ad3835212d6bcc3e`.

The final evidence manifest contains 459 verified entries and has SHA-256
`17bb756dfd0e81c39da3616f6f81edbcc59e299b295d422c1a701be603c02cc4`.
Decision-summary SHA-256 is
`d3d5678a31f494e194321951086ccdb7579c70f39a0b7cc311bc4f9c925666f0`.
All compatibility-cache warming, BSD-`mktemp`, report-routing, and PTY timing
corrections stayed beneath the authorized scratch root and were superseded
before claims; none changed repository source or concealed a candidate failure.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/bgentry/speakeasy v0.1.0` as one bounded dependency group. Resolve
canonical latest, authoritative source identity, release qualification,
default-branch history, and the highest qualified Go-1.18-floor-compatible
candidate from primary evidence. Make an exact dependency selection only if it
changes a selected version, preserves the retained floor through the complete
minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Perks v1.0.1. All earlier rejections,
no-change decisions, accepted closures, and evidence corrections remain final.
Do not revisit Perks or combine another module group. P8 remains queued.

Project MVS selects Speakeasy v0.1.0 at 2017-04-17T20:07:03Z although it is not
an explicit `go.mod` requirement. A minimal post-Perks survey finds exact
stable proxy versions v0.1.0 and v0.2.0. Exact Go `@latest`, `@v0`, and
`@master` resolve v0.2.0 at 2022-09-10T01:20:23Z, commit
`760eaf8b681647364e7a400b856e0921248728a5`. Both proxy module files declare
only exact path `github.com/bgentry/speakeasy`, with no Go directive or
requirements. Do not infer the complete closure floor from missing directives.

The public exact-path GitHub repository currently reports enabled, unarchived,
undisabled, non-fork status, default branch `master`, one branch, two tags, and
one GitHub Release. Master and v0.2.0 point to
`760eaf8b681647364e7a400b856e0921248728a5`; v0.1.0 points to
`4aabc24848ce5fd31929f7d1e4ea74d3709c14cd`. The sole non-draft,
non-prerelease GitHub Release is v0.2.0, published
2024-06-27T20:45:36Z for the older tag. Treat canonical source identity, tag
and commit signatures, release/tag history, retractions, all version and
branch declarations, complete closures, tests, consumers, loaded behavior,
and vulnerability effect as unknown until independently proved. Explicitly
distinguish tag/commit time from GitHub Release publication time.

# Measurements At Start

Latest dependency implementation remains Circbuf commit
`3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
`7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
`a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
`go.sum`. The Consul API, Go Metrics, Go Radix, and Perks evaluations made no
dependency commit. Operator-authorized lifecycle repair
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between the
earlier Repr implementation and direct-child Repr handoff. Preserve it.

Perks remains selected at v1.0.1, the exact-path signed stable latest and
default-branch tip. The enabled, unarchived exact-path repository is a fork,
but go-import and proxy identity prove it is the canonical published
`github.com/beorn7/perks` source; parent `github.com/bmizerany/perks` is a
distinct path with divergent untagged history. Exact selected-version get
changed no selection and only projected redundant metadata, which was not
applied. Perks evidence has 477 verified entries; manifest SHA-256 is
`aa77d2ab9cf676ecfb7ef544c1c5db76ecda86f580ff8ff8f03b8c78038677f0`;
decision-summary SHA-256 is
`6db4ca7260bde6b4affa48c62adb381bda20687f39f013c4b1cdc25397418754`.

Accepted measurements remain 234 selected modules, 3,581 graph edges, 429
native complete-test packages, 41 loaded modules, 1,045 `go.sum` lines, and a
361-line unapplied tidy projection. Relative to accepted go-cmp commit
`c314bcb`, accepted metadata adds exactly 29 checksum lines. The main module
retains Go 1.18 and toolchain Go 1.26.7. Ordinary and ignored status must be
empty.

Fresh primary vulnerability evidence contains 1,392 module records and no
exact or parent Perks record. Accepted projections retain exact 20-ID/22-trace
Darwin symbol, 30-ID Darwin module, and 20-ID/22-trace Windows symbol
populations.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate pinned tools beneath scratch as needed. Portable
receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
Treat rebuilt tool hashes as nonportable receipts; prove versions and
functionality.

# Role And Boundaries

From fresh external archives and caches, resolve every relevant exact-path
Speakeasy v0 version through the Go proxy and checksum database, authoritative
repository, go-import metadata, and primary Go vulnerability data. Record
selected and candidate commits/times, module Go declarations and requirements,
checksum pairs, source identity, tag and commit signatures, release history,
archived/deprecated status, and retractions. Explicitly distinguish stable
versions, prereleases, pseudo-versions, redirects, forks, alternate module
paths, branch heads, unreleased commits, and a delayed GitHub Release object.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete minimal closure, not from
missing directives. Measure old versus candidate modules, graph edges,
complete packages, checksums, loaded packages and paths, exact-get diff, and
`go mod tidy -diff`. Explain every selection, edge, and checksum change.
Identify real consumers and exercise the actually used symbols and terminal
behavior without requiring an interactive user or leaking terminal state.

Require candidate module complete, repeated, and race-enabled self-tests plus
vet; repository verify/build, complete tests, race, vet, Windows build, pinned
lint, byte-identical public help, API/CLI reports, and primary vulnerability
identity. Preserve the retained Go 1.18 floor across the complete changed
closure. Keep standalone test-only closure separate from project MVS.

Reject or retain without dependency edits if canonical qualification, source
identity, floor compatibility, closure, module tests, loaded behavior, or any
applicable quality contract fails. Do not manufacture metadata when exact get
changes no selected version. Record precise primary old/candidate
vulnerability IDs and traces.

# Required Reading

Verify branch, clean ordinary and ignored status, exact ancestry, reciprocal
archive history, `./codex-dev-start.sh --check`, and P7/P8 state before work.
Read this archive, the rolling handover, P7 roadmap, `go.mod`, `go.sum`, the
answered Perks, Go Radix, Go Metrics, Consul API, Circbuf, Optional, Template,
Colour, Chroma, Kong, Repr, Assert, and Units archives, retained
Kingpin/Resty/Errgo/Check/YAML outcomes, and the toolchain, quality, baseline,
compatibility, snapshot/Docker, acceptance, audit, and lifecycle contracts. Do
not reopen earlier decisions.

# Three Moves

If and only if v0.2.0 is qualified and its complete minimal closure preserves
Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/bgentry/speakeasy@v0.2.0` for one dependency-only commit.
Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, toolchain declaration, production source, quality apparatus,
and release input.

If v0.2.0 fails qualification, or exact get changes no selection, retain
v0.1.0 without adding a redundant direct or indirect requirement, main edge,
or checksum solely for metadata. A no-change decision gets no dependency
implementation commit.

After a changed selection, run the complete P7 dependency gate: module and
consumer tests, graph/path/checksum/tidy explanation, repository verify/build/
tests/race/vet/Windows/pinned lint, help/API/CLI, launcher and Make contracts,
preflight, host plus fresh snapshot/Docker meta and acceptance, audit meta,
focused and exact Q0-Q2 audits, separate full audit, vulnerability comparison,
empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 Q0-Q2 rows at L2 and zero held,
regressed, not-comparable, or dirty counts. Full audit may exit 1 only for
established queued L3 rows, never 2.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
create direct `/private/tmp/ply-*` roots, never bypass scratch cleanup, and
never run `go mod download all` inside a measured tree. Warm only exact needed
closures from a separate scratch archive. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Speakeasy decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
