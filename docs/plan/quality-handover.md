# Quality Upgrade Handover

Generated: 2026-09-13T04:13:50+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation remains Google UUID v1.4.0 commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions. Google Cloud Go
  Testing and GAX v2.7.0 were retained without dependency commits or metadata
  edits.
- Google Renameio v1.0.1 `394ec36`, Google Martian v3.3.2 `4644476`,
  Golang Snappy v1.0.0 `372f8e9`, Golang Protobuf v1.5.3 `6870e02`,
  Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1 `647d4fd`, Go Logfmt
  v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0
  `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain ancestors. Google
  pprof, Gogo Protobuf, GAX, and Google Cloud Go Testing remain retained
  without dependency edits. All earlier P7 decisions are final.
- The answered Google Cloud Go Testing archive and sole NEXT Enterprise
  Certificate Proxy archive link reciprocally. No `.agent-task/current.md` or
  repository `.quality/manual-evidence.json` exists. Do not push, merge,
  publish, release, stash, revert, launch a successor, bypass cleanup, or
  remove the worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, Google Renameio v1.0.1, and Google UUID
v1.4.0 moves. Google pprof
`v0.0.0-20210720184732-4bb14d4b1be1`, Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, historical root GLFW, Googleapis GAX Go v2 v2.7.0, and Google
Cloud Go Testing
`v0.0.0-20200911160855-bcd43fbb19e8` remain retained. P8 remains queued.
Do not reopen earlier groups or combine another dependency group.

The user's two 2026-09-12 GAX decisions are final. Retain v2.7.0 as an
inherited Go-floor exception and accept its known
`apierror.ParseError(err, false)` panic only while the complete project load
contains zero GAX packages. If any GAX package becomes loaded or directly
imported, the behavior exception expires and the owning checkpoint must stop
for a fresh dependency and product decision before merge. Do not patch or
replace GAX, add a root edge, change Viper or another parent, raise the Go
floor, or apply the GAX exception to another module.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## Google Cloud Go Testing Decision And Identity

Retain exact-path `github.com/googleapis/google-cloud-go-testing
v0.0.0-20200911160855-bcd43fbb19e8` without changing `go.mod` or
`go.sum`. No qualified exact-path release exists, proxy latest is an
unreleased archived branch head, and selected has no independent
disqualifier.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence agrees on the public,
archived, non-disabled, non-fork Apache-2.0 repository
`https://github.com/googleapis/google-cloud-go-testing`. It has zero tags and
zero GitHub releases. `RELEASING.md` explicitly says there are no releases and
consumers should use pseudo-versions. The stable proxy list is empty. There is
no redirect, deprecation directive, retraction, alternate exact-path release,
qualified fork, prerelease, or tag to promote.

Selected is branch `bcb-to-fb` tip commit
`bcd43fbb19e8d79524fce1b71f4a2145afbd6039`, parent
`8e1d251e947d1de4242ca2a81eb1a39917bf6b38`, tree
`44cc1490ddeb164a479884cd8beecc373b2ad60d`, dated
2020-09-11T16:08:55Z. Proxy latest is protected `master` tip
`1c9a4c676720af1d2d964c1b2c866ee500300a15`, parent
`1487aa9ec5b057debb42f236a1d1185f09f96804`, tree
`e1733414a158ce2090bf6eb073c4c6af8dceee52`, dated
2021-07-19T22:17:36Z. Selected/latest source/mod sums are
`h1:tlyzajkF3030q6M8SvmJSemC9DTHL/xaMa18b65+JM4=` /
`h1:dvDLG8qkwmyD9a/MJJN3XJcT3xFxOKAvTZGvuZmac9g=` and
`h1:zC34cGQu69FG7qzJ3WiKW244WfhDC3xxYMeNOX2gtUQ=` / the same mod sum.
Strict/full fsck and ancestry pass, and GitHub verifies both commit signatures.
Selected is two commits behind latest;
only README archive wording and SECURITY.md differ. Proxy/Git manifests match.

## Closure, API, Behavior, And Qualifications

Selected/latest module files are identical, declare Go 1.11, and require
Cloud Go v0.44.3, BigQuery v1.0.1, Datastore v1.0.0, Google API v0.9.0,
and tools-only vertices. Under Go 1.26.7 each isolated graph contains 36
modules/337 edges and 332/359 production/test package entries. Under Go
1.18.10 each has 36 modules/336 edges and 270/296 entries. The actual imported
closure uses 14 external modules and five own runtime packages. Its maximum
Go directive is 1.11; the whole graph including tools peaks at 1.12.

The root documentation package and `bqiface`, `dsiface`, `psiface`, and
`stiface` expose adapters, constrained interfaces, and shallow config shadows
for the pinned Cloud clients. There are no runtime commands, benchmarks, fuzz
targets, generated Go files, or testdata; `tools.go` is build-tag-only. Four
credentialed examples skip without their environment variables. Upstream
Kokoro targeted Go 1.12 with lint/staticcheck/tidy/race checks.

Adapters are value wrappers around embedded client pointers and normally
forward exact contexts, arguments, callbacks, errors, and results. Nil client
constructors return non-nil interfaces with nil embedded pointers, and later
calls generally panic. Successful nil upstream results can also become
non-nil interface wrappers. BigQuery config conversion shallow-aliases caller
data; several iterators drop results on error. Datastore can return wrappers
with errors and preserves callback/context identity. Pub/Sub message data and
attributes alias upstream storage, and Receive wraps messages while forwarding
the exact context. Storage builder setters mutate the embedded builder.
BigQuery copy, Pub/Sub publish, and Storage copy/compose assert concrete local
wrappers and panic for otherwise valid external fakes or typed nils. Mutable
builders/iterators are not concurrent-mutation safe; normal independent and
read-only reuse passes race tests. The wrappers add no retry, clock,
randomness, or network policy.

