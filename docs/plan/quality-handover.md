# Quality Upgrade Handover

Generated: 2026-09-09T13:14:38+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation remains Godbus D-Bus v5.1.0 commit
  `6472dce617eb80484ed022ae8a53cc350c8be6fe`, parent
  `262d97da7a6a50ecc1170bc1a2f33f02732c093e`, tree
  `250d55d4c374958fc9fab703fa7faa931ea4d4f8`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions.
- Gogo Protobuf v1.3.2 is retained without a dependency implementation commit
  or metadata edit. Go Stack v1.8.1 `647d4fd`, Go Logfmt v0.6.0 `3d4cfd`,
  Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0 `e5d6252`, and Speakeasy
  v0.2.0 `41f9561` remain ancestors. All earlier P7 decisions are final.
- The answered Gogo Protobuf archive and sole NEXT Golang Protobuf archive
  must link reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish,
  release, stash, revert, launch a successor, bypass cleanup, or remove the
  worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, and Godbus D-Bus v5.1.0 moves. Gogo Protobuf v1.3.2,
Crypt, OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage,
Fsnotify, Ghodss YAML, and historical root GLFW remain retained. P8 remains
queued. Do not reopen earlier groups or combine another dependency group.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## Gogo Protobuf Decision And Identity

Retain exact-path `github.com/gogo/protobuf v1.3.2`. It is the latest and
highest qualified stable exact-path release. It declares Go 1.15, and its
complete module plus package/test closure preserves Go 1.18. No redirect,
fork, alternate module path, lower vulnerable release, or unreleased commit
was promoted, and adding a redundant explicit root edge was rejected.

The proxy lists exactly eight stable releases, v1.0.0 through v1.3.2, and no
prereleases. Sumdb confirms v1.3.2 source/mod sums
`h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=` /
`h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=`. Its lightweight tag
resolves to unsigned commit
`b03c65ea87cdc3521ede29f62fe3ce239267c1bc`, tree
`56a9801e7f0e1b7da577e89c8be36d978d243a49`, parent
`550e88954e617545f49920b752c154d72abf1d8d`, at
2021-01-10T08:01:47Z. The tagged Git archive and proxy ZIP contain the same
770 byte-identical paths; both manifests hash to
`a259a75ae50454043cfbf66080ceda30be04fe2041bff5de3713eda48574dd5b`.

Go-import, proxy, sumdb, Git, and GitHub agree on public, enabled, unarchived,
non-fork repository `github.com/gogo/protobuf`, default branch `master`.
There are no retractions, exact-path v2 release, or repository redirect. The
release README sought new ownership and current master explicitly marks the
project deprecated. Its three post-v1.3.2 commits do not form a release and
change only README/generator ordering material. Cosmos gogoproto and
PlanetScale vtprotobuf are distinct module paths and have floors above Go
1.18; neither is an exact-path candidate or promoted successor.

V1.3.2 is also the first fixed release for GO-2021-0053/CVE-2021-3121. Its
tag commit adds the generated-unmarshal bounds check that rejects malicious
negative/overflowed skip positions. All earlier exact-path releases are
affected and therefore cannot displace selected.

## Closure, API, Behavior, And Qualifications

V1.3.2 resolves 12 modules under both SDKs. All transitive Go directives are
at most 1.15. It exposes 168 packages, including 12 commands and 110 packages
with tests. Complete `go list -deps -test` has 530 entries under Go 1.26.7
and 485 under Go 1.18.10; all 394 module-backed entries belong to Gogo itself,
so the historical errcheck/gotool/x-tools module graph loads no external
package. All 12 generator commands build under both SDKs. Public API export
for all packages is identical between v1.3.1 and v1.3.2, including internal
packages; generated code advertises `GoGoProtoPackageIsVersion3`.

The runtime covers reflection and generated fast paths for wire marshal,
unmarshal, size, merge, clone, equality, unknown-field discard, extensions,
maps, oneofs, nullable values, custom types, delimited streams, JSONPB, text,
well-known types, registration, and descriptor access. Buffer deterministic
mode sorts map and extension keys where supported. Errors cover nil messages,
required fields, bad wire types, integer overflow, truncation, and excessive
delimited size. There is no general decode recursion-depth limit, so callers
must bound hostile deeply recursive input. Global registration mutation is
intended for generated initialization, not concurrent runtime writes.

Independent fixtures exercise golden wire bytes, round trips, size, merge,
clone/equal, unknown retention/discard, deterministic maps, oneofs,
extensions, JSON/text, nil/zero values, malformed/truncated/overflow input,
delimited maximum size, registration, buffer reuse, and independent-message
concurrency. They pass `-race -count=20` under exact Go 1.26.7 and Go 1.18.10.
Filtered applicable native suites, vet, Linux/amd64, Windows/amd64, js/wasm,
and purego builds pass under both SDKs.

