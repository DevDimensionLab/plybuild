# Agent Session: Implement Cloud Services Loader Seam

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T145009+0200-implement-cloud-services-loader-seam`
Created: `2026-09-24T14:50:09+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `490b4424b2b0d279641784123fd72d41b4512911c0f7468899ff02b2699257e0`
Previous: [2026-09-24T143315+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T143315+0200-plan-next-p8-cloud-modernization-move.md)
Next: [2026-09-24T153032+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T153032+0200-plan-next-p8-cloud-modernization-move.md)
Outcome: the second bounded P8 implementation slice adds one private zero-value-safe services loader while preserving every protected public, refresh, caller, reader, fixture, Go-floor, dependency, graph, mutation, and quality contract

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the second bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the eager cached `services.json`
load behind one private zero-value-safe services loader while preserving the
closure-shaped `GitCloudConfig.Services` API and every refresh, caller, other-
reader, fixture, and dependency contract. Do not integrate or select
`ply-config`, and do not begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first P8
slice is complete at `57f9d5674157d238a2f93462b65161a17e3b5498`. This
checkpoint may change only `pkg/config/cloud.go` and
`pkg/config/cloud_test.go` for the private cached-services loader seam. It may
not change refresh, another cached reader, a constructor or caller, the
config-cloud mutation files, fixtures, dependency metadata, Spring, or
packaging.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, focused cache-probe implementation ancestry and exact
three-file shape, Google UUID ancestry, sole NEXT state, reciprocal 330-
archive chain, launcher mirror/check, ordinary and ignored cleanliness, and
unchanged source, fixtures, and dependency metadata. The protected project
remains Go 1.18 with 234 modules, 3,599 graph edges, 355 production entries,
429 complete-test entries, 197 module-backed entries over 41 loaded modules,
1,067 `go.sum` lines, and protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving private loader seam, not `ply-config`
integration, dependency evaluation, generic document-reader work, caller-
policy cleanup, construction refactoring, security work, Spring work, or
packaging work. Do not change any exported symbol, signature, or field; cache
layout; refresh dependency or branch; `Deprecated`, `ProjectDefaults`,
`GlobalCloudConfig`, Git hooks, templates, or examples; any command/context
caller; fixture; `go.mod`; or `go.sum`. Do not fetch, inspect, or select a new
dependency. Keep every disposable cache, report, copy, and build artifact
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch after
containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
answered cache-probe implementation archive, current P8 roadmap and rolling
handover, `docs/design/quality-lift.md`, `pkg/config/cloud.go`,
`pkg/config/cloud_test.go`, `pkg/config/cloud_refresh_test.go`,
`pkg/config/dir.go`, `pkg/file/file.go`'s `ReadJson` and `Open` paths, the
focused cloud fixture tests, and `scripts/mutate-config-cloud` with its meta-
test. Treat all caller, other-reader, real-network, and fixture gaps as guards.

# Three Moves

1. Add focused failing characterization for production loader selection,
   complete dependency delivery, the exact `services.json` request, eager one-
   call loading, stable repeated closure results and errors, existing fixture
   decode, exact `FilePath` error propagation, safe zero value, and non-empty
   recordings. Do not characterize another reader or caller.
2. Add one private services-loader interface and private dependency value. The
   production implementation accepts the complete `Directory`, calls
   `FilePath("services.json")`, and delegates to unchanged `file.ReadJson`.
   The private helper loads exactly once when `Services()` is called and closes
   over the exact value/error; a missing loader returns
   `filesystem.ErrNoFilesystem` without developer-path access. Preserve the
   exported wrapper, cache-only operation, and every other boundary exactly.
3. Run focused config/context/command and exact unchanged config-cloud mutation
   gates, then the full recorded exact-Go-1.26.7 gates under `umask 022`,
   readonly/offline module inputs, and managed scratch. Verify no API/CLI,
   source-scope, refresh, caller, other-reader, fixture, Go-floor, dependency,
   graph, mutation, or quality regression. Roll back the single two-file commit
   if any contract fails; do not compensate in another slice.

