# Agent Session: Evaluate Ghodss YAML Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-08T230017+0200-evaluate-ghodss-yaml-dependency`
Created: `2026-09-08T23:00:17+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `b5e5affb100636604c86f63e945cd5267a69dbd5e0f04c518115c65fd5d5b439`
Previous: [2026-09-08T210923+0200-evaluate-fsnotify-fsnotify-dependency.md](2026-09-08T210923+0200-evaluate-fsnotify-fsnotify-dependency.md)
Next: [2026-09-09T002351+0200-evaluate-go-gl-glfw-dependency.md](2026-09-09T002351+0200-evaluate-go-gl-glfw-dependency.md)
Outcome: Retained exact-path Ghodss YAML v1.0.0 without metadata changes; it is the sole stable release, preserves the Go 1.18 floor through its complete closure, is not loaded by Ply, and no unreleased or successor identity qualified as a replacement.

## Answer

### Decision And Exact Identity

Retain exact-path `github.com/ghodss/yaml v1.0.0` without editing `go.mod` or
`go.sum` and without manufacturing a dependency implementation commit. It is
both selected and the only proxy-listed stable release. No higher stable or
prerelease exists, and neither unreleased master nor the maintained distinct
module `sigs.k8s.io/yaml` may be promoted into this exact path.

Fresh proxy, checksum-database, go-import, Git, and GitHub evidence agrees:

- proxy list: exactly `v1.0.0`; `@latest` time
  `2017-03-27T23:54:44Z`;
- sumdb source/mod pair
  `h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=` /
  `h1:4dBDuWmgqj2HViK6kFavaiC9ZROes6MMH2rRYeMEF04=`;
- proxy ZIP SHA-256
  `c3f295d23c02c0b35e4d3b29053586e737cf9642df9615da99c0bda9bbacc624`,
  with all seven files byte-identical to the Git tag tree;
- go-import repository `https://github.com/ghodss/yaml.git`;
- lightweight unsigned tag and unsigned release commit
  `0ca9ea5df5451ffdf184b4428c902747c2c11cd7`, tree
  `252e285a136d503d2913ad39cb3b1669b9a76999`, parents
  `04f313413ffd65ce25f2541bfd2b2ceec5c0908c` and
  `a4f8cbd2fd05654d25f651b7e26d612ce3c98cc7`; and
- one GitHub release, non-draft and non-prerelease, published
  `2017-03-28T22:41:51Z` without assets.

The canonical public repository is enabled, unarchived, non-fork, and
defaults to `master`. The release is an ancestor of current master
`d8423dcdf3440d0a5baffc6f90a11e4128545620`, tree
`8df4d8facd33473e866d0f2e41e121c8732d5a6a`, parent
`1e4101787d1907800b0200eb38faa0f5041d8ee6`. There are 16 post-release
commits, but master resolves only as unreleased pseudo-version
`v1.0.1-0.20220118164431-d8423dcdf344`; it is not an exact stable upgrade.
There is no retraction, module deprecation, redirect, exact-path v2, or fork
identity. `sigs.k8s.io/yaml @latest v1.6.0` is a separately maintained fork
and module and was not substituted.

### Complete Floor And Source Quality

The v1.0.0 module file contains only `module github.com/ghodss/yaml`. Isolated
resolution under exact Go 1.26.7 and contained Go 1.18.10 proves the minimal
source closure is `gopkg.in/yaml.v2 v2.4.0`, which declares Go 1.15. Its own
historical test closure adds
`gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405`, with no Go directive
or further requirements. Native tests for all three components pass under
both SDKs when the historical test dependency is evaluated with Go 1.18
language semantics. The full closure therefore preserves the Go 1.18 floor.

The release package passes count-1, two independent count-10 runs, race, and
production cross-builds for Darwin/amd64, Linux/amd64, Linux/386,
Windows/amd64, and FreeBSD/amd64 under both SDKs. It is pure Go: no Cgo,
assembly, generated files, examples, testdata, fuzz tests, or property tests.
The release contains three Go files and one test file. Historical Travis CI
covered Go 1.3 and 1.4 with only test/build; master later covered Go 1.9-1.16
and has no current Actions workflow.

Two historical upstream source/test defects remain explicit. `go vet` under
both SDKs reports the unreachable return in `yaml.go`, and the v1.0.0 test
file does not compile for Linux/386 because it passes untyped `math.MaxInt64`
where `int` is required. The production package does compile for Linux/386,
and the first post-release commit fixed that test. Unreleased master is not a
quality upgrade: under Go 1.26 its native test command fails through vet
because `ExampleUnknown` names no identifier, and explicit vet still reports
the unreachable code. With no newer stable to select, these are retained,
characterized historical source/test limitations rather than concealed
candidate regressions.

