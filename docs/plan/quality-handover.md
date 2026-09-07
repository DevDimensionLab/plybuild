# Quality Upgrade Handover

Generated: 2026-09-07T21:46:38+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- Latest dependency implementation is Circbuf commit
  `3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
  `7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
  `a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
  `go.sum` with three insertions and no production-source change.
- Operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
  Repr implementation `6f4d02eb9c86ec2df8a488a85ed973aababe1f38` and
  direct-child Repr handoff `342c7ece82eae3cd726aa658462089cfb876fbe0`.
  All remain ancestors.
- The answered Circbuf archive and sole NEXT Consul API archive link
  reciprocally. Relative to accepted go-cmp commit `c314bcb`, accepted
  dependency metadata now adds exactly 29 checksum lines.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, successor launch, or
  worktree removal occurred.

## Lifecycle And Retained Roadmap

Every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Preserve the launcher lifecycle and scratch-
only mutation, compatibility, snapshot, Docker, acceptance, and audit rules.
Never run `go mod download all` in a measured tree.

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through accepted Circbuf pseudo-version
`v0.0.0-20190214190532-5111143e8da2`. All earlier accepted, rejected, and
no-change decisions and evidence corrections remain final. P8 remains queued.

## Accepted Armon Circbuf Group

Advance exact-path `github.com/armon/circbuf` from selected pseudo-version
`v0.0.0-20150827004946-bbbad097214e` to canonical latest and master-head
pseudo-version `v0.0.0-20190214190532-5111143e8da2`. The exact proxy stable
list is empty; exact Go and proxy `@latest` and `@master` resolve the candidate,
while `@v0` has no matching version. It is unreleased and must not be called a
stable release.

The authoritative repository is public, enabled, unarchived, and non-fork,
with default branch `master`, one branch, nine commits, zero tags, and zero
GitHub Releases. There are no stable versions, prereleases, retractions,
deprecation markers, or authoritative alternate `/v2`, `/v3`, renamed, or
gopkg.in module paths.

The old version is commit
`bbbad097214e2918d8543d5201d12bfd7bca254d` at
2015-08-27T00:49:46Z and is unsigned. The candidate is commit
`5111143e8da2e98b4ea6a8f32b9065ea1821c191`, tree
`2ab2d9cf2632ab7b549f7da7f081dbe868a697db`, at
2019-02-14T19:05:32Z. Its GitHub web-flow signature independently verifies
with fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; the signing key is currently
expired. There is no tag or tag signature.

Selected/candidate proxy ZIP SHA-256 values are
`3819cde26cd4b25c4043dc9384da7b0c1c29fd06e6e3a38604f4a6933fc017ed`
and `c8b7ba977844b5378a2413c123c3e55d0885fb67f64ad6cf06575a791a36b827`.
Their normalized source-manifest SHA-256 values are
`1bf21b9d070bf0d571b6f239ccbefdc4946472243cbe9993c5a7df937da5fdea`
and `9158ce75d0b4adcef1f783193bf7e00c4015d96e76160844ac29de83f235dc44`.
The only source delta is the candidate's one-line `go.mod`; every production
and test Go file is identical.

The old checksum pair is
`h1:QEF07wC0T1rKkctt1RINW/+RMTVmiwxETico2l3gxJA=` /
`h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=`. The candidate pair is
`h1:7Ip0wMmLHLRJdrloDxZfhMm0xrLXZS8+COSu2bXmEQs=` with the same `go.mod`
checksum. Independent checksum-database lookups agree.

Both versions declare no Go version and no requirements. The candidate's
complete minimal closure is Circbuf alone, so it cannot raise the retained Go
1.18 floor. Proxy and exact-commit forms independently verify/list, remain
unchanged, and pass complete count-1/count-10/race tests and vet under exact Go
1.26.7. The candidate also passes count-10/race/vet under exact Go 1.18.10.

The original graph requests Circbuf through Hashicorp Serf v0.8.2 and v0.9.6.
Circbuf loads in zero Ply packages, Ply has no repository import, and
`go mod why -m` says the main module does not need it. Historical Serf agent
code and selected Serf v0.10.1 source use `NewBuffer`, `Write`,
`TotalWritten`, `Size`, `String`, and `Bytes`. An explicitly external Go-1.18
fixture exercises those symbols and passes count-10, race, and vet; it is
consumer evidence, not a claim that Ply executes Circbuf.

Exact candidate get changes only Circbuf's selected version. Measurements are
234 -> 234 selected modules, 3,580 -> 3,581 graph edges, 429 -> 429 complete
packages, 41 -> 41 loaded modules, 1,043 -> 1,045 `go.sum` lines, and 356 ->
361 unapplied tidy-diff lines. The implementation adds an indirect candidate
requirement, one main-to-candidate edge, and the candidate full/module checksum
pair. Existing Serf edges remain. Tidy removes the explicit pin and the two
candidate checksum lines; inherited tidy debt is unchanged.

## Circbuf Quality And Evidence

- Repository module verification, build, complete count-1/count-10/race
  tests, vet, Windows-amd64 build, and pinned golangci-lint 2.12.2 pass.
  Root/status/upgrade/build help is byte-identical. API and CLI reports remain
  byte-identical at SHA-256
  `ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
  and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
