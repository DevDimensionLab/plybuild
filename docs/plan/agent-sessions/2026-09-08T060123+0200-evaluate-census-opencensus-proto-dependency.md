# Agent Session: Evaluate Census OpenCensus Proto Dependency

Status: NEXT
Session ID: `2026-09-08T060123+0200-evaluate-census-opencensus-proto-dependency`
Created: `2026-09-08T06:01:23+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `7810606564c924faf4a3e805361be07df5b5e87e49a5c9b7151c764f8412de45`
Previous: [2026-09-08T044755+0200-evaluate-bketelsen-crypt-dependency.md](2026-09-08T044755+0200-evaluate-bketelsen-crypt-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/census-instrumentation/opencensus-proto v0.3.0` as one bounded
dependency group. Resolve canonical latest, authoritative source identity,
archive and release qualification, divergent default-branch history, and the
highest qualified Go-1.18-floor-compatible candidate from primary evidence.
Make an exact dependency selection only if it changes a selected version,
preserves the retained floor through the complete minimal closure, and passes
every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, the accepted
Speakeasy v0.2.0 dependency move, and the retained Bketelsen Crypt
pseudo-version. All earlier rejections, no-change decisions, accepted
closures, and evidence corrections remain final. Do not revisit Crypt or
Speakeasy and do not combine another module group. P8 remains queued.

Project MVS selects OpenCensus Proto v0.3.0 through historical requirements
including `github.com/spf13/viper v1.10.1` and
`github.com/sagikazarmark/crypt v0.4.0`; it is not an explicit main-module
requirement and no OpenCensus Proto package is loaded by Ply. A minimal
post-Crypt survey finds stable proxy versions v0.0.1, v0.0.2, v0.1.0,
v0.2.0, v0.2.1, v0.3.0, v0.4.0, and v0.4.1. Exact Go `@latest` and `@v0`
resolve v0.4.1 at 2022-09-23T17:40:20Z, while exact `@master` resolves
unreleased pseudo-version
`v0.2.2-0.20230502190750-1664cc961550` at
2023-05-02T19:07:50Z. Treat canonical release precedence and the split
default-branch ancestry as facts to prove, not as interchangeable notions of
latest.

The public canonical repository is currently archived, non-fork, and has
default branch `master` at commit
`1664cc961550be8f3058ddd29390350242f44f1f`. Stable v0.4.1 is commit
`e53624a87b9b9b919147a9b4626c669a869ebb34` on the separate v0.4 release
line; selected v0.3.0 is commit
`4aa53e15cbf1a47c6e662018ef4587e12dd0c461`. Independently establish exact
commit/tree/parent/times, tag topology, releases and signatures, archive and
deprecation state, branch heads, ancestry, redirects, forks, and alternate
module paths.

Selected v0.3.0 declares only exact module path
`github.com/census-instrumentation/opencensus-proto`, with no Go directive or
requirements. Stable v0.4.1 declares Go 1.18 and requirements including
grpc-gateway/v2 v2.11.3, gRPC v1.49.0, protobuf v1.28.1, and 2022 x/net,
x/sys, x/text, and genproto modules. Do not infer the complete closure floor
or a minimal MVS move from either root declaration.

# Measurements At Start

The latest dependency implementation remains Speakeasy commit
`41f9561f6ea2f5b6395c5c4d9bcc56a54533133a`, exact parent
`1f55aaa31280a4610ed66c1c37ddda953bfa6a8e`, tree
`a50d4f256fe742e472f1f7e9cf0589a596e971b0`, changing only `go.mod` and
`go.sum` with three insertions. Crypt was retained without dependency edits.

Accepted measurements remain 234 selected modules, 3,582 graph edges, 429
native complete-test packages, 41 loaded modules, 197 loaded packages, 1,047
`go.sum` lines, and a 371-line unapplied tidy projection. Relative to accepted
go-cmp commit `c314bcb`, accepted metadata adds exactly 31 checksum lines. The
main module retains Go 1.18 and toolchain Go 1.26.7. Ordinary and ignored
status must be empty.

