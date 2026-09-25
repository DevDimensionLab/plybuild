# Agent Session: Evaluate Prometheus Client Golang Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-23T105846+0200-evaluate-prometheus-client-golang-dependency`
Created: `2026-09-23T10:58:46+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `45a2491eed08476facd0861778492eff9091b93883669a86257ab4d98e713fd7`
Previous: [2026-09-23T103802+0200-decide-prometheus-client-golang-guard-direction.md](2026-09-23T103802+0200-decide-prometheus-client-golang-guard-direction.md)
Next: [2026-09-23T124610+0200-decide-prometheus-client-golang-product-direction.md](2026-09-23T124610+0200-decide-prometheus-client-golang-product-direction.md)
Outcome: No canonical exact-path Go-1.18-compatible stable qualifies. Highest-floor v1.16.0 fails exact-SDK vet, Go 1.26.7 tests, and supported cross gates; its raw project get also changes protected go-conntrack and ordinary tidy returns selected v1.4.0. No dependency or source change was retained, and one reciprocal product decision was prepared.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with one fresh, complete evaluation of graph-selected transitive
exact `github.com/prometheus/client_golang v1.4.0`. Resolve its canonical
exact-path release line and highest fully qualified Go-1.18-compatible stable
from primary evidence. Implement one exact dependency-only changed selection
only if the candidate, its complete minimal closure, a genuine supported tidy-
stable project owner, and every earlier target-specific guard remain exact. Do
not combine another dependency group or begin P8.

# Defensive Scope

This is an ordinary dependency-quality evaluation. Use public metadata,
static source and repository records, project graph/build commands, upstream
tests, and only small bounded ordinary fixtures required by documented
behavior. Do not fuzz, stress, probe resource exhaustion, create oversized,
deeply nested, cyclic, malformed, adversarial, or escape-sequence payloads,
reproduce a security issue, or perform security or exploitability analysis.

Every disposable cache, tool, archive, report, project copy, fixture, or
advisory response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`.
Never write to `/private/tmp`, `/tmp`, a sibling of the managed root, or
another external root. Verify containment and remove task-owned scratch
evidence before handoff.

# Fresh Evaluation Boundary

The preceding guard decision selected option 1. The former NEXT archive's two
advisory hashes were transcription/provenance defects. The active guards are
the independently reproduced preceding-archive identities below. The stopped
evaluation's partial release, source, behavior, closure, projection, API,
test, cross-build, and advisory matrix is discarded. Begin from a newly
verified clean state with new contained evidence. Do not resume a row, reuse a
partial closure result, infer qualification from preliminary evidence, or
accept any previous candidate selection.

Exact selected `github.com/prometheus/client_golang v1.4.0` and its four graph
requests remain unchanged until this fresh evaluation completes:

- `github.com/armon/go-metrics v0.3.10` requests v1.4.0;
- `github.com/prometheus/common v0.9.1` requests v1.0.0;
- `github.com/prometheus/common v0.4.1` requests v0.9.1; and
- `github.com/prometheus/tsdb v0.7.1` requests v0.9.1.

The stopped turn observed genuine imports at all four requesters, negative
target/requester why, zero repository target imports, zero target loads and
runtime relevance, and no current or historical target root. It also observed
51 exact-path stables, six prereleases, 33 stables with declared floors no
higher than Go 1.18, and v1.16.0 as the highest by declared floor. These are
starting observations only. Independently reproduce them and do not treat
them as release, ownership, behavior, closure, projection, or qualification
results.

# Protected Starting State

The guard decision began clean on `codex/upgrade-quality` at stopped-evaluation
handoff HEAD `ff4bcc1b009662d30041e0489fcb2b01c35cbe52`, parent
`438178de84797dad6178f1011560c53e192f1e0c`, tree
`fc665bf6bfac864f2f4198dfca590e300650045c`. Its exact changed set was the
launcher, answered stopped evaluation, then-NEXT guard decision, rolling
handover, and roadmap. Google UUID implementation commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor. Verify the
new decision handoff, its reciprocal archive chain, sole NEXT state, launcher
mirror, exact changed set, ancestry, and ordinary and ignored cleanliness.

