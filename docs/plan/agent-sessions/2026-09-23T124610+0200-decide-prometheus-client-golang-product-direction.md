# Agent Session: Decide Prometheus Client Golang Product Direction

Status: NEXT
Session ID: `2026-09-23T124610+0200-decide-prometheus-client-golang-product-direction`
Created: `2026-09-23T12:46:10+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `9ff116a00bc51ee4197d8c52bd498952d1f56f16f094191a2b819647fb3c9ca4`
Previous: [2026-09-23T105846+0200-evaluate-prometheus-client-golang-dependency.md](2026-09-23T105846+0200-evaluate-prometheus-client-golang-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Choose exactly one bounded product direction for exact selected transitive
`github.com/prometheus/client_golang v1.4.0`. The fresh complete evaluation
found no qualifying canonical exact-path stable, no supported tidy-stable
owner for highest-floor v1.16.0, and a guard-crossing raw v1.16.0 projection.
Do not repeat the evaluation, change a dependency or product source, execute
an owner study, combine another dependency group, or begin P8 in this turn.

# Defensive Scope

This is an ordinary dependency-quality product decision. Use only the answered
fresh evaluation and narrow read-only guard revalidation required to choose
between the three options below. Do not fuzz, stress, probe resource
exhaustion, create oversized, deeply nested, cyclic, malformed, adversarial,
or escape-sequence payloads, reproduce a security issue, or perform security
or exploitability analysis.

Every disposable cache, tool, response, or report must remain beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`, `/tmp`, a
sibling of the managed root, or another external root. Verify containment and
remove task-owned scratch evidence before handoff.

# Completed Fresh Evaluation

The fresh evaluation began clean on `codex/upgrade-quality` at guard-decision
handoff HEAD `eea09d0057ad57278945ce9fbba02eba46017c10`, parent
`ff4bcc1b009662d30041e0489fcb2b01c35cbe52`, tree
`488b0b89d05c48f2ffbe390813b9302c50967cfd`. Google UUID implementation
commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remained an ancestor. The
reciprocal archive chain, sole NEXT state, launcher mirror, exact changed set,
ordinary and ignored cleanliness, official SDK identities, real-project
hashes/counts/tidy state, all target and earlier guards, and corrected narrow
advisory identities reproduced. No product source or dependency metadata was
changed.

The public enabled unarchived non-fork Apache-2.0 Prometheus repository owns
51 exact-path stables and six prereleases; `/v2` and `/v3` are absent. Exactly
33 stables declare a Go floor no higher than 1.18, through highest-floor
v1.16.0. Tags, commits, trees, signatures, ancestry, proxy/sumdb, archive-to-
Git content, module directives, licenses, packages, platform/build-tag
boundaries, exported API, and ordinary documented behavior were resolved for
every eligible stable. No replacement, retraction, deprecation, fork, patch,
vendor, branch, or alternate path is eligible.

The complete native matrix contains 660 rows across all 33 eligible stables,
both exact SDKs, module download/verification, build, count-one and count-ten
tests, race, vet, package enumeration, and API capture. Every stable fails at
least one mandatory gate. Selected v1.4.0 builds but fails tests, repetition,
race, and vet under both SDKs. Highest-floor v1.16.0 passes Go 1.18.10 build,
tests, repetition, and race but fails vet; under Go 1.26.7 it builds but fails
tests, repetition, race, and vet because its generated runtime-metric test
support ends before the current Go line. The 1,320-row cgo-disabled cross
build/test-compilation matrix covers Darwin amd64/arm64, Linux amd64/arm64/386,
Windows amd64/386, FreeBSD amd64, Plan 9 amd64, and js/wasm under both SDKs.
V1.16.0 additionally fails 32-bit test compilation and Plan 9 closure build.
No stable qualifies.

Exactly four requests remain: go-metrics v0.3.10 -> v1.4.0, Common v0.9.1 ->
v1.0.0, Common v0.4.1 -> v0.9.1, and TSDB v0.7.1 -> v0.9.1. Every requester
genuinely imports the target. Complete shortest routes run from direct
mvn-pom-mutator v0.2.3 through historical Viper v1.10.1/go-metrics/Common,
or directly through TSDB. Target and requester why are negative; repository
target imports, production/complete-test/module-backed loads, runtime
relevance, and current/history target roots are zero. Exact selected v1.4.0
is tidy-stable through the existing historical owner route. No requester asks
for v1.16.0, and ordinary tidy removes a manufactured v1.16.0 root.

