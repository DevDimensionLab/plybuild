# Agent Session: Evaluate Google Cloud Go Testing Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-13T005310+0200-evaluate-google-cloud-go-testing-dependency`
Created: `2026-09-13T00:53:10+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5170e6ae91da384631db67b2dc5281bc244146a255b43dae4d43f220421d3dbf`
Previous: [2026-09-10T054841+0200-evaluate-googleapis-gax-go-v2-dependency.md](2026-09-10T054841+0200-evaluate-googleapis-gax-go-v2-dependency.md)
Next: [2026-09-13T041350+0200-evaluate-enterprise-certificate-proxy-dependency.md](2026-09-13T041350+0200-evaluate-enterprise-certificate-proxy-dependency.md)
Outcome: Retained selected Google Cloud Go Testing pseudo-version without metadata changes; no qualified exact-path release exists, selected/latest code and API are identical, full source/test closure preserves Go 1.18, zero target or GAX packages load, and all applicable project quality contracts pass.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/googleapis/google-cloud-go-testing
v0.0.0-20200911160855-bcd43fbb19e8` as one bounded dependency group.
Resolve its repository and release identity, complete Go-floor closure,
package behavior and exported API, actual project loading, exact MVS effects,
and every applicable quality contract. Retain or select only an exact-path
release whose complete minimal source/test closure preserves Go 1.18 and whose
relevant behavior passes every contract.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7 and accepted
Speakeasy v0.2.0, XXHash v2.3.0, Fatih Color v1.15.0, Go Logfmt v0.6.0,
Go Stack v1.8.1, Godbus D-Bus v5.1.0, Golang Protobuf v1.5.3, Golang Snappy
v1.0.0, Google Martian v3.3.2, Google Renameio v1.0.1, and Google UUID
v1.4.0 moves. Google pprof remains retained at
`v0.0.0-20210720184732-4bb14d4b1be1`. Gogo Protobuf v1.3.2, Crypt,
OpenCensus Proto, Logex, Readline, Fnmatch, Imaging, ansimage, Fsnotify,
Ghodss YAML, historical root GLFW, and Googleapis GAX Go v2 v2.7.0 also
remain retained. All earlier decisions and lifecycle ancestry are final. Do
not revisit them or combine another dependency group. P8 remains queued.

GAX v2.7.0 is retained without dependency metadata changes under the user's
two explicit bounded decisions. Its inherited Go-1.19 floor exception and
known `apierror.ParseError(err, false)` panic exception remain valid only
while the complete project load contains zero GAX packages. If a future
checkpoint loads or directly imports GAX, that checkpoint must stop for a
fresh dependency and product decision before merge. Do not reopen GAX.

Project MVS selects Google Cloud Go Testing through
`github.com/spf13/afero v1.9.4`. The shortest current graph path is main ->
Afero v1.9.4 -> Google Cloud Go Testing selected. Historical Afero v1.8.2
declares the same selected pseudo-version. `go mod why -m` is negative, and
the complete project load contains no package from this module. Resolve this
ancestry precisely. Do not combine, upgrade, downgrade, remove, or
independently audit Afero, Cloud Go, BigQuery, Datastore, Pub/Sub, Storage,
Google API, or another dependency group. Candidate MVS movements caused by
the exact Google Cloud Go Testing edge remain in scope to measure, but not to
broaden into independent audits.

A fresh minimal survey found an empty stable proxy version list. Selected is
the pseudo-version above, dated 2020-09-11T16:08:55Z at commit
`bcd43fbb19e8`, with source/mod sums
`h1:tlyzajkF3030q6M8SvmJSemC9DTHL/xaMa18b65+JM4=` /
`h1:dvDLG8qkwmyD9a/MJJN3XJcT3xFxOKAvTZGvuZmac9g=`. Proxy latest is
`v0.0.0-20210719221736-1c9a4c676720`, dated 2021-07-19T22:17:36Z at commit
`1c9a4c676720`, with source/mod sums
`h1:zC34cGQu69FG7qzJ3WiKW244WfhDC3xxYMeNOX2gtUQ=` /
`h1:dvDLG8qkwmyD9a/MJJN3XJcT3xFxOKAvTZGvuZmac9g=`. Both module files
declare Go 1.11 and the same Cloud Go v0.44.3, BigQuery v1.0.1, Datastore
v1.0.0, Google API v0.9.0, and tool requirements. The minimal archive survey
observed only README changes and a new SECURITY.md in latest; verify rather
than trust every incoming fact.

An empty proxy version list does not qualify an arbitrary branch head or
pseudo-version as a release. Resolve repository tags, releases, branches,
archive/deprecation status, commits, ancestry, and the exact origin of both
pseudo-versions. Do not silently promote proxy latest, an unreleased commit,
arbitrary branch head, redirect, fork, alternate path, floor-ineligible
version, prerelease, or tag that does not version this module. If the module
has no qualified exact-path release newer than selected, retain selected
unless an independent disqualifier requires product direction; do not create
a dependency commit for an inert decision.

# Measurements At Start

The latest dependency implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, and tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`, changing only `go.mod`
and `go.sum` with three insertions and no deletions. GAX, Google pprof, and
Gogo Protobuf remain retained without dependency commits or metadata edits.