The unchanged real project has 234 modules, 3,599 graph edges, 355 production
entries, 429 complete-test entries, 197 module-backed entries across 41 loaded
modules, and 1,067 `go.sum` lines. `go.mod` is 74 lines at SHA-256
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`;
`go.sum` is SHA-256
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`;
the graph is SHA-256
`abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
Normal tidy has the established 432-line diff at SHA-256
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`
and common 52/948-line, 234-module/3,557-edge projection with `go.mod` /
`go.sum` hashes
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
No projection is retained. The module floor remains Go 1.18 and accepted
quality remains 27/27 Q0-Q2 PASS at L2.

The exact official SDK identities remain Go 1.18.10 archive/binary
`718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
`f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`
and Go 1.26.7 archive/binary
`020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Use only exact verified binaries for named SDK gates.

All 47 pre-Goe guarded selections and 276 incoming edges remain exact at
SHA-256 `7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
Preserve the separate Goe, pkg/errors, SFTP, go-difflib, and Complete guards,
plus ULID, go-conntrack, and every earlier selection, request, route, why/
import/load/runtime, graph, module, tidy, Go-floor, source, behavior, closure,
qualification, advisory, and expiry fact. No exception transfers or reopens.

# Corrected Advisory Guards

The active Go vulnerability-index guard is exactly 518,501 bytes and 1,402
records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. The active CVE-2026-14362 CNA guard is
PUBLISHED, exactly 2,807 bytes, at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`;
the endpoint returns `Cache-Control: no-store`. Revalidate them narrowly from
their exact primary URLs. A changed response is new input and requires a fresh
owning guard decision; do not repair another mismatch inside this evaluation.

The unchanged project advisory guard remains 30/22/20/20 module/package/
symbol/test-symbol findings. Gorilla WebSocket v1.4.2 retains
GHSA-w67g-5rqw-f597 / GO-2026-6278; go-retryablehttp v0.5.3 retains
GHSA-v6v8-xj6m-xwqh / GO-2024-2947; x/mod v0.14.0 retains GO-2026-6179 and
GO-2026-6180. Selected client_golang v1.4.0 has
GHSA-cg3q-j54f-5p7p / GO-2022-0322. Advisory absence for another version
cannot establish qualification.

# Evaluation Contract

Resolve the exact go-import owner, repository identity and state, canonical
exact-path stable/prerelease/replacement/retraction/deprecation lines, tags,
commits, trees, signatures, ancestry, proxy/sumdb, archive-to-Git, license,
module directives, Go floors, packages, build tags, cgo, generated/embed
boundaries, API, and documented ordinary behavior for every eligible stable.
Exclude branches, pseudo-versions, prereleases, replacements, alternate module
paths, forks, patches, and vendoring unless explicitly authorized; none is.

For every eligible stable, evaluate the complete minimal upstream module
closure under exact Go 1.18.10 and Go 1.26.7 with module verification, build,
complete count-one and repeated tests, race where supported, vet, and the
bounded supported cgo-disabled cross-build/test-compilation matrix. Use only
small ordinary deterministic fixtures where upstream coverage does not settle
documented behavior. Do not execute unrelated network services or destructive
paths. A release qualifies only if its complete closure and ordinary public
behavior qualify, it preserves the Go 1.18 floor, and it has a genuine
supported tidy-stable owner/request route that preserves every earlier guard.

Reproduce all exact target requests and requester import boundaries, complete
current and historical routes from main, owner support, `go mod why`,
repository imports, production/complete-test/module-backed loads, runtime
relevance, and current/history roots. In scratch-contained project copies
only, measure the selected exact get and at most the one highest otherwise-
qualified candidate get, followed by normal tidy, to establish exact closure
and projection. Do not retain a projection that removes or changes a direct
root, earlier guarded selection/request/route, product behavior, Go floor, or
accepted quality.

