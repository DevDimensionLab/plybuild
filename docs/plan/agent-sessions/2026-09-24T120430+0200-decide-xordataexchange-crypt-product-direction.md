# Agent Session: Decide Xordataexchange Crypt Product Direction

Status: NEXT
Session ID: `2026-09-24T120430+0200-decide-xordataexchange-crypt-product-direction`
Created: `2026-09-24T12:04:30+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `03dc1dd2ac0ce44e76d5f8b57cb0c30bbb41673ed6763947f804812d1592b7a9`
Previous: [2026-09-24T113445+0200-evaluate-xordataexchange-crypt-dependency.md](2026-09-24T113445+0200-evaluate-xordataexchange-crypt-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Make exactly one product-direction decision for exact selected indirect,
unloaded `github.com/xordataexchange/crypt
v0.0.3-0.20170626215501-b2862e3d0a77`. Choose one numbered direction below,
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
MIT `xordataexchange/crypt` repository on default branch `master`. The
canonical exact-path line has exactly two stables, v0.0.1 and v0.0.2, no
semver prerelease, replacement, retraction, deprecation, or alternate major.
Proxy latest and highest stable are v0.0.2. Selected is the later pseudo-
version of default-branch HEAD.

Both stable Git trees omit go.mod. Their synthesized path-only module files
declare no Go version or dependency even though source imports
armon/consul-api, coreos/go-etcd, and the retired
code.google.com/p/go.crypto/openpgp path. Exact Go 1.18.10 and Go 1.26.7 cannot
resolve that complete closure. Root-only verification passes, but build,
test compilation, count-one/count-ten/race, vet, and all 40 supported cgo-
disabled production/test-compilation rows fail before the sole bounded stable
test executes. Selected's undeclared closure also fails resolution. Passing
root verification cannot qualify an incomplete closure.

The exact owner, release chronology, unsigned tag/commit/parent/tree/ancestry,
proxy/sumdb, synthesized module, license, archive, source and package
identities are recorded. The 14/14/16 regular proxy entries byte-match Git;
no symlink, submodule, executable/special entry, build tag, cgo, generation,
generated/embed, or unsafe boundary exists.

The stable API exposes backend Store/Response, Consul and Etcd clients,
encrypted/plaintext configuration managers, Watch, and secconf Encode/Decode.
Recorded behavior covers reader and byte-slice ownership, OpenPGP/gzip/base64
serialization, crypto and network nondeterminism, synchronous Get/Set,
unbuffered Watch goroutines, shared unsynchronized wait indexes, stop-channel
and retry behavior, unclosed response channels, missing join/cleanup, direct
error propagation, filesystem/stdout/stderr/process-exit CLI boundaries, the
global-keyring helper use, and missing-key dereference. Selected adds list/set,
custom managers, mock backend, and newer Etcd/OpenPGP dependencies; its non-
empty List paths dereference nil result entries. These are ordinary API and
lifecycle facts, not security findings.

The sole current edge is direct, selected, loaded mvn-pom-mutator v0.2.3 ->
selected target. Requester why is positive through `cmd -> pkg/pom`, 22
project Go files import that package, and one requester package is loaded.
All eleven requester Go files have zero target imports. Target why,
repository imports, all 355/429/197/41 load populations, runtime relevance,
and root status are negative. The sum contains only the target module-file
checksum. The current request is metadata-only and has no target source owner.

All 88 follow-history commits reduce to 84 graphable go.mod checkpoints. Each
selects the exact pseudo-version and none roots it. Their 87 request instances
are Viper v1.4.0 at 36, co-pilot-cli mvn-pom-mutator v0.1.41 at four, and
DevDimensionLab mvn-pom-mutator v0.2.0/v0.2.1/v0.2.3 at 2/5/40. Viper v1.4.0
is never selected/rooted; its eight files contain one real target import, so
its requests are source-backed only in an unselected version. All nine
selected Viper versions and their 301 unique Go files contain zero target
imports. Each mutator requester is selected/rooted whenever present, and its
44 unique files contain zero target imports. Thirteen route epochs span six
historical main names. No checkpoint has a selected source-importing owner.

Exact-selected get manufactures only a redundant target root, source sum, and
main edge at 75/1,068 lines, 234 modules and 3,600 edges with unchanged loads.
Exact v0.0.2 get removes direct mvn-pom-mutator v0.2.3, adds the target root,
collapses to 74/1,069 lines and 157 modules/2,219 edges, and breaks every
source import of `pkg/pom`. Ordinary tidy restores the mutator, discards the
stable root, reselects the pseudo-version, and reproduces the exact common
52/948-line, 234-module/3,557-edge projection. No genuine supported tidy-
stable owner exists and no projection is retained.

Exact stable/selected OSV and narrow GitHub results are empty. Pinned
govulncheck v1.8.0 focal scans stop at unresolved closure. Advisory absence or
a stopped scan implies neither safety nor qualification. Corrected index/CNA,
unchanged project 30/22/20/20 without target or named protected traces, and
client_golang GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698 identities
reproduce.

The real project remains exact at 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, 1,067 sum lines, Go 1.18, and the protected module/sum/graph
hashes. Official SDK identities, common tidy projection, Google UUID ancestry,
Xiang90 Probing and every earlier qualified or excepted result, and accepted
27/27 Q0-Q2 PASS at L2 remain exact. Final exact-Go-1.26.7 unchanged-project
verify/build/count-one/race/vet gates pass. The evaluation grants no exception
and retains no dependency, source, root, projection, or study.

# Role And Boundaries

Choose exactly one:

1. Explicitly retain exact selected indirect, unloaded
   `github.com/xordataexchange/crypt
   v0.0.3-0.20170626215501-b2862e3d0a77` unchanged under a
   Xordataexchange-Crypt-specific, unqualified, non-transferable exception.
   State that it is not qualified, safe, or fixed. Preserve complete current/
   historical request, selection, source-import, route, load, runtime, root,
   owner/release/API/behavior/closure/test/cross/projection/advisory evidence
   and every earlier guard as expiry conditions. Do not transfer another
   exception.
2. Authorize exactly one later measurement-only **Mvn-Pom-Mutator
   Xordataexchange Crypt Elimination Study** to determine whether supported
   requester maintenance or replacement can remove the stale current target
   request while preserving product behavior, direct roots, Go 1.18, and
   every earlier guard. Do not run the study, grant an exception, change
   source or dependency metadata, or pre-authorize an implementation.
3. Stop P7 unresolved without changing source, dependency metadata,
   ownership, or any exception.

Option 1 is recommended because the target is unloaded and runtime-irrelevant,
and its sole current request is metadata-only. This is a bounded product
acceptance, not release qualification, a safety claim, or a fix.

# Protected Starting State

Start only from the clean Xordataexchange Crypt evaluation handoff on
`codex/upgrade-quality`. Verify its HEAD, parent, tree, exact five-file changed
set, branch and Google UUID ancestry, reciprocal archive chain, sole NEXT
state, launcher mirror/check, ordinary/ignored cleanliness, official SDK
identities, project counts/hashes, common tidy projection, Go floor, complete
target/requester history, advisory identities, every earlier guard, and final
unchanged-project gate result. Stop for a fresh owning decision if any
protected input changed.

Exact Xiang90 Probing, Ugorji Go, HTTP Unix, TMC, Testify, Objx, and every
earlier qualified or excepted selection remain closed only under their own
target-specific guards. Do not transfer an exception, run a rejected ownership
study, select a rejected candidate, or reopen a completed module. P8 remains
queued.

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
