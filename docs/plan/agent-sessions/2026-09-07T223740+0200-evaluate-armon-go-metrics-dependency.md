# Agent Session: Evaluate Armon Go Metrics Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-07T223740+0200-evaluate-armon-go-metrics-dependency`
Created: `2026-09-07T22:37:40+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f19f50eeb838b06f48ed0678af2c1fe4563d3b2bf8a20364cdb312b06e822834`
Previous: [2026-09-07T214638+0200-evaluate-armon-consul-api-dependency.md](2026-09-07T214638+0200-evaluate-armon-consul-api-dependency.md)
Next: [2026-09-07T234555+0200-evaluate-armon-go-radix-dependency.md](2026-09-07T234555+0200-evaluate-armon-go-radix-dependency.md)
Outcome: Retained Go Metrics v0.4.0 without dependency edits; v0.4.1 is the
  highest exact-path Go-1.18-compatible release but failed repeated module
  tests and vet, while v0.4.2 and later declare the renamed HashiCorp path.

# Answer

Retain `github.com/armon/go-metrics v0.4.0` without adding dependency
metadata. No dependency implementation commit was created.

## Canonical Identity And Qualification

The exact proxy contains 22 stable versions and no prereleases. Exact
`@latest` and `@v0` resolve v0.6.1 at 2026-07-29T13:06:13Z; `@master` resolves
unreleased pseudo-version `v0.6.2-0.20260907064447-465585286d74` at
2026-09-07T06:44:47Z. Both declare Go 1.25.0 and
`module github.com/hashicorp/go-metrics`, so neither is an exact-path or
retained-floor candidate.

The old repository URL redirects to authoritative public, enabled,
unarchived, non-fork `hashicorp/go-metrics`, with default `master`, 20
branches, 21 current tags, 15 stable GitHub Releases, and 263 default-branch
commits. Proxy v0.4.2 is the first renamed module declaration and aliases
v0.5.0 commit `aee7470331bc2a027cb2711f759f6f527a557ed8`; its alias tag is no
longer present and it has no Release. Exact get rejects v0.4.2's module-path
mismatch. The HashiCorp path is a migration, not an in-place exact-path update.

Therefore v0.4.1 is the highest exact-path candidate. It is stable Release
commit `b6d5c860c07ef6eeec89f4a662c7b452dd4d0c93`, tree
`d582c4e222a01bfe89e45a22b64273970b2008a8`, at
2022-09-08T12:00:52Z, directly after selected v0.4.0 commit
`129ee86de65934631a7fdbeb8c5aa0ec08bfdb6c`, tree
`dc8cd4f53669534104e66bf06f5ebaaf26d6feb4`. Both tags are lightweight and
have no tag-object signature; both merge commits cryptographically verify
with expired GitHub web-flow fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`. The sole retraction is v0.3.11
for an undocumented metrics-sink interface break.

Selected/candidate checksum pairs are
`h1:yCQqn7dwca4ITXb+CbubHmedzaQYHhNhrEXLYUeEe8Q=` /
`h1:E6amYzXo6aW1tqzoZGT755KkbgrJsSdpwZ+3JqfkOG4=` and
`h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=` /
`h1:E6amYzXo6aW1tqzoZGT755KkbgrJsSdpwZ+3JqfkOG4=`. Sumdb agrees. Each
31-file proxy archive matches its Git tag. Version v0.4.1 changes only the Prometheus
sink implementation and test: 86 insertions and four deletions.

## Closure, Tests, And Consumers

Both releases declare Go 1.12 and the same 11 requirements. Proxy and
exact-tag forms select the same 52-module standalone closure and list four
packages; the highest declared Go version anywhere in that closure is 1.12.
The complete closure therefore preserves Go 1.18.

Version v0.4.1 nevertheless fails mandatory qualification. Under exact Go 1.26.7,
proxy and exact-tag forms verify/list, but do not both pass count-1 or race and
both fail count-10 in flaky Statsite/DogStatsD UDP tests. Vet fails both forms
with 46 substantive diagnostics: 29 non-test-goroutine `Fatalf` calls, 12
unkeyed labels, one copied lock, and four empty appends. Exact Go 1.18.10
builds and passes one race run but reproduces count-10 failure and 44 vet
diagnostics. These dependency stop-rule failures reject the candidate.

Viper v1.15.0 supplies selected v0.4.0, but does not import it. Go Metrics
loads in zero Ply packages, Ply has no import, and `go mod why -m` says it is
not needed. Selected Serf v0.10.1 and Memberlist v0.3.0 are real historical
consumers. Their source covers labeled counters/timers/samples, labels,
in-memory/Statsite/StatsD/fanout/global setup, gauges, counters, and timers.
Both compile and vet against v0.4.1 under Go 1.18; focused coordinate and
awareness tests pass count-10 and race. Full legacy suites separately hit
Darwin loopback/test-harness limitations and are not represented as Ply
behavior or candidate-module waivers.

## Projection, Vulnerabilities, And Scope

Exact candidate get selects v0.4.1 but also materializes three indirect
requirements, 14 edges, and four full checksums. Measurements project 234 ->
234 modules, 3,581 -> 3,595 edges, 429 -> 429 packages, 41 -> 41 loaded
modules, zero -> zero loaded Metrics packages, 1,045 -> 1,049 sum lines, and
361 -> 386 tidy-diff lines. Tidy removes every projected requirement and
checksum and restores inherited v0.4.0. The 14 edges are the three main edges
plus v0.4.1's 11 requirements; the four sums are the candidate pair and full
sums for already-selected immutable-radix and golang-lru.

Fresh govulncheck v1.7.0 used primary data updated
2026-09-02T19:12:04Z. Its 1,392-record index has no old- or new-path Metrics
record or trace. Original/candidate results are identical after normalization:
20 IDs/22 traces for Darwin and Windows reachable symbols and 30 Darwin
module IDs. Because dependency stop rules failed, downstream repository,
snapshot/Docker, acceptance, and Q0-Q2 gates were inapplicable and were not
claimed. The accepted Circbuf baseline remains unchanged.

The sealed evidence manifest covers 336 entries at SHA-256
`5afb2b5fb26824c5e4c1db7497b8a6e6dbd24bc20194df7ca7fc5f4385e0dde1`;
decision-summary SHA-256 is
`c5b47a2a7d6f881e7a5aad5895d197027556e238c683bca685f45fff3b3c0592`.
Exact Go 1.26.7 retains binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
The initial Go-1.18 Xcode-resolver setup and invalid govulncheck module-mode
syntax were superseded by direct SDK/clang runs and correct no-pattern
`-scan=module` calls beneath session scratch; neither correction changed the
repository or concealed a failure.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/armon/go-metrics v0.4.0` as one bounded dependency group. Resolve
canonical latest, authoritative source identity, release qualification,
default-branch history, and the highest qualified Go-1.18-floor-compatible
candidate from primary evidence. Make an exact dependency selection only if it
changes a selected version, preserves the retained floor through the complete
minimal closure, and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained canonical-latest Consul API pseudo-
version `v0.0.0-20180202201655-eb2c6b5be1b6`. All earlier rejections,
no-change decisions, accepted closures, and evidence corrections remain final.
Do not revisit Consul API or combine another module group. P8 remains queued.