# Decision And Implementation Boundary

If and only if one canonical exact-path stable is the highest fully qualified
Go-1.18-compatible candidate, has a genuine supported tidy-stable project
owner, and its complete scratch projection preserves every guard, implement
exactly `go get github.com/prometheus/client_golang@<selected>` followed by
ordinary tidy and the full project gates. Retain only the exact measured
minimal dependency closure. Do not add a direct root, change source, patch,
vendor, fork, alter an unrelated requester, select an alternate path, or
combine another group.

If no stable qualifies, no supported tidy-stable owner exists, a projection
crosses any guard, or any starting input changed, retain no dependency or
source change. Record the completed fresh evaluation and prepare exactly one
reciprocal product decision for this target; do not silently retain it, grant
an exception, authorize an owner study, or move to another dependency in the
same turn.

# Required Reading And Checks

Read this archive; the answered guard decision and stopped evaluation; the
answered Complete decision and evaluation; rolling handover; P7/P8 roadmap;
`go.mod`; and `go.sum`. Verify branch/HEAD/parent/tree, exact changed set,
ancestry, ordinary and ignored cleanliness, reciprocal archive chain,
launcher check, exact SDK identities, project hashes/counts/tidy state, every
target request/route/why/import/load/runtime/root fact, all earlier guards, and
narrow fresh advisory identities. Stop for a fresh owning decision if any
protected input changed.

# Three Moves

First, independently complete this one target's fresh release/source/behavior/
closure/ownership/advisory evaluation and measure only authorized contained
projections. Second, either retain the one exact fully qualified dependency-
only selection or retain no change and prepare one reciprocal product
decision. Third, update roadmap and rolling handover, answer this archive,
prepare at most that one successor, verify containment and cleanup, run final
exact-Go project gates, and make the local handoff commit.

# Automatic Handoff

Do not launch a successor or study; push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, transfer or reopen an exception,
change another dependency group, write outside the managed scratch root, or
begin P8.
<!-- CODEX_SESSION_PROMPT_END -->

## Answer

No client_golang selection is authorized. All 33 canonical exact-path stable
versions whose declared Go floor does not exceed 1.18 fail at least one
mandatory complete-closure gate. Highest-floor v1.16.0 is therefore not fully
qualified; it also has no genuine supported tidy-stable project owner, and its
raw project projection changes the separately protected go-conntrack
pseudo-version. Product source, `go.mod`, and `go.sum` remain byte-exact, and
no projection was retained.

### Continuity and protected inputs

- Work began clean on `codex/upgrade-quality` at guard-decision handoff HEAD
  `eea09d0057ad57278945ce9fbba02eba46017c10`, parent
  `ff4bcc1b009662d30041e0489fcb2b01c35cbe52`, tree
  `488b0b89d05c48f2ffbe390813b9302c50967cfd`. That handoff changes exactly
  the launcher, answered stopped evaluation, then-NEXT guard decision, rolling
  handover, and roadmap. Google UUID implementation
  `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an ancestor.
- The reciprocal archive chain, sole NEXT state, launcher/archive prompt
  mirror and check, exact changed set, and ordinary and ignored cleanliness
  reproduced. No stopped-turn release, source, test, projection, API, or
  advisory result was reused.
- Fresh official Go 1.18.10 archive/binary SHA-256 is
  `718b32cb2c1d203ba2c5e6d2fc3cf96a6952b38e389d94ff6cdb099eb959dade` /
  `f96ea900187be55d1be92addc093c44af2a437f50b58e2d56b05bc8a37b29c74`;
  Go 1.26.7 is
  `020a1e8224811be75163e920bc77e0926a1390a6aeea19bdcf23f74b9d749f6d` /
  `9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
  Only those verified binaries ran named SDK gates.
