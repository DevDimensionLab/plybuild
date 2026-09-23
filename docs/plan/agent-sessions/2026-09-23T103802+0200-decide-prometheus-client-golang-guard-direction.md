# Agent Session: Decide Prometheus Client Golang Guard Direction

Status: NEXT
Session ID: `2026-09-23T103802+0200-decide-prometheus-client-golang-guard-direction`
Created: `2026-09-23T10:38:02+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e24e6a5b3484234fc3f1e15a2974f460001d24961b9f0b0ac09bb27d392d2d9c`
Previous: [2026-09-23T100201+0200-evaluate-prometheus-client-golang-dependency.md](2026-09-23T100201+0200-evaluate-prometheus-client-golang-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Choose exactly one bounded direction for the stopped
`github.com/prometheus/client_golang v1.4.0` evaluation. The evaluation could
not pass its mandatory starting guard because the NEXT archive's asserted Go
vulnerability-index and memberlist-CNA byte identities disagree with two fresh
primary responses and with the preceding answered archives. Do not resume or
complete the dependency evaluation, select a version, grant an exception, run
an owner study, combine another dependency group, or begin P8 in this turn.

# Defensive Scope

This is an ordinary dependency-quality guard decision. Use only the answered
stopped evaluation and narrow read-only provenance checks required to choose
between the three options below. Do not fuzz, stress, probe resource
exhaustion, create oversized, deeply nested, cyclic, malformed, adversarial,
or escape-sequence payloads, reproduce a security issue, or perform security
or exploitability analysis.

Every disposable cache, tool, response, or report must remain beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`, `/tmp`, a
sibling of the managed root, or another external root. Verify containment and
remove task-owned scratch evidence before handoff.

# Fixed Result At Start

The stopped evaluation began clean on branch `codex/upgrade-quality` at
Complete-decision handoff HEAD
`438178de84797dad6178f1011560c53e192f1e0c`, parent
`18e5105b90fd7380f4621ea82ed4e90ba59daf0e`, tree
`298f5959760cd5332ca85aa524d7083731ac70a0`. Google UUID implementation
commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. The
reciprocal archive chain, sole NEXT state, launcher mirror, exact changed set,
ordinary and ignored cleanliness, and exact official Go 1.18.10/1.26.7
archive and binary identities reproduced.

The real project remained byte-exact at 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. Its exact `go.mod`, `go.sum`,
graph, common-tidy, Go-1.18-floor, 47-selection/276-edge pre-Goe, separate
Goe/pkg-errors/SFTP/go-difflib/Complete, and every earlier target-specific
guard remained unchanged. No product source or dependency metadata changed.

The active NEXT archive asserts that the 518,501-byte, 1,402-record Go
vulnerability index has SHA-256
`bdd6a085321fce25b28e405543bd966376216594e386547494634986e43c282a`
and that the PUBLISHED 2,807-byte CVE-2026-14362 CNA response has SHA-256
`cacd856ff66c4aad5ee56c2c4f45c5a053f16540767cc71e843ec65cc115674c`.
Two independent primary fetches, including explicit no-cache requests, instead
returned those exact byte counts and states at the preceding archives' hashes:
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
The Go response retained Last-Modified 2026-09-17T17:29:18Z; the CNA endpoint
returned `Cache-Control: no-store`. The responses were pairwise byte-identical.
This is a protected starting-input mismatch, so the evaluation stopped.

Before the stop, bounded preliminary work reproduced the four exact requests
and showed that all four requesters genuinely import client_golang: go-metrics
v0.3.10 in production and tests, Prometheus Common v0.9.1 and v0.4.1 in
production version code, and TSDB v0.7.1 in production and tests. Target and
all four requester why results are negative; repository imports, loaded target
entries, runtime relevance, and current/history target roots are zero. Primary
metadata exposed 51 exact-path stables, six prereleases, no `/v2` or `/v3`,
and 33 stables whose declared floor does not exceed Go 1.18; v1.16.0 is the
highest such stable. Selected v1.4.0 has GHSA-cg3q-j54f-5p7p / GO-2022-0322;
v1.16.0's narrow OSV and GitHub global results were empty.

Those facts are preliminary only. The required all-eligible source, behavior,
closure, native/repeated/race/vet/cross-build, API, and advisory evaluation was
interrupted and is not a qualification result. A candidate v1.16.0 scratch get
also changed the protected go-conntrack pseudo-version and other selections;
ordinary tidy removed the candidate and returned the common projection with
v1.4.0 selected. No projection was retained. Do not use the partial matrix or
advisory absence to qualify, reject, retain, or except the target.

# Choose Exactly One

1. **Repair the handoff and reevaluate.** Treat the two asserted NEXT hashes as
   transcription/provenance defects, restore the two independently reproduced
   preceding-archive identities as the active guards, and authorize one fresh
   complete client_golang evaluation from the clean unchanged project. This
   option does not accept any preliminary result or dependency selection.
2. **Require guard provenance.** Preserve the two asserted NEXT hashes and
   authorize exactly one later measurement-only **Client Golang Advisory Guard
   Provenance Study** to identify the exact primary URLs, request conditions,
   or byte transformations that produced them. No dependency evaluation or
   implementation may run until one coherent guard pair is explicitly chosen.
3. **Stop P7 unresolved.** Make no further P7 dependency-quality claim and
   leave P8 queued.

# Decision Contract

Choose exactly one option and explain why it follows from the answered stopped
evaluation. Use only narrow guard-only revalidation; do not resume the partial
33-release matrix, repeat source/behavior/closure/projection work, or infer
qualification from the preliminary evidence.

Preserve exact selected `github.com/prometheus/client_golang v1.4.0` and all
four requests until a fresh authorized evaluation completes. Preserve the
Complete, go-difflib, SFTP, pkg/errors, Goe, ULID, go-conntrack, and every
earlier target-specific exception and qualified result without transfer or
reopening. Preserve all project/module/graph/tidy/Go-floor, route/relevance,
advisory, and expiry guards not limited to the disputed two response hashes.

If option 1 is selected, answer this archive and prepare exactly one fresh
bounded client_golang evaluation successor that discards all partial closure
results and begins from a newly verified clean state. If option 2 is selected,
prepare only the named provenance study. If option 3 is selected, record P7 as
unresolved and prepare no dependency successor. Do not execute the successor.

# Required Checks And Handoff

Read this archive, the answered stopped evaluation, the answered Complete
decision/evaluation, rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`.
Verify branch/HEAD/parent/tree, exact changed set, ancestry, reciprocal archive
chain, launcher check, ordinary and ignored cleanliness, exact real-project
hashes/counts, exact SDK identities, the two conflicting identity pairs, and
all unaffected earlier guards. Stop without choosing an option if any other
protected input changed.

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
