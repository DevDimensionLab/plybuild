# Quality Upgrade Handover

Generated: 2026-09-04T06:54:19+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Last P7 implementation commit:
  `33e187317c6c1be79ef8c5b64ddb2a0f8caec71f`.
- Its exact parent is
  `e1fd64e4a7bebd6dabf529675b1560172f12f1d5`; its clean tree is
  `6942bd34879b9935382b44ffab682606cf5ef5c9`.
- The Cobra/YAML handoff is the next commit after `33e1873`; require its exact
  parent to be `33e187317c6c1be79ef8c5b64ddb2a0f8caec71f` and begin with empty
  ordinary and ignored status.
- No push, merge, publication, release, stash, successor launch, retained
  evidence or image deletion, or worktree removal occurred.

## Continuity Checkpoint

P2A-P6 are complete. P7 remains active after the maintained-toolchain move,
terminal pair, logrus, Cobra v1.10.1 closure, pflag, Uniseg, and
Colorable/go-isatty groups and the rejected Viper, Emoji, and ini candidates.
P8 remains queued.

The answered Colorable archive links reciprocally to exactly one NEXT archive
for the measured Cobra v1.10.2/YAML v3.0.4 closure. Only the launcher's mutable
header and prompt regions may change; its stable executable skeleton must
remain byte-identical.

The tracked launcher/archive apparatus remains the task source. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`. Keep downloaded
tools, module/build caches, graph projections, vulnerability output, generated
artifacts, build contexts, schema-2 evidence, reports, and audit evidence
outside the worktree.

## Colorable Completion

Existing indirect `github.com/mattn/go-colorable v0.1.13` moved to latest
v0.1.15 with exact MVS companion `github.com/mattn/go-isatty v0.0.17` ->
v0.0.20. Fresh primary evidence confirms Colorable's lightweight v0.1.15 tag
at `8bf39a204f13f0cfcf86ab9b297c3d6e0668e54a`, Go 1.18 declaration, and
checksum pair:

- module: `h1:+u9SLTRGnXv73cEsnsmoZBom+dMU88B2M0aDcWy0/jY=`;
- go.mod: `h1:6LmQG8QLFO4G5z1gPvYEzlUgJ2wF+stgPZH1UqBm1s8=`.

Go-isatty v0.0.20 has lightweight tag commit
`a7c02353c47bc4ec6b30dc9628154ae4fe760c11`, declares Go 1.15, and has
checksum pair:

- module: `h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=`;
- go.mod: `h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=`.

Both states retain 233 selected modules, 3,551 graph edges, and the exact
429-package population. Only those two selections change; historical
`go mod tidy -diff` grows from 244 to 252 lines. The caller chain remains
`plybuild/cmd` -> `go-term-markdown` -> `fatih/color` -> `go-colorable`, and
existing behavior/help/API/CLI/artifact contracts cover the indirect boundary.

Exact `go get github.com/mattn/go-colorable@v0.1.15` produced implementation
commit `33e187317c6c1be79ef8c5b64ddb2a0f8caec71f`. Only `go.mod` and `go.sum`
changed. Production behavior, public API/CLI, toolchain and main-module
language declarations, Docker/release inputs, quality tools and thresholds,
numeric debt, compatibility allowlists, acceptance/mutation populations,
packaging, publishers, registries, credentials, and P8 code did not change.

Exact Go 1.26.7 build, focused callers, complete tests/race/vet, pinned lint,
byte-identical public help, API/CLI and entry/subprocess compatibility,
launcher/Make, complete preflight, host and fresh snapshot/Docker acceptance,
audit meta, focused and Q0-Q2 audits, vulnerability parity, empty-HOME
count-2, and final cleanliness all pass. Vulnerability sets remain exactly 22
Darwin symbol, 33 Darwin module, and 22 Windows symbol IDs.

## Toolchain And Method

The declared and verified toolchain remains exact Go 1.26.7. The retained
official executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; literal `go`
subprocesses must resolve to it. Keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, and do not inject ambient `GOFLAGS`.

Retained tools are golangci-lint 2.12.2 at
`/private/tmp/ply-http-gate.gzWxRY/bin/golangci-lint`, SHA-256
`3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`;
GoReleaser 2.17.1 at
`/private/tmp/ply-snapshot-probe.CiDxj0/tool/goreleaser`, SHA-256
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`;
apidiff at `/private/tmp/ply-http-gate.gzWxRY/bin/apidiff`, SHA-256
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`;
and govulncheck v1.7.0 at
`/private/tmp/ply-p7-terminal-selection.ctFhJk/tools/bin/govulncheck`, SHA-256
`0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

Warm fresh module caches from separate external Git archives. Never run
`go mod download all` inside a measured old/candidate tree or use a real-tree
download as a tidy surrogate. An aggregate coverage invocation can expand the
sparse historical `go.sum`; use an exact readonly Go invocation for that check
and retain any rejected projection outside the worktree.

## Accepted External Evidence

The Colorable selection root is
`/private/tmp/ply-p7-colorable-selection.e1fd64e.QlohU1`. Its verified
49,635-entry manifest SHA-256 is
`b8ec3adb51fc916820b300f0faf4ff0eddc7e1ae51a9f839b96ceed9c57cebb7`;
selection-summary SHA-256 is
`484e11cb3c6929dfd0a50c050a49008df4ed3908133ae3d87a3a3839f783e07c`.

Commit-bound schema-2 review evidence is retained at
`/private/tmp/ply-p7-colorable-quality-review.33e1873.eGDQR5`. Its verified
21-entry manifest SHA-256 is
`f6294878de4b3f19a8731b0b94485f8a203ab20137e37fd3e866f27a0afde24b`;
manual evidence SHA-256 is
`cfbe47aa455636613ec5f03c4ede4a12a4169e2a0da9483cf6fb0519e9a3ce66`,
and focused scorecard SHA-256 is
`566261a3e2f6e0bf3910a77850a1977430282dd2bd56fd3290d688c177848a11`.

