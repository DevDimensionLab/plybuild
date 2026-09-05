# Quality Upgrade Handover

Generated: 2026-09-05T23:12:14+02:00

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
- The incoming Resty v1 documentation session is commit
  `a855007ace19f1fc96302000a013cd7d802408ea`, exact parent
  `1899ff708ae165c43c283e97def420af2c11d108`, tree
  `c2293adc2d5e2cd33d1384430e111d7df4c5d5ab`. This handoff must be its
  direct child. The next dependency implementation, if any, must use the
  resulting documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Their minimal
  closures select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty
  v0.3.1, and rogpeppe/go-internal v1.9.0. Those four accepted groups add
  exactly 15 checksum lines.
- The answered Resty v1 archive and sole NEXT
  `gopkg.in/alecthomas/kingpin.v2` archive link reciprocally. Only the
  launcher's mutable header and prompt regions change during handoff. Ordinary
  and ignored status must end empty.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the maintained Go 1.26.7
toolchain and completed dependency groups through accepted cast v1.5.1,
rejected afero v1.10.0, retained canonical-latest jwalterweatherman v1.1.0,
rejected gotenv v1.6.0, x/image v0.16.0, x/net v0.25.0, x/text v0.15.0,
retained canonical-latest YAML v2 v2.4.0 and YAML v3 v3.0.1, rejected
canonical-latest Check v1 pseudo-version
`v1.0.0-20201130134442-10cb98267c6c`, retained canonical-latest Errgo v2
v2.1.0, and retained canonical-latest Resty v1 v1.12.0. All earlier recorded
dependency decisions and evidence corrections remain final. Do not revisit
them. P8 remains queued.

## Retained Gopkg Resty V1 Group

Fresh proxy evidence lists ten v1 semantic tags. Exact `go list` resolution
for `@latest`, `@v1`, `@master`, and v1.12.0 selects v1.12.0 at
2019-02-28T07:26:48Z, with no prerelease or retraction. V1.7.0 and v1.8.0
appear in the version list but fail exact download because their go.mod files
declare `github.com/go-resty/resty`. The older v1.0-v1.6 Git tags are
non-canonical shorthand versions and are absent from the module list. The
valid exact-path stable releases run from v1.9.0 through v1.12.0. Selected
v1.12.0 is canonical latest and the highest declaration-compatible stable
candidate for the retained Go 1.18 floor.

V1.12.0 is lightweight tag and commit
`fa5875c0caa5c260ab78acec5a244215a730247f`, tree
`5029acc2e860c8e9495d3b46fdc5d40d429808b8`, at
2019-02-28T07:26:48Z. GitHub reports its embedded commit signature verified
and valid. There is no separate tag object or tag signature, and local
cryptographic verification is unavailable because `gpg` is absent. The
selected checksum pair is
`h1:CuXP0Pjfw9rOuY6EP+UvtNvt5DSqHpIxILZKT/quCZI=` /
`h1:mDo4pnntr5jdWRML875a/NmxYqAlA73dVijT2AXvQQo=`; its proxy ZIP SHA-256 is
`43487bb0bb40626d16502b1fe9e719cf751e7a5b4e4233276971873e7863d3cf`.
The 30-file proxy archive and exact upstream tag manifest are byte-identical at
SHA-256
`9768c36ec0f94cf2d9a405cc5125e14e622a54fd83af175d21141683fd43e944`.

Authoritative gopkg metadata maps the exact module to
`https://gopkg.in/resty.v1` and source browsing to
`github.com/go-resty/resty/tree/v1.12.0`. The upstream repository is enabled,
unarchived, and not a fork. Its active default branch is v3, while its v1.x
branch and gopkg master both end exactly at v1.12.0. The v1.11-to-v1.12 range
has four commits, including `SetMultipartFields`; there are no later v1-branch
commits. The repo publishes a stable, non-draft v1.12.0 GitHub release and has
no module deprecation marker, though v1 has been dormant since 2019.

Exact resolution of pre-release commit
`0ecc38d58bec7c3a6c001c15e43a88e67e37ced1` yields the older unreleased
pseudo-version `v1.11.1-0.20190110224454-0ecc38d58bec`. Asking the v1 path for
v2 release commit `d467d573085e42a848cb49c1c28ba163b7aa5872` or current v3 commit
`2023fc4669a397817856115844d1dff89f6dede3` computes a later v1 pseudo-version
but rejects it because go.mod declares `github.com/go-resty/resty/v2` or
`resty.dev/v3`. Alternate v2.17.2 and v3.0.0-rc.3 both require Go 1.23 and use
different module paths; v3 remains a prerelease. They are not exact v1 upgrade
candidates.

V1.12.0 has no Go directive and requires only
`golang.org/x/net v0.0.0-20181220203305-927f97764cc3`, also without a Go
directive. The project already selects x/net v0.7.0 at Go 1.17. The
changed-selection closure is empty and preserves the Go 1.18 floor by
declaration, not by inference from modern tests.

An external exact `go get gopkg.in/resty.v1@v1.12.0` exits zero, but selected
module lists remain byte-identical at 234 entries. It projects only a redundant
indirect main-module requirement, one main-to-Resty graph edge, and the full
v1.12.0 checksum. Graph edges project 3,564 -> 3,565, go.sum lines
1,031 -> 1,032, and the unapplied tidy diff 332 -> 344 lines; complete packages
remain 429. Since no version selection changes, that explicit pin is not an
authorized implementation. Go.mod and go.sum remain unchanged and no
dependency commit was manufactured.

