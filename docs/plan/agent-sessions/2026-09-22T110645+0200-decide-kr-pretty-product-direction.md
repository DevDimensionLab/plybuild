# Agent Session: Decide Kr Pretty Product Direction

Status: NEXT
Session ID: `2026-09-22T110645+0200-decide-kr-pretty-product-direction`
Created: `2026-09-22T11:06:45+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `73327d467eb4fc0185e0f3e6fbfdf492a60c7406405c2753adf6c89d1827bea9`
Previous: [2026-09-22T100435+0200-evaluate-kr-pretty-dependency.md](2026-09-22T100435+0200-evaluate-kr-pretty-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with one bounded product decision for exact selected, inherited,
unloaded `github.com/kr/pretty v0.3.1`. Choose one authorized direction from
the completed evaluation: explicitly retain exact selected under a new target-
specific exception, authorize exactly one later measurement-only Cast owner/
request study, or stop P7 unresolved. Do not repeat the completed behavior or
release evaluation, reopen Cast, implement a parent or workaround, combine
another dependency group, or begin P8.

# Defensive Scope

This is an ordinary dependency product decision. Use public metadata, static
records, project graph/build commands, and public advisory evidence. Do not
fuzz, stress, probe resource exhaustion, generate oversized/deep/cyclic or
adversarial values, reproduce a security issue, or perform security or
exploitability analysis. The completed source, test-closure, and ordinary
pretty-formatting behavior evaluation is final and must not be repeated.

Every disposable cache, tool, archive, report, project copy, or advisory
response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to
`/private/tmp`, `/tmp`, a sibling of the managed root, or another external
root. Verify containment and remove task-owned scratch evidence before
handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and all target-specific decisions through exact selected, inherited, unloaded
kr/logfmt. Every earlier outcome and exception is final under its own guards.
P8 remains queued.

Option 1 for kr/logfmt retains only exact
`github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515` under its own
unqualified, non-transferable exception. No kr/logfmt, kr/fs, go-windows-
terminal-sequences, gotool, errcheck, httprouter, jtolds/gls, go-junit-report,
json-iterator, clockwork, demangle, strcase, memberlist, or other exception
transfers to kr/pretty. The accepted Cast v1.5.1 decision and its quicktest
v1.14.4, kr/pretty v0.3.1, and go-internal v1.9.0 test closure remain final;
this decision may not reopen or change Cast.

# Measurements At Start

The completed evaluation began from clean HEAD
`7fb81847f378ba7fbc9a758bacb507b42aea92be`, parent
`1017a890db0c7ac42c6ebb4bc57bc8c191ac718c`, tree
`a7e89ee68378fc8034290810ca60d6ac560a134b`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`. The evaluation made no product
source or dependency-metadata change and prepared this decision-only handoff
commit. Verify the new handoff HEAD, parent, tree, exact changed set, ancestry,
and clean ordinary and ignored status rather than assuming their values.

Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The project remains 234 selected modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, 1,067 sums, and the recorded 432-line tidy projection. All 36 earlier
guarded selections and their 228 incoming edges remain exact at sorted
SHA-256
`361d0c355a69842518a518e682f81c9728d37acfdeff64f430a4fb253929691e`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

# Completed Evaluation

No exact-path kr/pretty stable qualifies. The canonical exact-path line is
v0.1.0, v0.2.0, v0.2.1, v0.3.0, and latest/selected v0.3.1. Exact `/v2` and
`/v3` proxy lines do not exist. The continuous annotated-tag ancestry peels
respectively to commits `73f6ac0b30a98e433b289500d779f50c1a6f0712`,
`4e0886370c3a67530192c6a238cff68f56c141b0`,
`ead452280cd055b2ae8a7f0db5eb37a878d902f7`,
`a883a8422cd235c67c6c4fdcb7bbb022143e10b1`, and
`3cd153a126da607b78d1762779b1e1054f9889fc`. The last three commits are
GitHub-verified signed; the first two are unsigned. Annotated tags themselves
contain no signatures. Default `main`, proxy `@latest`, and v0.3.1 are the
same commit. There are no GitHub Releases.

Exact `go-import` metadata resolves without redirect to the public, enabled,
unarchived, non-fork MIT `kr/pretty` repository. Its owner remains `kr` and
default branch remains `main`. Proxy ZIP regular files match each peeled Git
tree byte for byte; there are no symlinks or submodules. Every module names
exact `github.com/kr/pretty`. V0.1.0 has no Go directive, v0.2.0/v0.2.1 use
`go 1.12`, and v0.3.0/v0.3.1 also use `go 1.12`; every stable preserves the
Go 1.18 floor. No stable has a retraction, deprecation, or replacement. The
selected closure requires kr/text v0.2.0 and go-internal v1.9.0; its maximum
module Go directive is 1.17.

