# Agent Session: Evaluate Alecthomas Chroma Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-07T133322+0200-evaluate-alecthomas-chroma-dependency`
Created: `2026-09-07T13:33:22+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `a7ff3e9fa68586cf4a197c44b39585798d189ccec160680d9a89ab80f495bae4`
Previous: [2026-09-07T125545+0200-evaluate-alecthomas-kong-dependency.md](2026-09-07T125545+0200-evaluate-alecthomas-kong-dependency.md)
Next: [2026-09-07T164119+0200-evaluate-alecthomas-colour-dependency.md](2026-09-07T164119+0200-evaluate-alecthomas-colour-dependency.md)
Outcome: Retained canonical-latest stable `github.com/alecthomas/chroma v0.10.0` without dependency metadata edits: exact selected-version get is a true no-op, while mandatory proxy/tag vet and the loaded consumer's project-MVS code-block golden contract fail reproducibly.

# Answer

The exact `github.com/alecthomas/chroma` proxy path lists 30 stable tags from
v0.1.0 through v0.10.0, with no prerelease or retraction. Proxy `@latest` and
exact Go `@latest`/`@v0` resolve canonical-latest stable v0.10.0 at
2022-01-12T10:49:38Z. Its lightweight tag is unsigned commit
`36bdd4b98823bd1d7be96767cde3dd575e60b406`, tree
`d27058989b845352d49b45e7538f7b0004c8d651`, parent
`b01c8fcab6f3d79534dadbc8fabd4c54cabaf4d8`. GitHub has a non-draft,
non-prerelease Release object for v0.10.0; 26 of the 30 v0 tags have Release
objects. The enabled, unarchived, undisabled, non-fork repository has no v0
deprecation marker.

Exact v0-path `@master` resolves the later unreleased pseudo-version
`v0.10.1-0.20220126230913-d491f1b5c1d2`, unsigned commit
`d491f1b5c1d2e20b85309ccf34fc778d950dc3f2`, tree
`e5bcdbd2d1db14d16da0df023888dd0ba849f8fe`, at
2022-01-26T23:09:13Z. The three post-v0.10.0 commits add a style, aliases, and
helper scripts; they were never tagged as a v0 release. The next commit changes
the module path to `github.com/alecthomas/chroma/v2`. The alternate `/v2`
line has 38 stable releases and four v2.0.0 alphas; its canonical latest is
v2.27.0 at unsigned lightweight-tag commit
`a6d00fe2cdfc88da0b91396e577da16c75c9c7fb`, declaring Go 1.25. Current
default-branch master instead declares the distinct `/v3` path and carries v3
alpha tags. Neither the v0 pseudo-version, `/v2`, nor `/v3` is an in-place
stable update of the selected exact path.

Selected v0.10.0 declares Go 1.13 and requires regexp2 v1.4.0 plus the
test-only testify v1.7.0 and indirect go-spew v1.1.1. Its standalone complete
test graph selects eight modules across nine edges: testify also selects
go-difflib v1.0.0, objx v0.1.0, YAML v3's 2020 pseudo-version, and Check's
2016 pseudo-version. Only Chroma and testify declare a Go version, both 1.13;
the rest have no Go directive. This complete standalone closure preserves the
retained Go 1.18 floor. Project MVS already selects Chroma v0.10.0 and later
regexp2/testify versions, so the changed-selection closure is empty.

The selected checksum pair is
`h1:7XDcGkCQopCNKjZHfYrNLraA+M7e0fMiJ/Mfikbfjek=` /
`h1:jtJATyUxlIORhUOFNA9NZDWGAQ8wpxQQqNSB4rjA/1s=`; independent checksum-
database lookup agrees. Its proxy ZIP SHA-256 is
`beb07b996ee33bc052fe039c93d1c0726e61bcc4819ca39f7bf63304f2ae8c49`.
All 696 proxy regular files match the exact tag commit. The proxy correctly
omits eight regular files in nested `cmd/chroma` and `cmd/chromad` modules and
seven tracked Hermit symlinks. The pseudo-version pair is
`h1:Gg09t2u+C08At6TYucNrD3Cbaq97SUHax84BzQwRTgU=` / the same go.mod sum;
its ZIP SHA-256 is
`033bd0de6588cb5c9508b9ced1eb0a0620eb9861fc43a11843d416ca3e245d79`,
and all 699 proxy regular files match its exact commit with the same justified
nested-module and symlink omissions.

Exact `go get github.com/alecthomas/chroma@v0.10.0` in an external clean
replay exits zero and changes no byte, requirement, selected version, edge, or
checksum. Old and replayed states are byte-identical at 234 selected modules,
3,580 graph edges, 429 complete packages, and 1,043 go.sum lines. The
unapplied `go mod tidy -diff` remains the same 356 lines. There is therefore no
authorized dependency implementation to manufacture.

Chroma is real loaded code rather than historical graph debt. Exactly 33
Chroma packages load through `plybuild/cmd -> go-term-markdown -> chroma`, and
`go mod why -m` confirms that path. Repository source has no direct Chroma
import. Go-term-markdown imports Chroma, `formatters`, `lexers`, and `styles`
and uses `Lexer`, `Coalesce`, `Formatter`, `lexers.Get`, `lexers.Analyse`,
`lexers.Fallback`, `formatters.Fallback`, `formatters.TTY8`,
`styles.Pygments`, `Lexer.Tokenise`, and `Formatter.Format` when rendering
code blocks.

The existing project Markdown reference-definition contract passes ten times.
The direct consumer's focused `TestRender/codeblock` golden, however, fails
all ten project-MVS repetitions: rendering completes, but selected Chroma
v0.10.0 together with the project's selected color closure produces different
ANSI intensity/color sequences and one-space code padding from the consumer's
v0.7.1-era golden. Its 613-line failure report SHA-256 is
`b6b6f24dcd9c1b97c9280e01b4277aa3614b5513e35bd8ad46e6cd7c9cea0d3c`.
This is an independently exercised loaded-behavior stop condition, not an
assertion that the alternate `/v2` path can replace v0 in place.

Writable proxy and exact-tag sources independently verify, list the same 34
packages, and pass complete count-1, count-10, and race tests without source
mutation. Mandatory `go vet ./...` exits 1 in both with the same normalized
7,175 diagnostics across 206 lexer files, all unkeyed
`github.com/alecthomas/chroma.Rule` literals. The sorted diagnostic SHA-256 is
`cbfbdc99bce5e357f83b20eec1a08644ddb86a16ed90f3a2c176cb501a68604e`.
This second stop-rule failure makes downstream changed-selection lint,
snapshot/Docker, exact-quality, and audit gates inapplicable.

Fresh govulncheck v1.7.0, built with Go 1.26.7, reports primary data updated
2026-09-02T19:12:04Z. Its 1,392-record primary module index has no Chroma
record. Exact old/replay IDs and normalized traces are identical at 20 unique
Darwin symbol IDs, 30 Darwin module IDs, and 20 unique Windows symbol IDs;
there is no Chroma finding or trace. No dependency metadata or implementation
commit was created.

Decision evidence beneath the active scratch root has 4,335 verified entries.
Evidence-manifest SHA-256 is
`bbd6a41e0cff3248bacc2fef6630ed3679a3299484cb21446fc2841a39a742d5`;
decision-summary SHA-256 is
`476f926b06cd8fbfce8ac22c9e0ed1102c8f5c8f8e2070bcb1636a825225927f`.
Exact Go 1.26.7 retained binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
The rebuilt govulncheck binary SHA-256
`30e6f51fc7d968d10f0d789c82dce56ff0da880c7dfed7c82e08e6feaa554f24`
is a nonportable receipt; its reported version, builder, database identity,
and functional scans are the proof.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/alecthomas/chroma v0.10.0` as one bounded dependency group.
Resolve canonical latest, release qualification, alternate-path history, and
the highest Go-1.18-floor-compatible candidate from primary evidence. Make an
exact dependency selection only if it changes a selected version, preserves
the retained floor through the complete minimal closure, and passes every
quality contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and bounded
dependency decisions through rejected Kong v1.3.0. All earlier rejections,
no-change decisions, accepted closures, and evidence corrections remain final.
Do not revisit Kong or combine another module group. P8 remains queued.

