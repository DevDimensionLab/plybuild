# Agent Session: Decide Prometheus Client Model Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T141334+0200-decide-prometheus-client-model-product-direction`
Created: `2026-09-23T14:13:34+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `6298a1f49efab8a20c2e729889d0e8c846578d18f4b3f52d7fefe5da5c1e6837`
Previous: [2026-09-23T131751+0200-evaluate-prometheus-client-model-dependency.md](2026-09-23T131751+0200-evaluate-prometheus-client-model-dependency.md)
Next: [2026-09-23T144751+0200-evaluate-prometheus-common-dependency.md](2026-09-23T144751+0200-evaluate-prometheus-common-dependency.md)
Outcome: Option 1 is final. Exact selected, inherited, indirect, unloaded `github.com/prometheus/client_model v0.2.0` remains unchanged under a client_model-specific, unqualified, non-transferable exception. It is neither qualified, supported, safe, nor fixed; the owner/request study and a dependency change were not authorized.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product decision for graph-selected
transitive exact `github.com/prometheus/client_model v0.2.0`. The completed
fresh evaluation found selected v0.2.0 unqualified because its exact declared
closure does not compile, while highest otherwise-qualified Go-1.18-compatible
stable v0.4.0 has no genuine supported tidy-stable project owner. Choose and
record exactly one of the three authorized directions below. Do not repeat the
completed evaluation, implement a dependency change, combine another group,
launch a study or successor, or begin P8.

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

The exact public Prometheus owner has nine canonical exact-path stable releases
v0.1.0 through v0.6.3, no prereleases, replacements, retractions, or
deprecations, and no `/v2` or `/v3` release line. Exactly v0.1.0 through v0.4.0
preserve the Go 1.18 floor. Exact tag, commit, tree, signature, ancestry,
proxy/sumdb, archive-to-Git, module, license, source, generated/build boundary,
exported API, and documented ordinary behavior evidence is complete.

Selected v0.2.0 verifies but fails build, count-one/count-ten tests, race, vet,
ordinary behavior compilation, and all 40 supported cgo-disabled cross rows
under exact Go 1.18.10 and Go 1.26.7: generated metrics.pb.go requires
`proto.ProtoPackageIsVersion3`, but its own exact go.mod selects
`github.com/golang/protobuf v1.2.0`, where that symbol does not exist. V0.1.0,
v0.3.0, and v0.4.0 pass complete native, repeated, race, vet, behavior, and
cross gates under both SDKs. Thus v0.4.0 is the highest otherwise-qualified
Go-1.18-compatible stable.

Exactly 15 requests and all genuine requester import boundaries reproduce:
go-metrics v0.3.10, client_golang v1.4.0, and Common v0.9.1 request v0.2.0;
client_golang v1.0.0 requests `fd36f4220a90`; Common v0.4.1 and TSDB v0.7.1
request `5c3871d89910`; nine historical go-control-plane vertices request
`14fe0d1b01d4`. All 24 complete shortest routes originate at the main module
and run through direct mvn-pom-mutator, its historical Viper/go-metrics/Common,
gRPC/cloud/genproto/client_golang/TSDB chains, or the recorded Google Martian
and Afero entries into those chains. Target/requester why is negative;
repository target import, production/complete-test/module-backed load, runtime
relevance, and current/history target-root state are zero. No requester asks
for v0.4.0.

A selected get only adds a redundant target root, source sum, and main edge;
tidy restores the common projection. A v0.4.0 get changes client_model v0.2.0
to v0.4.0 and google.golang.org/protobuf v1.28.1 to v1.30.0; ordinary tidy
removes the root and both changes and reselects v0.2.0. V0.4.0 therefore has
no genuine supported tidy-stable project owner. No projection, source,
`go.mod`, or `go.sum` change was retained.

Narrow client_model advisory responses are empty; advisory absence did not
establish qualification. The isolated v0.4.0 closure has only module-level
GO-2024-2611 against protobuf v1.30.0 and zero package/symbol findings. The
unchanged project remains exactly 30/22/20/20 module/package/symbol/test-symbol
findings with no target trace. Client_golang v1.4.0 still retains
GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698. The corrected Go-index
and CVE-2026-14362 CNA identities remain exact.

# Authorized Decision

Choose exactly one:

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/prometheus/client_model v0.2.0` unchanged under a client_model-
   specific, unqualified, non-transferable exception. State that it is not
   qualified, supported, safe, or fixed. Preserve all 15 exact requests,
   genuine requester import boundaries, 24 routes, negative why, zero target
   repository import/load/runtime/root state, exact release/source/behavior/
   closure/projection/advisory identities, and every earlier guard as expiry
   conditions. Do not transfer the client_golang or any earlier exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator Client Model
   Owner/Request Study** to determine whether a supported update of the direct
   root can establish genuine tidy-stable ownership of qualified v0.4.0 or
   eliminate all client_model requests while preserving product behavior,
   direct roots, Go 1.18, and every earlier guard. Do not run the study, change
   a dependency, grant an exception, or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because the target is unloaded and runtime-irrelevant,
