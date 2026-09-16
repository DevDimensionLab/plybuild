# Agent Session: Evaluate Hashicorp Go Immutable Radix Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-16T210129+0200-evaluate-hashicorp-go-immutable-radix-dependency`
Created: `2026-09-16T21:01:29+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `db59e7eeb349385d4781d3ccbc6087ece27e102d12e70d97fd42959c6938cb01`
Previous: [2026-09-15T234309+0200-decide-hashicorp-go-hclog-product-direction.md](2026-09-15T234309+0200-decide-hashicorp-go-hclog-product-direction.md)
Next: [2026-09-16T225521+0200-decide-hashicorp-go-immutable-radix-product-direction.md](2026-09-16T225521+0200-decide-hashicorp-go-immutable-radix-product-direction.md)
Outcome: No exact-path stable release qualifies; dependency metadata remains unchanged and P7 stops for a bounded product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/hashicorp/go-immutable-radix v1.3.1` as one bounded dependency
group. Resolve its complete repository and release identity, Go-floor closure,
package behavior and exported API, actual project loading, exact MVS effects,
vulnerability evidence, and every applicable quality contract. Retain or
select only a qualified exact-path stable release whose complete minimal
source/test closure preserves Go 1.18 and whose relevant behavior passes every
contract; otherwise stop for a fresh bounded product decision.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified exact go-cleanhttp
v0.5.2, and all recorded retained-module decisions through exact inherited,
unloaded go-hclog v1.2.0. Every earlier outcome and lifecycle ancestor is
final. Evaluate only Hashicorp go-immutable-radix in this session; do not
reopen or combine another dependency group. P8 remains queued.

The user's 2026-09-16 go-hclog option 1 decision retains exact selected,
inherited, unloaded `github.com/hashicorp/go-hclog v1.2.0` without dependency
metadata changes. It accepts only the five completed common behavior findings
and the characterized nil/panic, aliasing, allocation, concurrency, resource,
API, closure-vulnerability, and related qualification findings. Its exception
is go-hclog-specific and valid only while exact v1.2.0 and the sole selected-
version Viper v1.15.0 incoming edge remain unchanged, zero target packages
load, the module remains runtime-unreachable, and no new advisory or
independent defect appears. Direct import/loading, runtime reachability, a
target version or incoming-edge change, or a new advisory or independent defect
requires a fresh go-hclog dependency and product decision before merge.

Qualified go-cleanhttp v0.5.2 and the Errwrap v1.0.0, Consul SDK v0.8.0,
Consul API v1.18.0, Gateway v1.16.0, gRPC Prometheus v1.2.0, gRPC middleware
v1.0.0, Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 decisions remain separate under their recorded exact-selection,
incoming-edge, zero-load, runtime-unreachable, and no-new-finding guards.
Revalidate those guards and stop for the owning decision if any expires. Do
not transfer an earlier exception to go-immutable-radix or change Viper,
either Gateway parent, mvn-pom-mutator, GoConvey, historical Consul API,
go-multierror, Serf, or `sagikazarmark/crypt`.

# Measurements At Start

The go-hclog decision recording began from clean HEAD
`d536b6832a85123dc49cf4fbcf4876f4083458d4`, parent
`f7f3a65d645ce4ce4caf8d1ac3e1da77ccd19727`, tree
`cab76bf3e4948d54fdbf032ba0c4a788c7cafda7`. Verify the new handoff HEAD,
parent, tree, clean status, and exact handoff contents at start rather than
assuming their values. The latest dependency implementation remains exact
Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`.

Guard-only revalidation under exact Go 1.26.7 preserved all 12 guarded module
selections and recorded incoming edges. All 12 `go mod why -m` results remain
negative, repository imports are zero, and production and complete-test loads
contain zero guarded packages. The project remains 234 modules, 3,599 graph
edges, 355 production entries, 429 complete-test entries, 197 module-backed
complete-test entries across 41 loaded modules, 1,067 sum lines, and the
recorded 432-line unapplied tidy projection. Base `go.mod` and `go.sum` SHA-256
values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Fresh primary vulnerability data remains 1,399 records at SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
and Last-Modified 2026-09-15T18:54:35Z. Exact go-hclog v1.2.0 and every
guarded target except Gorilla have empty exact-version OSV results. Gorilla
retains only its recorded GO-2026-6278/GHSA-w67g-5rqw-f597 exact-version
result and existing GO-2020-0019 primary-index entry. No new guarded advisory
or independent defect appeared.

A bounded queue survey identifies selected Hashicorp go-immutable-radix
v1.3.1 in the current build list. That physical selection is not proof of
repository identity, release qualification, ancestry, floor, package loading,
behavior, vulnerability state, or suitability. Resolve those facts
independently and do not add a direct edge merely to alter MVS.

Use exact Go 1.26.7, verify binary SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
put it first in `PATH`, keep `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, inject no ambient `GOFLAGS`, and use
`LC_ALL=C LANG=C`. Recreate contained Go 1.18.10 and pinned tools beneath
`$CODEX_SESSION_SCRATCH_ROOT` as required. Preserve the known apidiff archive
reproducibility discrepancy, Python 3.14 Docker timestamp control, managed
bare-mktemp restriction, and nested launcher signal-retention timing race;
none is Hashicorp go-immutable-radix evidence.

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
testdata, generated file, platform or build-tag branch, and applicable API and
runtime boundary. Characterize deterministic behavior, errors and identity,
nil/panic behavior, mutation and aliasing, allocation, concurrency, global
state, resource cleanup, malformed inputs, and project consumers. Add
independent fixtures where useful and run source verification, package
listing, native complete tests, two independent repeats, race, vet, and
meaningful cross-builds under both SDKs. Classify every failure precisely.

