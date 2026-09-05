# Agent Session: Evaluate Gopkg Kingpin V2 Dependency

Status: NEXT
Session ID: `2026-09-05T231214+0200-evaluate-gopkg-kingpin-v2-dependency`
Created: `2026-09-05T23:12:14+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `500ead2f32e51632b360a769f4b734a014a3ebbf6d74021fdbb209afb1f85261`
Previous: [2026-09-05T221559+0200-evaluate-gopkg-resty-v1-dependency.md](2026-09-05T221559+0200-evaluate-gopkg-resty-v1-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected
`gopkg.in/alecthomas/kingpin.v2 v2.2.6` as one bounded dependency group.
Resolve canonical latest, release qualification, and the highest
floor-compatible candidate from primary evidence. Implement one exact changed
selection only if it preserves the retained Go 1.18 floor, has an explained
minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through retained canonical-latest
`gopkg.in/resty.v1 v1.12.0`. All earlier recorded rejections, no-change
decisions, accepted closures, and evidence corrections remain final. Do not
revisit them or combine another module group. P8 remains queued.

The current MVS build list selects
`gopkg.in/alecthomas/kingpin.v2 v2.2.6`. Treat its latest version, tag or
pseudo-version qualification, Go declaration, closure, loaded population,
consumer paths, source history, tests, and vulnerability effect as unknown
until independently resolved. This session may change only gopkg Kingpin v2's
exact required go.mod/go.sum metadata, its minimal MVS closure, and the
roadmap/handoff record. Do not change production Go, another dependency,
language or toolchain declarations, quality apparatus, Docker/release inputs,
packaging, publishers, or P8 code.

# Measurements At Start

Latest implementation remains dependency-only cast commit
`cf4fd4933251b2ca7844b6420fd3cf2a86f049a5`, exact parent
`17f7277fd2512ba06db76d1acb0af5f8b623243c`, clean tree
`485690c010cef3b02385d8f2471d3ef53611abe9`. The Resty v1 documentation
handoff must have exact parent
`a855007ace19f1fc96302000a013cd7d802408ea`. Relative to accepted go-cmp
commit c314bcb, accepted metadata changes remain go-colorful v1.2.0 -> v1.4.1,
go-runewidth v0.0.14 -> v0.0.17, go-toml/v2 v2.0.7 -> v2.2.2, and cast
v1.5.0 -> v1.5.1, with their recorded minimal closures and exactly 15 added
checksum lines. Ordinary and ignored status must be empty.

Current dependency measurements remain 234 selected modules, 3,564 graph
edges, 429 native complete-test packages, 1,031 go.sum lines, a 332-line
unapplied tidy projection, and exact Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations 20/30/20. The retained main module declares Go 1.18
and prefers toolchain Go 1.26.7.

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

Resty v1 no-change evidence is fully verified at
`/private/tmp/ply-p7-resty-v1-selection.a855007.Gg7QC9`, 110,457 entries and
manifest SHA-256
`72d739b05c507fab058737de6ab967b534e5fcdab1c9588582188532f1a926c2`;
decision-summary SHA-256 is
`008b8b9396b58c93e696a5d495037ac974224149fce9b34badb6c7786723a109`.
Selected v1.12.0 is canonical latest and the highest exact-path stable release
compatible with Go 1.18 by declaration. Exact selected-version get only
projects a redundant requirement and checksum, while the exact module suite
fails six legacy net/http and net/url error-string assertions under Go 1.26.7;
dependency metadata therefore remains unchanged. Preserve all prior verified
roots and the recorded go-colorful mutable telemetry, btree regression, cast
whitespace-path, source-archive normalization, Check history-table, Errgo
preflight-runner, and Resty committed-projection/signal-fixture corrections.

# Role And Boundaries

From fresh external archives and caches, resolve gopkg Kingpin v2 versions
through the Go proxy, checksum database, authoritative gopkg metadata and
upstream repository, and primary Go vulnerability data. Record exact tag and
pseudo-version commits/times, module Go declarations and requirements,
checksum pairs, source identity, tag and commit signature status, relevant
release history, and archived/deprecated state. Explicitly distinguish stable
semantic-version tags, prereleases, pseudo-versions, retractions, forks,
branch heads, alternate module paths, and unreleased commits. Do not call a
pseudo-version a stable release or treat a newer branch head or different
module path as an upgrade candidate without proving the version the Go tool
resolves for this exact module path.

Prove canonical latest and the highest qualified version compatible with Go
1.18 from declarations and the complete changed closure; do not infer
compatibility from a modern build. Measure old versus candidate selected
modules, graph edges, complete package population, checksums, loaded packages
and paths, explicit exact-get diff, and `go mod tidy -diff`. Require an
explained minimal selection/edge/checksum closure. Independently determine
whether gopkg Kingpin v2 is loaded by the main module, only by dependency
tests, or not at all in the complete project population. If loaded, identify
real consumers and exercise the actually used packages and symbols.

Also require candidate module complete tests, repeated tests, race, and vet;
repository build, complete tests/race/vet, pinned lint, byte-identical public
help, identical API/CLI reports, and exact Darwin and Windows vulnerability
populations. Resolve any test-only module requirements needed to reproduce the
module suite and separate that test apparatus from the project MVS closure.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity, module
self-tests, loaded behavior, or any repository quality contract fails. If the
selected version is already the exact floor-compatible decision and exact get
does not change a version selection, record that no-change decision without
manufacturing an explicit requirement or dependency commit. Do not assume
vulnerability identity; prove exact old/candidate IDs and traces.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, P7 in
the roadmap, go.mod/go.sum, the retained Resty and Errgo decisions, rejected
Check decision, retained YAML v3 and YAML v2 decisions, rejected x/text,
x/net, x/image, and gotenv decisions, retained jwalterweatherman decision,
accepted cast decision and quality evidence, rejected afero decision, the
accepted `go.yaml.in/yaml/v3` and Cobra/YAML closure decisions, prior bounded
dependency decisions, and the toolchain, compatibility, snapshot/Docker,
quality, baseline-reproduction, and audit contracts. Preserve every recorded
manifest correction.

# Three Moves

Only if every decision gate passes and the selection changes, use exact Go
1.26.7 and exact
`go get gopkg.in/alecthomas/kingpin.v2@<selected-version>` for one
dependency-only commit. Do not hand-edit module metadata and do not use tidy as
implementation. Preserve every retained dependency selection, especially
`gopkg.in/resty.v1 v1.12.0`, `gopkg.in/errgo.v2 v2.1.0`,
`gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
`gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
`go.yaml.in/yaml/v3 v3.0.5`, language/toolchain declarations, production
source, quality apparatus, and release input.

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
