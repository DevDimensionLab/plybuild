# Agent Session: Implement Cloud Template Lookup Loader Seam

Status: ANSWERED - HISTORY
Session ID: `2026-09-24T231158+0200-implement-cloud-template-lookup-loader-seam`
Created: `2026-09-24T23:11:58+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f099e835dbc518d958f6739c452dec874ecc090a71d20998403c20bffa68a050`
Previous: [2026-09-24T225524+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T225524+0200-plan-next-p8-cloud-modernization-move.md)
Next: [2026-09-24T235342+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T235342+0200-plan-next-p8-cloud-modernization-move.md)
Outcome: the ninth bounded P8 slice is complete at one focused two-file commit; `Template` now selects a private zero-value-safe Git-backed template-list loader while every eager walk, project load, lookup, caller, mutation, and public contract remains exact

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Implement exactly the ninth bounded P8 cloud-configuration slice selected by
the answered planning checkpoint: isolate the direct
`GitCloudConfig.Template` call to `Templates()` behind one private zero-value-
safe template-list loader while preserving full eager template walking,
project loading, lookup, errors, repeated calls, downstream callers, and every
public contract. Do not integrate or select `ply-config`, and do not begin a
second slice.

# Authorized Roadmap

P2A-P7 are complete. P8 is active in cloud configuration only. The first eight
P8 slices are complete at `57f9d5674157d238a2f93462b65161a17e3b5498`,
`c12307a078a29af74163df0c658d05c821482a33`,
`c999212d266d98868930769108e4e8a73070a0a2`,
`51da9bcfc0eb8ae871c98467c52c942db6b480bd`,
`f12344b115b6c732295d152a7f3d86b876a4d1d8`,
`76c3571f04969c304785338bc31c34a266e2f869`,
`9c1899d484633e0ce6cf112544ba33ca89bc4adb`, and
`d26d19232b18d70a4401586d75604d08f85bd8a8`. This checkpoint may change only
`pkg/config/cloud.go` and new focused
`pkg/config/cloud_template_lookup_test.go` for the private template-list
loader. It may not change `Templates`, `HasTemplate`, another reader or
constructor, a caller, fixture, mutation file, dependency metadata, Spring, or
packaging.

# Measurements At Start

Begin only from the clean reciprocal planning handoff on
`codex/upgrade-quality`. Verify its branch, exact five-file documentation and
launcher commit shape, answered planning archive, focused template-project-
loader implementation ancestry and exact two-file shape, all earlier P8 and
Google UUID ancestry, sole NEXT state, reciprocal 344-archive chain, launcher
mirror/check, ordinary and ignored cleanliness, and unchanged protected cloud,
project, profile, context, caller, fixture, mutation, and dependency inputs.
The eighth implementation remains commit
`d26d19232b18d70a4401586d75604d08f85bd8a8`, parent
`42857a2c39d252137bf873ab4bb761f514f5f070`, tree
`e6d265798b22dd8690f6f078d6df246332ea600c`, changing only
`pkg/config/cloud.go` and `pkg/config/cloud_templates_test.go`. The project
remains Go 1.18 with 234 modules, 3,599 graph edges, 355 production entries,
429 complete-test entries, 197 module-backed entries over 41 loaded modules,
1,067 `go.sum` lines, and protected `go.mod` / `go.sum` / graph hashes
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Stop for a fresh owning decision if a protected input changed.

# Role And Boundaries

