# Agent Session: Evaluate Kr Pretty Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-22T100435+0200-evaluate-kr-pretty-dependency`
Created: `2026-09-22T10:04:35+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `39f40cba60a059162c4143bcbfa9158fdeeea7e7f650bf657ec465c0088d6deb`
Previous: [2026-09-22T092932+0200-decide-kr-logfmt-product-direction.md](2026-09-22T092932+0200-decide-kr-logfmt-product-direction.md)
Next: [2026-09-22T110645+0200-decide-kr-pretty-product-direction.md](2026-09-22T110645+0200-decide-kr-pretty-product-direction.md)
Outcome: No exact-path stable qualifies. Product source and dependency
  metadata remain unchanged. A bounded reciprocal product decision was
  prepared but not executed.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected exact-path
`github.com/kr/pretty v0.3.1` as one bounded dependency group. Resolve the
canonical release line and highest qualified Go-1.18-compatible exact-path
stable from primary evidence. Implement one exact dependency-only changed
selection only if the candidate and its exact closure preserve every contract.
If no exact-path stable qualifies, leave source and dependency metadata
unchanged and prepare one bounded product decision. Do not combine kr/pty,
kr/text, another dependency group, or P8.

# Defensive Scope

This is an ordinary dependency quality evaluation. Use public metadata, static
source/repository records, project graph/build commands, upstream tests, and
small bounded ordinary formatting fixtures. Do not fuzz, stress, probe
resource exhaustion, generate oversized/deep/cyclic or adversarial values,
reproduce a security issue, or perform security/exploitability analysis.

Every disposable cache, tool, archive, report, project copy, fixture, or
advisory response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never
write to `/private/tmp`, `/tmp`, a sibling of the managed root, or another
external root. Verify containment and remove task-owned scratch evidence
before handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific decisions through exact selected, inherited, unloaded
kr/logfmt. P8 remains queued.

Option 1 explicitly retains exact
`github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515` as inherited and
unloaded under its own unqualified, non-transferable exception. Its three
ordinary behavior failures, three requests, genuine mvn-pom-mutator/Viper/
Prometheus/go-logfmt routes, requester-import facts, root/import/load/runtime,
graph/module/tidy/Go-floor, earlier-guard, and advisory/release/owner/
qualified-route conditions are final. No kr/logfmt, kr/fs, go-windows-
terminal-sequences, gotool, errcheck, httprouter, jtolds/gls, go-junit-report,
json-iterator, clockwork, demangle, strcase, memberlist, or other exception
transfers to kr/pretty.

# Measurements At Start

The kr/logfmt decision began from clean HEAD
`1017a890db0c7ac42c6ebb4bc57bc8c191ac718c`, parent
`f7d5b3a94833a92b5d5b5a48b09c33e2dcba9e18`, tree
`3bea69f9fa766139c87af1f50fbfcf627a5e71e3`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`. The decision changed no source or
dependency metadata and prepared this evaluation-only handoff. Verify the new
handoff HEAD, parent, tree, exact changed set, ancestry, and clean ordinary and
ignored status rather than assuming their values.

Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The project remains 234 selected modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, 1,067 sums, and the recorded 432-line tidy projection. All 36 guarded
selections including kr/logfmt and 228 incoming edges remain exact at sorted
SHA-256
`361d0c355a69842518a518e682f81c9728d37acfdeff64f430a4fb253929691e`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

Current queue observations select kr/pretty v0.3.1. Exact requests currently
observed are Consul SDK v0.8.0 -> v0.2.0; Prometheus client_golang v1.4.0,
sergi/go-diff v1.2.0, and errgo.v2 v2.1.0 -> v0.1.0; and selected Cast v1.5.1
-> v0.3.1. `go mod why -m` is positive only through the project config package
-> yaml.v2 -> yaml.v2 tests -> check.v1 -> kr/pretty. Repository imports and
production/complete-test package loads of kr/pretty are zero. These are queue
observations, not qualification; independently reproduce them and determine
all genuine current and historical owner routes and runtime relevance.

Cast v1.5.1 was previously accepted with the minimal quicktest v1.14.4,
kr/pretty v0.3.1, and rogpeppe/go-internal v1.9.0 test closure. That Cast
decision is final under its guards and does not qualify kr/pretty or authorize
reopening Cast.

# Evaluation Contract

Resolve exact go-import, repository/owner/status/license, version/tag/release,
commit/tree/signature/ancestry, module-path/directive/requirement/retraction/
deprecation/replacement, proxy/sumdb/archive-to-Git identity, default branch,
and exact-path major-line evidence. Consider only genuine exact-path stable
releases. Do not promote a fork, branch, pseudo-version, replacement,
alternate module path, or ownerless candidate.

