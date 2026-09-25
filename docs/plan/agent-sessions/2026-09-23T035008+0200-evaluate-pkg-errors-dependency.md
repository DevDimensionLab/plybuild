# Agent Session: Evaluate Pkg Errors Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T035008+0200-evaluate-pkg-errors-dependency`
Created: `2026-09-23T03:50:08+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `064b48d9109113325396cad3eb9f5632279e5425b3bc3d4e5dc77b487927e8c1`
Previous: [2026-09-23T031134+0200-decide-pascaldekloe-goe-product-direction.md](2026-09-23T031134+0200-decide-pascaldekloe-goe-product-direction.md)
Next: [2026-09-23T043812+0200-decide-pkg-errors-product-direction.md](2026-09-23T043812+0200-decide-pkg-errors-product-direction.md)
Outcome: No supported upgrade, replacement, removal, or exact-retention direction qualifies; selected v0.9.1 remains unchanged pending the prepared product decision.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with one bounded evaluation of exact selected
`github.com/pkg/errors v0.9.1`. Determine whether a genuine supported upgrade,
replacement, removal, or exact retention direction qualifies while preserving
the Go 1.18 compatibility floor and every earlier target-specific guard. Do
not make the product decision, implement a dependency change, evaluate another
dependency group, or begin P8.

# Defensive Scope

Use public metadata, static repository/source inspection, ordinary documented
behavior, and normal project/upstream graph, build, test, race, vet, and
vulnerability commands. Do not fuzz, stress, probe resource exhaustion,
construct oversized, deeply nested, cyclic, malformed, adversarial, or escape-
sequence payloads, reproduce a security issue, or perform security or
exploitability analysis.

Every disposable cache, tool, response, report, module copy, and project copy
must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to
`/private/tmp`, `/tmp`, a sibling of the managed root, or another external
root. Verify containment and remove task-owned scratch evidence before
handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
the nine final target-specific option-1 decisions through Goe, and qualified
no-selection-change modern-go/concurrent and modern-go/reflect2. P8 remains
queued.

Exact Goe v0.1.0, ULID, go-conntrack, mapstructure, go-homedir, Promptui,
emoji/v2, kr/text, kr/pty, kr/pretty, kr/logfmt, kr/fs, go-windows-terminal-
sequences, gotool, errcheck, httprouter, GLS, go-junit-report, json-iterator,
and clockwork exceptions are final, target-specific, unqualified, and non-
transferable. Every concurrent, reflect2, Cast, Viper, memberlist, earlier
selection, owner, request, route, why/import/load/runtime, graph/module/tidy/
Go-floor, advisory, source, behavior, closure, qualification, and expiry guard
remains final. Do not reopen, broaden, or transfer any decision.

The Goe decision explicitly retains exact selected, inherited, unloaded
v0.1.0 without product-source or dependency-metadata changes. Neither exact-
path stable qualifies: v0.1.0 fails the Go 1.26.7 example and complete race
gates; v0.1.1 still fails race and has no genuine supported tidy-stable
project owner. All four target requests, requester test-import/metadata-only
boundaries, every recorded route, no root, negative why/import/load/runtime
facts, repository/release/archive/source/behavior/closure identities, both
projections, advisories, and earlier guards are retention conditions. The
Armon Go-Metrics Goe Ownership Study is not authorized.

# Fixed Project Guards

The unchanged project has 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The 432-line tidy diff SHA-256 is
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes are
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.

All 47 pre-Goe guarded selections and 276 sorted incoming edges have SHA-256
`7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
Forty-one why results are negative; only kr/pretty, kr/text, emoji/v2,
Promptui, go-homedir, and mapstructure are positive. Promptui and go-homedir
are the only guarded repository imports. Emoji/v2, Promptui, go-homedir, and
mapstructure are the only production/complete-test loaded guarded modules.
Add the complete Goe selection/request/route/import/load snapshot as its own
guard; do not fold it into or replace the recorded pre-Goe identity.

Guard OSV remains limited to the recorded Gorilla WebSocket and go-
retryablehttp pairs; x/mod v0.14.0 retains GO-2026-6179 and GO-2026-6180. The
Go vulnerability index is 518,501 bytes/1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response is 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Advisory absence is not qualification.

# Evaluation Contract

