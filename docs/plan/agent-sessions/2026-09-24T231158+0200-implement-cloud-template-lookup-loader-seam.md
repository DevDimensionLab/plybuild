# Agent Session: Implement Cloud Template Lookup Loader Seam

Status: NEXT
Session ID: `2026-09-24T231158+0200-implement-cloud-template-lookup-loader-seam`
Created: `2026-09-24T23:11:58+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f099e835dbc518d958f6739c452dec874ecc090a71d20998403c20bffa68a050`
Previous: [2026-09-24T225524+0200-plan-next-p8-cloud-modernization-move.md](2026-09-24T225524+0200-plan-next-p8-cloud-modernization-move.md)
Next: none
Outcome: pending

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
