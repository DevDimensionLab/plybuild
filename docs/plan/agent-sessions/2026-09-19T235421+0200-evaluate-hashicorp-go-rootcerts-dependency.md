# Agent Session: Evaluate Hashicorp Go Rootcerts Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-19T235421+0200-evaluate-hashicorp-go-rootcerts-dependency`
Created: `2026-09-19T23:54:21+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c05529eedf4c44e43ea1a9e7eaba81e7fdb56b110781542abad6002107c1b86d`
Previous: [2026-09-19T233605+0200-decide-hashicorp-go-retryablehttp-product-direction.md](2026-09-19T233605+0200-decide-hashicorp-go-retryablehttp-product-direction.md)
Next: [2026-09-20T115926+0200-decide-hashicorp-go-rootcerts-product-direction.md](2026-09-20T115926+0200-decide-hashicorp-go-rootcerts-product-direction.md)
Outcome: No exact-path stable release qualifies: v1.0.0, v1.0.1, and latest v1.0.2 all lose filesystem error identity and can silently return an empty Darwin system-root pool, while their proxy archives omit a tracked symlink fixture required by an upstream test; dependency metadata remains unchanged and P7 stops for the reciprocal product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/go-rootcerts v1.0.2` as one bounded dependency group.
Resolve its complete repository and release identity, Go-floor closure,
package behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific retained-module decisions through exact inherited,
unloaded go-retryablehttp v0.5.3. Every earlier outcome and lifecycle
ancestor is final. Evaluate only Hashicorp go-rootcerts in this session; do
not reopen or combine another dependency group. P8 remains queued.

The authorized 2026-09-19 go-retryablehttp option 1 decision retains exact
selected, inherited, unloaded `github.com/hashicorp/go-retryablehttp v0.5.3`
without dependency metadata changes. Its target-specific, non-transferable
exception accepts only GO-2024-2947/GHSA-v6v8-xj6m-xwqh and the reproduced
Basic-auth URL credential disclosure; the x509, unsupported-scheme, and
redirect-limit retry-classification defects; final transport-error identity
loss; the test-only vet defect; and the completed API, body, cancellation,
status, backoff, hook, nil/panic, mutation, aliasing, allocation, concurrency,
global-state, resource, MVS, vulnerability, and related findings. It remains
valid only while exact v0.5.3 and the sole go-metrics v0.3.10 request remain
unchanged, repository imports and target loads remain zero, runtime
unreachability holds, all earlier guards remain intact, and no new advisory or
independent defect appears. Any change requires the owning fresh decision.

The go-multierror v1.1.0, go-msgpack v0.5.3, go-immutable-radix v1.3.1,
go-hclog v1.2.0, Errwrap v1.0.0, qualified go-cleanhttp v0.5.2, and every
other earlier exception or qualification remain separate under their exact
selection, incoming-edge, zero-load, runtime-unreachable, and no-new-finding
guards. Revalidate those guards and stop for the owning decision if any
expires. No earlier exception transfers to go-rootcerts. Do not change a
guarded parent, the Go floor, or an unrelated module.

Selected `github.com/hashicorp/go-rootcerts v1.0.2` is inherited through
exact requests from Viper v1.15.0, historical Viper v1.10.1,
`sagikazarmark/crypt v0.4.0`, and historical Consul API v1.12.0. Its current
`go mod why -m` result and repository import search are negative, and current
production and complete-test loads contain zero target packages. These queue
observations and the physical MVS selection are not proof of repository
identity, release qualification, ancestry, floor, behavior, vulnerability
state, or suitability. Resolve them independently and do not add a direct
edge merely to alter MVS.

# Measurements At Start

The go-retryablehttp decision recording began from clean handoff HEAD
`a4c24f29e25eb21cc229cde674ffdacdecc78f38`, parent
`24b1f33671362309526328170e67e68aaf753a3a`, tree
`ccca83e0256147f7e51d14eb07a2fb2472e2aabe`. Verify the new handoff HEAD,
parent, tree, exact changed-file set, clean ordinary and ignored status,
reciprocal archive history, and `./codex-dev-start.sh --check` rather than
assuming them. The latest dependency implementation remains exact Google UUID
v1.4.0 commit `cf53bc64eeb69471d35c7536d196bf1da15f3973`.

Guard-only decision revalidation under exact Go 1.26.7 preserved all 15
earlier guarded selections and recorded requests plus exact
go-retryablehttp v0.5.3 and its sole go-metrics request. All 16 why results
remain negative, repository imports are zero, and production and complete-
test loads contain zero guarded packages. The project remains 234 modules,
3,599 graph edges, 355 production entries, 429 complete-test entries, 197
module-backed entries across 41 loaded modules, 1,067 sum lines, and the
recorded 432-line unapplied tidy projection. Base `go.mod` and `go.sum`
SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. Exact-version OSV remains empty for
14 earlier guards, Gorilla retains only its recorded GO-2026-6278/
GHSA-w67g-5rqw-f597 result, and go-retryablehttp retains exactly its two
accepted identifiers. No new guarded advisory or independent defect appeared.

Use exact Go 1.26.7 binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, use `LC_ALL=C LANG=C`, and
run with `umask 022`. Recreate contained Go 1.18.10 and pinned tools beneath
`$CODEX_SESSION_SCRATCH_ROOT` as required. Preserve the known apidiff archive
reproducibility discrepancy, Python 3.14 Docker timestamp control, managed
bare-`mktemp` restriction, and nested launcher signal-retention timing race;
none is go-rootcerts evidence.

# Role And Boundaries

From fresh external archives and caches, resolve proxy, sumdb, `go-import`,
Git, and forge evidence for the exact module path: tags, releases, branches,
signatures, commits, parents, trees, times, ancestry, repository status,
licenses, retractions, deprecations, redirects, forks, alternate paths, major
module lines, and every serious exact-path stable candidate. Do not silently
promote a redirect, fork, alternate path, different major module path,
prerelease, non-versioning tag, unreleased branch head, or floor-ineligible
release.

Prove the complete minimal production and test closure under exact Go 1.26.7
and contained Go 1.18.10. Inspect imported source and test dependencies rather
than treating the module directive alone as floor proof. Separate isolated
source-time resolution from the project's selected MVS graph and from every
already-final guarded dependency decision.

Inspect every package, command, exported API, example, benchmark, fuzz target,
testdata, generated file, platform/build-tag branch, and applicable API and
runtime boundary. Characterize certificate-pool construction, PEM/file/path
handling, system roots, environment interaction, error identity, nil/panic
behavior, mutation and aliasing, allocation, concurrency, global state,
resource cleanup, malformed inputs, platform behavior, and actual project
consumers. Add independent fixtures where useful and run source verification,
package listing, native complete tests, two independent repeats, race, vet,
and meaningful cross-builds under both SDKs. Classify every failure precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, floor, architecture, or unrelated-
module change needs a fresh bounded decision rather than silent
implementation. Compare primary vulnerability results at module, package,
symbol, test-symbol, and reachable-trace levels.

# Required Reading

At start read this archive, the answered go-retryablehttp decision and
evaluation, the answered go-multierror and go-hclog decisions, rolling
handover, roadmap, `go.mod`, and `go.sum`. Verify the feature branch, clean
ordinary and ignored status, current ancestry, latest Google UUID
implementation identity, reciprocal archive history, P7/P8 state, every
unchanged qualification/exception guard, and `./codex-dev-start.sh --check`.
Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose
complete minimal source/test closure preserves Go 1.18. Do not promote an
unqualified or floor-ineligible identity. If no candidate satisfies the
existing contracts, preserve the evidence and stop for a fresh bounded
product decision.

Second, qualify selected and serious candidates across source, tests, API,
behavior, loading, MVS, vulnerability, and project contracts. For an
authorized changed selection, use exact Go 1.26.7 and exact `go get` for one
dependency-only commit, never tidy as implementation, then run the complete
P7 dependency gate. For a retained or blocked selection, prove the no-change
effect and run all applicable gates without manufacturing activity. Full
changed-selection quality must preserve 27/27 Q0-Q2 PASS at L2 with zero
held, regressed, non-comparable, or dirty counts.

Third, update the roadmap and rolling handover with exact evidence, outcome,
commit identity, limitations, and next boundary. Answer this archive and
prepare one reciprocal NEXT mission only after the bounded outcome is
coherent and committed. Do not execute the successor.

# Automatic Handoff

After a completed coherent result, make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

The evaluation began clean on branch `codex/upgrade-quality` at handoff HEAD
`888412a62d768cb11af6d9e132600cc4a34f0679`, parent
`a4c24f29e25eb21cc229cde674ffdacdecc78f38`, tree
`1892ece8798c3913cc46465cf12865281db20339`. The handoff changed exactly the
launcher, answered go-retryablehttp decision archive, this then-NEXT archive,
rolling handover, and roadmap. Ordinary and ignored status were empty; the
reciprocal archive chain and `./codex-dev-start.sh --check` passed. The latest
dependency implementation remains Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`.

