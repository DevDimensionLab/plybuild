# Agent Session: Plan Next P8 Cloud Modernization Move

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T190912+0200-plan-next-p8-cloud-modernization-move`
Created: `2026-09-24T19:09:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `66f73bd2834fc9d4178c5d28fc1aeb6db5e68374ee12ed4c43f8a68aa9db333a`
Previous: [2026-09-24T180445+0200-implement-cloud-global-config-loader-seam.md](2026-09-24T180445+0200-implement-cloud-global-config-loader-seam.md)
Next: [2026-09-24T192707+0200-implement-root-cloud-config-opener-seam.md](2026-09-24T192707+0200-implement-root-cloud-config-opener-seam.md)
Outcome: one private zero-value-safe root `Context.LoadProfile` cloud-config opener is fully owned as the sixth P8 implementation slice

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Plan exactly one next bounded P8 cloud-configuration modernization slice after
the completed private global-config loader seam. Re-characterize the resulting
cloud boundary, compare the smallest repository-owned next seams, and select
one implementation slice or stop unresolved. Do not implement the slice in
this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 remains active in cloud configuration only. The first
five P8 implementation slices are complete at focused commits
`57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`,
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`, and
`f12344b115b6c732295d152a7f3d86b876a4d1d8`. This checkpoint may plan only
one next cloud-configuration slice. It may not change source or dependency
metadata, integrate or select `ply-config`, begin Spring or packaging work,
reopen an earlier slice, or combine multiple implementation moves.

# Measurements At Start

Begin only from the clean reciprocal implementation handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused global-config implementation parent and exact
two-file shape, earlier deprecated, project-defaults, services, and cache-probe
ancestry, Google UUID ancestry, sole NEXT state, reciprocal 337-archive chain,
launcher mirror/check, ordinary and ignored cleanliness, and unchanged
protected readers, callers, fixtures, mutation files, and dependency metadata.
The protected project remains Go 1.18 with 234 modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries over
41 loaded modules, 1,067 `go.sum` lines, and protected `go.mod` / `go.sum` /
graph hashes
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

Read the answered global-config implementation archive and its planning
predecessor, the answered deprecated, project-defaults, services, and cache-
probe implementation archives, current P8 roadmap and rolling handover,
`docs/design/quality-lift.md`, `pkg/config/cloud.go`, all focused cloud tests,
the `CloudConfig` callers in `cmd` and `pkg/context`, and
`scripts/mutate-config-cloud` with its meta-test. Treat every recorded caller,
other-reader, real-network, environment, and fixture gap as a guard rather than
permission to combine work. Do not inspect `ply-config`.

# Three Moves

1. Confirm the completed global-config seam preserves production loader
   selection, complete `Directory` delivery, exact direct
   `global-config.yaml` path without `FilePath`, unchanged `file.Open`,
   environment-before-YAML ordering, one independent load per invocation,
   exact read/YAML/partial results, safe zero behavior, `SourceFor`, and all
   public, refresh, caller, reader, fixture, and mutation contracts.
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

The next P8 slice is owned and bounded: isolate only the root
`Context.LoadProfile` construction of `config.CloudConfig` behind one private,
zero-value-safe opener while preserving the existing concrete
`config.OpenGitCloudConfig` selection, exact `<profile>/cloud-config` cache
mapping, local-config lifecycle, project construction, callers, and public
surfaces. This is a repository-owned preparatory construction seam, not a
`ply-config` integration or a change to either production constructor.

### Protected Start And Focused Evidence

- The checkpoint began clean on `codex/upgrade-quality` at handoff HEAD
  `56fb53eaadf1ee8adfdde7b786eb3e82ae86c29b`, direct parent
  `f12344b115b6c732295d152a7f3d86b876a4d1d8`, tree
  `9137d60a9194462ea4054e155dca278a2918b7cd`, with exactly the launcher,
  answered global-config implementation archive, this then-NEXT archive,
  rolling handover, and roadmap changed. The global-config implementation has
  exact parent `82f100de434a37acd94ce16cd8bfba41f7f08e36`, tree
  `e4e28bf42fe484904bb92d8ff4f155878594591c`, and exact two-file set
  `pkg/config/cloud.go` and `pkg/config/cloud_test.go`.
- Cache-probe commit `57f9d5674157d238a2f93462b65161a17e3b5498`, services
  commit `c12307a078a29af74163df0c658d05c821482a33`, project-defaults
  commit `c999212d266d98868930769108e4e8a73070a0a2`, deprecated commit
  `51da9bcfc0eb8ae871c98467c52c942db6b480bd`, and Google UUID commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` are ancestral. The connected
  reciprocal archive graph had 337 records and exactly one NEXT record;
  launcher/archive prompt mirroring, launcher `--check`, shell syntax, and
  ordinary/ignored cleanliness passed. No protected reader, caller, fixture,
  mutation, or dependency file differs between the focused implementation and
  its documentation handoff.
