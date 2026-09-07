# Agent Session: Evaluate Armon Go Radix Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-07T234555+0200-evaluate-armon-go-radix-dependency`
Created: `2026-09-07T23:45:55+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `cdf0e0ed88ada0236974a6cdd2e680a35aa4d7170b1d5dcd83ba5dbe4dcace4b`
Previous: [2026-09-07T223740+0200-evaluate-armon-go-metrics-dependency.md](2026-09-07T223740+0200-evaluate-armon-go-metrics-dependency.md)
Next: [2026-09-08T005841+0200-evaluate-beorn7-perks-dependency.md](2026-09-08T005841+0200-evaluate-beorn7-perks-dependency.md)
Outcome: Retained Go Radix v1.0.0 without dependency edits; it is the sole
  stable exact-path version and canonical latest, while master is a tested but
  unreleased pseudo-version with no tag or GitHub Release.

# Answer

Retain `github.com/armon/go-radix v1.0.0` without adding dependency metadata.
No dependency implementation commit was created.

## Canonical Identity And Qualification

The exact Go proxy contains only stable v1.0.0. Exact `@latest` and `@v1`
resolve that selected version at 2018-08-24T02:57:28Z. Exact `@master`
resolves pseudo-version `v1.0.1-0.20221118154546-54df44f2176c` at
2022-11-18T15:45:46Z. It is an unreleased branch head, not a semantic tag or
GitHub Release.

Go-import metadata identifies `https://github.com/armon/go-radix.git` as the
canonical exact-path source. The public repository is enabled, unarchived,
and non-fork, with default `master`, one branch, one lightweight tag, and zero
GitHub Releases. Stable tag v1.0.0 identifies commit
`1a2de0c21c94309923825da3df33a4381872c795`, tree
`8c6d01daaee6076244d5f41247608c75a8ad4224`, parent
`7fddfc383310abc091d79a27f116d30cf0424032`. A stable tag is release-qualified
here even though the repository publishes no GitHub Release object.

The lightweight tag has no tag-object signature. The tagged commit's OpenPGP
signature verifies locally with fingerprint
`7A01BBD67E7E8ADD50E00714744E147AA52F5B0A`; GitHub currently reports
`unknown_key` rather than a verified identity binding. Master is merge commit
`54df44f2176c4a553657a4f0dbe6fdb108288be3`, tree
`87f38e748e5fc5c602ba3296793b25a782fbf02d`. GitHub reports its signature
valid; local cryptographic verification identifies expired web-flow key
fingerprint `5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`.

Version v1.0.0 is an ancestor of master. The ten later commits change only
`.travis.yml`, `radix.go`, and `radix_test.go`, with 86 insertions and 13
deletions. No proxy prereleases, alternate major paths, renamed authoritative
paths, retractions, or deprecation declarations were found. Newer master code
passes execution, but it has no authoritative release qualification and is
not selected merely because it is newer. V1.0.0 therefore remains the highest
qualified exact-path version.

## Source, Closure, And Tests

The selected module/checksum pair is
`h1:F4z6KzEeeQIMeLFa97iZU6vupzoecKdU5TX24SNppXI=` /
`h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=`. The master pseudo-version
pair is `h1:651/eoCRnQ7YtSjAnSzRucrJz+3iGEFt+ysraELS81M=` with the same
`go.mod` checksum. Sumdb agrees. Proxy ZIP SHA-256 values are respectively
`df93c816505baf12c3efe61328dc6f8fa42438f68f80b0b3725cae957d021c90`
and `f261a1141112f65564cec8f652ef6abf1d654228275e4fec6220b98337667c13`.

Each proxy ZIP's seven regular files match the corresponding exact Git
archive byte-for-byte. Normalized manifest SHA-256 values are
`f7801af436dc78ae44f1611f9160c6c37c3355ca3729c236c91f0fbd8834c518`
and `18854160a6271e82a14e65bcc15d142d83b640f78b087186b095176c162a1895`.

Neither form declares a Go version or requirements. Each complete minimal
module closure contains only Go Radix and one package. Missing directives were
not treated as compatibility proof: writable proxy and exact-Git forms pass
verification/listing, count-1, count-10, race, and vet under exact Go 1.26.7;
proxy forms also pass every gate under exact Go 1.18.10. The complete closures
therefore preserve the retained Go 1.18 floor by direct execution.

