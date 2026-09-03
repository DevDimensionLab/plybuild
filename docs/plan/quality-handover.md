# Quality Upgrade Handover

Generated: 2026-09-03T21:23:10+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P7 logrus implementation commit:
  `efc47ca02f9e02e3d31466ef4ddbb58c90772b9c`.
- Its exact parent is the launch continuity commit
  `30a0bfd2a7aa47d8ef6be4c9f7411280f886ce83`, whose exact parent is the terminal
  dependency implementation `0f93a527d1f95cc043a695101746d359c74c9c54`.
- The implementation tree is
  `f8b4d04e75c91981455df97d510fa715be267a47`.
- After this handoff, obtain the new continuity HEAD with `git rev-parse HEAD`;
  its exact parent must be `efc47ca`.
- The implementation and every accepted measurement had empty ordinary and
  ignored status. No push, merge, publication, release, stash, revert,
  successor launch, retained-evidence deletion, image deletion, or worktree
  removal occurred.

## Continuity Checkpoint

P2A-P6 are complete. P7 remains active after its maintained-toolchain move,
terminal dependency group, and logrus dependency group. P8 remains queued.

The answered logrus archive links reciprocally to exactly one NEXT archive for
the three-module Cobra closure. Only the launcher's mutable header and prompt
regions changed; its stable executable skeleton must remain byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. Downloaded tools,
module/build caches, reports, generated files, build contexts, and audit
evidence remain external.

## Logrus Decision

The bounded dependency selection is direct
`github.com/sirupsen/logrus v1.9.0` -> `v1.9.3`.

Independent primary-evidence replay corrected the incoming assertion that
v1.9.3 was the highest release compatible with Go 1.18. Logrus v1.9.4 was
published on 2026-01-15 and declares Go 1.17, but it raises
`github.com/stretchr/testify v1.8.1` to v1.10.0 in this graph. Logrus v1.9.3
declares Go 1.13 and is the highest release that preserves both the retained
`go 1.18` language floor and this session's one-module MVS closure. The v1.10.x
line declares Go 1.23.

The external old/candidate comparison selects 233 modules, 143 production
non-standard dependencies, 197 including tests, and 3,551 graph edges in both
states. Only logrus changes. The historical tidy projection changes from 211
to 217 lines only for the two v1.9.3 checksums. Candidate build, complete tests,
pinned lint, public help, API/CLI, snapshot/Docker meta, and distribution
contracts pass without drift.

GO-2025-4188 affects logrus Writer/WriterLevel scanner variants before v1.9.3.
The repository does not call those symbols. Host and Windows reachable findings
remain the same 22 IDs, while module findings improve from 34 to 33 by removing
only GO-2025-4188.

## Implementation And Compatibility

The exact Go command performed
`go get github.com/sirupsen/logrus@v1.9.3`; dependency metadata was not
hand-edited. Commit `efc47ca` changes only `go.mod` and `go.sum`, with three
insertions and one deletion. It does not change source, public Go API, CLI
output or semantics, toolchain declarations, Docker/release inputs, quality
tools or thresholds, baselines, compatibility allowlists, acceptance/mutation
populations, packaging, publishers, registries, credentials, or P8 code.

The candidate and committed implementation pass focused logrus callers,
complete tests, race, vet, pinned lint, public help, API/CLI and
entry/subprocess compatibility, launcher and Make contracts, complete
preflight, host acceptance, fresh snapshot and Docker meta/acceptance, audit
meta, focused/Q0-Q2/full audits, vulnerability comparison, and empty-HOME
count-2.

The declared and verified toolchain remains exact Go 1.26.7. The retained
official executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; literal `go`
subprocesses must resolve to it. Keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, and do not inject ambient `GOFLAGS`.

Warm fresh external module caches from Git archives outside the worktree. A
real-tree `go mod download all` materializes historical checksum debt and is
not an authorized tidy. The accepted regression includes current and v1.0.1
archive bootstraps for compatibility.

## Accepted External Evidence

The independent wider-release recheck root is
`/private/tmp/ply-p7-logrus-reverify.6YhT0c`. Its verified 12,458-entry manifest
SHA-256 is
`950d152aa8239ce625593a181d74a166f6d7c1c05a3ab34b132463c98362ab8f`.
It records the v1.9.4/testify closure and v1.10.2 Go 1.23 floor.

The accepted v1.9.3 selection root is
`/private/tmp/ply-p7-logrus-selection.30a0bfd.vNfo6P`. Its verified 44,364-entry
manifest SHA-256 is
`678daa930e157e38ffd9b1f08e88eebefc9d00cc150ba8ee153da4820a4f49ad`;
its selection-summary SHA-256 is
`9acb25cede390d34884c97f02f2b06a9ae1f8507271cdfb534c95365b399dda9`.