Prove exact project module, graph, package, checksum, tidy, API/CLI,
compatibility, acceptance, and vulnerability effects for selected and every
serious candidate in disposable trees. Explain why the target exists in MVS,
whether a target package actually loads, and preserve every unrelated module
selection. Any required parent, major-path, floor, architecture, or unrelated-
module change needs a fresh bounded decision rather than silent
implementation. Compare primary vulnerability results at module, package,
symbol, test-symbol, and reachable-trace levels.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, latest Google UUID implementation identity, reciprocal archive
history, P7/P8 state, every unchanged qualification/exception guard, and
`./codex-dev-start.sh --check`. Read this archive, the answered go-hclog
decision/evaluation, go-cleanhttp evaluation, Errwrap decision/evaluation,
Consul SDK and earlier guarded-decision archives, rolling handover, roadmap,
`go.mod`, `go.sum`, and every referenced quality, compatibility, release,
runner, evidence, and lifecycle contract. Earlier outcomes are final.

# Three Moves

First, determine the highest qualified exact-path stable release whose complete
minimal source/test closure preserves Go 1.18. Do not promote an unqualified or
floor-ineligible identity. If no candidate satisfies the existing contracts,
preserve the evidence and stop for a fresh bounded product decision.

Second, qualify selected and serious candidates across source, tests, API,
behavior, loading, MVS, vulnerability, and project contracts. For an authorized
changed selection, use exact Go 1.26.7 and exact `go get` for one dependency-
only commit, never tidy as implementation, then run the complete P7 dependency
gate. For a retained or blocked selection, prove the no-change effect and run
all applicable gates without manufacturing activity. Full changed-selection
quality must preserve 27/27 Q0-Q2 PASS at L2 with zero held, regressed, non-
comparable, or dirty counts.

Third, update the roadmap and rolling handover with exact evidence, outcome,
commit identity, limitations, and next boundary. Answer this archive and
prepare one reciprocal NEXT mission only after the bounded outcome is coherent
and committed. Do not execute the successor.

# Automatic Handoff

After a completed coherent result, make the required local
`docs: prepare next agent session` commit. Do not launch a successor, push,
merge, publish, release, stash, revert, bypass cleanup, remove the worktree,
combine another dependency group, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No exact-path stable `github.com/hashicorp/go-immutable-radix` release
qualifies under the existing contracts. Exact selected v1.3.1 remains
inherited and unloaded, but it is not accepted by this evaluation. No source,
`go.mod`, or `go.sum` change was made, and no dependency implementation commit
exists. P7 stops at the reciprocal decision-only session linked above.

The evaluation began from clean branch `codex/upgrade-quality` at HEAD
`c30c70043847f1f4ed2f40c100ab3d5d6dc1bc6e`, parent
`7afe8f607fef3067c71d3cdd219041a97303dd20`, tree
`220a90e23a0bb33f3d98a511a012be8974d6c422`. The latest dependency
implementation remains Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`, parent
`37ab9ece6b0e5b1d735d3da83c3bc76adce7c2e3`, tree
`b89afd4ec056133b1eefb5611f8f35c12c11824b`.

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves canonical
public active unarchived non-fork MPL-2.0 repository
`https://github.com/hashicorp/go-immutable-radix.git`. The exact-path v1 proxy
has five stable releases, v1.0.0 through `@latest` v1.3.1, with no retraction,
module deprecation, redirect, or newer exact-path stable. Releases v2.0.0 and
v2.1.0 belong to the distinct `/v2` module and were not promoted. Selected
v1.3.1 is commit `49d1d02c49a783de548d1ba8ae8fde466a20b9e6`, parent
`f63f49c0b598a5ead21c5015fb4d08fe7e3c21ea`, tree
`b7c0058c779f65f8fb81f1db6e17023cbe8bfe0c`, dated
2021-06-28T20:38:56Z. Strict Git verification and all proxy/Git byte
comparisons pass.

