# Agent Session: Implement Cloud Services Loader Seam

Status: NEXT
Session ID: `2026-09-24T145009+0200-implement-cloud-services-loader-seam`
Created: `2026-09-24T14:50:09+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `490b4424b2b0d279641784123fd72d41b4512911c0f7468899ff02b2699257e0`
Previous: [2026-09-24T143315+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T143315+0200-plan-next-p8-cloud-modernization-move.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the second bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the eager cached `services.json`
load behind one private zero-value-safe services loader while preserving the
closure-shaped `GitCloudConfig.Services` API and every refresh, caller, other-
reader, fixture, and dependency contract. Do not integrate or select
`ply-config`, and do not begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first P8
slice is complete at `57f9d5674157d238a2f93462b65161a17e3b5498`. This
checkpoint may change only `pkg/config/cloud.go` and
`pkg/config/cloud_test.go` for the private cached-services loader seam. It may
not change refresh, another cached reader, a constructor or caller, the
config-cloud mutation files, fixtures, dependency metadata, Spring, or
packaging.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused cache-probe implementation ancestry and exact
three-file shape, Google UUID ancestry, sole NEXT state, reciprocal 330-
archive chain, launcher mirror/check, ordinary and ignored cleanliness, and
unchanged source, fixtures, and dependency metadata. The protected project
remains Go 1.18 with 234 modules, 3,599 graph edges, 355 production entries,
429 complete-test entries, 197 module-backed entries over 41 loaded modules,
1,067 `go.sum` lines, and protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving private loader seam, not `ply-config`
integration, dependency evaluation, generic document-reader work, caller-
policy cleanup, construction refactoring, security work, Spring work, or
packaging work. Do not change any exported symbol, signature, or field; cache
layout; refresh dependency or branch; `Deprecated`, `ProjectDefaults`,
`GlobalCloudConfig`, Git hooks, templates, or examples; any command/context
caller; fixture; `go.mod`; or `go.sum`. Do not fetch, inspect, or select a new
dependency. Keep every disposable cache, report, copy, and build artifact
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch after
containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
answered cache-probe implementation archive, current P8 roadmap and rolling
handover, `docs/design/quality-lift.md`, `pkg/config/cloud.go`,
`pkg/config/cloud_test.go`, `pkg/config/cloud_refresh_test.go`,
`pkg/config/dir.go`, `pkg/file/file.go`'s `ReadJson` and `Open` paths, the
focused cloud fixture tests, and `scripts/mutate-config-cloud` with its meta-
test. Treat all caller, other-reader, real-network, and fixture gaps as guards.

# Three Moves

1. Add focused failing characterization for production loader selection,
   complete dependency delivery, the exact `services.json` request, eager one-
   call loading, stable repeated closure results and errors, existing fixture
   decode, exact `FilePath` error propagation, safe zero value, and non-empty
   recordings. Do not characterize another reader or caller.
2. Add one private services-loader interface and private dependency value. The
   production implementation accepts the complete `Directory`, calls
   `FilePath("services.json")`, and delegates to unchanged `file.ReadJson`.
   The private helper loads exactly once when `Services()` is called and closes
   over the exact value/error; a missing loader returns
   `filesystem.ErrNoFilesystem` without developer-path access. Preserve the
   exported wrapper, cache-only operation, and every other boundary exactly.
3. Run focused config/context/command and exact unchanged config-cloud mutation
   gates, then the full recorded exact-Go-1.26.7 gates under `umask 022`,
   readonly/offline module inputs, and managed scratch. Verify no API/CLI,
   source-scope, refresh, caller, other-reader, fixture, Go-floor, dependency,
   graph, mutation, or quality regression. Roll back the single two-file commit
   if any contract fails; do not compensate in another slice.

# Automatic Handoff

If and only if the one seam and every gate pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If a contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge, publish,
release, stash, revert, remove the worktree, integrate `ply-config`, or start a
second cloud, Spring, or packaging slice.
<!-- CODEX_SESSION_PROMPT_END -->
