# Agent Session: Evaluate Cespare XXHash V2 Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-08T074422+0200-evaluate-cespare-xxhash-v2-dependency`
Created: `2026-09-08T07:44:22+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `fa548066df8974bbc587c8fd96d27120a2b1361ed9375cf3a98e6577f4e44987`
Previous: [2026-09-08T060123+0200-evaluate-census-opencensus-proto-dependency.md](2026-09-08T060123+0200-evaluate-census-opencensus-proto-dependency.md)
Next: [2026-09-08T110936+0200-evaluate-chzyer-logex-dependency.md](2026-09-08T110936+0200-evaluate-chzyer-logex-dependency.md)
Outcome: Accepted exact XXHash v2.3.0 in one dependency-only commit; the complete closure preserves Go 1.18 and every dependency, consumer, platform, compatibility, vulnerability, and quality contract passed.

# Answer

Accept exact-path `github.com/cespare/xxhash/v2 v2.3.0`, canonical latest,
over historically selected v2.1.2. Exact Go 1.26.7 `go get` produced the
single dependency-only commit
`e5d6252825d7a1822c01819b9144050f345a6ad4`, tree
`5213ba55981d77d7c8061312915e237e80d29af8`, with only one `go.mod`
requirement and the v2.3.0 source/module checksum pair. The full minimal
closure remains below Go 1.18 and all applicable gates pass.

## Canonical Identity And Release Qualification

The exact proxy lists only stable v2.0.0, v2.1.0, v2.1.1, v2.1.2, v2.2.0,
and v2.3.0. Exact `@latest` is v2.3.0 at 2024-04-04T20:00:10Z. No release is
retracted and no module deprecation directive exists. Selected and latest
declare exact path `github.com/cespare/xxhash/v2` and Go 1.11. The v1 module
uses incompatible path `github.com/cespare/xxhash`; nested `xxhashbench` is a
separate benchmark/tool module, not an alternate v2 path.

Go-import metadata identifies `https://github.com/cespare/xxhash.git`. The
public repository is enabled, unarchived, non-fork, and defaults to `main`.
All six v2 tags are annotated and unsigned, all target commits are unsigned,
and there are no GitHub Release objects.

Selected tag object `7ae26c41ed6fb1f8a6c21e05eeff4d91b5e401c4`
dereferences commit `e7a6b52374f7e2abfb8abb27249d53a1997b09a7`, tree
`78bbde5aae5fc0324314789c69f39da97d995f60`, parent
`3b9a65d476075116e7baf9ca3c37294f637ea033`, at
2021-08-24T09:58:46Z. Latest tag object
`7438b35f14d771ee32d8bbcd9527d32a336e7dad` dereferences commit
`998dce232f17418a7a5721ecf87ca714025a3243`, tree
`9a415e99332a4ef14d2d590e9ed748c99ad4bc70`, parent
`21fc82c137a876186c8acb0349e941ddb280bf03`, at
2024-04-04T20:00:10Z. Default-branch head
`ab37246f3853501fb3e16d199556315b50889ad2` is one later CI-only update and
is not a release; every tag is its ancestor.

Every proxy ZIP matches exact Git source after excluding the nested benchmark
module. Normalized manifests for v2.0.0 through v2.3.0 are respectively
`04cbb9e84ae2b86cd24a46b500e0708e159c3e006fc70133c3d6cd6698d58032`,
`9667f1743cafabff992f020509d2c1538573540bb075b0dc80d0b83d97e965d1`,
`daab817c6c64cab7ae3126847d6d132fbb5aad7b7ea4ad3ef763174b182f84bc`,
`69a8ba25b194eb9b2c195ac76f8f81ab03595cecebf0699a59b244611591352f`,
`2b04a8ff379c136d984046f4db44a3dd2c057595869646190b009023620a796a`,
and `b604c7ab37f70e52a009950fbde4a6708f7372751450549f6e63aff0e7b8bd09`.
The retained evidence records every checksum-database pair.

## Closure, API, Platforms, And Consumers

Selected and candidate have no module requirements. Their three packages pass
verify/list/count-1/two count-10/vet and pure-Go tests under exact Go 1.26.7
and contained Go 1.18.10. Root race tests pass. The dynamic plugin test's
TestMain builds a non-race plugin that a race-enabled host cannot load; this is
identical in both versions, while normal repeated dynamic tests pass, so it is
not a candidate regression. The nested benchmark module's four-module closure
tops out at Go 1.13 and passes count/race/vet under both SDKs.

