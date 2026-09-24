# Agent Session: Implement Cloud Cache Probe Seam

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T131521+0200-implement-cloud-cache-probe-seam`
Created: `2026-09-24T13:15:21+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `307f948bec4ba4332baa6c515b4c491f5805aba70a7244b5676b70df4b3ecec9`
Previous: [2026-09-24T130019+0200-plan-first-p8-cloud-modernization-move.md](2026-09-24T130019+0200-plan-first-p8-cloud-modernization-move.md)
Next: [2026-09-24T143315+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T143315+0200-plan-next-p8-cloud-modernization-move.md)
Outcome: the first P8 implementation slice injects one private zero-value-safe cache probe while preserving the exact refresh branch and every protected public, caller, reader, fixture, Go-floor, dependency, graph, and quality contract

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the first bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: inject the private cache-presence probe used
by `GitCloudConfig.Refresh` while preserving the exact cache-first branch and
every public, caller, read, fixture, and dependency contract. Do not integrate
or select `ply-config`, and do not begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. This checkpoint
may implement only the private `<target>/.git` cache-probe seam in
`pkg/config/cloud.go`, its focused characterization in
`pkg/config/cloud_refresh_test.go`, and the corresponding exact mutation
expression in `scripts/mutate-config-cloud`. Spring and inactive packaging
remain out of scope.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, Google UUID ancestry, sole NEXT state, reciprocal
archive chain, launcher mirror/check, ordinary and ignored cleanliness, and
unchanged product source and dependency metadata. The protected project
remains Go 1.18 with 234 modules, 3,599 graph edges, 355 production entries,
429 complete-test entries, 197 module-backed entries over 41 loaded modules,
1,067 `go.sum` lines, and protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving internal seam, not direct `ply-config`
integration, dependency evaluation, caller-policy cleanup, reader migration,
security work, Spring work, or packaging work. Do not change any exported
symbol, signature, or field; construction or command callers; cached read
method; fixture; `go.mod`; or `go.sum`. Do not fetch, inspect, or select a new
dependency. Keep every disposable cache, report, copy, and build artifact
beneath `${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch after
containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract, the
current P8 roadmap and rolling handover, `docs/design/quality-lift.md`,
`pkg/config/cloud.go`, `pkg/config/cloud_refresh_test.go`, the focused cloud
fixture tests, `internal/adapter/filesystem/filesystem.go`, and
`scripts/mutate-config-cloud` with its meta-test. Treat the recorded caller and
reader gaps as guards, not invitations to broaden this slice.

# Three Moves

1. Add focused failing characterization for one exact cache probe, present and
   missing selection, production dependency selection, safe zero value, and
   local-config error ordering. Preserve all current refresh contracts.
2. Add the smallest private cache-presence dependency to the existing refresh
   dependency bundle. Production must delegate exactly once to the existing
   `file.Exists` behavior for `file.Path("%s/.git", target)`. Keep the existing
   Git dependency, URL/target order, logs, `shell.Output.FormatError`, public
   facade, callers, and cached reads unchanged. Update only the corresponding
   branch mutation expression; retain the same ten mutation IDs and meanings.
3. Run focused config/context/command tests and the exact config-cloud mutation
   harness, then the full recorded exact-Go-1.26.7 gates under `umask 022`,
   readonly/offline module inputs, and managed scratch. Verify no API/CLI,
   source-scope, fixture, Go-floor, dependency-metadata, graph, or quality
   regression. Roll back the single implementation commit if the private seam
   cannot preserve all contracts; do not compensate in another slice.

# Automatic Handoff

If and only if the one slice and all gates pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If any contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge, publish,
release, stash, revert, remove the worktree, integrate `ply-config`, or start a
second cloud, Spring, or packaging slice.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The bounded slice is complete at focused implementation commit
`57f9d5674157d238a2f93462b65161a17e3b5498` (parent
`d385616b3191be9a8a2a0b42ced411c3aacd9f11`, tree
`bbd931ef16a334079473e1cd35aa387ed6b61ef9`). It changes exactly
`pkg/config/cloud.go`, `pkg/config/cloud_refresh_test.go`, and
`scripts/mutate-config-cloud`. No second P8 slice, `ply-config` work, caller or
reader change, fixture change, dependency evaluation, Spring work, or
packaging work ran.

### Implementation And Characterization

- `refreshGitDependencies` now contains one private `refreshCacheProbe`. Its
  zero value reports absent, while `systemRefreshGitDependencies` selects the
  existing shell Git dependency and a private `fileRefreshCache` whose sole
  operation returns `file.Exists(path)`.
- `GitCloudConfig.Refresh` still parses local config first, derives the same
  target, composes exactly `file.Path("%s/.git", target)`, and probes once.
  Present selects one pull; absent selects one clone with URL before target.
  Logs, `shell.Output.FormatError`, return values, the public facade, callers,
  and cached reads are unchanged.
- The focused tests first failed because the cache field, production selector,
  adapter, and safe `Exists` method did not exist. They now record one exact
  probe for both branches, production adapter selection and present/missing
  behavior, safe zero-value operation without developer-path access, and
  local-config error return before any probe or Git call. Existing Git call,
  clone argument-order, and formatted-error contracts remain intact.
- Only the source expression of mutation `cloud-refresh-existing-branch`
  changed to the injected probe. All ten mutation IDs, selected tests, and
  meanings remain exact.

### Verification And Protected State

- Focused `pkg/config`, `pkg/context`, and `cmd/...` tests pass. The direct
  config-cloud harness reports `declared=10 killed=10 survived=0 unusable=0`;
  its T1-T10 meta-test passes with the same totals.
- Exact Go 1.26.7 module verification, build, uncached count-one tests, race
  count-one tests, and vet pass under `umask 022`, offline module resolution,
  and managed scratch. API/CLI compatibility, their meta-tests, the complete
  `make preflight`, `make test`, install, and 62-check launcher gates pass.
  Empty-HOME hermetic count-two tests pass.
- All 15 quality-audit meta-controls pass. The canonical audit has its
  documented finding exit 1, not audit-broken exit 2; authoritative structured
  Q0-Q2 contains no automated FAIL and no ratchet regression. The six manual
  Q0-Q2 contracts remain supported: this slice strengthens the recorder,
  default, and non-empty evidence; all eight mutation meta-suites retain
  10/10 kills; partial-failure and acceptance magnitude/read-only controls are
  unchanged. Accepted 27/27 Q0-Q2 PASS at L2 therefore has no regression.
- Discarded diagnostic preflight runs isolated environment mismatches rather
  than product findings: ambient `GOFLAGS=-mod=readonly` made the meta-test's T11
  before/after environment unequal, and resolving plain `go` to the Homebrew
  1.26.2 launcher made T15 reproduce only the wrong version string. The final
  gate used no ambient `GOFLAGS` as required, put the verified Go 1.26.7
  binary first on `PATH`, kept `GOPROXY=off`, and passed T1-T15.
- Go 1.18 and the protected 234 modules, 3,599 graph edges, 355 production
  entries, 429 complete-test entries, 197 module-backed entries across 41
  loaded modules, and 1,067 sum lines reproduce. `go.mod`, `go.sum`, and graph
  hashes remain respectively
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  Public API/CLI reports pass; fixtures and dependency metadata are byte-exact.
- The task-owned scratch root remained beneath managed session scratch. Its
  40,041 entries contained no special files and three intentional audit-fixture
  command symlinks; the non-symlink root and exact containment were verified
  before the complete task subtree was removed.

The implementation rollback boundary remains the single focused commit. The
next checkpoint is planning-only and may select at most one further owned
cloud slice; it may not execute that slice or inspect/select `ply-config`.
