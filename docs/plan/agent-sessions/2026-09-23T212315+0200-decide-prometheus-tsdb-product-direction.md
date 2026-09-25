# Agent Session: Decide Prometheus TSDB Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T212315+0200-decide-prometheus-tsdb-product-direction`
Created: `2026-09-23T21:23:15+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `7700486c79810493fe4706435fdf32a86ff15ed5e1309f5fd969e035ed3c7107`
Previous: [2026-09-23T195118+0200-evaluate-prometheus-tsdb-dependency.md](2026-09-23T195118+0200-evaluate-prometheus-tsdb-dependency.md)
Next: [2026-09-23T220817+0200-evaluate-rogpeppe-fastuuid-dependency.md](2026-09-23T220817+0200-evaluate-rogpeppe-fastuuid-dependency.md)
Outcome: Option 1 is final. Exact inherited indirect unloaded TSDB v0.7.1 remains unchanged only under its own unqualified, non-transferable exception; no study, dependency or source change, ownership claim, target root, transferred exception, other dependency group, or P8 work was authorized.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product decision for graph-selected
transitive exact `github.com/prometheus/tsdb v0.7.1`. The completed fresh
evaluation found no fully qualified canonical Go-1.18-compatible stable and
no genuine supported tidy-stable owner. Choose and record exactly one of the
three authorized directions below. Do not repeat the evaluation, implement a
dependency change, combine another group, launch a study or successor, or
begin P8.

# Authorized Roadmap

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

# Measurements At Start

Exact go-import metadata redirects the canonical exact module path to the
public, archived, enabled, non-fork Apache-2.0 `prometheus-junkyard/tsdb`
repository on `master`. The proxy exposes exactly 14 canonical stable releases
v0.1.0, v0.2.0, v0.3.0, v0.3.1, v0.4.0, v0.5.0, v0.6.0, v0.6.1, v0.7.0,
v0.7.1, v0.8.0, v0.9.0, v0.9.1, and v0.10.0. There is no exact-path
prerelease, replacement, retraction, deprecation, `/v2`, or `/v3` line; an
extra Git tag `0.8.0` is noncanonical. Every stable preserves the Go 1.18
floor; v0.10.0 is the highest compatible stable, but no stable qualifies.

Exact tag/commit/tree/signature/ancestry, proxy/sumdb/archive-to-Git,
module/license, source/build boundaries, packages, exported API, and documented
ordinary behavior are complete for all 14. Every exact proxy ZIP matches its
Git tree byte-for-byte. V0.1.0 and v0.2.0 have no Git `go.mod`; their proxy
module files synthesize only the module directive and do not provide a usable
standalone closure. All later original module files have no `go` directive.
The line is pure Go with no cgo, `go:generate`, `go:embed`, generated files,
symlinks, submodules, or special archive entries.

Selected v0.7.1 is valid signed annotated tag object
`249a0812a512567b7278c767772d4bd7bdddefcd`, peeled commit
`c20450564cc42983bf923c13f3fda42de709ac13`, tree
`4224bb77915350b75152cc12f213029050b4eaf8`, with source/mod sums
`h1:YZcsG11NqnK4czYLrWd9mpEuAJIHVQLwdrleYfszMAA=` /
`h1:qhTCs0VvXwvX/y3TZrWD7rabWM+ijKTux40TwIPHuXU=`. Highest v0.10.0 is valid
signed lightweight commit `7762249358193da791ec62e72b080d908f96e776`, tree
`e97090e2c1bf4e052d4f6816eb83232a278ebd1d`, with source/mod sums
`h1:If5rVCMTp6W2SiRAQFlbpJNgVlgMEd+U2GZckwK38ic=` /
`h1:oi49uRhEe9dPUTlS3JRZOwJuVi6tmh10QSgwXEyGCt4=`. Every release identity is
on master ancestry; v0.9.0 alone is unsigned.

Under exact Go 1.18.10 and Go 1.26.7, synthesized v0.1.0/v0.2.0 module
metadata verifies but cannot build its undeclared closure. Every v0.3.0-
v0.10.0 closure verifies and builds but fails mandatory `go vet` under both
SDKs on stable Seek-signature and unkeyed-literal diagnostics, among others.
Thus no buildable stable passes the mandatory native gate.

Static review also found conventionally named upstream `TestReaderFuzz`
functions that generate randomized WAL entries and corruption/incomplete-WAL
tests inside the ordinary suite. Active rows were stopped immediately when
that boundary became clear; no further TSDB test was executed and the
affected test evidence was not used. The defensive scope therefore prevents a
complete count-one/repetition/race qualification independently of the
universal vet failures. No extra behavior fixture or network service was used.

The complete cgo-disabled compile-only matrix covers all 14 releases, both
SDKs, and ten Darwin/Linux/Windows/FreeBSD/Plan-9/js-wasm targets: 280 rows,
each with production build and test-package compilation but no test execution.
V0.1.0/v0.2.0 fail every target from their incomplete closure. V0.3.0-
v0.10.0 pass Darwin, Linux, Windows, and FreeBSD rows, but fail Plan 9 on
incomplete fileutil helpers and js/wasm through the Unix-only x/sys closure.
No cross result rescues a native failure.

Selected v0.7.1 has a 27-module closure and 14 target packages. Exact Go
1.18.10 enumerates 182 production, 210 complete-test, and 62 module-backed
entries across 19 loaded modules; Go 1.26.7 gives 245/273/62 across the same
19. Its API snapshots contain 1,782 lines at
`9025635a356b58dc1209cd9bb2fbc3d4a6599f9c5b7f977c68dee10db8a70354` and
`9b9dddc0505621179b55d65575a82a523e85924677af3a2d00088be7ec55fd6f`.
V0.10.0 has a 40-module closure and 14 packages: Go 1.18.10 gives
180/208/60 entries across 20 loaded modules, Go 1.26.7 gives 243/271/60, and
its 1,846-line API snapshots are
`291355760ea24040a84c52ae1c004336d6b222dedf8d6eda1046b3a089aaaaff` /
`7797f41e3048afc5ced891ec5f828c1ab66428ed1d5c831a1c5b80fe8885280b`.

The full graph has exactly one request: direct
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3` requests v0.7.1 only in
its module metadata and has zero TSDB source imports. The sole complete route
is main -> direct mvn-pom-mutator v0.2.3 -> metadata-only TSDB v0.7.1.
Requester why is positive through main `cmd` and `pkg/pom`; target why is
negative. The repository has zero TSDB imports, and target production,
complete-test, module-backed load, and runtime relevance are zero. There is no
current or historical main-module TSDB root. The requester has direct-root
history, but no alternate genuine target route. `go.sum` has exactly selected
v0.7.1 source and module sums. No requester asks for v0.10.0.

