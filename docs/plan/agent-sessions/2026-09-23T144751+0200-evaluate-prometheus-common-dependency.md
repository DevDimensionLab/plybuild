# Agent Session: Evaluate Prometheus Common Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T144751+0200-evaluate-prometheus-common-dependency`
Created: `2026-09-23T14:47:51+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `336abee9d759aafbd0650b82d1e823c321642a56ba0533ad4dca3b80f23385ef`
Previous: [2026-09-23T141334+0200-decide-prometheus-client-model-product-direction.md](2026-09-23T141334+0200-decide-prometheus-client-model-product-direction.md)
Next: [2026-09-23T165623+0200-decide-prometheus-common-product-direction.md](2026-09-23T165623+0200-decide-prometheus-common-product-direction.md)
Outcome: No canonical stable fully qualifies. Selected v0.9.1 fails mandatory native, cross, and closure-advisory gates; highest Go-1.18-compatible stable v0.44.0 is native-clean but fails the Plan 9 closure rows, retains reachable closure findings, has no genuine supported tidy-stable owner, and crosses protected selections. No dependency or source change is retained; exactly one reciprocal product decision is prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one fresh, bounded evaluation of the next selected
alphabetical module, graph-selected transitive exact
`github.com/prometheus/common v0.9.1`. Resolve its canonical exact-path release
line and highest fully qualified Go-1.18-compatible stable from primary
evidence. Implement one exact dependency-only changed selection only if the
candidate, its complete minimal closure, a genuine supported tidy-stable
project owner, and every earlier guard remain exact. Do not combine another
dependency group or begin P8.

# Defensive Scope

This is an ordinary dependency-quality evaluation. Use public metadata,
static source and repository records, project graph/build commands, upstream
tests, and only small bounded ordinary fixtures required by documented
behavior. Do not fuzz, stress, probe resource exhaustion, create oversized,
deeply nested, cyclic, malformed, adversarial, or escape-sequence payloads,
reproduce a security issue, or perform security or exploitability analysis.

Every disposable cache, tool, archive, report, project copy, fixture, or
advisory response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`.
Never write to `/private/tmp`, `/tmp`, a sibling of the managed root, or
another external root. Set every tool temp/cache root explicitly, verify
containment, and remove task-owned scratch evidence before handoff.

# Authorized Roadmap And Closed Decisions

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
qualified no-selection-change modern-go/concurrent and modern-go/reflect2,
and all final target-specific decisions through Prometheus client_model. P8
remains queued.

Exact selected, inherited, indirect, unloaded
`github.com/prometheus/client_model v0.2.0` is retained only under its new
client_model-specific, unqualified, non-transferable exception. It is not
qualified, supported, safe, or fixed. Preserve all 15 requests and genuine
requester imports, all 24 complete shortest routes, negative target/requester
why, zero target repository import/load/runtime/root state, release/source/
behavior/API/closure/projection/advisory identities, and every expiry
condition. Do not run the rejected **Mvn-Pom-Mutator Client Model Owner/Request
Study**, select v0.4.0, add a target root, or transfer this exception.

Selected `github.com/prometheus/client_golang v1.4.0` separately remains only
under its own client_golang-specific, unqualified, non-transferable exception.
Do not run its rejected owner/request-removal study, select v1.16.0, add a
root, or transfer that exception. Preserve the separate Complete, go-difflib,
SFTP, pkg/errors, Goe, ULID, go-conntrack, and every earlier qualified or
excepted result under its own exact selection/request/route/relevance/source/
closure/projection/advisory/expiry guards. No exception transfers or reopens.

# Protected Starting State

