# Agent Session: Evaluate Fsnotify Fsnotify Dependency

Status: NEXT
Session ID: `2026-09-08T210923+0200-evaluate-fsnotify-fsnotify-dependency`
Created: `2026-09-08T21:09:23+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `b0aab51806035e0fb054a1a152cc8d93e8f981d358f073f9a32ee9f7e87a33e2`
Previous: [2026-09-08T170736+0200-evaluate-fatih-color-dependency.md](2026-09-08T170736+0200-evaluate-fatih-color-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/fsnotify/fsnotify v1.6.0` as one bounded dependency group. Resolve
every relevant release through stable v1.10.1, authoritative source and
release identity, complete Go-floor closure, operating-system watcher
behavior, actual Viper/Ply consumers, exact MVS effects, and all applicable
quality contracts. Select only the highest exact stable release whose complete
minimal closure preserves Go 1.18 and whose loaded behavior passes every
contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, and Fatih Color v1.15.0 moves. Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, and ansimage remain
retained. All earlier decisions and lifecycle ancestry are final. Do not
revisit them or combine another dependency group. P8 remains queued.

Ply loads Fsnotify in production through `plybuild/cmd -> github.com/spf13/viper
-> github.com/fsnotify/fsnotify`. Keep Fsnotify's declared closure, project
MVS's x/sys selection, Viper's configuration watching, loaded production/test
packages, and Ply's actual CLI behavior distinct. Do not turn this into an
x/sys, Viper, file-system, polling, container-runtime, or unrelated tidy
review.

A minimal post-Fatih survey finds 38 exact proxy versions. Selected v1.6.0
declares Go 1.16, requires x/sys
`v0.0.0-20220908164124-27713097b956`, and retracts accidental v1.5.3 plus
symlink-regressing v1.5.0. Its checksum pair is
`h1:n+5WquG0fcWoWp6xPWfHdbskMCQaFnG6PfBrh1Ky4HY=` /
`h1:sl3t1tCWJFWoRz9R8WJCbQihKKwmorjAbSClcnxKAGw=` at commit
`5f8c606accbcc6913853fe7e083ee461d181d88d`.

Stable v1.7.0, v1.8.0, and v1.9.0 declare Go 1.17. V1.9.0 requires x/sys
v0.13.0 and has checksum pair
`h1:2Ml+OJNzbYCTzsxtv8vKSFD9PbJjmhYF14k/jKC7S9k=` /
`h1:8jBTzvmWwFyi3Pb8djgCCO5IBqzKJ/Jwo8TRcHyHii0=` at commit
`ae0e7923765f64fb8061396db7edebb558cf6093`. Treat v1.9.0 as the
highest immediately visible floor-compatible candidate, not as prequalified.

Exact `@latest` is stable v1.10.1, released 2026-05-04 at commit
`76b01a6e8f502187fecedea8b025e79e5a86085c` and checksum pair
`h1:b0/UzAf9yR5rhf3RPm9gf3ehBPpf0oZKIjtpKrx59Ho=` /
`h1:TLheqan6HD6GBK6PrDWyDPBaEV8LspOxvPSjC+bVfgo=`. V1.10.0 and v1.10.1
declare Go 1.23 and cannot preserve the retained Go 1.18 floor. Do not select
them or move the main Go line.

Project MVS already selects x/sys v0.30.0, above v1.9.0's declared v0.13.0.
Prove old, v1.9.0, and ineligible v1.10.1 effects in disposable trees. Any
candidate change outside the explained exact Fsnotify edge and its authorized
MVS projection is a stop condition. Do not independently upgrade or remove
x/sys.

# Measurements At Start

The latest dependency implementation is Fatih Color v1.15.0 commit
`6ca672ef38688b7f6f505cf0cb273d07c4c2ba9a`, parent
`d181fcd6c11fa147e0b44dd007598872875fc9e6`, and tree
`9ba2fdc442553622028a4a8464536915d345510c`, changing only `go.mod` and
`go.sum` with three insertions and one deletion.

Accepted project measurements are 234 selected modules, 3,583 graph edges,
429 native complete-test packages, 41 loaded modules, 197 loaded packages,
1,051 `go.sum` lines, and a 383-line unapplied tidy projection. Relative to
accepted go-cmp commit `c314bcb`, metadata adds exactly 35 checksum lines and
removes zero. The main module retains Go 1.18 and toolchain Go 1.26.7.

Fresh primary vulnerability data contains 1,392 module records. Accepted
populations are 20 IDs/22 reachable traces for Darwin and Windows symbol
scans, 22 Darwin package findings, and 30 Darwin module findings. Fatih Color
has no record or trace. Do not attribute existing x/image/ansimage findings to
Fsnotify.

