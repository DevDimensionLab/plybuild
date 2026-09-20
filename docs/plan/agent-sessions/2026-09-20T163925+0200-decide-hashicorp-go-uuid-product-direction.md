# Agent Session: Decide Hashicorp Go UUID Product Direction

Status: NEXT
Session ID: `2026-09-20T163925+0200-decide-hashicorp-go-uuid-product-direction`
Created: `2026-09-20T16:39:25+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c2e32e320fd145aba954a7e9dde7676aea6a0156741398c2cf9470b77a35f078`
Previous: [2026-09-20T152123+0200-evaluate-hashicorp-go-uuid-dependency.md](2026-09-20T152123+0200-evaluate-hashicorp-go-uuid-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user-authorized automatic controller's
2026-09-20 explicit selection of option 1 for exact-path
`github.com/hashicorp/go-uuid`. Retain exact selected, inherited, unloaded
v1.0.1 without source or dependency metadata changes under the target-specific
exception below. Revalidate only the decision guards and record exactly this
choice. Do not repeat the audit, silently accept a new defect, implement a
dependency or source change, combine another dependency group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 is recording the bounded go-uuid option 1 decision
after exact Go 1.26.7, every accepted dependency move through Google UUID
v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded go-syslog v1.0.0. Every earlier outcome and lifecycle ancestor is
final. The user-authorized controller has explicitly chosen to retain exact
inherited, unloaded go-uuid v1.0.1 through all eight recorded requests. No
prior exception transfers to it. Preserve every guarded selection. This
session may record the choice and prepare one bounded follow-up; it may not
change a parent, the Go floor, product source, dependency metadata, or an
unrelated module. P8 remains queued.

Selected `github.com/hashicorp/go-uuid v1.0.1` exists through six exact
v1.0.1 requests from Consul API v1.1.0/v1.12.0, Consul SDK v0.1.1/v0.8.0,
and Serf v0.8.2/v0.9.6. Go-immutable-radix v1.0.0/v1.3.1 additionally
request v1.0.0. There is no direct main-module root or repository import; its
why result is negative, production and complete-test loads contain zero target
packages, and it is runtime-unreachable. Those facts bound risk but do not
constitute qualification or acceptance.

# Completed Evaluation

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves Hashicorp's
public active unarchived non-fork MPL-2.0 repository. The exact path has four
stable lightweight tags, v1.0.0 through proxy-latest v1.0.3, with no
prerelease, GitHub Release object, retraction, module deprecation, redirect,
alternate path, or `/v2` line. All proxy regular files match Git and sumdb
verifies each release. Current master
`f405b577e09f44ce7e9b484af2911859f9dab5a6` is 28 commits beyond v1.0.3,
declares Go 1.18, and is unreleased.

Each stable release has one package, no external module dependency, and a
complete minimal production/test closure of target plus standard library.
Every release passes verification, listing, native tests, independent repeats,
race, vet, and production/test cross-builds under exact Go 1.26.7 and contained
Go 1.18.10. V1.0.0/v1.0.1 expose four functions; v1.0.2/v1.0.3 compatibly add
two reader-injection functions. The package intentionally formats random bytes
as UUID text rather than setting RFC version or variant bits.

Every stable release loses underlying randomness-error identity by formatting
with `%v`. Under Go 1.18, v1.0.0/v1.0.1 return a non-wrapping error; under Go
1.26 their `crypto/rand.Read` path terminates the process when a failing global
reader is injected. V1.0.2/v1.0.3 add caller-owned reader APIs but retain the
same `%v` identity loss. All stable releases also panic for negative byte
counts despite returning an error. Unreleased master changes `%v` to `%w` and
passes error-identity fixtures under both SDKs, but retains the negative-size
panic. No stable candidate therefore qualifies.

Disposable direct v1.0.1 adds only one edge and content checksum. Direct
v1.0.2/v1.0.3 alter only the target selection and add one edge, but tidy
removes the manufactured root and restores selected v1.0.1. Direct v1.0.0
breaks the guarded graph by removing mvn-pom-mutator v0.2.3 and makes project
loads fail. No projection was applied. The project remains 234 modules, 3,599
edges, 355 production entries, 429 complete-test entries, 197 module-backed
entries across 41 loaded modules, 1,067 sum lines, and the recorded 432-line
tidy projection. `go.mod`/`go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Fresh primary vulnerability data has no target record. Exact OSV, repository
advisory, and exact global GitHub advisory queries are empty for all four
releases. Exact-Go isolated scans are empty; old-Go reachable findings belong
only to its standard library. Base and v1.0.1/v1.0.2/v1.0.3 project scans are
identical, with zero target assignment or trace. Every earlier guarded
selection/request/load/advisory condition remains exact, and no new guarded
finding appeared.

The exact-Go canonical project gate passes verify, build, count-one, two
count-ten repeats, race, vet, pinned lint, API/CLI, empty-HOME, four
cross-builds, full preflight, all 15 audit controls, all eight mutation
meta-stages, 80/80 kills, and host/snapshot/Docker acceptance. The fresh
clean-tree scorecard is 27/27 Q0-Q2 PASS at L2. No dependency implementation
commit exists.

# Authorized Product Decision

On 2026-09-20 the user-authorized automatic controller explicitly selected
option 1 with the recommended bounds: retain exact selected
`github.com/hashicorp/go-uuid v1.0.1` as an inherited, unloaded selection
without changing product source, `go.mod`, or `go.sum`. Accept only the
completed randomness-error identity, Go 1.26 fatal injected-reader,
negative-size panic, non-RFC generation, API/behavior, MVS, vulnerability,
and related recorded findings. This does not accept a new or independently
discovered defect.

The exception is go-uuid-specific and non-transferable. It is valid only
while exact v1.0.1, all six exact v1.0.1 requests from Consul API
v1.1.0/v1.12.0, Consul SDK v0.1.1/v0.8.0, and Serf v0.8.2/v0.9.6, and both
exact v1.0.0 requests from go-immutable-radix v1.0.0/v1.3.1 remain unchanged,
no direct main-module root or repository import is added, production and
complete-test target loads remain zero, the module remains runtime-
unreachable, every earlier guard remains intact, and no new advisory or
independent defect appears. Revalidate and record these guards. Direct import
or loading, runtime reachability, a target version or incoming-request change,
a new direct root, an earlier owning-guard change, or a new advisory or
independent defect expires the exception and requires the owning fresh
dependency and product decision before merge.

Do not change product source or dependency metadata, add a direct go-uuid
edge, select v1.0.3 or unreleased master, change or remove any Consul API,
Consul SDK, Serf, or go-immutable-radix parent, patch or fork source, authorize
a wrapper or replacement architecture, raise the Go floor, move unrelated
selections, alter another guarded dependency, or manufacture a dependency
implementation commit. Do not stop or ask for this same go-uuid decision
again while all guards hold. Preserve every earlier qualified result and
target-specific decision. The patch/fork/replacement and parent-chain removal
alternatives were not authorized.

# Required Reading And Guards

This is a decision-recording session, not a renewed audit or implementation.
The user-authorized automatic controller has explicitly selected and bounded
option 1. Reuse the completed evaluation; do not ask for the decision again.
At start read this archive, its answered go-uuid evaluation, the answered
go-syslog decision and evaluation, the answered go-sockaddr, go-rootcerts,
go-retryablehttp, go-multierror, and go-hclog decisions, rolling handover,
roadmap, `go.mod`, and `go.sum`. Verify branch, clean ordinary and ignored
status, handoff ancestry and exact changed-file set, reciprocal archive chain,
exact Go 1.26.7 identity, all 20 guarded selections/requests/why/import/load
conditions, module hashes, advisory guards, and
`./codex-dev-start.sh --check`. Stop if any guard differs.

# Three Moves

First, revalidate the narrow guards and stop for the owning decision if any
changed. Second, record exactly the controller's user-authorized option 1
exception, accepted findings, limits, expiry triggers, and no-change result in
the roadmap and rolling handover; preserve every unrelated selection. Third,
answer this archive, prepare exactly one reciprocal NEXT mission for the next
bounded P7 group, and commit the handoff. Do not execute the successor.

# Automatic Handoff

Only after guard revalidation, make only the required local
`docs: prepare next agent session` commit. Do not create a dependency commit,
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, remove the worktree, combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
