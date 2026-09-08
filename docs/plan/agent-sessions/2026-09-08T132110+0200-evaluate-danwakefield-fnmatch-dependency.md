# Agent Session: Evaluate Danwakefield Fnmatch Dependency

Status: NEXT
Session ID: `2026-09-08T132110+0200-evaluate-danwakefield-fnmatch-dependency`
Created: `2026-09-08T13:21:10+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `1209a621bfdabf5da6d8a5016eb038e3b5f65d5c09a68d62a1ae91410e77a404`
Previous: [2026-09-08T120754+0200-evaluate-chzyer-readline-dependency.md](2026-09-08T120754+0200-evaluate-chzyer-readline-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/danwakefield/fnmatch
v0.0.0-20160403171240-cbb64ac3d964` as one bounded dependency group. Resolve
canonical latest, authoritative source and release identity, complete Go-floor
closure, native matching behavior, historical and actual consumers, exact MVS
effects, and all applicable quality contracts. Make an exact dependency
selection only if a higher qualified exact stable release changes the
selection, preserves the retained Go 1.18 floor through the complete minimal
closure, and passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
Logex, and Readline selections. All earlier acceptances, rejections, no-change
decisions, evidence corrections, and lifecycle ancestry are final. Do not
revisit Readline, Logex, XXHash, OpenCensus Proto, Crypt, or Speakeasy and do
not combine another module group. P8 remains queued.

Project MVS selects Fnmatch through the main module's indirect requirement,
but `go mod why -m` says the main module does not need it and no Fnmatch package
is loaded. The graph retains a historical
`github.com/alecthomas/chroma v0.7.1 -> Fnmatch` edge while project MVS selects
Chroma v0.10.0, whose loaded production packages do not use Fnmatch. Preserve
the final Chroma decision and distinguish an old graph requirement from a
loaded production or test consumer.

A minimal post-Readline survey finds an empty exact proxy version list. Exact
`@latest` is the selected pseudo-version at 2016-04-03T17:12:40Z. Its proxy-
synthesized module file declares exact path `github.com/danwakefield/fnmatch`
without a Go directive or requirements. The authoritative repository has no
tags or GitHub Releases. Determine whether any redirected, forked, alternate-
path, tagged, prerelease, or unreleased identity can qualify without treating
an untagged commit or empty version list as a stable release.

The public `github.com/danwakefield/fnmatch` repository is currently enabled,
unarchived, non-fork, and defaults to `master`. Selected commit and master are
both `cbb64ac3d964b81592e64f957ad53df015803288`, parent
`eb9738ef552dd59a56a5953a4de6216f70564908`, dated
2016-04-03T17:12:40Z. It is unsigned according to GitHub. Treat these as a
minimal incoming survey to verify independently, including source ancestry and
whether repository activity timestamps reflect any source commit beyond the
four-commit master history.

# Measurements At Start

The latest dependency implementation remains exact XXHash v2.3.0 commit
`e5d6252825d7a1822c01819b9144050f345a6ad4`, exact parent
`a23ce0f60ad65aac4f4800d6c095d911c0f4e754`, tree
`5213ba55981d77d7c8061312915e237e80d29af8`, changing only `go.mod` and
`go.sum` with three insertions. Readline, Logex, Crypt, and OpenCensus Proto
remain retained without dependency edits.

Accepted project measurements are 234 selected modules, 3,583 graph edges,
429 native complete-test packages, 41 loaded modules, 197 loaded packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, current metadata adds exactly 33 checksum
lines. The main module retains Go 1.18 and toolchain Go 1.26.7. Ordinary and
ignored status must be empty.

