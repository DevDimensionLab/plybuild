# Agent Session: Plan Next P8 Cloud Modernization Move

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T143315+0200-plan-next-p8-cloud-modernization-move`
Created: `2026-09-24T14:33:15+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8390e2e850c7be6bd49cf8f2cfe9f0f2c065ed089fdae62d13df3b50b67f74bb`
Previous: [2026-09-24T131521+0200-implement-cloud-cache-probe-seam.md](2026-09-24T131521+0200-implement-cloud-cache-probe-seam.md)
Next: [2026-09-24T145009+0200-implement-cloud-services-loader-seam.md](2026-09-24T145009+0200-implement-cloud-services-loader-seam.md)
Outcome: one private zero-value-safe cached-services loader seam is fully owned as the second P8 implementation slice

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Plan exactly one next bounded P8 cloud-configuration modernization slice after
the completed private cache-probe seam. Re-characterize the resulting
`GitCloudConfig` boundary, compare the smallest repository-owned next seams,
and select one implementation slice or stop unresolved. Do not implement the
slice in this checkpoint.

# Authorized Roadmap

P2A-P7 are complete. P8 remains active in cloud configuration only. The first
P8 implementation slice is complete at focused commit
`57f9d5674157d238a2f93462b65161a17e3b5498`. This checkpoint may plan only
one next cloud-configuration slice. It may not change source or dependency
metadata, integrate or select `ply-config`, begin Spring or packaging work,
reopen P7, or combine multiple implementation moves.

# Measurements At Start

Begin only from the clean reciprocal implementation handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused implementation parent and three-file shape,
Google UUID ancestry, sole NEXT state, reciprocal archive chain, launcher
mirror/check, ordinary and ignored cleanliness, and unchanged fixtures and
dependency metadata. The protected project remains Go 1.18 with 234 modules,
3,599 graph edges, 355 production entries, 429 complete-test entries, 197
module-backed entries over 41 loaded modules, 1,067 `go.sum` lines, and
protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is a planning checkpoint, not implementation, dependency evaluation,
caller-policy cleanup, reader migration, security work, Spring work, or
packaging work. Use only repository source, tests, fixtures, design records,
and local Git history. Do not modify production code, tests, fixtures,
`go.mod`, or `go.sum`; do not fetch, inspect, or select a new dependency; and
do not claim real cloud compatibility beyond the recorded evidence. Keep all
disposable state beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`, verify containment
and entry types, and remove task-owned scratch before handoff.

# Required Reading

Read the answered cache-probe implementation archive and its planning
predecessor, the current P8 roadmap and rolling handover,
`docs/design/quality-lift.md`, `pkg/config/cloud.go`,
`pkg/config/cloud_refresh_test.go`, the focused cloud fixture tests, the
`CloudConfig` callers in `cmd` and `pkg/context`, and
`scripts/mutate-config-cloud` with its meta-test. Treat every recorded caller,
reader, real-network, and fixture gap as a guard rather than permission to fill
multiple gaps.

# Three Moves

1. Confirm the completed cache-probe seam preserves one exact
   `<target>/.git` probe, local-config-first ordering, present-to-pull and
   absent-to-clone selection, Git formatting, zero-value safety, and all
   public/caller/read/fixture contracts.
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

The next P8 slice is owned and bounded: isolate only the eager cached
`services.json` load behind a private, zero-value-safe services loader while
preserving the exported closure-shaped `GitCloudConfig.Services` contract.
This is a repository-owned preparatory seam, not a cached-reader migration to
`ply-config`, and it adds no public API, dependency, caller, fixture, or cache
change.

### Protected Start And Focused Evidence

- The checkpoint began clean on `codex/upgrade-quality` at handoff HEAD
  `b62fcc579f7efa5f0d050bcbf8c1b1b715747e68`, direct parent
  `57f9d5674157d238a2f93462b65161a17e3b5498`, with exactly the launcher,
  answered cache-probe implementation archive, this then-NEXT archive,
  rolling handover, and roadmap changed. The implementation has exact parent
  `d385616b3191be9a8a2a0b42ced411c3aacd9f11`, tree
  `bbd931ef16a334079473e1cd35aa387ed6b61ef9`, and the exact three-file set
  `pkg/config/cloud.go`, `pkg/config/cloud_refresh_test.go`, and
  `scripts/mutate-config-cloud`.