- Fresh normal and no-cache requests to
  `https://vuln.go.dev/index/modules.json` are pairwise byte-identical at
  518,501 bytes, 1,402 records, SHA-256
  `bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
  Last-Modified 2026-09-17T17:29:18Z. Equivalent requests to
  `https://cveawg.mitre.org/api/cve/CVE-2026-14362` are pairwise byte-identical
  PUBLISHED 2,807-byte responses at
  `cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`
  with `Cache-Control: no-store`. The mandatory repaired guard passes.
- The real project independently reproduces 234 modules, 3,599 graph edges,
  355 production entries, 429 complete-test entries, 197 module-backed
  entries across 41 loaded modules, and 1,067 sum lines. `go.mod`, `go.sum`,
  and graph SHA-256 remain
  `7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874`,
  `87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`,
  and `abdac9686ca5aeb51d5d413bf98d58e2cb7ea3244272d815784cdabe1c104cdf`.
  Normal tidy again produces the exact common 52/948-line files at
  `5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
  `b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`
  with 234 modules and 3,557 edges. The known local Git presentation is 431
  lines; applied bytes match the protected 432-line projection exactly.

### Canonical identity and release line

- Exact go-import metadata maps `github.com/prometheus/client_golang` to
  `https://github.com/prometheus/client_golang.git`. GitHub repository ID
  7823926 is public, enabled, unarchived, non-fork, Apache-2.0, owned by
  `prometheus`, and uses default branch `main`. Fresh `main` is commit
  `c946c36aec713ed74cfe5953fb8ff84777ef24f1`, tree
  `9f8f820d49615efb3dbe12ba458e29092cce48e5`.
- Proxy and `go list -m -versions` expose exactly 51 exact-path stables and six
  prereleases through v1.24.1. `/v2` and `/v3` return no module line. There
  are no replacements, retractions, or deprecations. Thirty-three stables
  have no declared floor or declare Go 1.11, 1.13, or 1.17; v1.16.0 is the
  highest whose declared floor does not exceed Go 1.18. Later stables declare
  Go 1.19 or newer and are ineligible. V1.13.0 through v1.16.0 exclude
  v1.12.1; this is neither a retraction nor replacement.
- All 33 eligible proxy archives match their mapped Git commit byte-for-byte
  across 40 to 139 regular files, with no proxy/Git mismatch, symlink, or
  submodule. Their sumdb source and module sums verify. Apache-2.0 `LICENSE`
  and `NOTICE` are byte-identical throughout at SHA-256
  `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4` /
  `5ff7ccec357cd44bb3fcf1b3b53a4eddee322eb383d487847d1cf1ccfbd94f46`.
- Thirty-one release commits have GitHub-valid signatures; v0.8.0 and v1.12.2
  are unsigned. V0.12.1 is an immutable proxy/sumdb release whose current Git
  tag has been deleted; its 115-file archive exactly matches valid-signed
  ancestor commit `77626d64fa02954e546be20b688e235617e42c7e`. V1.11.1 is a valid-signed
  release-branch commit not ancestral to current `main`. Every other eligible
  release commit is a `main` ancestor.

Exact eligible commit/tree identities are:

| Version | Commit | Tree |
| --- | --- | --- |
| v0.8.0 | `c5b7fccd204277076155f10851dad72b76a49317` | `dff442c299b943a1c0346f153ed4a8b6a2554034` |
| v0.9.0 | `1cafe34db7fdec6022e17e00e1c1ea501022f3e4` | `daa2fcc188433ce0a24f867d674c5c681ac284c5` |
| v0.9.1 | `abad2d1bd44235a26707c172eab6bca5bf2dbad3` | `396fc4534a012a89b6fde24cc17c62caf0be69fc` |
| v0.9.2 | `505eaef017263e299324067d40ca2c48f6a2cf50` | `73f5a72f1e4f252c89927b6c1670bdebb6774269` |
| v0.9.3 | `50c4339db732beb2165735d2cde0bff78eb3c5a5` | `aac7cf0632024449e1cb20f6f668d4b5e0ee29ca` |
| v0.9.4 | `2641b987480bca71fb39738eb8c8b0d577cb1d76` | `68870e1b488f76e0cbb9cacc91cdec71919e5907` |
| v0.12.1 | `77626d64fa02954e546be20b688e235617e42c7e` | `758aea79b0b2f50de494e822dad7d5dfe646801d` |
| v1.0.0 | `4ab88e80c249ed361d3299e2930427d9ac43ef8d` | `b0e0265d6fc3cdd9cca490ea93b91bec8fb755ba` |
| v1.1.0 | `170205fb58decfd011f1550d4cfb737230d7ae4f` | `1f52018fd54faa12741c8d9ef00dbabe5a3e10ea` |
| v1.2.0 | `9a2ab94c6aab59eb8382178905c1a9339c24d124` | `a49b7698fb766630c6db403a2afbb311cc5bd7c3` |
| v1.2.1 | `55450579111f95e3722cb93dec62fe9e847d6130` | `af74e5da06d23fdd0590717578a416eb1fcf666e` |
| v1.3.0 | `c42bebe5a5cddfc6b28cd639103369d8a75dfa89` | `340befb2cb7d7d45c078f51ee13c0e10843c9595` |
| v1.4.0 | `76dd6c581988366a52807465d426f83e776128ad` | `399889725112fc64b42014ef7249e2608722dacc` |
| v1.4.1 | `913f67ef0627596efc752f02d4014308e7345dbf` | `f4401f09129049b19c8bb974ba8c6b07f0c51b01` |
| v1.5.0 | `5538bedeb68bb9597b90500bdd485e1541aeb3cf` | `9540d50bd59e507a287ed80bc4af71bb096ddb82` |
| v1.5.1 | `aa9238db679fc02bf22cf4f2c27a980edcb5ada0` | `334a3e6249a9547123af818549e80210ee312140` |
| v1.6.0 | `6edbbd9e560190e318cdc5b4d3e630b442858380` | `5f196c241a44d18567dc567640f494eede3476a8` |
| v1.7.0 | `b05e50c9298cbadd2b4be3f1527f5c1b178acb55` | `9118f197bed131d8cd37f59de1862c02e7dc0df1` |
| v1.7.1 | `fe7bd95816c4bf8b6482e529a919a0957e631a6e` | `71b52682f5d4713bcbb246633116d611aba5971d` |
| v1.8.0 | `47cfdc9bb8ceaf18b12a4b7ba778b06e12d1a97a` | `5cd361803b1686312cb639310b59f7ac70a30c16` |
| v1.9.0 | `d89cf5af8849fcb6d69f2c16bdf3b910e03d341e` | `af2c333e80190962ffb06926f74486f3d03c64cc` |
| v1.10.0 | `efe7aa73021f40ba96d167c144941fa5fb7f43e2` | `159b7547c8f25eac96eb9d0c63887a3831d98c16` |
| v1.11.0 | `8184d76b3b0bd3b01ed903690431ccb6826bf3e0` | `7d8c9c5c6c149e5bdc891e287a700a9d76ce7709` |
| v1.11.1 | `989baa30fe956631907493ccee1f8e7708660d96` | `dba43f2e9826ec7eec4124941e3b7c5a4583f33b` |
| v1.12.0 | `01087964d02726ef7438fdb19c2672e6e5d0b48a` | `a540c198182c1a4345c6983907d172e38426ee9d` |
| v1.12.1 | `2e1c4818ccfdcf953ce399cadad615ff2bed968c` | `4a7685bd90b3b446fd188112f6e134d998a9c6c2` |
| v1.12.2 | `e203144f43306c1f344fbc548fd02c4b79962e30` | `fc1600cd1ccc0c286cec26045b7d58d12f331517` |
| v1.13.0 | `64435fc00ac419bb878a3f9c9658e8353c19a7cd` | `b5a522b2f96ab936b5a0cd3878f6171578a12123` |
| v1.13.1 | `53e51c4f5338f760a766232610e574b00ea720d8` | `42b3c3fff138f4736e0f0f450ec59809a5a9fda3` |
| v1.14.0 | `254e5468413f19fb75cdad45f5ddc0b8c975188c` | `479c3234f0308aebd2a32c5bebc9a6277855df46` |
| v1.15.0 | `d7896d4bd082b17e525c29055d79cc29484aa9cb` | `53f58baadfc7d31f6f35a0c33f907af5fb2473c9` |
| v1.15.1 | `4bbb297e54556fd1bcec279b355d049a50fb316e` | `4c2ea848e3456e01e2a6228954dc7866f35b33f9` |
| v1.16.0 | `3583c1e1d085b75cab406c78b015562d45552b39` | `5049d8142bf27f33903b95a9071e4f4f2173a7b5` |

