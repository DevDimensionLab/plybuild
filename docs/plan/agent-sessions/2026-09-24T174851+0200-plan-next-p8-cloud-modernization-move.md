# Agent Session: Plan Next P8 Cloud Modernization Move

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T174851+0200-plan-next-p8-cloud-modernization-move`
Created: `2026-09-24T17:48:51+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `32eb159e629eaf60c9f45cbee3aff3f455331f1b9c0b3b80a43e4a46e607fd7a`
Previous: [2026-09-24T170216+0200-implement-cloud-deprecated-loader-seam.md](2026-09-24T170216+0200-implement-cloud-deprecated-loader-seam.md)
Next: [2026-09-24T180445+0200-implement-cloud-global-config-loader-seam.md](2026-09-24T180445+0200-implement-cloud-global-config-loader-seam.md)
Outcome: one private zero-value-safe `GlobalCloudConfig` loader seam is fully owned as the fifth P8 implementation slice

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Plan exactly one next bounded P8 cloud-configuration modernization slice after
the completed private deprecated loader seam. Re-characterize the resulting
cloud boundary, compare the smallest repository-owned next seams, and select
one implementation slice or stop unresolved. Do not implement the slice in
this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 remains active in cloud configuration only. The first
four P8 implementation slices are complete at focused commits
`57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`, and
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`. This checkpoint may plan only
one next cloud-configuration slice. It may not change source or dependency
metadata, integrate or select `ply-config`, begin Spring or packaging work,
reopen an earlier slice, or combine multiple implementation moves.

# Measurements At Start

Begin only from the clean reciprocal implementation handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused deprecated implementation parent and exact
two-file shape, earlier project-defaults, services, and cache-probe ancestry,
Google UUID ancestry, sole NEXT state, reciprocal 335-archive chain, launcher
mirror/check, ordinary and ignored cleanliness, and unchanged protected
readers, callers, fixtures, mutation files, and dependency metadata. The
protected project remains Go 1.18 with 234 modules, 3,599 graph edges, 355
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

Read the answered deprecated implementation archive and its planning
predecessor, the answered project-defaults, services, and cache-probe
implementation archives, current P8 roadmap and rolling handover,
`docs/design/quality-lift.md`, `pkg/config/cloud.go`, all focused cloud tests,
the `CloudConfig` callers in `cmd` and `pkg/context`, and
`scripts/mutate-config-cloud` with its meta-test. Treat every recorded caller,
other-reader, real-network, environment, and fixture gap as a guard rather than
permission to combine work. Do not inspect `ply-config`.

# Three Moves

1. Confirm the completed deprecated seam preserves production loader
   selection, complete `Directory` delivery, exact `deprecated.json` and
   `file.ReadJson` behavior, one independent load per invocation, exact
   path/read/unmarshal and partial-result semantics, recursive and slice data,
   safe zero value, and all public, refresh, caller, reader, fixture, and
   mutation contracts.
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

The next P8 slice is owned and bounded: isolate only the cache-backed
`global-config.yaml` read behind a private, zero-value-safe loader while
preserving the exported `GitCloudConfig.GlobalCloudConfig` method and
`GlobalCloudConfig.SourceFor` behavior. This is a repository-owned preparatory
seam, not a reader migration to `ply-config`, and it adds no public API,
dependency, caller, fixture, cache, refresh, or mutation-manifest change.

### Protected Start And Focused Evidence

- The checkpoint began clean on `codex/upgrade-quality` at handoff HEAD
  `70543a3e7efa58c43c5d21fcc78f0c791f7dd1e7`, direct parent
  `51da9bcfc0eb8ae871c98467c52c942db6b480bd`, tree
  `1a316b9a9b760957dbf16288373eb180f0fb365f`, with exactly the launcher,
  answered deprecated implementation archive, this then-NEXT archive, rolling
  handover, and roadmap changed. The deprecated implementation has exact parent
  `deb6e42d333bb8621c4579ec6b658ab36a9af572`, tree
  `e095c58797d42b9fc4a1ba353c9ddc1916a70614`, and exact two-file set
  `pkg/config/cloud.go` and `pkg/config/cloud_test.go`.
- Cache-probe commit `57f9d5674157d238a2f93462b65161a17e3b5498`, services
  commit `c12307a078a29af74163df0c658d05c821482a33`, project-defaults
  commit `c999212d266d98868930769108e4e8a73070a0a2`, and Google UUID
  implementation commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`
  are ancestral. The connected reciprocal archive graph had 335 records and
  exactly one NEXT record; launcher/archive prompt mirroring, launcher
  `--check`, shell syntax, and ordinary/ignored cleanliness passed. No protected
  source, reader, caller, fixture, mutation, or dependency file differs between
  the focused implementation and its documentation handoff.
