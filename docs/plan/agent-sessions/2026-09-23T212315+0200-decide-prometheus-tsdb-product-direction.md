# Agent Session: Decide Prometheus TSDB Product Direction

Status: NEXT
Session ID: `2026-09-23T212315+0200-decide-prometheus-tsdb-product-direction`
Created: `2026-09-23T21:23:15+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `7700486c79810493fe4706435fdf32a86ff15ed5e1309f5fd969e035ed3c7107`
Previous: [2026-09-23T195118+0200-evaluate-prometheus-tsdb-dependency.md](2026-09-23T195118+0200-evaluate-prometheus-tsdb-dependency.md)
Next: none
Outcome: pending

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
