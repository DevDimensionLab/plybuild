# Agent Session: Evaluate Google Martian V3 Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-09T183155+0200-evaluate-google-martian-v3-dependency`
Created: `2026-09-09T18:31:55+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `1a8f809c0511941d49bc3a99683d66adb10fbff9d22eccf57f51b5b451fc2e3f`
Previous: [2026-09-09T162826+0200-evaluate-golang-snappy-dependency.md](2026-09-09T162826+0200-evaluate-golang-snappy-dependency.md)
Next: [2026-09-09T214949+0200-evaluate-google-pprof-dependency.md](2026-09-09T214949+0200-evaluate-google-pprof-dependency.md)
Outcome: Upgraded exact-path Google Martian from v3.2.1 to highest API-compatible qualified stable v3.3.2 at dependency-only commit `4644476`; complete identity, Go-floor closure, API and behavior, MVS/loading, vulnerability, project, and exact quality contracts pass. Latest v3.3.3 was rejected because it removes exported `martian.Init`.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/google/martian/v3 v3.2.1` as one bounded dependency group. Resolve
its complete release and repository identity, full Go-floor closure, package
behavior and public API, actual project loading, exact MVS effects, and every
applicable quality contract. Retain or select only an exact-path version whose
complete minimal closure preserves Go 1.18 and whose relevant behavior passes
every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, and Golang
Snappy v1.0.0 moves. Gogo Protobuf v1.3.2, Crypt, OpenCensus Proto, Logex,
Readline, Fnmatch, Imaging, ansimage, Fsnotify, Ghodss YAML, and historical
root GLFW remain retained. All earlier decisions and lifecycle ancestry are
final. Do not revisit them or combine another dependency group. P8 remains
queued.

Project MVS selects Google Martian v3.2.1. The module graph contains v3.2.1
declarations from eight historical `cloud.google.com/go` vertices v0.83.0,
v0.84.0, v0.87.0, v0.90.0, v0.93.3, v0.94.1, v0.97.0, and v0.99.0;
the project currently selects Cloud Go v0.105.0. `go mod why -m
github.com/google/martian/v3` says the main module does not need it, and the
complete project package load contains no Martian package. Resolve precisely
why the module remains selected. Do not combine, upgrade, remove, or
independently audit Cloud Go, Golang Snappy, Golang/Google Protobuf, gRPC,
X/Net, or another dependency group. Candidate MVS movements caused by the
exact Martian edge remain in scope to measure, but not to broaden into
independent audits.

A minimal post-Snappy survey finds five proxy versions: v3.0.0, v3.1.0,
selected v3.2.1, v3.3.2, and latest v3.3.3, with no prerelease. Selected
v3.2.1 is dated 2021-05-19T22:06:43Z, has source/mod sums
`h1:d8MncMlErDFTwQGBK1xhv026j9kqhvw1Qv9IbWT1VLQ=` /
`h1:oBOf6HBosgwRXnUGWUB05QECsc6uvmMiJ3+6W4l/CUk=`, declares Go 1.11,
and requires Golang Protobuf v1.5.2, Golang Snappy v0.0.3, historical X/Net,
gRPC v1.37.0, protoc-gen-go-grpc v1.1.0, and Google Protobuf v1.26.0.
Latest v3.3.3 is dated 2022-08-16T15:12:57Z at proxy origin commit
`0f7e6797a04da412118541344bbe0d65945e24c9`, has source/mod sums
`h1:DIhPTQrbPkgs2yJYdXU/eNACCG5DVQjySNRNlflZ9Fc=` /
`h1:iEPrYcgCF7jA9OtScMFQyAlZZ4YXTKEtJ1E6RWzmBA0=`, and declares Go 1.18.
Treat every incoming fact only as a survey to verify.

Resolve proxy, sumdb, go-import metadata, repository tags/releases/branches,
signatures, commits, times, trees, parents, ancestry, repository status,
deprecation, retractions, redirects, forks, alternate module paths, and every
serious exact-path candidate. Do not silently promote an unreleased commit,
redirect, fork, alternate path, floor-ineligible release, prerelease, or tag
that does not version this module.

# Measurements At Start

The latest dependency implementation is exact Golang Snappy v1.0.0 commit
`372f8e988d9bb86a4f426ae6f953d2e8226c05fe`, parent
`838fcd3dff204535d42b3b4078041688ad6de40d`, and tree
`044149546152b4ce44eb1189721b88f508ea9477`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Gogo Protobuf v1.3.2 remains
retained without a dependency commit or metadata edit.

