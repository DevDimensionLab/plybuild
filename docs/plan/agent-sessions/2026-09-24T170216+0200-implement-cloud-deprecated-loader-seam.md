# Agent Session: Implement Cloud Deprecated Loader Seam

Status: NEXT
Session ID: `2026-09-24T170216+0200-implement-cloud-deprecated-loader-seam`
Created: `2026-09-24T17:02:16+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d7ea6a34ef2384064d25fe685ca38cfe9f441fbdae77846bb627a3a3ca8dac24`
Previous: [2026-09-24T164219+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T164219+0200-plan-next-p8-cloud-modernization-move.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the fourth bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the cache-backed `deprecated.json`
read behind one private zero-value-safe loader while preserving the non-
memoized `GitCloudConfig.Deprecated` API and every refresh, caller, other-
reader, fixture, and dependency contract. Do not integrate or select
`ply-config`, and do not begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first three
P8 slices are complete at `57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`, and
`c999212d266d98868930769108e4e8a73070a0a2`. This checkpoint may change only
`pkg/config/cloud.go` and `pkg/config/cloud_test.go` for the private deprecated
loader seam. It may not change a caller, refresh, another reader, constructor,
tracked fixture, mutation file, dependency metadata, Spring, or packaging.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused project-defaults implementation ancestry and
exact two-file shape, services, cache-probe, and Google UUID ancestry, sole
NEXT state, reciprocal 334-archive chain, launcher mirror/check, ordinary and
ignored cleanliness, and unchanged protected source, callers, fixtures,
mutation files, and dependency metadata. The protected project remains Go 1.18
with 234 modules, 3,599 graph edges, 355 production entries, 429 complete-test
entries, 197 module-backed entries over 41 loaded modules, 1,067 `go.sum`
lines, and protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving private loader seam, not `ply-config`
integration, dependency evaluation, generic document-reader work, caller-
policy cleanup, construction refactoring, security work, Spring work, or
packaging work. Do not change any exported symbol, signature, or field; cache
layout; refresh; `Services`, `ProjectDefaults`, or `GlobalCloudConfig`; Git
hooks, templates, or examples; `ListDeprecated`; any Maven, command, or context
caller; tracked fixture; mutation file; `go.mod`; or `go.sum`. Do not fetch,
inspect, or select a new dependency. Keep every disposable cache, report, copy,
and test fixture beneath `${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-
owned scratch after containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
answered project-defaults, services, and cache-probe implementation archives,
current P8 roadmap and rolling handover, `docs/design/quality-lift.md`,
`pkg/config/cloud.go`, `pkg/config/cloud_test.go`, `pkg/config/dir.go`,
`pkg/config/types.go`, `pkg/file/file.go`'s `ReadJson` and `Open` paths,
`pkg/maven/deprecated.go` and its focused partial-failure test, the command
callers, the local acceptance fixture setup in `test/acceptance/common.sh`, and
`scripts/mutate-config-cloud` with its meta-test. Treat every caller, other-
reader, real-network, and fixture gap as a guard.

# Three Moves

1. Add focused failing characterization for production loader selection,
   complete `Directory` delivery, the exact `deprecated.json` request, exactly
   one load per invocation with independent repeated invocations,
   representative complete recursive JSON decode, exact dependency/path/read/
   unmarshal and partial-result semantics, safe zero value, and non-empty
   recordings. Do not characterize or change another reader or caller.
2. Add one private deprecated loader interface and private dependency value.
   The production implementation accepts the complete `Directory`, calls
   `FilePath("deprecated.json")`, and delegates to unchanged `file.ReadJson`.
   The private helper loads exactly once per call and returns the exact value/
   error; a missing loader returns `filesystem.ErrNoFilesystem` before
   developer-path access. Preserve the exported wrapper, non-memoized cache-
   only operation, and every other boundary exactly.
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
