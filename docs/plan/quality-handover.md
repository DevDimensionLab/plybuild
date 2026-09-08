# Quality Upgrade Handover

Generated: 2026-09-08T07:44:22+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- OpenCensus Proto and Crypt were retained without dependency implementation
  commits. The latest dependency implementation remains Speakeasy commit
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a`, exact parent
  `1f55aaa31280a4610ed66c1c37ddda953bfa6a8e`, tree
  `a50d4f256fe742e472f1f7e9cf0589a596e971b0`, changing only `go.mod` and
  `go.sum` with three insertions. Its handoff is commit
  `06edec90bed97f56b178df91364217c00e2778ac`.
- Earlier Circbuf implementation
  `3be2183ee310ccdc358ce4ed372c0785de25b88b` and operator-authorized
  lifecycle repair `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remain
  ancestors.
- The answered OpenCensus Proto archive and sole NEXT XXHash v2 archive must
  link reciprocally. No `.agent-task/current.md` or
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, the accepted
Speakeasy v0.2.0 move, and retained Crypt and OpenCensus Proto selections. P8
remains queued. All earlier acceptances, rejections, no-change decisions,
evidence corrections, and lifecycle ancestry are final; do not reopen them or
combine another group.

Every disposable cache, projection, archive, source, report, schema-2
document, generated artifact, evidence tree, and build context must remain
beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never create direct
`/private/tmp/ply-*` roots, never run `go mod download all` in a measured tree,
and preserve the launcher's scratch cleanup and exact reciprocal archive
contract.

## Retained OpenCensus Proto Group

Retain exact-path `github.com/census-instrumentation/opencensus-proto` at
selected v0.3.0 without changing `go.mod` or `go.sum`. Stable v0.4.1 is
canonical latest but its complete closure exceeds Go 1.18. Stable v0.4.0
preserves the floor but has a broken committed module and incompatible
generated API. No higher stable version passes every contract.

The exact proxy lists v0.0.1, v0.0.2, v0.1.0, v0.2.0, v0.2.1, v0.3.0,
v0.4.0, and v0.4.1. Exact Go `@latest` and `@v0` resolve v0.4.1 at
2022-09-23T17:40:20Z; exact `@master` resolves unreleased
`v0.2.2-0.20230502190750-1664cc961550` at 2023-05-02T19:07:50Z. No version is
retracted and no module deprecation directive exists.

Go-import identifies public
`https://github.com/census-instrumentation/opencensus-proto.git`. The enabled,
archived, non-fork repository has default branch `master`. OpenTelemetry is a
successor project, not an alternate exact module path.

Primary evidence corrects the incoming selected SHA. Tag v0.3.0 is commit
`4aa53e15cbf1a47bc9087e6cfdca214c1eea4e89`, tree
`ea0ad81b63231a53d01a5c0c90afec09692d87fb`, parent
`cfe95db9f268b874bdbdff0b1d018b11d6640abf`, at
2020-07-21T05:46:08Z. V0.4.1 is commit
`e53624a87b9b9b919147a9b4626c669a869ebb34`, tree
`6b25e84886ea951380b99771e4ee55f2a4c45d39`, parent
`576a4cae65940a353684090ac4a2ec456b88880d`. Master is commit
`1664cc961550be8f3058ddd29390350242f44f1f`, tree
`367e035582d24b89688e898f59548780d36362ce`, parent
`3619b5dda8bff26ff1974714c24de8f6d4953811`.

Master and v0.4.1 merge at
`a5f3b19cae5836c060a7646bc6db436bf5d784fc`, the parent of v0.4.0. The v0.4
release line and later master v0.2.1 line are descendants on different
branches. Canonical stable precedence and later default-branch commit time are
therefore not interchangeable.

All eight tags are annotated and unsigned. V0.4.0/v0.4.1 commits have valid
GitHub PGP verification; v0.4.1's packet issuer fingerprint is
`0700C92CD0CFBD86D9CB3D26563A85007BFA1BA2`, but the retired public key is
not currently available for independent replay. Master has a valid GitHub
web-flow signature; v0.3.0 is unsigned. GitHub Releases exist for all stable
tags except v0.2.1 and are non-draft/non-prerelease.

