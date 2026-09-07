# Quality Upgrade Handover

Generated: 2026-09-07T19:49:16+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- Latest dependency implementation remains Repr commit
  `6f4d02eb9c86ec2df8a488a85ed973aababe1f38`, exact parent
  `53475076e1c79f6d2181877e3238a1a5d389c246`, tree
  `82e7b1f5c503659082207481b8339e8113e38c10`, changing only `go.mod` and
  `go.sum`.
- Operator-authorized lifecycle repair
  `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
  that implementation and direct-child Repr handoff
  `342c7ece82eae3cd726aa658462089cfb876fbe0`. All three remain ancestors.
- The answered Optional archive and sole NEXT Circbuf archive link
  reciprocally. Relative to accepted go-cmp commit `c314bcb`, accepted
  dependency metadata still adds exactly 27 checksum lines. Template and
  Optional changed no metadata.
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
dependency decisions through retained Optional v1.0.0. All earlier accepted,
rejected, and no-change decisions and evidence corrections remain final. P8
remains queued.

## Retained Antihax Optional Group

Retain exact-path `github.com/antihax/optional v1.0.0` without a dependency
edit. The exact proxy list contains only v1.0.0. Exact Go `@latest`, `@v1`, and
`@v1.0.0` select it at 2019-10-10T23:37:20Z. There are no prereleases or
retractions. It is canonical latest, the sole stable exact-path release, and
the highest stable version compatible with retained Go 1.18.

The authoritative repository is public, enabled, unarchived, and non-fork. It
has one branch, one tag, one GitHub Release, and 14 downstream forks. The
v1.0.0 Release is named `Initial release`, is neither draft nor prerelease,
and was published at 2019-10-10T23:39:21Z. The lightweight tag resolves commit
`c3f0ba9c1a592b971d66b2787679af55b5c58f21`, tree
`b9328a8aa4526004bb36928dbc136c3acb8eec3a`. It has no tag object or tag
signature. The associated GitHub web-flow commit signature independently
verifies with fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`.

Exact `@master` resolves only as unreleased pseudo-version
`v1.0.1-0.20220101210036-407d38fabb55`, commit
`407d38fabb5592e58b0841e8adb93edd65ee1319`, tree
`3dffdae3ba5f11e0b140c2edd085e5d1a69df40d`, at
2022-01-01T21:00:36Z. Its web-flow signature also verifies, but this is not
release evidence. The two post-tag commits add/update only `README.md`; all 21
Go files and `go.mod` are byte-identical to v1.0.0. `/v2`, `/v3`, gopkg.in,
and plausible renamed-module probes find no alternate module.

V1.0.0 proxy ZIP SHA-256 is
`15ab4d41bdbb72ee0ac63db616cdefc7671c79e13d0f73b58355a6a88219c97f`;
its checksum pair is
`h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=` /
`h1:uupD/76wgC+ih3iEmQUL+0Ugr19nfwCT1kdvxnR2qWY=`. The pseudo-version ZIP
SHA-256 is
`7a4fe848cb95f2bd6305266a9a35b515933cf4144436d865a2569b7f1a56d92d`;
its source checksum is
`h1:AAnkpz6e/6bV9pphv0Hj/13169F89gXw0O4kOduU1Gc=` with the same `go.mod`
hash. Sumdb agrees.

All 23 v1.0.0 proxy regular files match the exact tag; all 24 pseudo-version
files match master. Their normalized manifest SHA-256 values are
`66452eb1ab11bc7dc7bfe9237c47c430810d77633add3a4ef7304e1c2db48c38`
and `916591eaf56bef23e0e97a616c3ecf5b73f7f01b1c502d319a808cf7acfb3b27`.
The tag/head 21-file Go manifest is identical at SHA-256
`781e71995c7d98ab7bba04b866c58011b57956db2fd15898f65d29d10919a761`.

Both versions declare Go 1.13 and no requirements. V1.0.0's complete declared
closure is the module alone and preserves Go 1.18. There are no test-only
requirements. Writable proxy and exact-tag copies independently verify/list,
remain unchanged, and pass complete count-1/count-10/race tests and vet. The
single package has no test files; the receipts prove buildability and absence
of failing upstream tests, not upstream behavioral coverage.

Exact external selected-version get exits zero but changes no selected module.
It projects 234 -> 234 modules, 3,580 -> 3,581 graph edges, 429 -> 429 complete
packages, 1,043 -> 1,044 `go.sum` lines, and a 356 -> 358-line unapplied tidy
diff. Module and package diffs are empty. The only effects are a redundant
indirect requirement, one main graph edge, and the selected full checksum;
tidy removes the requirement. No metadata was applied because this is not a
version upgrade.

The original graph has only grpc-gateway v1.16.0 -> Optional v1.0.0. Optional
loads in zero project packages, repository Go source has no import, and
`go mod why -m` says the main module does not need it. Four generated grpc-
gateway example clients use Optional Bool, Float64, Int32, Int64, Interface,
String, and Time wrappers with `IsSet` and `Value`, but those packages are not
loaded by Ply. An explicitly external Go-1.18 fixture exercises all seven plus
`Default` and passes count-10/race/vet. Do not mislabel it as project-loaded
behavior.

