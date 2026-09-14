# Agent Session: Decide Hashicorp Consul API Product Direction

Status: NEXT
Session ID: `2026-09-14T193201+0200-decide-hashicorp-consul-api-product-direction`
Created: `2026-09-14T19:32:01+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `bf5cbb42142253008c86ffc929bc8ea8a0d4f6ff3ca4186e18197ff18bee0930`
Previous: [2026-09-14T180109+0200-evaluate-hashicorp-consul-api-dependency.md](2026-09-14T180109+0200-evaluate-hashicorp-consul-api-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by recording the user's 2026-09-14 explicit selection of
option 1 for exact-path `github.com/hashicorp/consul/api`. Retain exact
selected, inherited, unloaded v1.18.0 without dependency metadata changes
under the bounded release-closure, native-test, TLS-fixture, response-metadata,
closure-vulnerability, nil/panic, mutation, and qualification exceptions
below. Do not repeat the technical audit, implement a dependency change,
evaluate another dependency group, or begin P8 in this decision-recording
session.

# Authorized Roadmap

P2A-P6 and every earlier P7 outcome are final. P7 remains active only for this
Consul API product decision; P8 remains queued. The user has explicitly chosen
to retain exact v1.18.0 solely through Viper v1.15.0 while zero Consul API
packages load and the module remains runtime-unreachable. This session may
record that one choice and prepare one bounded follow-up, but may not implement
the choice or combine another dependency group.

The user's Gateway option 1 decision and every earlier guarded decision remain
final, target-specific, and non-transferable. Revalidate their exact selection,
incoming-edge, zero-load, runtime-unreachable, and no-new-finding guards. Stop
for the owning decision if any guard expires. Do not change either Gateway
parent, mvn-pom-mutator, GoConvey, Viper, or any unrelated module.

# Measurements At Start

Fresh proxy, sumdb, `go-import`, Git, and forge evidence resolves canonical
public active unarchived non-fork repository
`https://github.com/hashicorp/consul.git`, with module subdirectory `api` and
MPL-2.0 at the selected source. The v1 proxy lists 96 entries: 84 stable exact-
semver releases and 12 prerelease/suffix entries. Selected v1.18.0 is
lightweight tag `api/v1.18.0`, commit
`13836d5ca84c71b35f201da06dd75f5ad6699c36`, parent
`18dffc51de5587f8fba5f188a636d1864bf2b98a`, repository tree
`1960701737d9ae23b0dc7d24a5087a4aa0541f83`, API subtree
`b088d6696716b2cb8b61c824299dca34cab203f4`, dated
2022-11-30T18:59:52Z with a GitHub-verified commit signature. Its proxy zip
SHA-256 is
`0dc6cfca8c71b05b3ba859726378d2ee611c15304fc85c2c030e3366ee068062`;
proxy and Git API source are byte-identical and sumdb agrees.

All 83 retrievable stable v1 go.mod files were inspected. V1.18.0 is the
highest stable exact-path release whose declaration does not exceed Go 1.18;
v1.18.1 is the first Go 1.19 release. Latest v1.34.5 declares Go 1.26.7 and has
five incompatible exported-API changes from v1.18.0. The distinct
`github.com/hashicorp/consul/api/v2 v2.0.0` path declares Go 1.26 and is outside
this exact-path decision.

The selected complete source/test closure imports 56 modules under both exact
Go 1.26.7 and Go 1.18.10; no imported module declares above Go 1.17. Source
verification, production build, vet, and four production cross-builds pass
under both SDKs. The watch package's count-one, two count-ten repeats, race,
and test cross-compilation pass under both SDKs. The selected release does not
qualify because:

- its root API test binary cannot link on Darwin or Linux under Go 1.26.7 due
  to old x/net's `syscall.recvmsg` reference; adjacent v1.17.0 reproduces it;
- its `consulent` test branch does not compile under either SDK because it
  lacks `defaultNamespace` and `defaultPartition` definitions;
- its proxy archive is not standalone: go.mod replaces the SDK with `../sdk`
  and tests need 14 certificate fixtures from `../test/client_certs`;
- four native TLS tests fail because the tagged certificates expired on
  2023-11-01; all other 213 tests pass repeated and race runs with Consul
  v1.14.2 and a scratch-only sandbox freeport shim;
- all response metadata parsing errors are discarded, so malformed Consul
  index/cache/hash/query-backend headers are silently accepted with defaults.

An independent consumer fixture covers configuration, request and option
encoding, auth/token/TLS inputs, cancellation/deadlines, error identity,
no-retry behavior, cleanup, malformed inputs, nil/panic boundaries, and
concurrency. It passes count-one, two count-ten repeats, race, and vet under
both SDKs. Pinned API comparisons show no incompatibility from v1.17.0 to
v1.18.0 or v1.18.0 to v1.18.1.

