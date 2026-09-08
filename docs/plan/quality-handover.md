# Quality Upgrade Handover

Generated: 2026-09-08T06:01:23+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- Crypt was retained without a dependency implementation commit. The latest
  dependency implementation remains Speakeasy commit
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a`, exact parent
  `1f55aaa31280a4610ed66c1c37ddda953bfa6a8e`, tree
  `a50d4f256fe742e472f1f7e9cf0589a596e971b0`, changing only `go.mod` and
  `go.sum` with three insertions. Its handoff is commit
  `06edec90bed97f56b178df91364217c00e2778ac`.
- Earlier Circbuf implementation
  `3be2183ee310ccdc358ce4ed372c0785de25b88b` and operator-authorized
  lifecycle repair `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remain ancestors.
- The answered Crypt archive and sole NEXT OpenCensus Proto archive must link
  reciprocally. No `.agent-task/current.md` or
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, the accepted
Speakeasy v0.2.0 move, and the retained Crypt pseudo-version. P8 remains
queued. All earlier acceptances, rejections, no-change decisions, evidence
corrections, and lifecycle ancestry are final; do not reopen them or combine
another group.

Every disposable cache, projection, archive, source, report, schema-2
document, generated artifact, evidence tree, and build context must remain
beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never create direct
`/private/tmp/ply-*` roots, never run `go mod download all` in a measured
tree, and preserve the launcher's scratch cleanup and exact reciprocal archive
contract.

## Retained Bketelsen Crypt Group

Retain exact-path `github.com/bketelsen/crypt` at selected pseudo-version
`v0.0.3-0.20200106085610-5cbc8cc4026c` without changing `go.mod` or
`go.sum`. Stable v0.0.5 is canonical latest, default-branch head, and the
highest stable version whose complete minimal closure preserves Go 1.18, but
it is not quality-qualified: exact Go 1.26.7 and Go 1.18.10 report the same
seven mandatory `go vet` failures. Stable v0.0.3 and v0.0.4 fail on those
same unkeyed `backend.Response` literals. No higher stable Crypt v0 version
passes every quality contract.

The exact proxy lists only v0.0.1 through v0.0.5. Exact `@latest`, `@v0`, and
`@master` all resolve v0.0.5 at 2021-10-08T10:39:19Z. No version is retracted
or deprecated. Go-import metadata maps the exact module path to the public,
enabled, unarchived `https://github.com/bketelsen/crypt.git`. That repository
is an exact-path fork of `xordataexchange/crypt`; the parent and later
`github.com/sagikazarmark/crypt` are distinct source/module identities.

The canonical repository's default `master` is v0.0.5 commit
`60c5f2086f0eae50f5275599096bcc6090d12cf8`, tree
`19f405d906832661efa1d2129d163f8efff16560`, with no later default-branch
commit. Its other 13 heads are feature or Dependabot branches, including
unreleased 2023 work. All five tags are lightweight. GitHub Releases exist
only for v0.0.3, v0.0.4, and v0.0.5 and are non-draft/non-prerelease.