### API And Conversion Semantics

V1.0.0 exports exactly `Marshal`, `Unmarshal`, `YAMLToJSON`, and
`JSONToYAML`; it does not export `UnmarshalStrict`. Unreleased master changes
Unmarshal to a variadic-options signature and adds `UnmarshalStrict`,
`YAMLToJSONStrict`, `JSONOpt`, and `DisallowUnknownFields`. Master strict mode
rejects duplicate keys, while unknown fields remain accepted unless the
explicit disallow option is supplied.

The release pipeline decodes YAML into interface values, recursively converts
map keys and scalar values into JSON-compatible forms, then passes through
`encoding/json`. Independent fixtures under both SDKs prove:

- ordinary string, integer, and boolean keys convert; composite and null keys
  return errors;
- integer, boolean, null, float, YAML 1.1 boolean, tag, anchor, and alias
  values convert as implemented;
- duplicate keys are last-wins, unknown struct fields are ignored, and only
  the first document of multi-document YAML is consumed;
- JSON-tagged and embedded structs, interfaces, custom JSON marshalers and
  unmarshalers, byte fields, scalar-to-string conversion, invalid inputs,
  ordering, escaping, and wrapped error propagation behave consistently;
- interface numbers become float64 after JSON unmarshal; JSONToYAML preserves
  number text and sorts ordinary map output; JSON subset input is canonicalized
  rather than passed through byte-for-byte;
- JSONToYAML accepts YAML-superset input because it begins with yaml.Unmarshal;
  and invalid `!!binary` bytes become Unicode replacement runes; and
- the shared struct-field cache is race-clean in 100-by-100 concurrent use.

Ordinary non-colliding conversion is deterministic, but the algorithm is
intentionally lossy at YAML/JSON boundaries. Distinct YAML keys `1` and
`"1"` both become JSON key `"1"`; repeated runs observe either surviving
value because the collision is resolved during Go map iteration. Duplicate
keys, tags, aliases, numeric types, later documents, and invalid binary bytes
also lose YAML-specific information. These boundaries were characterized,
not expanded into a replacement-library migration.

### Project MVS, Gates, And Vulnerabilities

The exact shortest declared graph ancestry is:

`github.com/devdimensionlab/plybuild ->
github.com/devdimensionlab/mvn-pom-mutator@v0.2.3 ->
github.com/spf13/viper@v1.10.1 -> go.etcd.io/etcd/api/v3@v3.5.1 ->
github.com/grpc-ecosystem/grpc-gateway@v1.16.0 ->
github.com/ghodss/yaml@v1.0.0`.

That is module-graph reachability, not package loading. `go mod why -m`
reports the main module does not need Ghodss YAML, and no Ghodss package is in
Ply's production or complete-test package graph. Ply directly imports
`gopkg.in/yaml.v2` and `gopkg.in/yaml.v3`; those runtime choices are separate.

Retention leaves the accepted project byte-exact at 234 modules, 3,583 graph
edges, 429 complete-test entries, 41 loaded modules, 197 loaded module-backed
packages, 1,051 `go.sum` lines, and 383 unapplied tidy-diff lines. The project
records only the Ghodss go.mod sum. A disposable exact selected `go get`
would add a redundant explicit indirect main edge and the unused source sum,
raising graph edges to 3,584, sum lines to 1,052, and tidy projection to 388
without loading a package. It was therefore not applied. A disposable master
projection likewise loads no package, while adding an unreleased direct edge,
two sums, and another declared YAML edge.

