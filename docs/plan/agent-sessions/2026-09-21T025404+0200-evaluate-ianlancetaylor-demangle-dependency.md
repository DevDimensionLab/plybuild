# Agent Session: Evaluate Ianlancetaylor Demangle Dependency

Status: NEXT
Session ID: `2026-09-21T025404+0200-evaluate-ianlancetaylor-demangle-dependency`
Created: `2026-09-21T02:54:04+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `31fb189fcdcd0efbdc461202a2dda7bcad94e8ad643e3f126399f0c20a7c8ac4`
Previous: [2026-09-21T022425+0200-decide-iancoleman-strcase-product-direction.md](2026-09-21T022425+0200-decide-iancoleman-strcase-product-direction.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/ianlancetaylor/demangle
v0.0.0-20200824232613-28f6c0f3b639` as one bounded dependency group.
Resolve its complete repository and release/pseudo-version identity, Go-floor
closure, exported API and demangling behavior, actual project loading, exact
MVS effects, vulnerability evidence, and every applicable quality contract.
Retain or select only a qualified exact-path version whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision. Do not reopen
strcase, protoc-gen-validate, memberlist, Serf, or POM work, combine another
dependency group, or begin P8.

# Defensive Retry Scope

This is an ordinary defensive dependency-quality review of the user's local
repository, not a security or abuse investigation. The first attempt from
clean HEAD `3e4354cebb1d18d189cf7c173773b19103b32c57` was automatically
stopped at a policy boundary after it began generating unusually deep and
non-ordinary parser inputs. It made no tracked worktree change. Resume this
same session under the narrower rules below instead of repeating the broad
evaluation.

Use only public repository, proxy, sumdb, release, and advisory metadata;
static source and API inspection; existing upstream tests; ordinary documented
valid examples; and normal project graph/build commands. Do not generate,
mutate, or test malformed, empty, deeply nested, oversized, randomized, or
adversarial symbol inputs. Do not fuzz, stress, probe resource exhaustion,
reproduce crashes, or perform security or exploitability analysis. The
already observed empty-argument command panic is final evidence: record it
without running it again. If a remaining conclusion would require prohibited
testing, state that limitation and stop for the bounded product decision.

Reuse these provisional results from the interrupted attempt: the selected
commit `28f6c0f3b63983aaa99575ca3b693afff7996387` and examined current head
`83e58baca7248962d58657affe41b3b6f27ee423` both declare Go 1.13 and have
standard-library-only closures; both passed existing upstream tests, count-10
repeats, race, and vet under exact Go 1.26.7 and contained Go 1.18.10, plus the
completed ordinary cross-builds. Static repository inspection found a public,
active, unarchived, non-fork BSD-3-Clause repository with no tags or GitHub
releases. Both examined `c++filt` commands panic on an empty argument. Confirm
only the exact current pseudo-version identity, static exported-API delta,
ordinary valid-input behavior, public advisory state, and project/MVS effects
still needed for the decision. Do not repeat completed closure or stress-style
checks.

Every disposable archive, clone, cache, tool, binary, report, and fixture must
be created beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write directly to
`/private/tmp`, `/tmp`, another external root, or a sibling of the managed
scratch directory. If a nested directory is needed, use
`mktemp -d "${CODEX_SESSION_SCRATCH_ROOT:?}/demangle.XXXXXX"`. The launcher
owns automatic cleanup of that managed root. Before handoff, verify that no
task-owned disposable path exists outside it; do not retain scratch evidence.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact selected,
inherited and unloaded strcase v0.2.0. Every earlier outcome remains final
under its own guards. P8 remains queued.

The 2026-09-21 strcase option-1 decision explicitly retains exact v0.2.0
without source or dependency metadata changes. V0.2.0 is not qualified. Its
target-specific, non-transferable exception accepts only the completed
concurrent acronym-map races/fatal crash, uppercase and Unicode/malformed
conversion loss, permanent global acronym state, and the characterized API,
digit, delimiter, allocation, MVS, vulnerability, and related findings. It is
bounded by exact v0.2.0, the sole protoc-gen-validate v0.6.2 request, both
historical incoming parent requests, no direct root/import/load/runtime
reachability, exact graph/module/tidy and earlier guards, and no new finding,
qualified stable release, genuine owner, or compatible route. Any change
expires that decision and requires its fresh owning decision. No strcase or
memberlist exception transfers to demangle.

Selected demangle is currently inherited only through thirteen historical
`github.com/google/pprof` requests. Selected pprof remains exact
`v0.0.0-20210720184732-4bb14d4b1be1`, and its earlier evaluation retained that
exact pseudo-version under its completed findings. Demangle has no direct root
or repository Go import, its `go mod why -m` result is negative, and production
and complete-test closures load zero target packages. These queue observations
and physical MVS selection are not release qualification or authorization to
retain the target. Resolve them independently and do not add a direct edge
merely to alter MVS.

