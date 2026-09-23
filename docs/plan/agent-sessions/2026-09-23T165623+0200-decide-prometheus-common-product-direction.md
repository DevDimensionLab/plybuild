# Agent Session: Decide Prometheus Common Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T165623+0200-decide-prometheus-common-product-direction`
Created: `2026-09-23T16:56:23+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `79d3f87016027b6acd87f61761431b02fa8318109effa8994b534e5b7d9e8d0d`
Previous: [2026-09-23T144751+0200-evaluate-prometheus-common-dependency.md](2026-09-23T144751+0200-evaluate-prometheus-common-dependency.md)
Next: [2026-09-23T172541+0200-evaluate-prometheus-procfs-dependency.md](2026-09-23T172541+0200-evaluate-prometheus-procfs-dependency.md)
Outcome: Option 1 is final. Exact selected, inherited, indirect, unloaded `github.com/prometheus/common v0.9.1` remains unchanged under a Common-specific, unqualified, non-transferable exception. It is neither qualified, supported, safe, nor fixed; no owner/request study or dependency change is authorized.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product decision for graph-selected
transitive exact `github.com/prometheus/common v0.9.1`. The completed fresh
evaluation found no fully qualified canonical stable. Choose and record exactly
one of the three authorized directions below. Do not repeat the completed
evaluation, implement a dependency change, combine another group, launch a
study or successor, or begin P8.

# Defensive Scope

This is an ordinary dependency product-direction decision. Use only the
completed static/release/build/graph/projection/advisory evidence and bounded
read-only guard checks. Do not fuzz, stress, probe resource exhaustion, create
oversized, deeply nested, cyclic, malformed, adversarial, or escape-sequence
payloads, reproduce a security issue, or perform security or exploitability
analysis.

Every disposable cache, report, project copy, or advisory response must remain
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`,
`/tmp`, a sibling of the managed root, or another external root. Set every
tool temp/cache root explicitly, verify containment, and remove task-owned
scratch evidence before handoff.

# Completed Evaluation

The exact public Prometheus owner has 88 canonical non-retracted exact-path
stable candidates after excluding three alternate-module v0 tags, retracted
v0.50.0 and the retracted accidental v1 line, plus one prerelease and no
replacement, deprecation, `/v2`, or `/v3` line. Exactly 51 stables preserve the
Go 1.18 floor through v0.44.0. Exact repository/tag/commit/tree/signature/
ancestry, proxy/sumdb/archive, module/license, source/build-tag, package, API,
and documented ordinary-behavior evidence is complete for all 51.

All 51 verify and build under exact Go 1.18.10 and Go 1.26.7. Selected v0.9.1
fails count-one/count-ten/race tests under both SDKs because its legacy
certificate fixtures and TLS expectations no longer match either supported
toolchain. Only v0.37.0-v0.44.0 are native-clean at the floor; v0.37.1 loses
that status under Go 1.26.7. Highest compatible v0.44.0 is native-clean under
both SDKs.

No release passes the complete bounded cgo-disabled cross matrix. All 51 fail
production build and test compilation on Plan 9 amd64 under both SDKs because
their exact go-conntrack closures reference unavailable
`syscall.ECONNREFUSED`; the other nine targets pass. Exact target-version OSV
and narrow GitHub responses are empty, but advisory absence is not
qualification. Selected v0.9.1's isolated closure has 25/3/0/1 non-stdlib
module/package/symbol/test-symbol findings. V0.44.0 has 20/6/5/5, and every
native-clean compatible release retains reachable closure findings. Thus no
canonical stable fully qualifies.

Exactly four target requests reproduce. TSDB v0.7.1 metadata-only requests
`v0.0.0-20181113130724-41aa239b4cce`; go-metrics v0.3.10 has a genuine test
import and requests v0.9.1; client_golang v1.4.0 has genuine production/test
imports and requests v0.9.1; client_golang v1.0.0 has genuine production/test
imports and requests v0.4.1. Their four complete shortest routes originate at
the main module through direct mvn-pom-mutator v0.2.3 and its TSDB or historical
Viper v1.10.1/go-metrics/Common/client_golang paths. Target/requester why is
negative; repository target import, target/requester production/complete-test
load, target module-backed load, runtime relevance, and current target roots are
zero. Historical target roots are nonzero. `go.sum` has three target go.mod
sums and no target source sum. No requester asks for v0.44.0.

