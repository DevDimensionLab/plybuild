# Agent Session: Implement Cloud Project Defaults Loader Seam

Status: NEXT
Session ID: `2026-09-24T154421+0200-implement-cloud-project-defaults-loader-seam`
Created: `2026-09-24T15:44:21+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8a0b189b8e07590cddf9d3a636ba4d35d87d2da93a332f02e841f77e7f67f23c`
Previous: [2026-09-24T153032+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T153032+0200-plan-next-p8-cloud-modernization-move.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the third bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the cache-backed
`project-defaults.json` read behind one private zero-value-safe loader while
preserving the non-memoized `GitCloudConfig.ProjectDefaults` API and every
refresh, caller, other-reader, fixture, and dependency contract. Do not
integrate or select `ply-config`, and do not begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first two
P8 slices are complete at `57f9d5674157d238a2f93462b65161a17e3b5498`
and `c12307a078a29af74163df0c658d05c821482a33`. This checkpoint may change
only `pkg/config/cloud.go` and `pkg/config/cloud_test.go` for the private
project-defaults loader seam. It may not change a caller, refresh, another
reader, constructor, tracked fixture, mutation file, dependency metadata,
Spring, or packaging.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused services implementation ancestry and exact
two-file shape, cache-probe and Google UUID ancestry, sole NEXT state,
reciprocal 332-archive chain, launcher mirror/check, ordinary and ignored
cleanliness, and unchanged protected source, callers, fixtures, mutation files,
and dependency metadata. The protected project remains Go 1.18 with 234
modules, 3,599 graph edges, 355 production entries, 429 complete-test entries,
197 module-backed entries over 41 loaded modules, 1,067 `go.sum` lines, and
protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving private loader seam, not `ply-config`
integration, dependency evaluation, generic document-reader work, caller-
policy cleanup, construction refactoring, security work, Spring work, or
packaging work. Do not change any exported symbol, signature, or field; cache
layout; refresh; `Services`, `Deprecated`, or `GlobalCloudConfig`; Git hooks,
templates, or examples; `Context.OnEachMavenProject` or
`MergeProjectDefaults`; any command/context caller; tracked fixture; mutation
file; `go.mod`; or `go.sum`. Do not fetch, inspect, or select a new dependency.
Keep every disposable cache, report, copy, and test fixture beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch after
containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
answered services and cache-probe implementation archives, current P8 roadmap
and rolling handover, `docs/design/quality-lift.md`, `pkg/config/cloud.go`,
`pkg/config/cloud_test.go`, `pkg/config/dir.go`, `pkg/config/types.go`,
`pkg/config/project.go`'s `MergeProjectDefaults`, `pkg/file/file.go`'s
`ReadJson` and `Open` paths, `pkg/context/context.go` and its focused defaults
test, the local acceptance fixture setup in `test/acceptance/common.sh`, and
`scripts/mutate-config-cloud` with its meta-test. Treat every caller,
other-reader, real-network, and fixture gap as a guard.

# Three Moves

1. Add focused failing characterization for production loader selection,
   complete `Directory` delivery, the exact `project-defaults.json` request,
   exactly one load per invocation with independent repeated invocations,
   representative complete JSON decode, exact dependency/path/read/unmarshal
   results, safe zero value, and non-empty recordings. Do not characterize or
   change another reader or caller.
2. Add one private project-defaults loader interface and private dependency
   value. The production implementation accepts the complete `Directory`,
   calls `FilePath("project-defaults.json")`, and delegates to unchanged
   `file.ReadJson`. The private helper loads exactly once per call and returns
   the exact value/error; a missing loader returns
   `filesystem.ErrNoFilesystem` before developer-path access. Preserve the
   exported wrapper, non-memoized cache-only operation, and every other
   boundary exactly.
3. Run focused config/context/Maven/command and exact unchanged config-cloud
   mutation gates, then the full recorded exact-Go-1.26.7 gates under
   `umask 022`, readonly/offline module inputs, and managed scratch. Verify no
   API/CLI, source-scope, refresh, caller, other-reader, fixture, Go-floor,
   dependency, graph, mutation, or quality regression. Roll back the single
   two-file commit if any contract fails; do not compensate in another slice.

# Automatic Handoff

If and only if the one seam and every gate pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If a contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge,
publish, release, stash, revert, remove the worktree, integrate `ply-config`,
or start a second cloud, Spring, or packaging slice.
<!-- CODEX_SESSION_PROMPT_END -->
