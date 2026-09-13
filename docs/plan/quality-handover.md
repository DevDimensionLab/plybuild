# Quality Upgrade Handover

Generated: 2026-09-13T17:27:33+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation remains Google UUID v1.4.0 commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
  `37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
  `b89afd4ec056133b1eefb5611f8f35c12c11824b`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions. GopherJS has no
  dependency implementation or metadata commit.
- Google Renameio v1.0.1 `394ec36`, Google Martian v3.3.2 `4644476`,
  Golang Snappy v1.0.0 `372f8e9`, Golang Protobuf v1.5.3 `6870e02`,
  Godbus D-Bus v5.1.0 `6472dce`, Go Stack v1.8.1 `647d4fd`, Go Logfmt
  v0.6.0 `3d4cfd`, Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0
  `e5d6252`, and Speakeasy v0.2.0 `41f9561` remain ancestors. All earlier
  decisions are final.
- The GopherJS evaluation is answered with a bounded product stop. The sole
  NEXT archive asks for that product direction. No `.agent-task/current.md`
  or repository `.quality/manual-evidence.json` exists. Do not execute the
  decision's successor, push, merge, publish, release, stash, revert, bypass
  cleanup, or remove the worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active. Exact Go 1.26.7 and the accepted
Speakeasy, XXHash, Fatih Color, Go Logfmt, Go Stack, Godbus D-Bus, Golang
Protobuf, Golang Snappy, Google Martian, Google Renameio, and Google UUID
moves remain final. Google pprof, Gogo Protobuf, Crypt, OpenCensus Proto,
Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify, Ghodss YAML, historical
root GLFW, Googleapis GAX Go v2 v2.7.0, Google Cloud Go Testing, and Enterprise
Certificate Proxy v0.2.1 remain retained. P8 is queued. Do not combine another
dependency with the GopherJS decision or begin P8.

The user's 2026-09-13 Enterprise Certificate Proxy option 1 decision remains
final and separate. Exact v0.2.1 retains its bounded Go-1.19 floor and recorded
release/behavior/safety exceptions only while it has zero loaded packages,
remains runtime-unreachable, and its exact version and sole Viper v1.15.0 edge
remain unchanged without a new advisory or independent disqualifier. The two
GAX exceptions likewise remain valid only for exact v2.7.0 while zero GAX
packages load. This evaluation revalidated zero packages for both targets; do
not transfer either exception to GopherJS.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, runtime/tool installation, and build context beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never run `go mod download all` in a measured
worktree or bypass launcher cleanup.

## GopherJS Product Stop

No qualified exact-path release satisfies the existing Go 1.18, release,
complete-closure, host, and behavior contracts. Do not change
`github.com/gopherjs/gopherjs` without a fresh bounded product decision.

Fresh proxy, sumdb, `go-import`, strict Git, and GitHub evidence resolve the
public, active, unarchived, non-fork BSD-2-Clause repository
`https://github.com/gopherjs/gopherjs.git`. The proxy lists thirteen versions.
GitHub publishes stable releases v1.17.2, v1.20.0-v1.20.2, and v1.21.0;
v1.18.0-beta1 through beta3 and v1.19.0-beta1/beta2 remain semver
prereleases even though GitHub's boolean prerelease field is false. Base tags
and `+go` tags for v1.17.2 and v1.18.0-beta3 point to the same commits. Tags
are lightweight and unsigned; the newer merge commits have valid GitHub
commit signatures. There is no retraction, deprecation, redirect, alternate
module path, qualified fork, or promotable unreleased branch head.

Selected pseudo-version `v0.0.0-20181017120253-0766667cb4d1` is unsigned
commit `0766667cb4d1cfb8d5fde1fe210ae41ead3cf589`, parent
`1babbf986f6fcb1156d0646cdba5c4f81bc32849`, tree
`dd58a45631ba0db649f1b0a69b403d92717bd7b3`, dated
2018-10-17T12:02:53Z. It is an ancestor of current master. The proxy created
only a synthetic module file containing the module path; source/mod sums are
`h1:EGx4pi6eqNxGaHF6qqu48+N2wcFQ5qg5FXgOdqsJ5d8=` and
`h1:wJfORRmW1u3UXTncJ5qlYoELFm8eSnnEO6hX4iZ3EWY=`.

