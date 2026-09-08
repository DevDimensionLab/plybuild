# Agent Session: Evaluate Eliukblau Pixterm Ansimage Dependency

Status: NEXT
Session ID: `2026-09-08T153002+0200-evaluate-eliukblau-pixterm-ansimage-dependency`
Created: `2026-09-08T15:30:02+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3202618e6930ca09e44599a1e5188c959a437cb059c44c824f4a289d67f1321e`
Previous: [2026-09-08T141743+0200-evaluate-disintegration-imaging-dependency.md](2026-09-08T141743+0200-evaluate-disintegration-imaging-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/eliukblau/pixterm/pkg/ansimage
v0.0.0-20191210081756-9fb6cf8c2f75` as one bounded dependency group. Resolve
canonical latest despite the empty exact proxy list and failing `@latest`,
authoritative nested-module and repository identity, complete Go-floor
closure, terminal-image behavior, actual production consumers, exact MVS
effects, and all applicable quality contracts. Make an exact dependency
selection only if a higher qualified exact stable release changes the
selection, preserves the retained Go 1.18 floor through the complete minimal
closure, and passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
Logex, Readline, Fnmatch, and Imaging selections. All earlier acceptances,
rejections, no-change decisions, evidence corrections, and lifecycle ancestry
are final. Do not revisit Imaging, Fnmatch, Readline, Logex, XXHash,
OpenCensus Proto, Crypt, or Speakeasy and do not combine another module group.
P8 remains queued.

Ply loads ansimage in production through `plybuild/cmd ->
github.com/MichaelMure/go-term-markdown ->
github.com/eliukblau/pixterm/pkg/ansimage`. Markdown opens local or HTTP image
readers, calls `NewScaledFromReader`, and renders terminal ANSI output.
Ansimage uses Imaging `Resize`, `Fit`, and `Fill`, go-colorful, and registered
x/image decoders. Keep the exact nested module's declared closure, project
MVS's selected closure, loaded production/test packages, and actual Ply
behavior distinct. Do not turn this into an Imaging, go-colorful, x/image,
terminal-color, or unrelated tidy review.

A minimal post-Imaging survey finds an empty exact proxy version list, while
the selected pseudo-version remains fetchable with checksum pair
`h1:vbix8DDQ/rfatfFr/8cf/sJfIL69i4BcZfjrVOxsMqk=` /
`h1:0gZuvTO1ikSA5LtTI6E13LEOdWQNjIo5MTQOvrV0eFg=`. Exact `@latest` returns
404 with “no matching versions”. The selected nested module declares exact
path `github.com/eliukblau/pixterm/pkg/ansimage`, Go 1.13, Imaging v1.6.2,
go-colorful v1.0.3, and x/image
v0.0.0-20191206065243-da761ea9ff43. Independently resolve every relevant
stable, prerelease, pseudo-version, proxy-absent tag, redirect, fork,
alternate-path, nested-module, and unreleased identity, including deprecation
and retractions.

The public repository `github.com/eliukblau/pixterm` is currently enabled,
unarchived, non-fork, and defaults to `master`. Selected unsigned commit
`9fb6cf8c2f75275ebcd4ac8c30a0e26930d497f7`, tree
`bdcdecab7b23ba3a6d18efd3ced9800826668031`, parent
`be34e524a7d8fbf6ab827fd06f669ad4e50943b0`, dated
2019-12-10T08:17:56Z, is also root-project stable Release/tag v1.3.0. Because
the tag is not prefixed for the nested module, the exact module resolves it as
the selected pseudo-version rather than v1.3.0.

The selected commit is the only repository commit containing
`pkg/ansimage/go.mod`. Its immediate successor
`9f095995d66abbf03a06cea8e8e9b7cfd679e06c`, tree
`ed6919823270b5ac3737c516ea8900c20bb4911e`, dated
2019-12-15T12:16:17Z, removes the nested module in favor of one repository-root
module. Root-project releases continue through v1.3.3 at current unsigned
master `24a1aedad1a99b2177808bfd72b23e27318989f6`, tree
`1e75a629c404807d2953ae70dd860d8f2c5f03ea`, dated
2026-08-01T21:30:42Z; current root go.mod declares Go 1.25.0. Treat this as a
minimal incoming survey to verify independently. Do not treat a root-project
tag, the consolidated root module, or a later directory tree as an exact-path
nested-module release without proving valid Go module identity.

# Measurements At Start

The latest dependency implementation remains exact XXHash v2.3.0 commit
`e5d6252825d7a1822c01819b9144050f345a6ad4`, exact parent
`a23ce0f60ad65aac4f4800d6c095d911c0f4e754`, tree
`5213ba55981d77d7c8061312915e237e80d29af8`, changing only `go.mod` and
`go.sum` with three insertions. Imaging, Fnmatch, Readline, Logex, Crypt, and
OpenCensus Proto remain retained without dependency edits.

Accepted project measurements are 234 selected modules, 3,583 graph edges,
429 native complete-test packages, 41 loaded modules, 197 loaded packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, current metadata adds exactly 33 checksum
lines and removes zero. The main module retains Go 1.18 and toolchain Go
1.26.7. Ordinary and ignored status must be empty.

Fresh primary vulnerability data contains 1,392 module records and no Imaging
record. Accepted current/no-op populations retain 20 IDs/22 reachable traces
for Darwin and Windows symbol scans, 22 Darwin package IDs/findings, and 30
Darwin module IDs/findings. Imaging evidence has 1,713 verified entries;
manifest SHA-256 is
`046746e0c4004d62ebac4838dac739ce37a0d4576a0fae3e5d1db987e4d47308`.
Decision-summary SHA-256 is
`507f403c1289ff6d698beffb31eea6c3a3c07835efcf5609d6bde475bb4cdc5e`.
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
ansimage version or commit through the Go proxy and checksum database,
authoritative repository, go-import metadata, upstream ancestry, nested-module
history, and primary Go vulnerability data. Record commits and times, module
declarations and requirements, checksum pairs, source identity, tag/commit
signatures, root versus nested tag semantics, release/default-branch history,
repository status, deprecation, and retractions. Distinguish exact nested
module releases from root-module releases, prereleases, pseudo-versions,
proxy-absent tags, redirects, forks, alternate paths, and unreleased heads.

Prove canonical latest when `@latest` currently fails and whether any higher
qualified exact stable candidate exists. Preserve Go 1.18 through the complete
minimal module and package/test closure; do not infer the closure floor from
only the selected Go 1.13 directive. Keep inherited Imaging, go-colorful, and
x/image MVS effects inside closure measurement without independently upgrading
those already separate dependency groups.

Inspect ansimage's exported API and reader/image constructors, scale modes,
dimension calculation, dithering and no-dithering, alpha/background handling,
ANSI/true-color/256-color output, terminal sizing, string/render methods,
malformed-input, small/zero image, resource-bound, and supported-format
behavior relevant to Ply. Inspect build tags, generated files, OS/architecture
code, Cgo/native surface, examples, testdata, fuzz/property coverage, and
release gaps. Add independent fixtures because the selected ansimage module
has no native `_test.go` files.

Preserve and characterize the Imaging-session consumer finding precisely:
selected ansimage no-dither `RenderExt` returns empty output for a
two-pixel-high scaled image, while go-term-markdown reports the image rendered
and emits only title/destination; four-pixel fixtures render normally. Decide
whether this is invariant selected behavior, a qualified-candidate change, or
a stop condition. Do not attribute it to Imaging.

Run source verification, package listing, independent complete tests, two
independent repeated-test passes, race where supported, vet, and relevant
cross-builds under exact Go 1.26.7 and a contained Go 1.18 SDK. Exercise the
exact go-term-markdown -> ansimage production path for local/readers and all
registered formats. Distinguish executable behavior from ansimage's absent
native suite and from go-term-markdown's PTY/NO_COLOR-sensitive ANSI goldens
and unrelated project-MVS Chroma golden drift.

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
identity, release qualification, complete closure floor, dependency behavior,
small-image handling, platform/terminal behavior, API compatibility,
production consumers, projection, or any quality contract fails. Do not
upgrade to the consolidated root module, later root tag, or unreleased commit
merely because it is newer, and do not remove the unchanged nested selection
as a side effect of tidy.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, XXHash implementation commit identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
`docs/plan/quality-handover.md`, `docs/plan/quality-upgrade.md`, `go.mod`,
`go.sum`, the answered Imaging, Fnmatch, Readline, Logex, XXHash, OpenCensus
Proto, Crypt, and Speakeasy archives, the earlier dependency archives named in
the handover, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

If and only if a higher exact stable ansimage release is qualified and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/eliukblau/pixterm/pkg/ansimage@<qualified-version>` for one
dependency-only commit. Do not hand-edit metadata, migrate import paths, or
use tidy as implementation. Preserve every retained version, Go 1.18,
toolchain Go 1.26.7, production source, quality apparatus, and release input.
Stop instead of applying an unexplained multi-selection or module-path move.

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

After the ansimage decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