Selected v1.4.0 source/module sums are
`h1:YVIb/fVcOTMSqtqZWSKnHpSLBxu8DKgxq8z6RuBZwqI=` /
`h1:e9GMxYsXl05ICDXkRhurwBS4Q3OK1iX/F2sw+iXX5zU=`; its 89-file proxy ZIP
SHA-256 is `7d43c60d49b3a110545074bdf57812f1f8d0484ce516eabbc29ee93bc8dc2e6b`.
V1.16.0 sums are `h1:yk/hx9hDbrGHovbci4BY+pRMfSuuat626eFsHb7tmT8=` /
`h1:Zsulrv/L9oM40tJ7T815tM89lFEugiJ9HzIqaAx4LKc=`; its 139-file ZIP
SHA-256 is `0167cee686b836da39815e4a7ea64ecc245f6a3fb9b3c3f729941ed55da7dd4f`.

### Requests, routes, ownership, and runtime boundary

- Exactly four graph requests reproduce: go-metrics v0.3.10 -> v1.4.0,
  Common v0.9.1 -> v1.0.0, Common v0.4.1 -> v0.9.1, and TSDB v0.7.1 ->
  v0.9.1. All are genuine imports. Go-metrics imports `prometheus` and
  `push` in production and `prometheus` in tests; both Common versions import
  `prometheus` in production `version/info.go`; TSDB imports `prometheus` and
  `promauto` in production and `prometheus/testutil` in tests.
- Every requester has one shortest route. Go-metrics is main -> direct
  mvn-pom-mutator v0.2.3 -> historical Viper v1.10.1 -> go-metrics v0.3.10.
  Common v0.9.1 extends that route through go-metrics; Common v0.4.1 extends
  it through Common v0.9.1 -> client_golang v1.0.0; TSDB is main -> direct
  mvn-pom-mutator -> TSDB v0.7.1. The selected target route is the go-metrics
  route followed by client_golang v1.4.0.
- Target and requester module why results are negative. The repository has no
  target import; target and requester packages are absent from all project
  production and complete-test loads. Target module-backed load and runtime
  relevance are zero. No current or historical main-module target root exists.
  `go.sum` contains only v0.9.1, v1.0.0, and selected v1.4.0 module-file sums,
  with no selected source sum.
- Exact selected v1.4.0 is genuinely requested and tidy-stable through the
  historical Viper/go-metrics route. No requester asks for v1.16.0. Adding it
  as a redundant main root does not make a supported owner: ordinary tidy
  removes it and reselects v1.4.0.

### Source, behavior, API, and complete closure

- Every eligible release is pure Go with no cgo or `go:embed`. Selected
  v1.4.0 has 71 Go files, 11 production package directories, and four raw
  build constraints. V1.16.0 has 112 Go files, 15 production package
  directories, four generated files, and a version/platform build-tag matrix
  for runtime collectors, js/wasm, Windows, Linux, and legacy Go releases.