The build list selects Chroma v0.10.0 as an indirect `go.mod` requirement. A
fresh post-Kong survey says exact-path proxy `@latest` remains stable v0.10.0
at 2022-01-12T10:49:38Z with Go 1.13, while exact `@master` resolves
unreleased pseudo-version `v0.10.1-0.20220126230913-d491f1b5c1d2` at
2022-01-26T23:09:13Z. Treat qualification, source identity, signatures,
release history, alternate `/v2` lineage, closure, self-tests, and
vulnerability effect as unknown until independently resolved. Thirty-three
Chroma packages load through `plybuild/cmd -> go-term-markdown -> chroma`,
while repository Go source has no direct Chroma import; treat those only as
starting observations.

# Measurements At Start

Latest dependency implementation remains Repr commit
`6f4d02eb9c86ec2df8a488a85ed973aababe1f38`, exact parent
`53475076e1c79f6d2181877e3238a1a5d389c246`, tree
`82e7b1f5c503659082207481b8339e8113e38c10`, changing only `go.mod` and
`go.sum`. Operator-authorized lifecycle repair
`f9f0f7669e635c8c7bb169aab816d0ebe1b16435` remains intentionally between
that implementation and direct-child Repr handoff
`342c7ece82eae3cd726aa658462089cfb876fbe0`. Preserve it.

