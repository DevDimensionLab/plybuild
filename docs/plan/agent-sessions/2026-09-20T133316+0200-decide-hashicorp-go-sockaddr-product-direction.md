# Agent Session: Decide Hashicorp Go Sockaddr Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-20T133316+0200-decide-hashicorp-go-sockaddr-product-direction`
Created: `2026-09-20T13:33:16+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `52d3f0c26d8304a2a3e96c465867080c5d196f1889138ae0fa6939b24ae8fd85`
Previous: [2026-09-20T121350+0200-evaluate-hashicorp-go-sockaddr-dependency.md](2026-09-20T121350+0200-evaluate-hashicorp-go-sockaddr-dependency.md)
Next: [2026-09-20T135311+0200-evaluate-hashicorp-go-syslog-dependency.md](2026-09-20T135311+0200-evaluate-hashicorp-go-syslog-dependency.md)
Outcome: Recorded authorized option 1 for exact inherited, unloaded go-sockaddr v1.0.0 after every guard passed; the target-specific exception accepts only the completed closure, behavior, race, vet, boundary, MVS, vulnerability, and related findings, dependency metadata stayed unchanged, and one bounded go-syslog evaluation was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one bounded product decision for exact-path
`github.com/hashicorp/go-sockaddr`. The completed independent evaluation found
no exact-path stable release that preserves Go 1.18 and passes every
qualification contract. Obtain or apply one explicit authorized choice from
the options below, record its exact accepted findings and expiry guards, and
stop. Do not repeat the audit, silently accept a defect, implement a dependency
change, evaluate another dependency group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 is blocked only on this go-sockaddr product decision
after exact Go 1.26.7, every accepted dependency move through Google UUID
v1.4.0, qualified go-cleanhttp v0.5.2, and all target-specific retained-module
decisions through exact inherited, unloaded go-rootcerts v1.0.2. Every earlier
outcome and lifecycle ancestor is final. Preserve every existing qualification
or exception; none transfers to go-sockaddr. P8 remains queued.

The evaluation left `go.mod` and `go.sum` unchanged. Exact selected
`github.com/hashicorp/go-sockaddr v1.0.0` remains inherited, unloaded, and
runtime-unreachable. This physical MVS selection is not qualification or risk
acceptance. Do not add a direct edge merely to alter MVS, change memberlist,
another guarded parent, the Go floor, or an unrelated module.

The go-rootcerts v1.0.2, go-retryablehttp v0.5.3, go-multierror v1.1.0,
go-msgpack v0.5.3, go-immutable-radix v1.3.1, go-hclog v1.2.0, Errwrap
v1.0.0, qualified go-cleanhttp v0.5.2, and every other recorded exception or
qualification remain separate under their exact selection, incoming-edge,
zero-load, runtime-unreachable, and no-new-finding guards. Stop for the owning
decision if any guard changes.

# Measurements At Start

The evaluation recording began from clean handoff HEAD
`a9289a1684e101399810bbec6294b62e4a41df92`, parent
`c530f7d8e39aa0181899d7cb2893f6a16555ebef`, tree
`13addaa0a9dbc39f9a4eb22cbb6b778e072cebd0`. The latest dependency
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
module to Hashicorp's public, active, unarchived, non-fork MPL-2.0 repository.
The exact path exposes only v1.0.0-v1.0.7 and v1.0.7 is latest. There are no
prereleases, retractions, module deprecations, redirects, alternate exact
paths, or `/v2` line. All version tags are lightweight, not signed tag
objects. The sole GitHub Release is v1.0.7. Current master is 63 commits beyond
v1.0.7, declares Go 1.25, and is not a release candidate. Proxy and Git
regular-file bytes agree; module zip rules omit only the tracked nested vendor
tree.

V1.0.0-v1.0.3 omit a `go` directive and are the only serious floor-eligible
candidates. V1.0.4 declares Go 1.19.0 and v1.0.5-v1.0.7 declare Go 1.19, so
they are ineligible under the retained floor. V1.0.1-v1.0.3 have a complete
13-module graph with 11 actually imported external modules, all preserving Go
1.18. Selected v1.0.0 instead omits four command dependencies from its
released `go.mod`; read-only production/test listing, native tests, race, vet,
command build, whole-module API export, and all seven cross-target loads fail
under both SDKs. Byte verification alone passes and is not closure proof.

V1.0.1-v1.0.3 load and cross-build under exact Go 1.26.7 and contained Go
1.18.10, but native repeats and race fail on current Darwin because upstream
tests hard-code old loopback flags and sample interfaces, the managed route
command exits 71, default-interface tie sorting is unstable, and Go 1.26
changes one equal-network-size ordering. Go 1.18 additionally cannot parse the
exported `HashiCorpDefault2016` template. Vet rejects every release for an
unreachable panic in `ifaddrs.go`. V1.0.3 and later invert the Windows
PowerShell availability check.

Independent fixtures reproduce four disqualifying defects in v1.0.0-v1.0.3
and v1.0.7 under both SDKs: IPv6 `ContainsAddress` rejects an interior host;
no-match `GetInterfaceIP` and `GetInterfaceIPs` return empty output with nil
error; all six exported attribute inventories return mutable global slices
and race under concurrent caller mutation; and `Host()` exposes the package-
global IPv6 host-mask backing array, allowing caller mutation to corrupt later
addresses. Valid address, Unix, RFC, JSON, CLI, enumeration, copy/non-aliasing,
allocation, error, nil/panic, global-state, environment, route-command,
resource, and concurrency boundaries are otherwise fully characterized.

Pinned API comparison finds v1.0.0-v1.0.3 root APIs identical. V1.0.7 only
adds compatible `ErrNoInterface` and `ErrNoRoute` sentinels. V1.0.1-v1.0.3
command smokes pass under both SDKs; v1.0.0's command cannot build read-only.

Selected v1.0.0 exists through exact requests from memberlist v0.1.3 and
v0.3.0. Its why result and repository import search are negative; both project
loads contain zero target packages. A redundant v1.0.0 root adds one edge and
one checksum line without changing selection. Exact v1.0.1-v1.0.3 projections
also add `mitchellh/go-wordwrap v1.0.0` and upgrade unrelated
`ryanuber/columnize` to v2.1.0+incompatible. V1.0.7 moves 19 module lines,
including guarded Errwrap v1.0.0. No projection was applied.

Fresh 1,402-record primary vulnerability data has SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Exact OSV, GitHub advisory, and
govulncheck module/package/symbol/test-symbol results are empty for the target.
Base and disposable v1.0.3 project scans are identical with no target frame or
reachable trace. All 18 target-plus-earlier guards retain negative why,
zero imports, zero production/test loads, runtime unreachability, and their
recorded advisory states.

Exact Go 1.26.7 project verify, build, count-one, two count-ten repeats, race,
vet, pinned lint, API/CLI, empty-HOME, four cross-builds, canonical preflight,
all script/meta populations, 80/80 live mutation kills, host/snapshot/Docker
acceptance, and all 15 audit meta-controls pass. Contained Go 1.18 retains only
the two accepted Darwin `pkg/shell` wording failures; its 26 unaffected
packages and applicable gates pass. No changed-selection scorecard applies.

# Required Product Decision

Choose exactly one. Option 1 is recommended because the target is unloaded and
runtime-unreachable and it is the only choice that preserves the current Go
floor, parents, API, graph, and every earlier guard. It is a new target-specific
risk acceptance, not a transfer of any earlier exception:

1. Retain exact selected, inherited, unloaded v1.0.0 without metadata changes
   under a go-sockaddr-specific exception. Accept only its incomplete released
   module closure and command unbuildability; the reproduced IPv6 containment,
   no-match nil-error, mutable global-slice/race, and exposed IPv6 mask defects;
   the characterized upstream Darwin, Go 1.18 template, vet, sorting, route-
   command, API/CLI, nil/panic, mutation, aliasing, allocation, concurrency,
   global-state, resource, environment, MVS, vulnerability, and related
   findings. Guard exact v1.0.0, both memberlist requests and versions, no
   direct root, zero target imports/load, runtime unreachability, every earlier
   guard, and no new advisory or independent defect.
2. Authorize a later dependency-only implementation selecting exact v1.0.2
   under a separate explicit exception. Accept the same independent behavior,
   race, vet, test, and boundary findings while repairing only released module
   closure, and authorize exactly the consequent addition of
   `mitchellh/go-wordwrap v1.0.0` and upgrade of `ryanuber/columnize` to
   v2.1.0+incompatible. The later implementation must preserve the Go floor,
   all parents and guarded modules, prove the exact diff, and pass the complete
   changed-selection gate. This decision session must not implement it.
3. Keep P7 blocked. Do not accept v1.0.0 or authorize v1.0.2. Any patch, fork,
   replacement, wrapper, or memberlist/parent-chain pruning proposal requires
   a new explicitly scoped architecture or graph decision before source or
   metadata changes.

Do not infer acceptance from physical selection or zero reachability. If no
explicit authorized choice is available, report the blocker and preserve this
decision as the sole next boundary without changing source or dependency
metadata.

# Role And Boundaries

This is a decision session, not a renewed audit or implementation. Reuse the
completed evaluation. Revalidate only exact selection and incoming requests,
negative why/imports, zero target and guarded package loads, project hashes,
runtime unreachability, every earlier guard, and current advisory state. Stop
for the owning decision if any guard changed. Do not broaden a choice into
direct use, unlisted MVS movement, parent changes, patching, forking,
replacement, a Go-floor change, another dependency group, or P8.

# Required Reading

At start read this archive, its answered go-sockaddr evaluation, the answered
go-rootcerts, go-retryablehttp, go-multierror, and go-hclog decisions, rolling
handover, roadmap, `go.mod`, and `go.sum`. Verify branch, clean ordinary and
ignored status, handoff ancestry and changed-file set, exact toolchain
identity, reciprocal archive history, all guarded selections/requests/load
results, and `./codex-dev-start.sh --check`. Earlier outcomes are final.

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
`github.com/hashicorp/go-sockaddr v1.0.0` as an inherited, unloaded module
without changing `go.mod` or `go.sum`. This is a new go-sockaddr-specific,
non-transferable exception. It accepts only the incomplete released module
closure and resulting read-only command unbuildability; the independently
reproduced IPv6 `ContainsAddress` failure, no-match nil-error behavior,
mutable global attribute slices and their caller-mutation race, and exposed
package-global IPv6 host-mask backing array; the characterized upstream
Darwin interface, route-command, and sorting failures; the Go 1.18 exported-
template parse failure; the unreachable-panic vet finding; and the completed
API/CLI, nil/panic, mutation, aliasing, allocation, concurrency, global-state,
resource, environment, MVS, vulnerability, and related qualification
findings. It accepts no new or independently discovered defect.

Guard-only revalidation began from clean ordinary and ignored worktree state
on branch `codex/upgrade-quality` at handoff HEAD
`5c957ce042e2c36589e0ec5f00fca5261ce76352`, parent
`a9289a1684e101399810bbec6294b62e4a41df92`, tree
`654abb0f7bf3261c0b2224f99d61ec968af826f2`. That handoff changed exactly
the launcher, answered go-sockaddr evaluation archive, this then-NEXT
decision archive, rolling handover, and roadmap. Its reciprocal archive chain,
Google UUID implementation ancestry, and launcher check passed.

A freshly downloaded official Go 1.26.7 archive and binary retained SHA-256
values
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
The exact binary ran first in `PATH` with `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, no ambient `GOFLAGS`, `GOPROXY` and `GOSUMDB` pinned to
their public services, `LC_ALL=C`, `LANG=C`, and `umask 022`.

