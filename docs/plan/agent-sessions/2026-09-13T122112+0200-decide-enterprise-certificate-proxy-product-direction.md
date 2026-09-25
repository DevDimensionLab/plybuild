# Agent Session: Decide Enterprise Certificate Proxy Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-13T122112+0200-decide-enterprise-certificate-proxy-product-direction`
Created: `2026-09-13T12:21:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `05da58d1df222345ae11a1fa73e1851fd756e262e4cda5829be7bd8e96882933`
Previous: [2026-09-13T041350+0200-evaluate-enterprise-certificate-proxy-dependency.md](2026-09-13T041350+0200-evaluate-enterprise-certificate-proxy-dependency.md)
Next: [2026-09-13T140912+0200-evaluate-gopherjs-dependency.md](2026-09-13T140912+0200-evaluate-gopherjs-dependency.md)
Outcome: Recorded the user's bounded option 1 decision, retained unloaded Enterprise Certificate Proxy v0.2.1 without metadata changes, revalidated its exact Viper edge and zero-load guard plus the separate GAX zero-load guard, and prepared the next bounded P7 group.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user's 2026-09-13 explicit selection of
option 1 for `github.com/googleapis/enterprise-certificate-proxy`. Retain the
exact selected, unloaded v0.2.1 without dependency metadata changes under the
bounded floor, release, behavior, and safety exceptions below. Do not perform
dependency implementation, change Viper or GAX, audit another dependency, or
begin P8 in this decision-recording move.

# Authorized Roadmap

P2A-P6 and all earlier P7 dependency outcomes are final. P7 is paused at this
product stop; P8 and every other dependency group remain queued. The existing
GAX v2.7.0 floor and behavior exceptions remain valid only for that exact
selection while zero GAX packages load and do not transfer to this target or a
different GAX release. Do not reopen earlier groups.

# Authorized Product Decision

On 2026-09-13 the user explicitly selected option 1 with the recommended
bounds: retain exact `enterprise-certificate-proxy v0.2.1` as an inherited,
unloaded selection. Keep the main module's Go 1.18 floor and accept v0.2.1's
Go 1.19 declaration only for this exact selected version and exact inherited
Viper v1.15.0 edge. Do not change `go.mod` or `go.sum`, add a direct target
edge, patch or replace the target, downgrade or independently audit Viper or
GAX, raise the Go floor, or manufacture a dependency commit.

The user separately accepts the known v0.2.1 release, behavior, and safety
findings already recorded by the answered evaluation: public-preview release
classification; the proxy zip's lost signer-fixture executable bit and raw
test failure; configuration and file-lifetime boundaries; unbounded signer
startup/RPC waits; incomplete failed-start cleanup; error-string-dependent,
non-idempotent Close behavior; nil/panic and caller-aliasing boundaries; the
process-global logger side effect; and the raw pointer/length, out-of-bounds,
short-buffer, conflated-error, and uncaught-panic C ABI boundaries. This does
not accept any new or independently discovered defect.

These exceptions are valid only while the complete project load contains zero
Enterprise Certificate Proxy packages and the selected version and sole
incoming Viper edge remain unchanged. Revalidate and record that zero-load
invariant in the final decision and handover. The exceptions expire if a
target package becomes loaded or directly imported, its version or incoming
edge changes, a new vulnerability/advisory or independent disqualifier
appears, or the target becomes reachable runtime behavior. The owning
checkpoint must then stop for a fresh dependency and product decision before
merge. Existing GAX exceptions remain separate and retain their own zero-load
guard. Do not stop or ask for this same Enterprise Certificate Proxy decision
again if these invariants still hold.

# Measurements At Start

Selected v0.2.1 declares Go 1.19 and had no inherited exception before the
authorized decision above. V0.2.0 is the highest exact-path tag with a complete
minimal source/test closure preserving Go 1.18, but GitHub classifies it as a
prerelease/public preview. Its proxy zip cannot pass its own tests without
restoring a lost executable bit. Its client
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
technical audit. The user has explicitly selected and bounded option 1; do not
ask for that decision again or broaden it into Viper, GAX, unrelated-module,
Go-floor, removal, direct-use, or future-version authorization.

# Required Reading

Read the answered Enterprise Certificate Proxy archive, rolling handover,
roadmap, and the answered GAX archive. They contain the complete repository and
release identity, floor closure, exported API, process, crypto, C ABI, exact
MVS/loading, vulnerability, quality, tool-control, and existing-exception
evidence. Verify the feature branch, lifecycle links, clean ordinary and
ignored status, and `./codex-dev-start.sh --check` before recording a choice.

# Three Moves

First, verify from the clean current HEAD that v0.2.1 and its sole Viper edge
remain selected, dependency metadata is untouched, `go mod why -m` remains
negative, and zero target and GAX packages load. Reuse the completed audit;
do not repeat the full dependency qualification.

Second, record the exact product decision, accepted known findings, invariant,
expiration triggers, and future stop condition in the roadmap and rolling
handover. Retain v0.2.1 without a dependency or implementation commit.

Third, answer this decision archive and prepare one bounded next P7 mission
under the existing queue. Do not execute that successor mission in this turn.

# Automatic Handoff

Follow the repository lifecycle contract after an explicit decision or a
continued product stop. Do not perform the resulting implementation in this
decision-recording move, launch a successor, push, merge, publish, release,
stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Recorded the user's 2026-09-13 option 1 selection and retained exact-path
`github.com/googleapis/enterprise-certificate-proxy v0.2.1` as an inherited,
unloaded module without changing `go.mod` or `go.sum` and without manufacturing
a dependency or implementation commit. The main module keeps its Go 1.18
floor. The bounded floor exception accepts v0.2.1's Go 1.19 declaration only
for this exact selection and its unchanged sole incoming Viper v1.15.0 edge.

The decision separately accepts only the findings already established in the
answered dependency evaluation: the public-preview release classification;
the proxy zip's lost signer-fixture executable bit and raw test failure;
configuration and file-lifetime boundaries; unbounded signer startup and RPC
waits; incomplete failed-start cleanup; error-string-dependent,
non-idempotent `Close`; nil/panic and caller-aliasing boundaries; the process-
global logger side effect; and the raw pointer/length, out-of-bounds,
short-buffer, conflated-error, and uncaught-panic C ABI boundaries. No new or
independently discovered defect is accepted.

The complete invariant was revalidated from clean decision HEAD `f671bc0` with
exact Go 1.26.7, `GOENV=off`, `GOWORK=off`, `GOTOOLCHAIN=local`, and no ambient
`GOFLAGS`. Project selection remains Enterprise Certificate Proxy v0.2.1,
Viper v1.15.0, and GAX v2.7.0. The only incoming target edge is
Viper v1.15.0 -> Enterprise Certificate Proxy v0.2.1. `go mod why -m` remains
negative. The 429-entry complete test load still contains 197 module-backed
packages across 41 loaded modules, with exactly zero Enterprise Certificate
Proxy packages and exactly zero GAX packages.

`go.mod` and `go.sum` are byte-identical to clean HEAD. Their SHA-256 values are
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`;
`go.sum` remains 1,067 lines. No direct target edge, patch, replacement, Viper
or GAX change, Go-floor change, or dependency metadata change was made.

The Enterprise Certificate Proxy exceptions remain valid only while zero
target packages load, the target stays unreachable at runtime, selected
v0.2.1 and its sole incoming Viper v1.15.0 edge remain unchanged, and no new
vulnerability/advisory or independent disqualifier appears. Loading or directly
importing a target package, runtime reachability, a target-version or incoming-
edge change, or a new vulnerability/advisory or independent disqualifier
expires the exceptions. The owning checkpoint must then stop for a fresh
dependency and product decision before merge. The existing GAX v2.7.0
exceptions remain separate under their own zero-load guard. Do not request the
same Enterprise Certificate Proxy decision again while all invariants hold.
