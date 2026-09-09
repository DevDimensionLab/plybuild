# Agent Session: Evaluate Godbus D-Bus v5 Dependency

Status: NEXT
Session ID: `2026-09-09T082209+0200-evaluate-godbus-dbus-v5-dependency`
Created: `2026-09-09T08:22:09+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `608ae7845227f40ce4c246a4e76cf58585b7b2bff098a4fcb90584a6411e6dae`
Previous: [2026-09-09T052020+0200-evaluate-go-stack-stack-dependency.md](2026-09-09T052020+0200-evaluate-go-stack-stack-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/godbus/dbus/v5 v5.0.4` as one bounded dependency group. Resolve
its complete release and repository identity, full Go-floor closure, package
behavior and public API, actual project loading, exact MVS effects, and every
applicable quality contract. Retain or select only an exact-path version whose
complete minimal closure preserves Go 1.18 and whose relevant behavior passes
every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0, and
Go Stack v1.8.1 moves. Crypt, OpenCensus Proto, Logex, Readline, Fnmatch,
Imaging, ansimage, Fsnotify, Ghodss YAML, and historical root GLFW remain
retained. All earlier decisions and lifecycle ancestry are final. Do not
revisit them or combine another dependency group. P8 remains queued.

Project MVS selects Godbus D-Bus v5 through the declared edge
`github.com/coreos/go-systemd/v22 v22.3.2 -> github.com/godbus/dbus/v5
v5.0.4`. `go mod why -m github.com/godbus/dbus/v5` says the main module
does not need it. Do not combine, upgrade, remove, or independently audit
go-systemd or any other dependency group.

A minimal post-Go-Stack survey finds eleven proxy versions: v5.0.0, v5.0.1,
v5.0.2, v5.0.3, selected v5.0.4, v5.0.5, v5.0.6, v5.1.0, v5.2.0, v5.2.1,
and v5.2.2, with no prereleases. Selected's source/mod checksum pair is
`h1:9349emZab16e7zQvpmsbtjc18ykshndd8y2PG3sgJbA=` /
`h1:xhWf0FNVPg57R7Z0UbKHbJfkEywrmjJnf7w5xrFpKfA=`; it declares Go 1.12
and has no requirements.

Latest v5.2.2 declares Go 1.20 and requires `golang.org/x/sys v0.27.0`, so
it and v5.2.0-v5.2.1 are initially floor-ineligible. Highest initial serious
candidate v5.1.0 declares Go 1.12, has no requirements, and has source/mod
sums `h1:4KLkAxT3aOY8Li4FRJe/KvhoNFFxo0m6fNuFUO8QJUk=` /
`h1:xhWf0FNVPg57R7Z0UbKHbJfkEywrmjJnf7w5xrFpKfA=`. Its lightweight tag
resolves to commit `e523abc905595cf17fb0001a7d77eaaddfaa216d`, tree
`af2a399b7f6f0c3a62c435e15a488c46d3dbf536`, parents
`b357b44b7ab3bf9e9b27c906fb31cb622b7a017e` and
`2c3cf657d630429b62f481c930c754a818d98327`, at
2022-02-27T11:53:47Z. Treat every incoming release fact only as a survey to
verify, including the proxy's anomalous v5.0.0/v5.0.1 module-file responses.

Resolve proxy, sumdb, go-import metadata, repository tags/releases/branches,
signatures, commits, times, trees, parents, ancestry, repository status,
deprecation, retractions, redirects, forks, alternate module paths, and every
serious exact-path candidate. Do not silently promote a redirect, fork,
alternate path, floor-ineligible release, or tag that does not version this
module.

# Measurements At Start

The latest dependency implementation is exact Go Stack v1.8.1 commit
`647d4fd71b226fbb1e916b4238b4c7ad87cd7975`, parent
`171ffd27ac9502c6018c31b2962a6163b2433a03`, and tree
`7c54d59e484d197e7290bc0eec1281c75f614fcd`, changing only `go.mod` and
`go.sum` with three insertions and no deletions.

Accepted project measurements are 234 selected modules, 3,585 graph edges,
429 native complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,055 `go.sum` lines, and a 396-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, metadata adds exactly 39
checksum lines and removes zero. The main module retains Go 1.18 and toolchain
Go 1.26.7.

Fresh primary vulnerability data has 1,392 module records. Accepted
populations remain 20 IDs/22 reachable traces for Darwin and Windows symbol
scans, 22 Darwin package findings, and 30 Darwin module findings. Go Stack
has no record, finding, or trace. Do not attribute inherited findings to
Godbus without exact evidence.

The accepted quality result has all 27 Q0-Q2 rows PASS at L2, with scorecard
SHA-256
`4e1e1b50c2d16666d18278efd0b9ded86df668503e32184c49a1024053dd227c`.
Go Stack decision-summary SHA-256 is
`26fca7f6b5b683f4d16fb6b18a9aa88e412948ffde04c7d71764a1a94c077b72`;
its 817-entry selected-evidence manifest SHA-256 is
`e3edb91f49616d36c93ea08e85eeba8ce5730d3f61dea8b668b4da49405a8421`.

Read the answered Go Stack archive and rolling handover for its complete
release, closure, API, MVS, vulnerability, quality, and evidence record. Do
not reopen Go Stack, Prometheus, go-systemd, or earlier groups.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

# Role And Boundaries

From fresh external archives and caches, resolve the exact stable/prerelease
population and authoritative repository. Verify every relevant checksum,
module file, tag, Git object, signature, timestamp, tree, parent, ancestry,
release status, retraction, deprecation, redirect, fork, and alternate path.
Keep repository history distinct from exact-path module release identity.

Prove the complete minimal module and package/test closure for selected and
every serious candidate under exact Go 1.26.7 and contained Go 1.18.10.
Inspect imported source and test dependencies rather than treating a Go
directive alone as floor proof. Keep isolated source-time resolution separate
from the project's selected graph.

Inspect every package and exported API. Characterize connection and transport
creation, session/system/bus-address discovery, authentication and negotiation,
message encoding/decoding and signatures, variants, Unix file descriptors,
object paths, names, match rules, exported objects, introspection, properties,
signals, calls/replies/errors, context cancellation and deadlines, disconnect
behavior, nil and invalid values, resource ownership, determinism, concurrency
and reuse, build tags, platform-specific Unix/Windows transports, generation,
examples, testdata, fuzz/property coverage, and upstream CI.

Add independent fixtures where useful for marshaling/signature validation,
messages and variants, connection state, cancellation, match rules, exported
method dispatch, signal delivery, invalid input, determinism, concurrency, and
compatibility. Avoid relying on a host session/system bus when hermetic
socket-pair or private-daemon coverage is possible. Run source verification,
package listing, native complete tests, two independent repeats, race, vet,
and meaningful cross-builds under both SDKs. Classify any toolchain, platform,
resource, daemon, timing, or test-design failure precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected
and every serious candidate in disposable trees. Explain why Godbus exists in
MVS while no package may be loaded, and preserve every unrelated module
selection. Any change outside the exact Godbus edge and its authorized MVS
projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, behavior, concurrency, tests, API,
loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, Go Stack implementation identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
the answered Go Stack archive, rolling handover, roadmap, `go.mod`, `go.sum`,
and every referenced quality, compatibility, release, runner, evidence, and
lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release. Do not
promote a redirect, fork, alternate path, floor-ineligible release, or
unreleased commit. If no higher release qualifies, retain selected without
hand-editing metadata or manufacturing a dependency commit.

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

After the Godbus D-Bus v5 decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
