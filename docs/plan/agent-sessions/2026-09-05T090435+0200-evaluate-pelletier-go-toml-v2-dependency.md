# Agent Session: Evaluate Pelletier Go-Toml V2 Dependency

Status: NEXT
Session ID: `2026-09-05T090435+0200-evaluate-pelletier-go-toml-v2-dependency`
Created: `2026-09-05T09:04:35+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `214d44bbd80fdb33c52e26a1b295f116626aada218bfa1b853389808e8ee4412`
Previous: [2026-09-05T082054+0200-evaluate-pelletier-go-toml-dependency.md](2026-09-05T082054+0200-evaluate-pelletier-go-toml-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`github.com/pelletier/go-toml/v2 v2.0.7` as one bounded dependency group.
Resolve the canonical latest release and highest floor-compatible candidate
from primary evidence. Implement one exact changed selection only if it
preserves the retained Go 1.18 floor, has an explained minimal closure, and
passes every quality contract; do not conflate the completed go-toml v1
decision with this distinct `/v2` module.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain,
completed dependency groups through accepted go-runewidth v0.0.17, and the
rejected/no-change HCL, mousetrap, properties, mapstructure, and go-toml v1
evaluations. Go-cmp v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3,
fatih/color v1.19.0, fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, and HCL, mousetrap, properties,
mapstructure, or go-toml v1 selection changes remain rejected or unnecessary
for their recorded floor, loaded-behavior, release-qualification, self-test,
exact-latest, or closure decisions. Do not revisit them or combine another
module group. P8 remains queued.

The current build list selects pelletier/go-toml/v2 v2.0.7. Treat the target
version, latest release, Go declaration, closure, loaded package population,
source history, and test status as unknown until independently resolved. This
session may change only go-toml/v2's exact required go.mod/go.sum metadata and
the roadmap/handoff record. Do not change go-toml v1, production Go, another
dependency, language/toolchain declarations, quality apparatus,
Docker/release inputs, packaging, publishers, or P8 code.

# Measurements At Start

Latest implementation is dependency-only go-runewidth commit
`ca19dcd4da9320112e6c9f0c5db507ccbe8b88e5`, exact parent
`5d0fa3ae0fff98784c349cc23e001018eb9637e1`, clean tree
`273c9b635e988ddd686d6c83846f7320bb9e4e33`. Go-toml v1 required no
implementation commit; the current documentation handoff must have exact
parent `09cc40bcab6f1b03273d726df7dd5ce832f92af6`, whose exact parent is
`f974b581baf36f40bb8fa31868371eb0db712faf`, whose exact parent is ca19dcd.
Relative to accepted go-cmp commit c314bcb, go.mod/go.sum change only
go-colorful v1.2.0 -> v1.4.1 and go-runewidth v0.0.14 -> v0.0.17, adding
exactly those two selected versions' checksum pairs. Ordinary and ignored
status must be empty.

Current dependency measurements are 234 selected modules, 3,557 graph edges,
429 native complete-test packages, a 313-line tidy projection, and exact
Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations 20/30/20.
The retained main module declares Go 1.18 and prefers toolchain Go 1.26.7.

Use exact Go 1.26.7 at
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in PATH, keep GOENV=off, GOWORK=off,
GOTOOLCHAIN=local, and inject no ambient GOFLAGS. Recovery evidence at
`/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` has verified two-entry manifest
SHA-256
`1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
Retain the verified golangci-lint 2.12.2, GoReleaser 2.17.1, apidiff, and
govulncheck v1.7.0 binaries and hashes recorded in the handover.

