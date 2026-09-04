# Quality Upgrade Handover

Generated: 2026-09-05T01:38:47+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation remains dependency-only go-cmp commit
  `c314bcb440b5f94871d71a249bca7ae7f87d5543`, exact parent
  `cbdb0a915d7d30bb6633cfbb2ccf0323aea1b77a`, clean tree
  `9a212b377cb08e7d6e3044fb6dae071058b6242b`.
- The HCL evaluation made no dependency or production change. Its final
  documentation handoff must have exact parent
  `99f508d927047f3803e96d082efc15e2dd2422b7`; `go.mod` and `go.sum` remain
  byte-identical to `c314bcb`. Ordinary and ignored status must be empty.
- The answered HCL archive and the inconshreveable/mousetrap NEXT archive must
  link reciprocally. Only the launcher's mutable header and prompt regions may
  change during handoff.
- No `.agent-task/current.md` or `.quality/manual-evidence.json` was created.
  No push, merge, publication, release, stash, revert, retained evidence/image
  deletion, successor launch, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain,
completed dependency groups through accepted `github.com/google/go-cmp
v0.6.0`, and the rejected/no-change Hashicorp HCL evaluation. Latest go-cmp
v0.7.0, Viper v1.16.0, Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0,
fsnotify v1.10.1, latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d`, and an HCL selection change remain
rejected for their recorded floor, loaded-behavior, release-qualification, or
self-test failures. P8 remains queued.

## Rejected/No-Change Hashicorp HCL Group

The canonical Go proxy `@latest`, `go list -m -versions`, and
`go list -m -u` agree that stable `github.com/hashicorp/hcl v1.0.0` remains
the selected module release; the update query contains no `Update` field.
Its proxy publication and commit time are 2018-08-26T00:51:36Z. The annotated,
unsigned tag object `5592d2526badd60c172ffa95c6a3b209bea9d1ee`, tagged at
2018-08-26T01:00:46Z, points to unsigned commit
`8cb6e5b959231cc1119e43259c4a608f9c51a241`. All 196 proxy files are
byte-identical to that commit. The module has no `go` directive and therefore
declares no floor above retained Go 1.18. Its sumdb pair is
`h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=` /
`h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=`.

Every higher v1 version string is an application-targeted prerelease. They are
floor-compatible but are not eligible as the canonical stable candidate:

| Version | Commit / UTC time | Go | Tag and commit signature | Proxy files and sumdb pair |
| --- | --- | --- | --- | --- |
| `v1.0.1-vault` | `809e678c39ec71ae0b37a792de60b7e10e571dfe`, 2020-10-15T18:49:41Z | none | original tag now absent; commit unsigned | 196 identical; `h1:UiJeEzCWAYdVaJr8Xo4lBkTozlW1+1yxVUnpbS1xVEk=` / `h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=` |
| `v1.0.1-vault-2` | `4ecc8eea909656efc5215b9d0c5a4f82232065be`, 2021-05-04T18:33:52Z | 1.15 | lightweight; valid GitHub commit signature | 199 identical; `h1:j0lTHGBdaU13Pc3GaTCdWjmsT22X98bsHnA+ShzIOtg=` / `h1:XYhtn6ijBSAj6n4YqAaf7RBPS4I06AItNorpy+MoQNM=` |
| `v1.0.1-vault-3` | `8a6d2ce4ce85ad6e2e8252089e7ba5fa9d7a6104`, 2021-06-11T21:31:12Z | 1.15 | lightweight; commit unsigned | 199 identical; `h1:V95v5KSTu6DB5huDSKiq4uAfILEuNigK/+qPET6H/Mg=` / `h1:XYhtn6ijBSAj6n4YqAaf7RBPS4I06AItNorpy+MoQNM=` |
| `v1.0.1-vault-4` | `6e9815dfdafe4c1a7565b9ef07400890386a53ee`, 2022-10-26T14:49:49Z | 1.15 | lightweight; valid GitHub commit signature | 200 identical; `h1:G9AZNqjH1d5P/vey5/l7HFwkWVJl7vVbLu/zRMcsM4g=` / `h1:XYhtn6ijBSAj6n4YqAaf7RBPS4I06AItNorpy+MoQNM=` |
| `v1.0.1-vault-5` | `e2a59886bba6f1fce980860c8206b8ec9ccd2d91`, 2022-10-13T20:27:02Z | 1.15 | lightweight; commit unsigned | 200 identical; `h1:kI3hhbbyzr4dldA8UdTb7ZlVVlI2DACdCfz31RPDgJM=` / `h1:XYhtn6ijBSAj6n4YqAaf7RBPS4I06AItNorpy+MoQNM=` |
| `v1.0.1-vault-6` | `9371994b9b055f46e21011505707d1168771cac9`, 2024-10-30T17:53:12Z | 1.15 | lightweight; valid GitHub commit signature | 200 identical; `h1:qThxNRouu5cv9LCLZ7pY43TroykqN+Uc7fT3f7tyYh4=` / `h1:XYhtn6ijBSAj6n4YqAaf7RBPS4I06AItNorpy+MoQNM=` |
| `v1.0.1-vault-7` | `02db4972906a1b43a46e2ffb0d2aae2c71875d94`, 2024-11-07T22:23:56Z | 1.15 | lightweight; valid GitHub commit signature | 201 identical; `h1:ag5OxFVy3QYTFTJODRzTKVZ6xvdfLLCA1cy/Y6xGI0I=` / `h1:XYhtn6ijBSAj6n4YqAaf7RBPS4I06AItNorpy+MoQNM=` |
| `v1.0.1-nomad-1` | `955ab59100a7eae83a214b343f67e414fb90ddae`, 2025-06-19T15:11:11Z | 1.14 | annotated SSH-signed tag object `5f0d5fd9`, tagged 2025-06-20T07:14:50Z; GitHub validates tag and commit, local verification lacks an allowed-signers file | 203 identical; `h1:0hOV+/m12cRBAfvHpVOgGdM68XU7uTxGafEuUB2UES8=` / `h1:gwlu9+/P9MmKtYrMsHeFRZPXj2CTPm11TDnMeaRHS7g=` |

The original `v1.0.1-vault` proxy archive identifies retained commit
`809e678` exactly despite the deleted tag. Vault-7 is the highest semantic
version only if these targeted prereleases are included. Its 21-commit branch
adds Vault-specific unused-key position, nested-JSON, and duplicate-key
behavior. Nomad-1 is a separate 25-commit targeted line with decode coercions,
unused-key handling, `hasKey`, and cherry-picked Vault behavior. There is no
stable v1 release after v1.0.0.

Upstream's latest v2 release is v2.24.0 at commit
`6b5068090eef06b1f127f61529db5ba0be7ed343`, published
2025-07-07T13:00:56Z. It is module `github.com/hashicorp/hcl/v2`, declares Go
1.23.0, and would require a production import/API migration. It is both
outside this exact module group and above the retained Go 1.18 floor.

Exact `go get github.com/hashicorp/hcl@v1.0.0` in the candidate archive is a
zero-output, zero-diff operation. Old and candidate states retain 234 selected
modules, 3,557 graph edges, 429 complete test packages, ten loaded HCL
packages, three main consumer packages, and an identical 308-line unapplied
tidy projection. No selection, edge, package, checksum, or tidy line differs.
The loaded path is `plybuild/cmd -> viper ->
viper/internal/encoding/hcl -> hcl`.

Viper's focused HCL codec package and Ply's three loaded consumer packages pass
at count 10 in both states. Required HCL v1.0.0 complete self-tests do not:
both `go test ./...` and `go test -race ./...` exit 1 under exact Go 1.26.7
because default vet rejects `hcl/parser/parser_test.go:243` for formatting an
`*ast.LiteralType` with `%s`. Diagnostic `-vet=off` tests and race pass all 12
packages, so this is historical upstream test drift, but the mission forbids
disabling the required self-test gate. Evaluation stopped before repository
quality execution; no implementation commit or dependency metadata edit was
made.

Independent old/candidate govulncheck v1.7.0 scans preserve exact identical
20 Darwin-symbol, 30 Darwin-module, and 20 Windows-symbol ID populations. The
primary Go vulnerability module index contains no HCL v1 or v2 entry.

## Evidence And Tool Identity

- Rejected/no-change HCL evidence root
  `/private/tmp/ply-p7-hcl-selection.99f508d.CFvyqD` has a fully verified
  21,427-entry manifest SHA-256
  `ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`;
  selection-summary SHA-256 is
  `c20b0d6354c35adf0e05916610d03a93c0210bb613ae889250316bffc34132d9`.
- All accepted prerequisites were freshly verified entry-by-entry before the
  decision: two-entry toolchain recovery `1b30193f...e83b`, 29,812-entry
  go-cmp selection `dcd18baa...1797`, 9,487-entry schema-2 review
  `24ede5a4...1727`, 252,077-entry exact quality `6d8104b4...f69`, and
  30,255-entry regression `00b736b5...3c5`. The accepted go-cmp exact-quality
  scorecard remains `6531ade3...49e1`: 21 stages and all 27 Q0-Q2 rows pass at
  L2, 80/80 mutations, all 8/8 mutation and 4/4 acceptance populations, and
  zero held, regressed, not-comparable, or dirty counts.
- Preserve the btree correction: its regression manifest retains the recorded
  digest and 24,197 entries, but 163 mutable cache/HOME entries no longer
  verify after later enumeration. Its stable regression summary is unchanged;
  the other three btree manifests verify, and fresh go-cmp evidence supersedes
  current regression claims.
- Tool identities remain Go 1.26.7
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
  golangci-lint 2.12.2
  `3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`,
  GoReleaser 2.17.1
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
  apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`,
  and govulncheck v1.7.0
  `0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.
- The first finalization contract run hit the previously observed nested
  partial-raw-log signal-fixture flake and ended non-zero. Its immediate
  complete rerun passed all 62 controls. Passing rerun evidence is sealed at
  `/private/tmp/ply-p7-hcl-handoff-contract.99f508d` with a verified
  three-entry manifest SHA-256
  `ccab44fd1f574e34b5ab040a7d91295ca0b506b745f0a4b95b8fe5f613008043`.

## Next Objective

Independently evaluate selected indirect
`github.com/inconshreveable/mousetrap v1.1.0` as exactly one bounded P7 module
group. Resolve the canonical latest release and every potentially compatible
release from fresh Go proxy, checksum-database, upstream, and primary Go
vulnerability evidence. Do not assume its latest version, Go floor, source
identity, closure, loaded status, or behavior.

Measure exact old/candidate selections, graph edges, complete package
population, checksums, exact-get diff, tidy projection, loaded packages/path,
focused behavior when loaded, candidate module tests, repository quality,
help/API/CLI identity, and vulnerability populations. Implement only an exact
floor-compatible selection with an explained minimal closure and every gate
passing; otherwise record rejection without editing dependency metadata.

Warm caches from separate Git archives. Never use tidy as implementation or
run `go mod download all` inside a measured tree. Stop before another module,
behavior/source changes, language/toolchain or quality-policy changes,
packaging, publication, or P8. Keep all evidence outside the worktree.
