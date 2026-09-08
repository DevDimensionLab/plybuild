# Quality Upgrade Handover

Generated: 2026-09-08T04:47:55+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- Latest dependency implementation is Speakeasy commit
  `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a`, exact parent
  `1f55aaa31280a4610ed66c1c37ddda953bfa6a8e`, tree
  `a50d4f256fe742e472f1f7e9cf0589a596e971b0`. It changes only `go.mod` and
  `go.sum`, with three insertions. The next documentation handoff must have
  this dependency commit as its exact parent.
- Earlier Circbuf implementation
  `3be2183ee310ccdc358ce4ed372c0785de25b88b` and operator-authorized
  lifecycle repair `f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remain
  ancestors. The lifecycle repair intentionally remains between the Repr
  implementation and direct-child Repr handoff.
- The answered Speakeasy archive and sole NEXT Crypt archive must link
  reciprocally. No `.agent-task/current.md` or
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and the accepted
Speakeasy v0.2.0 decision. P8 remains queued. All earlier acceptances,
rejections, no-change decisions, evidence corrections, and lifecycle ancestry
are final; do not reopen them or combine another group.

Every disposable cache, projection, archive, source, report, schema-2 document,
generated artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*` roots,
never run `go mod download all` in a measured tree, and preserve the launcher's
scratch cleanup and exact reciprocal archive contract.

## Accepted Speakeasy Group

Select exact-path `github.com/bgentry/speakeasy v0.2.0`. The proxy lists only
stable v0.1.0 and v0.2.0; exact `@latest`, `@v0`, `@master`, repository master,
and tag v0.2.0 all resolve commit
`760eaf8b681647364e7a400b856e0921248728a5` at
2022-09-10T01:20:23Z. Go-import maps the path to the public, enabled,
unarchived, undisabled, non-fork repository
`https://github.com/bgentry/speakeasy.git`, whose default branch has no later
commit. No prerelease, retraction, deprecation, or `/v2` module exists.

Selected v0.1.0 is unsigned commit
`4aabc24848ce5fd31929f7d1e4ea74d3709c14cd`, tree
`4905bd85608c760d3ebf869da51b5740731c50b1`, at
2017-04-17T20:07:03Z. Its signed annotated tag object
`12abe455afa69a18b1098acb18b1f004a4a6c33c` has later tagger time
2017-06-15T22:05:56Z and verifies with fingerprint
`757FD463E177A2F1CD1C89038B6EDBF713E83E69`; the key is now expired.

