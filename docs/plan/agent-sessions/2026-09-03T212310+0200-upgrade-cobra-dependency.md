# Agent Session: Upgrade Cobra Dependency Group

Status: NEXT
Session ID: `2026-09-03T212310+0200-upgrade-cobra-dependency`
Created: `2026-09-03T21:23:10+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3897f37dd49051769f2c4e6e1637dcf4718f358518d2853708f0479723a410dd`
Previous: [2026-09-02T142312+0200-upgrade-logrus-dependency.md](2026-09-02T142312+0200-upgrade-logrus-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its next smallest dependency group: upgrade direct
`github.com/spf13/cobra v1.6.1` to `v1.10.1` and only its existing-module MVS
closure, `github.com/cpuguy83/go-md2man/v2 v2.0.2` to `v2.0.6` and
`github.com/spf13/pflag v1.0.5` to `v1.0.9`. Independently reverify the release
decision from primary evidence, preserve behavior and every quality contract,
and finish with no unrelated drift.

# Authorized Roadmap

P2A-P6 are complete. The Go 1.26.7 toolchain, terminal pair, and logrus group
are complete; P7 remains active for dependency groups and P8 remains queued.
This session may change only the Cobra closure's required `go.mod`/`go.sum`
metadata, a focused executable contract if one is needed, and the roadmap
record.

The external probe selected v1.10.1 as the highest Cobra release that preserves
the smaller existing-module closure and the retained `go 1.18` language floor.
Cobra v1.10.1 declares Go 1.15. Latest v1.10.2 also declares Go 1.15, but adds
`go.yaml.in/yaml/v3 v3.0.4`, growing selection from 233 to 234 modules. Do not
silently widen this group; stop and record the decision if independent replay
does not reproduce the three-module closure.

# Measurements At Start

The logrus implementation commit is
`efc47ca02f9e02e3d31466ef4ddbb58c90772b9c`, exact parent
`30a0bfd2a7aa47d8ef6be4c9f7411280f886ce83`, clean tree
`f8b4d04e75c91981455df97d510fa715be267a47`. After the handoff, the continuity
HEAD must have exact parent `efc47ca`, and ordinary and ignored status must be
empty.

The declared and verified toolchain remains Go 1.26.7. The retained official
executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; audit
reproduction invokes literal `go`. Do not inject ambient `GOFLAGS`.

The accepted logrus selection root is
`/private/tmp/ply-p7-logrus-selection.30a0bfd.vNfo6P`, with verified
44,364-entry manifest SHA-256
`678daa930e157e38ffd9b1f08e88eebefc9d00cc150ba8ee153da4820a4f49ad`.
The next-group Cobra selection root is
`/private/tmp/ply-p7-next-cobra-selection.efc47ca-final.AmKPDY`, with verified
38,241-entry manifest SHA-256
`9259d539fae2470870eaeed6644b0728024ef01a2914de2e042c60e0df03eb91`.

The exact logrus quality root is
`/private/tmp/ply-p7-logrus-quality-gate.efc47ca.1944215232`, with verified
234,632-entry manifest SHA-256
`49c0d0158c7ed5248fa27c33258fe3fe191024f02c684b3050444e567848e40c`.
Its Q0-Q2 scorecard SHA-256 is
`52b675981b8beb2ca37ea7c56f96e99a76064e9dae8d4c510c749203483b20c0`:
all 27 criteria pass at L2, populations are 8/8 and 4/4, and held, regressed,
not-comparable, and dirty counts are zero.

The independent logrus regression root is
`/private/tmp/ply-p7-logrus-regression-gate.efc47ca-final.mAc8VU`, with verified
68,942-entry manifest SHA-256
`3afd845878b5b3b9874e13c0759aadc7e43d8df2aa5db3b32725c63c635c7c2b`.
All 40 stages pass. Its full report exits 1 only for queued L3 rows Q3.1,
Q3.3, Q3.4, and Q3.7; it is not a P7 dependency-group exit gate.

# Role And Boundaries

Verify time-sensitive Cobra, go-md2man, and pflag releases, module `go`
requirements, release metadata, and relevant vulnerability information from
primary Go module, repository, and Go vulnerability sources. Record why the
selected closure is compatible with pinned Go 1.26.7, golangci-lint 2.12.2,
GoReleaser 2.17.1, and the retained `go 1.18` language floor.

Do not change production Go behavior, public Go API, CLI semantics, toolchain
declarations, Docker or release inputs, any dependency outside the three-module
closure, quality-tool versions, quality thresholds, baseline numeric debt,
compatibility allowlists, acceptance or mutation populations, publishers,
registries, credentials, inactive packaging, or P8 domain code. Do not combine
a source fix with this dependency move; stop on an incompatible behavior
change.

Keep module/build caches, graph projections, vulnerability results, reports,
generated artifacts, build contexts, schema-2 evidence, and audit output
outside the worktree. Warm fresh module caches from external Git archives;
never use a real-tree `go mod download all` as a tidy surrogate. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint block before
editing. Verify the accepted manifests. Read the rolling handover, this
archive, the P7 roadmap, `go.mod`/`go.sum`, Cobra callers and tests, toolchain
and baseline-reproduction contracts, compatibility, snapshot/Docker, quality,
and audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, `go mod tidy -diff` projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, and a
vulnerability baseline. Probe v1.10.1 and v1.10.2 in separate external copies.
Compare old/new graphs exactly and require every changed module and checksum to
belong to Cobra, go-md2man/v2, or pflag.

Make one focused dependency-only implementation commit using the exact Go
tool, without hand-editing dependency metadata. Re-run focused Cobra callers
and tests, pinned lint, complete tests/race/vet, API/CLI and entry/subprocess
compatibility, launcher and Make contracts, complete preflight, host
acceptance, fresh snapshot/Docker meta and acceptance, audit meta, focused and
exact Q0-Q2 audits, the separate full audit, vulnerability comparison, and
empty-HOME count-2. Refresh external schema-2 evidence when commit binding
requires it.

Require exact `make quality` exit 0 at L2 with all 27 Q0-Q2 rows present,
8/8 mutation and 4/4 acceptance populations, and zero held, regressed,
not-comparable, or dirty counts. The full report may exit 1 only for the four
queued L3 rows, never 2. Retain immutable external evidence and verified
manifests for graph selection, quality, regression, and artifacts.

# Automatic Handoff

After this group succeeds, keep P7 active and P8 queued. Rewrite the rolling
handover and roadmap, answer this archive, create exactly one reciprocal NEXT
archive for the next evidence-selected small dependency group, replace only
the launcher's mutable regions, run launcher and handoff contracts, and make
the normal `docs: prepare next agent session` commit. Do not implement the next
group in this session.

Do not launch a successor, push, merge, publish, release, delete retained
evidence or local images, stash, revert, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