Exact Go 1.26.7 project mod verification, build, count-1, count-10, race, vet,
Windows/amd64 build, pinned golangci-lint 2.12.2, API/CLI compatibility,
byte-identical root/status/upgrade/build help, empty-HOME count-2, and complete
`make preflight` pass. All eight mutation families kill 80/80 mutants, the
launcher contract passes all 62 controls, and the audit self-test passes all
15 controls. Because the selection is unchanged, changed-selection-only
snapshot, Docker, and exact quality/audit work was inapplicable. The accepted
27/27 Q0-Q2 L2 scorecard remains
`dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.

A supplemental contained-Go-1.18 project build passes. Its full project test
diagnostic retains two pre-existing Darwin `pkg/shell` closed-file error-text
differences. Because neither retained nor master loads a Ghodss package and
their project source/package projections are identical, these inherited
failures are not Ghodss selection effects.

Fresh govulncheck v1.7.0 used the 1,392-record primary module index. Retained
and master projections are identical: 20 IDs/22 reachable traces for Darwin
and Windows symbol scans, 22 Darwin package findings, and 30 Darwin module
findings. Ghodss YAML has no module record, package/symbol finding, or trace.

The 92-entry selected evidence manifest verifies at SHA-256
`cfb5bfcc497db42efebda5323b9e187aa63520192a7354631fd93391970bd771`.
Decision-summary SHA-256 is
`724f7f2ac4394296cd36be360df90539deed1ecaf0b92386c6027b0f2c45f7b9`.
The worktree was clean before documentation handoff; no project metadata,
runtime source, manual evidence, or next-group implementation was created.
One superseded preflight invocation used the compatibility scripts' default
ignored `target/compatibility` report location and was stopped by the later
cleanliness guard. Those two generated JSON reports and their directories were
removed immediately; the valid complete rerun routed all reports beneath
session scratch and passed. The stopped attempt is not passing evidence.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/ghodss/yaml v1.0.0` as one bounded dependency group. Resolve its
complete exact-path release and repository identity, full Go-floor closure,
YAML-to-JSON behavior, historical and current maintenance state, actual
project loading, exact MVS effects, and every applicable quality contract.
Retain or select only an exact stable release whose complete minimal closure
preserves Go 1.18 and whose relevant behavior passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, and Fatih Color v1.15.0 moves. Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, and Fsnotify
remain retained. All earlier decisions and lifecycle ancestry are final. Do
not revisit them or combine another dependency group. P8 remains queued.

Project MVS selects Ghodss YAML only through the declared module edge
`github.com/grpc-ecosystem/grpc-gateway v1.16.0 -> github.com/ghodss/yaml
v1.0.0`. `go mod why -m github.com/ghodss/yaml` says the main module does not
need it, and no Ghodss YAML package is present in Ply's loaded production or
complete-test package graph. Keep the declared grpc-gateway edge, MVS
selection, Ghodss YAML's own source closure, direct Ply YAML libraries, and
runtime behavior distinct. Do not turn this into a grpc-gateway, etcd,
OpenTelemetry, gopkg.in/yaml.v2, go.yaml.in/yaml/v3, serialization migration,
or unrelated tidy review.

A minimal post-Fsnotify survey finds one proxy version: stable selected and
`@latest` v1.0.0, with proxy time 2017-03-27T23:54:44Z. Its lightweight tag is
Git commit `0ca9ea5df5451ffdf184b4428c902747c2c11cd7`. The checksum pair is
`h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=` /
`h1:4dBDuWmgqj2HViK6kFavaiC9ZROes6MMH2rRYeMEF04=`. The project currently
records only the go.mod checksum because the package is not loaded.

The published module file contains only
`module github.com/ghodss/yaml`: it has no `go` directive and no declared
requirements. Do not infer complete Go 1.18 compatibility from that file.
Resolve the imported source dependencies and minimal package/test closure
under contained Go 1.18, while preserving project MVS as a separate result.

Go-import metadata names `https://github.com/ghodss/yaml.git`. The public
repository is enabled, unarchived, non-fork, and defaults to `master`; current
master is `d8423dcdf3440d0a5baffc6f90a11e4128545620`, last pushed in 2023. Treat
the incoming facts only as a survey to verify. Resolve tags, stable,
prerelease, pseudo-version, proxy-absent, redirect, fork, alternate-path, and
unreleased identities; release/tag/commit signatures; commit time, tree,
parents and default-branch ancestry; repository status; deprecation and
retractions. A newer repository commit or a maintained successor package is
not an exact-path stable release and must not be substituted silently.

# Measurements At Start

The latest dependency implementation remains Fatih Color v1.15.0 commit
`6ca672ef38688b7f6f505cf0cb273d07c4c2ba9a`, parent
`d181fcd6c11fa147e0b44dd007598872875fc9e6`, and tree
`9ba2fdc442553622028a4a8464536915d345510c`, changing only `go.mod` and
`go.sum` with three insertions and one deletion. Fsnotify v1.6.0 was retained
without a dependency implementation commit because every newer
floor-compatible stable fails a native complete/repeat contract and v1.10.x
requires Go 1.23.

Accepted project measurements remain 234 selected modules, 3,583 graph
edges, 429 native complete-test entries, 41 loaded modules, 197 loaded
module-backed packages, 1,051 `go.sum` lines, and a 383-line unapplied tidy
projection. Relative to accepted go-cmp commit `c314bcb`, metadata adds
exactly 35 checksum lines and removes zero. The main module retains Go 1.18
and toolchain Go 1.26.7.