Highest floor-eligible stable v1.17.2 is unsigned commit
`fcf8e05a6f4fe7573b43c6bc65c2d1166fbd48cd`, parent
`3f8f90ca77315ddcfeb3fdae8db0262fcc86a758`, tree
`1e5857086864288b4a1cb1158384437da435b59e`, dated
2022-04-19T15:05:20Z. Its source/mod sums are
`h1:fQnZVsXk8uxXIStYb0N4bGk7jeyTalG/wsZjQ25dO0g=` and
`h1:pRRIvn/QzFLrKfvEz3qUuEhtE/zLCWfreZ6J5gM2i+k=`. Highest Go-1.18
candidate v1.18.0-beta3 is prerelease merge commit
`f1c1c698f546014a73f7d3f748201ddbe153c02f`. Stable v1.20.2 and latest
v1.21.0 are valid releases but declare Go 1.20 and Go 1.21 respectively.
Current master `490705b1d6fc7d5bd9202ac41888e146183328eb` is unreleased.

## Floor, Packages, API, And Behavior

The selected proxy source has 16 root packages but no dependency requirements.
Read-only resolution reports twelve missing providers, and the zip excludes
the vendored test package. Synthesizing its missing module metadata under Go
1.18 selects current dependencies whose `x/tools` closure imports post-1.18
standard packages; under either SDK the complete test closure still cannot
load. Its compiler also deliberately fails outside Go 1.11, so source API
export under Go 1.26.7 fails on the undefined
`___GOPHERJS_REQUIRES_GO_VERSION_1_11___` guard. It is neither a qualified
release nor a reproducible modern source/test closure.

V1.12.80 has explicit requirements but deliberately supports only Go 1.12;
Go 1.18 compilation fails both that version guard and the old
`testing.MainStart` closure. V1.17.2 declares Go 1.17. Its complete graph has
157 modules and peaks at Go 1.17; actual production/test loading uses 29
external packages across 14 external modules. Complete test loads contain 308
entries under Go 1.26.7 and 238 under Go 1.18.10. V1.18.0-beta3 declares Go
1.18, uses 33 external packages across 14 modules, and loads 318/247 entries.
Thus v1.17.2 is the highest stable release whose declared and imported closure
preserves the Go 1.18 floor, but floor qualification alone is insufficient.

V1.17.2 contains 172 Go files, 64 test files, 211 test functions, two
benchmarks, no fuzz targets, two examples, four generated Go files, and 103
build-constrained Go files. Its module exposes the root `gopherjs` command;
`build`, `build/cache`, `build/versionhack`; compiler analysis, AST, filter,
native, prelude, package, and type helpers; `js`; `nosync`; internal helpers;
and test packages. The CLI subcommands are `build`, `install`, `doc`, `get`,
`run`, `test`, `serve`, `version`, and `clean`. Malformed arguments, invalid
log levels, missing packages, and wrong target SDKs return exit 1.

The command requires an exact matching target distribution through
`GOPHERJS_GOROOT` (Go 1.17 for v1.17.2), uses the host default GOOS with only
linux/darwin supported, emulates 32-bit GOARCH `js`, and does not support cgo.
Build/install write JavaScript, source maps, cache entries, and GOPATH/GOBIN
artifacts; watch mode monitors files indefinitely. `get` and `doc` spawn the
Go tool, `run`/`test` spawn Node and use temporary files, and `serve` binds a
TCP listener without graceful shutdown. There are no contexts, input caps, or
command timeouts. Cache failures can be swallowed. Browser builds replace
filesystem, process, network, and runtime behavior with native shims; Node can
load the optional native syscall addon. No browser executable was available,
so browser behavior is source-characterized rather than claimed as an
executed browser pass.

