# Agent Session: Upgrade Uniseg Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-04T021635+0200-upgrade-uniseg-dependency`
Created: `2026-09-04T02:16:35+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `53f8d15f3a5a83cd96d030ece393ffb7d8fda50d5b62fa3b68f9d2cb96369472`
Previous: [2026-09-04T013713+0200-upgrade-viper-closure.md](2026-09-04T013713+0200-upgrade-viper-closure.md)
Next: [2026-09-04T045040+0200-upgrade-colorable-closure.md](2026-09-04T045040+0200-upgrade-colorable-closure.md)
Outcome: Uniseg v0.4.7 was independently confirmed as the latest stable release, exact one-selection move, Go 1.18-compatible dependency, and implemented at `a5ff9cf16d5daf2ed6a7578d23a859cfd73b5df3` with only `go.mod`/`go.sum` changes. Exact build, complete tests/race/vet, pinned lint, public help, API/CLI, launcher/Make, host/snapshot/Docker, audit, vulnerability-parity, and empty-HOME contracts pass; exact `make quality` passes all 27 Q0-Q2 rows at L2 with 80/80 killed mutations and 4/4 acceptance flows. The sealed selection, quality, and regression manifests contain 48,985, 234,985, and 98,323 entries. The next clean bounded replay selects Colorable v0.1.15 plus go-isatty v0.0.20 as the smallest compatible two-selection closure; Emoji v2.2.14 was rejected because it raises the main module directive to Go 1.21.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its next measured dependency group: upgrade existing indirect
github.com/rivo/uniseg v0.4.4 to latest v0.4.7 as an exact one-selection move.
Independently reverify the decision from primary evidence, preserve behavior and
every quality contract, and finish with no unrelated drift.

# Authorized Roadmap

P2A-P6 are complete. The Go 1.26.7 toolchain, terminal pair, logrus group,
Cobra closure, and pflag patch are complete. Viper v1.16.0 was rejected before
implementation because its exact closure violates the retained Go 1.18 floor.
P7 remains active for dependency groups and P8 remains queued.

This session may change only Uniseg's required go.mod/go.sum metadata, a focused
executable contract if one is genuinely needed, and the roadmap record. Do not
revisit Viper or raise the language floor in this group.

The current-tree external probe selects only Uniseg v0.4.7: both states retain
233 selected modules, 3,551 graph edges, and the exact 429-package test
population. The historical go mod tidy -diff projection grows from 240 to 244
lines. V0.4.7 declares Go 1.18 and has checksum pair
h1:WUdvkW8uEhrYfLC4ZzdpI2ztxP1I582+49Oc5Mq64VQ= /
h1:FN3SvrM+Zdj16jyLfmOkMNblXMcoc8DfTHruCPUcx88=. Exact Go 1.26.7 build,
complete tests, pinned lint, and byte-identical public help pass.

This was the smallest compatible candidate in the bounded probe. Emoji v2.2.14
also changes one selection but declares Go 1.21; colorable v0.1.15 changes two
selections; and gopkg.in/ini.v1 v1.67.3 changes three. Stop and record the
decision if independent replay changes any other selection, checksum, edge,
package, or declared minimum Go version.

# Measurements At Start

The last P7 implementation commit remains
360b2f3c792ca131d259f40852619f2840cefdf1, exact parent
4f0cc642e9c45960b641133c32a0ffc6b1b3a94b, clean tree
a1b941041f5ae685613dd44f4f2db20092cd472c. The Viper decision handoff is
27dbc7819afdcf37c4c52cdf3cde2b5b042ba89d, exact parent 360b2f3.
After this handoff, the continuity HEAD must have exact parent 27dbc78, and
ordinary and ignored status must be empty.

The declared and verified toolchain remains Go 1.26.7. The retained official
executable is
/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go, SHA-256
9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6.
Put its directory first in PATH as well as passing exact GO; audit reproduction
invokes literal go. Keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and do not
inject ambient GOFLAGS.

The rejected Viper decision and Uniseg selection evidence is retained at
/private/tmp/ply-p7-viper-decision.360b2f3.Wdvj52. Its verified 69,347-entry
manifest SHA-256 is
0b1e054f4fcfae3a630e5192603cb02f53d2bdd3b7623c3fe14bdb3b74eb757f;
decision-summary.json SHA-256 is
cd8cc895ca5979514e5e584b023ac70d9967cf4951ecf81ddb760267e10cbd7d.

The Viper replay matched exactly 31 changed selections, 233 -> 234 modules,
3,551 -> 3,561 edges, 429 -> 429 packages, and 240 -> 270 tidy lines. It was
rejected because eleven targets declare Go 1.19 and Testify v1.8.3 declares Go
1.20. Its technical build/test/lint/help/API/CLI/artifact checks pass and its
Darwin-symbol, Darwin-module, and Windows-symbol vulnerability ID populations
remain exactly 22/33/22. Do not implement or re-probe Viper in this session.

The independently replayed pflag selection root is
/private/tmp/ply-p7-pflag-selection.4f0cc64.58RgL0. Read its incident note
before relying on the original pre-cache manifest. Its authoritative final
44,617-entry manifest SHA-256 is
b64cad473c92a5551e997010f88243d46fafa2ddc65c11db5d2d6bb90786298c.

The exact accepted pflag quality root is
/private/tmp/ply-p7-pflag-quality-gate.360b2f3.1788472549.40415, with verified
234,799-entry manifest SHA-256
47971b2e6cf9399f085e233d9a264f9d22c2e5f45dd024375fbb8e151cbd36c6.
Its Q0-Q2 scorecard SHA-256 is
f50ee942007f4cf5b76b0d87f610469f9d95dc04c6363779e2d6e83a7c632934:
all 27 criteria pass at L2, populations are 8/8 and 4/4, and held, regressed,
not-comparable, and dirty counts are zero.

The independent pflag regression root is
/private/tmp/ply-p7-pflag-regression-gate.360b2f3.gb4Rvk, with verified
92,211-entry manifest SHA-256
ec5eb0c20975138fd96b4de3c10f855c0ceb11a5f1152039be2586025beec857.
All 40 stages pass. Its full report exits 1 only for queued L3 rows Q3.1, Q3.3,
Q3.4, and Q3.7; it is not a P7 dependency-group exit gate.

# Role And Boundaries

Reverify current Uniseg release metadata, latest stable tag and commit identity,
module go requirement, checksums, release changes, exact one-module MVS
selection, relevant callers and behavior, and vulnerability information from
primary Go module, repository, Go documentation, and Go vulnerability sources.
Record why the candidate is compatible with pinned Go 1.26.7, the retained Go
1.18 floor, golangci-lint 2.12.2, GoReleaser 2.17.1, Cobra v1.10.1, and pflag
v1.0.10.

Do not change production Go behavior, public Go API, CLI semantics, toolchain or
main-module language declarations, Docker or release inputs, any dependency
outside exact Uniseg v0.4.7, quality-tool versions, quality thresholds, baseline
numeric debt, compatibility allowlists, acceptance or mutation populations,
publishers, registries, credentials, inactive packaging, or P8 domain code. Do
not combine a source fix or another dependency group with this move.

Keep module/build caches, graph projections, vulnerability results, reports,
generated artifacts, build contexts, schema-2 evidence, and audit output
outside the worktree. Warm fresh module caches from a separate external Git
archive; never run go mod download all inside a measured old/candidate tree or
use a real-tree download as a tidy surrogate. Never create
.agent-task/current.md or .quality/manual-evidence.json.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher --check, and the P7/P8 checkpoint block before editing.
Verify the accepted manifests. Read the rolling handover, this archive, the P7
roadmap, go.mod/go.sum, Uniseg callers and tests, toolchain and
baseline-reproduction contracts, compatibility, snapshot/Docker, quality, and
audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, go mod tidy -diff projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, and a
vulnerability baseline. Reverify v0.4.7 against current proxy and repository
evidence. Reproduce the retained one-selection closure exactly and confirm that
its declarations preserve the Go 1.18 floor before editing.

Only if exact and compatible, make one focused dependency-only implementation
commit using exact Go 1.26.7 with go get github.com/rivo/uniseg@v0.4.7; do not
hand-edit dependency metadata. Re-run focused callers and tests, pinned lint,
complete tests/race/vet, API/CLI and entry/subprocess compatibility, launcher
and Make contracts, complete preflight, host acceptance, fresh snapshot/Docker
meta and acceptance, audit meta, focused and exact Q0-Q2 audits, the separate
full audit, vulnerability comparison, and empty-HOME count-2. Refresh external
schema-2 evidence when commit binding requires it.

Require exact make quality exit 0 at L2 with all 27 Q0-Q2 rows present, 8/8
mutation and 4/4 acceptance populations, and zero held, regressed,
not-comparable, or dirty counts. The full report may exit 1 only for the four
queued L3 rows, never 2. Retain immutable external evidence and verified
manifests for graph selection, quality, regression, and artifacts.

# Automatic Handoff

After this group succeeds, keep P7 active unless the measured dependency queue
is exhausted and keep P8 queued. Rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
evidence-selected dependency group, replace only the launcher's mutable
regions, run launcher and handoff contracts, and make the normal
docs: prepare next agent session commit. Do not implement the next group in
this session.

Do not launch a successor, push, merge, publish, release, delete retained
evidence or local images, stash, revert, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
