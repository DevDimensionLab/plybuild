# Agent Session: Plan Next P8 Cloud Modernization Move

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T213100+0200-plan-next-p8-cloud-modernization-move`
Created: `2026-09-24T21:31:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c521790fcb99719997e49489f70434d5ecb0817b93c6dd403433246af49d463c`
Previous: [2026-09-24T204758+0200-implement-project-cloud-config-opener-seam.md](2026-09-24T204758+0200-implement-project-cloud-config-opener-seam.md)
Next: [2026-09-24T215338+0200-implement-cloud-template-project-loader-seam.md](2026-09-24T215338+0200-implement-cloud-template-project-loader-seam.md)
Outcome: one private zero-value-safe template-project loader is fully owned as the eighth P8 implementation slice; current template walking, project/profile/POM construction, cache mapping, partial results, callers, and public behavior remain frozen

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Plan exactly one next bounded P8 cloud-configuration modernization slice after
the completed private project cloud-config opener seam. Re-characterize the
resulting cloud boundary, compare the smallest repository-owned next seams, and
select one implementation slice or stop unresolved. Do not implement the slice
in this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 remains active in cloud configuration only. The first
seven P8 implementation slices are complete at focused commits
`57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`,
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`,
`f12344b115b6c732295d152a7f3d86b876a4d1d8`,
`76c3571f04969c304785338bc31c34a266e2f869`, and
`9c1899d484633e0ce6cf112544ba33ca89bc4adb`. This checkpoint may plan only
one next cloud-configuration slice. It may not change source or dependency
metadata, integrate or select `ply-config`, begin Spring or packaging work,
reopen an earlier slice, or combine multiple implementation moves.

# Measurements At Start

Begin only from the clean reciprocal implementation handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused project-opener implementation parent and exact
two-file shape, earlier root-opener, global-config, deprecated, project-
defaults, services, and cache-probe ancestry, Google UUID ancestry, sole NEXT
state, reciprocal 341-archive chain, launcher mirror/check, ordinary and
ignored cleanliness, and unchanged protected readers, callers, fixtures,
mutation files, project construction, profile behavior, and dependency
metadata. The protected project remains Go 1.18 with 234 modules, 3,599 graph
edges, 355 production entries, 429 complete-test entries, 197 module-backed
entries over 41 loaded modules, 1,067 `go.sum` lines, and protected `go.mod` /
`go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is a planning checkpoint, not implementation, dependency evaluation,
profile-policy repair, caller-policy cleanup, generic reader or construction
refactoring, security work, Spring work, or packaging work. Use only repository
source, tests, fixtures, design records, and local Git history. Do not modify
production code, tests, fixtures, mutation files, `go.mod`, or `go.sum`; do not
fetch, inspect, or select a new dependency; and do not claim real cloud
compatibility beyond recorded evidence. Keep every disposable artifact beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`, verify containment and entry types, and
remove task-owned scratch before handoff.

# Required Reading

Read the answered project-opener implementation archive and its planning
predecessor, the answered root-opener, global-config, deprecated, project-
defaults, services, and cache-probe implementation archives, current P8 roadmap
and rolling handover, `docs/design/quality-lift.md`,
`pkg/config/project_init.go` and its focused tests, `pkg/config/profiles.go`,
`pkg/config/project.go`, `pkg/context/context.go` and its focused opener tests,
`pkg/config/cloud.go` and all focused cloud tests, the `CloudConfig` callers in
`cmd`, `pkg/context`, and `pkg/template`, and `scripts/mutate-config-cloud` with
its meta-test. Treat every recorded decoded-profile, non-ENOENT, migration,
caller-policy, other-reader, real-network, environment, and fixture gap as a
guard rather than permission to combine work. Do not inspect `ply-config`.

# Three Moves

1. Confirm the completed project opener preserves explicit Git production
   selection, exact complete resolved-profile-path delivery and one call,
   returned interface identity, active-profile and missing-profile migration
   cache mapping, current decoded-profile non-override behavior, zero-value
   safety, early configuration failure, independent repeated construction, and
   all public, root-opener, project/POM, caller, reader, fixture, mutation,
   Go-floor, and dependency contracts.
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
worktree, integrate `ply-config`, repair decoded-profile behavior, reopen an
earlier slice, or begin Spring or packaging work.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The next P8 slice is owned and bounded: isolate only the direct
`InitProjectFromDirectory` call inside `GitCloudConfig.templates` behind one
private, zero-value-safe template-project loader. The existing filesystem walk
seam, project constructor, profile and migration behavior, cache mapping,
template policy, and command callers remain unchanged.

### Protected Start And Focused Evidence

