# Agent Session: Evaluate Hashicorp Go Sockaddr Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T121350+0200-evaluate-hashicorp-go-sockaddr-dependency`
Created: `2026-09-20T12:13:50+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ed1874643357c9bcc2ca24e6ee16069a5e57c47638597f4a265d38625a10c984`
Previous: [2026-09-20T115926+0200-decide-hashicorp-go-rootcerts-product-direction.md](2026-09-20T115926+0200-decide-hashicorp-go-rootcerts-product-direction.md)
Next: [2026-09-20T133316+0200-decide-hashicorp-go-sockaddr-product-direction.md](2026-09-20T133316+0200-decide-hashicorp-go-sockaddr-product-direction.md)
Outcome: No exact-path stable go-sockaddr release qualified: v1.0.0 has an incomplete released module closure, v1.0.1-v1.0.3 retain independent behavior/race/vet failures and move unrelated MVS selections, and v1.0.4-v1.0.7 require Go 1.19. Dependency metadata stayed unchanged and one bounded go-sockaddr product decision was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/go-sockaddr v1.0.0` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure,
package behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2, and
all target-specific retained-module decisions through exact inherited,
unloaded go-rootcerts v1.0.2. Every earlier outcome and lifecycle ancestor is
final. Evaluate only Hashicorp go-sockaddr in this session; do not reopen or
combine another dependency group. P8 remains queued.

The authorized 2026-09-20 go-rootcerts option 1 decision retains exact
selected, inherited, unloaded `github.com/hashicorp/go-rootcerts v1.0.2`
without dependency metadata changes. Its target-specific, non-transferable
exception accepts only the completed filesystem-error-identity loss, Darwin
silent-empty system-root pool, proxy-archive symlink-fixture test failure,
v1.0.2 `Config` comparability change, and the completed certificate-pool,
PEM, file/path, precedence, system-root, environment, nil/panic, mutation,
aliasing, allocation, concurrency, global-state, resource, API, MVS,
vulnerability, and related findings. It remains valid only while exact v1.0.2
and all five recorded requests remain unchanged, repository imports and target
loads remain zero, runtime unreachability holds, every earlier guard remains
intact, and no new advisory or independent defect appears. Any change requires
the owning fresh decision.

The go-retryablehttp v0.5.3, go-multierror v1.1.0, go-msgpack v0.5.3,
go-immutable-radix v1.3.1, go-hclog v1.2.0, Errwrap v1.0.0, qualified
go-cleanhttp v0.5.2, and every other earlier exception or qualification remain
separate under their exact selection, incoming-edge, zero-load, runtime-
unreachable, and no-new-finding guards. Revalidate those guards and stop for
the owning decision if any expires. No earlier exception transfers to
go-sockaddr. Do not change a guarded parent, the Go floor, or an unrelated
module.

Selected `github.com/hashicorp/go-sockaddr v1.0.0` is inherited through exact
requests from memberlist v0.1.3 and v0.3.0. The current repository import
search and production and complete-test loads contain zero target packages.
These queue observations and the physical MVS selection are not proof of
repository identity, release qualification, ancestry, floor, behavior,
vulnerability state, or suitability. Resolve them independently and do not
add a direct edge merely to alter MVS.

# Measurements At Start

The go-rootcerts decision recording began from clean handoff HEAD
`c530f7d8e39aa0181899d7cb2893f6a16555ebef`, parent
`888412a62d768cb11af6d9e132600cc4a34f0679`, tree
`e0eeec6c821ba393d8ab30751cf411296b8ec0f6`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status,
reciprocal archive history, and `./codex-dev-start.sh --check` rather than
assuming them. The latest dependency implementation remains exact Google UUID
v1.4.0 commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`.

Guard-only decision revalidation under exact Go 1.26.7 preserved all 16
earlier guarded selections and recorded requests plus exact go-rootcerts
v1.0.2 and its five requests. All 17 why results remain negative, repository
imports are zero, and production and complete-test loads contain zero guarded
packages. The project remains 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, 1,067 sum lines, and the recorded 432-line unapplied tidy projection.
Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV remains empty for
14 earlier guards plus go-rootcerts, Gorilla retains only its recorded
GO-2026-6278/GHSA-w67g-5rqw-f597 result, and go-retryablehttp retains exactly
its two accepted identifiers. No new guarded advisory or independent defect
appeared.

