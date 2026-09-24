# Agent Session: Plan First P8 Cloud Modernization Move

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T130019+0200-plan-first-p8-cloud-modernization-move`
Created: `2026-09-24T13:00:19+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5b79ae79ae06e04d7f89a8ff62d71c777dc96e45376ccbd490f0da321c223a32`
Previous: [2026-09-24T122118+0200-complete-p7-quality-exit-gate.md](2026-09-24T122118+0200-complete-p7-quality-exit-gate.md)
Next: [2026-09-24T131521+0200-implement-cloud-cache-probe-seam.md](2026-09-24T131521+0200-implement-cloud-cache-probe-seam.md)
Outcome: the current boundary and evidence are mapped, and one private cache-probe seam is fully owned as the first P8 implementation slice

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Begin P8 with one bounded cloud-configuration modernization planning
checkpoint. Characterize the existing `GitCloudConfig` boundary, its
cache-first refresh behavior, and its compatibility fixtures; then select and
specify exactly one smallest owned implementation slice toward `ply-config`.
Do not implement the slice in this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 is active. Its ordered work is cloud configuration,
then Spring modernization under characterization tests, then any inactive
packaging only through a separate scope decision. This checkpoint may plan
only the first cloud-configuration slice. It may not start Spring or packaging
work, reopen P7, evaluate a dependency, run an ownership study, or change
source or dependency metadata.

# Measurements At Start

The P7 exit preserved Go 1.18, all direct roots, product behavior, accepted
27/27 Q0-Q2 PASS at L2, every target-specific P7 guard, and the exact
234-module / 3,599-edge / 355-production / 429-complete-test / 197-module-
backed / 41-loaded-module / 1,067-sum-line project state. Protected `go.mod`,
`go.sum`, and graph hashes remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Exact Go 1.26.7 unchanged-project verification, build, count-one tests, race
count-one tests, vet, compatibility contracts, and quality meta-controls pass.

Begin only from the clean reciprocal P7-exit handoff on
`codex/upgrade-quality`. Verify its branch, five-file documentation/launcher
commit shape, Google UUID ancestry, sole NEXT state, archive chain, launcher
mirror/check, ordinary and ignored cleanliness, and unchanged source and
dependency metadata. Stop for a fresh owning decision if a protected input
changed.

# Role And Boundaries

This is a product-boundary planning checkpoint, not implementation, dependency
evaluation, security investigation, or upstream ownership research. Use the
existing source, tests, design record, and local Git history. Do not add or
modify production code, tests, fixtures, `go.mod`, or `go.sum`; do not fetch or
select a `ply-config` dependency; and do not claim compatibility beyond the
recorded evidence.

Keep every disposable cache, report, project copy, and verification artifact
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Use only ordinary bounded
inspection and established tests. Verify containment and entry types, then
remove all task-owned scratch before handoff.

# Required Reading

Read the P8 roadmap and P7 exit answer, the rolling handover,
`docs/design/quality-lift.md`, `pkg/config/cloud.go`,
`pkg/config/cloud_refresh_test.go`, `pkg/config/cloud_test.go`, the focused
cloud fixture/contract tests, the `CloudConfig` callers in `cmd` and
`pkg/context`, and the existing config-cloud mutation harness. Treat those as
the planning inputs; do not repeat P7 dependency work.

# Three Moves

1. Map the current cloud-configuration construction, refresh, cache, read, and
   caller boundaries. Record the exact behaviors and public/API surfaces that
   a first slice must preserve, especially cache-first operation and current
   compatibility fixtures.
2. Run only focused ordinary tests needed to confirm that characterization.
   Identify any evidence gap, but do not fill it with production/test changes
   in this planning checkpoint. Compare candidate seams by size, ownership,
   reversibility, and compatibility risk.
3. Select exactly one smallest implementation slice, or stop unresolved.
   Specify its files, behavior contract, focused tests, full gates, rollback
   boundary, and expiry conditions. Do not implement or combine a second
   slice.

# Automatic Handoff

If and only if one implementation slice is fully owned and bounded, answer
this archive, update the roadmap and rolling handover, prepare exactly one
reciprocal NEXT archive that implements only that slice, replace only launcher
mutable regions, run launcher/handoff checks, and make one local handoff
commit. If ownership or behavior is unresolved, prepare one decision successor
instead. Do not execute the successor, push, merge, publish, release, stash,
revert, remove the worktree, or broaden P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The first P8 implementation slice is owned and bounded: inject only the
private cache-presence probe used by `GitCloudConfig.Refresh`, preserving its
exact `<target>/.git` decision and delegating the production probe to the
existing `file.Exists` behavior. This is a preparatory seam toward a later
`ply-config` adapter, not that adapter, and it requires no dependency fetch,
selection, or metadata change.

### Protected Start And Focused Evidence

- The checkpoint began clean on `codex/upgrade-quality` at HEAD
  `7179b28fa1b179fd66e0c6252f97aa7bbf81be41`, parent
  `59d026a0db3d5952c81b8168ed2c9fabf47b0104`, and tree
  `385df5e78d0a9750315c1f63e668b1f6394fb444`. Its changed set is exactly the
  launcher, answered P7-exit archive, this then-NEXT archive, rolling
  handover, and roadmap.
- Google UUID implementation commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` is ancestral. No product source
  changed after it. `go.mod` and `go.sum` retain protected SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
  and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
