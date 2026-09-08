# Agent Session: Evaluate Fatih Color Dependency

Status: NEXT
Session ID: `2026-09-08T170736+0200-evaluate-fatih-color-dependency`
Created: `2026-09-08T17:07:36+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5ebbbb5df7479ead6003c4292090ba95ee8d821c6df6efcf6795a0a6c3ce35a3`
Previous: [2026-09-08T153002+0200-evaluate-eliukblau-pixterm-ansimage-dependency.md](2026-09-08T153002+0200-evaluate-eliukblau-pixterm-ansimage-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/fatih/color v1.14.1` as one bounded dependency group. Resolve every
relevant exact release through current stable v1.19.0, authoritative source
and release identity, complete Go-floor closure, terminal-color behavior,
actual production consumers, exact MVS effects, and all applicable quality
contracts. Make an exact dependency selection only if the highest qualified
stable release preserves the retained Go 1.18 floor through the complete
minimal closure and passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0 and XXHash v2.3.0 moves, and retained Crypt, OpenCensus Proto,
Logex, Readline, Fnmatch, Imaging, and ansimage selections. All earlier
acceptances, rejections, no-change decisions, evidence corrections, and
lifecycle ancestry are final. Do not revisit ansimage, Imaging, Fnmatch,
Readline, Logex, XXHash, OpenCensus Proto, Crypt, or Speakeasy and do not
combine another module group. P8 remains queued.

Ply loads Fatih Color in production through `plybuild/cmd ->
github.com/MichaelMure/go-term-markdown -> github.com/fatih/color`.
Go-term-markdown creates Sprint functions for green, high-green, bold-green,
blue, blue-background italic, and red; it also branches on the package-global
`color.NoColor`. Image destinations and malformed-image fallbacks use Blue.
Keep Fatih Color's declared closure, project MVS's already-selected terminal
support modules, loaded production/test packages, and actual Ply behavior
distinct. Do not turn this into a go-colorable, go-isatty, x/sys, ansimage,
terminal emulator, or unrelated tidy review.

A minimal post-ansimage survey finds 22 exact proxy versions. Selected
v1.14.1 is a stable GitHub Release and declares Go 1.17 with go-colorable
v0.1.13, go-isatty v0.0.17, and x/sys v0.3.0. Its checksum pair is
`h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=` /
`h1:2oHN61fhTpgcxD3TSWCgKDiH1+x4OiDVVGH8WlgGZGg=`. Selected tag commit is
`3d5097c6b003cf3a784e670ddb79710cf46e9a07`.

Stable v1.15.0, v1.16.0, v1.17.0, and v1.18.0 all declare Go 1.17. V1.18.0
requires go-colorable v0.1.13, go-isatty v0.0.20, and x/sys v0.25.0; its tag
commit is `1c8d8706604ee5fb9a464e5097ba113101828a75` and checksum pair is
`h1:S8gINlzdQ840/4pfAwic/ZE0djQEH3wM94VfqLTZcOM=` /
`h1:4FelSpRwEGDpQ12mAdzqdOukCy4u8WUtOY6lkT/6HfU=`. Treat v1.18.0 as the
highest immediately visible floor-compatible candidate, not as prequalified.

Exact `@latest` is stable v1.19.0, released 2026-03-20, at tag commit
`ca25f6e17f118a5a259f3c2c0d395949d1103a5a`. It declares Go 1.25.0 and
requires go-colorable v0.1.14, go-isatty v0.0.20, and x/sys v0.42.0. Its
checksum pair is `h1:Zp3PiM21/9Ld6FzSKyL5c/BULoe/ONr9KlbYVOfG8+w=` /
`h1:zNk67I0ZUT1bEGsSGyCZYZNrHuTkJJB+r6Q9VuMi0LE=`. V1.19.0 does not
preserve the retained Go 1.18 floor and must not be selected. Still resolve
its identity, change boundary, and release facts accurately.

The public upstream `github.com/fatih/color` is currently enabled, unarchived,
non-fork, and defaults to `main`; stable GitHub Releases continue through
v1.19.0. Independently resolve stable, prerelease, pseudo-version,
proxy-absent tag, redirect, fork, alternate-path, and unreleased identities,
including deprecation, retractions, tag/commit signatures, commit times,
trees, parents, repository status, and default-branch ancestry. Treat the
incoming details only as a minimal survey to verify.

Project MVS already selects go-colorable v0.1.15, go-isatty v0.0.20, and
x/sys v0.30.0. These accepted selections are above or equal to v1.18.0's
declared requirements and below v1.19.0's x/sys requirement. Prove exact old,
v1.18.0, and ineligible v1.19.0 MVS effects in disposable trees without
combining their dependency groups. A candidate that changes anything beyond
the explained exact closure is a stop condition unless it remains strictly
inside this group's authorized MVS projection.

# Measurements At Start

The latest dependency implementation remains exact XXHash v2.3.0 commit
`e5d6252825d7a1822c01819b9144050f345a6ad4`, exact parent
`a23ce0f60ad65aac4f4800d6c095d911c0f4e754`, tree
`5213ba55981d77d7c8061312915e237e80d29af8`, changing only `go.mod` and
`go.sum` with three insertions. Ansimage, Imaging, Fnmatch, Readline, Logex,
Crypt, and OpenCensus Proto remain retained without dependency edits.

Accepted project measurements are 234 selected modules, 3,583 graph edges,
429 native complete-test packages, 41 loaded modules, 197 loaded packages,
1,049 `go.sum` lines, and a 381-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, current metadata adds exactly 33 checksum
lines and removes zero. The main module retains Go 1.18 and toolchain Go
1.26.7. Ordinary and ignored status must be empty.

Fresh primary vulnerability data contains 1,392 module records and no Fatih
Color or ansimage record. Accepted current/no-op populations retain 20 IDs/22
reachable traces for Darwin and Windows symbol scans, 22 Darwin package
IDs/findings, and 30 Darwin module IDs/findings. Ansimage appears only as a
call frame in 12 existing traces across 11 x/image IDs; do not misattribute
those inherited findings to Fatih Color.

Ansimage evidence has 73 verified entries; manifest SHA-256 is
`95d6ca0ac0d767c94f78b3a0a5c32f1008cd2e7b05a3ffc5f639b9af021dd112`.
Decision-summary SHA-256 is
`3f39b1aa5f16d3a1b398ce329e4a73621931b09d3b191dd7c0bbadc5743de0ed`.
The authoritative Q0-Q2 scorecard SHA-256 remains
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
Fatih Color release through the Go proxy and checksum database, authoritative
repository, go-import metadata, upstream ancestry, and primary Go
vulnerability data. Record version order, module declarations and full
requirements, checksum pairs, source identity, release/tag/commit signatures,
commits, times, trees, parents, default-branch history, repository status,
deprecation, and retractions. Distinguish stable releases, prereleases,
pseudo-versions, proxy-absent tags, redirects, forks, alternate module paths,
and unreleased heads.

Prove canonical latest and the highest qualified stable release that can
preserve Go 1.18 through the complete minimal module and package/test closure.
Do not infer the closure floor from Fatih Color's Go directive alone. Test
v1.14.1 and every serious floor-compatible candidate with a contained Go 1.18
SDK. Treat v1.19.0's Go 1.25 directive as an explicit incompatibility, not a
reason to move the main module's retained floor.

Inspect the exported API and semantics for attributes, colors, formatting,
Sprint/Sprintf/Sprintln, Print/Fprint variants, output writers, enable/disable
methods, format validation, `NoColor`, `Output`, `Error`, unrecognized
parameters, empty/newline input, nested formatting, concurrency, and global
state. Exercise TTY and non-TTY output plus `NO_COLOR`, `TERM=dumb`,
`CLICOLOR`, and `CLICOLOR_FORCE` where supported by each exact release.
Characterize Windows terminal support and all OS/build-tag branches without
reviewing go-colorable, go-isatty, or x/sys as independent groups.

Inspect build tags, generated files, OS/architecture code, Cgo/native surface,
examples, testdata, fuzz/property coverage, native suite gaps, and release
changes from v1.14.1 through v1.18.0. Add independent fixtures for
release-relevant output bytes, writer routing, terminal detection, environment
precedence, concurrent use, malformed format input, and Markdown's actual
consumer behavior.

Run source verification, package listing, native complete tests, two
independent repeated-test passes, race where supported, vet, and relevant
cross-builds under exact Go 1.26.7 and contained Go 1.18. Exercise the exact
go-term-markdown production path and distinguish Fatih Color behavior from
Markdown's PTY/NO_COLOR-sensitive ANSI goldens, ansimage's retained two-pixel
gap, and unrelated project-MVS Chroma output.

Measure selected modules, graph edges, native complete-test packages, loaded
modules/packages, checksum lines, exact dependency paths, old/candidate get
effects, and the unapplied tidy projection. Attribute every difference.
Project every qualified exact stable selection in a disposable scratch tree
and run repository verify, build, count-1/count-10/race/vet, Windows-amd64
build, pinned lint, byte-identical root/status/upgrade/build help, API/CLI
compatibility and reports, and empty-HOME count-2 before deciding whether
implementation is permissible.

Compare selected and candidate primary vulnerability results at module,
package, symbol, and reachable-trace levels. Reject or retain if canonical
identity, release qualification, complete closure floor, terminal behavior,
global-state/race behavior, native tests, platform support, API compatibility,
production consumers, MVS projection, or any quality contract fails. Do not
select Go-1.25-requiring v1.19.0, independently upgrade its support modules, or
remove an unchanged selection as a side effect of tidy.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, XXHash implementation commit identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
`docs/plan/quality-handover.md`, `docs/plan/quality-upgrade.md`, `go.mod`,
`go.sum`, the answered ansimage, Imaging, Fnmatch, Readline, Logex, XXHash,
OpenCensus Proto, Crypt, and Speakeasy archives, the earlier terminal-support
dependency archives named in the handover, and every referenced quality,
compatibility, release, runner, evidence, and lifecycle contract. Earlier
outcomes are final.

# Three Moves

Select only the highest exact stable Fatih Color release that independently
qualifies and whose complete minimal closure preserves Go 1.18. The incoming
survey makes v1.18.0 the highest plausible candidate and disqualifies v1.19.0
on its Go 1.25.0 directive. If v1.18.0 qualifies, use exact Go 1.26.7 and
exact `go get github.com/fatih/color@v1.18.0` for one dependency-only commit.
Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, Go 1.18, toolchain Go 1.26.7, production source, quality
apparatus, and release input. Stop instead of applying an unexplained
multi-selection move.

After a changed selection, run the complete P7 dependency gate: dependency
and consumer tests; graph/path/checksum/tidy proof; repository
verify/build/tests/race/vet/Windows/pinned lint; help/API/CLI; launcher and Make
contracts; preflight; host plus fresh snapshot/Docker meta and acceptance;
audit meta; focused and exact Q0-Q2 audits; separate full audit; vulnerability
comparison; empty-HOME count-2; and final ordinary/ignored cleanliness. Exact
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

After the Fatih Color decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