- Exact local Go 1.26.7 reproduced declared Go 1.18, 234 modules, 3,599 graph
  edges, 355 production entries, 429 complete-test entries, 197 module-backed
  entries across 41 loaded modules, and 1,067 sum lines. `go.mod`, `go.sum`,
  and graph SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
- Under `umask 022`, exact-Go-1.26.7 focused tests for `pkg/config`,
  `pkg/context`, `pkg/maven`, and `cmd/...` pass with offline, readonly module
  inputs and contained scratch. The unchanged config-cloud direct harness and
  T1-T10 meta-test each retain exact
  `declared=10 killed=10 survived=0 unusable=0` with unchanged IDs,
  selections, and meanings. One discarded focused run used the ambient umask
  and reproduced only the known local-config mode guard; the accepted rerun
  used the mandated `umask 022` without a product change.
- The updated reciprocal graph has 336 records and one NEXT implementation
  archive. Launcher/archive prompt mirroring, launcher `--check`, shell syntax,
  exact five-file documentation/launcher scope, diff checks, and launcher
  lifecycle controls pass. Task-owned scratch stayed beneath the managed
  session root; its 267 directories and 3,389 regular files had zero symlink
  or special entries and were removed before handoff.

### Resulting Cloud Boundary

- `CloudConfig` still exposes 15 methods. `GitCloudConfig`, exported `Impl`,
  `OpenGitCloudConfig`, `GlobalCloudConfig.SourceFor`, all public data types,
  and construction in `Context.LoadProfile` and `InitProjectFromDirectory`
  remain compatibility surfaces. No read method implicitly refreshes.
- Refresh remains separately private-effect-injected and preserves the exact
  local-config-first `<target>/.git` branch. `Services`, `ProjectDefaults`, and
  `Deprecated` retain their separate private loaders; Git hooks, examples, and
  templates retain their complete private filesystem dependencies. No selected
  work crosses any of those seven paths.
- The completed deprecated seam explicitly selects `fileDeprecatedLoader` for
  production, delivers the complete implementation `Directory`, requests only
  `deprecated.json`, and delegates to unchanged `file.ReadJson`. Its helper
  loads once on every invocation, so separate calls remain independent and
  non-memoized. Exact path, read, unmarshal, recursive, nil/empty-slice,
  partial-result, and safe-zero semantics are focused-characterized. The public
  wrapper, `ListDeprecated`, Maven/status/upgrade callers, refresh, fixtures,
  and mutation manifest remain unchanged.
- `GlobalCloudConfig` is the sole direct document reader. It composes
  `<implementation>/global-config.yaml` from `Directory.Dir()` using
  `file.Path`, intentionally does not call `Directory.FilePath`, reads through
  `file.Open`, expands the live process environment across the complete bytes,
  and only then unmarshals YAML. Each public invocation reads independently.
- `SourceFor` concatenates `RootUrl`, `RelativFileUrl`, the requested directory,
  and name with its existing slash pattern. Its only production consumers are
  template markdown output in `cmd/build_options.go` and tips output in
  `cmd/tips.go`. There is no tracked global-config fixture or focused reader
  test; real cloud/network behavior and caller refresh policy remain unproved
  guards rather than scope for this slice.

### Candidate Comparison And Selection

| Candidate | Files and public surface | Environment, dependency, and cache knowledge | Reversibility and characterization risk | Decision |
| --- | --- | --- | --- | --- |
| Private `GlobalCloudConfig` loader | `pkg/config/cloud.go` and `pkg/config/cloud_test.go`; no exported change | Uses only current `Directory.Dir`, `file.Path`, `file.Open`, `os.ExpandEnv`, and `yaml.Unmarshal`; remains cache-only and non-memoized | One-commit rollback; every distinct ordering/error/result contract can be covered with disposable local files and controlled environment values | **Selected** |
| Environment-expansion-only seam | Nominally the same production/test files and no exported change | Knows only `os.ExpandEnv`, leaving the sole direct document read in place | Small but incomplete; splits one reader into multiple future implementation moves without owning the boundary | Reject |
| `SourceFor` or command-consumer characterization | Tests and potentially command files; touches a public result consumer | No new dependency knowledge, but does not isolate the remaining reader | Consumer gaps are real, yet changing caller policy is forbidden and unnecessary for the private loader | Defer |
| Shared document loader | One production file plus tests spanning completed JSON readers and YAML | Would combine `FilePath` JSON semantics with direct-path YAML and process environment behavior | Larger rollback and guesses at a generic abstraction before `ply-config` is known | Reject |
| Constructor or refresh-policy seam | At least config/context/command files and tests | Cache-compatible but unrelated to the sole remaining direct read | Larger surface and explicitly separate policy/construction ownership | Defer |
| Direct `ply-config` integration | Unknown without forbidden inspection | Requires external dependency/API knowledge | Not repository-owned here | Not selected |

