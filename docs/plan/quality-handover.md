# Quality Upgrade Handover

Generated: 2026-09-09T01:50:29+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base master at `5635d50`.
- The latest dependency implementation remains exact Fatih Color v1.15.0
  commit `6ca672ef38688b7f6f505cf0cb273d07c4c2ba9a`, exact parent
  `d181fcd6c11fa147e0b44dd007598872875fc9e6`, and tree
  `9ba2fdc442553622028a4a8464536915d345510c`. It changes only `go.mod` and
  `go.sum`, with three insertions and one deletion.
- Root GLFW, Ghodss YAML, and Fsnotify were retained without dependency
  implementation commits. Exact XXHash v2.3.0 implementation
  `e5d6252825d7a1822c01819b9144050f345a6ad4` and Speakeasy v0.2.0
  implementation `41f9561f6ea2f5b6395c5c4d9bcc56a54533133a` remain ancestors.
  Ansimage, Imaging, Fnmatch, Readline, Logex, OpenCensus Proto, and Crypt
  remain retained without dependency edits.
- The answered root-GLFW archive and sole NEXT Go Logfmt archive link
  reciprocally. No `.agent-task/current.md` or repository
  `.quality/manual-evidence.json` was created. No push, merge, publication,
  release, stash, revert, successor launch, or worktree removal occurred.

## Lifecycle And Retained Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, accepted
Speakeasy v0.2.0, XXHash v2.3.0, and Fatih Color v1.15.0 moves, and retained
Crypt, OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, and root GLFW selections. P8 remains queued. Earlier outcomes and
lifecycle ancestry are final; do not reopen them or combine another module
group.

Every disposable cache, projection, archive, source, report, generated
artifact, evidence tree, and build context must remain beneath
`$CODEX_SESSION_SCRATCH_ROOT`. Never create direct `/private/tmp/ply-*`
roots, never run `go mod download all` in a measured tree, and preserve the
launcher's scratch cleanup and reciprocal archive contract.

## Root GLFW Decision

Retain exact-path `github.com/go-gl/glfw
v0.0.0-20190409004039-e6da0acd62b1` without editing dependency metadata. The
exact-root proxy list is empty and the repository has no root release tag.
Its only tags, `v3.4/glfw/v0.1.0-pre.1` and
`v3.4/glfw/v0.1.0-pre.2`, belong to a distinct nested module. There is no
root retraction, deprecation, redirect, or alternate stable identity.

Selected source/mod sums are
`h1:QbL/5oDUmRBzO9/Z7Seo6zf912W/a6Sr4Eu0G/3Jho0=` /
`h1:vR7hzQXu2zJy9AVAgeJqvqgH9Q5CA+iKCZ2gyEVpxRU=`. Proxy ZIP SHA-256
`96c694c42e7b866ea8e26dc48b612c4daa8582ce61fdeefbe92c1a4c46163169`
matches all 227 repository-backed files. The extra proxy file is the
synthesized requirement-free, directive-free root `go.mod`.

Unsigned selected commit `e6da0acd62b1b57ee2799d4d0a76a7d4514dc5bc`, tree
`17ab3d23b59cab5cffcafe1232da9c2af635a2a1`, parent
`39f94f8075907c0c6524d2791e05d625627fe268`, has committer time
`2019-04-09T00:40:39Z` and is on default `master` ancestry. Go-import resolves
the public, enabled, unarchived, non-fork
`https://github.com/go-gl/glfw.git` repository.

The last serious floor-compatible default-branch candidate is
`v0.0.0-20240506104042-037f3cc74f2a`, Go 1.12, with sum pair
`h1:FAC6eA052T8d4Lp5GR38Sxta8fnq//jjBGDtsT0TVAU=` /
`h1:wyvWpaEu9B/VQiV1jsPs7Mha9I7yto/HqIBw197ZAzk=`. GitHub verifies merge
commit `037f3cc74f2ab0b249928c4fc4b61e0f13befdb8`, tree
`d16f4bb5267aa50bdb30ea99992c3e0a2381e5b7`, parents
`a69d953ea14231b2f19998f0b434b6c1e6c6f710` and
`23bba3a4c89646aa52826de643bfde2a8e4cfe98`, at
`2024-05-06T10:40:42Z`. It fails complete Darwin/arm64 behavior.

