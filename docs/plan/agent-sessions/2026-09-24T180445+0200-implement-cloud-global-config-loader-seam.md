# Agent Session: Implement Cloud Global Config Loader Seam

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T180445+0200-implement-cloud-global-config-loader-seam`
Created: `2026-09-24T18:04:45+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `4125343c8df764a1421761111fea80358efd778a66461a86bd6ef218139a1265`
Previous: [2026-09-24T174851+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T174851+0200-plan-next-p8-cloud-modernization-move.md)
Next: [2026-09-24T190912+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T190912+0200-plan-next-p8-cloud-modernization-move.md)
Outcome: the fifth bounded P8 slice is complete at focused commit `f12344b115b6c732295d152a7f3d86b876a4d1d8`; the direct global-config read now has one private zero-value-safe loader with path, environment, YAML, result, formatting, caller, and dependency behavior unchanged

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the fifth bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the cache-backed
`global-config.yaml` read behind one private zero-value-safe loader while
preserving direct path, environment-expansion, YAML, result, source-formatting,
caller, and dependency contracts. Do not integrate or select `ply-config`, and
do not begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first four
P8 slices are complete at `57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`, and
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`. This checkpoint may change only
`pkg/config/cloud.go` and `pkg/config/cloud_test.go` for the private global-
config loader seam. It may not change a caller, refresh, another reader,
constructor, tracked fixture, mutation file, dependency metadata, Spring, or
packaging.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, answered planning archive, focused deprecated
implementation ancestry and exact two-file shape, earlier project-defaults,
services, and cache-probe ancestry, Google UUID ancestry, sole NEXT state,
reciprocal 336-archive chain, launcher mirror/check, ordinary and ignored
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
packaging work. Do not change any exported symbol, signature, or field;
`GlobalCloudConfig.SourceFor`; cache layout; refresh; `Services`,
`ProjectDefaults`, `Deprecated`, Git hooks, templates, or examples; any command
or context caller; `Directory` or `DirConfig.FilePath`; tracked fixture;
mutation file; `go.mod`; or `go.sum`. Do not fetch, inspect, or select a new
dependency. Keep every disposable cache, report, copy, and test fixture beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch after
containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
answered deprecated, project-defaults, services, and cache-probe implementation
archives, current P8 roadmap and rolling handover,
`docs/design/quality-lift.md`, `pkg/config/cloud.go`,
`pkg/config/cloud_test.go`, `pkg/config/dir.go`, `pkg/config/types.go`,
`pkg/file/file.go`'s `Open` path, the `GlobalCloudConfig` consumers in
`cmd/build_options.go` and `cmd/tips.go`, construction in `pkg/context` and
`pkg/config/project_init.go`, and `scripts/mutate-config-cloud` with its meta-
test. Treat every environment, consumer, real-network, and missing-fixture gap
as a guard rather than permission to broaden the slice.

# Three Moves

1. Add focused failing characterization for production loader selection,
   complete `Directory` delivery, exact direct `global-config.yaml` path
   composition without `FilePath`, one independent load per invocation, exact
   dependency/read/YAML and partial-result semantics, environment expansion
   before decode, safe zero value, unchanged `SourceFor`, and non-empty
   recordings. Do not characterize or change another reader or caller.
2. Add one private global-config loader interface and private dependency value.
   The production implementation accepts the complete `Directory`, composes
   exactly `file.Path("%s/global-config.yaml", directory.Dir())`, delegates to
   unchanged `file.Open`, applies `[]byte(os.ExpandEnv(string(b)))` before
   unchanged `yaml.Unmarshal`, and returns the exact value/error. The private
   helper loads once per invocation; a missing loader returns
   `filesystem.ErrNoFilesystem` before developer-path access. Preserve the
   exported wrapper, direct-path cache-only behavior, and every other boundary.
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
decision successor instead. Do not execute the successor, push, merge, publish,
release, stash, revert, remove the worktree, integrate `ply-config`, or start a
second cloud, Spring, or packaging slice.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The fifth bounded P8 cloud slice is complete at focused implementation commit
`f12344b115b6c732295d152a7f3d86b876a4d1d8` (parent
`82f100de434a37acd94ce16cd8bfba41f7f08e36`, tree
`e4e28bf42fe484904bb92d8ff4f155878594591c`). It changes exactly
`pkg/config/cloud.go` and `pkg/config/cloud_test.go`. No second reader, caller,
refresh path, constructor, tracked fixture, mutation file, dependency, Spring,
packaging, or `ply-config` work ran.

### Implementation And Characterization

- `globalConfigLoader` and `globalConfigDependencies` are private. Their zero
  value returns exact `filesystem.ErrNoFilesystem` with a zero
  `GlobalCloudConfig` before accessing the supplied `Directory`; production
  explicitly selects `fileGlobalConfigLoader`.
- The production loader receives the complete `Directory`, calls `Dir()` once,
  composes exact `file.Path("%s/global-config.yaml", directory.Dir())`, never
  calls `FilePath`, and delegates to unchanged `file.Open`. It applies exact
  `[]byte(os.ExpandEnv(string(b)))` to the complete successful read before
  unchanged `yaml.Unmarshal`, preserving raw read errors, decoded and partial
  values, YAML errors, and zero values.
- The exported `GlobalCloudConfig()` signature remains the production wrapper.
  Its private helper performs exactly one load per invocation; separate calls
  re-read the file and re-expand the live environment independently. No
  closure, memoization, refresh, existence probe, or cross-call state was
  introduced.
- Focused TDD first failed on the absent private interface, dependency
  selector, file loader, and helper. Final characterization covers production
  selection, complete `Directory` delivery, exact direct path and no
  `FilePath`, exact dependency results, raw read errors, complete environment-
  before-YAML decode, representative partial YAML value/error, independent
  repeated file/environment changes, safe zero behavior, byte-exact unchanged
  `SourceFor` slash formatting, and rejected empty recording populations.
- `SourceFor`, template/tips consumers, `Directory`, `DirConfig.FilePath`,
  refresh, `Services`, `ProjectDefaults`, `Deprecated`, Git hooks, templates,
  examples, construction, tracked fixtures, and the config-cloud mutation
  manifest remain byte-exact.

### Verification And Protected State

- Final-commit focused global-config tests and all `pkg/config`, `pkg/context`,
  `pkg/maven`, and `cmd/...` tests pass. The unchanged config-cloud direct
  harness and T1-T10 meta-test each retain exact
  `declared=10 killed=10 survived=0 unusable=0` with unchanged IDs,
  selections, and meanings.
- Direct exact Go 1.26.7 module verification, build, uncached count-one tests,
  race count-one tests, and vet pass under `umask 022`, offline module
  resolution, readonly project inputs, and managed scratch. Pinned API/CLI
  compatibility, both compatibility meta-tests, the CLI surface, complete
  `make preflight`, ordinary `make test`, `make test-install`, a real scratch-
  local install, fresh empty-HOME count-two tests, and all 62 launcher controls
  pass.
- All 15 quality-audit meta-controls pass. The canonical audit has its
  documented finding exit 1 rather than audit-broken exit 2. Its structured
  Q0-Q2 scope has all 21 automated criteria PASS and zero scoped ratchet
  regression. The copied baseline manual evidence is correctly stale; the sole
  overall ratchet regression remains the pre-existing Q3.4 documentation
  indicator. The six manual contracts remain supported, including stronger
  safe-default, argument-recording, non-empty, and partial-result evidence.
- Go 1.18 and protected 234 modules, 3,599 graph edges, 355 production entries,
  429 complete-test entries, 197 module-backed entries across 41 loaded
  modules, and 1,067 sum lines reproduce. `go.mod`, `go.sum`, and graph hashes
  remain respectively
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  API/CLI, cache layout, refresh, callers, other readers, fixtures, mutation
  files, dependency metadata, and the Go floor remain unchanged.
- Discarded environment-only attempts exposed no product defect: the Homebrew
  launcher resolved its underlying Go 1.26.2 for T15 when forced local, while
  the accepted preflight placed the installed Go 1.26.7 toolchain binary first
  on `PATH`; a first canonical post-check misclassified stale manual rows as
  automated before the structured criterion IDs were checked correctly.
  Sandbox-denied Xcode cache diagnostics from recording install doubles created
  no outside artifact and their exact leak guards passed.
- Every disposable artifact stayed beneath the managed task root. Before final
  removal, containment and entry types were verified across 5,314 directories
  and 52,672 regular files, with three intentional diagnostic symlinks and zero
  special entries.

The rollback boundary remains the single focused two-file implementation
commit. The reciprocal successor is planning-only and may select at most one
further owned cloud slice; it may not execute that slice, inspect or select
`ply-config`, combine readers, or begin Spring or packaging work.
