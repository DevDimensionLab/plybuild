# Agent Session: Evaluate Google Renameio Dependency

Status: NEXT
Session ID: `2026-09-09T232057+0200-evaluate-google-renameio-dependency`
Created: `2026-09-09T23:20:57+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ca10b944ae5cb376a8184cbcdff76b37fca8d0c3b6e4e627309f52985a4c8e81`
Previous: [2026-09-09T214949+0200-evaluate-google-pprof-dependency.md](2026-09-09T214949+0200-evaluate-google-pprof-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/google/renameio v0.1.0` as one bounded dependency group. Resolve
its complete repository and release identity, full Go-floor closure, package
behavior and public API, actual project loading, exact MVS effects, and every
applicable quality contract. Retain or select only an exact-path version whose
complete minimal closure preserves Go 1.18 and whose relevant behavior passes
every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, and Google Martian v3.3.2 moves. Google pprof remains retained at
`v0.0.0-20210720184732-4bb14d4b1be1`. Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, and historical root GLFW also remain retained. All earlier
decisions and lifecycle ancestry are final. Do not revisit them or combine
another dependency group. P8 remains queued.

Project MVS selects Renameio v0.1.0 through three historical HTools vertices:
`honnef.co/go/tools` v0.0.1-2019.2.3, v0.0.1-2020.1.3, and
v0.0.1-2020.1.4 each declare it. One exact path from the main module is
`mvn-pom-mutator v0.2.3 -> bketelsen/crypt -> Firestore v1.1.0 -> Cloud Go
v0.46.3 -> HTools v0.0.1-2019.2.3 -> Renameio v0.1.0`. MVS selects HTools
v0.0.1-2020.1.4. `go mod why -m github.com/google/renameio` is negative, and
the complete project load contains no Renameio package. Resolve this ancestry
precisely. Do not combine, upgrade, remove, or independently audit HTools,
Crypt, Firestore, Cloud Go, or another dependency group. Candidate MVS
movements caused by the exact Renameio edge remain in scope to measure, but
not to broaden into independent audits.

A minimal post-pprof survey finds exactly v0.1.0, v1.0.0, and v1.0.1 in the
Go proxy list. Selected is dated 2019-01-09T16:53:11Z, has source/mod sums
`h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=` /
`h1:KWCgfxg9yswjAJkECMjeO8J8rahYeXnNhOm40UhjYkI=`, and its module file
contains no Go directive or requirements. Proxy latest v1.0.1 is dated
2021-04-06T14:11:08Z, has source/mod sums
`h1:Lh/jXZmvZxb0BBeSY5VKEfidcbcbenKjZFzM/q0fSeU=` /
`h1:t/HQoYBZSsWSNK35C6CO/TpPLDVWvxOHboWUAweKUpk=`, and declares Go
1.13 with no requirements. An exact disposable v1.0.1 projection adds the
root edge, changes only Renameio selection, adds its source/mod checksum pair,
keeps 234 modules and 429 complete-test entries, changes graph edges from
3,597 to 3,598 and sum lines from 1,063 to 1,065, and loads no Renameio
package. Treat every incoming fact only as a survey to verify.

Resolve proxy, sumdb, go-import metadata, repository tags/releases/branches,
signatures, commits, times, trees, parents, ancestry, repository status,
deprecation, retractions, redirects, forks, alternate module paths, and every
serious exact-path stable candidate. Do not silently promote an unreleased
commit, arbitrary branch head, redirect, fork, alternate path,
floor-ineligible version, prerelease, or tag that does not version this
module.

# Measurements At Start

The latest dependency implementation remains exact Google Martian v3.3.2
commit `4644476aef78e24d0c7e6a1140baba4c3226a1cc`, parent
`8508a9d54f066985fb81ebd35fa31f45fa36da40`, and tree
`f20d9b3817fe29af4876745d09c14e9cbcac9152`, changing only `go.mod` and
`go.sum` with three insertions and no deletions. Google pprof and Gogo
Protobuf remain retained without dependency commits or metadata edits.

Accepted project measurements are 234 selected modules, 3,597 graph edges,
429 native complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,063 `go.sum` lines, and a 418-line unapplied tidy projection.
Relative to accepted go-cmp commit `c314bcb`, metadata adds exactly 47
checksum lines and removes zero. The main module retains Go 1.18 and toolchain
Go 1.26.7.

Fresh primary vulnerability data has 1,393 module records, modified
2026-09-09T17:56:34Z. Accepted project populations remain 30 module findings,
22 Darwin package findings, 23 Windows package findings, and 20 IDs/22
reachable traces for both Darwin and Windows symbol scans. Pprof has no exact
advisory record, package finding, or trace. Do not attribute inherited
findings to Renameio without exact evidence.

The accepted quality result is all 27 Q0-Q2 rows PASS at L2, with scorecard
SHA-256 `094a7668e43eb0a5b563bcd637b3b5a301e1bb5f5053903ad8934e2b2906ae7f`.
Google pprof decision-summary SHA-256 is
`20804b312d90e5390df74dde4c1b252e68d8b1183897e7e562f2c8a133ac5d08`;
its 188-entry selected-evidence manifest SHA-256 is
`e6090a35f4f669bfff0eccf3b179def466dcbc18645ec61ed5d2e81f92e4f0f7`.

Read the answered Google pprof archive and rolling handover for its complete
identity, closure, API, behavior, MVS, vulnerability, quality, and evidence
record. Do not reopen pprof, Martian, Cloud Go, HTools, Crypt, or earlier
groups.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256 `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and
inject no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools
beneath scratch as required. Portable receipts remain golangci-lint 2.12.2
archive `a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff source archive
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.

# Role And Boundaries

From fresh external archives and caches, prove the complete minimal module and
package/test closure for selected and every serious candidate under exact Go
1.26.7 and contained Go 1.18.10. Inspect imported source and test dependencies
rather than treating module Go directives as full closure proof. Keep
isolated source-time resolution separate from the project's selected graph.

Inspect every package and exported API. Characterize atomic replacement and
durability guarantees; temporary-file placement, permissions, ownership,
cleanup, close/sync/rename ordering, and error propagation; same-directory,
symlink, hard-link, directory, missing-parent, existing-destination,
read-only, cross-device, short/failing I/O, and concurrent-writer behavior;
resource ownership, deterministic names/content, reuse, platform semantics,
and Windows/Unix differences. Distinguish runtime libraries, internal
helpers, commands, examples, benchmarks, testdata, generated content,
fuzz/property coverage, and upstream CI.

Add independent fixtures where useful for successful atomic writes and
replacement, old-or-new visibility, file modes, cleanup after callback/write,
sync, close, or rename failures, short writes, pre-existing targets,
symlinks/hard links, temp-directory boundaries, concurrent writers, resource
reuse, and compatibility with the selected project graph. Run source
verification, package listing, native complete tests, two independent
repeats, race, vet, and meaningful cross-builds under both SDKs. Classify
every generator, filesystem, timing, platform, resource, or test-design
failure precisely.

Prove exact project module/graph/package/checksum/tidy effects for selected and
every serious candidate in disposable trees. Explain why Renameio exists in
MVS while no package is loaded, and preserve every unrelated module
selection. Any change outside the exact Renameio edge and its necessary
authorized MVS projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, behavior, concurrency, tests, API,
loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this
archive, the answered pprof archive, rolling handover, roadmap, `go.mod`,
`go.sum`, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release. Do not
promote a redirect, fork, alternate path, floor-ineligible release,
prerelease, non-versioning tag, or unreleased commit. If no higher release
qualifies, retain selected without hand-editing metadata or manufacturing a
dependency commit.

For a changed selection, use exact Go 1.26.7 and exact `go get` for one
dependency-only commit, never tidy as implementation, then run the complete
P7 dependency gate. For an inert/retained selection, prove the no-change
effect and run applicable dependency, closure, project, compatibility,
vulnerability, empty-HOME, and cleanliness gates without manufacturing
activity. Full changed-selection quality must preserve 27/27 Q0-Q2 PASS at L2
with zero held, regressed, non-comparable, or dirty counts.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`.
Never run `go mod download all` in a measured worktree or create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Google Renameio decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