A later Go-1.12 pseudo-version from nested-v3.4 development is not a released
or default root head and changes no root implementation. Its source/mod sums
are `h1:QuL5wPRo3euHS1jVtTKoBWdPZjCg/4AvIvmL1sqd1go=` /
`h1:wyvWpaEu9B/VQiV1jsPs7Mha9I7yto/HqIBw197ZAzk=`. Unsigned commit
`fa5a0d0838962e92f7846982429caf1716bb46b2`, tree
`e4defc2173880171b4b9ee2eea99df400a4b6aa9`, parent
`3db5e4c7a5901ef58ab50842a938a7f228139301`, is timed
`2026-02-27T15:11:18Z`. Exact latest
`v0.0.0-20260823155953-d41da22a9587` declares Go 1.19 and is floor-ineligible.
Its root packages are byte-identical to the rejected 2024 candidate. Correct
the prior survey: latest tree is
`fcbf95ef3a46826e0cbc29777b6f83d60265a57d` and parent is
`8fa725040a18feef1098d802d2108effa73ea652`. The GitHub-verified commit time
is `2026-08-23T15:59:53Z`.

## Closure, Native Boundaries, And API

The root module contains only Cgo packages `v3.0/glfw`, `v3.1/glfw`, and
`v3.2/glfw`, with no non-standard-library module dependency. Both exact Go
1.26.7 and contained Go 1.18.10 resolve the same closure. Selected has 46 Go
and 172 C/Objective-C/header files, but no `_test.go`, example, testdata,
fuzz, or property suite. The 2024 candidate's two testdata programs are not
tests. Upstream CI never covered all three root packages; current workflows
cover only nested v3.3/v3.4 modules.

The three package APIs contain 104, 133, and 150 exported declarations. The
2024 candidate changes no exported signature. V3.1/v3.2 install C-to-Go error
and event callbacks, keep global callback state, use unsafe pointers, and rely
on GLFW's OS-thread rules rather than independent locking. Callers must lock
the main OS thread before `Init`; `Terminate`, event processing, and most
window calls stay there. Contexts are calling-thread-affine, while
`PostEmptyEvent` is the explicit secondary-thread wakeup.

V3.0 requires an external GLFW 3.0 library. V3.1/v3.2 bundle native source.
Darwin links Cocoa/OpenGL/IOKit/CoreVideo; Linux uses X11/OpenGL by default and
v3.2 optionally Wayland/EGL/xkbcommon; Windows needs Win32 headers/libraries;
FreeBSD needs X11/Wayland libraries and v3.2 pkg-config GLFW. OpenBSD/NetBSD
have no supported native backend. `CGO_ENABLED=0` exposes no package on any
tested target, so there is no pure-Go/headless fallback.

With the real Xcode compiler/SDK and scratch-provisioned official GLFW 3.0.4
for v3.0, all packages pass Darwin/amd64 count-1, two count-10 repeats, and
race under both SDKs. Candidate vet passes; selected vet reports possible
`reflect.SliceHeader` misuse in v3.1/v3.2. Independent per-package value,
error, determinism, concurrency, and v3.2 lifecycle fixtures pass. Importing
all three historical packages into one binary is unsupported because their C
and exported callback symbols collide.

On Darwin/arm64, v3.0 and v3.2 pass, but default v3.1 fails selected and
candidate under both SDKs because its Cgo directives select no client-library
macro for arm64. A `gles2` compile probe is not a valid macOS runtime
substitute. Cross-target package selection passes both SDKs, but actual
Linux/Windows/FreeBSD compilation is not available in the contained Darwin
host without target C compilers, headers, display stacks, and native
libraries. No headless or cross-platform runtime pass was claimed.

## Project, Quality, And Vulnerability Measurements

Ply selects root GLFW only through:

