# Agent Session: Decide Manifoldco Promptui Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-22T171632+0200-decide-manifoldco-promptui-product-direction`
Created: `2026-09-22T17:16:32+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `990cb334e95c0e2e2488dfeee747f59f0e0e95fc63d5813a88434b054fe61751`
Previous: [2026-09-22T162048+0200-evaluate-manifoldco-promptui-dependency.md](2026-09-22T162048+0200-evaluate-manifoldco-promptui-dependency.md)
Next: [2026-09-22T180017+0200-evaluate-mitchellh-go-homedir-dependency.md](2026-09-22T180017+0200-evaluate-mitchellh-go-homedir-dependency.md)
Outcome: Option 1 explicitly retains exact selected direct runtime-relevant
  Promptui v0.9.0 under a target-specific unqualified exception; product and
  dependency metadata remain unchanged, and one reciprocal P7 evaluation was
  prepared without execution.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by recording exactly one bounded product decision for selected
direct exact-path `github.com/manifoldco/promptui v0.9.0`. The completed
evaluation found that no canonical exact-path stable qualifies: v0.9.0 is the
latest Go-1.18-compatible stable but loses an ordinary documented output error
and races internally during a single supported Prompt run; v0.3.2-v0.8.0
repeat those failures, while v0.1.0-v0.3.1 lack a complete standalone release
module closure. Choose only one of the three directions below, apply that
choice exactly, and prepare only its reciprocal successor if the choice
requires one. Do not repeat the dependency evaluation, evaluate another
dependency group, or begin P8.

# Defensive Decision Scope

This is an ordinary dependency-quality product decision. Use only the
completed public release/repository metadata, static source and graph facts,
small bounded ordinary prompt/text/terminal behavior results, project
projection, and advisory identities recorded here and in the answered
evaluation. Do not fuzz, stress, probe resource exhaustion, create oversized,
deeply nested, cyclic, malformed, adversarial, or escape-sequence payloads,
reproduce a security issue, or perform security or exploitability analysis.

Every disposable cache, tool, report, project copy, fixture, or advisory
response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to
`/private/tmp`, `/tmp`, a sibling of the managed root, or another external
root. Verify containment and remove task-owned scratch evidence before
handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and target-specific option-1 decisions through exact selected, direct-
indirect, runtime-relevant emoji/v2 v2.2.12. P8 remains queued.

Emoji/v2, kr/text, and kr/pty option 1 decisions are separately exact,
unqualified, target-specific, non-transferable, and final. Their selected and
historical requests, requester imports, complete routes, why/import/load/
runtime facts, release/repository/source identities, qualification results,
graph/module/tidy/Go-floor state, projections, advisories, all earlier guards,
and compatible-route conditions remain expiry guards. Any change requires the
corresponding fresh dependency and product decision before merge. Do not
transfer an exception or reopen emoji/v2, kr/text, kr/pty, kr/pretty, Cast, or
an earlier decision.

# Completed Evaluation Is Final

The evaluation began from clean branch `codex/upgrade-quality` at handoff HEAD
`7f8a3a8dbeaa07ff78a5d384f51d3df72cb538ef`, parent
`d3cb8122dcca2e06b29787facc3152bf2e38b947`, tree
`962beb97bf678f87c6bdd43c440522a17e5f82f2`. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. The
evaluation changed no product source or dependency metadata and prepared this
decision-only handoff. Verify the new handoff HEAD, parent, tree, exact changed
set, ancestry, reciprocal archive chain, and clean ordinary and ignored status
rather than assuming them.

Exact go-import metadata and module declarations resolve
`github.com/manifoldco/promptui` to public, enabled, unarchived, non-fork
BSD-3-Clause repository `manifoldco/promptui`, ID 107154646, owned by
`manifoldco`, defaulting to `master`, with no parent/source. The proxy exposes
exactly v0.1.0, v0.2.0, v0.2.1, v0.3.0, v0.3.1, v0.3.2, v0.4.0, v0.5.0,
v0.6.0, v0.7.0, v0.8.0, and v0.9.0; `@latest` is v0.9.0. Exact `/v2` and
`/v3` lines are absent. There are no retractions, module deprecations,
prereleases, or eligible alternate paths, forks, branches, pseudo-versions,
or replacements. Nine releases have GitHub Release objects; v0.3.1-v0.4.0
exist as genuine tags without them.

All tags are lightweight and form continuous ancestry. Selected/latest v0.9.0
is unsigned commit `c2e487d3597f59bcf76b24c9e80679740a72212b`, tree
`becf36d02c10091c82f8eeccba0e67bead88d418`. All twelve proxy ZIPs byte-match
the regular-file manifests of their exact Git tags. V0.9.0 source/module sums
are `h1:3V4HzJk1TtXW1MTZMP7mdlwbBpIinw3HztaIlYthEiA=` and
`h1:ka04sppxSGFAtxX0qhlYQjISsg9mR4GWtQEhdbn6Pgg=`. The answered evaluation
records every commit/tree/signature, archive, license, and manifest identity.

V0.1.0-v0.3.1 have no release `go.mod`, and their synthesized proxy module
files contain no requirements. V0.3.2-v0.4.0 include old tool requirements
and an exclude; v0.5.0 declares Go 1.11 and includes a go-i18n replacement;
v0.6.0 declares Go 1.11; v0.7.0-v0.9.0 declare Go 1.12. V0.9.0 requires
chzyer/logex v1.1.10, chzyer/readline
v0.0.0-20180603132655-2972be24d48e, chzyer/test
v0.0.0-20180213035817-a1ea475d72b1, and x/sys
v0.0.0-20181122145206-62eef0e2fa9b. Its loaded production/test package closure
spans Promptui plus that readline pseudo-version, while its module build list
retains all four requirements. Main's separate direct readline v1.5.1 request
wins MVS in the project.

V0.9.0 contains root, list, and screenbuf library packages plus examples and
exports Prompt, PromptTemplates, Select, SelectWithAdd, list/screen-buffer
helpers, and mutable style/key/search/icon/function/cursor globals. Callers own
global mutation/synchronization and caller-provided streams. Readline is the
external terminal boundary; there is no cgo, generated, embed, subprocess, or
network boundary. Windows/non-Windows source splits are supported; old
readline does not support js/wasm.

Exact Go 1.26.7 and contained Go 1.18.10 module verification, build, upstream
count-one/count-ten, race-count-ten, vet, and supported Darwin, Linux,
Windows, and FreeBSD cross-build/test gates pass v0.9.0. Its complete
production/test closure is 84/144 packages, 10 module-backed entries, and two
modules under Go 1.26.7, and 64/98/10/two under Go 1.18.10. V0.3.2-v0.8.0
pass the equivalent available upstream gates under both SDKs.

A bounded ordinary fixture uses only `ordinary\n` input, simple documented
templates, and a caller writer returning an ordinary error. Positive input,
validation, and deterministic rendering pass. Every v0.3.2-v0.9.0 release
nevertheless discards that writer failure and returns nil despite Prompt.Run's
documented execution-error result. A single Prompt run, without caller
concurrency or mutation, also races between readline's listener goroutine and
Prompt rendering under both SDKs. V0.1.0-v0.3.1 independently fail standalone
release build/test module resolution because their release archives and
synthesized module files provide no requirements. No canonical exact-path
stable qualifies.

The sole target request is main -> selected v0.9.0. Main first requested
v0.8.0, removed it, then reintroduced v0.9.0; the current source import and
prompt functions followed later. Repository source imports Promptui only in
`cmd/root.go`, configures PromptTemplates, and calls Prompt.Run in value and
confirmation paths. `go mod why -m` resolves main `cmd` -> Promptui. Root,
list, and screenbuf are production and complete-test loaded, with three
module-backed entries. The direct dependency is runtime relevant.

A disposable exact-Go-1.26.7
`go get github.com/manifoldco/promptui@v0.9.0` is a byte no-op. The unchanged
project remains 234 selected modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41
loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy projection remains
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes remain
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
No projection or implementation was retained.

All 40 earlier guarded selections remain exact. Their 239 sorted incoming
edges retain SHA-256
`deb453cac4dcb6df8a6b502f6f20aba72eaeb415ac2ff40f873a677113cf9d4a`.
Thirty-seven why results are negative; only closed kr/pretty, kr/text, and
emoji/v2 are positive. Guarded repository imports are zero and only emoji/v2
is production/complete-test loaded. Every emoji/v2, kr/text, and kr/pty expiry
guard remains exact.

Exact-version OSV is empty for every Promptui stable, Promptui's requested
readline pseudo-version, and project-selected readline v1.5.1; GitHub global
Promptui/readline and Promptui repository advisory results are empty. Pinned
isolated govulncheck v1.8.0 is empty for target module/package/symbol/test-
symbol scans; the unchanged project remains 30/22/20/20 with no target trace.
Guard OSV
retains only the recorded Gorilla WebSocket and go-retryablehttp pairs; x/mod
v0.14.0 retains GO-2026-6179 and GO-2026-6180. The 518,501-byte/1,402-record
Go index and PUBLISHED 2,807-byte memberlist CNA response remain byte-exact at
their recorded hashes. Advisory absence does not override ordinary behavior.

Final unchanged-project exact Go 1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass under `umask 022`. The Go 1.18,
source/API/CLI/help/launcher/Make/quality contracts and every earlier decision
remain exact. Accepted quality remains 27/27 Q0-Q2 PASS at L2. Treat all
completed release, source, behavior, closure, route, projection, advisory, and
final-gate findings as final.

# Choose Exactly One Direction

1. Explicitly retain exact selected, direct, production-loaded, runtime-
   relevant `github.com/manifoldco/promptui v0.9.0` without product-source or
   dependency-metadata changes under a Promptui-specific, non-transferable
   exception. Call it unqualified: no canonical exact-path stable qualifies.
   Accept only the completed release/source, ordinary behavior, exact owner/
   request, route, why/import/load/runtime, graph/tidy/Go-floor, earlier-guard,
   and advisory facts. Define the exact expiry guards below. This exception
   must not broaden or expire the closed emoji/v2, kr/text, or kr/pty
   exceptions.
2. Authorize exactly one later bounded measurement-only genuine owner/request
   study named `Ply Direct Promptui Ownership Study`. Its sole route is the
   existing main -> direct exact `github.com/manifoldco/promptui v0.9.0`
   request and the `cmd/root.go` value/confirmation Prompt.Run runtime use.
   Its exact question is whether that direct root and runtime route can be
   removed through a genuine supported project-owner change while preserving
   the public API/CLI behavior, Go 1.18 floor, all 40 earlier selections/239
   edges, and every emoji/v2, kr/text, and kr/pty exception. Do not run the
   study, change source or any selection, introduce a fork/replacement, or
   reopen an earlier decision in this decision-recording move; prepare one
   reciprocal measurement-only successor.
3. Stop P7 unresolved with no product-source, dependency-metadata, roadmap,
   or P8 implementation change and no successor execution.

Do not invent a fourth option or combine options. Physical selection,
directness, positive why, runtime loading, and advisory absence are not
qualification. Do not call v0.9.0 qualified, select an earlier release, add or
change a root, use a fork, replacement, branch, pseudo-version, or alternate
path, change readline, emoji/v2, kr/text, kr/pty, kr/pretty, Cast, the Go floor,
product source, dependency metadata, or another module, transfer an exception,
or begin P8.

# Option 1 Expiry Guards

If option 1 is selected, preserve and record at minimum:

- exact selected path/version `github.com/manifoldco/promptui v0.9.0`, main's
  sole direct request, target's four exact requirements, and main's direct
  readline v1.5.1 MVS selection;
- the historical main v0.8.0 request/removal and v0.9.0 reintroduction, sole
  `cmd/root.go` import, PromptTemplates and Prompt.Run uses, complete main
  `cmd` -> target route, positive target why, three production/complete-test
  target packages, and runtime relevance;
- the exact twelve-release line, absent `/v2` and `/v3`, repository/owner/
  status/license/default-branch identity, release/tag/commit/tree/signature/
  archive/sumdb/module facts, and no replacement/retraction/deprecation or
  eligible alternate;
- every stable remaining unqualified for the completed module-closure or
  ordinary output-error/race failures, and no future qualified
  Go-1.18-compatible canonical exact-path stable or supported route;
- baseline 234/3,599/355/429/197/41/1,067 state, exact module hashes, common
  tidy state, no-op v0.9.0 projection, and no dependency implementation;
- exact Go 1.18 floor, exact Go 1.26.7 identity, source/API/CLI/help/launcher/
  Make/quality contracts, and accepted 27/27 Q0-Q2 PASS at L2;
- all 40 earlier selections/239 incoming edges and their exact hash, plus
  every emoji/v2, kr/text, and kr/pty expiry guard; and
- no new target/closure advisory, independent defect, release, owner,
  qualified stable, supported owner, or compatible genuine route to a
  qualified Promptui, emoji/v2, kr/text, or kr/pty release.

Any target/request/import/owner-route, root/why/import/load/runtime, graph/
module/tidy/Go-floor, earlier or emoji/v2/kr/text/kr/pty guard, advisory/
finding, independent defect, repository/release/owner, qualification,
supported-owner, or compatible-route change expires retention and requires a
fresh Promptui dependency and product decision before merge. Any emoji/v2,
kr/text, or kr/pty guard change separately requires its own fresh dependency
and product decision. Option 1 authorizes no owner study, root change, branch
or pseudo-version, alternate path, dependency edit, workaround, unrelated
selection, implementation, or transferred exception.

# Role And Boundaries

This session is decision recording, not dependency implementation or a new
evaluation. Revalidate only the minimum continuity, exact graph/guard,
advisory-identity, exact-Go final-gate, archive, launcher, and containment
facts needed to make the one decision durable. Do not repeat completed
behavior fixtures, upstream gates, archive comparisons, candidate projections,
or govulncheck analysis. Do not run option 2's owner study during this move.

If option 1 is selected, update roadmap and rolling handover, answer this
archive, and prepare exactly one reciprocal next bounded P7 dependency
evaluation without executing it. If option 2 is selected, prepare exactly one
reciprocal measurement-only owner/request successor without executing it. If
option 3 is selected, record the stop and prepare no implementation successor.
None of the options authorizes P8.

# Required Reading

Read this archive and the answered Promptui evaluation; the answered emoji/v2,
kr/text, and kr/pty decisions/evaluations; the kr/pretty and Cast owner
records; rolling handover; P7/P8 roadmap; `go.mod`; and `go.sum`. Verify
branch, ancestry, clean ordinary/ignored state, reciprocal archive chain,
launcher check, exact Go identities, target requests/routes/loads, module
hashes/counts/tidy state, all 40 earlier selections/239 edges, all three closed
exceptions, and fresh advisory identities before recording the choice.

# Three Moves

First, choose exactly one authorized direction using only the final completed
evaluation. Second, record only that decision without changing product source
or dependency metadata or running an owner study. Third, update roadmap and
rolling handover, answer this archive, prepare only the one reciprocal
successor required by the chosen direction, verify containment, and make the
required local handoff commit without executing the successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, remove the worktree, transfer an exception, reopen Promptui,
emoji/v2, kr/text, kr/pty, kr/pretty, Cast, or an earlier decision, evaluate
another dependency group, write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 was selected on 2026-09-22. Exact selected, direct, production-
loaded, runtime-relevant `github.com/manifoldco/promptui v0.9.0` is explicitly
retained without changing product source, `go.mod`, or `go.sum`. It remains
unqualified and is not described as secure: v0.9.0 is the latest Go-1.18-
compatible canonical stable but discards an ordinary caller-provided output
error despite `Prompt.Run`'s documented execution-error result and races
internally during one supported Prompt run. V0.3.2-v0.8.0 share both failures;
v0.1.0-v0.3.1 lack a complete standalone release module closure. Physical
selection, directness, positive why, runtime loading, and advisory absence do
not qualify or implicitly accept the module.

The Promptui-specific, non-transferable exception accepts only the completed
release/repository/source, ordinary behavior, API/closure/platform, owner/
request, route, why/import/load/runtime, graph/tidy/Go-floor, earlier-guard,
emoji/v2/kr-text/kr-pty-guard, and advisory findings. It accepts no
uncharacterized behavior, new advisory, or independent defect. It neither
broadens nor expires the closed emoji/v2, kr/text, or kr/pty exceptions and
transfers no other exception.

The exception remains valid only while every one of these facts remains
exact:

- selected exact path/version `github.com/manifoldco/promptui v0.9.0` and
  main's sole direct request; target requirements chzyer/logex v1.1.10,
  chzyer/readline `v0.0.0-20180603132655-2972be24d48e`, chzyer/test
  `v0.0.0-20180213035817-a1ea475d72b1`, and x/sys
  `v0.0.0-20181122145206-62eef0e2fa9b`; and main's separate direct
  readline v1.5.1 request and MVS selection;
- main's historical v0.8.0 request/removal and v0.9.0 reintroduction; sole
  `cmd/root.go` import; PromptTemplates configuration and Prompt.Run value and
  confirmation uses; complete main `cmd` -> target why route; positive target
  why; root, list, and screenbuf production/complete-test loading; three
  target module-backed entries; and runtime relevance;
- public enabled, unarchived, non-fork BSD-3-Clause repository
  `manifoldco/promptui`, ID 107154646, owned by `manifoldco`, defaulting to
  `master`, with no parent/source; exactly v0.1.0, v0.2.0, v0.2.1, v0.3.0,
  v0.3.1, v0.3.2, v0.4.0, v0.5.0, v0.6.0, v0.7.0, v0.8.0, and v0.9.0;
  absent `/v2` and `/v3`; all completed release/tag/commit/tree/signature/
  ancestry/archive/sumdb/module/license identities; nine GitHub Release
  objects; and no retraction, deprecation, replacement, prerelease, eligible
  alternate path, fork, branch, or pseudo-version;
- every canonical stable remaining unqualified for the completed standalone-
  closure or ordinary output-error/race failures; the completed API, closure,
  platform, positive ordinary fixture, upstream-gate, and cross-build results
  remaining exact; and no future qualified Go-1.18-compatible canonical exact-
  path stable or supported route appearing;
- the baseline 234 selected modules, 3,599 graph edges, 355 production
  entries, 429 complete-test entries, 197 module-backed entries across 41
  loaded modules, and 1,067 sum lines;
- `go.mod` / `go.sum` SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  the 432-line tidy projection at
  `3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`,
  and common applied 52/948-line hashes
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`;
- the completed exact-v0.9.0 get remaining a byte no-op, with no projection or
  dependency implementation retained;