### Exact Implementation Contract

- Change only `pkg/config/cloud.go` and `pkg/config/cloud_test.go`. Add one
  private global-config loader interface, private dependency value, file-backed
  production implementation, and private helper. The exported
  `GlobalCloudConfig()` method remains the production wrapper with its exact
  signature; `GlobalCloudConfig.SourceFor` remains byte-exact.
- The production loader accepts the complete `Directory`, calls `Dir()` and
  composes exactly `file.Path("%s/global-config.yaml", directory.Dir())`. It
  must not call `FilePath` or add an existence probe. It delegates the read to
  unchanged `file.Open`, preserving raw read errors and zero results, applies
  `[]byte(os.ExpandEnv(string(b)))` to the complete successful read exactly
  before unchanged `yaml.Unmarshal`, and returns the exact decoded value,
  partial value, and error.
- The private helper invokes the selected loader exactly once per method call.
  Separate calls must re-read and re-expand independently; no closure,
  memoization, refresh, or cross-call state may be introduced. A missing loader
  returns exact `filesystem.ErrNoFilesystem` with a zero `GlobalCloudConfig`
  before any supplied developer directory is accessed.
- Preserve receiver-derived cache root, fixed filename, direct `Dir()` path
  semantics, `file.Open` result/error identity, environment lookup and
  expansion ordering, YAML tags and partial-result behavior, zero/missing
  values, `SourceFor` formatting, template/tips consumers, every exported
  symbol/signature/field and constructor, refresh, all other readers, Go 1.18,
  fixtures, and dependency metadata.
- Use TDD to record production loader selection, complete `Directory`
  delivery, exact direct filename composition with no `FilePath` call, one
  load per invocation and independent repeated invocations, exact success and
  dependency-error returns, representative complete YAML decoding with
  environment substitutions, exact raw read error, representative partial
  YAML result/error, safe zero behavior without developer-path access, exact
  unchanged `SourceFor` formatting, and non-empty recordings. Use only
  disposable local fixtures; do not add or change a tracked fixture, caller
  policy, another reader, network, or real-cloud test.
- The config-cloud mutation manifest and meta-test remain byte-unchanged; no
  current expression crosses `GlobalCloudConfig`. Both gates must retain exact
  `declared=10 killed=10 survived=0 unusable=0` with the same IDs, selected
  tests, and meanings.
- Focused verification covers the new global-config loader tests and all
  `pkg/config`, `pkg/context`, `pkg/maven`, and `cmd/...` tests, followed by
  both mutation gates. Full verification uses exact Go 1.26.7 under
  `umask 022`, offline/readonly modules, and managed scratch for module
  verification, build, count-one and race count-one tests, vet, API/CLI
  compatibility and meta-tests, preflight, ordinary test/install/launcher
  gates, empty-HOME count-two, all 15 audit meta-controls, and canonical audit.
  Reproduce protected source scope, fixtures, counts/hashes, and accepted
  Q0-Q2 state.
- Rollback is one focused two-file implementation commit. There is no public
  API, cache schema, fixture, user-data, caller-policy, or dependency migration
  to unwind.
- The slice expires if the reader stops using the direct fixed YAML path, starts
  using `FilePath`, changes read/environment/YAML/result semantics, stops
  loading independently per invocation, `SourceFor` or caller policy must
  change, another reader is needed, public construction or API changes,
  mutation ownership changes, `ply-config` knowledge becomes necessary, or any
  protected fixture, Go, dependency, graph, or quality input changes. Stop for
  a fresh owning decision rather than compensating in a second move.

No production code, tests, fixtures, mutation files, dependency metadata,
`ply-config` work, caller cleanup, reader implementation, Spring work,
packaging work, push, merge, publish, release, stash, revert, or worktree
removal ran in this planning checkpoint.
