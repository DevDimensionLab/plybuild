# Agent Session: Decide Hashicorp Go Retryablehttp Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-19T233605+0200-decide-hashicorp-go-retryablehttp-product-direction`
Created: `2026-09-19T23:36:05+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `44be42e6416b561758e17861aeb5672d14f068c5d5ebb4cd46ac3a021d822a35`
Previous: [2026-09-19T205837+0200-evaluate-hashicorp-go-retryablehttp-dependency.md](2026-09-19T205837+0200-evaluate-hashicorp-go-retryablehttp-dependency.md)
Next: [2026-09-19T235421+0200-evaluate-hashicorp-go-rootcerts-dependency.md](2026-09-19T235421+0200-evaluate-hashicorp-go-rootcerts-dependency.md)
Outcome: Recorded authorized option 1 for exact inherited, unloaded go-retryablehttp v0.5.3 after every guard passed; the target-specific exception accepts only the completed advisory and behavior findings, dependency metadata stayed unchanged, and one bounded go-rootcerts evaluation was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one bounded product decision for exact-path
`github.com/hashicorp/go-retryablehttp`. The completed independent evaluation
found no exact-path stable release that preserves the Go 1.18 floor and passes
every qualification contract. Obtain or apply one explicit authorized choice
from the options below, record its exact accepted findings and expiry guards,
and stop. Do not repeat the audit, silently accept the vulnerability, implement
a dependency change, evaluate another dependency group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 is blocked only on this go-retryablehttp product
decision after exact Go 1.26.7, every accepted dependency move through Google
UUID v1.4.0, qualified go-cleanhttp v0.5.2, and all retained-module decisions
through go-multierror v1.1.0. Every earlier outcome and lifecycle ancestor is
final. Preserve every target-specific qualification or exception; none
transfers to go-retryablehttp. P8 remains queued.

The evaluation left `go.mod` and `go.sum` unchanged. Exact selected
`github.com/hashicorp/go-retryablehttp v0.5.3` remains inherited through the
sole exact request from `github.com/armon/go-metrics v0.3.10`, unloaded, and
runtime-unreachable. This physical MVS selection is not qualification or risk
acceptance. Do not add a direct edge merely to alter MVS, change go-metrics,
Viper, mvn-pom-mutator, another guarded parent, the Go floor, or an unrelated
module.

The go-multierror v1.1.0, go-msgpack v0.5.3, go-immutable-radix v1.3.1,
go-hclog v1.2.0, Errwrap v1.0.0, and every other recorded exception or
qualification remain separate under their exact selection, incoming-edge,
zero-load, runtime-unreachable, and no-new-finding guards. Stop for the owning
decision if any guard changes. In particular, choosing v0.7.7 would move
guarded go-hclog v1.2.0 to v1.6.3 and fatih/color v1.15.0 to v1.16.0; neither
move can be accepted implicitly here.

# Measurements At Start

The evaluation recording began from clean handoff HEAD
`24b1f33671362309526328170e67e68aaf753a3a`, parent
`a6dba3eee16d9c3c622fdc68c075aa45595c27cc`, tree
`d00fe77d08624cb227aa61641b403ad42c9d4deb`. The latest dependency
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
1,067 sum lines, and the recorded 432-line unapplied tidy projection. The
authoritative no-change quality result is 27/27 Q0-Q2 PASS at L2, scorecard
SHA-256
`de13154181ae319fd80a29a98df78535762d70c2c01438f9e3cb735c2f624e81`.

# Completed Evaluation

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the exact
module to Hashicorp's public, active, unarchived, non-fork MPL-2.0 repository.
The proxy exposes 23 stable releases from v0.5.0 through latest v0.7.8. There
are no prereleases, retractions, module deprecations, redirects, alternate
paths, `/v2` line, or GitHub Release objects. Current main is 17 commits beyond
v0.7.8 and unreleased. Proxy and tagged-Git source bytes agree for all
releases; relevant tags are unsigned.

Selected v0.5.3 is verified commit
`357460732517ec3b57c05c51443296bdd6df1874`. Last declared-floor-compatible
v0.7.5 is commit `4165cf8897205a879a06b20d1ed0a2a76fbb6a17`.
First secure v0.7.7 contains fix commit
`a99f07beb3c5faaa0a283617e6eb6bcf25f5049a`; latest v0.7.8 is commit
`e1f5485fe84728709b857cb89e17088894c301d6`.

V0.5.0-v0.6.2 have no `go` directive, v0.6.3-v0.7.5 declare Go 1.13,
v0.7.6/v0.7.7 declare Go 1.19, and v0.7.8 declares Go 1.23. V0.7.5 is
therefore the highest candidate preserving the declared Go 1.18 floor, while
v0.7.7 is the first security-fixed release and is floor-ineligible. Actual Go
1.18 compilation of later directives does not convert their unsupported
declared floor into qualification.

All floor-compatible releases pass isolated upstream tests under contained Go
1.18.10. Under exact Go 1.26.7, v0.6.3-v0.7.5 fail the TLS retry-policy test
because modern Go wraps the x509 verification error. Selected v0.5.3 and
v0.6.2 pass native/repeated/race tests under both SDKs but retain a test-only
vet failure from `t.Fatalf` in non-test goroutines. V0.7.7 passes complete
exact-Go native/repeated/race/vet and five cross-build targets.

Every later release is export-incompatible with selected v0.5.3 because
`Client.Logger` changes from `Logger` to `interface{}`. The independent
behavior fixture SHA-256 is
`cf4ceef2094699c773cc5a44afb90ed992d4ac3ae222a3ee27ff79d95bea1a1c`.
Selected v0.5.3 fails permanent classification for x509 unknown-authority,
unsupported-scheme, and redirect-limit errors and loses `errors.Is` identity
for the final transport error. Body replay, response ownership, cancellation,
backoff/status, aliasing, nil/panic, allocation, concurrency, global-client,
and resource boundaries are fully characterized. V0.7.5 passes the fixture
but fails the modern-Go TLS branch.

Fresh 1,402-record primary vulnerability data has SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Reviewed GO-2024-2947 and
GHSA-v6v8-xj6m-xwqh affect v0.5.0-v0.7.6: Basic-auth credentials embedded in
a URL can be exposed in log and error paths. The independent fixture SHA-256
is `2b15b5812d21c3f91f9279105923423f6fbb311568c56ec8d3651e91fd952bcb`
and reproduces the credential disclosure through v0.7.6; v0.7.7/v0.7.8 pass.
Consumer module/package/symbol/test-symbol scans corroborate the affected and
fixed ranges.

Selected v0.5.3 exists only through go-metrics v0.3.10. Its why result and
repository import search are negative; production and complete-test loads
contain zero target packages. Disposable projections preserve zero target
load. V0.5.3/v0.6.2/v0.7.5 move no unrelated selection. V0.7.6/v0.7.7 also
upgrade guarded go-hclog and fatih/color; v0.7.8 makes those moves and raises
the main `go` directive to 1.23. No projection was applied.

Every earlier guard remains intact. Exact-Go project tests, pinned lint,
API/CLI, empty-HOME, cross-build, host/snapshot/Docker, all 17 script/meta,
80/80 mutation, and all 15 audit controls pass. The 547-entry complete-
evidence manifest SHA-256 is
`12f36165ac7c0ddb4a69a5e432c3b12ea57c977a9e6d7c7a338a359d36b4faa2`.

# Required Product Decision

Choose exactly one. Option 1 is recommended because the affected code is
unloaded and runtime-unreachable and it is the only option that preserves the
current Go floor, parents, API, and every existing guarded selection. It is a
new target-specific risk acceptance, not a transfer of any earlier exception:

1. Retain exact selected, inherited, unloaded v0.5.3 without metadata changes
   under a go-retryablehttp-specific exception. Accept only GO-2024-2947/
   GHSA-v6v8-xj6m-xwqh, the independently reproduced credential disclosure,
   the x509/unsupported-scheme/redirect retry classification and final-error-
   identity defects, the test-only vet defect, and the fully characterized
   API, body, cancellation, status, backoff, hook, nil/panic, mutation,
   aliasing, allocation, concurrency, global-state, resource, MVS, and related
   findings. Guard exact v0.5.3, the sole go-metrics v0.3.10 request, zero
   target imports/load, runtime unreachability, and no new advisory or
   independent defect.
2. Explicitly authorize a separately implemented exact v0.7.7 selection and a
   Go-floor/API migration. This is the first security-fixed release and passes
   the exact-Go behavior gates, but declares Go 1.19, changes exported API,
   upgrades guarded go-hclog v1.2.0 to v1.6.3 and fatih/color v1.15.0 to
   v1.16.0, and therefore requires fresh, separately sequenced Go-floor,
   go-hclog, and compatibility decisions before merge plus a later dependency-
   only implementation and complete changed-selection gate. This session must
   not implement or silently decide those scopes.
3. Keep P7 blocked and authorize a separately scoped parent-chain removal,
   upgrade, fork, patch, replacement, or architecture study involving
   go-metrics/historical Viper/mvn-pom-mutator. No such work is authorized in
   this decision session.

Do not infer acceptance from physical selection or zero reachability. If no
explicit authorized choice is available, report the blocker and preserve this
decision as the sole next boundary without changing source or dependency
metadata.

# Role And Boundaries

This is a decision session, not a renewed audit or implementation. Reuse the
completed evaluation. Revalidate only exact selection, the sole incoming
request, negative why/imports, zero target and guarded package loads, project
hashes, runtime unreachability, every earlier guard, and current advisory
state. Stop for the owning decision if any guard changed. Do not broaden a
choice into direct use, parent changes, patching, forking, replacement, a Go-
floor change, unrelated-module changes, another dependency group, or P8.

# Required Reading

At start read this archive, its answered go-retryablehttp evaluation, the
answered go-multierror and go-hclog decisions, rolling handover, roadmap,
`go.mod`, and `go.sum`. Verify branch, clean ordinary and ignored status,
handoff ancestry and changed-file set, exact toolchain identity, reciprocal
archive history, all guarded selections/requests/load results, and
`./codex-dev-start.sh --check`. Earlier outcomes are final.

# Three Moves

First, revalidate the narrow guards and stop for the owning decision if any
changed. Second, obtain or apply exactly one explicit choice above without
implementation, record the accepted findings and precise expiry conditions,
and preserve every unrelated selection. Third, update the roadmap and rolling
handover, answer this archive, prepare exactly one reciprocal NEXT mission for
the authorized bounded continuation, and commit the handoff. Do not execute
the successor.

# Automatic Handoff

After a coherent explicit decision, make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Recorded authorized option 1 exactly: retain selected
`github.com/hashicorp/go-retryablehttp v0.5.3` as an inherited, unloaded
module without changing `go.mod` or `go.sum`. This is a new
go-retryablehttp-specific, non-transferable exception. It accepts only
GO-2024-2947/GHSA-v6v8-xj6m-xwqh and the independently reproduced Basic-auth
URL credential disclosure through client and package-level log/error paths;
the x509 unknown-authority, unsupported-scheme, and redirect-limit permanent-
retry classification defects; loss of final transport-error `errors.Is`
identity; the test-only vet defect from `t.Fatalf` calls in non-test
goroutines; and the completed exported-API, body replay and response
ownership, cancellation/deadline, status, backoff, hook, nil/panic, mutation,
aliasing, allocation, concurrency, global-state, resource, MVS,
vulnerability, and related qualification findings. It accepts no new or
independently discovered defect.

Guard-only revalidation began from clean ordinary worktree state on branch
`codex/upgrade-quality` at handoff HEAD
`a4c24f29e25eb21cc229cde674ffdacdecc78f38`, parent
`24b1f33671362309526328170e67e68aaf753a3a`, tree
`ccca83e0256147f7e51d14eb07a2fb2472e2aabe`. That handoff changed exactly
the launcher, answered go-retryablehttp evaluation archive, this then-NEXT
decision archive, rolling handover, and roadmap. Its reciprocal archive
links, Google UUID implementation ancestry, and launcher lifecycle check
passed. The only ignored worktree entry observed around the checks was the
generated `target/` directory; it was removed before the final clean handoff.

Exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
was used first in `PATH` with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, no ambient `GOFLAGS`, `LC_ALL=C`, `LANG=C`, and
`umask 022`. Go-retryablehttp remains exact v0.5.3 through the sole exact
request from `github.com/armon/go-metrics v0.3.10`. Its `go mod why -m`
result remains negative, repository Go imports remain zero, and the 355-entry
production and 429-entry complete-test loads contain zero target packages.
The complete-test load retains 197 module-backed entries across 41 loaded
modules, so the target remains runtime-unreachable.

All 15 earlier guarded modules retain their exact selections and recorded
incoming requests. All 16 target-plus-earlier `go mod why -m` results remain
negative, repository Go imports remain zero, and production and complete-test
loads contain zero guarded packages. The project remains 234 modules, 3,599
graph edges, 1,067 `go.sum` lines, and the recorded 432-line unapplied tidy
projection. `go.mod` and `go.sum` retain SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
No parent, Go-floor, exported API, or unrelated selection changed.

Fresh primary vulnerability data remains byte-identical at 1,402 records,
SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV queries remain empty
for 14 earlier guards; Gorilla WebSocket v1.4.2 retains only its recorded
GO-2026-6278/GHSA-w67g-5rqw-f597 result. Exact go-retryablehttp v0.5.3
retains exactly GO-2024-2947 and GHSA-v6v8-xj6m-xwqh. All records remain
active. No new advisory or independently observed defect appeared.

The exception remains valid only while exact go-retryablehttp v0.5.3 and its
sole go-metrics v0.3.10 request remain unchanged, no direct main-module edge
or repository import is added, production and complete-test target loads
remain zero, the module remains runtime-unreachable, every earlier guarded
selection/request/load/advisory condition remains intact, and no new target
advisory or independent defect appears. Direct import or loading, runtime
reachability, a target version or incoming-request change, a new direct root,
an earlier owning-guard change, or a new advisory or independent defect
expires the exception and requires the owning fresh dependency and product
decision before merge.

No dependency implementation or metadata commit was created. Exact-Go module
verification, build, count-one tests, race tests, vet, and the launcher
lifecycle check pass; no changed-selection scorecard applies, and accepted
quality remains 27/27 Q0-Q2 PASS at L2. The sole reciprocal successor is the
bounded P7 evaluation of selected exact-path
`github.com/hashicorp/go-rootcerts v1.0.2`; it was prepared but not executed.