The commit-bound schema-2 review root is
`/private/tmp/ply-p7-logrus-quality-review.efc47ca.1rtWsz`. Its verified
42,550-entry manifest SHA-256 is
`436f8ef7d05a66aea0cf0473ae95c5f5f58a3bd952de8a6b4db9a9def82eaa29`.
The manual evidence SHA-256 is
`dd1360ea2a24750e96929a4bf465456bb69ca650742a686cb640283e5d3bf679`;
all 242 subjects and 207 declared files bind to `efc47ca`, and all six receipts
validate.

The exact complete quality apparatus root is
`/private/tmp/ply-p7-logrus-quality-gate.efc47ca.1944215232`. Its verified
234,632-entry manifest SHA-256 is
`49c0d0158c7ed5248fa27c33258fe3fe191024f02c684b3050444e567848e40c`.
Exact `make quality` exits 0. The Q0-Q2 scorecard SHA-256 is
`52b675981b8beb2ca37ea7c56f96e99a76064e9dae8d4c510c749203483b20c0`:
all 27 rows pass at L2, mutation and acceptance populations are 8/8 and 4/4,
and held, regressed, current-not-comparable, and dirty counts are zero. Its
Docker artifact's internal evidence manifest is verified at SHA-256
`b162d0f52c4c96e20560212748b7a7d0cd57a4bf72c4b2bb7c99aae9389a0871`.

The independent accepted regression root is
`/private/tmp/ply-p7-logrus-regression-gate.efc47ca-final.mAc8VU`. Its verified
68,942-entry manifest SHA-256 is
`3afd845878b5b3b9874e13c0759aadc7e43d8df2aa5db3b32725c63c635c7c2b`.
All 40 ledger stages pass. The exact Q0-Q2 scorecard has the same hash above;
the full scorecard SHA-256 is
`b512c477b75fc0db63ebfcb4e6016527137454d75cb24ca57d96c731996a1a7e`.
It exits the expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.
Its fresh Docker artifact's internal evidence manifest is verified at SHA-256
`d2faf89ffed2beba2c6d954dbe59477b798f76744e7d4f9ed375d7b5856010f7`.

One first quality attempt is retained as diagnostic evidence: it used an empty
HOME that hid Docker Desktop's buildx plugin and failed at Docker acceptance.
The complete rerun above used the normal Docker context and passed. Two
regression harness-invocation corrections are recorded before their ledger
stages: focused partial audits do not attain a whole level, and module-only
govulncheck accepts no package pattern. Neither was a product or dependency
failure; all 40 final ledger stages pass.

## Next Objective

Implement direct `github.com/spf13/cobra v1.6.1` -> `v1.10.1` and only its
existing-module MVS closure:

- `github.com/cpuguy83/go-md2man/v2 v2.0.2` -> `v2.0.6`;
- `github.com/spf13/pflag v1.0.5` -> `v1.0.9`.

The retained selection root is
`/private/tmp/ply-p7-next-cobra-selection.efc47ca-final.AmKPDY`; its verified
38,241-entry manifest SHA-256 is
`9259d539fae2470870eaeed6644b0728024ef01a2914de2e042c60e0df03eb91`
and its summary SHA-256 is
`f346d139106096cf31de4841efd89484da5c4a32499d73bd381b29a72ebe5ae5`.

As probed on 2026-09-03, Cobra v1.10.1 declares Go 1.15 and retains 233
selected modules, 3,551 graph edges, and the exact 429-package test population.
Only Cobra, go-md2man/v2, and pflag selections move; focused build/test/lint and
byte-identical public help pass. Latest v1.10.2 also declares Go 1.15 but adds
`go.yaml.in/yaml/v3 v3.0.4`, increasing selection to 234 modules and 3,552
edges. It is outside the next smaller group.

Make one dependency-only implementation commit. Stop if independent MVS replay
requires another module, or if the candidate changes behavior, API/CLI,
acceptance, artifacts, vulnerability reachability, or baseline debt. Do not
combine Viper, another dependency, a source fix, or P8 work with this group.

## Start And Stop

Confirm branch, exact ancestry, clean ordinary and ignored status, reciprocal
links, launcher `--check`, the P7/P8 queue, implementation commit/tree, and all
accepted manifests before editing. Read the active archive, this handover, the
P7 roadmap, module graph, Cobra callers/tests, toolchain contract,
compatibility, snapshot/Docker, quality, and audit contracts.

Stop before any dependency outside the three-module Cobra closure,
behavior/API/CLI change, quality-tool upgrade, P8 domain work, inactive
packaging work, publication, publisher/registry/credential change, or release.
Do not push, merge, publish, release, delete retained evidence or images, stash,
revert, launch a successor, or remove the worktree.