Fresh primary vulnerability data contains 1,392 module records and no Crypt
record. Accepted old/candidate comparisons retain exact 20-ID/22-trace Darwin
symbol, 30-ID Darwin module, and 20-ID/22-trace Windows symbol populations.

Crypt evidence has 10,694 verified entries; manifest SHA-256 is
`e0320cf4ef063c83cee9a5131fdad9ceae8b90730502617ab48d2f16c759c93b`;
decision-summary SHA-256 is
`6356ad761738c8d9e664693a0c550bbe4305db500bf9198e567ef7a3e1d33a75`.

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
OpenCensus Proto v0 version through the Go proxy and checksum database,
authoritative repository, go-import metadata, and primary Go vulnerability
data. Record selected/candidate commits and times, module declarations and
requirements, checksum pairs, source identity, tag and commit signatures,
release and branch history, repository archive/deprecation status, and
retractions. Distinguish stable tags, redirects, forks, alternate module
paths, the divergent default-branch pseudo-version, release-branch heads, and
unreleased commits.

Prove canonical latest and the highest qualified candidate that preserves the
Go 1.18 floor through the complete minimal module and package/test closure.
Explain how a release newer by semantic version can lie outside default-branch
ancestry, and qualify stable releases independently of the later master
pseudo-version. Do not upgrade to an unreleased branch head merely because it
is newer by commit time.

Measure old and candidate selected modules, complete graph edges, native
complete-test packages, loaded modules and packages, checksum lines, exact
dependency paths, exact-get effects, and the unapplied tidy projection.
Attribute every selection, edge, and checksum difference. Do not infer that a
root candidate update is minimal merely because MVS already selects newer
transitive modules, and do not combine another dependency group.

Identify the historical Viper and Sagikazarmark Crypt consumers and the exact
generated protobuf/gRPC packages and APIs they use. Exercise compatible
marshal/unmarshal and service/client behavior where meaningful without
inventing a Ply runtime path. Keep dependency-native generator/tool closures
and historical consumer closure distinct from the project selection and
loaded behavior.

Run dependency source verification, package listing, native complete tests,
two independent repeated-test passes, race where supported, and vet under
exact Go 1.26.7 and a contained Go 1.18 SDK. Treat generated-code vet findings,
missing tests, archived status, split history, or incompatible generators
precisely: establish whether each is a release disqualifier under the existing
contracts rather than silently waiving it.

Project any qualified exact selection in a disposable worktree and run
repository verify, build, count-1/count-10/race/vet, Windows-amd64 build,
pinned golangci-lint, byte-identical root/status/upgrade/build help, API/CLI
compatibility and reports, and empty-HOME count-2 before deciding whether an
implementation is permissible. Compare old and candidate primary
vulnerability results at module, package, symbol, and reachable-trace levels.

Reject or retain if canonical identity, archive/release qualification,
complete closure floor, dependency tests, generated API compatibility,
historical consumers, projection, or any quality contract fails. If exact
selected-version get changes no selected version, do not add a redundant
requirement or checksum.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, Speakeasy implementation commit identity, reciprocal archive
history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive,
`docs/plan/quality-handover.md`, `docs/plan/quality-upgrade.md`, `go.mod`,
`go.sum`, the answered Crypt and Speakeasy archives, the earlier dependency
archives named in the handover, and every referenced quality, compatibility,
release, runner, evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

If and only if a higher exact stable OpenCensus Proto version is qualified and
its complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/census-instrumentation/opencensus-proto@<qualified-version>`
for one dependency-only commit. Do not hand-edit metadata or use tidy as
implementation. Preserve every retained version, the Go 1.18 directive,
toolchain Go 1.26.7, production source, quality apparatus, and release input.
Stop instead of applying an unexplained multi-selection move.

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

After the OpenCensus Proto decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