- The checkpoint began clean on `codex/upgrade-quality` at handoff HEAD
  `47e81b2c9ef0d6184d8bf6476365ccdf70d48365`, direct parent
  `9c1899d484633e0ce6cf112544ba33ca89bc4adb`, tree
  `7517f25e00367d9ecfd82719728bbff867a3b047`, with exactly the launcher,
  answered project-opener implementation archive, this then-NEXT archive,
  rolling handover, and roadmap changed. The project-opener implementation has
  exact parent `c5fec9cce330094b0628b7ebc7a60e1c6f91d474`, tree
  `1083cb67f52b347a5f0760a4ddda7cf533c2a0d4`, and exact two-file set
  `pkg/config/project_init.go` and new
  `pkg/config/project_init_cloud_test.go`.
- Cache-probe commit `57f9d5674157d238a2f93462b65161a17e3b5498`, services
  commit `c12307a078a29af74163df0c658d05c821482a33`, project-defaults
  commit `c999212d266d98868930769108e4e8a73070a0a2`, deprecated commit
  `51da9bcfc0eb8ae871c98467c52c942db6b480bd`, global-config commit
  `f12344b115b6c732295d152a7f3d86b876a4d1d8`, root-opener commit
  `76c3571f04969c304785338bc31c34a266e2f869`, and Google UUID commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` are ancestral. The connected
  reciprocal graph had 341 archives and exactly one NEXT archive;
  launcher/archive prompt mirroring, launcher `--check`, shell syntax, and
  ordinary/ignored cleanliness passed.
- Every protected production, test, fixture, mutation, profile, caller, and
  dependency path is byte-identical between the focused project implementation
  and its documentation handoff. Exact Go 1.26.7 reproduced declared Go 1.18,
  234 modules, 3,599 graph edges, 355 production entries, 429 complete-test
  entries, 197 module-backed entries across 41 loaded modules, and 1,067
  `go.sum` lines. `go.mod`, `go.sum`, and graph SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
- Under `umask 022`, exact-Go-1.26.7 focused project-opener and template tests
  and all `pkg/config`, `pkg/context`, `pkg/maven`, and `cmd/...` tests pass
  with offline, readonly module inputs and contained scratch. The byte-unchanged
  config-cloud direct harness retains all ten exact IDs, selections, and
  meanings with `declared=10 killed=10 survived=0 unusable=0`; its T1-T10
  meta-test passes with the same totals.

### Resulting Cloud Boundary

- The completed project opener explicitly selects `gitProjectCloudConfigOpener`,
  passes the complete resolved profile path exactly once to unchanged
  `OpenGitCloudConfig`, and returns the same `CloudConfig` interface identity.
  Active-profile and missing-profile migration cache mapping, current decoded-
  profile non-override behavior, zero-value safety, early configuration
  failure, independent repeated construction, and project/POM ordering remain
  exact. The root opener retains its corresponding lifecycle contracts.
- The public 15-method `CloudConfig` interface, all exported cloud types,
  `GitCloudConfig.Impl`, `OpenGitCloudConfig`, readers, refresh, fixtures,
  callers, and cache layout are unchanged. Cache probing, Git refresh, four
  cached document loaders, Git-hook listing, example listing, and template
  filesystem walking are already private dependencies.
- One repository-owned construction dependency remains direct inside the cloud
  reader boundary: the template walk callback calls
  `InitProjectFromDirectory(relPath[0])` for each exact `ply.json` or legacy
  `co-pilot.json` match. That call imports active-profile/HOME selection,
  text-matched migration, decoded-profile behavior, cloud construction, and
  optional POM loading into the callback even though the walk itself is
  injected. It is the only such direct project-construction call in
  `pkg/config/cloud.go`.
- Pure domain and presentation methods remain repository-owned but effect-free:
  service link lookup, default-environment selection, deprecated logging,
  template lookup, `HasTemplate`, deduplication, and source URL formatting.
  Force-sync/warn/continue behavior, unconditional example refreshes, real
  network/cloud behavior, and missing caller or fixture evidence remain guards.

### Candidate Comparison And Selection

| Candidate | Files and public surface | Environment, dependency, and cache knowledge | Reversibility and characterization risk | Decision |
| --- | --- | --- | --- | --- |
| Private template-project loader | `pkg/config/cloud.go` and `pkg/config/cloud_templates_test.go`; no exported change | Receives only the already-discovered complete template directory and delegates to unchanged `InitProjectFromDirectory`; production retains all current profile/HOME/migration/POM/cache behavior behind the seam | One-commit rollback; focused recorders can freeze path, call count, complete project/error delivery, order, partial results, and zero behavior | **Selected** |
| Private template lookup/provider seam | `pkg/config/cloud.go` plus cloud tests; no exported change | Cache-neutral, but wraps already-injected `Templates()` calls in `HasTemplate`, `Template`, or `ValidTemplatesFrom` | Does not remove the direct construction dependency and risks changing repeated-walk policy | Defer |
| Pure service/deprecated/source policy seam | Usually `pkg/config/cloud.go` plus tests; no exported change | No environment or dependency knowledge | Repository-owned but does not isolate an effect or construction boundary | Defer |
| Shared root/project/template constructor | Context and config production/tests with a likely exported or cross-package factory | Would combine profile lifecycle, project construction, and template walking | Reopens completed slices and enlarges rollback and behavior risk | Reject |
| Caller refresh-policy cleanup | Multiple command/context files and tests | Requires command-specific force-sync, warning, continuation, and network knowledge | Recorded policy gaps are guards and behavior changes are out of scope | Reject |
| Direct `ply-config` integration | Unknown without forbidden inspection | Requires external dependency and API knowledge | Not repository-owned in this checkpoint | Not selected |

### Exact Implementation Contract

- Change only `pkg/config/cloud.go` and
  `pkg/config/cloud_templates_test.go`. Add one private template-project loader
  interface, extend the existing private `templatesDependencies` with that
  loader, and add one private production implementation that delegates to
  unchanged `InitProjectFromDirectory`. Replace only the direct constructor
  call in the existing walk callback with the private dependency call.
- `systemTemplatesDependencies` must continue to select the complete system
  filesystem dependency and must explicitly select the production project
  loader. For every exact current or legacy config match, production passes the
  complete already-derived `relPath[0]` exactly once to
  `InitProjectFromDirectory` and returns the complete `Project` value and exact
  error without wrapping, copying fields, refreshing, or inspection. A missing
  project loader returns zero `Project` plus exact
  `filesystem.ErrNoFilesystem` without invoking project construction.
- Preserve the exact receiver-derived `<implementation>/templates` walk root;
  callback order and incoming-error suppression; exact current/legacy filename
  matching; path splitting and relative-name derivation; one load per match;
  ordered append behavior; current and legacy project results; exact project-
  load logging/error propagation and ordered partial result; final walk-error
  and nil-result behavior; and independent repeated public invocations. Do not
  change `InitProjectFromDirectory`, either cloud opener, profile/migration/POM
  behavior, any public symbol, caller, other reader, fixture, or cache schema.
- Use TDD to record production loader and filesystem selection, exact arbitrary
  complete directory delivery and one call per matching entry, no project load
  for ignored or incoming-error callbacks, exact complete project and embedded
  interface identities, exact error/log/partial-result behavior, independent
  repeated invocations, safe missing-loader behavior with an injected walk,
  and rejected empty recorder populations. Retain the existing tracked current
  and legacy fixture characterization without adding a fixture.
- The config-cloud mutation files remain byte-exact. The two existing template
  mutations still target exact filename matching and relative-name selection;
  direct and meta gates must retain all ten IDs, selections, meanings, and
  exact 10/10 kills.
- Focused verification covers the template-project-loader tests and all
  `pkg/config`, `pkg/context`, `pkg/maven`, and `cmd/...` tests, followed by
  both mutation gates. Full verification uses exact Go 1.26.7 under
  `umask 022`, offline/readonly modules, and managed scratch for module
  verification, build, count-one and race count-one tests, vet, API/CLI
  compatibility and meta-tests, preflight, ordinary test/install/launcher
  gates, empty-HOME count-two, all 15 audit meta-controls, and canonical audit.
  Reproduce protected source scope, fixtures, counts/hashes, and accepted
  Q0-Q2 state.
- Rollback is one focused two-file implementation commit: revert the private
  loader/dependency addition and its focused tests. There is no public API,
  cache schema, profile, fixture, caller-policy, user-data, or dependency
  migration to unwind.
- The slice expires if template walk root, callback ordering, incoming-error
  suppression, matching, path/name derivation, project values or identities,
  loader call counts, project/profile/migration/POM/cache behavior, partial or
  final errors, logs, repeated-call behavior, or zero behavior changes; if the
  filesystem walk and project load cannot remain one bounded seam; if a caller,
  constructor, profile helper, fixture, exported/shared factory, or
  `ply-config` knowledge becomes necessary; or if any protected Go,
  dependency, graph, mutation, or quality input changes. Stop for a fresh
  owning decision rather than compensating in a second move.

No production code, tests, fixtures, mutation files, dependency metadata,
`ply-config` work, profile repair, caller cleanup, Spring work, packaging work,
push, merge, publish, release, stash, revert, or worktree removal ran in this
planning checkpoint.
