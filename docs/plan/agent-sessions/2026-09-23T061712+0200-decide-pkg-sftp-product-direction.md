# Agent Session: Decide Pkg SFTP Product Direction

Status: NEXT
Session ID: `2026-09-23T061712+0200-decide-pkg-sftp-product-direction`
Created: `2026-09-23T06:17:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e2dadcfc94881e4d9f48d52a81a415f8e6e2725f117279e56554fea740bcafdf`
Previous: [2026-09-23T051316+0200-evaluate-pkg-sftp-dependency.md](2026-09-23T051316+0200-evaluate-pkg-sftp-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by making exactly one product decision for exact selected,
inherited, indirect, unloaded `github.com/pkg/sftp v1.13.1` from the completed
evaluation. Choose only one direction offered below, preserve the explicit
pkg/errors exception and every earlier target-specific guard, update the
roadmap and rolling handover, answer this archive, prepare exactly one
reciprocal successor required by the chosen direction, and make the local
handoff commit. Do not reevaluate SFTP, implement a dependency change, execute
the named study, evaluate another dependency group, or begin P8.

# Defensive Scope

This is an ordinary product decision and guard-only revalidation. Use public
metadata, static repository records, and ordinary project graph/build
commands. Do not fuzz, stress, probe resource exhaustion, construct oversized,
deeply nested, cyclic, malformed, adversarial, or escape-sequence payloads,
reproduce a security issue, or perform security or exploitability analysis.

Every disposable cache, tool, response, report, or project copy must remain
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`,
`/tmp`, a sibling of the managed root, or another external root. Verify
containment and remove task-owned scratch evidence before handoff.

# Completed Evaluation Is Final

No genuine supported selection qualifies. Canonical exact-path v1.13.11 is
the latest of 35 repository stable tags and 22 proxy-listed canonical
`v1.8.0`-through-`v1.13.11` stables. The separate `/v2` path has only two
prereleases, `/v3` is absent, and there is no replacement, retraction,
deprecation, redirect, or eligible alternate. Do not promote a branch, fork,
pseudo-version, prerelease, alternate path, patch, or vendor copy.

V1.13.0 is the highest behavior-qualified Go-1.18-compatible stable. Its exact
commit is `5b7da38a9cdb1bb082343ddc7cff194b751ef665`, tree
`0a95b7d88a409566bc67cc0a45b19b29b6249205`; GitHub reports the merge commit
signature valid. Under both exact Go 1.18.10 and Go 1.26.7 it passes module
verification, complete build, count-one and count-ten tests, race, vet, the
small ordinary in-memory client/server lifecycle fixture, and ten cgo-disabled
cross build/test-compilation targets. It preserves the Go 1.18 project floor.

Selected v1.13.1 and every stable through v1.13.7 fail complete vet under both
SDKs because the internal RawPacket and RequestPacket `ReadFrom` methods have
nonstandard signatures; several versions add unreachable-code or lock-copy
findings. Selected v1.13.1 also fails 32-bit cross test compilation on
overflowing untyped test constants and Plan 9 compilation on undefined
`s_ISVTX`. V1.13.8 and v1.13.9 declare Go 1.15 but their exact x/crypto v0.31.0
closure requires `crypto/ecdh`, absent from Go 1.18. V1.13.10 declares Go
1.23; v1.13.11 declares Go 1.25. None can qualify. V1.13.0 predates selected
API additions including `File.ReadFromWithConcurrency` and
`RealPathFileLister`; that difference is recorded and is not an authorization
to change product behavior.

MVS selects v1.13.1 from genuine source imports in Afero v1.9.4 and v1.8.2.
Historical Afero v1.6.0 and v1.3.3 request and genuinely import SFTP v1.10.1.
The shortest selected route is main -> direct-indirect Afero v1.9.4 -> SFTP
v1.13.1. The recorded historical route is main -> direct mvn-pom-mutator
v0.2.3 -> Viper v1.10.1 -> Afero v1.6.0 -> SFTP v1.10.1. Current Viper
v1.15.0 only contributes its Afero request as metadata; main genuinely imports
Viper, while selected Afero root/internal/common/mem packages load. Afero's
SFTP adapter and the target do not load. Target why, repository imports,
production/complete-test/module-backed loads, runtime relevance, current main
root, and history root are all negative or zero.

An exact disposable selected get adds unused indirect roots for kr/fs,
pkg/errors, SFTP, and x/crypto plus source sums and main edges; tidy removes
them and returns the common projection. A v1.13.7 get changes SFTP, x/crypto,
x/net, x/text, x/mod, and x/tools and removes selected SFTP's genuine
pkg/errors request/import boundary; tidy reselects v1.13.1 and leaves changed
x/net/x/text state. Exact v1.13.0 removes direct mvn-pom-mutator, downgrades
direct Viper v1.15.0 to v1.10.1 and Afero v1.9.4 to v1.6.0, and cannot load at
127 modules/499 edges. Its tidy projection rediscovers mvn-pom-mutator but
reselects v1.13.1 in a non-equivalent 228-module/3,465-edge project with Viper
v1.10.1 and changed loads. It changes the closed pkg/errors SFTP owner/request
route and many earlier guards. No supported tidy-stable project owner requests
v1.13.0. No projection was retained.

The repository/release/tag/commit/tree/signature/ancestry/proxy/sumdb/module/
license/archive identities; API/documentation and source boundaries; globals,
ownership, mutation, determinism, concurrency, lifecycle, cleanup, and error
behavior; complete closures; build constraints and platform/cgo/generated/
embed boundaries; ordinary fixture; exact SDK/native/repeated/race/vet/cross
results; all projections; and advisory results are final. Do not repeat them
unless a narrow guard-only check proves an input changed; stop for a fresh
owning evaluation if it did.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
the ten final target-specific option-1 decisions through pkg/errors, and
qualified no-selection-change modern-go/concurrent and modern-go/reflect2.
P8 remains queued.

Exact selected, inherited, indirect, unloaded
`github.com/pkg/errors v0.9.1` is retained only under its explicit,
unqualified, pkg/errors-specific, non-transferable exception. Its selected
SFTP v1.13.1 genuine source owner, all ten requests and requester boundaries,
routes, root/why/import/load/runtime facts, source/behavior/closure identities,
projections, advisory identities, and expiry conditions remain exact. Do not
run the rejected Pkg Errors Owner-Route Removal Study, substitute standard
`errors`, patch/vendor upstream, change an owner/requester, or transfer the
exception.

Exact Goe, ULID, go-conntrack, mapstructure, go-homedir, Promptui, emoji/v2,
kr/text, kr/pty, kr/pretty, kr/logfmt, kr/fs, go-windows-terminal-sequences,
gotool, errcheck, httprouter, GLS, go-junit-report, json-iterator, and
clockwork exceptions remain final, target-specific, unqualified, and non-
transferable. Every concurrent, reflect2, Cast, Viper, memberlist, earlier
selection, owner, request, route, why/import/load/runtime, graph/module/tidy/
Go-floor, advisory, source, behavior, closure, qualification, and expiry guard
remains final. Do not reopen, broaden, or transfer any decision.

The unchanged project is 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The graph SHA-256 is
`abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
The 432-line tidy diff SHA-256 is
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes are
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

All 47 pre-Goe selections and 276 incoming edges remain exact at
`7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
The separate complete Goe guard and complete pkg/errors guard remain exact.
Forty-one guarded why results are negative; only kr/pretty, kr/text, emoji/v2,
Promptui, go-homedir, and mapstructure are positive. Promptui and go-homedir
are the only guarded repository imports; emoji/v2, Promptui, go-homedir, and
mapstructure are the only loaded guarded modules. Accepted quality remains
27/27 Q0-Q2 PASS at L2.

Exact-version SFTP OSV, GitHub global, and repository advisory results are
empty for selected and serious candidates. Pinned govulncheck v1.8.0 under
exact Go 1.26.7 records identical x/crypto-only populations for v1.13.0 and
v1.13.1, and a smaller but nonzero x/crypto-only population for v1.13.7;
advisory absence or population size is not qualification. Guard OSV remains
limited to Gorilla WebSocket and go-retryablehttp; x/mod v0.14.0 retains
GO-2026-6179 and GO-2026-6180. The Go vulnerability index remains 518,501
bytes/1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

# Measurements At Start

The SFTP evaluation began from clean branch `codex/upgrade-quality` at
handoff HEAD `a6d9bc87ba801899df9d88597100eace51d82488`, parent
`7b078f6d42a02c57feb7a45dd9a821e02d4b0f14`, tree
`cee9fe5c41bd2dbf664742ee851b84ef912aa589`. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an
ancestor. The evaluation changed no product source or dependency metadata and
prepared this decision-only handoff. Independently verify the new handoff
HEAD, parent, tree, exact changed set, ancestry, clean ordinary and ignored
status, reciprocal archive chain, launcher check, exact Go identities, fixed
project/earlier/Goe/pkg-errors/target guards, advisories, and final gates.

# Role And Boundaries

This turn may choose exactly one offered product direction after narrow
guard-only revalidation. It may record that decision and prepare its required
successor, but it may not repeat completed SFTP behavior/closure work, change
source or dependency metadata, run an owner/request study, or implement any
selection. Stop for a fresh owning evaluation if release, route, guard, or
advisory input has changed.

# Choose Exactly One Direction

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/pkg/sftp v1.13.1` without product-source or dependency-metadata
   changes under an SFTP-specific, unqualified, non-transferable exception.
   Record that v1.13.0 is the highest behavior-qualified Go-1.18-compatible
   stable but lacks a genuine supported tidy-stable project owner, while
   selected v1.13.1 fails complete vet and supported cross gates. Guard every
   request, requester import boundary and route, no main/history root,
   why/import/load/runtime facts, repository/release/source/behavior/closure
   identities, all projections, the pkg/errors owner route, all earlier
   decisions, and no new advisory or independent defect. If selected, prepare
   but do not execute the next bounded P7 queue item named by the roadmap.
2. Authorize exactly one later, measurement-only **Afero SFTP Owner/Request
   Study**. It may measure whether a genuine supported Afero owner/request
   route can select behavior-qualified v1.13.0, remove the unloaded SFTP route,
   or establish that neither is supportable while preserving Go 1.18. It may
   inspect the four recorded Afero requests/imports and current/historical
   Viper and mvn-pom-mutator routes. It may not change product source,
   dependency metadata, SFTP, Afero, Viper, mvn-pom-mutator, pkg/errors, a
   target root, any requester, or an earlier decision. Any possible future
   route change that affects pkg/errors requires a fresh pkg/errors decision;
   any implementation requires its own owning product decisions. If selected,
   prepare that named study as the sole successor but do not execute it.
3. Stop P7 unresolved, leave source/dependency metadata and all decisions
   unchanged, prepare no dependency evaluation or study, keep P8 queued, and
   mark the launcher COMPLETE only for the authorized roadmap.

Do not invent a fourth option, combine options, silently retain v1.13.1,
describe it as qualified, safe, fixed, or equivalent to v1.13.0, transfer the
pkg/errors or another exception, or use advisory absence as acceptance. This
decision authorizes no dependency implementation commit.

# Required Reading

Read this archive, the answered SFTP evaluation, answered pkg/errors and Goe
decisions/evaluations, answered ULID and go-conntrack decisions/evaluations,
reflect2 and concurrent evaluations, mapstructure, go-homedir, Promptui,
emoji/v2, kr/text, and kr/pty decisions/evaluations, rolling handover, P7/P8
roadmap, `go.mod`, and `go.sum`. Earlier evaluation facts are final. Refresh
only the narrow decision inputs: selection and four requests, routes and
requester boundaries, root/history/why/import/load/runtime facts, project
hashes/counts/tidy state, earlier/Goe/pkg-errors guards, repository release
status, exact target and guarded advisories, Go-index/memberlist identities,
and final exact-Go module verification, build, count-one tests, race count-one
tests, and vet. Stop for a fresh owning evaluation if a guard changed.

# Three Moves

First, guard-only revalidate the completed result and choose exactly one
offered direction. Second, record only that direction without dependency or
product implementation. Third, update roadmap and rolling handover, answer
this archive, prepare exactly one reciprocal successor required by the choice
or mark the authorized roadmap COMPLETE for option 3, verify scratch
containment, and make the local handoff commit without executing a successor.

# Automatic Handoff

Do not launch a successor, run the study, push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, change source or dependency
metadata, add a root, change an owner/requester, transfer an exception,
reevaluate SFTP, pkg/errors, or an earlier group, evaluate another dependency,
write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
