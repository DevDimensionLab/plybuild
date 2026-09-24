# Agent Session: Plan Next P8 Cloud Modernization Move

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T202938+0200-plan-next-p8-cloud-modernization-move`
Created: `2026-09-24T20:29:38+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `282d37ad128571492a409204effb3030f5216185b5dc1668be29aa22751f67f1`
Previous: [2026-09-24T192707+0200-implement-root-cloud-config-opener-seam.md](2026-09-24T192707+0200-implement-root-cloud-config-opener-seam.md)
Next: [2026-09-24T204758+0200-implement-project-cloud-config-opener-seam.md](2026-09-24T204758+0200-implement-project-cloud-config-opener-seam.md)
Outcome: one private zero-value-safe `InitProjectFromDirectory` cloud-config opener is fully owned as the seventh P8 implementation slice; current profile discovery, migration, decoded-profile, project, Maven, caller, cache, and public behavior remain frozen

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Plan exactly one next bounded P8 cloud-configuration modernization slice after
the completed private root cloud-config opener seam. Re-characterize the
resulting cloud boundary, compare the smallest repository-owned next seams, and
select one implementation slice or stop unresolved. Do not implement the slice
in this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 remains active in cloud configuration only. The first
six P8 implementation slices are complete at focused commits
`57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`,
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`,
`f12344b115b6c732295d152a7f3d86b876a4d1d8`, and
`76c3571f04969c304785338bc31c34a266e2f869`. This checkpoint may plan only
one next cloud-configuration slice. It may not change source or dependency
metadata, integrate or select `ply-config`, begin Spring or packaging work,
reopen an earlier slice, or combine multiple implementation moves.

# Measurements At Start

Begin only from the clean reciprocal implementation handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused root-opener implementation parent and exact
two-file shape, earlier global-config, deprecated, project-defaults, services,
and cache-probe ancestry, Google UUID ancestry, sole NEXT state, reciprocal
339-archive chain, launcher mirror/check, ordinary and ignored cleanliness,
and unchanged protected readers, callers, fixtures, mutation files, project
construction, and dependency metadata. The protected project remains Go 1.18
with 234 modules, 3,599 graph edges, 355 production entries, 429 complete-test
entries, 197 module-backed entries over 41 loaded modules, 1,067 `go.sum`
lines, and protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is a planning checkpoint, not implementation, dependency evaluation,
caller-policy cleanup, generic reader migration, construction refactoring,
security work, Spring work, or packaging work. Use only repository source,
tests, fixtures, design records, and local Git history. Do not modify production
code, tests, fixtures, mutation files, `go.mod`, or `go.sum`; do not fetch,
inspect, or select a new dependency; and do not claim real cloud compatibility
beyond recorded evidence. Keep every disposable artifact beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`, verify containment and entry types, and
remove task-owned scratch before handoff.

# Required Reading

Read the answered root-opener implementation archive and its planning
predecessor, the answered global-config, deprecated, project-defaults,
services, and cache-probe implementation archives, current P8 roadmap and
rolling handover, `docs/design/quality-lift.md`, `pkg/context/context.go` and
its focused tests, `pkg/config/cloud.go` and all focused cloud tests,
`pkg/config/project_init.go`, the `CloudConfig` callers in `cmd`, and
`scripts/mutate-config-cloud` with its meta-test. Treat every recorded project-
construction, caller-policy, other-reader, real-network, environment, and
fixture gap as a guard rather than permission to combine work. Do not inspect
`ply-config`.

# Three Moves

1. Confirm the completed root opener preserves explicit Git production
   selection, exact complete profile-path delivery and one call, returned
   interface identity, `<profile>/cloud-config` mapping, zero-value safety,
   local-before-cloud assignment, existence/touch/log ordering, repeated-call
   behavior, and all public, project-construction, caller, reader, fixture,
   mutation, Go-floor, and dependency contracts.
2. Map only the remaining cloud boundary and run focused read-only tests needed
   to validate candidate ownership. Compare candidates by file count, public
   surface, environment and dependency knowledge, cache compatibility,
   reversibility, and the recorded characterization gaps. Do not fetch or
   inspect `ply-config`.
3. Select exactly one smallest repository-owned implementation slice with its
   files, behavior contract, tests, mutation impact, full gates, rollback, and
   expiry conditions, or stop unresolved. Do not implement it or combine a
   second move.

# Automatic Handoff

If and only if one next slice is fully owned and bounded, answer this archive,
update the roadmap and rolling handover, prepare exactly one reciprocal NEXT
archive for that implementation, replace only launcher mutable regions, run
launcher/handoff checks, and make one local handoff commit. If ownership or
behavior is unresolved, prepare one decision successor instead. Do not execute
the successor, push, merge, publish, release, stash, revert, remove the
worktree, integrate `ply-config`, or begin Spring or packaging work.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The next P8 slice is owned and bounded: isolate only the final
`InitProjectFromDirectory` construction of `CloudConfig` behind one private,
zero-value-safe opener while preserving the existing Git selection and every
current project-construction behavior. The decoded project-profile branch is a
source-visible characterization gap and is explicitly not repaired in this
slice.

