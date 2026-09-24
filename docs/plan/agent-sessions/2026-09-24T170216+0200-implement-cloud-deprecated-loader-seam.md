# Agent Session: Implement Cloud Deprecated Loader Seam

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T170216+0200-implement-cloud-deprecated-loader-seam`
Created: `2026-09-24T17:02:16+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d7ea6a34ef2384064d25fe685ca38cfe9f441fbdae77846bb627a3a3ca8dac24`
Previous: [2026-09-24T164219+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T164219+0200-plan-next-p8-cloud-modernization-move.md)
Next: [2026-09-24T174851+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T174851+0200-plan-next-p8-cloud-modernization-move.md)
Outcome: the fourth bounded P8 slice is complete at focused commit `51da9bcfc0eb8ae871c98467c52c942db6b480bd`; the deprecated cache read now has one private zero-value-safe loader with its public non-memoized behavior unchanged

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

## Answer

The fourth bounded P8 cloud slice is complete at focused implementation commit
`51da9bcfc0eb8ae871c98467c52c942db6b480bd` (parent
`deb6e42d333bb8621c4579ec6b658ab36a9af572`, tree
`e095c58797d42b9fc4a1ba353c9ddc1916a70614`). It changes exactly
`pkg/config/cloud.go` and `pkg/config/cloud_test.go`. No second reader, caller,
refresh, constructor, fixture, mutation file, dependency metadata, Spring,
packaging, or `ply-config` work ran.

### Implementation And Characterization

- `deprecatedLoader` and `deprecatedDependencies` are private. Their zero value
  returns exact `filesystem.ErrNoFilesystem` before accessing the supplied
  `Directory`; production explicitly selects `fileDeprecatedLoader`.
- The production loader receives the complete `Directory`, requests exactly
  `deprecated.json`, preserves exact `FilePath` errors, and delegates to
  unchanged `file.ReadJson`. The exported `Deprecated()` signature remains the
  production wrapper. Its private helper performs exactly one load for each
  invocation, so repeated invocations remain independent, cache-only, and
  non-memoized.
- Focused TDD first failed on the absent private interface, dependency selector,
  production loader, and helper. The final tests cover production selection,
  complete dependency delivery, the exact filename, one load per invocation,
  independent repeated value/error results, complete recursive decode,
  explicit-empty and omitted slice semantics, exact dependency/path/read/
  unmarshal and partial-result behavior, safe zero defaults, and non-empty
  recording populations.
- `ListDeprecated`, Maven and command callers, refresh, `Services`,
  `ProjectDefaults`, `GlobalCloudConfig`, every other reader, construction,
  tracked fixtures, and the config-cloud mutation manifest remain byte-exact.

### Verification And Protected State

- Final-commit focused deprecated tests and all `pkg/config`, `pkg/context`,
  `pkg/maven`, and `cmd/...` tests pass. The unchanged config-cloud harness and
  T1-T10 meta-test each report exact
  `declared=10 killed=10 survived=0 unusable=0` with unchanged IDs,
  selections, and meanings.
- Exact Go 1.26.7 module verification, build, uncached count-one tests, race
  count-one tests, and vet pass under `umask 022`, offline module resolution,
  readonly project inputs, and managed scratch. Pinned API/CLI compatibility,
  both compatibility meta-tests, the CLI surface, complete `make preflight`,
  ordinary `make test`, `make test-install`, a real scratch-local install, and
  all 62 launcher controls pass. A fresh empty-HOME count-two run also passes.
- All 15 quality-audit meta-controls pass. The canonical audit has its
  documented finding exit 1 rather than audit-broken exit 2. Its structured
  Q0-Q2 scope has all 21 automated criteria PASS and zero Q0-Q2 ratchet
  regression. The copied baseline manual evidence is correctly rejected as
  stale because it is bound to the historical baseline commit; no false fresh
  manual PASS is claimed. The existing six manual contracts remain supported:
  this slice strengthens the safe-default, argument-recorder, non-empty, and
  partial-result evidence; all eight mutation harnesses retain 10/10 kills;
  acceptance magnitude and bad-input/read-only controls are unchanged.
- Go 1.18 and protected 234 modules, 3,599 graph edges, 355 production entries,
  429 complete-test entries, 197 module-backed entries across 41 loaded
  modules, and 1,067 sum lines reproduce. `go.mod`, `go.sum`, and graph hashes
  remain respectively
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  API/CLI, cache layout, refresh, callers, other readers, fixtures, mutation
  files, dependency metadata, and the Go floor remain unchanged.
- Environment-only discarded attempts reproduced known harness boundaries:
  the Homebrew launcher selected Go 1.26.2 in an empty HOME; macOS no-template
  `mktemp` ignored `TMPDIR`; passing a tool path through outer `make` polluted
  the Make meta-test through `MAKEFLAGS`; and system Xcode probes attempted
  sandbox-denied cache writes without creating artifacts. The accepted runs
  used the existing exact Go 1.26.7 binary, scratch-local pinned tools, the
  existing read-only module cache, a scratch-local no-template `mktemp`
  adapter, and direct Xcode compiler paths for race. Two ignored compatibility
  reports created by the first discarded preflight were verified as regular
  task-owned outputs and removed before the accepted clean rerun.
- All disposable evidence stayed beneath the managed task root. Before final
  removal, containment and entry types were verified across 8,043 directories
  and 46,534 regular files, with zero symlinks and zero special entries.

The rollback boundary remains the single focused two-file implementation
commit. The reciprocal successor is planning-only and may select at most one
further owned cloud slice; it may not execute that slice, inspect or select
`ply-config`, combine readers, or begin Spring or packaging work.
