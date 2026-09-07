# Agent Session: Evaluate Antihax Optional Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-07T183413+0200-evaluate-antihax-optional-dependency`
Created: `2026-09-07T18:34:13+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `4ccca27a6a4f5adabfd2d7cac70abbeb8c7e9926c94af46cd9cf670ae64dda7e`
Previous: [2026-09-07T174802+0200-evaluate-alecthomas-template-dependency.md](2026-09-07T174802+0200-evaluate-alecthomas-template-dependency.md)
Next: [2026-09-07T194916+0200-evaluate-armon-circbuf-dependency.md](2026-09-07T194916+0200-evaluate-armon-circbuf-dependency.md)
Outcome: Retained selected stable `github.com/antihax/optional v1.0.0` without dependency metadata edits: it is canonical latest and the highest stable exact-path Go-1.18-floor-compatible version; the later master pseudo-version is unreleased and changes no Go source, and exact selected-version get changes no selection.

# Answer

Retain exact-path `github.com/antihax/optional v1.0.0` without a dependency
edit. The exact stable proxy list contains only v1.0.0, and exact Go
`@latest`, `@v1`, and `@v1.0.0` select it at 2019-10-10T23:37:20Z. There are
no prereleases or retractions. V1.0.0 is therefore canonical latest, the sole
stable exact-path release, and the highest stable candidate compatible with
the retained Go 1.18 floor.

The authoritative repository is public, enabled, unarchived, and non-fork,
with one branch, one tag, and one GitHub Release. The v1.0.0 release is named
`Initial release`, is neither draft nor prerelease, and was published at
2019-10-10T23:39:21Z. The lightweight tag points to commit
`c3f0ba9c1a592b971d66b2787679af55b5c58f21`, tree
`b9328a8aa4526004bb36928dbc136c3acb8eec3a`. A lightweight tag has no tag
object or tag signature. The associated GitHub web-flow commit signature
independently verifies with fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`.

Exact `@master` instead resolves the unreleased pseudo-version
`v1.0.1-0.20220101210036-407d38fabb55` at 2022-01-01T21:00:36Z, commit
`407d38fabb5592e58b0841e8adb93edd65ee1319`, tree
`3dffdae3ba5f11e0b140c2edd085e5d1a69df40d`. Its commit signature also
independently verifies, but that does not qualify it as a stable release. The
two commits after v1.0.0 add and update only `README.md`; all 21 Go files and
`go.mod` are byte-identical. `/v2`, `/v3`, `gopkg.in`, and a plausible renamed
repository path resolve no authoritative alternate module. The repository has
14 downstream forks, but the canonical repository itself is not a fork.

The v1.0.0 proxy ZIP SHA-256 is
`15ab4d41bdbb72ee0ac63db616cdefc7671c79e13d0f73b58355a6a88219c97f`;
its checksum pair is
`h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=` /
`h1:uupD/76wgC+ih3iEmQUL+0Ugr19nfwCT1kdvxnR2qWY=`. The pseudo-version ZIP
SHA-256 is
`7a4fe848cb95f2bd6305266a9a35b515933cf4144436d865a2569b7f1a56d92d`;
its source checksum is
`h1:AAnkpz6e/6bV9pphv0Hj/13169F89gXw0O4kOduU1Gc=` and it shares the same
`go.mod` checksum. Independent checksum-database lookups agree. All 23 stable
proxy regular files match the exact tag, and all 24 pseudo-version files match
the exact master archive. Their normalized manifest SHA-256 values are
`66452eb1ab11bc7dc7bfe9237c47c430810d77633add3a4ef7304e1c2db48c38`
and `916591eaf56bef23e0e97a616c3ecf5b73f7f01b1c502d319a808cf7acfb3b27`.

Both versions declare Go 1.13 and no requirements. V1.0.0's complete declared
closure is therefore the module alone and preserves Go 1.18. There are no
test-only requirements. Writable proxy and exact-tag sources independently
verify and list, remain unchanged after testing, and pass complete count-1,
count-10, race-enabled tests, and vet. The sole package has no test files, so
the receipts prove buildability and absence of failing upstream tests rather
than upstream behavioral coverage.

An exact external `go get github.com/antihax/optional@v1.0.0` exits zero but
changes no selected module. Old/projected measurements are 234/234 modules,
3,580/3,581 graph edges, 429/429 complete packages, 1,043/1,044 `go.sum`
lines, and 356/358 unapplied tidy-diff lines. Module selections and package
sets are identical. The projection only adds a redundant indirect requirement,
one main-to-Optional graph edge, and the full v1.0.0 checksum; tidy removes the
requirement. Manufacturing that metadata would not upgrade a dependency, so
none was applied.

