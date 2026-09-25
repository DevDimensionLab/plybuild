# Agent Session: Evaluate Posener Complete Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T082105+0200-evaluate-posener-complete-dependency`
Created: `2026-09-23T08:21:05+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `f9a49422a5104c87bfd224d3e7bbffae4fa35192eda7c11a6c2e7429ff089098`
Previous: [2026-09-23T074006+0200-decide-pmezard-go-difflib-product-direction.md](2026-09-23T074006+0200-decide-pmezard-go-difflib-product-direction.md)
Next: [2026-09-23T092359+0200-decide-posener-complete-product-direction.md](2026-09-23T092359+0200-decide-posener-complete-product-direction.md)
Outcome: No canonical exact-path stable qualifies: v1.1.x has incomplete module metadata, and every buildable v1.2 stable leaks the temporary file created by exported Uninstall. Product source and dependency metadata remain unchanged; one bounded Complete product decision was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating the next unanswered selected queue
item, graph-selected transitive exact `github.com/posener/complete v1.2.3`, as
exactly one bounded dependency group. Resolve its canonical exact-path release
line and highest qualified Go-1.18-compatible stable from primary evidence.
Implement one exact dependency-only changed selection only if the candidate,
its complete minimal closure, and every earlier target-specific guard remain
exact. Do not combine another dependency group or begin P8.

# Defensive Scope

This is an ordinary dependency-quality evaluation. Use public metadata,
static source/repository records, project graph/build commands, upstream tests,
and only small bounded ordinary fixtures required by documented behavior. Do
not fuzz, stress, probe resource exhaustion, create oversized, deeply nested,
cyclic, malformed, adversarial, or escape-sequence payloads, reproduce a
security issue, or perform security or exploitability analysis.

Every disposable cache, tool, archive, report, project copy, fixture, or
advisory response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`.
Never write to `/private/tmp`, `/tmp`, a sibling of the managed root, or
another external root. Verify containment and remove task-owned scratch
evidence before handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
all final target-specific exceptions through go-difflib, and qualified no-
selection-change modern-go/concurrent and modern-go/reflect2. P8 remains
queued.

Exact selected, inherited, indirect, unloaded
`github.com/pmezard/go-difflib v1.0.0` is retained only under its explicit,
go-difflib-specific, unqualified, non-transferable exception. It is the sole
and latest stable but is neither qualified, supported, safe, nor fixed:
upstream ended maintenance, both exact SDKs fail complete vet, and its two
exported diff writers discard final buffered `Flush` errors. Preserve all 23
exact requests, Testify's nine genuine import boundaries, all 14 metadata-only
requester boundaries, every current/historical route, the positive target why
through direct go-term-markdown's dependency test, zero repository import/
load/runtime/root state, all release/source/behavior/closure/projection/
advisory identities, and every expiry condition. Do not run the rejected Go-
Term-Markdown Testify Go-Difflib Owner/Request Study, add a root, patch/vendor
upstream, change an owner/requester, or transfer the exception.

Exact selected, inherited, indirect, unloaded `github.com/pkg/sftp v1.13.1`
and `github.com/pkg/errors v0.9.1` remain retained only under their own
explicit unqualified non-transferable exceptions. Preserve their exact
requests, genuine owner/import boundaries, routes, root/why/import/load/
runtime facts, source/behavior/closure/projection/advisory identities, mutual
route guard, and expiry conditions. Do not run their rejected owner studies,
reopen them, or transfer either exception.

Exact Goe, ULID, go-conntrack, mapstructure, go-homedir, Promptui, emoji/v2,
kr/text, kr/pty, kr/pretty, kr/logfmt, kr/fs, go-windows-terminal-sequences,
gotool, errcheck, httprouter, GLS, go-junit-report, json-iterator, and
clockwork exceptions remain final, target-specific, unqualified, and non-
transferable. Every concurrent, reflect2, Cast, Viper, memberlist, earlier
selection, owner, request, route, why/import/load/runtime, graph/module/tidy/
Go-floor, advisory, source, behavior, closure, qualification, and expiry guard
remains final. Do not reopen, broaden, or transfer any decision.

# Measurements At Start

The go-difflib decision began from clean branch `codex/upgrade-quality` at
evaluation-handoff HEAD `4e6198c4c3ff5315928df9f6d0690de985bf9ce3`, parent
`632840058974b19122e27302fde780ed6739510b`, tree
`30b7ea878359482c0f2ff1495daf35c5f0741bc2`. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an
ancestor. The decision changed no product source or dependency metadata and
prepared this evaluation-only handoff. Verify the new handoff HEAD, parent,
tree, exact changed set, ancestry, reciprocal archive chain, and clean
ordinary and ignored status rather than assuming them.