Public `build.NewSession(nil)` panics. An independent fixture proves that a
valid call mutates caller-owned `Options` by filling GOPATH and forcing Verbose
for watch mode; the session maps/watcher and mutable caches are not a
concurrent public abstraction. `PackageData` derivatives reuse caller slices.
The `js` package relies on compiler intrinsics and has nil/malformed panic
boundaries; `nosync` intentionally panics on contention or misuse for the
single-threaded JavaScript model. A two-test session fixture passes count-10
and race under both host SDKs; its SHA-256 is
`6a75e8d98fa272673642d59800d5141c8c6551f2c037e75f800f5a292ecb02a0`.

Pinned apidiff can export v1.17.2 but cannot export selected. V1.17.2 to beta3
has incompatible public changes including `Options.GOROOT/GOPATH` removal,
`ImportCError` removal, `PackageData.JSFiles` type change, `XContext.GOOS`
removal, and native FS type change. Beta3 to v1.20.2 removes/reworks broad
compiler/build APIs and packages; v1.21.0 changes `WriteProgramCode`. Export
data SHA-256 values for v1.17.2, beta3, v1.20.2, and v1.21.0 are respectively
`5e1e34cedd45d06bac2f99201a4e5ab027a0e102b5114b3e1604ec5fcb25c63a`,
`8ec59b220fafb8216351b1acc220c754e3b51066892568675e8079a89e88d1af`,
`ae9b4943d192ef7079a207a0addefd02a16615e7821777d2f7ea1b93ce8c9947`,
and `3eb2a4567cec8df33769b95843dd92b4e5f862a09edb727eb65887bac5e16d5c`.

## Candidate Test Result

With exact Go 1.18.10 as host, contained Go 1.17.9 as target, upstream Node
12.22.12, and the native syscall addon, v1.17.2's native packages and main
GopherJS suite pass; its Go-repository compiler suite reproducibly records
421 pass, 62 known failures, and seven unexpected failures. All seven are
Darwin/arm64 generated-program crashes through
`syscall/zsyscall_darwin_arm64.go` into unimplemented
`internal/abi.FuncPCABI0`. A trivial independent goroutine/formatting fixture
also compiles but crashes at first output, and `gopherjs test` crashes on the
same path. Exact Go 1.26.7 reproduces those failures and additionally panics
on newer standard-library generic AST syntax and misparses Go 1.26's VERSION
file. Both required hosts therefore reject v1.17.2. Its upstream-supported vet
scope and race scope excluding the compiler repository suite pass; whole-tree
vet reports two deliberately unreachable test blocks.

For a Linux target, two independent builds produce byte-identical JavaScript
and source maps. Node prints `sum=42 goos=linux goarch=js`, and both source-map
files parse as version 3 with 107 sources. The test binary nevertheless exits
140. This limited deterministic generation does not cure the supported-host
Darwin failure. V1.18.0-beta3's complete suite passes under Go 1.18.10, but it
is a prerelease and its command cannot link under required Go 1.26.7 because
`github.com/visualfc/goembed/parser` references removed private symbol
`go/build.parseGoEmbed`. V1.20.2/v1.21.0 are floor-ineligible. No stable or
pseudo-version passes the complete contract.

## Project MVS, Vulnerability, And Quality

The unchanged graph path is main -> direct `mvn-pom-mutator v0.2.3` ->
`goconvey v1.6.4` -> selected GopherJS. The GoConvey edge is the target's only
incoming edge. `go mod why -m` is negative. The project remains at 234 selected
modules, 3,599 graph edges, 429 complete-test entries, 197 module-backed
packages across 41 loaded modules, 1,067 sum lines, and the established
432-line tidy projection. Exactly zero GopherJS, Enterprise Certificate Proxy,
or GAX packages load. `go.mod` and `go.sum` remain unchanged at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`
and `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.

