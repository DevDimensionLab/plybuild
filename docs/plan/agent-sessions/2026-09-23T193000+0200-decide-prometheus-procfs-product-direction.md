# Agent Session: Decide Prometheus Procfs Product Direction

Status: NEXT
Session ID: `2026-09-23T193000+0200-decide-prometheus-procfs-product-direction`
Created: `2026-09-23T19:30:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c057de9a8bc37f2358f719cea929287dc2d486286fe0161ce06e2a9ff2d5a012`
Previous: [2026-09-23T172541+0200-evaluate-prometheus-procfs-dependency.md](2026-09-23T172541+0200-evaluate-prometheus-procfs-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product decision for graph-selected
transitive exact `github.com/prometheus/procfs v0.0.8`. The completed fresh
evaluation found no fully qualified canonical Go-1.18-compatible stable.
Choose and record exactly one of the three authorized directions below. Do not
repeat the evaluation, implement a dependency change, combine another group,
launch a study or successor, or begin P8.

# Defensive Scope

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

# Completed Evaluation

Exact go-import metadata resolves the canonical module to the public, enabled,
unarchived, non-fork Apache-2.0 `prometheus/procfs` repository on `master`.
The exact path has 48 stable releases v0.0.1-v0.22.0 and no prerelease,
replacement, retraction, deprecation, `/v2`, or `/v3` line. Exactly 27 stables
preserve Go 1.18 through highest v0.9.0. Exact tag/commit/tree/signature/
ancestry, proxy/sumdb/archive-to-Git, module/license, package/build boundary,
API, and documented ordinary-behavior evidence is complete for all 27.

Every eligible stable verifies and builds under exact Go 1.18.10 and Go
1.26.7. Selected v0.0.8 verifies, builds, and vets but fails complete
count-one, count-ten, and race tests under both SDKs: its ordinary `TestVM`
fixture result is wrong on Darwin and Btrfs/sysfs suites return their
unsupported-platform errors. Its complete closure is the main module plus
test-only go-cmp v0.3.1 and x/sync `cd5d95a43a6e`.

Highest compatible v0.9.0 verifies, builds, and passes count-one, race, and
vet under both SDKs when its exact Git tree retains the seven test-fixture
symlinks that module ZIPs must omit. It fails count-ten under both SDKs because
`TestNetSoftnet` repeatedly supplies four columns where the parser requires at
least nine. Its complete closure adds go-cmp v0.5.9, x/sync v0.1.0, and x/sys
v0.3.0. No lower stable passes every mandatory native gate.

The complete cgo-disabled matrix covers all 27 releases, both SDKs, and ten
Darwin/Linux/Windows/FreeBSD/Plan-9/js-wasm targets. Selected v0.0.8 passes all
40 selected production-build and test-compilation cells. V0.9.0 passes all
production builds but fails Linux/386 and Windows/386 test compilation under
both SDKs because its tests assign `math.MinInt64`/`math.MaxInt64` to `int`.
Thus no canonical stable fully qualifies.

Exact-version OSV results for all 27 compatible releases and narrow selected/
v0.9.0 GitHub/repository responses are empty without implying qualification.
Pinned isolated govulncheck v1.8.0 is empty for v0.0.8. V0.9.0 has only
module-level GO-2026-5024 in x/sys v0.3.0 and no package, symbol, test-symbol,
or Procfs trace. No exploitability claim was made.

The graph has exactly four target requests. TSDB v0.7.1 and Common v0.4.1 are
metadata-only and request `v0.0.0-20181005140218-185b4288413d`.
Client_golang v1.4.0 genuinely imports Procfs in production and tests and
requests v0.0.8; client_golang v1.0.0 has the same genuine import boundary and
requests v0.0.2. The four complete shortest routes run from main through
direct mvn-pom-mutator v0.2.3 and TSDB or historical Viper v1.10.1/go-metrics/
Common/client_golang chains.

Target/requester why is negative. Repository target/requester imports,
target/requester production and complete-test loads, target module-backed
load, runtime relevance, current target roots, and historical target roots are
zero. `go.sum` has exactly three Procfs go.mod sums and no source sum. No
requester asks for v0.9.0.

A disposable selected exact get only adds a redundant indirect root, source
sum, and main edge; ordinary tidy removes all three and restores inherited
v0.0.8 plus the established common projection. No higher projection was
authorized because no higher stable was otherwise qualified. No projection,
dependency change, source change, owner, or exception was retained.

The unchanged real project remains exactly 234 modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
over 41 loaded modules, and 1,067 sum lines at the protected hashes and Go
1.18 floor. Project govulncheck remains 30/22/20/20 with no Procfs trace;
client_golang v1.4.0 retains GHSA-cg3q-j54f-5p7p / GO-2022-0322 /
CVE-2022-21698. Corrected Go-index/CNA identities, all 47 pre-Goe selections/
276 incoming edges, separate Goe/pkg-errors/SFTP/go-difflib/Complete/ULID/go-
conntrack/client_golang/client_model/Common guards, every earlier decision, and
accepted 27/27 Q0-Q2 PASS at L2 remain exact.

# Authorized Decision

Choose exactly one:

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/prometheus/procfs v0.0.8` unchanged under a Procfs-specific,
   unqualified, non-transferable exception. State that it is not qualified,
   supported, safe, or fixed. Preserve all four exact requests, requester
   import/metadata boundaries, four routes, negative why, zero repository
   import/load/runtime/current-and-history-root state, release/source/API/
   behavior/closure/native/cross/projection/advisory identities, and every
   earlier guard as expiry conditions. Do not transfer the Common,
   client_model, client_golang, or any earlier exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator Prometheus
   Procfs Owner/Request Study** to determine whether a supported update of the
   direct product root can eliminate every Procfs request or establish genuine
   tidy-stable ownership of a separately qualified future Go-1.18-compatible
   stable while preserving product behavior, direct roots, Go 1.18, and every
   earlier guard. Do not run the study, change a dependency, grant an
   exception, or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because the target is unloaded and runtime-irrelevant,
all qualification and ownership failures are explicit, and it preserves the
unchanged project without broadening a study or transferring an exception.
This is a product acceptance decision, not a qualification claim.

# Protected State And Checks

Start only from the clean Procfs evaluation handoff. Verify its HEAD, parent,
tree, exact changed set, branch/ancestry, reciprocal archive chain, sole NEXT
state, launcher mirror, ordinary and ignored cleanliness, official SDK
identities, real project counts and hashes, common tidy projection, Go 1.18
floor, all target request/import/route/relevance/root facts, target/closure/
project advisory identities, and every earlier guard. Stop for a fresh owning
decision if any protected input changed.

Preserve exact selected Common v0.9.1/four requests, client_model v0.2.0/15
requests, and client_golang v1.4.0/four requests only under their own separate
unqualified non-transferable exceptions. Do not run any rejected owner study,
select their rejected candidates, add a target root, or transfer an exception.
Preserve Complete, go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack, and
all earlier qualified or excepted results under their own exact guards. P8
remains queued.

Do only guard-level revalidation; do not repeat completed release, source,
behavior, API, closure, projection, native, cross, or advisory-matrix work.
Record the chosen product direction, update the roadmap and rolling handover,
answer this archive, verify containment and cleanup, run final exact-Go-1.26.7
project module verification, build, count-one tests, race count-one tests, and
vet, and make the local handoff commit.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, transfer or reopen an exception,
change a dependency or source file, combine another dependency group, or begin
P8.
<!-- CODEX_SESSION_PROMPT_END -->