The client_model product decision began clean on `codex/upgrade-quality` at
evaluation handoff HEAD `c339c4a32a42c81cc9be44209ffced75755e7faa`, parent
`f747850982e27aa784b64552c5ffedbdc04dc7c7`, tree
`5c37acd72962c82e9a9d5fc2d0bdfdc58ea6ddd5`. That handoff changes exactly
the launcher, answered client_model evaluation, then-NEXT decision, rolling
handover, and roadmap. Google UUID implementation commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. Verify the
new decision handoff, reciprocal archive chain, sole NEXT state, launcher
mirror, exact changed set, ancestry, and ordinary and ignored cleanliness.

The unchanged real project has 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, and 1,067 `go.sum` lines. `go.mod`, `go.sum`, and graph SHA-256 are
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Normal tidy has the established 432-line diff at
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`
and common 52/948-line, 234-module/3,557-edge projection with `go.mod` /
`go.sum` hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
No projection is retained. The module floor remains Go 1.18 and accepted
quality remains 27/27 Q0-Q2 PASS at L2.

Official Go 1.18.10 archive/binary SHA-256 remains
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`;
Go 1.26.7 remains
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Use only exact verified binaries for named SDK gates.

All 47 pre-Goe selections and 276 incoming edges remain exact at
`7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
Separate selections/request counts remain Goe v0.1.0/four, pkg/errors
v0.9.1/ten, SFTP v1.13.1/four, go-difflib v1.0.0/23, Complete v1.2.3/three,
ULID v1.3.1/one, go-conntrack `cc309e4a2223`/two, client_golang v1.4.0/four,
and client_model v0.2.0/15. Preserve every recorded why/import/load/runtime,
route, source, behavior, closure, qualification/exception, projection,
advisory, and expiry fact.

# Target Starting Observations

Read-only queue identification selects exact Prometheus Common v0.9.1
immediately after client_model. The current full graph has four target
requests: Prometheus TSDB v0.7.1 requests
`v0.0.0-20181113130724-41aa239b4cce`; go-metrics v0.3.10 and client_golang
v1.4.0 request v0.9.1; and client_golang v1.0.0 requests v0.4.1. Current
target why, repository import, production/complete-test load, runtime
relevance, and main root are zero; historical main-root state is nonzero.

These are starting observations only. Independently reproduce every exact
request, requester genuine-import or metadata-only boundary, complete current
and historical route from main, target/requester why, repository import,
production/complete-test/module-backed load, runtime relevance, and current/
historical root fact. Physical selection, a route, loading, or advisory
absence does not establish qualification or ownership.

The corrected Go vulnerability-index guard is exactly 518,501 bytes and 1,402
records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED CVE-2026-14362 CNA response is exactly 2,807 bytes at
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
The unchanged project advisory guard remains 30/22/20/20 module/package/
symbol/test-symbol findings. Client_golang v1.4.0 retains
GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698; client_model's eligible
responses remain empty without implying qualification. Revalidate target and
guard advisories narrowly; advisory absence cannot establish qualification.

# Evaluation Contract

Resolve the exact go-import owner, repository identity/status/license/default
branch, exact-path stable/prerelease/replacement/retraction/deprecation lines,
tags, commits, trees, signatures, ancestry, proxy/sumdb, archive-to-Git,
module directives, Go floors, packages, build tags, cgo, generated/embed
boundaries, exported API, and documented ordinary behavior for every eligible
stable. Exclude branches, pseudo-versions as candidates, prereleases,
replacements, alternate module paths, forks, patches, vendoring, and source
changes unless explicitly authorized; none is.

For selected and every eligible stable candidate, evaluate the complete
minimal upstream module closure under exact Go 1.18.10 and Go 1.26.7 with
module verification, build, complete count-one and repeated tests, race where
supported, vet, and a bounded supported cgo-disabled cross-build/test-
compilation matrix. Use only small ordinary deterministic fixtures where
upstream coverage does not settle documented behavior. Do not execute
unrelated network services or destructive paths. A release qualifies only if
its complete closure and ordinary public behavior qualify, it preserves the
Go 1.18 floor, and it has a genuine supported tidy-stable owner/request route
that preserves every earlier guard.

Use scratch-contained project copies only to measure the selected exact get
and at most the one highest otherwise-qualified candidate get, each followed
by normal tidy. Record the exact selection, closure, graph, imports/loads,
sums, Go-floor effect, genuine ownership, and every earlier guard. Do not add
or alter a target root in the real project except for the one exact final
implementation authorized below.

# Decision And Implementation Boundary

If and only if one canonical exact-path stable is the highest fully qualified
Go-1.18-compatible candidate, has a genuine supported tidy-stable project
owner, and its complete projection preserves every guard, implement exactly
`go get github.com/prometheus/common@<selected>` followed by ordinary tidy and
the full project gates. Retain only the exact measured minimal dependency
closure. Do not add a direct root, change source, patch, vendor, fork, alter an
unrelated requester, select an alternate path, or combine another group.

If no stable qualifies, no supported tidy-stable owner exists, a projection
crosses any guard, or any protected starting input changed, retain no
dependency or source change. Record the completed result and prepare exactly
one reciprocal product decision for this target. Do not silently retain an
unqualified target, grant or transfer an exception, authorize an owner study,
or move to another dependency in the same turn.

# Required Reading And Checks

Read this archive; the answered client_model product decision and evaluation;
the answered client_golang product decision and fresh evaluation; the answered
Complete decision/evaluation; rolling handover; P7/P8 roadmap; `go.mod`; and
`go.sum`. Verify branch/HEAD/parent/tree, exact changed set, ancestry,
reciprocal archive chain, launcher check, ordinary and ignored cleanliness,
exact SDK identities, project hashes/counts/tidy state, target requests/routes/
why/import/load/runtime/root facts, all earlier guards, and narrow advisory
identities. Stop for a fresh owning decision if any protected input changed.

Update the roadmap and rolling handover, answer this archive, prepare at most
one reciprocal product decision required by the result, verify containment
and cleanup, run final exact-Go project gates, and make the local handoff
commit.

# Automatic Handoff

Do not launch a successor or study; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, transfer or reopen an exception,
change another dependency group, write outside the managed scratch root, or
begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No dependency or source change is authorized. There is no fully qualified
canonical stable, and highest Go-1.18-compatible stable v0.44.0 additionally
lacks a genuine supported tidy-stable project owner and a guard-preserving
projection. The real project therefore remains on exact inherited indirect
`github.com/prometheus/common v0.9.1`; this evaluation grants no exception and
makes no safety or support claim.

### Canonical release and source line

- Exact go-import metadata resolves `github.com/prometheus/common` to the
  public, enabled, unarchived, non-fork Prometheus repository, GitHub ID
  37326770, Apache-2.0, default branch `main`. Evaluation-time HEAD is
  `042f3705950c873b84153dd8ebd95c6ea4c47ae5`, tree
  `36afceb2dc16bf51f455f82a3dfbb914437f4b8d`.
- The proxy exposes 92 v0 stable tags, one prerelease v0.9.0-pre1, and 36
  accidental v1 stable tags. V0.9.2-v0.9.4 declare the alternate
  `github.com/prometheus/client_golang` module, v0.50.0 is retracted, and the
  accidental v1 line is retracted by the exact-path v1.20.99 sentinel (the
  other v1 tags are also alternate-path artifacts). There are no replacements
  or deprecations and no `/v2` or `/v3` release line. Excluding mismatched,
  prerelease, and retracted entries leaves 88 canonical non-retracted stable
  candidates.
- Exactly 51 canonical candidates preserve the Go 1.18 floor, from v0.1.0
  through highest v0.44.0; v0.45.0 starts at Go 1.20. All 51 proxy archives
  verify against sumdb and match their exact Git commits byte-for-byte after
  the module proxy's required nested-module exclusions. They contain no
  symlink, submodule, or special-file boundary and share exact Apache-2.0
  license hash
  `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`.
- Of the 51 eligible tags, 37 are signed annotated tags, 14 are lightweight
  commit tags, and 13 of those commits are signed; only v0.4.0 is unsigned.
  Fifty are `main` ancestors. Official v0.37.1 is the sole release-branch
  exception. GitHub independently reports valid signatures for selected and
  candidate identities.
- Selected v0.9.1 is signed annotated tag object
  `cacacface9794af831b9694a3bc8970b7fddaffa`, commit
  `d978bcb1309602d68bb4ba69cf3f8ed900e07308`, tree
  `629ef783a44169b1bc150614186eef8a706f53e7`. Its source/mod sums are
  `h1:KOMtN28tlbam3/7ZKEYKHhKoJZYYj3gMH4uc62x7X7U=` /
  `h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8bs7vj7HSQ4=` and its proxy zip SHA-256
  is `51b17bde0d6eafecb33c6f248d239882b8ca4388b8fbbb565bcb04eeb437f999`.
- V0.44.0 is signed lightweight commit
  `94bf9828e56d9670579b28a9f78237d3cd8d0395`, tree
  `775c1c9d87de4f8908e160ac70c87f984e609881`. Its source/mod sums are
  `h1:+5BrQJwiBB9xsMygAB3TNvpQKOwlkc25LbISbrdOOfY=` /
  `h1:ofAIvZbQ1e/nugmZGz4/qCb9Ap1VoSTIO7x0VV9VvuY=` and its proxy zip SHA-256
  is `a89fdf5f749cf97c576c753f9cda6f05586376843706dcf1f6c0715b58d11cc6`.
- Every eligible release is pure Go with no cgo, `go:embed`, or generated-file
  boundary. Selected v0.9.1 has 51 Go files, 10 production package directories,
  22 test files, and seven raw build constraints. V0.44.0 has 56 Go files,
  nine production package directories, 24 tests, and four constraints. The
  ordinary public surface covers HTTP/TLS/OAuth configuration, Prometheus
  exposition encode/decode, metric/label/time/duration models, logging,
  routing, server helpers, and version reporting. Upstream tests exercise
  those documented behaviors; no extra fixture or unrelated service was
  needed.

### Complete closure, API, native, cross, and advisory results

- Selected v0.9.1 has an identical 34-module exact build list under both SDKs.
  Exact Go 1.18.10 gives 176 production and 203 complete-test entries, 51
  module-backed entries across 19 loaded modules; Go 1.26.7 gives 240/267 with
  the same 51/19 boundary. Its exported API/documentation snapshots are 1,251
  lines at
  `6131922e81b50eb698167668d6498a254be2e40f669b4f67d87fc3662f7d9320`
  and 1,249 lines at
  `26480b41d41532cc28727cf2cae3aef17a484951309c897a68f32b03cad945e4`.
- V0.44.0 has an identical 44-module build list under both SDKs. Exact Go
  1.18.10 gives 212 production and 239 complete-test entries, 86 module-backed
  entries across 22 loaded modules; Go 1.26.7 gives 276/303 with the same
  86/22 boundary. Its API/documentation snapshots are 1,449 lines at
  `562e51f1ecb9d5c81f1c5c0acad9b9c3ee201377ff86bcbea22db1f1e9859242`
  and 1,452 lines at
  `35ba6a2e60b4bf9e6477946fe5a6cbc9417bd60fa5cd4a461c66411bd650b28c`.
- The complete native matrix covers all 51 releases under exact Go 1.18.10
  and Go 1.26.7: 1,020 download, verification, build, count-one/count-ten,
  race, vet, enumeration, and API result cells. Every module download,
  verification, and build passes. Selected v0.9.1 fails count-one, repeated,
  and race tests under both SDKs because its legacy SHA1-RSA fixture chain and
  exact TLS configuration/error expectations no longer match either supported
  toolchain; vet passes. At the floor, only v0.37.0-v0.44.0 pass all native
  gates. Under Go 1.26.7, v0.37.1 has test/repeat failures while v0.37.0 and
  v0.38.0-v0.44.0 pass. Thus v0.44.0 is the highest native- and ordinary-
  behavior-clean compatible stable, not a fully qualified release.
- The bounded cgo-disabled matrix contains 1,020 release/SDK/target pairs and
  2,040 production-build/test-compilation results across Darwin amd64/arm64,
  Linux amd64/arm64/386, Windows amd64/386, FreeBSD amd64, Plan 9 amd64, and
  js/wasm. All 918 non-Plan-9 pairs pass. Every release fails both Plan 9 rows
  under both SDKs because its exact `github.com/mwitkow/go-conntrack` closure
  references unavailable `syscall.ECONNREFUSED`. Selected and v0.44.0 both
  retain that exact failure, so no stable passes the complete cross boundary.
- Exact-version OSV responses for all 51 eligible Common releases are empty;
  selected/candidate GitHub global and repository advisory responses are also
  empty. Absence did not establish qualification. Pinned govulncheck v1.8.0,
  built by the verified Go 1.26.7 binary against database timestamp
  2026-09-16T18:00:43Z, finds 25/3/0/1 non-stdlib module/package/symbol/test-
  symbol findings in selected v0.9.1's isolated closure and 20/6/5/5 in
  v0.44.0's. Every native-clean release retains reachable closure findings;
  v0.44.0's symbol set is GO-2024-2687, GO-2025-3503, GO-2026-4918,
  GO-2026-5026, and GO-2026-5970. No exploitability claim is made.

### Requests, routes, relevance, and projections

- The full graph has exactly four target requests: TSDB v0.7.1 requests
  `v0.0.0-20181113130724-41aa239b4cce`; go-metrics v0.3.10 and client_golang
  v1.4.0 request v0.9.1; client_golang v1.0.0 requests v0.4.1. TSDB has no
  Common source import and is metadata-only. Go-metrics imports `expfmt` in
  one test. Both client_golang vertices have genuine production and test
  imports of `model` and `expfmt`.
- The four complete shortest routes are: main -> direct mvn-pom-mutator v0.2.3
  -> TSDB v0.7.1 -> the target pseudo-version; main -> mvn-pom-mutator ->
  historical Viper v1.10.1 -> go-metrics v0.3.10 -> target v0.9.1; that same
  route through go-metrics -> client_golang v1.4.0 -> target v0.9.1; and that
  route through go-metrics -> Common v0.9.1 -> client_golang v1.0.0 -> target
  v0.4.1.
- Target and all three requester-module `go mod why -m` results are negative.
  Repository target imports, target/requester production and complete-test
  loads, target module-backed load, runtime relevance, and current main target
  roots are zero. History is nonzero: commit
  `73d4c0edfc375a3ad57d2f18a0dd02cb287a8640` added direct Common v0.4.0 and
  `c25ad9273c5de46c0674b9b8ebcc87d94b0d886f` removed it. `go.sum` retains
  exactly the pseudo-version, v0.4.1, and v0.9.1 go.mod sums and no target
  source sum. No requester asks for v0.44.0.
- Exact selected get adds only a redundant indirect root, source sum, and main
  edge. Raw state is 75/1,068 lines, 234 modules, 3,600 edges, unchanged
  355/429/197/41 loads, and negative why, with `go.mod` / `go.sum` / graph
  SHA-256
  `844f6d9346b1d9d1cbc12edef1dcde28a4cc1288d29338723562e5cb294e5915` /
  `0e286538b584b8f6f7ab78bac5c71b1e084d798b3a19030bff9bbc4bb86e9656` /
  `2910bc331fde3a6e7525214d6210734c1cf25b1d83890302872cd0df74c2c059`.
  Ordinary tidy removes the root/source sum/edge and restores the established
  52/948-line, 234-module/3,557-edge common projection at the protected
  go.mod/go.sum hashes.
- Exact v0.44.0 get yields 75/1,075 lines, 238 modules, 3,627 edges, unchanged
  loads, negative why, and hashes
  `18dde19586be2bfbde087d40db7aa12787709dce66517c3fa695f2b5b85c127c` /
  `045d67274977523948f42cbd92f70b32bdc51d4319e5a2b94a4d9c86ec21bf49` /
  `97fb7a2fcdeaf3e302b1df5ea978af0913c883f3099df8e3bf6196788e340940`.
  Its graph changes 15 existing selections, including protected go-conntrack,
  client_golang, and client_model, and adds four modules. Ordinary tidy removes
  the manufactured Common root and reselects v0.9.1, but retains six unrelated
  upgrades over the common projection. That 54/952-line, 234-module/3,564-edge
  state has hashes
  `86411e2b86ef1aaa5b5c63746ca3e8fdc5d374988dc5e04e3e14200be2c5db21` /
  `a7f9e21fd90ae1250fb0ac63c21a5aadc1623e36c084fee1a78b3e6c8d41ca9c` /
  `bf0935ecae88fe1cda71e1ebb3942f1a3d5add09fcaade8dcad2077f64f624d1`.
  The candidate is neither selected, tidy-stable, owned, nor guard-preserving.
  Neither projection is retained.

### Guards, decision, and handoff

- Fresh normal/no-cache primary responses reproduce the corrected 518,501-
  byte, 1,402-record Go index at
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
  and PUBLISHED 2,807-byte CVE-2026-14362 CNA response at
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
  The unchanged real project reproduces exactly 30/22/20/20 non-stdlib
  module/package/symbol/test-symbol findings with no Common trace.
  Client_golang v1.4.0 retains GHSA-cg3q-j54f-5p7p / GO-2022-0322 /
  CVE-2022-21698.
- Work began clean on `codex/upgrade-quality` at HEAD
  `e170896d0afaa5541ea2c0080160be520a3bde0c`, parent
  `c339c4a32a42c81cc9be44209ffced75755e7faa`, tree
  `8168466de3c5882e9ec274459bb45ae61305ff50`. It changes exactly the launcher,
  answered client_model decision, this then-NEXT evaluation, rolling handover,
  and roadmap. Google UUID implementation `cf53bc64e...` remains ancestral.
  The reciprocal archive chain, sole NEXT state, launcher mirror/check,
  branch, ancestry, and ordinary/ignored cleanliness all reproduced.
- Both official SDK archive/binary identities reproduced exactly. The real
  project remains 234 modules, 3,599 graph edges, 355 production entries, 429
  complete-test entries, 197 module-backed entries across 41 loaded modules,
  and 1,067 sum lines at the protected module/graph hashes and Go 1.18 floor.
  The exact full graph preserves all 47 pre-Goe selections/276 incoming edges
  and the separate Goe/pkg-errors/SFTP/go-difflib/Complete/ULID/go-conntrack/
  client_golang/client_model selections and request counts. Every earlier
  qualified or excepted result and accepted 27/27 Q0-Q2 PASS at L2 remains
  separate, exact, and untransferred.
- Because no stable fully qualifies, v0.44.0 has no supported tidy-stable owner,
  and its projection crosses protected guards, the implementation boundary
  forbids `go get` in the real project. Exactly one reciprocal product decision
  is prepared among Common-specific unqualified retention of v0.9.1, one later
  measurement-only owner/request study, or stopping P7 unresolved. It is not
  executed here. No exception, study, successor, other dependency group, or P8
  work begins.
- Final unchanged-project exact-Go-1.26.7 module verification, build, count-one
  tests, race count-one tests, and vet pass after the documentation handoff.
  Every task-owned SDK, repository, archive, cache, response, report, project
  copy, tool, and fixture was contained beneath the managed scratch root and
  removed; only launcher-owned scratch state remains.