Project MVS selects Go Metrics v0.4.0 at 2022-05-25T15:01:32Z although it is
not an explicit `go.mod` requirement. A minimal post-Consul survey finds 22
exact stable proxy versions. Exact Go `@latest` and `@v0` resolve v0.6.1 at
2026-07-29T13:06:13Z with `go 1.25.0`, above the retained floor. Exact
`@master` resolves unreleased pseudo-version
`v0.6.2-0.20260907064447-465585286d74` at 2026-09-07T06:44:47Z, also with
`go 1.25.0`. The GitHub request for `armon/go-metrics` currently resolves
repository metadata for `hashicorp/go-metrics`, reporting public, enabled,
unarchived, non-fork status, default branch `master`, 20 branches, 21 tags, and
15 GitHub Releases. Treat the redirect and exact module source identity,
release/tag/signature history, retractions, all version declarations, complete
closures, tests, consumers, and vulnerability effect as unknown until
independently proved. Do not select a Go-1.25 release or guess the highest
Go-1.18-compatible candidate.

# Measurements At Start

Latest dependency implementation remains Circbuf commit
`3be2183ee310ccdc358ce4ed372c0785de25b88b`, exact parent
`7a0ca4e2caba6d2fff20a9c169b181f9187b40dc`, tree
`a51846840b5f252a30359530dfc950f811398431`, changing only `go.mod` and
`go.sum`. The Consul API evaluation made no dependency commit. Operator-
authorized lifecycle repair `f9f0f7669e635c8c7bb169aab816d0ebe1b16435`
remains intentionally between the earlier Repr implementation and direct-child
Repr handoff. Preserve it.

Consul API's exact stable proxy list is empty. Proxy and exact Go `@latest` and
`@master` resolve the already-selected pseudo-version; exact `@v0` has no
match. It is commit `eb2c6b5be1b66bab83016e0b05f01b8d5496ffbd`, tree
`aeb2299aaf107d0823ce91f057798119b821e81b`, at
2018-02-02T20:16:55Z. It is an unreleased pseudo-version, not a stable release.
The exact path has no tags, GitHub Releases, prereleases, or retractions. Its
GitHub web-flow commit signature is cryptographically valid for fingerprint
`5DE3E0509C47EA3CF04A42D34AEE18F83AFDEB23`, whose key is currently expired.
The public repository is operationally unarchived but its README explicitly
deprecates the source in favor of distinct module/import path
`github.com/hashicorp/consul/api`; that path is a migration outside the exact-
path decision, and its current v1.34.4 latest declares Go 1.26.

