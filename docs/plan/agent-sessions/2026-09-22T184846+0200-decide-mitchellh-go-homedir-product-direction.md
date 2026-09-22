# Agent Session: Decide Mitchellh Go Homedir Product Direction

Status: NEXT
Session ID: `2026-09-22T184846+0200-decide-mitchellh-go-homedir-product-direction`
Created: `2026-09-22T18:48:46+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `6e83647c0e2b800f67344941c981055223c0496baa6dc166295c66558a974317`
Previous: [2026-09-22T180017+0200-evaluate-mitchellh-go-homedir-dependency.md](2026-09-22T180017+0200-evaluate-mitchellh-go-homedir-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by recording exactly one bounded product decision for selected
direct exact-path `github.com/mitchellh/go-homedir v1.1.0`. The completed
evaluation found that no canonical exact-path stable qualifies: v1.1.0 is
the latest Go-1.18-compatible stable, but its complete upstream suite fails
the ordinary documented HOME-unset discovery case on the supported Darwin
host under both exact SDKs; v1.0.0 repeats that failure and also cannot form
a guard-preserving project selection. Choose only one of the three directions
below, apply that choice exactly, and prepare only its reciprocal successor
if the choice requires one. Do not repeat the dependency evaluation, evaluate
another dependency group, or begin P8.

# Defensive Decision Scope

This is an ordinary dependency-quality product decision. Use only the
completed public release/repository metadata, static source and graph facts,
small bounded ordinary home-directory/path behavior results, project
projections, and advisory identities recorded here and in the answered
evaluation. Do not fuzz, stress, probe resource exhaustion, create oversized,
deeply nested, cyclic, malformed, adversarial, or escape-sequence payloads,
reproduce a security issue, or perform security or exploitability analysis.

Every disposable cache, tool, report, project copy, fixture, or advisory
response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to
`/private/tmp`, `/tmp`, a sibling of the managed root, or another external
root. Verify containment and remove task-owned scratch evidence before
handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and target-specific option-1 decisions through exact selected, direct,
production-loaded, runtime-relevant Promptui v0.9.0. P8 remains queued.

Promptui, emoji/v2, kr/text, and kr/pty option-1 decisions are separately
exact, unqualified, target-specific, non-transferable, and final. Their
selected and historical requests, requester imports, complete routes, why/
import/load/runtime facts, release/repository/source identities,
qualification results, graph/module/tidy/Go-floor state, projections,
advisories, all earlier guards, and compatible-route conditions remain
expiry guards. Any change requires the corresponding fresh dependency and
product decision before merge. Do not transfer an exception or reopen
Promptui, emoji/v2, kr/text, kr/pty, kr/pretty, Cast, or an earlier decision.

# Completed Evaluation Is Final

The evaluation began from clean branch `codex/upgrade-quality` at handoff
HEAD `35f8e0ee77586d93667420e00df1c6624c9ff4a5`, parent
`12881ba88e843657cd360c508937eeb00c38742c`, tree
`0ee5c4f31df0136618336879e584cc0fdfb00ba0`. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an
ancestor. The evaluation changed no product source or dependency metadata
and prepared this decision-only handoff. Verify the new handoff HEAD, parent,
tree, exact changed set, ancestry, reciprocal archive chain, and clean
ordinary and ignored status rather than assuming them.

Exact go-import metadata and the release module files resolve
`github.com/mitchellh/go-homedir` to public, enabled, archived, non-fork MIT
repository `mitchellh/go-homedir`, ID 23121609, owned by `mitchellh`,
defaulting to `main`, with no parent/source. The proxy exposes exactly v1.0.0
and v1.1.0; `@latest` is v1.1.0 at 2019-01-27T04:21:35Z. Exact `/v2` and
`/v3` lines are absent. There are no GitHub Releases, retractions, module
deprecations, replacements, prereleases, or eligible alternate paths, forks,
branches, or pseudo-versions.

Both tags are lightweight and form continuous ancestry. V1.0.0 is signed
commit `ae18d6b8b3205b561c79e8e5f69bff09736185f4`, tree
`30d387326d636916c152fe871e4305ee84f74b61`; GitHub cannot verify its old
key. V1.1.0 is validly signed commit
`af06845cf3004701891bf4fdb884bfe4920b3727`, tree
`c70b446f839c30d92fb3bbb491f0f0b540f47529`. Main contains v1.1.0. Both
proxy ZIPs byte-match their exact Git regular-file manifests, with no
symlink or submodule boundary. V1.1.0 source/module sums are
`h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=` and
`h1:SfyaCUpYCn1Vlf4IUYiD9fPX4A5wJrkLzIz1N1q0pr0=`. The answered evaluation
records both ZIP, manifest, license, sumdb, commit, tree, signature, and
ancestry identities.

Each release declares only the exact module path, with no Go directive,
requirements, retractions, or replacements. Both are standard-library-only,
single-package modules with no build tags, cgo, generated code, embed,
testdata, or network boundary. V1.0.0 exports DisableCache, Dir, and Expand;
v1.1.0 adds only Reset, which clears the cached home under the existing lock.
DisableCache is caller-managed mutable global state; the cached path and lock
are package-owned. HOME and operating-system user discovery are environment/
filesystem/process boundaries. Fixed environment and cache state are
deterministic; callers own no returned resource.

Under exact Go 1.26.7 and contained Go 1.18.10, both releases pass module
verification, build, vet, and supported Darwin amd64/arm64, Linux amd64/
arm64/386, Windows amd64/386, FreeBSD amd64, Plan 9 amd64, and js/wasm
cross-build/test gates. Their complete production/test closure is 61/125
packages, three module-backed entries, and one module under Go 1.26.7, and
42/81/three/one under Go 1.18.10.

A small ordinary fixture with a scratch-contained HOME verifies Dir; empty,
non-tilde, `~`, and `~/child` Expand behavior; rejection of unsupported
`~name`; caching; Reset where available; and four independent short
goroutines. It passes repeated and race runs under both SDKs. The complete
upstream count-one, count-ten, and race-count-ten suites nevertheless fail
for both releases under both SDKs. TestDir unsets HOME and expects discovery;
Darwin `dscl` returns `eServerError`, then the documented shell fallback
cannot `cd` because HOME is unset, so Dir returns `exit status 1`. Dir's
documentation allows a discovery error, but the P7 qualification contract
requires the complete upstream gate. No canonical exact-path stable
qualifies.

The graph contains seven exact target requests: main, mvn-pom-mutator v0.2.3,
and Viper v1.15.0 request v1.1.0; historical Viper v1.10.1 also requests
v1.1.0; go-rootcerts v1.0.2 and crypt v0.4.0 request v1.1.0; historical
go-rootcerts v1.0.0 requests v1.0.0. Main has requested direct v1.1.0 since
its initial module commit `d3543ac979220797a0ca0838e307ac12cd2b3557`.
Current repository production source imports the target only from
`pkg/config/profiles.go`, where `GetPlyHomePath` calls Dir. Three test files
import Reset. Main target why resolves through `pkg/config`; the package is
production and complete-test loaded with one module-backed entry and is
runtime relevant.

Mvn-pom-mutator imports and calls Dir from its own unloaded command while its
loaded `pkg/pom` owns the genuine current parent request. Viper has no target
source import at either graph version; those are metadata requests.
Go-rootcerts imports Dir only in Darwin source but is project-unloaded; crypt
has no target import and is unloaded. Mvn-pom-mutator and Viper have positive
project why, while go-rootcerts and crypt are negative. The shortest and
current supported owner route remains main's direct target request and
`pkg/config/profiles.go` runtime call.

A disposable exact-Go-1.26.7
`go get github.com/mitchellh/go-homedir@v1.1.0` is a byte no-op. The unchanged
project remains 234 selected modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41
loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy projection remains
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes remain
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

A disposable exact request for v1.0.0 cannot preserve project ownership or
guards: Go removes direct mvn-pom-mutator v0.2.3, downgrades Viper v1.15.0 to
v1.8.1, selects target v1.0.0, reduces the graph to 181 modules/2,551 edges,
and makes project loading fail because `cmd/build_add.go` can no longer import
mvn-pom-mutator/pkg/pom. Its raw go.mod/go.sum hashes are
`763cf59138d9c4e5e78d5b9dcac054808e1e8663cbb018e2497cb8af693602c7` /
`12e9395b66b4c87a1c4527f14a018adf38351cc9067232317f997fbf068b9379`.
`go mod tidy -diff` restores mvn-pom-mutator v0.2.3, target v1.1.0, and
Viper v1.10.1. No projection or implementation was retained.

All 41 earlier guarded selections remain exact. Their 240 sorted incoming
edges retain SHA-256
`a4e6042216821077b7a74472a6fcdc60d70644d07a805f6a6e99db1140a71563`.
Thirty-seven why results are negative; only closed kr/pretty, kr/text,
emoji/v2, and Promptui are positive. Promptui remains the only guarded
repository import, and only emoji/v2 and Promptui are production/complete-
test loaded. Every Promptui, emoji/v2, kr/text, and kr/pty expiry guard
remains exact.

Exact-version OSV and GitHub global queries are empty for both releases, and
the repository advisory endpoint is empty. Pinned isolated govulncheck v1.8.0
under exact Go 1.26.7 has no target module/package/symbol/test-symbol result
for either stable. Base and no-op-v1.1.0 project populations are identical at
30 module, 22 package, 20 symbol, and 20 test-symbol OSV IDs, with no target
trace. Guard OSV retains only Gorilla WebSocket GO-2026-6278/
GHSA-w67g-5rqw-f597 and go-retryablehttp GO-2024-2947/
GHSA-v6v8-xj6m-xwqh; x/mod v0.14.0 retains GO-2026-6179 and GO-2026-6180.
The Go index remains 518,501 bytes/1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. The PUBLISHED 2,807-byte memberlist CNA
response remains at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Advisory absence does not override the upstream gate.

Final unchanged-project exact Go 1.26.7 module verification, build,
count-one tests, race count-one tests, and vet pass under `umask 022`. The Go
1.18, source/API/CLI/help/launcher/Make/quality contracts and every earlier
decision remain exact. Accepted quality remains 27/27 Q0-Q2 PASS at L2.
Every task-owned scratch artifact was contained beneath the managed session
scratch root and removed; only the pre-existing launcher-owned Node compile
cache remains. Treat all completed release, source, behavior, closure, route,
projection, advisory, and final-gate findings as final.

# Choose Exactly One Direction

1. Explicitly retain exact selected, direct, production-loaded, runtime-
   relevant `github.com/mitchellh/go-homedir v1.1.0` without product-source
   or dependency-metadata changes under a go-homedir-specific,
   non-transferable exception. Call it unqualified: no canonical exact-path
   stable qualifies. Accept only the completed release/source, ordinary
   behavior, exact owner/request, route, why/import/load/runtime, graph/tidy/
   Go-floor, earlier-guard, and advisory facts. Define the exact expiry
   guards below. This exception must not broaden or expire the closed
   Promptui, emoji/v2, kr/text, or kr/pty exceptions.
2. Authorize exactly one later bounded measurement-only genuine owner/request
   study named `Ply Direct Go Homedir Ownership Study`. Its sole route is the
   existing main -> direct exact
   `github.com/mitchellh/go-homedir v1.1.0` request and the
   `pkg/config/profiles.go` Dir runtime use. Its exact question is whether
   that direct root and runtime route can be removed through a genuine
   supported project-owner change while preserving public API/CLI behavior,
   the Go 1.18 floor, all 41 earlier selections/240 edges, and every
   Promptui, emoji/v2, kr/text, and kr/pty exception. Do not run the study,
   change source or any selection, introduce a fork/replacement, or reopen an
   earlier decision in this decision-recording move; prepare one reciprocal
   measurement-only successor.
3. Stop P7 unresolved with no product-source, dependency-metadata, roadmap,
   or P8 implementation change and no successor execution.

Do not invent a fourth option or combine options. Physical selection,
directness, positive why, runtime loading, and advisory absence are not
qualification. Do not call v1.1.0 qualified, select v1.0.0, add or change a
root, use a fork, replacement, branch, pseudo-version, or alternate path,
change Promptui, emoji/v2, kr/text, kr/pty, kr/pretty, Cast, the Go floor,
product source, dependency metadata, or another module, transfer an
exception, or begin P8.

# Option 1 Expiry Guards

If option 1 is selected, preserve and record at minimum:

- exact selected path/version `github.com/mitchellh/go-homedir v1.1.0`,
  main's direct request, the other six exact current/historical requests, and
  all current requester selections;
- the historical direct-v1.1.0 root since the initial Go module, sole current
  production import and Dir use in `pkg/config/profiles.go`, three test Reset
  imports, complete main `pkg/config` -> target route, positive target why,
  one production/complete-test target package, and runtime relevance;
- the exact two-release line, absent `/v2` and `/v3`, repository/owner/status/
  license/default-branch identity, release/tag/commit/tree/signature/archive/
  sumdb/module facts, and no replacement/retraction/deprecation or eligible
  alternate;
- both stables remaining unqualified for the completed upstream-suite
  failure, and no future qualified Go-1.18-compatible canonical exact-path
  stable or supported route;
- baseline 234/3,599/355/429/197/41/1,067 state, exact module hashes, common
  tidy state, no-op v1.1.0 projection, rejected v1.0.0 projection, and no
  dependency implementation;
- exact Go 1.18 floor, exact Go 1.26.7 identity, source/API/CLI/help/launcher/
  Make/quality contracts, and accepted 27/27 Q0-Q2 PASS at L2;
- all 41 earlier selections/240 incoming edges and their exact hash, plus
  every Promptui, emoji/v2, kr/text, and kr/pty expiry guard; and
- no new target/closure advisory, independent defect, release, owner,
  qualified stable, supported owner, or compatible genuine route to a
  qualified go-homedir, Promptui, emoji/v2, kr/text, or kr/pty release.

Any target/request/requester-import/owner-route, root/why/import/load/runtime,
graph/module/tidy/Go-floor, source/behavior/closure, repository/release/owner,
advisory/finding, qualification, supported-owner, earlier guard, closed-
exception, or compatible-route change expires option 1 and requires a fresh
go-homedir dependency and product decision before merge. Any Promptui,
emoji/v2, kr/text, or kr/pty change separately requires its corresponding
fresh decision. The exception is not a security claim and cannot transfer.

# Decision Recording Contract

This successor records a product direction only. Reproduce the handoff,
chain, launcher, exact target/request/route, module/graph/tidy, closed-
exception, guard, owner/release, and fresh narrow advisory premises without
repeating the completed source, behavior, closure, upstream, cross-build,
projection, or govulncheck evaluation. Run final exact-Go unchanged-project
module verification, build, count-one tests, race count-one tests, and vet.

If option 1 is chosen, update the roadmap and rolling handover, answer this
archive, and prepare exactly one reciprocal evaluation of the next unanswered
selected queue item without executing it. If option 2 is chosen, prepare only
the named measurement-only successor. If option 3 is chosen, record the stop
without preparing a dependency evaluation or entering P8. Make only the
required local handoff documentation commit; do not create a dependency
implementation commit.

# Required Reading

Read this archive, the answered go-homedir evaluation, answered Promptui,
emoji/v2, kr/text, and kr/pty decisions/evaluations, kr/pretty and Cast owner
records, rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Verify
branch, ancestry, clean ordinary/ignored state, reciprocal archive chain,
launcher check, exact Go identities, module hashes/counts/tidy projection,
all seven target requests, source imports/uses, why/load/runtime state, all
41 earlier guards and the 240-edge snapshot, all four closed exception
boundaries, and fresh advisory identities before recording a choice.

# Three Moves

First, verify every completed premise without reopening the evaluation.
Second, choose and record exactly one option, make only its authorized
documentation/handoff changes, and run the unchanged-project final gate.
Third, verify scratch cleanup, reciprocal chain, single NEXT state, launcher
prompt mirror, exact changed set, and clean status; commit the local handoff
without executing its successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert,
bypass cleanup, remove the worktree, transfer an exception, reopen Promptui,
emoji/v2, kr/text, kr/pty, kr/pretty, Cast, or an earlier decision, repeat the
go-homedir evaluation, evaluate another dependency group, write outside the
managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
