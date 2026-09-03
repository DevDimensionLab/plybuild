# Quality Upgrade Handover

Generated: 2026-09-04T01:37:13+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- P7 pflag implementation commit:
  `360b2f3c792ca131d259f40852619f2840cefdf1`.
- Its exact parent is the pflag launch continuity commit
  `4f0cc642e9c45960b641133c32a0ffc6b1b3a94b`, whose exact parent is the Cobra
  implementation `e14ed5ecd0856893c282f37995077542845fb563`.
- The implementation tree is
  `a1b941041f5ae685613dd44f4f2db20092cd472c`.
- After this handoff, obtain the continuity HEAD with `git rev-parse HEAD`; its
  exact parent must be `360b2f3`.
- The implementation and accepted measurements had empty ordinary and ignored
  status. No push, merge, publication, release, stash, revert, successor
  launch, retained-evidence deletion, image deletion, or worktree removal
  occurred.

## Continuity Checkpoint

P2A-P6 are complete. P7 remains active after its maintained-toolchain move,
terminal pair, logrus group, Cobra closure, and pflag patch. P8 remains queued.

The answered pflag archive links reciprocally to exactly one NEXT archive for
the measured Viper v1.16.0 closure and its explicit compatibility decision.
Only the launcher's mutable header and prompt regions changed; its stable
executable skeleton must remain byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. Downloaded tools,
module/build caches, reports, generated files, build contexts, and audit
evidence remain external.

## Pflag Decision

Existing indirect `github.com/spf13/pflag v1.0.9` moved to latest `v1.0.10` as
the only selected-module change. Fresh Go proxy and primary repository evidence
confirmed the 2025-09-02 release, tag commit
`0491e5702ad2bb108bc519a5221bcc0f52aa9564`, verified GitHub signature,
Go 1.12 declaration, and exact checksum pair:

- module: `h1:4EBh2KAYBwaONj6b2Ye1GiHfwjqyROoF4RwYO+vPwFk=`;
- go.mod: `h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=`.

Both states retain 233 selected modules, 3,551 graph edges, and a byte-exact
429-package test population. The historical tidy projection changes from 238
to 240 lines only for the two v1.0.10 checksums. Build, complete tests, pinned
lint, help, API/CLI, snapshot/Docker meta, and distribution contracts pass.

The release changes tests and deprecation documentation and makes one
production compatibility correction: `errors.Is(err, ErrHelp)` becomes direct
`err == ErrHelp`, avoiding a Go 1.13 API in a module that declares Go 1.12. The
project does not directly invoke pflag outside the CLI compatibility exporter,
and every executable contract remains unchanged.

Darwin symbol, Darwin module, and Windows symbol vulnerability findings remain
exactly 22/33/22 IDs before and after. The primary Go vulnerability module
index has no pflag entry.

## Implementation And Compatibility

Exact Go 1.26.7 performed
`go get github.com/spf13/pflag@v1.0.10`; dependency metadata was not hand-edited.
Commit `360b2f3` changes only `go.mod` and `go.sum`, with three insertions and
one deletion. It does not change source, public Go API, CLI output or semantics,
toolchain declarations, Docker/release inputs, quality tools or thresholds,
baselines, compatibility allowlists, acceptance/mutation populations,
packaging, publishers, registries, credentials, or P8 code.

Focused pflag/Cobra callers, complete tests, race, vet, pinned lint, public
help, API/CLI and entry/subprocess compatibility, launcher and Make contracts,
complete preflight, host acceptance, fresh snapshot and Docker
meta/acceptance, audit meta, focused/Q0-Q2/full audits, vulnerability parity,
and empty-HOME count-2 pass. One independent preflight attempt retained a
transient signal-fixture failure; the isolated 62-check launcher suite and
complete preflight retry both passed before the stage was accepted.

The declared and verified toolchain remains exact Go 1.26.7. The retained
official executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; literal `go`
subprocesses must resolve to it. Keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, and do not inject ambient `GOFLAGS`.

Warm fresh external module caches from Git archives outside the worktree. A
real-tree `go mod download all` materializes historical checksum debt and is
not an authorized tidy. Keep measured old/candidate archives separate from the
bootstrap archive used to warm caches.

## Accepted External Evidence

Independent pflag selection evidence is retained at
`/private/tmp/ply-p7-pflag-selection.4f0cc64.58RgL0`. Its final verified
44,617-entry manifest SHA-256 is
`b64cad473c92a5551e997010f88243d46fafa2ddc65c11db5d2d6bb90786298c`;
its selection-summary SHA-256 is
`e47aa25c01cfaab27409b83b6cbf809e81a315e55725744c17a447df8f96bf79`.

The original selection manifest SHA-256
`c7712aafaf9d8440cd2661ce6ffaa4a638321ddbdb9719b49a187420276b430d`
verified fully when sealed. Follow-on implementation validation then reused
that Go build cache and rewrote 27 action-index files. Stable selection output
did not change. The original manifest and passing log, later mismatch log, an
incident note, and the final post-cache manifest are retained. Cite the final
manifest as current.

