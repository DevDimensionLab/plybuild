# Agent Session: Evaluate Hashicorp Go Multierror Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-19T190917+0200-evaluate-hashicorp-go-multierror-dependency`
Created: `2026-09-19T19:09:17+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `54b8add00636fc62471731095c1d4b6431baa1cdb20321df0f421b0cfc6fbd20`
Previous: [2026-09-19T185300+0200-decide-hashicorp-go-msgpack-product-direction.md](2026-09-19T185300+0200-decide-hashicorp-go-msgpack-product-direction.md)
Next: [2026-09-19T203710+0200-decide-hashicorp-go-multierror-product-direction.md](2026-09-19T203710+0200-decide-hashicorp-go-multierror-product-direction.md)
Outcome: No exact-path stable go-multierror release qualifies; preserved the unchanged inherited selection without accepting it, stopped P7 for a bounded product decision, and prepared that reciprocal decision session.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/go-multierror v1.1.0` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure, package
behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified exact go-cleanhttp
v0.5.2, and all recorded retained-module decisions through exact inherited,
unloaded go-msgpack v0.5.3. Every earlier outcome and lifecycle ancestor is
final. Evaluate only Hashicorp go-multierror in this session; do not reopen or
combine another dependency group. P8 remains queued.

The authorized 2026-09-19 go-msgpack option 1 decision retains exact selected,
inherited, unloaded `github.com/hashicorp/go-msgpack v0.5.3` without dependency
metadata changes. It accepts only the completed map-update panic, pointer-
`false` omission, malformed-build-comment/vet failure, upstream early-test-
flag failure, and characterized API, nil/panic, aliasing, allocation,
concurrency, global-state, resource, deterministic, malformed-input, error-
identity, RPC-cleanup, vulnerability, and related findings. Its exception is
go-msgpack-specific and valid only while exact v0.5.3 and all four requests
from memberlist v0.1.3/v0.3.0 and Serf v0.8.2/v0.9.6 remain unchanged, zero
target packages load, the module remains runtime-unreachable, and no new
advisory or independent defect appears. Direct import/loading, runtime
reachability, a target version or incoming-request change, or a new advisory or
independent defect requires a fresh go-msgpack dependency and product decision
before merge.

The go-immutable-radix v1.3.1 and go-hclog v1.2.0 exceptions, qualified
go-cleanhttp v0.5.2 result, and the Errwrap v1.0.0, Consul SDK v0.8.0, Consul
API v1.18.0, Gateway v1.16.0, gRPC Prometheus v1.2.0, gRPC middleware v1.0.0,
Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 decisions remain separate under their recorded exact-selection,
incoming-edge, zero-load, runtime-unreachable, and no-new-finding guards.
Revalidate those guards and stop for the owning decision if any expires. No
earlier exception transfers to go-multierror. Do not change memberlist, Serf,
Viper, either Gateway parent, mvn-pom-mutator, GoConvey, historical Consul API,
or another guarded parent.

Selected `github.com/hashicorp/go-multierror v1.1.0` is currently inherited.
Serf v0.9.6 requests v1.1.0, while memberlist v0.1.3/v0.3.0, mitchellh/cli
v1.0.0/v1.1.0, and posener/complete v1.2.3 request v1.0.0. This physical MVS
selection is not proof of repository identity, release qualification,
ancestry, floor, package loading, behavior, vulnerability state, or
suitability. Resolve those facts independently and do not add a direct edge
merely to alter MVS.

# Measurements At Start