- Google UUID implementation commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` is ancestral. The connected
  reciprocal archive graph has 329 records and exactly one NEXT record; the
  launcher/archive prompt mirror and launcher `--check` pass. Ordinary and
  ignored state were empty, and fixtures plus dependency metadata are
  unchanged from the focused implementation parent.
- Exact local Go 1.26.7 reproduced Go 1.18, 234 modules, 3,599 graph edges,
  355 production entries, 429 complete-test entries, 197 module-backed
  entries over 41 loaded modules, and 1,067 sum lines. `go.mod`, `go.sum`, and
  graph SHA-256 values remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
- Focused `pkg/config`, `pkg/context`, and `cmd` tests pass under `umask 022`,
  offline/readonly module inputs, and managed scratch. A discarded first
  measurement replaced `HOME` before pinning the already-installed toolchain,
  so Go's selector could not see local Go 1.26.7 and attempted a forbidden
  offline toolchain resolution; its zero pipeline counts are not evidence.
  The successful rerun invoked the existing exact binary directly and kept
  build, temporary, HOME, and config state beneath task scratch.

### Resulting Cloud Boundary

- `CloudConfig` still exposes 15 methods, and `GitCloudConfig`, its exported
  `Impl DirConfig`, `OpenGitCloudConfig`, `GlobalCloudConfig.SourceFor`, and
  all related data types remain compatibility surfaces. `Context.LoadProfile`
  and `InitProjectFromDirectory` remain the two concrete construction sites.
- Refresh is now fully private-effect-injected: it parses local config first,
  composes and probes exactly `<target>/.git` once, selects one pull when
  present or one clone with URL before target when absent, preserves
  `shell.Output.FormatError`, and is safe with a zero dependency bundle.
  None of those contracts is part of the selected next slice.
- `GitHookFiles`, `Examples`, and `Templates` already have private complete
  filesystem dependencies with focused exact-root, order, error, partial-
  result, production-selection, and zero-value contracts. They remain
  untouched.
- Four document readers remain direct. `Services` eagerly resolves
  `services.json` through `DirConfig.FilePath`, calls `file.ReadJson`, and
  returns a closure over the resulting `CloudServices` value and error.
  `Deprecated` and `ProjectDefaults` similarly read fixed JSON files, while
  `GlobalCloudConfig` reads `global-config.yaml`, expands environment values,
  and unmarshals YAML. No read method performs an implicit refresh.
- Caller policy is unchanged: status, upgrade, Git hooks, diagrams, and tips
  conditionally refresh through `ForceCloudSync` and generally warn/continue;
  build conditionally refreshes with its current default true; option-example
  and build-example flows refresh unconditionally; template listing and
  add-template consume the cache. Project defaults, deprecated data,
  services, templates, examples, tips, resources, Git hooks, and global
  source configuration retain their current cache locations and consumers.

### Remaining Evidence Gaps

- Refresh callers still lack direct focused contracts for force-sync gating,
  wrapped/continued errors, and unconditional example refreshes. Real cloud
  and network behavior remain unverified.
- `Deprecated`, `ProjectDefaults`, and `GlobalCloudConfig` have no direct
  repository fixture contract for their complete reader behavior. Their
  caller evidence does not authorize combining them with another reader.
- `Services` has the smallest positive starting contract: the tracked fixture
  proves successful decode and the pure service/default-environment methods
  have focused tests. Its exact eager one-load closure behavior, dependency
  error identity, fixed filename request, production loader selection, and
  safe zero dependency are not yet directly recorded. The selected slice may
  close only those `Services` gaps.

### Candidate Comparison And Selection

| Candidate | Files and public surface | Dependency/cache knowledge | Reversibility and gap risk | Decision |
| --- | --- | --- | --- | --- |
| Private cached-services loader | `pkg/config/cloud.go` and `pkg/config/cloud_test.go`; no exported change | Uses only the current `Directory`, `CloudServices`, and `file.ReadJson` contracts; remains cache-only and needs no `ply-config` knowledge | One-commit rollback; existing success fixture bounds the added eager/error/zero-value characterization | **Selected** |
| Private `Deprecated` or `ProjectDefaults` loader | Nominally the same two files and no exported change | Cache-only, but each direct reader lacks a focused fixture contract and has downstream caller behavior to separate first | Would combine first complete reader characterization with a new seam under weaker evidence | Defer |
| Private `GlobalCloudConfig` loader | Two files and no exported change | Must preserve environment expansion, YAML decoding, and source formatting | Highest single-reader semantic risk and no tracked focused fixture | Defer |
| Common loader across all document readers | One production file plus several test/caller surfaces | Guesses a shared adapter shape before `ply-config` is known | Combines four migrations and their distinct gaps | Reject for this checkpoint |
| Route both constructors through a new factory | At least `pkg/config`, `pkg/context`, project construction, and tests | Cache-compatible and dependency-neutral | Either adds exported API or leaves an incomplete package-local split; larger rollback/public surface | Defer |
| Direct `ply-config` integration | Unknown until separately evaluated | Requires expressly forbidden dependency/API knowledge | Not repository-owned in this checkpoint | Not selected |

### Exact Implementation Contract

- Change only `pkg/config/cloud.go` and `pkg/config/cloud_test.go`. Add a
  private services-loader interface, a private dependency value, a production
  file-backed implementation, and a private `services` helper. The exported
  `Services()` method remains the production wrapper and retains its exact
  signature.
- The production loader must accept the complete `Directory`, call
  `FilePath("services.json")`, then delegate to the unchanged
  `file.ReadJson` behavior. The private helper must invoke the selected loader
  exactly once, eagerly when `Services()` is called, and return a closure that
  replays the exact loaded `CloudServices` value and exact error without
  reloading or refreshing. A missing loader must return
  `filesystem.ErrNoFilesystem` without dereferencing or accessing a supplied
  developer path.
- Preserve fixed filename and receiver-derived cache root; present/missing and
  non-ENOENT `DirConfig.FilePath` semantics; `file.ReadJson` read/unmarshal
  formatting; successful fixture bytes and decoded fields; nil/empty values;
  closure timing; every exported symbol, signature, field, and concrete
  constructor; refresh behavior; all other readers; every caller and fixture;
  Go 1.18; and `go.mod`/`go.sum`.
- Use TDD to record production dependency selection, complete dependency
  delivery, exact `services.json` request, eager one-call behavior, stable
  repeated closure results for success and error, the existing tracked
  success fixture, exact `FilePath` error propagation, safe zero behavior,
  and rejection of empty recording populations. Do not add a caller-policy,
  other-reader, network, or real-cloud test in this slice.
- The config-cloud mutation manifest and meta-test remain byte-unchanged: no
  current mutation expression crosses `Services`. Both direct and meta
  harnesses must still report exact
  `declared=10 killed=10 survived=0 unusable=0` with the same IDs, selected
  tests, and meanings.
- Focused verification covers the new `Services` tests, all `pkg/config`,
  `pkg/context`, and `cmd` tests, and both config-cloud mutation gates. Full
  verification then uses exact Go 1.26.7 under `umask 022`, offline/readonly
  modules, and managed scratch for module verification, build, count-one and
  race count-one tests, vet, API/CLI compatibility plus meta-tests, preflight,
  ordinary test/install/launcher gates, empty-HOME count-two, all 15 audit
  meta-controls, and the canonical audit. Reproduce protected source scope,
  fixture and dependency bytes, counts/hashes, and accepted Q0-Q2 state.
- Rollback is one focused two-file implementation commit. There is no public
  API, cache schema, fixture, user-data, caller-policy, or dependency migration
  to unwind.
- The slice expires if `Services` stops being eager or closure-shaped, its
  filename/path/error/JSON/cache semantics change, another reader or caller is
  needed, public construction or API changes, mutation ownership changes,
  `ply-config` knowledge becomes necessary, or any protected fixture, Go,
  dependency, graph, or quality input changes. Stop for a fresh owning
  decision rather than compensating in a second move.

No production code, tests, fixtures, dependency metadata, `ply-config` work,
caller cleanup, reader migration, Spring work, packaging work, push, merge,
publish, release, stash, revert, or worktree removal ran in this planning
checkpoint.
