# Agent Session: Evaluate Hashicorp Go Syslog Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T135311+0200-evaluate-hashicorp-go-syslog-dependency`
Created: `2026-09-20T13:53:11+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `978558d5070821485bef07f244f2144e1f668065e33b5a7075b3d8138c189d08`
Previous: [2026-09-20T133316+0200-decide-hashicorp-go-sockaddr-product-direction.md](2026-09-20T133316+0200-decide-hashicorp-go-sockaddr-product-direction.md)
Next: [2026-09-20T145140+0200-decide-hashicorp-go-syslog-product-direction.md](2026-09-20T145140+0200-decide-hashicorp-go-syslog-product-direction.md)
Outcome: No exact-path stable go-syslog release qualified: sole release v1.0.0 preserves Go 1.18 but fails constructor availability, deadline-error, and priority-validation contracts; dependency metadata stayed unchanged and one bounded product decision was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/go-syslog v1.0.0` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure,
package behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded go-sockaddr v1.0.0. Every earlier outcome and lifecycle ancestor is
final. Evaluate only Hashicorp go-syslog in this session; do not reopen or
combine another dependency group. P8 remains queued.

The authorized 2026-09-20 go-sockaddr option 1 decision retains exact
selected, inherited, unloaded `github.com/hashicorp/go-sockaddr v1.0.0`
without dependency metadata changes. Its target-specific, non-transferable
exception accepts only the incomplete released module closure and command
unbuildability; the reproduced IPv6 containment, no-match nil-error, mutable
global-slice/race, and exposed IPv6 mask defects; the characterized upstream
Darwin, Go 1.18 template, vet, sorting, route-command, API/CLI, nil/panic,
mutation, aliasing, allocation, concurrency, global-state, resource,
environment, MVS, vulnerability, and related findings. It remains valid only
while exact v1.0.0 and both memberlist v0.1.3/v0.3.0 requests remain
unchanged, there is no direct root or repository import, target loads remain
zero, runtime unreachability holds, every earlier guard remains intact, and
no new advisory or independent defect appears. Any change requires the owning
fresh decision.

The go-rootcerts v1.0.2, go-retryablehttp v0.5.3, go-multierror v1.1.0,
go-msgpack v0.5.3, go-immutable-radix v1.3.1, go-hclog v1.2.0, Errwrap
v1.0.0, qualified go-cleanhttp v0.5.2, and every other earlier exception or
qualification remain separate under their exact selection, incoming-request,
zero-load, runtime-unreachable, and no-new-finding guards. Revalidate those
guards and stop for the owning decision if any expires. No earlier exception
transfers to go-syslog. Do not change a guarded parent, the Go floor, or an
unrelated module.

Selected `github.com/hashicorp/go-syslog v1.0.0` is inherited through exact
requests from Serf v0.8.2 and v0.9.6. The current repository import search and
production and complete-test loads contain zero target packages. These queue
observations and the physical MVS selection are not proof of repository
identity, release qualification, ancestry, floor, behavior, vulnerability
state, or suitability. Resolve them independently and do not add a direct edge
merely to alter MVS.

# Measurements At Start

The go-sockaddr decision recording began from clean handoff HEAD
`5c957ce042e2c36589e0ec5f00fca5261ce76352`, parent
`a9289a1684e101399810bbec6294b62e4a41df92`, tree
`654abb0f7bf3261c0b2224f99d61ec968af826f2`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status,
reciprocal archive history, and `./codex-dev-start.sh --check` rather than
assuming them. The latest dependency implementation remains exact Google UUID
v1.4.0 commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`.

Guard-only decision revalidation under exact Go 1.26.7 preserved all 18
guarded selections and recorded requests through exact go-sockaddr v1.0.0.
All 18 why results remain negative, repository imports are zero, and
production and complete-test loads contain zero guarded packages. The project
remains 234 modules, 3,599 graph edges, 355 production entries, 429 complete-
test entries, 197 module-backed entries across 41 loaded modules, 1,067 sum
lines, and the recorded 432-line unapplied tidy projection. Base `go.mod` and
`go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV remains empty for
16 guarded modules; Gorilla retains only its recorded
GO-2026-6278/GHSA-w67g-5rqw-f597 result, and go-retryablehttp retains exactly
its accepted GO-2024-2947/GHSA-v6v8-xj6m-xwqh result. Go-sockaddr's GitHub
repository and exact global-advisory queries remain empty. No new guarded
advisory or independent defect appeared.