All v1 tags are lightweight commit refs rather than signed tag objects. Exact
release identities are: v1.0.0 commit
`27df80928bb34bb1b0d6d0e01b9e679902e7a6b5`, parent
`7f3cd4390caab3250a57f30efdb2a65dd7649ecf`, tree
`14f852522e6b3d6c0572a26186f63ccb5b7275e0`, dated
2018-08-30T03:32:45Z; v1.1.0 commit
`7dd1121b595e4e1bd6dd5caa78e0f5c454740379`, parent
`5a8fb53eed279a0e955636bf3f03d4d19e6b49c4`, tree
`1b7b3a3854a456604e65afe5922156993c400bc9`, dated
2019-05-22T19:38:59Z; v1.2.0 merge commit
`e47f517af50356c81b117418760ac7dfd7b7af62`, parents
`0146a9aba1948ded4ed290cfd3fded2c15313f63` and
`56f5bef4b39cac512418ffbb2352d8ea69fbc343`, tree
`2b83140a1df8a20f8a07ab8a7caaf6777dc9f4c0`, dated
2020-03-18T18:06:31Z; and v1.3.0 commit
`57230d80bd59689cded197c5065835b39d890dd0`, parent
`67aa2d7a6ff40bae2f0ad354f3f6e7516a2921b8`, tree
`7862b0f94a6722b7b28da16472a563aae1cb021d`, dated
2020-09-17T14:26:30Z. GitHub marks v1.1.0, v1.2.0, and v1.3.1 commits verified;
v1.0.0 and v1.3.0 report `unknown_key`. V1.0.0 is an ancestor of v1.1.0;
v1.1.0 and v1.2.0 share merge base
`5a8fb53eed279a0e955636bf3f03d4d19e6b49c4` rather than direct ancestry due
the release-changelog fork; v1.2.0 -> v1.3.0 -> v1.3.1 -> v2.0.0 -> v2.1.0
is continuous ancestry. Current master and all other branch heads, including
the action-only `SECVULN-23537` branch, are unreleased or on `/v2` and were not
promoted.

Sumdb content sums are respectively
`h1:AKDB1HM5PWEA7i4nhcpwOrO2byshxBjXVn/J/3+z5/0=`,
`h1:vN9wG1D6KG6YHRTWr8512cxGOVgTMEfgEdSj/hr8MPc=`,
`h1:l6UW37iCXwZkZoAbEYnptSHVE/cQ5bOTPYG5W3vf9+8=`,
`h1:8exGP7ego3OmkfksihtSouGMZ+hQrhxx+FVELeXpVPE=`, and
`h1:DKHmCUm2hRBK510BaiZlwvpD40f8bJFeZnpfm2KLowc=`. Every release has the
same module-file sum
`h1:0y9vanUI8NX6FsYoO3zeMjhV/C5i9g4Q3DwcSNZ4P60=`.

Every v1 release omits a `go` directive and resolves only target production
dependency `github.com/hashicorp/golang-lru v0.5.0` plus test dependency
`github.com/hashicorp/go-uuid v1.0.0`. The complete imported closure preserves
Go 1.18. All five releases pass module verification, build, native count-one,
two count-ten repeats, race, vet, and Darwin AMD64, Linux AMD64/ARM64, Windows
AMD64, and js/wasm production/test cross-compilation under exact Go 1.26.7 and
contained Go 1.18.10.

Under Go 1.26.7 each release's production list has 47 packages with one
external production package; v1.0.0's complete-test list has 133 packages and
later releases have 134, with two external packages. Under Go 1.18.10 the
production list has 26 packages with one external package; v1.0.0 has 89
complete-test packages and later releases have 90, with two external packages.
The only imported external packages are `golang-lru/simplelru` in production
and `go-uuid` in tests.

Selected source contains one package, six production Go files, three tests,
and 16 files total, with no command, example, benchmark, modern fuzz target,
testdata, generated file, build-tag/platform branch, cgo, embed, generator, or
symlink. Pinned API comparison records compatible additions through v1.3.0;
v1.3.1 changes `ReverseIterator` from comparable to non-comparable. V1.0.0
removes six current exported APIs, including `Iterator.SeekLowerBound`, and is
not a compatible downgrade.

The independent behavior fixture SHA-256 is
`778d70a8c17156da17a0c6f32cf3f71b1b0033e4abb646d28516b4d3b089b32b`.
It passes deterministic ordering, CRUD/snapshot, prefix deletion, transaction
clone, watch lifecycle, arbitrary-byte key, immutable concurrent-read, nil/
panic, allocation, and aliasing checks under both SDKs. `Get` is allocation-
free in the measured case and `Insert` allocates 16 times under Go 1.26.7.
Values preserve caller identity; byte keys and returned keys alias caller or
internal storage and must remain immutable. Tree snapshots support concurrent
reads; transactions are explicitly not thread-safe, while separate clones are
independent. There is no process-global state or external resource lifecycle.