the exact failure and absence of a supported owner are now explicit, and it
preserves the unchanged project without broadening a study or transferring an
exception. This is a product acceptance decision, not a qualification claim.

# Protected State And Checks

Start only from the clean evaluation handoff. Verify its HEAD, parent, tree,
exact changed set, branch/ancestry, reciprocal archive chain, sole NEXT state,
launcher mirror, ordinary and ignored cleanliness, official SDK identities,
real 234-module/3,599-edge/355-production/429-complete-test/197-module-backed/
41-loaded/1,067-sum-line state, module/graph/common-tidy hashes, Go 1.18 floor,
all target requests/routes/relevance facts, advisory identities, and every
earlier guard. Stop for a fresh owning decision if any protected input changed.

All 47 pre-Goe selections and 276 incoming edges remain exact. Preserve the
separate Complete, go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack,
client_golang, and every earlier qualified or excepted result under its own
exact guards. Selected client_golang v1.4.0 remains only under its separate
unqualified non-transferable exception; do not run its rejected owner study,
select v1.16.0, add a root, or transfer that exception. Accepted quality
remains 27/27 Q0-Q2 PASS at L2. P8 remains queued.

Do only guard-level revalidation; do not repeat completed release, source,
behavior, API, closure, fixture, projection, native, cross, or advisory matrix
work. Record the chosen product direction, update the roadmap and rolling
handover, answer this archive, verify containment and cleanup, run final exact-
Go-1.26.7 project module verification, build, count-one tests, race count-one
tests, and vet, and make the local handoff commit.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, transfer or reopen an exception,
change a dependency or source file, combine another dependency group, or begin
P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is final. Exact selected, inherited, indirect, unloaded
`github.com/prometheus/client_model v0.2.0` is explicitly retained unchanged
under a client_model-specific, unqualified, non-transferable exception. It is
not qualified, supported, safe, or fixed. Its exact declared closure does not
compile: generated `metrics.pb.go` requires `proto.ProtoPackageIsVersion3`,
while the release's own `go.mod` selects `github.com/golang/protobuf v1.2.0`,
where that symbol is absent. Highest otherwise-qualified Go-1.18-compatible
stable v0.4.0 has no genuine supported tidy-stable project owner.

Option 2's **Mvn-Pom-Mutator Client Model Owner/Request Study** is not
authorized or run. Expanding this acceptance decision into a study of the
direct product root is not warranted while the target is unloaded and runtime-
irrelevant and every exact owner/request boundary is known. Option 3 is not
selected because the unchanged inherited state can be explicitly bounded
without making a qualification claim. No dependency or source change, target
root, ownership claim, implementation pre-authorization, transferred
exception, or P8 work is included.

### Exact exception and expiry boundary

Retention requires all 15 exact requests and their genuine source-import
boundaries to remain unchanged:

- go-metrics v0.3.10, client_golang v1.4.0, and Prometheus Common v0.9.1
  request v0.2.0;
- client_golang v1.0.0 requests
  `v0.0.0-20190129233127-fd36f4220a90`;