Fatih Color exact `make quality` records all 27 Q0-Q2 rows PASS at L2 with
scorecard SHA-256
`dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.
Its decision-summary SHA-256 is
`93a6e698e48fbcec2fc97d487de8a9a56c44a3ac2530feb28199c5c5d90d54f9`.
Read the answered archive and rolling handover for the sealed evidence
manifest.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate pinned tools beneath scratch as needed. Portable
receipts remain golangci-lint 2.12.2
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

# Role And Boundaries

From fresh external archives and caches, resolve every relevant exact-path
Fsnotify release through the Go proxy, checksum database, go-import metadata,
authoritative repository, release/tag ancestry, and primary Go vulnerability
data. Record version ordering, declarations and full requirements, checksum
pairs, commits, times, trees, parents, release/tag/commit signatures,
default-branch ancestry, repository status, deprecation, and retractions.
Distinguish stable, prerelease, pseudo-version, proxy-absent tag, redirect,
fork, alternate-path, and unreleased identities.

Prove the complete minimal module and package/test closure under contained Go
1.18 for selected v1.6.0 and every serious floor-compatible candidate. Do not
infer the floor from the root directive alone. Treat v1.10.x's Go 1.23
directive as explicit incompatibility.

Inspect the exported API and watcher semantics: constructors, Add/AddWith,
Remove, WatchList, Close, Events and Errors channels, event operations and
string formatting, duplicate paths, missing paths, files versus directories,
symlinks, renames, removes, chmod, creates, writes, recursive and non-recursive
behavior, path cleanup, descriptor limits, queue overflow, close ordering,
concurrent use, blocked consumers, and error propagation. Exercise Viper's
actual WatchConfig/OnConfigChange behavior and Ply's loaded initialization.

Inspect every OS/build-tag implementation for Linux/inotify, Darwin and BSD
kqueue, Windows, Solaris/illumos, and unsupported targets. Keep kernel and
platform differences explicit. Inspect Cgo/native surface, generated files,
examples, testdata, fuzz/property coverage, native suite gaps, CI coverage,
and release changes from v1.6.0 through v1.9.0. Add independent fixtures for
release-relevant events, lifecycle, concurrency, consumer behavior, and
platform-build coverage without pretending cross-compilation executes a
foreign kernel watcher.

Run source verification, package listing, native complete tests, two
independent repeated-test passes, race where supported, vet, and relevant
cross-builds under exact Go 1.26.7 and contained Go 1.18. Characterize flaky
or timing-sensitive native tests rather than hiding them with retries.

Measure selected modules, graph edges, native complete-test packages, loaded
modules/packages, checksum lines, exact dependency paths, old/candidate get
effects, and unapplied tidy projections. Project each qualified stable
selection in a disposable tree and run repository verify, build,
count-1/count-10/race/vet, Windows-amd64 build, pinned lint, byte-identical
root/status/upgrade/build help, API/CLI compatibility, and empty-HOME count-2
before implementation.

Compare old and candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, watcher semantics, concurrency,
native tests, platform support, API, consumers, MVS, or any quality contract
fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, Fatih implementation identity, reciprocal archive history,
P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive, the rolling
handover and roadmap, `go.mod`, `go.sum`, the answered Fatih, ansimage,
Imaging, Fnmatch, Readline, Logex, XXHash, OpenCensus Proto, Crypt, and
Speakeasy archives, and every referenced quality, compatibility, release,
runner, evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

Select only the highest exact stable Fsnotify release that independently
qualifies and whose complete minimal closure preserves Go 1.18. The initial
survey makes v1.9.0 the highest plausible candidate and disqualifies v1.10.x
on their Go 1.23 directives. If v1.9.0 qualifies, use exact Go 1.26.7 and
exact `go get github.com/fsnotify/fsnotify@v1.9.0` for one dependency-only
commit. Do not hand-edit metadata or use tidy as implementation.

After a changed selection, run the complete P7 dependency gate: dependency
and Viper/Ply consumer tests; graph/path/checksum/tidy proof; repository
verify/build/tests/race/vet/Windows/pinned lint; help/API/CLI; launcher and
Make contracts; preflight; host plus fresh snapshot/Docker meta and acceptance;
audit meta; focused and exact Q0-Q2 audits; separate full audit; vulnerability
comparison; empty-HOME count-2; and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 Q0-Q2 rows at L2 and zero held,
regressed, not-comparable, or dirty counts. Full audit may exit 1 only for the
established queued L3 rows, never 2.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
create direct `/private/tmp/ply-*` roots, never bypass scratch cleanup, and
never run `go mod download all` in a measured worktree. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Fsnotify decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