The module is one formatting/diff package. It exports Formatter and the
Errorf/Fprintf/Log/Print/Sprint families, Diff/Fdiff/Pdiff/Ldiff, Printfer, and
Logfer. Selected production has no cgo, generated source, embed, build tags,
platform files, network/subprocess/filesystem boundary, library-owned resource,
or mutable package global. Visited state is per call; callers own inputs,
writers, callbacks, and synchronization, while convenience Print/log helpers
use standard process globals.

Exact Go 1.26.7 and contained Go 1.18.10 both pass upstream build, count-one
and count-ten tests with vet disabled, supported test cross-builds, and bounded
positive ordinary formatting/non-mutation/concurrency checks for every stable.
Every stable nevertheless fails required qualification in both SDKs:

- Full upstream vet rejects the invalid `unsafe.Pointer(uintptr(1))` test
  construction; v0.1.0 and v0.2.0 have an additional int-to-string diagnostic.
- Upstream race count-ten aborts at that same test construction with checkptr
  reporting pointer arithmetic computed a bad pointer value.
- The documented Formatter pass-through contract fails for ordinary nil:
  `fmt.Sprintf("%v", pretty.Formatter(nil))` produces a Format-method panic
  diagnostic instead of fmt's `<nil>` representation.
- Diff of an ordinary four-entry map changes description order across ten
  calls because it traverses unsorted map keys, failing the required
  deterministic-output contract.

These results use only small bounded ordinary values and the upstream test
closure. No fuzzing, stress, oversized/deep/cyclic/adversarial fixture,
security reproduction, or exploitability analysis occurred. Passing ordinary
fixtures, physical selection, a why chain, zero loading, and advisory absence
do not override the failures.

Exact MVS has five target requests: selected Cast v1.5.1 requests v0.3.1;
Consul SDK v0.8.0 requests v0.2.0; and Prometheus client_golang v1.4.0,
sergi/go-diff v1.2.0, and errgo.v2 v2.1.0 request v0.1.0. Cast is the current
genuine tidy-stable owner: main has an indirect requirement on Cast v1.5.1,
the project loads Cast through direct Viper v1.15.0, and Cast's quicktest
v1.14.4 test closure imports kr/pretty. Cast itself and the four lower
requesters contain no Go import of kr/pretty; their direct module edges are
indirect, stale, or test-closure requirements.

The lower genuine historical routes are main -> direct mvn-pom-mutator
v0.2.3 -> historical Viper v1.10.1 -> Consul API v1.12.0 -> Consul SDK
v0.8.0; Viper -> go-metrics v0.3.10 -> client_golang v1.4.0, with parallel
routes through sagikazarmark/crypt v0.4.0; main/direct Assert v1.0.0 ->
sergi/go-diff v1.2.0 as well as main's direct sergi requirement; and three
historical Honnef-tools routes through go-internal v1.3.0 -> errgo.v2 v2.1.0.
The positive target why chain exists only through project config -> yaml.v2 ->
yaml.v2 tests -> check.v1 -> kr/pretty. Repository imports and target
production/complete-test loads are zero, so there is no project runtime
reachability.

Disposable exact `go get` projections were retained only for measurement.
V0.3.1 manufactured direct-indirect target/text/go-internal roots, added two
modules, seven graph edges, and five sums, but changed no package load. Tidy
removed the manufactured roots and converged to the common 52-line/948-line
state while Cast remained the genuine supported owner. V0.3.0, v0.2.1, and
v0.2.0 downgraded Cast/Viper or their accepted closures and many earlier
guards. V0.1.0 removed mvn-pom-mutator and made project loading fail; tidy
restored the owner but selected v0.2.0 rather than requested v0.1.0. No lower
candidate has an acceptable supported tidy-stable owner projection. No
projection, direct root, or dependency commit was retained.

The unchanged project remains 234/3,599/355/429/197/41/1,067. The 432-line
tidy projection retains SHA-256
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
its applied 52/948-line module hashes remain
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
All 36 earlier guarded selections retain exact versions, negative why, zero
repository imports and production/complete-test loads, and the 228-edge hash
above.

