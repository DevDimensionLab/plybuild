# Agent Session: Decide Prometheus Procfs Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T193000+0200-decide-prometheus-procfs-product-direction`
Created: `2026-09-23T19:30:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c057de9a8bc37f2358f719cea929287dc2d486286fe0161ce06e2a9ff2d5a012`
Previous: [2026-09-23T172541+0200-evaluate-prometheus-procfs-dependency.md](2026-09-23T172541+0200-evaluate-prometheus-procfs-dependency.md)
Next: [2026-09-23T195118+0200-evaluate-prometheus-tsdb-dependency.md](2026-09-23T195118+0200-evaluate-prometheus-tsdb-dependency.md)
Outcome: Option 1 is final. Exact inherited indirect unloaded procfs v0.0.8 remains unchanged only under its own unqualified, non-transferable exception; no study, dependency or source change, ownership claim, target root, transferred exception, or P8 work was authorized.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product decision for graph-selected
transitive exact `github.com/prometheus/procfs v0.0.8`. The completed fresh
evaluation found no fully qualified canonical Go-1.18-compatible stable.
Choose and record exactly one of the three authorized directions below. Do not
repeat the evaluation, implement a dependency change, combine another group,
launch a study or successor, or begin P8.

# Defensive Scope

This is an ordinary dependency product-direction decision. Use only the
completed public release/source/build/graph/projection/advisory evidence and
bounded read-only guard checks. Do not fuzz, stress, probe resource exhaustion,
create oversized, deeply nested, cyclic, malformed, adversarial, or escape-
sequence payloads, reproduce a security issue, or perform security or
exploitability analysis.

Every disposable cache, report, project copy, or advisory response must remain
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`,
`/tmp`, a sibling of the managed root, or another external root. Set every
tool temp/cache root explicitly, verify containment, and remove task-owned
scratch evidence before handoff.

# Completed Evaluation

Exact go-import metadata resolves the canonical module to the public, enabled,
unarchived, non-fork Apache-2.0 `prometheus/procfs` repository on `master`.
The exact path has 48 stable releases v0.0.1-v0.22.0 and no prerelease,
replacement, retraction, deprecation, `/v2`, or `/v3` line. Exactly 27 stables
preserve Go 1.18 through highest v0.9.0. Exact tag/commit/tree/signature/
ancestry, proxy/sumdb/archive-to-Git, module/license, package/build boundary,
API, and documented ordinary-behavior evidence is complete for all 27.

Every eligible stable verifies and builds under exact Go 1.18.10 and Go
1.26.7. Selected v0.0.8 verifies, builds, and vets but fails complete
count-one, count-ten, and race tests under both SDKs: its ordinary `TestVM`
fixture result is wrong on Darwin and Btrfs/sysfs suites return their
unsupported-platform errors. Its complete closure is the main module plus
test-only go-cmp v0.3.1 and x/sync `cd5d95a43a6e`.

Highest compatible v0.9.0 verifies, builds, and passes count-one, race, and
vet under both SDKs when its exact Git tree retains the seven test-fixture
symlinks that module ZIPs must omit. It fails count-ten under both SDKs because
`TestNetSoftnet` repeatedly supplies four columns where the parser requires at
least nine. Its complete closure adds go-cmp v0.5.9, x/sync v0.1.0, and x/sys
v0.3.0. No lower stable passes every mandatory native gate.

The complete cgo-disabled matrix covers all 27 releases, both SDKs, and ten
Darwin/Linux/Windows/FreeBSD/Plan-9/js-wasm targets. Selected v0.0.8 passes all
40 selected production-build and test-compilation cells. V0.9.0 passes all
production builds but fails Linux/386 and Windows/386 test compilation under
both SDKs because its tests assign `math.MinInt64`/`math.MaxInt64` to `int`.
Thus no canonical stable fully qualifies.

Exact-version OSV results for all 27 compatible releases and narrow selected/
v0.9.0 GitHub/repository responses are empty without implying qualification.
Pinned isolated govulncheck v1.8.0 is empty for v0.0.8. V0.9.0 has only
module-level GO-2026-5024 in x/sys v0.3.0 and no package, symbol, test-symbol,
or Procfs trace. No exploitability claim was made.

The graph has exactly four target requests. TSDB v0.7.1 and Common v0.4.1 are
metadata-only and request `v0.0.0-20181005140218-185b4288413d`.
Client_golang v1.4.0 genuinely imports Procfs in production and tests and
requests v0.0.8; client_golang v1.0.0 has the same genuine import boundary and
requests v0.0.2. The four complete shortest routes run from main through
direct mvn-pom-mutator v0.2.3 and TSDB or historical Viper v1.10.1/go-metrics/
Common/client_golang chains.

Target/requester why is negative. Repository target/requester imports,
target/requester production and complete-test loads, target module-backed
load, runtime relevance, current target roots, and historical target roots are
zero. `go.sum` has exactly three Procfs go.mod sums and no source sum. No
requester asks for v0.9.0.

A disposable selected exact get only adds a redundant indirect root, source
sum, and main edge; ordinary tidy removes all three and restores inherited
v0.0.8 plus the established common projection. No higher projection was
authorized because no higher stable was otherwise qualified. No projection,
dependency change, source change, owner, or exception was retained.

The unchanged real project remains exactly 234 modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
over 41 loaded modules, and 1,067 sum lines at the protected hashes and Go
1.18 floor. Project govulncheck remains 30/22/20/20 with no Procfs trace;
client_golang v1.4.0 retains GHSA-cg3q-j54f-5p7p / GO-2022-0322 /
CVE-2022-21698. Corrected Go-index/CNA identities, all 47 pre-Goe selections/
276 incoming edges, separate Goe/pkg-errors/SFTP/go-difflib/Complete/ULID/go-
conntrack/client_golang/client_model/Common guards, every earlier decision, and
accepted 27/27 Q0-Q2 PASS at L2 remain exact.

# Authorized Decision

Choose exactly one:

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/prometheus/procfs v0.0.8` unchanged under a Procfs-specific,
   unqualified, non-transferable exception. State that it is not qualified,
   supported, safe, or fixed. Preserve all four exact requests, requester
   import/metadata boundaries, four routes, negative why, zero repository
   import/load/runtime/current-and-history-root state, release/source/API/
   behavior/closure/native/cross/projection/advisory identities, and every
   earlier guard as expiry conditions. Do not transfer the Common,
   client_model, client_golang, or any earlier exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator Prometheus
   Procfs Owner/Request Study** to determine whether a supported update of the
   direct product root can eliminate every Procfs request or establish genuine
   tidy-stable ownership of a separately qualified future Go-1.18-compatible
   stable while preserving product behavior, direct roots, Go 1.18, and every
   earlier guard. Do not run the study, change a dependency, grant an
   exception, or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because the target is unloaded and runtime-irrelevant,
