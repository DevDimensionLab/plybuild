# Quality Upgrade Handover

Generated: 2026-09-04T02:16:35+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Last P7 implementation commit:
  `360b2f3c792ca131d259f40852619f2840cefdf1`.
- Its exact parent is
  `4f0cc642e9c45960b641133c32a0ffc6b1b3a94b`; its clean tree is
  `a1b941041f5ae685613dd44f4f2db20092cd472c`.
- The Viper decision handoff starts at
  `27dbc7819afdcf37c4c52cdf3cde2b5b042ba89d`, whose exact parent is
  `360b2f3`. After this handoff, obtain the continuity HEAD with
  `git rev-parse HEAD`; its exact parent must be `27dbc78`.
- The decision started with empty ordinary and ignored status. No push, merge,
  publication, release, stash, revert, successor launch, retained-evidence or
  image deletion, or worktree removal occurred.

## Continuity Checkpoint

P2A-P6 are complete. P7 remains active after its maintained-toolchain move,
terminal pair, logrus group, Cobra closure, pflag patch, and rejected Viper
v1.16.0 closure decision. P8 remains queued.

The answered Viper archive links reciprocally to exactly one NEXT archive for
the measured Uniseg v0.4.7 group. Only the launcher's mutable header and prompt
regions changed; its stable executable skeleton remains byte-identical.

The tracked launcher/archive apparatus remains the task source. Do not create
`.agent-task/current.md` or `.quality/manual-evidence.json`. Downloaded tools,
module/build caches, graph projections, vulnerability output, generated files,
build contexts, schema-2 evidence, reports, and audit evidence remain external.

## Viper Decision

Direct `github.com/spf13/viper v1.15.0` did not move to v1.16.0. Independent
exact-Go replay reproduced the retained candidate byte-for-byte: 233 -> 234
selected modules, 3,551 -> 3,561 graph edges, exactly 31 changed selections,
only `github.com/google/s2a-go v0.1.3` added, the exact 429-package population
retained, and the historical tidy projection growing from 240 to 270 lines.

The official Go module reference defines each module's `go` directive as its
minimum required Go version. P7's established compatibility rule retains a Go
1.18 dependency floor. Reading all 31 target module files found eleven modules
declaring Go 1.19 and `github.com/stretchr/testify v1.8.3` declaring Go 1.20.
The Go-1.19 modules are `cloud.google.com/go`, its `compute` and `longrunning`
modules, `github.com/googleapis/enterprise-certificate-proxy`,
`github.com/googleapis/gax-go/v2`, `github.com/hashicorp/consul/api`, the three
selected etcd v3 modules, `google.golang.org/api`, and
`google.golang.org/genproto`. Exact Go 1.26.7 build success cannot make the
lower Go 1.18 floor truthful, so the candidate was rejected before editing.

Fresh primary evidence confirms Viper v1.16.0's verified annotated tag object
`c19a006378aa5ee373af48ddee7fdf78d62c5c06`, peeled commit
`21a7fd828ed231bbe62068d6aafa5aa9f85dc79e`, 2023-05-30 release, Go 1.17
declaration, and checksum pair:

- module: `h1:rGGH0XDZhdUOryiDWjmIvUSWpbNqisK8Wk0Vyefw8hc=`;
- go.mod: `h1:yg78JgCJcbrQOvV9YLXgkLaZqUidkY9K+Dd1FofRzQg=`.

Exact Go 1.26.7 build, complete tests, pinned golangci-lint 2.12.2, public help,
API/CLI reports, module verification, toolchain declarations, distribution,
and snapshot/Docker meta-contracts pass for old and candidate states. API and
CLI report SHA-256 values remain
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
Fresh Go vulnerability data preserves exact 22/33/22 Darwin-symbol,
Darwin-module, and Windows-symbol ID populations; the primary module index has
no Viper or Testify entry.

No source, public Go API, CLI, dependency metadata, toolchain declaration,
Docker/release input, quality tool or threshold, numeric baseline, compatibility
allowlist, acceptance/mutation population, packaging, publisher, registry,
credential, or P8 file changed. There is no Viper implementation commit and the
main module still declares Go 1.18 with Cobra v1.10.1 and pflag v1.0.10.

## Toolchain And Method

The declared and verified toolchain remains exact Go 1.26.7. The retained
official executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; literal `go`
subprocesses must resolve to it. Keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, and do not inject ambient `GOFLAGS`.

