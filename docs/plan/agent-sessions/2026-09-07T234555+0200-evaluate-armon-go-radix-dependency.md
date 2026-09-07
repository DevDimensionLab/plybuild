# Agent Session: Evaluate Armon Go Radix Dependency

Status: NEXT
Session ID: `2026-09-07T234555+0200-evaluate-armon-go-radix-dependency`
Created: `2026-09-07T23:45:55+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `cdf0e0ed88ada0236974a6cdd2e680a35aa4d7170b1d5dcd83ba5dbe4dcace4b`
Previous: [2026-09-07T223740+0200-evaluate-armon-go-metrics-dependency.md](2026-09-07T223740+0200-evaluate-armon-go-metrics-dependency.md)
Next: none
Outcome: pending

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