all qualification and ownership failures are explicit, and it preserves the
unchanged project without broadening a study or transferring an exception.
This is a product acceptance decision, not a qualification claim.

# Protected State And Checks

Start only from the clean Procfs evaluation handoff. Verify its HEAD, parent,
tree, exact changed set, branch/ancestry, reciprocal archive chain, sole NEXT
state, launcher mirror, ordinary and ignored cleanliness, official SDK
identities, real project counts and hashes, common tidy projection, Go 1.18
floor, all target request/import/route/relevance/root facts, target/closure/
project advisory identities, and every earlier guard. Stop for a fresh owning
decision if any protected input changed.

Preserve exact selected Common v0.9.1/four requests, client_model v0.2.0/15
requests, and client_golang v1.4.0/four requests only under their own separate
unqualified non-transferable exceptions. Do not run any rejected owner study,
select their rejected candidates, add a target root, or transfer an exception.
Preserve Complete, go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack, and
all earlier qualified or excepted results under their own exact guards. P8
remains queued.

Do only guard-level revalidation; do not repeat completed release, source,
behavior, API, closure, projection, native, cross, or advisory-matrix work.
Record the chosen product direction, update the roadmap and rolling handover,
answer this archive, verify containment and cleanup, run final exact-Go-1.26.7
project module verification, build, count-one tests, race count-one tests, and
vet, and make the local handoff commit.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, transfer or reopen an exception,
change a dependency or source file, combine another dependency group, or begin
P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is final. Exact selected, inherited, indirect, unloaded
`github.com/prometheus/procfs v0.0.8` is explicitly retained unchanged under a
Procfs-specific, unqualified, non-transferable exception. It is not qualified,
supported, safe, or fixed. Selected v0.0.8 fails mandatory complete native
tests under both supported SDKs. Highest compatible stable v0.9.0 fails
mandatory repetition under both SDKs and supported Linux/386 and Windows/386
test compilation; no canonical stable fully qualifies.

Option 2's **Mvn-Pom-Mutator Prometheus Procfs Owner/Request Study** is not
authorized or run. A direct-product-root study is not warranted while Procfs
is unloaded and runtime-irrelevant and all exact qualification and ownership
failures are explicit. Option 3 is not selected because the unchanged
inherited state can be accepted within exact expiry guards without making a
qualification claim. No dependency or source change, target root, ownership
claim, implementation pre-authorization, transferred exception, other
dependency group, or P8 work is included.

### Exact exception and expiry boundary