Begin with guard-only reproduction of the new handoff and fixed state. Stop
for a fresh owning evaluation if any earlier guard changed. Then answer, with
reproducible evidence, only the questions needed for `github.com/pkg/errors`:

1. Resolve the exact physical selection, direct/indirect/root state, every
   current and relevant historical request, MVS outcome, requester selection,
   complete owner routes, requester imports or metadata-only boundaries,
   `go mod why`, repository imports, production/complete-test/module-backed
   loads, and runtime relevance.
2. Establish canonical module/repository ownership and status, license,
   default branch, exact stable/prerelease/major lines, replacements,
   retractions, deprecations, redirects, eligible alternates, tags, releases,
   commits, trees, ancestry, signatures, module files, proxy/sumdb identities,
   and archive-to-Git source identity. Do not treat a branch, fork, replacement,
   alternate path, or pseudo-version as a stable release.
3. Compare selected v0.9.1 with every genuine supported Go-1.18-compatible
   candidate. Inspect exported API, build constraints, cgo/generated/embed
   boundaries, dependency closure, and ordinary documented behavior that can
   materially affect this project. Use only small deterministic fixtures.
4. Run complete upstream verification, build, count-one and repeated ordinary
   tests where meaningful, race, and vet under exact Go 1.18.10 and Go 1.26.7,
   plus supported cgo-disabled cross-build/test-compile checks justified by
   the package boundary. Record every failure; a partial package pass does not
   qualify a release whose complete required gate fails.
5. In disposable project copies, measure selected/no-op, candidate upgrade,
   removal, and tidy projections that are supported by a real owner route.
   Record exact module/edge/load/module-file/sum changes, owner stability,
   Go-floor effects, earlier-guard preservation, and restoration behavior.
   Retain none of these projections.
6. Refresh exact target OSV, GitHub global and repository advisories, isolated
   pinned govulncheck, project govulncheck populations/traces, guard OSV/x/mod,
   Go-index, and memberlist-CNA identities. Advisory absence cannot qualify a
   candidate or erase an ordinary defect.

A qualifying direction must preserve the Go 1.18 floor, pass its complete
source/API/closure/behavior/build/test/race/vet contract, have a genuine
supported tidy-stable project owner where selection matters, preserve every
earlier guard, and introduce no unreviewed product or CLI behavior. If no
direction qualifies, say so plainly and prepare a bounded product decision;
do not silently retain, upgrade, remove, replace, patch, vendor, or add a root.

# Measurements At Start

The Goe decision began from clean branch `codex/upgrade-quality` at handoff
HEAD `413fb1a154a5c9a4acfe151dfb9518714a701ed1`, parent
`536df39770daedc4f841a6898e5c69b8605eed6f`, tree
`0f6a53bb3eaa70c43d1bc91047ca25fa399f6c0f`. Google UUID implementation
commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. The Goe
decision changed no product source or dependency metadata. Independently
verify the new handoff HEAD, parent, tree, exact changed set, ancestry, clean
ordinary and ignored status, reciprocal archive chain, launcher check, exact
Go identities, fixed project and Goe guards, advisories, and final gates.

# Role And Boundaries

Act only as the bounded pkg/errors evaluator. Measure and report the target's
real choices without selecting one, editing source or dependency metadata,
changing an owner/requester, or weakening a prior guard. Stop if the evidence
requires a broader product, architecture, security, or cross-dependency choice.

# Required Reading