The canonical owner is archived and the sole request is metadata-only, so
there is no genuine supported tidy-stable project owner. A disposable exact
selected get adds 13 redundant indirect roots without changing a selected
module. Raw state is 87/1,082 lines, 234 modules, 3,614 edges, unchanged
355/429/197/41 loads, negative target why, and `go.mod` / `go.sum` / graph
hashes
`74117df39c2e703bae000fb1bb195ea8a4a0d7d36128c96c275d45eef487838a` /
`5b655e3f7ae6afedf4812a80f33ea7284b84dce570ed05cf8111ce2647458dae` /
`a6235a49b5ab79982232d8b710bdc23386fef59eeb1b50471d88e459df55b771`.
Ordinary tidy removes the manufactured roots and returns the exact common
52/948-line, 234-module/3,557-edge projection with v0.7.1 selected. No
v0.10.0 projection was authorized because it is not otherwise qualified. No
projection, dependency change, source change, owner, or exception was retained.

Exact-version OSV for all 14 releases and selected/v0.10.0 narrow GitHub and
repository responses are empty without implying qualification. Pinned
govulncheck v1.8.0 finds three module findings in each isolated focal closure:
GO-2022-0322 in client_golang and GO-2022-0493 plus GO-2026-5024 in x/sys;
only GO-2022-0493 reaches a package, and no vulnerable symbol is found. There
is no TSDB advisory. No exploitability claim was made.

The unchanged real project remains exactly 234 modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
over 41 loaded modules, and 1,067 sum lines at the protected hashes and Go
1.18 floor. Project govulncheck remains 30/22/20/20 with no TSDB or Procfs
trace; client_golang v1.4.0 retains GHSA-cg3q-j54f-5p7p / GO-2022-0322 /
CVE-2022-21698. Corrected Go-index/CNA identities, all 47 pre-Goe selections/
276 incoming edges, separate Goe/pkg-errors/SFTP/go-difflib/Complete/ULID/go-
conntrack/client_golang/client_model/Common/Procfs guards, every earlier
decision, and accepted 27/27 Q0-Q2 PASS at L2 remain exact.