The go-msgpack decision recording began from clean handoff HEAD
`69017090043510b4754e8e81feab762a8011fe90`, parent
`13f63a2c7639e373b261ca74176955b2be6461b2`, tree
`0db64c05ce2166f40821aeb296e10576124ed349`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status, and
reciprocal archive history rather than assuming their values. The latest
dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`.

Guard-only revalidation under exact Go 1.26.7 preserved all 14 guarded module
selections and recorded incoming edges. All 14 `go mod why -m` results remain
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
and Last-Modified 2026-09-17T17:29:18Z. Exact go-msgpack v0.5.3 and every
guarded target except Gorilla have empty exact-version OSV results. Gorilla
retains only its recorded GO-2026-6278/GHSA-w67g-5rqw-f597 exact-version
result and existing GO-2020-0019 primary-index entry. No new guarded advisory
or independent defect appeared.

Use exact Go 1.26.7, verify binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, use `LC_ALL=C LANG=C`, and
run with `umask 022`. Recreate contained Go 1.18.10 and pinned tools beneath
`$CODEX_SESSION_SCRATCH_ROOT` as required. Preserve the known apidiff archive
reproducibility discrepancy, Python 3.14 Docker timestamp control, managed
bare-`mktemp` restriction, and nested launcher signal-retention timing race;
none is go-multierror evidence.

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
testdata, generated file, platform or build-tag branch, and applicable API and
runtime boundary. Characterize deterministic behavior, errors and identity,
nil/panic behavior, mutation and aliasing, allocation, concurrency, global
state, resource cleanup, malformed inputs, and project consumers. Add
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
`./codex-dev-start.sh --check`. Read this archive, the answered go-msgpack
decision/evaluation, go-immutable-radix and go-hclog decisions, rolling
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

No exact-path stable `github.com/hashicorp/go-multierror` release satisfies
the existing qualification contract. The repository and dependency metadata
remain unchanged; selected v1.1.0 is preserved only as the current physical
MVS selection, not accepted as qualified. P7 stops at a fresh bounded product
decision.

Fresh `go-import`, proxy, sumdb, Git, and GitHub evidence resolves the exact
module path to the public, active, unarchived, non-fork MPL-2.0 repository
`https://github.com/hashicorp/go-multierror.git`. The proxy exposes only
v1.0.0, v1.1.0, and `@latest` v1.1.1. There are no GitHub Release objects,
retractions, module deprecations, redirects, alternate exact-path module
lines, or prereleases. Current main is unreleased commit
`6d4d48630db25c3c83fa83ecd41dd8438b82963c`, 73 commits after v1.1.1,
and was not promoted.

All three tags are lightweight and linearly ancestral. V1.0.0 is commit
`886a7fbe3eb1c874d46f623bfa70af45f425b3d1`, parent
`aa5ed92cc00e3d41e4252f1049f8541daf9ad2ee`, tree
`5a0c0242d6303fea0383e022be19ef47cfdbe83f`. Selected v1.1.0 is commit
`2004d9dba6b07a5b8d133209244f376680f9d472`, parent
`a5e98147459549195562243aaf8caf5cedbb4864`, tree
`fcc486d753842d176b67384af19ca7c51157c43c`. Latest v1.1.1 is merge
commit `9974e9ec57696378079ecc3accd3d6f29401b3a0`, parents
`ab6846acb1cf402787cde041fcff5f01c0b8e362` and
`0023bb0ab1225a10509dca08627b42bed25d74a5`, tree
`93989e7fff83edc30d172b8a10d3dcb9c67b8cc1`. GitHub verifies the v1.1
commits; its v1.0.0 signature is `unknown_key`. The tags themselves have no
signed tag objects. Proxy and tagged-Git file manifests agree byte-for-byte.

Proxy ZIP SHA-256 values for v1.0.0/v1.1.0/v1.1.1 are respectively
`a66a1b9dff26a9a7fcaa5aa5e658c13f94c0daeb572536b1ecc7ebe51f4d0be7`,
`3c60b77bb039c4431734883c65d645436a57c2653d00ae4cd7ea4f0303bc38e2`,
and `972cd841ee51fdeac69c5a301e57f8ea27aebf15fddd7f621d5c240f28c3000c`.
Sumdb source/module sums match the proxy. V1.0.0 has no `go` directive;
v1.1.0 declares Go 1.14 and v1.1.1 declares Go 1.13. Every release requires
only exact `github.com/hashicorp/errwrap v1.0.0`.

Each release contains one package and no command, example, benchmark, modern
fuzz target, testdata, generated file, platform/build-tag branch, cgo, embed,
generator, or symlink. V1.0.0 has six production and six test files; v1.1.x
has seven of each. The complete imported closure is target plus Errwrap and
standard library: 62 production/125 complete-test entries under Go 1.26.7 and
42/80 under Go 1.18.10. No imported module declares a floor above Go 1.14.