- declared Go 1.18 floor; exact Go 1.26.7 darwin/arm64 official archive/
  binary SHA-256 identities
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  unchanged source/API/CLI/help/launcher/Make/quality contracts; and accepted
  27/27 Q0-Q2 PASS at L2;
- all 40 earlier guarded selections exact; 37 negative why results and only
  closed kr/pretty, kr/text, and emoji/v2 positive; zero guarded repository
  imports; only emoji/v2 production/complete-test loaded; and the 239 sorted
  incoming edges at SHA-256
  `deb453cac4dcb6df8a6b502f6f20aba72eaeb415ac2ff40f873a677113cf9d4a`;
- every emoji/v2 expiry guard, including its exact selected/request/importer/
  use/route/why/import/load/runtime/release/behavior/Go-floor/projection/
  advisory/qualification/compatible-route state; every kr/text guard,
  including exact requests/routes, wrapping failure, unreleased fix, zero
  load/runtime state, projection and compatible-route conditions; and every
  kr/pty guard, including exact selected v1.1.1, historical v0.1.0 request,
  unloaded `mc` import, route/why/import/load/runtime state, v1.1.4
  qualification and unsupported-owner result, later failures, projections,
  and compatible-route condition; and
- no new target, readline, requester, or closure advisory; independent
  defect; exact-path stable; repository/release/owner change; qualified
  stable; supported owner; or compatible genuine route to a qualified
  Promptui, emoji/v2, kr/text, or kr/pty release.

