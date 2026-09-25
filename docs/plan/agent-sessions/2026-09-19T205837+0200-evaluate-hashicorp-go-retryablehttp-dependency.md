# Agent Session: Evaluate Hashicorp Go Retryablehttp Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-19T205837+0200-evaluate-hashicorp-go-retryablehttp-dependency`
Created: `2026-09-19T20:58:37+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `273852eb9544da73de632c9fb35e0ba73e7d4a322f002b7c27a46c66dc5bfe72`
Previous: [2026-09-19T203710+0200-decide-hashicorp-go-multierror-product-direction.md](2026-09-19T203710+0200-decide-hashicorp-go-multierror-product-direction.md)
Next: [2026-09-19T233605+0200-decide-hashicorp-go-retryablehttp-product-direction.md](2026-09-19T233605+0200-decide-hashicorp-go-retryablehttp-product-direction.md)
Outcome: No exact-path stable release qualifies: every Go-1.18-floor-eligible release through v0.7.6 is affected by GO-2024-2947/GHSA-v6v8-xj6m-xwqh, while first-fixed v0.7.7 declares Go 1.19 and changes the exported API and unrelated guarded selections; dependency metadata remains unchanged and P7 stops for the reciprocal product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/go-retryablehttp v0.5.3` as one bounded dependency
group. Resolve its complete repository and release identity, Go-floor closure,
package behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2, and
all recorded retained-module decisions through exact inherited, unloaded
go-multierror v1.1.0. Every earlier outcome and lifecycle ancestor is final.
Evaluate only Hashicorp go-retryablehttp in this session; do not reopen or
combine another dependency group. P8 remains queued.

The authorized 2026-09-19 go-multierror option 1 decision retains exact
selected, inherited, unloaded `github.com/hashicorp/go-multierror v1.1.0`
without dependency metadata changes. It accepts only the completed Prefix
`errors.Is`/`errors.As` identity loss, typed-nil `WrappedErrors` panic, and
characterized API, formatting/order, nil/panic, mutation, aliasing, allocation,
concurrency, Group lifecycle, resource/global-state, MVS, vulnerability, and
related findings. Its exception is go-multierror-specific and valid only while
exact v1.1.0 and all six requests from Serf v0.9.6, memberlist v0.1.3/v0.3.0,
mitchellh/cli v1.0.0/v1.1.0, and posener/complete v1.2.3 remain unchanged,
zero target packages load, the module remains runtime-unreachable, the
separate Errwrap v1.0.0 guard remains intact, and no new advisory or
independent defect appears. Direct import/loading, runtime reachability, a
target version or incoming-request change, an Errwrap-guard change, or a new
advisory or independent defect requires the owning fresh decision before
merge.

The go-msgpack v0.5.3, go-immutable-radix v1.3.1, and go-hclog v1.2.0
exceptions, qualified go-cleanhttp v0.5.2 result, and the Errwrap v1.0.0,
Consul SDK v0.8.0, Consul API v1.18.0, Gateway v1.16.0, gRPC Prometheus
v1.2.0, gRPC middleware v1.0.0, Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 decisions remain separate under their recorded exact-selection,
incoming-edge, zero-load, runtime-unreachable, and no-new-finding guards.
Revalidate those guards and stop for the owning decision if any expires. No
earlier exception transfers to go-retryablehttp. Do not change
`armon/go-metrics`, Viper, mvn-pom-mutator, another guarded parent, or an
unrelated module.

Selected `github.com/hashicorp/go-retryablehttp v0.5.3` is inherited through
one exact request from `github.com/armon/go-metrics v0.3.10`. Its current
`go mod why -m` result and repository import search are negative, and
production and complete-test loads contain zero target packages. This physical
MVS selection and bounded queue survey are not proof of repository identity,
release qualification, ancestry, floor, behavior, vulnerability state, or
suitability. Resolve those facts independently and do not add a direct edge
merely to alter MVS.

# Measurements At Start

