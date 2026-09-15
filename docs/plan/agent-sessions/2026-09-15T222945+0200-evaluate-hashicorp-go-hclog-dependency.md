# Agent Session: Evaluate Hashicorp Go Hclog Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-15T222945+0200-evaluate-hashicorp-go-hclog-dependency`
Created: `2026-09-15T22:29:45+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `6855cd1c9d9f4c0e4fa1461e249b89623b41c4bb41888d9bcdbca29521673a86`
Previous: [2026-09-15T212106+0200-evaluate-hashicorp-go-cleanhttp-dependency.md](2026-09-15T212106+0200-evaluate-hashicorp-go-cleanhttp-dependency.md)
Next: [2026-09-15T234309+0200-decide-hashicorp-go-hclog-product-direction.md](2026-09-15T234309+0200-decide-hashicorp-go-hclog-product-direction.md)
Outcome: No exact-path stable go-hclog release qualifies; preserved the unchanged inherited selection without accepting it, stopped P7 for a bounded product decision, and prepared that reciprocal decision session.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/go-hclog v1.2.0` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure, package
behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified exact go-cleanhttp
v0.5.2, and all recorded retained-module decisions through Errwrap v1.0.0.
All earlier outcomes and lifecycle ancestry are final. Evaluate only Hashicorp
go-hclog in this session; do not reopen or combine another dependency group.
P8 remains queued.

The qualified go-cleanhttp result retains exact selected, inherited, unloaded
`github.com/hashicorp/go-cleanhttp v0.5.2` without dependency metadata changes.
It is valid while exact v0.5.2 and all three selected-version incoming edges
from Viper v1.15.0, historical Viper v1.10.1, and
`sagikazarmark/crypt v0.4.0` remain unchanged, zero packages load, the module
remains runtime-unreachable, and no new advisory or independently
disqualifying behavior appears. Direct import/loading, runtime reachability, a
target version or incoming-edge change, or a new advisory or defect requires a
fresh go-cleanhttp dependency decision before merge. Revalidate these guards
before work and stop for that owning decision if any fails. Do not transfer
this qualification or any earlier exception to go-hclog.

The user's 2026-09-15 Errwrap option 1 decision retains exact selected,
inherited, unloaded `github.com/hashicorp/errwrap v1.0.0` without dependency
metadata changes. It accepts only the completed concrete-type collision,
absent standard single/multi-error traversal, nil/panic, aliasing, allocation,
recursion, API, and related qualification findings. Its exception is Errwrap-
specific and valid only while exact v1.0.0 and the sole selected-version
go-multierror v1.1.0 incoming edge remain unchanged, zero Errwrap packages
load, the module remains runtime-unreachable, and no new advisory or
independent defect appears. Direct import/loading, runtime reachability, a
target version or incoming-edge change, or a new advisory or independent
defect requires a fresh Errwrap dependency and product decision before merge.
Do not change its parent chain.

The Consul SDK v0.8.0, Consul API v1.18.0, Gateway v1.16.0, gRPC Prometheus
v1.2.0, gRPC middleware v1.0.0, Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 decisions remain final and separate under their exact-
selection, recorded incoming-edge, zero-load, runtime-unreachable, and no-new-
finding guards. Revalidate those guards and stop for the fresh owning decision
if one expires. Do not change either Gateway parent, mvn-pom-mutator,
GoConvey, Viper, the historical Consul API parent, go-multierror, Serf,
`sagikazarmark/crypt`, or transfer an earlier exception.

# Measurements At Start

The go-cleanhttp evaluation began from clean handoff HEAD
`0d1d66480604da5b646c75593570b01ebec72dda`, parent
`f77cb8a75c6813b28656c85b654d17d7df464731`, tree
`46e8716418edbafc7095bb22b194ed79f67ffb83`. It retained qualified exact
v0.5.2 without a dependency implementation or metadata commit. The latest
dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`.

Go-cleanhttp v0.5.2 is proxy `@latest`, its complete stdlib-only source/test
closure preserves Go 1.18, and it passes source, native/repeated/race/vet,
cross-build, API, behavior, and vulnerability qualification under exact Go
1.26.7 and contained Go 1.18.10. Its why result is negative, repository source
imports are zero, and production and complete-test loads contain zero target
packages. The exact selected version exists through the three incoming edges
listed above. Do not add a direct edge merely to alter MVS.

Guard-only revalidation under exact Go 1.26.7 preserved all eleven guarded or
qualified versions and recorded incoming edges. All eleven `go mod why -m`
results remain negative, repository Go source contains zero guarded-path
occurrences, and production and complete-test loads contain zero guarded
packages. The project remains 234 selected modules, 3,599 graph edges, 429
complete-test entries, 197 module-backed entries across 41 loaded modules,
1,067 sum lines, and the recorded 432-line unapplied tidy projection. Its
`go.mod`/`go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data has 1,399 records at SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
and Last-Modified 2026-09-15T18:54:35Z. It has no exact go-cleanhttp,
Errwrap, or other new guarded-target record; exact go-cleanhttp
v0.5.0/v0.5.1/v0.5.2 OSV responses are empty. Gorilla retains only
GO-2020-0019 and unwithdrawn GO-2026-6278 at record SHA-256
`bb7e1fb07bb85bab3169ecf4a108e25a61aeee50ed8175f7dd6f8134d9bfb228`.

