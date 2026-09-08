# Agent Session: Evaluate Cespare XXHash V2 Dependency

Status: NEXT
Session ID: `2026-09-08T074422+0200-evaluate-cespare-xxhash-v2-dependency`
Created: `2026-09-08T07:44:22+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `fa548066df8974bbc587c8fd96d27120a2b1361ed9375cf3a98e6577f4e44987`
Previous: [2026-09-08T060123+0200-evaluate-census-opencensus-proto-dependency.md](2026-09-08T060123+0200-evaluate-census-opencensus-proto-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/cespare/xxhash/v2 v2.1.2` as one bounded dependency group.
Resolve canonical latest, authoritative source identity, release qualification,
complete Go-floor closure, native and historical-consumer behavior, exact MVS
effects, and all applicable quality contracts. Make an exact dependency
selection only if a higher exact stable version changes the selection,
preserves the retained Go 1.18 floor through the complete minimal closure, and
passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, the accepted
Speakeasy v0.2.0 dependency move, and retained Crypt and OpenCensus Proto
selections. All earlier acceptances, rejections, no-change decisions, evidence
corrections, and lifecycle ancestry are final. Do not revisit OpenCensus Proto,
Crypt, or Speakeasy and do not combine another module group. P8 remains queued.

Project MVS selects XXHash v2.1.2 through historical requirements from
`github.com/spf13/viper v1.10.1` and
`github.com/sagikazarmark/crypt v0.4.0`. It is not an explicit main-module
requirement, no XXHash package is loaded by Ply, and Ply has no direct import.
A minimal post-OpenCensus survey finds exact proxy versions v2.0.0, v2.1.0,
v2.1.1, v2.1.2, v2.2.0, and v2.3.0. Exact `@latest` is v2.3.0 at
2024-04-04T20:00:10Z. Selected and latest both declare exact module path
`github.com/cespare/xxhash/v2` and Go 1.11.

The public canonical repository is currently enabled, unarchived, non-fork,
and has default branch `main`. Selected v2.1.2 is annotated tag object
`7ae26c41ed6fb1f8a6c21e05eeff4d91b5e401c4` dereferencing commit
`e7a6b52374f7e2abfb8abb27249d53a1997b09a7`. Latest v2.3.0 is annotated tag
object `7438b35f14d771ee32d8bbcd9527d32a336e7dad` dereferencing commit
`998dce232f17418a7a5721ecf87ca714025a3243`. Treat these as survey facts to
verify independently; establish exact trees, parents, times, signatures,
release objects, default-branch ancestry, redirects, forks, alternate module
paths, deprecation, retractions, and any post-release commits from primary
evidence.

# Measurements At Start

The latest dependency implementation remains Speakeasy commit
`41f9561f6ea2f5b6395c5c4d9bcc56a54533133a`, exact parent
`1f55aaa31280a4610ed66c1c37ddda953bfa6a8e`, tree
`a50d4f256fe742e472f1f7e9cf0589a596e971b0`, changing only `go.mod` and
`go.sum` with three insertions. Crypt and OpenCensus Proto were retained
without dependency edits.

Accepted project measurements remain 234 selected modules, 3,582 graph edges,
429 native complete-test packages, 41 loaded modules, 197 loaded packages,
1,047 `go.sum` lines, and a 371-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, accepted metadata adds exactly 31 checksum
lines. The main module retains Go 1.18 and toolchain Go 1.26.7. Ordinary and
ignored status must be empty.

Fresh primary vulnerability data contains 1,392 module records and no exact
OpenCensus Proto record. Accepted comparisons retain 20 IDs/22 reachable
traces for Darwin and Windows symbol scans and 30 Darwin module IDs.
OpenCensus Proto evidence has 459,623 verified entries; manifest SHA-256 is
`55c7090ab788c633f20666ba9d70e8f4d5f4bcb446b9e5ef8f530cd78e3b074a`.
Decision-summary SHA-256 is
`f7f4ed0dd015d7d7da58141a48f63537b8552c34761c05bc90de5b3efa22564a`.
An initial scratch-only API/CLI warm used prohibited `go mod download all`;
it touched no measured worktree and none of its results is retained as proof.
Fresh exact-get/graph/package-list cache warming and the fully offline replay
are the authoritative compatibility evidence.

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
XXHash v2 version through the Go proxy and checksum database, authoritative
repository, go-import metadata, and primary Go vulnerability data. Record
selected/candidate commits and times, module declarations and requirements,
checksum pairs, source identity, tag and commit signatures, release/default-
branch history, repository status, deprecation, and retractions. Distinguish
stable releases, prereleases, redirects, forks, alternate module paths,
unreleased branch heads, and the incompatible v1 module path.

Prove canonical latest and the highest qualified exact stable candidate that
preserves Go 1.18 through the complete minimal module and package/test closure.
Do not infer the closure floor from the root's Go 1.11 declaration. Inspect
architecture-specific implementations, pure-Go fallbacks, assembly, generated
files, and build tags. Prove behavior on supported Darwin-arm64 and
Windows-amd64 paths and identify any other release-relevant architecture
surface without expanding into another dependency group.

Measure old and candidate selected modules, complete graph edges, native
complete-test packages, loaded modules and packages, checksum lines, exact
dependency paths, exact-get effects, and the unapplied tidy projection.
Attribute every selection, edge, and checksum difference. If exact selected-
version get changes no selected version, do not add a redundant requirement or
checksum.

Identify the exact historical Viper/Crypt consumer chain and the XXHash
packages, constructors, digest methods, checksum helpers, interfaces, and
architecture paths it actually exercises. Use focused consumer and independent
known-vector/streaming fixtures where meaningful without inventing a Ply
runtime path. Keep dependency-native generator/tool closure, historical
consumer closure, and Ply's unloaded project behavior distinct.

Run source verification, package listing, native complete tests, two
independent repeated-test passes, race where supported, vet, and relevant
cross-builds under exact Go 1.26.7 and a contained Go 1.18 SDK. Treat assembly,
missing tests, build tags, vet findings, flaky tests, release gaps, or checksum
differences precisely; determine whether each is a release disqualifier under
the retained contracts rather than silently waiving it.

Project any qualified exact selection in a disposable worktree and run
repository verify, build, count-1/count-10/race/vet, Windows-amd64 build,
pinned golangci-lint, byte-identical root/status/upgrade/build help, API/CLI
compatibility and reports, and empty-HOME count-2 before deciding whether an
implementation is permissible. Compare old and candidate primary
vulnerability results at module, package, symbol, and reachable-trace levels.

Reject or retain if canonical identity, release qualification, complete
closure floor, dependency tests, platform implementation, API compatibility,
historical consumers, projection, or any quality contract fails. Do not
upgrade to an unreleased branch head merely because its commit time is newer.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, Speakeasy implementation commit identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
`docs/plan/quality-handover.md`, `docs/plan/quality-upgrade.md`, `go.mod`,
`go.sum`, the answered OpenCensus Proto, Crypt, and Speakeasy archives, the
earlier dependency archives named in the handover, and every referenced
quality, compatibility, release, runner, evidence, and lifecycle contract.
Earlier outcomes are final.

# Three Moves

If and only if a higher exact stable XXHash v2 version is qualified and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/cespare/xxhash/v2@<qualified-version>` for one dependency-
only commit. Do not hand-edit metadata or use tidy as implementation. Preserve
every retained version, the Go 1.18 directive, toolchain Go 1.26.7, production
source, quality apparatus, and release input. Stop instead of applying an
unexplained multi-selection move.

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

After the XXHash v2 decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