The go-multierror decision recording began from clean handoff HEAD
`a6dba3eee16d9c3c622fdc68c075aa45595c27cc`, parent
`ea7bc704e1160eca51b3b35d00411db28babcac8`, tree
`99c8ac8c6c063ec11cc66232af33f69d48c7e054`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status,
reciprocal archive history, and `./codex-dev-start.sh --check` rather than
assuming them. The latest dependency implementation remains exact Google UUID
v1.4.0 commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`.

Guard-only revalidation under exact Go 1.26.7 preserved all 15 guarded module
selections and recorded incoming edges. All 15 `go mod why -m` results remain
negative, repository imports are zero, and production and complete-test loads
contain zero guarded packages. The project remains 234 modules, 3,599 graph
edges, 355 production entries, 429 complete-test entries, 197 module-backed
entries across 41 loaded modules, 1,067 sum lines, and the recorded 432-line
unapplied tidy projection. Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Exact go-multierror v1.1.0 and every
guarded target except Gorilla retain empty exact-version OSV results. Gorilla
retains only its recorded GO-2026-6278/GHSA-w67g-5rqw-f597 result and existing
GO-2020-0019 primary-index entry. No new guarded advisory or independent
defect appeared.

Use exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, use `LC_ALL=C LANG=C`, and
run with `umask 022`. Recreate contained Go 1.18.10 and pinned tools beneath
`$CODEX_SESSION_SCRATCH_ROOT` as required. Preserve the known apidiff archive
reproducibility discrepancy, Python 3.14 Docker timestamp control, managed
bare-`mktemp` restriction, and nested launcher signal-retention timing race;
none is go-retryablehttp evidence.

# Role And Boundaries

From fresh external archives and caches, resolve proxy, sumdb, `go-import`,
Git, and forge evidence for the exact module path: tags, releases, branches,
signatures, commits, parents, trees, times, ancestry, repository status,
licenses, retractions, deprecations, redirects, forks, alternate paths, major
module lines, and every serious exact-path stable candidate. Do not silently
promote a redirect, fork, alternate path, different major module path,
prerelease, non-versioning tag, unreleased branch head, or floor-ineligible
release.

Prove the complete minimal production and test closure under exact Go 1.26.7
and contained Go 1.18.10. Inspect imported source and test dependencies rather
than treating the module directive alone as floor proof. Separate isolated
source-time resolution from the project's selected MVS graph and from every
already-final guarded dependency decision.

Inspect every package, command, exported API, example, benchmark, fuzz target,
testdata, generated file, platform/build-tag branch, and applicable API and
runtime boundary. Characterize retry policy, backoff, redirects, bodies,
cancellation/deadlines, error identity, status handling, hooks, nil/panic
behavior, mutation and aliasing, allocation, concurrency, global state,
resource cleanup, malformed inputs, and actual project consumers. Add
independent fixtures where useful and run source verification, package listing,
native complete tests, two independent repeats, race, vet, and meaningful
cross-builds under both SDKs. Classify every failure precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, floor, architecture, or unrelated-
module change needs a fresh bounded decision rather than silent
implementation. Compare primary vulnerability results at module, package,
symbol, test-symbol, and reachable-trace levels.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, latest Google UUID implementation identity, reciprocal archive
history, P7/P8 state, every unchanged qualification/exception guard, and
`./codex-dev-start.sh --check`. Read this archive, the answered go-multierror
decision/evaluation, the answered Errwrap decision/evaluation, rolling
handover, roadmap, `go.mod`, `go.sum`, and every referenced quality,
compatibility, release, runner, evidence, and lifecycle contract. Earlier
outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose
complete minimal source/test closure preserves Go 1.18. Do not promote an
unqualified or floor-ineligible identity. If no candidate satisfies the
existing contracts, preserve the evidence and stop for a fresh bounded product
decision.

Second, qualify selected and serious candidates across source, tests, API,
behavior, loading, MVS, vulnerability, and project contracts. For an authorized
changed selection, use exact Go 1.26.7 and exact `go get` for one dependency-
only commit, never tidy as implementation, then run the complete P7 dependency
gate. For a retained or blocked selection, prove the no-change effect and run
all applicable gates without manufacturing activity. Full changed-selection
quality must preserve 27/27 Q0-Q2 PASS at L2 with zero held, regressed, non-
comparable, or dirty counts.

Third, update the roadmap and rolling handover with exact evidence, outcome,
commit identity, limitations, and next boundary. Answer this archive and
prepare one reciprocal NEXT mission only after the bounded outcome is coherent
and committed. Do not execute the successor.

# Automatic Handoff

After a completed coherent result, make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The bounded result is blocked: no stable exact-path
`github.com/hashicorp/go-retryablehttp` release satisfies the complete current
contracts. No `go get`, dependency implementation, `go.mod`, or `go.sum`
change was retained. P7 stops at the reciprocal product decision rather than
silently accepting a vulnerability, raising the Go floor, breaking API, or
changing another guarded module.

### Continuity And Identity

The evaluation began clean on branch `codex/upgrade-quality` at handoff HEAD
`24b1f33671362309526328170e67e68aaf753a3a`, parent
`a6dba3eee16d9c3c622fdc68c075aa45595c27cc`, tree
`d00fe77d08624cb227aa61641b403ad42c9d4deb`. The handoff changed exactly the
launcher, answered go-multierror decision archive, this then-NEXT archive,
rolling handover, and roadmap. Ordinary and ignored status were empty; the
reciprocal archive chain and `./codex-dev-start.sh --check` passed. The latest
dependency implementation remains Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`.

