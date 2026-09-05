# Quality Upgrade Handover

Generated: 2026-09-05T20:51:13+02:00

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
- The incoming Check v1 documentation session is commit
  `a8e707ce125b6ecb72c203925f6978b09a85406f`, exact parent
  `cdb5b4ca02fa5db54efebf741d2b415608c053b2`, tree
  `cb258bc31914dc1e9fef76a756eb99a0f1da69be`. This handoff must be its
  direct child. The next dependency implementation, if any, must use the
  resulting documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Their minimal
  closures select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty
  v0.3.1, and rogpeppe/go-internal v1.9.0. Those four accepted groups add
  exactly 15 checksum lines.
- The answered Check v1 archive and sole NEXT `gopkg.in/errgo.v2` archive
  link reciprocally. Only the launcher's mutable header and prompt regions
  changed during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1,
rejected afero v1.10.0, retained canonical-latest jwalterweatherman v1.1.0,
rejected gotenv v1.6.0, x/image v0.16.0, x/net v0.25.0, x/text v0.15.0,
retained canonical-latest YAML v2 v2.4.0 and YAML v3 v3.0.1, and rejected
canonical-latest Check v1 pseudo-version
`v1.0.0-20201130134442-10cb98267c6c`. All earlier recorded dependency
decisions and evidence corrections remain final. Do not revisit them. P8
remains queued.

## Rejected Gopkg Check V1 Group

Fresh Go proxy and checksum-database evidence shows an empty `@v/list`: this
exact module path has no listed stable semantic-version or prerelease tags.
Proxy `@latest`, exact `go list ...@latest`, and the authoritative gopkg
`master` query all resolve
`v1.0.0-20201130134442-10cb98267c6c` at
2020-11-30T13:44:42Z. It is a pseudo-version, not a stable release. There are
no Git tags, retractions, or GitHub Release objects; `@v1` has no matching
version. The candidate is canonical latest and the highest qualified version
compatible with the retained Go 1.18 floor.

The selected pseudo-version
`v1.0.0-20190902080502-41f04d3bba15` maps to commit
`41f04d3bba152ddec2103e299fed053415705330`, tree
`d293fcb19bdb3deacaf589ad37f4534d1342c8ae`, at
2019-09-02T08:05:02Z. Its synthetic module file has no Go declaration. The
candidate maps to commit
`10cb98267c6cb43ea9cd6793f29ff4089c306974`, tree
`b5ee34e90064a88d7614bf6f6714d32bfac33e66`, and declares Go 1.11 with
only `github.com/kr/pretty v0.2.1`. Candidate checksums are
`h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=` /
`h1:JHkPIbrfpd72SG/EVd6muEfDQjcINNoR0C8j2r3qZ4Q=`; selected checksums are
`h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=` /
`h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=`.

Authoritative gopkg metadata maps `gopkg.in/check.v1` to
`https://gopkg.in/check.v1` and source browsing to
`github.com/go-check/check/tree/v1`. Fresh gopkg and GitHub clones have
identical v1 histories and exact objects. All seven resolved pseudo-version
snapshots match their upstream Git commits file-for-file. The 27-file
candidate proxy ZIP has SHA-256
`f555684e5c5dacc2850dddb345fef1b8f93f546b72685589789da6d2b062710e`.
The GitHub repository is unarchived, enabled, and not a fork; its default v1
branch is the candidate. GitHub's `master` head
`163297374fe15df72fc217b2bdfabb5894c73059` is an older divergent 2014
line, not an upgrade candidate. There is no module deprecation marker or
README deprecation notice.

The candidate has no tag or tag signature. Its commit carries an embedded
signature that GitHub reports as verified and valid; local cryptographic
verification was unavailable because `gpg` is absent. The selected commit is
unsigned. Exactly four commits separate selected from candidate: a temporary
kr/pretty fork, documentation correction, `ioutil.TempDir` repair, and the
candidate module/CI/Windows-test update that restores upstream kr/pretty. Net
history changes seven files by +67/-16.

An exact external projection of
`go get gopkg.in/check.v1@v1.0.0-20201130134442-10cb98267c6c` exits zero.
It projects 234 -> 236 selected modules: Check advances, while
`github.com/creack/pty v1.1.9` and
`github.com/pkg/diff v0.0.0-20210226163009-20ebb0f2a09e` enter the build
list. Ten graph edges are added, none removed, for 3,564 -> 3,574. Complete
project packages remain exactly 429. The projection adds four indirect
go.mod requirements for Check, already-selected kr/pretty v0.3.1, kr/text
v0.2.0, and rogpeppe/go-internal v1.9.0, plus exactly eight go.sum lines.
Those closure modules declare no Go version or Go 1.12, 1.13, 1.15, or 1.17,
so the complete changed closure preserves Go 1.18. The unapplied tidy
projection grows from 332 to 353 lines. Project module verification passes in
both states.