The unchanged project remains 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The graph SHA-256 is
`abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
The 432-line tidy diff SHA-256 is
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes are
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

All 47 pre-Goe guarded selections and 276 incoming edges remain exact at
`7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
The separate complete Goe, pkg/errors, SFTP, and go-difflib guards remain
exact. Preserve every earlier selection and target-specific qualification or
exception boundary. Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Guard OSV remains limited to Gorilla WebSocket and go-retryablehttp; x/mod
v0.14.0 retains GO-2026-6179 and GO-2026-6180. The Go vulnerability index
remains 518,501 bytes/1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Advisory absence is not qualification.

The queue identifies exact `github.com/posener/complete v1.2.3` as the next
selected module after go-difflib. Current guard-only observations show one
exact selected request, `github.com/hashicorp/serf v0.9.6` -> target, and the
shortest selected route main -> direct mvn-pom-mutator v0.2.3 -> historical
Viper v1.10.1 -> Serf v0.9.6 -> target. Target why is negative; repository
imports, loaded target entries, runtime relevance, and current/history roots
are zero. Treat these as starting observations only. Independently reproduce
every request, current or historical route, requester genuine-import or
metadata-only boundary, target/requester why, repository import, production/
complete-test/module-backed load, runtime relevance, and root/history-root
fact before selecting a candidate. Do not infer qualification, ownership, or
runtime relevance from physical graph selection or a route.

# Role And Boundaries

Resolve exact go-import and module-path identity; repository owner/status/
license/default branch; version/tag/release, commit/tree/signature/ancestry,
proxy/sumdb/archive-to-Git identity; module directives and requirements;
retractions, deprecation, replacements, and exact-path major lines. Consider
only genuine exact-path stable releases. Do not promote a fork, branch,
pseudo-version, prerelease, replacement, alternate module path, or ownerless
candidate as a stable.

For selected and every serious stable candidate, inspect the complete module
and test closure, exported API and documentation, Go-floor compatibility,
platform/build-tag/cgo/generated/embed boundaries, globals, ownership and
mutation, determinism, concurrency, lifecycle, cleanup, and error behavior.
Exercise only small bounded ordinary values needed to verify documented
behavior. Run upstream build, tests, repeated tests, race, vet, and supported
cross-builds under exact Go 1.26.7 and a contained Go 1.18 toolchain. A
release qualifies only if every applicable ordinary documented contract and
every project guard pass.

Map every target MVS request and genuine current or historical route.
Reproduce target and requester why, repository imports, production and
complete-test loads, module-backed entries, runtime relevance, graph counts,
hashes, tidy projection, all earlier guarded selections, the 276-edge pre-Goe
snapshot, and the separate Goe, pkg/errors, SFTP, and go-difflib guards.
Physical selection, transitive presence, a route, loading, or advisory absence
is not qualification.

Use disposable project copies beneath the managed scratch root to measure
exact candidate projections. Never add or alter a target root in the real
project outside the one final exact dependency implementation authorized
below. Record selection, closure, graph, imports/loads, sums, tidy result,
Go-floor effect, genuine supported tidy-stable ownership, and every earlier
guard for each projection. Do not retain a projection unless the candidate
qualifies and the normal dependency implementation contract authorizes it.

Refresh exact-version OSV and GitHub advisory evidence, repository advisories,
the guarded advisory population, x/mod guard, Go-index identity, memberlist
CNA identity, and a pinned isolated govulncheck comparison. Advisory absence
cannot override ordinary behavior, an upstream gate, ownership, or an earlier-
guard failure. Stay within the defensive scope.

# Decision And Implementation Boundary

Select only the highest qualified Go-1.18-compatible exact-path stable with a
genuine supported project owner. If that exact selection changes and every
earlier guard remains exact, use exact Go 1.26.7 and exact
`go get github.com/posener/complete@<selected-version>` for one dependency-
only commit; do not hand-edit metadata and do not use tidy as the
implementation. Explain and verify the minimal exact transitive closure.

If no genuine supported changed selection qualifies, do not retain a direct
root, redundant source sum, downgrade, replacement, branch, fork, pseudo-
version, prerelease, alternate path, patch, vendor copy, workaround, or
unrelated metadata churn. Record the completed result and prepare exactly one
reciprocal product-decision successor offering only target-specific
unqualified retention, one precisely bounded later measurement-only owner/
request study if justified by evidence, or stopping P7 unresolved. Do not
silently retain an unqualified target or transfer an exception.

# Required Reading And Verification