# Measurements At Start

The strcase decision began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at HEAD
`5071025e1061015964c752dbd1004a07c72857c3`, parent
`c0d76f12ae3c2f4b132406eabb3f16d2714ed6f6`, tree
`f1625d3179abc6e4febe53203cae4fca71a5b57f`. That handoff changed exactly the
launcher, answered strcase evaluation archive, then-NEXT strcase decision
archive, rolling handover, and roadmap. Its reciprocal 232-archive chain,
latest Google UUID implementation ancestry, and launcher check passed. Verify
the new handoff rather than assuming these facts.

Exact Go 1.26.7 binary SHA-256 is
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
The unchanged project has 234 modules, 3,599 graph edges, 355 production and
429 complete-test entries, 197 module-backed complete-test entries across 41
loaded modules, 1,067 sum lines, and a 432-line tidy projection. `go.mod` and
`go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
No dependency implementation or metadata commit exists after Google UUID;
accepted quality remains 27/27 Q0-Q2 PASS at L2.

With strcase included, all 25 guarded selections and their 167 incoming graph
edges are exact at sorted snapshot SHA-256
`b6e0bcee17f83b0cc88a370c2f86a51c1872506b113075960c497bf420c48fc9`.
All 26 demangle-plus-guarded why results are negative; repository imports are
zero; and production and complete-test closures load zero target or guarded
packages. The thirteen demangle requests remain only the recorded Google pprof
vertices; the selected request is
`google/pprof@v0.0.0-20210720184732-4bb14d4b1be1 ->
ianlancetaylor/demangle@v0.0.0-20200824232613-28f6c0f3b639`.

Fresh primary memberlist evidence remains the byte-identical PUBLISHED
HashiCorp CNA response at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
with memberlist below v0.6.0 affected. The Go vulnerability index remains
1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. Earlier exact-version results retain only
the recorded Gorilla and go-retryablehttp pairs. Revalidate target and guard
advisories without reopening completed strcase or memberlist work.

# Role And Boundaries

From fresh public archives and caches, resolve proxy, sumdb, `go-import`, Git,
and forge evidence for the exact demangle module path: versions, tags,
pseudo-versions, commits, ancestry, repository status, license, retractions,
deprecation, redirects, alternate paths, major lines, and every serious
exact-path candidate. Do not silently promote a fork, branch, prerelease,
redirect, alternate path, version-masquerading replacement, or floor-
ineligible candidate.

Treat the completed Go 1.26.7 and Go 1.18.10 source/test-closure, native-test,
repeat, race, vet, and cross-build results as provisional evidence requiring
only narrow consistency checks. Inspect packages, exported API, supported
formats and options, examples, tests, generated files, build tags, platform
branches, and actual project consumers statically. Any additional execution
must use existing upstream tests or ordinary documented valid inputs only. Do
not add independent parser fixtures or expand this dependency qualification
into robustness, security, crash, resource-limit, or hostile-input research.

Measure exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects in disposable trees.
Explain why demangle exists in MVS and whether any package loads. Preserve
strcase, memberlist, Google pprof, and every earlier guarded decision. A
parent, Go-floor, unrelated-selection, or non-exact-path change requires its
own fresh bounded decision; do not manufacture a direct dependency owner.

# Required Reading

Read this archive, the answered strcase decision and evaluation, the answered
memberlist ownership decision and migration/study/evaluation chain, the
answered Google pprof evaluation, relevant retained-module records, rolling
handover, roadmap, `go.mod`, and `go.sum`. Verify branch, clean ordinary and
ignored state, handoff HEAD/parent/tree and changed set, reciprocal archive
chain, latest Google UUID implementation ancestry, exact Go identity, all
target/parent and guarded selection/request/why/import/load/advisory
conditions, module hashes, and `./codex-dev-start.sh --check`. Earlier
outcomes are final.

# Three Moves

First, revalidate the starting guards and finish only the unresolved safe
identity, static API, ordinary behavior, public advisory, loading, and exact
MVS facts while reusing the completed evidence above. Second, account for the
recorded empty-argument panic without reproducing it. If and only if one exact-
path candidate preserves Go 1.18 and passes every applicable ordinary contract,
implement that dependency-only selection and run the normal changed-selection
gate; otherwise leave metadata unchanged and stop for one bounded product
decision. Third, update the roadmap and rolling handover, answer this archive,
prepare exactly one reciprocal NEXT mission for the authorized result, verify
scratch containment, and commit the handoff without executing the successor.

# Automatic Handoff

After one coherent bounded outcome, make any separate dependency-only commit
first if a qualified selection was implemented, then make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, reopen strcase/protoc-gen-validate/
memberlist/Serf/POM work, write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
