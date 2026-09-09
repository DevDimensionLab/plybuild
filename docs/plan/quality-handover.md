# Quality Upgrade Handover

Generated: 2026-09-09T05:20:20+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation is exact Go Logfmt v0.6.0 commit
  `3d4cfdbae0a67e757d37022be7eeedaf32c72772`, parent
  `f81dfd617a1f26211fd21213fbe96b8abd16d34f`, tree
  `34a91aaedab07e855022c454a16000a079f0dbbd`. It changes only `go.mod` and
  `go.sum`, with three insertions and no deletions.
- Fatih Color v1.15.0 implementation
  `6ca672ef38688b7f6f505cf0cb273d07c4c2ba9a`, XXHash v2.3.0
  `e5d6252825d7a1822c01819b9144050f345a6ad4`, and Speakeasy v0.2.0
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remain ancestors. Crypt,
  OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
  Ghodss YAML, and historical root GLFW remain retained without dependency
  edits.
- The answered Go Logfmt archive and sole NEXT Go Stack archive must link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` exists. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, and Go Logfmt v0.6.0
moves. P8 remains queued. Earlier decisions and lifecycle ancestry are final;
do not reopen them or combine another module group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*` roots,
never run `go mod download all` in a measured tree, and preserve launcher
scratch cleanup and reciprocal archive contracts.

## Go Logfmt Decision

Upgrade exact-path `github.com/go-logfmt/logfmt` from MVS-selected v0.4.0 to
v0.6.0. V0.6.0 is the highest qualified stable release whose complete minimal
closure preserves Go 1.18. Latest v0.6.1 declares Go 1.21 and is ineligible.
No redirect, fork, alternate path, floor-ineligible version, or unreleased
commit was promoted.

Go-import maps the exact module directly to the public
`https://github.com/go-logfmt/logfmt.git` repository, which is enabled,
unarchived, not disabled, not a fork, and defaults to its sole `main` branch.
The proxy exposes exactly eight stable versions, v0.1.0 through v0.6.1, and no
prereleases. All eight are annotated unsigned Git tags and non-draft,
non-prerelease GitHub Releases whose peeled commits are on main ancestry.
There is no nested module, retraction, deprecation, redirect, or exact-path
alternate identity. Proxy ZIPs v0.4.0-v0.6.1 byte-match the corresponding Git
trees, and sumdb confirms all evaluated checksum pairs.

Qualified v0.6.0 has source/mod sums
`h1:wGYYu3uicYdqXVgoYbvnkrPVXkuLM1p1ifugDMEdRi4=` /
`h1:WYhtIu8zTZfxdn5+rREduYbwxfcBr/Vr6KEVveWlfTs=` and proxy ZIP SHA-256
`a49c00cff30c02d9c09a4974ce91215bfe37f528a74f129576697869a1b8c630`.
Unsigned commit `76262ea710c6213a336b12b0356fec81341935e1`, tree
`53970738989e78890c85f821ca7541104d74cc8f`, parent
`5a3c9dc1265bdc95f1f72c912b61bb013e197d7d`, has commit time
2023-01-31T03:55:27Z, tag time 03:57:11Z, and release time 03:59:29Z.

Selected v0.4.0 has no Go directive and requires historical
`github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515`; that closure has no
further non-standard module. V0.5.0 declares Go 1.13, while v0.5.1 and v0.6.0
declare Go 1.17 and have no requirements. V0.6.1 declares Go 1.21 and requires
go-cmp v0.7.0. Source and imported-test inspection confirms the directives do
not hide an additional floor constraint.

## Package Behavior, API, And Gates

The exact module is one pure-Go package with no cgo, production build tags,
generation, or testdata. Its public surface consists of four sentinel errors,
`MarshalKeyvals`, decoder scan/accessor APIs, encoder encode/end/reset APIs,
`MarshalerError`, and `SyntaxError`. V0.5.x makes no public API change;
v0.6.0 adds only compatible `NewDecoderSize`.

Independent fixtures cover grammar and round trips, deterministic ordered
encoding, quoting/escaping, invalid and Unicode input, error identity and
offsets, streaming/record boundaries, default and explicit scanner limits,
marshaler precedence/errors, odd key/value lists, validation atomicity,
partial/short reads and writes, reset/reuse, and independent-instance
concurrency. Encoder and Decoder are deliberately stateful and unguarded;
separate instances are race-safe, while shared concurrent mutation is not a
supported use.

V0.6.1 fixes DEL U+007F handling: values are quoted/escaped and keys reject
DEL. V0.6.0 preserves selected's raw-DEL behavior. Logfmt has no formal
standard, no Ply package is loaded, and the change is recorded as an upstream
bug fix rather than made a new compatibility contract. The Go 1.21 directive
independently excludes v0.6.1.

