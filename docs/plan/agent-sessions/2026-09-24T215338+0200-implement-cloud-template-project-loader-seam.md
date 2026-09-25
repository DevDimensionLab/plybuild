# Agent Session: Implement Cloud Template Project Loader Seam

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T215338+0200-implement-cloud-template-project-loader-seam`
Created: `2026-09-24T21:53:38+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `111ad125c4d40dfb64c7a8600080d15e81d9f063a7c28411d15963687438b3b9`
Previous: [2026-09-24T213100+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T213100+0200-plan-next-p8-cloud-modernization-move.md)
Next: [2026-09-24T225524+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T225524+0200-plan-next-p8-cloud-modernization-move.md)
Outcome: Completed the eighth bounded P8 cloud slice at focused commit `d26d19232b18d70a4401586d75604d08f85bd8a8`; the reciprocal successor is planning-only.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the eighth bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the direct
`InitProjectFromDirectory` call inside `GitCloudConfig.templates` behind one
private zero-value-safe template-project loader while preserving template
walking, project/profile/POM construction, cache mapping, partial results,
callers, and every public contract. Do not integrate or select `ply-config`, and
do not begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first seven
P8 slices are complete at `57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`,
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`,
`f12344b115b6c732295d152a7f3d86b876a4d1d8`,
`76c3571f04969c304785338bc31c34a266e2f869`, and
`9c1899d484633e0ce6cf112544ba33ca89bc4adb`. This checkpoint may change only
`pkg/config/cloud.go` and `pkg/config/cloud_templates_test.go` for the private
template-project loader seam. It may not change the project or root opener,
profile behavior, another reader, caller, fixture, mutation file, dependency
metadata, Spring, or packaging.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, answered planning archive, focused project-opener
implementation ancestry and exact two-file shape, earlier root-opener,
global-config, deprecated, project-defaults, services, and cache-probe ancestry,
Google UUID ancestry, sole NEXT state, reciprocal 342-archive chain, launcher
mirror/check, ordinary and ignored cleanliness, and unchanged protected cloud,
project, profile, context, caller, fixture, mutation, and dependency inputs. The
protected project remains Go 1.18 with 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, 1,067 `go.sum` lines, and protected `go.mod` / `go.sum` /
graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving private construction dependency seam, not
`ply-config` integration, dependency evaluation, profile-policy repair, generic
reader or construction refactoring, a shared/exported factory, caller-policy
cleanup, security work, Spring work, or packaging work. Do not change any
exported symbol, signature, or field; `InitProjectFromDirectory`;
`OpenGitCloudConfig`; either cloud opener; profile/migration/POM behavior;
cache layout; refresh; another reader; any command, context, or template caller;
tracked fixture; mutation file; `go.mod`; or `go.sum`. Do not fetch, inspect, or
select a new dependency. Keep every disposable cache, report, copy, and test
fixture beneath `${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch
after containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
answered project-opener and root-opener implementation archives, the earlier
global-config, deprecated, project-defaults, services, and cache-probe
implementation archives, current P8 roadmap and rolling handover,
`docs/design/quality-lift.md`, `pkg/config/cloud.go`,
`pkg/config/cloud_templates_test.go`, all other focused cloud tests,
`pkg/config/project_init.go` and its focused opener tests,
`pkg/config/profiles.go`, `pkg/config/project.go`, `pkg/context/context.go` and
its focused opener tests, the `CloudConfig` callers in `cmd`, `pkg/context`, and
`pkg/template`, and `scripts/mutate-config-cloud` with its meta-test. Treat every
decoded-profile, non-ENOENT, migration, caller-policy, other-reader,
real-network, environment, and fixture gap as a guard rather than permission to
broaden this slice. Do not inspect `ply-config`.

# Three Moves

1. Add focused failing characterization for production project-loader and
   filesystem selection, exact arbitrary complete template-directory delivery
   and one call per current or legacy config match, no load for ignored or
   incoming-error callbacks, complete project and embedded interface identity,
   exact error/log/partial-result behavior, independent repeated invocations,
   safe missing-loader behavior with an injected walk, and non-empty
   recordings. Preserve all current template-walk and project behavior.
2. Add one private template-project loader interface and production
   implementation, extend only the existing private `templatesDependencies`,
   and replace only the callback's direct constructor call. Production must
   delegate exactly once per match to unchanged
   `InitProjectFromDirectory(relPath[0])` and return its complete `Project` and
   exact error. A missing loader returns zero `Project` plus exact
   `filesystem.ErrNoFilesystem` without invoking construction. Preserve walk
   root, callback order, matching, path/name derivation, append order, logging,
   partial/final results, repeated-call behavior, profile/POM/cache behavior,
   public facade, and every other boundary.
3. Run focused config/context/Maven/command and exact unchanged config-cloud
   mutation gates, then the full recorded exact-Go-1.26.7 gates under
   `umask 022`, readonly/offline module inputs, and managed scratch. Verify no
   API/CLI, source-scope, opener, profile, project/POM, cache, caller, reader,
   fixture, Go-floor, dependency, graph, mutation, or quality regression. Roll
   back the single two-file commit if any contract fails; do not compensate in
   another slice.

# Automatic Handoff

If and only if the one seam and every gate pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If any contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge, publish,
release, stash, revert, remove the worktree, integrate `ply-config`, repair
profile behavior, reopen an earlier slice, alter caller policy, or start a
second cloud, Spring, or packaging slice.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The eighth bounded P8 cloud slice is complete at focused implementation commit
`d26d19232b18d70a4401586d75604d08f85bd8a8` (parent
`42857a2c39d252137bf873ab4bb761f514f5f070`, tree
`e6d265798b22dd8690f6f078d6df246332ea600c`). It changes exactly
`pkg/config/cloud.go` and `pkg/config/cloud_templates_test.go`, with 324
insertions and 47 deletions. No public symbol, constructor, opener, caller,
fixture, mutation file, dependency file, Spring path, or packaging path changed.

### Implementation And Characterization

- `templateProjectLoader` and `initTemplateProjectLoader` are private. The
  production selector extends only `templatesDependencies`, keeps the complete
  system filesystem dependency, and explicitly selects the production project
  loader. Its one method delegates exactly once to unchanged
  `InitProjectFromDirectory(directory)` and returns the complete `Project` and
  exact error.
- The existing walk callback changed only its direct constructor call to the
  private dependency call. A missing loader returns zero `Project` plus exact
  `filesystem.ErrNoFilesystem`; no project construction is attempted. Walk
  root, callback order, incoming-error suppression, current/legacy matching,
  path and relative-name derivation, append order, logging, partial/final
  results, and independent repeated calls remain unchanged.
- Focused TDD first failed because the private interface, production
  implementation, dependency field, and zero-safe dispatch did not exist.
  Final characterization covers production filesystem and loader selection;
  tracked current and legacy projects; arbitrary complete directory delivery
  exactly once per match; no load for ignored or incoming-error callbacks;
  complete project, cloud-config, and project-type interface identity; exact
  error, log, stop-order, partial-result, and final-walk-error behavior;
  independent repeated invocations; missing filesystem and missing loader;
  and rejected empty recordings.

### Verification And Protected State

- Focused exact-Go-1.26.7 template/config/context/Maven/command tests pass,
  including count-two config/context coverage. The byte-unchanged config-cloud
  direct harness and its T1-T10 meta-test each retain all ten exact IDs,
  selections, meanings, and `declared=10 killed=10 survived=0 unusable=0`.
- Exact Go 1.26.7 module verification, build, uncached count-one tests, race
  count-one tests, and vet pass under `umask 022`, offline module resolution,
  readonly project inputs, direct Xcode compilers, and managed scratch. Pinned
  API/CLI compatibility, both compatibility meta-tests, the CLI surface,
  complete `make preflight`, ordinary `make test`, explicit
  `make test-install`, a real scratch-local install/help smoke, fresh empty-HOME
  count-two tests, and all 62 launcher controls pass.
- All 15 quality-audit meta-controls pass. The canonical audit returned its
  documented findings exit 1. Its structured Q0-Q2 scope has all 21 automated
  criteria PASS and zero scoped ratchet regression; all six manual criteria
  remain correctly unclaimed. The sole overall ratchet regression remains the
  pre-existing Q3.4 documentation indicator.
- Go 1.18 and the protected 234 modules, 3,599 graph edges, 355 production
  entries, 429 complete-test entries, 197 module-backed entries across 41
  loaded modules, and 1,067 `go.sum` lines reproduce. `go.mod`, `go.sum`, and
  graph hashes remain respectively
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  API/CLI, root and project openers, project/profile/migration/POM behavior,
  cache layout, refresh, callers, other readers, fixtures, mutation files,
  dependency metadata, and the Go floor remain unchanged.
- Two discarded preflight invocations exposed existing harness environment
  contracts rather than product defects. One inherited
  `GOFLAGS=-mod=readonly`, so T11 compared that receipt with its intentional
  post-`unset GOFLAGS` receipt; their sole JSON difference was
  `/tool/go_build/goflags`. The next passed the pinned lint binary as a Make
  command-line override, which propagated into and defeated the lint
  meta-test's deliberate missing-binary probe. The accepted full rerun kept
  offline resolution, began with empty `GOFLAGS`, supplied pinned tools as
  ordinary environment inputs, and passed unchanged. No attempt changed
  repository state.
- The task-owned scratch subtree was containment and entry-type audited before
  removal: 8,509 directories, 43,890 regular files, zero symbolic links, and
  zero special entries. Its resolved path remained a strict child of the
  launcher-provided session scratch root, and only that subtree was removed.

The rollback boundary remains the single focused two-file implementation
commit. The reciprocal successor is planning-only and may select at most one
further smallest repository-owned cloud slice or stop unresolved; it may not
execute a slice, inspect or select `ply-config`, repair profile behavior, reopen
an earlier seam, alter caller policy, or begin Spring or packaging work.
