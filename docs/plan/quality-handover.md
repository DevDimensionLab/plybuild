# Quality Upgrade Handover

Generated: 2026-09-03T23:32:49+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P7 Cobra implementation commit:
  `e14ed5ecd0856893c282f37995077542845fb563`.
- Its exact parent is the Cobra launch continuity commit
  `b10076effe353edfae1020e0585175332cc2ea5d`, whose exact parent is the logrus
  implementation `efc47ca02f9e02e3d31466ef4ddbb58c90772b9c`.
- The implementation tree is
  `f17f7cf92e3c52be38f1a5180e283023a76aa75a`.
- After this handoff, obtain the new continuity HEAD with `git rev-parse HEAD`;
  its exact parent must be `e14ed5e`.
- The implementation and every accepted measurement had empty ordinary and
  ignored status. No push, merge, publication, release, stash, revert,
  successor launch, retained-evidence deletion, image deletion, or worktree
  removal occurred.

## Continuity Checkpoint

P2A-P6 are complete. P7 remains active after its maintained-toolchain move,
terminal dependency group, logrus dependency group, and Cobra dependency
group. P8 remains queued.

The answered Cobra archive links reciprocally to exactly one NEXT archive for
the one-module pflag patch. Only the launcher's mutable header and prompt
regions changed; its stable executable skeleton must remain byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. Downloaded tools,
module/build caches, reports, generated files, build contexts, and audit
evidence remain external.

## Cobra Decision

The bounded dependency selection is direct
`github.com/spf13/cobra v1.6.1` -> `v1.10.1` plus exactly:

- `github.com/cpuguy83/go-md2man/v2 v2.0.2` -> `v2.0.6`;
- `github.com/spf13/pflag v1.0.5` -> `v1.0.9`.

Independent primary repository, Go proxy, module, and tag evidence replayed the
incoming choice. Cobra v1.10.1 was published 2025-09-01 and declares Go 1.15;
go-md2man v2.0.6 and pflag v1.0.9 declare Go 1.12. Cobra v1.10.2 remains the
latest release and also declares Go 1.15, but it swaps the existing YAML module
path for `go.yaml.in/yaml/v3 v3.0.4`. That increases selection from 233 to 234
modules and graph edges from 3,551 to 3,552, so v1.10.1 is the highest release
preserving this session's smaller existing-module closure.

Old and selected states retain 233 selected modules, 3,551 graph edges, and a
byte-identical 429-package test population. The historical tidy projection
changes from 217 to 238 lines only for authorized checksums. Candidate build,
complete tests, pinned lint, public help, API/CLI, snapshot/Docker meta, and
distribution contracts pass without drift.

Host and Windows symbol-level vulnerability findings stay at the same 22 IDs;
host module findings stay at the same 33 IDs. The primary Go vulnerability
module index has no entry for Cobra, go-md2man, or pflag at the selected
versions.

## Implementation And Compatibility

The exact Go command performed
`go get github.com/spf13/cobra@v1.10.1`; dependency metadata was not hand-edited.
Commit `e14ed5e` changes only `go.mod` and `go.sum`, with seven insertions and
two deletions. It does not change source, public Go API, CLI output or
semantics, toolchain declarations, Docker/release inputs, quality tools or
thresholds, baselines, compatibility allowlists, acceptance/mutation
populations, packaging, publishers, registries, credentials, or P8 code.

The candidate and committed implementation pass focused Cobra callers,
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

The independently replayed Cobra selection root is
`/private/tmp/ply-p7-cobra-selection.b10076e.94XcUY`. Its verified 43,656-entry
manifest SHA-256 is
`5e8a8cce3fe2495620987bce5bc9ebb98ac6db871aa4d7e14875a00f45f05e8f`;
its selection-summary SHA-256 is
`d93ceaa82a5cc2128574b075bb6ee2d7bb26a1d2d0525e42d55919ff5e18b135`.

The commit-bound schema-2 review root is
`/private/tmp/ply-p7-cobra-quality-review.e14ed5e.SQfQ5c`. Its verified
449-entry manifest SHA-256 is
`a2becf7b8f273162520bccdad21a2ecaeaca867008fd3d1c1e79d1d9b1d9691e`.
The canonical evidence SHA-256 is
`ff89b536faa0ca7f2edf4e3835858d62e62a54697aba961e7ae6526a321fc5b9`;
all declared source digests were independently refreshed and all six receipts
validate. A deliberately invalid digest probe is retained separately inside
the root.