# Role And Boundaries

Choose exactly one:

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/prometheus/tsdb v0.7.1` unchanged under a TSDB-specific,
   unqualified, non-transferable exception. State that it is not qualified,
   supported, safe, or fixed. Preserve the exact sole metadata-only request,
   sole route, target/requester why boundary, zero repository import/load/
   runtime/current-and-history-root state, two sums, release/source/API/
   behavior/closure/native/cross/projection/advisory identities, and every
   earlier guard as expiry conditions. Do not transfer the Procfs, Common,
   client_model, client_golang, or any earlier exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator Prometheus
   TSDB Owner/Request Study** to determine whether a supported update of the
   direct product root can eliminate the stale metadata-only TSDB request or
   establish genuine tidy-stable ownership of a separately qualified future
   Go-1.18-compatible stable while preserving product behavior, direct roots,
   Go 1.18, and every earlier guard. Do not run the study, change a dependency,
   grant an exception, or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because the target is unloaded and runtime-irrelevant,
all qualification and ownership failures are explicit, and it preserves the
unchanged project without broadening a study or transferring an exception.
This is a product acceptance decision, not a qualification claim.

# Required Reading

Start only from the clean TSDB evaluation handoff. Verify its HEAD, parent,
tree, exact changed set, branch/ancestry, reciprocal archive chain, sole NEXT
state, launcher mirror, ordinary and ignored cleanliness, official SDK
identities, real project counts and hashes, common tidy projection, Go 1.18
floor, all target request/import/route/relevance/root facts, target/closure/
project advisory identities, and every earlier guard. Stop for a fresh owning
decision if any protected input changed.

Preserve exact selected Procfs v0.0.8/four requests, Common v0.9.1/four,
client_model v0.2.0/15, and client_golang v1.4.0/four only under their own
separate unqualified non-transferable exceptions. Do not run any rejected
owner study, select their rejected candidates, add a target root, or transfer
an exception. Preserve Complete, go-difflib, SFTP, pkg/errors, Goe, ULID,
go-conntrack, and all earlier qualified or excepted results under their own
exact guards. P8 remains queued.

Do only guard-level revalidation; do not repeat completed release, source,
behavior, API, closure, projection, native, cross, or advisory-matrix work.
Record the chosen product direction, update the roadmap and rolling handover,
answer this archive, verify containment and cleanup, run final exact-Go-1.26.7
project module verification, build, count-one tests, race count-one tests, and
vet, and make the local handoff commit.

# Three Moves

Choose only one of the three numbered directions above. Option 1 is the
recommended bounded acceptance; option 2 authorizes only a later measurement
study; option 3 stops unresolved. None authorizes work in another dependency
group or P8.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, transfer or reopen an exception,
change a dependency or source file, combine another dependency group, or begin
P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 is final. Exact selected, inherited, indirect, unloaded
`github.com/prometheus/tsdb v0.7.1` is explicitly retained unchanged under a
TSDB-specific, unqualified, non-transferable exception. It is not qualified,
supported, safe, or fixed. Every buildable canonical stable fails mandatory
`go vet` under both supported SDKs, the defensive boundary prevents complete
test/repetition/race qualification, and the archived canonical repository plus
sole metadata-only request provides no genuine supported tidy-stable owner.

Option 2's **Mvn-Pom-Mutator Prometheus TSDB Owner/Request Study** is not
authorized or run. Option 3 is not selected because the unchanged inherited,
unloaded, runtime-irrelevant state can be accepted within exact expiry guards
without a qualification or ownership claim. No dependency or source change,
target root, projection, implementation pre-authorization, transferred
exception, other dependency group, or P8 work is included.

### Exact exception and expiry boundary

Retention requires the sole exact request and route to remain unchanged:
direct `github.com/devdimensionlab/mvn-pom-mutator v0.2.3` requests TSDB
v0.7.1 only in module metadata and retains zero TSDB source imports; the only
complete route is main -> direct mvn-pom-mutator v0.2.3 -> metadata-only TSDB
v0.7.1. Requester `go mod why -m` must remain positive through main `cmd` and
requester `pkg/pom`, while target why remains negative. Repository TSDB
imports, target production/complete-test/module-backed loads, runtime
relevance, and current and historical main-module TSDB roots must all remain
zero. The requester retains only its recorded direct-root history and there is
no alternate genuine target route or request for v0.10.0. Any changed request,
route, requester import/metadata boundary, why/import/load/runtime/root fact,
or newly supported owner route expires retention and requires a fresh owning
evaluation and explicit product decision before merge.

The selected release's two verified sumdb identities remain source sum
`h1:YZcsG11NqnK4czYLrWd9mpEuAJIHVQLwdrleYfszMAA=` and module sum
`h1:qhTCs0VvXwvX/y3TZrWD7rabWM+ijKTux40TwIPHuXU=`. Guard-only inspection also
clarifies the protected real-project metadata at its unchanged hash: `go.sum`
contains one TSDB line, the v0.7.1 `/go.mod` sum, rather than a retained source
sum. This corrects the evaluation handoff's `go.sum` wording without changing
the selected release, dependency metadata, or either verified release sum.
Any change to either release identity or to that exact project sum boundary
expires the exception.

Every canonical release and source identity recorded in the answered TSDB
evaluation remains an expiry guard: the exact-path redirect to the public,
enabled, archived, non-fork Apache-2.0 `prometheus-junkyard/tsdb` repository on
`master`; exactly 14 canonical stables v0.1.0-v0.10.0; absent prerelease,
replacement, retraction, deprecation, `/v2`, and `/v3` lines; the noncanonical
Git tag `0.8.0`; and Go-1.18 compatibility for every stable through highest
v0.10.0. Exact tag/commit/tree/signature/ancestry, proxy/sumdb/archive-to-Git,
module/license, source/build boundary, package, exported-API, and documented
ordinary-behavior evidence for all 14 remains incorporated by reference as a
non-transferable expiry condition.

Selected v0.7.1 remains signed annotated tag object
`249a0812a512567b7278c767772d4bd7bdddefcd`, peeled commit
`c20450564cc42983bf923c13f3fda42de709ac13`, tree
`4224bb77915350b75152cc12f213029050b4eaf8`, with the two sums above. Highest
compatible v0.10.0 remains signed lightweight commit
`7762249358193da791ec62e72b080d908f96e776`, tree
`e97090e2c1bf4e052d4f6816eb83232a278ebd1d`, with source/mod sums
`h1:If5rVCMTp6W2SiRAQFlbpJNgVlgMEd+U2GZckwK38ic=` /
`h1:oi49uRhEe9dPUTlS3JRZOwJuVi6tmh10QSgwXEyGCt4=`. Every release remains on
master ancestry and v0.9.0 alone remains unsigned.

The complete closure/native/cross identities remain exact. V0.1.0/v0.2.0
retain synthesized module-directive-only metadata and cannot build their
undeclared standalone closure. Every v0.3.0-v0.10.0 closure verifies and
builds under exact Go 1.18.10 and Go 1.26.7 but fails mandatory vet under both.
The upstream randomized-WAL and corruption-test boundary remains recorded; no
TSDB test was run in this decision session. The complete 280-row compile-only
matrix remains incorporated by reference: v0.1.0/v0.2.0 fail every target;
v0.3.0-v0.10.0 pass Darwin/Linux/Windows/FreeBSD and fail Plan 9/js-wasm.
Any newly qualified canonical stable expires this exception.

Selected v0.7.1 retains its 27-module/14-package closure, exact
182/210/62 and 245/273/62 SDK load populations across 19 modules, and 1,782-
line API hashes
`9025635a356b58dc1209cd9bb2fbc3d4a6599f9c5b7f977c68dee10db8a70354` /
`9b9dddc0505621179b55d65575a82a523e85924677af3a2d00088be7ec55fd6f`.
V0.10.0 retains its 40-module/14-package closure, 180/208/60 and 243/271/60
populations across 20 modules, and 1,846-line API hashes
`291355760ea24040a84c52ae1c004336d6b222dedf8d6eda1046b3a089aaaaff` /
`7797f41e3048afc5ced891ec5f828c1ab66428ed1d5c831a1c5b80fe8885280b`.

The selected disposable projection remains an expiry guard. Exact selected
get has 87/1,082 metadata lines, 234 modules, 3,614 edges, unchanged
355/429/197/41 loads, negative target why, and `go.mod` / `go.sum` / graph
hashes
`74117df39c2e703bae000fb1bb195ea8a4a0d7d36128c96c275d45eef487838a` /
`5b655e3f7ae6afedf4812a80f33ea7284b84dce570ed05cf8111ce2647458dae` /
`a6235a49b5ab79982232d8b710bdc23386fef59eeb1b50471d88e459df55b771`.
Ordinary tidy must continue removing all 13 manufactured roots and restoring
the exact common 52/948-line, 234-module/3,557-edge projection with v0.7.1
selected. No v0.10.0 projection was authorized or retained.

Exact empty target OSV and selected/v0.10.0 GitHub/repository responses remain
guards, and absence is not qualification. Pinned govulncheck v1.8.0 at
database timestamp 2026-09-16T18:00:43Z must retain 3/1/0/0 non-stdlib
module/package/symbol/test-symbol populations for both focal closures:
GO-2022-0322 in client_golang and GO-2022-0493 plus GO-2026-5024 in x/sys,
with only GO-2022-0493 package-reachable and no TSDB or vulnerable-symbol
finding. The unchanged project must remain 30/22/20/20 without a TSDB or
Procfs trace. Client_golang v1.4.0 retains its separate
GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698 identity and exception.

This TSDB exception does not transfer the separate Procfs v0.0.8/four-request,
Common v0.9.1/four-request, client_model v0.2.0/15-request, or client_golang
v1.4.0/four-request unqualified, non-transferable exceptions. It also does
not transfer, reopen, or alter Complete, go-difflib, SFTP, pkg/errors, Goe,
ULID, go-conntrack, or any earlier qualified or excepted result. Each remains
governed only by its own exact guards.

### Guard-only verification and handoff

- Work began clean on `codex/upgrade-quality` at TSDB evaluation handoff HEAD
  `96fe50c72d9b5d85a86e6eb93e0abbdc0d2cfa94`, parent
  `960f1531c6ecf20579e317e0721e693953d7160c`, tree
  `ce68f1086a71a6d05d3671cdf9b1a4ebb14560ec`. It changes exactly the
  launcher, answered TSDB evaluation, this then-NEXT decision, rolling
  handover, and roadmap. Google UUID implementation
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains ancestral.
- The reciprocal 296-archive chain, sole NEXT state, launcher/archive prompt
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
- The unchanged project reproduces 234 modules, 3,599 graph edges, 355
  production entries, 429 complete-test entries, 197 module-backed entries
  across 41 loaded modules, and 1,067 sum lines. `go.mod`, `go.sum`, and graph
  hashes remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  The Go floor remains 1.18. Contained ordinary tidy reaches the exact common
  projection at
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`
  and reselects v0.7.1. No projection is retained.
