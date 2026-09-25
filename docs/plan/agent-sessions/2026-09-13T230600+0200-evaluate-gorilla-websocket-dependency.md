# Agent Session: Evaluate Gorilla WebSocket Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-13T230600+0200-evaluate-gorilla-websocket-dependency`
Created: `2026-09-13T23:06:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `11cfe304526d66a79d06d34ca9b840b3605e26b7c81ddeac1079cce25f725a02`
Previous: [2026-09-13T172733+0200-decide-gopherjs-product-direction.md](2026-09-13T172733+0200-decide-gopherjs-product-direction.md)
Next: [2026-09-14T001519+0200-decide-gorilla-websocket-product-direction.md](2026-09-14T001519+0200-decide-gorilla-websocket-product-direction.md)
Outcome: Rejected every published exact-path stable candidate under the combined Go 1.18, release, source-security, and full-test contracts; preserved exact selected v1.4.2 without metadata changes and stopped for a bounded product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/gorilla/websocket v1.4.2` as one bounded dependency group. Resolve
its complete repository and release identity, full Go-floor closure, package
and connection behavior, exported API, actual project loading, exact MVS
effects, vulnerability evidence, and every applicable quality contract. Retain
or select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, Google Renameio v1.0.1, and Google UUID
v1.4.0 moves. Google pprof, Gogo Protobuf, Crypt, OpenCensus Proto, Logex,
Readline, Fnmatch, Imaging, ansimage, Fsnotify, Ghodss YAML, historical root
GLFW, Googleapis GAX Go v2 v2.7.0, Google Cloud Go Testing, Enterprise
Certificate Proxy v0.2.1, and exact selected GopherJS remain retained. All
earlier decisions and lifecycle ancestry are final. Do not revisit them or
combine another dependency group. P8 remains queued.

The user's 2026-09-13 GopherJS option 1 decision is final. Exact
`github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1` retains
only its documented release, closure, compiler, command, and generated-runtime
exceptions while zero GopherJS packages load, it remains runtime-unreachable,
its exact version and sole incoming GoConvey v1.6.4 edge remain unchanged, and
no new advisory or independent disqualifier appears. Direct import/loading,
runtime reachability, a version or incoming-edge change, or a new advisory or
independent disqualifier expires the exception. Revalidate these guards before
work and stop for a fresh dependency and product decision if any fails. Do not
reopen GopherJS, change mvn-pom-mutator or GoConvey, or transfer its exceptions
to Gorilla WebSocket.

The Enterprise Certificate Proxy v0.2.1 exceptions remain valid only for its
exact version and sole Viper v1.15.0 edge while zero target packages load, it
remains runtime-unreachable, and no new advisory or independent disqualifier
appears. The two GAX v2.7.0 exceptions remain valid only while zero GAX
packages load. Revalidate those zero-load guards and stop for a fresh owning
decision if one expires. Do not reopen or transfer either exception.

# Measurements At Start

The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. GopherJS, Enterprise
Certificate Proxy, Google Cloud Go Testing, GAX, Google pprof, and Gogo
Protobuf were retained without dependency implementation commits or metadata
edits.

The clean GopherJS decision revalidated exact selected GopherJS and its sole
GoConvey v1.6.4 edge, a negative `go mod why -m` result, and zero source imports.
The unchanged project has 234 selected modules, 3,599 graph edges, 429 complete-
test entries, 197 module-backed packages, 41 loaded modules, 1,067 `go.sum`
lines, and a 432-line unapplied tidy projection. Exactly zero GopherJS,
Enterprise Certificate Proxy, or GAX packages load. `go.mod` and `go.sum`
SHA-256 values are
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2 with scorecard SHA-256
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.

Fresh primary vulnerability data contains 1,398 module records, was last
modified 2026-09-10T16:28:28Z, and has no exact GopherJS record. Selected
Gorilla WebSocket v1.4.2 is only a starting MVS selection, not proof of
canonical repository, release qualification, ancestry, floor, loading,
behavior, vulnerability state, or suitability. Independently resolve its
incoming graph paths, `go mod why -m` result, loaded-package population, and
every serious exact-path stable candidate. Do not add a direct edge merely to
alter MVS.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, and use
`LC_ALL=C LANG=C`. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and the recorded apidiff source receipt. Preserve the known apidiff archive
reproducibility discrepancy and Python 3.14 Docker timestamp control; neither
is Gorilla WebSocket evidence.

# Role And Boundaries

From fresh external archives and caches, resolve proxy, sumdb, `go-import`,
Git, and forge evidence for the exact module path: tags, releases, branches,
signatures, commits, parents, trees, times, ancestry, repository status,
licenses, retractions, deprecations, redirects, forks, alternate paths, and
every serious exact-path stable candidate. Do not silently promote a redirect,
fork, alternate path, prerelease, non-versioning tag, unreleased branch head,
or floor-ineligible release.

Prove the complete minimal production and test closure under exact Go 1.26.7
and contained Go 1.18.10. Inspect imported source and test dependencies rather
than treating the module directive alone as floor proof. Separate isolated
source-time resolution from the project's selected MVS graph.