Selected `github.com/mitchellh/cli v1.1.0` is the actual graph consumer and
uses `New`, `Insert`, `Get`, `Walk`, `WalkPrefix`, `LongestPrefix`, `Tree`, and
`WalkFn`. Serf v0.9.6/v0.10.1 only declare Radix indirectly through CLI.
Focused nested-command, help, autocomplete, and subcommand consumer tests pass
count-1, count-10, race, and vet against both Radix versions. Go Radix loads in
zero Ply packages; Ply has no source import and `go mod why -m` says it is not
needed by the main module.

## Projection, Quality, And Vulnerabilities

Accepted measurements remain 234 selected modules, 3,581 graph edges, 429
complete packages, 41 loaded modules, zero loaded Radix packages, 1,045
`go.sum` lines, and 361 tidy-diff lines. Exact get of selected v1.0.0 changes
no selection; it only projects a redundant indirect requirement, one main
edge, and the full module checksum: 234/3,582/429/41/0, 1,046 sum lines, and
367 tidy-diff lines. Tidy removes that materialized metadata.

Exact get of master changes only Radix's selected version, but also adds an
explicit main edge and its checksum pair: 234 modules, 3,582 edges, 429
packages, 41 loaded modules, zero loaded Radix packages, 1,047 sum lines, and
369 tidy-diff lines. Tidy removes the pin and pair and restores inherited
v1.0.0. No projection was applied.

Baseline and master projections pass repository verification/build, complete
count-1/count-10/race tests, vet, Windows-amd64 build, and pinned
golangci-lint 2.12.2. Root, status, upgrade, and build help are byte-identical.
API and CLI reports are byte-identical at SHA-256
`ce39e1c47389f8a3b9699e0f6b165f974f2b0b0f8a3a1553c4b287e1ece005f4`
and `955f1dda0de5d1e52f7ddc8368426bfe9a32b6e42716f1608428ca83dad645b2`.
Corrected scratch-local preflight passes all 62 launcher controls, every
distribution/lint/install/toolchain and script meta-contract, and all 15 audit
controls. Empty-HOME count-2 passes.

Fresh govulncheck v1.7.0 used primary data updated
2026-09-02T19:12:04Z. The 1,392-record module index has no Go Radix entry or
trace. Baseline/master results are byte-identical: 20 IDs/22 traces for Darwin
and Windows reachable-symbol scans and 30 Darwin module IDs.

Canonical qualification produced no changed selection. Consequently the
changed-selection-only host/snapshot/Docker acceptance, exact `make quality`,
focused/Q0-Q2/full audits, and dependency commit were inapplicable and are not
claimed. Source-prefix normalization, historical compatibility-cache warming,
reachable-trace filtering, and scratch-local BSD-mktemp/MAKEOVERRIDES runner
corrections changed no repository file and concealed no failure.

The sealed evidence manifest contains 390 verified entries and has SHA-256
`6be6b857e77c7097b402ea5f4fe0ce849e15382b7167a004448cd85d82926eaf`;
decision-summary SHA-256 is
`493e832f2d0e6456bb64462104e1bd5b80f4b6eb32d401b9b098acdd0dd95e6a`.
Exact Go 1.26.7 retained required binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/armon/go-radix v1.0.0` as one bounded dependency group. Resolve
canonical latest, authoritative source identity, release qualification,
default-branch history, and the highest qualified Go-1.18-floor-compatible
candidate from primary evidence. Make an exact dependency selection only if it
changes a selected version, preserves the retained floor through the complete
minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Go Metrics v0.4.0. All earlier
rejections, no-change decisions, accepted closures, and evidence corrections
remain final. Do not revisit Go Metrics or combine another module group. P8
remains queued.

Project MVS selects Go Radix v1.0.0 at 2018-08-24T02:57:28Z although it is not
an explicit `go.mod` requirement. A minimal post-Metrics survey finds one
exact stable proxy version. Exact Go `@latest` and `@v1` resolve the selected
v1.0.0 with no exposed Go declaration. Exact `@master` resolves unreleased
pseudo-version `v1.0.1-0.20221118154546-54df44f2176c` at
2022-11-18T15:45:46Z, also with no exposed Go declaration. Do not infer floor
compatibility from a missing directive or select an unreleased branch head
without qualification.

The public GitHub repository currently reports enabled, unarchived, non-fork
status, default branch `master`, one branch, one tag, and zero GitHub Releases.
Master is commit `54df44f2176c4a553657a4f0dbe6fdb108288be3`; tag v1.0.0 points to
`1a2de0c21c94309923825da3df33a4381872c795`. Treat canonical source identity,
tag and commit signatures, release/tag history, retractions, all version and
branch declarations, complete closures, tests, consumers, loaded behavior,
and vulnerability effect as unknown until independently proved. Explicitly
distinguish a stable tag from a GitHub Release and an unreleased pseudo-version.

# Measurements At Start

Latest dependency implementation remains Circbuf commit
`3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
`7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
`a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
`go.sum`. The Consul API and Go Metrics evaluations made no dependency commit.
Operator-authorized lifecycle repair `f9f0f7669e635c8c7bb169aab816d0ebe1b16435`
remains intentionally between the earlier Repr implementation and direct-child
Repr handoff. Preserve it.