Use exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, use `LC_ALL=C LANG=C`, and
run with `umask 022`. Recreate contained Go 1.18.10 and pinned tools beneath
`$CODEX_SESSION_SCRATCH_ROOT` as required. Preserve the known apidiff archive
reproducibility discrepancy, Python 3.14 Docker timestamp control, managed
bare-`mktemp` restriction, and nested launcher signal-retention timing race;
none is go-sockaddr evidence.

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
runtime boundary. Characterize address construction and parsing, Unix and
platform-specific behavior, interfaces and enumeration, error identity,
malformed inputs, nil/panic behavior, mutation and aliasing, allocation,
concurrency, global state, resources, environment and network interaction, and
actual project consumers. Add independent fixtures where useful and run source
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

At start read this archive, the answered go-rootcerts decision and evaluation,
the answered go-retryablehttp, go-multierror, and go-hclog decisions, rolling
handover, roadmap, `go.mod`, and `go.sum`. Verify the recorded handoff identity,
changed-file set, clean ordinary and ignored status, exact toolchain identity,
reciprocal archive history, every guarded selection/request/load result, and
`./codex-dev-start.sh --check`. Earlier outcomes are final.

# Three Moves

First, resolve go-sockaddr identity, release line, source/test closure, API,
behavior, vulnerability, load, and exact MVS facts without changing the
worktree. Second, if and only if one exact-path stable release preserves Go
1.18 and passes every applicable contract, implement that exact dependency-
only selection and run the complete changed-selection gate; otherwise leave
metadata unchanged and stop for a bounded product decision. Third, update the
roadmap and rolling handover, answer this archive, prepare exactly one
reciprocal NEXT mission for the authorized result, and commit the handoff. Do
not execute the successor.

# Automatic Handoff

After one coherent bounded outcome, make any separate dependency-only commit
first if a qualified selection was implemented, then make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No exact-path stable `github.com/hashicorp/go-sockaddr` release qualifies
under the current Go 1.18, behavior, test, vet, and MVS contracts. The
evaluation therefore stopped without a dependency implementation commit or
any `go.mod`/`go.sum` change and prepared one decision-only successor.

The session began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at handoff HEAD
`a9289a1684e101399810bbec6294b62e4a41df92`, parent
`c530f7d8e39aa0181899d7cb2893f6a16555ebef`, tree
`13addaa0a9dbc39f9a4eb22cbb6b778e072cebd0`. That handoff changed exactly
the launcher, answered go-rootcerts decision archive, this then-NEXT
evaluation archive, rolling handover, and roadmap. Its reciprocal history,
Google UUID implementation ancestry, and launcher check passed.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the exact
module to Hashicorp's public, active, unarchived, non-fork MPL-2.0 repository
at `https://github.com/hashicorp/go-sockaddr.git`. The exact path exposes only
stable v1.0.0-v1.0.7; v1.0.7 is latest. There is no prerelease, retraction,
module deprecation, redirect, alternate exact path, or `/v2` line. The sole
GitHub Release is v1.0.7. All tags are lightweight rather than signed tag
objects. The mirror exposes nine remote branches. Current master is 63 commits
beyond v1.0.7, declares Go 1.25, and is not a release candidate. Proxy and Git
regular-file bytes agree; module zip rules correctly omit the tracked nested
`cmd/sockaddr/vendor` tree.

The exact release identities are:

