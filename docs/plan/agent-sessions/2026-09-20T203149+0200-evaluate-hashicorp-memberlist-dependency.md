# Agent Session: Evaluate Hashicorp Memberlist Dependency

Status: NEXT
Session ID: `2026-09-20T203149+0200-evaluate-hashicorp-memberlist-dependency`
Created: `2026-09-20T20:31:49+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e75f031049356b10eba461bd03e6bc1e21c183155940c4e60e8b2147d89f26b5`
Previous: [2026-09-20T201536+0200-decide-hashicorp-mdns-product-direction.md](2026-09-20T201536+0200-decide-hashicorp-mdns-product-direction.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/memberlist v0.3.0` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure,
package behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Defensive Retry Scope

This is an authorized defensive software-supply-chain review of the user's
local repository. A first attempt from clean HEAD
`a1dee3a81a71e0e6656de09fef603b94cc135bab` was automatically stopped by a
cybersecurity policy check after collecting public advisory and release facts;
it made no tracked worktree change. Resume this same bounded session under the
narrower rules below rather than restarting the broad investigation.

Use only public release metadata, the public vendor advisory, static source and
remediation-diff inspection, the project module graph, and ordinary existing
build or test commands with non-adversarial inputs in disposable local trees.
Do not create, reproduce, simulate, operationalize, or optimize an exploit or
proof of concept. Do not craft malformed packets, generate custom attack
traffic, scan or contact hosts, exercise production or non-public systems,
inspect credentials or key material, or test a security boundary. Do not
recreate the advisory's failure condition. If completing the decision would
require any such activity, stop and record that limitation instead.

The interrupted attempt retained these provisional facts for narrow
confirmation: public vendor advisory `HCSEC-2026-18` / `CVE-2026-14362` reports
all memberlist releases through v0.5.4 as affected; v0.6.0 is the first fixed
release; and v0.6.0 declares Go 1.25. Confirm only the affected/fixed range,
the candidate Go directive and source/test closure, and the static remediation
identity from public primary material. Do not investigate exploitability. If
confirmed, no fixed memberlist release preserves the project's Go 1.18 floor,
so leave dependency metadata unchanged and prepare the bounded product
decision. Treat missing temporary output from the interrupted attempt as a
reason to recheck these few facts, not to repeat the broad evaluation.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded mdns v1.0.4. Every earlier outcome and lifecycle ancestor is final.
Evaluate only Hashicorp memberlist in this session; do not reopen or combine
another dependency group. P8 remains queued.

The authorized 2026-09-20 mdns option 1 decision retains exact selected,
inherited, unloaded `github.com/hashicorp/mdns v1.0.4` without product source
or dependency metadata changes. Its target-specific, non-transferable
exception accepts only the completed exact-Go linker/test, Go-1.18 IPv6,
protocol, packet/channel, timeout/cancellation, lifecycle/resource, nil/panic,
mutation/aliasing, concurrency, environment, allocation, API, MVS,
vulnerability, and related recorded findings. It remains valid only while
exact v1.0.4, the Serf-v0.9.6-to-v1.0.4 and historical
Serf-v0.8.2-to-v1.0.0 requests, the mdns-v1.0.0-to-go.net-v0.0.1 request, no
direct root or repository import, zero target loads, runtime unreachability,
every earlier guard, and the no-new-advisory-or-defect condition remain exact.
Any change requires the owning fresh decision.

The golang-lru v0.5.4, go.net v0.0.1, go-uuid v1.0.1, go-syslog v1.0.0,
go-sockaddr v1.0.0, go-rootcerts v1.0.2, go-retryablehttp v0.5.3,
go-multierror v1.1.0, go-msgpack v0.5.3, go-immutable-radix v1.3.1,
go-hclog v1.2.0, Errwrap v1.0.0, qualified go-cleanhttp v0.5.2, and every
other recorded exception or qualification remain separate under their exact
selection, incoming-request, zero-load, runtime-unreachable, and no-new-
finding guards. Revalidate those guards and stop for the owning decision if
any expires. No earlier exception transfers to memberlist. Do not change a
guarded parent, the Go floor, or an unrelated module.

Selected `github.com/hashicorp/memberlist v0.3.0` is inherited through the
exact request from `github.com/hashicorp/serf v0.9.6`; historical Serf v0.8.2
requests memberlist v0.1.3. The current repository import search and
production and complete-test loads contain zero target packages, its why
result is negative, and there is no direct main-module root. These queue
observations and the physical MVS selection are not proof of repository
identity, release qualification, ancestry, floor, behavior, vulnerability
state, or suitability. Resolve them independently and do not add a direct
edge merely to alter MVS.

# Measurements At Start

The mdns decision recording began from clean ordinary and ignored state at
HEAD `9edb38c51e0067672e90224d8665bae580f92a05`, parent
`b641b55ec01f18a79ecf57b8e5e866310dd44562`, tree
`63a43a0d26649532e760560a9e114fea47feb8d1`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status,
reciprocal archive chain, latest Google UUID implementation ancestry, and
`./codex-dev-start.sh --check` rather than assuming them.

Guard-only decision revalidation under exact Go 1.26.7 preserved all 22
earlier guarded selections and 162 incoming graph edges plus exact mdns
v1.0.4 and its two incoming Serf requests; the mdns-v1.0.0-to-go.net request
also remains exact. All guarded why results remain negative, repository
imports are zero, and production and complete-test loads contain zero guarded
packages. The project remains 234 modules, 3,599 graph
edges, 355 production entries, 429 complete-test entries, 197 module-backed
entries across 41 loaded modules, 1,067 sum lines, and the recorded 432-line
unapplied tidy projection. Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV remains empty for
mdns and every guarded target except the recorded Gorilla
GO-2026-6278/GHSA-w67g-5rqw-f597 and go-retryablehttp
GO-2024-2947/GHSA-v6v8-xj6m-xwqh pairs. Mdns's exact GitHub global and
repository advisory queries remain empty. No new guarded advisory or
independent defect appeared.

Use exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, use `LC_ALL=C LANG=C`, and
run with `umask 022`. Recreate contained Go 1.18.10 and pinned tools beneath
`$CODEX_SESSION_SCRATCH_ROOT` as required. Preserve the known apidiff archive
reproducibility discrepancy, Python 3.14 Docker timestamp control, managed
bare-`mktemp` restriction, Make-variable inheritance boundary, and nested
launcher signal-retention timing race; none is memberlist evidence.

# Role And Boundaries

From fresh external archives and caches, resolve proxy, sumdb, `go-import`,
Git, and forge evidence for the exact module path: tags, releases, branches,
signatures, commits, parents, trees, times, ancestry, repository status,
licenses, retractions, deprecations, redirects, forks, alternate paths, major
module lines, and every serious exact-path stable candidate. Do not silently
promote a redirect, fork, alternate path, different major module path,
prerelease, non-versioning tag, unreleased branch head, or floor-ineligible
release.

Prove the complete minimal production and test closure under exact Go 1.26.7
and contained Go 1.18.10. Inspect imported source and test dependencies rather
than treating the module directive alone as floor proof. Separate isolated
source-time resolution from the project's selected MVS graph and from every
already-final guarded dependency decision.

Inspect public source, packages, exported API, platform/build-tag branches, and
ordinary runtime lifecycle only as needed for the dependency decision. Reuse
the interrupted attempt's conclusions where recorded. Run existing upstream
tests and non-adversarial compatibility checks without modifying them, plus
source verification, package listing, race, vet, and meaningful cross-builds
under both SDKs where still necessary. Do not add security-focused fixtures or
custom network, packet, cryptographic, credential, malformed-input, fuzzing, or
failure-condition tests. Do not expand ordinary functional review into a
protocol-security assessment. Classify any ordinary build or test failure
precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, floor, architecture, or unrelated-
module change needs a fresh bounded decision rather than silent
implementation. For the advisory, record only public primary affected/fixed
version facts, the fixed release's Go floor, the project's zero-load and
runtime-unreachable state, and a static remediation-diff identity. Do not
perform exploitability analysis, dynamic reproduction, or custom reachability
testing.

# Required Reading

At start read this archive, the answered mdns decision and evaluation, the
answered golang-lru decision and evaluation, the answered go.net decision and
evaluation, rolling handover, roadmap, `go.mod`, and `go.sum`. Verify the
recorded handoff identity, changed-file set, clean ordinary and ignored
status, exact toolchain identity, reciprocal archive history, every guarded
selection/request/load result, and `./codex-dev-start.sh --check`. Earlier
outcomes are final.

# Three Moves

First, narrowly confirm the retained public advisory range, first fixed
release, Go-floor incompatibility, repository identity, project zero-load, and
exact MVS facts without changing the worktree or reproducing the advisory.
Second, if those provisional facts hold, leave metadata unchanged and prepare
one bounded product decision comparing at least: raising the Go floor and
selecting the fixed line; retaining exact unloaded v0.3.0 under an explicit
target-specific risk decision; and removing or changing the owning parent
edge in a separately measured study. Do not choose or implement an option in
this session. Only if primary evidence disproves the provisional facts may a
Go-1.18-compatible fixed exact-path stable candidate proceed through ordinary
qualification. Third, update the roadmap and rolling handover, answer this
archive, prepare exactly one reciprocal NEXT decision mission, and commit the
handoff. Do not execute the successor.

# Automatic Handoff

After one coherent bounded outcome, make any separate dependency-only commit
first if a qualified selection was implemented, then make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
