# Agent Session: Decide Spaolacci Murmur3 Product Direction

Status: NEXT
Session ID: `2026-09-24T054158+0200-decide-spaolacci-murmur3-product-direction`
Created: `2026-09-24T05:41:58+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `60074379bc051f5fa9d50cdb06765d0d1e79aa828c951b8c2ddf2eb83289b3de`
Previous: [2026-09-24T051502+0200-evaluate-spaolacci-murmur3-dependency.md](2026-09-24T051502+0200-evaluate-spaolacci-murmur3-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with exactly one explicit product-direction decision for exact
selected indirect `github.com/spaolacci/murmur3
v0.0.0-20180118202830-f09979ecbc72`. Choose one of the three authorized moves
below. This is a decision-only session: do not rerun the completed evaluation,
execute a study, change a dependency, add a root, combine another dependency
group, or begin P8.

# Closed Evaluation

The fresh bounded Murmur3 evaluation is complete. The canonical exact-path
line contains exactly two stable releases, v1.0.0 and v1.1.0; v1.1.0 is the
highest stable compatible with the Go 1.18 floor. Both tags and the selected
pseudo-version resolve to the same unsigned commit
`f09979ecbc725b9e6d41a297405f65e7e8804acc` and tree
`6dc7b6cbff8cabb8a11d92892c729c806fc1bac8`. There is no canonical
prerelease, replacement, retraction, deprecation, or alternate-major line.
The public, enabled, unarchived, non-fork spaolacci/murmur3 repository and
exact go-import, release chronology, GitHub Releases, tag/commit/tree/
signature/ancestry, proxy/sumdb/archive-to-Git, module/license, source, and
build identities are resolved.

No canonical stable fully qualifies. Both v1.0.0 and v1.1.0 have identical
source and pass module verification, production build, native test
compilation, complete count-one and count-ten tests, and every supported
cgo-disabled production/test-compilation row under exact Go 1.18.10 and Go
1.26.7. Mandatory vet fails under both SDKs at `murmur32.go:129` for possible
misuse of `unsafe.Pointer`. The ordinary upstream race suite also fails under
Go 1.18.10 with a checkptr pointer-arithmetic fatal. Passing safe evidence is
partial and does not qualify either release.

The exact current route is main -> direct mvn-pom-mutator v0.2.3 -> selected
TSDB v0.7.1 -> selected cespare/xxhash v1.1.0 -> selected Murmur3 pseudo-
version. Mvn-pom-mutator's TSDB edge is metadata-only; TSDB genuinely imports
xxhash in production; xxhash genuinely imports Murmur3 only in its tests and
directly requests the selected pseudo-version. Target and requester why are
negative. Both are indirect and unloaded; Murmur3 is absent from project
imports and all project load populations, runtime-irrelevant, selected through
one request, and never a main root.

All 84 reconstructed historical `go.mod` checkpoints retain the same selected
pseudo-version, sole xxhash v1.1.0 request, and no target or requester root.
Their complete shortest-route epochs pass through Viper/client_golang/TSDB,
Viper/crypt-or-Firestore-or-cloud/grpc, or old/new mvn-pom-mutator/TSDB before
the same xxhash -> Murmur3 test-import boundary. No current or historical
requester asks for a canonical Murmur3 stable.

A disposable selected-version get manufactures one indirect target root and
source sum; ordinary tidy removes them and restores the exact common
projection and pseudo-version. A disposable v1.1.0 get changes only Murmur3
among selected modules and adds its source/module sums, but ordinary tidy
discards v1.1.0 because the genuine xxhash test owner still requests the
pseudo-version. It again restores the exact common projection. No genuine
supported tidy-stable requester owns v1.1.0, and no projection is retained.

Exact-version OSV and narrow GitHub global/repository responses are empty for
the selected pseudo-version and both stables. Pinned govulncheck v1.8.0 focal
module/package/symbol/test-symbol results are 0/0/0/0 for v1.1.0. Advisory
absence is not a safety or qualification claim. The unchanged project remains
30/22/20/20 without a Murmur3 or protected named trace; the corrected Go
index, PUBLISHED CNA, client_golang, and every earlier advisory guard remain
exact.

# Protected State And Closed Guards

Start only from the clean Murmur3 evaluation handoff on
`codex/upgrade-quality`. Verify its HEAD, parent, tree, exact changed set,
branch and Google UUID ancestry, reciprocal archive chain, sole NEXT state,
launcher mirror/check, and ordinary and ignored cleanliness. Stop for a fresh
owning decision if a protected input changed.

The unchanged real project remains exactly 234 modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
over 41 loaded modules, and 1,067 `go.sum` lines. Protected `go.mod`, `go.sum`,
graph, normal-tidy diff, common 52/948-line 234-module/3,557-edge projection,
Go 1.18 floor, 27/27 Q0-Q2 PASS at L2, and exact Go 1.18.10/Go 1.26.7 archive
and binary identities remain exact. No projection is retained.

Preserve all 47 pre-Goe selections and 276 incoming edges, every separate
selected version and request count, every qualified result, and every target-
specific unqualified non-transferable exception through Cmux. Rogpeppe/go-
internal remains closed as Cast v1.5.1's minimal test closure. Do not run a
rejected study, reopen a completed module, select a rejected candidate,
transfer an exception, or relax the Go floor.

Preserve Murmur3's exact owner/repository/release chronology, Go-floor,
GitHub Release, tag/commit/tree/signature/ancestry, proxy/sumdb/archive-to-Git,
module/license, source/build, public API/behavior/caller ownership and
mutation, determinism/concurrency/lifecycle/cleanup/error, unsafe/native-
endian/build-tag/cgo/generate/embed, symlink/submodule, closure/native/vet/
test/race/cross, projection, and advisory identities. Preserve the exact
current and 84-checkpoint historical routes; selected sums; target/requester
why, imports, loads, ownership, runtime, and root facts; and every earlier
guard.

The corrected Go vulnerability index remains exactly 518,501 bytes and 1,402
records at its protected SHA-256. The PUBLISHED CNA response remains exactly
2,807 bytes at its protected SHA-256. The project advisory guard remains
30/22/20/20 without a Murmur3 or named protected trace, while client_golang
retains GHSA-cg3q-j54f-5p7p / GO-2022-0322 / CVE-2022-21698. Advisory absence
is not a safety or qualification claim.

# Authorized Product Moves

Choose exactly one:

1. **Retain the selected pseudo-version under a Murmur3-specific exception
   (recommended).** Keep exact selected indirect
   `github.com/spaolacci/murmur3
   v0.0.0-20180118202830-f09979ecbc72` unchanged under a new unqualified,
   non-transferable exception. State precisely that no canonical stable fully
   qualifies and no changed stable has a genuine supported tidy-stable project
   owner. The bounded acceptance may rely on the identical selected/stable
   source, exact xxhash test-import owner, unloaded/runtime-irrelevant target,
   unchanged 84-checkpoint selection, passing safe evidence, and exact
   advisory boundary, but it is not a release qualification or safety claim.
   Expire the exception on any owner, request, route, import/load/runtime,
   root, source/API/behavior, unsafe implementation, closure, vet/test/race,
   Go-floor, projection, or advisory change. Make no dependency or source
   change.
2. **Authorize a later Cespare XXHash Murmur3 Ownership Study.** Define one
   separate measurement-only study that may determine whether a supported
   xxhash release can genuinely request a fully qualified Murmur3 stable while
   preserving every project guard. Do not run the study now, pre-authorize a
   dependency change, add a target root, combine the requester with this
   group, waive vet or race failure, or imply that v1.1.0 qualifies. P7 remains
   blocked on the later study and another explicit owning decision.
3. **Stop P7 unresolved.** Make no dependency, source, root, exception, study,
   other-group, or P8 change. Record Murmur3 as the unresolved blocker.

# Decision Contract

Make exactly one choice. Preserve the selected pseudo-version and every
earlier guard unless option 1 explicitly grants only the target-specific
retention above. Do not transfer another exception. A study authorization is
not permission to run it or to change xxhash, Murmur3, TSDB, mvn-pom-mutator,
or any other module now.

Answer this archive, update the roadmap and rolling handover, verify contained
scratch cleanup, run final exact-Go-1.26.7 project module verification, build,
count-one tests, race count-one tests, and vet, and make one local handoff
commit. Do not rerun the completed Murmur3 evaluation unless a guard-only
check proves a protected input changed.

# Required Reading

Read the answered Murmur3 evaluation and all incorporated Cmux, Goconvey,
Assertions, sanitized_anchor_name, go-diff, Blackfriday, fastuuid, TSDB, and
earlier guards before deciding. Do not reopen completed work or infer
qualification from identical source, passing safe rows, physical selection,
genuine test ownership, or empty advisory results.

# Automatic Handoff

Do not launch a study or successor; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, change a dependency or product
source, add a target root, transfer or reopen an exception, combine another
dependency group, or begin P8. If option 1 closes Murmur3, prepare but do not
launch exactly one next bounded P7 queue item. If option 2 is chosen, prepare
but do not launch exactly the named study. If option 3 is chosen, prepare no
successor.
<!-- CODEX_SESSION_PROMPT_END -->