- Exact local Go 1.26.7 reproduced declared Go 1.18, 234 modules, 3,599 graph
  edges, 355 production entries, 429 complete-test entries, 197 module-backed
  entries across 41 loaded modules, and 1,067 `go.sum` lines. `go.mod`,
  `go.sum`, and graph SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
- Under `umask 022`, exact-Go-1.26.7 tests for `pkg/config`, `pkg/context`,
  `pkg/maven`, and `cmd/...` pass with offline, readonly module inputs and
  contained scratch. The unchanged config-cloud direct harness retains all ten
  exact IDs, selections, and meanings with
  `declared=10 killed=10 survived=0 unusable=0`; its T1-T10 meta-test passes
  with the same totals.
- The updated reciprocal graph has 338 records and one NEXT implementation
  archive. Launcher/archive prompt mirroring, launcher `--check`, shell syntax,
  exact five-file documentation/launcher scope, diff checks, and all 62
  launcher lifecycle controls pass. Task-owned scratch stayed beneath the
  managed session root; its 271 directories and 3,391 regular files had zero
  symlink or special entries and were removed before handoff.

### Resulting Cloud Boundary

- The completed global-config seam explicitly selects
  `fileGlobalConfigLoader`, supplies the complete `Directory`, calls `Dir()`
  once, composes exact `file.Path("%s/global-config.yaml", directory.Dir())`,
  never calls `FilePath`, and delegates to unchanged `file.Open`. It expands
  the complete successful bytes through the live environment before unchanged
  YAML decode and preserves raw read errors, complete or partial values, YAML
  errors, safe zero behavior, and one independent load per invocation.
  `SourceFor` slash formatting and its template/tips consumers are unchanged.
- Eight cloud read/effect paths are privately injected: the refresh cache
  probe, Git hooks, examples, templates, services, project defaults,
  deprecated data, and global config. The refresh Git dependency remains
  private as well. No read refreshes implicitly, and the recorded force-sync,
  warning/continue, unconditional example refresh, real-network, and fixture
  gaps remain guards.
- The public `CloudConfig` interface still has 15 methods. `GitCloudConfig`,
  its exported `Impl`, `OpenGitCloudConfig`, all exported cloud data types,
  `GlobalCloudConfig.SourceFor`, and every method signature remain
  compatibility surfaces. The five completed P8 slices changed none of them.
- Exactly two production sites construct cloud configurations. Root profile
  loading in `pkg/context.Context.LoadProfile` assigns
  `config.OpenGitCloudConfig(profilePath)` after assigning the matching local
  config, while `pkg/config.InitProjectFromDirectory` independently selects an
  active or project-specific profile and calls `OpenGitCloudConfig`. Both map
  to `<profile>/cloud-config`. The root path has focused characterization for
  assignment, missing-file creation, existing-byte preservation, exact mode,
  logs, and touch failure; project construction also carries profile discovery,
  migration, project override, configuration, and Maven concerns.

### Candidate Comparison And Selection