Fresh `go-import` metadata resolves the exact path to
`https://github.com/hashicorp/go-rootcerts.git`. GitHub identifies a public,
active, unarchived, enabled, non-fork MPL-2.0 repository with protected
default branch `master`, 12 branches, and no GitHub Release objects. The Go
proxy exposes exactly three stable versions, v1.0.0, v1.0.1, and latest
v1.0.2. There is no prerelease, retraction, module deprecation, redirect,
alternate exact path, or `/v2` line. Current master
`95e2cf6a7dc27942f8dba98f69db96f2d3169efa` is 34 commits beyond v1.0.2;
its pseudo-version declares Go 1.23 and is unreleased, so it is not promoted.

Strict mirror verification passes. V1.0.0 is lightweight unsigned commit
`63503fb4e1eca22f9ae0f90b49c5d5538a0e87eb`, tree
`26bab66...`, dated 2019-01-18. V1.0.1 is lightweight unsigned commit
`df8e78a645e18d56ed7bb9ae10ffb8174ab892e2`, tree `74bf9fa...`, dated
2019-06-10. V1.0.2 is unsigned annotated tag object
`dafef5a8f4d570697f148431a2004909428ca75b` over GitHub-verified commit
`98fadc2a5ba2ad2a534a179b352ecdfd1f4259aa`, parents `df8e78a...` and
`bc065...`, tree `7e2fd4...`, dated 2019-12-10. All tags are ancestors of
master.