A bounded queue survey identifies selected Hashicorp go-hclog v1.2.0 in the
current build list. That selection is not proof of repository identity, release
qualification, ancestry, floor, package loading, behavior, vulnerability
state, or suitability. Resolve those facts independently and do not add a
direct edge merely to alter MVS.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, and use
`LC_ALL=C LANG=C`. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`
and GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`.
Preserve the known apidiff archive reproducibility discrepancy, Python 3.14
Docker timestamp control, managed bare-mktemp restriction, and nested launcher
signal-retention timing race; none is Hashicorp go-hclog evidence.

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
state, resource cleanup, and malformed inputs. Distinguish tools, examples,
optional packages, and test-only helpers from behavior actually loaded by this
project.

Add independent fixtures where useful for deterministic behavior, error
identity, nil/panic boundaries, malformed inputs, mutation, aliasing, supported
concurrency, resource cleanup, and selected-project consumers. Run source
verification, package listing, native complete tests, two independent repeats,
race, vet, and meaningful cross-builds under both SDKs. Classify every failure
precisely.

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
ancestry, latest Google UUID dependency implementation identity, reciprocal
archive history, P7/P8 state, every unchanged qualification/exception guard,
and `./codex-dev-start.sh --check`. Read this archive, the answered
go-cleanhttp evaluation, the Errwrap decision/evaluation, the Consul SDK and
earlier guarded-decision archives, rolling handover, roadmap, `go.mod`,
`go.sum`, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose complete
minimal source/test closure preserves Go 1.18. Do not promote an unqualified or
floor-ineligible identity. If no candidate satisfies the existing contracts,
preserve the evidence and stop for a fresh bounded product decision.

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

No exact-path stable `github.com/hashicorp/go-hclog` release qualifies under
the current behavior and compatibility contracts. The selected v1.2.0 remains
physically selected only because the project graph is unchanged; it was not
qualified, retained by exception, or replaced. No dependency implementation or
metadata commit was created. P7 stops at a fresh bounded go-hclog product
decision.

The evaluation began from clean handoff HEAD
`452fffbd865dcae811fba3b02821cf7fd48f774c`, parent
`0d1d66480604da5b646c75593570b01ebec72dda`, tree
`4858e0dd3c216b772476058324a265aef0cb005a`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`. The worktree, branch,
ancestry, reciprocal archive chain, ignored status, P7/P8 boundary, and
launcher continuity contract were clean at entry.

Fresh `go-import` metadata maps the exact path to
`https://github.com/hashicorp/go-hclog.git`. GitHub reports a public, active,
unarchived, non-fork repository with protected default branch `main` and MIT
license. The proxy lists 31 valid stable releases, v0.7.0-v1.6.3, with no
prerelease in that set, retraction, module deprecation, redirect, or `/v2`
module. Git also exposes non-semver alias tags `v1.0` and `v0.17`, nested
`hclogvet` module tags, and branch `f-v2`; none is an exact-path target release.
Current unreleased main is commit
`e4c86b4cdbc417b598e03d5e4344ddff3419aea1`, parent
`7c1ccdc645c8fda4d57c2beb24567c89cf8d6b99`, tree
`55a509603684758007d981c67f816c33a16f6435`, dated
2026-08-25T15:03:33Z, declares Go 1.25, and was not promoted.