Use exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, use `LC_ALL=C LANG=C`, and
run with `umask 022`. Recreate contained Go 1.18.10 and pinned tools beneath
`$CODEX_SESSION_SCRATCH_ROOT` as required. Preserve the known apidiff archive
reproducibility discrepancy, Python 3.14 Docker timestamp control, managed
bare-`mktemp` restriction, and nested launcher signal-retention timing race;
none is go-syslog evidence.

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
runtime boundary. Characterize syslog connection and protocol behavior,
Unix/network handling, framing and formatting, dial/write/close ownership,
error identity, malformed inputs, nil/panic behavior, mutation and aliasing,
allocation, concurrency, global state, resources, environment and network
interaction, and actual project consumers. Add independent fixtures where
useful and run source verification, package listing, native complete tests,
two independent repeats, race, vet, and meaningful cross-builds under both
SDKs. Classify every failure precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, floor, architecture, or unrelated-
module change needs a fresh bounded decision rather than silent
implementation. Compare primary vulnerability results at module, package,
symbol, test-symbol, and reachable-trace levels.

# Required Reading

At start read this archive, the answered go-sockaddr decision and evaluation,
the answered go-rootcerts, go-retryablehttp, go-multierror, and go-hclog
decisions, rolling handover, roadmap, `go.mod`, and `go.sum`. Verify the
recorded handoff identity, changed-file set, clean ordinary and ignored status,
exact toolchain identity, reciprocal archive history, every guarded selection/
request/load result, and `./codex-dev-start.sh --check`. Earlier outcomes are
final.

# Three Moves

First, resolve go-syslog identity, release line, source/test closure, API,
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

No exact-path stable `github.com/hashicorp/go-syslog` release qualifies. The
evaluation left `go.mod` and `go.sum` unchanged, created no dependency commit,
and stopped at the required fresh product-decision boundary.

The session began from clean ordinary and ignored worktree state on branch
`codex/upgrade-quality` at handoff HEAD
`2ee06f1db429101b908d764de8c37a9cd2989229`, parent
`5c957ce042e2c36589e0ec5f00fca5261ce76352`, tree
`4c96957b3a4391d2379732962243f12f50928d95`. That handoff changed exactly
`codex-dev-start.sh`, the answered go-sockaddr decision archive, this then-NEXT
go-syslog evaluation archive, the rolling handover, and the roadmap. Its
reciprocal archive chain, latest Google UUID dependency commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, and launcher check all passed.

A freshly downloaded official Go 1.26.7 archive and binary had SHA-256 values
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Fresh Go 1.18.10 archive and binary values were
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
Each exact SDK ran first in `PATH` with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, no ambient `GOFLAGS`, public proxy and sumdb, `LC_ALL=C`,
`LANG=C`, and `umask 022`.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the exact
module to the public, active, unarchived, non-fork MIT Hashicorp repository.
The exact proxy exposes only v1.0.0; it is also `@latest`. There are no other
stable or prerelease versions, retractions, module deprecations, redirects,
alternate exact paths, or `/v2` release line. GitHub has no Release object for
the tag. Six remote branches and 21 returned fork entries are nonrelease or
different-path sources and are not exact-path stable candidates.

V1.0.0 is a lightweight tag at GitHub-verified signed commit
`8d1874e3e8d1862b74e0536851e218c4571066a5`, parent
`326bf4a7f709d263f964a6a96558676b103f3534`, tree
`a81f242fd29ca869a5c2fb666ff9e4dd7a0af570`, dated
2019-01-18T15:19:36Z. Its proxy ZIP SHA-256 is
`a0ca8b61ea365e9ecdca513b94f200aef3ff68b4c95d9dabc88ca25fcb33bce6`
and proxy `go.mod` SHA-256 is
`628ea5197c3c1c9f1e04c7354d9323d12ccb79b1f8a851f817e3bd04f6828aa6`.
All eight regular proxy files match Git byte-for-byte, with no symlink or
special-file discrepancy. Sumdb records
`KaodqZuhUoZereWVIYmpUgZysurB1kBLX2j0MwMrUAE=` for the module and
`qPfqrKkXGihmCqbJM2mZgkZGvKG1dFdvsLplgctolz4=` for its `go.mod`.

Current master is verified commit
`40240a543e78c0ccbdb437b470650ea7acaf3da8`, parent
`8ce6c70daeff97d502321b3a7d712946cdc0c813`, tree
`a1eb97fb750d9bf7cba5b18ee9ba7607d94d0a51`, dated
2026-08-27T06:31:31Z. It is 48 commits beyond v1.0.0, declares Go 1.23.0,
and is unreleased and floor-ineligible. Pinned API comparison reports no
exported API change from v1.0.0; export archives were byte-identical at
SHA-256
`88bff3612eaaca603f8c5399a61fc8dc46d6e363f584e1b1fe5650ccea0b070a`.
Master fixes only the reproduced deadline-error defect below and retains the
constructor and priority defects, so it would not qualify even if it were a
release and floor-eligible.

V1.0.0 has no `go` directive, no external module dependencies, and a complete
minimal production/test closure consisting only of the target and standard
library. It contains one `gsyslog` package, four production Go files, no native
test files, commands, examples, benchmarks, fuzz targets, testdata, generated
files, cgo, embeds, or `go:generate` directives. Source verification, package
listing, native count-one and two count-ten repeats, race, and vet pass under
both SDKs; the native runs correctly report `[no test files]`. Meaningful
production and external-consumer cross-builds were exercised under both SDKs.