Accepted project measurements are 234 selected modules, 3,590 graph edges,
429 native complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,061 `go.sum` lines, and a 407-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, metadata adds exactly 45
checksum lines and removes zero. The main module retains Go 1.18 and toolchain
Go 1.26.7.

Fresh primary vulnerability data has 1,392 module records. Accepted project
populations are 30 Darwin module findings, 22 Darwin package findings, 23
Windows package findings, and 20 IDs/22 reachable traces for both Darwin and
Windows symbol scans. Snappy has no finding or trace. Do not attribute
inherited findings to Martian without exact evidence.

The accepted quality result is all 27 Q0-Q2 rows PASS at L2, with scorecard
SHA-256
`7320f432090c578973bcfe4abc8736ef68a83f492a5bd057922ec150ce2a03c5`.
Golang Snappy decision-summary SHA-256 is
`3d606c0a6ca3f0cd523c4603de9ded354311d1ac6135c9f679ee58371d3ebc79`;
its 677-entry selected-evidence manifest SHA-256 is
`e4c81eed9a21f737f3a8427360908835b7f195f9981bd1be9a5bdbd953091449`.

Read the answered Golang Snappy archive and rolling handover for its complete
release, closure, API, behavior, MVS, vulnerability, quality, and evidence
record. Do not reopen Snappy, Golang Protobuf, Cloud Go, Google Protobuf, Go
CMP, or earlier groups.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and
inject no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools
beneath scratch as required. Portable receipts remain golangci-lint 2.12.2
archive `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff source archive
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

# Role And Boundaries

From fresh external archives and caches, prove the complete minimal module
and package/test closure for selected and every serious candidate under exact
Go 1.26.7 and contained Go 1.18.10. Inspect imported source and test
dependencies rather than treating a module Go directive as full closure
proof. Keep isolated source-time resolution separate from the project's
selected graph.

Inspect every package and exported API. Characterize proxy construction and
lifecycle; HTTP request/response modifiers; Martian context and modifier
groups; filter, header, cookie, URL/query, body, and method behavior; HAR
recording; HTTP/2 and gRPC support; CONNECT tunneling, WebSockets, and MITM
TLS CA/key/certificate issuance, caching, host/SNI, and verification; listener
and port behavior; close/reset/error semantics; short or failing I/O;
buffering and resource limits; determinism; concurrency and reuse; build
tags/platform behavior; commands, examples, benchmarks, testdata,
fuzz/property coverage, and upstream CI. Distinguish runtime packages,
generated protocol code, commands, and test fixtures.

Add independent fixtures where useful for proxy lifecycle, chained and
filtered request/response modification, headers/cookies/URLs/bodies, context
isolation, HTTP and HTTPS CONNECT/MITM flows, certificate properties and
caching, WebSocket/tunnel passthrough, HAR/protocol round trips, malformed or
truncated traffic, short/failing I/O, shutdown/reset, resource boundaries,
reuse, concurrency, and compatibility with the selected project graph. Run
source verification, package listing, native complete tests, two independent
repeats, race, vet, and meaningful cross-builds under both SDKs. Classify
every generator, network, TLS, timing, platform, resource, or test-design
failure precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected
and every serious candidate in disposable trees. Explain why Martian exists
in MVS while no package is loaded, and preserve every unrelated module
selection. Any change outside the exact Martian edge and its necessary
authorized MVS projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, behavior, concurrency, tests, API,
loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this
archive, the answered Snappy archive, rolling handover, roadmap, `go.mod`,
`go.sum`, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release. Do not
promote a redirect, fork, alternate path, floor-ineligible release,
prerelease, or unreleased commit. If no higher release qualifies, retain
selected without hand-editing metadata or manufacturing a dependency commit.

For a changed selection, use exact Go 1.26.7 and exact `go get` for one
dependency-only commit, never tidy as implementation, then run the complete
P7 dependency gate. For an inert/retained selection, prove the no-change
effect and run applicable dependency, closure, project, compatibility,
vulnerability, empty-HOME, and cleanliness gates without manufacturing
activity. Full changed-selection quality must preserve 27/27 Q0-Q2 PASS at
L2 with zero held, regressed, non-comparable, or dirty counts.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Google Martian decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Upgrade exact-path `github.com/google/martian/v3` from v3.2.1 to v3.3.2.
V3.3.2 is the highest stable exact-path release whose complete minimal closure
preserves Go 1.18, whose exported API remains compatible with selected, and
whose applicable behavior, project, vulnerability, and quality contracts pass.
Exact Go 1.26.7 `go get github.com/google/martian/v3@v3.3.2` produced
dependency-only commit `4644476aef78e24d0c7e6a1140baba4c3226a1cc`, parent
`8508a9d54f066985fb81ebd35fa31f45fa36da40`, tree
`f20d9b3817fe29af4876745d09c14e9cbcac9152`, changing only `go.mod` and
`go.sum` with three insertions and no deletions.

The proxy exposes exactly five stable v3 releases: v3.0.0, v3.1.0, v3.2.1,
v3.3.2, and v3.3.3, with no prerelease or retraction. Proxy, sumdb, parent
go-import metadata, tagged Git, and GitHub resolve to the public, non-fork,
Apache-2.0 `https://github.com/google/martian` repository. It is archived and
its default `master` branch equals v3.3.3; there is no newer unreleased branch
commit. The module is not deprecated and has no qualified redirect, fork, or
alternate path. The repository also has invalid module-version tags/releases
`3.2` and `v3.2`; neither was considered. V3 tags are lightweight, so there
are no tag-object signatures. GitHub reports valid commit signatures for all
five exact v3 release commits; the applicable releases are neither draft nor
prerelease.

