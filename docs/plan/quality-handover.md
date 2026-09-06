# Quality Upgrade Handover

Generated: 2026-09-06T04:47:30+02:00

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
- The incoming Kingpin v2 documentation session is commit
  `c44013020b1622cae0fdada815c445e311b7b86f`, exact parent
  `a855007ace19f1fc96302000a013cd7d802408ea`, tree
  `d6fef9dbd58de24b0313e566fc7dba6b8e2f5a81`. This handoff must be its
  direct child. The next dependency implementation, if any, must use the
  resulting documentation commit as its exact parent.
- Relative to accepted go-cmp commit `c314bcb`, accepted dependency metadata
  moves remain go-colorful v1.2.0 -> v1.4.1, go-runewidth v0.0.14 -> v0.0.17,
  go-toml/v2 v2.0.7 -> v2.2.2, and cast v1.5.0 -> v1.5.1. Their minimal
  closures select testify v1.9.0, objx v0.5.2, quicktest v1.14.4, kr/pretty
  v0.3.1, and rogpeppe/go-internal v1.9.0. Those four accepted groups add
  exactly 15 checksum lines.
- The answered Kingpin v2 archive and sole NEXT
  `github.com/alecthomas/units` archive link reciprocally. Only the
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
v2.1.0, retained canonical-latest Resty v1 v1.12.0, and retained highest
stable exact-path Kingpin v2 v2.2.6. All earlier recorded dependency decisions
and evidence corrections remain final. Do not revisit them. P8 remains queued.

## Retained Gopkg Kingpin V2 Group

Fresh proxy evidence lists exactly 40 stable v2 semantic versions, with no
prereleases or retractions. Proxy and exact Go resolution identify v2.4.0 at
2023-09-30T22:59:49Z as canonical latest. It is not an in-place upgrade for
the selected path: v2.4.0 declares `github.com/alecthomas/kingpin/v2`, and
exact `go get gopkg.in/alecthomas/kingpin.v2@v2.4.0` fails. V2.3.0 has an
invalid non-`.v2` GitHub declaration; v2.3.1-v2.4.0 use the alternate `/v2`
path. Proxy module metadata for the 36 stable tags through selected v2.2.6
declares the exact gopkg path.
V2.2.6 is therefore the highest qualified stable exact-path release and has
no Go directive or requirements, so its empty changed closure preserves Go
1.18 by declaration.

V2.2.6 is lightweight tag and unsigned commit
`947dcec5ba9c011838740e680966fd7087a71d0d`, tree
`7b5bf03129479f410ca39ea2647dd190bddee927`, at
2017-12-17T18:08:38Z. Its checksum pair is
`h1:jMFz6MfLP0/4fUyZle81rXUoxOBFi19VUFKVDOQfozc=` /
`h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=`. Proxy ZIP SHA-256 is
`638080591aefe7d2642f2575b627d534c692606f02ea54ba89f42db112ba8839`;
all 39 proxy files match the exact tag at manifest SHA-256
`5d4858f09550b4dca7afb2412d315025a18e2f48b613e9267f1561879ec5e908`.
All 40 v2 tags are lightweight commit refs, so none has a tag-object
signature; local commit status is `N` throughout. GitHub reports v2.2.6
unsigned and v2.4.0's commit signature valid.

Authoritative gopkg metadata maps the selected path to its gopkg Git endpoint
and browses the v2.4.0 GitHub tree. The gopkg mirror pins `master` to v2.4.0,
whereas current GitHub master is commit `177e1b9ba430164dc4d48fe6f6d7613ccebfe629`
at 2026-08-24T07:50:00Z, 18 commits later. The repository is enabled,
unarchived, and not a fork, but its README says `CONTRIBUTIONS ONLY` and names
the GitHub `/v2` module as current stable. The exact gopkg module has no
deprecation marker. Its only v2 GitHub Release object is stable v2.4.0.

The exact path also resolves three unreleased historical snapshots:
`v2.2.7-0.20181031024914-c2ca6a1e4f86`, module-conversion
`v2.2.7-0.20181107222045-102f372a17d4`, and divergent `v3-unstable` head
`v2.1.12-0.20191105091915-95d230a53780`. They download, but none is a stable
release; the last is below v2.2.6's semantic line. Current GitHub master forms
alternate-path pseudo-version
`v2.4.1-0.20260824075000-177e1b9ba430`; exact gopkg-path get rejects its
GitHub `/v2` declaration. No pseudo-version qualifies as the stable candidate.

External exact `go get gopkg.in/alecthomas/kingpin.v2@v2.2.6` exits zero but
changes no selected version. It adds only redundant indirect requirements on
already-selected template, units, and Kingpin modules, three main graph edges,
and their three full checksum lines. Selected modules remain 234 and complete
packages remain 429; graph edges project 3,564 -> 3,567, go.sum lines
1,031 -> 1,034, and the unapplied tidy diff 332 -> 353 lines. The changed
selection closure is empty. An explicit redundant pin is not an authorized
selection change, so go.mod/go.sum remain unchanged and no dependency commit
was manufactured.