This is one behavior-preserving private lookup dependency seam, not
`ply-config` integration, dependency evaluation, a `Templates` rewrite,
`HasTemplate` work, shared catalog design, caching or lookup optimization,
profile-policy repair, generic reader or constructor refactoring, caller-policy
cleanup, security work, Spring work, or packaging work. Do not change any
exported symbol, signature, or field; the filesystem walk or project loader;
either cloud opener; project/profile/migration/POM behavior; cache layout;
refresh; `HasTemplate`; `ValidTemplatesFrom`; another reader; any command,
context, Maven, or template caller; tracked fixture; mutation file; `go.mod`;
or `go.sum`. Do not fetch, inspect, or select a new dependency. Keep every
disposable cache, report, copy, and test fixture beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}` and remove task-owned scratch after
containment and entry-type verification.

# Required Reading

Read the answered planning archive and its exact selected-slice contract; the
answered template-project-loader implementation and preceding completed P8
implementation archives; current P8 roadmap and rolling handover;
`docs/design/quality-lift.md`; `pkg/config/cloud.go`,
`pkg/config/cloud_templates_test.go`, `pkg/config/cloud_test.go`, and every
other focused cloud test; `pkg/config/project_init.go` and its focused opener
tests; `pkg/config/profiles.go`; `pkg/config/project.go`;
`pkg/context/context.go` and its focused opener tests; active `Template` and
`ValidTemplatesFrom` callers in `cmd`, `pkg/maven`, and `pkg/template`; and
`scripts/mutate-config-cloud` with its meta-test. Treat every decoded-profile,
non-ENOENT, migration, caller-policy, `HasTemplate`, other-reader,
real-network, environment, and fixture gap as a guard rather than permission
to broaden this slice. Do not inspect `ply-config`.

# Three Moves

1. Add focused failing characterization for production loader selection,
   complete arbitrary receiver delivery and one call, complete list and exact
   error delivery, eager list-error precedence over a matching partial list,
   first exact case-sensitive match with complete project and interface
   identities, empty/not-found behavior, independent repeated calls, safe
   missing-loader behavior without path access, and non-empty recordings.
   Retain current/legacy production lookup, downstream partial results, and all
   completed template behavior.
2. Add one private template-list loader interface, one private dependency
   value, one Git-backed production implementation, and one small private
   dependency-taking lookup helper. Production delegates exactly once per
   `Template` invocation to unchanged `gitCfg.Templates()` and returns its
   complete slice and exact error. A missing loader returns nil plus exact
   `filesystem.ErrNoFilesystem` before receiver access. Preserve full eager
   walking/project loading, ordering, identities, exact error precedence,
   first-match/not-found behavior, fresh repeated loads, public wrapper,
   `HasTemplate`, `ValidTemplatesFrom`, and every other boundary.
3. Run focused lookup/config/context/Maven/template/command tests and both
   byte-unchanged config-cloud mutation gates, then the full recorded exact-Go-
   1.26.7 gates under `umask 022`, readonly/offline module inputs, and managed
   scratch. Verify no API/CLI, source-scope, completed-seam, profile,
   project/POM, caller, reader, fixture, Go-floor, dependency, graph, mutation,
   or quality regression. Roll back the single two-file commit if any contract
   fails; do not compensate in another slice.

# Automatic Handoff

If and only if the one seam and every gate pass, record its focused commit and
evidence, answer this archive, update the roadmap and rolling handover, prepare
exactly one reciprocal NEXT archive for a fresh bounded P8 cloud planning
checkpoint, replace only launcher mutable regions, run launcher/handoff checks,
and make one local handoff commit. If any contract is unresolved, prepare one
decision successor instead. Do not execute the successor, push, merge,
publish, release, stash, revert, remove the worktree, inspect or integrate
`ply-config`, repair profile behavior, route `HasTemplate`, alter caller policy,
reopen an earlier slice, or begin a second cloud, Spring, or packaging move.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The ninth bounded P8 cloud slice is complete at focused implementation commit
`96275513bf0e1992ba6864b30c17e35bf2b27447` (parent
`3d98042f2fa10efe2386822fda0deebf7d541365`, tree
`3c82a32c2a67eafd659140f6a92625769b307330`). It changes exactly
`pkg/config/cloud.go` and new focused
`pkg/config/cloud_template_lookup_test.go`, with 305 insertions and one
deletion. No public symbol, `Templates`, `HasTemplate`, caller, constructor,
opener, fixture, mutation file, dependency file, Spring path, packaging path,
or `ply-config` input changed.

### Implementation And Characterization

- `templateListLoader`, `templateLookupDependencies`, and
  `gitTemplateListLoader` are private. Production explicitly selects the
  Git-backed loader, passes the complete `GitCloudConfig` receiver, delegates
  exactly once to unchanged `gitCfg.Templates()`, and returns its complete
  slice and exact error without filtering, copying, caching, refreshing, or
  inspection. A missing loader returns nil plus exact
  `filesystem.ErrNoFilesystem` before receiver-path access.
- Public `GitCloudConfig.Template(name string)` retains its signature and is
  the production wrapper around one small private dependency-taking helper.
  The helper preserves zero `CloudTemplate` and eager error precedence for any
  list error, including a matching partial list; first exact case-sensitive
  match; complete project and embedded interface identities; exact not-found
  text; nil/empty behavior; and one fresh complete load on every invocation.
- Focused TDD first failed only on the absent loader, dependency selector, and
  helper symbols. Final tests cover Git production selection and tracked
  current/legacy lookups; arbitrary complete receiver delivery and one call;
  complete slice and exact error delivery; matching-partial error precedence;
  first exact match and interface identities; nil, empty, case-mismatch, empty-
  name, and exact not-found results; independent repeated loads; safe missing-
  loader behavior without path access; and rejected empty recordings.
- `Templates`, both completed template seams, full eager filesystem walking,
  project construction, ordering, partial/final errors, `HasTemplate`, and
  `ValidTemplatesFrom` remain unchanged. Existing ordered downstream partial
  results, Maven replacement continuation, build/add-template callers, and
  template-package integrations remain exact.

### Verification And Protected State

- The checkpoint began clean on `codex/upgrade-quality` at planning handoff
  HEAD `3d98042f2fa10efe2386822fda0deebf7d541365`, parent
  `d297ccb8a5cdbd4d92c97b9c322f27196148a944`, tree
  `988f3e284e4a18a76ce428e635425f9d52ff341a`, with exactly the launcher,
  answered planning archive, this then-NEXT implementation archive, rolling
  handover, and roadmap changed. The reciprocal graph had 344 archives and one
  NEXT; all earlier P8 and Google UUID commits were ancestral; launcher mirror/
  check, shell syntax, protected inputs, and ordinary/ignored cleanliness
  passed.
- Focused final-commit lookup/config/context/Maven/template/command tests pass.
  The byte-unchanged config-cloud direct harness retains all ten exact IDs,
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
- One discarded compatibility invocation used the Homebrew launcher with
  automatic toolchain resolution while checksum verification was disabled;
  it stopped before comparison. Accepted compatibility and all subsequent
  gates invoked the resolved Go 1.26.7 binary directly with
  `GOTOOLCHAIN=local`. Sandbox-denied `xcrun` cache diagnostics during Make
  recording probes created no outside artifact and their leak guards passed.
- All disposable evidence stayed beneath the managed task root. Its exact
  containment and ordinary entry types were verified, and the task-owned
  subtree was removed before the reciprocal handoff commit.

The rollback boundary remains the single focused two-file implementation
commit. The reciprocal successor is planning-only and may select at most one
further smallest repository-owned cloud slice or stop unresolved; it may not
execute a slice, inspect or select `ply-config`, repair profile behavior, route
`HasTemplate`, reopen an earlier seam, alter caller policy, or begin Spring or
packaging work.