The Darwin exported API is `Priority`, severity constants `LOG_EMERG` through
`LOG_DEBUG`, `Syslogger` with `WriteLevel`, `Write`, and `Close`, and
constructors `NewLogger` and `DialLogger`. There is no CLI. Independent
behavior fixtures, SHA-256
`6d9eb0b8fbfefc1556630519301609b6d343cb37ec2b3b5c46d6d6c2378d46d0`,
pass count-one, two count-ten repeats, race, and vet under both SDKs for UDP
and TCP priority/formatting, local-daemon environment behavior, connection
ownership, idempotent close, reconnect after close, buffer mutation and
non-aliasing, concurrent shared use, case-insensitive facilities, invalid
facility/network error identity, allocation, nil-receiver panics, and embedded
newline behavior. The package starts no goroutines and has no mutable package
global state. It reads `os.Args[0]`, hostname and time, performs network I/O,
and probes `/dev/log`, `/var/run/syslog`, and `/var/run/log`; local and remote
write deadlines are 20ms and 50ms. Embedded newlines in messages and tags are
transmitted without escaping, so one call may contain multiple newline-framed
records; this is a characterized malformed-input/framing boundary.

Three independent defects disqualify the sole stable release:

1. The package itself builds on nominal targets, but an external API consumer
   cannot compile because `NewLogger` and `DialLogger` are absent on AIX and
   `js/wasm` under both SDKs and on `wasip1/wasm` under Go 1.26.7. The exact
   external fixture SHA-256 is
   `b50b035d93a2244fbe575ee2d0dbb34ec1d70b9ca501bde3817f517781f8b58a`.
   Darwin, Linux, Windows, FreeBSD, Solaris, Android, and Plan 9 applicable
   builds pass. This contradicts the package's cross-compilation/runtime-error
   purpose rather than representing an unsupported host-only smoke.
2. `writeString` ignores `SetWriteDeadline` errors and writes anyway. A
   controlled connection returns a sentinel deadline error, observes the
   subsequent write, and receives nil from the package under both SDKs. This
   loses error identity and invalidates the stated blocking-avoidance
   guarantee. Unreleased master returns the wrapped deadline error.
3. `DialLogger` accepts out-of-range severity `Priority(8)` with `LOCAL0` and
   silently recomposes it as `LOG_EMERG` instead of rejecting it. The focused
   blocker reproduces under both SDKs and still fails on current master.

The project selects exact v1.0.0 solely through the selected-version requests
`github.com/hashicorp/serf@v0.8.2 -> github.com/hashicorp/go-syslog@v1.0.0`
and `github.com/hashicorp/serf@v0.9.6 ->
github.com/hashicorp/go-syslog@v1.0.0`. Its `go mod why -m` result is negative,
repository Go imports are zero, and the 355-entry production and 429-entry
complete-test loads contain zero target packages. All 18 earlier guarded
targets retain exact selections and requests; every why result and repository
import remains negative, and both project loads contain zero guarded packages.
The 197 module-backed complete-test entries still span 41 modules, so the
target and every earlier guard remain runtime-unreachable.

A disposable redundant v1.0.0 root changes no selected module, adds only the
indirect root, its content checksum, and one graph edge, and leaves both loads
at zero target packages. Its 441-line tidy projection removes the manufactured
root and converges byte-identically with the base 432-line tidy projection:
tidied `go.mod`/`go.sum` SHA-256 values are
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
No projection was applied. The project remains 234 modules, 3,599 graph
edges, 1,067 sum lines, and retains base `go.mod`/`go.sum` SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z, with no target entry. The exact OSV
query, repository advisory feed, and exact global GitHub advisory query are
empty. Go 1.26.7 isolated module/package/symbol/test-symbol scans have zero
findings. Go 1.18.10 reports only old-SDK standard-library findings: 89 module
IDs, ten package IDs, and three reachable IDs, with no advisory assigned to
go-syslog. Base and redundant-root project scans are identical at 30 module,
22 package, 20 production-symbol, and 20 test-symbol IDs, with zero target
package or trace occurrence. Earlier guarded advisory state remains exact.

The unchanged project passes exact Go 1.26.7 module verification, build,
count-one, two independent count-ten repeats, race, vet, pinned
golangci-lint 2.12.2, API/CLI compatibility, an empty-HOME run, four production
cross-builds, canonical full preflight, every script/meta population, all
eight mutation meta-stages with 80/80 kills, host acceptance, pinned snapshot
acceptance, daemon-backed Docker acceptance under physical Python 3.14.6, and
all 15 audit controls. Initial superseded attempts reproduced only the known
cold offline API-base cache, bare-`mktemp` managed sandbox, repository-local
compatibility-report, and Make command-line-variable inheritance boundaries;
scratch-contained canonical reruns pass. No changed-selection scorecard
applies, so accepted quality remains 27/27 Q0-Q2 PASS at L2.

The sole reciprocal successor is the bounded go-syslog product decision. It
was prepared but not executed; no other dependency group or P8 work began.
