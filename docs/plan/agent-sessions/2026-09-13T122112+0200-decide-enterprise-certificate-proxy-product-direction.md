# Agent Session: Decide Enterprise Certificate Proxy Product Direction

Status: NEXT
Session ID: `2026-09-13T122112+0200-decide-enterprise-certificate-proxy-product-direction`
Created: `2026-09-13T12:21:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `9aebb4096b6877ce800bf98f3ac85418a6469ca215d6e3392d1243cd020fb642`
Previous: [2026-09-13T041350+0200-evaluate-enterprise-certificate-proxy-dependency.md](2026-09-13T041350+0200-evaluate-enterprise-certificate-proxy-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by obtaining and recording the user's fresh bounded product
decision for `github.com/googleapis/enterprise-certificate-proxy`. The
independent evaluation found no qualified, technically acceptable, in-scope
exact-path selection under the current Go 1.18 and behavior contracts. Do not
implement, retain by implication, grant an exception, change Viper or GAX,
audit another dependency, or begin P8 until the user explicitly chooses and
bounds a direction.

# Authorized Roadmap

P2A-P6 and all earlier P7 dependency outcomes are final. P7 is paused at this
product stop; P8 and every other dependency group remain queued. The existing
GAX v2.7.0 floor and behavior exceptions remain valid only for that exact
selection while zero GAX packages load and do not transfer to this target or a
different GAX release. Do not reopen earlier groups.

# Measurements At Start

Selected v0.2.1 declares Go 1.19 and has no inherited exception. V0.2.0 is the
highest exact-path tag with a complete minimal source/test closure preserving
Go 1.18, but GitHub classifies it as a prerelease/public preview. Its proxy zip
cannot pass its own tests without restoring a lost executable bit. Its client
has unbounded signer startup/RPC waits, incomplete failed-start cleanup,
error-string-dependent and non-idempotent Close behavior, typed-nil panics, and
caller aliasing. Its C-shared ABI trusts raw pointer/length pairs, permits
out-of-bounds access, conflates failures as zero, and allows invalid inputs to
escape as uncaught Go panics. Selected v0.2.1 shares those boundaries and also
can permanently discard the process-global logger.

Viper v1.15.0 supplies the only project edge to selected v0.2.1. `go mod why
-m` is negative, and zero target packages load. An exact v0.2.0 request
necessarily downgrades Viper to v1.14.0, GAX to v2.6.0, and changes 20 module
selections total. Exact v0.1.0 is API-incompatible, downgrades Viper to v1.13.0,
and disappears after tidy.

No dependency metadata changed. The unchanged project remains 234 modules,
3,599 graph edges, 429 complete-test entries, 197 module-backed packages, 41
loaded modules, 1,067 sum lines, and a 432-line tidy projection. Target and GAX
each remain unloaded. Accepted quality remains 27/27 Q0-Q2 PASS at L2. Fresh
primary vulnerability data contains no target record, and direct serious-
candidate scans are zero.

The 562-entry selected-evidence manifest SHA-256 is
`59868d9547ece628b950faf9200377108a0aab55fbbfab731d5f0b5fbcd10033`;
decision-summary SHA-256 is
`276efebfe75500789ce6501630e8ecdb0dbbf8241ae20c74f57bccfbc348180c`.

# Role And Boundaries

This is a decision-recording move, not a dependency implementation or renewed
technical audit. Present the measured consequences accurately, distinguish a
floor exception from behavior and safety acceptance, and require the user to
bound any Viper, GAX, unrelated-module, Go-floor, or removal scope explicitly.
Do not convert a recommendation or absence of objection into authorization.

# Required Reading

Read the answered Enterprise Certificate Proxy archive, rolling handover,
roadmap, and the answered GAX archive. They contain the complete repository and
release identity, floor closure, exported API, process, crypto, C ABI, exact
MVS/loading, vulnerability, quality, tool-control, and existing-exception
evidence. Verify the feature branch, lifecycle links, clean ordinary and
ignored status, and `./codex-dev-start.sh --check` before recording a choice.

# Three Moves

First, ask the user to choose and bound exactly one direction:

1. Retain unloaded v0.2.1 by explicitly defining new Enterprise Certificate
   Proxy Go-1.19 and identified behavior/safety exceptions, including their
   invariant and stop condition.
2. Authorize a broader Viper/MVS dependency group and separately decide
   whether v0.2.0's public-preview, packaging, process-lifecycle, and C-ABI
   failures may be risk-accepted. This must address changed GAX and all other
   selections and does not itself authorize implementation without a new
   bounded audit plan.
3. Authorize a different Viper/removal/Go-floor strategy that eliminates the
   target or permits a later release, with explicit dependency, compatibility,
   repository, and exception boundaries.

Second, if the user does not choose, preserve the stop and ask for direction.
If the user does choose, record the exact scope, exceptions, invariants, and
stop conditions in the roadmap and rolling handover without implementation.

Third, answer this decision archive and prepare one bounded next mission that
faithfully implements or evaluates only the authorized choice.

# Automatic Handoff

Follow the repository lifecycle contract after an explicit decision or a
continued product stop. Do not perform the resulting implementation in this
decision-recording move, launch a successor, push, merge, publish, release,
stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