Fresh `go-import` metadata resolves the exact module path to
`https://github.com/hashicorp/go-retryablehttp.git`. GitHub identifies a
public, active, unarchived, non-disabled, non-fork MPL-2.0 repository with
default branch `main`; it has no GitHub Release objects. The proxy exposes 23
stable tags only, v0.5.0-v0.5.4, v0.6.0-v0.6.8, and v0.7.0-v0.7.8, with
v0.7.8 at `@latest`. There is no prerelease, retraction, module deprecation,
redirect, alternate module path, or `/v2` line. Current main
`fd004584...` is 17 commits beyond v0.7.8 and remains unreleased.

All relevant tags are in one ancestry and all proxy/Git file manifests agree
byte-for-byte. Tags are lightweight and unsigned except annotated, unsigned
v0.7.1. Strict Git verification passes. Selected v0.5.3 is commit
`357460732517ec3b57c05c51443296bdd6df1874`, parent
`73489d0...`, tree `8a3d77e...`, dated 2019-03-25. Last floor-eligible v0.7.5
is merge commit `4165cf8897205a879a06b20d1ed0a2a76fbb6a17`, tree
`b825...`; first-fixed v0.7.7 contains unsigned fix commit
`a99f07beb3c5faaa0a283617e6eb6bcf25f5049a`; latest v0.7.8 is commit
`e1f5485fe84728709b857cb89e17088894c301d6`, tree `c96...`.

Selected v0.5.3 has sumdb identity
`h1:QlWt0KvWT0lq8MFppF9tsJGF+ynG7ztc2KIPhzRGk7s=` and proxy-ZIP SHA-256
`1560044c4deed91fa2a27874216ed4580afbabd37f53232d2364b131c915d94f`.
V0.7.5, v0.7.7, and v0.7.8 sums are respectively
`h1:bJj+Pj19UZMIweq/iie+1u5YCdGrnxCT9yvm0e+Nd5M=`,
`h1:C8hUCYzor8PIfXHa4UrZkU4VvK8o9ISHxT2Q8+VepXU=`, and
`h1:ylXZWnqa7Lhqpk0L1P1LzDtGcCR0rPVUrx/c8Unxc48=`. Their proxy-ZIP SHA-256
values are `0782b9e874735dafc5c3a7f6c7d73f3c7022b25bbb34e8c135fa5f502e89a735`,
`29fa3655abd0a7c91f3a1eb186d309ec24b580afb10d39a3cdfa8cde0e99e0a9`,
and `9e2c175e4af37cdfeba28ed89250ad1fd691eda0fab167fafa9928222646b71e`.

