# Agent Session: Evaluate Spf13 Cast Dependency

Status: NEXT
Session ID: `2026-09-05T115915+0200-evaluate-spf13-cast-dependency`
Created: `2026-09-05T11:59:15+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e5e1f4efe7720e0bcb28d42e5adb1d5e7dd2aa30d9e16dbbdb7b2346a3fed549`
Previous: [2026-09-05T110550+0200-evaluate-spf13-afero-dependency.md](2026-09-05T110550+0200-evaluate-spf13-afero-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`github.com/spf13/cast v1.5.0` as one bounded dependency group. Resolve the
canonical latest stable release and highest floor-compatible candidate from
primary evidence. Implement one exact changed selection only if it preserves
the retained Go 1.18 floor, has an explained minimal closure, and passes every
quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through accepted pelletier/go-toml/v2 v2.2.2,
plus the rejected afero v1.10.0 evaluation. Go-cmp v0.7.0, Viper v1.16.0,
Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, latest
gomarkdown `v0.0.0-20260824154242-13c5cf49db8d`, and HCL, mousetrap,
properties, mapstructure, go-toml v1, or afero changes remain rejected or
unnecessary for their recorded floor, loaded-behavior, release-qualification,
self-test, exact-latest, or closure decisions. Do not revisit them or combine
another module group. P8 remains queued.

The current build list selects spf13/cast v1.5.0. Treat its latest release,
Go declaration, release qualification, closure, loaded population, consumer
path, source history, and tests as unknown until independently resolved. This
session may change only cast's exact required go.mod/go.sum metadata, its
minimal MVS closure, and the roadmap/handoff record. Do not change production
Go, jwalterweatherman, gotenv, afero, another dependency,
language/toolchain declarations, quality apparatus, Docker/release inputs,
packaging, publishers, or P8 code.

# Measurements At Start

Latest implementation remains dependency-only go-toml/v2 commit
`8cce879d285f08a3f8710a1ea35358fb057a9bfc`, exact parent
`34f1314af51de0548fdb8ac3214dbd4481aabb24`, clean tree
`8a88c3ff919b5275599d8461c1a176457997d2e7`. The documentation handoff
carrying this archive must have exact parent
`5a4653999ac6a569a95529034a43d457a41dab1b`. Relative to accepted go-cmp
commit c314bcb, go.mod/go.sum change only go-colorful v1.2.0 -> v1.4.1,
go-runewidth v0.0.14 -> v0.0.17, and go-toml/v2 v2.0.7 -> v2.2.2, with
testify v1.9.0 and objx v0.5.2 selected by the last accepted group's minimal
closure. Those three accepted groups add exactly ten checksum lines. Ordinary
and ignored status must be empty.

Current dependency measurements remain 234 selected modules, 3,565 graph
edges, 429 native complete-test packages, a 321-line unapplied tidy projection,
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

Afero decision evidence is fully verified at
`/private/tmp/ply-p7-afero-selection.5a46539.udvGSh`, 79,866 entries and
manifest SHA-256
`f0ad8ce41436a5f8bf48222106e79f57b32a9533b6a9a15b290615a54f162405`;
decision-summary SHA-256 is
`c3301a1c2b13f36956551ed1eb09b1a83bf2ac51cf236c6614f46eaf70ea65de`.
Canonical latest afero v1.15.0 declares Go 1.23.0. Highest floor-compatible
v1.10.0 declares Go 1.16, but its seven-package complete count-10 self-test
deterministically fails in tarfs/TestRead with 27 diagnostics. Selected v1.9.4
normalizes identically, and count-1/race/vet, focused Viper and Ply behavior,
and 20/30/20 vulnerability identity pass; the explicit module-self-test stop
rule nevertheless rejects the change and retains v1.9.4 without metadata edits.

Accepted go-toml/v2 selection, review, exact-quality, and regression evidence
remains fully verified at 17,093/8/236,599/3,743 entries and the exact roots
and manifest hashes recorded in the handover. Go-toml v1 and mapstructure
evidence remains fully verified at
37,719/`e3e516511d60caf7d2004ee1feb29261ecc45b35339dac8e44d19f286606e1b8`
and 46,352/`f15fbc74ed6efd7ee0b3c5ce91b0a9ded4104c5d363aca1272848ec5fca25ecf`.
Accepted runewidth selection/review/exact-quality/regression evidence remains
fully verified at 28,870/15/236,329/4,888 entries and the hashes recorded in
the handover. Properties, mousetrap, and HCL evidence remains fully verified
at 17,240/`941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`,
29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`,
and 21,427/`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`.
Preserve the recorded go-colorful mutable telemetry and btree regression
manifest corrections.

# Role And Boundaries

From fresh external archives and caches, resolve cast releases through the Go
proxy, checksum database, upstream repository, and primary Go vulnerability
data. Record exact tag commits/times, module Go declarations and requirements,
checksum pairs, source identity, tag and commit signature status, relevant
release history, and archived/deprecated status. Explicitly distinguish stable
releases, prereleases, retractions, forks or later major paths, and unreleased
upstream commits. Prove canonical latest and highest stable release compatible
with Go 1.18 from declarations; do not infer compatibility from a modern
toolchain build.

Measure old versus candidate selected modules, graph edges, complete package
population, checksums, loaded packages and path, explicit exact-get diff, and
`go mod tidy -diff`. Require an explained minimal selection/edge/checksum
closure. Independently prove whether cast is loaded; if it is, identify its
real consumers and run focused conversion/configuration behavior that covers
the actually used surface, in addition to candidate module complete tests,
repository build, complete tests/race/vet, pinned lint, byte-identical public
help, and identical API/CLI reports.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity, module
self-tests, loaded behavior, or any repository quality contract fails. If the
selected release is already the exact floor-compatible decision and exact get
is a no-op, record that no-change decision without forcing an implementation
commit. Measure the exact vulnerability population; do not assume parity.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, P7 in
the roadmap, go.mod/go.sum, the rejected afero decision evidence, accepted
go-toml/v2 evidence, prior bounded dependency decisions, and the toolchain,
compatibility, snapshot/Docker, quality, baseline-reproduction, and audit
contracts. Preserve both manifest corrections.

# Three Moves

Only if all decision evidence passes and the selection changes, use exact Go
1.26.7 and exact `go get github.com/spf13/cast@<selected-version>` for one
dependency-only commit. Do not hand-edit module metadata and do not use tidy
as implementation. Preserve every retained selection, including afero v1.9.4,
go-toml/v2 v2.2.2, go-toml v1.9.5, mapstructure v1.5.0, runewidth v0.0.17,
properties v1.8.7, go-colorful v1.4.1, mousetrap v1.1.0, HCL v1.0.0,
go-cmp v0.6.0, btree v1.1.3, gomarkdown fallback, regexp2 v1.12.0, Cobra
v1.10.2, YAML v3.0.5, go-md2man v2.0.7, Blackfriday v2.1.0, pflag v1.0.10,
Uniseg v0.4.7, Colorable v0.1.15, and go-isatty v0.0.20.

After a changed selection, re-run the complete P7 dependency gate: focused
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
