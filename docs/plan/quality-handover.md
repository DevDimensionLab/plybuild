# Quality Upgrade Handover

Generated: 2026-09-09T21:49:49+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation is Google Martian v3.3.2 commit
  `4644476aef78e24d0c7e6a1140baba4c3226a1cc`, parent
  `8508a9d54f066985fb81ebd35fa31f45fa36da40`, tree
  `f20d9b3817fe29af4876745d09c14e9cbcac9152`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions.
- Golang Snappy v1.0.0 `372f8e9`, Golang Protobuf v1.5.3 `6870e02`,
  Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1 `647d4fd`, Go Logfmt
  v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0
  `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain ancestors. Gogo
  Protobuf v1.3.2 remains retained without a dependency edit. All earlier P7
  decisions are final.
- The answered Google Martian archive and sole NEXT Google PProf archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish,
  release, stash, revert, launch a successor, bypass cleanup, or remove the
  worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, and Google Martian v3.3.2 moves. Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, and historical root GLFW remain retained. P8 remains queued. Do
not reopen earlier groups or combine another dependency group.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## Google Martian Decision And Identity

Upgrade exact-path `github.com/google/martian/v3` from v3.2.1 to v3.3.2.
V3.3.2 is the highest qualified stable exact-path release. The proxy exposes
exactly v3.0.0, v3.1.0, v3.2.1, v3.3.2, and v3.3.3, with no prerelease,
retraction, or deprecation. Proxy, sumdb, parent go-import metadata, tagged
Git, and GitHub agree on public, non-fork, Apache-2.0 repository
`https://github.com/google/martian`. It is archived, its default branch is
`master`, and master equals v3.3.3. There are no newer unreleased commits or
qualified alternate paths. Invalid tags/releases `3.2` and `v3.2` do not
version the module.

V3 tags are lightweight and have no signed tag objects. GitHub reports valid
commit signatures for the five exact v3 release commits. The serious
candidates are:

- V3.2.1: commit `7e75073889cd2324f33b959c4fb4545440da046c`, tree
  `b39ce198afdb87dadbac7e9ca36082ed5b8a043e`, parent
  `19163e1b8984e2577ea366d8aa97f0f75c6f6ac6`, at
  2021-05-19T22:06:43Z.
- V3.3.2: commit `b41869697484173d2e936754152cd64e40a046d4`, tree
  `a1680fb4f3f79b23ee2fde53a119d61841393f37`, parent
  `8087be7bb3737e2de9a70007fbc0dc6c74c0fb1c`, at
  2022-03-14T19:34:56Z.
- V3.3.3: commit `0f7e6797a04da412118541344bbe0d65945e24c9`, tree
  `0ad4a12b3b27833d80140eb0a54c2c1ba8d137b7`, parent
  `d6ef5c8f4bee8c1322bdbeb59505bca21d5d5f0a`, at
  2022-08-16T15:12:57Z.

V3.3.2 source/mod sums are
`h1:IqNFLAmvJOgVlpdEBiQbDc2EwKW77amAycfTuWKdfvw=` /
`h1:oBOf6HBosgwRXnUGWUB05QECsc6uvmMiJ3+6W4l/CUk=`. Selected and
candidate module files are byte-identical and declare Go 1.11. V3.3.3
declares Go 1.18 and removes the generator module. Proxy and Git manifests are
byte-identical for all releases; their ordered SHA-256 values are
`4c092e2a900cf633c6b081e20eb42a692fe9ea14ba61407e856eff8a98f796ac`,
`e37b2abf94dc9614c09f70a10be61d2c59b174f6a6b0f71aea614cee733abc6e`,
`50b3e3cffa269e8fa34ebba7cd748a2f96ffccd97c7c50271767ba434736e03c`,
`3f84285118e5c6daa8d19ea631505d04cd52be080bef1fd73f3a8345ec17fbe6`,
and `5a08b12494eca88f9c98c3dbc096e4c997f37e76988c5f7cc0558c95fdbbf997`.

## Closure, API, Behavior, And Qualifications

Selected and v3.3.2 resolve 37-module standalone source-and-test closures;
latest resolves 36. Complete test loading has 428 entries under Go 1.26.7 and
365 under Go 1.18.10, including 215 module-backed entries and 46 Martian
packages. Actual compilation and execution under contained Go 1.18.10 prove
the complete Go floor.

The 46 packages cover HTTP proxying, request/response modifier chains,
context and modifier groups, filters, headers, cookies, query/body/method
modifiers, HAR, HTTP/2, gRPC, MITM TLS, traffic shaping, helpers, generated
protocol code, and the proxy and marbl commands. Apidiff is empty from v3.2.1
to v3.3.2. The candidate moves global flag parsing out of library
`martian.Init` into the proxy command and recognizes wrapped
`net.ErrClosed`, improving composability and shutdown handling. V3.3.3 adds
`martian/log.Logger` and `SetLogger` but incompatibly removes exported
`martian.Init` and inherited `cmd/marbl -v`, so latest is rejected.

A 387-line independent fixture covers modifier chaining/filtering,
headers/cookies/queries/bodies, context/session isolation, HTTP proxy
lifecycle, MITM certificate SAN/organization/cache properties, opaque CONNECT
passthrough, malformed traffic, HAR duplicate/orphan/export/reset/gzip,
failing reads, and concurrency. Its behavior and Init-contract SHA-256 values
are `d41342168319236bc936b4e4977db4460fd0c754195f955d1ad5342c8a21bd96`
and `6262ef23586fd72e15f430ef8f2b0f83efe8dbba940864b13f40ec4d49404c10`.
Candidate race count-10 passes under both SDKs.

