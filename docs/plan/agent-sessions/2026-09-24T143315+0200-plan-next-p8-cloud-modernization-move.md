# Agent Session: Plan Next P8 Cloud Modernization Move

Status: NEXT
Session ID: `2026-09-24T143315+0200-plan-next-p8-cloud-modernization-move`
Created: `2026-09-24T14:33:15+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8390e2e850c7be6bd49cf8f2cfe9f0f2c065ed089fdae62d13df3b50b67f74bb`
Previous: [2026-09-24T131521+0200-implement-cloud-cache-probe-seam.md](2026-09-24T131521+0200-implement-cloud-cache-probe-seam.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Plan exactly one next bounded P8 cloud-configuration modernization slice after
the completed private cache-probe seam. Re-characterize the resulting
`GitCloudConfig` boundary, compare the smallest repository-owned next seams,
and select one implementation slice or stop unresolved. Do not implement the
slice in this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 remains active in cloud configuration only. The first
P8 implementation slice is complete at focused commit
`57f9d5674157d238a2f93462b65161a17e3b5498`. This checkpoint may plan only
one next cloud-configuration slice. It may not change source or dependency
metadata, integrate or select `ply-config`, begin Spring or packaging work,
reopen P7, or combine multiple implementation moves.

# Measurements At Start

Begin only from the clean reciprocal implementation handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused implementation parent and three-file shape,
Google UUID ancestry, sole NEXT state, reciprocal archive chain, launcher
mirror/check, ordinary and ignored cleanliness, and unchanged fixtures and
dependency metadata. The protected project remains Go 1.18 with 234 modules,
3,599 graph edges, 355 production entries, 429 complete-test entries, 197
module-backed entries over 41 loaded modules, 1,067 `go.sum` lines, and
protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is a planning checkpoint, not implementation, dependency evaluation,
caller-policy cleanup, reader migration, security work, Spring work, or
packaging work. Use only repository source, tests, fixtures, design records,
and local Git history. Do not modify production code, tests, fixtures,
`go.mod`, or `go.sum`; do not fetch, inspect, or select a new dependency; and
do not claim real cloud compatibility beyond the recorded evidence. Keep all
disposable state beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`, verify containment
and entry types, and remove task-owned scratch before handoff.

# Required Reading

Read the answered cache-probe implementation archive and its planning
predecessor, the current P8 roadmap and rolling handover,
`docs/design/quality-lift.md`, `pkg/config/cloud.go`,
`pkg/config/cloud_refresh_test.go`, the focused cloud fixture tests, the
`CloudConfig` callers in `cmd` and `pkg/context`, and
`scripts/mutate-config-cloud` with its meta-test. Treat every recorded caller,
reader, real-network, and fixture gap as a guard rather than permission to fill
multiple gaps.

# Three Moves

1. Confirm the completed cache-probe seam preserves one exact
   `<target>/.git` probe, local-config-first ordering, present-to-pull and
   absent-to-clone selection, Git formatting, zero-value safety, and all
   public/caller/read/fixture contracts.
2. Map only the remaining cloud boundary and run focused read-only tests needed
   to validate candidate ownership. Compare candidates by file count, public
   surface, dependency knowledge, cache compatibility, reversibility, and the
   recorded characterization gaps. Do not fetch or inspect `ply-config`.
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
