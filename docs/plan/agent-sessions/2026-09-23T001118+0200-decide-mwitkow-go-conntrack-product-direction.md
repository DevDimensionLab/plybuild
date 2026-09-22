# Agent Session: Decide Mwitkow Go-Conntrack Product Direction

Status: NEXT
Session ID: `2026-09-23T001118+0200-decide-mwitkow-go-conntrack-product-direction`
Created: `2026-09-23T00:11:18+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `2a28c3fcb8c7cbed504e3b1449753ff449ee550bd911bdbae13d9983ca868591`
Previous: [2026-09-22T231926+0200-evaluate-mwitkow-go-conntrack-dependency.md](2026-09-22T231926+0200-evaluate-mwitkow-go-conntrack-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by recording exactly one product direction for graph-selected,
transitive, unloaded, unqualified exact
`github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223`.
Choose only one of the three authorized options below. This is one bounded
decision-recording group, not dependency implementation, owner research, a
new evaluation, another dependency group, or P8.

# Authorized Evaluation Result

The canonical exact-path stable release line is empty. Exact go-import maps to
public enabled unarchived non-fork Apache-2.0 repository
`mwitkow/go-conntrack`, but it has zero tags and zero GitHub Releases, and the
Go proxy stable list is empty. Selected
`v0.0.0-20161129095857-cc309e4a2223` and proxy `@latest`
`v0.0.0-20190716064945-2f068394615f` are pseudo-versions, not stables. Do not
promote either pseudo-version, `master`, a pull-request ref, fork, replacement,
or alternate path as a stable. No canonical stable exists or qualifies, so no
selection or dependency metadata changed.

Selected commit `cc309e4a22231782e8893f3c35ced0967807a33e`, tree
`8e653f1c675ad96001b221496975aa5858f0df09`, is an unsigned ancestor of
signed current master/latest commit
`2f068394615f73e460c2f3d2c158b0ad9321cadb`, tree
`a8074d12ad49ac88f1078e3942c1a95c8bfbc5f2`. Both use synthetic exact-path
module files without a Go directive, requirements, replacement, retraction,
or deprecation. There is no `/v2` or `/v3`. The selected proxy archive exactly
matches its 18-file Git tree. Preserve all recorded owner, license, commit,
tree, signature, ancestry, proxy, sumdb, archive, source, and module identities.

The selected source's historical Prometheus closure builds and passes count-
one tests and vet under exact Go 1.18.10 and Go 1.26.7, but repeated tests fail
under both SDKs because global metric state accumulates. Go 1.26.7 race testing
also finds concurrent historical `x/net/trace` event-map access. With the
project's actual Prometheus client v1.4.0/common v0.9.1 closure, full upstream
build/tests/race/vet fail under both SDKs at removed `prometheus.Handler` use;
the two library packages and all nine supported cross-build targets compile.
Preserve the recorded exported API, documentation, ownership/mutation,
determinism, concurrency, lifecycle, cleanup, errors, globals, platform,
network/TLS/file, cgo/build-tag/generated/embed, and Go-floor findings.

Exactly two graph edges request the selected pseudo-version: Prometheus common
v0.9.1 and historical common v0.4.1. Both genuinely import and use the target.
The current route is main -> direct mvn-pom-mutator v0.2.3 -> historical Viper
v1.10.1 -> go-metrics v0.3.10 -> common v0.9.1 -> target, with a second branch
through Prometheus client v1.4.0. The historical route continues common v0.9.1
-> Prometheus client v1.0.0 -> common v0.4.1 -> target. Target and requester
why results are negative, repository target imports are zero, and target
production, complete-test, module-backed load, and runtime relevance are zero.

Disposable exact gets are not retainable. The selected get keeps 234 modules
but promotes the target and seven already-selected package dependencies as
main roots, producing 82/1,076 metadata lines and 3,609 graph edges. The latest
pseudo get additionally selects `jpillora/backoff v1.0.0`, producing 235
modules, 83/1,079 metadata lines, and 3,610 edges. Both remain why-negative and
unloaded and tidy back to the exact common 52/948-line projection. Real state
remains 234/3,599/355/429/197/41/1,067 with the recorded module/tidy hashes and
Go 1.18 floor.

All 45 earlier guarded selections and 273 incoming edges remain exact at
`7e5c820da743ac628fb18da129b4d91428a36fd374ed248760d1318f2694c525`.
Every reflect2, concurrent, mapstructure, go-homedir, Promptui, emoji/v2,
kr/text, kr/pty, kr/pretty, Cast, Viper, closed-exception, owner/request,
why/import/load/runtime, graph/module/tidy/Go-floor, advisory, qualification,
and compatible-route fact remains final under its own recorded guard.

Exact target OSV, GitHub global, repository advisory, and pinned govulncheck
target findings are empty. The unchanged project retains exact 30/22/20/20
non-stdlib govulncheck populations. Guard OSV retains only the recorded Gorilla
WebSocket and go-retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and
GO-2026-6180. The 518,501-byte/1,402-record Go index and PUBLISHED 2,807-byte
memberlist CNA identities remain exact. Advisory absence is not qualification.
Final exact-Go verify/build/count-one/race/vet passes and accepted quality
remains 27/27 Q0-Q2 PASS at L2.

# Choose Exactly One Direction

1. Explicitly retain exact selected, inherited, unloaded
   `github.com/mwitkow/go-conntrack
   v0.0.0-20161129095857-cc309e4a2223` without product-source or dependency-
   metadata changes under a go-conntrack-specific, non-transferable exception.
   Call it unqualified: the exact-path stable line is empty, the selection is
   a pseudo-version, and its ordinary upstream closure gates are not all green.
   Accept only the completed evaluation and define the exact expiry guards
   below.
2. Authorize exactly one later bounded measurement-only study named
   **Prometheus Common Go-Conntrack Ownership Study**. Its sole question is
   whether a later genuine supported tidy-stable selection along the existing
   Prometheus common owner/request route removes go-conntrack or requests a new
   qualified Go-1.18-compatible exact-path stable while preserving every
   project and earlier-decision guard. Do not run the study, evaluate or change
   Prometheus common/client/go-metrics/Viper/mvn-pom-mutator, add a root, change
   a selection, or reopen an earlier decision in this move; prepare exactly one
   reciprocal measurement-only successor.
3. Stop P7 unresolved with no product-source, dependency-metadata, roadmap, or
   P8 implementation change and no successor execution.

Do not invent a fourth option or combine options. Physical MVS selection,
transitive presence, negative why, zero loading, an active repository, or
advisory absence is not qualification. Do not call either pseudo-version
stable or qualified, select latest/master, add a direct target root, change a
requester, transfer an exception, or begin P8.

# Option 1 Expiry Guards

If option 1 is selected, preserve and record at minimum:

- exact selected path/pseudo-version and both common v0.9.1/v0.4.1 requests;
- both genuine requester imports, complete current/historical routes, negative
  why/import/load/runtime facts, zero target root, and no supported stable
  owner/request route;
- the empty exact-path stable line, absent major lines, zero tags/releases,
  exact repository/owner/status/license/default branch, pseudo-version/commit/
  tree/signature/ancestry/archive/sumdb/module identities, and no replacement,
  retraction, deprecation, eligible alternate owner, or later stable;
- every selected-source API/behavior/closure/upstream/repeated/race/vet/cross-
  build/Go-floor finding, including actual-project Prometheus incompatibility;
- exact baseline and disposable-projection counts/hashes, tidy restoration,
  Go 1.18 floor, exact Go 1.26.7 identity, and all source/API/CLI/help/launcher/
  Make/quality contracts;
- all 45 earlier selections/273 incoming edges, all closed decisions and their
  separate expiry boundaries, and accepted 27/27 Q0-Q2 PASS at L2; and
- no new advisory, independent defect, release, owner, qualified exact-path
  stable, supported tidy-stable requester, or compatible genuine route.

Any path/version, request, requester import, owner route, root/why/import/load/
runtime, source/behavior/closure, graph/module/tidy/Go-floor, earlier guard,
advisory/finding, repository/release/owner, qualification, supported-owner, or
compatible-route change expires the exception and requires a fresh target
evaluation and product decision before merge. Option 1 authorizes no owner
study, direct root, dependency edit, workaround, or implementation.

# Role And Boundaries

This session records one decision only. Revalidate only the minimum continuity,
target/graph/advisory identities, launcher, and unchanged-project exact-Go
checks needed to record the chosen option safely. Do not repeat network/TLS
fixtures, upstream behavior work, or race reproduction. Do not launch a study
or successor. Preserve a clean worktree except for the bounded documentation/
launcher handoff, then make the required local handoff commit.

# Required Reading

Read this decision archive, the answered go-conntrack evaluation, answered
reflect2/concurrent evaluations, mapstructure, go-homedir, Promptui, emoji/v2,
kr/text, and kr/pty decisions/evaluations, kr/pretty and Cast owner records,
rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Verify the handoff
HEAD/parent/tree, exact changed set and ancestry, reciprocal archive chain,
clean ordinary/ignored state, launcher check, exact Go identities, module
hashes/counts/tidy projection, target requests/routes/why/import/load facts,
all 45 earlier guards/273 edges, closed-exception boundaries, and fresh
advisory identities.

# Three Moves

First, ask for or apply exactly one explicit option without reopening the
evaluation. Second, record only that option and its exact guards; change no
product source or dependency metadata. Third, update roadmap and rolling
handover, answer this archive, prepare exactly one reciprocal successor only
if the chosen option requires one, verify containment, and make the required
local handoff commit without executing a successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert,
bypass cleanup, remove the worktree, transfer an exception, reopen any earlier
decision, evaluate another dependency group, write outside the managed scratch
root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
