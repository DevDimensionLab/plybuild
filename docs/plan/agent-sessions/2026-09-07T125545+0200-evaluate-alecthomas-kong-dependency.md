# Agent Session: Evaluate Alecthomas Kong Dependency

Status: NEXT
Session ID: `2026-09-07T125545+0200-evaluate-alecthomas-kong-dependency`
Created: `2026-09-07T12:55:45+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `b40fae22339cbf01df6b8dd7be98e15ec82323d36da67ff5b34976035b3c5470`
Previous: [2026-09-06T162839+0200-evaluate-alecthomas-repr-dependency.md](2026-09-06T162839+0200-evaluate-alecthomas-repr-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected `github.com/alecthomas/kong`
`v0.2.1-0.20190708041108-0548c6b1afae` as one bounded dependency group.
Resolve canonical latest, release qualification, and the highest floor-
compatible candidate from primary evidence. Implement one exact changed
selection only if it preserves the retained Go 1.18 floor, has an explained
minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and completed
dependency groups through accepted canonical-latest stable
`github.com/alecthomas/repr v0.5.4`. All earlier rejections, no-change
decisions, accepted closures, and evidence corrections remain final. Do not
revisit them or combine another module group. P8 remains queued.

The current MVS build list selects Kong
`v0.2.1-0.20190708041108-0548c6b1afae`. A fresh post-Repr survey reports
update v1.16.1 at 2026-08-09T07:05:31Z, but treat canonical latest, tag and
release qualification, Go declaration, closure, loaded population, consumer
paths, source history, tests, and vulnerability effect as unknown until
independently resolved. Kong is not a direct current `go.mod` requirement;
go.sum has the old go.mod checksum, repository Go source has no Kong import,
and `go mod why -m` says the main module does not need it. Treat those only as
starting observations. This session may change only Kong's exact selected
go.mod/go.sum metadata, its minimal MVS closure, and the roadmap/handoff
record. Do not change production Go, another dependency, language or toolchain
declarations, quality apparatus, Docker/release inputs, packaging, publishers,
or P8 code.

# Measurements At Start

Latest dependency implementation is Repr commit
`6f4d02eb9c86ec2df8a488a85ed973aababe1f38`, exact parent
`53475076e1c79f6d2181877e3238a1a5d389c246`, tree
`82e7b1f5c503659082207481b8339e8113e38c10`, changing only `go.mod` and
`go.sum`. Operator-authorized lifecycle repair commit
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435`, parent `6f4d02e`, tree
`81ca3d7a6537ce2b9c26337c2bf3d4d67c28223a`, is intentionally retained
between that implementation and its documentation handoff. The Repr docs
handoff must be the direct child of f9f0f76. Relative to accepted go-cmp commit
c314bcb, accepted metadata changes add exactly 27 checksum lines. Ordinary and
ignored status must be empty.