Selected v3.2.1 is commit
`7e75073889cd2324f33b959c4fb4545440da046c`, tree
`b39ce198afdb87dadbac7e9ca36082ed5b8a043e`, parent
`19163e1b8984e2577ea366d8aa97f0f75c6f6ac6`, at
2021-05-19T22:06:43Z. Candidate v3.3.2 is commit
`b41869697484173d2e936754152cd64e40a046d4`, tree
`a1680fb4f3f79b23ee2fde53a119d61841393f37`, parent
`8087be7bb3737e2de9a70007fbc0dc6c74c0fb1c`, at
2022-03-14T19:34:56Z. Latest v3.3.3 is commit
`0f7e6797a04da412118541344bbe0d65945e24c9`, tree
`0ad4a12b3b27833d80140eb0a54c2c1ba8d137b7`, parent
`d6ef5c8f4bee8c1322bdbeb59505bca21d5d5f0a`, at
2022-08-16T15:12:57Z. Normalized proxy and tagged-Git manifests match
byte-for-byte for every release; their ordered SHA-256 values are
`4c092e2a900cf633c6b081e20eb42a692fe9ea14ba61407e856eff8a98f796ac`,
`e37b2abf94dc9614c09f70a10be61d2c59b174f6a6b0f71aea614cee733abc6e`,
`50b3e3cffa269e8fa34ebba7cd748a2f96ffccd97c7c50271767ba434736e03c`,
`3f84285118e5c6daa8d19ea631505d04cd52be080bef1fd73f3a8345ec17fbe6`,
and `5a08b12494eca88f9c98c3dbc096e4c997f37e76988c5f7cc0558c95fdbbf997`.

V3.3.2 source/mod sums are
`h1:IqNFLAmvJOgVlpdEBiQbDc2EwKW77amAycfTuWKdfvw=` and
`h1:oBOf6HBosgwRXnUGWUB05QECsc6uvmMiJ3+6W4l/CUk=`. Selected and
candidate module files are byte-identical, declare Go 1.11, and resolve a
37-module complete standalone source-and-test closure. Latest declares Go
1.18, drops the generator module, and resolves 36 modules. Complete test
loading has 428 entries under Go 1.26.7 and 365 under Go 1.18.10, including
215 module-backed entries and 46 Martian packages. Compiling and running the
actual closure under contained Go 1.18.10 proves the project floor rather than
relying on module directives alone.

The 46 packages contain HTTP proxy and modifier libraries, context and group
types, filters, headers, cookies, query/body/method modifiers, HAR, HTTP/2 and
gRPC protocol support, MITM TLS, traffic shaping, generated protocol code,
helpers, and the `cmd/proxy` and `cmd/marbl` commands. V3.2.1 to v3.3.2 has no
exported API change. Its six-insertion/one-deletion source change moves global
flag parsing from library `martian.Init` into the proxy command and recognizes
wrapped `net.ErrClosed` during proxy shutdown. This removes a library
composability side effect and improves lifecycle error classification while
retaining API compatibility. V3.3.3 adds compatible `martian/log.Logger` and
`SetLogger`, but incompatibly removes exported `martian.Init` and the inherited
`cmd/marbl -v` flag; latest was therefore rejected.