# Automatic Handoff

If and only if the one seam and every gate pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If a contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge, publish,
release, stash, revert, remove the worktree, integrate `ply-config`, or start a
second cloud, Spring, or packaging slice.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The bounded cached-services slice is complete at focused implementation commit
`c12307a078a29af74163df0c658d05c821482a33` (parent
`9bc26b9dc80b8ed059e09597b1a82675abd184af`, tree
`1edb19744f0abd6a4f6256991e14e8aa561711cd`). It changes exactly
`pkg/config/cloud.go` and `pkg/config/cloud_test.go`. No second reader, caller,
refresh path, constructor, fixture, mutation file, dependency, Spring, or
packaging change ran.

### Implementation And Characterization

- `servicesLoader` and `servicesDependencies` are private. Their zero value
  returns exact `filesystem.ErrNoFilesystem` before accessing the supplied
  `Directory`; production explicitly selects `fileServicesLoader`.
- The production loader receives the complete `Directory`, requests exactly
  `services.json`, preserves the exact `FilePath` error, and delegates decode
  to unchanged `file.ReadJson`. `GitCloudConfig.Services` keeps its exported
  closure-shaped signature and delegates to a private helper that loads once,
  eagerly, then replays the exact value and error without reloading.
- Focused tests first failed because the private loader, dependency selector,
  and helper did not exist. They now record production selection, complete
  dependency delivery, the exact filename, one eager load, stable repeated
  success/error closure results, populated fixture decode, exact path-error
  propagation, safe zero behavior without developer-path access, and rejected
  empty recording populations.

### Verification And Protected State

- Focused services tests and all `pkg/config`, `pkg/context`, and `cmd/...`
  tests pass. The byte-unchanged config-cloud direct harness and its T1-T10
  meta-test each report exact
  `declared=10 killed=10 survived=0 unusable=0` with unchanged IDs, selections,
  and meanings.
- Exact Go 1.26.7 module verification, build, count-one tests, race count-one
  tests, and vet pass under `umask 022`, offline module resolution, readonly
  project inputs, and managed scratch. Pinned API/CLI compatibility, both
  compatibility meta-tests, the CLI surface, complete `make preflight`,
  ordinary `make test`, install, and all 62 launcher checks pass. A fresh
  empty-HOME count-two run also passes.
- All 15 quality-audit meta-controls pass. The canonical audit has its
  documented finding exit 1 rather than audit-broken exit 2. Its structured
  27-criterion Q0-Q2 scope contains zero automated FAIL and zero ratchet
  regressions; the single reported overall ratchet is the pre-existing Q3.4
  documentation indicator outside the protected scope. The six manual Q0-Q2
  contracts and accepted 27/27 L2 state remain supported without regression.
- The first attempted full-gate environment selected the Homebrew launcher's
  local Go 1.26.2 after `GOTOOLCHAIN=local`; it was stopped before acceptance.
  The successful gate invoked the existing exact Go 1.26.7 binary directly.
  A first preflight attempt also exposed macOS `make`/`mktemp` attempts outside
  the authorized scratch subtree; the successful complete rerun used the same
  system make binary directly plus a task-local no-template `mktemp` adapter.
  Exact pinned apidiff and golangci-lint binaries were rebuilt only from the
  existing local module download cache into task scratch, with no network or
  project dependency change.
- Go 1.18 and the protected 234 modules, 3,599 graph edges, 355 production
  entries, 429 complete-test entries, 197 module-backed entries across 41
  loaded modules, and 1,067 sum lines reproduce. `go.mod`, `go.sum`, and graph
  hashes remain respectively
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  API/CLI, refresh, callers, other readers, fixtures, mutation files, and
  dependency metadata remain unchanged.
- All disposable evidence stayed beneath managed session scratch. Containment
  and entry types were verified before the complete task-owned subtree was
  removed.

The rollback boundary is the single focused two-file implementation commit.
The reciprocal successor is planning-only and may select at most one further
owned cloud slice; it may not execute that slice, inspect or select
`ply-config`, or begin Spring or packaging work.
