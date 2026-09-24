# Agent Session: Plan First P8 Cloud Modernization Move

Status: NEXT
Session ID: `2026-09-24T130019+0200-plan-first-p8-cloud-modernization-move`
Created: `2026-09-24T13:00:19+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5b79ae79ae06e04d7f89a8ff62d71c777dc96e45376ccbd490f0da321c223a32`
Previous: [2026-09-24T122118+0200-complete-p7-quality-exit-gate.md](2026-09-24T122118+0200-complete-p7-quality-exit-gate.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Begin P8 with one bounded cloud-configuration modernization planning
checkpoint. Characterize the existing `GitCloudConfig` boundary, its
cache-first refresh behavior, and its compatibility fixtures; then select and
specify exactly one smallest owned implementation slice toward `ply-config`.
Do not implement the slice in this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 is active. Its ordered work is cloud configuration,
then Spring modernization under characterization tests, then any inactive
packaging only through a separate scope decision. This checkpoint may plan
only the first cloud-configuration slice. It may not start Spring or packaging
work, reopen P7, evaluate a dependency, run an ownership study, or change
source or dependency metadata.

# Measurements At Start

The P7 exit preserved Go 1.18, all direct roots, product behavior, accepted
27/27 Q0-Q2 PASS at L2, every target-specific P7 guard, and the exact
234-module / 3,599-edge / 355-production / 429-complete-test / 197-module-
backed / 41-loaded-module / 1,067-sum-line project state. Protected `go.mod`,
`go.sum`, and graph hashes remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Exact Go 1.26.7 unchanged-project verification, build, count-one tests, race
count-one tests, vet, compatibility contracts, and quality meta-controls pass.

Begin only from the clean reciprocal P7-exit handoff on
`codex/upgrade-quality`. Verify its branch, five-file documentation/launcher
commit shape, Google UUID ancestry, sole NEXT state, archive chain, launcher
mirror/check, ordinary and ignored cleanliness, and unchanged source and
dependency metadata. Stop for a fresh owning decision if a protected input
changed.

# Role And Boundaries

This is a product-boundary planning checkpoint, not implementation, dependency
evaluation, security investigation, or upstream ownership research. Use the
existing source, tests, design record, and local Git history. Do not add or
modify production code, tests, fixtures, `go.mod`, or `go.sum`; do not fetch or
select a `ply-config` dependency; and do not claim compatibility beyond the
recorded evidence.

Keep every disposable cache, report, project copy, and verification artifact
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Use only ordinary bounded
inspection and established tests. Verify containment and entry types, then
remove all task-owned scratch before handoff.

# Required Reading

Read the P8 roadmap and P7 exit answer, the rolling handover,
`docs/design/quality-lift.md`, `pkg/config/cloud.go`,
`pkg/config/cloud_refresh_test.go`, `pkg/config/cloud_test.go`, the focused
cloud fixture/contract tests, the `CloudConfig` callers in `cmd` and
`pkg/context`, and the existing config-cloud mutation harness. Treat those as
the planning inputs; do not repeat P7 dependency work.

# Three Moves

1. Map the current cloud-configuration construction, refresh, cache, read, and
   caller boundaries. Record the exact behaviors and public/API surfaces that
   a first slice must preserve, especially cache-first operation and current
   compatibility fixtures.
2. Run only focused ordinary tests needed to confirm that characterization.
   Identify any evidence gap, but do not fill it with production/test changes
   in this planning checkpoint. Compare candidate seams by size, ownership,
   reversibility, and compatibility risk.
3. Select exactly one smallest implementation slice, or stop unresolved.
   Specify its files, behavior contract, focused tests, full gates, rollback
   boundary, and expiry conditions. Do not implement or combine a second
   slice.

# Automatic Handoff

If and only if one implementation slice is fully owned and bounded, answer
this archive, update the roadmap and rolling handover, prepare exactly one
reciprocal NEXT archive that implements only that slice, replace only launcher
mutable regions, run launcher/handoff checks, and make one local handoff
commit. If ownership or behavior is unresolved, prepare one decision successor
instead. Do not execute the successor, push, merge, publish, release, stash,
revert, remove the worktree, or broaden P8.
<!-- CODEX_SESSION_PROMPT_END -->