Fresh primary vulnerability data has 1,392 module records. Accepted
populations remain 20 IDs/22 reachable traces for Darwin and Windows symbol
scans, 22 Darwin package findings, and 30 Darwin module findings. Fsnotify has
no record, finding, or trace. Do not attribute existing x/image/ansimage or
other inherited findings to Ghodss YAML.

The unchanged accepted quality baseline has all 27 Q0-Q2 rows PASS at L2,
with scorecard SHA-256
`dae9e51e26353f72d9026e1d6eecbef697bcfa3b06cdc1905164c6d20057dafb`.
Fsnotify decision-summary SHA-256 is
`6ce37518efc54801eea666954aa6f850142909d99c347cad21d75aa408133a34`;
its 899-entry selected evidence manifest SHA-256 is
`e03b23b7a574cb7d4ec4213fc7a28d785fc5a70517a864a338fa64f227db11db`.

Read the answered Fsnotify archive and rolling handover for its complete
release, OS watcher, test-flake, MVS, vulnerability, and evidence record.

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
database, go-import metadata, authoritative repository, tag ancestry,
repository status, and primary Go vulnerability identity. Record the exact
source/mod checksums, Git identity, release time, tree, parents, signature
state, default-branch ancestry, deprecation, retractions, and all other tags
or pseudo-versions that matter to the exact path.

Prove the complete minimal module and package/test closure under exact Go
1.26.7 and contained Go 1.18.10. Inspect every imported dependency rather
than treating the absent root `go` directive as proof. Keep source-time
dependency resolution separate from the project's already-selected graph.

Inspect exported `Marshal`, `Unmarshal`, `UnmarshalStrict`, `YAMLToJSON`, and
`JSONToYAML` APIs and the actual conversion pipeline. Cover maps with string
and non-string keys, integer/bool/null scalars, numbers, tags, anchors and
aliases, duplicate keys, unknown fields, embedded/custom marshalers,
JSON-tagged structs, interfaces, byte/string inputs, invalid YAML/JSON,
multi-document input, ordering, escaping, and error propagation. Characterize
lossy YAML-to-JSON behavior and strict-mode boundaries without expanding the
task into a replacement-library migration.

Inspect Cgo/native surface, generated files, examples, testdata, fuzz/property
coverage, upstream CI coverage, and differences between v1.0.0 and unreleased
master. Add independent fixtures for release-relevant conversion, strictness,
errors, determinism, concurrency, and package/build coverage. Run source
verification, package listing, native complete tests, two independent repeated
passes, race where supported, vet, and relevant cross-builds under both SDKs.

Prove exact project module/graph/package/checksum/tidy effects for selected
v1.0.0 and any serious exact-path candidate in disposable trees. Explain why
the selected module exists in MVS while no package is loaded. Any change
outside the exact Ghodss YAML edge and its authorized MVS projection is a stop
condition. Do not independently upgrade or remove grpc-gateway or YAML
implementations.

Compare old/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject or retain if canonical identity,
release qualification, complete floor, conversion semantics, concurrency,
native tests, API, loading, MVS, or any applicable quality contract fails.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, Fatih implementation identity, reciprocal archive history,
P7/P8 state, and `./codex-dev-start.sh --check`. Read this archive, the
answered Fsnotify archive, rolling handover, roadmap, `go.mod`, `go.sum`, and
every referenced quality, compatibility, release, runner, evidence, and
lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine whether any exact-path stable release higher than selected
v1.0.0 exists and independently qualifies. Do not promote an unreleased
commit, redirect, fork, maintained successor, or alternate module path into an
exact stable upgrade. If selected/latest v1.0.0 is the only qualified stable,
retain it without hand-editing metadata or manufacturing a dependency commit.

For a changed selection, use exact Go 1.26.7 and exact `go get` for one
dependency-only commit, never tidy as implementation, then run the complete
P7 dependency gate. For an inert/retained selection, prove the no-change
effect and run the applicable dependency, closure, project, compatibility,
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

After the Ghodss YAML decision, rewrite the rolling handover and roadmap,
answer this archive, create exactly one reciprocal NEXT archive for the next
single P7 group, replace only launcher mutable regions, run launcher/handoff
contracts, and make the normal `docs: prepare next agent session` commit. Do
not implement the next group, launch a successor, push, merge, publish,
release, stash, revert, bypass cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
