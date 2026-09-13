# Agent Session: Evaluate Enterprise Certificate Proxy Dependency

Status: NEXT
Session ID: `2026-09-13T041350+0200-evaluate-enterprise-certificate-proxy-dependency`
Created: `2026-09-13T04:13:50+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `72cfd252fdc1ebbe5642afa664924dc4f1114118dda60ecf9427707d3d05c4ff`
Previous: [2026-09-13T005310+0200-evaluate-google-cloud-go-testing-dependency.md](2026-09-13T005310+0200-evaluate-google-cloud-go-testing-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/googleapis/enterprise-certificate-proxy v0.2.1` as one bounded
dependency group. Resolve its repository and release identity, complete
Go-floor closure, package/command behavior and exported API, actual project
loading, exact MVS effects, and every applicable quality contract. The selected
release declares Go 1.19 and has no inherited exception. Retain or select only
an exact-path release whose complete minimal source/test closure preserves Go
1.18 and whose relevant behavior passes every contract; if no in-scope exact
selection can satisfy that rule, stop for a fresh product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, Google Renameio v1.0.1, and Google UUID
v1.4.0 moves. Google pprof remains retained at
`v0.0.0-20210720184732-4bb14d4b1be1`. Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, historical root GLFW, Googleapis GAX Go v2 v2.7.0, and Google
Cloud Go Testing
`v0.0.0-20200911160855-bcd43fbb19e8` also remain retained. All earlier
decisions and lifecycle ancestry are final. Do not revisit them or combine
another dependency group. P8 remains queued.

GAX v2.7.0 remains under the user's two explicit bounded decisions. Its
inherited Go-1.19 floor exception and known
`apierror.ParseError(err, false)` panic exception remain valid only while the
complete project load contains zero GAX packages. If this checkpoint loads or
directly imports GAX, stop for a fresh dependency and product decision before
merge. Do not reopen GAX or apply its exceptions to this target.

Project MVS selects Enterprise Certificate Proxy only through
`github.com/spf13/viper v1.15.0`; the shortest and only incoming path is main
-> Viper v1.15.0 -> Enterprise Certificate Proxy v0.2.1. `go mod why -m` is
negative, and the complete project load contains no package from this module.
Resolve this ancestry precisely. Do not combine, upgrade, downgrade, remove,
or independently audit Viper or another dependency group. Candidate MVS
movements caused by the exact Enterprise Certificate Proxy edge remain in
scope to measure, but any required Viper or unrelated selection change needs
fresh product direction rather than silent implementation.

A fresh minimal proxy survey found 29 stable exact-path versions. Selected
v0.2.1 is the tagged commit `80592736477602cc7992372d4280f819dc7e4cbf`,
dated 2022-12-07T03:49:14Z, with source/mod sums
`h1:RY7tHKZcRlk788d5WSo/e83gOyyy742E8GSs771ySpg=` /
`h1:AwSRAtLfXpU5Nm3pW+v7rGDHp09LsPtGY9MduiEsR9k=` and a Go 1.19
directive. V0.1.0 and v0.2.0 declare Go 1.18; v0.2.1 through v0.3.4 declare
Go 1.19; later releases declare Go 1.23 or newer. Thus v0.2.0 is the initial
highest declaration-eligible candidate, but its full closure, release identity,
behavior, API, and MVS feasibility are unproven and must be evaluated rather
than assumed.

Proxy latest is qualified tag v0.3.22 at commit
`62d25fa2858321169ff476109479205309fa3a68`, dated
2026-09-09T18:17:35Z, with source/mod sums
`h1:NU4XpII6jD+Dxcot94fqjE+AfJoE/lQP9q3faYGzC/c=` /
`h1:L3D/IQExI6LqEjBdXcZQ1WluSgigQmSwBboFstVPM4w=`. It declares Go
1.25.0 and toolchain Go 1.26.5 and is floor-ineligible. Verify every incoming
identity and floor fact. Do not silently promote latest, an unreleased commit,
arbitrary branch head, redirect, fork, alternate path, prerelease, or tag that
does not version this module.

# Measurements At Start

