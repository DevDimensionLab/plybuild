# Agent Session: Evaluate Alecthomas Colour Dependency

Status: NEXT
Session ID: `2026-09-07T164119+0200-evaluate-alecthomas-colour-dependency`
Created: `2026-09-07T16:41:19+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `4acac96a4efd6d709997ab9cb5b2cc14649985348db0c25107fdf3dda28d1c2f`
Previous: [2026-09-07T133322+0200-evaluate-alecthomas-chroma-dependency.md](2026-09-07T133322+0200-evaluate-alecthomas-chroma-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/alecthomas/colour v0.1.0` as one bounded dependency group. Resolve
canonical latest, release qualification, alternate-path history, and the
highest Go-1.18-floor-compatible candidate from primary evidence. Make an
exact dependency selection only if it changes a selected version, preserves
the retained floor through the complete minimal closure, and passes every
quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through retained Chroma v0.10.0. All earlier rejections,
no-change decisions, accepted closures, and evidence corrections remain final.
Do not revisit Chroma or combine another module group. P8 remains queued.

The build list selects Colour v0.1.0 as an indirect `go.mod` requirement. A
fresh post-Chroma survey finds only stable v0.1.0 in the exact-path proxy list.
Proxy `@latest` and exact Go `@latest`, `@v0`, and `@master` all resolve
v0.1.0 at 2019-11-01T02:47:59Z; Go reports no module Go declaration. Treat
qualification, source identity, signatures, release and branch history,
closure, self-tests, loaded consumers, actual symbols, and vulnerability
effect as unknown until independently resolved.

# Measurements At Start

Latest dependency implementation remains Repr commit
`6f4d02eb9c86ec2df8a488a85ed973aababe1f38`, exact parent
`53475076e1c79f6d2181877e3238a1a5d389c246`, tree
`82e7b1f5c503659082207481b8339e8113e38c10`, changing only `go.mod` and
`go.sum`. Operator-authorized lifecycle repair
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
that implementation and direct-child Repr handoff
`342c7ece82eae3cd726aa658462089cfb876fbe0`. Preserve it.

Chroma v0.10.0 was retained without dependency edits. It is the exact-path
canonical latest stable release, and exact get is a byte-for-byte no-op with
an empty changed-selection closure. Its complete standalone closure preserves
Go 1.18, but proxy and exact-tag sources both fail mandatory vet with the same
normalized 7,175 unkeyed-Rule-literal diagnostics across 206 lexer files. The
loaded go-term-markdown code-block golden also fails all ten project-MVS
repetitions. Chroma's 4,335-entry evidence-manifest SHA-256 is
`bbd6a41e0cff3248bacc2fef6630ed3679a3299484cb21446fc2841a39a742d5`;
decision-summary SHA-256 is
`476f926b06cd8fbfce8ac22c9e0ed1102c8f5c8f8e2070bcb1636a825225927f`.

Current measurements remain 234 selected modules, 3,580 graph edges, 429
native complete-test packages, 1,043 go.sum lines, a 356-line unapplied tidy
projection, and exact 20/30/20 Darwin-symbol/Darwin-module/Windows-symbol
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

From fresh external archives and caches, resolve Colour versions through the
Go proxy, checksum database, authoritative upstream repository, and primary Go
vulnerability data. Record exact tag and branch-head commits/times, module Go
declarations and requirements, checksum pairs, source identity, tag and commit
signature status, release history, and archived/deprecated state. Explicitly
distinguish stable tags, prereleases, pseudo-versions, retractions, forks,
alternate module paths, branch heads, and unreleased commits. Do not treat an
unreleased commit or alternate module path as an in-place stable update.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete closure. Measure old versus
candidate modules, graph edges, complete packages, checksums, loaded packages
and paths, explicit exact-get diff, and `go mod tidy -diff`. Explain every
selection, edge, and checksum change. Identify real consumers and exercise the
actually used Colour packages and symbols, including terminal renderer
behavior where applicable.

Require candidate module complete tests, repeated tests, race, and vet;
repository build, complete tests/race/vet, pinned lint, byte-identical public
help, identical API/CLI reports, and exact Darwin and Windows vulnerability
populations. Keep module test-only requirements separate from project MVS.

Stop and record rejection without dependency edits if canonical resolution,
floor compatibility, exact closure, source identity, module self-tests, loaded
behavior, or any repository quality contract fails. If selected v0.1.0 is
already the exact floor-compatible decision and exact get changes no selected
version, record no change without manufacturing metadata or a dependency
commit. Prove exact old/candidate vulnerability IDs and traces.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and P7/P8 checkpoint before editing. Read
this archive, rolling handover, P7 roadmap, go.mod/go.sum, answered Chroma,
Kong, Repr, Assert, and Units decisions, retained Kingpin/Resty/Errgo/YAML
decisions, earlier accepted and rejected bounded dependencies, and toolchain,
compatibility, snapshot/Docker, quality, baseline-reproduction, and audit
contracts. Preserve every recorded manifest and setup correction.

# Three Moves

Only if every decision gate passes and a version selection changes, use exact
Go 1.26.7 and exact
`go get github.com/alecthomas/colour@<selected-version>` for one dependency-
only commit. Do not hand-edit metadata or use tidy as implementation. Preserve
every retained selection, especially Repr v0.5.4, Assert v1.0.0, Units'
2024 pseudo-version, Chroma v0.10.0, the selected Kong pseudo-version, Kingpin
v2.2.6, Resty v1.12.0, Errgo v2.1.0, Check's 2019 pseudo-version, all three
retained YAML paths, language/toolchain declarations, production source,
quality apparatus, and release input.

After a changed selection, run the complete P7 dependency gate: focused
behavior, graph/path, tests/race/vet, pinned lint, help/API/CLI, launcher and
Make contracts, preflight, host plus fresh snapshot/Docker meta and acceptance,
audit meta, focused and exact Q0-Q2 audits, separate full audit, vulnerability
comparison, empty-HOME count-2, and final ordinary/ignored cleanliness. Exact
`make quality` must exit 0 with all 27 rows at L2 and zero held, regressed,
not-comparable, or dirty counts. Full audit may exit 1 only for established
queued L3 rows, never 2.

Put every disposable cache, projection, report, generated artifact, build
context, evidence tree, and audit output beneath
`$CODEX_SESSION_SCRATCH_ROOT`; never create direct `/private/tmp/ply-*` roots.
The launcher deletes scratch after every turn. Warm caches only from a separate
archive beneath that root and never run `go mod download all` inside a measured
tree. Retain only compact decisions, digests, and receipts in tracked docs.
Preserve the lifecycle repair and bounded-scratch policy. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Automatic Handoff

After the decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next measured P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
the next group, launch a successor, push, merge, publish, release, stash,
revert, bypass scratch cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
