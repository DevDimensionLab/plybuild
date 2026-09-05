# Agent Session: Evaluate Golang X Text Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-05T175806+0200-evaluate-golang-x-text-dependency`
Created: `2026-09-05T17:58:06+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5c7043f310fa91685443274ba9f6be935da99021cbe88bedceff8a57486cd5ee`
Previous: [2026-09-05T170519+0200-evaluate-golang-x-net-dependency.md](2026-09-05T170519+0200-evaluate-golang-x-net-dependency.md)
Next: [2026-09-05T184249+0200-evaluate-gopkg-yaml-v2-dependency.md](2026-09-05T184249+0200-evaluate-gopkg-yaml-v2-dependency.md)
Outcome: rejected floor-compatible v0.15.0; retained v0.7.0 because mandatory module tests and vet fail

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`golang.org/x/text v0.7.0` as one bounded dependency group. Resolve the
canonical latest stable release and highest floor-compatible candidate from
primary evidence. Implement one exact changed selection only if it preserves
the retained Go 1.18 floor, has an explained minimal closure, and passes every
quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and
completed dependency groups through retained x/net v0.7.0. All earlier
recorded rejections, no-change decisions, accepted closures, and evidence
corrections remain final. Do not revisit them or combine another module group.
P8 remains queued.

The current build list selects golang.org/x/text v0.7.0. Treat its latest
release, Go declaration, release qualification, closure, loaded population,
consumer paths, source history, tests, and vulnerability effect as unknown
until independently resolved. This session may change only x/text's exact
required go.mod/go.sum metadata, its minimal MVS closure, and the roadmap/
handoff record. Do not change production Go, another dependency, language or
toolchain declarations, quality apparatus, Docker/release inputs, packaging,
publishers, or P8 code.

# Measurements At Start

Latest implementation remains dependency-only cast commit
`cf4fd4933251b2ca7844b6420fd3cf2a86f049a5`, exact parent
`17f7277fd2512ba06db76d1acb0af5f8b623243c`, clean tree
`485690c010cef3b02385d8f2471d3ef53611abe9`. The x/net documentation
handoff must have exact parent
`e5ea7ae4ddb5bc4d58085ff3fcfe5b472e6951ce`. Relative to accepted go-cmp
commit c314bcb, accepted metadata changes remain go-colorful v1.2.0 -> v1.4.1,
go-runewidth v0.0.14 -> v0.0.17, go-toml/v2 v2.0.7 -> v2.2.2, and cast
v1.5.0 -> v1.5.1, with their recorded minimal closures and exactly 15 added
checksum lines. Ordinary and ignored status must be empty.

Current dependency measurements remain 234 selected modules, 3,564 graph
edges, 429 native complete-test packages, a 332-line unapplied tidy projection,
and exact Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations
20/30/20. The retained main module declares Go 1.18 and prefers toolchain Go
1.26.7.

