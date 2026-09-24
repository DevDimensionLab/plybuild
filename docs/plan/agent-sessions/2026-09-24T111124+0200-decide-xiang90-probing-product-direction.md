# Agent Session: Decide Xiang90 Probing Product Direction

Status: NEXT
Session ID: `2026-09-24T111124+0200-decide-xiang90-probing-product-direction`
Created: `2026-09-24T11:11:24+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `1e8ec92f1306a44e65b0053619a350380ed3a3fcfc62f34ab7fc0b76d21d54cb`
Previous: [2026-09-24T103516+0200-evaluate-xiang90-probing-dependency.md](2026-09-24T103516+0200-evaluate-xiang90-probing-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Make exactly one product-direction decision for exact selected indirect,
unloaded `github.com/xiang90/probing
v0.0.0-20190116061207-43a291ad63a2`. Choose one numbered direction below,
record it as final, preserve every measured expiry guard, and make one local
documentation-only handoff commit. Do not evaluate another dependency, run an
ownership study, change source or dependency metadata, or begin P8.

# Defensive Scope

This is an ordinary product decision based only on the completed evaluation.
Do not perform security or exploitability analysis; fuzz, stress, probe
resource exhaustion; create oversized, deeply nested, cyclic, malformed,
adversarial, or escape-sequence payloads; or reproduce a security issue.
Keep every disposable cache, report, and verification artifact beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`, verify containment, and remove task-owned
scratch before handoff.

# Completed Evaluation

Primary go-import evidence resolves the public, enabled, unarchived, non-fork
MIT `xiang90/probing` repository on default branch `master`. It has two
lightweight tags and GitHub Releases named `0.0.1` and `0.0.2`, but both omit
the mandatory `v` prefix. The canonical exact-path stable and prerelease lines
are therefore empty; `/v2` and `/v3` are absent. Selected is the pseudo-version
of the `0.0.2` commit. Proxy latest is later pseudo-version
`v0.0.0-20221125231312-a49e3df8f510`. There is no canonical stable to select.

Selected and latest have path-only synthesized module files, no declared Go
version or requirements, and one-module standard-library closures. Both pass
exact Go 1.18.10 and Go 1.26.7 verification, build, test compilation, complete
count-one/count-ten/race, vet, and all 40 supported cgo-disabled production/
test-compilation rows. Their seven regular archive entries are byte-identical
to Git; no symlink, submodule, generated/embed, build-tag, cgo, or unsafe
boundary exists. Passing rows cannot turn a pseudo-version into a stable.

The public API, ordinary HTTP/JSON behavior, caller-owned transport and
endpoint slice, live status aliasing, mutex boundaries, nondeterministic
network/timer/clock/goroutine behavior, response-body handling, errors,
rotation, reset, stop-channel/ticker cleanup, and lack of in-flight request
cancellation are recorded exactly. Latest changes only private first-sample
SRTT initialization. The bounded three-test upstream suite passes but is suite
evidence, not proof of correctness.

The sole current edge is direct, loaded mvn-pom-mutator v0.2.3 requesting the
target indirectly. Requester why is positive through `cmd -> pkg/pom`, but its
11 Go files have zero target imports. Target why, repository imports, all
355/429/197/41 load populations, runtime relevance, and root status are
negative. The protected sum contains only the target module-file checksum.

All 88 follow-history commits reduce to 84 graphable module checkpoints. Every
checkpoint selects the exact pseudo-version and none roots it. The 123 request
instances are Viper v1.4.0/v1.7.0/v1.7.1 at 36/1/35, co-pilot-cli
mvn-pom-mutator v0.1.41 at four, and DevDimensionLab mvn-pom-mutator
v0.2.0/v0.2.1/v0.2.3 at 2/5/40. Exact requester selection/root facts,
thirteen route epochs, six historical main names, and all 68 requester Go
files' zero-import boundary are recorded. Every current and historical request
is metadata-only; no genuine target source owner exists.

Selected exact-get manufactures only a redundant target root/source sum/main
edge. Latest exact-get changes exactly the target selection plus the same
manufactured root and sums. Both keep the protected loads exact with zero
target entries. Ordinary tidy removes either root, discards latest, retains
selected only through requester metadata, and restores the exact common
52/948-line, 234-module/3,557-edge projection. No projection is retained and
latest has no genuine supported tidy-stable owner.

Exact selected/latest OSV and narrow GitHub results are empty; pinned
govulncheck v1.8.0 reports 0/0/0/0 for both focal closures. Advisory absence
does not imply safety or qualification. Corrected index/CNA, unchanged project
30/22/20/20 without named protected traces, and client_golang
GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698 identities reproduce.

The real project remains exact at 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, 1,067 sum lines, Go 1.18, and the protected module/sum/graph
hashes. Official SDK identities, common tidy projection, Google UUID ancestry,
every earlier qualified or excepted result, and accepted 27/27 Q0-Q2 PASS at
L2 remain exact. Final exact-Go-1.26.7 unchanged-project verify/build/
count-one/race/vet gates pass. The evaluation grants no exception and retains
no dependency, source, root, projection, or study.

# Role And Boundaries

Choose exactly one:

1. Explicitly retain exact selected indirect, unloaded
   `github.com/xiang90/probing
   v0.0.0-20190116061207-43a291ad63a2` unchanged under a Xiang90-Probing-
   specific, unqualified, non-transferable exception. State that it is not
   qualified, safe, or fixed. Preserve the complete current/historical
   metadata-only request and route identities, negative target facts, source/
   release/API/behavior/closure/test/cross/projection/advisory evidence, and
   every earlier guard as expiry conditions. Do not transfer another
   exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator Xiang90
   Probing Elimination Study** to determine whether supported requester
   maintenance or replacement can remove the stale target request while
   preserving product behavior, direct roots, Go 1.18, and every earlier
   guard. Do not run the study, grant an exception, change source or dependency
   metadata, or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata,
   ownership, or any exception.

Option 1 is recommended because every observed request is metadata-only and
the target is unloaded and runtime-irrelevant. This is a bounded product
acceptance, not release qualification, a safety claim, or a fix.

# Protected Starting State

Start only from the clean Xiang90 Probing evaluation handoff on
`codex/upgrade-quality`. Verify its HEAD, parent, tree, exact five-file changed
set, branch and Google UUID ancestry, reciprocal archive chain, sole NEXT
state, launcher mirror/check, ordinary/ignored cleanliness, official SDK
identities, project counts/hashes, common tidy projection, Go floor, complete
target/requester history, advisory identities, every earlier guard, and final
unchanged-project gate result. Stop for a fresh owning decision if any
protected input changed.

Exact Ugorji Go, HTTP Unix, TMC, Testify, Objx, and every earlier qualified or
excepted selection remain closed only under their own target-specific guards.
Do not transfer an exception, run a rejected ownership study, select a rejected
candidate, or reopen a completed module. P8 remains queued.

# Required Handoff

Do only guard-level revalidation. Record the chosen direction, answer this
archive, update the roadmap and rolling handover, verify containment/cleanup,
run final exact-Go-1.26.7 unchanged-project module verification, build,
count-one tests, race count-one tests, and vet, replace only launcher mutable
regions, run launcher/handoff checks, and make one local handoff commit.

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