Selected v1.18.0 has exactly one selected-version incoming edge, from Viper
v1.15.0. `go mod why -m` is negative, repository imports are zero, zero target
packages occur in the 429-entry complete project load, and it is runtime-
unreachable. Exact disposable v1.18.0 `go get` manufactures eight indirect
roots and 26 sum lines but tidy removes them and converges byte-identically to
the base tidy projection. V1.17.0 downgrades Viper; v1.18.1 upgrades x/net and
x/text and exceeds the floor; v1.34.5 raises the main Go declaration and moves
many unrelated selections. No projection was applied.

Fresh primary vulnerability data remains 1,398 records at SHA-256
`cde9b02ce42b801cbd683fb39c85b61acfbb3ea7cf896f95d4999003bed4fdbf`
and Last-Modified 2026-09-10T16:28:28Z, with no Consul API module record.
Exact OSV queries for v1.17.0, v1.18.0, v1.18.1, and v1.34.5 are empty. The
selected production closure has one inherited module-only ID, GO-2026-5024 in
old x/sys, with no vulnerable loaded package or called symbol. The test closure
adds inherited GO-2022-0603 in yaml.v3 and loads that test package, but reaches
no vulnerable symbol. Normal project scans have zero target occurrence.

The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and a
432-line tidy projection. `go.mod`/`go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact Go 1.26.7 applicable project gates pass; Go 1.18.10 retains only the two
accepted Darwin shell wording assertions. Accepted quality remains 27/27 Q0-
Q2 PASS at L2. Buildx is now present, but Docker acceptance reaches the known
Python 3.14 timestamp-control incompatibility after building and testing the
image. One completed-handoff lifecycle run passes all 62 controls; three later
final-text runs pass outer controls 1-50 and reproduce only the known nested
signal-log timing race at control 51. This remains target-independent.

The 9,780-entry external evidence manifest and decision-summary SHA-256 values
are `966263870f7529dcdd4713eba204dc1a8a00e17901a5f702557467b3c79ac46e`
and `5957399af9cf4564b9864e4b6dea8c84e8912db3846e446c953345dadbddb08a`.

# Authorized Product Decision

On 2026-09-14 the user explicitly selected option 1 with the recommended
bounds: retain exact selected `github.com/hashicorp/consul/api v1.18.0` as an
inherited, unloaded selection without changing `go.mod` or `go.sum`. Accept
only the documented non-standalone release closure, Go 1.26 Unix native-test
link incompatibility, broken `consulent` test branch, expired TLS fixtures,
silently discarded response-metadata errors, inherited closure-only
vulnerability findings, documented nil/panic and mutation boundaries, and
related completed qualification findings. This does not accept a new or
independently discovered defect.

The exception is target-specific and non-transferable. It is valid only while
exact v1.18.0 and its sole Viper v1.15.0 selected-version incoming edge remain
unchanged, the complete project load contains zero Consul API packages, the
module remains runtime-unreachable, and no new advisory or independent
disqualifier appears. Revalidate and record those guards. Direct import or
loading, runtime reachability, a target version or incoming-edge change, or a
new advisory or independent defect expires the exception and requires a fresh
Consul API dependency and product decision before merge.

Do not add a direct target edge, select another v1 release, move to `/api/v2`,
change Viper, raise the Go floor, authorize parent modernization/removal or
replacement architecture work, move unrelated selections, or manufacture a
dependency implementation commit. The Gateway decision and every earlier
exception remain separate. Do not stop or ask for this same Consul API decision
again while all guards hold.

# Role And Boundaries

This is a decision-recording session, not a renewed audit or implementation.
The user has explicitly selected and bounded option 1. Reuse the completed
evidence; do not ask for the decision again or broaden it into direct use,
another version, `/api/v2`, Viper or floor modernization, parent removal,
replacement architecture, or unrelated-module authorization.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status,
reciprocal archive history, latest Google UUID dependency implementation,
P7/P8 state, every guarded invariant, and `./codex-dev-start.sh --check`. Read
this archive, the answered Consul API evaluation, rolling handover, roadmap,
`go.mod`, `go.sum`, and referenced contracts. Do not repeat the audit.

# Three Moves

First, revalidate only the exact selection, incoming edge, negative why result,
zero target and guarded package loads, module hashes, runtime unreachability,
and current advisory state; reuse the completed audit and do not broaden it.
Second, record the exact option 1 exception, accepted findings, guards, expiry
triggers, and no-change result in the roadmap and rolling handover. Third,
answer this archive and prepare exactly one reciprocal NEXT mission for the
next bounded P7 group without executing it.

# Automatic Handoff

After recording an explicit choice, run the applicable no-change lifecycle
gates and create the required local `docs: prepare next agent session` commit.
Do not launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