Proxy-ZIP SHA-256 values for v1.0.0/v1.0.1/v1.0.2 are respectively
`4393b0b9cd741e00de5624d5124cf054bf50c57231d4b1caff84c8a4d16c6a47`,
`3f558b1a436ed6fb15872383545109227f9552bf5daa95583e9402bbd3a24fff`, and
`864a48e642e87a273fb5ef60bb3575bd74a7090510f93143163fa6700be31948`.
Sumdb identities are
`h1:Rqb66Oo1X/eSV1x66xbDccZjhJigjg0+e82kpwzSwCI=`,
`h1:DMo4fmknnz0E0evoNYnV48RjWndOsmd6OW+09R3cEP8=`, and
`h1:jzhAVGtqPKbwpyCPELlgNWhE1znq+qwJtW5Oi2viEzc=`. V1.0.1/v1.0.2 share
go.mod sum `h1:pqUvnprVnM5bf7AOirdbb01K4ccR319Vf4pU3K5EGc8=`. Proxy and Git bytes
match for every regular file. Go module ZIP rules omit the two tracked
`test-fixtures/capath-with-symlinks` symlinks, which is behaviorally material
to the archived upstream test suite.

### Go Floor, Source, Tests, And API

V1.0.0 has no Go directive; v1.0.1/v1.0.2 declare Go 1.12. The complete
minimal non-standard-library closure is target plus
`github.com/mitchellh/go-homedir`: v1.0.0 selects homedir v1.0.0 and the later
tags select v1.1.0. Homedir has no Go directive. Imported production and test
source for all three candidates therefore preserves Go 1.18. Exact Go 1.26.7
resolves 190 complete-test entries and contained Go 1.18.10 resolves 129;
target and homedir are the only external modules.

Each tag contains one library package with three production Go files, one
documentation file, two test files, certificate testdata, and Darwin versus
`!darwin` system-root branches. There are no commands, examples, benchmarks,
fuzz targets, generated files, cgo, embeds, or generators. Tagged Git source
for all three versions passes module verification, package listing, native
count-one, two independent count-ten repeats, race, vet, and Darwin AMD64,
Linux AMD64/ARM64, Windows AMD64, and js/wasm production/test builds under
both SDKs. Testing directly from each proxy archive fails only
`TestLoadCACertsFromDirWithSymlinks` because the standard module ZIP omitted
the symlink fixture directory; tagged Git source passes the same test. This is
a reproducible release-archive/test-fixture defect, not a production compile
failure.

Pinned `apidiff` finds no exported API change from v1.0.0 to v1.0.1. V1.0.2
compatibly adds `AppendCertificate` and `Config.CACertificate`, but the added
`[]byte` field makes `Config` non-comparable and is therefore an incompatible
type-property change. Selected v1.0.2 exports `Config` plus `ConfigureTLS`,
`LoadCACerts`, `LoadCAFile`, `AppendCertificate`, `LoadCAPath`, and
`LoadSystemCAs`. There is no later stable API candidate.

### Independent Behavior And Security

The shared independent behavior fixture SHA-256 is
`512a4a1df5b96659c975ca6863c8172832a37a4fe9038da015165b28cad1acd1`;
the selected-only extension is
`4f0c62af97695e5015db568558cf401e85d03d41efca8954e9d6c3f727d97768`.
All applicable characterization tests pass under both SDKs: valid file/path
pool construction, empty-directory behavior, CAFile/CACertificate/CAPath
precedence, ConfigureTLS success and failure mutation, nil target handling,
multi-certificate bundles, malformed and trailing data, caller-PEM mutation,
x509 verification, allocations, concurrent loads, and resource release.

