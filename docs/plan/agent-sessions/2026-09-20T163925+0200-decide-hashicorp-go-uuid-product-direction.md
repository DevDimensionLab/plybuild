# Agent Session: Decide Hashicorp Go UUID Product Direction

Status: NEXT
Session ID: `2026-09-20T163925+0200-decide-hashicorp-go-uuid-product-direction`
Created: `2026-09-20T16:39:25+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f6a7356f3ef40cba9eb8ea1549d5915b41c9fa94ae68c54517496c087f06e9eb`
Previous: [2026-09-20T152123+0200-evaluate-hashicorp-go-uuid-dependency.md](2026-09-20T152123+0200-evaluate-hashicorp-go-uuid-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by obtaining and recording one bounded product decision for
exact-path `github.com/hashicorp/go-uuid`. No exact-path stable release passes
every applicable behavior contract, so the preceding evaluation made no
dependency change. Choose among the bounded options below; do not silently
accept the findings, implement a dependency or source change, combine another
dependency group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded go-syslog v1.0.0. Every earlier outcome and lifecycle ancestor is
final. The evaluated go-uuid selection is exact inherited v1.0.1; no prior
exception transfers to it. Preserve every guarded selection and do not change
a parent, the Go floor, product source, or an unrelated module without the
corresponding fresh authorization.

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

# Decision Required

Choose exactly one direction:

1. **Retain exact inherited, unloaded v1.0.1 under a target-specific
   exception (recommended).** Make no product-source or dependency-metadata
   change. Accept only the completed randomness-error identity, Go 1.26 fatal
   injected-reader, negative-size panic, non-RFC generation, API/behavior,
   MVS, vulnerability, and related recorded findings. The exception is
   go-uuid-specific and non-transferable. It remains valid only while exact
   v1.0.1, all six v1.0.1 requests and both v1.0.0 requests remain unchanged,
   no direct root or repository import is added, both target loads remain
   zero, runtime unreachability holds, every earlier guard remains intact, and
   no new advisory or independent defect appears. Any change requires the
   owning fresh dependency and product decision.
2. **Authorize a separately evaluated Go-1.18-compatible patch, fork, or
   replacement.** Scope a new bounded design/evaluation to preserve the
   required API, wrap reader errors, define non-panicking negative-size
   behavior, and prove provenance, maintenance, MVS, license, vulnerability,
   compatibility, and full quality effects before implementation. Do not
   implement it in this decision session.
3. **Authorize a separately evaluated parent-chain removal.** Scope a fresh
   graph/product study across the exact Consul API, Consul SDK, Serf, and
   go-immutable-radix requesters. This may affect earlier guarded decisions and
   cannot be performed or assumed here.

Selecting v1.0.3 alone is not a qualifying fourth option: it retains both
behavior blockers, introduces a direct root that tidy removes, and provides no
loaded project consumer benefit. An unreleased master, redirect, fork, or
alternate path is likewise not an authorized substitute.

# Required Reading And Guards

At start read this archive, its answered go-uuid evaluation, the answered
go-syslog decision and evaluation, the answered go-sockaddr, go-rootcerts,
go-retryablehttp, go-multierror, and go-hclog decisions, rolling handover,
roadmap, `go.mod`, and `go.sum`. Verify branch, clean ordinary and ignored
status, handoff ancestry and exact changed-file set, reciprocal archive chain,
exact Go 1.26.7 identity, all 20 guarded selections/requests/why/import/load
conditions, module hashes, advisory guards, and
`./codex-dev-start.sh --check`. Stop if any guard differs.

# Three Moves

First, obtain one explicit product choice and do nothing else while it is
pending. Second, revalidate the narrow guards and record exactly the authorized
choice, limits, expiry triggers, and no-change result; a patch/replacement or
parent-removal choice authorizes only its next bounded evaluation, not an
implementation here. Third, update the roadmap and rolling handover, answer
this archive, prepare exactly one reciprocal NEXT mission for the authorized
result, and commit the handoff. Do not execute the successor.

# Automatic Handoff

After one coherent bounded outcome, make only the required local
`docs: prepare next agent session` commit. Do not create a dependency commit,
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, remove the worktree, combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
