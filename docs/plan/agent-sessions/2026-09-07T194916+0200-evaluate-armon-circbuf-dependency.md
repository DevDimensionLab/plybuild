# Agent Session: Evaluate Armon Circbuf Dependency

Status: NEXT
Session ID: `2026-09-07T194916+0200-evaluate-armon-circbuf-dependency`
Created: `2026-09-07T19:49:16+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d9b7f88352bd85c10f926385df6e8284e3f62dda7ff03d8af14a7489fc2d1dc4`
Previous: [2026-09-07T183413+0200-evaluate-antihax-optional-dependency.md](2026-09-07T183413+0200-evaluate-antihax-optional-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e` as one bounded
dependency group. Resolve canonical latest, release qualification, default-
branch history, and the highest Go-1.18-floor-compatible candidate from primary
evidence. Make an exact dependency selection only if it changes a selected
version, preserves the retained floor through the complete minimal closure,
and passes every quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Antihax Optional v1.0.0. All earlier
rejections, no-change decisions, accepted closures, and evidence corrections
remain final. Do not revisit Optional or combine another module group. P8
remains queued.

Project MVS selects Circbuf's 2015 pseudo-version although it is not an explicit
`go.mod` requirement. A fresh post-Optional survey finds an empty exact stable
proxy list. Proxy and exact Go `@latest` and `@master` resolve the newer
pseudo-version `v0.0.0-20190214190532-5111143e8da2` at
2019-02-14T19:05:32Z without a reported Go declaration; exact `@v0` has no
matching version. The authoritative repository currently reports public,
enabled, unarchived, non-fork status, default branch `master`, zero tags, and
zero GitHub Releases. Treat canonical qualification, source identity,
signatures, release and alternate-path history, complete closure, self-tests,
loaded consumers, actual symbols, and vulnerability effect as unknown until
independently resolved. Do not call an untagged pseudo-version a stable release.

# Measurements At Start

Latest dependency implementation remains Repr commit
`6f4d02eb9c86ec2df8a488a85ed973aababe1f38`, exact parent
`53475076e1c79f6d2181877e3238a1a5d389c246`, tree
`82e7b1f5c503659082207481b8339e8113e38c10`, changing only `go.mod` and
`go.sum`. Operator-authorized lifecycle repair
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
that implementation and direct-child Repr handoff
`342c7ece82eae3cd726aa658462089cfb876fbe0`. Preserve it.

Optional v1.0.0 was retained without dependency edits. It is the exact path's
sole stable release, canonical latest, and highest stable Go-1.18-floor-
compatible version. Its lightweight tag resolves commit
`c3f0ba9c1a592b971d66b2787679af55b5c58f21`, tree
`b9328a8aa4526004bb36928dbc136c3acb8eec3a`, at
2019-10-10T23:37:20Z. The later master pseudo-version
`v1.0.1-0.20220101210036-407d38fabb55` is unreleased and changes only README
material; every Go file and `go.mod` is identical. Both declare Go 1.13 with
no requirements. The release's complete closure is Optional alone and passes
complete, repeated, and race-enabled tests plus vet in both proxy and exact-tag
forms, without source mutation.

Exact selected-version get changes no module selection and only projects a
redundant indirect requirement, main graph edge, and full checksum. Therefore
no dependency implementation commit was created. Optional is unloaded; Ply has
no import and `go mod why` says the main module does not need it. The historical
grpc-gateway example consumer and an explicitly external Go-1.18 fixture cover
the actually referenced wrapper types and methods. Untouched and projected
repository verify/build/tests/race/vet, Windows build, pinned lint, help,
API/CLI compatibility, CLI surface, and clean full preflight pass.

Optional's 540-entry evidence-manifest SHA-256 is
`345a29c62bfadf332f14aa234b1265a9dc9142739d212534d83ace531fb39afa`;
decision-summary SHA-256 is
`945b7ca8fe9d9faa605a6a80d8d416655925f95cd727071babf813bcd4caa27b`.
Preserve its separately warmed v1.0.1 compatibility cache, exact source-
manifest and vulnerability normalization, scratch-local GnuPG keyring and
BSD-mktemp/preflight corrections. They changed no repository file.

Current accepted measurements remain 234 selected modules, 3,580 graph edges,
429 native complete-test packages, 1,043 `go.sum` lines, a 356-line unapplied
tidy projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
vulnerability populations. The main module retains Go 1.18 and toolchain Go
1.26.7. Relative to accepted go-cmp commit c314bcb, accepted metadata adds
exactly 27 checksum lines. Ordinary and ignored status must be empty.

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
Treat rebuilt golangci-lint and govulncheck hashes as nonportable receipts;
prove versions and functionality.

# Role And Boundaries

From fresh external archives and caches, resolve Circbuf versions through the
Go proxy, checksum database, authoritative upstream repository, and primary Go
vulnerability data. Record exact selected and candidate commits/times, module
Go declarations and requirements, checksum pairs, source identity, tag and
commit signature status, release history, and archived/deprecated state.
Explicitly distinguish stable versions, prereleases, pseudo-versions,
retractions, forks, alternate module paths, branch heads, and unreleased
commits. Do not treat the newer pseudo-version as a stable release without
authoritative release evidence.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete closure. Measure old versus
candidate modules, graph edges, complete packages, checksums, loaded packages
and paths, explicit exact-get diff, and `go mod tidy -diff`. Explain every
selection, edge, and checksum change. Identify real consumers and exercise the
actually used symbols.

Require candidate module complete, repeated, and race-enabled self-tests plus
vet; repository verify/build, complete tests, race, vet, Windows build, pinned
lint, byte-identical public help, API/CLI reports, and primary vulnerability
identity. Preserve the retained Go 1.18 floor across the complete changed
closure. Keep standalone test-only closure separate from project MVS.

Reject or retain without dependency edits if canonical qualification, source
identity, floor compatibility, closure, module tests, loaded behavior, or any
applicable quality contract fails. Evaluate an unreleased canonical pseudo-
version under the recorded dependency policy without relabeling it as stable.
Do not manufacture metadata when exact get changes no selected version. Record
precise primary old/candidate vulnerability IDs and traces.

# Required Reading

Verify branch, clean ordinary and ignored status, exact ancestry, reciprocal
archive history, `./codex-dev-start.sh --check`, and P7/P8 state before work.
Read this archive, the rolling handover, P7 roadmap, `go.mod`, `go.sum`, the
answered Optional, Template, Colour, Chroma, Kong, Repr, Assert, and Units
archives, retained Kingpin/Resty/Errgo/Check/YAML outcomes, and the toolchain,
quality, baseline, compatibility, snapshot/Docker, acceptance, audit, and
lifecycle contracts. Do not reopen earlier decisions.

# Three Moves

If and only if a qualified candidate changes the selected version and its
complete minimal closure preserves Go 1.18, use exact Go 1.26.7 and exact
`go get github.com/armon/circbuf@<candidate>` for one dependency-only commit.
Do not hand-edit metadata or use tidy as implementation. Preserve every
retained version, the toolchain declarations, production source, quality
apparatus, and release inputs.

After a changed selection, run the complete P7 dependency gate: module and
consumer tests, graph/path/checksum/tidy explanation, repository verify/build/
tests/race/vet/Windows/pinned lint, help/API/CLI, launcher and Make contracts,
preflight, host plus fresh snapshot/Docker meta and acceptance, audit meta,
focused and exact Q0-Q2 audits, separate full audit, vulnerability comparison,
empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 rows at L2 and zero held, regressed,
not-comparable, or dirty counts. Full audit may exit 1 only for established
queued L3 rows, never 2.

Put every disposable cache, projection, source, report, generated artifact,
evidence tree, and build context beneath `$CODEX_SESSION_SCRATCH_ROOT`. Never
create direct `/private/tmp/ply-*` roots, never bypass scratch cleanup, and
never run `go mod download all` inside a measured tree. Warm only exact needed
closures from a separate scratch archive. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the Circbuf decision, rewrite the rolling handover and roadmap, answer
this archive, create exactly one reciprocal NEXT archive for the next single
P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
