# Agent Session: Evaluate Go GL GLFW Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-09T002351+0200-evaluate-go-gl-glfw-dependency`
Created: `2026-09-09T00:23:51+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c9dea36852bba2a8e37f1c7c920f783b841ed278bfe21c772457a8a83e589606`
Previous: [2026-09-08T230017+0200-evaluate-ghodss-yaml-dependency.md](2026-09-08T230017+0200-evaluate-ghodss-yaml-dependency.md)
Next: [2026-09-09T015029+0200-evaluate-go-logfmt-logfmt-dependency.md](2026-09-09T015029+0200-evaluate-go-logfmt-logfmt-dependency.md)
Outcome: retained exact root GLFW pseudo-version without dependency metadata changes

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/go-gl/glfw v0.0.0-20190409004039-e6da0acd62b1` as one bounded
dependency group. Resolve its complete root-module release and repository
identity, full Go-floor closure, Cgo/native and platform behavior, actual
project loading, exact MVS effects, and every applicable quality contract.
Retain or select only an exact root-path version whose complete minimal
closure preserves Go 1.18 and whose relevant behavior passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, and Fatih Color v1.15.0 moves. Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify, and
Ghodss YAML remain retained. All earlier decisions and lifecycle ancestry are
final. Do not revisit them or combine another dependency group. P8 remains
queued.

Project MVS selects this historical root module through the declared edge
`golang.org/x/exp v0.0.0-20191030013958-a1ab85dbe136 -> github.com/go-gl/glfw
v0.0.0-20190409004039-e6da0acd62b1`. `go mod why -m github.com/go-gl/glfw`
says the main module does not need it. Keep this exact root module distinct
from the separately selected nested module
`github.com/go-gl/glfw/v3.3/glfw`; do not combine, upgrade, remove, or audit
that nested path, x/exp, Gio, graphics stacks, windowing libraries, or other
dependency groups.

A minimal post-Ghodss survey finds an empty proxy version list. Selected is
the fetchable pseudo-version
`v0.0.0-20190409004039-e6da0acd62b1`, with source/mod checksum pair
`h1:QbL/5oDUmRBzO9/Z7Seo6zf912W/a6Sr4Eu0G/3Jho0=` /
`h1:vR7hzQXu2zJy9AVAgeJqvqgH9Q5CA+iKCZ2gyEVpxRU=`. Its synthesized module
file contains only `module github.com/go-gl/glfw`, without a Go directive or
requirements. Selected commit is
`e6da0acd62b1b57ee2799d4d0a76a7d4514dc5bc`, tree
`17ab3d23b59cab5cffcafe1232da9c2af635a2a1`, parent
`39f94f8075907c0c6524d2791e05d625627fe268`, unsigned, with commit time
2019-04-09T00:40:39Z.

Exact-path `@latest` is the newer pseudo-version
`v0.0.0-20260823155953-d41da22a9587`, checksum pair
`h1:OWknICoxrl3cDP3NtbCnTgntY+0CM5RNam8IXHK0NlU=` /
`h1:fOxQgJvH6dIDHn5YOoXiNC8tUMMNuCgbMK2yZTlZVQA=`. Its module file declares
Go 1.19, so it is provisionally ineligible for the retained Go 1.18 floor.
Commit `d41da22a9587f777098f96d37014f6cdd35d1afb` is GitHub-verified, has tree
`fcbf95c11882251c94a9ab80ab2205c1edfb7733`, parent
`8fa725d95c7913e898bcb58962caa953fefb151a`, and time
2026-08-23T15:59:53Z. Treat all incoming facts only as a survey to verify.

Go-import metadata names `https://github.com/go-gl/glfw.git`. The public
repository is enabled, unarchived, non-fork, defaults to `master`, and was
pushed in 2026. Resolve the historical root module, repository tags and
branches, proxy-absent stable/prerelease identities, pseudo-version ancestry,
root-versus-nested module boundaries, signatures, commit time/tree/parents,
default-branch ancestry, repository status, deprecation, retractions,
redirects, forks, alternate paths, and any serious exact-root candidate. Do
not silently promote a repository tag belonging to a nested module or a
floor-ineligible unreleased commit.

