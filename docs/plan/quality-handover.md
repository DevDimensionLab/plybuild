# Quality Upgrade Handover

Generated: 2026-09-13T12:21:12+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation remains Google UUID v1.4.0 commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions. Enterprise Certificate
  Proxy, Google Cloud Go Testing, GAX v2.7.0, Google pprof, and Gogo Protobuf
  have no dependency implementation or metadata commit.
- Google Renameio v1.0.1 `394ec36`, Google Martian v3.3.2 `4644476`,
  Golang Snappy v1.0.0 `372f8e9`, Golang Protobuf v1.5.3 `6870e02`,
  Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1 `647d4fd`, Go Logfmt
  v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0
  `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain ancestors. All earlier
  dependency decisions are final and must not be reopened.
- The Enterprise Certificate Proxy evaluation is answered, and the user has
  supplied option 1 with bounded exceptions in the sole NEXT decision archive.
  No `.agent-task/current.md` or repository `.quality/manual-evidence.json`
  exists. Do not push, merge, publish, release, stash, revert, launch a
  successor, bypass cleanup, or remove the worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 is recording the user's bounded Enterprise Certificate
Proxy option 1 decision after exact Go 1.26.7 and accepted Speakeasy v0.2.0,
XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0, Go Stack v1.8.1,
Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy v1.0.0, Google
Martian v3.3.2, Google Renameio v1.0.1, and Google UUID v1.4.0 moves. Google pprof
`v0.0.0-20210720184732-4bb14d4b1be1`, Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, historical root GLFW, Googleapis GAX Go v2 v2.7.0, and Google
Cloud Go Testing `v0.0.0-20200911160855-bcd43fbb19e8` remain retained. P8
remains queued. Do not start another dependency group until the decision is
recorded in a committed handoff.

The user's two 2026-09-12 GAX decisions remain final. V2.7.0 retains its
inherited Go-1.19 floor exception and known
`apierror.ParseError(err, false)` panic exception only while the complete
project load contains zero GAX packages. If any GAX package becomes loaded or
directly imported, the behavior exception expires and the owning checkpoint
must stop for a fresh dependency and product decision. Those GAX exceptions do
not transfer to Enterprise Certificate Proxy, which now has its own separately
bounded decision below.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
run `go mod download all` in a measured worktree or bypass launcher cleanup.

## Enterprise Certificate Proxy Product Decision

There is no qualified, technically acceptable, in-scope exact-path release of
`github.com/googleapis/enterprise-certificate-proxy` under the current Go 1.18
and behavior contracts. On 2026-09-13 the user explicitly selected option 1
with the recommended bounds: retain exact selected, unloaded v0.2.1 without a
dependency edit. Keep the main Go 1.18 floor and accept v0.2.1's Go 1.19
declaration only for this exact version and the unchanged inherited Viper
v1.15.0 edge.

The decision also accepts only the already documented v0.2.1 release,
behavior, and safety findings: public-preview classification; lost executable
bit and raw proxy-test failure; configuration and file-lifetime boundaries;
unbounded signer/RPC waits; incomplete cleanup; error-string-dependent and
non-idempotent close; nil/panic and aliasing boundaries; global-logger side
effects; and the recorded unsafe C ABI boundaries. It does not accept a new or
independently discovered defect.

The exceptions remain valid only while the complete project load contains zero
target packages and selected v0.2.1 plus its sole incoming Viper edge remain
unchanged. Revalidate and record zero loading in the committed decision. Direct
import or loading, runtime reachability, a target version or incoming-edge
change, or a new vulnerability/advisory or independent disqualifier expires
the exceptions and requires a fresh dependency and product decision before
merge. Do not modify metadata, add a direct edge, patch or replace the target,
downgrade or independently audit Viper/GAX, raise the Go floor, or ask for this
same decision again while the invariants hold.

The exact path resolves to the public, active, unarchived, non-fork Apache-2.0
repository `https://github.com/googleapis/enterprise-certificate-proxy.git`.
The proxy contains 29 stable-semver versions. Git tag `0.3.10` lacks the `v`
prefix and does not version this module. GitHub release metadata marks v0.1.0,
v0.2.0, v0.2.1, and v0.3.15 as prereleases even though their proxy versions
have no semver prerelease suffix.