The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod`
and `go.sum` with three insertions and no deletions. Google Cloud Go Testing,
GAX, Google pprof, and Gogo Protobuf remain retained without dependency
commits or metadata edits.

Accepted project measurements remain 234 selected modules, 3,599 graph edges,
429 native complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,067 `go.sum` lines, and a 432-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, metadata adds exactly 51 checksum
lines and removes zero. The main module retains Go 1.18 and toolchain Go
1.26.7. Google Cloud Go Testing and GAX each have zero loaded packages.

Fresh primary vulnerability data has 1,398 module records, index Last-Modified
2026-09-10T16:28:28Z, and scanner DB update 2026-09-10T14:48:42Z. Accepted
project populations remain 30 module findings, 22 Darwin package findings, 23
Windows package findings, and 20 IDs/22 reachable traces on both Darwin and
Windows symbol scans. Do not attribute inherited findings to Enterprise
Certificate Proxy without exact evidence.

The accepted quality result remains all 27 Q0-Q2 rows PASS at L2, with seven
improved and zero held, regressed, non-comparable, or dirty counts. Scorecard
SHA-256 is
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
The fresh unchanged-project confirmation is also 27/27 PASS with scorecard
SHA-256
`a706ee72aa1483217bdcced30b36523da95841fb3874ef017645ffcf74051829`.
Google Cloud Go Testing decision-summary SHA-256 is
`9652c69c7971189749d2fd4590dc39b49a11f6f159cb80c2287767251312db2a`;
its 788-entry selected-evidence manifest SHA-256 is
`34e3574d2ced8aef523df1607300134664907cf7634cb40f89c3f615a151f106`.

Read the answered Google Cloud Go Testing and GAX archives plus the rolling
handover for their complete identity, closure, API, behavior, MVS,
vulnerability, quality, exceptions, and evidence records. Do not reopen Google
Cloud Go Testing, GAX, UUID, Renameio, pprof, Martian, Viper, or earlier groups.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff source archive
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
Docker timestamps with more than six fractional digits require the real Python
3.14 interpreter ahead of `/usr/bin/python3`; preserve this as an external
execution control.

# Role And Boundaries

From fresh external archives and caches, prove the complete minimal module and
package/test closure for selected, v0.2.0, and every other serious candidate
under exact Go 1.26.7 and contained Go 1.18.10. Inspect imported source and
test dependencies rather than treating module Go directives as full closure
proof. Separate isolated source-time resolution from the project's selected
graph, and report each module's declared floor and the floor of every
dependency in its minimal production and test closure.

Inspect every package, command, and exported API. The selected archive initially
contains `client`, `client/util`, the `cshared` command, and an internal test
signer. Characterize configuration discovery/parsing and path precedence;
signer executable launch and RPC framing; stdin/stdout/stderr handling;
process lifecycle, Close, errors, exits, malformed responses, timeouts and
cancellation; `crypto.Signer` semantics, public key and certificate-chain
parsing, hash options and digest validation; nil/zero inputs, typed nils,
aliasing, concurrency and platform paths; and the C-shared ABI's pointer,
length, ownership, output-buffer, status-code, and panic boundaries.
Distinguish runtime libraries, commands, internal helpers, examples,
benchmarks, testdata, generated content, fuzz/property coverage, platform
files, native/cgo requirements, external executables, credentials, and
upstream CI. Do not infer guarantees from declarations alone or run an
untrusted signer outside contained scratch fixtures.

Add independent fixtures where useful for exact config resolution, relative
and absolute paths, subprocess argument/environment and JSON identities,
successful and failing signer exchanges, partial/malformed output, close and
double-close, digest/hash mismatch, key/certificate parsing, caller buffer
aliasing, cancellation or unbounded waits, concurrent calls, native/cgo
exports, and supported platforms. Run source verification, package listing,
native complete tests, two independent repeats, race, vet, and meaningful
cross-builds under both SDKs. Classify every subprocess, filesystem, network,
credential, cgo, compiler, platform, resource, timeout, or test-design failure
precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected and
every serious candidate in disposable trees. Explain why Enterprise
Certificate Proxy exists in MVS while no package is loaded. An explicit
v0.2.0 request may be overridden by Viper's higher minimum or may force
unrelated movements; measure this exactly. Do not add a replace directive,
patch fork, or direct root edge merely to evade MVS. Any required Viper or
unrelated module change is a stop for product direction, not authorization to
broaden the group.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject, retain, or stop if canonical
identity, release qualification, complete floor, behavior, concurrency,
subprocess/C-ABI safety, tests, API, loading, MVS, or any applicable quality
contract fails. Zero project-loaded packages does not excuse source-time
floor, closure, behavior, or vulnerability evidence.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this
archive, the answered Google Cloud Go Testing and GAX archives, rolling
handover, roadmap, `go.mod`, `go.sum`, and every referenced quality,
compatibility, release, runner, evidence, and lifecycle contract. Earlier
outcomes are final.

# Three Moves

First, determine the highest qualified exact-path release whose full minimal
closure preserves Go 1.18. Do not promote a redirect, fork, alternate path,
floor-ineligible release, prerelease, non-versioning tag, or unreleased commit.
Selected v0.2.1 has no inherited Go-floor exception. If exact-path v0.2.0 is
qualified and technically sound but MVS cannot select it without changing
Viper or another group, preserve the evidence and stop for a fresh bounded
product decision; do not silently retain v0.2.1, manufacture an exception, or
implement out-of-scope metadata.

For an authorized changed selection, use exact Go 1.26.7 and exact `go get`
for one dependency-only commit, never tidy as implementation, then run the
complete P7 dependency gate. For a retained or blocked selection, prove the
no-change effect and run applicable dependency, closure, project,
compatibility, vulnerability, empty-HOME, and cleanliness gates without
manufacturing activity. Full changed-selection quality must preserve 27/27
Q0-Q2 PASS at L2 with zero held, regressed, non-comparable, or dirty counts.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, signer fixture, C-shared output, and build context beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never run `go mod download all` in a measured
worktree or create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the Enterprise Certificate Proxy decision or bounded product stop,
rewrite the rolling handover and roadmap, answer this archive, and follow the
repository lifecycle contract. Do not implement another dependency group,
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