- The reciprocal 327-archive graph, sole NEXT state, launcher prompt mirror,
  launcher `--check`, ordinary cleanliness, and ignored cleanliness all pass.
  Exact Go 1.26.7 is active.
- Under isolated task scratch, offline module resolution, and no source
  mutation, focused cloud/config tests, focused context construction tests,
  and focused command-consumer tests pass. The first context run inherited a
  restrictive umask and observed mode `0600`; the recorded project environment
  `umask 022` rerun passed its established `0644` contract. That first run is
  an environment mismatch, not product evidence.
- The updated launcher/archive mirror and all 62 launcher continuity controls
  pass. The task-owned scratch subtree contained 3,480 ordinary file/directory
  entries, zero symlinks, and zero special entries; it was verified beneath the
  managed session root and removed before handoff.

### Current Boundary

- `OpenGitCloudConfig(profilePath)` returns the public concrete
  `GitCloudConfig` and maps it to `<profilePath>/cloud-config` through its
  exported `Impl DirConfig` field. `Context.LoadProfile` constructs the root
  `CloudConfig`; `InitProjectFromDirectory` independently constructs each
  project value from the active or project-selected profile.
- The public `CloudConfig` interface has 15 methods: `Implementation`,
  `Refresh`, `Services`, `LinkFromService`, `DefaultServiceEnvironmentUrl`,
  `Deprecated`, `ProjectDefaults`, `GitHookFiles`, `ListDeprecated`,
  `HasTemplate`, `ValidTemplatesFrom`, `Templates`, `Template`, `Examples`,
  and `GlobalCloudConfig`. `GlobalCloudConfig.SourceFor`, the related exported
  data types, `GitCloudConfig`, and `OpenGitCloudConfig` are also compatibility
  surfaces. The selected slice changes none of them.
- `Refresh` first parses `LocalConfigFile.Config`, derives the implementation
  directory, and probes exactly `<target>/.git`. Presence selects one
  `Pull(target)`; absence selects one `Clone(configuredURL, target)`. The URL
  precedes the target, Git failures retain `shell.Output.FormatError`, and
  success returns nil. The existing private, zero-value-safe Git dependency
  seam delegates production calls to `shell.GitPull`/`shell.GitClone`.
- Cache-first operation is split across this branch and its callers. Existing
  cache state selects pull instead of clone, while every read method consumes
  files from the local implementation directory and never refreshes itself.
  Status, upgrade, Git-hook, diagrams, and tips flows refresh only through
  `SyncActiveProfileCloudConfig` when `ForceCloudSync` is true and otherwise
  read cached content; most of those callers warn and continue on refresh
  failure. Build refreshes when its inherited `--cloud-sync` value is true
  (currently default true). Build-example and option-example flows refresh
  unconditionally before reading. Template listing and add-template read the
  cache without an internal refresh. These caller policies are preserved, not
  normalized, by the selected slice.
- Cached reads cover `services.json`, `deprecated.json`,
  `project-defaults.json`, `global-config.yaml`, `git-hooks`, `templates`,
  `examples`, `tips`, and `resources`. Template and downstream Maven/context
  consumers retain the same `CloudConfig` interface boundary.