Read this archive, the answered Goe decision/evaluation, answered ULID and go-
conntrack decisions/evaluations, reflect2 and concurrent evaluations,
mapstructure, go-homedir, Promptui, emoji/v2, kr/text, and kr/pty decisions/
evaluations, rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`. Treat
completed facts as final unless a narrow guard proves an input changed.

# Three Moves

First, reproduce the handoff and fixed guards. Second, complete and record only
the bounded pkg/errors evaluation. Third, update the roadmap and rolling
handover, answer this archive, prepare exactly one reciprocal successor for the
product decision required by the evidence, verify scratch containment, and
make the local handoff commit. Do not make that decision or execute the
successor. If an earlier guard changed, stop with no successor and request the
necessary owning reevaluation.

# Automatic Handoff

Do not push, merge, publish, release, stash, revert, bypass cleanup, remove the
worktree, change product source or dependency metadata, add a root, change a
requester, transfer an exception, evaluate another dependency, begin P8, or
write outside the managed scratch root.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No genuine supported direction qualifies. Exact selected, inherited, indirect,
unloaded `github.com/pkg/errors v0.9.1` remains byte-unchanged only because this
turn was an evaluation, not because retention passed. V0.9.1 is the canonical
latest stable and has genuine selected owner routes, but it fails the complete
required Go 1.26.7 upstream count-one, repeated, race, and vet gates. There is
no later stable upgrade, supported replacement, or guard-preserving removal.
No product source or dependency metadata changed, no exception was created,
and the reciprocal product decision was prepared but not executed.

### Handoff And Fixed Guards

Evaluation began clean on `codex/upgrade-quality` at handoff HEAD
`bf71bc950a601e38233d20e33a50c45f44f0e9c9`, parent
`413fb1a154a5c9a4acfe151dfb9518714a701ed1`, tree
`ab1d295120cac85d6f7c4d633c54d5e919cbafa2`. That handoff changes exactly
the launcher, answered Goe decision archive, this then-NEXT archive, rolling
handover, and roadmap. Google UUID implementation commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973` and the final Goe decision remain
ancestors. Branch, ordinary/ignored cleanliness, reciprocal archive chain,
single NEXT tail, and the contained launcher check reproduced.