`golang.org/x/exp@v0.0.0-20191030013958-a1ab85dbe136 ->
github.com/go-gl/glfw@v0.0.0-20190409004039-e6da0acd62b1`.

`go mod why -m` says the main module does not need root GLFW. Neither root
GLFW nor nested v3.3 packages are loaded. MVS nevertheless retains the root
because x/exp's module graph declares it; nested module selection is a
separate identity.

Retention preserves 234 modules, 3,583 graph edges, 429 complete-test entries,
41 loaded modules, 197 loaded module-backed packages, 1,051 sum lines, and a
383-line unapplied tidy projection. Relative to accepted go-cmp commit
`c314bcb`, sums remain exactly +35/-0. Exact selected `go get` would add only
a redundant main edge and source sum, yielding 3,584 edges, 1,052 sums, and a
389-line tidy projection. Candidate/latest projections change only root GLFW,
retain all load counts, add two sums, and produce 392 tidy lines. None was
applied.

Project mod verify, build, count-1, two count-10 repeats, race, vet,
Windows-amd64 build, pinned lint, byte-identical core help, API/CLI
compatibility, empty-HOME count-2, and complete preflight pass. Preflight
includes 62 launcher checks, 80/80 killed mutants, and 15/15 audit self-tests.
Changed-selection-only snapshot, Docker, and exact quality work is
inapplicable. The accepted 27/27 Q0-Q2 L2 scorecard remains
`dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.

A supplemental contained-Go-1.18 project build passes after removing the
newer toolchain directive in a disposable projection. Full tests retain two
pre-existing Darwin `pkg/shell` closed-file error-text failures; GLFW is
unloaded and unrelated.

Fresh govulncheck v1.7.0 used 1,392 primary records. Root and nested GLFW have
no record, finding, or trace. Selected/candidate results normalize identically
to 20 IDs/22 Darwin and Windows reachable traces, 22 Darwin package findings,
and 30 Darwin module findings. Existing x/image/ansimage findings remain
unrelated.

The 220-entry selected evidence manifest SHA-256 is
`73a24b7220ca0c8df778188840fdbb5ed31195b8509f4aa61916ca3931f8f696`;
decision-summary SHA-256 is
`3e4009ecf940c8627c5fdbea7daa6ec94d100e70529a9bb3dbddb265a5143e4e`.

## Tools And Retained Decisions

- Exact Go 1.26.7 binary/archive SHA-256 values are
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`
  and `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d`.
  Contained Go 1.18.10 binary/archive SHA-256 values are
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
  and `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade`.
- Pinned golangci-lint 2.12.2 archive SHA-256 is
  `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`.
  GoReleaser 2.17.1 and portable apidiff receipts remain
  `f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`
  and `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
- Retain accepted Fatih Color v1.15.0, XXHash v2.3.0, and Speakeasy v0.2.0
  moves and every earlier exact decision. Answered archives are authoritative
  detail.

## Next Objective

Independently evaluate exact-path `github.com/go-logfmt/logfmt v0.4.0` as the
next single P7 group. Do not combine Prometheus, `kr/logfmt`, or any other
dependency group.

MVS selects v0.4.0 through Prometheus Common v0.9.1 while Prometheus TSDB
v0.7.1 requests v0.3.0; `go mod why -m` says the main module does not need
Go Logfmt. The proxy exposes eight stable tags from v0.1.0 through v0.6.1.
Selected v0.4.0 has checksum pair
`h1:MP4Eh7ZCb31lleYCFuwm0oe4/YGak+5l1vA2NOE80nA=` /
`h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=` and a directive-free module
file requiring historical `github.com/kr/logfmt`.

Latest v0.6.1 declares Go 1.21 and is provisionally floor-ineligible. V0.6.0
declares Go 1.17 and is the initial serious candidate; v0.5.1 also declares Go
1.17 and v0.5.0 declares Go 1.13. Resolve complete release/repository identity,
closure, behavior, API, loading, MVS, project gates, and vulnerability results
without treating this survey as accepted evidence.
