# Agent Session: Decide Memberlist Ownership Blocker

Status: NEXT
Session ID: `2026-09-21T012527+0200-decide-memberlist-ownership-blocker`
Created: `2026-09-21T01:25:27+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `9ab40546b7b57363334ab3deb736c76d2fbef8903fc77f930750da942545c8e9`
Previous: [2026-09-21T004448+0200-integrate-go125-fixed-memberlist.md](2026-09-21T004448+0200-integrate-go125-fixed-memberlist.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one fresh product decision after the authorized
integrated Go-1.25/fixed-memberlist migration found no genuine supported
tidy-stable owner. Choose exactly one bounded direction: explicit guarded
retention of affected unloaded memberlist v0.3.0; a separately scoped POM
subsystem redesign intended to remove its owning graph; or stopping P7 with the
known affected selection unresolved. Do not implement the chosen direction,
revive the rejected projection, combine another dependency group, or begin P8.

# Authorized Roadmap

The integrated migration authorization is exhausted. First-fixed v0.6.0 was
preferred over v0.7.0, but direct memberlist and Serf roots are removed by
tidy, replacements are unsupported version masquerading, exclusions do not
select a fixed version, and the only syntactically durable projection declares
`github.com/hashicorp/serf/cmd/serf@v0.10.4` as a Go tool. This repository does
not use that CLI. Adding it solely to pin memberlist would manufacture
ownership and combine broad unrelated dependency changes with the remediation.

Do not treat Go's acceptance of a `tool` directive as project ownership. Do
not reopen that route unless a separate independently real product requirement
for the Serf CLI exists and is explicitly supplied. A build-tag, blank-import,
test-only, or other unused anchor is likewise manufactured. Do not create a
direct root, replacement, fork, patch, exclusion-based omission, or version
masquerade. Do not silently remove or redesign the POM parent, and do not
silently retain affected v0.3.0.

# Measurements At Start

The failed migration began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at HEAD
`7bb2629681b30b9b841f324f10b33f05690ff692`, parent
`e1ff4d49490ed94513d54807b5e19cebd6faaacc`, tree
`feb1073de931445bedd5272d9b9d4cf5541b577e`. Its handoff predecessor changed
exactly the launcher, answered Serf-study archive, final-direction archive,
then-NEXT integration archive, rolling handover, and roadmap. Verify the new
handoff rather than assuming its branch, clean state, reciprocal archive
chain, implementation ancestry, or launcher result.

The retained base is unchanged: exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
234 modules; 3,599 graph edges; 355 production and 429 complete-test entries;
197 module-backed complete-test entries across 41 loaded modules; 1,067 sum
lines; and a 432-line tidy projection. `go.mod`/`go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

All 23 earlier guarded selections and 164 incoming edges remain exact at
snapshot SHA-256
`f42304a43a7ade746c5c2bc238dfa9a6daec62594834fd5b606fdf245d44ea2e`.
All 24 target-plus-earlier why results are negative; repository imports and
production/complete-test guarded loads are zero. Selected Serf is v0.10.1 and
memberlist is v0.3.0. The five incoming Serf requests and both historical
Serf-to-memberlist requests remain exact. Memberlist has no direct root,
import, load, or runtime reachability.

The PUBLISHED HashiCorp CNA response remains byte-identical at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`:
memberlist below v0.6.0 is affected and v0.6.0 is first fixed. The Go
vulnerability module index remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. Exact OSV remains empty for memberlist and
the Serf candidates; earlier guards retain only the recorded Gorilla and
go-retryablehttp pairs. Empty secondary feeds do not override the primary CNA
record.

# Failed Projection To Reuse

The disposable Go-tool projection retained Go 1.25.0, toolchain go1.26.7,
Serf v0.10.4, and memberlist v0.6.0 across tidy, but it existed only because
of the unused Serf CLI declaration. It produced 261 modules, 3,826 edges,
1,031 sum lines, 355 production entries, 402 complete-test entries, and 170
module-backed complete-test entries across 41 modules. Its tool closure alone
contained 270 packages. It removed 51 graph edges, added 278, moved six
guarded selections, and added 31 guarded incoming edges, yielding guarded
snapshot SHA-256
`bef5bb840a3f67db60fd3fdc2460cbd1dd4e10c70bf7da1065aabe1ec82f9616`.
No delta was accepted and no projection or implementation commit was retained.

The six moved guarded selections were errwrap v1.1.0, go-multierror v1.1.1,
go-retryablehttp v0.7.7, go-sockaddr v1.0.7, golang-lru v1.0.2, and mdns
v1.0.7. The broader selection churn included unrelated CLI, metrics,
compression, serialization, Prometheus, crypt, protobuf, terminal, and x/*
modules. Because the owner and focused-change gates failed first, no owning
decision was granted for these changes and no Go-floor or stricter-format
source integration was attempted.

The completed memberlist evaluation, Serf/owning-graph study, candidate
comparison, and failed owner qualification are final. Do not repeat release,
security, exploitability, attack-traffic, host, credential, or security-
boundary investigations. Revalidate only the exact decision guards and fresh
primary advisory identity needed to choose under current premises.

# Role And Boundaries

Choose exactly one:

1. Explicitly retain exact selected, inherited, unloaded memberlist v0.3.0
   under a new target-specific, non-transferable affected-version risk
   decision. State that v0.3.0 is affected and not secure. Bound the decision
   by zero import/load/runtime reachability; exact Serf, parent, graph, module,
   advisory, affected/fixed-range, and earlier-guard identities; and immediate
   expiry on any changed guard, new finding, real fixed owner, or compatible
   fixed route. Authorize no source, dependency, floor, parent, or P8 change.
2. Authorize one separately scoped POM subsystem redesign successor. It must
   preserve the public POM behavior used across the recorded 21 source/test
   files while removing or replacing `mvn-pom-mutator v0.2.3`, measure every
   API/CLI/runtime/fixture and dependency-graph consequence, make fresh owning
   decisions for every guarded delta, and retain implementation only after its
   own complete gates pass. This is a product redesign, not a memberlist pin or
   permission to manufacture an owner.
3. Stop P7 explicitly with selected v0.3.0 still affected, unresolved, and not
   accepted. Prepare no implementation successor and do not start P8.

If none is acceptable, choose option 3. Do not combine options or reinterpret
zero current loading as proof that the affected release is secure.

# Required Reading

Read this archive, its answered migration, final-direction decision, Serf
study, memberlist decision/evaluation, relevant retained-module and POM
ownership records, rolling handover, roadmap, `go.mod`, and `go.sum`. Reuse
their final measurements and decisions; do not reopen completed studies.

# Three Moves

First, verify branch, clean ordinary and ignored state, handoff HEAD/parent/tree
and changed set, reciprocal chain, latest Google UUID implementation ancestry,
exact Go identity, retained graph/import/load guards, module hashes, fresh
primary advisory identity, and `./codex-dev-start.sh --check`. Stop for
investigation if any decision premise differs.

Second, choose and record exactly one option with explicit ownership, expiry,
and successor bounds. This is a documentation-only product decision. Do not
edit product source, `go.mod`, or `go.sum`; run a migration; accept guarded
selection changes; repair stricter-format calls; or execute a POM redesign.

Third, update the roadmap and rolling handover, answer this archive, and
prepare exactly one reciprocal NEXT mission matching the decision, or a
COMPLETE state if option 3 ends the authorized roadmap. Run all applicable
final checks and make only the required local documentation handoff commit.
Do not execute the successor.

# Automatic Handoff

Make only the required local `docs: prepare next agent session` commit. Do not
create an implementation commit, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, remove the worktree, combine another
dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
