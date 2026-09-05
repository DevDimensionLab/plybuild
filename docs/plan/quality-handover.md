# Quality Upgrade Handover

Generated: 2026-09-05T22:15:59+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation remains dependency-only cast commit
  `cf4fd4933251b2ca7844b6420fd3cf2a86f049a5`, exact parent
  `17f7277fd2512ba06db76d1acb0af5f8b623243c`, clean tree
  `485690c010cef3b02385d8f2471d3ef53611abe9`. It changes only `go.mod` and
  `go.sum`.
- The incoming Errgo v2 documentation session is commit
  `1899ff708ae165c43c283e97def420af2c11d108`, exact parent
  `a8e707ce125b6ecb72c203925f6978b09a85406f`, tree
  `46a0328a8e910132a35814cadc291a99f6c514b7`. This handoff must be its
  direct child. The next dependency implementation, if any, must use the
  resulting documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Their minimal
  closures select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty
  v0.3.1, and rogpeppe/go-internal v1.9.0. Those four accepted groups add
  exactly 15 checksum lines.
- The answered Errgo v2 archive and sole NEXT `gopkg.in/resty.v1` archive link
  reciprocally. Only the launcher's mutable header and prompt regions change
  during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1,
rejected afero v1.10.0, retained canonical-latest jwalterweatherman v1.1.0,
rejected gotenv v1.6.0, x/image v0.16.0, x/net v0.25.0, x/text v0.15.0,
retained canonical-latest YAML v2 v2.4.0 and YAML v3 v3.0.1, rejected
canonical-latest Check v1 pseudo-version
`v1.0.0-20201130134442-10cb98267c6c`, and retained canonical-latest Errgo v2
v2.1.0. All earlier recorded dependency decisions and evidence corrections
remain final. Do not revisit them. P8 remains queued.

## Retained Gopkg Errgo V2 Group

Fresh proxy evidence lists exactly three stable semantic versions: v2.0.0,
v2.0.1, and v2.1.0. It contains no prerelease or retracted version. Proxy
`@latest`, exact `go list -m -versions -retracted ...@latest`, `@v2`,
`@master`, and authoritative gopkg source metadata all select v2.1.0 at
2018-08-23T07:15:05Z. Selected v2.1.0 is therefore canonical latest and the
highest qualified stable release compatible with the retained Go 1.18 floor.

The three releases use annotated, unsigned tags:

- v2.0.0: tag `4c49f8a4dbfc9d09c2cab69965377f39e247e53a`, commit
  `82cd8dd47e1c988e185c8c89d0d05aff1c3dac50` at
  2018-08-18T17:56:26Z.
- v2.0.1: tag `946adb7488b1a3af293c1ecaab3f83fb3d5a0807`, commit
  `16491ea1e140f01d877d5daa06d9cebb81ca8cf1` at
  2018-08-22T21:36:27Z.
- v2.1.0: tag `635edbc13741bd819966931e8b599c690b4f074d`, commit
  `f768c5ab0476c50e978b039312180859c10fe8c0`, tree
  `cb540a0ae1ad20359a13e238d412878b955e6099`, at
  2018-08-23T07:15:05Z.

GitHub reports all three peeled merge-commit signatures verified and valid;
the tag objects carry no signatures, and local cryptographic verification was
unavailable because `gpg` is absent. The selected checksum pair is
`h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=` /
`h1:hNsd1EY+bozCKY1Ytp96fpM3vjJbqLJn88ws8XvfDNI=`. All three proxy archives
match the exact upstream tag commits file-for-file. The v2.0.0, v2.0.1, and
v2.1.0 proxy ZIP SHA-256 values are respectively
`04f40f54a7f7c6d71a26bf69d6dc5db935bd11af7c5c6a824bd51bab70c05b4e`,
`011b9480231acccf099b674348ae236f8fa8300244b708496688df76db3ae4d1`,
and `6b8954819a20ec52982a206fd3eb94629ff53c5790aa77534e6d8daf7de01bee`.

Authoritative gopkg metadata maps the exact module to
`https://gopkg.in/errgo.v2` and source browsing to
`github.com/go-errgo/errgo/tree/v2.1.0`. Fresh gopkg and GitHub clones share
the same relevant tag, v2, and v1 objects. Gopkg synthesizes its `master` ref
to the v2.1.0 commit, whereas GitHub's actual `master` is an older divergent
2014 line; neither is a later candidate. The GitHub repository is enabled,
unarchived, and not a fork, has no GitHub Release objects, and carries no
module or README deprecation notice. Its v2 line is dormant after August 2018.

