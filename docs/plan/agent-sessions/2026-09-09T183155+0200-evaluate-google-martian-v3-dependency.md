# Agent Session: Evaluate Google Martian V3 Dependency

Status: NEXT
Session ID: `2026-09-09T183155+0200-evaluate-google-martian-v3-dependency`
Created: `2026-09-09T18:31:55+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `1a8f809c0511941d49bc3a99683d66adb10fbff9d22eccf57f51b5b451fc2e3f`
Previous: [2026-09-09T162826+0200-evaluate-golang-snappy-dependency.md](2026-09-09T162826+0200-evaluate-golang-snappy-dependency.md)
Next: none
Outcome: pending

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