Complete-package enumeration finds zero Kingpin packages, repository Go source
has zero imports, and `go mod why -m` says the main module does not need it.
Only old Prometheus tsdb/common graph edges select Kingpin, reached through
mvn-pom-mutator. It is historical MVS graph debt and is not loaded in the
complete project population.

The exact proxy source and upstream tag remain byte-identical and source-clean
after tests. Their separate seven-module test apparatus is not project MVS.
Count-1, ten fresh count-1 processes, race, and vet pass. Both exact-source
replays fail `go test ./... -count=10`: repetitions 2-10 fail
`TestRequiredArgWithEnvarMissingErrors` and
`TestRequiredWithEnvarMissingErrors`, 18 failures total, because sibling tests
set `TEST_ARG_ENVAR` and `TEST_ENVAR` without unsetting process-global state.
The mandatory repeated-suite stop rule independently rejects even the
redundant exact-selected projection.

Old and projected repository verification, build, complete tests/race/vet,
Windows build, pinned lint, CLI surface, API/CLI reports, and public help pass
or remain byte-identical. A reproducible `-trimpath -buildvcs=false` binary is
also byte-identical. Empty-HOME count-2 passes. The corrected full preflight
uses the committed projection plus external API/CLI reports and exits zero
with all 62 launcher checks, 80/80 killed mutants, and 15 audit meta-controls.
Earlier retained attempts expose the established signal-retention timing
fixture and the external-report cleanliness requirement; neither changes the
final passing result. Changed-selection quality, snapshot/Docker, and focused/
full audit gates are inapplicable because no version changed and the mandatory
module gate failed.

Govulncheck v1.7.0 uses the primary database updated
2026-09-02T19:12:04Z. Its fresh 1,392-entry module index has no Kingpin
record, finding, or trace. Old and projected normalized outputs and exact IDs
are byte-identical for Darwin symbol, Darwin module, and Windows symbol scans,
preserving populations 20/30/20.

## Kingpin V2 Evidence

- Selection evidence is sealed at
  `/private/tmp/ply-p7-kingpin-v2-selection.c440130.2WGm6X`. Its fully verified
  97,185-entry manifest SHA-256 is
  `ce50700fe0df4e208ce1a679d7ef8c9b97fdca45398843d392351673fa70837c`;
  decision-summary SHA-256 is
  `d9080ecc8983da371aacfe1ef234a1290e7e577738c4b5dc9d99078365bba54f`.
- The manifest covers all inherited-root verification, fresh proxy/sumdb/
  gopkg/Git/GitHub evidence, the complete 40-tag history, pseudo- and alternate-
  path qualification, source/signature/repository identity, declaration-based
  floor proof, exact-get/tidy closure, loaded population, module rejection,
  old/projected repository gates, public compatibility identity, corrected
  full preflight, empty-HOME, and primary vulnerability identity.

## Inherited Evidence And Tool Identity

- All 34 immutable prerequisite roots were independently recomputed before
  measurement: Resty's 33-root prerequisite table plus Resty itself. The two
  recorded mutable correction roots were also replayed with their exact
  expected mismatch populations. The clean 36-row PASS table has SHA-256
  `6dfb19ebfb58dc83001427bd49fc59915938fbecdc5d18c7ba48b8d1d8a2e41f`
  and is sealed in Kingpin evidence.
- Resty v1 no-change evidence remains fully verified at
  `/private/tmp/ply-p7-resty-v1-selection.a855007.Gg7QC9`,
  110,457/`72d739b05c507fab058737de6ab967b534e5fcdab1c9588582188532f1a926c2`;
  its decision summary is
  `008b8b9396b58c93e696a5d495037ac974224149fce9b34badb6c7786723a109`.
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

Independently evaluate selected
`github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4` as exactly one
bounded P7 module group. The fresh update survey reports
`v0.0.0-20240927000941-0f3dac36c52b`, but treat canonical latest, stable-tag
and pseudo-version qualification, branch history, declarations and complete
floor-compatible closure, source identity, signatures, repository state,
loaded population, real consumers, self-tests, and vulnerability data as
unknown until independently proved from primary evidence.

The current build list and graph measurements remain 234 selected modules,
3,564 graph edges, 429 complete packages, 1,031 go.sum lines, a 332-line
unapplied tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/
Windows-symbol vulnerability populations. Use exact Go 1.26.7 with GOENV off,
GOWORK off, GOTOOLCHAIN local, and no ambient GOFLAGS.

Implement only an exact floor-compatible changed selection with an explained
minimal closure and every applicable gate passing. Do not call a pseudo-version
a stable release or select an unreleased branch merely because it is newer. If
the selected version is already the exact decision, do not manufacture an
explicit requirement or dependency commit. Stop on any mandatory failure. Do
not revisit Kingpin, Resty, Errgo, Check, either gopkg YAML group, another
dependency, or P8.