A selected exact get only adds a redundant root, source sum, and main edge;
ordinary tidy restores the established common projection and v0.9.1. A
v0.44.0 get changes 15 existing selections, including protected go-conntrack,
client_golang, and client_model, and adds four modules. Ordinary tidy removes
the Common root and reselects v0.9.1 but retains six unrelated upgrades over
the common projection. V0.44.0 therefore has no genuine supported tidy-stable
owner and is not guard-preserving. Neither projection was retained.

The unchanged real project remains exactly 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, and 1,067 sum lines at the protected hashes and Go 1.18
floor. Project govulncheck remains 30/22/20/20 with no Common trace;
client_golang v1.4.0 retains GHSA-cg3q-j54f-5p7p / GO-2022-0322 /
CVE-2022-21698. Corrected Go-index/CNA identities, all 47 pre-Goe selections/
276 incoming edges, separate Goe/pkg-errors/SFTP/go-difflib/Complete/ULID/go-
conntrack/client_golang/client_model guards, every earlier decision, and
accepted 27/27 Q0-Q2 PASS at L2 remain exact.

# Authorized Decision

Choose exactly one:

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/prometheus/common v0.9.1` unchanged under a Common-specific,
   unqualified, non-transferable exception. State that it is not qualified,
   supported, safe, or fixed. Preserve all four exact requests, requester
   import/metadata boundaries, four routes, negative why, zero repository
   import/load/runtime/current-root state, nonzero historical-root state,
   release/source/API/behavior/closure/native/cross/projection/advisory
   identities, and every earlier guard as expiry conditions. Do not transfer
   the client_golang, client_model, or any earlier exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator Prometheus
   Common Owner/Request Study** to determine whether a supported update of the
   direct product root can eliminate every Common request or establish genuine
   tidy-stable ownership of a separately qualified future Go-1.18-compatible
   stable while preserving product behavior, direct roots, Go 1.18, and every
   earlier guard. Do not run the study, change a dependency, grant an exception,
   or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because the target is unloaded and runtime-irrelevant,
all qualification and ownership failures are explicit, and it preserves the
unchanged project without broadening a study or transferring an exception.
This is a product acceptance decision, not a qualification claim.

# Protected State And Checks

Start only from the clean Common evaluation handoff. Verify its HEAD, parent,
tree, exact changed set, branch/ancestry, reciprocal archive chain, sole NEXT
state, launcher mirror, ordinary and ignored cleanliness, official SDK
identities, real project counts and hashes, common tidy projection, Go 1.18
floor, all target request/import/route/relevance/root facts, target/closure/
project advisory identities, and every earlier guard. Stop for a fresh owning
decision if any protected input changed.

Preserve exact selected v0.9.1 and its four requests. Preserve client_model
v0.2.0/15 requests and client_golang v1.4.0/four requests only under their own
separate unqualified non-transferable exceptions; do not run either rejected
owner study, select v0.4.0 or v1.16.0, add a target root, or transfer an
exception. Preserve Complete, go-difflib, SFTP, pkg/errors, Goe, ULID, go-
conntrack, and all earlier qualified or excepted results under their own exact
guards. P8 remains queued.

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
`github.com/prometheus/common v0.9.1` is explicitly retained unchanged under a
Common-specific, unqualified, non-transferable exception. It is not qualified,
supported, safe, or fixed. Selected v0.9.1 fails the mandatory native,
complete-cross, and closure-advisory gates; highest native-clean Go-1.18-
compatible stable v0.44.0 also fails the complete cross and closure-advisory
gates, has no genuine supported tidy-stable project owner, and is not guard-
preserving.

Option 2's **Mvn-Pom-Mutator Prometheus Common Owner/Request Study** is not
authorized or run. Expanding this acceptance decision into a direct-product-
root study is not warranted while Common is unloaded and runtime-irrelevant
and all exact request and ownership failures are known. Option 3 is not
selected because the unchanged inherited state can be explicitly bounded
without making a qualification claim. No dependency or source change, target
root, ownership claim, implementation pre-authorization, transferred
exception, other dependency group, or P8 work is included.

### Exact exception and expiry boundary

Retention requires all four exact requests and requester boundaries to remain
unchanged:

- Prometheus TSDB v0.7.1 metadata-only requests
  `v0.0.0-20181113130724-41aa239b4cce` and has no Common source import;
- go-metrics v0.3.10 requests v0.9.1 and retains its genuine Common test
  import;
- client_golang v1.4.0 requests v0.9.1 and retains genuine Common production
  and test imports; and
- client_golang v1.0.0 requests v0.4.1 and retains genuine Common production
  and test imports.

All four complete shortest routes from the main module must remain exact:
main -> direct mvn-pom-mutator v0.2.3 -> TSDB v0.7.1 -> the Common pseudo-
version; main -> mvn-pom-mutator -> historical Viper v1.10.1 -> go-metrics
v0.3.10 -> Common v0.9.1; that same route through go-metrics ->
client_golang v1.4.0 -> Common v0.9.1; and that route through go-metrics ->
Common v0.9.1 -> client_golang v1.0.0 -> Common v0.4.1.

The exception additionally requires negative target and requester `go mod why
-m`; zero repository Common imports; zero target/requester production and
complete-test load; zero Common module-backed load and runtime relevance; zero
current main Common roots; the exact nonzero historical root boundary where
`73d4c0edfc375a3ad57d2f18a0dd02cb287a8640` added v0.4.0 and
`c25ad9273c5de46c0674b9b8ebcc87d94b0d886f` removed it; exactly three
Common go.mod sums and no source sum; and no request for v0.44.0. Any changed
request, requester import/metadata boundary, route, why/import/load/runtime/
root/sum fact, or newly supported owner route expires retention and requires a
fresh owning evaluation and explicit product decision before merge.

Every canonical-release and source identity recorded in the answered Common
evaluation remains an expiry guard: the exact public enabled, unarchived,
non-fork Apache-2.0 Prometheus owner on `main`; 88 canonical non-retracted
exact-path stables after excluding the three alternate-module v0 tags,
retracted v0.50.0 and accidental v1 line, and the one prerelease; absent
replacement, deprecation, `/v2`, and `/v3` lines; and exactly 51 Go-1.18-
compatible candidates through v0.44.0. Exact repository/tag/commit/tree/
signature/ancestry, proxy/sumdb/archive-to-Git, module/license, source/build-
tag, package, API, and documented ordinary-behavior evidence for all 51 must
remain unchanged.

Selected v0.9.1 remains signed annotated tag object
`cacacface9794af831b9694a3bc8970b7fddaffa`, commit
`d978bcb1309602d68bb4ba69cf3f8ed900e07308`, tree
`629ef783a44169b1bc150614186eef8a706f53e7`, with source/mod sums
`h1:KOMtN28tlbam3/7ZKEYKHhKoJZYYj3gMH4uc62x7X7U=` /
`h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8bs7vj7HSQ4=`. Candidate v0.44.0
remains signed lightweight commit
`94bf9828e56d9670579b28a9f78237d3cd8d0395`, tree
`775c1c9d87de4f8908e160ac70c87f984e609881`, with source/mod sums
`h1:+5BrQJwiBB9xsMygAB3TNvpQKOwlkc25LbISbrdOOfY=` /
`h1:ofAIvZbQ1e/nugmZGz4/qCb9Ap1VoSTIO7x0VV9VvuY=`. The exact archive,
license, package, build-boundary, API, and behavior identities in the
evaluation remain incorporated by reference as non-transferable expiry
conditions.

The complete closure/native/cross identities remain exact. Selected v0.9.1
has the same 34-module build list under both SDKs, exact 176/203 and 240/267
production/complete-test entries with 51 module-backed entries over 19 loaded
modules, and the two recorded API snapshot hashes. It continues to verify and
build but fail count-one, count-ten, and race tests under exact Go 1.18.10 and
Go 1.26.7 on its legacy certificate fixtures and TLS expectations. V0.44.0
retains its 44-module build list, 212/239 and 276/303 entry counts with 86
module-backed entries over 22 loaded modules, and its two API snapshot hashes;
it remains native-clean under both SDKs. Only v0.37.0-v0.44.0 remain native-
clean at the floor, with v0.37.1 losing that status under Go 1.26.7.

All 51 releases must retain the exact complete cgo-disabled cross result: both
Plan 9 amd64 production-build and test-compilation rows fail under both SDKs
because their go-conntrack closures reference unavailable
`syscall.ECONNREFUSED`, while the other nine targets pass. Exact empty target
OSV and narrow GitHub responses remain guards, and their absence is not
qualification. Pinned govulncheck v1.8.0 at database timestamp
2026-09-16T18:00:43Z must retain selected v0.9.1's 25/3/0/1 and v0.44.0's
20/6/5/5 non-stdlib module/package/symbol/test-symbol closure findings; every
native-clean candidate must retain reachable closure findings. Any newly
qualified release or supported guard-preserving ownership route expires this
exception rather than silently broadening it.

Both disposable projections remain expiry guards. Exact selected get has
75/1,068 lines, 234 modules, 3,600 edges, unchanged 355/429/197/41 loads, and
`go.mod`/`go.sum`/graph hashes
`844f6d9346b1d9d1cbc12edef1dcde28a4cc1288d29338723562e5cb294e5915` /
`0e286538b584b8f6f7ab78bac5c71b1e084d798b3a19030bff9bbc4bb86e9656` /
`2910bc331fde3a6e7525214d6210734c1cf25b1d83890302872cd0df74c2c059`;
tidy removes the redundant root, source sum, and edge and restores the common
projection. Exact v0.44.0 get has 75/1,075 lines, 238 modules, 3,627 edges,
the recorded 15 changed and four added modules, unchanged loads, and hashes
`18dde19586be2bfbde087d40db7aa12787709dce66517c3fa695f2b5b85c127c` /
`045d67274977523948f42cbd92f70b32bdc51d4319e5a2b94a4d9c86ec21bf49` /
`97fb7a2fcdeaf3e302b1df5ea978af0913c883f3099df8e3bf6196788e340940`.
Its tidy result reselects v0.9.1 but retains six unrelated upgrades in the
54/952-line, 234-module/3,564-edge state at the three recorded hashes. Neither
projection is retained.

This Common exception does not transfer the separate client_model v0.2.0/15-
request or client_golang v1.4.0/four-request unqualified, non-transferable
exceptions. It also does not transfer, reopen, or alter Complete, go-difflib,
SFTP, pkg/errors, Goe, ULID, go-conntrack, or any earlier qualified or excepted
result. Each remains governed only by its own exact selection/request/route/
relevance/source/closure/projection/advisory and expiry guards.

### Guard-only verification and handoff

- Work began clean on `codex/upgrade-quality` at Common evaluation handoff
  HEAD `2eb35fcdd1d4e1cbc7375acefe31df2470e08ea3`, parent
  `e170896d0afaa5541ea2c0080160be520a3bde0c`, tree
  `75a82918f7dc220e2f5ba7d816ab5b5badf73ae3`. It changes exactly the
  launcher, answered Common evaluation, this then-NEXT decision, rolling
  handover, and roadmap. Google UUID implementation
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor.
- The reciprocal archive chain, sole NEXT state, launcher/archive prompt
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
  355 production entries, 429 complete-test entries, 197 module-backed
  entries across 41 loaded modules, and 1,067 sum lines. `go.mod`, `go.sum`,
  and graph SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  The Go floor remains 1.18. Guard-only tidy in a contained project copy
  reaches the exact common 52/948-line, 234-module/3,557-edge projection at
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`
  and reselects v0.9.1. No projection is retained.