For selected and every serious stable candidate, inspect the complete module
and test closure, exported API and documentation, Go-floor compatibility,
platform/build-tag/cgo/generated-source boundaries, globals, ownership and
mutation, determinism, concurrency, lifecycle, external boundaries, and
ordinary error behavior. Exercise only small bounded ordinary pretty-
formatting values needed to verify documented behavior. Run upstream build,
tests, repeated tests, race, vet, and supported cross-builds under exact Go
1.26.7 and a contained Go 1.18 toolchain. A release qualifies only if all
ordinary documented contracts and every project guard pass.

Map every MVS request and genuine route, including the accepted Cast v1.5.1
closure and lower historical requests. Reproduce why, repository imports,
production and complete-test loads, module-backed entries, runtime relevance,
graph counts, hashes, tidy projection, all 36 earlier guarded selections, and
the 228-edge snapshot. Physical selection, a why chain, zero loading, or
advisory absence is not qualification.

Use disposable project copies under the managed scratch root to measure exact
candidate projections. Never add a target root to the real project. For each
projection record exact selection, closure, graph, imports/loads, sums, tidy
result, Go-floor effect, and whether a genuine supported tidy-stable owner
exists. Do not retain a projection unless the exact candidate qualifies and
the normal dependency implementation contract authorizes it.

Refresh exact-version OSV and GitHub advisory evidence, repository advisories,
the guarded advisory population, x/mod guard, Go vulnerability index identity,
memberlist CNA identity, and pinned isolated govulncheck comparison. Advisory
absence cannot override an ordinary behavior failure. Stay within the
defensive scope.

# Decision And Implementation Boundary

Select only the highest qualified Go-1.18-compatible exact-path stable with a
genuine supported project owner. If that exact selection changes, use exact Go
1.26.7 and exact `go get github.com/kr/pretty@<selected-version>` for one
dependency-only commit; do not hand-edit metadata and do not use tidy as the
implementation. Explain and verify the minimal exact transitive closure and
preserve every earlier selection and exception.

If selected v0.3.1 is already the highest qualified exact decision and exact
get is a no-op, record the no-change decision without forcing a dependency
commit. If no stable qualifies, or no qualified candidate has a genuine
supported tidy-stable owner, retain no projection and prepare one reciprocal
product-decision archive offering only: a target-specific unqualified
exception, exactly one named later measurement-only genuine owner/request
study, or stop P7 unresolved. Do not choose that product direction during the
evaluation.

After any changed selection, run the complete P7 dependency gate required by
the roadmap. For an unchanged result, run focused target/graph/advisory guards
and final exact-Go module verification, build, count-one tests, race count-one
tests, and vet. In either case preserve the Go 1.18 floor, source/API/CLI/help/
launcher/Make/quality contracts, all earlier guards, and accepted 27/27 Q0-Q2
PASS at L2.

# Required Reading

Read this archive, the answered kr/logfmt decision/evaluation, kr/fs,
go-windows-terminal-sequences, gotool, and errcheck decisions/evaluations, the
accepted Cast v1.5.1 evaluation and related requester/retained-module records,
rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Verify branch,
ancestry, clean ordinary/ignored state, reciprocal archive chain, launcher
check, exact Go identities, module hashes/counts/tidy projection, all 36
earlier guards and 228-edge snapshot, and fresh advisory identities before any
implementation.

# Three Moves

First, independently evaluate only kr/pretty and choose the exact qualified
stable or the bounded no-qualified result. Second, make at most the one
authorized dependency-only selection change and verify it, or leave source and
metadata unchanged. Third, update roadmap and rolling handover, answer this
archive, prepare exactly one reciprocal successor required by the result,
verify containment, and make the required local handoff commit without
executing the successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert,
bypass cleanup, remove the worktree, transfer an exception, reopen Cast or an
earlier decision, evaluate another dependency group, write outside the managed
scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Outcome

No exact-path `github.com/kr/pretty` stable qualifies. Exact selected v0.3.1
therefore remains physically inherited and unloaded, but this evaluation does
not accept or qualify it. Product source, `go.mod`, and `go.sum` are unchanged;
no dependency implementation commit or changed-selection scorecard exists.

### Canonical Release And Module Evidence

