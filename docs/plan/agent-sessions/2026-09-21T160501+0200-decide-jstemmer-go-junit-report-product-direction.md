# Agent Session: Decide Jstemmer Go JUnit Report Product Direction

Status: NEXT
Session ID: `2026-09-21T160501+0200-decide-jstemmer-go-junit-report-product-direction`
Created: `2026-09-21T16:05:01+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5736691c212e4c93b65e973ff505e92384a56c18bcb21df0a7ffd201c38ac915`
Previous: [2026-09-21T143249+0200-evaluate-jstemmer-go-junit-report-dependency.md](2026-09-21T143249+0200-evaluate-jstemmer-go-junit-report-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one bounded product decision for exact-path
`github.com/jstemmer/go-junit-report`. The completed independent evaluation
found no stable exact-path release that preserves Go 1.18 and passes every
applicable ordinary output/resource-lifecycle contract. Choose exactly one
direction below, record its ownership and expiry guards, and stop. Do not
repeat the writer-error fixture, silently accept selected v0.9.1, implement a
dependency or parent change, evaluate another dependency group, reopen any
earlier decision, or begin P8.

# Defensive Scope

This is a documentation-only product decision for an ordinary dependency-
quality review. Reuse the completed public metadata, static source/API,
admissible upstream-test, bounded deterministic fixture, project-graph, and
advisory evidence. Do not fuzz, stress, probe resource exhaustion, generate
oversized or deeply nested input, or perform security or exploitability
analysis. The recorded ordinary destination-writer failure is final; do not
reproduce it.

Every disposable archive, cache, tool, report, and fixture must remain beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`, `/tmp`, a
sibling of the managed root, or another external root. Verify containment and
remove task-owned scratch evidence before handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 is blocked only on this go-junit-report decision after
exact Go 1.26.7, every accepted dependency move through Google UUID v1.4.0,
qualified go-cleanhttp v0.5.2, and all target-specific retained-module
decisions through exact inherited unloaded json-iterator v1.1.12. Every
earlier outcome is final under its own guards. P8 remains queued.

The evaluation left product source, `go.mod`, and `go.sum` unchanged. Selected
go-junit-report remains exact stable v0.9.1, inherited through 25 historical
Cloud Go/storage requests. It has negative `go mod why -m`, zero repository
imports, zero production or complete-test package loads, and no runtime
reachability. Physical MVS selection and zero target loading are not
qualification or risk acceptance. No json-iterator, clockwork, demangle,
strcase, memberlist, or other exception transfers.

# Measurements At Start

The go-junit-report evaluation began from clean ordinary and ignored state on
branch `codex/upgrade-quality` at handoff HEAD
`5946d0958a1538bdea281bb66c59647d0e1a2824`, parent
`19c500237af16f658eb99951e56e042e0d3ea4ee`, tree
`dd70c74d68e652d0c5e8f1ab80d777c1809a28d6`. That handoff changed exactly the
launcher, answered json-iterator decision archive, then-NEXT go-junit-report
evaluation archive, rolling handover, and roadmap. Verify the new handoff,
reciprocal archive chain, latest Google UUID implementation ancestry, exact Go
identity, and launcher check rather than assuming these facts.

The unchanged project has 234 modules, 3,599 graph edges, 355 production and
429 complete-test entries, 197 module-backed complete-test entries across 41
loaded modules, 1,067 sum lines, and the recorded 432-line tidy projection.
`go.mod` and `go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
All 28 earlier guarded selections and their 188 incoming graph edges retain
sorted snapshot SHA-256
`ca6a8b9d4a0d4aa5c6c436d9edeb27cf6d603e36e29edccef9bc2d9d1b81622e`;
all guarded why results are negative and guarded imports/loads are zero.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

# Completed Evaluation

Fresh `go-import` metadata resolves the exact path without redirect to the
public, unarchived, non-fork MIT repository `jstemmer/go-junit-report`. The
exact-path proxy exposes exactly stable v0.9.0, v0.9.1, and v1.0.0. They form
one ancestry, have annotated tags, no retraction, deprecation, replacement, or
alternate exact path, and their proxy/Git regular files and sumdb identities
agree. V0 tags are unsigned; GitHub verifies the v1.0.0 tag signature.

Latest stable v1.0.0 is commit
`16c7efad77dbe9cb51f1a8c896c652c8ecdaa20e`, tree
`2cd81297935f0f139b09b32a7f47f60429208986`. The unreleased v1 branch is two
test/workflow-only commits later and retains the production defect. Current
master declares the different valid `/v2` module path; stable `/v2` releases
and its beta cannot be promoted as the exact-path target.

Each stable exact-path release has a one-module, standard-library-only minimal
production/test closure that resolves and passes count-one, count-ten, race,
vet, build, and five-platform test cross-compilation under exact Go 1.26.7 and
contained Go 1.18.10. All three share the same 121-line exported parser and
formatter API. Bounded command fixtures pass valid/failed/empty report XML,
determinism, header and version options, stdin/stdout, help, positional-error,
and default/opt-in exit-code behavior. Parser/formatter fixtures pass ordinary
report, coverage, benchmark, reader-error, XML, determinism, and concurrent
render behavior. V1.0.0 adds only the `-version` flag to the command surface.

No stable release qualifies. `formatter.JUnitReportXML` wraps the supplied
writer with `bufio.Writer` but ignores errors from its writes and `Flush`, then
returns nil. A small deterministic writer that returns an ordinary write error
therefore yields nil under both SDKs for v0.9.0, v0.9.1, and v1.0.0. This
violates the exported function's error and output resource-lifecycle contract.
The unreleased v1 branch has no production change and retains the defect.

The trees have no build tags, platform branches, cgo, generated source, embed,
symlink, actual example/benchmark/fuzz function, network, subprocess, or
production filesystem boundary. Parser results own their parsed strings,
slices, and maps; rendering allocates its XML and buffered writer; ordering is
deterministic from caller input. There is no mutable package-global runtime
state beyond immutable regex/flag/version values and no internal closeable
resource. Callers retain reader/writer and command-stdio ownership.

MVS selects v0.9.1 through 25 historical `cloud.google.com/go` and storage
requests. A disposable direct v0.9.1 root leaves selection unchanged and gives
234 modules, 3,600 edges, and 1,068 sum lines. A direct v1.0.0 root changes
only the target selection and gives 234/3,600/1,069; target loading remains
zero, all earlier guards remain exact, compatibility/CLI/project gates pass,
and all four host acceptance scripts pass. Tidy removes either manufactured
root, restores v0.9.1, and returns the common baseline 52-line/948-line hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
No genuine tidy-stable owner was identified and no projection was applied.

Fresh stable/v1-branch OSV, GitHub global/repository, Go-index, and pinned
govulncheck v1.8.0 isolated module/package/symbol/test-symbol evidence is
empty. Base and v1.0.0 project scans have byte-identical normalized streams,
identical 30/22/20/20 populations, and no target trace. The 1,402-record Go
index remains at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
Guard OSV retains only the recorded Gorilla/go-retryablehttp pairs, and the
2,807-byte PUBLISHED memberlist CNA response remains exact at
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

# Required Product Decision

Choose exactly one. Option 1 is recommended because the target is completely
unloaded, no stable exact-path release qualifies, and retaining the exact
inherited selection preserves every genuine request and earlier guard:

1. Explicitly retain exact selected, inherited, unloaded v0.9.1 without
   metadata changes under a go-junit-report-specific, non-transferable
   exception. Accept only the completed writer-error, repository/archive,
   API/command, behavior, ownership/lifecycle, Go-floor, MVS, loading,
   repeatability, vulnerability, and related findings. Guard exact v0.9.1,
   all 25 requests and genuine parent identities, negative why/import/load
   results, runtime unreachability, exact graph/module/tidy state, every
   earlier guard, and no new advisory, independent defect, qualified stable
   release, genuine supported tidy-stable owner, or compatible qualified
   route.
2. Authorize exactly one separate measurement-only owning-parent/request
   study. It may identify whether genuine compatible Cloud Go/storage request
   changes remove go-junit-report and measure every API, behavior, Go-floor,
   graph, guarded-selection, vulnerability, and project consequence. It may
   recommend a later route but may not implement one, add a direct target
   root, combine parent upgrades, or alter a finalized guard without its fresh
   owning decision.
3. Stop P7 explicitly with selected v0.9.1 unresolved and unaccepted. Prepare
   no implementation or dependency-evaluation successor and do not begin P8.

If none is acceptable, choose option 3. Do not infer acceptance from zero
loading, describe selected as qualified, promote a branch/pseudo-version,
select `/v2`, add a direct root, or transfer another dependency's exception.

# Required Reading And Moves

Read this archive and its answered evaluation, the answered json-iterator,
clockwork, and demangle decisions/evaluations, pprof and strcase records,
memberlist ownership chain, relevant retained-module records, rolling
handover, roadmap, `go.mod`, and `go.sum`. Reuse the evaluation; revalidate
only exact decision guards and fresh advisory identities needed for the
choice.

First verify branch, clean ordinary/ignored state, handoff identity/changed
set, reciprocal chain, exact Go identity, module hashes, all 25 target
requests and genuine parents, why/import/load state, every guarded selection/
edge, memberlist CNA identity, scratch containment, and
`./codex-dev-start.sh --check`. Stop if a premise changed.

Second, choose and record exactly one option with explicit ownership and
expiry bounds. This is documentation-only. Do not edit product source,
`go.mod`, or `go.sum`; repeat the writer-error fixture; run the option-2 study;
add a root; implement a workaround; alter a parent/guard; or execute P8.

Third, update the roadmap and rolling handover, answer this archive, and
prepare exactly one reciprocal NEXT mission matching the decision, or a
COMPLETE state if option 3 ends the authorized roadmap. Run applicable final
checks, verify scratch containment, and make only the required local
documentation handoff commit. Do not execute a successor.

# Automatic Handoff

Make only the required local `docs: prepare next agent session` commit. Do not
create an implementation commit, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, remove the worktree, combine another
dependency group, reopen json-iterator/clockwork/mvn-pom-mutator/demangle/
pprof/strcase/memberlist work, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