- All four exact requests, requester import/metadata boundaries, complete
  route edges, negative target/requester why, zero target repository import/
  load/runtime/current-root state, nonzero historical-root state, and target
  sum identities reproduce. The exact graph preserves all 47 pre-Goe
  selections and 276 incoming edges at
  `7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
  Separate selection/request guards reproduce: Goe v0.1.0/four, pkg/errors
  v0.9.1/ten, SFTP v1.13.1/four, go-difflib v1.0.0/23, Complete v1.2.3/three,
  ULID v1.3.1/one, go-conntrack `cc309e4a2223`/two, client_golang v1.4.0/four,
  and client_model v0.2.0/15. Every earlier qualified or excepted result
  remains closed, separate, and untransferred; accepted quality remains 27/27
  Q0-Q2 PASS at L2.
- Fresh narrow Common v0.9.1/v0.44.0 OSV and Common GitHub responses remain
  empty; client_golang v1.4.0 returns exactly GHSA-cg3q-j54f-5p7p /
  GO-2022-0322 / CVE-2022-21698. Pinned govulncheck v1.8.0 at database
  timestamp 2026-09-16T18:00:43Z reproduces the unchanged-project 30/22/20/20
  population with no Common trace. The Go index remains 518,501 bytes/1,402
  records at
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  the PUBLISHED CNA response remains 2,807 bytes at
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

Final unchanged-project exact-Go-1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass. Product source, `go.mod`, and
`go.sum` remain byte-exact. Every task-owned SDK, cache, response, report,
tool, and project copy was contained beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`
and removed before handoff.

P7 remains active only with one prepared bounded evaluation of the next
selected alphabetical module, exact `github.com/prometheus/procfs v0.0.8`.
Its four graph requests are starting observations only. The successor is
prepared but not launched; the rejected Common, client_model, and
client_golang studies are neither prepared nor run, and P8 remains queued.