The unchanged project reproduces 234 selected modules, 3,599 graph edges,
355 production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` remain
SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`;
the full graph is
`abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
The 432-line tidy diff remains
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes remain
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
Those exact project inputs reproduce the pre-Goe 47-selection/276-edge guard
at `7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`
and the separate complete Goe selection/request/route/import/load snapshot.
No earlier selection, owner, request, route, why/import/load/runtime fact,
exception, qualified result, or expiry boundary changed.

Official Darwin arm64 Go 1.26.7 archive/binary SHA-256 identities are
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`;
contained Go 1.18.10 identities are
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`.

### Selection, Requests, And Runtime Boundary

There is no main root or repository-history `go.mod` request for the target.
MVS selects v0.9.1 as an indirect graph module because two selected requesters
ask for that maximum; `go.sum` initially contains only the v0.8.0, v0.8.1,
and v0.9.1 synthetic module-file sums, not the v0.9.1 source sum. The ten
exact current and relevant historical requests are:

1. Kong pseudo-version `0548c6b1afae` -> v0.8.1.
2. Consul SDK v0.1.1 -> v0.8.1.
3. Consul SDK v0.8.0 -> v0.8.1.
4. SFTP v1.10.1 -> v0.8.1.
5. Selected SFTP v1.13.1 -> v0.9.1.
6. Prometheus Common v0.4.1 -> v0.8.0.
7. Prometheus Common v0.9.1 -> v0.8.1.
8. Prometheus TSDB v0.7.1 -> v0.8.0.
9. Selected Viper v1.15.0 -> v0.9.1.
10. Zap v1.17.0 -> v0.8.1.

Kong, both Consul SDK versions' `testutil` packages, both SFTP versions, both
Prometheus Common versions, TSDB, and Zap genuinely import pkg/errors. The
selected Viper v1.15.0 request is metadata-only. SFTP v1.13.1 is the genuine
selected source owner that makes v0.9.1 the MVS maximum; Viper's equal request
does not create a source import.

The complete shortest routes are main -> go-term-markdown v0.1.4 -> Chroma
v0.7.1 -> Kong; main -> direct mvn-pom-mutator -> bketelsen/crypt -> Consul
API v1.1.0 -> Consul SDK v0.1.1; main -> mvn-pom-mutator -> historical Viper
v1.10.1 -> Consul API v1.12.0 -> Consul SDK v0.8.0; main -> historical Viper
-> Afero v1.6.0 -> SFTP v1.10.1; main -> selected Afero v1.9.4 -> selected
SFTP v1.13.1; main -> historical Viper -> go-metrics -> Prometheus Common
v0.9.1 -> Prometheus Client v1.0.0 -> Common v0.4.1; main -> historical Viper
-> go-metrics -> Common v0.9.1; main -> mvn-pom-mutator -> TSDB v0.7.1;
main -> direct selected Viper v1.15.0; and main -> historical Viper -> etcd
client/pkg v3.5.1 -> Zap v1.17.0.

Only selected Viper has a positive requester `go mod why`; target and every
other requester are negative. The repository has no pkg/errors import. The
355 production, 429 complete-test, and 197 module-backed entries across 41
loaded modules contain no target package. Thus current project runtime
relevance is zero even though selected SFTP supplies a genuine source owner.

### Canonical Ownership And Release Identity

Go-import metadata, module declarations, and proxy records resolve the public,
enabled, unarchived, non-fork BSD-2-Clause `pkg/errors` repository, GitHub ID
48643510, default branch `master`, with no parent/source redirect. Its README
places the package in maintenance mode, rejects new functionality, welcomes
bug fixes and issue reports, and still describes a final 1.0 roadmap. Current
signed master is `87f8819acf6dc28bf5d3c14b334268236d686f48`, tree
`60652f0e917d39e5d310641579b61c4682d64164`; only CI/workflow commits follow
v0.9.1, with no later released source.

The exact path has 13 stable tags from v0.1.0 through proxy-latest v0.9.1,
12 GitHub Release objects, and no prerelease, `/v2`, `/v3`, replacement,
retraction, deprecation, redirect, or eligible alternate path. No branch,
fork, pseudo-version, or standard-library migration is an exact-path stable
release. All tags form continuous ancestry. V0.1.0-v0.8.1 are unsigned
annotated tags/commits, v0.9.0 is an unsigned lightweight commit, and selected
v0.9.1 is validly signed commit
`614d223910a179a466c1767a985424175c39b465`, parent
`49f8f617296114c890ae0b7ac18c5953d2b1ca0f`, tree
`6dd01fd9b7f97a850cc87788579cfc01fd6431fd`, timestamp
2020-01-14T19:47:44Z.

Every release lacks a repository `go.mod`; the proxy supplies the same
synthetic `module github.com/pkg/errors` file without a Go directive,
requirement, replacement, or retraction. Selected sumdb identities are source
`h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=` and module
`h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=`. Its proxy ZIP is SHA-256
`d4c36b8bcd0616290a3913215e0f53b931bd6e00670596f2960df1b44af2bd07`;
the normalized 16-file manifest is
`3447aa1b782b7fda22d0510849eeaaf9f54d352ec8d9107d214e48f6d1a0019b`.
All 13 proxy archives byte-match their exact Git regular-file manifests and
contain no symlink or submodule.

### API, Behavior, Closure, And Complete Gates

Selected v0.9.1 exposes As, Cause, Errorf, Is, New, Unwrap, WithMessage,
WithMessagef, WithStack, Wrap, Wrapf, Frame, and StackTrace plus the exported
format/marshal methods. The line evolves from the five-function v0.1.0 API;
v0.8.1 adds the final message/stack helpers and v0.9.0 adds Go 1.13 Is/As/
Unwrap support. V0.9.1 intentionally restores legacy Cause behavior: Cause
follows only a `Cause() error` chain, while Is/As/Unwrap follow standard error
chains. New/Errorf/Wrap/Wrapf capture stacks; wrapping nil returns nil; `%+v`
renders message/cause stacks. The ordinary upstream fixtures cover these
documented nil, chain, format, JSON, frame, and stack behaviors deterministically.

The single production package is standard-library-only with no module
dependency, cgo, generated, embed, network, filesystem, or subprocess
boundary. Legacy build constraints isolate Go-1.13 compatibility and a
benchmark. Exact Go 1.18.10 sees 41 production/80 complete-test standard-
library packages; Go 1.26.7 sees 61/125. Module verification passes under
both SDKs. With cgo disabled, package build and test compilation pass for
Darwin amd64/arm64, Linux amd64/arm64/386, Windows amd64/386, FreeBSD amd64,
and js/wasm under both SDKs.

Under exact Go 1.18.10, selected v0.9.1 passes build, count-one tests,
count-ten tests, race count-one, and vet. Under exact Go 1.26.7, build and all
18 cross gates pass, but the required complete native gates do not:

- ordinary test and race commands fail first because vet rejects two upstream
  non-constant calls to Wrapf and WithMessagef;
- vet alone reports those same two failures;
- even with vet disabled, count-one, every count-ten repetition, and the race
  run fail TestStackTrace because Go 1.26.7 reports the nested function as
  `TestStackTrace.TestStackTrace.func2.func3`, while the fixture requires
  `TestStackTrace.func2.1`.

No race report appears when vet is disabled; the required race gate still
fails because its complete test population fails. These ordinary upstream
defects are sufficient to reject exact retention. V0.9.1 is already latest,
so there is no genuine supported Go-1.18-compatible upgrade candidate. Older
stables are downgrades, not upgrade or replacement candidates. Replacing the
API with standard `errors` would discard pkg/errors stack/Cause/format
semantics and require external requester source changes; it is not an eligible
supported graph operation and introduces unreviewed behavior.

### Disposable Projections

An exact disposable v0.9.1 get preserves all 234 selections and 355/429/197/41
loads but is not a byte no-op because main has no root. It adds an indirect
root, the selected source sum, and one main edge, producing 3,600 edges,
75/1,068 module-file lines, and hashes
`d75a97c0e0f397eca75b2092d5a70944a0febada9f818f91321f5925a300f672` /
`28f51d0099dea28b2c1674161c04fce8fe300c681bbdf9cbb326936269985a56`.
Tidy removes that manufactured root and restores the common 52/948-line
projection, with v0.9.1 still selected. No later stable exists, so no fictional
upgrade projection was run.

The exact removal command succeeds only by changing real owners: it removes
direct mvn-pom-mutator, downgrades direct go-term-markdown v0.1.4 -> v0.1.2,
direct Viper v1.15.0 -> v1.3.2, and Afero v1.9.4 -> v1.2.2. Its raw state has
109 modules, 430 edges, 73/1,091 module-file lines, no target, hashes
`cde5f075021c6ac7d22bbf939201c260fb9b38019c56cb78630b43be868d85ef` /
`e9e0ec07eb1aa817efb6323ead7dd6256e6a41ebbe63567a1920deb39f48ba21`,
and cannot load because the product still imports mvn-pom-mutator while the
downgraded markdown route lacks a sum.

Tidy must rediscover direct mvn-pom-mutator and reintroduces pkg/errors v0.9.1,
but it does not restore the project: it leaves go-term-markdown v0.1.2, Viper
v1.10.1, Afero v1.8.2, only 228 modules/3,465 edges, and 346/420/188/41 loads,
with 53/941 metadata lines and hashes
`9ab84f943f7e79a8c5c73089d501a5f7f0e68a9e4f8741bf4347072a66a9242c` /
`353a9297c9efd86104737eac6a404d2b1a8395e491183003a562643a0f984203`.
That crosses direct-owner and numerous earlier-selection guards. Restoring the
original two module files returns the exact fixed state. Neither projection
was retained, and the Go 1.18 directive/toolchain lines never changed.

### Advisories And Final Verification

Fresh exact v0.9.1 OSV and GitHub global results are empty, as is the
repository advisory endpoint. Pinned govulncheck v1.8.0, built with exact Go
1.26.7 at binary SHA-256
`1e968c1a82faea943a15a01a29b014ac9df8cbb62c87c603b8189ee1fd3f2783`,
finds no target module, package, symbol, or test-symbol result. Unchanged-
project populations reproduce exactly 30 module, 22 vulnerable-package, 20
called-symbol, and 20 test-symbol OSVs with no pkg/errors trace.

Guard OSV remains limited to Gorilla WebSocket `GO-2026-6278` /
`GHSA-w67g-5rqw-f597` and go-retryablehttp `GO-2024-2947` /
`GHSA-v6v8-xj6m-xwqh`; x/mod v0.14.0 retains `GO-2026-6179` and
`GO-2026-6180`. The Go vulnerability index remains 518,501 bytes and 1,402
records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`.
The PUBLISHED memberlist CNA response remains 2,807 bytes at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Advisory absence did not override the ordinary upstream failures.

Final unchanged-project exact Go 1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass under ordinary `umask 022`. A first
restrictive-launcher-umask run produced only the known three mode-assertion
failures; the ordinary rerun is clean. Source/API/CLI/help/launcher/Make/
quality contracts and accepted 27/27 Q0-Q2 PASS at L2 remain unchanged.

All disposable SDKs, caches, tools, responses, reports, module/repository/
project copies, and advisory data remained below the managed session scratch
root and were removed before handoff; only the pre-existing launcher-owned
Node compile cache remains. No successor was executed, no other dependency
was evaluated, and P8 was not begun.