# Measurements At Start

The latest dependency implementation remains Fatih Color v1.15.0 commit
`6ca672ef38688b7f6f505cf0cb273d07c4c2ba9a`, parent
`d181fcd6c11fa147e0b44dd007598872875fc9e6`, and tree
`9ba2fdc442553622028a4a8464536915d345510c`, changing only `go.mod` and
`go.sum` with three insertions and one deletion. Ghodss YAML v1.0.0 was
retained without a dependency implementation commit because it is the only
exact-path stable release; unreleased master was not substituted.

Accepted project measurements remain 234 selected modules, 3,583 graph
edges, 429 native complete-test entries, 41 loaded modules, 197 loaded
module-backed packages, 1,051 `go.sum` lines, and a 383-line unapplied tidy
projection. Relative to accepted go-cmp commit `c314bcb`, metadata adds
exactly 35 checksum lines and removes zero. The main module retains Go 1.18
and toolchain Go 1.26.7.

Fresh primary vulnerability data has 1,392 module records. Accepted
populations remain 20 IDs/22 reachable traces for Darwin and Windows symbol
scans, 22 Darwin package findings, and 30 Darwin module findings. Ghodss YAML
has no record, finding, or trace. Do not attribute existing x/image/ansimage
or other inherited findings to the root GLFW module.

