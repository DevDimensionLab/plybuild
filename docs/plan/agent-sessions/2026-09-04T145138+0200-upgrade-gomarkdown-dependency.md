# Agent Session: Upgrade Gomarkdown Dependency

Status: NEXT
Session ID: `2026-09-04T145138+0200-upgrade-gomarkdown-dependency`
Created: `2026-09-04T14:51:38+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d34836648fac466d13642b8ee20822f4fa743c7daa5c7d004be0749d9172b661`
Previous: [2026-09-04T124955+0200-upgrade-regexp2-dependency.md](2026-09-04T124955+0200-upgrade-regexp2-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its next measured dependency group: upgrade selected indirect
github.com/gomarkdown/markdown v0.0.0-20221013030248-663e2500819c to latest
v0.0.0-20260824154242-13c5cf49db8d as an exact one-selection move.
Independently reverify the decision from primary evidence, preserve loaded
Markdown parsing/rendering behavior and every quality contract, retain the
measured security improvement, and finish with no unrelated drift.

# Authorized Roadmap

P2A-P6 are complete. The Go 1.26.7 toolchain, terminal pair, logrus, Cobra
v1.10.1 closure, pflag, Uniseg, Colorable/go-isatty, Cobra v1.10.2/YAML v3.0.4,
YAML v3.0.5, go-md2man v2.0.7, and regexp2 v1.12.0 groups are complete. Viper
v1.16.0, Emoji v2.2.14, and ini v1.67.3 remain rejected because their closures
violate the retained Go 1.18 floor. Current primary module metadata also rejects
fatih/color v1.19.0 at Go 1.25.0 and fsnotify v1.10.1 at Go 1.23. P7 remains
active for dependency groups and P8 remains queued.

This session may change only the exact gomarkdown latest-pseudoversion pin's
required go.mod/go.sum metadata, a focused executable Markdown contract if one
is genuinely needed, and the roadmap record. Do not revisit rejected groups,
raise the language floor, or combine another dependency group.

The clean offline current-tree probe changes exactly one selection:
github.com/gomarkdown/markdown v0.0.0-20221013030248-663e2500819c ->
v0.0.0-20260824154242-13c5cf49db8d. Both states retain 234 selected modules,
3,556 graph edges, and the exact 429-package test population. Exact go get
replaces the existing indirect requirement and adds only the new module and
go.mod checksum lines while retaining the old selected checksum pair. The
historical go mod tidy -diff projection grows from 282 to 285 lines; it retains
the explicit new requirement and projects removal of the old selected checksum
pair. Record this consequence rather than using tidy as the implementation
command.

Four gomarkdown packages are loaded through
plybuild/cmd -> go-term-markdown -> gomarkdown/markdown. The candidate declares
Go 1.12 and the main module remains Go 1.18. Exact Go 1.26.7 build, complete
main-module tests, gomarkdown self-tests, focused go-term-markdown rendering
tests, pinned lint, byte-identical public help, identical API/CLI reports, and
artifact meta-contracts pass.

The candidate intentionally improves the current vulnerability population:
Darwin symbol / Darwin module / Windows symbol counts move from 22/33/22 to
20/30/20 with no additions. GO-2023-2074 and GO-2024-3205 disappear from both
symbol scans, and those two plus GO-2026-5208 disappear from the module scan.
Stop and record the decision if independent replay changes any other selection,
checksum, edge count, package population, declared minimum Go version, API/CLI
report, public help, loaded path, or expected vulnerability delta.

# Measurements At Start

The last P7 implementation commit is
3fd6684941c57bf004ac7d09b73f5e716c12689b, exact parent
d043241b0f092235edb63353a83daf1f57996bd5, clean tree
8d08fa15c7bd85a229380d8336ca25a3312f19b1. After the regexp2 documentation
handoff, continuity HEAD must have exact parent 3fd6684, and ordinary and
ignored status must be empty.

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

The authoritative regexp2 implementation replay root is
/private/tmp/ply-p7-regexp2-replay.d043241.GjPoeU. Its verified 46,723-entry
manifest SHA-256 is
5b614df2c43b7769643aab1bd3b96a7bcbb959efc3f63ba3df1ac02d00c5f1fa;
selection-summary.json SHA-256 is
6df19722cc55cf8099f99ccd18b3396836ab58c0dc4ab072eb23c8ad71af5024.

The exact accepted regexp2 schema-2 review root is
/private/tmp/ply-p7-regexp2-quality-review.3fd6684.l0sSof. Its verified 24-entry
manifest SHA-256 is
62d0e6c4e10db44fcd646bd3683568cf77711b2ece9fb04dfb6bd7dca68a5a62;
manual evidence SHA-256 is
6a39bf49fe11affd7e75bf09451fd42bb81ec3fdb7a06e2bdc124caa5dc6fa5f,
and all six focused rows pass.

The exact accepted regexp2 quality root is
/private/tmp/ply-p7-regexp2-quality-parent.3fd6684.1BZGVe/quality-gate. Its
verified 235,425-entry manifest SHA-256 is
9bcc902d7eac1085e8f453989faf4b4d432522eb5a97928454177faba2280115.
Its Q0-Q2 scorecard SHA-256 is
0e7139dc25cc2cd125807536e24b18e27f91f48084cd5ddf644381ee292dd2a8:
exact make quality passes all 27 rows at L2, 80/80 mutations, 8/8 mutation and
4/4 acceptance populations, with held, regressed, not-comparable, and dirty
counts zero.

The independent regexp2 regression root is
/private/tmp/ply-p7-regexp2-regression-gate.3fd6684.YQdxua. Its verified
70,973-entry manifest SHA-256 is
2443aa2ba177642c2416f534ab66b83ac1be498080396e75c0ec87e40f3fcce7.
All 40 stages pass. Its regression-summary.json SHA-256 is
dcce4656a650406e9ee4dbcfb31ea8a69991267c917ce8f96b9c7e94170ab8b7;
full scorecard SHA-256 is
c95201c96297a198b7bec9b5a031443e42c0fe4be5f17a4d383362b9dce5f527.
The full audit exits 1 only for queued L3 rows Q3.1, Q3.3, Q3.4, and Q3.7 and
is not a P7 dependency-group exit gate. Current vulnerability populations are
exactly 22/33/22.

The authoritative gomarkdown selection root is
/private/tmp/ply-p7-gomarkdown-selection.3fd6684.VhoTOi. Its verified
46,966-entry manifest SHA-256 is
85c97bd06674c2b1b76cca98de6d203a2442273ae86afbe2f258e66506892e26;
selection-summary.json SHA-256 is
d6451deffdb8b70fc7112e3acdeae44aea1876ef1cb686cfe0ccf7e8e917db87.

The Go proxy exposes no tagged gomarkdown module versions and resolves @latest
to pseudoversion v0.0.0-20260824154242-13c5cf49db8d, observed at
2026-08-24T15:42:42Z. It declares Go 1.12 and resolves to unsigned commit
13c5cf49db8d0bd189d6d04337e618ae01a83c8f. Its checksum pair is
h1:8VtgBGEPLZ2Yn0Fuh6Pwmy3qF6indeaqy8mrBMbUKRQ= /
h1:JDGcbDT52eL4fju3sZ4TeHGsQwhG9nbDV21aMyhwPoA=.

The 119-commit, 87-file delta extends parser/renderer hooks and Markdown/HTML
attribute support; changes newline, link, list, code-block, deep-nesting, and
round-trip behavior; and improves large-document performance. It fixes
out-of-bounds and short-URL panics, unterminated-code and empty-definition-list
loops, quadratic list/link parsing, unsafe heading and code-info attributes,
and the three indexed vulnerabilities, with broad regression, fuzz, example,
and documentation coverage.

# Role And Boundaries

Reverify current gomarkdown module metadata, @latest resolution, absence of
tagged module versions, commit identity, module Go requirement, checksums,
release changes, exact one-selection MVS closure, loaded package/path evidence,
explicit-pin and tidy behavior, and vulnerability information from primary Go
module, repository, Go documentation, and Go vulnerability sources. Record why
the selection is compatible with pinned Go 1.26.7, the retained Go 1.18 floor,
regexp2 v1.12.0, Cobra v1.10.2, YAML v3.0.5, go-md2man v2.0.7, Blackfriday
v2.1.0, golangci-lint 2.12.2, GoReleaser 2.17.1, pflag v1.0.10, Uniseg v0.4.7,
Colorable v0.1.15, and go-isatty v0.0.20.

Do not change production Go behavior, public Go API, CLI semantics, toolchain
or main-module language declarations, Docker or release inputs, any dependency
selection outside the exact gomarkdown pseudoversion, quality-tool versions,
quality thresholds, baseline numeric debt, compatibility allowlists,
acceptance or mutation populations, publishers, registries, credentials,
inactive packaging, or P8 domain code. Keep all retained versions and both
gopkg.in YAML selections fixed; do not combine a source fix or another
dependency group.

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
roadmap, go.mod/go.sum, gomarkdown graph and loaded-path evidence, toolchain and
baseline-reproduction contracts, compatibility, snapshot/Docker, quality, and
audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, go mod tidy -diff projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, loaded
gomarkdown packages/path, and a vulnerability baseline. Reverify the exact
pseudoversion against current proxy and repository evidence. Reproduce the
retained one-selection closure and expected security delta exactly, and confirm
its Go 1.12 declaration preserves the Go 1.18 floor before editing.

Only if exact and compatible, make one focused dependency-only implementation
commit using exact Go 1.26.7 with
go get github.com/gomarkdown/markdown@v0.0.0-20260824154242-13c5cf49db8d;
do not hand-edit dependency metadata and do not run tidy as the implementation
command. Re-run focused dependency graph, loaded-path and Markdown parsing/
rendering checks, pinned lint, complete tests/race/vet, API/CLI and entry/
subprocess compatibility, launcher and Make contracts, complete preflight,
host acceptance, fresh snapshot/Docker meta and acceptance, audit meta,
focused and exact Q0-Q2 audits, the separate full audit, vulnerability
comparison, and empty-HOME count-2. Refresh external schema-2 evidence when
commit binding requires it.

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
