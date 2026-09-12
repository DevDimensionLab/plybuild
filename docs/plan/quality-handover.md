# Quality Upgrade Handover

Generated: 2026-09-13T00:53:10+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation remains Google UUID v1.4.0 commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions. GAX v2.7.0 was
  retained without a dependency commit or metadata edit.
- Google Renameio v1.0.1 `394ec36`, Google Martian v3.3.2 `4644476`,
  Golang Snappy v1.0.0 `372f8e9`, Golang Protobuf v1.5.3 `6870e02`,
  Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1 `647d4fd`, Go Logfmt
  v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0
  `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain ancestors. Google
  pprof, Gogo Protobuf, and GAX remain retained without dependency edits. All
  earlier P7 decisions are final.
- The answered GAX archive and sole NEXT Google Cloud Go Testing archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish,
  release, stash, revert, launch a successor, bypass cleanup, or remove the
  worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, Google Renameio v1.0.1, and Google UUID
v1.4.0 moves. Google pprof
`v0.0.0-20210720184732-4bb14d4b1be1`, Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, historical root GLFW, and Googleapis GAX Go v2 v2.7.0 remain
retained. P8 remains queued. Do not reopen earlier groups or combine another
dependency group.

The user's two 2026-09-12 GAX decisions are final. Retain v2.7.0 as an
inherited Go-floor exception and accept its known
`apierror.ParseError(err, false)` panic only while the complete project load
contains zero GAX packages. If any GAX package becomes loaded or directly
imported, the behavior exception expires and the owning checkpoint must stop
for a fresh dependency and product decision before merge. Do not patch or
replace GAX, add a root edge, change Viper or another parent, raise the Go
floor, or ask for either decision again.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## GAX Decision And Identity

Retain exact-path `github.com/googleapis/gax-go/v2 v2.7.0` without changing
`go.mod` or `go.sum`. No higher qualified stable exact-path release preserves
Go 1.18, and every exact root alternative violates the bounded MVS scope.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence agrees on the public,
active, unarchived, non-fork BSD-3-Clause repository
`https://github.com/googleapis/gax-go` and exact `/v2` submodule. The proxy
lists 42 v2-looking stable tags, but v2.0.0 lacks `v2/go.mod` and cannot
version this path; 41 valid stable modules remain through v2.24.1. There is no
deprecation, retraction, redirect, qualified fork, alternate-path release, or
prerelease. Stable tags are lightweight and their commits are GitHub-
verified. Main has three unreleased commits after latest and was not promoted.

Selected v2.7.0 is commit
`2592e2286ac291a2d9e7c06e1f31b3a6d772131c`, parent
`9dbd96d59b9d54ceb7c025513aa8c1a9d727382f`, tree
`f4e76bc2c50c9078de2e08237977a9f60b03934f`, at
2022-11-02T20:02:53Z. Latest v2.24.1 is commit
`269185f57eafcc619f159ffe8857eee6073e955e`, parent
`d0623a1d3c843a64e85ad8b4d4ea37695a668868`, tree
`57ab5ad8af3a8c6027ac9ae0b8b729f1a37c9b1c`, at
2026-09-03T20:21:21Z. Selected/latest source/mod sums are
`h1:IcsPKeInNvYi7eqSaDjiZqDDKu5rsmunY0Y1YupQSSQ=` /
`h1:TEop28CZZQ2y+c0VxMUmu1lV+fQx57QpBWsYpwqHJx8=` and
`h1:AtqTN21IXMMWo99LiEVAiBfNNQmO40d8xUfZI640mc0=` /
`h1:bWeBei0NVwaNZKb2y1HUBS7gLXIF3/Tu3pq7j8D2Tb0=`. Strict fsck and
release ancestry pass. Proxy/Git bytes match, including latest after
documented module-zip LICENSE copying and mode normalization.

## Closure, API, Behavior, And Qualifications

V2.5.1 is the last release declaring Go 1.18. Every v2.6.0-or-newer stable
declares at least Go 1.19; v2.12.5 declares Go 1.20 and v2.24.1 declares Go
1.25. Exact graph sizes for v2.5.1, v2.6.0, v2.7.0, v2.12.4, v2.12.5, and
v2.24.1 are 89, 142, 84, 45, 46, and 178 modules; production/test package
populations are 301/334, 301/334, 301/342, 301/345, 301/345, and 361/419.
Selected's imported production/test closure contains 8/9 modules and three
Go 1.19 declarations: GAX, Google API v0.102.0, and Genproto.

All six serious versions pass exact Go 1.26.7 verification, native repeats,
race, vet, six production cross-builds, and Linux/Windows test cross-builds.
V2.5.1, v2.6.0, and v2.7.0 also load/test on Go 1.18.10; deep selected and
v2.5.1 gates pass. V2.12.4/v2.12.5 fail on Go 1.18 in the gRPC generic-
atomic closure and latest is parser-ineligible. Selected's build success does
not erase the declared floors; it remains only through the authorized
inherited exception.

Selected has public runtime packages `gax` and `apierror`, internal helpers
and generated protos, examples/testdata, and no command, benchmark, or fuzz
target. Upstream Kokoro omits `apierror` from its root test command. Apidiff
finds later compatible additions apart from the Version constant value;
v2.5.1 removes selected API.

`Invoke` applies later options last, calls once before checking cancellation,
creates one retryer lazily, adapts gRPC/HTTP errors, and has no implicit
timeout or attempt cap. Sleep is cancelable. `OnCodes` with selected gRPC does
not classify `%w`-wrapped statuses; `OnHTTPCodes` recognizes wrapped HTTP
errors. Backoff is mutable full jitter over global `math/rand`, can exceed Max
on the first pause, panics for negative Initial, and caps overflow. Retryer/
Backoff instances are not safe for concurrent sharing; fresh-per-Invoke
concurrency passes.

