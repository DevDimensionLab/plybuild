# Quality Upgrade Handover

Generated: 2026-09-05T19:23:59+02:00

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
- The incoming YAML v2 documentation session is commit
  `1969442891c30b5750da632032312cb1b96624e9`, exact parent
  `6ab5945de5fef2ee6f7deff10f4bc51e94596b30`, tree
  `4aec00f4704f56b26262ee0662afd3d20615c28d`. This handoff must be its direct
  child. The next dependency implementation, if any, must use the resulting
  documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Their minimal
  closures select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty
  v0.3.1, and rogpeppe/go-internal v1.9.0. Those four accepted groups add
  exactly 15 checksum lines.
- The answered YAML v2 archive and sole NEXT `gopkg.in/yaml.v3` archive must
  link reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff. Ordinary and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1,
rejected afero v1.10.0, retained canonical-latest jwalterweatherman v1.1.0,
rejected gotenv v1.6.0, x/image v0.16.0, x/net v0.25.0, and x/text v0.15.0,
and retained canonical-latest YAML v2 v2.4.0. Go-cmp v0.7.0, Viper v1.16.0,
Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, latest
gomarkdown `v0.0.0-20260824154242-13c5cf49db8d`, HCL, mousetrap, properties,
mapstructure, go-toml v1, afero, jwalterweatherman, gotenv, x/image, x/net,
x/text, and YAML v2 changes remain rejected or unnecessary for their recorded
floor, loaded-behavior, release-qualification, self-test, exact-latest, or
closure decisions. Do not revisit them. P8 remains queued.

## Retained Gopkg YAML V2 Group

Fresh Go proxy, checksum-database, gopkg metadata, original and successor Git
repositories, GitHub API, and primary Go vulnerability evidence enumerate
exactly 14 stable releases, v2.0.0 through v2.4.0, with no prereleases or
retractions. Proxy `@latest` and exact `go list` resolve v2.4.0 at
2020-11-17T15:46:20Z. It declares Go 1.15 and requires
`gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405`. No earlier release
declares above Go 1.18, so v2.4.0 is both canonical latest and the highest
stable release compatible with the retained floor for this exact module path.
Its checksum pair is `h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=` /
`h1:RDklbk79AGWmwhnvt/jBztapEOGDOx6ZbXqjP6csGnQ=`.

V2.4.0 is lightweight tag commit
`7649d4548cb53a614db133b2a8ac1f31859dda8c`, tree
`32e8d4cd33ca7feab3e2b4081202299f7556af4e`. The tag has no separate tag
object or signature, and the commit is unsigned. All 24 proxy files match the
exact original `go-yaml/yaml` tag and successor `yaml/go-yaml` mirror tag;
proxy ZIP SHA-256 is
`ede49e27c4cca6cdd2ec719aed8ea4d363710cceb3d411e7a786fbdec0d391fd`.
The original repository is archived but not disabled or a fork, has no commit
after v2.4.0 on its v2 branch, publishes no GitHub Release objects, and places
no formal deprecation marker on the module.

Successor tags v2.4.1-v2.4.4 belong to different module path
`go.yaml.in/yaml/v2`; `gopkg.in/yaml.v2@v2.4.1` is an unknown revision.
`gopkg.in/yaml.v3` v3.0.1, `go.yaml.in/yaml/v3` v3.0.5, and prerelease-only
`go.yaml.in/yaml/v4` v4.0.0-rc.6 are distinct module paths or majors and were
not candidates. Relevant v2.3.0-v2.4.0 history contains the public
`FutureLineWrap` change and line-wrap revert, Go 1.14 CI, and Marshal docs;
unreleased successor history does not change the exact-path decision.

Exact projected `go get gopkg.in/yaml.v2@v2.4.0` produces empty stdout/stderr
and zero-byte go.mod/go.sum diffs. Both states retain 234 selected modules,
3,564 graph edges, 429 complete packages, 1,031 go.sum lines, and an identical
332-line unapplied tidy projection. Selected
`gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15` already dominates YAML
v2's declared 2016 pseudo-version, so the changed MVS closure has zero
selections, edges, or checksums. Project module verification passes.

Exactly one YAML v2 package loads through
`github.com/devdimensionlab/plybuild/pkg/config -> gopkg.in/yaml.v2`.
`pkg/config` uses `yaml.Marshal` and `yaml.Unmarshal` for local and cloud
configuration. Focused `go test ./pkg/config -count=10` passes in both states
with normalized-identical output.

The module source remains byte-identical, verifies, lists, passes complete
count-1 tests, and passes race count 1. Both independent old and no-op
candidate replays fail count 10 with nine `S.TestLineWrapping` failures:
`yaml.FutureLineWrap()` permanently changes package-global state after the
first iteration. Both vet replays fail with exactly 27 identical legacy
malformed struct-tag diagnostics in decode and encode tests. These mandatory
module failures stop the group before repository-wide build, tests/race/vet,
lint, help/API/CLI, snapshot/Docker, quality, and audit acceptance. No
dependency metadata or implementation commit was made.