Go-sockaddr remains exact v1.0.0 through exactly the two recorded requests
from memberlist v0.1.3 and v0.3.0, with no direct main-module edge. Its
`go mod why -m` result remains negative, repository Go imports remain zero,
and the 355-entry production and 429-entry complete-test loads contain zero
target packages. All 18 target-plus-earlier guarded modules retain their exact
selections and recorded incoming requests. All 18 why results remain negative,
repository imports remain zero, and production and complete-test loads contain
zero guarded packages. The complete-test load retains 197 module-backed
entries across 41 loaded modules, so go-sockaddr and every earlier guarded
target remain runtime-unreachable.

The project remains 234 modules, 3,599 graph edges, and 1,067 `go.sum` lines.
`go.mod` and `go.sum` retain SHA-256 values
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
No parent, Go-floor, exported API, source, direct root, or unrelated selection
changed.

Fresh primary vulnerability data remains byte-identical at 1,402 records,
SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV responses remain
empty for go-sockaddr v1.0.0 and every guarded target except Gorilla
WebSocket v1.4.2 and go-retryablehttp v0.5.3. Gorilla retains only its
recorded GO-2026-6278/GHSA-w67g-5rqw-f597 results; go-retryablehttp retains
only its accepted GO-2024-2947/GHSA-v6v8-xj6m-xwqh results. Both the
go-sockaddr repository advisory feed and exact GitHub global-advisory query
remain empty. No new advisory or independently observed defect appeared.

The exception remains valid only while exact go-sockaddr v1.0.0 and both
memberlist v0.1.3/v0.3.0 requests remain unchanged, no direct main-module
edge or repository import is added, production and complete-test target loads
remain zero, the module remains runtime-unreachable, every earlier guarded
selection, request, load, runtime, and advisory condition remains intact, and
no new target advisory or independent defect appears. Direct import or
loading, runtime reachability, a target version or incoming-request change, a
new direct root, an earlier owning-guard change, or a new advisory or
independent defect expires the exception and requires the owning fresh
dependency and product decision before merge.

No dependency implementation or metadata commit was created. Exact Go 1.26.7
module verification and the applicable no-change project and launcher gates
pass; no changed-selection scorecard applies, and accepted quality remains
27/27 Q0-Q2 PASS at L2. The sole reciprocal successor is the bounded P7
evaluation of selected exact-path `github.com/hashicorp/go-syslog v1.0.0`;
it was prepared but not executed.