### Go Floor, Source, Tests, And API

V0.5.0-v0.6.2 have no `go` directive, v0.6.3-v0.7.5 declare Go 1.13,
v0.7.6/v0.7.7 declare Go 1.19, and v0.7.8 declares Go 1.23. Therefore
v0.7.5 is the highest release that preserves the declared Go 1.18 floor;
v0.7.7 is the first secure release and is floor-ineligible. Go 1.18.10 can
compile empty test selections for v0.7.6-v0.7.8, but that observation does not
override their declared unsupported floors.

Selected v0.5.3 has a two-module minimal closure and no external production
dependency. V0.7.5's six-module graph imports the target, go-cleanhttp v0.5.2,
and go-hclog v1.2.0; its imported production/test closure declares at most Go
1.13. V0.7.6 and later add hclog v1.6.3, fatih/color v1.16.0, and x/sys
v0.20.0; only the target directive exceeds Go 1.18. Serious releases contain
one library package and no commands, examples, fuzz targets, testdata,
generated files, cgo, embeds, generators, or symlinks. V0.7 adds a round-
tripper; v0.7.6 adds two version/build-tag certificate branches.

Pinned API comparison finds every release after v0.5.3 incompatible because
`Client.Logger` changes from `Logger` to `interface{}`. V0.7.5 also adds
`StandardClient`, `SetBody`, response handlers, `WriteTo`,
`ErrorPropagatedRetryPolicy`, `FromRequest`, `LeveledLogger`,
`NewRequestWithContext`, and `RoundTripper`; v0.7.6 adds `PrepareRetry`, and
v0.7.8 adds `RateLimitLinearJitterBackoff`. No later release is compatible
with the selected exported API.

All 20 floor-eligible releases pass isolated upstream tests under contained Go
1.18.10. Under exact Go 1.26.7, v0.5.0-v0.6.2 pass tests with vet disabled;
v0.6.3-v0.7.5 fail the upstream TLS policy test because modern Go wraps the
certificate verification error and the policy relies on direct type/string
matching; v0.7.6-v0.7.8 pass. Selected v0.5.3 and v0.6.2 pass native
count-one, two independent count-two repeats, and race under both SDKs, but
vet fails in test code because `t.Fatalf` is called from non-test goroutines.
V0.7.5 passes the complete Go 1.18.10 battery yet repeats the focused TLS
failure under Go 1.26.7. First-fixed v0.7.7 passes exact-Go count-one, two
count-two repeats, race, vet, and five production/test cross-builds. Selected,
v0.6.2, and v0.7.5 pass Darwin AMD64, Linux AMD64/ARM64, Windows AMD64, and
js/wasm test compilation under both SDKs.

### Independent Behavior And Security

The independent behavior fixture SHA-256 is
`cf4ceef2094699c773cc5a44afb90ed992d4ac3ae222a3ee27ff79d95bea1a1c`.
Selected v0.5.3 and v0.6.2 fail permanent classification for x509 unknown-
authority, unsupported-scheme, and redirect-limit errors and lose
`errors.Is` identity for the final transport error. They pass body replay,
retry/final-response close ownership, cancellation identity and stop,
status/backoff, BodyBytes copy characterization, nil-panic characterization,
concurrent shared-client requests, and allocation logging under both SDKs.
V0.6.3 repairs direct x509 and redirect classification but not unsupported
schemes or final-cause identity. V0.7.5 passes this fixture but fails the
actual modern-Go TLS branch; selected-source body and logger aliasing,
package-global client use, retry-body draining, nil-response panic, and
PassthroughErrorHandler body ownership are fully characterized in the retained
evidence.