Selected v1.2.0 is lightweight unsigned tag/commit
`b6b55671f4e5b82443139ee3e9f4417603c4cd72`, parents
`81033451e6eb54da74139737d45a39b6412c7f34` and
`fc772a82149bb181261310935e431f89667288c3`, tree
`45f0da0a3b7eb001522fc9891504cee03d951ef0`, committed
2022-03-03T03:51:25Z. Its GitHub Release was published
2022-03-08T22:18:43Z. Proxy zip SHA-256 is
`3df5039c586056534758eecb3164c72c3b10e47383b2f669de2bb329baa78d79`;
sumdb source/mod sums are
`h1:La19f8d7WIlm4ogzNHB0JGqs5AUDAZ2UfCY4sJXcJdM=` /
`h1:whpDNt7SSdeAju8AWKIWsul05p54N/39EeqMAyrmvFQ=`.

Last API-compatible v1.3.1 is commit
`0d6179fa10233c02ec090700a92b25c4cdf60c45`, parents
`9846caeae374c1be1e9463cb22af4e2ae8e5087c` and
`33175ef984a644083bad9603634507d39c983d31`, tree
`51caa6ad3358f9d143e22cbb51ed53fa2c97af9e`, dated
2022-09-21T03:33:53Z. Its proxy zip SHA-256 is
`5c88941c557a5be4419d3a79dfaebf8830b7b5be7d1590ce6d9bb3fb186eb702`.
First API-breaking v1.4.0 is commit
`8b7499ad6ad46ce583e825d324d008032b644703`. Latest and proxy `@latest`
v1.6.3 is commit `d12136aa2e51933c460084f5083b6d5bb9d41960`, parents
`5dbb615f9aa8587fce14c2ab180aa7369f0ee703` and
`cb8687a9e2d8bab634ddff8412b3e03a7d60c068`, tree
`26c5c9247ad5b3a6f930ea2b1fc0a1ee532a7804`, committed
2024-04-01T20:03:54Z. Its proxy zip SHA-256 is
`ebab3136e5327ad17606485f633c2d033c61eadb843b2f3629b6d65d6fbb1400`;
sumdb source/mod sums are
`h1:Qr2kF+eVWjTiYmU7Y31tYlP1h0q/X3Nl3tPGdaB11/k=` /
`h1:W4Qnvbt70Wk/zYJryRzDRU/4r0kIg0PVHBcfoyhpF5M=`. Strict Git verification
passes, serious tags form the expected ancestry, and all inspected proxy
archives byte-match their Git trees after excluding the nested module.

All valid releases either omit a Go directive or declare Go 1.13. Selected
v1.2.0 resolves eight modules and 10/9 graph edges under Go 1.26.7/1.18.10,
with 76/55 production and 217/154 complete-test entries. Latest v1.6.3
resolves 11 modules and 20/19 edges, with 78/57 production and 220/157
complete-test entries. Imported production/test dependencies declare no more
than Go 1.17. The complete stable line therefore preserves Go 1.18; module
directives were not used as sole floor proof.

Each serious release has one root package, 25 regular module-archive files, 12
production Go files, seven tests, Unix/Windows color build branches, and one
benchmark. There are no commands, examples, fuzz targets, testdata, generated
files, cgo, embeds, `go:generate` directives, or symlinks. Selected v1.2.0,
v1.2.1, v1.2.2, v1.3.0, v1.3.1, v1.4.0, and latest v1.6.3 verify, build, pass
native count-one, count-ten repetition, race, and vet under exact Go 1.26.7 and
Go 1.18.10. Selected/latest additionally pass two independent count-ten runs
and production/test cross-builds for Darwin AMD64, Linux AMD64/ARM64, Windows
AMD64, and js/wasm under both SDKs. One first v1.2.0 run from a directory named
`candidate-v1.2.0` failed only its hard-coded caller-location expectation;
the canonical `go-hclog` basename passes, classifying this as an upstream test
path sensitivity rather than production behavior.

Pinned apidiff export SHA-256 values are
`17e8b795d8e83f4383b007c45ce24b80743531a0eb98e4fd5a48bb2e940bd77f`
for v1.2.0-v1.2.2,
`b102c0c649922fd9be13726340db9e369a2c3e16039d90999301dbe5f8ffc3a8`
for v1.3.0, and
`cb221e7cd0eb66564a5faa98f93332230cf62a3a13967deacac1f2cdc4eef159`
for v1.3.1. V1.3.0/v1.3.1 add only compatible
`LoggerOptions.ColorHeaderAndFields`. V1.4.0 export SHA-256
`01b5927d94c58d52da16ef7a5807fcd1728ebc78166886f6c93c4c4f63525aab`
adds `Logger.GetLevel` to the exported `Logger` interface, which is
incompatible for external implementations. Latest v1.6.3 export SHA-256
`79c2328c654f72dd74f2bceb3c6e723d135c7fe882fbfad29033e31e910fb221`
retains that break and otherwise adds JSON-escape, sublogger-hook,
sync-parent-level, and color-support API. V1.3.1 is the last API-compatible
candidate.