The selected disposable get adds only a redundant root, source sum, and edge,
then tidy restores the common projection. The v1.16.0 get upgrades nine
selections, including protected go-conntrack from `cc309e4a2223` to
`2f068394615f`; tidy removes all nine changes and reselects v1.4.0. Thus
v1.16.0 is neither qualified, guard-preserving, nor genuinely supported and
tidy-stable. Neither projection was retained.

Fresh advisory evidence identifies GHSA-cg3q-j54f-5p7p / GO-2022-0322 /
CVE-2022-21698 for selected v1.4.0 and 22 other eligible releases through
v1.11.0. V1.16.0's narrow OSV and GitHub global responses are empty, but
advisory absence is not qualification. The unchanged project remains exactly
30/22/20/20 module/package/symbol/test-symbol non-stdlib govulncheck findings
with no target trace. The corrected 518,501-byte/1,402-record Go vulnerability
index and PUBLISHED 2,807-byte CNA identities remain respectively
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.

# Choose Exactly One

1. **Retain exact selected v1.4.0.** Explicitly retain exact inherited,
   indirect, unloaded `github.com/prometheus/client_golang v1.4.0` unchanged
   under a target-specific, unqualified, non-transferable exception bounded by
   every recorded request, route, relevance, source, behavior, closure,
   projection, earlier-guard, advisory, and expiry fact. This option does not
   call the selection qualified, supported, safe, or fixed.
2. **Authorize one owner/request-removal study.** Prepare exactly one later
   measurement-only **Mvn-Pom-Mutator Client Golang Owner/Request Removal
   Study** to determine whether a supported stable update of the existing
   direct mvn-pom-mutator root can eliminate all four target requests and
   routes while preserving behavior, every direct root, the Go 1.18 floor,
   and every earlier guard. The study may not implement a change or select any
   client_golang version.
3. **Stop P7 unresolved.** Make no further P7 dependency-quality claim and
   leave P8 queued.

# Decision Contract

Choose exactly one option and explain why it follows from the completed fresh
evaluation. Do not rerun release/source/behavior/closure/projection/API/test/
cross-build work, infer qualification from advisory absence, silently retain
v1.4.0, transfer an earlier exception, or authorize a different study.

Preserve exact selected v1.4.0 and all four requests unless a later explicitly
authorized study and product decision removes them. Preserve the Complete,
go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack, and every earlier
qualified or excepted result without transfer or reopening. Preserve all
project/module/graph/tidy/Go-floor, route/relevance, release/source/closure,
advisory, and expiry guards.

If option 1 is selected, answer this archive and prepare at most the next one
bounded P7 target authorized by the roadmap. If option 2 is selected, prepare
only the named measurement study. If option 3 is selected, record P7 as
unresolved and prepare no dependency successor. Do not execute a successor.

# Required Checks And Handoff

Read this archive, the answered fresh evaluation and guard decision, the
answered Complete decision/evaluation, rolling handover, P7/P8 roadmap,
`go.mod`, and `go.sum`. Verify branch/HEAD/parent/tree, exact changed set,
ancestry, reciprocal archive chain, launcher check, ordinary and ignored
cleanliness, exact project hashes/counts/tidy state, SDK identities, four
requests/import boundaries/routes, why/import/load/runtime/root boundary,
v1.4.0/v1.16.0 projection guards, all earlier guards, and narrow advisory
identities. Stop without choosing an option if a protected input changed.

Update the roadmap and rolling handover, answer this archive, prepare at most
the one successor authorized by the chosen option, verify scratch containment
and cleanup, run final exact-Go unchanged-project module verification, build,
count-one tests, race count-one tests, and vet, and make the local handoff
commit.

# Automatic Handoff

Do not launch a successor or study; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, change a dependency or product
source, add a target root, transfer or reopen an exception, combine another
dependency group, write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