- Prometheus Common v0.4.1 and TSDB v0.7.1 request
  `v0.0.0-20180712105110-5c3871d89910`; and
- go-control-plane v0.10.1, v0.9.0,
  `v0.9.1-0.20191026205805-5f8ba28d4473`, v0.9.4, v0.9.7,
  `v0.9.9-0.20201210154907-fd9021fe5dad`,
  `v0.9.9-0.20210217033140-668b12f5399d`,
  `v0.9.9-0.20210512163311-63b5d3c536b0`, and
  `v0.9.10-0.20210907150352-cf90f659a021` request
  `v0.0.0-20190812154241-14fe0d1b01d4`.

Go-metrics and TSDB retain test imports; both client_golang and both Common
vertices retain production and test imports; all nine go-control-plane
vertices retain generated production imports. Every one of the 24 recorded
complete shortest routes from the main module must remain exact. Those routes
enter through direct mvn-pom-mutator v0.2.3 and its historical Viper/go-
metrics/Common, gRPC/CNCF xDS/UDPA/genproto/cloud/API/Firestore/GAX,
client_golang, or TSDB chains, apart from the recorded Google Martian and
Afero entries into the same historical gRPC/cloud chains. Every route ends at
one of the exact requester vertices and requests above.

The exception additionally requires negative target and requester `go mod
why -m`; zero repository target imports; zero target/requester production and
complete-test load; zero target module-backed load and runtime relevance; zero
current and historical main target roots; four target go.mod sums and no
selected target source sum; and no request for v0.4.0. Any changed request,
requester import, route, why/import/load/runtime/root/sum fact, or newly
supported owner route expires retention and requires a fresh owning evaluation
and explicit product decision before merge.

The exact public enabled, unarchived, non-fork Apache-2.0 Prometheus owner,
`master` branch, nine canonical exact-path stables v0.1.0 through v0.6.3,
absent prerelease/replacement/retraction/deprecation and `/v2`/`/v3` lines,
and four Go-1.18-compatible candidates through v0.4.0 remain guards. Eligible
tag commit/tree identities remain:

- v0.1.0 `d1d2010b5beead3fa1c5f271a5cf626e40b3ad6e` /
  `1560b2bb69ddc589807ab7243377a72fe3e5bb35`;
- v0.2.0 `7bc5445566f0fe75b15de23e6b93886e982d7bf9` /
  `b607c60177c234f9010729ee2d6f019f34e7ae4b`;
- v0.3.0 `63fb9822ca3ba7a4ba5184071fb8f2ea000a99ef` /
  `67f55441b06e9a732c2843b71518a1d4aa09e68f`; and
- v0.4.0 `91c3945f2cfbfb9040e34a0b6764d804b5a5a490` /
  `5521419e66ced766970708248f7161937cd1e30f`.

Exact signatures, ancestry, proxy/sumdb/archive-to-Git equivalence, module and
license identities, the sole generated `github.com/prometheus/client_model/go`
production package, absent upstream-test/cgo/build-constraint/embed boundaries,
exported API, and documented ordinary behavior all remain expiry guards. The
complete closure facts also remain exact: v0.1.0 and v0.2.0 select protobuf
v1.2.0 plus X/Sync; v0.3.0 selects protobuf v1.3.5; and v0.4.0 selects protobuf
v1.5.0, google.golang.org/protobuf v1.30.0, go-cmp v0.5.5, and X/Xerrors.
V0.1.0, v0.3.0, and v0.4.0 must continue to pass their complete native,
repeated, race, vet, behavior, and 120 supported cross rows under exact Go
1.18.10 and Go 1.26.7; selected v0.2.0 must retain its exact compile defect and
40 failed cross rows. A newly qualified or newly supported compatible route
expires this exception rather than silently upgrading its claim.

