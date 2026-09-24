# Agent Session: Plan Next P8 Cloud Modernization Move

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T164219+0200-plan-next-p8-cloud-modernization-move`
Created: `2026-09-24T16:42:19+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `82b6c90a9c7909be4d16cde068acde461c6017f3d3ed814165cd557d0b0f3f2c`
Previous: [2026-09-24T154421+0200-implement-cloud-project-defaults-loader-seam.md](2026-09-24T154421+0200-implement-cloud-project-defaults-loader-seam.md)
Next: [2026-09-24T170216+0200-implement-cloud-deprecated-loader-seam.md](2026-09-24T170216+0200-implement-cloud-deprecated-loader-seam.md)
Outcome: one private zero-value-safe `deprecated.json` loader seam is fully owned as the fourth P8 implementation slice

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Plan exactly one next bounded P8 cloud-configuration modernization slice after
the completed private project-defaults loader seam. Re-characterize the
resulting `GitCloudConfig` boundary, compare the smallest repository-owned next
seams, and select one implementation slice or stop unresolved. Do not implement
the slice in this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 remains active in cloud configuration only. The first
three P8 implementation slices are complete at focused commits
`57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`, and
`c999212d266d98868930769108e4e8a73070a0a2`. This checkpoint may plan only
one next cloud-configuration slice. It may not change source or dependency
metadata, integrate or select `ply-config`, begin Spring or packaging work,
reopen an earlier slice, or combine multiple implementation moves.

# Measurements At Start

Begin only from the clean reciprocal implementation handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused project-defaults implementation parent and exact
two-file shape, earlier cache-probe and services implementation ancestry,
Google UUID ancestry, sole NEXT state, reciprocal 333-archive chain, launcher
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

Read the answered project-defaults implementation archive and its planning
predecessor, the answered services and cache-probe implementation archives,
current P8 roadmap and rolling handover, `docs/design/quality-lift.md`,
`pkg/config/cloud.go`, all focused cloud tests, the `CloudConfig` callers in
`cmd` and `pkg/context`, and `scripts/mutate-config-cloud` with its meta-test.
Treat every recorded caller, other-reader, real-network, and fixture gap as a
guard rather than permission to combine work. Do not inspect `ply-config`.

# Three Moves

1. Confirm the completed project-defaults seam preserves production loader
   selection, complete `Directory` delivery, exact `project-defaults.json` and
   `file.ReadJson` behavior, one independent load per invocation, exact
   path/read/unmarshal and partial-result semantics, safe zero value, and all
   public, refresh, caller, reader, fixture, and mutation contracts.
2. Map only the remaining cloud boundary and run focused read-only tests needed
   to validate candidate ownership. Compare candidates by file count, public
   surface, dependency knowledge, cache compatibility, reversibility, and the
   recorded characterization gaps. Do not fetch or inspect `ply-config`.
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
`deprecated.json` read behind a private, zero-value-safe loader while
preserving the exported `GitCloudConfig.Deprecated` contract and its existing
one-read-per-invocation behavior. This is a repository-owned preparatory seam,
not a reader migration to `ply-config`, and it adds no public API, dependency,
caller, fixture, cache, refresh, or mutation-manifest change.

### Protected Start And Focused Evidence

- The checkpoint began clean on `codex/upgrade-quality` at handoff HEAD
  `e3dbf55ce219a30887e4a1140f0edef8199036b7`, direct parent
  `c999212d266d98868930769108e4e8a73070a0a2`, tree
  `da7c07deba78d3bb1884bf04577e2b16ebaa63e9`, with exactly the launcher,
  answered project-defaults implementation archive, this then-NEXT archive,
  rolling handover, and roadmap changed. The project-defaults implementation
  has exact parent `f45d6c4ab63b1c6e76fa31a1b710fbc07ae31f38`, tree
  `7914ce896206532040fc66fc0dd7f366926aa891`, and exact two-file set
  `pkg/config/cloud.go` and `pkg/config/cloud_test.go`.