The original graph has one Optional request:
`github.com/grpc-ecosystem/grpc-gateway@v1.16.0 ->
github.com/antihax/optional@v1.0.0`. Optional loads in zero project packages,
Ply has no repository import, and `go mod why -m` says the main module does not
need it. Grpc-gateway has four generated example clients that import Optional
and use its Bool, Float64, Int32, Int64, Interface, String, and Time wrappers
with `IsSet` and `Value`; those example packages are not loaded by Ply. A
deliberately external Go-1.18 fixture exercises those seven wrappers plus
`Default` and passes count-10, race, and vet. It is consumer evidence, not a
claim that Ply executes Optional.

Untouched and projected repository copies pass module verification, build,
complete count-1/count-10/race tests, vet, Windows build, and pinned
golangci-lint 2.12.2. API and CLI reports are byte-identical at SHA-256
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
Root, status, and upgrade help are byte-identical; compatibility against
v1.0.1 and the CLI surface contract pass. A clean complete preflight passes
all 62 launcher controls, distribution/lint/install/toolchain contracts,
snapshot/Docker meta-contracts, mutation and verification meta-tests, and all
15 quality-audit controls. Changed-selection host/snapshot/Docker/exact-quality
runs are inapplicable because no selected version changed.

Fresh primary vulnerability data was updated 2026-09-02T19:12:04Z and contains
1,392 module records with no Optional record. Sorted old/projected IDs and
traces are identical: 20 Darwin reachable-symbol IDs with 22 trace events, 30
Darwin module IDs, and 20 Windows reachable-symbol IDs with 22 trace events.
No finding or trace contains Optional. The reachable-symbol ID SHA-256 is
`58889cc567e8def8340f36f3855fc510c3565677e81a81c517ff6cc115e4309b`;
the reachable-trace SHA-256 is
`f97e3357875d2ad04e8ae069770281375955495153cdfbd8763ab1f31b6e2190`;
the Darwin module-ID SHA-256 is
`be304d82c9db9981b650fd0effa3a214aed50b2ca1ba79b1ad3835212d6bcc3e`.

Exact Go 1.26.7 was recreated beneath scratch. Its binary SHA-256 is
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
the official archive SHA-256 is
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
It remained first in PATH with GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and no
ambient GOFLAGS. The portable tool receipts remain intact; rebuilt apidiff and
govulncheck hashes are nonportable while exact version metadata and functional
gates pass. A separately warmed v1.0.1 cache, scratch-local GnuPG keyring, and
BSD-mktemp adapter corrected runner setup only. Source-manifest, checksum-delta,
and vulnerability-order corrections are explicitly recorded; none changed the
repository.