| Version | Commit | Parent(s) | Tree | Commit time |
| --- | --- | --- | --- | --- |
| v1.0.0 | `a6a0d2df398f7e0e9f6e43f589c8b51cec0eb6b0` | `e92cdb5343bbaf42b0a596937ae0f382270d6759` | `9c45698e106164445d4361773979ea241c73f01a` | 2019-01-18T10:16:41-05:00 |
| v1.0.1 | `3aed17b5ee41761cc2b04f2a94c7107d428967e5` | v1.0.0 | `8127f7e7e070cd2ccf00fbadd56d1c9bc3e46aa9` | 2019-01-18T10:59:32-05:00 |
| v1.0.2 | `c7188e74f6acae5a989bdc959aa779f8b9f42faf` | v1.0.1 | `bbba8b5032324aa052bad8a1b996b2675167c347` | 2019-03-08T08:53:01-05:00 |
| v1.0.3 | `21bd71244b1c754622d9dafb735b7b287c40d17e` | `1bfef1479d8d663aa6659919194a16d5c0afa3b3`, `cc581bb788ed4f44fd65111b5aef4ed1eaa366f6` | `38cc6645a51ae3960b93d29c7d4fae1e15b90381` | 2023-06-09T14:36:52-04:00 |
| v1.0.4 | `5b3b245ad6781cc5142fded60ea8618e0f5696c6` | v1.0.3 | `d2b51dbdaf742fa4e12deb2c155547fceeffcdf3` | 2023-09-01T15:06:03-04:00 |
| v1.0.5 | `fbafcc89f0a6475c72629678aa4c01f904aea824` | v1.0.4 | `4a774d632fd02b203721f51452f15e334c241cce` | 2023-09-06T10:43:13-04:00 |
| v1.0.6 | `081a518b8abca02d4190e80f206993a6cf19a425` | `fbafcc89f0a6475c72629678aa4c01f904aea824`, `eeae47b90fe92dadef85b18ff97ff282fe418ac3` | `91f30a1e9e6001a18d69daf631cc7529095d2f5d` | 2023-11-10T10:01:56-08:00 |
| v1.0.7 | `b74dd36f318ed2ac4e01c93f02ef99739f454ed6` | `ff38b0bfd4772bd19ee23451bb3f65cb18f902e5` | `46f7386f8b92c0d98b88b3bcdad1e91ac86ac00c` | 2024-09-19T10:47:04+01:00 |

GitHub marks the v1.0.0, v1.0.1, v1.0.3, v1.0.6, and v1.0.7 commits
verified and v1.0.2, v1.0.4, and v1.0.5 unverified; none has a signed tag
object. Fresh sumdb lookups bind every release. Selected v1.0.0 has module and
go.mod sums `h1:GeH6tui99pF4NJgfnhp+L6+FfobzVW3Ah46sLo0ICXs=` /
`h1:7Xibr9yA9JjQq1JpNB2Vw7kxv8xerXegt+ozgdvDeDU=`; latest v1.0.7 has
`h1:G+pTkSO01HpR5qCxg7lxfsFEZaG+C0VssTy/9dbT+Fw=` /
`h1:FZQbEYa1pxkQ7WLpyXJ6cbjpT8q0YgQaK/JakXqGyWw=`.

V1.0.0-v1.0.3 omit a `go` directive and are the only serious Go-1.18-floor
candidates. V1.0.4 declares Go 1.19.0 and v1.0.5-v1.0.7 declare Go 1.19, so
they are floor-ineligible even though tested source compiles with Go 1.18.
V1.0.1-v1.0.3 have a 13-module graph and 11 actually imported external
modules; every imported module omits a Go directive and the complete source
and test closure compiles under exact Go 1.26.7 and contained Go 1.18.10.
V1.0.0's released `go.mod` names only the module and omits four command
dependencies. Read-only production/test listing, native tests, race, vet,
command build, API module export, and all seven production/test cross-target
loads therefore fail under both SDKs. Verification of the bytes alone passes
but is not closure proof.

The source contains four packages including one command: the root address
package, template helpers, `cmd/sockaddr`, and its command helper. Across the
candidate line there are no examples, benchmarks, fuzz targets, testdata,
generated files, cgo, embeds, or `go:generate` directives. V1.0.1-v1.0.3 pass
read-only package and closure listing under both SDKs, and all seven Darwin,
Linux, Windows, FreeBSD, Solaris, and Android cross-loads. Their native
count-one, two independent count-ten repeats, and race runs fail on current
Darwin because upstream tests hard-code old loopback flags and sample
interfaces, `/sbin/route -n get default` exits 71 in the managed environment,
unstable default-interface ties reorder results, and Go 1.26 changes an
equal-network-size sort result. Go 1.18 additionally cannot parse the exported
`HashiCorpDefault2016` template. Vet rejects every tested release because
`ifaddrs.go` contains an unreachable panic. These environment-sensitive test
failures do not erase the independent defects below.

The independent fixture fails the same four contracts for v1.0.0-v1.0.3 and
v1.0.7 under both SDKs. `ContainsAddress` incorrectly rejects an interior IPv6
host within a containing network. `GetInterfaceIP` and `GetInterfaceIPs`
return empty output and nil error when no interface matches. All six exported
attribute inventory functions return mutable package-global slices; a
targeted Go 1.26 race run reports a data race. `Host()` exposes the package-
global IPv6 host-mask backing array, so caller mutation corrupts later parsed
addresses. V1.0.3 and later also invert the Windows `exec.LookPath` PowerShell
check, selecting the legacy path when PowerShell exists and attempting
PowerShell when it is absent.

