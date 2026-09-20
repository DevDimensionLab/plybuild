# Agent Session: Decide Hashicorp Go Syslog Product Direction

Status: NEXT
Session ID: `2026-09-20T145140+0200-decide-hashicorp-go-syslog-product-direction`
Created: `2026-09-20T14:51:40+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `69b4575a8633e48d7f6b549df29f1bc9b9c5c577eae6ae4c9639970fb9951496`
Previous: [2026-09-20T135311+0200-evaluate-hashicorp-go-syslog-dependency.md](2026-09-20T135311+0200-evaluate-hashicorp-go-syslog-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user-authorized automatic controller's
2026-09-20 explicit selection of option 1 for exact-path
`github.com/hashicorp/go-syslog`. Retain exact selected, inherited, unloaded
v1.0.0 without source or dependency metadata changes under the target-specific
exception below. Revalidate only the decision guards and record exactly this
choice. Do not repeat the audit, silently accept a new defect, implement a
dependency change, evaluate another dependency group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 is recording the bounded go-syslog option 1 decision
after exact Go 1.26.7, every accepted dependency move through Google UUID
v1.4.0, qualified go-cleanhttp v0.5.2, and all target-specific retained-module
decisions through exact inherited, unloaded go-sockaddr v1.0.0. Every earlier
outcome and lifecycle ancestor is final. The user-authorized controller has
explicitly chosen to retain exact inherited, unloaded go-syslog v1.0.0 through
both exact Serf requests. Preserve every existing qualification or exception;
none transfers to go-syslog. This session may record the choice and prepare one
bounded follow-up; it may not implement a dependency change or combine another
group. P8 remains queued.

The evaluation left `go.mod` and `go.sum` unchanged. Exact selected
`github.com/hashicorp/go-syslog v1.0.0` remains inherited, unloaded, and
runtime-unreachable. This physical MVS selection is not qualification or risk
acceptance. Do not add a direct edge merely to alter MVS, change Serf, another
guarded parent, the Go floor, or an unrelated module.

The go-sockaddr v1.0.0, go-rootcerts v1.0.2, go-retryablehttp v0.5.3,
go-multierror v1.1.0, go-msgpack v0.5.3, go-immutable-radix v1.3.1, go-hclog
v1.2.0, Errwrap v1.0.0, qualified go-cleanhttp v0.5.2, and every other
recorded exception or qualification remain separate under their exact
selection, incoming-request, zero-load, runtime-unreachable, and no-new-
finding guards. Stop for the owning decision if any guard changes.

# Measurements At Start

The go-syslog evaluation recording began from clean handoff HEAD
`2ee06f1db429101b908d764de8c37a9cd2989229`, parent
`5c957ce042e2c36589e0ec5f00fca5261ce76352`, tree
`4c96957b3a4391d2379732962243f12f50928d95`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status,
reciprocal archive chain, and `./codex-dev-start.sh --check` rather than
assuming them.

Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The project remains 234 modules, 3,599 graph edges, 355 production entries,
429 complete-test entries, 197 module-backed entries across 41 loaded modules,
1,067 sum lines, and the recorded 432-line unapplied tidy projection. Because
the evaluation changed no production source or dependency metadata, no
changed-selection scorecard applies; accepted quality remains 27/27 Q0-Q2
PASS at L2.

# Completed Evaluation

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the exact
module to Hashicorp's public, active, unarchived, non-fork MIT repository. The
proxy exposes only v1.0.0, also `@latest`; there are no other stable or
prerelease versions, releases, retractions, module deprecations, redirects,
alternate exact paths, or `/v2` line. V1.0.0 is a lightweight tag at verified
commit `8d1874e3e8d1862b74e0536851e218c4571066a5`, parent
`326bf4a7f709d263f964a6a96558676b103f3534`, tree
`a81f242fd29ca869a5c2fb666ff9e4dd7a0af570`. Proxy and Git regular files
agree and sumdb verifies the release.

Current master is verified commit
`40240a543e78c0ccbdb437b470650ea7acaf3da8`, 48 commits beyond v1.0.0,
declares Go 1.23.0, and is an unreleased, floor-ineligible branch head. It has
no exported API change, fixes only the deadline-error defect below, and
retains the constructor and priority defects. No branch, fork, redirect,
alternate path, or unreleased commit is an authorized substitute.

V1.0.0 has no `go` directive or external dependency; the complete minimal
closure is the target plus standard library and passes native source,
listing, repeat, race, vet, and applicable cross-build checks under exact Go
1.26.7 and contained Go 1.18.10. The release has one package, four production
files, no native tests, commands, examples, benchmarks, fuzz targets,
testdata, generated files, cgo, or embeds. Its API exposes `Priority`, eight
severity constants, `Syslogger`, `NewLogger`, and `DialLogger`; there is no
CLI.

Independent fixtures under both SDKs characterize UDP/TCP format and
priority, local Unix-daemon behavior, framing, dial/write/close ownership,
idempotent close and reconnect, error identity, malformed input, nil panics,
buffer mutation/non-aliasing, allocation, concurrency, global state,
resources, and environment/network interaction. Embedded newlines are sent
verbatim and may form multiple records. Passing fixture SHA-256 is
`6d9eb0b8fbfefc1556630519301609b6d343cb37ec2b3b5c46d6d6c2378d46d0`;
external cross-platform API fixture SHA-256 is
`b50b035d93a2244fbe575ee2d0dbb34ec1d70b9ca501bde3817f517781f8b58a`.

The sole stable release fails three independent qualification contracts:

1. `NewLogger` and `DialLogger` are absent for AIX and `js/wasm` under both
   SDKs and for `wasip1/wasm` under Go 1.26.7. The package itself nominally
   builds there, but an external consumer cannot compile, contradicting the
   library's cross-compilation/runtime-error purpose.
2. `writeString` ignores a `SetWriteDeadline` error, writes anyway, returns
   nil, loses error identity, and defeats the documented blocking-avoidance
   guarantee. The controlled blocker reproduces under both SDKs. Unreleased
   master fixes this one defect.
3. `DialLogger` accepts out-of-range severity `Priority(8)` with `LOCAL0` and
   silently remaps it to `LOG_EMERG` rather than rejecting it. The blocker
   reproduces under both SDKs and remains on current master.

Selected v1.0.0 exists in MVS through exactly Serf v0.8.2 and v0.9.6 requests.
Its why result is negative, repository imports are zero, and production and
complete-test loads contain zero target packages. A redundant direct root
adds only one edge and one checksum line, changes no selection or load, and is
removed by tidy; base and candidate tidy projections converge byte-identically.
All 18 earlier guarded modules retain exact requests, negative why/imports,
zero load, runtime unreachability, and their recorded advisory states.

Fresh primary vulnerability data contains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. It has no target entry; exact OSV,
repository, and exact global GitHub advisory queries are empty. Exact-Go
isolated scans have zero target findings. Old-SDK isolated findings belong
only to its standard library, not go-syslog. Base and redundant-root project
scans are identical at 30 module, 22 package, and 20 symbol/test-symbol IDs,
with zero target package or trace occurrence.

Exact Go 1.26.7 no-change project verification, build, count-one, two
count-ten repeats, race, vet, pinned lint, API/CLI compatibility, empty-HOME,
four cross-builds, full preflight, every script/meta population, all eight
mutation meta-stages with 80/80 kills, host/snapshot/Docker acceptance, and
all 15 audit controls pass. Known apidiff-cache, bare-`mktemp`, Python 3.14
timestamp, Make-variable, and launcher timing boundaries remain independent
of go-syslog. No changed-selection scorecard applies.

# Authorized Product Decision

On 2026-09-20 the user-authorized automatic controller explicitly selected
option 1 with the recommended bounds: retain exact selected
`github.com/hashicorp/go-syslog v1.0.0` as an inherited, unloaded selection
without changing product source, `go.mod`, or `go.sum`. Accept only the
completed constructor-availability, deadline-error, priority-validation,
newline-framing, API/behavior, platform, MVS, vulnerability, and related
recorded findings. This does not accept a new or independently discovered
defect.

The exception is go-syslog-specific and non-transferable. It is valid only
while exact v1.0.0 and both exact requests from Serf v0.8.2 and v0.9.6 remain
unchanged, no direct main-module root or repository import is added, production
and complete-test target loads remain zero, the module remains runtime-
unreachable, every earlier guard remains intact, and no new advisory or
independent defect appears. Revalidate and record these guards. Direct import
or loading, runtime reachability, a target version or incoming-request change,
a new direct root, an earlier owning-guard change, or a new advisory or
independent defect expires the exception and requires the owning fresh
dependency and product decision before merge.

Do not change product source or dependency metadata, add a direct go-syslog
edge, select an unreleased branch, change or remove either Serf parent, patch or
fork source, authorize a wrapper or replacement architecture, raise the Go
floor, move unrelated selections, alter another guarded dependency, or
manufacture a dependency implementation commit. Do not stop or ask for this
same go-syslog decision again while all guards hold. Preserve every earlier
qualified result and target-specific decision.

# Role And Boundaries

This is a decision-recording session, not a renewed audit or implementation.
The user-authorized automatic controller has explicitly selected and bounded
option 1. Reuse the completed evaluation; do not ask for the decision again or
broaden it into direct use, unlisted MVS movement, parent changes, patching,
forking, wrapping, replacement, a Go-floor change, unrelated-module
authorization, another dependency group, or P8. Revalidate only exact
selection and incoming requests, negative why/imports, zero target and guarded
package loads, project hashes, runtime unreachability, every earlier guard, and
current advisory state. Stop for the owning decision if any guard changed.

# Required Reading

At start read this archive, its answered go-syslog evaluation, the answered
go-sockaddr, go-rootcerts, go-retryablehttp, go-multierror, and go-hclog
decisions, rolling handover, roadmap, `go.mod`, and `go.sum`. Verify branch,
clean ordinary and ignored status, handoff ancestry and changed-file set,
exact toolchain identity, reciprocal archive history, all guarded selections/
requests/load results, and `./codex-dev-start.sh --check`. Earlier outcomes
are final.

# Three Moves

First, revalidate the narrow guards and stop for the owning decision if any
changed. Second, record exactly the controller's user-authorized option 1
exception, accepted findings, limits, expiry triggers, and no-change decision
result in the roadmap and rolling handover; preserve every unrelated selection.
Third, answer this archive, prepare exactly one reciprocal NEXT mission for the
next bounded P7 group, and commit the handoff. Do not execute the successor.

# Automatic Handoff

Only after guard revalidation, make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
