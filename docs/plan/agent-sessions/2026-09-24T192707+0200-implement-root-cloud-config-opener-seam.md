# Agent Session: Implement Root Cloud Config Opener Seam

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T192707+0200-implement-root-cloud-config-opener-seam`
Created: `2026-09-24T19:27:07+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `997909146ef2e03f8162465aaf57a85a2408b08705ef1bc93898762e94470c0a`
Previous: [2026-09-24T190912+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T190912+0200-plan-next-p8-cloud-modernization-move.md)
Next: [2026-09-24T202938+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T202938+0200-plan-next-p8-cloud-modernization-move.md)
Outcome: the sixth bounded P8 slice is complete at focused commit `76c3571f04969c304785338bc31c34a266e2f869`; root profile loading now selects its unchanged Git cloud configuration through one private zero-value-safe opener with every protected lifecycle and caller contract unchanged

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the sixth bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the root `Context.LoadProfile`
construction of `config.CloudConfig` behind one private zero-value-safe opener
while preserving concrete Git selection, cache mapping, local-config lifecycle,
project construction, caller policy, and every public contract. Do not integrate
or select `ply-config`, and do not begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first five
P8 slices are complete at `57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`,
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`, and
`f12344b115b6c732295d152a7f3d86b876a4d1d8`. This checkpoint may change only
`pkg/context/context.go` and `pkg/context/context_test.go` for the private root
cloud-config opener seam. It may not route project construction, add an exported
factory, change a caller or reader, modify mutation files or fixtures, change
dependency metadata, or begin Spring or packaging work.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, answered planning archive, focused global-config
implementation ancestry and exact two-file shape, earlier deprecated, project-
defaults, services, and cache-probe ancestry, Google UUID ancestry, sole NEXT
state, reciprocal 338-archive chain, launcher mirror/check, ordinary and ignored
cleanliness, and unchanged protected source, callers, fixtures, mutation files,
and dependency metadata. The protected project remains Go 1.18 with 234 modules,
3,599 graph edges, 355 production entries, 429 complete-test entries, 197
module-backed entries over 41 loaded modules, 1,067 `go.sum` lines, and
protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving private construction seam, not `ply-config`
integration, dependency evaluation, a shared/exported factory, project-
construction refactoring, caller-policy cleanup, reader migration, security
work, Spring work, or packaging work. Do not change any exported symbol,
signature, or field; `config.OpenGitCloudConfig`; `GitCloudConfig`;
`CloudConfig`; `InitProjectFromDirectory`; cache layout; local-config lifecycle;
refresh; any reader or command caller; fixture; mutation file; `go.mod`; or
`go.sum`. Do not fetch, inspect, or select a new dependency. Keep every
disposable cache, report, copy, and test fixture beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch after
containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
answered global-config, deprecated, project-defaults, services, and cache-probe
implementation archives, current P8 roadmap and rolling handover,
`docs/design/quality-lift.md`, `pkg/context/context.go` and its focused tests,
`pkg/config/cloud.go`, `pkg/config/project_init.go`, the `CloudConfig` callers
in `cmd`, and `scripts/mutate-config-cloud` with its meta-test. Treat every
project-construction, caller-policy, real-network, environment, and fixture gap
as a guard rather than permission to broaden the slice. Do not inspect
`ply-config`.

# Three Moves

1. Add focused failing characterization for production opener selection, exact
   complete profile-path delivery and one call, exact returned interface
   identity, public production cache mapping, safe zero behavior, and non-empty
   recordings. Retain every existing local-config assignment, creation,
   content, mode, log, and touch-error contract.
2. Add one private cloud-config opener interface and private dependency value
   in `pkg/context`, one Git-backed production implementation, and a private
   `loadProfile` helper. Production delegates exactly once to unchanged
   `config.OpenGitCloudConfig(profilePath)` and assigns the returned complete
   interface. A missing opener returns nil without touching the supplied path.
   Preserve the exported wrapper, exact local-before-cloud assignment and
   existence/touch/log ordering, project construction, callers, and every other
   boundary.
3. Run focused context/config/Maven/command and exact unchanged config-cloud
   mutation gates, then the full recorded exact-Go-1.26.7 gates under
   `umask 022`, readonly/offline module inputs, and managed scratch. Verify no
   API/CLI, source-scope, project-construction, caller, reader, fixture, Go-floor,
   dependency, graph, mutation, or quality regression. Roll back the single
   two-file commit if any contract fails; do not compensate in another slice.

# Automatic Handoff

If and only if the one seam and every gate pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If a contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge, publish,
release, stash, revert, remove the worktree, integrate `ply-config`, route the
project constructor, or start a second cloud, Spring, or packaging slice.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The sixth bounded P8 cloud slice is complete at focused implementation commit
`76c3571f04969c304785338bc31c34a266e2f869` (parent
`75b2e249e6c9665e280902c709238fbfe412aa42`, tree
`3fc1887b116fbb0b56ca44a26a4b8eff7bd64898`). It changes exactly
`pkg/context/context.go` and `pkg/context/context_test.go`. Project
construction, callers, readers, fixtures, mutation files, dependency metadata,
Spring, packaging, and `ply-config` remain untouched.

### Implementation And Characterization

- `cloudConfigOpener`, `loadProfileDependencies`, and
  `gitCloudConfigOpener` are private. Production explicitly selects the
  Git-backed opener, which delegates directly to unchanged
  `config.OpenGitCloudConfig(profilePath)` and returns the complete result as
  `config.CloudConfig`. A zero dependency returns nil before invoking an
  opener or touching the supplied path.
- Public `Context.LoadProfile(profilePath string)` keeps its exact signature
  and is the production wrapper around private `loadProfile`. The helper still
  assigns `config.OpenLocalConfig(profilePath)` first and the cloud interface
  second, then preserves exact existence, debug-log, `TouchFile`, and error-log
  ordering. Missing-parent errors, default bytes and mode, existing bytes, and
  repeated public construction behavior remain unchanged.
- Focused TDD first failed on the absent production dependency selector,
  Git-backed implementation, helper, and safe dependency value. Final tests
  cover concrete production selection, exact public `<profile>/cloud-config`
  mapping, complete arbitrary profile-path delivery exactly once, exact
  returned interface identity, safe zero behavior, and rejected empty recorder
  populations while retaining every earlier local-config lifecycle case.
- The first committed test draft exposed its own package-logger leak under the
  required count-two gate. The new test was moved into the existing
  state-restoring subtest wrapper and the same focused commit was amended before
  acceptance. Final `pkg/context -count=2` and the fresh full empty-HOME
  count-two suite pass.
- `config.OpenGitCloudConfig`, `GitCloudConfig`, the 15-method `CloudConfig`
  interface, `InitProjectFromDirectory`, project-selected profiles, cache
  layout, refresh, every command caller, and every cached reader are byte-exact
  against the implementation parent. Pinned API and CLI compatibility reports
  pass.

### Verification And Protected State

- Final-commit `pkg/context`, `pkg/config`, `pkg/maven`, and `cmd/...` tests
  pass. The byte-unchanged config-cloud direct harness retains all ten exact
  IDs, selections, and meanings with
  `declared=10 killed=10 survived=0 unusable=0`; its T1-T10 meta-test passes
  with the same totals.
- Direct exact Go 1.26.7 module verification, build, uncached count-one tests,
  race count-one tests, and vet pass under `umask 022`, offline module
  resolution, readonly project inputs, direct Xcode compilers, and managed
  scratch. Pinned API/CLI compatibility and both compatibility meta-tests,
  complete `make preflight`, ordinary `make test`, explicit
  `make test-install`, a real scratch-local install, fresh empty-HOME count-two
  tests, and all 62 launcher controls pass.
- All 15 quality-audit meta-controls pass. The canonical audit has its expected
  finding exit 1 rather than audit-broken exit 2. Its structured Q0-Q2 scope
  has all 21 automated criteria PASS and zero scoped ratchet regression; the
  six manual criteria remain correctly unclaimed. The sole overall regression
  remains the pre-existing Q3.4 documentation indicator.
- Go 1.18 and protected 234 modules, 3,599 graph edges, 355 production entries,
  429 complete-test entries, 197 module-backed entries across 41 loaded
  modules, and 1,067 sum lines reproduce. `go.mod`, `go.sum`, and graph hashes
  remain respectively
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
- Discarded environment-only attempts were not accepted as evidence: the
  Homebrew launcher reported its underlying Go 1.26.2 under
  `GOTOOLCHAIN=local`, two default compatibility reports briefly appeared in
  ignored `target/` before exact removal and scratch redirection, and an audit
  wrapper initially required a markdown artifact that the valid finding exit
  deliberately invalidates. None changed product state or the final gates.
- The contained 57,730-entry task root was audited before removal: 8,065
  directories, 49,665 regular files, zero symlinks, and zero special entries.
  Its resolved root stayed beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`; the exact
  task-owned subtree was then removed completely before the handoff commit.

The rollback boundary remains the single focused two-file implementation
commit. The reciprocal successor is planning-only and may select at most one
further owned cloud slice; it may not execute that slice, inspect or select
`ply-config`, route the project constructor, combine construction moves, alter
caller policy, or begin Spring or packaging work.
