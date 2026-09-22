# Agent Session: Decide Kr Text Product Direction

Status: ANSWERED - HISTORY
Session ID: `2026-09-22T141112+0200-decide-kr-text-product-direction`
Created: `2026-09-22T14:11:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `eccc3b598e1d611534489cc48b552989be70b267dd6dbf345daf29c2d43fc5aa`
Previous: [2026-09-22T131156+0200-evaluate-kr-text-dependency.md](2026-09-22T131156+0200-evaluate-kr-text-dependency.md)
Next: [2026-09-22T144301+0200-evaluate-kyokomi-emoji-v2-dependency.md](2026-09-22T144301+0200-evaluate-kyokomi-emoji-v2-dependency.md)
Outcome: Option 1 selected; exact inherited, unloaded kr/text v0.2.0 is retained unqualified under a target-specific exception.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by recording exactly one bounded product decision for selected
exact-path `github.com/kr/text v0.2.0`. The completed evaluation found that no
exact-path stable qualifies: v0.1.0 and v0.2.0 both fail an ordinary documented
wrapping contract, and the upstream fix exists only after the latest stable.
Choose only one of the three directions below, apply that choice exactly, and
prepare only its reciprocal successor if the choice requires one. Do not repeat
the dependency evaluation, evaluate another dependency group, or begin P8.

# Defensive Decision Scope

This is an ordinary dependency-quality product decision. Use only the completed
public release/repository metadata, static source and graph facts, small bounded
ordinary text-behavior results, project projections, and advisory identities
recorded here and in the answered evaluation. Do not fuzz, stress, probe
resource exhaustion, create oversized, deeply nested, cyclic, malformed, or
adversarial inputs, reproduce a security issue, or perform security or
exploitability analysis.

Every disposable cache, tool, report, project copy, fixture, or advisory
response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to
`/private/tmp`, `/tmp`, a sibling of the managed root, or another external root.
Verify containment and remove task-owned scratch evidence before handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2, and
target-specific decisions through exact selected, inherited, unloaded kr/pty
v1.1.1. P8 remains queued.

Kr/pty option 1 is exact, unqualified, target-specific, non-transferable, and
final. The historical kr/text v0.1.0 -> kr/pty v1.1.1 request, requester import,
complete route families, why/import/load/runtime facts, graph/module/tidy/Go-
floor state, earlier guards, release/repository/advisory identities,
qualification, and compatible-route condition remain expiry guards. Any change
requires a fresh kr/pty dependency and product decision before merge. Do not
transfer the kr/pty exception to kr/text or reopen kr/pty, kr/pretty, Cast, or an
earlier decision.

# Completed Evaluation Is Final

The evaluation began from clean branch `codex/upgrade-quality` at HEAD
`0781a229dae9852e496c80a10b81075384e6fab8`, parent
`91c6f272a9a71e57b72a5519d9a79f657ded1951`, tree
`a7b7c64cece1d80a1b71b51185d23ee4b82b1b52`. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. The
evaluation changed no product source or dependency metadata and prepared this
decision-only handoff. Verify the new handoff HEAD, parent, tree, exact changed
set, ancestry, and clean ordinary and ignored status rather than assuming them.

The exact go-import record maps `github.com/kr/text` to
`https://github.com/kr/text.git`. Public repository ID 4450535 is enabled,
unarchived, not a fork, MIT licensed, owned by `kr`, and defaults to `main`.
There are no GitHub Releases. The exact-path proxy line contains only v0.1.0
and v0.2.0; `/v2` and `/v3` lines do not exist. Neither module is retracted,
deprecated, or replaced, and neither declares a Go version.

Unsigned annotated tag v0.1.0 peels to unsigned commit
`e2ffdb16a802fe2bb95e2e35ff34f0e53aeef34f`, tree
`82ddc7d7e0af2da48325deab3c7c7a0e16167151`. Unsigned annotated tag v0.2.0
peels to GitHub-verified commit
`702c74938df48b97370179f33ce2107bd7ff3b3e`, its direct descendant, tree
`87941db3a7fd2fb33a98b3bf508eb2f352c53061`. Proxy archive, Git regular-file,
and sumdb identities agree. No fork, branch, pseudo-version, alternate module
path, replacement, or ownerless candidate was promoted.

