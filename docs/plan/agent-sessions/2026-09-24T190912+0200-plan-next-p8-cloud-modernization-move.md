# Agent Session: Plan Next P8 Cloud Modernization Move

Status: NEXT
Session ID: `2026-09-24T190912+0200-plan-next-p8-cloud-modernization-move`
Created: `2026-09-24T19:09:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `66f73bd2834fc9d4178c5d28fc1aeb6db5e68374ee12ed4c43f8a68aa9db333a`
Previous: [2026-09-24T180445+0200-implement-cloud-global-config-loader-seam.md](2026-09-24T180445+0200-implement-cloud-global-config-loader-seam.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Plan exactly one next bounded P8 cloud-configuration modernization slice after
the completed private global-config loader seam. Re-characterize the resulting
cloud boundary, compare the smallest repository-owned next seams, and select
one implementation slice or stop unresolved. Do not implement the slice in
this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 remains active in cloud configuration only. The first
five P8 implementation slices are complete at focused commits
`57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`,
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`, and
`f12344b115b6c732295d152a7f3d86b876a4d1d8`. This checkpoint may plan only
one next cloud-configuration slice. It may not change source or dependency
metadata, integrate or select `ply-config`, begin Spring or packaging work,
reopen an earlier slice, or combine multiple implementation moves.

# Measurements At Start

Begin only from the clean reciprocal implementation handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused global-config implementation parent and exact
two-file shape, earlier deprecated, project-defaults, services, and cache-probe
ancestry, Google UUID ancestry, sole NEXT state, reciprocal 337-archive chain,
launcher mirror/check, ordinary and ignored cleanliness, and unchanged
protected readers, callers, fixtures, mutation files, and dependency metadata.
The protected project remains Go 1.18 with 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, 1,067 `go.sum` lines, and protected `go.mod` / `go.sum` /
graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is a planning checkpoint, not implementation, dependency evaluation,
caller-policy cleanup, generic reader migration, construction refactoring,
security work, Spring work, or packaging work. Use only repository source,
tests, fixtures, design records, and local Git history. Do not modify production
code, tests, fixtures, mutation files, `go.mod`, or `go.sum`; do not fetch,
inspect, or select a new dependency; and do not claim real cloud compatibility
beyond recorded evidence. Keep every disposable artifact beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`, verify containment and entry types, and
remove task-owned scratch before handoff.

# Required Reading

Read the answered global-config implementation archive and its planning
predecessor, the answered deprecated, project-defaults, services, and cache-
probe implementation archives, current P8 roadmap and rolling handover,
`docs/design/quality-lift.md`, `pkg/config/cloud.go`, all focused cloud tests,
the `CloudConfig` callers in `cmd` and `pkg/context`, and
`scripts/mutate-config-cloud` with its meta-test. Treat every recorded caller,
other-reader, real-network, environment, and fixture gap as a guard rather than
permission to combine work. Do not inspect `ply-config`.

# Three Moves

1. Confirm the completed global-config seam preserves production loader
   selection, complete `Directory` delivery, exact direct
   `global-config.yaml` path without `FilePath`, unchanged `file.Open`,
   environment-before-YAML ordering, one independent load per invocation,
   exact read/YAML/partial results, safe zero behavior, `SourceFor`, and all
   public, refresh, caller, reader, fixture, and mutation contracts.
2. Map only the remaining cloud boundary and run focused read-only tests needed
   to validate candidate ownership. Compare candidates by file count, public
   surface, environment and dependency knowledge, cache compatibility,
   reversibility, and the recorded characterization gaps. Do not fetch or
   inspect `ply-config`.
3. Select exactly one smallest repository-owned implementation slice with its
   files, behavior contract, tests, mutation impact, full gates, rollback, and
   expiry conditions, or stop unresolved. Do not implement it or combine a
   second move.

# Automatic Handoff

If and only if one next slice is fully owned and bounded, answer this archive,
update the roadmap and rolling handover, prepare exactly one reciprocal NEXT
archive for that implementation, replace only launcher mutable regions, run
launcher/handoff checks, and make one local handoff commit. If ownership or
behavior is unresolved, prepare one decision successor instead. Do not execute
the successor, push, merge, publish, release, stash, revert, remove the
worktree, integrate `ply-config`, or begin Spring or packaging work.
<!-- CODEX_SESSION_PROMPT_END -->