Selected, v0.5.0, v0.5.1, v0.6.0, v0.6.1 source, and historical kr/logfmt
pass source verification, native tests, two count-10 repeats, race, vet, and
applicable Linux/Windows/Darwin/FreeBSD cross-builds under exact Go 1.26.7 and
contained Go 1.18.10. Selected's build-tagged gofuzz package builds, but
combining it with the normal suite duplicates legacy test names; classify that
as historical test design, not production failure.

All project build, native test, repeat, race, vet, offline-list, Windows,
pinned-lint, compatibility, launcher, Make, preflight, and empty-HOME gates
pass. The contained-Go-1.18 project build/mod verification pass; only two
inherited `pkg/shell` assertions fail on that SDK's closed-pipe error text,
with Go Logfmt unloaded and unrelated.

Exact changed-selection `make quality` exits zero after preflight, all 8
mutation meta-stages, 80/80 killed mutations, host acceptance, GoReleaser
snapshot acceptance, real Docker acceptance, and authoritative audit. All 27
Q0-Q2 rows PASS at L2, manual evidence is valid with six receipts, and ratchet
counts are seven improved and zero held/regressed/not-comparable. Scorecard
SHA-256 is
`0cff5be4fb1609d296664f13a2de5567c7bdb8f574e61b6b05c0d92f44c4eb8f`.
The separate full audit exits expected 1 only for queued Q3.1, Q3.3, Q3.4,
and Q3.7; scorecard SHA-256 is
`736c6e7dbc2e55df075c8fc15d58e9331d138835d875350f2a12d619657b6c54`.

## Project And Vulnerability Measurements

Before implementation, Prometheus Common v0.9.1 requested v0.4.0 and
Prometheus TSDB v0.7.1 requested v0.3.0, so MVS selected v0.4.0.
`go mod why -m` says the main module does not need Go Logfmt. No package from
go-logfmt or kr/logfmt is loaded; their modules remain because MVS retains
declared edges from otherwise unused Prometheus roots.

Exact `go get github.com/go-logfmt/logfmt@v0.6.0` adds only the main module's
indirect edge and exact checksum pair. Current measurements are 234 modules,
3,584 graph edges, 429 complete-test entries, 41 loaded modules, 197 loaded
module-backed packages, 1,053 `go.sum` lines, and a 392-line unapplied tidy
projection. Relative to accepted go-cmp commit `c314bcb`, sums are +37/-0.
No unrelated selection moved; historical kr/logfmt remains selected via
Prometheus TSDB. Main Go 1.18 and toolchain Go 1.26.7 are unchanged.

Fresh govulncheck v1.7.0 data has 1,392 module records. Selected and v0.6.0
normalize identically to 30 Darwin module findings, 22 Darwin package
findings, and 20 IDs/22 reachable traces for Darwin and Windows symbol scans.
Neither go-logfmt nor kr/logfmt has a record, finding, or trace.

Decision-summary SHA-256 is
`d72e3498eaafb37925188496a131d744d717b2de33cd4aaffeb62e57f7f88dc2`;
the 595-entry evidence manifest SHA-256 is
`ac964598ef5933db2136fb7797738ad5595c81d94e092cf0a58c7e4e9266e4b9`;
the 127-entry quality subset manifest SHA-256 is
`e86964d24c081b3619c489f8cb42c813b3ddef1b11e64185d340dd33e672daf8`;
manual-evidence SHA-256 is
`6d0cee48eb3f835cd409fb57af985a06d3c2d443fc8f3e148722089386869346`.

## Tools And Retained Decisions

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Contained Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Pinned golangci-lint 2.12.2 archive, GoReleaser 2.17.1 binary, and portable
  apidiff receipts remain
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
- Retain all accepted or retained P7 outcomes. Answered archives are the
  authoritative detail.

## Next Objective

Independently evaluate exact-path `github.com/go-stack/stack v1.8.0` as the
next single P7 group. Do not combine Prometheus Common, Prometheus TSDB, or any
other dependency group.

MVS selects v1.8.0 through Prometheus TSDB v0.7.1 and Prometheus Common
v0.4.1; `go mod why -m` says the main module does not need it. The initial
proxy survey exposes ten stable versions from v1.4.0 through v1.8.1. Selected
v1.8.0 has no Go directive or requirements and sum pair
`h1:5SgMzNM5HxrEjV0ww2lTmX6E2Izsfxas4+YHWRs3Lsk=` /
`h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=`.

Latest v1.8.1 declares Go 1.17 and has sums
`h1:ntEHSVwIt7PNXNpgPmVfMrNhLtgjlmnZha2kOpuRiDw=` /
`h1:dcoOX6HbPZSZptuspn9bctJ+N/CnF5gGygcUP3XYfe4=`. Its annotated unsigned
tag peels to commit `93c7c7e3550c72bc91dead1452a0020142e2a902`, tree
`970d18b6c7c8b790ab52ddd5d82cbaee98261044`, parents
`2fee6af1a9795aafbe0253a0cfbdf668e1fb8a9a` and
`473edce91b111d1f6c7b946691b24af71f5a5b25`, at
2021-08-18T18:48:21Z. Treat these only as incoming survey facts to verify.