Selected v0.2.1 is commit
`80592736477602cc7992372d4280f819dc7e4cbf`, parent
`bee115d4cb1a6e7db50513cbdcb000697b431169`, tree
`d1d2b29329ca44130f6e3f333551b89d11648811`, dated
2022-12-07T03:49:14Z. Its source/mod sums are
`h1:RY7tHKZcRlk788d5WSo/e83gOyyy742E8GSs771ySpg=` and
`h1:AwSRAtLfXpU5Nm3pW+v7rGDHp09LsPtGY9MduiEsR9k=`. Highest declaration-
eligible v0.2.0 is commit `c5f65f94cd903cc16bba58964863229eb35ec41c`,
tree `26607aaea8d3ebf96ab4f6f98041ed7ec5e149f5`, dated
2022-09-28T23:03:15Z, with source/mod sums
`h1:y8Yozv7SZtlU//QXbezB6QkpuE6jMD2/gfzk4AftXjs=` and
`h1:8C0jb7/mgJe/9KK8Lm7X9ctZC2t60YyIpYEI16jx0Qg=`.

Proxy latest v0.3.22 is commit
`62d25fa2858321169ff476109479205309fa3a68`, dated
2026-09-09T18:17:35Z, and current `main` equals it. Its source/mod sums are
`h1:NU4XpII6jD+Dxcot94fqjE+AfJoE/lQP9q3faYGzC/c=` and
`h1:L3D/IQExI6LqEjBdXcZQ1WluSgigQmSwBboFstVPM4w=`. It declares Go 1.25.0
and toolchain Go 1.26.5 and is ineligible.

## Floor, API, Behavior, And ABI Evidence

V0.1.0 and v0.2.0 declare Go 1.18; v0.2.1 through v0.3.4 declare Go 1.19;
all later releases declare Go 1.23 or newer. V0.1.0, v0.2.0, and v0.2.1 have
no external module requirements, so their complete minimal production/test
closure is the root module plus standard library. Under Go 1.26.7 the closure
has 202 production and 221 complete-test entries; under Go 1.18.10 it has 137
and 157. The four root-module packages are `client`, `client/util`, command
`cshared`, and the internal test signer command. Repository platform signer
submodules have independent `go.mod` files and are excluded from the root zip.

V0.2.0 and v0.2.1 export identical Go APIs. V0.1.0 is an incompatible
fallback because v0.2.0 removed `client/util.Libs.SignerBinary` and added
`Libs.ECP`. Both v0.2 module zips lose the executable bit on
`client/testdata/signer.sh`, so raw proxy tests fail under both SDKs with
`permission denied`. A scratch-only mode correction makes upstream tests,
two independent repeats, race, vet, native builds, and broad cross-builds
pass; it characterizes source but does not repair release qualification.

Configuration accepts an explicit path or platform gcloud default, does not
expand environment variables or tilde, resolves relative signer paths from
process cwd, accepts unknown and duplicate JSON fields, has no read-size cap,
and does not close the opened config file. Client startup inherits cwd/env,
passes only the config path argument, reserves stdin/stdout for gob net/rpc,
and writes signer stderr to parent stderr. There is no context, deadline, or
timeout. Startup and signing can block forever, and failed initialization does
not close, kill, or wait for the child.

`Close` is non-idempotent. V0.2.0 relies on exact platform-dependent process
and RPC error strings and can panic on nil errors; v0.2.1 retains exact RPC
string and nil risks. Certificate chains alias internal mutable slices, typed-
nil hash options panic, and nil receivers have panic boundaries. Concurrent
successful signing passes a 32-call race fixture. V0.2.1 additionally sets the
process-global Go logger to discard output when its logging variable is first
unset and never restores it; upstream later removed that behavior because it
suppressed unrelated application logging.