All three stable releases fail the same error-identity contract.
`LoadCAFile`, `LoadCAPath` through a dangling symlink, and `ConfigureTLS`
stringify underlying I/O errors with `%s`; `errors.Is(err, fs.ErrNotExist)` is
false. They also fail the Darwin system-root contract. When the three
keychain commands return success with no PEM bytes, as reproduced with a
fresh empty HOME in the managed Darwin environment, `LoadSystemCAs` returns a
non-nil but empty `*x509.CertPool` and nil error. The source ignores the false
result from `AppendCertsFromPEM`; the upstream Darwin test checks only the nil
error and misses the empty trust result. Current master retains both defects.

The package has no mutable globals. File and path functions create local
pools; returned pools remain caller-mutable, while input PEM bytes do not
alias the pool. `filepath.Walk` is lexical, follows symlinked files, aborts
and discards the pool on a malformed or unreadable entry, and returns an empty
pool for an empty directory. `ConfigureTLS` precedence is file, in-memory
certificate, path, then system roots; a nil TLS target returns nil before
validation. The Darwin-only homedir dependency adds its documented cached
HOME/environment lookup, protected by a mutex.

Fresh primary vulnerability data contains 1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`
and Last-Modified 2026-09-17T17:29:18Z. It has no rootcerts record; exact OSV
queries for all three tags and GitHub repository advisories are empty.
Govulncheck v1.8.0 module, package, symbol, and test-symbol scans are empty for
each isolated tag. Base and redundant-v1.0.2 project scans have identical
finding sets at every level, with no target frame or reachable trace.

### Project MVS And Quality Effects

Selected v1.0.2 exists through five graph requests: Viper v1.15.0, historical
Viper v1.10.1, `sagikazarmark/crypt v0.4.0`, and historical Consul API v1.12.0
request v1.0.2; historical Consul API v1.1.0 requests v1.0.0. The shortest
root path is main -> direct Viper v1.15.0 -> target. `go mod why -m` and the
repository import search are negative; the 355-entry production and 429-entry
complete-test loads contain zero target packages, so current target code is
runtime-unreachable.

The base remains 234 modules, 3,599 graph edges, 197 module-backed complete-
test entries across 41 loaded modules, and 1,067 sum lines. A disposable exact
v1.0.2 `go get` adds only a redundant main-module indirect edge, one graph
edge, and the full source checksum; selections and loads are unchanged, and
tidy converges byte-identically with the base projection. Exact v1.0.1 or
v1.0.0 instead forces Viper v1.15.0 to v1.8.1, removes mvn-pom-mutator, drops
the graph to 181 modules and approximately 2,553 edges, changes numerous
earlier guarded selections, and makes project packages unloadable. No
projection was applied.

All 16 earlier guarded selections and recorded requests remain exact. Their
why results and the target why result are negative; repository imports and
both package loads contain zero guarded packages. Fresh exact-version OSV
revalidation is empty except for Gorilla WebSocket's recorded
GO-2026-6278/GHSA-w67g-5rqw-f597 and go-retryablehttp's recorded
GO-2024-2947/GHSA-v6v8-xj6m-xwqh. No owning guard expired.

Exact Go 1.26.7 project verification, build, count-one, two independent
count-ten repeats, race, vet, pinned lint, API/CLI compatibility, empty-HOME,
four cross-builds, host/snapshot/Docker acceptance, full preflight, all 17
script/meta pairs, all eight mutation meta-stages, 80/80 live mutation kills,
and all 15 audit meta-controls pass. The contained Go 1.18 projection removes
only the unsupported toolchain line, resolves 366 complete-test entries, and
reproduces only the two accepted Darwin `pkg/shell` closed-file wording
failures; all 26 unaffected packages pass count-one, both count-ten repeats,
race, and vet, while host acceptance and four cross-builds pass.

The prior external manual evidence is correctly rejected as stale because it
is commit-bound to `24b1f336...`, while this docs-only handoff began at
`888412a...`. A changed-selection scorecard does not apply: production source,
dependency metadata, and measured behavior remain unchanged. The accepted
authoritative result therefore remains 27/27 Q0-Q2 PASS at L2, scorecard
SHA-256
`de13154181ae319fd80a29a98df78535762d70c2c01438f9e3cb735c2f624e81`,
with zero held, regressed, non-comparable, or dirty counts. The initial cold
offline API cache and bare-`mktemp` attempts were superseded by the established
scratch cache and wrapper; they are known environment boundaries, not target
findings.

Final `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Because every stable release shares independently disqualifying behavior and
the older tags also violate project graph guards, no release qualifies. The
next bounded action is the prepared go-rootcerts product decision. P8 and
every other dependency group remain untouched.