Read this archive, the answered go-difflib decision and evaluation, answered
SFTP and pkg/errors decisions/evaluations, answered Goe, ULID, and go-
conntrack decisions/evaluations, reflect2 and concurrent evaluations,
mapstructure, go-homedir, Promptui, emoji/v2, kr/text, and kr/pty decisions/
evaluations, rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Earlier
evaluation facts are final unless a guard-only check proves an input changed;
stop for a fresh owning evaluation if one did.

Before implementation, verify branch/HEAD/parent/tree, exact changed set,
ancestry, ordinary and ignored cleanliness, reciprocal archive chain,
launcher check, exact Go identities, project hashes/counts/tidy state,
selection/requests/routes/import/load/runtime facts, all earlier guards, and
fresh advisories. After any changed selection, run focused consumer contracts,
module verification, build, count-one and repeated tests, race, vet, supported
cross-builds, pinned lint, API/CLI/help/launcher/Make/quality contracts,
accepted Q0-Q2 quality, vulnerability comparison, and exact tidy analysis.

# Three Moves

First, independently resolve canonical releases, qualification, ownership,
routes, closures, projections, advisories, and earlier guards. Second, make at
most one exact dependency-only implementation if the highest supported
candidate fully qualifies; otherwise leave source and dependency metadata
unchanged. Third, update roadmap and rolling handover, answer this archive,
prepare exactly one reciprocal successor required by the result, verify
scratch containment and cleanup, and make the local handoff commit.

# Automatic Handoff

Do not launch a successor, run an owner study, push, merge, publish, release,
stash, revert, bypass cleanup, remove the worktree, change another dependency,
transfer an exception, reopen go-difflib, SFTP, pkg/errors, or an earlier
group, write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No Posener Complete selection is authorized. Exact selected/latest v1.2.3 is
the highest canonical exact-path stable and has a genuine project owner, but
it is not qualified: its exported uninstall path leaves the temporary file it
creates after an ordinary successful operation. Every buildable v1.2 stable
has the same cleanup failure, while v1.1.1 and v1.1.2 additionally have
incomplete synthesized module metadata and fail the complete closure gates.
No product source, `go.mod`, or `go.sum` change was retained.

### Continuity and exact project state

- The evaluation began clean on `codex/upgrade-quality` at
  `7c126d36f09a63f25a48f120f6987d7cba87d1a8`, parent
  `4e6198c4c3ff5315928df9f6d0690de985bf9ce3`, tree
  `8d8f348847000689280619b47e054698b953e29b`. That handoff changed exactly
  the launcher, answered go-difflib decision, this then-NEXT archive, rolling
  handover, and roadmap. Google UUID implementation commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. Ordinary
  and ignored status were clean; the reciprocal archive chain and launcher
  check passed.
- The real project remains 234 selected modules, 3,599 graph edges, 355
  production entries, 429 complete-test entries, 197 module-backed entries
  across 41 loaded modules, and 1,067 sum lines. `go.mod`, `go.sum`, and graph
  SHA-256 remain respectively
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  Exact tidy output remains the common 52/948-line state with hashes
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479`
  and `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`;
  it is the recorded 432-line unapplied projection. The differing local
  431-line `git diff` presentation did not change either applied file byte and
  therefore was not treated as input drift.
- Official Go 1.26.7 Darwin-arm64 archive/binary SHA-256 is
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
  Contained official Go 1.18.10 archive/binary SHA-256 is
  `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.

### Canonical identity and release line

- Exact go-import metadata maps `github.com/posener/complete` to
  `https://github.com/posener/complete.git`. GitHub repository ID 90418143 is
  public, enabled, unarchived, non-fork, owned by `posener`, MIT licensed, and
  uses default branch `v1`. That branch explicitly says current development
  occurs on `master` and remains the default for older libraries and
  compilers. The selected MIT `LICENSE.txt` hashes to
  `42707f6d0ca72916d2a48dbc977aab20a5e23270a7651bd18227e992f9d41548`.
- Proxy and `go list -m -versions` expose exactly six canonical exact-path
  stables: v1.1.1, v1.1.2, v1.2.0, v1.2.1, v1.2.2, and selected/latest
  v1.2.3. The repository also has old `v1.0` and `v1.1` release objects whose
  two-component tags are not canonical module versions. The separate
  canonical `/v2` line has stable v2.0.0 and v2.1.0 plus prereleases and is an
  ineligible alternate module path; `/v3` is absent. There are no exact-path
  prereleases, replacements, retractions, or deprecations.