All three releases pass module verification, native count-one, two independent
count-ten repeats, race, vet, and source-integrity checks under exact Go 1.26.7
and contained Go 1.18.10. Production and test compilation pass for Darwin
AMD64, Linux AMD64/ARM64, Windows AMD64, and js/wasm under both SDKs. Native
and cross ledgers have SHA-256
`c11d08f718625ac42fe06a7b25abd83d00512b278eef908e9ce6873bb8717164` /
`d69b2b816d239eefdfa4bd4fa1451c27be340dfa118e7a076e0a6ba1c99b685b`.

V1.0.0's exported API snapshot SHA-256 is
`57bc475905f8d183cc704a630527648517ff9eec4bc77b2df73449ef783fe2d3`;
v1.1.0/v1.1.1 are
`d42649d2cd42d2918d681681995837a84b9e48a72b3bc2012ee2f76b981bf3c3` /
`61144850ca974bd048f95288c4ec12d45216647593ca0962a260bcd6efcfe4cc`.
V1.1.0 compatibly adds `Group` and `(*Error).Unwrap`; v1.1.0 and v1.1.1
have no API diff. V1.0.0 therefore is not a current-API-compatible downgrade.
The rebuilt pinned apidiff binary retained the known nonportable archive/binary
reproducibility discrepancy; tool module/version identity was exact and the
discrepancy is not target evidence.

The independent passing behavior fixture SHA-256 is
`e5fd467b66758b1d2c79d3a5b02f8eaea0c28e4c688e234e13d45fd1c87806c3`.
Every release passes its applicable fixture at count one, two independent
count-ten repeats, race, and vet under both SDKs. Formatting and slice order
are deterministic. `Append` filters nil and typed-nil members, flattens one
level, and mutates/returns its first `*Error`; `WrappedErrors` aliases the
exported `Errors` slice. `Flatten` recursively produces an independent slice
while preserving constituent error identity. `Prefix` mutates a multierror in
place. `ErrorOrNil` is nil-safe and allocation-free for an empty error.

Nil `Error.Error`, `GoString`, `Flatten` on a typed nil, and sorting a nil
error or nil member panic. V1.0.0 and selected v1.1.0 also panic when
`WrappedErrors` is called through a typed-nil receiver; v1.1.1 fixes only that
panic. Immutable Error reads are race-free. Multi-member `Unwrap` allocates a
shallow chain snapshot. `Group` waits for all work, collects non-nil failures
concurrently, and exposes completion-order rather than deterministic launch
order. The package owns no global mutable state or external resource.

No stable release qualifies because `Prefix` delegates to Errwrap v1.0.0 and
destroys standard error identity. The focused blocker fixture SHA-256 is
`1620577fe5db5a1e76e50d68498165924d54513848b8093f94c94518d87872ac`;
its four-run ledger SHA-256 is
`ff2ab517da88f32f44ee132cacefa436ee2da0b2ef4dbeec55ab2ee6e016ba36`.
Under both SDKs, v1.1.0 and v1.1.1 fail the same three assertions:
`errors.Is` cannot traverse Prefix on a plain error or multierror member, and
`errors.As` cannot traverse a prefixed multierror member. V1.0.0 uses the same
Prefix design. This conflicts with the package's stated standard Is/As/Unwrap
compatibility. Upstream issue 56 and unmerged PR 155 head
`8224dca4d7fdce17091ec99615804e283e17bf13` corroborate the `%w` fix, but
that commit is unreleased and cannot be promoted.

Selected v1.1.0 exists through six exact requests: Serf v0.9.6 requests
v1.1.0; memberlist v0.1.3/v0.3.0, mitchellh/cli v1.0.0/v1.1.0, and
posener/complete v1.2.3 request v1.0.0. The shortest path is main ->
mvn-pom-mutator v0.2.3 -> historical Viper v1.10.1 -> Serf v0.9.6 ->
go-multierror. `go mod why -m` is negative, repository Go imports are zero,
and production and complete-test loads contain zero target packages, so the
target is runtime-unreachable.