The selected prerelease-form pseudo-version is untagged merge commit
`5cbc8cc4026c0c1d3bf9c5d4e5a30398f99c99a9`, tree
`55cf4568fd528c73254e6b1f03089d8a61bd39a5`, at
2020-01-06T08:56:10Z. Selected plus v0.0.3, v0.0.4, and v0.0.5 commits carry
GitHub web-flow signatures independently verified with fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`; the key is now expired.
V0.0.1 and v0.0.2 commits are unsigned.

Every proxy ZIP matches exact Git source. Normalized manifest SHA-256 values
for v0.0.1, v0.0.2, selected, v0.0.3, v0.0.4, and v0.0.5 are respectively
`a9f8748b51390bfdedb38eee4c23c9e53c9f6809168ab64fb6c9a3f46c5b98d7`,
`e0c379854d08d226f1bc592ad859d5cabe1fd2fedfc5012ff7b4bcdc6b422b6b`,
`911a080413d659aedef61922bbbf77c6a0ddf7d5e4a387fde44578bec82c7959`,
`d4318e058c2d911504b84f2e5512f05ccf241040fb3999195b29dbb694eb3938`,
`8a616cc41db6134fd0bba0471a9fd7ad94b87f8464761c12ea8aa8dd3e20e8ed`,
and `de3b02f2771c24908e9f55d09b6e8ff61752586beb73fec61c04c7b272585900`.
Sumdb source/module pairs and archive hashes are in the answered archive.

V0.0.3 onward, including selected, declare exact path and Go 1.12. V0.0.5's
complete minimal closure contains 151 modules, 2,384 exact-Go graph edges,
eight Crypt packages, and 426 package dependencies including tests. Its
highest declared Go version is 1.17. The proxy and exact-Git candidate both
pass verify, package listing, count-1, two count-10 runs, and race under exact
Go 1.26.7; contained Go 1.18.10 count-1/count-10/race also pass. Vet alone is
the decisive release-qualification failure.

`github.com/devdimensionlab/mvn-pom-mutator v0.2.3` requests the selected
pseudo-version but imports no Crypt package. Crypt is unloaded by Ply and
`go mod why -m` reports that the main module does not need it. Historical
Viper remote support is the real consumer: v1.7.x requests selected, v1.8.x
requests v0.0.4, then v1.9.0 changes module identity. Viper v1.7.1's remote
package compiles and vets against selected and v0.0.5; Crypt's mock-backed
tests exercise standard/encrypted Set/Get/List/Watch manager behavior.

## Projection, Quality, And Vulnerability Measurements

The accepted project remains 234 modules, 3,582 graph edges, 429 native
complete-test packages, 41 loaded modules, 197 loaded packages, 1,047
`go.sum` lines, and a 371-line unapplied tidy projection. Relative to accepted
go-cmp commit `c314bcb`, accepted metadata still adds exactly 31 checksum
lines.

A v0.0.5 project projection changes only Crypt and has 234 modules, 3,700
edges, 429 packages, 41 loaded modules, 197 loaded packages, 1,065 checksum
lines, and a 461-line tidy projection. Its 18 checksum additions are the
candidate pair plus 16 transitive module-file hashes. Tidy removes the
explicit pin and candidate checksums and restores inherited selected Crypt.
The projection passes project verify/build/count-1/count-10/race/vet,
Windows-amd64 build, pinned lint, byte-identical help, API/CLI compatibility,
and empty-HOME count-2, but none can override the dependency's own vet
failure.

An exact selected-version projection changes no selected module. It merely
adds a redundant indirect root requirement, one main graph edge, and the
selected source checksum: 234 modules, 3,583 edges, 1,048 checksum lines, and
a 373-line tidy projection. It was not applied.

Fresh govulncheck v1.7.0 uses primary data updated
2026-09-02T19:12:04Z. The 1,392-record module index contains no Crypt record.
Old and candidate normalized results are identical: 20 IDs/22 reachable
traces on Darwin and Windows and 30 Darwin module IDs. Crypt appears in no
finding or trace. Both platforms' normalized reachable SHA-256 is
`7757df547ed97c0b709bfe8cbbbf807356331727c870bdeb7f0602b1f393ba79`.

Crypt evidence contains 10,694 verified entries; manifest SHA-256 is
`e0320cf4ef063c83cee9a5131fdad9ceae8b90730502617ab48d2f16c759c93b`.
Decision-summary SHA-256 is
`6356ad761738c8d9e664693a0c550bbe4305db500bf9198e567ef7a3e1d33a75`.

## Tools And Corrections

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Exact Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Portable receipts remain golangci-lint 2.12.2 archive `a9c54498...`,
  GoReleaser 2.17.1 binary `f5f08a77...`, and apidiff `0c55d9e3...`.
- Scratch-only corrections were newline-safe GitHub signature extraction,
  a writable historical-consumer replay, separate historical API cache
  warming, evidence-directory-relative manifest verification, and
  vulnerability return-code normalization. Initial superseded attempts are
  retained in evidence; none changed repository source or concealed a
  candidate failure.

## Retained Decisions

Retain the accepted Speakeasy v0.2.0 move and all earlier exact selections and
outcomes, including Repr v0.5.4, Assert v1.0.0, Units pseudo-version
`0f3dac36c52b`, Chroma v0.10.0, Colour v0.1.0, Template pseudo-version
`fb15b899a751`, Optional v1.0.0, Circbuf pseudo-version `5111143e8da2`,
Consul API `eb2c6b5be1b6`, Go Metrics v0.4.0, Go Radix v1.0.0, Perks v1.0.1,
Kingpin v2.2.6, Resty v1.12.0, Errgo v2.1.0, Check pseudo-version
`41f04d3bba15`, YAML v2.4.0/v3.0.1, and `go.yaml.in/yaml/v3 v3.0.5`.
Preserve the accepted Kong pseudo-version and all other decisions recorded in
the roadmap and answered archives.

## Next Objective

Independently evaluate selected exact-path
`github.com/census-instrumentation/opencensus-proto v0.3.0` as the next single
P7 group. The proxy lists stable v0.0.1 through v0.4.1; exact `@latest` and
`@v0` resolve v0.4.1 at 2022-09-23T17:40:20Z, while exact `@master` resolves
unreleased pseudo-version `v0.2.2-0.20230502190750-1664cc961550` because
default-branch history split after the v0.4 release line. The canonical
repository is public and archived.

Selected v0.3.0 declares only its module path. Stable v0.4.1 declares Go 1.18
and materially newer grpc-gateway, gRPC, protobuf, x/net, x/sys, x/text, and
genproto requirements. Independently prove the split release/default-branch
history, archive and source identity, signatures, complete closure floor,
tests, historical Viper and Sagikazarmark consumers, project projection, and
vulnerability effect. Do not implement or combine another group until this
bounded decision is complete.