Valid IPv4, IPv6, Unix, RFC, JSON, CLI, interface-enumeration, copying, and
non-aliasing controls otherwise pass. The independent boundary suite also
characterizes the intentional or source-defined panics for unknown/cast
address types, empty comparator lists with multiple values, nil marshaler
receivers, nil attribute inputs, and Must-style malformed input. Unix
construction opens no resource. Interface helpers enumerate the host, while
default-route helpers spawn platform commands without context or timeout;
route failures can be returned or deferred by sorting. The package starts no
goroutines. Template FuncMaps are mutable globals. `NewIPv4Addr` allocates
seven objects per run under Go 1.26.7 and six under Go 1.18.10.

Pinned API comparison finds no root API change from v1.0.0 through v1.0.3;
v1.0.0 whole-module export fails only because its declared closure is
incomplete. V1.0.7 compatibly adds `ErrNoInterface` and `ErrNoRoute`. The
v1.0.1-v1.0.3 command builds and its version/help/RFC-list and positive/
negative RFC membership smokes pass under both SDKs; it identifies itself as
`sockaddr 0.2.0-dev`. V1.0.0 cannot build the command read-only.

Selected v1.0.0 remains inherited through exactly memberlist v0.1.3 and
v0.3.0 requests. Its why result is negative, repository Go imports are zero,
and the 355-entry production and 429-entry complete-test loads contain zero
target packages. The project remains 234 modules, 3,599 graph edges, 197
module-backed entries across 41 loaded modules, and 1,067 sum lines.
V1.0.0 as a redundant root leaves selection unchanged but adds one graph edge
and checksum line. Exact v1.0.1-v1.0.3 projections produce 235 modules and
3,604 edges and also add `mitchellh/go-wordwrap v1.0.0` and upgrade unrelated
`ryanuber/columnize` from its selected 2016 pseudo-version to
v2.1.0+incompatible. V1.0.7 produces 243 modules and 3,623 edges and moves 19
module lines, including guarded Errwrap v1.0.0 to v1.1.0, CLI v1.1.0 to
v1.1.5, and `x/crypto` to v0.17.0. Every projection retains zero target loads;
none was applied. Tidy converges each disposable projection to the recorded
unapplied projection.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z, with no target record. Exact OSV and
GitHub advisory queries are empty for every candidate. Isolated govulncheck
module, package, symbol, and test-symbol scans are empty. Base and disposable
v1.0.3 project scans are byte-identical at all four precision levels, with no
target frame or reachable trace. All 18 target-plus-earlier why results remain
negative; imports and production/test guarded-package loads remain zero, so
every earlier owning guard remains intact.

Exact Go 1.26.7 project verification, build, count-one, two count-ten repeats,
race, vet, pinned golangci-lint 2.12.2, API/CLI compatibility, empty-HOME
count-two, four cross-builds, host/snapshot/Docker acceptance, the canonical
full preflight, all script/meta populations, all eight mutation meta-stages,
80/80 live mutation kills, and all 15 audit meta-controls pass. Contained Go
1.18.10 resolves 366 complete-test entries; its 26 unaffected packages pass
count-one, both count-ten repeats, race, vet, and four cross-builds. The full
floor suite retains only the two accepted Darwin `pkg/shell` closed-file
wording failures. A raw clean-tree audit has 21 automated PASS, zero FAIL, and
exactly six manual-evidence-bound UNMEASURABLE rows; no changed-selection
scorecard applies, so accepted quality remains 27/27 Q0-Q2 PASS at L2. After
the handoff edit, one launcher lifecycle replay reproduced the documented
control-26 signal/log-retention timing race; an independent repeat passed all
62 controls, launcher-auto passed 24 assertions, and `--check` passed.

V1.0.0 is disqualified by its incomplete released source/test closure and the
independent behavior, aliasing, and race defects. V1.0.1-v1.0.3 repair module
closure but retain those defects, fail the test/vet contracts, and require
unrelated MVS movement; v1.0.3 adds the Windows defect. V1.0.4-v1.0.7 are
additionally floor-ineligible. Master is unreleased and requires Go 1.25.
There is therefore no authorized implementation. The sole reciprocal
successor is a bounded go-sockaddr product decision offering guarded v1.0.0
retention, a separately implemented exact v1.0.2 exception, or a continued P7
block; it was prepared but not executed.
