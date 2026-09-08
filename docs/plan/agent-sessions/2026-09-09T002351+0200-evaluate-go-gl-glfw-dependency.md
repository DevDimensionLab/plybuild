# Agent Session: Evaluate Go GL GLFW Dependency

Status: NEXT
Session ID: `2026-09-09T002351+0200-evaluate-go-gl-glfw-dependency`
Created: `2026-09-09T00:23:51+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c9dea36852bba2a8e37f1c7c920f783b841ed278bfe21c772457a8a83e589606`
Previous: [2026-09-08T230017+0200-evaluate-ghodss-yaml-dependency.md](2026-09-08T230017+0200-evaluate-ghodss-yaml-dependency.md)
Next: none
Outcome: pending

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