The unchanged accepted quality baseline has all 27 Q0-Q2 rows PASS at L2,
with scorecard SHA-256
`dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.
Ghodss YAML decision-summary SHA-256 is
`724f7f2ac4394296cd36be360df90539deed1ecaf0b92386c6027b0f2c45f7b9`;
its 92-entry selected evidence manifest SHA-256 is
`cfb5bfcc497db42efebda5323b9e187aa63520192a7354631fd93391970bd771`.

Read the answered Ghodss YAML archive and rolling handover for its complete
release, conversion, closure, MVS, vulnerability, and evidence record.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

# Role And Boundaries

From fresh external archives and caches, resolve the exact proxy, checksum
database, go-import metadata, authoritative repository, root-module history,
and primary Go vulnerability identity. Record exact source/mod checksums, Git
identity, timestamps, trees, parents, signature state, ancestry, repository
status, deprecation/retraction state, and every tag or pseudo-version that
matters to the exact root path.

Prove complete minimal module and package/test closure under exact Go 1.26.7
and contained Go 1.18.10. Inspect imported source and test dependencies rather
than treating the synthesized root module file as floor proof. Keep isolated
source-time resolution separate from the project's selected graph.

Inspect every root-module package and its public API, including versioned
directories within that root module. Resolve Cgo directives, bundled C/header
sources, generated bindings, callbacks, unsafe use, thread-affinity rules,
init/terminate lifecycle, error propagation, concurrency, build tags,
examples, testdata, fuzz/property coverage, and upstream CI. Characterize
Darwin frameworks, Linux X11/Wayland or other system-library requirements,
Windows behavior, unsupported targets, headless behavior, and which checks
can truthfully run in the contained host. Do not hide native prerequisites or
turn unavailable display/system libraries into a passing runtime claim.

Add independent fixtures where useful for release-relevant package/build,
constants/value behavior, pure-Go or native boundaries, errors, determinism,
and concurrency. Run source verification, package listing, native complete
tests, two independent repeats, race where supported, vet, and relevant
cross-build or compile probes under both SDKs. Classify toolchain, OS,
display-server, and system-library failures precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected
and every serious exact-root candidate in disposable trees. Explain why the
root module exists in MVS while no package is loaded, and keep any nested
v3.3 package loading separate. Any change outside the exact root GLFW edge
and its authorized MVS projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, native behavior, concurrency, tests,
API, loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, Fatih implementation identity, reciprocal archive history,
P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive, the
answered Ghodss YAML archive, rolling handover, roadmap, `go.mod`, `go.sum`,
and every referenced quality, compatibility, release, runner, evidence, and
lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine whether any qualified exact-root stable or pseudo-version
higher than selected exists. Do not promote the distinct nested v3.3 module,
a redirect, fork, alternate path, floor-ineligible commit, or repository tag
that does not version the root module. If no higher exact-root candidate
qualifies, retain selected without hand-editing metadata or manufacturing a
dependency commit.

For a changed selection, use exact Go 1.26.7 and exact `go get` for one
dependency-only commit, never tidy as implementation, then run the complete
P7 dependency gate. For an inert/retained selection, prove the no-change
effect and run applicable dependency, closure, project, compatibility,
vulnerability, empty-HOME, and cleanliness gates without rerunning
changed-selection-only snapshot/Docker/audit work merely to manufacture
activity. Exact changed-selection `make quality` must exit 0 with all 27
Q0-Q2 rows PASS at L2 and zero held, regressed, not-comparable, or dirty
counts; full audit may exit 1 only for established queued L3 rows, never 2.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
create direct `/private/tmp/ply-*` roots, never bypass scratch cleanup, and
never run `go mod download all` in a measured worktree. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the root GLFW decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Retain exact-path `github.com/go-gl/glfw
v0.0.0-20190409004039-e6da0acd62b1` without editing `go.mod` or `go.sum` and
without manufacturing a dependency implementation commit. No higher
exact-root candidate satisfies the combined release-identity, Go 1.18 floor,
and complete native-package contracts. This is retention of an unloaded
historical MVS selection, not a claim that its obsolete native packages are a
recommended graphics dependency.

### Release And Repository Identity

Fresh proxy, sumdb, go-import, Git, and GitHub evidence agree on selected
source/mod sums
`h1:QbL/5oDUmRBzO9/Z7Seo6zf912W/a6Sr4Eu0G/3Jho0=` /
`h1:vR7hzQXu2zJy9AVAgeJqvqgH9Q5CA+iKCZ2gyEVpxRU=`. The proxy ZIP SHA-256 is
`96c694c42e7b866ea8e26dc48b612c4daa8582ce61fdeefbe92c1a4c46163169`;
all 227 repository-backed files match commit
`e6da0acd62b1b57ee2799d4d0a76a7d4514dc5bc` byte for byte. The remaining
proxy file is the synthesized root `go.mod`, which contains only `module
github.com/go-gl/glfw`.

Selected has tree `17ab3d23b59cab5cffcafe1232da9c2af635a2a1`, parent
`39f94f8075907c0c6524d2791e05d625627fe268`, committer/pseudo-version time
`2019-04-09T00:40:39Z`, and no signature. It is an ancestor of current
default branch `master`. Go-import resolves
`https://github.com/go-gl/glfw.git`; the public repository is enabled,
unarchived, non-fork, and not redirected.

The exact-root proxy list is empty. The repository has no root-module stable
or prerelease tag and no GitHub release. Its only two tags,
`v3.4/glfw/v0.1.0-pre.1` and `v3.4/glfw/v0.1.0-pre.2`, belong to the distinct
nested `v3.4/glfw` module. There is no retraction, module deprecation, root
redirect, or exact-path alternate stable identity. The separately selected
`github.com/go-gl/glfw/v3.3/glfw` module was neither combined nor audited.

The last serious Go-floor-compatible default-branch root candidate is
`v0.0.0-20240506104042-037f3cc74f2a`, with sums
`h1:FAC6eA052T8d4Lp5GR38Sxta8fnq//jjBGDtsT0TVAU=` /
`h1:wyvWpaEu9B/VQiV1jsPs7Mha9I7yto/HqIBw197ZAzk=` and ZIP SHA-256
`fc50e87ad1f4bc66dd991bea3de354f799e7fd7aef496dbc7e19b864e4ef0e23`.
Its root module declares Go 1.12. GitHub verifies merge commit
`037f3cc74f2ab0b249928c4fc4b61e0f13befdb8`, tree
`d16f4bb5267aa50bdb30ea99992c3e0a2381e5b7`, parents
`a69d953ea14231b2f19998f0b434b6c1e6c6f710` and
`23bba3a4c89646aa52826de643bfde2a8e4cfe98`, at
`2024-05-06T10:40:42Z`. It fails the Darwin/arm64 contract described below.

