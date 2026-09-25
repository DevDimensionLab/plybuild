# Agent Session: Evaluate Beorn7 Perks Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-08T005841+0200-evaluate-beorn7-perks-dependency`
Created: `2026-09-08T00:58:41+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `85007d115e55d755724cc145e800f8dac7492b64df725fe83b2a16e551429e8a`
Previous: [2026-09-07T234555+0200-evaluate-armon-go-radix-dependency.md](2026-09-07T234555+0200-evaluate-armon-go-radix-dependency.md)
Next: [2026-09-08T023421+0200-evaluate-bgentry-speakeasy-dependency.md](2026-09-08T023421+0200-evaluate-bgentry-speakeasy-dependency.md)
Outcome: Retained Perks v1.0.1 without dependency edits; it is the canonical
  exact-path stable latest, while repository ancestry points to a divergent
  parent with a distinct module identity.

# Answer

Retain `github.com/beorn7/perks v1.0.1` without adding dependency metadata.
No dependency implementation commit was created.

## Canonical Identity And Release Qualification

The exact proxy lists only stable v1.0.0 and v1.0.1, dated
2019-04-14T22:11:40Z and 2019-07-31T12:00:54Z. Exact `@latest`, `@v1`,
`@master`, the exact-path default branch, and tag v1.0.1 all resolve selected
commit `37c8de3658fcb183f997c4e13e8337516ab753e6`. No proxy prerelease,
retraction, deprecation, `/v2` module, later default-branch commit, or renamed
authoritative exact path was found. V1.0.1 is therefore the canonical latest
and highest release-qualified exact-path version.

Go-import metadata maps `github.com/beorn7/perks` to
`https://github.com/beorn7/perks.git`. GitHub reports that repository enabled,
unarchived, undisabled, defaulting to `master`, and a fork whose parent and
source are `bmizerany/perks`. The ancestry fact does not replace the published
module identity: `github.com/bmizerany/perks` has its own go-import mapping and
is a distinct module path. Its current untagged master
`03f9df79da1edead2cdf5f8b4cf4d4f831d6e2d1` diverges from the exact-path
fork after 2014, has no `go.mod`, and resolves parent pseudo-version
`v0.0.0-20230307044200-03f9df79da1e`. Parent history does not produce an
in-place exact-path candidate.

The exact fork has five branches: `master`, `broken-but-fast`, `histo`,
`opt/pool-for-sample`, and `test/error-calculation`. The two queryable named
heads resolve old v0 pseudo-versions; the slash-named heads are 2014 Git refs,
one ancestral and one divergent. None is a later v1 candidate. The parent has
three branches and no tags. The fork has two tags and zero GitHub Releases;
the parent has zero tags and zero Releases. A stable signed Git tag is not a
GitHub Release object.

Both exact tags are signed annotated objects, while their target commits are
unsigned. Tag v1.0.0 object `4ded152d4a3e...` targets commit
`4b2b341e8d7715fae06375aa633dbb6e91b3fb46`, tree
`de22b44fded2eabaabd0ac413a3b04b44c8b3490`, parent
`4cf9acfa...`. Local verification succeeds with fingerprint
`5C69F212D616C4340FA8DD8504ABA6153ADA0C25`; GitHub reports `unknown_key`.
Tag v1.0.1 object `c49ff274687...` targets commit
`37c8de3658fcb183f997c4e13e8337516ab753e6`, tree
`77aeb432cabb3c9b94297378a2fe7697c29a5a88`, parent
`4b2b341e8d7715fae06375aa633dbb6e91b3fb46`. Local verification succeeds
with fingerprint `A100A34F34DEC17EE5EEF14C851C3DA17D748D03`; GitHub reports the
signature valid and verified. V1.0.1 changes only the `go.mod` Go declaration
from 1.12 to 1.11.

## Source, Closure, Tests, And Consumer

V1.0.0 declares exact module path `github.com/beorn7/perks`, Go 1.12, and no
requirements. V1.0.1 declares the same path, Go 1.11, and no requirements.
Each complete minimal closure is one module and three packages, so neither can
raise the retained Go 1.18 floor.

