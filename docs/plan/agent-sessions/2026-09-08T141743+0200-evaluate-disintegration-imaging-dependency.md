# Agent Session: Evaluate Disintegration Imaging Dependency

Status: NEXT
Session ID: `2026-09-08T141743+0200-evaluate-disintegration-imaging-dependency`
Created: `2026-09-08T14:17:43+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `6016d94ee09144381bfa8ecea3bd6221aab2ea1fcfca90a6d352d6d5c1c25254`
Previous: [2026-09-08T132110+0200-evaluate-danwakefield-fnmatch-dependency.md](2026-09-08T132110+0200-evaluate-danwakefield-fnmatch-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/disintegration/imaging v1.6.2` as one bounded dependency group.
Resolve canonical latest, authoritative source and release identity, complete
Go-floor closure, native image behavior, actual production consumers, exact
MVS effects, and all applicable quality contracts. Make an exact dependency
selection only if a higher qualified exact stable release changes the
selection, preserves the retained Go 1.18 floor through the complete minimal
closure, and passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
Logex, Readline, and Fnmatch selections. All earlier acceptances, rejections,
no-change decisions, evidence corrections, and lifecycle ancestry are final.
Do not revisit Fnmatch, Readline, Logex, XXHash, OpenCensus Proto, Crypt, or
Speakeasy and do not combine another module group. P8 remains queued.

Ply loads Imaging in production through `plybuild/cmd ->
github.com/MichaelMure/go-term-markdown ->
github.com/eliukblau/pixterm/pkg/ansimage -> Imaging`. The graph also retains
the main module's indirect Imaging requirement. Selected Imaging v1.6.2
requires `golang.org/x/image
v0.0.0-20191009234506-e7c1f5e7dbb8`, while project MVS selects x/image v0.5.0.
Keep Imaging's declared edge, MVS's selected closure, loaded production/test
packages, and actual Ply behavior distinct. Do not turn this into an x/image
review or unrelated tidy cleanup.

A minimal post-Fnmatch survey finds 15 exact proxy versions from v1.0.0
through v1.6.2. Exact `@latest` is the selected stable v1.6.2 at
2019-11-16T20:43:25Z. Its module file declares exact path
`github.com/disintegration/imaging`, no Go directive, and the old x/image
pseudo-version above. Independently resolve every relevant stable,
prerelease, pseudo-version, redirected, forked, alternate-path, and unreleased
identity, including deprecation and retractions.

The public repository is currently enabled, unarchived, non-fork, and defaults
to `master`. Selected v1.6.2 is a stable GitHub Release and lightweight tag at
verified commit `acabd8315e63bfcaac97d52d68a7a0b88d2eea93`, tree
`6584cbb2a26e4d38bfec8f2587633f234500810e`, with parents
`9aab30e6aa535fe3337b489b76759ef97dfaf362` and
`675e3c209ff3e9bbee22db0bffe990d3abace4ce`, dated
2019-11-16T20:43:25Z. Master is later verified commit
`d40f48ce0f098c53ab1fcd6e0e402da682262da5`, tree
`cfae2d9af62546482831388535eb24023ceac298`, dated
2020-12-18T19:30:11Z, and resolves only as unreleased
`v1.6.3-0.20201218193011-d40f48ce0f09`. Treat these as a minimal incoming
survey to verify independently; do not select an unreleased branch head.

# Measurements At Start

The latest dependency implementation remains exact XXHash v2.3.0 commit
`e5d6252825d7a1822c01819b9144050f345a6ad4`, exact parent
`a23ce0f60ad65aac4f4800d6c095d911c0f4e754`, tree
`5213ba55981d77d7c8061312915e237e80d29af8`, changing only `go.mod` and
`go.sum` with three insertions. Fnmatch, Readline, Logex, Crypt, and OpenCensus
Proto remain retained without dependency edits.

Accepted project measurements are 234 selected modules, 3,583 graph edges,
429 native complete-test packages, 41 loaded modules, 197 loaded packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, current metadata adds exactly 33 checksum
lines and removes zero. The main module retains Go 1.18 and toolchain Go
1.26.7. Ordinary and ignored status must be empty.