- All six canonical tags are lightweight valid-signed commits and ancestors
  of `v1` head `9a4745ac49b29530e07dc2581745a218b646b7a3`. Their exact commits/trees are:
  v1.1.1 `98eb9847f27ba2008d380a32c98be474dea55bdf` /
  `026be229c487263c2d437ff4f486e9c841b81a9c`; v1.1.2
  `dcda3199365ca2a5f24aea4c42aa56f6a197d117` /
  `82f634a35a9bf2c0ad95e48f50886f1e37ccc60b`; v1.2.0
  `ffc2cf5e958af5ac4156a886153c7cadd24a521a` /
  `1b4936b7ce044bd52e1ae3f9e4cf7c41a2a89a35`; v1.2.1
  `3ef9b31a6a0613ae832e7ecf208374027c3b2343` /
  `e5efaab6e25c864267e8297ffe3da8e5684e654d`; v1.2.2
  `98a0c28ec7908620d626d072df3a49c34eadc2b6` /
  `f2a8af9bf26cb8fca022da1246268969adc2d7ef`; and v1.2.3
  `05b68ffc813dd10c420993cb1cf927b346c057b8` /
  `58ea2366a478ec75ecf97d4d30097a3f16709fab`.
- Selected sumdb source/module hashes are
  `h1:NP0eAhjcjImqslEwo/1hq7gpajME0fTLTezBKDqfXqo=` and
  `h1:WZIdtGGp+qx0sLrYKtIRAruyNpv6hFCicSgv7Sy7s/s=`. Its proxy ZIP SHA-256 is
  `88b48005b995dc6592fa6fda08130488c83f63bcaa4ccb0fb8e926fee63112ec`.
  All 41 regular files match Git exactly; the normalized archive/Git manifest
  hashes to
  `2288a98d97813d7aebf8fd745949edd58c137119587d53e949d497f16df71ebb`.
  No symlink, submodule, or non-regular Git object is present.
- V1.1.x receives only proxy-synthesized module metadata and omits its
  imported go-multierror dependency. V1.2.0/v1.2.1 require go-multierror
  v1.0.0 and have no Go directive. V1.2.2/v1.2.3 declare Go 1.13 and require
  go-multierror v1.0.0 plus Testify v1.4.0. V1.2.3 differs from v1.2.2 only by
  deleting the nested `gocomplete` module files, returning that command to the
  main module.

### Requests, routes, and runtime boundary

- Exactly three current/historical requests target the module path:
  `hashicorp/serf@v0.9.6 -> complete@v1.2.3`,
  `mitchellh/cli@v1.1.0 -> complete@v1.1.1`, and
  `mitchellh/cli@v1.0.0 -> complete@v1.1.1`. Serf's request is metadata-only;
  both CLI versions genuinely import Complete from production command and
  autocomplete source, with an additional test import in v1.1.0.
- The shortest selected route is main -> direct mvn-pom-mutator v0.2.3 ->
  historical Viper v1.10.1 -> Serf v0.9.6 -> target. One CLI route continues
  from that Serf vertex through CLI v1.1.0; the other is main -> direct
  mvn-pom-mutator -> crypt -> Consul API v1.1.0 -> Serf v0.8.2 -> CLI v1.0.0
  -> target. Every path therefore enters through the existing direct
  mvn-pom-mutator owner.
- Target, selected Serf v0.10.1, and selected CLI v1.1.0 why results are
  negative. The repository has zero target, Serf, or CLI imports. Their
  packages are absent from all 355 production and 429 complete-test entries;
  the target contributes zero module-backed entries and has zero runtime
  relevance. Neither `go.mod` nor repository history contains a target main-
  module root.

### Source, closure, and qualification

- Selected v1.2.3 has a nine-module build list. Its loaded external closure is
  errwrap v1.0.0, go-multierror v1.0.0, go-spew v1.1.0, go-difflib v1.0.0,
  Testify v1.4.0, and yaml.v2 v2.2.2; objx and check.v1 remain build-list-only.
  Exact Go 1.18.10 closure counts are 76 production, 173 complete-test, and 18
  module-backed entries across seven loaded modules; exact Go 1.26.7 gives
  104/237/18 across the same seven. Its exported API/documentation snapshot is
  316 lines at SHA-256
  `f8dea804265e7dda3261d21129dd5e667f09b5de8c674d5dec010d2611df10b2`.
- The module is pure Go with no cgo, build tags, generated source, or embed.
  It uses caller-owned command maps and writers plus exported mutable `Log`;
  callers must synchronize shared mutation. Map and filesystem enumeration
  make candidate ordering unspecified. Install helpers own shell-configuration
  edits, file predictors read the filesystem, and the command intentionally
  exits the process. These boundaries are not project runtime-relevant because
  the module is unloaded.