Use exact Go 1.26.7 at
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in PATH, keep GOENV=off, GOWORK=off,
GOTOOLCHAIN=local, and inject no ambient GOFLAGS. Recovery evidence at
`/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` retains its verified two-entry
manifest SHA-256
`1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
Retain the verified golangci-lint 2.12.2, GoReleaser 2.17.1, apidiff, and
govulncheck v1.7.0 binaries and hashes recorded in the handover.

X/net rejection evidence is fully verified at
`/private/tmp/ply-p7-x-net-selection.e5ea7ae.cP57X8`, 71,887 entries and
manifest SHA-256
`88e4be87483322792ad3da0064bcef4ced5ee21c52821bb1f80c143f80a3aaf1`;
decision-summary SHA-256 is
`6e32698abcde67604e159488ba21f6e780fd61dfba2ed9299644b5f011940dff`.
V0.58.0 is canonical x/net latest and v0.25.0 is the highest complete
Go-1.18-floor-compatible release, but v0.25.0 complete tests and race fail
under Go 1.26.7. X/net v0.7.0 remains selected and dependency metadata is
unchanged. Preserve all prior verified roots and the recorded go-colorful
mutable telemetry, btree regression, and cast whitespace-path corrections.

# Role And Boundaries

From fresh external archives and caches, resolve x/text releases through the
Go proxy, checksum database, authoritative upstream repository and mirror,
and primary Go vulnerability data. Record exact tag commits/times, module Go
declarations and requirements, checksum pairs, source identity, tag and commit
signature status, relevant release history, and archived/deprecated status.
Explicitly distinguish stable releases, prereleases, retractions, forks,
nested or later major paths, and unreleased upstream commits. Prove canonical
latest and the highest stable release compatible with Go 1.18 from declarations
and complete changed closure; do not infer compatibility from a modern build.

Measure old versus candidate selected modules, graph edges, complete package
population, checksums, loaded packages and paths, explicit exact-get diff, and
`go mod tidy -diff`. Require an explained minimal selection/edge/checksum
closure. Independently prove whether x/text is loaded; if it is, identify its
real consumers and run focused behavior over the actually used packages and
symbols. Also require candidate module complete tests, repository build,
complete tests/race/vet, pinned lint, byte-identical public help, identical
API/CLI reports, and exact Darwin and Windows vulnerability populations.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity, module
self-tests, loaded behavior, or any repository quality contract fails. If the
selected release is already the exact floor-compatible decision and exact get
is a no-op, record that no-change decision without manufacturing a dependency
commit. Do not assume vulnerability identity; prove exact old/candidate IDs
and traces.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, P7 in
the roadmap, go.mod/go.sum, the rejected x/net, x/image, and gotenv decisions,
retained jwalterweatherman decision, accepted cast decision and quality
evidence, rejected afero decision, prior bounded dependency decisions, and the
toolchain, compatibility, snapshot/Docker, quality, baseline-reproduction, and
audit contracts. Preserve every recorded manifest correction.

# Three Moves

Only if every decision gate passes and the selection changes, use exact Go
1.26.7 and exact `go get golang.org/x/text@<selected-version>` for one
dependency-only commit. Do not hand-edit module metadata and do not use tidy as
implementation. Preserve every retained dependency selection, language/
toolchain declaration, production source, quality apparatus, and release input.

After a changed selection, run the complete P7 dependency gate: focused
behavior, graph/path, tests/race/vet, pinned lint, help/API/CLI, launcher and
Make contracts, preflight, host plus fresh snapshot/Docker meta and acceptance,
audit meta, focused and exact Q0-Q2 audits, separate full audit, vulnerability
comparison, empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
make quality must exit 0 with all 27 rows at L2 and zero held, regressed,
not-comparable, or dirty counts. The full audit may exit 1 only for established
queued L3 rows, never 2.

Keep all caches, projections, reports, generated artifacts, build contexts,
schema-2 evidence, and audit output outside the worktree. Warm caches from a
separate external Git archive; never run `go mod download all` inside a
measured tree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next measured P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
that next group, launch a successor, push, merge, publish, release, stash,
revert, delete retained evidence/images, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Rejected `golang.org/x/text v0.15.0` and retained selected v0.7.0 without
editing dependency metadata. Fresh primary evidence resolves v0.41.0 as the
canonical latest stable release and v0.15.0 as the highest stable release
whose complete changed closure preserves the retained Go 1.18 floor. The
candidate then fails the mandatory complete module test and vet stop gates
under exact Go 1.26.7, so repository-wide candidate gates and an implementation
commit are inapplicable.

The Go proxy enumerates 49 stable semantic releases from v0.1.0 through
v0.41.0, with no prereleases or retractions. V0.41.0 is authoritative tag
commit `acdba6655fd45cdb5ab73c9d6a8981333bd65a39` at
2026-08-11T15:22:47Z and declares Go 1.25.0. Releases v0.14.0-v0.22.0 directly
declare Go 1.18, but v0.16.0 first requires x/tools pseudo-version
`v0.21.1-0.20240508182429-e35e4ccd0d2d`, which declares Go 1.19; that
above-floor requirement remains through v0.22.0. V0.23.0 and later directly
declare above Go 1.18. V0.15.0 instead requires x/tools v0.6.0, x/mod v0.8.0,
and x/sys v0.5.0, all declaring Go 1.18 or lower. It is therefore the complete
closure boundary, not merely the highest directly compatible declaration.

V0.15.0 is unsigned lightweight tag commit
`8d533a0c40adec778a7d09ac6c8aa640d3c883f4` at
2024-04-15T18:14:38Z. Its checksum pair is
`h1:h1V/4gjBv8v9cjcR6+AR5+/cIYK5N/WAgiv4xlsEtAk=` /
`h1:18ZOQIKpY8NJVqYksKHtTdi31H5itFRjB5/qKTNYzSU=`. All 542 proxy files
match the exact authoritative tag; ZIP SHA-256 is
`13faee7e46c8a18c8a28f3eceebf15db6d724b9a108c3c0482a6d2e58ba73a73`.
Selected v0.7.0 is unsigned lightweight tag commit
`71a9c9afc4cd710b9412f7f99f0d8e35b10e488a` at
2023-01-31T16:01:06Z. Its checksum pair is
`h1:4BRB4x83lYWy72KwLD/qYDuTu7q9PjSagHvijDw7cLo=` /
`h1:mrYo+phRRbMaCq/xk9113O4dZlRixOauAjOtrjsXDZ8=`; all 530 proxy files
match, with ZIP SHA-256
`4d017493c58addadf3c753056b921b47ae386a4cfd10eab2d90ed1252c6ba0e4`.
The authoritative Go repository and unarchived, undisabled, non-fork GitHub
mirror have identical master and all tag refs. No nested or later major module
path exists. The three historical GitHub Release objects are stable; higher
releases are qualified by proxy and matching authoritative tags. Master
`f53c31601f90c1b840c0703c5537c3ce1e4b6f5c` is 12 commits beyond v0.41.0,
declares Go 1.26.0, and is unreleased.

Exact projected `go get golang.org/x/text@v0.15.0` changes only x/text in
go.mod and adds exactly the candidate checksum pair. MVS additionally changes
x/mod from pseudo-version `v0.6.0-dev.0.20220419223038-86c51ed26bb4` to
v0.8.0 and x/tools v0.1.12 to v0.6.0; selected x/sys v0.30.0 dominates the
candidate's v0.5.0 requirement, and x/sync remains v0.1.0. Modules stay 234,
complete loaded test packages stay byte-identical at 429, graph edges change
3,564 -> 3,567 through one removal and four additions, go.sum projects 1,031
-> 1,033 lines, and the unapplied tidy projection changes 332 -> 333 lines.
This is the explained minimal selection, edge, checksum, and metadata closure.

Three x/text packages load: `runes`, `transform`, and `unicode/norm`. The real
path is `plybuild/cmd -> spf13/viper -> spf13/afero -> x/text/runes`; Afero
constructs `transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)),
norm.NFC)` and calls `transform.String`. Focused loaded-package, Afero, and Ply
cmd tests pass at count 10 in both projections with normalized-identical
output.

Candidate source remains byte-identical and module verification passes, but
complete `go test ./... -count=1`, `-count=10`, race, and vet all exit 1. The
default and race runs reject three stale Example identifiers and the
message/pipeline test panics in x/tools v0.6.0's SSA builder on a Go 1.26
range-over-function construct. Repeated tests additionally expose persistent
`cases.TestShortBuffersAndOverflow` failures. Vet reports the same three
Example errors plus 20 unkeyed literals, four unreachable-code diagnostics,
and one unused `currency.Unit.String` result. Selected v0.7.0 reproduces the
same failure classes. The Example-name repair first ships in v0.18.0, whose
closure already requires the Go-1.19 x/tools pseudo-version, so it cannot
repair the retained floor line.

Govulncheck v1.7.0 preserves exact old/candidate Darwin-symbol,
Darwin-module, and Windows-symbol populations at 20/30/20. Both versions have
the same sole x/text module finding, GO-2026-5970, fixed in v0.39.0; neither
symbol scan has an x/text-reachable finding or trace. The fresh 1,392-entry
primary index contains four x/text records, while the other three were fixed
before v0.7.0. No vulnerability effect overrides the mandatory module gate.

Evidence is sealed at
`/private/tmp/ply-p7-x-text-selection.6ab5945.3qn5Bj`. Its fully verified
82,629-entry manifest SHA-256 is
`48c872da4ab796ef1115003fcf0a226a07bf5791a163d36bc9c946c9d9f92ed6`;
decision-summary SHA-256 is
`50d6bee90ea7cea060e23400e820b1537742ec2167bee04f188070e2e5eb39c9`.
All 27 inherited manifests were independently verified before measurement,
with the recorded go-colorful mutable telemetry, btree regression, and cast
whitespace-path corrections preserved. `go.mod`, `go.sum`, production code,
quality apparatus, and release inputs remain unchanged, and there is no
dependency implementation commit.
