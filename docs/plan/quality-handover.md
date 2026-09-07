# Quality Upgrade Handover

Generated: 2026-09-07T18:34:13+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree /Users/perottochristensen/github/ply/upgrade-quality, branch
  codex/upgrade-quality, base master at 5635d50.
- Latest dependency implementation remains Repr commit
  6f4d02eb9c86ec2df8a488a85ed973aababe1f38, exact parent
  53475076e1c79f6d2181877e3238a1a5d389c246, tree
  82e7b1f5c503659082207481b8339e8113e38c10, changing only go.mod and go.sum.
- Operator-authorized lifecycle repair
  f9f0f7669e635c8c7bb169aab816d0ebe1b16435 remains intentionally between
  that implementation and direct-child Repr handoff
  342c7ece82eae3cd726aa658462089cfb876fbe0. Preserve it.
- The answered Template archive and sole NEXT Optional archive link
  reciprocally. Relative to accepted go-cmp commit c314bcb, accepted
  dependency metadata still adds exactly 27 checksum lines. Chroma, Colour,
  and Template changed no metadata.
- No .agent-task/current.md or .quality/manual-evidence.json was created.
  No push, merge, publication, release, stash, revert, successor launch, or
  worktree removal occurred.

## Lifecycle And Retained Roadmap

Every disposable cache, projection, external source, report, evidence tree,
and build context must remain beneath $CODEX_SESSION_SCRATCH_ROOT. Preserve
the launcher lifecycle controls and scratch-only mutation, snapshot, Docker,
and acceptance cleanup.

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Template's selected pseudo-version. All
earlier accepted, rejected, and no-change decisions and evidence corrections
remain final. P8 remains queued.

## Retained Alecthomas Template Group

Retain exact-path github.com/alecthomas/template
v0.0.0-20190718012654-fb15b899a751 without a dependency edit. The exact-path
proxy list is empty: there is no stable version or prerelease. Proxy @latest
and exact Go @latest and @master select the 2019 pseudo-version; @v0 reports
no match. There are no retractions. This is canonical latest and the highest
exact-path version compatible by declaration with Go 1.18, but it is neither a
stable release nor quality-qualified under current tests.

The authoritative public repository is enabled, unarchived, undisabled, and
non-fork. Its sole branch master remains at commit
fb15b899a75114aa79cc930e33c46b577cc664b1, tree
9658e953ba71f92dcf44f2d39cc5f90a27a0b88b, time
2019-07-18T01:26:54Z. The repository has 11 commits, zero tags, zero GitHub
Releases, and no deprecation marker. The only newer fetched ref is closed,
unmerged Renovate pull request 8 at
56c872bab4cb136e6659445e77c6934f09515a6e. It is not a branch head or release
candidate. /v2, /v3, and gopkg.in alternate-path probes find no module.

There is no tag object or tag signature. GitHub reports the selected commit's
embedded web-flow signature valid; independent detached gpgv verification
succeeds with fingerprint
5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23. Historical commit
a0175ee3bccc567396460bf5acd36800cb10c49c is unsigned. The selected commit
only adds one-line go.mod relative to that 2016 commit; every source and test
file is unchanged.

The selected checksum pair is
h1:JYp7IbQjafoB+tBA3gMyHYHrpOtNuDiK/uB5uXxq5wM= /
h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=. The old pseudo-version
has h1:cAKDfWh5VpdgMhJosfJnn5/FoN2SRZ4p7fJNX58YPaU= and the same go.mod
checksum. Sumdb agrees. Selected/old ZIP SHA-256 values are
25e3be7192932d130d0af31ce5bcddae887647ba4afcfb32009c3b9b79dbbdb3 and
86de3337a475e323a0fb54ef03386a4e495682201f42795bd7be646c05298692.
All 22 selected proxy regular files match upstream; normalized manifest
SHA-256 is
057c4c0d5af3039d317b83daee33c29ba753f4efd2afba9d0eb5003e6fb56cc7.

Raw selected metadata declares only the exact module path, with no Go
directive and no requirements. Its complete declared closure is Template
alone and cannot raise Go 1.18. Go 1.26.7's implicit Go 1.16 version when
testing this source as a main module is not upstream metadata. The module has
two packages and no test-only external requirement.

Writable proxy and exact-commit source independently verify and list, remain
unchanged after testing, and produce identical normalized results. Mandatory
complete count-1, count-10, and race tests all fail in TestJSEscaping; vet
passes. The test expects U+FDFF escaped, but the implementation follows
current unicode.IsPrint and emits it literally. Count-1/race normalized
diagnostic SHA-256 is
5d9c2becd6177fcb0227e4d8052c7d55ade12e3d2fab24a125dea9169cffdeb5;
count-10 is
06b388a002e0f7e9891b8f3fa28023af002ab173cf2d29dcadbdc6eda87e56ef.
This mandatory module-self-test failure triggers the stop rule.