V1.0.0's source/module checksum pair is
`h1:HWo1m869IqiPhD389kmkxeTalrjNbbJTC8LXupb+sl0=` /
`h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=`. V1.0.1's pair is
`h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=` /
`h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=`; sumdb agrees. Proxy
ZIP SHA-256 values are
`a7ec6164e31ea8e10c601abb9793753ec43cb218283b226800c134fb23cea409`
and `25bd9e2d94aca770e6dbc1f53725f84f6af4432f631d35dd2c46f96ef0512f1a`.
Every one of each archive's 15 regular files matches the corresponding exact
tag. Normalized source-manifest hashes are
`5e623be139334393e642a7e5e1096bb3616f4e21a787959d563a9bb4f55b1153`
and `80fa7a8e5aa4189f70e1896740da50eab8d4ad61c5464adbc3a4574d0a141c2f`.

Final writable proxy and exact-tag matrices for both releases pass
verification/listing, count-1, two count-10 runs, race, vet, and source
immutability under exact Go 1.26.7 and Go 1.18.10. V1.0.1 also passes focused
`topk` count-100 in both forms and toolchains. An earlier v1.0.0 proxy
Go-1.26.7 count-10 run failed once in `topk.TestTopK` at index 9 (`want "9",
got "b"`). Four final matrices and count-100 stress did not reproduce it.
Because TopK sorts only by count, equal-count/map-order nondeterminism is a
plausible cause; the observation remains recorded and is not attributed to
v1.0.1's go.mod-only delta.

Selected `github.com/prometheus/client_golang v1.4.0` is the actual graph
consumer. Its Summary imports only `perks/quantile` and uses `Stream`,
`NewTargeted`, `Insert`, `Count`, `Query`, and `Reset`. Focused Summary tests
pass count-1, count-10, and race under exact Go 1.18.10. Its separate
standalone test closure contains 41 modules and 167 packages, with no declared
Go version above 1.12. Default consumer test/vet is separately blocked by six
historical int-to-string diagnostics in gauge, histogram, and Summary test
formatting; execution with implicit test vet isolated passes. Mandatory Perks
and project vet both pass.

## Projection, Quality, And Vulnerabilities

Accepted measurements remain 234 selected modules, 3,581 graph edges, 429
complete-test packages, 41 loaded modules, zero loaded Perks packages, 1,045
`go.sum` lines, and a 361-line tidy projection. Exact selected-version get
changes no selection. It projects only a redundant indirect requirement, one
main-to-Perks edge, and the full v1.0.1 checksum: 234/3,582/429/41/0, 1,046
sum lines, and 370 tidy-diff lines. Tidy removes the two projected metadata
lines and restores the inherited selection; the remaining 361 lines are the
accepted historical debt. No projection was applied.

The inherited requests are TSDB v0.7.1 -> old pseudo-version, client_golang
v1.4.0 -> v1.0.1, client_golang v1.0.0 -> v1.0.0, and Prometheus common
v0.4.1 -> old pseudo-version. Perks has no outgoing graph edges. It loads in
no Ply package, Ply has no source import, and `go mod why -m` says the main
module does not need Perks.

Baseline and projected trees pass repository verification/build, complete
count-1/count-10/race tests, vet, Windows-amd64 build, pinned golangci-lint
2.12.2, empty-HOME count-2, and byte-identical root/status/upgrade/build help.
API and CLI reports are byte-identical at SHA-256
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
Corrected preflight passes all 62 launcher controls, distribution/lint/install/
toolchain and acceptance-script meta-contracts, 80/80 mutation-harness
controls, and all 15 audit controls.

Fresh govulncheck v1.7.0 used primary data updated
2026-09-02T19:12:04Z. The 1,392-record module index contains neither exact
nor parent Perks. Baseline/projected findings are byte-identical after
normalization: Darwin and Windows each have 20 reachable IDs in 22 traces;
the Darwin module scan has 30 IDs. Reachable IDs are `GO-2023-1989`,
`GO-2023-1990`, `GO-2024-2937`, `GO-2024-3333`, `GO-2025-3595`,
`GO-2026-4440`, `GO-2026-4441`, `GO-2026-4815`, `GO-2026-4961`,
`GO-2026-5025`, `GO-2026-5027`, `GO-2026-5028`, `GO-2026-5029`,
`GO-2026-5030`, `GO-2026-5031`, `GO-2026-5032`, `GO-2026-5061`,
`GO-2026-5062`, `GO-2026-5066`, and `GO-2026-6222`.