- Exact `go-import` metadata resolves without redirect to
  `https://github.com/kr/pretty.git`. GitHub repository ID 1253164 is public,
  enabled, unarchived, non-fork, MIT licensed, owned by `kr`, and defaults to
  protected `main`. It has no GitHub Releases. The v0.3.1 commit and default
  branch are identical. The MIT license SHA-256 is
  `283c28b781bc7487bd7bcaf453cbd811f5ee8d7ce625cfcf8dcacae6ba1f4bb5`.
- The canonical proxy exposes exactly v0.1.0, v0.2.0, v0.2.1, v0.3.0, and
  v0.3.1; proxy `@latest` is v0.3.1 at
  `3cd153a126da607b78d1762779b1e1054f9889fc`. Exact `/v2` and `/v3`
  endpoints are absent. No branch, pseudo-version, fork, replacement, or
  alternate path was promoted.
- Each stable is an annotated unsigned tag object. Tag object -> peeled commit
  identities are v0.1.0
  `c21352ee25c82043780a7b54cc9b2fbe60411457` ->
  `73f6ac0b30a98e433b289500d779f50c1a6f0712`; v0.2.0
  `c18e17649ccc296e05303470af2af293597bd9c9` ->
  `4e0886370c3a67530192c6a238cff68f56c141b0`; v0.2.1
  `326f0883a47c4b88a699dc12c4ca20164aadd512` ->
  `ead452280cd055b2ae8a7f0db5eb37a878d902f7`; v0.3.0
  `6259398fb07372ee10ef52d88375fdab383bcf9d` ->
  `a883a8422cd235c67c6c4fdcb7bbb022143e10b1`; and v0.3.1
  `5da5474ed8ef553af2f50611599087057f3c82d4` ->
  `3cd153a126da607b78d1762779b1e1054f9889fc`. The peeled commits form
  continuous ancestry. GitHub reports the final three commits verified and
  the first two unsigned.
- Commit tree identities in release order are
  `6100dd6119afb15d1a87c990b440abe6e8295b0c`,
  `d2811c37c4687092b2848146f50ac853eb09883c`,
  `21980971d1498b4f1837c51bd1b84a9050550f38`,
  `fa6e8883156bdc9f2deb59112a2a2ed28b11f956`, and
  `cc4e9ed3d6f83719e804f38b1712decf004f43a9`. Their parents are
  `cfb55aafdaf3ec08f0db22699ab822c50091b1c4`,
  `71e7e49937503c662b9b636fd6b2c14b1aa818a5`, the v0.2.0 commit,
  `59b4212b823be760588121cf7a7aee9f6d6c0cad`, and
  `d8c7eb13a5c499cc7ccbeb58b3421b42c599bfae` respectively.
- Proxy ZIP SHA-256 values in release order are
  `06063d21457e06dc2aba4a5bd09771147ec3d8ab40b224f26e55c5a76089ca43`,
  `006c8a1a9fbd487942eb43ce663e6e2770dbeea8e79efc1c6f7dbeaf3ca52e19`,
  `80af0452082052d1b3265d7cb8985d464d4be222c27e14658e95632c222761e5`,
  `3ac65e185f956d889d77485173fadcc30e959b6bcfdaa8acafaec5f4dac5cd48`,
  and
  `ecf5a4af24826c3ad758ce06410ca08e2d58e4d95053be3b9dde2e14852c0cdc`.
  All regular files are byte-identical to the corresponding Git tree; there
  are no symlinks or submodules.
- Sumdb source/module identities are v0.1.0
  `h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=` /
  `h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=`; v0.2.0
  `h1:s5hAObm+yFO5uHYt5dYjxi2rXrsnmRpJx4OYvIWUaQs=` /
  `h1:ipq/a2n7PKx3OHsz4KJII5eveXtPO4qwEXGdVfWzfnI=`; v0.2.1
  `h1:Fmg33tUaq4/8ym9TJN1x7sLJnHVwhP33CNkpYV/7rwI=` with the same
  module-file identity as v0.2.0; v0.3.0
  `h1:WgNl7dwNpEZ6jJ9k1snq4pZsg7DOEN8hP9Xw0Tsjwk0=` /
  `h1:640gp4NfQd8pI5XOwp5fnNeVWj67G7CFk/SaSQn7NBk=`; and v0.3.1
  `h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=` /
  `h1:hoEshYVHaxMs3cyo3Yncou5ZscifuDolrwPKZanG3xk=`.
- Every release declares exact module path `github.com/kr/pretty`. V0.1.0 has
  no Go directive and requires kr/text v0.1.0. V0.2.0/v0.2.1 use `go 1.12`
  and kr/text v0.1.0. V0.3.0 uses `go 1.12`, kr/text v0.2.0, and go-internal
  v1.6.1. V0.3.1 uses `go 1.12`, kr/text v0.2.0, and go-internal v1.9.0.
  There are no retractions, deprecation directives, or replacements. Every
  stable and closure preserves the Go 1.18 floor; selected's closure has a
  maximum Go directive of 1.17.