Accepted project measurements remain 234 selected modules, 3,599 graph
edges, 429 native complete-test entries, 41 loaded modules, 197 loaded
module-backed packages, 1,067 `go.sum` lines, and a 432-line unapplied tidy
projection. Relative to accepted go-cmp commit `c314bcb`, metadata adds
exactly 51 checksum lines and removes zero. The main module retains Go 1.18
and toolchain Go 1.26.7.

Fresh GAX-session primary vulnerability data has 1,398 module records, index
Last-Modified 2026-09-10T16:28:28Z, and scanner DB update
2026-09-10T14:48:42Z. Accepted project populations remain 30 module
findings, 22 Darwin package findings, 23 Windows package findings, and 20
IDs/22 reachable traces on both Darwin and Windows symbol scans. No GAX
module, package, or trace enters the project. Do not attribute inherited
findings to Google Cloud Go Testing without exact evidence.

The accepted quality result remains all 27 Q0-Q2 rows PASS at L2, with seven
improved and zero held, regressed, non-comparable, or dirty counts. Scorecard
SHA-256 is
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.
GAX decision-summary SHA-256 is
`26a2a6a7ad4d06a2c059959d7eaf8df64ee9baeb1688c03927b08b1fd650bfe1`;
its 863-entry selected-evidence manifest SHA-256 is
`8b8ea0d9b842b8b3eecc1d7c7b8f8754c24662a74dd21d9d1cef0304bbd18b5b`.

Read the answered GAX archive and rolling handover for its complete identity,
closure, API, behavior, MVS, vulnerability, quality, exceptions, and evidence
record. Do not reopen GAX, UUID, Renameio, pprof, Martian, Afero, or earlier
groups.

