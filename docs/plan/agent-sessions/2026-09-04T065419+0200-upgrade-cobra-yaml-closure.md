# Agent Session: Upgrade Cobra/YAML Dependency Closure

Status: NEXT
Session ID: `2026-09-04T065419+0200-upgrade-cobra-yaml-closure`
Created: `2026-09-04T06:54:19+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `58889b034aab39e17d13957317c0e56de0f9398b3860769ec323f3cdbb73ee9d`
Previous: [2026-09-04T045040+0200-upgrade-colorable-closure.md](2026-09-04T045040+0200-upgrade-colorable-closure.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its next measured dependency group: upgrade direct
github.com/spf13/cobra v1.10.1 to latest v1.10.2 with exact newly selected MVS
companion go.yaml.in/yaml/v3 v3.0.4 as a two-selection move. Independently
reverify the decision from primary evidence, preserve behavior and every
quality contract, and finish with no unrelated drift.

# Authorized Roadmap

P2A-P6 are complete. The Go 1.26.7 toolchain, terminal pair, logrus, Cobra
v1.10.1 closure, pflag, Uniseg, and Colorable/go-isatty groups are complete.
Viper v1.16.0, Emoji v2.2.14, and ini v1.67.3 are rejected because their
closures violate the retained Go 1.18 floor. P7 remains active for dependency
groups and P8 remains queued.

This session may change only the exact Cobra/YAML closure's required
go.mod/go.sum metadata, a focused executable contract if one is genuinely
needed, and the roadmap record. Do not revisit Viper, Emoji, or ini, raise the
language floor, or combine another dependency group.

The clean offline current-tree probe changes exactly two selections:
github.com/spf13/cobra v1.10.1 -> v1.10.2 and newly selects
go.yaml.in/yaml/v3 v3.0.4. It moves 233 -> 234 selected modules and 3,551 ->
3,552 graph edges while retaining the exact 429-package population. The
historical go mod tidy -diff projection grows from 252 to 254 lines. Cobra
declares Go 1.15 and YAML declares Go 1.16; the main module remains Go 1.18.
Exact Go 1.26.7 build, complete tests, pinned lint, byte-identical public help,
identical API/CLI reports, and 22/33/22 vulnerability parity pass.

Stop and record the decision if independent replay changes any other
selection, checksum, edge, package, declared minimum Go version, API/CLI
report, public help, or vulnerability population.

# Measurements At Start

The last P7 implementation commit is
33e187317c6c1be79ef8c5b64ddb2a0f8caec71f, exact parent
e1fd64e4a7bebd6dabf529675b1560172f12f1d5, clean tree
6942bd34879b9935382b44ffab682606cf5ef5c9. After the Colorable handoff, the
continuity HEAD must have exact parent 33e1873, and ordinary and ignored status
must be empty.

The declared and verified toolchain remains Go 1.26.7. The retained official
executable is
/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go, SHA-256
9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6.
Put its directory first in PATH as well as passing exact GO; audit reproduction
invokes literal go. Keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and do not
inject ambient GOFLAGS.

Retain golangci-lint 2.12.2 at
/private/tmp/ply-http-gate.gzWxRY/bin/golangci-lint, SHA-256
3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f;
GoReleaser 2.17.1 at
/private/tmp/ply-snapshot-probe.CiDxj0/tool/goreleaser, SHA-256
f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c;
apidiff at /private/tmp/ply-http-gate.gzWxRY/bin/apidiff, SHA-256
0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20;
and govulncheck v1.7.0 at
/private/tmp/ply-p7-terminal-selection.ctFhJk/tools/bin/govulncheck, SHA-256
0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5.

The accepted Colorable selection root is
/private/tmp/ply-p7-colorable-selection.e1fd64e.QlohU1. Its verified
49,635-entry manifest SHA-256 is
b8ec3adb51fc916820b300f0faf4ff0eddc7e1ae51a9f839b96ceed9c57cebb7;
selection-summary.json SHA-256 is
484e11cb3c6929dfd0a50c050a49008df4ed3908133ae3d87a3a3839f783e07c.

The exact accepted Colorable schema-2 review root is
/private/tmp/ply-p7-colorable-quality-review.33e1873.eGDQR5. Its verified
21-entry manifest SHA-256 is
f6294878de4b3f19a8731b0b94485f8a203ab20137e37fd3e866f27a0afde24b;
manual evidence SHA-256 is
cfbe47aa455636613ec5f03c4ede4a12a4169e2a0da9483cf6fb0519e9a3ce66,
and all six focused rows pass.

The exact accepted Colorable quality root is
/private/tmp/ply-p7-colorable-quality-parent.33e1873.31b10h/quality-gate-retry3.
Its verified 235,059-entry manifest SHA-256 is
6695ee2f2a181c46fbcf953518c2fd9581902f3f3be28bd761602eecab0785ae.
Its Q0-Q2 scorecard SHA-256 is
55cf7f1a4867f54bba31108f2d7ef104f1b0901b44a5ac29aea11e50aec007e6:
exact make quality passes all 27 rows at L2, 80/80 mutations, 8/8 mutation and
4/4 acceptance populations, with held, regressed, not-comparable, and dirty
counts zero.

The independent Colorable regression root is
/private/tmp/ply-p7-colorable-regression-gate.33e1873.1Yb6le. Its verified
98,451-entry manifest SHA-256 is
cdf2963affe5e2395ab8c7911efeaa525e4f39e700ede3945e484cfacb63802c.
All 40 stages pass. Its full scorecard SHA-256 is
c6e7a71ac8cf187654a83f043a3967d53f8e51035c69bf717909117cd7b8ff1b;
it exits 1 only for queued L3 rows Q3.1, Q3.3, Q3.4, and Q3.7 and is not a P7
dependency-group exit gate. Vulnerability populations remain exactly 22/33/22.

The rejected ini decision root is
/private/tmp/ply-p7-next-selection.33e1873.mRehC2. Its verified 39,374-entry
manifest SHA-256 is
c7af82cbe2327de51dc93fcf48c61e65e6b38e5d51a3a7556c07d6d753016181;
decision-summary.json SHA-256 is
04081e0b5f0b3324e49ca321f5edf0a2be5b9b938bce31f68c30aabb89c8a90d.
Do not implement it: selected github.com/stretchr/objx v0.5.2 declares Go 1.20.

The authoritative Cobra/YAML selection root is
/private/tmp/ply-p7-next-cobra-selection.33e1873.ZPppZP. Its verified
43,788-entry manifest SHA-256 is
d522c38eafcd59bd8b172233f06d79a6f32f2410f959e9474c0b8bd0fb8c7a38;
selection-summary.json SHA-256 is
3bf4208b4b8d3fd1d300976274dc6e26f26322aec365554bc7c7fd0e4a377b85.

Cobra v1.10.2's observed proxy time is 2025-12-03T23:51:15Z; it declares Go
1.15, its lightweight tag resolves to commit
88b30ab89da2d0d0abb153818746c5a2d30eccec, and its checksum pair is
h1:DMTTonx5m65Ic0GOoRY2c16WCbHxOOw6xxezuLaBpcU= /
h1:7C1pvHqHw5A4vrJfjNwvOdzYu0Gml16OCs2GRiTUUS4=.

Go.yaml.in/yaml/v3 v3.0.4 declares Go 1.16, its lightweight tag resolves to
commit c3552c15f996075a7634df5159d9161c67bf3d76, and its checksum pair is
h1:tfq32ie2Jv2UxXFdLJdh3jXuOzWiL1fo0bu/FbuKpbc= /
h1:DhzuOOF2ATzADvBadXxruRBLzYTpT36CKvDb3+aBEFg=.

# Role And Boundaries

Reverify current Cobra and YAML release metadata, latest stable tags and commit
identities, module Go requirements, checksums, release changes, exact
two-selection MVS closure, relevant direct callers and behavior, and
vulnerability information from primary Go module, repository, Go
documentation, and Go vulnerability sources. Record why the closure is
compatible with pinned Go 1.26.7, the retained Go 1.18 floor, golangci-lint
2.12.2, GoReleaser 2.17.1, pflag v1.0.10, Uniseg v0.4.7, Colorable v0.1.15,
and go-isatty v0.0.20.

Do not change production Go behavior, public Go API, CLI semantics, toolchain
or main-module language declarations, Docker or release inputs, any dependency
outside exact Cobra v1.10.2/go.yaml.in/yaml/v3 v3.0.4, quality-tool versions,
quality thresholds, baseline numeric debt, compatibility allowlists,
acceptance or mutation populations, publishers, registries, credentials,
inactive packaging, or P8 domain code. Do not change the retained
gopkg.in/yaml.v2 or gopkg.in/yaml.v3 selections and do not combine a source fix
or another dependency group.

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
roadmap, go.mod/go.sum, Cobra callers and tests, toolchain and
baseline-reproduction contracts, compatibility, snapshot/Docker, quality, and
audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, go mod tidy -diff projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, and a
vulnerability baseline. Reverify v1.10.2 and v3.0.4 against current proxy and
repository evidence. Reproduce the retained two-selection closure exactly and
confirm both declarations preserve the Go 1.18 floor before editing.

Only if exact and compatible, make one focused dependency-only implementation
commit using exact Go 1.26.7 with
go get github.com/spf13/cobra@v1.10.2; do not hand-edit dependency metadata.
Re-run focused Cobra callers and tests, pinned lint, complete tests/race/vet,
API/CLI and entry/subprocess compatibility, launcher and Make contracts,
complete preflight, host acceptance, fresh snapshot/Docker meta and acceptance,
audit meta, focused and exact Q0-Q2 audits, the separate full audit,
vulnerability comparison, and empty-HOME count-2. Refresh external schema-2
evidence when commit binding requires it.

Require exact make quality exit 0 at L2 with all 27 Q0-Q2 rows present, 8/8
mutation and 4/4 acceptance populations, and zero held, regressed,
not-comparable, or dirty counts. The full report may exit 1 only for the four
queued L3 rows, never 2. Retain immutable external evidence and verified
manifests for graph selection, quality, regression, and artifacts.

# Automatic Handoff

After this group succeeds, keep P7 active unless the measured dependency queue
is exhausted and keep P8 queued. Rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
evidence-selected compatible dependency group, replace only the launcher's
mutable regions, run launcher and handoff contracts, and make the normal
docs: prepare next agent session commit. Do not implement the next group in
this session.

Do not launch a successor, push, merge, publish, release, delete retained
evidence or local images, stash, revert, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