Fresh primary vulnerability data contains 1,392 module records and no Fnmatch
record. Accepted current/no-op populations retain 20 IDs/22 reachable traces
for Darwin and Windows symbol scans, 22 Darwin package IDs/findings, and 30
Darwin module IDs/findings. Fnmatch evidence has 72 verified entries; manifest
SHA-256 is
`97130774da6c5f883f1c8fd0afa8190490be6aa3f8c998af4deaa9b7c79d44f1`.
Decision-summary SHA-256 is
`96f4a6d499ac460e99d3c4b7b9395b2ed826374d79fca4602e89c359b676eae1`.
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
Imaging version through the Go proxy and checksum database, authoritative
repository, go-import metadata, upstream ancestry, and primary Go
vulnerability data. Record commits and times, module declarations and
requirements, checksum pairs, source identity, tag/commit signatures,
release/default-branch history, repository status, deprecation, and
retractions. Distinguish stable releases, prereleases, pseudo-versions,
proxy-absent tags, redirects, forks, alternate module paths, and unreleased
branch heads.

Prove canonical latest and whether any qualified higher exact stable candidate
exists. Preserve Go 1.18 through the complete minimal module and package/test
closure; do not infer the closure floor from Imaging's module file merely
because it lacks a Go directive. Keep inherited x/image selection effects
inside closure measurement without selecting or independently upgrading that
next dependency group.

Inspect Imaging's exported API and the decoding, orientation, resizing,
resampling, cropping, compositing, adjustment, convolution, encoding, color,
alpha, malformed-input, resource-bound, and image-format behavior relevant to
the selected release and Ply's terminal-image consumer. Inspect build tags,
generated files, OS/architecture-specific code, native/Cgo surface, examples,
testdata, fuzz/property coverage, and release gaps. Add independent fixtures
where the native suite leaves release-relevant behavior weak.

Run source verification, package listing, native complete tests, two
independent repeated-test passes, race where supported, vet, and relevant
cross-builds under exact Go 1.26.7 and a contained Go 1.18 SDK. Exercise the
exact go-term-markdown -> pixterm/ansimage production path and distinguish
Imaging production behavior from dependency tests. Classify absent, flaky,
platform-limited, golden-image, or toolchain-sensitive tests precisely rather
than treating compilation as behavioral evidence.

Measure selected modules, graph edges, native complete-test packages, loaded
modules/packages, checksum lines, exact dependency paths, exact selected-get
or candidate-get effects, and the unapplied tidy projection. Attribute every
difference. Project every qualified exact selection in a disposable scratch
tree and run repository verify, build, count-1/count-10/race/vet,
Windows-amd64 build, pinned lint, byte-identical root/status/upgrade/build
help, API/CLI compatibility and reports, and empty-HOME count-2 before
deciding whether implementation is permissible.

Compare selected and candidate primary vulnerability results at module,
package, symbol, and reachable-trace levels. Reject or retain if canonical
identity, release qualification, complete closure floor, dependency tests,
image/platform behavior, API compatibility, production consumers, projection,
or any quality contract fails. Do not upgrade to unreleased master because its
timestamp is newer, and do not remove an unchanged selection as a side effect
of tidy.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, XXHash implementation commit identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
`docs/plan/quality-handover.md`, `docs/plan/quality-upgrade.md`, `go.mod`,
`go.sum`, the answered Fnmatch, Readline, Logex, XXHash, OpenCensus Proto,
Crypt, and Speakeasy archives, the earlier dependency archives named in the
handover, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

If and only if a higher exact stable Imaging release is qualified and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/disintegration/imaging@<qualified-version>` for one
dependency-only commit. Do not hand-edit metadata or use tidy as
implementation. Preserve every retained version, Go 1.18, toolchain Go 1.26.7,
production source, quality apparatus, and release input. Stop instead of
applying an unexplained multi-selection move. The incoming survey finds
selected already stable latest, so any changed selection requires new primary
stable-release evidence, not unreleased master.

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

After the Imaging decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