Recreate exact Go 1.26.7 beneath `$CODEX_SESSION_SCRATCH_ROOT`, verify binary
SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in PATH, keep GOENV=off, GOWORK=off, GOTOOLCHAIN=local, and inject
no ambient GOFLAGS. Recreate contained Go 1.18.10 and pinned tools beneath
scratch as required. Portable receipts remain golangci-lint 2.12.2 archive
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29`,
GoReleaser 2.17.1 binary
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
and apidiff source archive
`0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`.
Docker timestamps with more than six fractional digits require the real
Python 3.14 interpreter ahead of `/usr/bin/python3`; preserve this as an
external execution control.

# Role And Boundaries

From fresh external archives and caches, prove the complete minimal module and
package/test closure for selected and every serious candidate under exact Go
1.26.7 and contained Go 1.18.10. Inspect imported source and test dependencies
rather than treating module Go directives as full closure proof. Separate
isolated source-time resolution from the project's selected graph, and report
both each module's declared floor and the floor of every dependency in its
minimal production and test closure.

Inspect every package and exported API. The initial source tree contains the
root documentation package plus BigQuery `bqiface`, Datastore `dsiface`,
Pub/Sub `psiface`, and Storage `stiface` adapters/interfaces. Characterize
all wrappers, forwarding, embedded interfaces, constructors, return values,
nil/zero inputs, nil embedded clients, typed nils, aliasing, iterator and
stream behavior, callbacks, context propagation and cancellation, errors and
panics, concurrent reuse, generated/manual content, examples, platform
differences, and compatibility with their pinned Cloud client APIs.
Distinguish runtime packages, tools-only requirements, internal helpers,
commands, examples, benchmarks, testdata, fuzz/property coverage, and upstream
CI. Do not infer implementation guarantees from interface declarations alone.

Add independent fixtures where useful for exact method delegation, argument
and context identity, return/error preservation, interface satisfaction,
nil/panic boundaries, iterator/page behavior, stream receive/send behavior,
callback invocation, mock substitution, concurrency, deterministic inputs,
and compatibility with the selected project graph. Run source verification,
package listing, native complete tests, two independent repeats, race, vet,
and meaningful cross-builds under both SDKs. Classify every clock, timer,
randomness, network, emulator, credential, platform, resource, or test-design
failure precisely. Do not require live cloud services for unit qualification;
classify integration prerequisites and keep them contained.

Prove exact project module/graph/package/checksum/tidy effects for selected and
every serious candidate in disposable trees. Explain why Google Cloud Go
Testing exists in MVS while no package is loaded, and preserve every unrelated
module selection. Any change outside the exact target edge and its necessary
authorized MVS projection is a stop condition.

Compare selected/candidate primary vulnerability results at module, package,
symbol, and reachable-trace levels. Reject, retain, or stop if canonical
identity, release qualification, complete floor, behavior, concurrency,
tests, API, loading, MVS, or any applicable quality contract fails. Do not
treat absence from the loaded package set as permission to ignore source-time
floor, closure, behavior, or vulnerability evidence.

# Required Reading

At start, verify the feature branch, clean ordinary and ignored status,
current ancestry, latest dependency implementation identity, reciprocal
archive history, P7/P8 state, and `./codex-dev-start.sh --check`. Read this
archive, the answered GAX archive, rolling handover, roadmap, `go.mod`,
`go.sum`, and every referenced quality, compatibility, release, runner,
evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path release. Do not promote a
redirect, fork, alternate path, floor-ineligible version, prerelease,
non-versioning tag, or unreleased commit. In particular, proxy `@latest` does
not itself qualify the latest pseudo-version as a release. If no higher
release qualifies, retain selected without hand-editing metadata or
manufacturing a dependency commit.

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

After the Google Cloud Go Testing decision, rewrite the rolling handover and
roadmap, answer this archive, create exactly one reciprocal NEXT archive for
the next single P7 group, replace only launcher mutable regions, run
launcher/handoff contracts, and make the normal
`docs: prepare next agent session` commit. Do not implement the next group,
launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Retain exact-path
`github.com/googleapis/google-cloud-go-testing
v0.0.0-20200911160855-bcd43fbb19e8` without changing `go.mod` or
`go.sum` and without manufacturing a dependency commit. The repository has
never published a qualified release, proxy latest is merely its final archived
master-branch commit, and selected has no independent floor, behavior, API,
loading, MVS, vulnerability, or quality disqualifier.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves the public,
Apache-2.0 repository
`https://github.com/googleapis/google-cloud-go-testing`. GitHub marks it
archived, not disabled and not a fork, with protected default branch `master`.
It has zero tags and zero GitHub releases. `RELEASING.md` explicitly says the
repository has no releases and consumers should use pseudo-versions. The
stable proxy version list is empty. Therefore neither proxy `@latest`, an
arbitrary branch head, nor either pseudo-version is independently qualified as
a release.

Selected is commit
`bcd43fbb19e8d79524fce1b71f4a2145afbd6039`, parent
`8e1d251e947d1de4242ca2a81eb1a39917bf6b38`, tree
`44cc1490ddeb164a479884cd8beecc373b2ad60d`, dated
2020-09-11T16:08:55Z, with subject `chore: add CODEOWNERS (#38)`. It is the
tip of branch `bcb-to-fb`. Its source/mod sums are
`h1:tlyzajkF3030q6M8SvmJSemC9DTHL/xaMa18b65+JM4=` /
`h1:dvDLG8qkwmyD9a/MJJN3XJcT3xFxOKAvTZGvuZmac9g=`. Proxy latest
`v0.0.0-20210719221736-1c9a4c676720` is commit
`1c9a4c676720af1d2d964c1b2c866ee500300a15`, parent
`1487aa9ec5b057debb42f236a1d1185f09f96804`, tree
`e1733414a158ce2090bf6eb073c4c6af8dceee52`, dated
2021-07-19T22:17:36Z, with subject `chore: archive project (#43)`. It is the
tip of `master`; its sums are
`h1:zC34cGQu69FG7qzJ3WiKW244WfhDC3xxYMeNOX2gtUQ=` /
`h1:dvDLG8qkwmyD9a/MJJN3XJcT3xFxOKAvTZGvuZmac9g=`.