### Protected Start And Focused Evidence

- The checkpoint began clean on `codex/upgrade-quality` at handoff HEAD
  `64ebd292672c80bad207fc7325345aefc78e42c9`, direct parent
  `76c3571f04969c304785338bc31c34a266e2f869`, tree
  `2ce627a2f845c616987ecfc1083532dd8ab0c0d1`, with exactly the launcher,
  answered root-opener implementation archive, this then-NEXT archive,
  rolling handover, and roadmap changed. The root-opener implementation has
  exact parent `75b2e249e6c9665e280902c709238fbfe412aa42`, tree
  `3fc1887b116fbb0b56ca44a26a4b8eff7bd64898`, and exact two-file set
  `pkg/context/context.go` and `pkg/context/context_test.go`.
- Cache-probe commit `57f9d5674157d238a2f93462b65161a17e3b5498`, services
  commit `c12307a078a29af74163df0c658d05c821482a33`, project-defaults
  commit `c999212d266d98868930769108e4e8a73070a0a2`, deprecated commit
  `51da9bcfc0eb8ae871c98467c52c942db6b480bd`, global-config commit
  `f12344b115b6c732295d152a7f3d86b876a4d1d8`, and Google UUID commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` are ancestral. The connected
  reciprocal graph had 339 archives and exactly one NEXT archive;
  launcher/archive prompt mirroring, launcher `--check`, shell syntax, and
  ordinary/ignored cleanliness passed.
- Every protected reader, caller, fixture, mutation file, project-construction
  file, context file, and dependency file is byte-identical between the
  focused root implementation and its documentation handoff. Exact Go 1.26.7
  reproduced declared Go 1.18, 234 modules, 3,599 graph edges, 355 production
  entries, 429 complete-test entries, 197 module-backed entries over 41 loaded
  modules, and 1,067 `go.sum` lines. `go.mod`, `go.sum`, and graph SHA-256
  values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
- Under `umask 022`, exact-Go-1.26.7 tests for `pkg/context`, `pkg/config`,
  `pkg/maven`, and `cmd/...` pass with offline, readonly module inputs and
  contained scratch. The byte-unchanged config-cloud direct harness retains
  all ten exact IDs, selections, and meanings with
  `declared=10 killed=10 survived=0 unusable=0`; its T1-T10 meta-test passes
  with the same totals.
- The prepared reciprocal graph has 340 archives and one NEXT implementation
  archive. Launcher/archive prompt mirroring, launcher `--check`, shell syntax,
  exact five-file documentation/launcher scope, diff checks, and all 62
  launcher lifecycle controls pass. The task-owned scratch subtree remained
  contained and comprised 259 directories and 3,382 regular files with zero
  symlinks or special entries before complete removal.

### Resulting Cloud Boundary

- The completed root opener explicitly selects `gitCloudConfigOpener`, passes
  the complete profile path exactly once to unchanged
  `config.OpenGitCloudConfig`, returns the same `config.CloudConfig` interface
  identity, and retains exact `<profile>/cloud-config` mapping. Its zero value
  returns nil without path access. `LoadProfile` still assigns local before
  cloud configuration and preserves existence, touch, logs, errors, contents,
  mode, and repeated-call behavior.
- Eight cache read/effect paths and the refresh Git dependency remain private
  seams. The public `CloudConfig` interface still has 15 methods; all exported
  cloud types, `GitCloudConfig.Impl`, `OpenGitCloudConfig`,
  `GlobalCloudConfig.SourceFor`, readers, refresh, fixtures, and command policy
  are unchanged.
- Exactly one production call site still selects a concrete cloud
  implementation: `pkg/config.InitProjectFromDirectory`. It loads and
  populates project configuration, resolves the active profile, migrates and
  retries only when the error text contains `no such file or directory`, then
  calls `OpenGitCloudConfig(profilePath)` before its POM probe and final project
  field assignments. A missing `.ply` home retains the current legacy
  `.co-pilot` fallback through `GetPlyHomePath`.
- Source re-characterization corrects an earlier overstatement without changing
  behavior: the decoded value is held in local `projectConfig`, but the branch
  checks `project.Config.Profile` before `project.Config = projectConfig` at the
  end. A non-empty `profile` in the decoded project therefore does not
  currently override the active profile. Non-ENOENT active-profile errors also
  remain part of the legacy named-error/control flow. These construction gaps
  are guards; this slice must neither fix nor normalize them.

### Candidate Comparison And Selection

| Candidate | Files and public surface | Environment, dependency, and cache knowledge | Reversibility and characterization risk | Decision |
| --- | --- | --- | --- | --- |
| Private project cloud-config opener | `pkg/config/project_init.go` plus new focused `pkg/config/project_init_cloud_test.go`; no exported change | Receives only the already-resolved profile path and delegates to unchanged `OpenGitCloudConfig`; preserves current HOME/profile/migration logic and `<profile>/cloud-config` | One-commit rollback; focused tests can freeze exact path, identity, call count, active/migration behavior, and the decoded-profile gap without changing them | **Selected** |
| Private template project loader | `pkg/config/cloud.go` and `pkg/config/cloud_templates_test.go`; no exported change | Would wrap `InitProjectFromDirectory` during template walking | The filesystem path is already injected and characterized; this does not isolate the remaining concrete cloud implementation selection | Defer |
| Pure domain lookup/logging seam | Usually `pkg/config/cloud.go` plus tests; no exported change | Cache-neutral and repository-owned | Does not isolate construction or choose an implementation boundary | Defer |
| Shared opener across root and project constructors | Context and config production/tests, with a likely exported cross-package shape | Cache-compatible but combines the completed root move with the project move | Reopens the sixth slice and enlarges rollback/public-surface risk | Reject |
| Caller refresh-policy cleanup | Command/context files and tests | Requires command-specific policy and error knowledge | Recorded force-sync/warn/continue gaps are guards and behavior changes are out of scope | Reject |
| Direct `ply-config` integration | Unknown without forbidden inspection | Requires external dependency/API knowledge | Not repository-owned here | Not selected |

### Exact Implementation Contract

- Change only `pkg/config/project_init.go` and a new focused
  `pkg/config/project_init_cloud_test.go`. Add one private cloud-config opener
  interface, private dependency value, Git-backed production opener, and a
  small private dependency-taking opener helper. Keep the body and public
  signature of `InitProjectFromDirectory` in place; replace only its direct
  `OpenGitCloudConfig(profilePath)` assignment with the private production
  selection.
- Production must explicitly select the Git-backed opener and delegate exactly
  once to unchanged `OpenGitCloudConfig(profilePath)`. The helper must pass the
  complete already-resolved profile path once and return the exact resulting
  `CloudConfig` interface without wrapping, copying, refreshing, or inspecting
  it. A missing opener returns nil before touching the supplied path.
- Preserve exact project-configuration load/populate ordering; active-profile
  lookup; text-matched missing-file migration and retry; current legacy-home
  choice; the pre-assignment `project.Config.Profile` branch and therefore the
  current decoded-profile non-override behavior; named error propagation;
  cloud assignment before POM probing; POM parse warning/partial behavior; and
  final `ConfigFile`, `Path`, and `Config` assignments. Do not change
  `profiles.go`, `project.go`, `cloud.go`, `Context.LoadProfile`, any caller,
  reader, fixture, or public symbol.
- Use TDD to record production opener selection, exact arbitrary profile-path
  delivery and one call, exact returned interface identity, public active-
  profile and missing-profile migration cache mapping, current decoded-profile
  non-override behavior, safe zero behavior, early project-configuration error
  before opener selection, repeated independent construction, and rejected
  empty recorder populations. Tests must use disposable HOME/project inputs,
  restore global homedir/logger state, and must not add tracked fixtures.
- The config-cloud mutation manifest and meta-test remain byte-unchanged because
  their ten expressions are confined to `pkg/config/cloud.go`. Both gates must
  retain the same IDs, selections, meanings, and exact 10/10 kills.
- Focused verification covers the new project-opener tests and all
  `pkg/config`, `pkg/context`, `pkg/maven`, and `cmd/...` tests, followed by
  both mutation gates. Full verification uses exact Go 1.26.7 under
  `umask 022`, offline/readonly modules, and managed scratch for module
  verification, build, count-one and race count-one tests, vet, API/CLI
  compatibility and meta-tests, preflight, ordinary test/install/launcher
  gates, empty-HOME count-two, all 15 audit meta-controls, and canonical audit.
  Reproduce protected source scope, fixtures, counts/hashes, and accepted
  Q0-Q2 state.
- Rollback is one focused two-file implementation commit: revert
  `pkg/config/project_init.go` and remove the new focused test. There is no
  public API, cache schema, fixture, caller-policy, user-data, or dependency
  migration to unwind.
- The slice expires if constructor signature or ordering, active-profile or
  migration behavior, legacy-home selection, current decoded-profile behavior,
  named errors, POM/project results, concrete Git selection, cache mapping, or
  repeated-call behavior changes; if the root seam must reopen or both
  constructors must move atomically; if caller policy, another reader, an
  exported/shared factory, or `ply-config` knowledge becomes necessary; or if
  any protected fixture, Go, dependency, graph, mutation, or quality input
  changes. Stop for a fresh owning decision rather than compensating in a
  second move.

No production code, tests, fixtures, mutation files, dependency metadata,
`ply-config` work, caller cleanup, construction-behavior repair, Spring work,
packaging work, push, merge, publish, release, stash, revert, or worktree
removal ran in this planning checkpoint.
