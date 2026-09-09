# Quality Upgrade Handover

Generated: 2026-09-09T11:19:20+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation is Godbus D-Bus v5.1.0 commit
  `6472dce617eb80484ed022ae8a53cc350c8be6fe`, parent
  `262d97da7a6a50ecc1170bc1a2f33f02732c093e`, tree
  `250d55d4c374958fc9fab703fa7faa931ea4d4f8`. It changes only `go.mod`
  and `go.sum`, with three insertions and no deletions.
- Go Stack v1.8.1 `647d4fd71b226fbb1e916b4238b4c7ad87cd7975`,
  Go Logfmt v0.6.0 `3d4cfdbae0a67e757d37022be7eeedaf32c72772`,
  Fatih Color v1.15.0 `6ca672e`, XXHash v2.3.0 `e5d6252`, and
  Speakeasy v0.2.0 `41f9561` remain ancestors. All accepted or retained P7
  decisions are final and documented by their answered archives.
- The answered Godbus archive and sole NEXT Gogo Protobuf archive must link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. Do not push, merge, publish,
  release, stash, revert, launch a successor, bypass cleanup, or remove the
  worktree.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, and Godbus D-Bus v5.1.0 moves. Crypt, OpenCensus Proto,
Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify, Ghodss YAML, and
historical root GLFW remain retained. P8 remains queued. Do not reopen earlier
groups or combine another dependency group.

Keep every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or bypass launcher
scratch cleanup.

## Godbus D-Bus v5 Decision

Upgrade exact-path `github.com/godbus/dbus/v5` from v5.0.4 to highest
qualified stable v5.1.0. It declares Go 1.12, has no requirements, preserves
the complete Go 1.18 floor, and passes every applicable contract. No redirect,
fork, alternate module path, floor-ineligible release, or unreleased commit
was promoted.

The proxy lists eleven tags and no prereleases. V5.0.0 and v5.0.1 are invalid
list-only `/v5` anomalies because those tags declare the pre-v5 module path;
the valid exact-path population is v5.0.2 through v5.2.2. V5.2.0-v5.2.2
declare Go 1.20 and require `golang.org/x/sys v0.27.0`; exact Go 1.18.10
compilation fails on `strings.CutPrefix` and `unsafe.String`, so all three
are floor-ineligible.

Go-import, proxy, sumdb, Git, and GitHub resolve to the public enabled,
unarchived, non-fork `github.com/godbus/dbus` repository with default
`master`. There are no retractions, deprecation notices, redirects, or
alternate nested v5 module. Historical pre-v5 paths remain distinct.

V5.1.0 sums are
`h1:4KLkAxT3aOY8Li4FRJe/KvhoNFFxo0m6fNuFUO8QJUk=` /
`h1:xhWf0FNVPg57R7Z0UbKHbJfkEywrmjJnf7w5xrFpKfA=`. Its lightweight tag
resolves to signed merge commit
`e523abc905595cf17fb0001a7d77eaaddfaa216d`, tree
`af2a399b7f6f0c3a62c435e15a488c46d3dbf536`, parents
`b357b44b7ab3bf9e9b27c906fb31cb622b7a017e` and
`2c3cf657d630429b62f481c930c754a818d98327`, at
2022-02-27T11:53:47Z. GitHub reports the commit signature valid. The proxy ZIP
and tagged Git archive are byte-identical.

## Closure, API, Behavior, And Gates

V5.0.4-v5.0.6 contain complete package/test closures of 102 entries under Go
1.18.10 and 148 under Go 1.26.7. V5.1.0 has 103/149 entries. All use only the
module itself, seven nonstandard closure entries, and three packages:
`dbus`, `introspect`, and `prop`.

API diffs through v5.1.0 contain compatible additions only. The package covers
connections and transports, bus discovery, authentication/negotiation,
messages and signatures, variants, Unix FDs, paths and names, match rules,
exported objects, introspection/properties, signals, calls/replies/errors,
contexts/deadlines, and close state. V5.1.0 adds or inherits validation,
FD-aware codec, no-autostart, property-introspection, address-escaping,
invalid-message, cancellation, variant, and lifecycle fixes.