Go Metrics remains selected at v0.4.0. V0.4.1 was the highest exact-path
Go-1.18-compatible stable candidate, but its proxy and tag forms failed
repeated module tests and vet. Proxy v0.4.2 and later declare distinct module
path `github.com/hashicorp/go-metrics`; current v0.6.1 and master also declare
Go 1.25.0. No dependency metadata was changed. Its 336-entry evidence-
manifest SHA-256 is
`5afb2b5fb26824c5e4c1db7497b8a6e6dbd24bc20194df7ca7fc5f4385e0dde1`;
decision-summary SHA-256 is
`c5b47a2a7d6f881e7a5aad5895d197027556e238c683bca685f45fff3b3c0592`.

Accepted measurements remain 234 selected modules, 3,581 graph edges, 429
native complete-test packages, 41 loaded modules, 1,045 `go.sum` lines, and a
361-line unapplied tidy projection. Relative to accepted go-cmp commit
`c314bcb`, accepted metadata adds exactly 29 checksum lines. The main module
retains Go 1.18 and toolchain Go 1.26.7. Ordinary and ignored status must be
empty.

Fresh primary vulnerability evidence contains 1,392 module records and no Go
Metrics record or trace. Accepted projections retain exact 20-ID/22-trace
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

From fresh external archives and caches, resolve every relevant Go Radix v1
version through the Go proxy and checksum database, the exact-path repository,
and primary Go vulnerability data. Record selected and candidate commits/times,
module Go declarations and requirements, checksum pairs, source identity, tag
and commit signatures, release history, archived/deprecated state, and
retractions. Explicitly distinguish stable versions, prereleases, pseudo-
versions, redirects, forks, alternate module paths, branch heads, and
unreleased commits.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete minimal closure, not from a
single module's directive or the absence of one. Determine whether the branch
head is a legitimate candidate despite lacking a release. Measure old versus
candidate modules, graph edges, complete packages, checksums, loaded packages
and paths, exact-get diff, and `go mod tidy -diff`. Explain every selection,
edge, and checksum change. Identify real consumers and exercise the actually
used symbols.

Require candidate module complete, repeated, and race-enabled self-tests plus
vet; repository verify/build, complete tests, race, vet, Windows build, pinned
lint, byte-identical public help, API/CLI reports, and primary vulnerability
identity. Preserve the retained Go 1.18 floor across the complete changed
closure. Keep standalone test-only closure separate from project MVS.

Reject or retain without dependency edits if canonical qualification, source
identity, floor compatibility, closure, module tests, loaded behavior, or any
applicable quality contract fails. Do not select `@master` merely because it is
newer, and do not manufacture metadata when exact get changes no selected
version. Record precise primary old/candidate vulnerability IDs and traces.

# Required Reading

Verify branch, clean ordinary and ignored status, exact ancestry, reciprocal
archive history, `./codex-dev-start.sh --check`, and P7/P8 state before work.
Read this archive, the rolling handover, P7 roadmap, `go.mod`, `go.sum`, the
answered Go Metrics, Consul API, Circbuf, Optional, Template, Colour, Chroma,
Kong, Repr, Assert, and Units archives, retained Kingpin/Resty/Errgo/Check/YAML
outcomes, and the toolchain, quality, baseline, compatibility, snapshot/Docker,
acceptance, audit, and lifecycle contracts. Do not reopen earlier decisions.

# Three Moves

If and only if a qualified candidate changes the selected version and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/armon/go-radix@<candidate>` for one dependency-only commit.
Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, the toolchain declarations, production source, quality
apparatus, and release inputs.

If canonical qualification confirms v1.0.0 as the highest qualified selection,
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

After the Go Radix decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