Disposable exact v1.17.2 requests produce 239 modules, 3,627 edges, 1,080
sum lines, and six changed module paths: target plus five added closure
modules. Beta3 produces 238/3,626/1,080 and five changed paths. Both retain
the 429/197/41 load population with zero target packages, and both tidy to the
byte-identical base tidy projection, restoring selected. V1.20.2 produces
239/3,631/1,083, changes twelve module paths, and leaves five unrelated `x/*`
upgrades after tidy. V1.21.0 produces 240/3,636/1,083, changes thirteen paths,
raises the main Go directive to 1.21, and leaves six unrelated `x/*` upgrades
after tidy. The two modern raw requests fail complete project loading because
needed new `x/net` sums are absent until broader metadata work. No request is
an authorized target-only implementation, and no direct edge was added.

Fresh primary vulnerability data has 1,398 module records, index Last-Modified
2026-09-10T16:28:28Z, scanner update 2026-09-10T14:48:42Z, and no exact
GopherJS record. OSV exact-version queries are zero for selected, v1.17.2,
beta3, v1.20.2, and v1.21.0. V1.17.2's loaded package scan reports inherited
GO-2025-4188 in Logrus and GO-2022-0493 in `x/sys/unix`; beta3 reports only
GO-2025-4188. Neither has a reachable symbol or test-symbol trace. The
unchanged project remains at 30 module findings, 22 Darwin and 23 Windows
package findings, and 20 IDs/22 reachable traces on both symbol platforms,
with zero GopherJS occurrences.

Exact Go 1.26.7 project verification/load/build, two independent native test
runs, race, vet, pinned lint, empty-HOME count-2, Linux/Windows builds, API/CLI
compatibility, and authoritative full preflight pass. The Go 1.18.10
projection removes only the unsupported toolchain line and passes verification,
load, build, vet, Linux/Windows builds, two count-10 repeats and race for the
26 compatible packages; full count-1 retains only the two accepted
`pkg/shell` closed-file wording failures. No dependency metadata changed, so a
changed-selection scorecard is inapplicable. Accepted quality remains 27/27
Q0-Q2 PASS at L2 with scorecard SHA-256
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.

## Tools And Next Decision

- Exact Go 1.26.7 binary/archive hashes are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`;
  Go 1.18.10 values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
  Target Go 1.17.9 archive/binary hashes are
  `1f8a0312bdf991d63734c2fd5693b06b053dfb0cce8f204c7c790a50b22cab03`
  and `45bfab07387c803098fa492be4f88409c30db4032155192e2d10d9f01b96c5b8`.
- Node 12.22.12 x64 archive SHA-256 is
  `32927913ed549ce01685a6f9f4697567a64592c7fd1e9a845ac8a10efa1475e6`;
  it ran under Rosetta. Ambient Node 24 failed the legacy node-gyp install and
  is classified as a tool incompatibility, not dependency behavior.
- Official golangci-lint 2.12.2 and GoReleaser 2.17.1 reproduce portable
  hashes `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`
  and `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`.
  Pinned apidiff was rebuilt at SHA-256
  `8609399efbeeed6b9c4d5422cd087601044043f7b198406dc103d55a84c05a8f`;
  retain the known source-archive reproducibility discrepancy.
- Use `LC_ALL=C LANG=C`; inherited `C.UTF-8` is unavailable. Preflight needs
  base API dependency cache seeding, scratch report paths, a scratch-only
  bare-`mktemp` wrapper, and tool paths supplied as environment variables so
  Make's negative contract tests can override them.

The next bounded session must obtain one product direction. Recommended option
1 retains exact selected GopherJS unchanged only under a new, non-transferable
zero-load exception covering its unqualified pseudo-release, incomplete modern
source/test closure, and known command/runtime incompatibility. Option 2 may
authorize a broader `mvn-pom-mutator`/GoConvey parent or removal group to
eliminate the edge. Option 3 may authorize a specified floor/prerelease/fork or
patched-source strategy, including its broader MVS/API/behavior consequences.
Revalidate the target's sole edge, negative why result, zero target/ECP/GAX
loads, checksums, and fresh advisory state before recording a choice. Do not
infer a choice or implement beyond the selected direction.
