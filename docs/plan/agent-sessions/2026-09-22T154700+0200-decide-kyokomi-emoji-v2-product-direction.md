# Agent Session: Decide Kyokomi Emoji V2 Product Direction

Status: NEXT
Session ID: `2026-09-22T154700+0200-decide-kyokomi-emoji-v2-product-direction`
Created: `2026-09-22T15:47:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `44b39ecd1791d1a236166c545355e8318df61300eb044800d2a6733a65af6ced`
Previous: [2026-09-22T144301+0200-evaluate-kyokomi-emoji-v2-dependency.md](2026-09-22T144301+0200-evaluate-kyokomi-emoji-v2-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by recording exactly one bounded product decision for selected
exact-path `github.com/kyokomi/emoji/v2 v2.2.12`. The completed evaluation
found that no canonical exact-path stable qualifies: every stable fails two
ordinary documented formatting contracts, while latest v2.2.14 also requires
Go 1.21. Choose only one of the three directions below, apply that choice
exactly, and prepare only its reciprocal successor if the choice requires one.
Do not repeat the dependency evaluation, evaluate another dependency group,
or begin P8.

# Defensive Decision Scope

This is an ordinary dependency-quality product decision. Use only the
completed public release/repository metadata, static source and graph facts,
small bounded ordinary emoji/text behavior results, project projections, and
advisory identities recorded here and in the answered evaluation. Do not fuzz,
stress, probe resource exhaustion, create oversized, deeply nested, cyclic,
malformed, or adversarial inputs, reproduce a security issue, or perform
security or exploitability analysis.

Every disposable cache, tool, report, project copy, fixture, or advisory
response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to
`/private/tmp`, `/tmp`, a sibling of the managed root, or another external
root. Verify containment and remove task-owned scratch evidence before
handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2, and
target-specific option-1 decisions through exact selected, inherited,
unloaded kr/text v0.2.0. P8 remains queued.

Kr/text and kr/pty option 1 decisions are separately exact, unqualified,
target-specific, non-transferable, and final. Their selected and historical
requests, requester imports, complete routes, why/import/load/runtime facts,
release/repository/source identities, qualification results, graph/module/
tidy/Go-floor state, projections, advisories, all earlier guards, and
compatible-route conditions remain expiry guards. Any change requires the
corresponding fresh dependency and product decision before merge. Do not
transfer either exception or reopen kr/text, kr/pty, kr/pretty, Cast, or an
earlier decision.

# Completed Evaluation Is Final

The evaluation began from clean branch `codex/upgrade-quality` at handoff HEAD
`89b6c413eca06d5ddcc6e9d424da70d3bb324447`, parent
`448692f193ecdd9d7f0e608472251ce780ff62cb`, tree
`13f2bb7b54f9cee6a2488d991992502952b21516`. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. The
evaluation changed no product source or dependency metadata and prepared this
decision-only handoff. Verify the new handoff HEAD, parent, tree, exact changed
set, ancestry, reciprocal archive chain, and clean ordinary and ignored status
rather than assuming them.

Go-import and the exact module declaration resolve
`github.com/kyokomi/emoji/v2` to public enabled, unarchived, non-fork MIT
repository `kyokomi/emoji`, ID 21064634, owned by `kyokomi`, defaulting to
`main`, with no parent/source repository. The proxy enumerates earlier v2
tags, but v2.0.0-v2.2.4 lack the `/v2` module path and belong to the
unsuffixed `+incompatible` line. The canonical resolvable exact-path line is
exactly v2.2.5-v2.2.14; `/v3` is absent. There are no retractions,
deprecations, replacements, prereleases, or eligible alternate paths, forks,
branches, or pseudo-versions. GitHub has release objects for every canonical
stable except v2.2.6.

V2.2.5-v2.2.13 declare Go 1.14 with no requirements; v2.2.14 declares Go
1.21 with no requirements. Selected v2.2.12 is verified commit
`c37c65064ac285cd1d62ec6300c94a781d8a8486`, tree
`d61267efe2a6039bc285deab78fe6a69b806681a`. Highest floor-compatible v2.2.13
is unsigned commit `6143718b8c151d51c5a3388d3ebf2947887be8ec`, tree
`31aaf9a84b4e742c02b0ddc987446c9490688220`. Latest v2.2.14 is verified
commit `a659fe56a640109bc44bfecdd0a9641642971f27`, tree
`92215f23d95d0ddac5b6df477e66852ea8000dfd`. Continuous tag ancestry,
proxy/sumdb, regular-file archive/Git manifests, signatures, module files, and
license identities agree as recorded in the evaluation.

The standard-library-only module contains one library and one example main
package; the repository generator is a separate nested module. There is no
cgo, embed, build-tag, platform-specific, or runtime network boundary.
Generated emoji maps and `ReplacePadding` are intentionally exposed mutable
globals; callers own mutation and synchronization. Callers also own writers,
writer errors, and stdout use. Ordinary conversion is deterministic and no
resource lifecycle exists.

Exact Go 1.26.7 and contained Go 1.18.10 upstream build, count-one/count-ten
tests, race-count-ten tests, vet, and supported cross-builds pass for the
owner-requested, selected, and highest floor-compatible serious candidates;
exact Go 1.26.7 gates also pass v2.2.14. Positive bounded fixtures for known
and unknown shortcodes, determinism, error propagation, and independent
concurrent writers pass. Every canonical stable nevertheless fails the same
two documented ordinary contracts under both SDKs and race:
`Fprintln(&buf, "left", "right")` returns `"leftright\n"` instead of the
documented `fmt.Fprintln` result `"left right\n"`, and `Errorf` constructs a
plain error so `%w` does not support `errors.Is` as documented `fmt.Errorf`
does. V2.2.14 independently violates the Go 1.18 floor. No exact-path stable
qualifies.

Exact MVS has two target requests: main -> selected v2.2.12 and direct
supported `github.com/MichaelMure/go-term-markdown v0.1.4` -> v2.2.8. Main
historically introduced target v2.2.8 with that requester before raising only
its target root to v2.2.12. The requester imports the target in `renderer.go`
and applies `emoji.Sprint` to Markdown text. Target and requester why results
are positive through main `cmd` -> direct go-term-markdown -> target. Project
source has no direct target import but directly imports the requester; both
requester and target are production and complete-test loaded and runtime
relevant.

Disposable exact gets were measured and not retained. V2.2.8 changes only the
target root and preserves 234 modules/3,599 edges/355 production/429 complete-
test/197 module-backed/41 loaded-module/1,067-sum state. V2.2.12 is an exact
no-op. V2.2.13 changes only the target root and two checksum lines for
234/3,599/355/429/197/41/1,069; its tidy projection retains 234 modules/3,557
edges/948 sums and swaps only target identity. V2.2.14 raises main `go 1.18`
to `go 1.21`, produces 234/3,601/355/429/197/41/1,069, and tidies to
234/3,559/952. The empty target requirement closure makes the exact target the
only dependency selection delta. Every raw candidate preserves all 39 earlier
guard selections and their exact 237-edge snapshot; every comparable tidy
projection preserves the kr/text and kr/pty paths, versions, requests, routes,
why/import/load/runtime facts, and their comparable guarded tidy state. No projection was
retained.

The unchanged project remains 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The recorded 432-line tidy projection remains exact; the common applied
52/948-line hashes remain
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

All 39 earlier guarded selections remain exact. Their 237 sorted incoming
edges retain SHA-256
`151e72c0b5444acc59ffab45d6ce8fe43e80821b42df6a49a11902b3a6632803`;
37 why results are negative and only closed kr/pretty and kr/text are
positive, with zero guarded repository imports and production/complete-test
loads. Every kr/text and kr/pty expiry guard remains exact.

Exact-version OSV and GitHub global results for all ten canonical stables and
repository advisories are empty. Pinned exact-Go govulncheck v1.8.0 isolated
module/package/symbol/test-symbol scans are empty for selected, v2.2.13, and
v2.2.14; base and disposable-v2.2.13 project populations are identical at
30/22/20/20 with no target finding. Guard OSV retains only the recorded
Gorilla WebSocket and go-retryablehttp pairs; x/mod v0.14.0 retains
GO-2026-6179 and GO-2026-6180. The 518,501-byte/1,402-record Go module index
and PUBLISHED 2,807-byte memberlist CNA response remain byte-exact at their
recorded hashes. Advisory absence does not override ordinary behavior.

Final unchanged-project exact Go 1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass under `umask 022`. The Go 1.18,
source/API/CLI/help/launcher/Make/quality contracts and every earlier decision
remain exact. Accepted quality remains 27/27 Q0-Q2 PASS at L2. Treat all
completed release, source, behavior, closure, route, projection, advisory, and
final-gate findings as final.

# Choose Exactly One Direction

1. Explicitly retain exact selected, direct-indirect, runtime-relevant
   `github.com/kyokomi/emoji/v2 v2.2.12` without product-source or dependency-
   metadata changes under an emoji/v2-specific, non-transferable exception.
   Call it unqualified: no canonical exact-path stable qualifies because all
   ten fail the documented Fprintln and Errorf contracts, and v2.2.14 also
   violates the Go floor. Accept only the completed release/source, ordinary
   behavior, exact owner/request, route, why/import/load/runtime, graph/tidy/
   Go-floor, earlier-guard, and advisory facts. Define the exact expiry guards
   below. This exception must not broaden or expire the closed kr/text or
   kr/pty exceptions.
2. Authorize exactly one later bounded measurement-only genuine owner/request
   study. Name the existing shortest supported route: main -> direct exact
   `github.com/MichaelMure/go-term-markdown v0.1.4` -> exact selected
   `github.com/kyokomi/emoji/v2 v2.2.12` by MVS over its v2.2.8 request. State
   the exact question: whether a later genuine supported tidy-stable
   go-term-markdown selection removes emoji/v2 or requests a future qualified
   Go-1.18-compatible canonical exact-path stable while preserving the direct
   owner contract, all 39 earlier selections/237 edges, both kr/text and
   kr/pty exceptions, the Go floor, and every project contract. Do not run the
   study, change any selection, or reopen an earlier decision in this decision-
   recording move; prepare one reciprocal measurement-only successor.
3. Stop P7 unresolved with no product-source, dependency-metadata, roadmap,
   or P8 implementation change and no successor execution.

Do not invent a fourth option or combine options. Physical selection, positive
why, runtime loading, and advisory absence are not qualification. Do not call
v2.2.12 or v2.2.13 qualified, select v2.2.14, add another direct root, change
go-term-markdown, kr/text, kr/pty, kr/pretty, Cast, the Go floor, product
source, dependency metadata, or another module, transfer an exception, or
begin P8.

# Option 1 Expiry Guards

If option 1 is selected, preserve and record at minimum:

- exact selected path/version `github.com/kyokomi/emoji/v2 v2.2.12`, main's
  direct-indirect request, go-term-markdown v0.1.4 -> v2.2.8, and the
  historical main v2.2.8 -> v2.2.12 change;
- the requester's `renderer.go` import and `emoji.Sprint` use, the complete
  main `cmd` -> direct requester -> target route, no direct repository target
  import, positive target/requester why results, production/complete-test
  loads, and runtime relevance;
- the exact ten-release canonical `/v2` line, separation from the older
  unsuffixed `+incompatible` tags, absent `/v3`, repository/owner/status/
  license/default-branch identity, release/tag/commit/tree/archive/sumdb
  facts, and no replacement/retraction/deprecation or eligible alternate;
- every stable remaining unqualified for the completed Fprintln/Errorf
  failures, v2.2.14 remaining Go-1.21 floor-ineligible, and no future
  qualified Go-1.18-compatible canonical exact-path stable or supported route;
- baseline 234/3,599/355/429/197/41/1,067 state, exact module hashes, common
  tidy state, and all completed disposable-projection results;
- exact Go 1.18 floor, exact Go 1.26.7 identity, source/API/CLI/help/launcher/
  Make/quality contracts, and accepted 27/27 Q0-Q2 PASS at L2;
- all 39 earlier selections/237 incoming edges and their exact hash, plus
  every kr/text and kr/pty expiry guard; and
- no new target/requester advisory, independent defect, release, owner,
  qualified stable, supported tidy-stable owner, or compatible genuine route
  to a qualified target, kr/text, or kr/pty release.

Any target/request/requester-import/owner-route, root/why/import/load/runtime,
graph/module/tidy/Go-floor, earlier or kr/text/kr/pty guard, advisory/finding,
independent defect, repository/release/owner, qualification, supported-owner,
or compatible-route change expires retention and requires a fresh emoji/v2
dependency and product decision before merge. Any kr/text or kr/pty guard
change separately requires its own fresh dependency and product decision.
Option 1 authorizes no owner study, direct-root change, branch or pseudo-
version, alternate path, dependency edit, workaround, unrelated selection,
implementation, or transferred exception.

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

Read this archive and the answered emoji/v2 evaluation; the answered kr/text
and kr/pty decisions/evaluations; the kr/pretty and Cast owner records;
rolling handover; P7/P8 roadmap; `go.mod`; and `go.sum`. Verify branch,
ancestry, clean ordinary/ignored state, reciprocal archive chain, launcher
check, exact Go identities, target requests/routes/loads, module hashes/counts/
tidy state, all 39 earlier selections/237 edges, both closed exceptions, and
fresh advisory identities before recording the choice.

# Three Moves

First, choose exactly one authorized direction using only the final completed
evaluation. Second, record only that decision without changing product source
or dependency metadata or running an owner study. Third, update roadmap and
rolling handover, answer this archive, prepare only the one reciprocal
successor required by the chosen direction, verify containment, and make the
required local handoff commit without executing the successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, remove the worktree, transfer an exception, reopen emoji/v2, kr/text,
kr/pty, kr/pretty, Cast, or an earlier decision, evaluate another dependency
group, write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
