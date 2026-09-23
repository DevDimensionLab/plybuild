# Agent Session: Decide Prometheus Client Model Product Direction

Status: NEXT
Session ID: `2026-09-23T141334+0200-decide-prometheus-client-model-product-direction`
Created: `2026-09-23T14:13:34+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `6298a1f49efab8a20c2e729889d0e8c846578d18f4b3f52d7fefe5da5c1e6837`
Previous: [2026-09-23T131751+0200-evaluate-prometheus-client-model-dependency.md](2026-09-23T131751+0200-evaluate-prometheus-client-model-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product decision for graph-selected
transitive exact `github.com/prometheus/client_model v0.2.0`. The completed
fresh evaluation found selected v0.2.0 unqualified because its exact declared
closure does not compile, while highest otherwise-qualified Go-1.18-compatible
stable v0.4.0 has no genuine supported tidy-stable project owner. Choose and
record exactly one of the three authorized directions below. Do not repeat the
completed evaluation, implement a dependency change, combine another group,
launch a study or successor, or begin P8.

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

The exact public Prometheus owner has nine canonical exact-path stable releases
v0.1.0 through v0.6.3, no prereleases, replacements, retractions, or
deprecations, and no `/v2` or `/v3` release line. Exactly v0.1.0 through v0.4.0
preserve the Go 1.18 floor. Exact tag, commit, tree, signature, ancestry,
proxy/sumdb, archive-to-Git, module, license, source, generated/build boundary,
exported API, and documented ordinary behavior evidence is complete.

Selected v0.2.0 verifies but fails build, count-one/count-ten tests, race, vet,
ordinary behavior compilation, and all 40 supported cgo-disabled cross rows
under exact Go 1.18.10 and Go 1.26.7: generated metrics.pb.go requires
`proto.ProtoPackageIsVersion3`, but its own exact go.mod selects
`github.com/golang/protobuf v1.2.0`, where that symbol does not exist. V0.1.0,
v0.3.0, and v0.4.0 pass complete native, repeated, race, vet, behavior, and
cross gates under both SDKs. Thus v0.4.0 is the highest otherwise-qualified
Go-1.18-compatible stable.

Exactly 15 requests and all genuine requester import boundaries reproduce:
go-metrics v0.3.10, client_golang v1.4.0, and Common v0.9.1 request v0.2.0;
client_golang v1.0.0 requests `fd36f4220a90`; Common v0.4.1 and TSDB v0.7.1
request `5c3871d89910`; nine historical go-control-plane vertices request
`14fe0d1b01d4`. All 24 complete shortest routes originate at the main module
and run through direct mvn-pom-mutator, its historical Viper/go-metrics/Common,
gRPC/cloud/genproto/client_golang/TSDB chains, or the recorded Google Martian
and Afero entries into those chains. Target/requester why is negative;
repository target import, production/complete-test/module-backed load, runtime
relevance, and current/history target-root state are zero. No requester asks
for v0.4.0.

A selected get only adds a redundant target root, source sum, and main edge;
tidy restores the common projection. A v0.4.0 get changes client_model v0.2.0
to v0.4.0 and google.golang.org/protobuf v1.28.1 to v1.30.0; ordinary tidy
removes the root and both changes and reselects v0.2.0. V0.4.0 therefore has
no genuine supported tidy-stable project owner. No projection, source,
`go.mod`, or `go.sum` change was retained.

Narrow client_model advisory responses are empty; advisory absence did not
establish qualification. The isolated v0.4.0 closure has only module-level
GO-2024-2611 against protobuf v1.30.0 and zero package/symbol findings. The
unchanged project remains exactly 30/22/20/20 module/package/symbol/test-symbol
findings with no target trace. Client_golang v1.4.0 still retains
GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698. The corrected Go-index
and CVE-2026-14362 CNA identities remain exact.

# Authorized Decision

Choose exactly one:

1. Explicitly retain exact selected, inherited, indirect, unloaded
   `github.com/prometheus/client_model v0.2.0` unchanged under a client_model-
   specific, unqualified, non-transferable exception. State that it is not
   qualified, supported, safe, or fixed. Preserve all 15 exact requests,
   genuine requester import boundaries, 24 routes, negative why, zero target
   repository import/load/runtime/root state, exact release/source/behavior/
   closure/projection/advisory identities, and every earlier guard as expiry
   conditions. Do not transfer the client_golang or any earlier exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator Client Model
   Owner/Request Study** to determine whether a supported update of the direct
   root can establish genuine tidy-stable ownership of qualified v0.4.0 or
   eliminate all client_model requests while preserving product behavior,
   direct roots, Go 1.18, and every earlier guard. Do not run the study, change
   a dependency, grant an exception, or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because the target is unloaded and runtime-irrelevant,
the exact failure and absence of a supported owner are now explicit, and it
preserves the unchanged project without broadening a study or transferring an
exception. This is a product acceptance decision, not a qualification claim.

# Protected State And Checks

Start only from the clean evaluation handoff. Verify its HEAD, parent, tree,
exact changed set, branch/ancestry, reciprocal archive chain, sole NEXT state,
launcher mirror, ordinary and ignored cleanliness, official SDK identities,
real 234-module/3,599-edge/355-production/429-complete-test/197-module-backed/
41-loaded/1,067-sum-line state, module/graph/common-tidy hashes, Go 1.18 floor,
all target requests/routes/relevance facts, advisory identities, and every
earlier guard. Stop for a fresh owning decision if any protected input changed.

All 47 pre-Goe selections and 276 incoming edges remain exact. Preserve the
separate Complete, go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack,
client_golang, and every earlier qualified or excepted result under its own
exact guards. Selected client_golang v1.4.0 remains only under its separate
unqualified non-transferable exception; do not run its rejected owner study,
select v1.16.0, add a root, or transfer that exception. Accepted quality
remains 27/27 Q0-Q2 PASS at L2. P8 remains queued.

Do only guard-level revalidation; do not repeat completed release, source,
behavior, API, closure, fixture, projection, native, cross, or advisory matrix
work. Record the chosen product direction, update the roadmap and rolling
handover, answer this archive, verify containment and cleanup, run final exact-
Go-1.26.7 project module verification, build, count-one tests, race count-one
tests, and vet, and make the local handoff commit.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, transfer or reopen an exception,
change a dependency or source file, combine another dependency group, or begin
P8.
<!-- CODEX_SESSION_PROMPT_END -->