Disposable exact v1.1.0/v1.1.1 gets retain all 234 module selections, every
unrelated version, the 355/429 production/complete-test entries, and zero
target load. V1.1.0 manufactures indirect main-module roots for both target
and Errwrap, two graph edges, and two source checksums without changing any
selection; its after-get `go.mod`/`go.sum` SHA-256 values are
`1eb3c764803cfb6c79f0af5e45e582153e35c0572caf7ca8889faba2b11a5b6b` /
`0efad21b342fc156968e82381659b8131a54ea12d36f6bf41d81deaaf1849a29`.
V1.1.1 makes the same roots, moves only the target selection, and adds the
v1.1.1-to-Errwrap graph edge; its hashes are
`1139aa03b205c0e7b20ea784e454e8a087fab217c25b6b99fdb042a93e61eef6` /
`3e9cbab8551f4ab69bc43bf959cf25c5b6bbeb0934f2c43a88c37083c2c04de9`.
That new graph edge and root expire the separate Errwrap guard.

Both viable projections pass project verify, build, count-one, two count-ten
repeats, race, and vet. Tidy removes their manufactured roots and returns both
copies byte-identically to base projection hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`,
selecting v1.1.0 again. Exact v1.0.0 instead removes mvn-pom-mutator v0.2.3,
broadly rewrites the graph, and makes the project unloadable because
`github.com/devdimensionlab/mvn-pom-mutator/pkg/pom` has no provider. No
projection was applied.

Fresh primary vulnerability data contains 1,402 records, SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z, with no target record. Exact OSV queries
for v1.0.0-v1.1.1 are empty. Exact govulncheck v1.8.0 module, package,
symbol, and test-symbol scans under Go 1.26.7 have zero findings. Go 1.18.10
reports only that old SDK's 89 module, seven package, and two reachable
standard-library findings; no advisory is assigned to Multierror. Base and
v1.1.1 project scan outputs are byte-identical at 30/22/20/20 findings with
zero target advisory or trace.

All 14 earlier guarded modules retain their exact selections and recorded
incoming edges, negative why results, zero repository imports, zero production
or complete-test target loads, runtime unreachability, and prior advisory
state. Gorilla alone retains its recorded GO-2026-6278 exact-version result
and GO-2020-0019 primary-index entry. The guard ledger SHA-256 is
`37c915bc95f9d24184e6155c454d7ef6952d36b36f400fc6708a21614f84a95b`.

The unchanged project remains 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, 1,067 sum lines, and a 432-line unapplied native tidy diff at SHA-256
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`.
Base `go.mod`/`go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Exact Go 1.26.7 verify, build, count-one, two independent count-ten repeats,
race, vet, pinned golangci-lint 2.12.2, API/CLI compatibility, empty-HOME
count-two, and four production cross-builds pass. Full preflight, all 17
script/meta stages, 80/80 mutation kills, host, pinned-snapshot, and daemon-
backed Docker acceptance, and all 15 quality-audit controls pass. The
authoritative scorecard is 27/27 Q0-Q2 PASS at L2 with zero held, regressed,
non-comparable, or dirty counts; scorecard SHA-256 is
`f440896c323544637a5841a2115eb8cbedf94c020d29e665742bbcac20548207`.
Initial failures reproduced only the known scratch-cache, cold offline API,
managed bare-`mktemp`, and nested launcher signal-retention boundaries;
scratch-contained canonical reruns passed, including three independent
62-control launcher runs. None is target evidence.

Exact Go 1.26.7 archive/binary SHA-256 remains
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
Go 1.18.10 remains
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
Portable golangci-lint 2.12.2 and GoReleaser 2.17.1 receipts match their
required hashes. The 469-entry selected-evidence manifest SHA-256 is
`7db323255b0c8361371237d1e3e7405e697e3ad7f53167678f756535dc66b7f3`.

No dependency implementation or metadata commit exists. The sole successor
is the bounded go-multierror product-decision session offering guarded
selected-v1.1.0 retention, an explicit v1.1.1 exception/upgrade that also
requires a fresh Errwrap decision, or a separately scoped remediation/block.
It was prepared but not executed.
