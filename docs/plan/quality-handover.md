# Quality Upgrade Handover

Generated: 2026-09-05T20:06:53+02:00

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
- The incoming YAML v3 documentation session is commit
  `cdb5b4ca02fa5db54efebf741d2b415608c053b2`, exact parent
  `1969442891c30b5750da632032312cb1b96624e9`, tree
  `9a46386d766eb0a46ffc94dc398850b228fe3654`. This handoff must be its direct
  child. The next dependency implementation, if any, must use the resulting
  documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Their minimal
  closures select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty
  v0.3.1, and rogpeppe/go-internal v1.9.0. Those four accepted groups add
  exactly 15 checksum lines.
- The answered YAML v3 archive and sole NEXT `gopkg.in/check.v1` archive must
  link reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1,
rejected afero v1.10.0, retained canonical-latest jwalterweatherman v1.1.0,
rejected gotenv v1.6.0, x/image v0.16.0, x/net v0.25.0, and x/text v0.15.0,
and retained canonical-latest YAML v2 v2.4.0 and YAML v3 v3.0.1. Go-cmp
v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0,
fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, HCL, mousetrap, properties,
mapstructure, go-toml v1, afero, jwalterweatherman, gotenv, x/image, x/net,
x/text, YAML v2, and YAML v3 changes remain rejected or unnecessary for their
recorded floor, loaded-behavior, release-qualification, self-test,
exact-latest, or closure decisions. Do not revisit them. P8 remains queued.

## Retained Gopkg YAML V3 Group

Fresh Go proxy, checksum-database, gopkg metadata, original and successor Git
repositories, GitHub API, and primary Go vulnerability evidence enumerate
exactly two stable exact-path releases, v3.0.0 and v3.0.1, with no prereleases
or retractions. Proxy `@latest` and exact `go list` resolve v3.0.1 at
2022-05-27T08:35:30Z. Both releases omit a `go` declaration and require
`gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405`, which also omits a
Go declaration. V3.0.1 is therefore both canonical latest and the highest
stable release compatible with the retained Go 1.18 floor for this exact
module path. Its checksum pair is
`h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=` /
`h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=`.

V3.0.1 is lightweight tag commit
`f6f7691b1fdeb513f56608cd2c32c51f8194bf51`, tree
`1cb2e60a039c6b3cfdfbd33cc2d049bf68c7db00`. The tag has no separate tag
object or signature, and the commit is unsigned. All 24 proxy files match the
exact archived original `go-yaml/yaml` tag and active `yaml/go-yaml` successor
mirror tag; proxy ZIP SHA-256 is
`aab8fbc4e6300ea08e6afe1caea18a21c90c79f489f52c53e2f20431f1a9a015`.
Both repositories report un-disabled, non-fork GitHub metadata, while the
successor describes itself as the maintained YAML-org fork. The original v3
branch has one later unreleased README-only commit marking the project
unmaintained. The selected module has no formal deprecation marker, and both
repositories have zero GitHub Release objects.

Successor v3.0.2-v3.0.5 tags use distinct path `go.yaml.in/yaml/v3`; that path
retracts the invalid inherited v3.0.0-v3.0.1 tags and retains already accepted
v3.0.5. Retained `gopkg.in/yaml.v2 v2.4.0` is a different major, and
`go.yaml.in/yaml/v4` exposes only v4.0.0-rc.1 through rc.6. None is a candidate
for the exact gopkg YAML v3 group. The v3.0.0-v3.0.1 delta is one commit,
adding nil-token parser guards and an invalid-input regression. The primary
GO-2022-0603 fix commit `8f96da9f5d5e` precedes stable v3.0.0.

Exact projected `go get gopkg.in/yaml.v3@v3.0.1` produces empty stdout/stderr
and zero-byte go.mod/go.sum diffs. Both states retain 234 selected modules,
3,564 graph edges, 429 complete packages, 1,031 go.sum lines, and an identical
332-line unapplied tidy projection. Selected
`gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15` already dominates the
declared 2016 pseudo-version, so the changed MVS closure has zero selections,
edges, or checksums. Project module verification passes.

