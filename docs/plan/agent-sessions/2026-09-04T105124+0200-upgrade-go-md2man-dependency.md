# Agent Session: Upgrade Go-md2man Dependency

Status: NEXT
Session ID: `2026-09-04T105124+0200-upgrade-go-md2man-dependency`
Created: `2026-09-04T10:51:24+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `718cca634d724dc53662afb492dee3034079af47accf8a02b09a7313b0400480`
Previous: [2026-09-04T083743+0200-upgrade-yaml-v3-dependency.md](2026-09-04T083743+0200-upgrade-yaml-v3-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its next measured dependency group: upgrade selected indirect
github.com/cpuguy83/go-md2man/v2 v2.0.6 to latest v2.0.7 as an exact
one-selection move. Independently reverify the decision from primary evidence,
preserve behavior and every quality contract, and finish with no unrelated
drift.

# Authorized Roadmap

P2A-P6 are complete. The Go 1.26.7 toolchain, terminal pair, logrus, Cobra
v1.10.1 closure, pflag, Uniseg, Colorable/go-isatty, Cobra v1.10.2/YAML v3.0.4,
and YAML v3.0.5 groups are complete. Viper v1.16.0, Emoji v2.2.14, and ini
v1.67.3 are rejected because their closures violate the retained Go 1.18 floor.
P7 remains active for dependency groups and P8 remains queued.

This session may change only the exact go-md2man v2.0.7 pin's required
go.mod/go.sum metadata, including the unchanged selected Blackfriday v2.1.0
metadata recorded by exact go get, a focused executable contract if one is
genuinely needed, and the roadmap record. Do not revisit Viper, Emoji, or ini,
raise the language floor, or combine another dependency group.

The clean offline current-tree probe changes exactly one selection:
github.com/cpuguy83/go-md2man/v2 v2.0.6 -> v2.0.7. Both states retain 234
selected modules and the exact 429-package population. Exact go get adds
explicit indirect requirements for go-md2man v2.0.7 and its already-selected
github.com/russross/blackfriday/v2 v2.1.0, so the graph gains three edges,
3,553 -> 3,556, without changing Blackfriday's selection. The historical
go mod tidy -diff projection grows from 271 to 278 lines. No loaded main-module
package imports go-md2man or Blackfriday, so projected tidy removes the two
explicit requirements and their newly materialized checksum lines; record this
known metadata consequence rather than using tidy as the implementation
command. Go-md2man declares Go 1.12, Blackfriday has no go directive, and the
main module remains Go 1.18. Exact Go 1.26.7 build, complete tests, pinned lint,
byte-identical public help, identical API/CLI reports, artifact meta-contracts,
and 22/33/22 vulnerability parity pass.

Stop and record the decision if independent replay changes any other selection,
checksum, edge, package, declared minimum Go version, API/CLI report, public
help, or vulnerability population.

# Measurements At Start

The last P7 implementation commit is
cfdcb370e98ba4e4536f0bf08d5b1abbb01f1856, exact parent
fdd9986d55d0aee27679ccc3675710e03f0d1e1a, clean tree
7244802bfe8e229ab128c122c0b6cb1a49fc7729. After the YAML v3.0.5 handoff, the
continuity HEAD must have exact parent cfdcb37, and ordinary and ignored status
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

The accepted YAML v3.0.5 selection root is
/private/tmp/ply-p7-yaml-selection.fdd9986.cf5t7R. Its verified 31,004-entry
manifest SHA-256 is
1274e8d593cc6dd526ef2c857168212204660653eb6045835ae7943cfc94d986;
selection-summary.json SHA-256 is
6dcd8a236217e7b620b7408a4d79270d36003e1364563ead6b875d361950825b.

The exact accepted YAML schema-2 review root is
/private/tmp/ply-p7-yaml-quality-review.cfdcb37.3xOd0m. Its verified 22-entry
manifest SHA-256 is
ea61323dafe61f94c4d60846817e4c9e7918a28bf086ce0bb8ea5cac01045dd0;
manual evidence SHA-256 is
19d5ea397a9db90479596f40db04ca3ab3932b245ccc8334eccb9dd3fdf4e818,
and all six focused rows pass.

The exact accepted YAML quality root is
/private/tmp/ply-p7-yaml-quality-parent.cfdcb37.xxdcLh/quality-gate. Its
verified 235,241-entry manifest SHA-256 is
919b44aa5ec811e2df49dc7784d02c7f70820e52939e8b8545d8347859db0417.
Its Q0-Q2 scorecard SHA-256 is
9f8e6e70705a7c3d27456a78c8cdf7bf23f5d7ca47dbdfabb83371f15c886c95:
exact make quality passes all 27 rows at L2, 80/80 mutations, 8/8 mutation and
4/4 acceptance populations, with held, regressed, not-comparable, and dirty
counts zero.

The independent YAML regression root is
/private/tmp/ply-p7-yaml-regression-gate.cfdcb37.dISqmD. Its verified
59,161-entry manifest SHA-256 is
2f0206fcc142c17858347fd814ffbb4b996d20c3a3b70f9a335aa141c7cc3cdd.
All 40 stages pass. Its regression-summary.json SHA-256 is
8552e0df7ed87a8ad1a327b72e1967a3adb295039126edbf97a5f312e56ecd95;
full scorecard SHA-256 is
1fbccda6eccc11fe9f3628a6dcad906d4d7fab920c7acccedf5c6d2bcf799b68.
The full audit exits 1 only for queued L3 rows Q3.1, Q3.3, Q3.4, and Q3.7 and
is not a P7 dependency-group exit gate. Vulnerability populations remain
exactly 22/33/22.

The rejected ini decision root remains
/private/tmp/ply-p7-next-selection.33e1873.mRehC2. Its verified 39,374-entry
manifest SHA-256 is
c7af82cbe2327de51dc93fcf48c61e65e6b38e5d51a3a7556c07d6d753016181.
Do not implement it: selected github.com/stretchr/objx v0.5.2 declares Go 1.20.

The authoritative go-md2man selection root is
/private/tmp/ply-p7-next-selection.cfdcb37.LkLUsE. Its verified 26,533-entry
manifest SHA-256 is
527a0d1d869845326c4b3b6ba16dc6c6b8b5794728b4e56d50145db88e9e6ebc;
selection-summary.json SHA-256 is
83a016e8807b8737cdf26895f382e67f8efaff4ae2371c7fccd5db5f55e15f9d.

Go-md2man v2.0.7's observed proxy time is 2025-04-24T23:51:24Z. It declares
Go 1.12, its lightweight tag resolves to verified commit
061b6c7cbecd6752049221aa15b7a05160796698, and its checksum pair is
h1:zbFlGlXEAKlwXpmvle3d8Oe3YnkKIK4xSRTd3sHPnBo= /
h1:oOW0eioCTA6cOiMLiUPZOpcVxMig6NIQQ7OS05n1F4g=. Existing selected
Blackfriday v2.1.0 has no go directive; exact go get newly materializes its
module checksum h1:JIOH55/0cWyOuilr9/qlrm0BSXldqnqwMsf35Ld67mk= without changing
its version.

The v2.0.6...v2.0.7 release delta has 17 commits. The material production
change fixes roff table rendering for long final-row cells; other production
edits are package documentation and lint-comment formatting. Tests add the
long-table regression and simplify an unused helper signature; CI, lint, README,
and installation documentation are refreshed. The Go vulnerability module
index has no go-md2man or Blackfriday entry.

# Role And Boundaries

Reverify current go-md2man release metadata, latest stable tag and commit
identity, module Go requirement, checksums, release changes, exact one-selection
MVS closure, loaded-package absence, explicit-pin and tidy behavior, and
vulnerability information from primary Go module, repository, Go documentation,
and Go vulnerability sources. Record why the selection is compatible with
pinned Go 1.26.7, the retained Go 1.18 floor, Cobra v1.10.2, YAML v3.0.5,
golangci-lint 2.12.2, GoReleaser 2.17.1, pflag v1.0.10, Uniseg v0.4.7,
Colorable v0.1.15, and go-isatty v0.0.20.

Do not change production Go behavior, public Go API, CLI semantics, toolchain
or main-module language declarations, Docker or release inputs, any dependency
selection outside exact go-md2man v2.0.7, quality-tool versions, quality
thresholds, baseline numeric debt, compatibility allowlists, acceptance or
mutation populations, publishers, registries, credentials, inactive packaging,
or P8 domain code. Keep Blackfriday selected at v2.1.0, Cobra at v1.10.2, YAML
at v3.0.5, and the retained gopkg.in/yaml.v2 and gopkg.in/yaml.v3 selections;
do not combine a source fix or another dependency group.

Keep module/build caches, graph projections, vulnerability results, reports,
generated artifacts, build contexts, schema-2 evidence, and audit output outside
the worktree. Warm fresh module caches from a separate external Git archive;
never run go mod download all inside a measured old/candidate tree or use a
real-tree download as a tidy surrogate. Never create .agent-task/current.md or
.quality/manual-evidence.json.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher --check, and the P7/P8 checkpoint block before editing.
Verify the accepted manifests. Read the rolling handover, this archive, the P7
roadmap, go.mod/go.sum, go-md2man graph and loaded-package evidence, toolchain
and baseline-reproduction contracts, compatibility, snapshot/Docker, quality,
and audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, go mod tidy -diff projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, and a
vulnerability baseline. Reverify v2.0.7 against current proxy and repository
evidence. Reproduce the retained one-selection closure exactly and confirm its
Go 1.12 declaration preserves the Go 1.18 floor before editing.

Only if exact and compatible, make one focused dependency-only implementation
commit using exact Go 1.26.7 with
go get github.com/cpuguy83/go-md2man/v2@v2.0.7; do not hand-edit dependency
metadata and do not run tidy as the implementation command. Re-run focused
dependency graph and loaded-package checks, pinned lint, complete
tests/race/vet, API/CLI and entry/subprocess compatibility, launcher and Make
contracts, complete preflight, host acceptance, fresh snapshot/Docker meta and
acceptance, audit meta, focused and exact Q0-Q2 audits, the separate full audit,
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