Both disposable projection identities remain expiry guards. The selected get
is the exact 75/1,068-line, 234-module/3,600-edge state at `go.mod`/`go.sum`/
graph hashes
`ba927aeaee03ecbb3828731f68b398160132e76237388a0ea35774553eb79434` /
`a819601278ec695d13014b39549e6fc0f6d79773750bba9fb3c50a1c66ddb9b1` /
`a874a905f5798d87066da2ed6972f4ea4e9dbad8ad963d484e62b3b3b0e13a00`;
tidy removes its redundant root, source sum, and main edge. The v0.4.0 get is
the exact 75/1,069-line, 234-module/3,601-edge state at
`24bd7ae26f272e6ba7aa81f06b60d96199a511497c65e16401b7cc6b012291bc` /
`a1f86062cc269c863d871d4b4df9ecbe25666ddeb57c544c32222c2a315b2867` /
`560733563bdafdb471e8bb36f3b36ffd3a0eb23b1c3a21c28825ac5f593d7382`;
tidy removes the root and the client_model/protobuf changes and reselects
v0.2.0. Neither projection is retained.

Exact empty target OSV and GitHub advisory responses remain guards; their
absence is not qualification. The isolated v0.4.0 closure must retain only
module-level GO-2024-2611 against protobuf v1.30.0 and zero package, symbol,
or test-symbol findings. The unchanged project remains 30/22/20/20 non-stdlib
module/package/symbol/test-symbol findings with no client_model trace.
Client_golang v1.4.0 separately retains GHSA-cg3q-j54f-5p7p /
GO-2022-0322 / CVE-2022-21698 under its own unqualified non-transferable
exception. No client_golang or earlier exception transfers to client_model.

### Guard-only verification and handoff

- Work began clean on `codex/upgrade-quality` at evaluation handoff HEAD
  `c339c4a32a42c81cc9be44209ffced75755e7faa`, parent
  `f747850982e27aa784b64552c5ffedbdc04dc7c7`, tree
  `5c37acd72962c82e9a9d5fc2d0bdfdc58ea6ddd5`. It changes exactly the
  launcher, answered client_model evaluation, this then-NEXT decision,
  rolling handover, and roadmap. Google UUID implementation
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor.
- The reciprocal archive chain, sole NEXT state, launcher/archive prompt
  mirror and check, branch, exact changed set, and ordinary and ignored
  cleanliness reproduce. Completed release, source, behavior, API, closure,
  fixture, projection, native, cross, and advisory-matrix work was not
  repeated.
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
  and reselects v0.2.0. No projection is retained.
- All 15 exact requests, every requester import boundary, all 24 complete
  shortest routes, negative target/requester why, and zero target repository
  import/load/runtime/root state reproduce. The exact graph and module
  identities preserve all 47 pre-Goe selections and 276 incoming edges at
  `7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
  Separate selection/request guards reproduce: Goe v0.1.0/four, pkg/errors
  v0.9.1/ten, SFTP v1.13.1/four, go-difflib v1.0.0/23, Complete v1.2.3/three,
  ULID v1.3.1/one, go-conntrack `cc309e4a2223`/two, and client_golang
  v1.4.0/four. Every earlier qualified or excepted result remains closed,
  separate, and untransferred; accepted quality remains 27/27 Q0-Q2 PASS at
  L2.
- Fresh narrow target OSV responses remain empty and client_golang v1.4.0
  returns exactly GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698.
  Pinned govulncheck v1.8.0 at database timestamp 2026-09-16T18:00:43Z
  reproduces the unchanged-project 30/22/20/20 population with no target
  trace. The Go index remains 518,501 bytes/1,402 records at
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  the PUBLISHED CNA response remains 2,807 bytes at
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

Final unchanged-project exact-Go-1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass. Product source, `go.mod`, and
`go.sum` remain byte-exact. Every task-owned SDK, cache, response, report,
tool, and project copy was contained beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`
and removed before handoff.

P7 remains active only with one prepared bounded evaluation of the next
selected alphabetical module, exact `github.com/prometheus/common v0.9.1`.
Its current four requests and zero current why/import/load/runtime/root facts
are starting observations only. The successor was prepared but not launched;
no Common result is authorized here, the rejected client_model study was not
prepared or run, and P8 remains queued.
