# Agent Session: Evaluate Pkg Errors Dependency

Status: NEXT
Session ID: `2026-09-23T035008+0200-evaluate-pkg-errors-dependency`
Created: `2026-09-23T03:50:08+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `064b48d9109113325396cad3eb9f5632279e5425b3bc3d4e5dc77b487927e8c1`
Previous: [2026-09-23T031134+0200-decide-pascaldekloe-goe-product-direction.md](2026-09-23T031134+0200-decide-pascaldekloe-goe-product-direction.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with one bounded evaluation of exact selected
`github.com/pkg/errors v0.9.1`. Determine whether a genuine supported upgrade,
replacement, removal, or exact retention direction qualifies while preserving
the Go 1.18 compatibility floor and every earlier target-specific guard. Do
not make the product decision, implement a dependency change, evaluate another
dependency group, or begin P8.

# Defensive Scope

Use public metadata, static repository/source inspection, ordinary documented
behavior, and normal project/upstream graph, build, test, race, vet, and
vulnerability commands. Do not fuzz, stress, probe resource exhaustion,
construct oversized, deeply nested, cyclic, malformed, adversarial, or escape-
sequence payloads, reproduce a security issue, or perform security or
exploitability analysis.

Every disposable cache, tool, response, report, module copy, and project copy
must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to
`/private/tmp`, `/tmp`, a sibling of the managed root, or another external
root. Verify containment and remove task-owned scratch evidence before
handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
the nine final target-specific option-1 decisions through Goe, and qualified
no-selection-change modern-go/concurrent and modern-go/reflect2. P8 remains
queued.

Exact Goe v0.1.0, ULID, go-conntrack, mapstructure, go-homedir, Promptui,
emoji/v2, kr/text, kr/pty, kr/pretty, kr/logfmt, kr/fs, go-windows-terminal-
sequences, gotool, errcheck, httprouter, GLS, go-junit-report, json-iterator,
and clockwork exceptions are final, target-specific, unqualified, and non-
transferable. Every concurrent, reflect2, Cast, Viper, memberlist, earlier
selection, owner, request, route, why/import/load/runtime, graph/module/tidy/
Go-floor, advisory, source, behavior, closure, qualification, and expiry guard
remains final. Do not reopen, broaden, or transfer any decision.

The Goe decision explicitly retains exact selected, inherited, unloaded
v0.1.0 without product-source or dependency-metadata changes. Neither exact-
path stable qualifies: v0.1.0 fails the Go 1.26.7 example and complete race
gates; v0.1.1 still fails race and has no genuine supported tidy-stable
project owner. All four target requests, requester test-import/metadata-only
boundaries, every recorded route, no root, negative why/import/load/runtime
facts, repository/release/archive/source/behavior/closure identities, both
projections, advisories, and earlier guards are retention conditions. The
Armon Go-Metrics Goe Ownership Study is not authorized.

# Fixed Project Guards

The unchanged project has 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy diff SHA-256 is
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes are
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

All 47 pre-Goe guarded selections and 276 sorted incoming edges have SHA-256
`7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
Forty-one why results are negative; only kr/pretty, kr/text, emoji/v2,
Promptui, go-homedir, and mapstructure are positive. Promptui and go-homedir
are the only guarded repository imports. Emoji/v2, Promptui, go-homedir, and
mapstructure are the only production/complete-test loaded guarded modules.
Add the complete Goe selection/request/route/import/load snapshot as its own
guard; do not fold it into or replace the recorded pre-Goe identity.

Guard OSV remains limited to the recorded Gorilla WebSocket and go-
retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and GO-2026-6180. The
Go vulnerability index is 518,501 bytes/1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response is 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Advisory absence is not qualification.

# Evaluation Contract

Begin with guard-only reproduction of the new handoff and fixed state. Stop
for a fresh owning evaluation if any earlier guard changed. Then answer, with
reproducible evidence, only the questions needed for `github.com/pkg/errors`:

1. Resolve the exact physical selection, direct/indirect/root state, every
   current and relevant historical request, MVS outcome, requester selection,
   complete owner routes, requester imports or metadata-only boundaries,
   `go mod why`, repository imports, production/complete-test/module-backed
   loads, and runtime relevance.