The common independent fixture has 183 lines and SHA-256
`b62af8023aa4e9d2b84a20db69627992998bfa29e7c268832025b345dcf5bb26`.
It passes characterization count-one, two count-ten repeats, race, and vet for
v1.2.0, v1.3.1, and v1.6.3 under both SDKs; v1.2.1, v1.2.2, and v1.3.0 also
pass characterization and directly reproduce the blockers. It verifies
plain/JSON formatting, sorted fields, odd keys, stdlib adaptation, level
routing, exact reset/flush error identity, caller-owned output cleanup,
nil/panic behavior, aliasing, allocations, and supported immutable concurrent
logging. `New` allocates four times in v1.2.0/v1.3.1 and five in v1.6.3;
`NewNullLogger` allocates zero under both SDKs.

Every selected/API-compatible candidate fails the same five independent
contract tests under both SDKs, and source inspection confirms the responsible
implementations remain through latest:

- JSON output silently drops the full record for NaN because recovery handles
  `json.UnsupportedTypeError` but not `json.UnsupportedValueError`.
- Caller fields overwrite core JSON metadata; the fixture emits exactly
  `{"@level":"trace","@message":"overwritten"}`.
- `DeregisterSink` decrements the counter even when the sink is absent, so a
  subsequent registered sink receives zero calls.
- `Accept` executes while the sink registry mutex is held, so self-deregistration
  reproducibly deadlocks at the fixture 200 ms boundary.
- `SetDefault(nil)` makes `Default` and `FromContext` nil despite
  `FromContext` documenting that the result is guaranteed non-nil.

The latest 201-line extension, SHA-256
`15e4bc71cde744a69ee0cecc18cea1ec6730acb595a87f11ea3de6d92d2ebd02`,
also finds the `SyncParentLevel` concurrent `SetLevel`/`GetLevel` race: Go
1.26.7 reports eight race warnings and Go 1.18.10 reports four, involving
`level`, `setEpoch`, and `ownEpoch`. Other recorded boundaries include
`FromContext(nil)` and nil-options `FromStandardLogger` panics, silent dropping
of unmatched `With(nil)`, `ImpliedArgs` slice aliasing, caller-map aliasing in
`LeveledWriter`, retaining the old output on flush error, exact reset error
identity, and caller ownership of output closure. The five common defects are
independent disqualifiers; the latest race is additional.

No higher stable fixes those defects. V1.3.1 therefore cannot qualify despite
compatible API, and v1.4.0+ add an independent interface break. Exact v1.1.0
or lower cannot be selected inside this group: disposable v1.1.0 downgrades
Viper v1.15.0 to v1.10.1, leaving 229 modules and 3,540 edges; v0.9.2 further
removes mvn-pom-mutator, leaving 182 modules and 2,555 edges. Those prohibited
parent changes were not applied.

Selected v1.2.0 has exactly one selected-version incoming edge, Viper v1.15.0
-> go-hclog v1.2.0. Historical Viper v1.10.1 and
`sagikazarmark/crypt v0.4.0` request v1.0.0; historical Consul API v1.12.0 and
Consul SDK v0.8.0 request v0.12.0. The shortest selected path is main -> Viper
v1.15.0 -> go-hclog v1.2.0. `go mod why -m` is negative, repository Go source
contains zero target occurrences, and production and complete-test loads
contain zero target packages, so the target is runtime-unreachable.

Disposable exact gets have these effects, none applied:

- v1.2.0 retains the selected version and 234 modules, changes edges
  3,599 -> 3,606 and sums 1,067 -> 1,069, and yields `go.mod`/`go.sum`
  SHA-256
  `5d9923ea2007e98df5a26ef15b0d5f32c4ad5f01c3fab90c99e7464650fa26d9` /
  `f18b2fae4bb3460a58b20b987d7fcbc6f7b0881ee2cb9cbad35526b0b64f9862`.
