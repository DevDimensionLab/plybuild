# Agent Session: Decide Prometheus Common Product Direction

Status: NEXT
Session ID: `2026-09-23T165623+0200-decide-prometheus-common-product-direction`
Created: `2026-09-23T16:56:23+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `79d3f87016027b6acd87f61761431b02fa8318109effa8994b534e5b7d9e8d0d`
Previous: [2026-09-23T144751+0200-evaluate-prometheus-common-dependency.md](2026-09-23T144751+0200-evaluate-prometheus-common-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product decision for graph-selected
transitive exact `github.com/prometheus/common v0.9.1`. The completed fresh
evaluation found no fully qualified canonical stable. Choose and record exactly
one of the three authorized directions below. Do not repeat the completed
evaluation, implement a dependency change, combine another group, launch a
study or successor, or begin P8.

# Defensive Scope

This is an ordinary dependency product-direction decision. Use only the
completed static/release/build/graph/projection/advisory evidence and bounded
read-only guard checks. Do not fuzz, stress, probe resource exhaustion, create
oversized, deeply nested, cyclic, malformed, adversarial, or escape-sequence
payloads, reproduce a security issue, or perform security or exploitability
analysis.

Every disposable cache, report, project copy, or advisory response must remain
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`,
`/tmp`, a sibling of the managed root, or another external root. Set every
tool temp/cache root explicitly, verify containment, and remove task-owned
scratch evidence before handoff.

# Completed Evaluation

The exact public Prometheus owner has 88 canonical non-retracted exact-path
stable candidates after excluding three alternate-module v0 tags, retracted
v0.50.0 and the retracted accidental v1 line, plus one prerelease and no
replacement, deprecation, `/v2`, or `/v3` line. Exactly 51 stables preserve the
Go 1.18 floor through v0.44.0. Exact repository/tag/commit/tree/signature/
ancestry, proxy/sumdb/archive, module/license, source/build-tag, package, API,
and documented ordinary-behavior evidence is complete for all 51.

All 51 verify and build under exact Go 1.18.10 and Go 1.26.7. Selected v0.9.1
fails count-one/count-ten/race tests under both SDKs because its legacy
certificate fixtures and TLS expectations no longer match either supported
toolchain. Only v0.37.0-v0.44.0 are native-clean at the floor; v0.37.1 loses
that status under Go 1.26.7. Highest compatible v0.44.0 is native-clean under
both SDKs.

No release passes the complete bounded cgo-disabled cross matrix. All 51 fail
production build and test compilation on Plan 9 amd64 under both SDKs because
their exact go-conntrack closures reference unavailable
`syscall.ECONNREFUSED`; the other nine targets pass. Exact target-version OSV
and narrow GitHub responses are empty, but advisory absence is not
qualification. Selected v0.9.1's isolated closure has 25/3/0/1 non-stdlib
module/package/symbol/test-symbol findings. V0.44.0 has 20/6/5/5, and every
native-clean compatible release retains reachable closure findings. Thus no
canonical stable fully qualifies.

Exactly four target requests reproduce. TSDB v0.7.1 metadata-only requests
`v0.0.0-20181113130724-41aa239b4cce`; go-metrics v0.3.10 has a genuine test
import and requests v0.9.1; client_golang v1.4.0 has genuine production/test
imports and requests v0.9.1; client_golang v1.0.0 has genuine production/test
imports and requests v0.4.1. Their four complete shortest routes originate at
the main module through direct mvn-pom-mutator v0.2.3 and its TSDB or historical
Viper v1.10.1/go-metrics/Common/client_golang paths. Target/requester why is
negative; repository target import, target/requester production/complete-test
load, target module-backed load, runtime relevance, and current target roots are
zero. Historical target roots are nonzero. `go.sum` has three target go.mod
sums and no target source sum. No requester asks for v0.44.0.

A selected exact get only adds a redundant root, source sum, and main edge;
ordinary tidy restores the established common projection and v0.9.1. A
v0.44.0 get changes 15 existing selections, including protected go-conntrack,
client_golang, and client_model, and adds four modules. Ordinary tidy removes
the Common root and reselects v0.9.1 but retains six unrelated upgrades over
the common projection. V0.44.0 therefore has no genuine supported tidy-stable
owner and is not guard-preserving. Neither projection was retained.

The unchanged real project remains exactly 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, and 1,067 sum lines at the protected hashes and Go 1.18
floor. Project govulncheck remains 30/22/20/20 with no Common trace;
client_golang v1.4.0 retains GHSA-cg3q-j54f-5p7p / GO-2022-0322 /
CVE-2022-21698. Corrected Go-index/CNA identities, all 47 pre-Goe selections/
276 incoming edges, separate Goe/pkg-errors/SFTP/go-difflib/Complete/ULID/go-
conntrack/client_golang/client_model guards, every earlier decision, and
accepted 27/27 Q0-Q2 PASS at L2 remain exact.

# Authorized Decision

Choose exactly one:

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/prometheus/common v0.9.1` unchanged under a Common-specific,
   unqualified, non-transferable exception. State that it is not qualified,
   supported, safe, or fixed. Preserve all four exact requests, requester
   import/metadata boundaries, four routes, negative why, zero repository
   import/load/runtime/current-root state, nonzero historical-root state,
   release/source/API/behavior/closure/native/cross/projection/advisory
   identities, and every earlier guard as expiry conditions. Do not transfer
   the client_golang, client_model, or any earlier exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator Prometheus
   Common Owner/Request Study** to determine whether a supported update of the
   direct product root can eliminate every Common request or establish genuine
   tidy-stable ownership of a separately qualified future Go-1.18-compatible
   stable while preserving product behavior, direct roots, Go 1.18, and every
   earlier guard. Do not run the study, change a dependency, grant an exception,
   or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because the target is unloaded and runtime-irrelevant,
all qualification and ownership failures are explicit, and it preserves the
unchanged project without broadening a study or transferring an exception.
This is a product acceptance decision, not a qualification claim.

# Protected State And Checks

Start only from the clean Common evaluation handoff. Verify its HEAD, parent,
tree, exact changed set, branch/ancestry, reciprocal archive chain, sole NEXT
state, launcher mirror, ordinary and ignored cleanliness, official SDK
identities, real project counts and hashes, common tidy projection, Go 1.18
floor, all target request/import/route/relevance/root facts, target/closure/
project advisory identities, and every earlier guard. Stop for a fresh owning
decision if any protected input changed.

Preserve exact selected v0.9.1 and its four requests. Preserve client_model
v0.2.0/15 requests and client_golang v1.4.0/four requests only under their own
separate unqualified non-transferable exceptions; do not run either rejected
owner study, select v0.4.0 or v1.16.0, add a target root, or transfer an
exception. Preserve Complete, go-difflib, SFTP, pkg/errors, Goe, ULID, go-
conntrack, and all earlier qualified or excepted results under their own exact
guards. P8 remains queued.

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