### Compatibility Evidence And Gaps

- Refresh contracts cover exact clone URL/target order, existing `.git`
  pull selection, formatted clone/pull errors, safe zero-value Git behavior,
  and non-empty call populations. The ten-mutation config-cloud harness binds
  four refresh decisions plus Git-hook, template, example, service, and
  template-deduplication decisions with zero survivors in the recorded state.
- Focused fixtures cover service/default-environment resolution; current
  `ply.json` and legacy nested `co-pilot.json` templates; exact walk order,
  names, loaded projects, partial results, and errors; and virtual filesystem
  contracts for examples and Git hooks. Host acceptance constructs cached
  `deprecated.json` and `project-defaults.json`. P7's unchanged API/CLI
  compatibility evidence remains applicable because this checkpoint changed
  no product or test file.
- There is no direct caller test for `ForceCloudSync` gating, refresh-error
  wrapping/continue policy, or the unconditional example refreshes. There is
  no focused fixture contract for `GlobalCloudConfig`, `Deprecated`, or
  `ProjectDefaults`, and real cloud/network behavior remains unverified. The
  existing refresh test distinguishes present from missing `.git` but does
  not independently record the exact probe or non-ENOENT stat semantics.
  These gaps are recorded; the chosen seam does not change the uncovered
  caller or reader behavior.

### Candidate Comparison And Selected Slice

| Candidate | Size and ownership | Reversibility and risk | Decision |
| --- | --- | --- | --- |
| Inject the private `.git` cache probe | One production file, one focused test file, and one mutation-manifest expression; entirely repository-owned | Private, mechanically delegates to existing `file.Exists`, one-commit rollback, no caller/API/fixture/dependency change | Selected |
| Add a new construction factory and route both constructors | Crosses `pkg/config`, `pkg/context`, and project construction; an exported factory would add API | More callers and compatibility surface than needed before an adapter exists | Defer |
| Migrate a cached reader or caller refresh policy | Crosses fixture formats or command behavior with recorded evidence gaps | Higher compatibility risk and not required for the first seam | Defer |
| Integrate `ply-config` directly | Its dependency API and selection were intentionally outside this checkpoint | Not owned without forbidden fetch/evaluation and metadata changes | Not selected |

### Exact Implementation Contract

- Modify only `pkg/config/cloud.go`, `pkg/config/cloud_refresh_test.go`, and
  the affected branch mutation expression in `scripts/mutate-config-cloud`.
  `scripts/test-mutate-config-cloud` is a required gate but should not need a
  semantic change.
- Add a private, zero-value-safe cache-presence dependency to the existing
  refresh dependency bundle. Production must select an implementation that
  calls the existing `file.Exists` once with exactly
  `file.Path("%s/.git", target)`. A nil cache dependency reports absent and
  cannot touch developer paths.
- Preserve local-config parsing order, target derivation, presence-to-pull and
  absence-to-clone selection, single Git attempt, argument order, logs,
  `FormatError`, return values, public symbols/signatures/fields, all reads,
  every caller, fixture bytes, Go 1.18, and `go.mod`/`go.sum`.
- Focused verification must cover production dependency selection, exact
  probe bytes and one probe, present/missing branch selection, local-config
  error before probe/Git, clone/pull error formatting, safe zero values, all
  current cloud fixtures, context construction, and the exact ten-mutation
  harness. Then run the roadmap's exact-Go-1.26.7 module verification, build,
  count-one, race count-one, vet, API/CLI compatibility and meta-tests,
  preflight/test/install/launcher gates, hermetic count-two test, quality
  meta-controls, and canonical audit. Use `umask 022`, offline/readonly module
  inputs, and managed scratch.
- Rollback is the single implementation commit across those three owned files;
  there is no cache schema, user-data, fixture, dependency, or public-API
  migration to unwind.
- The slice expires if cache location or presence semantics, local-config
  ordering, Git output/error behavior, refresh callers, public API, fixture
  formats, mutation ownership, Go/dependency metadata, or the intended
  `ply-config` integration contract changes. Any request to select or fetch
  that dependency, normalize caller refresh policy, migrate a read method, or
  touch a second P8 slice requires a fresh owning checkpoint.

No production code, tests, fixtures, dependency metadata, Spring work,
packaging work, dependency evaluation, or ownership study ran here.