Two v2 commits follow the selected release. Linear commit
`28806950c76b10f5cff5bcaf2d6b5cddce8502a0` removes gocheck, and verified
merge commit `81e25171e4fb5331d1f784276d93cb70a0a31475` rewrites the tests. Exact
commit resolution produces
`v2.1.1-0.20180823084306-81e25171e4fb`, with checksum pair
`h1:/VNQHfUDtkG/Q76yWdItpo8Uzu5Cx16cYVzWU8PoqCY=` /
`h1:e2+3YDPGA8XJ8Lk+bKAvy+3jAR1cSuuG8+hJ3GeqXI8=`. This is an unreleased
pseudo-version, not a stable release and not the version selected by canonical
proxy or gopkg resolution. It is excluded from the qualified candidate set.

V2.1.0 has no Go directive. It declares only kr/pretty v0.1.0 and Check's 2018
pseudo-version, both of which omit Go declarations. The project already
selects kr/pretty v0.3.1 at Go 1.12 and the retained 2019 Check pseudo-version
with no Go directive; the selected go-internal v1.9.0 declares Go 1.17. The
empty changed-selection closure therefore preserves Go 1.18 by declarations,
not merely because tests run under a modern toolchain.

An external exact `go get gopkg.in/errgo.v2@v2.1.0` exits zero, but selected
module lists remain byte-identical at 234 entries. It only projects a redundant
indirect main-module requirement, one main-to-Errgo graph edge, and the full
v2.1.0 module checksum. Graph edges project 3,564 -> 3,565, go.sum lines
1,031 -> 1,032, and the unapplied tidy diff 332 -> 344 lines; complete packages
remain 429. Since no exact version selection changes, this redundant explicit
pin is not an authorized implementation. Go.mod and go.sum remain unchanged
and no dependency commit was manufactured.

Exact `go list -test -deps -json ./...` finds zero Errgo packages in both
states, repository Go source has zero Errgo imports, and `go mod why -m` says
the main module does not need Errgo. The version is historical unloaded MVS
graph debt contributed by an old go-internal v1.3.0 requirement edge, not a
main-module or dependency-test consumer path.

Two independent writable module replays from the proxy source and the exact
upstream tag remain byte-identical, verify, list two packages, and pass
count-1, count-10, race, and vet. Their standalone build list contains five
modules and is test apparatus separate from the project closure. Old and
projected repository build, complete tests/race/vet, Windows build, pinned
lint, CLI surface, and public root/status/upgrade/build help pass with
byte-identical output. Corrected API and CLI reports are byte-identical.
Empty-HOME count-2 passes.

The final uncontaminated full preflight exits zero with all 62 launcher checks,
80 mutation controls, and 15 audit meta-controls passing. Several prior full
runs retain the established primary or nested signal-fixture timing diagnostic;
independent direct and Make launcher suites pass 62/62. A separate preflight
attempt correctly reached the lint meta-suite but was contaminated by outer
Make command-line tool overrides propagating through `MAKEFLAGS`; the final
environment-only pinned-tool invocation removes that runner error. The first
API compatibility attempts also failed closed only because the external
historical v1.0.1 archive cache lacked go-isatty v0.0.16; after warming that
separate archive, both corrected retries passed. No measured tree was used for
`go mod download all`.

Govulncheck v1.7.0 uses a primary database updated
2026-09-02T19:12:04Z. The fresh 1,392-entry module index has no Errgo record,
finding, or trace. Old and projected outputs and exact ID sets are
byte-identical for Darwin symbol, Darwin module, and Windows symbol scans,
preserving populations 20/30/20.

## Errgo V2 Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-errgo-v2-selection.1899ff7.HoEmPo`. Its fully verified
  152,333-entry manifest SHA-256 is
  `e9a3c3b1b9c311a1dd96f1804167d57359c2af02615dde2ffb544690dac45ab2`;
  decision-summary SHA-256 is
  `aa8a30c6f9e261a5c29be3bd59d87cb96e5b07f85582b6f1d55a4f78fcf3a3ec`.
- The manifest covers all 32 prerequisite-root verifications, fresh
  proxy/sumdb/gopkg/Git/GitHub resolution, exact stable tags and later v2
  commits, source and repository identity, signatures, Go-floor declarations,
  exact-get/tidy and minimal projection diffs, loaded-population classification,
  two independent module replays, old/projected repository gates,
  compatibility/help identity, the complete preflight and empty-HOME gate,
  and byte-identical primary vulnerability scans.

## Inherited Evidence And Tool Identity