The exact complete quality apparatus root is
`/private/tmp/ply-p7-cobra-quality-gate.e14ed5e.1788465147.56571`. Its verified
234,708-entry manifest SHA-256 is
`a59de1e99bdb27b1bcd1e364f2a5feb1d9ba7c1dfec1a1fc381fd93327749e71`.
Exact `make quality` exits 0. The Q0-Q2 scorecard SHA-256 is
`3756e807d0fb56aee5f49b4a43d60e80526c6ab7f1b8456b388cccedfb20b9b8`:
all 27 rows pass at L2 with valid manual evidence, 8/8 mutation and 4/4
acceptance populations, 80/80 killed mutations, and zero held, regressed,
current-not-comparable, or dirty counts. Its Docker artifact's internal
evidence manifest SHA-256 is
`26127866f5f1d2cce2c4a4a3d36ae2770b0df8f415f86c8b13126b1fd8525a3a`.

The independent accepted regression root is
`/private/tmp/ply-p7-cobra-regression-gate.e14ed5e.H7N4M5`. Its verified
97,882-entry manifest SHA-256 is
`ec18888b35837e6353825cecec7edeb4bd1daa25a7cc428c4ace567130fb80f1`.
All 40 ledger stages pass. The exact Q0-Q2 scorecard has the same hash above;
the full scorecard SHA-256 is
`0f9011741f9b421abdf30f42da16deb4b19b3c6641ceb0a886ca0ce86bb60b8c`.
It exits the expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.
Fresh Docker acceptance built once and ran ten containers without publication.
Its initial attempt inherited the runner's empty HOME and could not discover
Docker Desktop buildx before any build; that attempt is retained verbatim. The
canonical retry used the operator HOME only for plugin discovery, then the
acceptance contract bound the actual build to a fresh external Docker config.

The older logrus selection manifest file still has its recorded SHA-256 and
44,364 entries, but its HOME telemetry count file later changed and now causes
one material verification failure. No tracked or accepted result changed. The
fresh independent Cobra selection above supersedes it for current decisions;
do not cite the older root as currently fully material-verified.

## Next Objective

Implement only existing indirect `github.com/spf13/pflag v1.0.9` -> latest
`v1.0.10`.

The retained selection root is
`/private/tmp/ply-p7-next-pflag-selection.e14ed5e-final.IJI34Q`; its verified
38,164-entry manifest SHA-256 is
`16b28f1dd4fd329b7404d16cb00367a8610f9dd354b3bfed87c9ac4eebf0c5d4`
and its selection-summary SHA-256 is
`cc04bd6d258e132b6ea2fe7c560cde11ead7da43b7673bda0d876441c45cbdde`.

As probed on 2026-09-03, pflag v1.0.10 is latest, was published 2025-09-02,
and declares Go 1.12. It retains 233 selected modules, 3,551 graph edges, and
the exact 429-package population. Only pflag moves; the tidy projection grows
from 238 to 240 lines for its two checksums. Candidate build, complete tests,
pinned lint, byte-identical public help, and exact 22/33/22 vulnerability
populations pass.

Make one dependency-only implementation commit. Stop if independent MVS replay
requires another module or checksum, or if the candidate changes behavior,
API/CLI, acceptance, artifacts, vulnerability reachability, or baseline debt.
Do not combine Viper, another dependency, a source fix, or P8 work with this
group.

## Deferred Viper Probe

A wider probe is retained at
`/private/tmp/ply-p7-next-viper-selection.e14ed5e-final.YYuSN8`; its verified
52,381-entry manifest SHA-256 is
`b17256b1e57146dbb341206667828cce831823b8a0270d97152b9734a225ffe4`
and its summary SHA-256 is
`052e9613b0b7410d70cf9239631768bac93819021bd296070ade3af0a510153a`.

Viper v1.16.0 is the next release and preserves the 429-package population,
but changes 31 selected modules, adds `github.com/google/s2a-go`, grows graph
edges from 3,551 to 3,561, and upgrades `testify` to a module declaring Go
1.20. Viper v1.17.0 further grows selection to 245 modules, graph edges to
3,611, and packages to 437. Viper requires a separate closure/language-floor
decision and is not the next small group.

## Start And Stop

Confirm branch, exact ancestry, clean ordinary and ignored status, reciprocal
links, launcher `--check`, the P7/P8 queue, implementation commit/tree, and all
accepted manifests before editing. Read the active archive, this handover, the
P7 roadmap, module graph, pflag/Cobra callers and tests, toolchain contract,
compatibility, snapshot/Docker, quality, and audit contracts.

Stop before any dependency other than pflag v1.0.10, behavior/API/CLI change,
quality-tool upgrade, P8 domain work, inactive packaging work, publication,
publisher/registry/credential change, or release. Do not push, merge, publish,
release, delete retained evidence or images, stash, revert, launch a successor,
or remove the worktree.