Retained qualifications are explicit. Proxy close is not idempotent and does
not own the listener; raw CONNECT copy loops can wait for peer shutdown.
Response modifier errors become warnings, upstream failure becomes 502, and
response write errors are logged. MITM uses RSA-2048, random serials, one
reused leaf key, and an unbounded locked per-host certificate cache. HTTP/2
preface forwarding uses one `Read`. gRPC frame length/decompression and HAR
body/cache growth are not application-bounded; a truncated terminal gRPC frame
can return nil. HAR exports are not deep copies. Malformed body ranges can
panic through an out-of-bounds slice, and traffic-shape map iteration is
nondeterministic. These are inherited across all serious releases; callers
must bound hostile inputs and coordinate shutdown.

Module verification, package listing, native tests with package vet disabled,
two independent repeats, race, commands, production vet, and 1,380 cross-build
contexts pass every serious candidate under both SDKs. Full vet reports only
inherited test-only Fatalf-from-goroutine issues and a Go 1.26 RSA formatting
finding. Saturated concurrent repeats expose only proxy startup timeouts and
one traffic-shape sub-millisecond boundary miss; serial repeats pass. Upstream
has 271 tests, one example, no benchmarks, and no native fuzz target.

## Project, Vulnerability, And Quality Measurements

Eight historical `cloud.google.com/go` graph vertices, v0.83.0, v0.84.0,
v0.87.0, v0.90.0, v0.93.3, v0.94.1, v0.97.0, and v0.99.0, each declare
Martian v3.2.1. Requirements of encountered old vertices remain graph edges
even though MVS selects Cloud Go v0.105.0. This explains selection alongside a
negative `go mod why -m` and zero loaded Martian packages. The exact root
edge now selects v3.3.2 and changes no unrelated module selection.

Current measurements are 234 modules, 3,597 graph edges, 429 complete-test
entries, 41 loaded modules, 197 loaded module-backed packages, zero Martian
packages, 1,063 sum lines, and a 418-line unapplied tidy projection. Relative
to accepted go-cmp commit `c314bcb`, sums are +47/-0. The v3.3.3 projection
has 3,599 edges but identical loaded and checksum populations. Main Go 1.18
and preferred toolchain Go 1.26.7 remain unchanged.

Fresh vulnerability data contains 1,392 module records and no Martian record.
Direct candidate results are identical: 27 inherited module findings, 17
inherited package findings, and five inherited IDs/81 traces; none belong to
Martian. Project base/candidate populations are identical: 30 module
findings, 22 Darwin package findings, 23 Windows package findings, and 20
IDs/22 traces on both symbol platforms, with no Martian finding or trace.

Exact Go 1.26.7 verification, load, build, full tests, two repeats, race, vet,
empty-HOME, and Linux/Windows builds pass. The Go 1.18.10 projection passes
verification, load, build, vet, and cross-builds. Full tests retain only the
two accepted `pkg/shell` closed-file error-wording assertions; the compatible
population passes repeats and race.

The exact 21-stage `make quality` exits zero with 27/27 Q0-Q2 PASS at L2,
80/80 mutations killed, seven improvements, and zero held, regressed,
not-comparable, or dirty counts. Scorecard SHA-256 is
`094a7668e43eb0a5b563bcd637b3b5a301e1bb5f5053903ad8934e2b2906ae7f`.
The 509-entry selected-evidence manifest SHA-256 is
`614921201697ce6f3accdf0c6e619c86ddb3502f23c3b03ee6314a2347a9a649`;
decision-summary SHA-256 is
`374c235a5ebeeb7a97254e8eab57b3dc60224fc7f103611beceafff701f5756e`.

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
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`.
  Rebuilt apidiff has exact required module/version identity; govulncheck is
  v1.7.0.

## Next Objective

Independently evaluate exact-path
`github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1` as the next
single P7 group. Do not combine Cloud Go, Demangle, Readline, Logex, X/Sys, or
another dependency.

The selected pseudo-version exists through exact graph edge
`cloud.google.com/go v0.90.0 -> github.com/google/pprof`, despite selected
Cloud Go v0.105.0. `go mod why -m` is negative and no PProf package is
loaded. The initial proxy list is empty. Selected is dated
2021-07-20T18:47:32Z, declares Go 1.14, and has source/mod sums
`h1:K6RDEckDVWvDI9JAJYCmNdQXq6neHJOYx3V6jnqNEec=` /
`h1:kpwsk12EmLew5upagYY7GY0pfYCcupk39gWOCRROcvE=`.

Proxy latest is `v0.0.0-20260906184651-6331bc6350fe`, dated
2026-09-06T18:46:51Z at origin commit
`6331bc6350fe55a6fec2957299e0581dd7510e36`, with source/mod sums
`h1:QAinXoAFJdGQYztXn3VpFey7KCwpedbZ/EkzbplQ0cY=` /
`h1:jl5iWTm0/hd5PjEYEOuwAJ57L/CibdZfrqZ5XA5GrCk=`. Its Go 1.25.0
directive is floor-ineligible. Treat every survey fact as incoming evidence to
verify, and do not infer permission to choose an arbitrary pseudo-version from
the repository's lack of semver tags.