2. Establish canonical module/repository ownership and status, license,
   default branch, exact stable/prerelease/major lines, replacements,
   retractions, deprecations, redirects, eligible alternates, tags, releases,
   commits, trees, ancestry, signatures, module files, proxy/sumdb identities,
   and archive-to-Git source identity. Do not treat a branch, fork, replacement,
   alternate path, or pseudo-version as a stable release.
3. Compare selected v0.9.1 with every genuine supported Go-1.18-compatible
   candidate. Inspect exported API, build constraints, cgo/generated/embed
   boundaries, dependency closure, and ordinary documented behavior that can
   materially affect this project. Use only small deterministic fixtures.
4. Run complete upstream verification, build, count-one and repeated ordinary
   tests where meaningful, race, and vet under exact Go 1.18.10 and Go 1.26.7,
   plus supported cgo-disabled cross-build/test-compile checks justified by
   the package boundary. Record every failure; a partial package pass does not
   qualify a release whose complete required gate fails.
5. In disposable project copies, measure selected/no-op, candidate upgrade,
   removal, and tidy projections that are supported by a real owner route.
   Record exact module/edge/load/module-file/sum changes, owner stability,
   Go-floor effects, earlier-guard preservation, and restoration behavior.
   Retain none of these projections.
6. Refresh exact target OSV, GitHub global and repository advisories, isolated
   pinned govulncheck, project govulncheck populations/traces, guard OSV/x/mod,
   Go-index, and memberlist-CNA identities. Advisory absence cannot qualify a
   candidate or erase an ordinary defect.

A qualifying direction must preserve the Go 1.18 floor, pass its complete
source/API/closure/behavior/build/test/race/vet contract, have a genuine
supported tidy-stable project owner where selection matters, preserve every
earlier guard, and introduce no unreviewed product or CLI behavior. If no
direction qualifies, say so plainly and prepare a bounded product decision;
do not silently retain, upgrade, remove, replace, patch, vendor, or add a root.

# Measurements At Start

The Goe decision began from clean branch `codex/upgrade-quality` at handoff
HEAD `413fb1a154a5c9a4acfe151dfb9518714a701ed1`, parent
`536df39770daedc4f841a6898e5c69b8605eed6f`, tree
`0f6a53bb3eaa70c43d1bc91047ca25fa399f6c0f`. Google UUID implementation
commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. The Goe
decision changed no product source or dependency metadata. Independently
verify the new handoff HEAD, parent, tree, exact changed set, ancestry, clean
ordinary and ignored status, reciprocal archive chain, launcher check, exact
Go identities, fixed project and Goe guards, advisories, and final gates.

# Role And Boundaries

Act only as the bounded pkg/errors evaluator. Measure and report the target's
real choices without selecting one, editing source or dependency metadata,
changing an owner/requester, or weakening a prior guard. Stop if the evidence
requires a broader product, architecture, security, or cross-dependency choice.

# Required Reading

Read this archive, the answered Goe decision/evaluation, answered ULID and go-
conntrack decisions/evaluations, reflect2 and concurrent evaluations,
mapstructure, go-homedir, Promptui, emoji/v2, kr/text, and kr/pty decisions/
evaluations, rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Treat
completed facts as final unless a narrow guard proves an input changed.

# Three Moves

First, reproduce the handoff and fixed guards. Second, complete and record only
the bounded pkg/errors evaluation. Third, update the roadmap and rolling
handover, answer this archive, prepare exactly one reciprocal successor for the
product decision required by the evidence, verify scratch containment, and
make the local handoff commit. Do not make that decision or execute the
successor. If an earlier guard changed, stop with no successor and request the
necessary owning reevaluation.

# Automatic Handoff

Do not push, merge, publish, release, stash, revert, bypass cleanup, remove the
worktree, change product source or dependency metadata, add a root, change a
requester, transfer an exception, evaluate another dependency, begin P8, or
write outside the managed scratch root.
<!-- CODEX_SESSION_PROMPT_END -->