Any target path/version, request, import, owner identity or route, root,
why/import/load/runtime fact, graph, module hash, tidy state, Go floor, earlier
or emoji/v2/kr-text/kr-pty guard, advisory, independent finding, repository/
release/owner, qualification, supported owner, or compatible-route change
expires this retention and requires a fresh Promptui dependency and product
decision before merge. Any emoji/v2, kr/text, or kr/pty guard change
separately requires its own fresh dependency and product decision. This
decision authorizes no owner study, root change, branch or pseudo-version,
alternate path, dependency edit, workaround, unrelated selection,
implementation, or transferred exception.

Guard-only revalidation began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at handoff HEAD
`12881ba88e843657cd360c508937eeb00c38742c`, parent
`7f8a3a8dbeaa07ff78a5d384f51d3df72cb538ef`, tree
`a4f1996976f2d9c2ddaaf217b081b82f1d0eed80`. That handoff changes exactly
`codex-dev-start.sh`, the answered Promptui evaluation archive, this then-NEXT
decision archive, rolling handover, and roadmap. The reciprocal archive chain,
sole NEXT state, launcher/archive prompt mirror, exact changed set, clean
status, and launcher check passed. Exact Google UUID v1.4.0 dependency commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor.

A freshly downloaded official Go 1.26.7 archive and binary reproduce the
identities above. Under `GOENV=off`, `GOWORK=off`, `GOTOOLCHAIN=local`, empty
`GOFLAGS`, `LC_ALL=C`, `LANG=C`, scratch-contained caches, and canonical
`umask 022`, the baseline counts, module hashes, tidy projection and applied
tidy hashes, sole target request, four target requirements, target route and
three loaded packages, all 40 earlier selections, and the 239-edge snapshot
reproduce exactly. Including Promptui gives 41 guarded selections and 240
incoming edges at SHA-256
`a4e6042216821077b7a74472a6fcdc60d70644d07a805f6a6e99db1140a71563`.
All three closed exception requests, why/import/load boundaries, and exact
selected versions reproduce without guard drift.