Inspect every package, exported API, example, benchmark, fuzz target, testdata,
generated file, and platform or build-tag branch. Characterize client and
server handshakes, URL/header/cookie/auth handling, origin and subprotocol
selection, proxies and TLS, compression negotiation, frame validation and
fragmentation, UTF-8 rules, control frames, close codes, ping/pong handlers,
deadlines, read limits, prepared messages, JSON helpers, buffers and pools,
error identity, nil and malformed inputs, aliasing, concurrency constraints,
connection lifecycle, and network/resource cleanup. Distinguish documented
single-reader/single-writer requirements from independently unsafe behavior.

Add independent fixtures where useful for handshake success/failure, malformed
and oversized frames, fragmentation, control-frame interleaving, close
round-trips, deadlines/cancellation, compression, subprotocol/origin policy,
buffer reuse, concurrent supported use, leak cleanup, deterministic outputs,
and selected-project compatibility. Run source verification, package listing,
native complete tests, two independent repeats, race, vet, and meaningful
cross-builds under both SDKs. Classify every network, TLS, proxy, clock,
timeout, filesystem, subprocess, platform, resource, or test-design failure
precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent or unrelated module change needs a fresh bounded
decision rather than silent implementation. Compare primary vulnerability
results at module, package, symbol, test-symbol, and reachable-trace levels.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, unchanged exception guards, and
`./codex-dev-start.sh --check`. Read this archive, the answered GopherJS
decision/evaluation, answered Enterprise Certificate Proxy decision/evaluation,
answered GAX archive, rolling handover, roadmap, `go.mod`, `go.sum`, and every
referenced quality, compatibility, release, runner, evidence, and lifecycle
contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose complete
minimal source/test closure preserves Go 1.18. Do not promote an unqualified or
floor-ineligible identity. If no candidate satisfies the existing contracts,
preserve the evidence and stop for a fresh bounded product decision.

Second, qualify selected and serious candidates across source, tests, API,
behavior, loading, MVS, vulnerability, and project contracts. For an authorized
changed selection, use exact Go 1.26.7 and exact `go get` for one dependency-
only commit, never tidy as implementation, then run the complete P7 dependency
gate. For a retained or blocked selection, prove the no-change effect and run
all applicable gates without manufacturing activity. Full changed-selection
quality must preserve 27/27 Q0-Q2 PASS at L2 with zero held, regressed, non-
comparable, or dirty counts.

Third, record the exact decision and evidence in the roadmap and rolling
handover, answer this archive, and follow the lifecycle contract for one next
bounded authorized mission. Do not execute the successor in this turn.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, fixture, runtime/tool installation, and build context beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never run `go mod download all` in a measured
worktree or create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the Gorilla WebSocket decision or bounded product stop, rewrite the
rolling handover and roadmap, answer this archive, and follow the repository
lifecycle contract. Do not implement another dependency group, launch a
successor, push, merge, publish, release, stash, revert, bypass cleanup, or
remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Stopped for a fresh bounded Gorilla WebSocket product decision without changing
`go.mod` or `go.sum`. No published exact-path stable release satisfies the
combined release, complete Go-1.18 closure, source-security, and full-test
contracts. No dependency implementation commit was manufactured.

Fresh proxy, sumdb, `go-import`, strict Git, and GitHub evidence resolves
canonical `https://github.com/gorilla/websocket.git`, a public active
unarchived non-fork BSD-2-Clause repository. The proxy lists eleven versions.
Serious stable tags are v1.4.2, v1.5.0, v1.5.1, and v1.5.3. GitHub marks
v1.5.2 prerelease. Its proxy record preserves original commit
`1bddf2e0dba6f35492b0f5614905b291cd0ab88d`, while the current forge tag was
overwritten to `9ec25ca502ec881a7e873d3cdf35da60eab22037`; unreleased main
later retracts v1.5.2 for that tag accident. Strict fsck and serious-tag
ancestry pass. There is no redirect, alternate module path, qualified fork,
deprecation, or promotable branch head.

Selected v1.4.2 is commit
`b65e62901fc1c0d968042419e74789f6af455eb9`, parent `8c288dca`, tree
`2a570dce`. V1.5.0 is `9111bb834a68b893cebbbaed5060bdbc1d9ab7d2`;
v1.5.1 is `ac0789be11725ab2285233e9a3800c2312cff4fc`; and v1.5.3 is
`ce903f6d1d961af3a8602f2842c8b1c3fca58c4d`, parent `9ec25ca`, tree
`0cef0944`. Security commit
`d67f41855da42d7bccd9ef050c49f7e54e783b95` was committed after v1.5.3
and is only on unreleased main. Proxy and Git archives are byte-identical for
v1.4.2, v1.5.0, and v1.5.3. V1.5.1/v1.5.2 zips retain
`vendor/modules.txt` but omit vendored packages, breaking default vendor mode.

