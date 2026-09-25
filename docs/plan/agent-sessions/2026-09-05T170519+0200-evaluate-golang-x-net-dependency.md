# Agent Session: Evaluate Golang X Net Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-05T170519+0200-evaluate-golang-x-net-dependency`
Created: `2026-09-05T17:05:19+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8b0f151cc4c4ba3d8cbc149fd72b298dd0a7e48bada5278969ed1702dfa67c54`
Previous: [2026-09-05T161421+0200-evaluate-golang-x-image-dependency.md](2026-09-05T161421+0200-evaluate-golang-x-image-dependency.md)
Next: [2026-09-05T175806+0200-evaluate-golang-x-text-dependency.md](2026-09-05T175806+0200-evaluate-golang-x-text-dependency.md)
Outcome: rejected v0.25.0; retained v0.7.0 because mandatory complete module tests fail

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`golang.org/x/net v0.7.0` as one bounded dependency group. Resolve the
canonical latest stable release and highest floor-compatible candidate from
primary evidence. Implement one exact changed selection only if it preserves
the retained Go 1.18 floor, has an explained minimal closure, and passes every
quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted spf13/cast v1.5.1, rejected
afero v1.10.0, gotenv v1.6.0, and x/image v0.16.0, and retained
canonical-latest jwalterweatherman v1.1.0. All earlier recorded rejections,
no-change decisions, accepted closures, and evidence corrections remain final.
Do not revisit them or combine another module group. P8 remains queued.

The current build list selects golang.org/x/net v0.7.0. Treat its latest
release, Go declaration, release qualification, closure, loaded population,
consumer paths, source history, tests, and vulnerability effect as unknown
until independently resolved. This session may change only x/net's exact
required go.mod/go.sum metadata, its minimal MVS closure, and the roadmap/
handoff record. Do not change production Go, another dependency, language or
toolchain declarations, quality apparatus, Docker/release inputs, packaging,
publishers, or P8 code.

# Measurements At Start

Latest implementation remains dependency-only cast commit
`cf4fd4933251b2ca7844b6420fd3cf2a86f049a5`, exact parent
`17f7277fd2512ba06db76d1acb0af5f8b623243c`, clean tree
`485690c010cef3b02385d8f2471d3ef53611abe9`. The x/image documentation
handoff must have exact parent
`9e9ef26ff356f9c1c2d051992eb92a10948ec807`. Relative to accepted go-cmp
commit c314bcb, accepted metadata changes remain go-colorful v1.2.0 -> v1.4.1,
go-runewidth v0.0.14 -> v0.0.17, go-toml/v2 v2.0.7 -> v2.2.2, and cast
v1.5.0 -> v1.5.1, with their recorded minimal closures and exactly 15 added
checksum lines. Ordinary and ignored status must be empty.

Current dependency measurements remain 234 selected modules, 3,564 graph
edges, 429 native complete-test packages, a 332-line unapplied tidy projection,
and exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
20/30/20. The retained main module declares Go 1.18 and prefers toolchain Go
1.26.7.