Fresh exact-version OSV and GitHub global results are empty for all five
stables, and the repository advisory endpoint is empty. Pinned isolated
govulncheck v1.8.0 module/package/symbol/test-symbol scans have no finding for
any stable. Base and disposable direct-v0.3.1 project populations are
identical at 30 module, 22 package, 20 symbol, and 20 test-symbol OSVs with no
target trace. Guard OSV retains only Gorilla WebSocket `GO-2026-6278` /
`GHSA-w67g-5rqw-f597` and go-retryablehttp `GO-2024-2947` /
`GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0 retains `GO-2026-6179` and
`GO-2026-6180`. The 518,501-byte/1,402-record Go module index and the
PUBLISHED 2,807-byte memberlist CNA response remain byte-exact. Advisory
absence does not qualify the ordinary behavior or failed upstream gate.

No source, `go.mod`, `go.sum`, parent, Go-floor, earlier guard, or accepted
quality result changed. No changed-selection scorecard applies. Final
unchanged-project exact-Go module verification, build, count-one tests, race
count-one tests, and vet pass under canonical `umask 022`.

# Required Product Decision

Choose exactly one:

1. Explicitly retain exact selected, inherited, unloaded v0.3.1 without source
   or dependency-metadata changes under a kr/pretty-specific, non-transferable
   exception. Call it unqualified. Accept only the completed nil Formatter,
   map Diff determinism, upstream vet, and upstream race/checkptr failures plus
   the completed related evaluation. Define exact expiry guards for path/
   version, all five requests and genuine routes, requester/import facts,
   Cast/quicktest closure, why/import/load/runtime, graph/module/tidy/Go-floor,
   every earlier guard, advisories/findings/releases/owners, and any compatible
   genuine route to a qualified exact-path release.
2. Authorize exactly one later bounded measurement-only owner/request study.
   Name the existing shortest genuine owner route: main -> exact selected Cast
   v1.5.1 -> exact selected kr/pretty v0.3.1. State the exact question: whether
   a later genuine supported tidy-stable Cast selection removes kr/pretty or
   requests a qualified exact-path stable while preserving the accepted Cast
   decision, Go floor, and every contract. Do not run it, reopen or change
   Cast, implement a parent change, add a direct target root, transfer an
   exception, or imply approval of graph/source/metadata changes.
3. Stop P7 unresolved, record the blocker, and prepare no dependency or P8
   implementation.

The decision must not call physical MVS selection, a why chain, zero loading,
or advisory absence qualification. It must not select another stable without
a qualified candidate and genuine supported owner, promote a branch, pseudo-
version, fork, replacement, or alternate path, add a direct target edge,
change Cast, Viper, quicktest, a lower historical requester, the Go floor,
product source, dependency metadata, or another selection, reopen an earlier
decision, or begin P8.

# Role And Boundaries

This is a decision-recording session, not a renewed target audit or an
implementation. Reuse the completed evaluation and make exactly one choice.
Do not broaden it into direct use, ownerless selection, alternate-path
promotion, parent removal, patching, forking, wrapping, replacement, a Go-
floor change, unrelated-module authorization, another dependency group, or
P8.

# Required Reading

Read this decision archive, the answered target evaluation, the accepted Cast
v1.5.1 evaluation and related requester/retained-module records, kr/logfmt,
kr/fs, go-windows-terminal-sequences, gotool, and errcheck decisions/
evaluations, rolling handover, roadmap, `go.mod`, and `go.sum`. Verify the new
handoff HEAD/parent/tree and exact changed set, reciprocal archive chain,
latest Google UUID implementation ancestry, exact Go 1.26.7 identity, launcher
check, module hashes/counts/tidy projection, all five requests and genuine
routes, target/requester why-import-load state, all 36 guarded selections and
228-edge snapshot, and fresh target/guard advisory identities. Treat the
completed behavior, test-closure, release, and disposable-project measurements
as final.

# Three Moves

First, verify decision prerequisites and choose exactly one authorized option.
Second, record only that decision; leave source, dependency metadata, parents,
Go floor, Cast, and every earlier guard unchanged. Third, update roadmap and
rolling handover, answer this archive, prepare the one reciprocal successor
the choice requires, verify containment, and commit the documentation handoff
without executing the successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert,
bypass cleanup, remove the worktree, transfer an exception, repeat the
evaluation, run the Cast owner study, reopen Cast or an earlier decision,
evaluate another dependency group, write outside the managed scratch root, or
begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