Current dependency measurements are 234 selected modules, 3,580 graph edges,
429 native complete-test packages, 1,043 go.sum lines, a 356-line unapplied
tidy projection, and exact Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations 20/30/20. The retained main module declares Go 1.18
and prefers toolchain Go 1.26.7.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT` and verify the
binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate pinned tools beneath scratch as needed. Verified
portable/exact receipts are golangci-lint 2.12.2 official archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
The fresh golangci-lint Darwin arm64 binary was
`691b9100ce968ff0009b6b7757ef6a585e31ae9ab11dfe0340ebb6e8e21fdc3d`;
freshly rebuilt govulncheck v1.7.0 was
`0c536fb25db0d42aff2e5f4496d37a2a89b9f6c449466e2e493975db2652bb2e`.
Their older recorded binary hashes are nonportable build receipts and did not
reproduce; preserve them as history, but prove versions and functionality
rather than requiring byte identity from a new build.

Repr acceptance evidence verified 447 entries; manifest SHA-256 is
`aba1a6ac45d6b8a51a664e0d11f40ccccf82017240c91beaa59ec36866e26711`
and decision-summary SHA-256 is
`900e96c7001a89857ca7c16e1f2572ea0b90db9f2ab4d4e4e6a1c3153808bbcf`.
Canonical latest is stable lightweight tag v0.5.4 at unsigned commit
`9b3680b9bb4e172c4fe539347ac5d0ddf5de0aa9`; its no-requirement module
declares Go 1.18. Exact get changed only Repr's selection, relabeled one
existing main edge, and added its checksum pair. Module and repository gates,
exact quality, and 20/30/20 vulnerability identity pass. Preserve all prior
verified receipts and recorded corrections, including Repr's separately
warmed caches, scratch-local BSD-mktemp adapter, uninherited Docker config,
and scratch-relative manifest verification.

# Role And Boundaries

From fresh external archives and caches, resolve Kong versions through the Go
proxy, checksum database, authoritative upstream repository, and primary Go
vulnerability data. Record exact tag and pseudo-version commits/times, module
Go declarations and requirements, checksum pairs, source identity, tag and
commit signature status, relevant release history, and archived/deprecated
state. Explicitly distinguish stable semantic tags, prereleases, pseudo-
versions, retractions, forks, branch heads, alternate module paths, and
unreleased commits. Do not call a pseudo-version a stable release or select a
newer branch head without proving what exact-path Go queries resolve.

Prove canonical latest and the highest qualified version compatible with Go
1.18 from declarations and the complete changed closure; do not infer
compatibility from a modern build. Measure old versus candidate selected
modules, graph edges, complete package population, checksums, loaded packages
and paths, explicit exact-get diff, and `go mod tidy -diff`. Explain every
selection, edge, and checksum change in the minimal closure. Independently
determine whether Kong is loaded by the main module, only by dependency tests,
or not at all in the complete project population. If loaded, identify real
consumers and exercise the actually used packages and symbols.

Also require candidate module complete tests, repeated tests, race, and vet;
repository build, complete tests/race/vet, pinned lint, byte-identical public
help, identical API/CLI reports, and exact Darwin and Windows vulnerability
populations. Resolve any test-only requirements needed to reproduce the module
suite and keep that apparatus separate from project MVS closure.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity, module
self-tests, loaded behavior, or any repository quality contract fails. If the
selected version is already the exact floor-compatible decision and exact get
does not change a version selection, record no change without manufacturing a
direct requirement or dependency commit. Do not assume vulnerability identity;
prove exact old/candidate IDs and traces.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Treat recorded manifest hashes as historical receipts and recreate only the
evidence required for this decision. Read this archive, rolling handover, P7
roadmap, go.mod/go.sum, accepted Repr, Assert, and Units decisions, retained
Kingpin, Resty, Errgo, YAML v2/v3, go.yaml.in/yaml/v3, jwalterweatherman, cast,
and Cobra/YAML decisions, rejected Check, x/text, x/net, x/image, gotenv, and
afero decisions, prior bounded dependency decisions, and the toolchain,
compatibility, snapshot/Docker, quality, baseline-reproduction, and audit
contracts. Preserve every recorded manifest and setup correction.

# Three Moves

Only if every decision gate passes and the selection changes, use exact Go
1.26.7 and exact `go get github.com/alecthomas/kong@<selected-version>` for
one dependency-only commit. Do not hand-edit module metadata and do not use
tidy as implementation. Preserve every retained dependency selection,
especially `github.com/alecthomas/repr v0.5.4`,
`github.com/alecthomas/assert v1.0.0`,
`github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
`gopkg.in/alecthomas/kingpin.v2 v2.2.6`, `gopkg.in/resty.v1 v1.12.0`,
`gopkg.in/errgo.v2 v2.1.0`,
`gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
`gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
`go.yaml.in/yaml/v3 v3.0.5`, plus language/toolchain declarations, production
source, quality apparatus, and release input.

After a changed selection, run the complete P7 dependency gate: focused
behavior, graph/path, tests/race/vet, pinned lint, help/API/CLI, launcher and
Make contracts, preflight, host plus fresh snapshot/Docker meta and acceptance,
audit meta, focused and exact Q0-Q2 audits, separate full audit, vulnerability
comparison, empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 rows at L2 and zero held, regressed,
not-comparable, or dirty counts. The full audit may exit 1 only for established
queued L3 rows, never 2.

Put every disposable cache, projection, report, generated artifact, build
context, evidence tree, and audit output beneath
`$CODEX_SESSION_SCRATCH_ROOT`; never create direct `/private/tmp/ply-*` roots.
The launcher deletes scratch after every turn, including failure and
interruption. Warm caches only from a separate archive beneath that root and
never run `go mod download all` inside a measured tree. Retain only compact
decisions, digests, and receipts in tracked documents. Preserve the operator-
authorized lifecycle repair and carry the bounded-scratch policy forward.
Never create `.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next measured P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
that next group, launch a successor, push, merge, publish, release, stash,
revert, bypass scratch cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