### Closure, API, And Qualification Result

- Selected is one package. Production files are `diff.go`, `formatter.go`,
  `pretty.go`, and `zero.go`; its tests are `diff_test.go`,
  `formatter_test.go`, and `example_test.go`. It exports Formatter,
  Errorf/Fprintf, Log/Logf/Logln, Print/Printf/Println, Sprint/Sprintf,
  Diff/Fdiff/Pdiff/Ldiff, Printfer, and Logfer.
- Under Go 1.26.7 the selected closure has 67 production and 129 complete-test
  entries; Go 1.18.10 has 45/83. The five-module graph is target, kr/text
  v0.2.0, go-internal v1.9.0, creack/pty v1.1.9, and pkg/diff at its recorded
  pseudo-version. Only target, kr/text, go-internal/fmtsort, and test variants
  load; pty and pkg/diff are unloaded module/test-apparatus requirements.
- There is no cgo, generated source, embed, build tag, platform source,
  testdata, network, subprocess, filesystem/resource boundary, or production
  mutable global. Visited maps are per call. Callers own inputs, writers,
  callbacks, and synchronization; convenience Print/log functions use
  standard stdout/logger globals.
- Exact Go 1.26.7 and contained Go 1.18.10 pass build, upstream count-one and
  count-ten tests with vet disabled, and test cross-builds for darwin/amd64,
  linux amd64/arm64, windows amd64/386, freebsd/amd64, and js/wasm for every
  stable. Small ordinary positive formatting/diff, non-mutation, ten-repeat,
  and four-independent-call race fixtures also pass.
- All five stables fail full upstream vet in both SDKs because `diff_test.go`
  constructs `unsafe.Pointer(uintptr(1))`; v0.1.0/v0.2.0 additionally retain
  an int-to-string diagnostic. Every upstream race-count-ten run aborts at the
  same test construction with `fatal error: checkptr: pointer arithmetic
  computed bad pointer value`.
- Every stable also fails two bounded ordinary contracts under both SDKs.
  `fmt.Sprintf("%v", pretty.Formatter(nil))` produces
  `%!v(PANIC=Format method: reflect: call of reflect.Value.Interface on zero
  Value)` rather than fmt's ordinary `<nil>` pass-through representation.
  `Diff` of an ordinary four-entry map changes returned-description order
  across ten calls because `diff.go` traverses unsorted `MapKeys`.
- No stable passes all applicable documented ordinary behavior, determinism,
  upstream test-closure, race, and vet contracts. The highest qualified
  stable therefore does not exist. The evaluation used only small bounded
  ordinary values and upstream tests; it did not fuzz, stress, create
  oversized/deep/cyclic/adversarial values, reproduce a security issue, or
  analyze exploitability.

### MVS, Ownership, And Project Effects

- Exact MVS requests are Cast v1.5.1 -> v0.3.1; Consul SDK v0.8.0 -> v0.2.0;
  and Prometheus client_golang v1.4.0, sergi/go-diff v1.2.0, and errgo.v2
  v2.1.0 -> v0.1.0. None of the five requester source trees imports the
  target. Cast's requirement represents its quicktest v1.14.4 test closure,
  whose production code imports kr/pretty; the other four requirements are
  indirect/stale test-closure metadata.
- The shortest genuine current owner route is main -> selected Cast v1.5.1 ->
  selected v0.3.1. Main loads Cast through direct Viper v1.15.0, but neither
  Cast's quicktest tests nor kr/pretty load. Historical lower routes run from
  direct mvn-pom-mutator through Viper/Consul API/Consul SDK or Viper/go-
  metrics/client_golang, from direct sergi and direct Assert through sergi,
  and from three Honnef-tools routes through go-internal v1.3.0/errgo.v2.
- `go mod why -m` is positive only through project `pkg/config` -> yaml.v2 ->
  yaml.v2 tests -> check.v1 -> kr/pretty. Repository imports and target
  production/complete-test package loads are zero. The target has no runtime
  reachability. Physical selection, why, and unloaded state are not
  qualification.
- Baseline measurements reproduce exactly: 234 selected modules, 3,599 graph
  edges, 355 production entries, 429 complete-test entries, 197 module-backed
  entries across 41 loaded modules, and 1,067 sum lines. Real `go.mod` and
  `go.sum` hashes remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
  The 432-line tidy projection remains
  `3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
  applied tidy is 52/948 lines at
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`
  and retains selected v0.3.1 through Cast.
