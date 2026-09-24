# Agent Session: Plan Next P8 Cloud Modernization Move

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T153032+0200-plan-next-p8-cloud-modernization-move`
Created: `2026-09-24T15:30:32+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `b6c40a0cf745776fcb4aebbc1ca8558976550796aced836f68cfa5aff780d34c`
Previous: [2026-09-24T145009+0200-implement-cloud-services-loader-seam.md](2026-09-24T145009+0200-implement-cloud-services-loader-seam.md)
Next: [2026-09-24T154421+0200-implement-cloud-project-defaults-loader-seam.md](2026-09-24T154421+0200-implement-cloud-project-defaults-loader-seam.md)
Outcome: one private zero-value-safe `project-defaults.json` loader seam is fully owned as the third P8 implementation slice

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Plan exactly one next bounded P8 cloud-configuration modernization slice after
the completed private cached-services loader seam. Re-characterize the
resulting `GitCloudConfig` boundary, compare the smallest repository-owned next
seams, and select one implementation slice or stop unresolved. Do not implement
the slice in this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 remains active in cloud configuration only. The first
two P8 implementation slices are complete at focused commits
`57f9d5674157d238a2f93462b65161a17e3b5498` and
`c12307a078a29af74163df0c658d05c821482a33`. This checkpoint may plan only
one next cloud-configuration slice. It may not change source or dependency
metadata, integrate or select `ply-config`, begin Spring or packaging work,
reopen an earlier slice, or combine multiple implementation moves.

# Measurements At Start

Begin only from the clean reciprocal implementation handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused cached-services implementation parent and exact
two-file shape, cache-probe and Google UUID ancestry, sole NEXT state,
reciprocal 331-archive chain, launcher mirror/check, ordinary and ignored
cleanliness, and unchanged protected readers, callers, fixtures, mutation
files, and dependency metadata. The protected project remains Go 1.18 with 234
modules, 3,599 graph edges, 355 production entries, 429 complete-test entries,
197 module-backed entries over 41 loaded modules, 1,067 `go.sum` lines, and
protected `go.mod` / `go.sum` / graph hashes
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

Read the answered cached-services implementation archive and its planning
predecessor, the answered cache-probe implementation archive, current P8
roadmap and rolling handover, `docs/design/quality-lift.md`,
`pkg/config/cloud.go`, all focused cloud tests, the `CloudConfig` callers in
`cmd` and `pkg/context`, and `scripts/mutate-config-cloud` with its meta-test.
Treat every recorded caller, other-reader, real-network, and fixture gap as a
guard rather than permission to combine work. Do not inspect `ply-config`.

# Three Moves

1. Confirm the completed services seam preserves production loader selection,
   complete `Directory` delivery, exact `services.json` and `file.ReadJson`
   behavior, one eager load, stable closure replay, exact path errors, safe zero
   value, and all public, refresh, caller, reader, fixture, and mutation
   contracts.
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
`project-defaults.json` read behind a private, zero-value-safe loader while
preserving the exported `GitCloudConfig.ProjectDefaults` contract and its
existing one-read-per-invocation behavior. This is a repository-owned
preparatory seam, not a reader migration to `ply-config`, and it adds no public
API, dependency, caller, fixture, cache, or refresh change.

### Protected Start And Focused Evidence

- The checkpoint began clean on `codex/upgrade-quality` at handoff HEAD
  `d6382a80e14e9bcee48a058a74940e310e15fc48`, direct parent
  `c12307a078a29af74163df0c658d05c821482a33`, tree
  `a015048e12f63ddc4d472b0285de30be201de207`, with exactly the launcher,
  answered cached-services implementation archive, this then-NEXT archive,
  rolling handover, and roadmap changed. The services implementation has exact
  parent `9bc26b9dc80b8ed059e09597b1a82675abd184af`, tree
  `1edb19744f0abd6a4f6256991e14e8aa561711cd`, and exact two-file set
  `pkg/config/cloud.go` and `pkg/config/cloud_test.go`.