`APIError` preserves originals through `FromError`/`Unwrap`, parses standard
details, retains unknown details, and exposes shallow detail/metadata views.
Malformed detail payloads are ignored. Selected `ParseError(err, false)` can
produce an `APIError` whose `Error()` panics. Repair commit `22c16e7bff` first
shipped in v2.12.5, whose Go 1.20 floor is ineligible. This is the bounded
zero-load behavior exception. Other inherited qualifications include odd-
argument `XGoogHeader` panic, nil content-reader panic, at-most-512-byte
content sniff/replay, non-concurrent ProtoJSONStream, and non-idempotent stream
Close.

The independent 499-line contract fixture SHA-256 is
`884faf82e80a7589977fd613e767a990610ce0b9676e47043c9119525880bb5f`.
It covers option order, retry/no-retry, cancellation/deadlines, error
conversion/details/wrapping, backoff/timers, headers/content, stream
lifecycle, malformed/nil/panic boundaries, and concurrent Invoke. Selected
passes test, two count-10 repeats, race, and vet under both SDKs with the
expected panic assertion; v2.12.5 passes on Go 1.26.7 with repaired behavior.

## Project, Vulnerability, And Quality Measurements

Viper v1.15.0 supplies shortest path main -> Viper -> GAX v2.7.0.
Historical Cloud/Google/Viper/Crypt vertices retain v2.0.4-v2.1.1. MVS keeps
these graph edges without a package import; `go mod why -m` is negative and
exactly zero GAX packages load. Exact v2.5.1 downgrades Viper and many
unrelated selections; v2.12.5/latest broadly upgrade unrelated selections,
and latest raises the main Go line. A redundant selected-version root edge is
pointless. Every projection keeps the loaded GAX count at zero.

The unchanged project remains 234 modules, 3,599 edges, 429 complete-test
entries, 41 loaded modules, 197 loaded module-backed packages, 1,067 sum
lines, and a 432-line unapplied tidy projection. The main module stays Go
1.18 with toolchain Go 1.26.7. Relative to go-cmp commit `c314bcb`, sums
remain +51/-0.

Raw v2.5.1/explicit-v2.7.0/v2.12.5/latest projections measure module/edge/
sum counts as 234/4,088/1,158, 346/3,777/1,077, 248/3,674/1,084, and
260/3,763/1,106. All retain 429 packages, 41 loaded modules, 197 module-
backed packages, and zero GAX packages. The base tidy control is
234/3,557/948; tidied candidates are 231/3,555/950, byte-identical base,
234/3,561/949, and 234/3,582/958. Their module-list/go.mod/go.sum deltas
against base tidy are 136/19/22, 0/0/0, 29/13/24, and 85/58/128 lines.

Fresh primary vulnerability data has 1,398 module records, index Last-
Modified 2026-09-10T16:28:28Z, and no exact GAX record. Project populations
remain 30 module findings, 22 Darwin and 23 Windows package findings, and 20
IDs/22 traces per symbol platform, with no GAX occurrence. Direct selected-
source scans record 24 module findings, 14 package findings/11 IDs, and 23
traces/three IDs inherited from `x/net` HTTP/2 and protobuf protojson. These
do not enter the unchanged project, but reinforce the future-loading stop.

Exact Go 1.26.7 project verify/load/build, native tests, two repeats, race,
vet, pinned lint, empty-HOME, Linux/Windows, compatibility, complete preflight,
80/80 mutation controls, audits, host, snapshot, and real Docker acceptance
pass. Applicable Go 1.18.10 gates pass with only the two accepted closed-file
wording failures. The unchanged scorecard remains all 27 Q0-Q2 rows PASS at
L2, seven improved, and zero held/regressed/non-comparable/dirty, SHA-256
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
The 863-entry GAX evidence manifest and decision summary hash to
`8b8ea0d9b842b8b3eecc1d7c7b8f8754c24662a74dd21d9d1cef0304bbd18b5b`
and `26a2a6a7ad4d06a2c059959d7eaf8df64ee9baeb1688c03927b08b1fd650bfe1`.

## Tools

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
- Contained Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Official golangci-lint 2.12.2 archive SHA-256 is
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  GoReleaser 2.17.1 binary SHA-256 is
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`;
  the accepted apidiff archive SHA-256 is
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
  The apidiff mirror was HTTP 429 on final refresh; exact built source identity
  is `golang.org/x/exp v0.0.0-20260709172345-9ea1abe57597` commit
  `9ea1abe5759723e90ff9d6e93de05faedacb461a`.
- Complete preflight used a scratch-only bare-`mktemp` shim. Docker acceptance
  used Buildx v0.33.0 and real Python 3.14.6 ahead of `/usr/bin/python3` for
  timestamps with more than six fractional digits.

## Next Objective

Independently evaluate exact-path
`github.com/googleapis/google-cloud-go-testing
v0.0.0-20200911160855-bcd43fbb19e8` as the next single P7 group. Do not
combine Afero, Cloud Go, BigQuery, Datastore, Pub/Sub, Storage, Google API, or
another group.

Afero v1.9.4 declares selected and supplies shortest path main -> Afero ->
target; historical Afero v1.8.2 declares the same selected edge. `go mod why
-m` is negative and zero target packages load. The stable proxy list is empty.
Proxy latest is `v0.0.0-20210719221736-1c9a4c676720`; selected/latest are
dated 2020-09-11/2021-07-19, and both module files declare Go 1.11 with the
same requirements. Resolve whether the repository has any qualified exact-
path release; do not promote proxy latest or an arbitrary branch head merely
because it is newer. If no newer release qualifies and no independent
disqualifier exists, retain selected without metadata changes or a dependency
commit.