Kong v1.3.0 was rejected before implementation. Canonical latest v1.16.1 and
every stable release from v1.4.0 declare Go 1.20. V1.3.0's complete closure
preserves Go 1.18, but proxy and exact-tag sources both fail mandatory vet
with the same 77 malformed-struct-tag diagnostics. No Kong metadata changed.
Its 211-entry evidence manifest SHA-256 is
`6148c409545f453a78ffdc3a0a10934b75e082c1681f2805c29c8253a23d01ae`;
decision-summary SHA-256 is
`028c0fd0025a5f5de433b65d63a5258e18c6e959b440446e9834b03925f98b36`.

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

From fresh external archives and caches, resolve Chroma versions through the
Go proxy, checksum database, authoritative upstream repository, and primary Go
vulnerability data. Record exact tag and pseudo-version commits/times, module
Go declarations and requirements, checksum pairs, source identity, tag and
commit signature status, release history, and archived/deprecated state.
Explicitly distinguish stable tags, prereleases, pseudo-versions, retractions,
forks, branch heads, the `/v2` module, and unreleased commits. Do not treat a
pseudo-version or alternate module path as an in-place stable update.

Prove canonical latest and the highest qualified exact-path version compatible
with Go 1.18 from declarations and the complete closure. Measure old versus
candidate modules, graph edges, complete packages, checksums, loaded packages
and paths, explicit exact-get diff, and `go mod tidy -diff`. Explain every
selection, edge, and checksum change. Identify real consumers and exercise the
actually used Chroma packages and symbols, including renderer behavior.

Require candidate module complete tests, repeated tests, race, and vet;
repository build, complete tests/race/vet, pinned lint, byte-identical public
help, identical API/CLI reports, and exact Darwin and Windows vulnerability
populations. Keep module test-only requirements separate from project MVS.

Stop and record rejection without dependency edits if canonical resolution,
floor compatibility, exact closure, source identity, module self-tests, loaded
behavior, or any repository quality contract fails. If selected v0.10.0 is
already the exact floor-compatible decision and exact get changes no selected
version, record no change without manufacturing metadata or a dependency
commit. Prove exact old/candidate vulnerability IDs and traces.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and P7/P8 checkpoint before editing. Read
this archive, rolling handover, P7 roadmap, go.mod/go.sum, answered Kong, Repr,
Assert, and Units decisions, retained Kingpin/Resty/Errgo/YAML decisions,
earlier accepted and rejected bounded dependencies, and toolchain,
compatibility, snapshot/Docker, quality, baseline-reproduction, and audit
contracts. Preserve every recorded manifest and setup correction.

# Three Moves

Only if every decision gate passes and a version selection changes, use exact
Go 1.26.7 and exact
`go get github.com/alecthomas/chroma@<selected-version>` for one dependency-
only commit. Do not hand-edit metadata or use tidy as implementation. Preserve
every retained selection, especially Repr v0.5.4, Assert v1.0.0, Units'
2024 pseudo-version, the selected Kong pseudo-version, Kingpin v2.2.6, Resty
v1.12.0, Errgo v2.1.0, Check's 2019 pseudo-version, all three retained YAML
paths, language/toolchain declarations, production source, quality apparatus,
and release input.

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