Exact `go list -test -deps ./...` finds zero Check packages in either state,
and repository source has zero Check imports. The complete package population
is unchanged. `go mod why -m` crosses
`pkg/config -> gopkg.in/yaml.v2 -> gopkg.in/yaml.v2.test ->
gopkg.in/check.v1`; this is a dependency-test path, not a main-module
consumer. Check is selected by dependency requirements but is not loaded by
the complete native project population.

The candidate standalone module build list contains Check, kr/pretty v0.2.1,
kr/text v0.1.0, and kr/pty v1.1.1, with no extra test-only requirements. Two
independent source replays verify, list the sole package, and pass complete
count-1 and race tests without source mutation. Both count-10 runs fail
identically: retained suite registration state makes
`EmbeddedS.TestMethod` observe true rather than false, then 16 suites run
where eight are expected. Both vet runs also fail identically on four unkeyed
`go/printer` and `go/ast` composite literals. These mandatory module
failures reject the candidate and stop the group before repository
build/tests/race/vet, pinned lint, help/API/CLI, snapshot/Docker, quality, and
audit acceptance. No dependency metadata or implementation commit was made.

Govulncheck v1.7.0 uses a primary database updated
2026-09-02T19:12:04Z. Its fresh 1,392-entry module index has no Check record,
finding, or trace. Old and projected-candidate outputs and exact ID sets are
byte-identical for Darwin symbol, Darwin module, and Windows symbol scans,
preserving populations 20/30/20. They are also byte-identical to the retained
YAML v3 evidence.

## Check V1 Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-check-v1-selection.a8e707c.DD2k28`. Its fully verified
  18,182-entry manifest SHA-256 is
  `d45958adc2d3be460e85b3f6e379053a8b011af7cd221d4026d9cb0e1bfd57df`;
  decision-summary SHA-256 is
  `93d0b2833dc78ea2cddb7ea0454713ef96a1ee1879d2ab38692d44bf9e8c7e51`.
- The manifest covers all 31 prerequisite-root verifications, fresh
  proxy/sumdb/gopkg/Git/GitHub resolution, seven pseudo-version snapshots,
  source and branch identity, signature/repository state, selected-to-latest
  history, Go-floor declarations, exact-get/tidy and minimal closure diffs,
  loaded-package and consumer classification, two independent candidate
  module replays, source-mutation checks, and byte-identical primary
  vulnerability scans. The final pseudo-version table records exact
  `go list` times rather than the time-less `go mod download` JSON.

## Inherited Evidence And Tool Identity

- All 31 prerequisite manifest roots were independently recomputed before
  measurement: the 30 inherited roots recorded by YAML v3 and the YAML v3 root
  itself. The exact PASS table is sealed in Check evidence.
- YAML v3 no-change evidence remains fully verified at
  `/private/tmp/ply-p7-yaml-v3-selection.cdb5b4c.bbbTdr`,
  40,121/`84bd12b7b4cbff806abcbff213b4b7bc5bac230a7371c3fe0b3e2def643ed741`;
  decision summary is
  `95d5aa003b77880ecf72770ad66c061e0013dfd022fc222eb600860bc9a7799f`.
- YAML v2, x/text, x/net, x/image, gotenv, and jwalterweatherman roots remain
  fully verified at
  29,219/`8dd34cc54f56e4bb1370b2f7f1f124aa8b501a4dcca177effe4f6377c3919e14`,
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
  `/private/tmp/ply-p7-afero-selection.5a46539.udvGSh`,
  79,866/`f0ad8ce41436a5f8bf48222106e79f57b32a9533b6a9a15b290615a54f162405`.
- Preserve the go-colorful correction: its manifest retains 47,633 entries and
  SHA-256
  `a47facc21231accc5c2ed73980c9e3a7542538c6a41d7778de5cfbb410312adc`;
  47,632 stable entries verify and exactly one mutable telemetry counter still
  differs. Preserve the btree correction: its 24,197-entry regression manifest
  retains SHA-256
  `9cca917b12dd761c58bf91652e78b3e999f55eeb1ffc39f9f3bbf56a52aef463`
  and exactly 163 recorded mutable cache/HOME mismatches still reproduce;
  later evidence supersedes them.
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

Independently evaluate selected `gopkg.in/errgo.v2 v2.1.0` as exactly one
bounded P7 module group. Resolve canonical latest, release and pseudo-version
history, Go declarations and complete floor-compatible closure, exact source
identity, signature and repository state, loaded package population and real
consumers, module self-tests, and primary vulnerability data from fresh
evidence before selecting anything.

The current build list and graph measurements remain 234 selected modules,
3,564 graph edges, 429 complete packages, 1,031 go.sum lines, a 332-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Use exact Go 1.26.7 with GOENV off,
GOWORK off, GOTOOLCHAIN local, and no ambient GOFLAGS.

Implement only an exact floor-compatible changed selection with an explained
minimal closure and every applicable gate passing. If canonical latest is a
pseudo-version, prove its authoritative commit/time and qualification rather
than treating a branch head as a release. If exact selected is already the
decision, record no-change without manufacturing a dependency commit. Stop on
any mandatory failure. Do not revisit Check, either gopkg YAML group, the
accepted `go.yaml.in/yaml/v3` group, another dependency, or P8.