Later floor-compatible pseudo-version
`v0.0.0-20260227151118-fa5a0d083896` is from nested-v3.4 development rather
than a released/default root head and changes no root package implementation
relative to the 2024 candidate. Its source/mod sums are
`h1:QuL5wPRo3euHS1jVtTKoBWdPZjCg/4AvIvmL1sqd1go=` /
`h1:wyvWpaEu9B/VQiV1jsPs7Mha9I7yto/HqIBw197ZAzk=`. Unsigned commit
`fa5a0d0838962e92f7846982429caf1716bb46b2`, tree
`e4defc2173880171b4b9ee2eea99df400a4b6aa9`, parent
`3db5e4c7a5901ef58ab50842a938a7f228139301`, has committer time
`2026-02-27T15:11:18Z`. Exact-root latest
`v0.0.0-20260823155953-d41da22a9587` is independently ineligible because its
module declares Go 1.19. Its root package implementations are also identical
to the rejected 2024 candidate.

The incoming latest Git identities were incorrect. Fresh primary evidence
gives tree `fcbf95ef3a46826e0cbc29777b6f83d60265a57d` and parent
`8fa725040a18feef1098d802d2108effa73ea652`, not the surveyed values. Commit
`d41da22a9587f777098f96d37014f6cdd35d1afb` is GitHub-verified at
`2026-08-23T15:59:53Z`; the checksum pair remains the surveyed
`h1:OWknICoxrl3cDP3NtbCnTgntY+0CM5RNam8IXHK0NlU=` /
`h1:fOxQgJvH6dIDHn5YOoXiNC8tUMMNuCgbMK2yZTlZVQA=`.

### Closure, API, And Native Behavior

The selected root module contains exactly three packages:
`v3.0/glfw`, `v3.1/glfw`, and `v3.2/glfw`. Both exact Go 1.26.7 and contained
Go 1.18.10 resolve the same packages and no non-standard-library Go module
dependency. Selected contains 46 Go and 172 C, Objective-C, or header files,
but zero `_test.go` files, examples, testdata, fuzz targets, or property tests.
The 2024 candidate has no real package tests either; its two `testdata`
directories are standalone custom-cursor programs. Upstream Travis tested
only v3.2 in selected and v3.2 plus the then-nested v3.3 in the candidate;
current GitHub workflows test only nested v3.3/v3.4 modules.

All 104, 133, and 150 exported declaration signatures for v3.0, v3.1, and
v3.2 are unchanged in the 2024 candidate. Only documentation references
change. V3.1/v3.2 install a global C error callback, buffer one last error,
return designated recoverable errors, log platform errors, and panic on
unexpected/programmer errors. C-to-Go callbacks, global callback holders,
unsafe pointers, and selected's `reflect.SliceHeader` conversions rely on
GLFW's documented main-thread/event-loop discipline rather than independent
locking. `Init`, `Terminate`, event processing, and most window operations
must execute on an OS-locked main thread; context operations are bound to the
calling thread. `PostEmptyEvent` is the documented secondary-thread wakeup.

V3.0 links an externally installed GLFW 3.0 library on every supported OS.
V3.1/v3.2 bundle GLFW C/header source. Darwin uses Cocoa, OpenGL, IOKit, and
CoreVideo frameworks; Windows requires Win32 headers/libraries (and v3.0's
GLFW DLL import library); Linux defaults to X11/OpenGL and v3.2 has a Wayland
tag requiring Wayland/EGL/xkbcommon libraries. FreeBSD adds its own X11 or
Wayland libraries and v3.2 `pkg-config: glfw3`. OpenBSD and NetBSD select
generic Go Cgo files but have no native backend and are unsupported. With
`CGO_ENABLED=0`, all three packages disappear on every tested target; there
is no pure-Go or headless fallback.

