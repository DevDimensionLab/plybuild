# Agent Session: Upgrade Viper Dependency Closure

Status: ANSWERED - HISTORY
Session ID: `2026-09-04T013713+0200-upgrade-viper-closure`
Created: `2026-09-04T01:37:13+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c4ea0273672360b60ab76be524e2a72d5352ce038964fd4fd8dde8d2d5993198`
Previous: [2026-09-03T233249+0200-upgrade-pflag-dependency.md](2026-09-03T233249+0200-upgrade-pflag-dependency.md)
Next: [2026-09-04T021635+0200-upgrade-uniseg-dependency.md](2026-09-04T021635+0200-upgrade-uniseg-dependency.md)
Outcome: Viper v1.16.0 was rejected before implementation because exact replay found 12 changed selections declaring a minimum Go version above the retained Go 1.18 floor: eleven declare Go 1.19 and Testify v1.8.3 declares Go 1.20. The replay otherwise matched all 31 selections and passed build, tests, lint, help, API/CLI, artifact-meta, module-verification, and 22/33/22 vulnerability-parity checks. No dependency metadata or production behavior changed; the sealed 69,347-entry decision root selects Uniseg v0.4.7 as the next measured one-selection group.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its next measured dependency group: evaluate and, only after a
truthful compatibility decision, upgrade direct
`github.com/spf13/viper v1.15.0` to `v1.16.0` with its exact 31-selection MVS
closure. Preserve
behavior and every quality contract, and finish with no unrelated drift.

# Authorized Roadmap

P2A-P6 are complete. The Go 1.26.7 toolchain, terminal pair, logrus group,
Cobra closure, and pflag patch are complete; P7 remains active for dependency
groups and P8 remains queued. This session may change only the required
Viper-closure `go.mod`/`go.sum` metadata, a focused executable contract if one
is needed, and the roadmap record.

The current-tree external probe selected v1.16.0 as the smallest next stable
Viper release. It changes exactly 31 selected modules, adds only
`github.com/google/s2a-go v0.1.3`, grows the selection from 233 to 234 modules
and the graph from 3,551 to 3,561 edges, and preserves the exact 429-package
test population. The historical `go mod tidy -diff` projection grows from 240
to 270 lines. Do not silently widen this group; stop and record the decision if
independent replay changes any other selection, checksum, edge, or package.

Viper v1.16.0 declares Go 1.17, but its closure selects
`github.com/stretchr/testify v1.8.3`, whose module declares Go 1.20. Exact Go
1.26.7 builds the candidate while the main module retains `go 1.18`. Before
implementation, independently decide whether this transitive declaration is
compatible with the project's stated Go 1.18 language floor. If it is not,
stop and hand off the evidence without raising the floor or implementing a
partial closure.

Latest Viper v1.21.0 declares Go 1.23.0. The retained wider probe found that
already v1.17.0 expands selection to 245 modules and the package population to
437. Newer Viper releases and a toolchain/language-floor move are outside this
group.

# Measurements At Start

The pflag implementation commit is
`360b2f3c792ca131d259f40852619f2840cefdf1`, exact parent
`4f0cc642e9c45960b641133c32a0ffc6b1b3a94b`, clean tree
`a1b941041f5ae685613dd44f4f2db20092cd472c`. After the handoff, the continuity
HEAD must have exact parent `360b2f3`, and ordinary and ignored status must be
empty.

The declared and verified toolchain remains Go 1.26.7. The retained official
executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; audit
reproduction invokes literal `go`. Keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, and do not inject ambient `GOFLAGS`.

The independently replayed pflag selection root is
`/private/tmp/ply-p7-pflag-selection.4f0cc64.58RgL0`. Its final verified
44,617-entry manifest SHA-256 is
`b64cad473c92a5551e997010f88243d46fafa2ddc65c11db5d2d6bb90786298c`;
read its retained incident note before relying on the original pre-cache
manifest.