Every proxy source matches exact Git source. Normalized manifest SHA-256
values for v0.0.1, v0.0.2, v0.1.0, v0.2.0, v0.2.1, v0.3.0, v0.4.0, and
v0.4.1 are respectively
`15162dc2a357c6ca2357dea3638a0f7b03c4b4925c7e3143abc601c11019e7c1`,
`58378e75065483b0c81f90d85e655efa25220fffa0e6c64535dc263045434082`,
`161da69d3a204102656e9949914391880937ed3cb6d6aba4125f4307a22fb698`,
`adc3ef7eea686d49f2407d0c5c650233673579e06d642142c092ecd83e696097`,
`4c5ccb2758b7af1b9b087860b807c35c9e82e260291da56154b1cd251b6b1ed0`,
`fb39a13734d16781ba5dabe26c1bf5ff035ee6ebabfc78090d85bf070318b3bf`,
`d3f301255eb92f55cc9191d385d37e5f283a070c203765714b410f7c17fcfcc2`,
and `1f7b8f040437b1b3c84e1dd8b843bd4d03c93fce44cd098778f1921ab2ce87ef`.
The answered archive records every sumdb source/module checksum pair.

V0.4.1's standalone closure has 28 modules and 64 graph edges. Root Go 1.18
is not sufficient: selected genproto pseudo-version `9e6da59bd2fc` declares
Go 1.19, so the complete closure violates the retained floor. Its proxy and
exact-Git source otherwise pass verify/list/count-1/two count-10/race/vet
under Go 1.26.7 and contained Go 1.18.10. It has seven packages and zero
native test files.

V0.4.0's standalone closure has 36 modules, 208 graph edges, and no declaration
above Go 1.18. It fails package/test/race/vet in proxy and exact-Git form under
both SDKs because its generated gateway files import grpc-gateway/v2 while
`go.mod` requires grpc-gateway v1.16.0. The v0.4.1 repair explicitly says tidy
was forgotten after proto updates. Pinned apidiff also reports eight
incompatible MetricsService/TraceService gateway handler signatures from the
v1-to-v2 `ServeMux` type change.

Selected v0.3.0's seven zero-test packages pass verify/list/count-1/two
count-10/race/vet in the accepted project-selected closure under both SDKs. A
scratch compatibility harness passes populated `TraceConfig` marshal/
unmarshal and an in-memory `TraceService.Export` round trip against v0.3.0 and
v0.4.1 under both SDKs.

Historical Viper v1.10.1 and Sagikazarmark Crypt v0.4.0 reach only
`gen-go/trace/v1` through Firestore, Google API gRPC transport, gRPC xDS, and
Envoy trace configuration. Focused tests and race are invariant and pass
against v0.3.0/v0.4.0/v0.4.1. Viper vet passes. Crypt's invariant unkeyed
`backend.Response` vet finding at `backend/firestore/firestore.go:110` is
historical consumer debt, not a candidate regression.

## Projection, Quality, And Vulnerability Measurements

The accepted project remains 234 modules, 3,582 graph edges, 429 native
complete-test packages, 41 loaded modules, 197 loaded packages, 1,047
`go.sum` lines, and a 371-line unapplied tidy projection. Relative to accepted
go-cmp commit `c314bcb`, accepted metadata still adds exactly 31 checksum
lines.

Exact selected-version get changes no selection and projects only a redundant
indirect root requirement, one main edge, and the selected source checksum:
234 modules, 3,583 edges, 429/41/197, 1,048 checksum lines, and a 375-line
tidy diff. It was not applied.

Exact v0.4.0 get changes only OpenCensus Proto: 234 modules, 3,591 edges,
429/41/197, 1,049 checksum lines, and a 377-line tidy diff. Eight root
requirements plus the main edge explain the edge delta; its checksum pair is
the only checksum addition. Exact v0.4.1 get also selects grpc-gateway/v2
v2.11.3: 235 modules, 3,591 edges, 429/41/197, 1,049 checksum lines, and a
377-line tidy diff. The gateway checksum was already retained.