Using the real Xcode compiler/SDK and a scratch-built universal static GLFW
3.0.4 prerequisite for v3.0, all three selected and candidate packages pass
Darwin/amd64 count-1, two independent count-10 repeats, and race under both
SDKs. Candidate vet passes both SDKs; selected vet reports the same possible
`reflect.SliceHeader` misuse in v3.1/v3.2. Independent per-package fixtures
pass constants, version values, v3.1/v3.2 error formatting, repeated and
concurrent read-only version queries, and v3.2 `Init`/`Terminate` under both
SDKs. The contained macOS lifecycle initializes despite sandbox notification
warnings; it does not prove window creation, rendering, or Linux headless
display behavior. Importing all three historical versions into one binary
fails with duplicate bundled C and exported callback symbols.

Darwin/arm64 is the candidate stop condition. V3.0 and v3.2 compile and run,
but default v3.1 fails selected and candidate under both SDKs with `No
supported client library selected`: its build directives define OpenGL only
for 386/amd64 and GLES for arm/explicit tags, not default arm64. An explicit
`gles2` compile probe succeeds, but macOS does not supply OpenGL ES and this is
not a valid runtime substitute. Cross-target package selection succeeds for
Linux, Windows, and FreeBSD under both SDKs, while truthful Cgo compilation is
unavailable on the contained Darwin host because it has no target C compiler,
headers, or native libraries. Those probes fail on missing pthread/standard,
Windows, or target-native headers; they are not runtime passes.

### Project, Vulnerability, And Quality Effects

Project MVS retains the root only through
`golang.org/x/exp@v0.0.0-20191030013958-a1ab85dbe136 ->
github.com/go-gl/glfw@v0.0.0-20190409004039-e6da0acd62b1`. `go mod why -m`
reports that the main module does not need it. Neither the root nor nested
v3.3 package is in the Darwin complete-test load set. MVS still selects the
root because module graphs retain declared requirement edges even when no
package import uses them.

Retention preserves 234 selected modules, 3,583 graph edges, 429 complete-test
entries, 41 loaded modules, 197 loaded module-backed packages, 1,051 checksum
lines, and the 383-line unapplied tidy projection. Relative to accepted
go-cmp commit `c314bcb`, `go.sum` remains exactly +35/-0 lines. An exact
selected `go get` would manufacture a redundant explicit main edge and source
checksum: 3,584 edges, 1,052 checksum lines, and 389 tidy-diff lines without
loading a package. The 2024/latest projections change only the root selection,
retain all module/package counts, add the explicit edge plus two sums, and
produce 392 tidy-diff lines. They were not applied.

Fresh govulncheck v1.7.0 used all 1,392 primary module records. No GLFW root
or nested path has a record. Selected and candidate normalized findings are
identical: 30 Darwin module findings, 22 Darwin package findings, 20 IDs/22
reachable Darwin symbol traces, and 20 IDs/22 reachable Windows traces. No
trace names GLFW; inherited x/image/ansimage findings remain unrelated.

Exact Go 1.26.7 mod verification, build, count-1, two count-10 repeats, race,
vet, Windows-amd64 build, pinned golangci-lint 2.12.2, API/CLI compatibility,
byte-identical core help, empty-HOME count-2, and complete preflight pass.
Preflight includes all 62 launcher checks, 80/80 killed mutants, and 15/15
audit self-tests. A contained-Go-1.18 project build passes in a disposable
floor projection after removing the newer toolchain directive; its full tests
retain only two inherited Darwin `pkg/shell` closed-file error-text failures.
GLFW is unloaded in both cases. Changed-selection-only snapshot, Docker, and
quality work is inapplicable; the accepted 27/27 Q0-Q2 L2 scorecard remains
`dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.

The 220-entry evidence manifest SHA-256 is
`73a24b7220ca0c8df778188840fdbb5ed31195b8509f4aa61916ca3931f8f696`;
decision-summary SHA-256 is
`3e4009ecf940c8627c5fdbea7daa6ec94d100e70529a9bb3dbddb265a5143e4e`.