Exact complete-package enumeration finds zero Resty packages in either state,
repository Go source has zero Resty imports, and `go mod why -m` says the main
module does not need it. The only graph path is historical requirement debt
from mvn-pom-mutator v0.2.3, whose source also has zero Resty imports. Resty is
therefore not loaded by the main module or dependency tests.

The exact proxy source verifies, lists one package, and remains source-clean.
Its standalone two-module test apparatus needs only the declared x/net
requirement and two external checksums; it adds no test-only requirement and
is separate from the unchanged project closure. Count-1 and full count-2 tests
reproduce six failures: `TestClientRedirectPolicy`,
`TestClientRetryWithSetContext`, `TestNoAutoRedirect`,
`TestHTTPAutoRedirectUpTo10`, `TestIncorrectURL`, and `TestClientRetryGet`.
They assert legacy pre-modern net/http and net/url error strings, including
unquoted URLs. Race reproduces the assertions without a data-race diagnostic;
vet passes. This mandatory module-suite failure rejects any metadata projection
even though v1.12.0 remains the exact selected canonical version.

Old and projected repository verification, build, complete tests/race/vet,
Windows build, pinned lint, CLI surface, binary build, and API/CLI checks all
pass; reports and root/status/upgrade/build help remain byte-identical.
Empty-HOME count-2 and the direct launcher suite pass. A disposable committed
projection was required for audit meta-tests to receive a clean repository,
preserving the established preflight-runner correction. One primary launcher
signal-fixture timing attempt failed closed; the direct and final Make suites
then pass 62/62. Final full preflight exits zero with all 62 launcher checks,
80 mutation controls, and 15 audit meta-controls passing. The changed-selection
quality, snapshot/Docker, and focused/full audit gates are inapplicable because
no version changed and the mandatory module gate failed.

Govulncheck v1.7.0 uses a primary database updated
2026-09-02T19:12:04Z. The fresh 1,392-entry module index has no exact
`gopkg.in/resty.v1` record, finding, or trace. Its only Resty record is
GO-2023-2328 on the alternate v2 path. Old and projected outputs and exact ID
sets are byte-identical for Darwin symbol, Darwin module, and Windows symbol
scans, preserving populations 20/30/20.

## Resty V1 Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-resty-v1-selection.a855007.Gg7QC9`. Its fully verified
  110,457-entry manifest SHA-256 is
  `72d739b05c507fab058737de6ab967b534e5fcdab1c9588582188532f1a926c2`;
  decision-summary SHA-256 is
  `008b8b9396b58c93e696a5d495037ac974224149fce9b34badb6c7786723a109`.
- The manifest covers 33 prerequisite-root verifications, fresh
  proxy/sumdb/gopkg/Git/GitHub resolution, complete v1 tag history, exact-path
  and alternate-path qualification, source/signature/repository identity,
  Go-floor declarations, exact-get/tidy closure, loaded-population
  classification, module rejection tests, old/projected repository gates,
  compatibility/help identity, full preflight, empty-HOME, and byte-identical
  primary vulnerability scans.

## Inherited Evidence And Tool Identity

- All 33 prerequisite manifest roots were independently recomputed before
  measurement: the 32 roots recorded by Errgo plus the Errgo root itself. The
  exact PASS table has SHA-256
  `f0ff7a7f14d209438f61e8fc127825891928d99711c7ccf9ab9d05b7e1234017`
  and is sealed in Resty evidence.
- Errgo v2 no-change evidence remains fully verified at
  `/private/tmp/ply-p7-errgo-v2-selection.1899ff7.HoEmPo`,
  152,333/`e9a3c3b1b9c311a1dd96f1804167d57359c2af02615dde2ffb544690dac45ab2`;
  its decision summary is
  `aa8a30c6f9e261a5c29be3bd59d87cb96e5b07f85582b6f1d55a4f78fcf3a3ec`.
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

Independently evaluate selected `gopkg.in/alecthomas/kingpin.v2 v2.2.6` as
exactly one bounded P7 module group. Resolve canonical latest,
stable/prerelease and pseudo-version qualification, exact-path and
alternate-path history, Go declarations and complete floor-compatible closure,
source identity, signature and repository state, loaded package population and
real consumers, module self-tests, and primary vulnerability data from fresh
evidence before selecting anything.

The current build list and graph measurements remain 234 selected modules,
3,564 graph edges, 429 complete packages, 1,031 go.sum lines, a 332-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Use exact Go 1.26.7 with GOENV off,
GOWORK off, GOTOOLCHAIN local, and no ambient GOFLAGS.

Implement only an exact floor-compatible changed selection with an explained
minimal closure and every applicable gate passing. If a later branch or module
path exists, prove whether it qualifies for exact
`gopkg.in/alecthomas/kingpin.v2` rather than treating it as an in-place
release. If exact selected is already the decision, do not manufacture an
explicit requirement or dependency commit. Stop on any mandatory failure. Do
not revisit Resty, Errgo, Check, either gopkg YAML group, the accepted
`go.yaml.in/yaml/v3` group, another dependency, or P8.
