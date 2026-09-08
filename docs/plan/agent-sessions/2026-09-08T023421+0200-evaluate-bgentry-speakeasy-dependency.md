# Agent Session: Evaluate Bgentry Speakeasy Dependency

Status: NEXT
Session ID: `2026-09-08T023421+0200-evaluate-bgentry-speakeasy-dependency`
Created: `2026-09-08T02:34:21+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ad9291ae115c91e01fa403345b73a551c4d3e1f7321c91ed5f31a13e2beaa808`
Previous: [2026-09-08T005841+0200-evaluate-beorn7-perks-dependency.md](2026-09-08T005841+0200-evaluate-beorn7-perks-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/bgentry/speakeasy v0.1.0` as one bounded dependency group. Resolve
canonical latest, authoritative source identity, release qualification,
default-branch history, and the highest qualified Go-1.18-floor-compatible
candidate from primary evidence. Make an exact dependency selection only if it
changes a selected version, preserves the retained floor through the complete
minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Perks v1.0.1. All earlier rejections,
no-change decisions, accepted closures, and evidence corrections remain final.
Do not revisit Perks or combine another module group. P8 remains queued.

Project MVS selects Speakeasy v0.1.0 at 2017-04-17T20:07:03Z although it is not
an explicit `go.mod` requirement. A minimal post-Perks survey finds exact
stable proxy versions v0.1.0 and v0.2.0. Exact Go `@latest`, `@v0`, and
`@master` resolve v0.2.0 at 2022-09-10T01:20:23Z, commit
`760eaf8b681647364e7a400b856e0921248728a5`. Both proxy module files declare
only exact path `github.com/bgentry/speakeasy`, with no Go directive or
requirements. Do not infer the complete closure floor from missing directives.

The public exact-path GitHub repository currently reports enabled, unarchived,
undisabled, non-fork status, default branch `master`, one branch, two tags, and
one GitHub Release. Master and v0.2.0 point to
`760eaf8b681647364e7a400b856e0921248728a5`; v0.1.0 points to
`4aabc24848ce5fd31929f7d1e4ea74d3709c14cd`. The sole non-draft,
non-prerelease GitHub Release is v0.2.0, published
2024-06-27T20:45:36Z for the older tag. Treat canonical source identity, tag
and commit signatures, release/tag history, retractions, all version and
branch declarations, complete closures, tests, consumers, loaded behavior,
and vulnerability effect as unknown until independently proved. Explicitly
distinguish tag/commit time from GitHub Release publication time.

# Measurements At Start

Latest dependency implementation remains Circbuf commit
`3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
`7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
`a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
`go.sum`. The Consul API, Go Metrics, Go Radix, and Perks evaluations made no
dependency commit. Operator-authorized lifecycle repair
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between the
earlier Repr implementation and direct-child Repr handoff. Preserve it.

Perks remains selected at v1.0.1, the exact-path signed stable latest and
default-branch tip. The enabled, unarchived exact-path repository is a fork,
but go-import and proxy identity prove it is the canonical published
`github.com/beorn7/perks` source; parent `github.com/bmizerany/perks` is a
distinct path with divergent untagged history. Exact selected-version get
changed no selection and only projected redundant metadata, which was not
applied. Perks evidence has 477 verified entries; manifest SHA-256 is
`aa77d2ab9cf676ecfb7ef544c1c5db76ecda86f580ff8ff8f03b8c78038677f0`;
decision-summary SHA-256 is
`6db4ca7260bde6b4affa48c62adb381bda20687f39f013c4b1cdc25397418754`.

Accepted measurements remain 234 selected modules, 3,581 graph edges, 429
native complete-test packages, 41 loaded modules, 1,045 `go.sum` lines, and a
361-line unapplied tidy projection. Relative to accepted go-cmp commit
`c314bcb`, accepted metadata adds exactly 29 checksum lines. The main module
retains Go 1.18 and toolchain Go 1.26.7. Ordinary and ignored status must be
empty.

Fresh primary vulnerability evidence contains 1,392 module records and no
exact or parent Perks record. Accepted projections retain exact 20-ID/22-trace
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
Treat rebuilt tool hashes as nonportable receipts; prove versions and
functionality.

# Role And Boundaries

From fresh external archives and caches, resolve every relevant exact-path
Speakeasy v0 version through the Go proxy and checksum database, authoritative
repository, go-import metadata, and primary Go vulnerability data. Record
selected and candidate commits/times, module Go declarations and requirements,
checksum pairs, source identity, tag and commit signatures, release history,
archived/deprecated status, and retractions. Explicitly distinguish stable
versions, prereleases, pseudo-versions, redirects, forks, alternate module
paths, branch heads, unreleased commits, and a delayed GitHub Release object.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete minimal closure, not from
missing directives. Measure old versus candidate modules, graph edges,
complete packages, checksums, loaded packages and paths, exact-get diff, and
`go mod tidy -diff`. Explain every selection, edge, and checksum change.
Identify real consumers and exercise the actually used symbols and terminal
behavior without requiring an interactive user or leaking terminal state.

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
answered Perks, Go Radix, Go Metrics, Consul API, Circbuf, Optional, Template,
Colour, Chroma, Kong, Repr, Assert, and Units archives, retained
Kingpin/Resty/Errgo/Check/YAML outcomes, and the toolchain, quality, baseline,
compatibility, snapshot/Docker, acceptance, audit, and lifecycle contracts. Do
not reopen earlier decisions.

# Three Moves

If and only if v0.2.0 is qualified and its complete minimal closure preserves
Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/bgentry/speakeasy@v0.2.0` for one dependency-only commit.
Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, toolchain declaration, production source, quality apparatus,
and release input.

If v0.2.0 fails qualification, or exact get changes no selection, retain
v0.1.0 without adding a redundant direct or indirect requirement, main edge,
or checksum solely for metadata. A no-change decision gets no dependency
implementation commit.

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

After the Speakeasy decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