The independent credential-redaction fixture SHA-256 is
`2b15b5812d21c3f91f9279105923423f6fbb311568c56ec8d3651e91fd952bcb`.
Fresh primary vulnerability data contains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Reviewed GO-2024-2947 and
GHSA-v6v8-xj6m-xwqh affect every release through v0.7.6: Basic-auth
credentials embedded in a URL can be written to logs and returned errors from
client and package-level request paths. The fixture reproduces the exact
`https://alice:pasten@...` disclosure in v0.5.3, v0.6.2, v0.7.5, and
v0.7.6 under applicable SDKs; v0.7.7/v0.7.8 pass. Exact-version OSV results
for v0.5.0-v0.7.6 contain both IDs, while v0.7.7/v0.7.8 are empty.
Independent consumer `govulncheck` module/package/symbol/test-symbol scans
report GO-2024-2947 for v0.5.3 and v0.7.5 and are empty for v0.7.7/v0.7.8.

### Project MVS And Quality Effects

The target exists solely because `github.com/armon/go-metrics v0.3.10`
requests exact v0.5.3. The shortest root path is main ->
`mvn-pom-mutator v0.2.3` -> historical Viper v1.10.1 -> go-metrics v0.3.10
-> target. `go mod why -m` and repository import search are negative; the 355-
entry production and 429-entry complete-test loads contain zero target
packages, so selected code is runtime-unreachable.

The base remains 234 modules, 3,599 graph edges, 197 module-backed complete-
test entries across 41 loaded modules, and 1,067 sum lines. Disposable exact
`go get` projections preserve zero target load. V0.5.3, v0.6.2, and v0.7.5
move no unrelated selected version; v0.7.6/v0.7.7 also upgrade guarded
go-hclog v1.2.0 to v1.6.3 and fatih/color v1.15.0 to v1.16.0; v0.7.8 makes
the same moves and raises the main `go` directive to 1.23. All exact-Go
projection verify/build/count-one/vet gates pass only because the target is
unloaded. Base and v0.5.3/v0.6.2/v0.7.5 tidy projections are byte-identical;
v0.7.6/v0.7.7 retain the color upgrade, and v0.7.8 retains that upgrade plus
Go 1.23. No projection was applied.

All 15 earlier guarded selections and recorded incoming edges remain exact;
all 15 why results and repository import searches are negative, and both load
sets contain zero guarded packages. Fresh exact-version OSV revalidation is
empty for 14 guards; Gorilla WebSocket v1.4.2 retains only its recorded
GO-2026-6278/GHSA-w67g-5rqw-f597 result. No owning guard expired.

Exact Go 1.26.7 project verify, build, count-one, two count-ten repeats, race,
vet, pinned lint, API/CLI, empty-HOME, four cross-build, host/snapshot/Docker,
all 17 script/meta, all 80 mutation, and all 15 audit controls pass. Contained
Go 1.18.10 retains only the two accepted Darwin shell closed-file wording
failures in the full suite; 26 unaffected packages, vet, host acceptance, and
four cross-builds pass. The authoritative 21-stage scorecard has all 27 Q0-Q2
rows PASS at L2, zero held/regressed/non-comparable/dirty counts, SHA-256
`de13154181ae319fd80a29a98df78535762d70c2c01438f9e3cb735c2f624e81`.
The external manual evidence SHA-256 is
`a6d429808a82a05a7243be740635d5b511d39b489ca48c5a291601d093155206`;
the 547-entry complete-evidence manifest SHA-256 is
`12f36165ac7c0ddb4a69a5e432c3b12ea57c977a9e6d7c7a338a359d36b4faa2`.
Known cold offline API-cache, nested signal timing, bare-`mktemp`, and Python
3.9 timestamp attempts were superseded by the required scratch/cache/timing
and Python 3.14 controls and are not target findings.

Final `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Because all Go-1.18-floor-eligible releases are vulnerable and the first fixed
release requires an explicit floor/API/guarded-selection decision, the next
bounded action is the prepared go-retryablehttp product decision. P8 and every
other dependency group remain untouched.