The commit-bound schema-2 review root is
`/private/tmp/ply-p7-pflag-quality-review.360b2f3.57B9Z6`. Its verified
447-entry manifest SHA-256 is
`caeace41486ddc6797b5fc737ccd5a4757d2052eb3a5c3451f7da0bcc39f04b8`.
The canonical evidence SHA-256 is
`d3c5d725e99fb633186277cd66fc9dc1f17b1c8f0705926b70c449f990526c69`;
all declared source digests were refreshed and all six receipts validate.

The exact complete quality apparatus root is
`/private/tmp/ply-p7-pflag-quality-gate.360b2f3.1788472549.40415`. Its verified
234,799-entry manifest SHA-256 is
`47971b2e6cf9399f085e233d9a264f9d22c2e5f45dd024375fbb8e151cbd36c6`.
Exact `make quality` exits 0. The Q0-Q2 scorecard SHA-256 is
`f50ee942007f4cf5b76b0d87f610469f9d95dc04c6363779e2d6e83a7c632934`:
all 27 rows pass at L2 with valid manual evidence, 8/8 mutation and 4/4
acceptance populations, 80/80 killed mutations, and zero held, regressed,
current-not-comparable, or dirty counts. Fresh Docker acceptance built once,
ran ten containers, and published zero times; its internal evidence manifest
SHA-256 is
`3d0f1e342dbd9e218d371e29914f7aafdec93d81920b9cfd3f06a3f77ff094d5`.

The independent regression root is
`/private/tmp/ply-p7-pflag-regression-gate.360b2f3.gb4Rvk`. Its verified
92,211-entry manifest SHA-256 is
`ec5eb0c20975138fd96b4de3c10f855c0ceb11a5f1152039be2586025beec857`.
All 40 ordered stages pass. The exact Q0-Q2 scorecard has the same hash above;
the full scorecard SHA-256 is
`4f72c988701326e4a21749306a6155fa849ac3c31ba4e0b293ec6a9f3f92bf08`.
It exits the expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.
Fresh regression Docker acceptance built once, ran ten containers, published
zero times, and recorded image ID
`sha256:7425c458a37c57140ec1a114b7cea332513621384772dff4b5939e6207885867`.

## Next Objective

Evaluate and, only after a truthful compatibility decision, implement direct
`github.com/spf13/viper v1.15.0` -> `v1.16.0` with its exact measured MVS
closure. Do not combine any source fix, another dependency group, a language
floor change, or P8 work.

The current-tree selection root is
`/private/tmp/ply-p7-next-viper-selection.360b2f3-final.oGxITF`. Its verified
40,078-entry manifest SHA-256 is
`1390745861503e9f57169fff6054db637e14aa6f7ada9b708c9b9155879253a4`;
its selection-summary SHA-256 is
`d3b2704feeab8883cbc67ab98d5aebdb880ee601fde8157aba4fc10ca1d340f5`.

Viper v1.16.0 is the smallest next stable release. Fresh proxy and repository
evidence resolves tag commit
`21a7fd828ed231bbe62068d6aafa5aa9f85dc79e`; the module declares Go 1.17.
Against `360b2f3`, it changes 31 selected modules, adds only
`github.com/google/s2a-go v0.1.3`, grows 233/3,551 modules/edges to 234/3,561,
and preserves the exact 429-package population. The tidy projection grows from
240 to 270 lines. Candidate build, complete tests, pinned lint, byte-identical
help, and exact 22/33/22 vulnerability populations pass.

The closure selects `github.com/stretchr/testify v1.8.3`, whose module declares
Go 1.20. Exact Go 1.26.7 builds it while Viper and main-module production
sources remain at Go 1.17 and the retained `go 1.18` language declaration.
Independently decide and record whether that transitive module declaration is
compatible with the project's stated floor. Stop and hand off without an
implementation if preserving the floor is not truthful; do not silently raise
the main-module `go` line.

Latest Viper is v1.21.0 and declares Go 1.23.0. The retained wider probe shows
that already v1.17.0 expands selection to 245 modules and the package population
to 437. Keep this group fixed at v1.16.0 and its exact 31-selection closure.

## Start And Stop

Confirm branch, exact ancestry, clean ordinary and ignored status, reciprocal
links, launcher `--check`, the P7/P8 queue, implementation commit/tree, and all
accepted manifests before editing. Read the active archive, this handover, the
P7 roadmap, module graph, Viper/config callers and tests, toolchain and
baseline-reproduction contracts, compatibility, snapshot/Docker, quality, and
audit contracts.

Stop before any dependency outside the exact measured Viper v1.16.0 closure,
behavior/API/CLI change, main-module language-floor change, quality-tool
upgrade, P8 domain work, inactive packaging work, publication,
publisher/registry/credential change, or release. Do not push, merge, publish,
release, delete retained evidence or images, stash, revert, launch a successor,
or remove the worktree.