- The public surface provides concurrency-safe counters, gauges, histograms,
  summaries, registries, gathering, vector deletion/reset/currying, process
  and Go collectors, HTTP exposition/instrumentation, Pushgateway operations,
  and test utilities. `Must*` constructors and registration helpers document
  panic paths; ordinary constructors return errors. Shared default registries
  and exported vector/registry state are mutable; collectors, descriptors,
  label ownership, gather ordering, timers, exemplar paths, HTTP error paths,
  body closure, and push behavior are covered by upstream tests. No additional
  fixture was needed and no unrelated network service was run.
- Selected v1.4.0 has a 41-module build list. Exact Go 1.18.10 closure counts
  are 11 target packages, 158 production entries, 193 complete-test entries,
  and 47 module-backed entries across 11 loaded modules; Go 1.26.7 gives
  220/256/47 across the same 11. Its exported API/documentation snapshots are
  2,902 lines at
  `8d273b3a62bd5a67dba03be6894a3e4badf65ba48c1c3f67328642eb89e059fc`
  and 2,885 lines at
  `3ce80697d7b19995c2c38da447bd7ea6329348eee15caa6ad227c73a17010afe`.
- V1.16.0 has a 46-module build list. Go 1.18.10 counts are 15 target
  packages, 194 production, 258 complete-test, and 104 module-backed entries
  across 20 loaded modules; Go 1.26.7 gives 256/321/104 across the same 20.
  API/documentation snapshots are 4,009 lines at
  `ac1c375463ca7f17114f11f823a9e8d1ad188ea03c9a5d0c98d7007b153c621e`
  and 3,997 lines at
  `2ee803af7c5706198f48b6e2bf43f900edb5da80591f6ea9e4f858db636d5af6`.
- The complete native matrix contains 660 rows: for every eligible release
  under both exact SDKs, module download/verification, build, count-one and
  count-ten tests, race, vet, production/complete-test enumeration, and API
  documentation were attempted. All module downloads and verifications pass,
  but every release fails at least one mandatory gate; there is no fully
  qualified stable.
- Selected v1.4.0 builds under both SDKs but fails tests, repeated tests,
  race, and vet under both. Go 1.18.10 includes a reflect2/json-iterator crash
  and an obsolete HTTP error-string expectation; Go 1.26.7 includes API test
  panics, the same error-string mismatch, and vet/example failures.
- Highest-floor v1.16.0 passes Go 1.18.10 build, count-one/count-ten tests, and
  race but fails vet on invalid example identifiers, unkeyed external
  literals, and integer-to-string conversions. Under Go 1.26.7 it builds but
  tests, repeat, race, and vet fail because its generated runtime-metric test
  support ends before the current Go line, leaving collector expectation
  identifiers undefined. Thus v1.16.0 is not qualified independently of
  ownership or projection.
- The cgo-disabled cross matrix contains 1,320 production-build and test-
  compilation rows for Darwin amd64/arm64, Linux amd64/arm64/386, Windows
  amd64/386, FreeBSD amd64, Plan 9 amd64, and js/wasm under both SDKs. V1.8.0
  alone passes every cross row but fails native tests/race/vet. V1.16.0's
  32-bit tests overflow untyped millisecond constants, its procfs closure does
  not build on Plan 9, and every Go 1.26.7 test compilation hits the missing
  runtime-metric expectations. No cross result rescues a native failure.

### Projections, guards, and advisories

- Disposable exact selected get adds only a redundant indirect target root,
  one main edge, and the source sum. Raw state is 234 modules, 3,600 edges,
  75/1,068 lines, unchanged 355/429/197/41 loads, negative target why, and
  `go.mod` / `go.sum` / graph SHA-256
  `e9ba9100af4cbab440527fc52f283e6854f520c8a160ec80953d99d217a3ae60` /
  `7f1d5c01c64db9623689b294853ad417bb9f5bd9702dc8f65322f8d55864b24b` /
  `8012776a70d08ddd52dab8cb2f384c06d805e8243205bde2c0772b66c9e95c3a`.
  Tidy removes it and restores the common projection with v1.4.0 selected.