- Complete preflight passes all 62 launcher controls and all distribution,
  install, lint, toolchain, snapshot/Docker, mutation/acceptance, and 15 audit
  meta-contracts. Host, fresh snapshot, and fresh Docker acceptance pass.
  Snapshot report SHA-256 is
  `a82c4524f47650f3521cd95ea4b1256b2022cae9ebbda35ed12dcbfa3109a157`;
  Docker report SHA-256 is
  `77f81c8f81e48557763e5c28a00d5f168f102abf8f5ef115d1a737b0bf875461`.
- Exact `make quality` exits 0 with all 27 Q0-Q2 rows at L2, 80/80 mutations
  killed, and zero held, regressed, not-comparable, or dirty counts. Its
  scorecard SHA-256 is
  `855329bc3e39600a5704cf0b1267002cace378f3fc9cb139b5ff3784f1b7cb3c`.
  Full audit exits 1, never 2, only for established queued L3 rows Q3.1,
  Q3.3, Q3.4, and Q3.7; scorecard SHA-256 is
  `1c99edae5e9783a83bb7f480df067ba00ac83687f3fae9b56bab28da1536873b`.
- Fresh primary vulnerability data updated 2026-09-02T19:12:04Z has 1,392
  module records and no Circbuf record or trace. Sorted old/candidate results
  are identical at exact 20 Darwin reachable-symbol IDs/22 traces, 30 Darwin
  module IDs, and 20 Windows reachable-symbol IDs/22 traces. Normalized
  summary SHA-256 is
  `806aad56f111b194bcd3aa3f5d45730d11666ef1e85e3acf8136ad32ecaa3f88`.
- Decision evidence contains 430 verified entries. Evidence-manifest SHA-256
  is `2684f8c269376320bcfca3404bfbf10c638c81753e6af36b5be59c4df2370d62`;
  decision-summary SHA-256 is
  `17441163b7924e2cd61cb167b739562da7672fa5abfa10089da36ef906f07ff9`.

## Tools And Corrections

- Exact Go 1.26.7 binary SHA-256 is
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  official archive SHA-256 is
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Exact Go 1.18.10 archive SHA-256 is
  `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
  Go 1.26.7 remained first in PATH with GOENV=off, GOWORK=off,
  GOTOOLCHAIN=local, and no ambient GOFLAGS.
- Preserve portable golangci-lint 2.12.2 archive receipt
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  GoReleaser 2.17.1 receipt
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and apidiff receipt
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
  Govulncheck v1.7.0 and rebuilt-tool hashes are nonportable receipts; versions
  and functionality were proved.
- Preserve the corrected source-prefix normalization, fixture-local checksum
  population and fail-fast rerun, concatenated vulnerability-JSON decoding and
  trace-depth filter, exact compatibility-cache warming, resolved Docker
  binary path, scratch-local BSD-mktemp adapter, and C-locale/exact-Go-1.18.10
  final manifest rerun. Initial failures were superseded runner/setup checks;
  none changed repository files or hid a candidate failure.

## Retained Decisions

- Retain exact selections `github.com/alecthomas/repr v0.5.4`,
  `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
  `github.com/alecthomas/chroma v0.10.0`,
  `github.com/alecthomas/colour v0.1.0`,
  `github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751`,
  `github.com/antihax/optional v1.0.0`, accepted Circbuf candidate,
  `gopkg.in/alecthomas/kingpin.v2 v2.2.6`, `gopkg.in/resty.v1 v1.12.0`,
  `gopkg.in/errgo.v2 v2.1.0`,
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
  `gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
  `go.yaml.in/yaml/v3 v3.0.5`. Retain selected Kong pseudo-version and every
  other prior bounded decision.
- Preserve all recorded corrections from prior dependency groups. The
  deliberately purged former two-entry recovery manifest remains historical
  SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.

## Next Objective

Independently evaluate selected exact-path
`github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6` as the next
single P7 group. It is selected by project MVS but is not an explicit `go.mod`
requirement. A minimal post-Circbuf survey finds an empty exact stable proxy
list. Proxy and exact Go `@latest` and `@master` select the same pseudo-version
at 2018-02-02T20:16:55Z; exact `@v0` has no match. Its selected `go.mod`
checksum is `h1:grANhF5doyWs3UAsr3K4I6qtAmlQcZDesFNEHPZAzj8=`.

The authoritative `armon/consul-api` repository reports public, enabled,
unarchived, non-fork status, default branch `master`, one branch, zero tags,
and zero GitHub Releases; its recorded push time matches the selected commit
time. Treat canonical qualification, source/signatures, release and alternate-
path history, complete closure, module tests, loaded packages, actual symbols,
and vulnerability effect as unknown until independently proved. Do not call
the pseudo-version stable.

Current measurements are 234 selected modules, 3,581 graph edges, 429 native
complete-test packages, 41 loaded modules, 1,045 `go.sum` lines, a 361-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Change only Consul API's exact
selected version and explained minimal MVS closure if a qualified changed
selection preserves the floor and passes every applicable gate. If the
selected version is canonical latest, do not manufacture redundant metadata.
Do not combine Circbuf, another module group, or P8.
