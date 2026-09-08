# Agent Session: Evaluate Bketelsen Crypt Dependency

Status: NEXT
Session ID: `2026-09-08T044755+0200-evaluate-bketelsen-crypt-dependency`
Created: `2026-09-08T04:47:55+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `dd9227a21c68734836532255961f9eb75da50a592dfef3059b5be686e1d2cf3f`
Previous: [2026-09-08T023421+0200-evaluate-bgentry-speakeasy-dependency.md](2026-09-08T023421+0200-evaluate-bgentry-speakeasy-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/bketelsen/crypt v0.0.3-0.20200106085610-5cbc8cc4026c` as one
bounded dependency group. Resolve canonical latest, authoritative source
identity, release qualification, default-branch history, and the highest
qualified Go-1.18-floor-compatible candidate from primary evidence. Make an
exact dependency selection only if it changes a selected version, preserves
the retained floor through the complete minimal closure, and passes every
quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and the accepted
Speakeasy v0.2.0 dependency move. All earlier rejections, no-change decisions,
accepted closures, and evidence corrections remain final. Do not revisit
Speakeasy or combine another module group. P8 remains queued.

Project MVS selects Crypt pseudo-version
`v0.0.3-0.20200106085610-5cbc8cc4026c` at 2020-01-06T08:56:10Z through
`github.com/devdimensionlab/mvn-pom-mutator v0.2.3`; it is not an explicit
main-module requirement and is not loaded by Ply. A minimal post-Speakeasy
survey finds stable proxy versions v0.0.1 through v0.0.5. Exact Go `@latest`,
`@v0`, and `@master` resolve v0.0.5 at 2021-10-08T10:39:19Z. Treat source
identity, commits, signatures, releases, branches, retractions, declarations,
complete closures, tests, consumers, loaded behavior, project projection, and
vulnerability effect as unknown until independently proved.

The selected pseudo-version and v0.0.5 both declare exact path
`github.com/bketelsen/crypt` and Go 1.12, but their requirements differ
materially. Selected history requires the old CoreOS etcd/Consul API and 2019
Google/crypto/grpc modules; v0.0.5 requires `go.etcd.io/etcd/client/v2
v2.305.0`, Consul API v1.11.0, and newer 2021 Google/crypto/grpc modules. Do
not infer the complete closure floor or a minimal MVS move from the root Go
directive.

# Measurements At Start

The Speakeasy dependency implementation is commit
`41f9561f6ea2f5b6395c5c4d9bcc56a54533133a`, exact parent
`1f55aaa31280a4610ed66c1c37ddda953bfa6a8e`, tree
`a50d4f256fe742e472f1f7e9cf0589a596e971b0`, changing only `go.mod` and
`go.sum` with three insertions. It selects exact stable latest/default-branch
Speakeasy v0.2.0 and preserves Go 1.18 through its one-module closure.

Accepted measurements are 234 selected modules, 3,582 graph edges, 429 native
complete-test packages, 41 loaded modules, 197 loaded packages, 1,047
`go.sum` lines, and a 371-line unapplied tidy projection. Relative to accepted
go-cmp commit `c314bcb`, accepted metadata adds exactly 31 checksum lines. The
main module retains Go 1.18 and toolchain Go 1.26.7. Ordinary and ignored
status must be empty.

Fresh primary vulnerability data contains 1,392 module records and no
Speakeasy record. Accepted old/candidate comparisons retain exact 20-ID/
22-trace Darwin symbol, 30-ID Darwin module, and 20-ID/22-trace Windows symbol
populations.

Speakeasy evidence has 459 verified entries; manifest SHA-256 is
`17bb756dfd0e81c39da3616f6f81edbcc59e299b295d422c1a701be603c02cc4`;
decision-summary SHA-256 is
`d3d5678a31f494e194321951086ccdb7579c70f39a0b7cc311bc4f9c925666f0`.
Exact quality passed all 21 stages and 27 Q0-Q2 rows at L2; its scorecard
SHA-256 is
`48decac359a9ebab23e59c29680e63682ab6d1a65141d13911644143a8db2a01`.

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

# Role And Boundaries

From fresh external archives and caches, resolve every relevant exact-path
Crypt v0 version through the Go proxy and checksum database, authoritative
repository, go-import metadata, and primary Go vulnerability data. Record
selected/candidate commits and times, module declarations and requirements,
checksum pairs, source identity, tag and commit signatures, release and branch
history, archive/deprecation status, and retractions. Distinguish stable tags,
the selected prerelease-form pseudo-version, redirects, forks, alternate
module paths, branch heads, and unreleased commits.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete minimal closure. Measure old
versus candidate modules, graph edges, complete packages, checksums, loaded
packages and paths, exact-get diff, and `go mod tidy -diff`. Explain every
selection, edge, and checksum change; do not silently absorb another retained
dependency decision into the Crypt group.

Identify the historical Crypt consumers and exercise the actually used
configuration/client/backend symbols. Keep standalone test-only closure
separate from project MVS. Require candidate module complete, repeated, and
race-enabled self-tests plus vet; repository verify/build, complete tests,
race, vet, Windows build, pinned lint, byte-identical public help, API/CLI
reports, and primary vulnerability identity.

Reject or retain without dependency edits if canonical qualification, source
identity, floor compatibility, closure, module tests, consumer behavior, or
any applicable quality contract fails. If exact selected-version get changes
no selection, do not manufacture a direct or indirect requirement, main edge,
or checksum solely for metadata.

# Required Reading

Verify branch, clean ordinary and ignored status, exact ancestry, reciprocal
archive history, `./codex-dev-start.sh --check`, and P7/P8 state before work.
Read this archive, the rolling handover, P7 roadmap, `go.mod`, `go.sum`, the
answered Speakeasy, Perks, Go Radix, Go Metrics, Consul API, and Circbuf
archives, all earlier retained outcomes named by the handover, and the
toolchain, quality, baseline, compatibility, snapshot/Docker, acceptance,
audit, and lifecycle contracts. Do not reopen earlier decisions.

# Three Moves

If and only if a higher exact Crypt version is qualified and its complete
minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/bketelsen/crypt@<qualified-version>` for one dependency-only
commit. Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, toolchain declaration, production source, quality apparatus,
and release input. Stop instead of applying an unexplained multi-selection
move.

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

After the Crypt decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next single P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
the next group, launch a successor, push, merge, publish, release, stash,
revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