The exact accepted quality root is
`/private/tmp/ply-p7-colorable-quality-parent.33e1873.31b10h/quality-gate-retry3`.
Its verified 235,059-entry manifest SHA-256 is
`6695ee2f2a181c46fbcf953518c2fd9581902f3f3be28bd761602eecab0785ae`.
Exact `make quality` exits 0 across the 21-stage ledger. The Q0-Q2 scorecard
SHA-256 is
`55cf7f1a4867f54bba31108f2d7ef104f1b0901b44a5ac29aea11e50aec007e6`:
all 27 rows pass at L2 with valid evidence, 80/80 killed mutations, 8/8
mutation and 4/4 acceptance populations, and zero held, regressed,
not-comparable, or dirty counts.

Two earlier quality output roots are retained as noncanonical. An over-narrow
PATH forced the system Python and exposed the launcher signal-log race; the
isolated 62-check launcher contract and the full corrected-PATH retry passed.

The independent regression root is
`/private/tmp/ply-p7-colorable-regression-gate.33e1873.1Yb6le`. Its verified
98,451-entry manifest SHA-256 is
`cdf2963affe5e2395ab8c7911efeaa525e4f39e700ede3945e484cfacb63802c`.
All 40 ordered stages pass. Its full scorecard SHA-256 is
`c6e7a71ac8cf187654a83f043a3967d53f8e51035c69bf717909117cd7b8ff1b`;
it exits expected 1, never 2, only for queued Q3.1, Q3.3, Q3.4, and Q3.7. The
root retains the exact 208-line sparse-checksum incident patch and its clean
readonly replay.

## Rejected Ini Decision

`gopkg.in/ini.v1 v1.67.0` -> latest v1.67.3 is rejected before
implementation. It selects Testify v1.11.1 and Objx v0.5.2; Objx declares Go
1.20, exceeding the retained Go 1.18 floor. The technical candidate still
builds and tests, passes pinned lint and byte-identical help, and retains
22/33/22 vulnerability parity, but those passes do not replace the floor rule.

The sealed decision root is
`/private/tmp/ply-p7-next-selection.33e1873.mRehC2`. Its verified 39,374-entry
manifest SHA-256 is
`c7af82cbe2327de51dc93fcf48c61e65e6b38e5d51a3a7556c07d6d753016181`;
decision-summary SHA-256 is
`04081e0b5f0b3324e49ca321f5edf0a2be5b9b938bce31f68c30aabb89c8a90d`.

## Next Objective

Independently reverify and, only if exact and compatible, upgrade direct
`github.com/spf13/cobra v1.10.1` to latest v1.10.2 with newly selected exact
MVS companion `go.yaml.in/yaml/v3 v3.0.4`. Do not combine a source fix, other
dependency group, language-floor change, Viper/Emoji/ini retry, or P8 work.

The authoritative selection root is
`/private/tmp/ply-p7-next-cobra-selection.33e1873.ZPppZP`. Its verified
43,788-entry manifest SHA-256 is
`d522c38eafcd59bd8b172233f06d79a6f32f2410f959e9474c0b8bd0fb8c7a38`;
selection-summary SHA-256 is
`3bf4208b4b8d3fd1d300976274dc6e26f26322aec365554bc7c7fd0e4a377b85`.

The clean offline replay changes exactly Cobra v1.10.1 -> v1.10.2 and adds
YAML v3.0.4, moves 233 -> 234 modules and 3,551 -> 3,552 graph edges,
preserves the exact 429-package population, and grows tidy from 252 to 254
lines. Cobra declares Go 1.15; YAML declares Go 1.16; the main module remains
Go 1.18. Exact build, complete tests, pinned lint, byte-identical help,
identical API/CLI reports, and vulnerability parity pass.

Cobra's lightweight v1.10.2 tag is exact commit
`88b30ab89da2d0d0abb153818746c5a2d30eccec`; its checksums are:

- module: `h1:DMTTonx5m65Ic0GOoRY2c16WCbHxOOw6xxezuLaBpcU=`;
- go.mod: `h1:7C1pvHqHw5A4vrJfjNwvOdzYu0Gml16OCs2GRiTUUS4=`.

YAML v3.0.4 has exact tag commit
`c3552c15f996075a7634df5159d9161c67bf3d76`; its checksums are:

- module: `h1:tfq32ie2Jv2UxXFdLJdh3jXuOzWiL1fo0bu/FbuKpbc=`;
- go.mod: `h1:DhzuOOF2ATzADvBadXxruRBLzYTpT36CKvDb3+aBEFg=`.

The release delta migrates Cobra's YAML documentation dependency to the new
module path, documents repeated flags, applies const-only refactors, and
updates lint/CI. Independently reverify all primary evidence and direct Cobra
callers before editing.

## Start And Stop

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
links, launcher `--check`, P7/P8 queue, implementation commit/tree, and every
accepted manifest before editing. Read the active archive, this handover, the
P7 roadmap, `go.mod`/`go.sum`, Cobra callers and tests, toolchain/baseline
reproduction, compatibility, snapshot/Docker, quality, and audit contracts.

Stop before any dependency outside the exact Cobra/YAML closure,
behavior/API/CLI change, main-module language-floor change, quality-tool
upgrade, rejected-candidate retry, P8 domain work, inactive packaging,
publication, publisher/registry/credential change, or release. Do not push,
merge, publish, release, delete retained evidence or images, stash, revert,
launch a successor, or remove the worktree.