Request and response modifier chains, filters, headers, cookies, URLs,
queries, bodies, methods, context/session isolation, HTTP proxying, shutdown,
raw CONNECT tunnelling, MITM CONNECT, certificate SAN/organization/cache
properties, HAR export/reset/protocol round trips, malformed traffic, failing
body reads, and concurrent reuse passed the independent 387-line fixture.
Its behavior and Init-contract SHA-256 values are
`d41342168319236bc936b4e4977db4460fd0c754195f955d1ad5342c8a21bd96`
and `6262ef23586fd72e15f430ef8f2b0f83efe8dbba940864b13f40ec4d49404c10`.
Candidate race count-10 passes under both SDKs.

The retained implementation is intentionally qualified at hostile-input and
lifecycle boundaries. Response modifier errors become `Warning` and
processing continues; upstream failures produce 502, while response write
errors are logged. `Proxy.Close` is not idempotent and does not own the
listener; raw CONNECT copy loops can await peer shutdown. MITM uses RSA-2048,
random certificate serials, one reused leaf key, and a locked but unbounded
per-host certificate cache. HTTP/2 preface forwarding uses one `Read`, not
`ReadFull`. gRPC frame length and decompression are not application-bounded,
and a truncated terminal frame can be accepted without an error. HAR body and
cache growth are unbounded and exported objects are not deep-copied. A
malformed body range can panic through an out-of-bounds slice. Traffic-shape
map iteration is nondeterministic. These behaviors are inherited by selected,
candidate, and latest; callers must bound hostile protocol/body inputs and
coordinate shutdown.

All serious candidates pass module verification, package listing, native
tests with package vet disabled, two independent repeats, race, commands, and
production vet under both SDKs. Upstream full vet reports only inherited
test-only Fatalf-from-goroutine findings, plus Go 1.26 formatting of an RSA
private-key test value. Saturated parallel repeats produced only classified
proxy startup timeouts and one sub-millisecond traffic-shape boundary miss;
serial repeats pass. All 1,380 package cross-build contexts pass across the
chosen Darwin, Linux, Windows, FreeBSD, and js/wasm targets. Upstream contains
271 tests, one example, no benchmark, and no native fuzz target.

Martian remains selected despite negative `go mod why -m` and zero loaded
Martian packages because eight encountered historical `cloud.google.com/go`
vertices, v0.83.0, v0.84.0, v0.87.0, v0.90.0, v0.93.3, v0.94.1, v0.97.0,
and v0.99.0, each declare v3.2.1. The module graph retains requirements of
encountered historical vertices even though MVS selects Cloud Go v0.105.0.
The explicit root edge now selects v3.3.2. No unrelated selection moves.

The accepted project has 234 selected modules, 3,597 graph edges, 429
complete-test entries, 41 loaded modules, 197 loaded module-backed packages,
1,063 `go.sum` lines, and a 418-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, sums are +47/-0. Exact v3.3.3 would have
3,599 edges but the same loaded and checksum populations. Main Go 1.18 and
toolchain Go 1.26.7 remain unchanged. Exact Go 1.26.7 verification, package
load, build, tests, repeats, race, vet, empty-HOME, and Linux/Windows builds
pass. The Go 1.18.10 projection passes verification, load, build, vet, and
cross-builds; full tests retain only the two accepted `pkg/shell` closed-file
error-wording assertions, while the compatible population passes repeats and
race. Martian is not loaded in either project population.

Fresh vulnerability data contains 1,392 module records and no Martian record.
The direct serious-candidate results are identical after separating patched
stdlib scanner metadata: 27 inherited module findings, 17 inherited package
findings, and five inherited IDs/81 reachable traces, none attributed to
Martian. Project base and v3.3.2 results are identical: 30 module findings, 22
Darwin package findings, 23 Windows package findings, and 20 IDs/22 reachable
traces on each symbol platform, with no Martian finding or trace.

Exact 21-stage `make quality` exits zero with 27/27 Q0-Q2 PASS at L2, 80/80
mutations killed, seven improvements, and zero held, regressed,
not-comparable, or dirty counts. Scorecard SHA-256 is
`094a7668e43eb0a5b563bcd637b3b5a301e1bb5f5053903ad8934e2b2906ae7f`.
Superseded attempts exposed only a cold offline Logrus cache, a launcher
signal-log race, BSD `mktemp` sandbox placement, and Docker's seven-digit
fractional timestamp exceeding Apple Python 3.9 parsing. Session-contained
caches and scratch-only shims resolved them without repository changes.

The 509-entry selected-evidence manifest SHA-256 is
`614921201697ce6f3accdf0c6e619c86ddb3502f23c3b03ee6314a2347a9a649`;
decision-summary SHA-256 is
`374c235a5ebeeb7a97254e8eab57b3dc60224fc7f103611beceafff701f5756e`.