Candidate v0.2.0 is a lightweight tag with no independent tag object, tagger
time, or tag signature. Its commit has tree
`b2198079290577be7ff0fb5732f778664e28fd64`, parent
`a4057f540bab4628fbdc15a38dfbccdd34e6fa5d`, and a valid GitHub web-flow
signature verified locally with expired fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`. The sole non-draft,
non-prerelease GitHub Release was published 2024-06-27T20:45:36Z for that old
2022 tag. Release publication time is distinct from commit/tag time.

Both eight-file proxy ZIPs match their exact Git tag sources. Sumdb pairs are:

- v0.1.0 source `h1:ByYyxL9InA1OWqxJqqp2A5pYHUrCiAL6K3J+LKSsQkY=` and
  module `h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=`;
- v0.2.0 source `h1:tgObeVOf8WAvtuAX6DhJ4xks4CFNwPDZiqzGqIHE51E=` and
  the same module hash.

Neither root declares Go or requirements. The complete minimal closure is one
module/two packages with only standard-library imports. Proxy and exact-Git
forms of both versions pass verify/list, count-1, repeated count-10, race, and
vet under exact Go 1.26.7 and Go 1.18.10. This execution, not missing
directives, proves retained-floor compatibility.

Selected Mitchellh CLI v1.1.0 is the actual consumer. `BasicUi.AskSecret`
calls only `speakeasy.Ask("")`. Focused count-1/count-10/race and vet pass
against both versions under Go 1.18. A scratch PTY probe waits for ECHO to turn
off before sending a secret, proves the secret never appears in terminal
output, and verifies ECHO restoration. The first immediate-write probe was a
harness scheduling race and was superseded without repository changes.

## Selection, Quality, And Vulnerability Measurements

Exact get changes only Speakeasy v0.1.0 -> v0.2.0. Accepted measurements are
234 modules, 3,582 graph edges, 429 native complete-test packages, 41 loaded
modules, 197 loaded packages, 1,047 `go.sum` lines, and a 371-line unapplied
tidy projection. Relative to accepted go-cmp commit `c314bcb`, accepted
metadata adds exactly 31 checksum lines.

The only new edge is main -> Speakeasy v0.2.0. Mitchellh CLI v1.0.0/v1.1.0
retain their v0.1.0 requests. The only new checksum lines are the candidate
pair. Tidy removes the explicit candidate pin/pair, restores inherited v0.1.0,
and retains all historical debt; it was never used as implementation.
Speakeasy loads in zero Ply packages and `go mod why -m` says the main module
does not need it.

Repository verify/build/count-1/count-10/race/vet, Windows-amd64 build, pinned
golangci-lint 2.12.2, byte-identical root/status/upgrade/build help, API/CLI,
all 62 launcher controls, Make and production-script contracts, complete
preflight, host/fresh snapshot/fresh Docker acceptance, audit meta, focused
manual audit, exact Q0-Q2 audit, expected full audit, and empty-HOME count-2
pass. API/CLI hashes remain `ce39e1c...` and `955f1dda...`.

Exact `make quality` exits 0 with 21/21 stages, 27/27 Q0-Q2 PASS at L2,
80/80 killed mutations, and zero held/regressed/not-comparable/dirty counts.
Q0-Q2 scorecard SHA-256 is
`48decac359a9ebab23e59c29680e63682ab6d1a65141d13911644143a8db2a01`.
Full audit exits expected 1, never 2, only for Q3.1, Q3.3, Q3.4, Q3.7;
scorecard SHA-256 is
`53f8d233daa28ad7c5777bc79c8fa46434ce9989fc0c15719ddd3614a2453c15`.
Snapshot/Docker report hashes are `ceed704a...` and `11a03f09...`.

Fresh govulncheck v1.7.0 uses primary data updated
2026-09-02T19:12:04Z. The 1,392-record index has no Speakeasy module. Old and
candidate results are identical: 20 IDs/22 traces for Darwin and Windows
reachable symbols and 30 Darwin module IDs. Normalized reachable SHA-256 is
`7757df547ed97c0b709bfe8cbbbf807356331727c870bdeb7f0602b1f393ba79`.
Exact IDs and 22 endpoint traces are recorded in the answered archive.

Speakeasy evidence contains 459 verified entries; manifest SHA-256 is
`17bb756dfd0e81c39da3616f6f81edbcc59e299b295d422c1a701be603c02cc4`.
Decision-summary SHA-256 is
`d3d5678a31f494e194321951086ccdb7579c70f39a0b7cc311bc4f9c925666f0`.

## Tools And Runner Corrections

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Exact Go 1.18.10 binary/archive values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Portable receipts remain golangci-lint 2.12.2 archive `a9c54498...`,
  GoReleaser 2.17.1 binary `f5f08a77...`, and apidiff `0c55d9e3...`.
  Rebuilt binaries are nonportable; version/build metadata and functionality
  are the proof.
- Scratch-only corrections were exact historical compatibility graph warming,
  external report routing, golangci cache relocation, BSD-`mktemp` adaptation,
  cleared recursive Make overrides, vulnerability normalization, and terminal-
  state-aware PTY timing. A known launcher partial-log timing probe flaked once;
  an immediate standalone repeat and full preflight repeat passed all 62
  controls. None changed repository source or concealed a candidate failure.

## Retained Decisions

Retain all earlier exact selections and outcomes, including Repr v0.5.4,
Assert v1.0.0, Units pseudo-version `0f3dac36c52b`, Chroma v0.10.0, Colour
v0.1.0, Template pseudo-version `fb15b899a751`, Optional v1.0.0, Circbuf
pseudo-version `5111143e8da2`, Consul API `eb2c6b5be1b6`, Go Metrics v0.4.0,
Go Radix v1.0.0, Perks v1.0.1, Kingpin v2.2.6, Resty v1.12.0, Errgo v2.1.0,
Check pseudo-version `41f04d3bba15`, YAML v2.4.0/v3.0.1, and
`go.yaml.in/yaml/v3 v3.0.5`. Preserve the accepted Kong pseudo-version and all
other decisions recorded in the roadmap and answered archives.

## Next Objective

Independently evaluate selected exact-path `github.com/bketelsen/crypt`
pseudo-version `v0.0.3-0.20200106085610-5cbc8cc4026c` as the next single P7
group. MVS selects it through `github.com/devdimensionlab/mvn-pom-mutator
v0.2.3`; it is neither an explicit main requirement nor loaded by Ply.

A minimal survey finds stable proxy versions v0.0.1 through v0.0.5. Exact
`@latest`, `@v0`, and `@master` resolve v0.0.5 at
2021-10-08T10:39:19Z. Both selected and latest roots declare exact path
`github.com/bketelsen/crypt` and Go 1.12. Their requirements differ materially:
v0.0.5 replaces the old CoreOS etcd/Consul and 2019 Google/crypto/grpc roots
with `go.etcd.io/etcd/client/v2 v2.305.0`, Consul API v1.11.0, and newer 2021
Google/crypto/grpc modules. Independently prove source and release identity,
signatures, complete closure floor, tests, real consumers, projection, and
vulnerability effect. Do not implement or combine another group until Crypt's
bounded decision is complete.