- v1.3.1, v1.4.0, and v1.6.3 move only the target, retain 234 modules, grow to
  3,610 edges and 1,071 sums, preserve 355 production and 429 complete-test
  entries, and load zero target packages. Their hash pairs are respectively
  `bb78ce6754b53e77e82c56ccff817f430de59fb08993f8065208f1fbd2f3e870` /
  `a69dacd7f750f62c56b854cf3e02caa84dcca3d4b803a97d14a040332fb399b7`,
  `d86146809cd99a9f4c43a256b788e56601725069889eb69e14948a6be6097812` /
  `e35d4c1b0631b8877adb32a1cefafa6bfe4433a6ecc4e85ee7a37983ae6ac4cf`,
  and
  `89658d9f015442ad53ea882a6f4bc9497bd0e13ab7cfb39cfbac287bd1a026c6` /
  `1ea7f4007288ad6e92a80bb886dd3e33ead6c2e0cc8fcfefeea4f2225f10f5cf`.

Every tidy copy removes the manufactured root and returns exactly to selected
v1.2.0 and base projection hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
Tidy was not used as implementation.

Fresh primary vulnerability data has 1,399 records at SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
and Last-Modified 2026-09-15T18:54:35Z. It has no go-hclog record; the GitHub
repository advisory list is empty; exact OSV queries for v1.2.0, v1.3.1,
v1.4.0, and v1.6.3 are empty. Selected v1.2.0 has zero Go 1.26.7
module/package/symbol/test-symbol findings. V1.3.1, v1.4.0, and v1.6.3 each
report one module-only GO-2026-5024 finding through
`x/sys v0.0.0-20220503163025-988cb79eb6c6`; the advisory affects only
`x/sys/windows.NewNTUnicodeString`, and scans contain no vulnerable-package,
called-symbol, go-hclog, or reachable target trace.

Under Go 1.18.10 selected source has 89 IDs and 90/99/117/265
module/package/symbol/test-symbol finding paths. Production symbol traces have
18 go-hclog frames across two old-standard-library IDs; tests have 108 target
frames across 30 IDs. Latest has one additional x/sys ID/path at every level
and otherwise the same target-trace counts. No advisory is assigned to
go-hclog. Normal Go 1.26.7 project scans remain 30 IDs and 30/52/74/74 Darwin
finding paths plus 75 Windows symbol paths, with zero target SBOM package,
symbol, test-symbol, or reachable-trace occurrence because no target package
loads. All eleven prior guarded targets retain their exact versions, recorded
selected-version edges, negative why results, zero repository Go occurrences,
zero production/test package loads, runtime unreachability, and no new target
advisory. Gorilla alone retains its two previously recorded IDs.

The project remains byte-for-byte unchanged at 234 selected modules, 3,599
graph edges, 355 production entries, 429 complete-test entries, 197
module-backed entries across 41 loaded modules, 1,067 sum lines, and the
recorded 432-line tidy projection. `go.mod` and `go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Exact Go 1.26.7 module verification, build, count-one, two count-ten repeats,
race, vet, pinned golangci-lint 2.12.2, API/CLI compatibility, empty-HOME
count-two, four production cross-builds, host acceptance, pinned GoReleaser
snapshot acceptance, all 17 script meta-tests, all eight 10-mutant populations
with 80/80 kills, and all 15 audit controls pass. Contained Go 1.18.10 loads
366 complete-test entries; 26 unaffected packages pass count-one, both
count-ten repeats, race, and vet, and all four production cross-builds pass.
The full suite retains only the two accepted Darwin `pkg/shell` closed-file
wording failures.

Full preflight twice passed its substantive stages before reproducing the
known nested launcher signal-retention timing race. An independent launcher
run passed all 62 controls. A scratch-only `mktemp` wrapper was required for
legacy meta-tests because the managed sandbox denied bare macOS temp creation.
These are pre-recorded harness boundaries, not go-hclog evidence. Exact Go
1.26.7 and Go 1.18.10 archive/binary receipts remain
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
and
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
Portable lint and GoReleaser receipts match their required hashes. Pinned
apidiff binary SHA-256 is
`71c0ecf63ff09cdf1a2d53d8df566baf91e8fa0185c6dee7a58d3fa4a4a11200`;
the known archive reproducibility discrepancy remains unchanged. The 1,328-
entry disposable evidence manifest SHA-256 is
`e6647301b469e2c915ca07bd1f2b0f7856e6b55df75c76a350943cc780783a9f`.

Because no source or dependency metadata changed, no changed-selection quality
run applies. The accepted 27/27 Q0-Q2 PASS at L2 remains authoritative. The
reciprocal successor is a decision-only go-hclog session offering guarded
v1.2.0 retention, a separately scoped remediation study, or a P7 block. It was
prepared but not executed.