Untouched and projected repository copies pass module verification, build,
complete count-1/count-10/race tests, vet, Windows build, and pinned
golangci-lint 2.12.2. API/CLI reports and root/status/upgrade help are byte-
identical. Compatibility against v1.0.1 and CLI surface pass. A clean full
preflight passes all 62 launcher controls, distribution/lint/install/toolchain
contracts, snapshot/Docker meta-contracts, mutation and verification meta-
tests, and all 15 quality-audit controls. Changed-selection host, fresh
snapshot/Docker, exact-quality, focused-audit, and empty-HOME gates are
inapplicable because no selected version changed.

Fresh primary vulnerability data updated 2026-09-02T19:12:04Z has 1,392 module
records and no Optional record. Sorted old/projected IDs and traces are
identical at exact 20 Darwin reachable-symbol IDs/22 traces, 30 Darwin module
IDs, and 20 Windows reachable-symbol IDs/22 traces. No trace mentions Optional.

## Optional Evidence, Tools, And Corrections

- Decision evidence has 540 verified entries. Evidence-manifest SHA-256 is
  `345a29c62bfadf332f14aa234b1265a9dc9142739d212534d83ace531fb39afa`;
  decision-summary SHA-256 is
  `945b7ca8fe9d9faa605a6a80d8d416655925f95cd727071babf813bcd4caa27b`.
- Exact Go 1.26.7 binary SHA-256 is
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  official archive SHA-256 is
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  It remained first in PATH with GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and
  no ambient GOFLAGS.
- Preserve portable golangci-lint 2.12.2 archive receipt
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
  GoReleaser 2.17.1 receipt
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  and apidiff receipt
  `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
  Freshly rebuilt apidiff and govulncheck hashes are nonportable; exact
  versions and functionality were proved.
- Preserve the separately warmed v1.0.1 compatibility cache; scratch-local
  GnuPG home and exact detached-payload correction; source-manifest root
  correction; project `rg` PATH and checksum-delta corrections; sorted
  vulnerability normalization; and scratch-local BSD-mktemp/preflight
  correction. Initial failures were runner setup or report-order issues and
  were superseded by passing exact reruns. None changed the repository.

## Retained Decisions

- Retain exact selections `github.com/alecthomas/repr v0.5.4`,
  `github.com/alecthomas/assert v1.0.0`,
  `github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b`,
  `github.com/alecthomas/chroma v0.10.0`,
  `github.com/alecthomas/colour v0.1.0`,
  `github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751`,
  `github.com/antihax/optional v1.0.0`,
  `gopkg.in/alecthomas/kingpin.v2 v2.2.6`, `gopkg.in/resty.v1 v1.12.0`,
  `gopkg.in/errgo.v2 v2.1.0`,
  `gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15`,
  `gopkg.in/yaml.v2 v2.4.0`, `gopkg.in/yaml.v3 v3.0.1`, and
  `go.yaml.in/yaml/v3 v3.0.5`. Retain selected Kong pseudo-version and every
  other prior bounded decision.
- Preserve all recorded corrections, including go-colorful mutable telemetry,
  btree regression, cast whitespace-path, source-archive normalization, Check
  history, Errgo preflight runner, Resty committed projection/signal fixture,
  Kingpin external report, Units launcher timing, Assert cache/TMPDIR, Repr
  cache/mktemp/Docker/manifest, Chroma source/consumer normalization, Colour
  pre-module/cache/mktemp, Template GnuPG keyring, and Optional corrections.
- The deliberately purged former two-entry recovery manifest remains
  historical SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.

## Next Objective

Independently evaluate selected exact-path
`github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e` as the next
single P7 group. It is selected by project MVS but is not an explicit `go.mod`
requirement. A fresh post-Optional survey finds the exact stable proxy list
empty. Proxy and exact Go `@latest` and `@master` select newer pseudo-version
`v0.0.0-20190214190532-5111143e8da2` at 2019-02-14T19:05:32Z without a
reported Go declaration; exact `@v0` has no matching version. The public,
enabled, unarchived, non-fork repository reports default branch `master`, zero
tags, and zero GitHub Releases. Do not call the pseudo-version a stable release.
Treat canonical qualification, source/signatures, release and alternate-path
history, complete closure, module tests, loaded packages, actual symbols, and
vulnerability effect as unknown until independently proved.

Current measurements remain 234 selected modules, 3,580 graph edges, 429
native complete-test packages, 1,043 `go.sum` lines, a 356-line unapplied tidy
projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations. The main module retains Go 1.18 and toolchain Go
1.26.7. Change only Circbuf's exact selected version and explained minimal MVS
closure if a qualified changed selection preserves the floor and passes every
applicable gate. Do not combine Optional, another module group, or P8.