No-change go-toml v1 evidence is fully verified at
`/private/tmp/ply-p7-go-toml-selection.09cc40b.jGVQz6`, with 37,719 entries
and manifest SHA-256
`e3e516511d60caf7d2004ee1feb29261ecc45b35339dac8e44d19f286606e1b8`.
Its selection-summary SHA-256 is
`85deac527b69ebbfd2608a970684b3f915563d48a538644f3da357961adc268e`.
It proves the original v1 path's canonical latest and highest
Go-1.18-compatible stable release are both selected v1.9.5, with an empty
exact-get closure, no loaded v1 packages, passing repository decision gates,
and no implementation commit. It observes `/v2` only to distinguish module
lineage; independently resolve every `/v2` selection fact. Preserve go-toml
v1.9.5.

No-change mapstructure evidence remains fully verified at
`/private/tmp/ply-p7-mapstructure-selection.f974b58.GbAMTJ`, with 46,352
entries and manifest SHA-256
`f15fbc74ed6efd7ee0b3c5ce91b0a9ded4104c5d363aca1272848ec5fca25ecf`.
Accepted go-runewidth selection, review, exact-quality, and regression evidence
remains fully verified at 28,870/15/236,329/4,888 entries and the manifest
SHA-256 values recorded in the handover. Preserve mapstructure v1.5.0 and
go-runewidth v0.0.17.

Properties decision evidence remains fully verified at
`/private/tmp/ply-p7-properties-selection.469049f.iT6qcB`,
17,240/`941e805f6edf4e7c4d07a8912dcf8e5203dd2cea48530aa3fe6995cb19da097f`.
No-change mousetrap and rejected HCL evidence remain fully verified at
29,283/`71f337b114d1383cd103ef53856a5e767d519eda0b5635e9598b26654a78a152`
and 21,427/`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`.
Preserve the recorded go-colorful mutable telemetry correction and btree
regression-manifest correction.

# Role And Boundaries

From fresh external archives and caches, resolve go-toml/v2 releases through
the Go proxy, checksum database, and upstream repository. Record exact tag
commits/times, module Go declarations and requirements, checksum pairs, source
identity, tag and commit signature status, relevant release history, and
primary Go vulnerability data. Prove which stable `/v2` release is canonical
latest and which is highest compatible with Go 1.18. Distinguish stable
releases, prereleases, retractions, v1 lineage, any later major module path,
and archived/deprecated status explicitly. Do not infer compatibility from a
successful modern-toolchain build.

Measure old versus candidate selected modules, graph edges, complete package
population, checksums, loaded packages/path, explicit exact-get diff, and
`go mod tidy -diff`. Require an explained minimal selection/edge/checksum
closure. Because `/v2` is expected to be loaded, independently prove its
actual package population and path, then run focused Viper TOML decoding and
encoding behavior in addition to candidate module complete tests, repository
build, complete tests/race/vet, pinned lint, byte-identical public help, and
identical API/CLI reports.

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
the roadmap, go.mod/go.sum, go-toml v1 decision evidence, mapstructure and
accepted runewidth evidence, properties decision evidence, accepted
go-colorful evidence, mousetrap and HCL evidence, and the toolchain,
compatibility, snapshot/Docker, quality, baseline-reproduction, and audit
contracts. Preserve both manifest corrections.

# Three Moves

Only if all decision evidence passes and the selection changes, use exact Go
1.26.7 and exact
`go get github.com/pelletier/go-toml/v2@<selected-version>` for one
dependency-only commit. Do not hand-edit module metadata and do not use tidy
as implementation. Preserve every retained selection, including go-toml
v1.9.5, mapstructure v1.5.0, runewidth v0.0.17, properties v1.8.7,
go-colorful v1.4.1, mousetrap v1.1.0, HCL v1.0.0, go-cmp v0.6.0, btree
v1.1.3, gomarkdown fallback, regexp2 v1.12.0, Cobra v1.10.2, YAML v3.0.5,
go-md2man v2.0.7, Blackfriday v2.1.0, pflag v1.0.10, Uniseg v0.4.7,
Colorable v0.1.15, and go-isatty v0.0.20.

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
and make the normal `docs: prepare next agent session` commit. Do not
implement that next group, launch a successor, push, merge, publish, release,
stash, revert, delete retained evidence/images, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