Strict and full Git fsck pass. Selected is an ancestor of latest with exactly
two intervening commits. GitHub verifies both commit signatures. Git and proxy source bytes
match: selected has 37 files with manifest SHA-256
`ac0e86d2adf82603aaf6d39a402f8869c2129dff810d773f6cad0dca9bcd661c`;
latest has 38 files with manifest SHA-256
`f6d9a5c0008cfd35eb2b621d13141245d91b378c421a95e9be0a72e351335ded`.
The only selected-to-latest changes are archive/status wording in README and a
new SECURITY.md. No Go source, module metadata, checksum, test, or CI behavior
changed.

Selected and latest declare Go 1.11 and identical requirements: Cloud Go
v0.44.3, BigQuery v1.0.1, Datastore v1.0.0, Google API v0.9.0, plus build-tagged
tool requirements. Under exact Go 1.26.7 each isolated graph has 36 modules
and 337 edges, with 332 production and 359 test package entries. Under
contained Go 1.18.10 each has 36 modules and 336 edges, with 270 production
and 296 test entries. Both actually import 14 external modules and expose five
runtime packages. The maximum Go directive in the imported production/test
closure is 1.11; the maximum across the whole graph, including tools-only
vertices, is 1.12. Complete minimal closure therefore preserves Go 1.18.

The root package is documentation plus an example. Runtime packages are
BigQuery `bqiface`, Datastore `dsiface`, Pub/Sub `psiface`, and Storage
`stiface`. Their exported surface consists of adapters and interface/config
shadows: BigQuery has `AdaptClient`, thirteen interfaces and seven wrapper
structs; Datastore has `AdaptClient` and four interfaces; Pub/Sub has
`AdaptClient`, `AdaptMessage`, five interfaces and `SubscriptionConfig`;
Storage has `AdaptClient` and ten interfaces. Each interface embeds an
unexported method to constrain production implementations to the supplied
wrappers while still allowing fakes by embedding the interface. There are no
runtime commands, benchmarks, fuzz/property targets, generated Go files, or
testdata. `tools.go` is tools-build-tag-only. Upstream Kokoro uses Go 1.12 and
runs lint, staticcheck, tidy checks, and race tests. Four live examples require
`BQIFACE_PROJECT`, `DATASTORE_PROJECT_ID`, `PSIFACE_TOPIC`, or
`STIFACE_BUCKET` and skip cleanly without credentials.

All adapters are small value wrappers around embedded Cloud client pointers.
Calls generally forward the exact context, arguments, callback, and upstream
return values. Constructors given nil clients return non-nil typed interface
values containing nil embedded pointers; most method calls then panic.
Storage can construct bucket/object handles before an I/O method reaches the
nil client and panics. Successful upstream nil pointers can likewise be
wrapped as non-nil interfaces.

BigQuery query/load/copy/extract config conversion is shallow: slices, maps,
pointers, and nested values alias caller storage, and builder setters mutate
the wrapped upstream object. `SetQueryConfig` deliberately preserves the
embedded destination table when its shadow `Dst` is nil; `SetCopyConfig`
appends all sources. Iterator/job/result wrappers frequently return nil result
interfaces on error and wrap nil results on nil error. Metadata conversion
panics for nil metadata or nil access entries. BigQuery copy and Storage
copy/compose require the package's concrete wrappers and panic for otherwise
valid external fakes.

Datastore preserves transaction, commit, iterator, callback, context, and
error identities; transaction/commit/run methods can return a wrapper along
with a non-nil upstream error. Pub/Sub directly forwards message data and
attributes without copying, so aliases are preserved. Ack/Nack/Get and
subscription receive forward directly; Receive passes the exact context and
wraps each message for the callback. Pub/Sub publish asserts its concrete
message wrapper and panics for external or typed-nil alternatives. Storage
builder setters and callback hooks mutate the embedded upstream builders and
preserve upstream return/error behavior. Mutable builders and iterators are
not promised safe for concurrent mutation; independent/read-only wrapper use
passes race tests. The wrappers add no retry, clock, randomness, or network
policy of their own; those remain properties of the pinned Cloud clients.

Selected and latest produce byte-identical exported API blobs, SHA-256
`1a9dbaa10ba5d79c1b207cd5acaefe6c1903c3ecceb868c2fe6d8605370a2e00`,
and apidiff is empty in both directions. The independent 375-line behavior
fixture, SHA-256
`10bea668cbc7b6ec438f6fa15ff7b8681e6901ef4c4076f8909b052abe5d991c`,
covers exact delegation and identities, config aliasing, concrete-assertion
panics, nil and typed-nil boundaries, iterator paging, HTTP cancellation,
Datastore callbacks, local Pub/Sub streaming through `pstest`, Storage
builders, errors, and concurrent reuse. It passes native, count-10, race, and
vet runs for both versions under both SDKs.

