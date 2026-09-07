# Agent Session: Evaluate Armon Consul API Dependency

Status: NEXT
Session ID: `2026-09-07T214638+0200-evaluate-armon-consul-api-dependency`
Created: `2026-09-07T21:46:38+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `45a46dd3188cfabbad82b9ff2bf7c6e295924f1ae11876d46ab479059b3f7760`
Previous: [2026-09-07T194916+0200-evaluate-armon-circbuf-dependency.md](2026-09-07T194916+0200-evaluate-armon-circbuf-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6` as one
bounded dependency group. Resolve canonical latest, release qualification,
default-branch history, and the highest Go-1.18-floor-compatible candidate
from primary evidence. Make an exact dependency selection only if it changes a
selected version, preserves the retained floor through the complete minimal
closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through accepted Circbuf pseudo-version
`v0.0.0-20190214190532-5111143e8da2`. All earlier rejections, no-change
decisions, accepted closures, and evidence corrections remain final. Do not
revisit Circbuf or combine another module group. P8 remains queued.

Project MVS already selects Consul API's 2018 pseudo-version although it is not
an explicit `go.mod` requirement. A minimal post-Circbuf survey finds an empty
exact stable proxy list. Proxy and exact Go `@latest` and `@master` resolve the
same selected pseudo-version at 2018-02-02T20:16:55Z; exact `@v0` has no
matching version. The authoritative repository currently reports public,
enabled, unarchived, non-fork status, default branch `master`, one branch, zero
tags, and zero GitHub Releases. Treat canonical qualification, source identity,
signatures, release and alternate-path history, complete closure, self-tests,
loaded consumers, actual symbols, and vulnerability effect as unknown until
independently resolved. Do not call an untagged pseudo-version a stable release.

# Measurements At Start

Latest dependency implementation is Circbuf commit
`3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
`7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
`a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
`go.sum`. Operator-authorized lifecycle repair
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
the earlier Repr implementation and direct-child Repr handoff. Preserve it.

Circbuf advanced from 2015 pseudo-version
`v0.0.0-20150827004946-bbbad097214e` to canonical latest and master-head
`v0.0.0-20190214190532-5111143e8da2`, commit
`5111143e8da2e98b4ea6a8f32b9065ea1821c191`, tree
`2ab2d9cf2632ab7b549f7da7f081dbe868a697db`, at
2019-02-14T19:05:32Z. It is an unreleased pseudo-version, not a stable release.
The exact path has no stable or prerelease versions, tags, GitHub Releases,
retractions, or authoritative alternate module path. The candidate adds only
one `go.mod` line relative to the selected source; every Go source and test
file is identical. It declares no Go version and no requirements, so its
complete minimal closure is Circbuf alone and preserves Go 1.18.

Candidate proxy and exact-commit forms independently pass verify/list,
count-1/count-10/race tests, and vet without mutation. Exact Go 1.18.10 also
passes count-10/race/vet. Circbuf loads in zero Ply packages; the historical
and selected Serf consumers plus an external Go-1.18 fixture exercise
`NewBuffer`, `Write`, `TotalWritten`, `Size`, `String`, and `Bytes`.

The exact get added one indirect requirement, one main graph edge, and the
candidate checksum pair. Current accepted measurements are 234 selected
modules, 3,581 graph edges, 429 native complete-test packages, 41 loaded
modules, 1,045 `go.sum` lines, and a 361-line unapplied tidy projection. Tidy
removes the explicit Circbuf pin and its two candidate checksum lines while
leaving inherited debt unchanged. Relative to accepted go-cmp commit c314bcb,
accepted metadata adds exactly 29 checksum lines. The main module retains Go
1.18 and toolchain Go 1.26.7. Ordinary and ignored status must be empty.

Circbuf's 430-entry evidence-manifest SHA-256 is
`2684f8c269376320bcfca3404bfbf10c638c81753e6af36b5be59c4df2370d62`;
decision-summary SHA-256 is
`17441163b7924e2cd61cb167b739562da7672fa5abfa10089da36ef906f07ff9`.
Exact quality passed all 27 Q0-Q2 rows at L2 with 80/80 mutations killed and
zero held, regressed, not-comparable, or dirty counts. Full audit exited 1 only
for queued L3 rows Q3.1/Q3.3/Q3.4/Q3.7. Fresh primary vulnerability evidence
has no Circbuf record or trace and preserves exact 20/30/20 Darwin-symbol/
Darwin-module/Windows-symbol populations.

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

From fresh external archives and caches, resolve Consul API versions through
the Go proxy, checksum database, authoritative upstream repository, and
primary Go vulnerability data. Record exact selected and candidate commits/
times, module Go declarations and requirements, checksum pairs, source
identity, tag and commit signature status, release history, and archived/
deprecated state. Explicitly distinguish stable versions, prereleases, pseudo-
versions, retractions, forks, alternate module paths, branch heads, and
unreleased commits. Do not treat the selected pseudo-version as a stable
release without authoritative release evidence.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete closure. Measure old versus
candidate modules, graph edges, complete packages, checksums, loaded packages
and paths, explicit exact-get diff, and `go mod tidy -diff`. Explain every
selection, edge, and checksum change. Identify real consumers and exercise the
actually used symbols.

Require candidate module complete, repeated, and race-enabled self-tests plus
vet; repository verify/build, complete tests, race, vet, Windows build, pinned
lint, byte-identical public help, API/CLI reports, and primary vulnerability
identity. Preserve the retained Go 1.18 floor across the complete changed
closure. Keep standalone test-only closure separate from project MVS.

Reject or retain without dependency edits if canonical qualification, source
identity, floor compatibility, closure, module tests, loaded behavior, or any
applicable quality contract fails. Evaluate an unreleased canonical pseudo-
version under the recorded dependency policy without relabeling it as stable.
Do not manufacture metadata when exact get changes no selected version. Record
precise primary old/candidate vulnerability IDs and traces.

# Required Reading

Verify branch, clean ordinary and ignored status, exact ancestry, reciprocal
archive history, `./codex-dev-start.sh --check`, and P7/P8 state before work.
Read this archive, the rolling handover, P7 roadmap, `go.mod`, `go.sum`, the
answered Circbuf, Optional, Template, Colour, Chroma, Kong, Repr, Assert, and
Units archives, retained Kingpin/Resty/Errgo/Check/YAML outcomes, and the
toolchain, quality, baseline, compatibility, snapshot/Docker, acceptance,
audit, and lifecycle contracts. Do not reopen earlier decisions.

# Three Moves

If and only if a qualified candidate changes the selected version and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/armon/consul-api@<candidate>` for one dependency-only
commit. Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, the toolchain declarations, production source, quality
apparatus, and release inputs.

If canonical qualification confirms the already-selected version and exact
get changes no selection, retain it without adding a redundant direct or
indirect requirement, main edge, or full checksum solely for metadata. A no-
change decision gets no dependency implementation commit.

After a changed selection, run the complete P7 dependency gate: module and
consumer tests, graph/path/checksum/tidy explanation, repository verify/build/
tests/race/vet/Windows/pinned lint, help/API/CLI, launcher and Make contracts,
preflight, host plus fresh snapshot/Docker meta and acceptance, audit meta,
focused and exact Q0-Q2 audits, separate full audit, vulnerability comparison,
empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 rows at L2 and zero held, regressed,
not-comparable, or dirty counts. Full audit may exit 1 only for established
queued L3 rows, never 2.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
create direct `/private/tmp/ply-*` roots, never bypass scratch cleanup, and
never run `go mod download all` inside a measured tree. Warm only exact needed
closures from a separate scratch archive. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Consul API decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