The 22 Darwin/Windows trace endpoints are exact across seven groups: seven
`x/image/tiff.Decode -> cmd.ExecuteE` traces (15 frames), nine
`x/net/html.Parse -> cmd.ExecuteE` traces (13), two
`x/image/webp.Decode -> cmd.ExecuteE` traces (15), and one each for
`x/image/bmp.Decode -> cmd.ExecuteE` (15),
`x/image/vp8l.Decode -> cmd.ExecuteE` (17),
`x/image/tiff.(*buffer).ReadAt -> pkg/http.GetJsonWithAccessToken` (5), and
`x/image/webp.init -> cmd.init` (4). Normalized Darwin-symbol, Windows-symbol,
and Darwin-module finding hashes are respectively
`fcf9d449a455a00562126bca8b863b882ede61d7c087a967ac9694c30978c859`,
`b7d80a548ee7476342d7a0bd2b0efd87524a5efa0c0da960e802c4e186fb3e8f`,
and `390f7bf685e19fc5d10602da7e35394c97f8fd460fe40f92c4e6726c00627d29`.

Canonical qualification produced no changed selection. Consequently the
changed-selection-only host/snapshot/Docker acceptance, exact `make quality`,
focused/Q0-Q2/full audits, and dependency commit were inapplicable and are not
claimed. Exact Go 1.26.7 retained required binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
The sealed evidence manifest contains 477 verified entries and has SHA-256
`aa77d2ab9cf676ecfb7ef544c1c5db76ecda86f580ff8ff8f03b8c78038677f0`;
decision-summary SHA-256 is
`6db4ca7260bde6b4affa48c62adb381bda20687f39f013c4b1cdc25397418754`.

Scratch-only corrections remain explicit: initial environment and source-
prefix fixes, a removed empty slash-branch query, vulnerability-order
normalization, exact historical compatibility-cache warming, golangci cache
relocation, BSD-mktemp adaptation, cleared recursive make overrides, and a
200 ms signal-fixture scheduling shim after two raw-log readiness races. None
changed a repository file or concealed a dependency failure.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/beorn7/perks v1.0.1` as one bounded dependency group. Resolve
canonical source identity across the exact-path fork and its parent, canonical
latest, release qualification, default-branch history, and the highest
qualified Go-1.18-floor-compatible candidate from primary evidence. Make an
exact dependency selection only if it changes a selected version, preserves
the retained floor through the complete minimal closure, and passes every
quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Go Radix v1.0.0. All earlier rejections,
no-change decisions, accepted closures, and evidence corrections remain final.
Do not revisit Go Radix or combine another module group. P8 remains queued.

Project MVS selects Perks v1.0.1 at 2019-07-31T12:00:54Z although it is not an
explicit `go.mod` requirement. A minimal post-Radix survey finds exact stable
proxy versions v1.0.0 and v1.0.1. Exact Go `@latest`, `@v1`, and `@master` all
resolve selected v1.0.1; its proxy module declares Go 1.11. Do not infer the
complete closure floor from this one directive or assume a selected latest
version needs materialized metadata.

The public exact-path GitHub repository currently reports enabled, unarchived,
fork status, default branch `master`, five branches, two tags, and zero GitHub
Releases. It identifies `bmizerany/perks` as both parent and source. Exact-path
master and tag v1.0.1 point to
`37c8de3658fcb183f997c4e13e8337516ab753e6`; v1.0.0 points to
`4b2b341e8d7715fae06375aa633dbb6e91b3fb46`. Treat canonical source identity,
the fork/parent relationship, tag and commit signatures, release/tag history,
retractions, all version and branch declarations, complete closures, tests,
consumers, loaded behavior, and vulnerability effect as unknown until
independently proved. Explicitly distinguish a stable tag from a GitHub
Release and the exact module identity from repository ancestry.

# Measurements At Start

Latest dependency implementation remains Circbuf commit
`3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
`7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
`a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
`go.sum`. The Consul API, Go Metrics, and Go Radix evaluations made no
dependency commit. Operator-authorized lifecycle repair
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between the
earlier Repr implementation and direct-child Repr handoff. Preserve it.

Go Radix remains selected at v1.0.0, the sole stable exact-path proxy version
and canonical `@latest`/`@v1`. Master pseudo-version
`v1.0.1-0.20221118154546-54df44f2176c` passed module and consumer execution
but is an unreleased branch head with no tag or GitHub Release, so it was not
selected merely because it is newer. Exact selected-version get changed no
selection and only projected redundant metadata, which was not applied. Go
Radix evidence has 390 verified entries; manifest SHA-256 is
`6be6b857e77c7097b402ea5f4fe0ce849e15382b7167a004448cd85d82926eaf`;
decision-summary SHA-256 is
`493e832f2d0e6456bb64462104e1bd5b80f4b6eb32d401b9b098acdd0dd95e6a`.

