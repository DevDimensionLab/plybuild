# Agent Session: Decide Gorilla WebSocket Product Direction

Status: NEXT
Session ID: `2026-09-14T001519+0200-decide-gorilla-websocket-product-direction`
Created: `2026-09-14T00:15:19+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3349e16ee74218d95380b4a3dcfbc2ec6b60c7830f19e50aa39a2482017059d6`
Previous: [2026-09-13T230600+0200-evaluate-gorilla-websocket-dependency.md](2026-09-13T230600+0200-evaluate-gorilla-websocket-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by obtaining and then recording one bounded product decision
for exact-path `github.com/gorilla/websocket`. The completed evaluation found
no published exact-path stable release whose complete minimal source/test
closure both preserves Go 1.18 and contains the client-mask security fix. Do
not implement a dependency change, audit another group, or begin P8 until the
user chooses one option below.

# Authorized Roadmap

P2A-P6 and all earlier P7 outcomes remain final. P7 is active only for this
Gorilla product decision; P8 remains queued. This session may record one user
choice and prepare one bounded follow-up, but may not implement that choice or
combine another dependency group.

# Measurements At Start

Fresh proxy, sumdb, `go-import`, strict Git, and GitHub evidence resolves the
canonical public active unarchived non-fork BSD-2-Clause repository. Serious
stable releases are v1.4.2, v1.5.0, v1.5.1, and v1.5.3. GitHub marks v1.5.2 a
prerelease. Its proxy identity is original commit `1bddf2e0`, but the forge tag
was overwritten to `9ec25ca`; unreleased main later retracts v1.5.2 for that
reason.

V1.4.2, v1.5.0, and v1.5.3 declare Go 1.12 and have standard-library-only
complete minimal source/test closures under exact Go 1.26.7 and Go 1.18.10.
V1.5.1 and v1.5.2 declare Go 1.20 and have `x/net` closures; their module zips
also have broken default vendor mode. V1.5.2 imports Go-1.20-only
`http.NewResponseController` and fails Go 1.18 compilation. V1.5.3 is the
highest stable floor-preserving release, but it is not security-qualified.

Fresh primary GO-2026-6278/GHSA-w67g-5rqw-f597 evidence describes weak PRNG
use for WebSocket client masks, marks versions before v1.5.3 affected, and
labels v1.5.3 fixed. Release source disproves that boundary: v1.5.3 explicitly
uses `math/rand`, and a seed-controlled independent fixture proves repeatable
mask keys. Security commit
`d67f41855da42d7bccd9ef050c49f7e54e783b95` changes production to
`crypto/rand` only after v1.5.3 on unreleased main. The original proxy v1.5.2
tree contains the fix despite the advisory range, but is a prerelease and
floor-ineligible. Selecting v1.5.3 merely makes the current range green without
fixing the behavior and is not an offered solution.

All packages, examples, tests, benchmarks, generated/build-tag files, exported
APIs, and relevant connection behavior were inspected. Independent fixtures
cover handshake success/failure, headers/auth/cookies, origin/subprotocol,
compression, fragmentation/control interleaving, malformed/oversized frames,
close/deadline/error behavior, prepared messages, pools, deterministic
masking, supported concurrency, and cleanup under both SDKs. A repeated full-
source Go 1.26 race selection finds an upstream cross-test lifecycle defect in
v1.4.2 and v1.5.3: a test handler logs after its owning test returns. The
independent supported-use fixture is race-clean.

Selected remains exact v1.4.2 solely through direct
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3`; `go mod why -m` is
negative, source imports are zero, and zero target packages load. Exact
v1.4.2/v1.5.0/v1.5.3 projections add only a redundant target requirement and
tidy byte-identically to base; v1.5.1/v1.5.2 additionally alter unrelated
`x/*` selections. No candidate was implemented.

The unchanged project has 234 selected modules, 3,599 graph edges, 429
complete-test entries, 197 module-backed packages, 41 loaded modules, 1,067
`go.sum` lines, and a 432-line tidy projection. `go.mod` and `go.sum` remain
byte-identical at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact Go 1.26.7 project and full-preflight gates pass; applicable Go 1.18.10
gates retain only the two accepted shell wording failures. Accepted quality
remains 27/27 Q0-Q2 PASS at L2 with scorecard
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.

# Role And Boundaries

1. Recommended: retain exact selected, inherited, unloaded v1.4.2 without
   dependency metadata changes. Add a non-transferable exception accepting
   only GO-2026-6278's weak `math/rand` client-mask behavior and the documented
   upstream full-source Go 1.26 cross-test lifecycle race. The exception is
   valid only while exact v1.4.2 and its sole mvn-pom-mutator v0.2.3 edge remain
   unchanged, zero target packages load, it remains runtime-unreachable, and
   no new advisory or independent disqualifier appears. Direct import/loading,
   runtime reachability, a target version or incoming-edge change, or a new
   advisory/independent disqualifier expires it and requires a fresh decision.
2. Authorize a broader owning-parent/removal evaluation. This may change or
   remove mvn-pom-mutator and therefore explicitly reopens the final GopherJS
   incoming-edge guard. It is a new joint product boundary, not a Gorilla-only
   dependency move.
3. Authorize evaluation of a maintained patch, fork, replacement, or unreleased
   exact commit. This requires a new provenance, release qualification,
   long-term maintenance, vulnerability, API, MVS, and Go-floor policy before
   any implementation.

Ask the user to choose option 1, 2, or 3 if no choice accompanied this session.
Do not infer acceptance from the recommended label.

## Guarded Decisions

P2A-P6 and all earlier P7 outcomes are final. The user's 2026-09-13 GopherJS
option 1 decision remains valid only for exact
`v0.0.0-20181017120253-0766667cb4d1`, its sole GoConvey v1.6.4 edge, zero
loaded packages, runtime unreachability, and no new advisory or independent
disqualifier. Enterprise Certificate Proxy v0.2.1 retains only its separate
exception for its sole Viper v1.15.0 edge under the same zero-load and
unreachable guards. GAX v2.7.0 retains only its two exceptions while zero GAX
packages load. Revalidate these guards before recording a choice and stop for
a fresh owning decision if any expired.

Do not silently select v1.5.3, add a root edge, change mvn-pom-mutator or
GoConvey, raise the Go floor, transfer another dependency's exception, or make
an implementation commit. Option 1 authorizes only a decision record and no
dependency edit. Options 2 and 3 authorize only the newly bounded evaluation,
not its implementation in the decision-recording turn.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, latest Google UUID dependency implementation identity, reciprocal
archive history, P7/P8 state, guarded exceptions, and
`./codex-dev-start.sh --check`. Read this archive, the answered Gorilla
evaluation, rolling handover, roadmap, `go.mod`, `go.sum`, and referenced
release, quality, compatibility, evidence, and lifecycle contracts. Earlier
outcomes are final.

# Three Moves

First, obtain an explicit option 1, 2, or 3 selection without inferring consent.
Second, revalidate the unchanged graph/load/guard facts and record exactly that
decision without implementing it. Third, rewrite the roadmap and rolling
handover, answer this archive, and prepare one reciprocal NEXT archive for the
authorized bounded follow-up.

# Automatic Handoff

Run the launcher contract and applicable no-change gates, create the required
local `docs: prepare next agent session` commit, and stop. Do not launch the
successor, push, merge, publish, release, stash, revert, bypass cleanup, or
remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