OpenCensus Proto loads in zero Ply packages and `go mod why -m` says the main
module does not need it. The v0.4.0 project projection nevertheless passes
mod verify/build/count-1/count-10/race/vet, Windows build, pinned lint,
byte-identical help, identical API/CLI reports, compatibility, CLI surface,
and empty-HOME count-2. Changed-selection acceptance/quality/audit gates are
inapplicable after the dependency stop rule.

Fresh primary data updated 2026-09-02T19:12:04Z has 1,392 module records and
no OpenCensus Proto record. Old/v0.4.0/v0.4.1 normalized results are identical:
20 IDs/22 reachable traces for Darwin and Windows symbol scans, 22 IDs/22
Darwin package findings, and 30 Darwin module IDs. No finding names the target.

OpenCensus Proto evidence contains 459,623 verified entries; manifest SHA-256
is `55c7090ab788c633f20666ba9d70e8f4d5f4bcb446b9e5ef8f530cd78e3b074a`.
Decision-summary SHA-256 is
`f7f4ed0dd015d7d7da58141a48f63537b8552c34761c05bc90de5b3efa22564a`.

## Tools And Corrections

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Exact Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Portable receipts remain golangci-lint 2.12.2 archive `a9c54498...`,
  GoReleaser 2.17.1 binary `f5f08a77...`, and apidiff `0c55d9e3...`.
- Scratch-only corrections were exact source normalization, exported
  environment reruns after a shell-prefix error, contained SDK/race/vet paths,
  a vet-corrected compatibility assertion, historical API/CLI cache warming,
  and a concurrent-lint-lock replay. The initial scratch-only API/CLI warm used
  prohibited `go mod download all`; it touched no measured worktree and none of
  its results is relied upon. Fresh clones/cache warmed only by exact get,
  graph, and package-list operations passed the final checks fully offline.
  Superseded evidence remains sealed; none changed repository source or
  concealed a candidate failure.

## Retained Decisions

Retain the accepted Speakeasy v0.2.0 move and all earlier exact selections and
outcomes, including Repr v0.5.4, Assert v1.0.0, Units pseudo-version
`0f3dac36c52b`, Chroma v0.10.0, Colour v0.1.0, Template pseudo-version
`fb15b899a751`, Optional v1.0.0, Circbuf pseudo-version `5111143e8da2`,
Consul API `eb2c6b5be1b6`, Go Metrics v0.4.0, Go Radix v1.0.0, Perks v1.0.1,
Kingpin v2.2.6, Resty v1.12.0, Errgo v2.1.0, Check pseudo-version
`41f04d3bba15`, YAML v2.4.0/v3.0.1, and `go.yaml.in/yaml/v3 v3.0.5`.
Preserve the accepted Kong pseudo-version and all other roadmap/archive
decisions.

## Next Objective

Independently evaluate selected exact-path
`github.com/cespare/xxhash/v2 v2.1.2` as the next single P7 group. It is
selected through Viper v1.10.1 and Sagikazarmark Crypt v0.4.0, but no XXHash
package is loaded by Ply. The exact proxy lists v2.0.0, v2.1.0, v2.1.1,
v2.1.2, v2.2.0, and v2.3.0; canonical latest v2.3.0 is dated
2024-04-04T20:00:10Z. Both selected and latest declare exact path and Go 1.11.

The public `github.com/cespare/xxhash` repository is enabled, unarchived,
non-fork, and defaults to `main`. Selected/latest annotated tag refs dereference
commits `e7a6b52374f7e2abfb8abb27249d53a1997b09a7` and
`998dce232f17418a7a5721ecf87ca714025a3243`. Independently establish release
identity/signatures, complete closure floor, architecture-specific and pure-Go
behavior, historical consumers, exact project projection, and vulnerability
effect before considering one exact dependency move. Do not implement or
combine another group until that bounded decision is complete.
