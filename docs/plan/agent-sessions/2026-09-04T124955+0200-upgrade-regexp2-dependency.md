# Agent Session: Upgrade Regexp2 Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-04T124955+0200-upgrade-regexp2-dependency`
Created: `2026-09-04T12:49:55+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `03166e3c53f8147223715a66747a1672e1cf23ebb8f4a4fb7e983dd90f848e5e`
Previous: [2026-09-04T105124+0200-upgrade-go-md2man-dependency.md](2026-09-04T105124+0200-upgrade-go-md2man-dependency.md)
Next: [2026-09-04T145138+0200-upgrade-gomarkdown-dependency.md](2026-09-04T145138+0200-upgrade-gomarkdown-dependency.md)
Outcome: Upgraded only selected indirect `github.com/dlclark/regexp2` from v1.8.1 to v1.12.0 in dependency-only commit `3fd6684`; exact selection, Markdown/highlighting, quality, regression, artifact, vulnerability-parity, and clean-tree gates pass, and P7 remains active for the measured gomarkdown latest-pseudoversion group.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its next measured dependency group: upgrade selected indirect
github.com/dlclark/regexp2 v1.8.1 to latest v1.12.0 as an exact one-selection
move. Independently reverify the decision from primary evidence, preserve the
loaded Markdown/highlighting behavior and every quality contract, and finish
with no unrelated drift.

# Authorized Roadmap

P2A-P6 are complete. The Go 1.26.7 toolchain, terminal pair, logrus, Cobra
v1.10.1 closure, pflag, Uniseg, Colorable/go-isatty, Cobra v1.10.2/YAML v3.0.4,
YAML v3.0.5, and go-md2man v2.0.7 groups are complete. Viper v1.16.0, Emoji
v2.2.14, and ini v1.67.3 are rejected because their closures violate the
retained Go 1.18 floor. P7 remains active for dependency groups and P8 remains
queued.

This session may change only the exact regexp2 v1.12.0 pin's required
go.mod/go.sum metadata, a focused executable Markdown/highlighting contract if
one is genuinely needed, and the roadmap record. Do not revisit Viper, Emoji,
or ini, raise the language floor, or combine another dependency group.

The clean offline current-tree probe changes exactly one selection:
github.com/dlclark/regexp2 v1.8.1 -> v1.12.0. Both states retain 234 selected
modules, 3,556 graph edges, and the exact 429-package test population. Exact
go get replaces the existing indirect requirement and adds only the v1.12.0
module/checksum pair while retaining the old v1.8.1 checksum pair. The
historical go mod tidy -diff projection grows from 278 to 282 lines; it retains
the explicit v1.12.0 requirement and projects removal of the old v1.8.1
checksum pair. Record this known metadata consequence rather than using tidy as
the implementation command.

Regexp2 and regexp2/syntax are loaded through plybuild/cmd -> go-term-markdown
-> Chroma -> regexp2. Regexp2 v1.12.0 declares Go 1.13 and the main module
remains Go 1.18. Exact Go 1.26.7 build, complete tests, pinned lint,
byte-identical public help, identical API/CLI reports, artifact meta-contracts,
and 22/33/22 vulnerability parity pass.

Stop and record the decision if independent replay changes any other selection,
checksum, edge count, package population, declared minimum Go version, API/CLI
report, public help, loaded path, or vulnerability population.

# Measurements At Start

The last P7 implementation commit is
55dc69dda20c4d8f96b6dbcdd70cfef76467cc85, exact parent
6fe18fe01326f77bee375d3f8834b44a0416907d, clean tree
321e39594a4bbe0b7c6afda4e9e5a07cf395fe47. After the go-md2man handoff, the
continuity HEAD must have exact parent 55dc69d, and ordinary and ignored status
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

The authoritative go-md2man implementation selection root is
/private/tmp/ply-p7-md2man-selection.6fe18fe.fBmfEU. Its verified 44,219-entry
manifest SHA-256 is
77c1ec0a2fd1b60e4369d87c3c387efc39bcb1297cc6cf698fac6af3b42d6807;
selection-summary.json SHA-256 is
7c5cab15a7a0d813df60b176666fae4b9bd536c2f081566032509da18c1873dd.

The exact accepted go-md2man schema-2 review root is
/private/tmp/ply-p7-md2man-quality-review.55dc69d.BdRjpS. Its verified 22-entry
manifest SHA-256 is
af36ebc39220021cb78b7995f975d68c35b46b52e94dc854d8f5bce17e8fe6c5;
manual evidence SHA-256 is
e8b1e54604b073ea71d0bf2a568f0e69b60cd22ef90a530eb1cfa233d561c68f,
and all six focused rows pass.

The exact accepted go-md2man quality root is
/private/tmp/ply-p7-md2man-quality-parent.55dc69d.7VciIk/quality-gate. Its
verified 235,335-entry manifest SHA-256 is
79765081255f37f2223a4cfa10b37e900cdb2842f3b60589cac985b747ad64ed.
Its Q0-Q2 scorecard SHA-256 is
798ab4be7e9642cd822bb855e5010f6477fa3cd3524c06422da2feef45643c71:
exact make quality passes all 27 rows at L2, 80/80 mutations, 8/8 mutation and
4/4 acceptance populations, with held, regressed, not-comparable, and dirty
counts zero.

The independent go-md2man regression root is
/private/tmp/ply-p7-md2man-regression-gate.55dc69d.bD0Kbx. Its verified
68,729-entry manifest SHA-256 is
bd07a6975b0a495c872bde02c15a3c745e5527768baaf4edded40362af072fa3.
All 40 stages pass. Its regression-summary.json SHA-256 is
e7b852b82883a12d37a1977f4d170dcb16f19520af5d4c05a5215aa9f80cba63;
full scorecard SHA-256 is
ff12e563966050557bc2755e1e60541eebb8cef01cccc66db8cd0ea3b7e7d0b3.
The full audit exits 1 only for queued L3 rows Q3.1, Q3.3, Q3.4, and Q3.7 and
is not a P7 dependency-group exit gate. Vulnerability populations remain
exactly 22/33/22.

The rejected ini decision root remains
/private/tmp/ply-p7-next-selection.33e1873.mRehC2. Its verified 39,374-entry
manifest SHA-256 is
c7af82cbe2327de51dc93fcf48c61e65e6b38e5d51a3a7556c07d6d753016181.
Do not implement it: selected github.com/stretchr/objx v0.5.2 declares Go 1.20.

The authoritative regexp2 selection root is
/private/tmp/ply-p7-next-selection.55dc69d.8kz8Cw. Its verified 49,090-entry
manifest SHA-256 is
a2d03742f02396ede46ac73a19b00191812d84cfcad2b416dcf1f19b460b3a10;
selection-summary.json SHA-256 is
f5ee0c23f405f47fb6ae9627ae9a8faabc6a8112713c3691c9f8bcb95d6927e1.

Regexp2 v1.12.0's observed proxy time is 2026-04-18T22:12:01Z. It declares
Go 1.13, and its lightweight tag resolves to unsigned commit
3d5df45b703801b3fe51eb3f5c0dd302e8b0d676. Its checksum pair is
h1:0j4c5qQmnC6XOWNjP3PIXURXN2gWx76rd3KvgdPkCz8= /
h1:DHkYz0B9wPfa6wondMfaivmHpzrQ3v9q8cnmRbL6yW8=. The 12-commit
v1.8.1...v1.12.0 delta fixes lazy-loop stack corruption and termination,
timeout false positives, runner text retention, and ECMAScript plus Singleline
dot behavior; it adds single-letter Unicode property syntax, timeout test
controls, and text marshaling, plus their regressions and documentation. The Go
vulnerability module index has no regexp2 entry.

# Role And Boundaries

Reverify current regexp2 release metadata, latest stable tag and commit
identity, module Go requirement, checksums, release changes, exact one-selection
MVS closure, loaded package/path evidence, explicit-pin and tidy behavior, and
vulnerability information from primary Go module, repository, Go documentation,
and Go vulnerability sources. Record why the selection is compatible with
pinned Go 1.26.7, the retained Go 1.18 floor, Cobra v1.10.2, YAML v3.0.5,
go-md2man v2.0.7, Blackfriday v2.1.0, golangci-lint 2.12.2, GoReleaser 2.17.1,
pflag v1.0.10, Uniseg v0.4.7, Colorable v0.1.15, and go-isatty v0.0.20.

Do not change production Go behavior, public Go API, CLI semantics, toolchain
or main-module language declarations, Docker or release inputs, any dependency
selection outside exact regexp2 v1.12.0, quality-tool versions, quality
thresholds, baseline numeric debt, compatibility allowlists, acceptance or
mutation populations, publishers, registries, credentials, inactive packaging,
or P8 domain code. Keep Cobra at v1.10.2, YAML at v3.0.5, go-md2man at v2.0.7,
Blackfriday at v2.1.0, and retained gopkg.in/yaml.v2 and gopkg.in/yaml.v3
selections; do not combine a source fix or another dependency group.

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
roadmap, go.mod/go.sum, regexp2 graph and loaded-path evidence, toolchain and
baseline-reproduction contracts, compatibility, snapshot/Docker, quality, and
audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, go mod tidy -diff projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, loaded
regexp2 packages/path, and a vulnerability baseline. Reverify v1.12.0 against
current proxy and repository evidence. Reproduce the retained one-selection
closure exactly and confirm its Go 1.13 declaration preserves the Go 1.18 floor
before editing.

Only if exact and compatible, make one focused dependency-only implementation
commit using exact Go 1.26.7 with
go get github.com/dlclark/regexp2@v1.12.0; do not hand-edit dependency metadata
and do not run tidy as the implementation command. Re-run focused dependency
graph, loaded-path and Markdown/highlighting checks, pinned lint, complete
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
