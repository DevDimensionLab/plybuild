# Agent Session: Implement Cloud Valid Templates Loader Seam

Status: ANSWERED - HISTORY
Session ID: `2026-09-25T001748+0200-implement-cloud-valid-templates-loader-seam`
Created: `2026-09-25T00:17:48+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c0b4e2c3ec301e7c5483b4a49ded99c7ada7a1901714d14778e9b432713adbab`
Previous: [2026-09-24T235342+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T235342+0200-plan-next-p8-cloud-modernization-move.md)
Next: [2026-09-25T005818+0200-plan-next-p8-cloud-modernization-move.md](2026-09-25T005818+0200-plan-next-p8-cloud-modernization-move.md)
Outcome: Completed the tenth bounded P8 cloud slice at focused commit `fd69ac206216260b4f1966214062149264e3e347`; the reciprocal successor is planning-only.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the tenth bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the direct
`GitCloudConfig.ValidTemplatesFrom` call to `Template` behind one private zero-
value-safe per-name template loader while preserving first-occurrence
deduplication, one fresh complete template load per unique name, ordered
results, partial-result and error behavior, the active build caller, and every
public contract. Do not inspect, integrate, or select `ply-config`, and do not
begin a second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first nine
bounded P8 slices are complete at
`57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`,
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`,
`f12344b115b6c732295d152a7f3d86b876a4d1d8`,
`76c3571f04969c304785338bc31c34a266e2f869`,
`9c1899d484633e0ce6cf112544ba33ca89bc4adb`,
`d26d19232b18d70a4401586d75604d08f85bd8a8`, and
`96275513bf0e1992ba6864b30c17e35bf2b27447`. This checkpoint may change only
`pkg/config/cloud.go` and new focused
`pkg/config/cloud_valid_templates_test.go` for the private per-name loader. It
may not change `unique`, `Template`, `Templates`, `HasTemplate`, another reader
or constructor, a caller, fixture, mutation file, dependency metadata, Spring,
or packaging.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, answered planning archive, focused template-list-loader
implementation ancestry and exact two-file shape, all earlier P8 and Google
UUID ancestry, sole NEXT state, reciprocal 346-archive chain, launcher mirror/
check, ordinary and ignored cleanliness, and unchanged protected cloud,
project, profile, context, caller, fixture, mutation, and dependency inputs.
The ninth implementation remains commit
`96275513bf0e1992ba6864b30c17e35bf2b27447`, parent
`3d98042f2fa10efe2386822fda0deebf7d541365`, tree
`3c82a32c2a67eafd659140f6a92625769b307330`, changing only
`pkg/config/cloud.go` and `pkg/config/cloud_template_lookup_test.go`. The
project remains Go 1.18 with 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries over 41 loaded
modules, 1,067 `go.sum` lines, and protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving private composition seam, not `ply-config`
integration, dependency evaluation, lookup optimization, shared/exported
catalog design, `unique` refactoring, `HasTemplate` work, profile-policy repair,
generic reader or constructor refactoring, caller-policy cleanup, security
work, Spring work, or packaging work. Do not change any exported symbol,
signature, or field; the completed template-list or template-project loaders;
filesystem walking; either cloud opener; project/profile/migration/POM
behavior; cache layout; refresh; another reader; any command, context, Maven,
or template caller; tracked fixture; mutation file; `go.mod`; or `go.sum`. Do
not fetch, inspect, or select a new dependency. Keep every disposable cache,
report, copy, and test fixture beneath `${CODEX_SESSION_SCRATCH_ROOT:?}` and
remove task-owned scratch after containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract; the
answered template-list-loader and template-project-loader implementation
archives and preceding completed P8 implementation archives; current P8
roadmap and rolling handover; `docs/design/quality-lift.md`;
`pkg/config/cloud.go`, `pkg/config/cloud_valid_templates_test.go` once created,
`pkg/config/cloud_template_lookup_test.go`, `pkg/config/cloud_templates_test.go`,
`pkg/config/cloud_test.go`, and every other focused cloud test;
`pkg/config/project_init.go` and its focused opener tests;
`pkg/config/profiles.go`; `pkg/config/project.go`;
`pkg/context/context.go` and its focused opener tests; active `Template` and
`ValidTemplatesFrom` callers in `cmd`, `pkg/maven`, and `pkg/template`; and
`scripts/mutate-config-cloud` with its meta-test. Treat every decoded-profile,
non-ENOENT, migration, caller-policy, `HasTemplate`, other-reader, real-network,
environment, and fixture gap as a guard rather than permission to broaden this
slice. Do not inspect `ply-config`.

# Three Moves

1. Add focused failing characterization for production loader selection;
   complete arbitrary receiver and name delivery with one call; exact template,
   embedded interface identity, and error delivery; first-occurrence unique
   order; duplicate suppression; stop-at-first-error with the exact prior
   partial result and discarded failing value; nil/empty behavior; independent
   repeated invocations; safe missing-loader behavior without path access; and
   non-empty recordings. Retain the tracked production partial-result test.
2. Add one private per-name template-loader interface, one private dependency
   value, one Git-backed production implementation, and one small private
   dependency-taking helper. Production delegates exactly once per first-
   occurrence unique name to unchanged `gitCfg.Template(name)` and returns its
   complete value and exact error without inspection, copying, caching,
   refreshing, or normalization. A missing loader returns zero `CloudTemplate`
   plus exact `filesystem.ErrNoFilesystem` before receiver access. Preserve
   unchanged `unique`, fresh full eager walking/project loading per unique name,
   append order, error precedence, public wrapper, callers, and every completed
   seam.
3. Run focused valid-template/config/context/Maven/template/command tests and
   both byte-unchanged config-cloud mutation gates, then the full recorded exact-
   Go-1.26.7 gates under `umask 022`, readonly/offline module inputs, and managed
   scratch. Verify no API/CLI, source-scope, completed-seam, profile,
   project/POM, caller, reader, fixture, Go-floor, dependency, graph, mutation,
   or quality regression. Roll back the single two-file implementation commit
   if any contract fails; do not compensate in another slice.

# Automatic Handoff

If and only if the one seam and every gate pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If any contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge, publish,
release, stash, revert, remove the worktree, inspect or integrate `ply-config`,
repair profile behavior, route `HasTemplate`, change `unique`, alter caller
policy, reopen an earlier slice, or begin a second cloud, Spring, or packaging
move.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The tenth bounded P8 cloud slice is complete at focused implementation commit
`fd69ac206216260b4f1966214062149264e3e347` (parent
`7be6dbf91e5cfbc8e6e0b00c7400d27bcecc8e63`, tree
`43fb3cb60255482590737012691b8c67ce8c9116`). It changes exactly
`pkg/config/cloud.go` and new focused
`pkg/config/cloud_valid_templates_test.go`, with 310 insertions and one
deletion. No public symbol, `unique`, `Template`, `Templates`, `HasTemplate`,
caller, opener, project/profile/POM path, fixture, mutation file, dependency
file, Spring path, packaging path, or `ply-config` input changed.

### Implementation And Characterization

- `perNameTemplateLoader`, `validTemplatesDependencies`, and
  `gitPerNameTemplateLoader` are private. Production explicitly selects the
  Git-backed loader, passes the complete `GitCloudConfig` receiver and exact
  name, delegates exactly once to unchanged `gitCfg.Template(name)`, and
  returns its complete `CloudTemplate` and exact error without inspection,
  copying, caching, refreshing, or normalization. A missing loader returns
  zero `CloudTemplate` plus exact `filesystem.ErrNoFilesystem` before receiver
  access.
- Public `GitCloudConfig.ValidTemplatesFrom(list []string)` retains its
  signature and is the production wrapper around one small private dependency-
  taking helper. The helper keeps unchanged `unique(list)`, case-sensitive
  first-occurrence order, one load per unique name, exact ordered successes,
  and first-error precedence. It returns the exact prior partial result,
  discards any value delivered with the failing error, and does not call later
  names. Nil and empty inputs return nil and nil without a load; separate
  invocations remain independent.
- Focused TDD first failed only on the absent private selector, dependency,
  production loader, and helper. Final tests cover production Git selection
  and the tracked complete template; arbitrary receiver/name delivery and one
  call; exact value/error and embedded interface identities; case-sensitive
  first-occurrence order including duplicate empty names; stop-at-first-error,
  exact partial result, discarded failing value, and no later load; nil/empty
  inputs; independent repeated invocations; safe zero behavior without path
  access; and rejected empty recordings. The existing tracked production
  partial-result test remains unchanged.
- Each first-occurrence unique name still enters public `Template` once, so the
  ninth-slice list loader and both earlier template seams still perform one
  fresh complete eager walk and project load per unique name. Build validation,
  add-template, Maven replacement continuation, `HasTemplate`, every other
  reader, and both cloud openers retain their previous behavior.

### Verification And Protected State

- The checkpoint began clean on `codex/upgrade-quality` at planning handoff
  HEAD `7be6dbf91e5cfbc8e6e0b00c7400d27bcecc8e63`, parent
  `2fbed1cc16d64aef53b1d0048109fab8a498d85f`, tree
  `02e10945e0262ea9b98647f576c0b5ae6f44312e`, with exactly the launcher,
  answered planning archive, this then-NEXT implementation archive, rolling
  handover, and roadmap changed. The ninth implementation retains its exact
  parent/tree/two-file shape. All ten preceding P8/Google UUID ancestry checks,
  the reciprocal 346-archive/sole-NEXT state, launcher mirror/check, shell
  syntax, protected inputs, and ordinary/ignored cleanliness passed.
- Focused valid-template/config/context/Maven/template/command tests pass. The
  byte-unchanged config-cloud direct harness retains all ten exact IDs,
  selections, and meanings with
  `declared=10 killed=10 survived=0 unusable=0`; its T1-T10 meta-test passes
  with the same totals.
- Direct exact Go 1.26.7 module verification, build, count-one tests, race
  count-one tests, and vet pass under `umask 022`, offline module resolution,
  readonly inputs, direct Xcode compilers, and managed scratch. Pinned API/CLI
  compatibility and both meta-tests, the CLI surface, complete preflight,
  ordinary `make test`, explicit `make test-install`, a real scratch-local
  install/help smoke, fresh empty-HOME count-two tests, and all 62 launcher
  controls pass.
- All 15 quality-audit meta-controls pass. The canonical audit returned its
  documented findings exit 1 rather than audit-broken exit 2. Its structured
  Q0-Q2 scope has all 21 automated criteria PASS; the six manual criteria
  remain correctly unclaimed. Seven comparable ratchets improve, and the sole
  overall regression remains the pre-existing Q3.4 documentation indicator.
- Go 1.18 and protected 234 modules, 3,599 graph edges, 355 production entries,
  429 complete-test entries, 197 module-backed entries across 41 loaded
  modules, and 1,067 `go.sum` lines reproduce. `go.mod`, `go.sum`, and graph
  hashes remain respectively
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  API/CLI, completed seams, project/profile/migration/POM behavior, cache
  layout, callers, other readers, fixtures, mutation files, dependency
  metadata, and the Go floor remain unchanged.
- Discarded environment-only attempts reproduced known boundaries rather than
  product defects: forcing the Homebrew launcher local exposed its underlying
  Go 1.26.2 during initial measurement, and the first empty-HOME install smoke
  lacked an explicit existing `GOMODCACHE`. Accepted measurements and gates
  invoked the resolved Go 1.26.7 binary directly, supplied the existing read-
  only module cache where HOME was empty, and changed no repository input.
- All disposable evidence stayed beneath the managed task root. Its exact
  containment and ordinary entry types were verified, and the task-owned
  subtree was removed before the reciprocal handoff commit.

The rollback boundary remains the single focused two-file implementation
commit. The reciprocal successor is planning-only and may select at most one
further smallest repository-owned cloud slice or stop unresolved; it may not
execute a slice, inspect or select `ply-config`, repair profile behavior,
combine or reopen completed seams, alter caller policy, or begin Spring or
packaging work.