- V1.1.1 and v1.1.2 pass module verification but fail build, tests, race, and
  vet under both exact SDKs because their synthesized `go.mod` omits imported
  `github.com/hashicorp/go-multierror`. V1.2.0 through v1.2.3 pass module
  verification, build, count-one/count-ten tests, race, and vet under both
  exact SDKs. Each also passes production and test compilation with cgo
  disabled for Darwin amd64/arm64, Linux amd64/arm64/386, Windows amd64/386,
  FreeBSD amd64, Plan 9 amd64, and js/wasm.
- Every buildable v1.2 release has byte-identical uninstall implementation.
  Exported `cmd/install.Uninstall` calls `removeContentToTempFile`, which uses
  `ioutil.TempFile("/tmp", "complete-")` and returns its name. The caller
  copies that file over the original and removes only the backup, never the
  temporary file. A successful documented operation therefore leaks an owned
  temporary file. Close errors are also discarded. The leaking path was not
  executed because its hard-coded destination lies outside managed scratch;
  static source identity is complete evidence. This ordinary cleanup failure
  disqualifies v1.2.0 through v1.2.3 despite their passing upstream gates.

### Projection, guards, advisories, and handoff

- Disposable exact `go get github.com/posener/complete@v1.2.3` changes no
  selected version but manufactures redundant indirect roots/main edges for
  Complete v1.2.3, go-multierror v1.1.0, and errwrap v1.0.0 plus their three
  source sums. The raw projection has 234 modules, 3,602 edges, unchanged
  355/429/197/41 loads, 1,070 sum lines, negative target why, and `go.mod` /
  `go.sum` / graph hashes
  `1b19b99f39fa76820e40902d92005b59eb1e4661ddf392e361834450484088dd`,
  `d667b87d804b85e47f340f9cf40293f90436894aa59aff96e10cad89aa80d1b1`,
  and `83a4fd27ba9d3fc203706f4e448ff26b201c4714d752ed2713e84e6be0fc94d8`.
  Normal tidy removes the manufactured state and returns the exact common
  234-module/3,557-edge 52/948-line projection. Nothing was retained.
- All 47 pre-Goe selected versions remain exact; their 276 sorted incoming
  edges reproduce SHA-256
  `7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
  Forty-one guarded why results remain negative; only kr/pretty, kr/text,
  emoji/v2, Promptui, go-homedir, and mapstructure are positive. Promptui and
  go-homedir remain the only guarded repository imports; emoji/v2, Promptui,
  go-homedir, and mapstructure remain the only loaded guarded modules. The
  separate exact Goe v0.1.0/four-request, pkg/errors v0.9.1/ten-request, SFTP
  v1.13.1/four-request, and go-difflib v1.0.0/23-request guards and every
  earlier final decision remain exact.
- Exact-version OSV is empty for all six canonical target stables at the
  two-byte `{}` SHA-256
  `44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a`.
  Exact GitHub global and repository advisory responses are empty at the
  two-byte `[]` SHA-256
  `4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945`.
  Pinned govulncheck v1.8.0 under exact Go 1.26.7 uses database timestamp
  2026-09-16T18:00:43Z. Its isolated target module/package/symbol/test-symbol
  scans are empty; the unchanged project reproduces 30 module, 22 package, 20
  symbol, and 20 test-symbol IDs with no target occurrence.
- Guard OSV retains Gorilla WebSocket `GHSA-w67g-5rqw-f597` /
  `GO-2026-6278` and go-retryablehttp `GHSA-v6v8-xj6m-xwqh` /
  `GO-2024-2947` at their recorded response hashes; x/mod 0.14.0 retains
  `GO-2026-6179` and `GO-2026-6180` at
  `1398fe7aa760157dca6f051fd6e8ac6b760627923283bb61c53cb14f1409536e`.
  The Go vulnerability index remains 518,501 bytes/1,402 records at
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  Last-Modified 2026-09-17T17:29:18Z. The PUBLISHED memberlist CNA response
  remains 2,807 bytes at
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
  Advisory absence was not used as qualification.
- Final unchanged-project exact-Go-1.26.7 module verification, build,
  count-one tests, race count-one tests, and vet pass under ordinary umask 022.
  Exact source, API/CLI/help/launcher/Make/quality contracts and accepted 27/27
  Q0-Q2 PASS at L2 remain unchanged.
- Exactly one reciprocal successor is prepared and not executed. It offers
  only Complete-specific unqualified retention of exact v1.2.3, one later
  measurement-only **Mvn-Pom-Mutator Complete Owner/Request Study**, or
  stopping P7 unresolved. P8 remains queued.