The upstream source has three examples named `Example_AdaptClient`. Modern Go
vet reports those suffixes as referring to missing identifiers. Exact Go
1.26.7 default `go test` and standalone vet fail only for those names; exact
Go 1.18.10 default tests pass but standalone vet reports the same defect.
Runtime tests pass with vet disabled. Scratch-only renaming to
`ExampleAdaptClient` makes default tests and vet pass for both versions and
both SDKs without altering production source. This is a historical test-name
defect, not a runtime or release-selection disqualifier. With that controlled
correction, all production builds pass for Darwin/amd64, Linux/amd64,
Linux/arm64, Windows/amd64, FreeBSD/amd64, and js/wasm; representative Linux
and Windows test compilation also passes.

Afero v1.9.4 supplies the only current incoming edge and shortest MVS path:
main -> Afero v1.9.4 -> selected Google Cloud Go Testing. Historical Afero
v1.8.2 declares the same selected pseudo-version. `go mod why -m` is negative,
and the complete project load contains exactly zero target packages and zero
GAX packages. A raw explicit selected request keeps 234 modules and all
selections, raises graph edges from 3,599 to 3,605, and adds six source sums
for the target's otherwise-pruned tools vertices. Both raw requests temporarily
add explicit requirements for the target, BurntSushi TOML, `x/lint`, `x/mod`,
`x/tools`, and `honnef.co/go/tools`. A raw latest request keeps 234 modules,
changes only the target selection, raises edges to 3,612, and adds seven source
sums. Neither
changes the 429 complete-test entries, 41
loaded modules, 197 loaded module-backed packages, target/GAX zero-load result,
or application behavior.

Normal tidy makes base, selected, and latest projections byte-identical in
`go.mod` and `go.sum`, restores the Afero-selected pseudo-version, and measures
234 modules, 3,557 graph edges, and 948 sum lines. The accepted unchanged
project remains 234 modules, 3,599 edges, 429 complete-test entries, 41 loaded
modules, 197 loaded module-backed packages, 1,067 sum lines, and a 432-line
unapplied tidy projection. No unrelated module was selected, downgraded,
upgraded, removed, or independently audited.

Fresh primary vulnerability data contains 1,398 module records, index
Last-Modified 2026-09-10T16:28:28Z, scanner DB update
2026-09-10T14:48:42Z, and zero exact target records. Project populations stay
at 30 module findings, 22 Darwin and 23 Windows package findings, and 20 IDs /
22 reachable traces on each symbol platform, with no target or GAX occurrence.
Direct selected/latest source-time scans are identical: 31 inherited module
findings; 18 package findings across 16 IDs; 51 production traces and 66 test
traces across three IDs. GO-2023-1571 and GO-2024-2687 affect old `x/net`
HTTP/2, and GO-2026-6061 affects old gRPC transport. Target code appears only
as caller frames; the affected modules are `x/net` and gRPC, so these findings
are not attributed to Google Cloud Go Testing.

Exact Go 1.26.7 project verify/load/build, count-1, two count-10 repeats, race,
vet, pinned lint, empty-HOME count-2, Linux/Windows builds, API/CLI
compatibility, complete preflight, 80/80 mutation controls, audits, and host,
snapshot, and real Docker acceptance pass. Docker built once, ran ten
containers, made 54 calls, and published zero times. The contained Go 1.18.10
projection passes applicable corrected gates; the full project has only the
two already accepted `pkg/shell` closed-file error-wording differences. The
fresh independent quality run is 27/27 Q0-Q2 PASS at L2, seven improved, zero
held/regressed/non-comparable/dirty, with scorecard SHA-256
`a706ee72aa1483217bdcced30b36523da95841fb3874ef017645ffcf74051829`.
Because dependency metadata is unchanged, the accepted scorecard remains
`576c6e9f666d89adf533bb0e24503e966b9935d3e70183f7ddf52bad84b1edac`.

The 788-entry selected-evidence manifest SHA-256 is
`34e3574d2ced8aef523df1607300134664907cf7634cb40f89c3f615a151f106`;
decision-summary SHA-256 is
`9652c69c7971189749d2fd4590dc39b49a11f6f159cb80c2287767251312db2a`.