The unmodified upstream suite has four precisely bounded qualifications:

- `protoc-gen-gogo.TestGolden` expects three `testdata/multi/*.pb.go` outputs
  that the release intentionally ignores and omits; Go 1.26 gofmt also changes
  historical generated comment spacing.
- Random `TestStdTypesGoString` can render nil `*time.Time` as `<nil>`, which
  is not a Go expression. Seed 11 reproduces it. Two Go 1.26 repeats pass;
  one of two Go 1.18 repeats observes this inherited randomized defect.
- `proto.TestRace` reproducibly finds a real inherited race under both SDKs
  when binary sizing atomically writes `XXX_sizecache` while JSON reflection
  reads the same generated message. All other race packages and the proto
  package without that exact mixed-shared-message case pass. Concurrent use
  of independently owned messages and buffers passes.
- Vet reports only `test/stdtypes/concurrency_test.go` calling `T.Fatal` from
  a worker goroutine. The corrected applicable vet population passes.

These defects are recorded rather than normalized away. They affect generator
goldens, randomized GoString output, or unsupported mixed concurrent access to
one mutable generated message. No Gogo package or generated Gogo message is
loaded by Ply, and v1.3.2 does not regress them from v1.3.1; the relevant
project and independently owned-message behavior passes.

## Project, Vulnerability, And Quality Measurements

Viper v1.15.0 and etcd/api v3.5.1 declare v1.3.2, while Prometheus
TSDB/Common declare v1.1.1. MVS therefore retains v1.3.2 even though `go mod
why -m` says the main module does not need it and complete package loading
contains zero Gogo packages. Disposable exact `go get ...@v1.3.2` adds only a
redundant indirect root edge and the selected source checksum, changes no
selection, and raises the graph by that one artificial edge; it was not
applied. A v1.3.1 root requirement remains selected at v1.3.2 through MVS.

Current measurements remain 234 modules, 3,586 graph edges, 429 complete-test
entries, 41 loaded modules, 197 loaded module-backed packages, 1,057 sum
lines, and the accepted 400-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, sums are +41/-0. Main Go 1.18, preferred
Go 1.26.7, every unrelated selection, and repository metadata are unchanged.

Fresh vulnerability data contains 1,392 module records. Direct v1.3.1
harnesses produce one GO-2021-0053 module, package, and symbol finding and an
exact reachable trace to `(*unmarshal).Generate`; v1.3.2 produces zero at all
three levels. Project results remain 30 Darwin module findings, 22 Darwin
package findings, and 20 IDs/22 reachable traces for Darwin and Windows,
with no Gogo finding or trace.

Exact Go 1.26.7 project verification, build, count-1, two count-10 repeats,
race, vet, pinned lint, offline load, Linux/Windows build, API/CLI
compatibility, empty-HOME replay, and launcher/Make contract components pass.
The read-only Go 1.18.10 projection removes only the toolchain directive;
verification, build, vet, applicable repeats/race, and cross-builds pass, with
only the two accepted `pkg/shell` closed-file wording assertions failing.
Because selection is inert, changed-selection-only snapshot/Docker/mutation
and full-quality activity was not manufactured. The accepted result remains
27/27 Q0-Q2 PASS at L2 with scorecard SHA-256
`1756c66aecfdd5c6d246c894d9a63630f39dc5dd789532c488e8a3ddbb247bcb`.

The 2,163-entry selected-evidence manifest SHA-256 is
`b92185ef4d63011df510369fb75a8f35a07d02f11c464efeedc6b0739d38f433`;
decision-summary SHA-256 is
`964fa19e8793454bc1e1d2db71da02cf3a5cc7c4b706e0555ce46dc9b6a12deb`.

## Tools

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
- Contained Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- The official golangci-lint 2.12.2 archive remains
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  Rebuilt apidiff has exact required module/version metadata; fresh
  govulncheck is v1.7.0 with database update 2026-09-02T19:12:04Z. Protoc
  3.14.0 supplied generator context only.

## Next Objective

Independently evaluate exact-path `github.com/golang/protobuf v1.5.2` as the
next single P7 group. Do not combine declaring modules, Go CMP, the successor
Google Protobuf module, or any other dependency group.

The initial proxy survey exposes 18 stable releases, four v1.4.0 release
candidates, and latest v1.5.4. Selected v1.5.2 declares Go 1.9 and requires
Go CMP v0.5.5 plus Google Protobuf v1.26.0; latest v1.5.4 declares Go 1.17,
requires Google Protobuf v1.33.0, and marks the module deprecated in favor of
that successor path. MVS selects v1.5.2 through sixteen selected-module
declarations, `go mod why -m` remains negative, and no package is loaded.
Treat these only as incoming survey facts to verify.
