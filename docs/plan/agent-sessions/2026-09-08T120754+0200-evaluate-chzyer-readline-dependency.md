# Agent Session: Evaluate Chzyer Readline Dependency

Status: NEXT
Session ID: `2026-09-08T120754+0200-evaluate-chzyer-readline-dependency`
Created: `2026-09-08T12:07:54+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `bab4255185a96b081ed86903c6b1572fe18edea6c831ea49cf7029738545d8b4`
Previous: [2026-09-08T110936+0200-evaluate-chzyer-logex-dependency.md](2026-09-08T110936+0200-evaluate-chzyer-logex-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/chzyer/readline v1.5.1` as one bounded dependency group. Resolve
canonical latest, authoritative source and release identity, complete Go-floor
closure, native terminal/platform behavior, Promptui production consumers,
exact MVS effects, and all applicable quality contracts. Make an exact
dependency selection only if a higher exact stable version changes the
selection, preserves the retained Go 1.18 floor through the complete minimal
closure, and passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
and Logex selections. All earlier acceptances, rejections, no-change decisions,
evidence corrections, and lifecycle ancestry are final. Do not revisit Logex,
XXHash, OpenCensus Proto, Crypt, or Speakeasy and do not combine another module
group. P8 remains queued.

Project MVS selects Readline v1.5.1 through the main module's indirect
requirement. Ply reaches it in production through
`github.com/manifoldco/promptui v0.9.0`. Readline v1.5.1 requires
`github.com/chzyer/test v1.0.0`, `github.com/chzyer/logex v1.2.1`, and
`golang.org/x/sys v0.0.0-20220310020820-b874c991c1a5`. Preserve the final
Logex decision and distinguish Promptui's production use of Readline from
Readline's test-only chzyer/test -> Logex path.

A minimal post-Logex survey finds exact proxy versions v1.5.0 and v1.5.1.
Exact `@latest` is already selected v1.5.1 at 2022-07-15T12:48:48Z. It
declares exact module path `github.com/chzyer/readline` and Go 1.15. Repository
tags v1.0 through v1.4 also exist but are absent from the proxy version list;
determine their module/release status rather than treating them as equivalent
proxy releases.

The public repository is currently enabled, unarchived, non-fork, and defaults
to `main`. Selected v1.5.1 is annotated tag object
`704f339125f222987e1fde71641f3185f6eda206` targeting commit
`7f93d88cd5ffa0e805d58d2f9fc3191be15ec668`, tree
`d842017d1ed9d9fd529cce8e199c3a3a69e68e0c`, parent
`8e4bd417b9169c9482a55f3faaeef208b5bf7eb4`, dated
2022-07-15T12:48:48Z. Main is
`9dfc369f8652ba9013dadffd2d2efeada64fe44d`, tree
`c0ed5f5684075d6df7c6e1eb34e15e567e11d3a2`, parent
`fcb4d79af3fbe295b4cb6360e14b8c0b8337353f`, dated
2025-06-20T03:33:30Z. It has three post-release commits. Treat these as survey
facts to verify independently; do not select an unreleased branch head.

# Measurements At Start

The latest dependency implementation remains exact XXHash v2.3.0 commit
`e5d6252825d7a1822c01819b9144050f345a6ad4`, exact parent
`a23ce0f60ad65aac4f4800d6c095d911c0f4e754`, tree
`5213ba55981d77d7c8061312915e237e80d29af8`, changing only `go.mod` and
`go.sum` with three insertions. Logex, Crypt, and OpenCensus Proto remain
retained without dependency edits.

Accepted project measurements are 234 selected modules, 3,583 graph edges,
429 native complete-test packages, 41 loaded modules, 197 loaded packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, current metadata adds exactly 33 checksum
lines. The main module retains Go 1.18 and toolchain Go 1.26.7. Ordinary and
ignored status must be empty.

Fresh primary vulnerability data contains 1,392 module records and no Logex
record. Accepted comparisons retain 20 IDs/22 reachable traces for Darwin and
Windows symbol scans, 22 Darwin package IDs/findings, and 30 Darwin module
IDs. Logex evidence has 636 verified entries; manifest SHA-256 is
`aa33eeb8ea36b17f6830ae7f063fbe1c39ec9ac5d3c9e2417fa2fe844ceb8516`.
Decision-summary SHA-256 is
`6720115c3316e14b89a3bedcb40a59faba4a4f898727051177a682bb1435b49d`.
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
Readline version through the Go proxy and checksum database, authoritative
repository, go-import metadata, and primary Go vulnerability data. Record
commits and times, module declarations and requirements, checksum pairs,
source identity, tag/commit signatures, release/default-branch history,
repository status, deprecation, and retractions. Distinguish stable releases,
prereleases, proxy-absent tags, redirects, forks, alternate module paths, and
unreleased branch heads.

Prove canonical latest and whether any qualified higher exact stable candidate
exists. Preserve Go 1.18 through the complete minimal module and package/test
closure; do not infer the closure floor from the root Go 1.15 declaration.
Inspect Unix/BSD/Linux/Windows source splits, syscall and console paths, build
tags, generated files, and any release-relevant native surface without
expanding into another dependency group.

Measure selected modules, graph edges, native complete-test packages, loaded
modules/packages, checksum lines, exact dependency paths, exact selected-get
effects, and the unapplied tidy projection. Attribute every difference. If
exact selected-version get changes no selected version, do not add a redundant
requirement or checksum.

Identify the exact Promptui -> Readline production consumer chain and the
Readline constructors, interfaces, terminal operations, completers, history,
masking, editing, and OS paths it exercises. Use focused Promptui historical-
consumer and independent behavior fixtures where meaningful. Keep Readline's
native test closure, Promptui production use, Readline's dependency-test-only
chzyer/test/Logex chain, and Ply's actual loaded behavior distinct.

Run source verification, package listing, native complete tests, two
independent repeated-test passes, race where supported, vet, and relevant
cross-builds under exact Go 1.26.7 and a contained Go 1.18 SDK. Existing
Readline v1.5.1 tests pass under both SDKs; vet reports historical nonstandard
`WriteTo` and `ReadRune` signatures. Verify and classify those facts rather
than silently waiving or misattributing them. Treat terminal availability,
missing tests, build tags, flaky tests, release gaps, and checksum differences
precisely.

Project every qualified exact selection in a disposable worktree and run
repository verify, build, count-1/count-10/race/vet, Windows-amd64 build,
pinned lint, byte-identical root/status/upgrade/build help, API/CLI
compatibility and reports, and empty-HOME count-2 before deciding whether an
implementation is permissible. Compare selected and candidate primary
vulnerability results at module, package, symbol, and reachable-trace levels.

Reject or retain if canonical identity, release qualification, complete
closure floor, dependency tests, platform behavior, API compatibility,
historical consumers, projection, or any quality contract fails. Do not
upgrade to an unreleased main head merely because its commit time is newer.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, XXHash implementation commit identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
`docs/plan/quality-handover.md`, `docs/plan/quality-upgrade.md`, `go.mod`,
`go.sum`, the answered Logex, XXHash, OpenCensus Proto, Crypt, and Speakeasy
archives, the earlier dependency archives named in the handover, and every
referenced quality, compatibility, release, runner, evidence, and lifecycle
contract. Earlier outcomes are final.

# Three Moves

If and only if a higher exact stable Readline version is qualified and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/chzyer/readline@<qualified-version>` for one dependency-only
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

After the Readline decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