The separate minimal fixture SHA-256 is
`8e9476f3989e6ab4bed98ec99cc4409401b0d727982362b5de2374e4b92b6318`.
It proves that a valid exported sequence on an empty tree—missing
`SeekPrefix`, then `SeekLowerBound` on the same iterator—nil-dereferences and
panics under both SDKs. V1.1.0, v1.2.0, v1.3.0, and v1.3.1 all fail. Open
upstream issue #50 reports the same defect. V1.1.0-v1.3.0 additionally retain
the older issue-#28 prefix-key lower-bound panic, while v1.0.0 lacks the API.
Therefore there is no qualified exact-path stable candidate.

Selected v1.3.1 has selected-version incoming edges from Viper v1.15.0,
historical Viper v1.10.1, and `sagikazarmark/crypt v0.4.0`; lower requests
come from go-metrics and memberlist. Its why result is negative, repository Go
imports are zero, and the 355-entry production and 429-entry complete-test
loads contain zero target packages, so it is runtime-unreachable. A disposable
exact v1.3.1 get preserves all 234 selections and zero target load but
manufactures target and golang-lru roots plus two checksums; tidy returns the
exact base projection. Exact v1.0.0-v1.3.0 gets instead downgrade Viper to
v1.9.0, move golang-lru, remove mvn-pom-mutator, collapse the graph to 180
modules, and break project loading. None was applied.

Fresh primary vulnerability data remains byte-identical at 1,399 records,
SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`,
and Last-Modified 2026-09-15T18:54:35Z. It has no target record, and exact OSV
queries for all five v1 releases are empty. Go 1.26.7 isolated module/package/
symbol/test-symbol scans are empty. Go 1.18.10 reports only old-SDK standard-
library advisories. Base and disposable-v1.3.1 project scan populations are
identical with zero target package, symbol, test-symbol, or reachable trace.

All twelve existing guarded modules retain exact selections, recorded incoming
edges, negative why results, zero repository imports, zero production/test
loads, runtime unreachability, and prior advisory state. Gorilla alone retains
its recorded GO-2026-6278/GHSA-w67g-5rqw-f597 exact-version result and
GO-2020-0019 primary-index entry.

The project remains byte-for-byte unchanged at 234 modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, 1,067 sum lines, and the 432-line tidy projection.
`go.mod`/`go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Exact Go 1.26.7 verification, build, count-one, two count-ten repeats, race,
vet, pinned lint, API/CLI compatibility, empty-HOME count-two, four cross-
builds, host and snapshot acceptance, all 17 script meta-tests, all eight
10-mutant populations with 80/80 kills, all 15 audit controls, launcher-auto,
and related meta-controls pass. Contained Go 1.18.10 loads 366 complete-test
entries; its 26 unaffected packages pass count-one, both repeats, race, vet,
and four cross-builds, while the full suite retains only the two accepted
Darwin `pkg/shell` wording failures.

Exact Go 1.26.7 archive/binary SHA-256 values are
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
Go 1.18.10 values are
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
The golangci-lint 2.12.2 archive and GoReleaser 2.17.1 binary retain required
portable SHA-256 values
`a9c54498731b3128f79e090be6110f3e5fffccc617b08142ed244d4126c73f29` /
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`.
Freshly rebuilt govulncheck v1.8.0 SHA-256 is
`6cd9dbf44dc2668b90e8396b4fab9f4ba166a916f1f832cf5bcdffbd0919872a`;
its build receipt is nonportable. The rebuilt apidiff SHA-256 is
`206c308c3b84c57f8437cf7e1d82103e045922286c368ecbbcd03be760934a81`,
preserving the known archive-reproducibility discrepancy from the prior
portable receipt.

Six launcher lifecycle attempts reproduced only the known signal/log-retention
timing race after controls 1-25; one nested attempt also reported its known
second-generation-header consequence. `./codex-dev-start.sh --check` and the
independent launcher-auto contract pass. The initial inherited `umask 077`,
the scratch-only bare-`mktemp` wrapper, the apidiff build-receipt discrepancy,
and this launcher timing race are preserved environment/harness limitations,
not target findings. Because dependency and source metadata did not change,
no changed-selection quality run applies and accepted 27/27 Q0-Q2 PASS at L2
remains authoritative. The 435-entry disposable evidence manifest SHA-256 is
`8a30047f3bf22c7ad7e33974aee395c26dc5c10cbd37c9f6dd03daafc1e4b939`.