Both stables expose the same root text API, colwriter package, and `agg` and
`mc` commands. They have no cgo, generated source, build tags, embed, network,
or shared library global state. Library operations are deterministic; caller
owns inputs, output writers, mutable writer instances, synchronization, and
colwriter Flush lifecycle. Commands own their process-local state and ordinary
stdin/stdout/file boundaries. The `mc` command alone owns a PTY boundary.

Exact Go 1.26.7 and contained Go 1.18.10 upstream build, count-one/count-ten
tests, race-count-ten tests, vet, and supported Darwin/Linux/FreeBSD production
cross-builds pass for both stables. Their complete closures preserve Go 1.18.
Windows and js/wasm fail only at the PTY-backed `mc` command boundary. Small
ordinary positive fixtures for indentation, input nonmutation, independent
wrapping, determinism, and concurrency pass under both SDKs.

Both stables nevertheless fail the same documented ordinary contract. The API
states that a line can exceed the limit only when a single word exceeds it,
but `Wrap("overlong overlong foo", 4)` returns
`"overlong overlong\nfoo"`, joining two over-limit words. The failure repeats
under both SDKs and race. Upstream verified commit
`838204404ccb967534580e2a6efb24062880dd0f`, tree
`fd73252dc53a2d59bd6b7c128858d17e63172c2a`, merged PR #9 to fix adjacent long
words and add its fixture, but no later stable exists. Branch/main and that
commit are not eligible release candidates. Therefore no exact-path stable
qualifies.

Exact MVS selects v0.2.0. Selected Cast v1.5.1 genuinely requests v0.2.0 from
its supported test closure. Historical kr/pretty v0.1.0 and v0.2.0 each request
v0.1.0. V0.1.0 alone requests kr/pty v1.1.1, imported only by its unloaded `mc`
command; v0.2.0 instead requests creack/pty v1.1.9. The complete current and
historical owner/request routes are final as recorded in the evaluation.

Target why is positive only through project `pkg/config` -> yaml.v2 -> yaml.v2
tests -> check.v1 -> kr/pretty -> kr/text. Repository imports and production/
complete-test loads for kr/text, kr/pty, and creack/pty are zero. Kr/pty and
creack/pty why results are negative. The target and both PTY boundaries are not
runtime reachable.

A disposable exact v0.2.0 get manufactures direct main -> kr/text v0.2.0 and
kr/text v0.2.0 -> creack/pty v1.1.9 edges, producing 235 modules, 3,601 edges,
1,069 sums, and unchanged 355/429/197/41 load facts. Tidy removes both roots and
restores exact base selection, the historical kr/pty request, and common tidy
state. The raw graph/module change would expire a kr/pty guard, so none was
retained.

A disposable v0.1.0 get downgrades kr/text, Cast, Viper, and kr/pretty,
producing 229 modules, 3,585 edges, 389 production entries, 425 complete-test
entries, and 1,073 sums. Tidy retains those selection and graph changes. It
violates the closed Cast, kr/pretty, and kr/pty guards. No projection was
retained. No exact `go get` was run in the real project.

The unchanged project remains 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod`/`go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy projection remains
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes remain
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

All 38 earlier guarded selections remain exact. Their 234 sorted incoming
edges retain SHA-256
`d31b2afdb5a11dadac236eec600408cc9ba0ee3deca2793ae2adf7ffc449c320`;
37 why results are negative and only the closed kr/pretty why is positive,
with zero guarded repository imports and production/complete-test loads.

Fresh exact-version OSV, GitHub global, and repository advisory results for
both stables are empty. Pinned govulncheck v1.8.0 isolated module/package/
symbol/test-symbol scans are empty for both; base and disposable direct-v0.2.0
project populations are identical at 30/22/20/20 with no target trace. Guard
OSV retains only the Gorilla WebSocket and go-retryablehttp pairs; x/mod v0.14
retains GO-2026-6179 and GO-2026-6180. The 518,501-byte, 1,402-record Go index
and PUBLISHED 2,807-byte memberlist CNA response retain their recorded hashes.
Advisory absence does not override the behavior failure.