Exact external selected-version get exits zero and changes no selected module.
It projects 234 -> 234 modules, 3,580 -> 3,581 graph edges, 429 -> 429 complete
packages, 1,043 -> 1,044 go.sum lines, and a 356 -> 358-line unapplied tidy
diff. The only effects are a redundant indirect requirement, main graph edge,
and selected full checksum; tidy removes the requirement and retains the sum.
No metadata was applied because that projection is not a version upgrade.

Template loads in zero project packages, repository source has no direct
import, and go mod why says the main module does not need it. Existing graph
requests arrive from prometheus/tsdb v0.7.1 and prometheus/common v0.4.1 for
the 2016 pseudo-version and prometheus/common v0.9.1 for the selected version.
MVS selects the latter. There is no actual project symbol to exercise.

Fresh primary vulnerability data was updated 2026-09-02T19:12:04Z and has
1,392 module records with no Template record. Normalized old/candidate IDs and
traces are identical at exact 20/30/20 Darwin-symbol/Darwin-module/Windows-
symbol populations; no trace contains Template.

The canonical no-change decision plus mandatory self-test stop makes repository
build/tests/race/vet, lint, help/API/CLI, preflight, snapshot/Docker, quality,
and audit gates for a changed selection inapplicable. No dependency
implementation commit was made.

## Template Evidence And Tools

- Decision evidence beneath the active scratch root has 1,702 verified
  entries. Evidence-manifest SHA-256 is
  2bfa08ca73085154ab5f0833814872efe481088703e6d7480edf1e2a40853068;
  decision-summary SHA-256 is
  4acc26e68fa5ad432a44d64ade0d399464f9b07eb2c6244b29711e048ae24170.
- Exact Go 1.26.7 was recreated beneath scratch. Binary SHA-256 is
  9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6;
  official archive SHA-256 is
  020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d.
  It was first in PATH with GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and no
  ambient GOFLAGS.
- Fresh govulncheck v1.7.0 reports exact Go 1.26.7. Rebuilt binary SHA-256 is
  6c0e51a045c0ccefaea5e780d11a2680f9a6133e9ebc19d95993971607a22320;
  this is a nonportable receipt, so version and functionality are durable.
- Preserve portable golangci-lint 2.12.2 archive receipt
  a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29,
  GoReleaser 2.17.1 receipt
  f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c,
  and apidiff receipt
  0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20.
- Detached commit verification used a short scratch-local keyring after the
  user GnuPG home and long agent-socket path proved unsuitable. This setup
  correction changed no repository file.

## Retained Decisions

- Retain exact selections github.com/alecthomas/repr v0.5.4,
  github.com/alecthomas/assert v1.0.0,
  github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b,
  github.com/alecthomas/chroma v0.10.0,
  github.com/alecthomas/colour v0.1.0,
  github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751,
  gopkg.in/alecthomas/kingpin.v2 v2.2.6, gopkg.in/resty.v1 v1.12.0,
  gopkg.in/errgo.v2 v2.1.0,
  gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15,
  gopkg.in/yaml.v2 v2.4.0, gopkg.in/yaml.v3 v3.0.1, and
  go.yaml.in/yaml/v3 v3.0.5. Retain selected Kong pseudo-version and every
  other prior bounded decision.
- Preserve all recorded corrections, including go-colorful mutable telemetry,
  btree regression, cast whitespace-path, source-archive normalization, Check
  history table, Errgo preflight runner, Resty committed projection/signal
  fixture, Kingpin external report, Units launcher timing, Assert external
  cache/TMPDIR, Repr cache/mktemp/Docker/manifest, Chroma source/consumer
  normalization, Colour pre-module/cache/mktemp, and Template GnuPG keyring
  corrections.
- The deliberately purged former two-entry recovery manifest remains
  historical SHA-256
  1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b.

## Next Objective

Independently evaluate selected exact-path github.com/antihax/optional v1.0.0
as the next single P7 group. It is selected by project MVS but is not an
explicit go.mod requirement. A fresh post-Template survey finds the stable
proxy list contains only v1.0.0. Proxy and exact Go @latest select that stable
version at 2019-10-10T23:37:20Z with Go 1.13. Exact @master selects the newer
unreleased pseudo-version v1.0.1-0.20220101210036-407d38fabb55 at
2022-01-01T21:00:36Z, also with Go 1.13. Do not treat the branch head as a
stable update. Treat qualification, source/signatures, release and alternate-
path history, closure, module tests, loaded packages, actual symbols, and
vulnerability effect as unknown until independently proved.

Current measurements remain 234 selected modules, 3,580 graph edges, 429
native complete-test packages, 1,043 go.sum lines, a 356-line unapplied tidy
projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations. The main module retains Go 1.18 and toolchain Go
1.26.7. Change only Optional's exact selected metadata and explained minimal
MVS closure if a qualified changed selection preserves the floor and passes
every applicable gate. Do not combine Template, another group, or P8.