Warm fresh external module caches from separate Git archives outside the
worktree. Never run `go mod download all` inside a measured old/candidate tree
or use a real-tree download as a tidy surrogate. A separate bootstrap archive
may materialize historical checksum debt only to warm an offline cache; retain
and diff that materialization.

Pinned tools retained compatibility: golangci-lint 2.12.2, GoReleaser 2.17.1,
Cobra v1.10.1, and pflag v1.0.10. The preferred tool hashes and paths are in the
sealed decision root.

## Accepted External Evidence

The Viper decision and next-group selection root is
`/private/tmp/ply-p7-viper-decision.360b2f3.Wdvj52`. Its verified 69,347-entry
manifest SHA-256 is
`0b1e054f4fcfae3a630e5192603cb02f53d2bdd3b7623c3fe14bdb3b74eb757f`;
`decision-summary.json` SHA-256 is
`cd8cc895ca5979514e5e584b023ac70d9967cf4951ecf81ddb760267e10cbd7d`.

Independent pflag selection evidence remains at
`/private/tmp/ply-p7-pflag-selection.4f0cc64.58RgL0`. Read its incident note
before using it: later build-cache reuse rewrote 27 action-index files, so its
final post-cache 44,617-entry manifest is authoritative. That manifest's
SHA-256 is
`b64cad473c92a5551e997010f88243d46fafa2ddc65c11db5d2d6bb90786298c`.

The exact accepted pflag quality root is
`/private/tmp/ply-p7-pflag-quality-gate.360b2f3.1788472549.40415`. Its verified
234,799-entry manifest SHA-256 is
`47971b2e6cf9399f085e233d9a264f9d22c2e5f45dd024375fbb8e151cbd36c6`.
Exact `make quality` exits 0. The Q0-Q2 scorecard SHA-256 is
`f50ee942007f4cf5b76b0d87f610469f9d95dc04c6363779e2d6e83a7c632934`:
all 27 rows pass at L2 with 8/8 mutation and 4/4 acceptance populations and
zero held, regressed, current-not-comparable, or dirty counts.

The independent regression root is
`/private/tmp/ply-p7-pflag-regression-gate.360b2f3.gb4Rvk`. Its verified
92,211-entry manifest SHA-256 is
`ec5eb0c20975138fd96b4de3c10f855c0ceb11a5f1152039be2586025beec857`.
All 40 ordered stages pass. The full scorecard SHA-256 is
`4f72c988701326e4a21749306a6155fa849ac3c31ba4e0b293ec6a9f3f92bf08`;
it exits expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7.

## Next Objective

Evaluate and, only after independent replay and compatibility review, upgrade
indirect `github.com/rivo/uniseg v0.4.4` to latest v0.4.7 as one exact selected
module. Do not combine any source fix, another dependency group, language-floor
change, Viper reconsideration, or P8 work.

The current-tree external probe selects only Uniseg v0.4.7. Both states retain
233 selected modules, 3,551 graph edges, and the exact 429-package population;
the historical tidy projection grows from 240 to 244 lines. V0.4.7 declares Go
1.18, exact Go 1.26.7 build, complete tests, pinned lint, and byte-identical
help pass, and its checksum pair is:

- module: `h1:WUdvkW8uEhrYfLC4ZzdpI2ztxP1I582+49Oc5Mq64VQ=`;
- go.mod: `h1:FN3SvrM+Zdj16jyLfmOkMNblXMcoc8DfTHruCPUcx88=`.

This was the smallest compatible next candidate in the bounded probe. Emoji
v2.2.14 changes one selection but declares Go 1.21; colorable v0.1.15 changes
two selections; and `gopkg.in/ini.v1 v1.67.3` changes three. Independently
reverify Uniseg release/tag identity, module declarations, checksums, exact
selection, behavior, artifacts, and vulnerability evidence before editing.

## Start And Stop

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
links, launcher `--check`, P7/P8 queue, implementation commit/tree, and all
accepted manifests before editing. Read the active archive, this handover, the
P7 roadmap, module graph and Uniseg callers, toolchain/baseline reproduction,
compatibility, snapshot/Docker, quality, and audit contracts.

Stop before any dependency beyond exact Uniseg v0.4.7, behavior/API/CLI change,
main-module language-floor change, quality-tool upgrade, Viper retry, P8 domain
work, inactive packaging, publication, publisher/registry/credential change,
or release. Do not push, merge, publish, release, delete retained evidence or
images, stash, revert, launch a successor, or remove the worktree.