Final unchanged-project exact Go 1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass under `umask 022`. Accepted quality
remains 27/27 Q0-Q2 PASS at L2. Treat the completed release, source, behavior,
closure, route, projection, advisory, and final-gate findings as final.

# Choose Exactly One Direction

1. Explicitly retain exact selected, inherited, unloaded
   `github.com/kr/text v0.2.0` without product-source or dependency-metadata
   changes under a kr/text-specific, non-transferable exception. Call it
   unqualified: neither exact-path stable qualifies, and the behavior fix is
   unreleased. Accept only the completed adjacent-overlong-word failure,
   release/closure/platform facts, exact owner/request routes, why/import/load/
   runtime boundary, graph/tidy/Go-floor state, earlier guards, and advisory
   identities. Define the exact expiry guards below. This exception must not
   broaden or expire the closed kr/pty exception.
2. Authorize exactly one later bounded measurement-only genuine owner/request
   study. Name the existing shortest supported route: main -> exact selected
   Cast v1.5.1 -> exact selected kr/text v0.2.0. State the exact question:
   whether a later genuine supported tidy-stable Cast selection removes
   kr/text or requests a future qualified exact-path kr/text stable while
   preserving Go 1.18, the closed Cast/kr/pretty/kr/pty decisions, the exact
   historical v0.1.0 -> kr/pty v1.1.1 boundary, and every project contract. Do
   not run the study or change any selection in this decision-recording move;
   prepare one reciprocal measurement-only successor.
3. Stop P7 unresolved with no product-source, dependency-metadata, roadmap, or
   P8 implementation change and no successor execution.

Do not invent a fourth option or combine options. Physical selection, a why
chain, zero loading, and advisory absence are not qualification. Do not call
v0.2.0 qualified, promote the unreleased fix, add a target root, select v0.1.0,
change Cast, Viper, kr/pretty, kr/pty, a historical requester, the Go floor,
product source, dependency metadata, or another module, transfer an exception,
or begin P8.

# Option 1 Expiry Guards

If option 1 is selected, preserve and record at minimum:

- exact selected path/version `github.com/kr/text v0.2.0`, selected Cast
  v1.5.1 -> v0.2.0, historical kr/pretty v0.1.0/v0.2.0 -> v0.1.0, and
  historical v0.1.0 -> kr/pty v1.1.1;
- complete current and historical genuine owner/request routes, requester
  imports, Cast test-closure ownership, and no direct target root;
- the positive yaml.v2-tests/check.v1/kr/pretty why chain only, zero target and
  PTY repository imports and production/complete-test loads, negative kr/pty
  and creack/pty why results, and runtime unreachability;
- the exact two-release line, missing `/v2` and `/v3`, repository/owner/status/
  license/default-branch identity, tag/commit/tree/archive/sumdb facts, no
  replacement/retraction/deprecation, and no later exact-path stable;
- both releases remaining unqualified for the completed adjacent-overlong-word
  behavior failure, the post-release fix remaining unreleased, the completed
  API/closure/platform/ordinary fixture results, and no qualified supported
  route;
- baseline 234/3,599/355/429/197/41/1,067 state, exact module hashes and tidy
  state, both disposable-projection results, and every kr/pty expiry guard;
- exact Go 1.18 floor, exact Go 1.26.7 identity, source/API/CLI/help/launcher/
  Make/quality contracts, and accepted 27/27 Q0-Q2 PASS at L2;
- all 38 earlier guarded selections and their 234 incoming-edge snapshot; and
- no new target/requester advisory, independent defect, exact-path stable,
  repository/release/owner change, qualified stable, supported tidy-stable
  owner, or compatible genuine route to a qualified kr/text or kr/pty release.