Independent private-daemon fixtures exercise valid and invalid wire behavior,
dispatch, match/signal delivery, cancellation/deadlines, disconnects,
determinism, and concurrent reuse. V5.0.4-v5.0.6 incorrectly transmit an
already-cancelled call; v5.1.0 does not. V5.1.0 fixtures pass two count-10
repeats and race under both SDKs.

Unmodified Linux tests pass both SDKs. Upstream Darwin tests omit a Darwin
exclusion for their `execCommand` helper, while a source-preserving platform
projection passes. Vet reports two test-only `Fatal`-from-goroutine
findings. A nonce-TCP count-100 stress exposes an immediate-close timing panic
on Go 1.18.10; formal repeats and race pass. Production cross-builds pass
Darwin, Linux, Windows, NetBSD, and OpenBSD; FreeBSD/DragonFly Unix
credentials require cgo. These are classified upstream test/platform/timing
limitations, not loaded Ply behavior.

Project module verification, build, count-1, two count-10 repeats, race, vet,
offline listing, Windows build, pinned lint, API/CLI compatibility, launcher,
empty-HOME, and contained-Go-1.18 gates pass. The Go 1.18 projection removes
only the later toolchain line; its only test failures are the same two
inherited `pkg/shell` closed-file wording assertions.

Exact 21-stage `make quality` exits 0, including 80/80 killed mutations,
host, snapshot, Docker, and audit acceptance. All 27 Q0-Q2 rows PASS at L2
with seven improved and zero held, regressed, not-comparable, or dirty counts.
Scorecard SHA-256 is
`1756c66aecfdd5c6d246c894d9a63630f39dc5dd789532c488e8a3ddbb247bcb`.
Full audit exits expected 1 only for queued Q3.1/Q3.3/Q3.4/Q3.7; its
scorecard SHA-256 is
`b2f36c83e3021f7280dbd26c5816cbc943ce102c1851a71f534bbd2ef669ac38`.

## Project And Vulnerability Measurements

Go-systemd v22.3.2 declares v5.0.4. No loaded package imports Godbus and
`go mod why -m` says the main module does not need it, but MVS retains that
declared edge. The new exact root v5.1.0 edge selects the higher version,
changes no unrelated selection, and loads no package.

Current measurements are 234 selected modules, 3,586 graph edges, 429
complete-test entries, 41 loaded modules, 197 loaded module-backed packages,
1,057 `go.sum` lines, and a 400-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, sums are +41/-0. Main Go 1.18 and toolchain
Go 1.26.7 remain unchanged.

Fresh primary vulnerability data has 1,392 module records. Selected and
candidate normalize identically to 30 Darwin module findings, 22 Darwin
package findings, and 20 IDs/22 reachable traces for Darwin and Windows
symbol scans. Godbus has no record, finding, symbol, or reachable trace.

Decision-summary SHA-256 is
`8973eda3a30a7a1ec4311178f4ee4c701ba114d53ff7454ab8e947f8fcd754d2`;
the 576-entry selected-evidence manifest SHA-256 is
`a665eda7b6763d2e4b0815fc19269ea3d1424eee65e19761ea897a13dd075ac9`;
manual-evidence SHA-256 is
`7d07ea50ec7b968fa8a506873d3c800e53ed0524712a4ad2c81fec57b9a26580`.

## Tools

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
- Contained Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Pinned golangci-lint 2.12.2 archive, GoReleaser 2.17.1 binary, and apidiff
  receipts remain
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

## Next Objective

Independently evaluate exact-path `github.com/gogo/protobuf v1.3.2` as the
next single P7 group. Do not combine Viper, etcd, Prometheus, or any other
dependency group.

MVS selects v1.3.2 through Viper v1.15.0 and etcd/api v3.5.1; older v1.1.1
declarations remain through Prometheus TSDB/Common. `go mod why -m` says the
main module does not need it. The initial proxy survey exposes eight stable
versions, v1.0.0 through selected/latest v1.3.2, and no prereleases. V1.3.2
declares Go 1.15 and requires historical indirect errcheck, gotool, and
x/tools versions. Its source/mod sums are
`h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=` /
`h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=`. Treat these only as
incoming survey facts to verify.