Retention requires all four exact requests and requester boundaries to remain
unchanged:

- Prometheus TSDB v0.7.1 metadata-only requests
  `v0.0.0-20181005140218-185b4288413d` and has no Procfs source import;
- Prometheus Common v0.4.1 metadata-only requests that same pseudo-version and
  has no Procfs source import;
- client_golang v1.4.0 requests v0.0.8 and retains one genuine Procfs
  production import and one genuine test import; and
- client_golang v1.0.0 requests v0.0.2 and retains the same one-production /
  one-test genuine import boundary.

All four complete shortest routes from the main module must remain exact:
main -> direct mvn-pom-mutator v0.2.3 -> TSDB v0.7.1 -> the Procfs pseudo-
version; main -> mvn-pom-mutator -> historical Viper v1.10.1 -> go-metrics
v0.3.10 -> Common v0.9.1 -> client_golang v1.0.0 -> Common v0.4.1 -> the
pseudo-version; main -> mvn-pom-mutator -> Viper -> go-metrics ->
client_golang v1.4.0 -> v0.0.8; and main -> mvn-pom-mutator -> Viper -> go-
metrics -> Common v0.9.1 -> client_golang v1.0.0 -> v0.0.2.

The exception additionally requires negative target and requester `go mod why
-m`; zero repository target/requester imports; zero target/requester
production and complete-test loads; zero target module-backed load and runtime
relevance; zero current and historical main-module Procfs roots; exactly the
pseudo-version, v0.0.2, and v0.0.8 go.mod sums with no Procfs source sum; and
no request for v0.9.0. Any changed request, requester import/metadata boundary,
route, why/import/load/runtime/root/sum fact, or newly supported owner route
expires retention and requires a fresh owning evaluation and explicit product
decision before merge.

Every canonical release and source identity recorded in the answered Procfs
evaluation remains an expiry guard: the exact public, enabled, unarchived,
non-fork Apache-2.0 Prometheus repository on `master`; exactly 48 stable exact-
path releases v0.0.1-v0.22.0; absent prerelease, replacement, retraction,
deprecation, `/v2`, and `/v3` lines; and exactly 27 Go-1.18-compatible stables
through v0.9.0. Exact tag/commit/tree/signature/ancestry, proxy/sumdb/archive-
to-Git, module/license, package/build boundary, API, and documented ordinary-
behavior evidence for all 27 remains incorporated by reference as a non-
transferable expiry condition.

Selected v0.0.8 remains signed lightweight commit
`6d489fc7f1d9cd890a250f3ea3431b1744b9623f`, tree
`b333b5669051a6fb51b220bd142650f66731fbe8`, with source/mod sums
`h1:+fpWZdT24pJBiqJdAwYBjPSk+5YmQzYNPYzQsdzLkt8=` /
`h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=`. Highest compatible v0.9.0
remains signed lightweight commit
`bb7727a9ca9b3cd1559b42ffa74e13ef7efb6283`, tree
`942019e4c3a443c18bf5430166a43475fe42bc2b`, with source/mod sums
`h1:wzCHvIvM5SxWqYvwgVL7yJY8Lz3PKn49KQtpgMYJfhI=` /
`h1:+pB4zwohETzFnmlpe6yd2lSc+0/46IYZRB/chUwxUZY=`. Their exact ZIP, archive,
license, package, build-boundary, API, and behavior identities remain guarded.

The complete closure/native/cross identities remain exact. Selected v0.0.8's
closure remains the module plus test-only go-cmp v0.3.1 and x/sync
`cd5d95a43a6e`; it verifies, builds, and vets but fails count-one, count-ten,
and race tests under exact Go 1.18.10 and Go 1.26.7 on its Darwin VM fixture
and Btrfs/sysfs unsupported-platform results. Its complete 40-cell selected
cross boundary remains passing. V0.9.0's closure remains go-cmp v0.5.9,
x/sync v0.1.0, and x/sys v0.3.0; its exact Git tree passes count-one, race,
and vet under both SDKs but fails count-ten on the four-column softnet fixture
and fails Linux/386 and Windows/386 test compilation under both SDKs. Any newly
qualified canonical stable expires this exception.

The selected disposable projection remains an expiry guard. Exact selected
get has 75/1,068 lines, 234 modules, 3,600 edges, unchanged 355/429/197/41
loads, and `go.mod` / `go.sum` / graph hashes
`d3f6b5f5b31869a6f7dc23060d27a75014c6297714b7e4a41ac89c0d4431d621` /
`3ee61692fbd1b28df65bd666b8e0256719fcd9b245f2a5358a0a0351e37c8543` /
`27748f6be14bea6db88fae500a8789a0b96ab469675be7dd949b9c1e7293486d`.
Ordinary tidy must continue removing the redundant root, source sum, and main
edge and restoring inherited v0.0.8 plus the established common projection.
No higher projection was authorized or retained.