| Candidate | Files and public surface | Environment, dependency, and cache knowledge | Reversibility and characterization risk | Decision |
| --- | --- | --- | --- | --- |
| Private root `Context.LoadProfile` cloud opener | `pkg/context/context.go` and `pkg/context/context_test.go`; no exported change | Receives only the complete profile path and delegates production to existing `config.OpenGitCloudConfig`; preserves `<profile>/cloud-config` and needs no new dependency knowledge | One-commit rollback; existing focused lifecycle coverage bounds exact calls, identity, assignment, and errors | **Selected** |
| Private project-construction opener | `pkg/config/project_init.go` plus focused tests; no necessary exported change | Must preserve active-profile discovery, missing-profile migration, project profile override, and HOME-sensitive behavior | Same nominal file count but materially wider environment and error-order surface | Defer |
| One factory routed through both constructors | At least config/context production and test files; likely a new exported cross-package entry point | Cache-compatible, but chooses a shared factory shape before an alternate implementation is known | Larger rollback and public-surface risk; combines two construction moves | Reject |
| Service/template/deprecated domain helper seam | Usually `pkg/config/cloud.go` and tests; no exported change | Repository-owned and cache-neutral | Pure lookup/logging characterization does not isolate the remaining implementation-selection boundary | Defer |
| Caller refresh-policy cleanup | Command/context files and tests | Requires command-specific policy and error knowledge | Recorded gaps are guards and behavior changes are outside this slice | Defer |
| Direct `ply-config` integration | Unknown without forbidden inspection | Requires external dependency/API knowledge | Not repository-owned here | Not selected |

### Exact Implementation Contract

- Change only `pkg/context/context.go` and `pkg/context/context_test.go`. Add one
  private cloud-config opener interface, a private dependency value, the
  Git-backed production opener, and a private `loadProfile` helper. The public
  `Context.LoadProfile(profilePath string)` method remains the production
  wrapper with its exact signature and no return value.
- Production must explicitly select the Git-backed opener and delegate exactly
  once to unchanged `config.OpenGitCloudConfig(profilePath)`, returning that
  complete value as the existing `config.CloudConfig` interface. The helper
  must deliver the complete profile path once and assign the returned interface
  without wrapping, copying, refreshing, or inspecting it. A missing opener
  returns nil without touching the supplied path.
- Preserve exact ordering and behavior in `LoadProfile`: assign
  `config.OpenLocalConfig(profilePath)`, assign the cloud config, check the
  local config's existence, and only when absent log the same debug message and
  call `TouchFile`, logging the same exact error without returning it. Preserve
  missing-parent behavior, existing local bytes, default contents and mode,
  and repeated-call behavior.
- Keep `config.OpenGitCloudConfig`, `GitCloudConfig`, `CloudConfig`,
  `InitProjectFromDirectory`, project-selected profiles, all commands, refresh,
  readers, cache layout, fixtures, mutation files, Go 1.18, and dependency
  metadata byte-exact. Do not route the project constructor or add an exported
  factory in this slice.
- Use TDD to record production opener selection, exact complete profile-path
  delivery and one call, exact returned interface identity, the public
  production result and `<profile>/cloud-config` mapping, safe zero behavior,
  and rejected empty call populations. Retain the existing focused missing,
  existing, touch-error, content, mode, log, and assignment cases; do not add
  caller-policy, project-construction, real-network, or external-library tests.
- The config-cloud mutation manifest and meta-test remain byte-unchanged because
  their ten expressions are confined to `pkg/config/cloud.go`. Both gates must
  retain the same IDs, selections, meanings, and exact 10/10 kills.
- Focused verification covers the new opener tests and all `pkg/context`,
  `pkg/config`, `pkg/maven`, and `cmd/...` tests, followed by both mutation
  gates. Full verification uses exact Go 1.26.7 under `umask 022`, offline/
  readonly modules, and managed scratch for module verification, build,
  count-one and race count-one tests, vet, API/CLI compatibility and meta-tests,
  preflight, ordinary test/install/launcher gates, empty-HOME count-two, all 15
  audit meta-controls, and canonical audit. Reproduce protected source scope,
  fixtures, counts/hashes, and accepted Q0-Q2 state.
- Rollback is one focused two-file implementation commit. There is no public
  API, cache schema, fixture, caller-policy, user-data, or dependency migration
  to unwind.
- The slice expires if `LoadProfile`'s signature, local/cloud assignment order,
  touch/log behavior, concrete Git selection, profile-to-cache mapping, or
  repeated-call semantics change; if both constructors must move atomically;
  if project construction or caller policy must change; if an exported factory
  or `ply-config` knowledge becomes necessary; or if any protected fixture,
  Go, dependency, graph, mutation, or quality input changes. Stop for a fresh
  owning decision rather than compensating in a second move.

No production code, tests, fixtures, mutation files, dependency metadata,
`ply-config` work, caller cleanup, constructor implementation, Spring work,
packaging work, push, merge, publish, release, stash, revert, or worktree
removal ran in this planning checkpoint.
