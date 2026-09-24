# Agent Session: Implement Project Cloud Config Opener Seam

Status: NEXT
Session ID: `2026-09-24T204758+0200-implement-project-cloud-config-opener-seam`
Created: `2026-09-24T20:47:58+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `b92375f465a2132fff86c4a38e6e7b5b97e2c24da7b44e1434f72fcf4c90ccc9`
Previous: [2026-09-24T202938+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T202938+0200-plan-next-p8-cloud-modernization-move.md)
Next: none
Outcome: pending

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
