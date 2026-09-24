# Agent Session: Implement Cloud Cache Probe Seam

Status: NEXT
Session ID: `2026-09-24T131521+0200-implement-cloud-cache-probe-seam`
Created: `2026-09-24T13:15:21+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `307f948bec4ba4332baa6c515b4c491f5805aba70a7244b5676b70df4b3ecec9`
Previous: [2026-09-24T130019+0200-plan-first-p8-cloud-modernization-move.md](2026-09-24T130019+0200-plan-first-p8-cloud-modernization-move.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the first bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: inject the private cache-presence probe used
by `GitCloudConfig.Refresh` while preserving the exact cache-first branch and
every public, caller, read, fixture, and dependency contract. Do not integrate
or select `ply-config`, and do not begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. This checkpoint
may implement only the private `<target>/.git` cache-probe seam in
`pkg/config/cloud.go`, its focused characterization in
`pkg/config/cloud_refresh_test.go`, and the corresponding exact mutation
expression in `scripts/mutate-config-cloud`. Spring and inactive packaging
remain out of scope.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, Google UUID ancestry, sole NEXT state, reciprocal
archive chain, launcher mirror/check, ordinary and ignored cleanliness, and
unchanged product source and dependency metadata. The protected project
remains Go 1.18 with 234 modules, 3,599 graph edges, 355 production entries,
429 complete-test entries, 197 module-backed entries over 41 loaded modules,
1,067 `go.sum` lines, and protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving internal seam, not direct `ply-config`
integration, dependency evaluation, caller-policy cleanup, reader migration,
security work, Spring work, or packaging work. Do not change any exported
symbol, signature, or field; construction or command callers; cached read
method; fixture; `go.mod`; or `go.sum`. Do not fetch, inspect, or select a new
dependency. Keep every disposable cache, report, copy, and build artifact
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch after
containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
current P8 roadmap and rolling handover, `docs/design/quality-lift.md`,
`pkg/config/cloud.go`, `pkg/config/cloud_refresh_test.go`, the focused cloud
fixture tests, `internal/adapter/filesystem/filesystem.go`, and
`scripts/mutate-config-cloud` with its meta-test. Treat the recorded caller and
reader gaps as guards, not invitations to broaden this slice.

# Three Moves

1. Add focused failing characterization for one exact cache probe, present and
   missing selection, production dependency selection, safe zero value, and
   local-config error ordering. Preserve all current refresh contracts.
2. Add the smallest private cache-presence dependency to the existing refresh
   dependency bundle. Production must delegate exactly once to the existing
   `file.Exists` behavior for `file.Path("%s/.git", target)`. Keep the existing
   Git dependency, URL/target order, logs, `shell.Output.FormatError`, public
   facade, callers, and cached reads unchanged. Update only the corresponding
   branch mutation expression; retain the same ten mutation IDs and meanings.
3. Run focused config/context/command tests and the exact config-cloud mutation
   harness, then the full recorded exact-Go-1.26.7 gates under `umask 022`,
   readonly/offline module inputs, and managed scratch. Verify no API/CLI,
   source-scope, fixture, Go-floor, dependency-metadata, graph, or quality
   regression. Roll back the single implementation commit if the private seam
   cannot preserve all contracts; do not compensate in another slice.

# Automatic Handoff

If and only if the one slice and all gates pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If any contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge, publish,
release, stash, revert, remove the worktree, integrate `ply-config`, or start a
second cloud, Spring, or packaging slice.
<!-- CODEX_SESSION_PROMPT_END -->