The selected synthetic `go.mod` has no Go directive or requirements. Its
complete declared minimal closure is Consul API alone and preserves Go 1.18.
All 21 regular proxy files match the exact commit archive and Git tree. The
checksum pair is
`h1:G1bPvciwNyF7IUmKXNt9Ak3m6u9DE1rF+RmtIkBpVdA=` /
`h1:grANhF5doyWs3UAsr3K4I6qtAmlQcZDesFNEHPZAzj8=`, independently confirmed by
sum.golang.org.

Proxy and exact-commit forms verify and list under exact Go 1.26.7, but their
complete count-1, count-10, and race suites fail identically with 33, 330, and
33 tests unable to connect to an external Consul agent at 127.0.0.1:8500. Vet
also fails on two `testing.T.Fatalf` calls from non-test goroutines. Exact Go
1.18.10 compiles the package but reproduces the test and vet failures. These
are dependency stop-rule failures. A historical Crypt consumer and external
Go-1.18 fixture exercise the actually used configuration, client, KV get/list/
put, query, metadata, key, and value symbols; the fixture passes count-10,
race, and vet. Consul API loads in zero Ply packages and Ply imports none of it.

Exact selected-version get changes no selection. Its projection adds only a
redundant indirect requirement, one main graph edge, and the full checksum;
tidy removes all three. Accepted measurements therefore remain 234 selected
modules, 3,581 graph edges, 429 native complete-test packages, 41 loaded
modules, zero loaded Consul API packages, 1,045 `go.sum` lines, and a 361-line
unapplied tidy projection. Relative to accepted go-cmp commit c314bcb,
accepted metadata adds exactly 29 checksum lines. The main module retains Go
1.18 and toolchain Go 1.26.7. Ordinary and ignored status must be empty.

Fresh primary vulnerability evidence contains 1,392 module records and no
Consul API record or trace. Unchanged/exact-get projections preserve exact
20-ID/22-trace Darwin symbol, 30-ID Darwin module, and 20-ID/22-trace Windows
symbol populations. Consul API's 315-entry evidence-manifest SHA-256 is
`1dbe7063e72dade7bc31e9c8967da78a60ef97859b68562ffa1a07b75a0b3b0b`;
decision-summary SHA-256 is
`0f5413813a179949c2dbce2a029becfceb9f650f7a603ac8b19755c3f4d48733`.

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

From fresh external archives and caches, resolve every relevant Go Metrics v0
version through the Go proxy and checksum database, the exact-path repository
redirect and authoritative upstream repository, and primary Go vulnerability
data. Record selected and candidate commits/times, module Go declarations and
requirements, checksum pairs, source identity, tag and commit signatures,
release history, archived/deprecated state, and retractions. Explicitly
distinguish stable versions, prereleases, pseudo-versions, redirects, forks,
alternate module paths, branch heads, and unreleased commits.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete minimal closure, not from a
single module's directive. Measure old versus candidate modules, graph edges,
complete packages, checksums, loaded packages and paths, exact-get diff, and
`go mod tidy -diff`. Explain every selection, edge, and checksum change.
Identify real consumers and exercise the actually used symbols.

Require candidate module complete, repeated, and race-enabled self-tests plus
vet; repository verify/build, complete tests, race, vet, Windows build, pinned
lint, byte-identical public help, API/CLI reports, and primary vulnerability
identity. Preserve the retained Go 1.18 floor across the complete changed
closure. Keep standalone test-only closure separate from project MVS.

Reject or retain without dependency edits if canonical qualification, source
identity, floor compatibility, closure, module tests, loaded behavior, or any
applicable quality contract fails. Do not select v0.6.1 merely because it is
stable when its Go declaration exceeds the retained floor. Do not manufacture
metadata when exact get changes no selected version. Record precise primary
old/candidate vulnerability IDs and traces.

# Required Reading

Verify branch, clean ordinary and ignored status, exact ancestry, reciprocal
archive history, `./codex-dev-start.sh --check`, and P7/P8 state before work.
Read this archive, the rolling handover, P7 roadmap, `go.mod`, `go.sum`, the
answered Consul API, Circbuf, Optional, Template, Colour, Chroma, Kong, Repr,
Assert, and Units archives, retained Kingpin/Resty/Errgo/Check/YAML outcomes,
and the toolchain, quality, baseline, compatibility, snapshot/Docker,
acceptance, audit, and lifecycle contracts. Do not reopen earlier decisions.

# Three Moves

If and only if a qualified candidate changes the selected version and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/armon/go-metrics@<candidate>` for one dependency-only
commit. Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, the toolchain declarations, production source, quality
apparatus, and release inputs.

If canonical qualification confirms v0.4.0 as the highest qualified selection,
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

After the Go Metrics decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