Any path/version, request, requester import, owner route, root, why/import/load/
runtime fact, graph, module hash, tidy state, Go floor, earlier or kr/pty guard,
advisory, independent finding, repository/release/owner, qualification,
supported owner, or compatible-route change expires the exception and requires
a fresh kr/text dependency and product decision before merge. Any kr/pty guard
change separately requires a fresh kr/pty dependency and product decision.
Option 1 authorizes no owner study, direct root, unreleased commit, alternate
path, dependency edit, workaround, or implementation.

# Role And Boundaries

This session is decision recording, not dependency implementation or a new
evaluation. Revalidate only the minimum continuity, exact graph/guard,
advisory identity, launcher, and final unchanged-project checks necessary to
record the chosen option safely. Do not rerun text behavior fixtures, PTY
behavior, candidate projections, or an owner study. Preserve a clean worktree
except for the bounded documentation/launcher handoff, then make the required
local handoff commit.

# Required Reading

Read this decision archive, the answered kr/text evaluation, answered kr/pty
decision/evaluation, answered kr/pretty decision/evaluation, kr/logfmt, kr/fs,
go-windows-terminal-sequences, gotool, and errcheck decisions/evaluations,
accepted Cast v1.5.1 and requester records, rolling handover, P7/P8 roadmap,
`go.mod`, and `go.sum`. Verify the new handoff HEAD/parent/tree and exact
changed set, reciprocal archive chain, latest Google UUID implementation
ancestry, exact Go identity, launcher check, module hashes/counts/tidy
projection, target requests/routes/why/import/load state, all 38 earlier
guards and the 234-edge snapshot, and fresh advisory identities.

# Three Moves

First, ask for or apply exactly one explicit option without reopening the
completed evaluation. Second, record only the chosen bounded direction and run
its minimum unchanged-project guards; do not implement an owner study or
dependency change. Third, update roadmap and rolling handover, answer this
archive, prepare at most the one reciprocal successor authorized by the
choice, verify containment, and make the required local handoff commit without
executing a successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert, bypass
cleanup, remove the worktree, select a dependency, add a target root, transfer
an exception, reopen kr/pty, kr/pretty, Cast, or an earlier decision, evaluate
another dependency group, write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

Option 1 was selected on 2026-09-22. Exact selected
`github.com/kr/text v0.2.0` is explicitly retained as inherited and unloaded
without changing product source, `go.mod`, or `go.sum`. It remains unqualified
and is not described as secure: neither exact-path stable qualifies, and the
matching upstream behavior fix remains unreleased. Physical MVS selection, a
test-only why chain, zero loading, and advisory absence bound current exposure;
none qualifies or implicitly accepts the module.

The kr/text-specific, non-transferable exception accepts only the completed
adjacent-overlong-word failure and the completed repository, release, archive,
module, Go-floor, API, closure, platform, positive ordinary fixture, owner/
request, why/import/load/runtime, graph/tidy, projection, earlier-guard, and
advisory findings. It accepts no uncharacterized behavior, new advisory, or
independent defect. It does not broaden, replace, transfer, or expire the
closed kr/pty exception, and it transfers no other exception.

The exception remains valid only while every one of these facts remains exact:

- selected exact path/version `github.com/kr/text v0.2.0`; selected Cast
  v1.5.1 -> v0.2.0; historical kr/pretty v0.1.0/v0.2.0 -> v0.1.0; and
  historical kr/text v0.1.0 -> kr/pty v1.1.1;
- the current shortest supported owner/request route, main -> selected Cast
  v1.5.1 -> selected kr/text v0.2.0, including main's indirect Cast
  requirement, the loaded direct Viper v1.15.0 consumer route, and Cast's
  supported Go-1.18 quicktest v1.14.4 / kr/pretty v0.3.1 /
  rogpeppe/go-internal v1.9.0 test closure; Cast retaining no target source
  import and no direct target root being present;