Exactly one YAML v3 package loads through `plybuild/cmd -> pkg/config ->
gopkg.in/yaml.v3`. `cmd/tips.go` and `cmd/build_options.go` call
`GitCloudConfig.GlobalCloudConfig`, which uses `yaml.Unmarshal`. An
external-only focused fixture exercises that production path and validates the
decoded fields and `SourceFor`; ten old and candidate repetitions pass with
normalized-identical output.

The module source remains byte-identical, verifies, lists one package, and
passes complete count-1, count-10, and race tests in repeated independent old
and no-op-candidate runs. Both vet replays fail with exactly 32 identical
legacy malformed struct-tag diagnostics across decode, encode, and node tests.
This mandatory module gate stops the group before repository build, complete
tests/race/vet, pinned lint, help/API/CLI, snapshot/Docker, quality, and audit
acceptance. No dependency metadata or implementation commit was made.

Govulncheck v1.7.0 preserves byte-identical Darwin-symbol/Darwin-module/
Windows-symbol populations at 20/30/20 with no YAML v3 finding or trace. The
fresh 1,392-entry primary index contains exactly one YAML v3 record,
GO-2022-0603, fixed at
`v3.0.0-20220521103104-8f96da9f5d5e`; both stable releases follow it.

## YAML V3 Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-yaml-v3-selection.cdb5b4c.bbbTdr`. Its fully verified
  40,121-entry manifest SHA-256 is
  `84bd12b7b4cbff806abcbff213b4b7bc5bac230a7371c3fe0b3e2def643ed741`;
  decision-summary SHA-256 is
  `95d5aa003b77880ecf72770ad66c061e0013dfd022fc222eb600860bc9a7799f`.
- The manifest covers fresh proxy/sumdb/gopkg/Git/GitHub release history, both
  stable releases and declarations, exact source identity across both
  repositories, tag/commit signature status, archive/unmaintained/deprecation
  state, distinct-path and prerelease exclusions, exact-get and tidy diffs,
  selection/edge/package/checksum populations, loaded consumers and focused
  behavior, repeated independent module gates, source-mutation checks, and
  byte-identical primary vulnerability scans.

## Inherited Evidence And Tool Identity

- All 30 prerequisite manifest roots were independently recomputed before
  measurement: the 28 inherited roots recorded by YAML v2, the YAML v2 root
  itself, and the recovery root. The exact PASS table is sealed in YAML v3
  evidence.
- YAML v2 no-change evidence remains fully verified at
  `/private/tmp/ply-p7-yaml-v2-selection.1969442.0BpgSy`,
  29,219/`8dd34cc54f56e4bb1370b2f7f1f124aa8b501a4dcca177effe4f6377c3919e14`;
  decision summary is
  `61d8132e932d4613aea071c5f3d55a9a38c2873e2428011b83231611517885d6`.
- X/text, x/net, x/image, gotenv, and jwalterweatherman selection roots remain
  fully verified respectively at
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

Independently evaluate selected
`gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15` as exactly one bounded
P7 module group. Resolve canonical latest, release and pseudo-version history,
Go declarations and complete floor-compatible closure, exact source identity,
signature and repository state, loaded package population and real consumers,
module self-tests, and primary vulnerability data from fresh evidence before
selecting anything.

The current build list and graph measurements remain 234 selected modules,
3,564 graph edges, 429 complete packages, 1,031 go.sum lines, a 332-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Use exact Go 1.26.7 with GOENV off,
GOWORK off, GOTOOLCHAIN local, and no ambient GOFLAGS. Determine independently
whether gocheck is loaded by the main module or only selected by dependency
test requirements; do not infer behavior from YAML's declared edge.

Implement only an exact floor-compatible changed selection with an explained
minimal closure and every applicable gate passing. If canonical latest is a
pseudo-version, prove its authoritative commit/time and qualification rather
than treating a branch head as a release. If exact selected is already the
decision, record no-change without manufacturing a dependency commit. Stop on
any mandatory failure. Do not revisit either gopkg YAML group, the accepted
`go.yaml.in/yaml/v3` group, another dependency, or P8.