Independent fixtures pass five known vectors, whole-buffer/string helpers,
streaming digest, reset, size/block size, `hash.Hash64`, and binary marshal
interfaces; v2.3.0 seeded construction/reset also passes. Pinned apidiff finds
only compatible `NewWithSeed` and `(*Digest).ResetWithSeed` additions.

V2.1.2 uses pure Go on Darwin-arm64; v2.3.0 adds Darwin-arm64 assembly. Both
use Windows/Linux amd64 assembly and retain pure-Go, appengine, and non-gc
fallbacks. Exact Go 1.26.7/1.18.10 cross-builds pass for Windows-amd64, Linux
amd64/arm64/386/arm, `xxhsum`, purego, and appengine. Static Linux-amd64
containers execute repeated assembly behavior for both versions.

Ply reaches Viper v1.10.1 through mvn-pom-mutator v0.2.3. Viper requires Crypt
v0.4.0 and XXHash v2.1.2; Crypt also requires v2.1.2. Lower v2.1.1 requests
come from grpc v1.43.0 and Prometheus client v1.4.0. Ply imports and loads no
XXHash package. Viper/Crypt test closures do not load it. The actual historical
consumer, grpc xDS ringhash/resolver, uses `Sum64String`; focused count-1,
count-10, race, and vet pass against selected and candidate under both SDKs.

## MVS, Quality, And Vulnerability Results

Baseline measurements are 234 modules, 3,582 edges, 429 complete-test
packages, 41 loaded modules, 197 loaded module-backed packages, 1,047 checksum
lines, and 371 unapplied tidy lines. Exact selected get changes no selection
and adds only a redundant requirement/main edge/full checksum, so it was not
applied. Exact candidate get projects 234 modules, 3,583 edges, 429/41/197,
1,049 checksums, and 381 tidy lines. Only XXHash changes selection, the main
edge is the only edge addition, and the v2.3.0 pair is the only checksum
addition. Tidy would remove the pin/pair and restore inherited v2.1.2.
Relative to accepted go-cmp commit `c314bcb`, current `go.sum` adds exactly 33
lines and removes zero.

The projected and committed candidate pass mod verify, complete online/offline
package listing, build, count-1/count-10/race/vet, Windows-amd64 build, pinned
lint, byte-identical help, API/CLI compatibility, CLI surface, and empty-HOME
count-2. Exact `make quality` passes all 21 stages: 27/27 Q0-Q2 rows at L2,
80/80 mutants killed, 4/4 acceptance, six valid manual receipts, and zero
held/regressed/not-comparable/dirty counts. The separate audit-meta suite
passes 15 controls. Full audit exits 1 only for queued L3 rows Q3.1, Q3.3,
Q3.4, and Q3.7, attains L2, and never exits 2.

Fresh vulnerability data contains 1,392 records and no XXHash record.
Baseline/candidate results are identical: 20 IDs/22 reachable Darwin and
Windows traces, 22 Darwin package IDs/findings, and 30 Darwin module IDs. No
module, package, symbol, or trace names XXHash.

Evidence has 608 verified entries; manifest SHA-256 is
`6f57be6cbf177c6617ebf62e0eea7b38f6b3b41eb00873060e160ddf02f1cf34`.
Decision-summary SHA-256 is
`4e2f0559dbefd98f22e8f1efe34a1e538ecb8e1dd2da0749af8478c88c0a1a20`.
Authoritative Q0-Q2 scorecard SHA-256 is
`579e5b135db2403904943bfe71c3ced987f59f30007cbdc0dec88a46d19b42aa`.

Scratch-only corrections supplied an explicit template for BSD `mktemp -d`
and placed Python 3.14 ahead of Apple Python 3.9 so Docker 29's eight-digit
fractional creation timestamp could be parsed. A fresh standalone Docker pass
proved the environment correction before the entire quality gate was rerun.
Superseded attempts are not proof. No project source or quality apparatus was
changed, all evaluation images were cleaned, and ordinary/ignored status was
empty before handoff.

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
