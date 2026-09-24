# Agent Session: Decide Ugorji Go Product Direction

Status: NEXT
Session ID: `2026-09-24T101913+0200-decide-ugorji-go-product-direction`
Created: `2026-09-24T10:19:13+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `63fbe687f18188db18d15f44245cc62d3aec3356ab2303edf13e3385e80db310`
Previous: [2026-09-24T093713+0200-evaluate-ugorji-go-dependency.md](2026-09-24T093713+0200-evaluate-ugorji-go-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one bounded product-direction decision for exact
selected indirect `github.com/ugorji/go v1.1.4`. The completed fresh
evaluation found no fully qualified canonical Go-1.18-compatible stable and
no genuine supported tidy-stable owner for a changed selection. Choose and
record exactly one of the three authorized directions below. Do not repeat the
evaluation, implement a dependency or source change, combine another group,
launch a study or successor, or begin P8.

# Defensive Scope

This is an ordinary dependency product-direction decision. Use only completed
public release/source/build/graph/projection/advisory evidence and bounded
read-only guard checks. Do not fuzz, stress, probe resource exhaustion, create
oversized, deeply nested, cyclic, malformed, adversarial, or escape-sequence
payloads, reproduce a security issue, or perform security or exploitability
analysis.

Every disposable cache, report, project copy, or advisory response must remain
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Set every tool temp/cache root
explicitly, verify containment, and remove task-owned scratch before handoff.

# Completed Evaluation

Primary go-import metadata maps the exact path to the public, enabled,
unarchived, non-fork `ugorji/go` MIT repository on default branch `master`.
The proxy exposes 27 canonical stable releases from v1.1.1 through v1.2.14
and six prereleases. V1.2.14 is the highest stable; `/v2` and `/v3` do not
exist. There are no replacements, retractions, deprecations, or alternate-
major lines. Every stable preserves the Go 1.18 floor. Exact tags, commits,
trees, signatures, ancestry, sums, archive-to-Git identities, module/license
files, packages, and source boundaries are recorded in the answered
evaluation.

V1.2.14 is a one-file compatibility module with no exported API and a blank
import of split module `github.com/ugorji/go/codec v1.2.14`; codec is therefore
its complete minimal closure. Safe native verification/build/test-compilation/
vet rows, all 40 cgo-disabled production/test-compilation rows across both
exact SDKs and ten supported targets, and small ordinary codec fixtures pass.
Those are partial results only.

Every one of the 27 stable releases contains mandatory complete tests using
explicit cyclic-reference and large-container inputs; later releases also
contain depth-limit and malformed-CBOR cases. Those payload classes are
forbidden by the governing defensive scope. Complete count-one, repeated, and
race tests were therefore stopped before execution under both SDKs. Partial
evidence cannot qualify a release, so no canonical stable fully qualifies.
Selected v1.1.4 also fails native build/test-compilation/vet and every cross
row because codecgen imports undeclared `golang.org/x/tools/go/packages`.

The current sole incoming edge is direct, loaded
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3 -> target v1.1.4`.
Requester why is positive through `cmd -> pkg/pom`, but all 11 requester Go
files have zero target imports. Target why is negative; the repository has no
target import; target is absent from production, complete-test, and module-
backed loads, runtime-irrelevant, indirect, and not a root. The protected sum
contains only its module-file sum.

All 84 distinct historical go.mod checkpoints graph successfully and select
v1.1.4. They contain 87 direct target requests: Viper v1.4.0 in 36,
co-pilot-cli/mvn-pom-mutator v0.1.41 in four, and devdimensionlab/mvn-pom-
mutator v0.2.0/v0.2.1/v0.2.3 in 2/5/40. Every requester version has zero
target imports, so every request is metadata-only. Thirty-seven route epochs
span all six historical main-module names. Target is never a root; exact
requester-root history and all routes are recorded in the evaluation.

A disposable exact v1.2.14 get adds both target v1.2.14 and codec v1.2.14,
producing 235 modules and 3,602 edges and violating the one-selection rule.
Ordinary tidy discards both manufactured roots and returns exactly to the
protected common 234-module/3,557-edge projection, retaining selected v1.1.4
only through requester metadata. No changed candidate has a genuine supported
tidy-stable target owner. No projection was retained.

Selected/candidate exact-version OSV and narrow GitHub/repository results are
empty; pinned govulncheck v1.8.0 reports 0/0/0/0 for candidate and codec.
Those absences do not imply safety or qualification. The protected Go module
index, CNA response, project 30/22/20/20 advisory populations, and
client_golang advisory identity remain exact with no target or named protected
trace.

The real project remains exactly 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, 1,067 sum lines, Go 1.18, and all protected hashes. Exact
SDK identities, common tidy projection, every earlier decision and exception,
and accepted 27/27 Q0-Q2 PASS at L2 remain exact. Final exact-Go-1.26.7
verify/build/count-one/race/vet gates pass. No dependency, source, root,
projection, exception, or study was retained by the evaluation.

# Role And Boundaries

Choose exactly one:

1. Explicitly retain exact selected indirect, unloaded
   `github.com/ugorji/go v1.1.4` unchanged under an Ugorji-Go-specific,
   unqualified, non-transferable exception. State that it is not qualified,
   safe, or fixed. Preserve the exact current and historical requests/routes,
   metadata-only requester boundaries, negative target why/import/load/runtime/
   root facts, sums, release/source/API/behavior/closure/native/test-scope/
   cross/projection/advisory identities, and every earlier guard as expiry
   conditions. Do not transfer another exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator Ugorji Go
   Elimination Study** to determine whether supported requester maintenance or
   replacement can remove the stale target request while preserving product
   behavior, direct roots, Go 1.18, and every earlier guard. Do not run the
   study, grant an exception, change a dependency or source file, or
   pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata, ownership,
   or any exception.

Option 1 is recommended because every observed request is metadata-only and
the target is unloaded and runtime-irrelevant. This is a product acceptance
decision, not a qualification or safety claim.

# Required Reading And Handoff

Start only from the clean Ugorji Go evaluation handoff. Verify its protected
HEAD/parent/tree/changed set, branch and UUID ancestry, reciprocal archive
chain, sole NEXT state, launcher mirror/check, ordinary/ignored cleanliness,
official SDK identities, real project counts/hashes, common tidy projection,
Go floor, complete target/requester facts, advisories, and every earlier guard.
Stop for a fresh owning decision if any protected input changed.

Preserve HTTP Unix, TMC, Testify, Objx, and every earlier qualified or excepted
result only under its own exact guards. Do not run rejected ownership studies,
select a rejected candidate, reopen a completed module, or transfer an
exception. P8 remains queued.

Do only guard-level revalidation. Record the chosen direction, answer this
archive, update the roadmap and rolling handover, verify containment/cleanup,
run final exact-Go-1.26.7 unchanged-project module verification, build,
count-one tests, race count-one tests, and vet, and make one local handoff
commit.

# Three Moves

Choose only one numbered direction. Option 1 is the recommended bounded
acceptance; option 2 authorizes only a later measurement study; option 3 stops
unresolved. None authorizes another dependency group or P8.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, transfer or reopen an exception,
change a dependency or source file, combine another dependency group, or begin
P8.
<!-- CODEX_SESSION_PROMPT_END -->