Accepted measurements remain 234 selected modules, 3,581 graph edges, 429
native complete-test packages, 41 loaded modules, 1,045 `go.sum` lines, and a
361-line unapplied tidy projection. Relative to accepted go-cmp commit
`c314bcb`, accepted metadata adds exactly 29 checksum lines. The main module
retains Go 1.18 and toolchain Go 1.26.7. Ordinary and ignored status must be
empty.

Fresh primary vulnerability evidence contains 1,392 module records and no Go
Radix record or trace. Accepted projections retain exact 20-ID/22-trace
Darwin symbol, 30-ID Darwin module, and 20-ID/22-trace Windows symbol
populations.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate pinned tools beneath scratch as needed. Portable
receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
Treat rebuilt golangci-lint and govulncheck hashes as nonportable receipts;
prove versions and functionality.

# Role And Boundaries

From fresh external archives and caches, resolve every relevant exact-path
Perks v1 version through the Go proxy and checksum database, exact-path
repository, its declared parent/source repositories, go-import metadata, and
primary Go vulnerability data. Record selected and candidate commits/times,
module Go declarations and requirements, checksum pairs, source identity, tag
and commit signatures, release history, archived/deprecated status, and
retractions. Explicitly distinguish stable versions, prereleases, pseudo-
versions, redirects, forks, alternate module paths, branch heads, and
unreleased commits.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete minimal closure, not from a
single module directive. Decide whether the exact-path fork is the canonical
published source and whether parent history changes qualification. Measure old
versus candidate modules, graph edges, complete packages, checksums, loaded
packages and paths, exact-get diff, and `go mod tidy -diff`. Explain every
selection, edge, and checksum change. Identify real consumers and exercise the
actually used symbols.

Require candidate module complete, repeated, and race-enabled self-tests plus
vet; repository verify/build, complete tests, race, vet, Windows build, pinned
lint, byte-identical public help, API/CLI reports, and primary vulnerability
identity. Preserve the retained Go 1.18 floor across the complete changed
closure. Keep standalone test-only closure separate from project MVS.

Reject or retain without dependency edits if canonical qualification, source
identity, floor compatibility, closure, module tests, loaded behavior, or any
applicable quality contract fails. Do not manufacture metadata when exact get
changes no selected version. Record precise primary old/candidate
vulnerability IDs and traces.

# Required Reading

Verify branch, clean ordinary and ignored status, exact ancestry, reciprocal
archive history, `./codex-dev-start.sh --check`, and P7/P8 state before work.
Read this archive, the rolling handover, P7 roadmap, `go.mod`, `go.sum`, the
answered Go Radix, Go Metrics, Consul API, Circbuf, Optional, Template, Colour,
Chroma, Kong, Repr, Assert, and Units archives, retained
Kingpin/Resty/Errgo/Check/YAML outcomes, and the toolchain, quality, baseline,
compatibility, snapshot/Docker, acceptance, audit, and lifecycle contracts. Do
not reopen earlier decisions.

# Three Moves

If and only if a qualified candidate changes the selected version and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/beorn7/perks@<candidate>` for one dependency-only commit.
Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, toolchain declaration, production source, quality apparatus,
and release input.

If canonical qualification confirms v1.0.1 as the highest qualified selection,
or exact get changes no selection, retain it without adding a redundant direct
or indirect requirement, main edge, or checksum solely for metadata. A no-
change decision gets no dependency implementation commit.

After a changed selection, run the complete P7 dependency gate: module and
consumer tests, graph/path/checksum/tidy explanation, repository verify/build/
tests/race/vet/Windows/pinned lint, help/API/CLI, launcher and Make contracts,
preflight, host plus fresh snapshot/Docker meta and acceptance, audit meta,
focused and exact Q0-Q2 audits, separate full audit, vulnerability comparison,
empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 Q0-Q2 rows at L2 and zero held,
regressed, not-comparable, or dirty counts. Full audit may exit 1 only for
established queued L3 rows, never 2.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
create direct `/private/tmp/ply-*` roots, never bypass scratch cleanup, and
never run `go mod download all` inside a measured tree. Warm only exact needed
closures from a separate scratch archive. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Perks decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next single P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
the next group, launch a successor, push, merge, publish, release, stash,
revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