Fresh proxy metadata still exposes exactly the twelve stable releases, latest
v0.9.0, and no `/v2` or `/v3` line. GitHub still reports the exact repository
identity and nine Release objects. Exact-version OSV is empty for all twelve
Promptui stables, Promptui's requested readline pseudo-version, and selected
readline v1.5.1; GitHub global Promptui/readline and Promptui repository
advisory results are empty. Fresh selected emoji/v2, kr/text, and kr/pty OSV
and GitHub global/repository advisory results are also empty. Guard OSV
retains only Gorilla WebSocket `GO-2026-6278` /
`GHSA-w67g-5rqw-f597` and go-retryablehttp `GO-2024-2947` /
`GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0 retains `GO-2026-6179` and
`GO-2026-6180`. The Go module index remains 518,501 bytes/1,402 records at
SHA-256 `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. The PUBLISHED 2,807-byte memberlist CNA
response remains SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
updated 2026-07-08T19:40:16.119Z.

The completed Promptui behavior, release/source/closure, upstream-gate,
archive-comparison, candidate-projection, and govulncheck evaluations were not
repeated; option 2 was neither authorized nor run. No source, dependency
metadata, parent, toolchain declaration, or earlier guard changed. No changed-
selection scorecard applies; accepted quality remains 27/27 Q0-Q2 PASS at L2.
Final unchanged-project exact-Go module verification, build, count-one tests,
race count-one tests, and vet pass.

Every task-owned SDK, cache, archive, report, project copy, fixture, and
advisory response remained beneath the exact managed session scratch root and
was removed before handoff; only the pre-existing launcher-owned Node compile
cache remains.

P7 continues only with the linked reciprocal bounded evaluation of the next
unanswered selected queue item, direct exact
`github.com/mitchellh/go-homedir v1.1.0`. That successor was prepared but not
executed and may not reopen this Promptui decision, emoji/v2, kr/text, kr/pty,
kr/pretty, Cast, or any earlier decision. P8 remains queued.