- Cache-probe commit `57f9d5674157d238a2f93462b65161a17e3b5498`, services
  commit `c12307a078a29af74163df0c658d05c821482a33`, and Google UUID
  implementation commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`
  are ancestral. The connected reciprocal archive graph has 333 records and
  exactly one NEXT record; launcher/archive prompt mirror and launcher
  `--check` pass. Ordinary and ignored state were empty. No protected source,
  reader, caller, fixture, mutation, or dependency file differs between the
  focused implementation and its documentation handoff.
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
  selections, and meanings.
- The updated reciprocal graph has 334 records and one NEXT implementation
  archive. Launcher/archive prompt mirror, launcher `--check`, shell syntax,
  exact five-file documentation/launcher scope, diff checks, and all 62
  launcher lifecycle controls pass. The contained task root had 267
  directories, 3,396 regular files, zero symlinks, and zero special entries;
  it was verified beneath managed session scratch and removed before handoff.

### Resulting Cloud Boundary

- `CloudConfig` still exposes 15 methods. `GitCloudConfig`, exported `Impl`,
  `OpenGitCloudConfig`, `GlobalCloudConfig.SourceFor`, all data types, and the
  two concrete construction sites in `Context.LoadProfile` and
  `InitProjectFromDirectory` remain compatibility surfaces.
- Refresh remains separately private-effect-injected and preserves the exact
  local-config-first `<target>/.git` branch. `Services` and `ProjectDefaults`
  now have separate private loaders. The former retains exact eager one-load
  closure replay; the latter retains one independent load per invocation.
  Both receive the complete `Directory`, request their exact JSON names,
  delegate to unchanged `file.ReadJson`, and fail safely when zero-valued.
- The completed project-defaults seam explicitly selects
  `fileProjectDefaultsLoader` for production, delivers the complete
  implementation `Directory`, requests only `project-defaults.json`, and
  delegates to unchanged `file.ReadJson`. Its helper loads exactly once on
  every invocation, so separate calls remain independent and non-memoized.
  Exact path, read, unmarshal, partial-result, and safe-zero semantics are
  focused-characterized. The public wrapper, context merge/error/log/order
  contract, refresh, other readers, fixtures, and mutation manifest remain
  unchanged.
- `GitHookFiles`, `Examples`, and `Templates` retain their private complete
  filesystem dependencies and focused root/order/error/partial-result/default
  contracts. No selected work crosses any of those five injected read paths.
- Two direct document readers remain. `Deprecated` resolves exact
  `deprecated.json` through `Directory.FilePath` and calls `file.ReadJson` for
  every invocation. `ListDeprecated` logs its dependency, associated-
  dependency, and replacement-template data; Maven `RemoveDeprecated` and
  `StatusDeprecated` feed upgrade/status flows. The acceptance profile supplies
  an empty valid JSON document, while focused Maven tests use a cloud double for
  ordered partial and replacement-error behavior. There is no direct reader
  fixture contract yet.
- `GlobalCloudConfig` directly composes `global-config.yaml`, calls
  `file.Open`, expands the process environment, unmarshals YAML, and feeds
  template/tips `SourceFor` output. It has no focused reader fixture contract.
  No read method implicitly refreshes. Force-sync gating, warning/continuation,
  unconditional example refreshes, and real cloud/network gaps remain guards.

### Candidate Comparison And Selection

| Candidate | Files and public surface | Dependency and cache knowledge | Reversibility and characterization risk | Decision |
| --- | --- | --- | --- | --- |
| Private `Deprecated` loader | `pkg/config/cloud.go` and `pkg/config/cloud_test.go`; no exported change | Uses only current `Directory`, `CloudDeprecated`, and `file.ReadJson`; remains cache-only and non-memoized | One-commit rollback; exact reader semantics can be focused-characterized without changing wider callers | **Selected** |
| Private `GlobalCloudConfig` loader | Nominally the same two files and no exported change | Must preserve direct path composition, `file.Open`, environment expansion, YAML decode, and source formatting | More distinct semantics, two command consumers, process-environment input, and no focused fixture | Defer |
| Shared document loader | One production file plus tests for multiple result types | Would guess a generic JSON/YAML abstraction before `ply-config` is known | Combines filenames, JSON/YAML, environment, result, and caller contracts | Reject |
| Constructor or caller-policy seam | At least config/context/command files and tests | Cache-compatible but unrelated to the remaining direct read | Larger rollback and either public construction or command policy change | Defer |
| Direct `ply-config` integration | Unknown without forbidden inspection | Requires external dependency/API knowledge | Not repository-owned here | Not selected |

### Exact Implementation Contract

- Change only `pkg/config/cloud.go` and `pkg/config/cloud_test.go`. Add one
  private deprecated loader interface, private dependency value, file-backed
  production implementation, and private helper. The exported `Deprecated()`
  method remains the production wrapper with its exact signature.
- The production loader accepts the complete `Directory`, requests exactly
  `deprecated.json`, and delegates to unchanged `file.ReadJson`. The private
  helper invokes the selected loader exactly once per method call. Separate
  calls must re-read independently; no closure, memoization, refresh, or cross-
  call state may be introduced. A missing loader returns exact
  `filesystem.ErrNoFilesystem` with a zero `CloudDeprecated` before any
  supplied developer path is accessed.
- Preserve receiver-derived cache root, fixed filename, `DirConfig.FilePath`
  presence and exact error behavior, `file.ReadJson` read/unmarshal formatting
  and any partial result, all recursive `CloudDeprecated` fields and JSON tags,
  nil/empty slices, `ListDeprecated`, Maven and command caller behavior,
  acceptance fixture bytes, every exported symbol/signature/field and
  constructor, refresh, all other readers, Go 1.18, and dependency metadata.
- Use TDD to record production loader selection, complete `Directory`
  delivery, exact filename, one load per invocation and independent repeated
  invocations, exact success/value and dependency-error returns,
  representative complete recursive JSON decoding plus exact read/unmarshal
  error formatting, exact `FilePath` error propagation, safe zero behavior
  without developer-path access, and non-empty recordings. Do not add or
  change a caller-policy, another reader, tracked fixture, network, or real-
  cloud test.
- The config-cloud mutation manifest and meta-test remain byte-unchanged; no
  current expression crosses `Deprecated`. Both gates must retain exact
  `declared=10 killed=10 survived=0 unusable=0` with the same IDs, selected
  tests, and meanings.
- Focused verification covers the new deprecated-loader tests and all
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
- The slice expires if `Deprecated` stops using the fixed cached JSON, stops
  reading once per invocation, its path/read/unmarshal/result semantics change,
  caller policy must change, another reader is needed, public construction or
  API changes, mutation ownership changes, `ply-config` knowledge becomes
  necessary, or any protected fixture, Go, dependency, graph, or quality input
  changes. Stop for a fresh owning decision rather than compensating in a
  second move.

No production code, tests, fixtures, mutation files, dependency metadata,
`ply-config` work, caller cleanup, reader migration, Spring work, packaging
work, push, merge, publish, release, stash, revert, or worktree removal ran in
this planning checkpoint.