Govulncheck v1.7.0 preserves exact Darwin-symbol/Darwin-module/Windows-symbol
populations at 20/30/20 with identical IDs and no YAML v2 module finding or
symbol trace. The fresh 1,392-entry primary index contains exactly three YAML
v2 records: GO-2020-0036 fixed in v2.2.8, GO-2021-0061 fixed in v2.2.3, and
GO-2022-0956 fixed in v2.2.4. Selected v2.4.0 is after every fix.

## YAML V2 Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-yaml-v2-selection.1969442.0BpgSy`. Its fully verified
  29,219-entry manifest SHA-256 is
  `8dd34cc54f56e4bb1370b2f7f1f124aa8b501a4dcca177effe4f6377c3919e14`;
  decision-summary SHA-256 is
  `61d8132e932d4613aea071c5f3d55a9a38c2873e2428011b83231611517885d6`.
- The manifest covers fresh proxy/sumdb/gopkg/Git/GitHub release history, all
  14 releases and declarations, exact source identity, tag/commit signature
  status, archive/deprecation state, related-path exclusion, exact-get and
  tidy diffs, selection/edge/package/checksum populations, loaded consumers
  and focused behavior, independent module test replays, source-mutation
  checks, and exact primary vulnerability scans. Every entry was recomputed
  successfully after sealing.

## Inherited Evidence And Tool Identity

- All 28 complete inherited manifest roots were independently verified at
  session start, including immediate x/text plus x/net, x/image, gotenv,
  jwalterweatherman, cast, afero, go-toml/v2, go-toml v1, mapstructure,
  runewidth, properties, mousetrap, HCL, go-colorful, btree, and toolchain
  roots. The verification table is sealed inside the YAML v2 evidence.
- X/text rejection evidence remains fully verified at
  `/private/tmp/ply-p7-x-text-selection.6ab5945.3qn5Bj`,
  82,629/`48c872da4ab796ef1115003fcf0a226a07bf5791a163d36bc9c946c9d9f92ed6`;
  decision summary is
  `50d6bee90ea7cea060e23400e820b1537742ec2167bee04f188070e2e5eb39c9`.
- X/net rejection evidence remains fully verified at
  `/private/tmp/ply-p7-x-net-selection.e5ea7ae.cP57X8`,
  71,887/`88e4be87483322792ad3da0064bcef4ced5ee21c52821bb1f80c143f80a3aaf1`;
  decision summary is
  `6e32698abcde67604e159488ba21f6e780fd61dfba2ed9299644b5f011940dff`.
- X/image rejection evidence remains fully verified at
  `/private/tmp/ply-p7-x-image-selection.9e9ef26.9J9x1F`,
  52,533/`84972465aef5f19f888102c7991aadd67446146a7056658b006d005666bd82bf`;
  decision summary is
  `2b70234a63a199ab6cf11050212645162a78c4836aa9aafa4cb00dad764eb936`.
- Gotenv rejection evidence remains fully verified at
  `/private/tmp/ply-p7-gotenv-selection.3a8cadc.Dai0IG`,
  48,662/`ab2bd356b863d529085cff9b918d004aa2a03a121a7ea360a66767a2820cc88c`;
  decision summary is
  `014c050a3ccefd6fa501f4c5d04aded2ee93d02ad4039a9935d6f4d478c21e11`.
- Jwalterweatherman no-change evidence remains fully verified at
  `/private/tmp/ply-p7-jwalterweatherman-selection.2dd8fdf.h4EKeQ`,
  76,374/`54040aefec38df7415f56da66af04d48c8256dd46ae6de8e06a7c4bbe5ce4d8e`;
  decision summary is
  `342330b73509caf414a406dda3a14c026bfc59512d09bf6393f710127974dc52`.
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
- Tool identities remain Go 1.26.7
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
  golangci-lint 2.12.2
  `3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`,
  GoReleaser 2.17.1
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`,
  and govulncheck v1.7.0
  `0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

## Next Objective

Independently evaluate selected `gopkg.in/yaml.v3 v3.0.1` as exactly one
bounded P7 module group. Resolve canonical latest, every stable release that
could satisfy the retained Go 1.18 floor through its complete changed closure,
source identity, release qualification, requirements, loaded package
population and behavior, module self-tests, and primary vulnerability data
from fresh evidence before selecting anything.

Treat the already accepted `go.yaml.in/yaml/v3 v3.0.5` decision as final and
the path as distinct. Explicitly distinguish gopkg YAML v2, the successor
`go.yaml.in/yaml/v3` path, prerelease YAML v4, forks, and unreleased commits;
none is an upgrade candidate for the gopkg YAML v3 group. Measure exact
old/candidate modules, edges, packages, checksums, exact-get and tidy diffs,
loaded paths and symbols, focused consumer behavior, module and repository
quality, help/API/CLI identity, and exact Darwin and Windows vulnerability
populations. Implement only an exact floor-compatible changed release with an
explained minimal closure and every gate passing; otherwise record rejection
or no-change without dependency metadata edits. Do not combine YAML v2,
go.yaml.in YAML v3, YAML v4, another dependency group, or P8.