Use exact Go 1.26.7 at
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in PATH, keep GOENV=off, GOWORK=off,
GOTOOLCHAIN=local, and inject no ambient GOFLAGS. Recovery evidence at
`/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` retains its verified two-entry
manifest SHA-256
`1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
Retain the verified golangci-lint 2.12.2, GoReleaser 2.17.1, apidiff, and
govulncheck v1.7.0 binaries and hashes recorded in the handover.

X/image rejection evidence is fully verified at
`/private/tmp/ply-p7-x-image-selection.9e9ef26.9J9x1F`, 52,533 entries and
manifest SHA-256
`84972465aef5f19f888102c7991aadd67446146a7056658b006d005666bd82bf`;
decision-summary SHA-256 is
`2b70234a63a199ab6cf11050212645162a78c4836aa9aafa4cb00dad764eb936`.
V0.45.0 is canonical latest and v0.16.0 is the highest complete
Go-1.18-floor-compatible release, but v0.16.0 complete tests and vet fail
under Go 1.26.7. X/image v0.5.0 remains selected and dependency metadata is
unchanged. Preserve all prior verified roots and the recorded go-colorful
mutable telemetry, btree regression, and cast whitespace-path corrections.

# Role And Boundaries

From fresh external archives and caches, resolve x/net releases through the Go
proxy, checksum database, authoritative upstream repository and mirror, and
primary Go vulnerability data. Record exact tag commits/times, module Go
declarations and requirements, checksum pairs, source identity, tag and commit
signature status, relevant release history, and archived/deprecated status.
Explicitly distinguish stable releases, prereleases, retractions, forks,
nested or later major paths, and unreleased upstream commits. Prove canonical
latest and the highest stable release compatible with Go 1.18 from declarations
and complete changed closure; do not infer compatibility from a modern build.

Measure old versus candidate selected modules, graph edges, complete package
population, checksums, loaded packages and paths, explicit exact-get diff, and
`go mod tidy -diff`. Require an explained minimal selection/edge/checksum
closure. Independently prove whether x/net is loaded; if it is, identify its
real consumers and run focused behavior over the actually used packages and
symbols. Also require candidate module complete tests, repository build,
complete tests/race/vet, pinned lint, byte-identical public help, identical
API/CLI reports, and exact Darwin and Windows vulnerability populations.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity, module
self-tests, loaded behavior, or any repository quality contract fails. If the
selected release is already the exact floor-compatible decision and exact get
is a no-op, record that no-change decision without manufacturing a dependency
commit. Do not assume vulnerability identity: x/net is present in current
findings, so prove exact old/candidate IDs and traces.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, P7 in
the roadmap, go.mod/go.sum, the rejected x/image and gotenv decisions, retained
jwalterweatherman decision, accepted cast decision and quality evidence,
rejected afero decision, prior bounded dependency decisions, and the
toolchain, compatibility, snapshot/Docker, quality, baseline-reproduction, and
audit contracts. Preserve every recorded manifest correction.

# Three Moves

Only if every decision gate passes and the selection changes, use exact Go
1.26.7 and exact `go get golang.org/x/net@<selected-version>` for one
dependency-only commit. Do not hand-edit module metadata and do not use tidy as
implementation. Preserve every retained dependency selection, language/
toolchain declaration, production source, quality apparatus, and release input.

After a changed selection, run the complete P7 dependency gate: focused
behavior, graph/path, tests/race/vet, pinned lint, help/API/CLI, launcher and
Make contracts, preflight, host plus fresh snapshot/Docker meta and acceptance,
audit meta, focused and exact Q0-Q2 audits, separate full audit, vulnerability
comparison, empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
make quality must exit 0 with all 27 rows at L2 and zero held, regressed,
not-comparable, or dirty counts. The full audit may exit 1 only for established
queued L3 rows, never 2.

Keep all caches, projections, reports, generated artifacts, build contexts,
schema-2 evidence, and audit output outside the worktree. Warm caches from a
separate external Git archive; never run `go mod download all` inside a
measured tree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next measured P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
that next group, launch a successor, push, merge, publish, release, stash,
revert, delete retained evidence/images, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Canonical latest stable is v0.58.0 at unsigned lightweight tag commit
`acc78e0d2b2c855c0c4fbdcfe5f42a9e3d0f9778`, published
2026-08-12T17:41:32Z and declaring Go 1.25.0. The proxy lists exactly 58
stable semantic releases and no prereleases or retractions; the active module
has no deprecation marker, later major path, nested module, or GitHub Release
objects. Master is 22 commits later, unreleased, and declares Go 1.26.0.

The highest stable candidate whose complete changed closure preserves the
retained Go 1.18 floor is v0.25.0. Although x/net itself continues to declare
Go 1.18 through v0.35.0, v0.26.0 is the first release whose x/text requirement
selects x/tools pseudo-version
`v0.21.1-0.20240508182429-e35e4ccd0d2d`, which declares Go 1.19. V0.25.0
instead selects x/text v0.15.0, x/tools v0.6.0, x/mod v0.8.0, and x/crypto
v0.23.0; every changed selection declares Go 1.18 or lower.

V0.25.0 is unsigned lightweight tag commit
`d27919b57fa8dd03198f85ca9e675e1a09babd7d` at
2024-05-06T16:24:48Z. Its checksum pair is
`h1:d/OCCoBEUq33pjydKrGQhw7IlUPI2Oylr+8qLx49kac=` /
`h1:JkAGAh7GEvH74S6FOH42FLoXpXbE/aqXSrIQjXgsiwM=`. All 778 proxy files
match the exact authoritative upstream tag; the proxy ZIP SHA-256 is
`7fd8464681c3011736f2c75beb20f88fff553a17f4f574325bce5ca5dc1fcf83`.
The authoritative repository and unarchived, undisabled, non-fork GitHub
mirror have identical master and all 58 tag refs. All tags are lightweight
commit refs without independent signatures, and relevant commits are unsigned.

The exact projection changes x/net v0.7.0 -> v0.25.0 and x/text v0.7.0 ->
v0.15.0 in go.mod and adds exactly those two checksum pairs. MVS also selects
x/crypto v0.23.0, x/mod v0.8.0, and x/tools v0.6.0; existing x/sys v0.30.0
and x/term v0.29.0 dominate the candidate's lower requirements. Modules and
complete packages remain 234 and 429; graph edges change 3,564 -> 3,568
through five removals and nine additions, and the unapplied tidy projection
changes 332 -> 335 lines.

Two x/net packages load: html and html/atom. The path is `plybuild/cmd ->
go-term-markdown -> x/net/html`; the consumer calls `html.Parse` and walks
nodes, node-type constants, and attributes. Both loaded packages and the Ply
markdown contract pass at count 10 in both projections. The upstream consumer
test reproduces the same pre-existing ANSI expected-text mismatch in both
states with normalized-identical output, so there is no x/net behavior delta.

The candidate is rejected. Under exact Go 1.26.7, complete module tests fail
at count 1, count 10, and race because `route.TestRouteMessage` cannot create
an AF_ROUTE raw socket in the measured Darwin sandbox and reports `operation
not permitted`; count 10 produces exactly ten such diagnostics. Selected
v0.7.0 fails the same route test, and the fatal socket behavior remains on
unreleased master. Candidate vet and module verification pass, and its source
remains unchanged, but the mandatory module self-test stop rule still forbids
the dependency edit and all downstream candidate repository gates.

The candidate leaves Darwin and Windows symbol populations identical at
20/20, including the same nine symbol-reachable x/net IDs and traces. Its
Darwin module population improves 30 -> 27 by removing GO-2023-1988,
GO-2023-2102, and GO-2024-2687. That module-only improvement cannot override
the failed quality contract. `go.mod` and `go.sum` remain unchanged and there
is no implementation commit.

Evidence is sealed at
`/private/tmp/ply-p7-x-net-selection.e5ea7ae.cP57X8`; its fully verified
71,887-entry manifest SHA-256 is
`88e4be87483322792ad3da0064bcef4ced5ee21c52821bb1f80c143f80a3aaf1`.
Decision-summary SHA-256 is
`6e32698abcde67604e159488ba21f6e780fd61dfba2ed9299644b5f011940dff`.
