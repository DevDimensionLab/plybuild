# Agent Session: Evaluate Alecthomas Assert Dependency

Status: NEXT
Session ID: `2026-09-06T150339+0200-evaluate-alecthomas-assert-dependency`
Created: `2026-09-06T15:03:39+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `a81ae99c986ee2424abf9827e06a224fea964b57a38b834d0db00af9e16b94a5`
Previous: [2026-09-06T044730+0200-evaluate-alecthomas-units-dependency.md](2026-09-06T044730+0200-evaluate-alecthomas-units-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected
`github.com/alecthomas/assert v0.0.0-20170929043011-405dbfeb8e38` as one
bounded dependency group. Resolve canonical latest, release qualification,
and the highest floor-compatible candidate from primary evidence. Implement
one exact changed selection only if it preserves the retained Go 1.18 floor,
has an explained minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted canonical-latest
`github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`. All earlier
recorded rejections, no-change decisions, accepted closures, and evidence
corrections remain final. Do not revisit them or combine another module group.
P8 remains queued.

The current MVS build list selects
`github.com/alecthomas/assert v0.0.0-20170929043011-405dbfeb8e38`. A fresh
post-Units survey reports update `v1.0.0` at 2021-12-01T05:52:01Z, but treat
canonical latest, tag and release qualification, Go declaration, closure,
loaded population, consumer paths, source history, tests, and vulnerability
effect as unknown until independently resolved. This session may change only
alecthomas/assert's exact required go.mod/go.sum metadata, its minimal MVS
closure, and the roadmap/handoff record. Do not change production Go, another
dependency, language or toolchain declarations, quality apparatus, Docker/
release inputs, packaging, publishers, or P8 code.

# Measurements At Start

Latest implementation is dependency-only Units commit
`1c874fda104baa36c3a32aaef2d2d2be9689a942`, exact parent
`ae5e731121c17338d08dba4c74a99a1c95a36c49`, clean tree
`744d61deae7d7c3dc25fe1fd1e13a1fc94798981`. The Units documentation handoff
must be its direct child. Relative to accepted go-cmp commit c314bcb,
accepted metadata changes are go-colorful v1.2.0 -> v1.4.1, go-runewidth
v0.0.14 -> v0.0.17, go-toml/v2 v2.0.7 -> v2.2.2, cast v1.5.0 -> v1.5.1,
and Units' 2019 -> 2024 pseudo-version, with recorded minimal closures and
exactly 17 added checksum lines. Ordinary and ignored status must be empty.

Current dependency measurements are 234 selected modules, 3,566 graph edges,
429 native complete-test packages, 1,033 go.sum lines, a 341-line
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

Units acceptance evidence is fully verified at
`/private/tmp/ply-p7-units-selection.ae5e731.smI4xp`, 72,950 entries and
manifest SHA-256
`2cdb64e74a40304c86ef35f6b8e479c49190d4c2d0e6678770bee4e6486493ec`;
decision-summary SHA-256 is
`e30405aeac0eecf6356570f8a49677de1715ad9a7aee15e8efdc9739f55c4cd0`.
Canonical latest is the unreleased default-branch pseudo-version, not a stable
release. It declares Go 1.15 and requires already-selected testify v1.9.0 at
Go 1.17. Its exact one-selection, two-edge, two-checksum closure passed module
and repository gates, exact quality, acceptance, and 20/30/20 vulnerability
identity. Preserve all prior verified roots and the recorded go-colorful
mutable telemetry, btree regression, cast whitespace-path, source-archive
normalization, Check history-table, Errgo preflight-runner, Resty committed-
projection/signal-fixture, Kingpin external-report, and Units launcher-timing
corrections.

# Role And Boundaries

From fresh external archives and caches, resolve alecthomas/assert versions
through the Go proxy, checksum database, and authoritative upstream repository,
plus primary Go vulnerability data. Record exact tag and pseudo-version
commits/times, module Go declarations and requirements, checksum pairs, source
identity, tag and commit signature status, relevant release history, and
archived/deprecated state. Explicitly distinguish stable semantic-version tags,
prereleases, pseudo-versions, retractions, forks, branch heads, alternate
module paths, and unreleased commits. Do not call a pseudo-version a stable
release or select a newer branch head without proving what the Go tool resolves
for this exact module path.

Prove canonical latest and the highest qualified version compatible with Go
1.18 from declarations and the complete changed closure; do not infer
compatibility from a modern build. Measure old versus candidate selected
modules, graph edges, complete package population, checksums, loaded packages
and paths, explicit exact-get diff, and `go mod tidy -diff`. Require an
explained minimal selection/edge/checksum closure. Independently determine
whether alecthomas/assert is loaded by the main module, only by dependency
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
the roadmap, go.mod/go.sum, the retained Kingpin, Resty, and Errgo decisions,
rejected Check decision, retained YAML v3 and YAML v2 decisions, rejected
x/text, x/net, x/image, and gotenv decisions, retained jwalterweatherman
decision, accepted cast decision and quality evidence, rejected afero
decision, the accepted `go.yaml.in/yaml/v3` and Cobra/YAML closure decisions,
prior bounded dependency decisions, and the toolchain, compatibility,
snapshot/Docker, quality, baseline-reproduction, and audit contracts. Preserve
every recorded manifest correction.

# Three Moves

Only if every decision gate passes and the selection changes, use exact Go
1.26.7 and exact
`go get github.com/alecthomas/assert@<selected-version>` for one dependency-only
commit. Do not hand-edit module metadata and do not use tidy as implementation.
Preserve every retained dependency selection, especially accepted
`github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
`gopkg.in/alecthomas/kingpin.v2 v2.2.6`, `gopkg.in/resty.v1 v1.12.0`,
`gopkg.in/errgo.v2 v2.1.0`,
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