- every historical route family: direct sergi/go-diff v1.2.0 and direct
  Assert v1.0.0 through sergi reaching kr/pretty v0.1.0; direct
  mvn-pom-mutator v0.2.3 through historical Viper v1.10.1, Consul API
  v1.12.0, and Consul SDK v0.8.0 reaching kr/pretty v0.2.0; Viper or
  sagikazarmark/crypt v0.4.0 through go-metrics v0.3.10 and client_golang
  v1.4.0 reaching kr/pretty v0.1.0; and all three historical Honnef-tools
  vertices through go-internal v1.3.0 and errgo.v2 v2.1.0 reaching kr/pretty
  v0.1.0, with every applicable lower route continuing through kr/text
  v0.1.0 to kr/pty v1.1.1;
- the positive target why chain only through project `pkg/config` -> yaml.v2
  -> yaml.v2 tests -> check.v1 -> kr/pretty -> kr/text; negative kr/pty and
  creack/pty why results; zero target and PTY repository imports and
  production/complete-test loads; and no target or PTY runtime reachability;
- exactly two exact-path releases, v0.1.0 and v0.2.0, with absent `/v2` and
  `/v3` lines; public enabled, unarchived, non-fork MIT repository ID 4450535,
  owned by `kr`, defaulting to `main`, with no GitHub Releases; no retraction,
  deprecation, replacement, later exact-path stable, or eligible alternate
  path, fork, branch, or pseudo-version;
- unsigned annotated v0.1.0 tag object
  `a90d266dd68b224558779a7ed518f29481427f41` peeling to unsigned commit
  `e2ffdb16a802fe2bb95e2e35ff34f0e53aeef34f` and tree
  `82ddc7d7e0af2da48325deab3c7c7a0e16167151`; unsigned annotated v0.2.0 tag
  object `0e5f52c28dd72ab84daeb81b5a51f20fdc35f9c5` peeling to its direct
  descendant, GitHub-verified commit
  `702c74938df48b97370179f33ce2107bd7ff3b3e` and tree
  `87941db3a7fd2fb33a98b3bf508eb2f352c53061`;
- proxy ZIP SHA-256 values
  `9363a4c8f1f3387a36014de51b477b831a13981fc59a5665f9d21609bea9e77c` /
  `368eb318f91a5b67be905c47032ab5c31a1d49a97848b1011a0d0a2122b30ba4`,
  normalized regular-file archive/Git manifest hashes
  `5960b8daf161168d9c1b4c933e5387b50b2fda028a7d45a8280a21e7bd37c733` /
  `7e59725cfbda17e0daeec72fd1cb375a0e118697a28a250ca17c0e76f1ec2616`,
  and the recorded v0.1.0 and v0.2.0 sumdb source/module identities;
- both releases remaining unqualified because
  `Wrap("overlong overlong foo", 4)` joins the two adjacent over-limit words;
  verified fix commit `838204404ccb967534580e2a6efb24062880dd0f`, tree
  `fd73252dc53a2d59bd6b7c128858d17e63172c2a`, remaining post-release and
  unreleased; the completed API/closure/platform results and positive bounded
  fixtures remaining exact; and no qualified supported route appearing;
- the baseline 234 selected modules, 3,599 graph edges, 355 production
  entries, 429 complete-test entries, 197 module-backed entries across 41
  loaded modules, and 1,067 sum lines;