The exact pflag quality root is
`/private/tmp/ply-p7-pflag-quality-gate.360b2f3.1788472549.40415`, with verified
234,799-entry manifest SHA-256
`47971b2e6cf9399f085e233d9a264f9d22c2e5f45dd024375fbb8e151cbd36c6`.
Its Q0-Q2 scorecard SHA-256 is
`f50ee942007f4cf5b76b0d87f610469f9d95dc04c6363779e2d6e83a7c632934`:
all 27 criteria pass at L2, populations are 8/8 and 4/4, and held, regressed,
not-comparable, and dirty counts are zero.

The independent pflag regression root is
`/private/tmp/ply-p7-pflag-regression-gate.360b2f3.gb4Rvk`, with verified
92,211-entry manifest SHA-256
`ec5eb0c20975138fd96b4de3c10f855c0ceb11a5f1152039be2586025beec857`.
All 40 stages pass. Its full report exits 1 only for queued L3 rows Q3.1,
Q3.3, Q3.4, and Q3.7; it is not a P7 dependency-group exit gate.

The current-tree Viper selection root is
`/private/tmp/ply-p7-next-viper-selection.360b2f3-final.oGxITF`, with verified
40,078-entry manifest SHA-256
`1390745861503e9f57169fff6054db637e14aa6f7ada9b708c9b9155879253a4`
and selection-summary SHA-256
`d3b2704feeab8883cbc67ab98d5aebdb880ee601fde8157aba4fc10ca1d340f5`.
The candidate build, complete tests, pinned lint, public help, exact package
population, and 22/33/22 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations pass without drift.

# Role And Boundaries

Verify time-sensitive Viper release metadata, tag identity, module `go`
requirement, checksums, release changes, the exact 31-module closure, every
relevant closure `go` declaration, and vulnerability information from primary
Go module, repository, Go documentation, and Go vulnerability sources. Record
why the decision is or is not compatible with pinned Go 1.26.7,
golangci-lint 2.12.2, GoReleaser 2.17.1, Cobra v1.10.1, pflag v1.0.10, and the
retained `go 1.18` language floor.

Do not change production Go behavior, public Go API, CLI semantics, toolchain
or main-module language declarations, Docker or release inputs, any dependency
outside the exact v1.16.0 MVS closure, quality-tool versions, quality
thresholds, baseline numeric debt, compatibility allowlists, acceptance or
mutation populations, publishers, registries, credentials, inactive
packaging, or P8 domain code. Do not combine a source fix or another dependency
group with this move; stop on an incompatible behavior or floor decision.

Keep module/build caches, graph projections, vulnerability results, reports,
generated artifacts, build contexts, schema-2 evidence, and audit output
outside the worktree. Warm fresh module caches from a separate external Git
archive; never run `go mod download all` inside a measured old/candidate tree
or use a real-tree download as a tidy surrogate. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint block before
editing. Verify the accepted manifests. Read the rolling handover, this
archive, the P7 roadmap, `go.mod`/`go.sum`, Viper/config callers and tests,
toolchain and baseline-reproduction contracts, compatibility, snapshot/Docker,
quality, and audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, `go mod tidy -diff` projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, and a
vulnerability baseline. Reverify v1.16.0 and every closure change against
current proxy and repository evidence. Reproduce the retained 31-selection
closure exactly and decide the Go 1.18-floor question before editing.

Only if compatible, make one focused dependency-only implementation commit
using the exact Go tool, without hand-editing dependency metadata. Re-run
focused Viper/config callers and tests, pinned lint, complete tests/race/vet,
API/CLI and entry/subprocess compatibility, launcher and Make contracts,
complete preflight, host acceptance, fresh snapshot/Docker meta and
acceptance, audit meta, focused and exact Q0-Q2 audits, the separate full
audit, vulnerability comparison, and empty-HOME count-2. Refresh external
schema-2 evidence when commit binding requires it.

Require exact `make quality` exit 0 at L2 with all 27 Q0-Q2 rows present, 8/8
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
`docs: prepare next agent session` commit. Do not implement the next group in
this session.

Do not launch a successor, push, merge, publish, release, delete retained
evidence or local images, stash, revert, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