- All 32 prerequisite manifest roots were independently recomputed before
  measurement: the 31 inherited roots recorded by Check and the Check root
  itself. The exact PASS table is sealed in Errgo evidence.
- Check rejection evidence remains fully verified at
  `/private/tmp/ply-p7-check-v1-selection.a8e707c.DD2k28`,
  18,182/`d45958adc2d3be460e85b3f6e379053a8b011af7cd221d4026d9cb0e1bfd57df`;
  its decision summary is
  `93d0b2833dc78ea2cddb7ea0454713ef96a1ee1879d2ab38692d44bf9e8c7e51`.
- YAML v3 and YAML v2 evidence remain fully verified at
  40,121/`84bd12b7b4cbff806abcbff213b4b7bc5bac230a7371c3fe0b3e2def643ed741`
  and
  29,219/`8dd34cc54f56e4bb1370b2f7f1f124aa8b501a4dcca177effe4f6377c3919e14`.
- X/text, x/net, x/image, gotenv, and jwalterweatherman evidence remain fully
  verified at
  82,629/`48c872da4ab796ef1115003fcf0a226a07bf5791a163d36bc9c946c9d9f92ed6`,
  71,887/`88e4be87483322792ad3da0064bcef4ced5ee21c52821bb1f80c143f80a3aaf1`,
  52,533/`84972465aef5f19f888102c7991aadd67446146a7056658b006d005666bd82bf`,
  48,662/`ab2bd356b863d529085cff9b918d004aa2a03a121a7ea360a66767a2820cc88c`,
  and 76,374/`54040aefec38df7415f56da66af04d48c8256dd46ae6de8e06a7c4bbe5ce4d8e`.
- Cast selection/review/exact-quality/regression evidence remains fully
  verified at 82,586/12/236,782/607 entries with manifest SHA-256 values
  `dfa699f75d0b222d7d64e8fd7a09d05306de98ef4ac5cbf6d253fd33b2dddac0`,
  `850d40f319d13024275b35909e8ff84ce2f9c5ce5959cf8d707461a3b5e5d9ce`,
  `b3d410a70d89e2176daf1a013f97d53e704ef3887f02d67bd6c2564ab0a6d02a`,
  and `97cafb6719f470ab7d057a603b522f0d10418deda2d997e6d5723883efd2e24a`.
  Preserve its NUL-delimited whitespace-path correction.
- Afero rejection evidence remains fully verified at
  79,866/`f0ad8ce41436a5f8bf48222106e79f57b32a9533b6a9a15b290615a54f162405`.
- Preserve the go-colorful correction: 47,632 stable entries pass and exactly
  one mutable telemetry counter differs from its 47,633-entry manifest. Keep
  manifest SHA-256
  `a47facc21231accc5c2ed73980c9e3a7542538c6a41d7778de5cfbb410312adc`.
  Preserve the btree correction: exactly 163 mutable cache/HOME paths differ
  from its 24,197-entry manifest, whose SHA-256 is
  `9cca917b12dd761c58bf91652e78b3e999f55eeb1ffc39f9f3bbf56a52aef463`.
- Recovery `/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` remains fully
  verified at two entries and manifest SHA-256
  `1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
- Tool identities were recomputed as Go 1.26.7
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
  golangci-lint 2.12.2
  `3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`,
  GoReleaser 2.17.1
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`,
  and govulncheck v1.7.0
  `0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

## Next Objective

Independently evaluate selected `gopkg.in/resty.v1 v1.12.0` as exactly one
bounded P7 module group. Resolve canonical latest, stable/prerelease and
pseudo-version qualification, exact-path and alternate-path history, Go
declarations and complete floor-compatible closure, source identity,
signature and repository state, loaded package population and real consumers,
module self-tests, and primary vulnerability data from fresh evidence before
selecting anything.

The current build list and graph measurements remain 234 selected modules,
3,564 graph edges, 429 complete packages, 1,031 go.sum lines, a 332-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Use exact Go 1.26.7 with GOENV off,
GOWORK off, GOTOOLCHAIN local, and no ambient GOFLAGS.

Implement only an exact floor-compatible changed selection with an explained
minimal closure and every applicable gate passing. If a later branch or module
path exists, prove whether it qualifies for exact `gopkg.in/resty.v1` rather
than treating it as an in-place release. If exact selected is already the
decision, do not manufacture an explicit requirement or dependency commit.
Stop on any mandatory failure. Do not revisit Errgo, Check, either gopkg YAML
group, the accepted `go.yaml.in/yaml/v3` group, another dependency, or P8.
