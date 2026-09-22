# Agent Session: Decide Mitchellh Mapstructure Product Direction

Status: NEXT
Session ID: `2026-09-22T202511+0200-decide-mitchellh-mapstructure-product-direction`
Created: `2026-09-22T20:25:11+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3098013c91be7cda8e4cd2c4d14bee151aee65360ae14ad32dc4de7cda018b85`
Previous: [2026-09-22T192114+0200-evaluate-mitchellh-mapstructure-dependency.md](2026-09-22T192114+0200-evaluate-mitchellh-mapstructure-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by recording exactly one bounded product decision for selected
indirect exact-path `github.com/mitchellh/mapstructure v1.5.0`. The completed
evaluation found that no canonical exact-path stable qualifies: every one of
the 17 Go-1.18-compatible stable releases retains stale destination slice
elements during an ordinary documented decode, and v1.2.2 also fails its
repeated upstream suite. The exact-path repository is archived and its owner
states that it will never tag another release; the blessed active fork uses
alternate module path `github.com/go-viper/mapstructure/v2` and is not an
eligible exact-path release. Choose only one of the three directions below,
apply it exactly, and prepare only its reciprocal successor if the choice
requires one. Do not repeat the dependency evaluation, evaluate another
dependency group, or begin P8.

# Defensive Decision Scope

This is an ordinary dependency-quality product decision. Use only the
completed public release/repository metadata, static source and graph facts,
small bounded ordinary map-decoding results, project projections, and
advisory identities recorded here and in the answered evaluation. Do not
fuzz, stress, probe resource exhaustion, create oversized, deeply nested,
cyclic, malformed, adversarial, or escape-sequence payloads, reproduce a
security issue, or perform security or exploitability analysis.

Every disposable cache, tool, report, project copy, fixture, or advisory
response must remain beneath `${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write
to `/private/tmp`, `/tmp`, a sibling of the managed root, or another external
root. Verify containment and remove task-owned scratch evidence before
handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after exact Go 1.26.7, every accepted
dependency move through Google UUID v1.4.0, qualified go-cleanhttp v0.5.2,
and target-specific option-1 decisions through exact selected, direct,
production-loaded, runtime-relevant go-homedir v1.1.0. P8 remains queued.

Go-homedir, Promptui, emoji/v2, kr/text, and kr/pty option-1 decisions are
separately exact, unqualified, target-specific, non-transferable, and final.
Their selected and historical requests, requester imports, complete routes,
why/import/load/runtime facts, release/repository/source identities,
qualification results, graph/module/tidy/Go-floor state, projections,
advisories, all earlier guards, and compatible-route conditions remain
expiry guards. Any change requires the corresponding fresh dependency and
product decision before merge. Do not transfer an exception or reopen
go-homedir, Promptui, emoji/v2, kr/text, kr/pty, kr/pretty, Cast, Viper, or an
earlier decision.

# Completed Evaluation Is Final

The evaluation began from clean branch `codex/upgrade-quality` at handoff
HEAD `bff26c941f30ef51a4f5a35c064a558ef1b45d16`, parent
`0f2bf5aa44b1c0f57003cea4f17dad96740e55d2`, tree
`6771799cf803b4591bcb56d086aeed8c534e2019`. Exact Google UUID v1.4.0
dependency commit `cf53bc64eeb69471d35c7536d196bf1da15f3973` remains an
ancestor. The evaluation changed no product source or dependency metadata
and prepared this decision-only handoff. Verify the new handoff HEAD, parent,
tree, exact changed set, ancestry, reciprocal archive chain, and clean
ordinary and ignored status rather than assuming them.

Exact go-import metadata and all release module files resolve
`github.com/mitchellh/mapstructure` to public, enabled, archived, non-fork
MIT repository `mitchellh/mapstructure`, ID 10166531, owned by `mitchellh`
and defaulting to `main`. The proxy exposes exactly 17 stable releases from
v1.0.0 through v1.5.0; `@latest` is v1.5.0 at
2022-04-20T22:31:31Z. Exact `/v2` and `/v3` lines are absent. There are no
GitHub Releases, retractions, module deprecations, replacements, or
prereleases. All 17 tags are lightweight and form continuous ancestry;
v1.0.0 has an unverifiable old key and the later 16 commits are validly
signed. V1.0.0-v1.1.2 have no Go directive and v1.2.0-v1.5.0 declare Go
1.14. All have no requirements, replacements, or retractions.

Selected v1.5.0 is validly signed commit
`ab69d8d93410fce4361f4912bb1ff88110a81311`, tree
`5a1166013faa55170e1acce958c1a3863faf3f80`, parent
`bd687ea300c090473812a1a5730c3a335fbb5b72`. Its proxy source/module sums
are `h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=` and
`h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=`. Its proxy ZIP has
SHA-256 `118d5b2cb65c50dba967fb6d708f450a9caf93f321f8fc99080675b2ee374199`
and byte-matches its exact 13-file Git regular-file manifest at
`7c590ec67ae92865b8628f7f8c7d15487f31c7b83e4846941482a051920bcf09`.
Sixteen proxy archives byte-match current tag manifests. V1.2.0 is the sole
historical tag-mutation exception: its immutable proxy/sumdb archive matches
historical commit `047abd31f2839526c057afa2966f1bc4a374d025`, while the current
tag points to signed descendant `9e4011917e467353a5cacb6b11f3991b381055c8`
and differs only by two CHANGELOG lines. The selected archive is unaffected.

The original owner states that archived repositories will accept no further
issues or pull requests and will never receive another tag, and names
`go-viper/mapstructure` as the blessed fork. That active repository is a fork
of the exact-path repository but its current module declaration is
`github.com/go-viper/mapstructure/v2` with Go 1.18. It is an alternate path,
not an eligible release or owner for the exact module line. No canonical
exact-path stable therefore has a genuine supported owner.

V1.5.0 is a standard-library-only, single-package 13-file module with no
build tags, cgo, generated code, embed, platform-specific implementation,
symlink, submodule, package-owned mutable global, external I/O, or cleanup
boundary. Its exported surface comprises Decode/WeakDecode variants,
NewDecoder and Decoder.Decode, DecoderConfig and Metadata, Error, decode-hook
types/execution/composition, and conversion hooks. Decoder/config/output/
metadata/hook state is caller-owned; callers must not reuse a config after
NewDecoder and must synchronize shared state. Independent decoders are
deterministic and race-free in the completed bounded fixtures. Errors are
sorted for deterministic rendering. The reflection and conversion boundary
is the core behavior under evaluation.

Under exact Go 1.26.7 and contained Go 1.18.10, all 17 releases pass module
verification, build, count-one tests, race count-one tests, vet, and Darwin
amd64/arm64, Linux amd64/arm64/386, Windows amd64/386, FreeBSD amd64, Plan 9
amd64, and js/wasm cross-builds. Repeated count-ten tests pass except at v1.2.2,
whose upstream metadata-order assertions fail under both SDKs. Selected
v1.5.0 has complete production/test closures of 76/133 packages under Go
1.26.7 and 57/88 under Go 1.18.10, with three module-backed entries and one
loaded module. Small bounded ordinary deterministic, invalid-output, and
four-independent-decoder concurrency fixtures pass count-ten and race-count-
ten under both SDKs.

The decisive ordinary fixture decodes a one-element string slice into a
pre-populated two-element destination. Every stable from v1.0.0 through
v1.5.0 returns `["new" "stale"]` instead of replacing it with the documented
one-element result `["new"]` under both exact SDKs.
Upstream merged PR #266 three days after the selected tag and commit
`33d262eb9c0f5add618477081144108325cbc55c` adds the missing destination
slice truncation; main records it for unreleased v1.5.1. Later unreleased
commits also fix decode-hook error wrapping and a named-string
TextUnmarshaller panic. Those branch fixes are not eligible releases. No
canonical exact-path stable qualifies.

The graph contains nine exact target requests. Main and Viper v1.15.0 request
v1.5.0. Historical Viper v1.10.1 and crypt v0.4.0 request v1.4.3; historical
Kong `v0.2.1-0.20190708041108-0548c6b1afae` and Consul API v1.1.0/v1.12.0
request v1.1.2; Serf v0.8.2/v0.9.6 request pseudo-version
`v0.0.0-20160808181253-ca63d7c062ee`. Main's explicit request is indirect.
Current repository source has no target import. Viper v1.15.0 imports the
target in production and test source; project `cmd` imports Viper, so target
why is positive through `cmd -> Viper -> mapstructure`, and one target
package is production/complete-test loaded. The project calls Viper config
setup/read APIs but has no repository-local mapstructure Decode/Unmarshal
call. Historical Viper, Consul, and Serf versions genuinely import the
target; Kong and crypt are metadata-only requesters. Only current Viper is
loaded among target requesters; the historical requester why results are
negative.

A disposable exact-Go-1.26.7
`go get github.com/mitchellh/mapstructure@v1.5.0` is a byte no-op. The
unchanged project remains 234 selected modules, 3,599 graph edges, 355
production entries, 429 complete-test entries, 197 module-backed entries
across 41 loaded modules, and 1,067 sum lines. `go.mod` / `go.sum` SHA-256 is
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
The unchanged 432-line tidy projection retains SHA-256
`3309bc33637f6868a75b6f3006f333e547d8481645da22f238763297bbedc708`;
the common applied 52/948-line hashes remain
`5881324093819c0c824281bab7ee950a9eac9a888e9ae50377af386119872479` /
`b01164dfb62d3a2b7a049ee46d79ef1f48fd6e5a3527715729a1911e2b6dff8a`.
No lower projection was retained because the ordinary failure spans the full
line.

All 42 earlier guarded selections remain exact. Their 247 sorted incoming
edges retain SHA-256
`4caf9e803c1b324d2afcd17b9ba664e6d2c8d220cd74b264484107ea3dd93f42`.
Thirty-seven why results remain negative; only closed kr/pretty, kr/text,
emoji/v2, Promptui, and go-homedir are positive. Promptui and go-homedir are
the only guarded repository imports; emoji/v2, Promptui, and go-homedir are
the only loaded guarded modules. Every go-homedir, Promptui, emoji/v2,
kr/text, and kr/pty expiry guard remains exact.

Exact-version OSV and GitHub global queries are empty for all 17 releases,
and the repository advisory endpoint is empty. Pinned isolated govulncheck
v1.8.0 under exact Go 1.26.7 has no target module/package/symbol/test-symbol
result for any stable. Base and no-op-v1.5.0 project populations are
identical at 30 module, 22 package, 20 symbol, and 20 test-symbol OSV IDs,
with no target trace. Guard OSV retains only Gorilla WebSocket GO-2026-6278/
GHSA-w67g-5rqw-f597 and go-retryablehttp GO-2024-2947/
GHSA-v6v8-xj6m-xwqh; x/mod v0.14.0 retains GO-2026-6179 and GO-2026-6180.
The Go index remains 518,501 bytes/1,402 records at SHA-256
`bdd6fef3e1c488176122a98315a4c8c6874d09c6ad0d65ae69fa0fb44fa2c4cd`,
Last-Modified 2026-09-17T17:29:18Z. The PUBLISHED 2,807-byte memberlist CNA
response remains at SHA-256
`cacd85de49cc8685c4eb5a378d1825408bae27e152b0db60c6068e784082c659`.
Advisory absence does not override the ordinary failure or unsupported owner.

Final unchanged-project exact Go 1.26.7 module verification, build,
count-one tests, race count-one tests, and vet pass under `umask 022`. The Go
1.18, source/API/CLI/help/launcher/Make/quality contracts and every earlier
decision remain exact. Accepted quality remains 27/27 Q0-Q2 PASS at L2.
Every task-owned scratch artifact was contained beneath the managed session
scratch root and removed; only the pre-existing launcher-owned Node compile
cache remains. Treat all completed release, source, behavior, closure,
route, projection, advisory, and final-gate findings as final.

# Choose Exactly One Direction

1. Explicitly retain exact selected indirect, production-loaded
   `github.com/mitchellh/mapstructure v1.5.0` without product-source or
   dependency-metadata changes under a mapstructure-specific,
   non-transferable exception. Call it unqualified: no canonical exact-path
   stable qualifies and no exact-path stable has a genuine supported owner.
   Accept only the completed release/source, ordinary behavior, exact owner/
   request, route, why/import/load/runtime, graph/tidy/Go-floor, earlier-
   guard, and advisory facts. Define the exact expiry guards below. This
   exception must not broaden or expire any closed exception.
2. Authorize exactly one later bounded measurement-only genuine owner/request
   study named `Ply Viper Mapstructure Ownership Study`. Its sole scope is
   the existing main indirect v1.5.0 request, current `cmd -> Viper v1.15.0
   -> mapstructure v1.5.0` import/load route, historical requester provenance,
   and the owner-declared `go-viper/mapstructure` alternate-path fork. Its
   exact question is whether a genuine supported Go-1.18-compatible owner
   route can remove the unsupported exact-path request while preserving
   public API/CLI behavior, all 42 earlier selections/247 edges, and every
   closed exception. Do not run the study, change source or any selection,
   reopen Viper or an earlier decision, or introduce a fork/replacement in
   this decision-recording move; prepare one reciprocal measurement-only
   successor.
3. Stop P7 unresolved with no product-source, dependency-metadata, roadmap,
   or P8 implementation change and no successor execution.

Do not invent a fourth option or combine options. Physical selection,
indirectness, positive why, production loading, and advisory absence are not
qualification. Do not call v1.5.0 qualified, select a lower release, add or
change a root, use the unreleased branch, a fork, replacement, pseudo-version,
or alternate path, change Viper or any earlier guard, transfer an exception,
or begin P8.

# Option 1 Expiry Guards

If option 1 is selected, preserve and record at minimum:

- exact selected path/version `github.com/mitchellh/mapstructure v1.5.0`,
  main's indirect request, the other eight exact current/historical requests,
  all current requester selections, and every current/historical requester
  import or metadata-only boundary;
- zero repository target imports, the current `cmd -> Viper -> target` why
  route, Viper production/test imports and project config calls, one
  production/complete-test target package, and the exact load/runtime-use
  boundary;
- the exact 17-release line, absent exact-path `/v2` and `/v3`, archived
  repository/owner/status/license/default-branch identity, blessed alternate-
  path fork identity, release/tag/commit/tree/signature/archive/sumdb/module
  facts, the v1.2.0 tag-mutation record, and no replacement/retraction/
  deprecation or eligible exact-path alternate;
- all 17 stables remaining unqualified for the completed ordinary stale-slice
  result, the v1.2.2 repeated-suite result, no future qualified Go-1.18-
  compatible exact-path stable, and no genuine supported exact-path owner;
- baseline 234/3,599/355/429/197/41/1,067 state, exact module hashes, common
  tidy state, no-op v1.5.0 projection, and no dependency implementation;
- exact Go 1.18 floor, exact Go 1.26.7 identity, source/API/CLI/help/launcher/
  Make/quality contracts, and accepted 27/27 Q0-Q2 PASS at L2;
- all 42 earlier selections/247 incoming edges and their exact hash, plus
  every go-homedir, Promptui, emoji/v2, kr/text, and kr/pty expiry guard; and
- no new target/closure advisory, independent defect, release, owner,
  qualified stable, supported exact-path owner, or compatible genuine route
  to a qualified target or closed-exception release.

Any target/request/requester-import/owner-route, root/why/import/load/runtime,
graph/module/tidy/Go-floor, source/behavior/closure, repository/release/owner,
advisory/finding, qualification, supported-owner, earlier guard, closed-
exception, or compatible-route change expires option 1 and requires a fresh
mapstructure dependency and product decision before merge. Any closed-
exception change separately requires its corresponding fresh decision. The
exception is not a security claim and cannot transfer.

# Decision Recording Contract

This successor records a product direction only. Reproduce the handoff,
chain, launcher, exact target/request/route, module/graph/tidy, closed-
exception, guard, owner/release, and fresh narrow advisory premises without
repeating the completed source, behavior, closure, upstream, cross-build,
projection, archive, or govulncheck evaluation. Run final exact-Go unchanged-
project module verification, build, count-one tests, race count-one tests,
and vet.

If option 1 is chosen, update the roadmap and rolling handover, answer this
archive, and prepare exactly one reciprocal evaluation of the next
unanswered selected queue item without executing it. If option 2 is chosen,
prepare only the named measurement-only successor. If option 3 is chosen,
record the stop without preparing a dependency evaluation or entering P8.
Make only the required local handoff documentation commit; do not create a
dependency implementation commit.

# Required Reading

Read this archive, the answered mapstructure evaluation, answered go-homedir,
Promptui, emoji/v2, kr/text, and kr/pty decisions/evaluations, kr/pretty and
Cast owner records, rolling handover, P7/P8 roadmap, `go.mod`, and `go.sum`.
Verify branch, ancestry, clean ordinary/ignored state, reciprocal archive
chain, launcher check, exact Go identities, module hashes/counts/tidy
projection, all nine target requests, requester imports and routes, why/load/
runtime state, all 42 earlier guards and the 247-edge snapshot, all five
closed exception boundaries, and fresh advisory identities before recording
a choice.

# Three Moves

First, verify every completed premise without reopening the evaluation.
Second, choose and record exactly one option, make only its authorized
documentation/handoff changes, and run the unchanged-project final gate.
Third, verify scratch cleanup, reciprocal chain, single NEXT state, launcher
prompt mirror, exact changed set, and clean status; commit the local handoff
without executing its successor.

# Automatic Handoff

Do not launch a successor, push, merge, publish, release, stash, revert,
bypass cleanup, remove the worktree, transfer an exception, reopen
go-homedir, Promptui, emoji/v2, kr/text, kr/pty, kr/pretty, Cast, Viper, or an
earlier decision, repeat the mapstructure evaluation, evaluate another
dependency group, write outside the managed scratch root, or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