Selected/latest exported API blobs are identical, SHA-256
`1a9dbaa10ba5d79c1b207cd5acaefe6c1903c3ecceb868c2fe6d8605370a2e00`,
and apidiff is empty both ways. The independent 375-line contract fixture,
SHA-256
`10bea668cbc7b6ec438f6fa15ff7b8681e6901ef4c4076f8909b052abe5d991c`,
covers forwarding, identity, aliasing, nil/panic boundaries, paging,
cancellation, callbacks, local Pub/Sub streaming, builders, errors, and
concurrency. It passes test, count-10 repeats, race, and vet under both SDKs.

Three upstream `Example_AdaptClient` names are malformed for modern vet. Raw
Go 1.26.7 default tests/vet and Go 1.18.10 standalone vet report only this
test-name defect; vet-disabled runtime tests pass. Scratch-only correction to
`ExampleAdaptClient` makes default tests and vet pass under both SDKs. Corrected
production cross-builds pass Darwin/amd64, Linux/amd64, Linux/arm64,
Windows/amd64, FreeBSD/amd64, and js/wasm, with Linux/Windows test compilation.
This historical example-name issue is not a runtime disqualifier.

## Project, Vulnerability, And Quality Measurements

Afero v1.9.4 supplies the only current incoming edge and shortest path main ->
Afero -> selected. Historical Afero v1.8.2 declares the same selected
pseudo-version. `go mod why -m` is negative; zero target and zero GAX packages
load. Raw explicit selected/latest requests change no loaded population or
unrelated selection. Both temporarily add root requirements for the target and
five tools vertices. Selected adds six pruned-tool source sums and raises edges
to 3,605; latest changes only target, adds seven source sums, and raises edges
to 3,612. Tidy makes base/selected/latest byte-identical, restores selected,
and measures 234 modules/3,557 edges/948 sum lines.

The unchanged accepted project remains 234 modules, 3,599 graph edges, 429
complete-test entries, 41 loaded modules, 197 loaded module-backed packages,
1,067 sum lines, and a 432-line unapplied tidy projection.

Fresh primary vulnerability data contains 1,398 module records, index
Last-Modified 2026-09-10T16:28:28Z, scanner DB update
2026-09-10T14:48:42Z, and no exact target advisory. Project populations remain
30 module findings, 22 Darwin and 23 Windows package findings, and 20 IDs/22
traces per symbol platform, with no target or GAX occurrence. Direct
selected/latest scans are identical: 31 inherited module findings, 18 package
findings/16 IDs, and 51 production or 66 test traces/three IDs. Those three
IDs affect old `x/net` HTTP/2 and gRPC transport, not the target module.

Exact Go 1.26.7 project verify/load/build, native tests, two repeats, race,
vet, pinned lint, empty-HOME, Linux/Windows, compatibility, complete preflight,
80/80 mutation controls, audits, host, snapshot, and real Docker acceptance
pass. Applicable Go 1.18.10 gates pass with only the two accepted shell
closed-file wording differences. The fresh independent scorecard is 27/27
Q0-Q2 PASS at L2, seven improved and zero held/regressed/non-comparable/dirty,
SHA-256
`a706ee72aa1483217bdcced30b36523da95841fb3874ef017645ffcf74051829`.
Because metadata did not change, the accepted scorecard remains
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
The 788-entry evidence manifest and decision summary hash to
`34e3574d2ced8aef523df1607300134664907cf7634cb40f89c3f615a151f106`
and `9652c69c7971189749d2fd4590dc39b49a11f6f159cb80c2287767251312db2a`.

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
  accepted apidiff archive SHA-256 is
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
  The apidiff mirror refresh returned HTTP 429; exact pinned module source was
  rebuilt from cache.
- Complete preflight used a scratch-only bare-`mktemp` shim. Docker acceptance
  used scratch-discovered Buildx v0.33.0 and real Python 3.14.6 ahead of
  `/usr/bin/python3` for long-fraction timestamps.

## Next Objective

Independently evaluate exact-path
`github.com/googleapis/enterprise-certificate-proxy v0.2.1` as the next
single P7 group. Do not combine or independently audit Viper or another
module.

Viper v1.15.0 supplies the sole incoming edge and shortest path main -> Viper
-> target. `go mod why -m` is negative and zero target packages load. A fresh
proxy survey found 29 stable versions: selected v0.2.1 is tagged commit
`80592736477602cc7992372d4280f819dc7e4cbf`, declares Go 1.19, and has no
inherited floor exception. V0.1.0 and v0.2.0 declare Go 1.18; v0.2.1 through
v0.3.4 declare Go 1.19; later releases require Go 1.23 or newer. Latest
v0.3.22 is tagged commit `62d25fa2858321169ff476109479205309fa3a68`
and declares Go 1.25.0/toolchain Go 1.26.5.

Resolve the complete v0.2.0 closure, identity, API, behavior, and exact MVS
feasibility. If the highest eligible release is sound but cannot be selected
without changing Viper or another group, stop for a fresh bounded product
decision. Do not silently retain the Go-1.19 target, invent an exception, add
a replace/fork, or broaden the checkpoint.