- `go.mod` / `go.sum` SHA-256 values
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  the 432-line tidy projection at
  `3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`,
  and the common applied 52/948-line hashes
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`;
- both completed disposable projections: exact v0.2.0 manufacturing target and
  creack/pty roots for 235 modules/3,601 edges/1,069 sums before tidy restores
  the base, while exact v0.1.0 downgrades Cast, Viper, and kr/pretty for
  229 modules/3,585 edges/389 production/425 complete-test/1,073 sums and
  retains closed-guard drift after tidy; neither projection being retained;
- every kr/pty expiry guard, including exact selected v1.1.1, its historical
  v0.1.0 request and unloaded `mc` import, complete route families, negative
  why, zero import/load/runtime state, graph/module/tidy/Go-floor state,
  release/repository/advisory identities, v1.1.4 qualification and missing
  supported owner, later failures, and compatible-route condition;
- declared Go 1.18 floor; exact Go 1.26.7 darwin/arm64 official archive/binary
  SHA-256 identities
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
  unchanged source/API/CLI/help/launcher/Make/quality contracts; and accepted
  27/27 Q0-Q2 PASS at L2;
- all 38 earlier guarded selections exact, 37 negative why results and only
  closed kr/pretty positive, zero guarded repository imports and production/
  complete-test loads, and their 234 sorted incoming edges at SHA-256
  `d31b2afdb5a11dadac236eec600408cc9ba0ee3deca2793ae2adf7ffc449c320`;
  including kr/text yields 39 selections and 237 incoming edges at SHA-256
  `151e72c0b5444acc59ffab45d6ce8fe43e80821b42df6a49a11902b3a6632803`;
  and
- no new target/requester advisory, independent defect, exact-path stable,
  repository/release/owner change, qualified stable, supported tidy-stable
  owner, or compatible genuine route to a qualified kr/text or kr/pty release.

Any path/version, request, requester import, owner identity or route, root,
why/import/load/runtime fact, graph, module hash, tidy state, Go floor, earlier
or kr/pty guard, advisory, independent finding, repository/release/owner,
qualification, supported owner, or compatible-route change expires this
exception and requires a fresh kr/text dependency and product decision before
merge. Any kr/pty guard change separately requires a fresh kr/pty dependency
and product decision. This decision authorizes no owner study, direct root,
unreleased commit, alternate path, dependency edit, workaround, unrelated
selection, or implementation.

Guard-only revalidation began from clean ordinary and ignored state on branch
`codex/upgrade-quality` at handoff HEAD
`448692f193ecdd9d7f0e608472251ce780ff62cb`, parent
`0781a229dae9852e496c80a10b81075384e6fab8`, tree
`d449b3b805c4db0e58d422283ce3a018950f87d3`. That handoff changes exactly
`codex-dev-start.sh`, the answered kr/text evaluation archive, this then-NEXT
decision archive, rolling handover, and roadmap. Its reciprocal 260-archive
chain, sole NEXT state, launcher/archive prompt mirror, exact changed set,
clean status, and launcher check passed. Exact Google UUID v1.4.0 dependency
commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor.

A freshly downloaded official Go 1.26.7 archive and its exact binary reproduce
the identities above. Under `GOENV=off`, `GOWORK=off`,
`GOTOOLCHAIN=local`, empty `GOFLAGS`, `LC_ALL=C`, `LANG=C`, scratch-contained
caches, and canonical `umask 022`, the baseline counts, module hashes, tidy
projection and applied tidy hashes, target requests, why/import/load state, all
38 earlier selections, and both edge snapshots reproduce exactly.

Fresh proxy metadata still exposes only v0.1.0 and v0.2.0, resolves latest to
v0.2.0, and has no `/v2` or `/v3` line. GitHub still reports repository ID
4450535 as enabled, unarchived, non-fork, MIT, owned by `kr`, defaulting to
`main`, with no Releases. Exact-version OSV and GitHub global results for both
stables and the repository advisory result remain empty. Guard OSV retains only
Gorilla WebSocket `GO-2026-6278` / `GHSA-w67g-5rqw-f597` and
go-retryablehttp `GO-2024-2947` / `GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0
retains `GO-2026-6179` and `GO-2026-6180`. The Go module vulnerability index
remains 518,501 bytes/1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. The PUBLISHED 2,807-byte memberlist CNA
response remains SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`,
updated 2026-07-08T19:40:16.119Z.

The completed text and PTY behavior, release/source/closure, and disposable-
project evaluations were not repeated; option 2 was neither authorized nor
run. No source, dependency metadata, parent, toolchain declaration, or earlier
guard changed. No changed-selection scorecard applies; accepted quality
remains 27/27 Q0-Q2 PASS at L2. Final unchanged-project exact-Go module
verification, build, count-one tests, race count-one tests, and vet pass.

P7 continues only with the linked bounded evaluation of the next selected
queue item, `github.com/kyokomi/emoji/v2 v2.2.12`. That reciprocal successor
was prepared but not executed and may not reopen this kr/text decision, the
kr/pty decision, or any earlier decision. P8 remains queued.
