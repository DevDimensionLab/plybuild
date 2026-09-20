# Agent Session: Decide Hashicorp Go Rootcerts Product Direction

Status: NEXT
Session ID: `2026-09-20T115926+0200-decide-hashicorp-go-rootcerts-product-direction`
Created: `2026-09-20T11:59:26+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `80a19eb740097e5933af06a87545f1d2ab359da2ce1572d15c136bb3ed42adba`
Previous: [2026-09-19T235421+0200-evaluate-hashicorp-go-rootcerts-dependency.md](2026-09-19T235421+0200-evaluate-hashicorp-go-rootcerts-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one bounded product decision for exact-path
`github.com/hashicorp/go-rootcerts`. The completed independent evaluation
found no exact-path stable release that preserves Go 1.18 and passes every
qualification contract. Obtain or apply one explicit authorized choice from
the options below, record its exact accepted findings and expiry guards, and
stop. Do not repeat the audit, silently accept a defect, implement a dependency
change, evaluate another dependency group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 is blocked only on this go-rootcerts product decision
after exact Go 1.26.7, every accepted dependency move through Google UUID
v1.4.0, qualified go-cleanhttp v0.5.2, and every target-specific retained-
module decision through exact inherited, unloaded go-retryablehttp v0.5.3.
Every earlier outcome and lifecycle ancestor is final. Preserve every existing
qualification or exception; none transfers to go-rootcerts. P8 remains queued.

The evaluation left `go.mod` and `go.sum` unchanged. Exact selected
`github.com/hashicorp/go-rootcerts v1.0.2` remains inherited, unloaded, and
runtime-unreachable. This physical MVS selection is not qualification or risk
acceptance. Do not add a direct edge merely to alter MVS, change Viper,
`sagikazarmark/crypt`, Consul API, mvn-pom-mutator, another guarded parent, the
Go floor, or an unrelated module.

The go-retryablehttp v0.5.3, go-multierror v1.1.0, go-msgpack v0.5.3,
go-immutable-radix v1.3.1, go-hclog v1.2.0, Errwrap v1.0.0, qualified
go-cleanhttp v0.5.2, and every other recorded exception or qualification remain
separate under their exact selection, incoming-edge, zero-load, runtime-
unreachable, and no-new-finding guards. Stop for the owning decision if any
guard changes.

# Measurements At Start

The evaluation recording began from clean handoff HEAD
`888412a62d768cb11af6d9e132600cc4a34f0679`, parent
`a4c24f29e25eb21cc229cde674ffdacdecc78f38`, tree
`1892ece8798c3913cc46465cf12865281db20339`. The latest dependency
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
the evaluation changed no production source or dependency metadata, no changed-
selection scorecard applies; accepted quality remains 27/27 Q0-Q2 PASS at L2.

# Completed Evaluation

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the exact
module to Hashicorp's public, active, unarchived, non-fork MPL-2.0 repository.
The proxy exposes only v1.0.0, v1.0.1, and latest v1.0.2. There are no
prereleases, retractions, module deprecations, redirects, alternate exact paths,
`/v2` line, or GitHub Release objects. Current master is 34 commits beyond
v1.0.2; its unreleased pseudo-version declares Go 1.23 and is not a candidate.
Proxy and Git regular-file bytes agree. The v1.0.0/v1.0.1 tags are lightweight
and unsigned; v1.0.2 is an unsigned annotated tag over GitHub-verified commit
`98fadc2a5ba2ad2a534a179b352ecdfd1f4259aa`.

V1.0.0 has no Go directive; v1.0.1/v1.0.2 declare Go 1.12. Their complete
minimal imported closure is target plus go-homedir v1.0.0 or v1.1.0 and
preserves Go 1.18. Tagged Git sources pass verification, native/repeated/race/
vet, and five production/test cross-builds under exact Go 1.26.7 and contained
Go 1.18.10. Standard proxy archives omit two tracked symlink fixtures, so all
three archived source suites fail `TestLoadCACertsFromDirWithSymlinks` because
the fixture directory is absent.

Pinned API comparison finds no v1.0.0-to-v1.0.1 change. V1.0.2 adds
`AppendCertificate` and `Config.CACertificate`; the `[]byte` field makes
`Config` non-comparable. There is no later stable candidate.

Independent fixtures reproduce two disqualifying behaviors in every stable
release under both SDKs. File, path, and ConfigureTLS errors stringify their
underlying I/O causes with `%s`, so `errors.Is(err, fs.ErrNotExist)` is false.
On Darwin, successful keychain commands that emit no PEM produce a non-nil
empty `*x509.CertPool` and nil error because `AppendCertsFromPEM`'s false
result is ignored. The upstream Darwin test checks only the error and misses
the empty trust pool. Current master retains both defects. Valid pool
construction, precedence, mutation/aliasing, malformed input, nil handling,
allocation, concurrency, resource, and global-state boundaries are otherwise
fully characterized.

Fresh 1,402-record primary vulnerability data has SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Exact OSV, GitHub advisory, and
govulncheck module/package/symbol/test-symbol scans are empty for all three
tags. Base and redundant-v1.0.2 project findings are identical with no target
frame or reachable trace.

Selected v1.0.2 exists through exact requests from Viper v1.15.0, historical
Viper v1.10.1, `sagikazarmark/crypt v0.4.0`, and historical Consul API v1.12.0;
historical Consul API v1.1.0 also requests v1.0.0. Its why result and repository
import search are negative, both project loads contain zero target packages,
and selected code is runtime-unreachable. A redundant exact v1.0.2 root adds
only one edge and one checksum line, then tidy converges to the base projection.
Exact v1.0.1 or v1.0.0 instead downgrades Viper to v1.8.1, removes
mvn-pom-mutator, disrupts earlier guarded selections, and makes project
packages unloadable. No projection was applied.

Every earlier guard remains exact. Exact-Go project build, native/repeated/
race/vet, pinned lint, API/CLI, empty-HOME, four cross-builds, full preflight,
all 17 script/meta pairs, all eight mutation meta-stages, 80/80 live mutation
kills, host/snapshot/Docker acceptance, and all 15 audit meta-controls pass.
Contained Go 1.18 retains only the two accepted Darwin `pkg/shell` wording
failures; its 26 unaffected packages and applicable gates pass.

# Required Product Decision

Choose exactly one. Option 1 is recommended because the target is unloaded and
runtime-unreachable and it is the only choice that preserves the current Go
floor, parents, API, graph, and every earlier guard. It is a new target-specific
risk acceptance, not a transfer of any earlier exception:

1. Retain exact selected, inherited, unloaded v1.0.2 without metadata changes
   under a go-rootcerts-specific exception. Accept only the loss of underlying
   filesystem error identity; the reproduced Darwin silent-empty system-root
   pool; the proxy-archive symlink-fixture test failure; the v1.0.2 Config
   comparability change; and the fully characterized certificate-pool, PEM,
   file/path, precedence, system-root, environment, nil/panic, mutation,
   aliasing, allocation, concurrency, global-state, resource, API, MVS,
   vulnerability, and related findings. Guard exact selected v1.0.2, all five
   recorded incoming requests and versions, zero target imports/load, runtime
   unreachability, every earlier guard, and no new advisory or independent
   defect.
2. Keep P7 blocked and authorize a separately designed Go-1.18-compatible
   patch, fork, replacement, or wrapper that preserves filesystem error
   identity, rejects an empty Darwin root pool, and supplies a proxy-safe test
   fixture. This requires a later implementation session plus full source,
   API, MVS, vulnerability, and changed-selection qualification; this decision
   session must not implement or silently design it.
3. Keep P7 blocked and authorize a separately scoped parent-chain and graph-
   pruning study involving Viper, `sagikazarmark/crypt`, historical Consul API,
   mvn-pom-mutator, and the unapplied tidy projection. Those parents and their
   guarded closures are outside this decision and cannot be changed here.

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
direct use, parent changes, patching, forking, replacement, a Go-floor change,
unrelated-module changes, another dependency group, or P8.

# Required Reading

At start read this archive, its answered go-rootcerts evaluation, the answered
go-retryablehttp, go-multierror, and go-hclog decisions, rolling handover,
roadmap, `go.mod`, and `go.sum`. Verify branch, clean ordinary and ignored
status, handoff ancestry and changed-file set, exact toolchain identity,
reciprocal archive history, all guarded selections/requests/load results, and
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