Decision evidence has 540 verified entries. Evidence-manifest SHA-256 is
`345a29c62bfadf332f14aa234b1265a9dc9142739d212534d83ace531fb39afa`;
decision-summary SHA-256 is
`945b7ca8fe9d9faa605a6a80d8d416655925f95cd727071babf813bcd4caa27b`.
No dependency implementation commit was created. Accepted metadata still adds
exactly 27 checksum lines relative to accepted go-cmp commit c314bcb, and all
retained selections and the lifecycle repair remain unchanged.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/antihax/optional v1.0.0` as one bounded dependency group. Resolve
canonical latest, stable-release qualification, the later default-branch
history, and the highest Go-1.18-floor-compatible candidate from primary
evidence. Make an exact dependency selection only if it changes a selected
version, preserves the retained floor through the complete minimal closure,
and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Alecthomas Template. All earlier
rejections, no-change decisions, accepted closures, and evidence corrections
remain final. Do not revisit Template or combine another module group. P8
remains queued.

Project MVS selects Optional v1.0.0 although it is not an explicit `go.mod`
requirement. A fresh post-Template survey finds the exact stable proxy list
contains only v1.0.0. Proxy and exact Go `@latest` resolve v1.0.0 at
2019-10-10T23:37:20Z with a Go 1.13 declaration, while exact Go `@master`
resolves the later unreleased pseudo-version
`v1.0.1-0.20220101210036-407d38fabb55` at 2022-01-01T21:00:36Z, also
reporting Go 1.13. Do not treat the unreleased branch head as an in-place
stable update. Treat canonical qualification, source identity, signatures,
release and alternate-path history, closure, self-tests, loaded consumers,
actual symbols, and vulnerability effect as unknown until independently
resolved.

# Measurements At Start

Latest dependency implementation remains Repr commit
`6f4d02eb9c86ec2df8a488a85ed973aababe1f38`, exact parent
`53475076e1c79f6d2181877e3238a1a5d389c246`, tree
`82e7b1f5c503659082207481b8339e8113e38c10`, changing only `go.mod` and
`go.sum`. Operator-authorized lifecycle repair
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
that implementation and direct-child Repr handoff
`342c7ece82eae3cd726aa658462089cfb876fbe0`. Preserve it.

Template's selected 2019 pseudo-version was retained without dependency edits.
It is canonical latest and the highest exact-path Go-1.18-floor-compatible
version, but it is not a stable release. Exact get changes no selected version
and only projects a redundant indirect requirement, one graph edge, and the
full checksum. Template has no loaded project package or repository import.
Both exact source forms fail mandatory complete, repeated, and race-enabled
self-tests in `TestJSEscaping` under exact Go 1.26.7; vet passes. Template's
1,702-entry evidence-manifest SHA-256 is
`2bfa08ca73085154ab5f0833814872efe481088703e6d7480edf1e2a40853068`;
decision-summary SHA-256 is
`4acc26e68fa5ad432a44d64ade0d399464f9b07eb2c6244b29711e048ae24170`.

Current measurements remain 234 selected modules, 3,580 graph edges, 429
native complete-test packages, 1,043 go.sum lines, a 356-line unapplied tidy
projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations. The main module retains Go 1.18 and toolchain Go
1.26.7. Relative to accepted go-cmp commit c314bcb, accepted metadata adds
exactly 27 checksum lines. Ordinary and ignored status must be empty.

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
Treat rebuilt golangci-lint and govulncheck hashes as nonportable receipts;
prove versions and functionality.

# Role And Boundaries

From fresh external archives and caches, resolve Optional versions through the
Go proxy, checksum database, authoritative upstream repository, and primary Go
vulnerability data. Record exact tag and branch-head commits/times, module Go
declarations and requirements, checksum pairs, source identity, tag and commit
signature status, release history, and archived/deprecated state. Explicitly
distinguish the stable v1.0.0 tag from prereleases, pseudo-versions,
retractions, forks, alternate module paths, branch heads, and unreleased
commits. Do not treat the newer master commit as a stable release without
authoritative release evidence.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete closure. Measure old versus
candidate modules, graph edges, complete packages, checksums, loaded packages
and paths, explicit exact-get diff, and `go mod tidy -diff`. Explain every
selection, edge, and checksum change. Identify real consumers and exercise the
actually used Optional packages and symbols where applicable.

Require candidate module complete tests, repeated tests, race, and vet;
repository build, complete tests/race/vet, pinned lint, byte-identical public
help, identical API/CLI reports, and exact Darwin and Windows vulnerability
populations. Keep module test-only requirements separate from project MVS.

Stop and record rejection without dependency edits if canonical resolution,
floor compatibility, exact closure, source identity, module self-tests, loaded
behavior, or any repository quality contract fails. If v1.0.0 remains the
highest qualified stable exact-path decision and exact get changes no selected
version, record no change without manufacturing metadata or a dependency
commit. Prove exact old/candidate vulnerability IDs and traces.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and P7/P8 checkpoint before editing. Read
this archive, rolling handover, P7 roadmap, go.mod/go.sum, answered Template,
Colour, Chroma, Kong, Repr, Assert, and Units decisions, retained Kingpin/
Resty/Errgo/YAML decisions, earlier accepted and rejected bounded dependencies,
and toolchain, compatibility, snapshot/Docker, quality, baseline-reproduction,
and audit contracts. Preserve every recorded manifest and setup correction.

# Three Moves

Only if every decision gate passes and a version selection changes, use exact
Go 1.26.7 and exact `go get github.com/antihax/optional@<selected-version>`
for one dependency-only commit. Do not hand-edit metadata or use tidy as
implementation. Preserve every retained selection, especially Repr v0.5.4,
Assert v1.0.0, Units' 2024 pseudo-version, Colour v0.1.0, Chroma v0.10.0,
Template's selected pseudo-version, the selected Kong pseudo-version, Kingpin
v2.2.6, Resty v1.12.0, Errgo v2.1.0, Check's 2019 pseudo-version, all three
retained YAML paths, language/toolchain declarations, production source,
quality apparatus, and release input.

After a changed selection, run the complete P7 dependency gate: focused
behavior, graph/path, tests/race/vet, pinned lint, help/API/CLI, launcher and
Make contracts, preflight, host plus fresh snapshot/Docker meta and acceptance,
audit meta, focused and exact Q0-Q2 audits, separate full audit, vulnerability
comparison, empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 rows at L2 and zero held, regressed,
not-comparable, or dirty counts. Full audit may exit 1 only for established
queued L3 rows, never 2.

Put every disposable cache, projection, report, generated artifact, build
context, evidence tree, and audit output beneath
`$CODEX_SESSION_SCRATCH_ROOT`; never create direct `/private/tmp/ply-*` roots.
The launcher deletes scratch after every turn. Warm caches only from a separate
archive beneath that root and never run `go mod download all` inside a measured
tree. Retain only compact decisions, digests, and receipts in tracked docs.
Preserve the lifecycle repair and bounded-scratch policy. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next measured P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
the next group, launch a successor, push, merge, publish, release, stash,
revert, bypass scratch cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