V1.4.2, v1.5.0, and v1.5.3 declare Go 1.12 and have standard-library-only
minimal closures: one module, 199/221 production/test load entries under exact
Go 1.26.7 and 135/158 under Go 1.18.10. V1.5.1/v1.5.2 declare Go 1.20 and
use six-module graphs with an actual two-package/one-module `x/net` closure.
V1.5.2 imports `http.NewResponseController` and fails Go 1.18 compilation.
V1.5.3 is therefore the highest stable floor-preserving release, but floor
qualification does not cure its source-security failure.

Fresh primary data contains 1,398 module records and was last modified
2026-09-10T16:28:28Z. GO-2026-6278, alias GHSA-w67g-5rqw-f597, identifies weak
PRNG use for client WebSocket mask keys, lists 34 affected symbols, marks
versions before v1.5.3 affected, and labels v1.5.3 fixed. Actual release source
contradicts that range: v1.5.3 explicitly restores `math/rand`, and a seed-
controlled independent fixture proves repeatable `newMaskKey` output. Commit
`d67f418...` changes production to `crypto/rand` only after v1.5.3. Conversely,
the original proxy v1.5.2 source contains that fix, but v1.5.2 is prerelease and
Go-floor-ineligible. V1.4.2, v1.5.0, v1.5.1, and v1.5.3 are source-vulnerable;
there is no qualified published stable release.

The reviewed govulncheck v1.8.0 snapshot updated
2026-09-10T14:48:42Z omits the unreviewed GO record, so its direct scans show
zero target records. That omission does not override fresh primary or inspected
source evidence. Project target package, symbol, test-symbol, and reachable
traces are zero because no Gorilla package loads, while primary module evidence
identifies selected v1.4.2 as affected.

All five candidate trees, packages, examples, tests, benchmarks, generated and
build-tag files, APIs, and relevant connection behavior were inspected. Pinned
apidiff records the added `Dialer.NetDialTLSContext` and `Conn.NetConn`; v1.5.2
adds and v1.5.3 removes `FormatMessageType`. V1.5.3 remains API-compatible with
selected. The source review covers client/server handshake, URL/header/cookie/
auth handling, origins and subprotocols, proxy/TLS paths, compression, frame
validation/fragmentation/UTF-8/control frames, close/error identity, deadlines,
read limits, prepared messages, JSON, buffers/pools, aliasing, documented
one-reader/one-writer constraints, lifecycle, and cleanup.

The independent fixture at SHA-256
`f922190c239b907108adaf48c251b56e6fda5460f0a66f432a6731e1e04e3c8f`
exercises real handshake success/failure, headers/auth/cookies, origin and
subprotocol policy, compression, echo/close round-trips, fragmentation with
interleaved controls, malformed/unmasked and oversized frames, UTF-8/close
handling, deadlines, prepared-message copying/concurrent reuse, buffer pools,
deterministic masks, and transport cleanup. It passes selected and v1.5.3
count-1/count-10/race/vet under both SDKs. A repeated full-source Go 1.26 race
selection independently exposes an upstream cross-test lifecycle defect in
both: `cstHandler` logs after its test has returned. The affected test passes
alone, the connection group reproduces it, and the fixture is race-clean; this
is test-design cleanup, not unsupported concurrent connection use.

The sole selected incoming edge is direct
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3`; `go mod why -m` is
negative, repository source imports are zero, and zero target packages load.
Exact v1.4.2/v1.5.0/v1.5.3 requests create a redundant direct requirement but
no unrelated selection; tidy returns byte-identically to base and restores
v1.4.2. V1.5.1/v1.5.2 requests add five unrelated `x/*` selection changes
that survive tidy. No projection was applied.

The unchanged project remains 234 modules, 3,599 graph edges, 429 complete-test
entries, 197 module-backed packages, 41 loaded modules, 1,067 sum lines, and a
432-line unapplied tidy projection. `go.mod` and `go.sum` remain byte-identical
at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
GopherJS, Enterprise Certificate Proxy, and GAX exact versions, incoming-edge
and zero-load guards remain valid and separate.

Exact Go 1.26.7 project verify/load/build, count-1, two independent count-10
runs, race, vet, pinned lint, empty-HOME count-2, Linux/Windows builds, API/CLI
compatibility, and full authoritative preflight pass. The Go 1.18.10 projection
passes applicable verify/load/build/vet/cross-build/repeat/race gates and retains
only the two accepted `pkg/shell` closed-file wording failures. Accepted quality
remains 27/27 Q0-Q2 PASS at L2 with scorecard SHA-256
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
The selected 1,074-file evidence digest is
`cc15144b88edbf330e1300dda430eead2419087883d97788edb7d781bb48b912`;
the decision-summary digest is
`4f36e9a15e7e17cf9a09b79ef4f875e2dbecbb7f171682d803aca9514e8d7f91`.

The next bounded decision offers three directions: recommended retention of
exact inherited unloaded v1.4.2 under a new exact-version/sole-edge/zero-load
security and test-design exception; a broader owning-parent removal evaluation
that explicitly reopens the final GopherJS edge guard; or evaluation of a
maintained patch/fork/unreleased strategy under a new provenance, release,
maintenance, vulnerability, API, MVS, and floor policy. It explicitly rejects
selecting v1.5.3 merely to silence the currently incorrect advisory boundary.