- The sole request/route, metadata-only requester boundary, positive requester
  why, negative target why, zero target import/load/runtime/current-and-
  history-root state, exact sums, and absent v0.10.0 request reproduce. The
  exact graph preserves all 47 pre-Goe selections/276 incoming edges and the
  separate Goe/pkg-errors/SFTP/go-difflib/Complete/ULID/go-conntrack/
  client_golang/client_model/Common/Procfs selections and request counts.
  Every earlier decision remains closed, separate, and untransferred;
  accepted quality remains 27/27 Q0-Q2 PASS at L2.
- Fresh narrow TSDB, guard OSV, GitHub, repository, focal-closure, and project
  advisory checks reproduce the recorded identities. Normal/no-cache Go
  module indexes remain pairwise byte-identical at 518,501 bytes/1,402 records
  and SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`;
  PUBLISHED memberlist CNA responses remain pairwise byte-identical at 2,807
  bytes and
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

Final unchanged-project exact-Go-1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass under ordinary umask 022. Product
source, `go.mod`, and `go.sum` remain byte-exact. Every task-owned SDK, cache,
response, report, tool, and project copy was contained beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}` and removed before handoff.

P7 remains active with one prepared, unlaunched bounded evaluation of the next
unevaluated alphabetical module, exact `github.com/rogpeppe/fastuuid v1.2.0`.
Its sole grpc-gateway request is a starting observation only. No owner study,
other dependency evaluation, or P8 work was run; P8 remains queued.