Fresh primary vulnerability data contains 1,392 module records and no Fnmatch
record. Accepted comparisons retain 20 IDs/22 reachable traces for Darwin and
Windows symbol scans, 22 Darwin package IDs/findings, and 30 Darwin module
IDs. Readline evidence has 745 verified entries; manifest SHA-256 is
`8ab036d5f96a8d92ebf682ed0fef1d5f116621dff18d670bb200172581158a00`.
Decision-summary SHA-256 is
`435eafc93ae6df466970eb1af57a8127e5f8c389863c82fffadec236ca6bbf8a`.
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
Fnmatch version or commit identity through the Go proxy and checksum database,
authoritative repository, go-import metadata, upstream ancestry, and primary
Go vulnerability data. Record commits and times, module declarations and
requirements, checksum pairs, source identity, commit signatures, default-
branch history, repository status, deprecation, and retractions. Distinguish
stable releases, prereleases, pseudo-versions, proxy-absent tags, redirects,
forks, alternate module paths, and unreleased branch heads.

Prove canonical latest and whether any qualified higher exact stable candidate
exists. Preserve Go 1.18 through the complete minimal module and package/test
closure; do not infer the closure floor from a proxy-synthesized module file
with no Go directive. Inspect build tags, generated files, OS/architecture-
specific code, native or Cgo surface, and release-relevant matching behavior
without expanding into another dependency group.

Measure selected modules, graph edges, native complete-test packages, loaded
modules/packages, checksum lines, exact dependency paths, exact selected-get
effects, and the unapplied tidy projection. Attribute every difference. Keep
the root indirect requirement, the historical Chroma v0.7.1 graph edge,
selected Chroma v0.10.0, loaded packages, and actual Ply behavior distinct. Do
not use tidy as an implementation or silently convert this dependency review
into an unrelated stale-requirement cleanup.

Identify Fnmatch's exported API, flag semantics, wildcard/bracket/escape/
separator/case behavior, malformed-pattern handling, and platform assumptions.
Run independent behavior fixtures where the native suite is missing or weak.
Find exact historical consumers and prove whether any selected production or
test package loads Fnmatch. Preserve the accepted Chroma outcome.

Run source verification, package listing, native complete tests, two
independent repeated-test passes, race where supported, vet, and relevant
cross-builds under exact Go 1.26.7 and a contained Go 1.18 SDK. Classify absent
tests, lack of terminal/platform applicability, build tags, flaky tests,
release gaps, and checksum differences precisely rather than treating a
successful compile as complete behavioral evidence.

Project every qualified exact selection in a disposable scratch tree and run
repository verify, build, count-1/count-10/race/vet, Windows-amd64 build,
pinned lint, byte-identical root/status/upgrade/build help, API/CLI
compatibility and reports, and empty-HOME count-2 before deciding whether an
implementation is permissible. Compare selected and candidate primary
vulnerability results at module, package, symbol, and reachable-trace levels.

Reject or retain if canonical identity, release qualification, complete
closure floor, dependency tests, platform behavior, API compatibility,
historical consumers, projection, or any quality contract fails. Do not
upgrade to an untagged or unreleased commit merely because its timestamp is
newer, and do not remove an unchanged selection as a side effect of tidy.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, XXHash implementation commit identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
`docs/plan/quality-handover.md`, `docs/plan/quality-upgrade.md`, `go.mod`,
`go.sum`, the answered Readline, Logex, XXHash, OpenCensus Proto, Crypt, and
Speakeasy archives, the earlier dependency archives named in the handover, and
every referenced quality, compatibility, release, runner, evidence, and
lifecycle contract. Earlier outcomes are final.

# Three Moves

If and only if a higher exact stable Fnmatch release is qualified and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/danwakefield/fnmatch@<qualified-version>` for one
dependency-only commit. Do not hand-edit metadata or use tidy as
implementation. Preserve every retained version, Go 1.18, toolchain Go 1.26.7,
production source, quality apparatus, and release input. Stop instead of
applying an unexplained multi-selection move. The incoming survey finds no
tagged or listed release, so any changed selection requires new primary
release evidence, not repository metadata activity or an untagged commit.

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

After the Fnmatch decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
