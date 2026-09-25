# Agent Session: Implement Project Cloud Config Opener Seam

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T204758+0200-implement-project-cloud-config-opener-seam`
Created: `2026-09-24T20:47:58+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `b92375f465a2132fff86c4a38e6e7b5b97e2c24da7b44e1434f72fcf4c90ccc9`
Previous: [2026-09-24T202938+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T202938+0200-plan-next-p8-cloud-modernization-move.md)
Next: [2026-09-24T213100+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T213100+0200-plan-next-p8-cloud-modernization-move.md)
Outcome: the seventh bounded P8 slice is complete at focused commit `9c1899d484633e0ce6cf112544ba33ca89bc4adb`; project construction now selects unchanged Git cloud configuration through one private zero-value-safe opener while every protected profile, migration, decoded-profile, project/POM, caller, cache, mutation, and public contract remains unchanged

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the seventh bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the final
`InitProjectFromDirectory` construction of `CloudConfig` behind one private
zero-value-safe opener while preserving current profile discovery, migration,
decoded-profile behavior, project/POM behavior, callers, cache mapping, and
every public contract. Do not integrate or select `ply-config`, and do not begin
a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first six
P8 slices are complete at `57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`,
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`,
`f12344b115b6c732295d152a7f3d86b876a4d1d8`, and
`76c3571f04969c304785338bc31c34a266e2f869`. This checkpoint may change only
`pkg/config/project_init.go` and a new focused
`pkg/config/project_init_cloud_test.go` for the private project cloud-config
opener seam. It may not reopen the root opener, change profile behavior, alter a
caller or reader, modify fixtures or mutation files, change dependency
metadata, or begin Spring or packaging work.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, answered planning archive, focused root-opener
implementation ancestry and exact two-file shape, earlier global-config,
deprecated, project-defaults, services, and cache-probe ancestry, Google UUID
ancestry, sole NEXT state, reciprocal 340-archive chain, launcher mirror/check,
ordinary and ignored cleanliness, and unchanged protected context, cloud,
project, profile, caller, fixture, mutation, and dependency inputs. The
protected project remains Go 1.18 with 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, 1,067 `go.sum` lines, and protected `go.mod` / `go.sum` /
graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving private construction seam, not `ply-config`
integration, dependency evaluation, profile-policy repair, generic construction
refactoring, a shared/exported factory, caller-policy cleanup, reader migration,
security work, Spring work, or packaging work. Do not change any exported
symbol, signature, or field; `OpenGitCloudConfig`; `GitCloudConfig`;
`CloudConfig`; `Context.LoadProfile`; `profiles.go`; project/profile/POM
ordering or results; cache layout; refresh; any command caller; fixture;
mutation file; `go.mod`; or `go.sum`. In particular, preserve the current
source behavior in which decoded `projectConfig.Profile` is not assigned to
`project.Config` before the existing profile check; do not fix or normalize
that gap. Do not fetch, inspect, or select a new dependency. Keep every
disposable cache, report, copy, and test fixture beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch after
containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
answered root-opener, global-config, deprecated, project-defaults, services,
and cache-probe implementation archives, current P8 roadmap and rolling
handover, `docs/design/quality-lift.md`, `pkg/config/project_init.go`,
`pkg/config/profiles.go`, `pkg/config/project.go`, existing project and focused
cloud tests, `pkg/context/context.go` and its focused opener tests,
`pkg/config/cloud.go`, the `InitProjectFromDirectory` and `CloudConfig` callers
in `cmd`, `pkg/context`, and `pkg/template`, and
`scripts/mutate-config-cloud` with its meta-test. Treat every current decoded-
profile, non-ENOENT, migration, caller-policy, other-reader, real-network,
environment, and fixture gap as a guard rather than permission to broaden the
slice. Do not inspect `ply-config`.

# Three Moves

1. Add focused failing characterization for Git production selection, exact
   complete resolved-profile-path delivery and one call, exact returned
   interface identity, public active-profile and missing-profile migration
   cache mapping, current decoded-profile non-override behavior, safe zero
   behavior, early project-configuration failure before opening, independent
   repeated construction, and non-empty recordings. Preserve existing project,
   POM, profile, migration, error, and caller behavior.
2. Add one private cloud-config opener interface and private dependency value
   in `pkg/config`, one Git-backed production implementation, and one small
   private dependency-taking opener helper. Production delegates exactly once
   to unchanged `OpenGitCloudConfig(profilePath)` and assigns the returned
   complete `CloudConfig` interface. A missing opener returns nil without
   touching the supplied path. Replace only the direct construction assignment;
   do not move or normalize the surrounding constructor flow.
3. Run focused config/context/Maven/command and exact unchanged config-cloud
   mutation gates, then the full recorded exact-Go-1.26.7 gates under
   `umask 022`, readonly/offline module inputs, and managed scratch. Verify no
   API/CLI, source-scope, root-opener, profile, project/POM, caller, reader,
   fixture, Go-floor, dependency, graph, mutation, or quality regression. Roll
   back the single two-file commit if any contract fails; do not compensate in
   another slice.

# Automatic Handoff

If and only if the one seam and every gate pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If a contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge, publish,
release, stash, revert, remove the worktree, integrate `ply-config`, repair
decoded-profile behavior, reopen the root opener, or start a second cloud,
Spring, or packaging slice.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The seventh bounded P8 cloud slice is complete at focused implementation
commit `9c1899d484633e0ce6cf112544ba33ca89bc4adb` (parent
`c5fec9cce330094b0628b7ebc7a60e1c6f91d474`, tree
`1083cb67f52b347a5f0760a4ddda7cf533c2a0d4`). It changes exactly
`pkg/config/project_init.go` and new focused
`pkg/config/project_init_cloud_test.go`. The root opener, cloud readers,
profile helpers, project/POM types, callers, fixtures, mutation files,
dependency metadata, Spring, packaging, and `ply-config` remain untouched.

### Implementation And Characterization

- `projectCloudConfigOpener`, `projectCloudConfigDependencies`, and
  `gitProjectCloudConfigOpener` are private. Production explicitly selects the
  Git-backed opener, which delegates exactly once to unchanged
  `OpenGitCloudConfig(profilePath)` and returns the complete result as
  `CloudConfig`. The small dependency-taking helper returns nil before path
  access when its opener is absent.
- `InitProjectFromDirectory` retains its body, public signature, and exact
  ordering. Only the direct cloud-construction assignment now calls the private
  production helper. Configuration load/populate, active-profile lookup,
  text-matched missing-file migration/retry, legacy-home selection, named
  errors, cloud-before-POM assignment, POM warning/partial behavior, and final
  `ConfigFile`, `Path`, and `Config` assignments remain in place.
- Focused TDD first failed because the private interface, dependency value,
  production selector, and helper did not exist. Final characterization covers
  Git production selection, arbitrary complete resolved-profile-path delivery
  exactly once, exact returned interface identity, public active-profile and
  missing-profile migration cache mapping, safe zero behavior without path
  access, early project-configuration failure before cloud assignment,
  independent repeated construction, and rejected empty recorder populations.
- The decoded project configuration still remains local until the final
  assignment. The existing branch therefore still checks zero
  `project.Config.Profile`; a decoded project profile is returned in
  `project.Config` but does not override the active profile used for cloud
  cache mapping. The focused test freezes this source-visible gap without
  repairing or normalizing it.

### Verification And Protected State

- Final-commit focused config/context/Maven/command tests pass, including
  count-two config/context runs. The byte-unchanged config-cloud direct harness
  retains all ten exact IDs, selections, and meanings with
  `declared=10 killed=10 survived=0 unusable=0`; its T1-T10 meta-test passes
  with the same totals.
- Direct exact Go 1.26.7 module verification, build, uncached count-one tests,
  race count-one tests, and vet pass under `umask 022`, offline module
  resolution, readonly project inputs, direct Xcode compilers, and managed
  scratch. Pinned API/CLI compatibility, both compatibility meta-tests, the CLI
  surface, complete preflight, ordinary `make test`, explicit
  `make test-install`, a real scratch-local install/help smoke, fresh empty-HOME
  count-two tests, and all 62 launcher controls pass.
- All 15 quality-audit meta-controls pass. The accepted canonical audit has its
  documented findings exit 1 rather than audit-broken exit 2. Its structured
  Q0-Q2 scope has all 21 automated criteria PASS and zero scoped ratchet
  regression; the six manual criteria remain correctly unclaimed. The sole
  overall ratchet regression remains the pre-existing Q3.4 documentation
  indicator.
- Go 1.18 and protected 234 modules, 3,599 graph edges, 355 production entries,
  429 complete-test entries, 197 module-backed entries across 41 loaded
  modules, and 1,067 sum lines reproduce. `go.mod`, `go.sum`, and graph hashes
  remain respectively
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  API/CLI, root opener, project/profile/POM behavior, cache layout, refresh,
  callers, readers, fixtures, mutation files, dependency metadata, and the Go
  floor remain unchanged.
- Two discarded environment-only attempts exposed existing guard behavior, not
  product defects: the first preflight used macOS no-template `mktemp`, which
  ignored managed `TMPDIR` and was denied before creating an outside artifact;
  the accepted rerun used the scratch-local adapter and direct Xcode `make`.
  The first canonical audit correctly rejected ambient `GOFLAGS=-mod=readonly`
  as a baseline build-context mismatch with exit 2; the fresh accepted run
  kept offline resolution, matched the baseline's empty `GOFLAGS`, and returned
  the required findings exit 1. Neither attempt changed repository state.
- A first handoff-only rerun used a restricted `PATH` that hid the installed
  Codex executable and then hit the signal-retention control's timing
  assertion after 24 passes. Restoring the executable path and immediately
  rerunning the unmodified handoff produced a passing launcher check and all
  62 contract controls; repository state remained unchanged.
- The contained task root was audited before removal: 7,782 directories,
  45,343 regular files, zero symbolic links, and zero special entries. Its
  resolved path remained a strict `project-cloud-opener.*` child of the
  launcher-provided session scratch root; the exact subtree was then removed.

The rollback boundary remains the single focused two-file implementation
commit. The reciprocal successor is planning-only and may select at most one
further owned cloud slice; it may not execute that slice, inspect or select
`ply-config`, repair decoded-profile behavior, reopen an earlier slice, combine
construction moves, alter caller policy, or begin Spring or packaging work.