The C-shared ABI exports `GetCertPemForPython` and `SignForPython`. Valid calls
work with caller-owned buffers under both SDKs, but the boundary trusts raw
pointer/length pairs through `unsafe.Slice`, silently truncates short
certificate buffers, conflates failures as status zero, and has no panic
recovery. Isolated child probes show nil-positive and negative lengths escaping
as uncaught Go panics; oversized lengths permit out-of-bounds access. Go
1.18.10 native linking needs `-ldflags=-w` only because its Darwin debug output
is incompatible with current Xcode `dsymutil`.

## Project MVS, Vulnerability, And Quality

The unchanged project has 234 selected modules, 3,599 graph edges, 429 native
complete-test entries, 197 module-backed packages, 41 loaded modules, 1,067
`go.sum` lines, and a 432-line unapplied tidy projection. The sole target edge
is Viper v1.15.0 -> Enterprise Certificate Proxy v0.2.1, reached from the main
module's direct Viper edge. `go mod why -m` is negative and zero target packages
load because MVS follows module requirements independently of package imports.
GAX also remains at zero loaded packages.

Exact `go get github.com/googleapis/enterprise-certificate-proxy@v0.2.0` in a
disposable tree downgrades Viper v1.15.0 to v1.14.0, GAX v2.7.0 to v2.6.0,
and changes 20 selections total: target, Viper, and 18 unrelated modules. The
raw projection has 233 modules, 3,601 edges, unchanged 429/197/41 loading
counts, 1,069 sum lines, and a 433-line tidy projection. Target and GAX remain
unloaded. Tidy still leaves 21 selection differences from the base tidy
control. Exact v0.1.0 downgrades Viper to v1.13.0 and tidy removes the target
entirely. Neither request is an authorized target-only dependency move.

Fresh primary vulnerability data contains 1,398 module records, index
Last-Modified 2026-09-10T16:28:28Z, scanner DB update
2026-09-10T14:48:42Z, and no target record. Direct v0.1.0/v0.2.0/v0.2.1 scans
have zero module, Darwin/Windows package, symbol, or test-symbol findings. The
unchanged project and exact raw v0.2.0 projection match at 30 module findings,
22 Darwin package findings, 23 Windows package findings, and 20 IDs/22
reachable traces on both symbol platforms, with no target or GAX occurrence.

The unchanged exact Go 1.26.7 project passes verify/load/build, native tests,
two count-10 repeats, race, vet, pinned lint, empty-HOME, Darwin/Linux/Windows
build coverage, and authoritative full preflight. The Go 1.18.10 projection
passes compatible gates with only the two already accepted `pkg/shell`
closed-file wording failures. The disposable v0.2.0 MVS projection has the
same result under both SDKs. No dependency metadata changed, so a changed-
selection scorecard is inapplicable; accepted quality remains 27/27 Q0-Q2
PASS at L2 with scorecard SHA-256
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.

The 562-entry selected-evidence manifest and decision summary hash to
`59868d9547ece628b950faf9200377108a0aab55fbbfab731d5f0b5fbcd10033`
and `276efebfe75500789ce6501630e8ecdb0dbbf8241ae20c74f57bccfbc348180c`.

## Tools And Execution Controls

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Contained Go 1.18.10 values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Official golangci-lint 2.12.2 reproduces archive SHA-256
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  Pinned apidiff version `v0.0.0-20260709172345-9ea1abe57597` was rebuilt and
  used successfully. Current proxy/GitHub source archives do not reproduce the
  previously recorded apidiff archive digest `0c55d9e...`; keep this as an
  external portable-receipt discrepancy.
- Use `LC_ALL=C LANG=C` because the inherited `C.UTF-8` is unavailable and can
  panic Perl/Go helpers. Complete preflight used scratch-only `make`/`mktemp`
  wrappers so every bare temporary directory stayed below scratch. Keep real
  Python 3.14 ahead of `/usr/bin/python3` for Docker timestamps with long
  fractional seconds.

## Next Objective

Record the user's explicit option 1 decision without dependency implementation.
Revalidate unchanged v0.2.1/Viper selection, negative `go mod why -m`, untouched
metadata, and zero loaded target and GAX packages. Then answer the decision
archive and prepare exactly one next bounded P7 dependency mission. Do not
execute that successor mission in the decision-recording turn.
