# Agent Session: Evaluate Chzyer Logex Dependency

Status: NEXT
Session ID: `2026-09-08T110936+0200-evaluate-chzyer-logex-dependency`
Created: `2026-09-08T11:09:36+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5db5fb2785e935339a637b6ecdf63596c9f1f785775d75897866583851c915a4`
Previous: [2026-09-08T074422+0200-evaluate-cespare-xxhash-v2-dependency.md](2026-09-08T074422+0200-evaluate-cespare-xxhash-v2-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/chzyer/logex v1.2.1` as one bounded dependency group. Resolve
canonical latest, authoritative source and release identity, complete Go-floor
closure, native and historical-consumer behavior, exact MVS effects, and all
applicable quality contracts. Make an exact dependency selection only if a
higher exact stable version changes the selection, preserves the retained Go
1.18 floor through the complete minimal closure, and passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, the accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt and OpenCensus
Proto selections. All earlier acceptances, rejections, no-change decisions,
evidence corrections, and lifecycle ancestry are final. Do not revisit XXHash,
OpenCensus Proto, Crypt, or Speakeasy and do not combine another module group.
P8 remains queued.

Project MVS selects Logex v1.2.1 through direct historical requirements from
`github.com/chzyer/readline v1.5.1` and `github.com/chzyer/test v1.0.0`.
Promptui v0.9.0 and historical Google pprof versions request lower v1.1.10.
Ply reaches Readline through `github.com/manifoldco/promptui v0.9.0`; Logex is
needed by Readline's test closure but no Logex package is loaded by Ply's own
package/test list.

A minimal post-XXHash survey finds exact proxy versions v1.1.1 through v1.1.10
plus v1.2.0 and v1.2.1. Exact `@latest` is already selected v1.2.1 at
2022-04-24T13:13:51Z. It declares exact module path
`github.com/chzyer/logex` and Go 1.15. Repository tags v1.0 and v1.1 also exist
but are absent from the proxy version list; determine their module/release
status rather than treating them as equivalent proxy releases.

Go-import metadata maps the module to
`https://github.com/chzyer/logex.git`. The public repository is currently
enabled, unarchived, non-fork, and defaults to `master`. Selected v1.2.1 is
lightweight tag/commit `2f95bdde8c3c97bfbf6d016fcc410669a895b9e7`, tree
`36fcd9ac7d56d659872b2a6576ca6f66385c938a`, exact parent
`a21c317abc1e9a4f23ed3455107a4d20375735cc`, dated
2022-04-24T13:13:51Z. Master is
`5a7e37d2e8a8bbe3ef54984ab949eebaa948b8b4`, tree
`f9cc17fbf471a8558b15bdc07ee3e1a9dba4631d`, with two parents, dated
2024-04-02T15:49:33Z. It has five post-release test/CI commits. Treat these as
survey facts to verify independently; do not select an unreleased branch head.

# Measurements At Start

The latest dependency implementation is exact XXHash v2.3.0 commit
`e5d6252825d7a1822c01819b9144050f345a6ad4`, exact parent
`a23ce0f60ad65aac4f4800d6c095d911c0f4e754`, tree
`5213ba55981d77d7c8061312915e237e80d29af8`, changing only `go.mod` and
`go.sum` with three insertions. Crypt and OpenCensus Proto remain retained
without dependency edits.

Accepted project measurements are 234 selected modules, 3,583 graph edges,
429 native complete-test packages, 41 loaded modules, 197 loaded packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, current metadata adds exactly 33 checksum
lines. The main module retains Go 1.18 and toolchain Go 1.26.7. Ordinary and
ignored status must be empty.

Fresh primary vulnerability data contains 1,392 module records and no XXHash
record. Accepted post-XXHash comparisons retain 20 IDs/22 reachable traces for
Darwin and Windows symbol scans, 22 Darwin package IDs/findings, and 30 Darwin
module IDs. XXHash evidence has 608 verified entries; manifest SHA-256 is
`6f57be6cbf177c6617ebf62e0eea7b38f6b3b41eb00873060e160ddf02f1cf34`.
Decision-summary SHA-256 is
`4e2f0559dbefd98f22e8f1efe34a1e538ecb8e1dd2da0749af8478c88c0a1a20`.
The authoritative Q0-Q2 scorecard SHA-256 is
`579e5b135db2403904943bfe71c3ced987f59f30007cbdc0dec88a46d19b42aa`.

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
Logex version through the Go proxy and checksum database, authoritative
repository, go-import metadata, and primary Go vulnerability data. Record
commits and times, module declarations and requirements, checksum pairs,
source identity, tag/commit signatures, release/default-branch history,
repository status, deprecation, and retractions. Distinguish stable releases,
prereleases, proxy-absent tags, redirects, forks, alternate module paths, and
unreleased branch heads.

Prove canonical latest and whether any qualified higher exact stable candidate
exists. Preserve Go 1.18 through the complete minimal module and package/test
closure; do not infer the closure floor from the root Go 1.15 declaration.
Inspect OS/architecture-specific source, build tags, generated files, and any
release-relevant native surface without expanding into another dependency
group.

Measure selected modules, graph edges, native complete-test packages, loaded
modules/packages, checksum lines, exact dependency paths, exact selected-get
effects, and the unapplied tidy projection. Attribute every difference. If
exact selected-version get changes no selected version, do not add a redundant
requirement or checksum.

Identify the exact Promptui -> Readline -> test -> Logex consumer chain and the
Logex packages, constructors, logging methods, helpers, interfaces, and OS
paths it exercises. Use focused historical-consumer and independent behavior
fixtures where meaningful without inventing a Ply runtime path. Keep Logex's
native test/tool closure, historical dependency-test consumers, and Ply's
unloaded project behavior distinct.

Run source verification, package listing, native complete tests, two
independent repeated-test passes, race where supported, vet, and relevant
cross-builds under exact Go 1.26.7 and a contained Go 1.18 SDK. Treat missing
tests, build tags, vet findings, flaky tests, release gaps, and checksum
differences precisely; decide whether each is a release disqualifier rather
than silently waiving it.

Project every qualified exact selection in a disposable worktree and run
repository verify, build, count-1/count-10/race/vet, Windows-amd64 build,
pinned lint, byte-identical root/status/upgrade/build help, API/CLI
compatibility and reports, and empty-HOME count-2 before deciding whether an
implementation is permissible. Compare selected and candidate primary
vulnerability results at module, package, symbol, and reachable-trace levels.

Reject or retain if canonical identity, release qualification, complete
closure floor, dependency tests, platform behavior, API compatibility,
historical consumers, projection, or any quality contract fails. Do not
upgrade to an unreleased master head merely because its commit time is newer.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, XXHash implementation commit identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
`docs/plan/quality-handover.md`, `docs/plan/quality-upgrade.md`, `go.mod`,
`go.sum`, the answered XXHash, OpenCensus Proto, Crypt, and Speakeasy archives,
the earlier dependency archives named in the handover, and every referenced
quality, compatibility, release, runner, evidence, and lifecycle contract.
Earlier outcomes are final.

# Three Moves

If and only if a higher exact stable Logex version is qualified and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/chzyer/logex@<qualified-version>` for one dependency-only
commit. Do not hand-edit metadata or use tidy as implementation. Preserve
every retained version, Go 1.18, toolchain Go 1.26.7, production source,
quality apparatus, and release input. Stop instead of applying an unexplained
multi-selection move. The incoming survey finds selected already canonical
latest, so any changed selection requires new primary evidence, not a branch
head or redundant pin.

After a changed selection, run the complete P7 dependency gate: dependency and
consumer tests; graph/path/checksum/tidy proof; repository verify/build/tests/
race/vet/Windows/pinned lint; help/API/CLI; launcher and Make contracts;
preflight; host plus fresh snapshot/Docker meta and acceptance; audit meta;
focused and exact Q0-Q2 audits; separate full audit; vulnerability comparison;
empty-HOME count-2; and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 Q0-Q2 rows at L2 and zero held,
regressed, not-comparable, or dirty counts. Full audit may exit 1 only for the
established queued L3 rows, never 2.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
create direct `/private/tmp/ply-*` roots, never bypass scratch cleanup, and
never run `go mod download all` in a measured worktree. Warm only exact needed
closures from a separate scratch archive. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Logex decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next single P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
the next group, launch a successor, push, merge, publish, release, stash,
revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