- Cache-probe commit `57f9d5674157d238a2f93462b65161a17e3b5498` and Google
  UUID implementation commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`
  are ancestral. The connected reciprocal archive graph has 331 records and
  exactly one NEXT record; launcher/archive prompt mirror and launcher
  `--check` pass. Ordinary and ignored state were empty. Protected readers,
  callers, fixtures, mutation files, and dependency metadata have no change
  after the focused services implementation.
- Exact local Go 1.26.7 reproduced declared Go 1.18, 234 modules, 3,599 graph
  edges, 355 production entries, 429 complete-test entries, 197 module-backed
  entries across 41 loaded modules, and 1,067 sum lines. `go.mod`, `go.sum`,
  and graph SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
- Under `umask 022`, exact-Go-1.26.7 focused tests for `pkg/config`,
  `pkg/context`, `pkg/maven`, and `cmd/...` pass with offline, readonly module
  inputs and contained scratch. A discarded first package run inherited a
  restrictive umask and failed only the existing `0644` local-config metadata
  assertion; the recorded rerun used the project umask and passed. The
  unchanged config-cloud direct harness and T1-T10 meta-test each retain exact
  `declared=10 killed=10 survived=0 unusable=0` with unchanged IDs,
  selections, and meanings.

### Resulting Cloud Boundary

- `CloudConfig` still exposes 15 methods. `GitCloudConfig`, exported `Impl`,
  `OpenGitCloudConfig`, `GlobalCloudConfig.SourceFor`, all data types, and the
  two concrete construction sites in `Context.LoadProfile` and
  `InitProjectFromDirectory` remain compatibility surfaces.
- Refresh remains separately private-effect-injected and preserves the exact
  local-config-first `<target>/.git` branch. `Services` now has a private
  loader with exact production selection, complete `Directory` delivery,
  `FilePath("services.json")`, unchanged `file.ReadJson`, one eager load,
  stable closure replay, exact path error, and safe zero-value contracts.
- `GitHookFiles`, `Examples`, and `Templates` retain their private complete
  filesystem dependencies and focused root/order/error/partial-result/default
  contracts. No selected work crosses any of those four injected paths.
- Three direct document readers remain. `ProjectDefaults` resolves
  `project-defaults.json` through `Directory.FilePath` and calls
  `file.ReadJson` for every invocation. `Context.OnEachMavenProject` is its
  only concrete caller: it reads once per delivered non-nil cloud config,
  warns and continues on error, merges even the returned zero/partial value,
  then runs the existing project flow. Its focused context test records call
  counts, merge-before-job behavior, error continuation, and logs; the local
  acceptance profile supplies a cache-only JSON fixture.
- `Deprecated` directly reads `deprecated.json`; its consumers span
  `ListDeprecated`, Maven removal/status behavior, and status/upgrade command
  flows. `GlobalCloudConfig` directly reads `global-config.yaml`, expands the
  environment, unmarshals YAML, and feeds template/tips source formatting.
  Neither reader has a direct focused repository fixture contract. No read
  method implicitly refreshes.
- Force-sync gating, warning/continuation, and unconditional example-refresh
  gaps remain caller-policy guards. Real cloud/network behavior remains
  unverified. None is part of the selected slice.

### Candidate Comparison And Selection

| Candidate | Files and public surface | Dependency and cache knowledge | Reversibility and characterization risk | Decision |
| --- | --- | --- | --- | --- |
| Private `ProjectDefaults` loader | `pkg/config/cloud.go` and `pkg/config/cloud_test.go`; no exported change | Uses only current `Directory`, `CloudProjectDefaults`, and `file.ReadJson`; remains cache-only and non-memoized | One-commit rollback; one caller has focused call/order/error/merge coverage and acceptance supplies valid JSON | **Selected** |
| Private `Deprecated` loader | Nominally the same two files and no exported change | Same JSON/cache primitives | Wider `ListDeprecated`/Maven/status/upgrade surface and no direct reader fixture contract | Defer |
| Private `GlobalCloudConfig` loader | Two files and no exported change | Must retain direct path composition, environment expansion, YAML decode, and source formatting | Highest single-reader semantic risk; two command consumers and no focused fixture | Defer |
| Shared document loader | One production file plus tests for multiple result types | Would guess a generic adapter before `ply-config` is known | Combines distinct JSON/YAML, environment, error, and caller contracts | Reject |
| Constructor or caller-policy seam | At least config/context/command files and tests | Cache-compatible but unrelated to the direct read boundary | Larger rollback and either public construction or policy change | Defer |
| Direct `ply-config` integration | Unknown without forbidden inspection | Requires external dependency/API knowledge | Not repository-owned here | Not selected |

### Exact Implementation Contract

- Change only `pkg/config/cloud.go` and `pkg/config/cloud_test.go`. Add one
  private project-defaults loader interface, private dependency value,
  file-backed production implementation, and private helper. The exported
  `ProjectDefaults()` method remains the production wrapper with its exact
  signature.
- The production loader accepts the complete `Directory`, requests exactly
  `project-defaults.json`, and delegates to unchanged `file.ReadJson`. The
  private helper invokes the selected loader exactly once per method call.
  Separate calls must re-read independently; no closure, memoization, refresh,
  or cross-call state may be introduced. A missing loader returns exact
  `filesystem.ErrNoFilesystem` with a zero `CloudProjectDefaults` before any
  supplied developer path is accessed.
- Preserve receiver-derived cache root, fixed filename, `DirConfig.FilePath`
  presence and exact error behavior, `file.ReadJson` read/unmarshal formatting
  and any partial result, all `CloudProjectDefaults` fields and JSON tags,
  nil/empty slices, caller merge/error/log/order behavior, acceptance fixture
  bytes, every exported symbol/signature/field and constructor, refresh,
  services, every other reader/caller, Go 1.18, and dependency metadata.
- Use TDD to record production loader selection, complete `Directory`
  delivery, exact filename, one load per invocation and independent repeated
  invocations, exact success/value and dependency-error returns, representative
  complete JSON decoding plus exact read/unmarshal error formatting, exact
  `FilePath` error propagation, safe zero behavior without developer-path
  access, and non-empty recordings. Do not add or change a caller-policy,
  another reader, tracked fixture, network, or real-cloud test.
- The config-cloud mutation manifest and meta-test remain byte-unchanged; no
  current expression crosses `ProjectDefaults`. Both gates must retain exact
  `declared=10 killed=10 survived=0 unusable=0` with the same IDs, selected
  tests, and meanings.
- Focused verification covers the new project-defaults tests and all
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
- The slice expires if `ProjectDefaults` stops using the fixed cached JSON,
  stops reading once per invocation, its path/read/unmarshal/result semantics
  change, caller or merge policy must change, another reader is needed, public
  construction or API changes, mutation ownership changes, `ply-config`
  knowledge becomes necessary, or any protected fixture, Go, dependency,
  graph, or quality input changes. Stop for a fresh owning decision rather
  than compensating in a second move.

No production code, tests, fixtures, mutation files, dependency metadata,
`ply-config` work, caller cleanup, reader migration, Spring work, packaging
work, push, merge, publish, release, stash, revert, or worktree removal ran in
this planning checkpoint.