- Disposable v1.16.0 get yields 235 modules, 3,622 edges, 75/1,069 lines,
  unchanged loads, negative target why, and hashes
  `dc7337944ce0ca77e223ae209a34d25d67e62c39b78745b27382f12220346a3e` /
  `d896cd82b32a95e23cf34e79d080436598d26f4fcca75ce30e761d7b51aed644` /
  `521a4eb53dd630539ca54f9b8ef42142316429a66251f89b52e2f3653fb69758`.
  It adds jpillora/backoff and upgrades client_model, Common, procfs,
  protobuf extensions, oauth2, protobuf, and protected go-conntrack from
  `cc309e4a2223` to `2f068394615f`. Tidy removes all nine changes and returns
  the exact common projection with v1.4.0. It is neither guard-preserving nor
  tidy-stable. Neither projection was retained.
- The exact full project graph proves the 47-selection/276-edge pre-Goe guard
  remains at
  `7f2d3038c8376a7de6790d21b2a496017182f3171ec3335af11387b542c70d3d`.
  Separate selections/request counts reproduce: Goe v0.1.0/four,
  pkg/errors v0.9.1/ten, SFTP v1.13.1/four, go-difflib v1.0.0/23, Complete
  v1.2.3/three, ULID v1.3.1/one, and go-conntrack `cc309e4a2223`/two. Their
  recorded why/import/load/runtime, route, source, behavior, closure,
  projection, qualification/exception, advisory, and expiry guards remain
  unchanged and untransferred.
- Fresh exact-version OSV identifies GHSA-cg3q-j54f-5p7p /
  GO-2022-0322 / CVE-2022-21698 for 23 eligible releases through v1.11.0,
  including selected v1.4.0. GitHub global and repository advisory endpoints
  identify the same sole repository advisory for v1.4.0. V1.16.0's narrow OSV
  and GitHub global responses are empty; absence is not qualification.
- Guard OSV reproduces Gorilla WebSocket GHSA-w67g-5rqw-f597 /
  GO-2026-6278, go-retryablehttp GHSA-v6v8-xj6m-xwqh / GO-2024-2947, and
  x/mod GO-2026-6179 and GO-2026-6180. Pinned govulncheck v1.8.0 is built by
  exact Go 1.26.7 and uses database timestamp 2026-09-16T18:00:43Z. The
  unchanged project reproduces exactly 30 module, 22 package, 20 symbol, and
  20 test-symbol non-stdlib findings with no target trace.

### Decision and handoff

No canonical stable qualifies, and the highest-floor candidate additionally
lacks supported tidy-stable ownership and crosses an earlier guard. Therefore
no `go get`, tidy, dependency implementation, direct root, source change,
exception, or owner study is authorized or retained.

Final unchanged-project exact-Go-1.26.7 module verification, build, count-one
tests, race count-one tests, and vet pass after documentation handoff. The
launcher/archive mirror, reciprocal chain, sole NEXT state, diff checks,
project hashes, and clean ordinary/ignored state pass. Every task-owned SDK,
archive, repository, cache, project copy, response, report, tool, and fixture
was contained under `${CODEX_SESSION_SCRATCH_ROOT:?}` and removed; only
launcher-owned scratch state remains.

Exactly one reciprocal client_golang product-decision successor is prepared
and not executed. It may choose target-specific unqualified retention of exact
v1.4.0, authorize one later measurement-only Mvn-Pom-Mutator client_golang
owner/request-removal study, or stop P7 unresolved. It may not repeat this
evaluation, change a dependency or source, transfer an exception, combine
another group, or begin P8.