- A disposable exact v0.3.1 `go get` is not a metadata no-op because the real
  module has no target root. It manufactures indirect roots for target,
  kr/text v0.2.0, and go-internal v1.9.0, adds creack/pty and pkg/diff to the
  build list, and yields 236 modules/3,606 edges/1,072 sums while package
  loads remain 355/429/197/41 with zero target load. Tidy removes the roots and
  returns byte-exact common tidy state. Cast is a genuine supported tidy-stable
  owner, but v0.3.1 is unqualified.
- Disposable v0.3.0, v0.2.1, and v0.2.0 requests downgrade Cast, Viper or the
  accepted Cast closure and dozens of earlier selections. Their raw project
  counts are 235/3,607/1,071; 229/3,587/1,074; and 229/3,586/1,073. Lower
  loads fall to 351/425/193/41. Tidy retains broad unrelated selection drift,
  so no supported guard-preserving owner exists for those projections.
- The v0.1.0 request removes mvn-pom-mutator, collapses to 146 modules/2,208
  edges, and makes project loading fail. Tidy restores mvn-pom-mutator but
  selects v0.2.0 instead of requested v0.1.0 and retains broad guard drift.
  Every raw projection retains the main `go 1.18`/toolchain declaration, but
  none passes all earlier guards. No projection was retained.
- All 36 earlier guarded selections reproduce exactly with negative why and
  zero repository imports and production/complete-test loads. Their 228
  incoming path-version edges retain sorted SHA-256
  `361d0c355a69842518a518e682f81c9728d37acfdeff64f430a4fb253929691e`.

### Advisory, Continuity, And Final Evidence

- Fresh exact-version OSV and GitHub global queries are empty for every stable;
  the repository advisory endpoint is an empty five-byte JSON array. Pinned
  isolated govulncheck v1.8.0 module/package/symbol/test-symbol scans have no
  target finding for any stable. Base and disposable direct-v0.3.1 project ID
  sets are identical at 30 module, 22 package, 20 symbol, and 20 test-symbol
  OSVs, with no target trace.
- Exact OSV across all 36 earlier guards retains only Gorilla WebSocket
  `GO-2026-6278` / `GHSA-w67g-5rqw-f597` and go-retryablehttp
  `GO-2024-2947` / `GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0 retains
  `GO-2026-6179` and `GO-2026-6180`. The Go module vulnerability index remains
  518,501 bytes and 1,402 records at SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  Last-Modified 2026-09-17T17:29:18Z. The PUBLISHED memberlist CNA response
  remains 2,807 bytes at SHA-256
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
  updated 2026-07-08T19:40:16.119Z. Advisory absence is not qualification.
- Evaluation began from clean ordinary and ignored state on branch
  `codex/upgrade-quality` at HEAD
  `7fb81847f378ba7fbc9a758bacb507b42aea92be`, parent
  `1017a890db0c7ac42c6ebb4bc57bc8c191ac718c`, tree
  `a7e89ee68378fc8034290810ca60d6ac560a134b`. That handoff changes exactly
  the launcher, answered kr/logfmt decision, then-NEXT target archive, rolling
  handover, and roadmap. Google UUID implementation commit
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor.
- Official contained Go 1.26.7 archive/binary SHA-256 values are
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  Go 1.18.10 values are
  `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.
  Required commands used those binaries with `GOENV=off`, `GOWORK=off`,
  `GOTOOLCHAIN=local`, empty `GOFLAGS`, `LC_ALL=C`, `LANG=C`, and `umask 022`.
- Final unchanged-project exact-Go module verification, build, count-one
  tests, race count-one tests, and vet pass. No source, dependency metadata,
  Go-floor, parent, prior selection, or accepted 27/27 Q0-Q2 PASS at L2
  changed.
- The reciprocal 256-archive chain, single NEXT state, launcher/archive prompt
  mirror, exact documentation-only changed set, diff checks, and launcher
  check pass after handoff editing. Every task-owned clone, SDK, cache, tool,
  report, fixture, advisory response, and project copy was contained beneath
  the managed session scratch root and removed; only the pre-existing
  launcher-owned Node compile cache remains.

P7 continues only with the linked bounded product decision among a kr/pretty-
specific unqualified exception, one later measurement-only study of the exact
main -> Cast v1.5.1 -> kr/pretty v0.3.1 owner/request route, or stopping P7
unresolved. The successor was prepared but not executed. Cast and every
earlier decision remain closed; P8 remains queued.
