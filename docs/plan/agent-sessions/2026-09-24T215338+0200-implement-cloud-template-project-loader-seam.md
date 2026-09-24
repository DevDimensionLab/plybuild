# Agent Session: Implement Cloud Template Project Loader Seam

Status: NEXT
Session ID: `2026-09-24T215338+0200-implement-cloud-template-project-loader-seam`
Created: `2026-09-24T21:53:38+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `111ad125c4d40dfb64c7a8600080d15e81d9f063a7c28411d15963687438b3b9`
Previous: [2026-09-24T213100+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T213100+0200-plan-next-p8-cloud-modernization-move.md)
Next: none
Outcome: pending

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