Exact empty target OSV and selected/v0.9.0 GitHub/repository responses remain
guards, and their absence is not qualification. Pinned govulncheck v1.8.0 at
database timestamp 2026-09-16T18:00:43Z must keep selected v0.0.8's isolated
closure empty. V0.9.0 must retain only module-level GO-2026-5024 in x/sys
v0.3.0, with no package, symbol, test-symbol, or Procfs trace. The unchanged
project must remain 30/22/20/20 with no Procfs trace; client_golang v1.4.0
must retain GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698. Advisory
absence does not establish qualification, support, safety, or a fix.

This Procfs exception does not transfer the separate Common v0.9.1/four-
request, client_model v0.2.0/15-request, or client_golang v1.4.0/four-request
unqualified, non-transferable exceptions. It also does not transfer, reopen,
or alter Complete, go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack, or
any earlier qualified or excepted result. Each remains governed only by its
own exact selection/request/route/relevance/source/closure/projection/advisory
and expiry guards.

### Guard-only verification and handoff

- Work began clean on `codex/upgrade-quality` at Procfs evaluation handoff
  HEAD `c1810321150709acf53fbe9e0e80ea87659e4176`, parent
  `64f06865d2d263847b67a7a9737477e98571f4e7`, tree
  `18018307d18adaa33b2366d422d01cacaac6953f`. It changes exactly the launcher,
  answered Procfs evaluation, this then-NEXT decision, rolling handover, and
  roadmap. Google UUID implementation
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor.
- The reciprocal 294-archive chain, sole NEXT state, launcher/archive prompt
  mirror and check, branch, exact changed set, and ordinary and ignored
  cleanliness reproduce. Completed release, source, behavior, API, closure,
  projection, native, cross, and advisory-matrix work was not repeated.
- Fresh contained official Darwin arm64 SDKs reproduce Go 1.18.10 archive /
  binary SHA-256
  `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and Go 1.26.7
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
- The unchanged project reproduces exactly 234 modules, 3,599 graph edges,
  355 production entries, 429 complete-test entries, 197 module-backed entries
  across 41 loaded modules, and 1,067 sum lines. `go.mod`, `go.sum`, and graph
  SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  The Go floor remains 1.18. Guard-only tidy in a contained project copy
  reaches the exact common 52/948-line, 234-module/3,557-edge projection at
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`
  and reselects v0.0.8. No projection is retained.
- All four requests, requester import/metadata boundaries, route edges,
  negative target/requester why, zero repository import/load/runtime/current-
  and-history-root state, and target sums reproduce. The byte-exact full graph
  preserves the 47 pre-Goe selections and 276 incoming edges at
  `7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
  Separate selections/request counts reproduce: Goe v0.1.0/four, pkg/errors
  v0.9.1/ten, SFTP v1.13.1/four, go-difflib v1.0.0/23, Complete v1.2.3/three,
  ULID v1.3.1/one, go-conntrack `cc309e4a2223`/two, client_golang v1.4.0/four,
  client_model v0.2.0/15, and Common v0.9.1/four. Every earlier qualified or
  excepted result remains closed, separate, and untransferred; accepted
  quality remains 27/27 Q0-Q2 PASS at L2.
- Fresh narrow Procfs v0.0.8/v0.9.0 OSV and GitHub global/repository responses
  remain empty. Pinned govulncheck reproduces empty v0.0.8 isolation,
  v0.9.0's sole module-level GO-2026-5024, and unchanged-project 30/22/20/20
  populations without a Procfs trace. Guard OSV retains the recorded Gorilla,
  retryablehttp, x/mod, and client_golang findings. Normal/no-cache Go-index
  responses remain pairwise byte-identical at 518,501 bytes/1,402 records and
  SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  the PUBLISHED memberlist CNA responses remain pairwise byte-identical at
  2,807 bytes and
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

Final unchanged-project exact-Go-1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass under ordinary umask 022. Product
source, `go.mod`, and `go.sum` remain byte-exact. Every task-owned SDK, cache,
response, report, tool, and project copy was contained beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}` and removed before handoff.

P7 remains active only with one prepared bounded evaluation of the next
selected alphabetical module, exact `github.com/prometheus/tsdb v0.7.1`.
Its single graph request from direct mvn-pom-mutator v0.2.3 is a starting
observation only. The successor is prepared but not launched; the rejected
Procfs, Common, client_model, client_golang, and earlier studies are neither
prepared nor run, and P8 remains queued.
