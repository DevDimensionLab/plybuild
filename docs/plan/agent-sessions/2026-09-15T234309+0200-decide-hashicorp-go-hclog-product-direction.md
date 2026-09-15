# Agent Session: Decide Hashicorp Go Hclog Product Direction

Status: NEXT
Session ID: `2026-09-15T234309+0200-decide-hashicorp-go-hclog-product-direction`
Created: `2026-09-15T23:43:09+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3258dcf1534f7509a5ff1c22c1e7d942fd2a1849f1ba75e5d5f373b47357ff02`
Previous: [2026-09-15T222945+0200-evaluate-hashicorp-go-hclog-dependency.md](2026-09-15T222945+0200-evaluate-hashicorp-go-hclog-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by obtaining and recording one bounded product decision for
exact-path `github.com/hashicorp/go-hclog`. Reuse the completed evaluation: no
exact-path stable release satisfies the existing behavior and compatibility
contracts. Present the three authorized options below, wait for the user to
select one, then record exactly that choice. Do not repeat the audit, implement
a dependency change, evaluate another group, or begin P8.

# Authorized Roadmap

P2A-P6 are complete. P7 is blocked only on this go-hclog product decision after
exact Go 1.26.7, every accepted dependency move through Google UUID v1.4.0,
qualified exact go-cleanhttp v0.5.2, and all retained-module decisions through
Errwrap v1.0.0. Every earlier result and lifecycle ancestor is final. This
session may record one explicit user choice and prepare one bounded follow-up;
it may not implement that choice or combine another dependency group. P8
remains queued.

The completed evaluation left `go.mod` and `go.sum` unchanged. Exact selected
`github.com/hashicorp/go-hclog v1.2.0` remains inherited only because Viper
v1.15.0 requests it; this physical selection is not a qualification or risk
acceptance. Do not add a direct edge merely to alter MVS.

The qualified go-cleanhttp v0.5.2 result and the Errwrap v1.0.0, Consul SDK
v0.8.0, Consul API v1.18.0, Gateway v1.16.0, gRPC Prometheus v1.2.0, gRPC
middleware v1.0.0, Gorilla WebSocket v1.4.2, GopherJS
`v0.0.0-20181017120253-0766667cb4d1`, Enterprise Certificate Proxy v0.2.1,
and GAX v2.7.0 decisions remain separate under their recorded exact-selection,
incoming-edge, zero-load, runtime-unreachable, and no-new-finding guards.
Revalidate those guards before recording a choice and stop for the owning
decision if one expired. Do not transfer an earlier exception to go-hclog or
change Viper, either Gateway parent, mvn-pom-mutator, GoConvey, historical
Consul API, go-multierror, Serf, or `sagikazarmark/crypt`.

# Measurements At Start

The completed evaluation began from clean handoff HEAD
`452fffbd865dcae811fba3b02821cf7fd48f774c`, parent
`0d1d66480604da5b646c75593570b01ebec72dda`, tree
`4858e0dd3c216b772476058324a265aef0cb005a`. The latest dependency
implementation remains exact Google UUID v1.4.0 commit
`cf53bc64eeb69471d35c7536d196bf1da15f3973`. The evaluation made no source
or dependency metadata change and prepared this decision-only handoff commit.

Base `go.mod` and `go.sum` SHA-256 values remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2. Verify the new handoff HEAD,
parent, tree, and clean status at start rather than assuming their values.

# Completed Evaluation

Fresh proxy, sumdb, `go-import`, Git, and GitHub evidence resolves exact module
path `github.com/hashicorp/go-hclog` to the public, active, unarchived, non-fork
MIT repository `https://github.com/hashicorp/go-hclog.git`. The proxy exposes
31 valid stable releases from v0.7.0 through v1.6.3. There are no retractions,
module deprecation, redirect, `/v2` module, or qualified alternate path.
Non-semver alias tags, nested `hclogvet` tags, branch `f-v2`, and unreleased
main were not promoted. Current main declares Go 1.25.

Selected v1.2.0 is lightweight unsigned tag/commit
`b6b55671f4e5b82443139ee3e9f4417603c4cd72`, parents
`81033451e6eb54da74139737d45a39b6412c7f34` and
`fc772a82149bb181261310935e431f89667288c3`, tree
`45f0da0a3b7eb001522fc9891504cee03d951ef0`, dated
2022-03-03T03:51:25Z. Latest v1.6.3 is commit
`d12136aa2e51933c460084f5083b6d5bb9d41960`, parents
`5dbb615f9aa8587fce14c2ab180aa7369f0ee703` and
`cb8687a9e2d8bab634ddff8412b3e03a7d60c068`, tree
`26c5c9247ad5b3a6f930ea2b1fc0a1ee532a7804`, dated
2024-04-01T20:03:54Z. Strict Git verification and proxy/Git byte comparison
pass. V1.6.3 is proxy `@latest`.

All valid releases declare no more than Go 1.13. Complete imported production
and test closures for selected v1.2.0 and latest v1.6.3 declare no more than Go
1.17 and execute under contained Go 1.18.10. Thus the full stable line is
floor-eligible; the blocker is behavior and API qualification.

Each serious v1 release has one root package, 12 production Go files, seven
tests, only Unix/Windows color build branches, and one benchmark. There are no
commands, examples, fuzz targets, testdata, generated files, cgo, embeds,
`go:generate` directives, or symlinks. Selected, v1.2.1, v1.2.2, v1.3.0,
v1.3.1, v1.4.0, and latest v1.6.3 verify, build, pass native count-one,
repeated count-ten, race, and vet under exact Go 1.26.7 and Go 1.18.10.
Selected/latest production and test cross-builds pass for Darwin AMD64, Linux
AMD64/ARM64, Windows AMD64, and js/wasm under both SDKs.

Pinned API comparison makes v1.3.1 the last API-compatible upgrade: v1.2.1 and
v1.2.2 are identical to selected, while v1.3.0/v1.3.1 add only
`LoggerOptions.ColorHeaderAndFields`. V1.4.0 adds `Logger.GetLevel`, an
incompatible method addition to the exported interface; v1.4.0-v1.6.3 retain
that break. Latest adds otherwise compatible options and `SupportsColor`.

Independent fixtures characterize deterministic plain/JSON output, sorting,
odd keys, stdlib adaptation, level routing, exact output-reset error identity,
caller-owned output cleanup, nil/panic boundaries, aliasing, allocations, and
supported immutable concurrent logging. Selected through v1.3.1 each fail the
same five release-blocking contracts under both SDKs, and source history shows
the same implementations through latest:

1. JSON logging silently drops the complete record for unsupported values such
   as NaN because only `json.UnsupportedTypeError` is recovered.
2. Caller fields named `@message` or `@level` overwrite core JSON metadata.
3. Deregistering an absent sink underflows the sink count and disables a later
   real sink.
4. Sink callbacks run while the registry mutex is held, so a sink that
   deregisters itself deadlocks.
5. `SetDefault(nil)` makes `Default` and `FromContext` return nil despite
   `FromContext` documenting a guaranteed non-nil logger.

Latest v1.6.3 also races when `SyncParentLevel` observes concurrent
`SetLevel`/`GetLevel`: the race detector reports writes and reads of `level`,
`setEpoch`, and `ownEpoch` under both SDKs. Other characterized boundaries are
not blockers: `New`, `NewNullLogger`, and `NewSinkAdapter` accept nil options;
`FromContext(nil)` and `FromStandardLogger` with nil options panic;
`With(nil)` drops the unmatched value; `ImpliedArgs` and leveled-writer maps
alias caller-visible state; reset preserves old output on flush error; and
callers own output closure.

No higher version fixes the five common defects. V1.3.1 therefore cannot be
selected despite API compatibility, and v1.4.0 or later adds an independent
public-interface break. Exact v1.1.0 or lower cannot replace selected without
changing the owning parent: an exact v1.1.0 get downgrades Viper v1.15.0 to
v1.10.1, while v0.9.2 additionally removes mvn-pom-mutator. Those parent moves
are outside this dependency group.

Selected v1.2.0 exists through exactly one selected-version edge from Viper
v1.15.0. Historical Viper v1.10.1 and `sagikazarmark/crypt v0.4.0` request
v1.0.0; historical Consul API v1.12.0 and Consul SDK v0.8.0 request v0.12.0.
The shortest path is main -> Viper v1.15.0 -> go-hclog v1.2.0. `go mod why -m`
is negative, repository Go imports are zero, and production and complete-test
loads contain zero target packages, so go-hclog is runtime-unreachable.

Disposable exact v1.2.0/v1.3.1/v1.4.0/v1.6.3 gets retain 234 modules, all
unrelated selections, 355 production entries, 429 complete-test entries, and
zero target load. They manufacture a direct indirect root and target checksum;
later choices move only the target. Tidy removes the manufactured root and
returns every projection to selected v1.2.0 and the same base projection. No
projection was applied.

Fresh primary vulnerability data has 1,399 records at SHA-256
`033d823ed9fd985c165ae9e65500cda18910e17cf35126866a12d94bbc5a805f`
and Last-Modified 2026-09-15T18:54:35Z, with no go-hclog record. Exact OSV
queries for v1.2.0, v1.3.1, v1.4.0, and v1.6.3 are empty. Selected v1.2.0 has
zero Go 1.26.7 module/package/symbol/test-symbol findings. V1.3.1 and later
carry module-only GO-2026-5024 through old `x/sys`; it affects
`x/sys/windows.NewNTUnicodeString`, but scans contain no vulnerable-package,
called-symbol, go-hclog, or reachable target trace. Go 1.18.10 scans report
only the old SDK and that later closure module population; no advisory is
assigned to go-hclog. The project retains its normalized 30 advisory IDs and
zero target SBOM package, symbol, test-symbol, or reachable trace.

The unchanged project remains 234 modules, 3,599 graph edges, 429 complete-test
entries, 197 module-backed entries across 41 loaded modules, 1,067 sum lines,
and a 432-line unapplied tidy projection. `go.mod` and `go.sum` SHA-256 remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
Applicable Go 1.26.7 and contained Go 1.18.10 project gates pass; Go 1.18 keeps
only the two accepted Darwin `pkg/shell` wording failures. Accepted quality
remains 27/27 Q0-Q2 PASS at L2. No dependency implementation or metadata
commit was created.

# Authorized Product Options

Present these options without silently choosing one:

1. **Retain exact v1.2.0 under a bounded exception (recommended).** Keep exact
   inherited, unloaded v1.2.0 and the sole selected-version Viper v1.15.0 edge
   without changing `go.mod` or `go.sum`. Accept only the five completed common
   behavior findings, characterized nil/panic, aliasing, allocation,
   concurrency, resource, API, closure-vulnerability, and related qualification
   findings. The exception is go-hclog-specific and valid only while exact
   v1.2.0 and that edge remain unchanged, zero target packages load, the module
   remains runtime-unreachable, and no new advisory or independent defect
   appears. Direct import/loading, runtime reachability, a target version or
   incoming-edge change, or a new advisory or independent defect requires a
   fresh go-hclog dependency and product decision before merge.
2. **Authorize a separate remediation study.** Keep this session no-change and
   prepare a new bounded mission to evaluate parent removal, an upstream patch,
   fork, wrapper, or architecture replacement. This option does not itself
   authorize implementation, a Viper change, direct target use, a Go-floor
   change, or unrelated module movement.
3. **Block P7.** Leave go-hclog unresolved and stop the upgrade before merge
   until an exact-path stable release independently satisfies the contracts.

Do not offer v1.3.1 as qualified, promote v1.4.0+, downgrade Viper, change or
remove a parent, add a direct go-hclog edge, patch or fork source, raise the Go
floor, alter another guarded dependency, or manufacture a dependency commit
without a separately explicit choice and scope.

# Role And Boundaries

This is a decision session, not a renewed audit or implementation. Ask for one
explicit option selection if none was provided with the session invocation.
Elapsed time is not approval. Once the user chooses, revalidate only exact
selection, incoming edge, negative why, zero target and guarded package loads,
project hashes, runtime unreachability, and current advisory state. Record the
choice without broadening it.

# Required Reading

At start verify the feature branch, clean ordinary and ignored status, current
ancestry, reciprocal archive history, latest Google UUID implementation
identity, P7/P8 state, all guarded invariants, and
`./codex-dev-start.sh --check`. Read this archive, the answered go-hclog
evaluation, rolling handover, roadmap, `go.mod`, `go.sum`, and referenced
lifecycle contracts. Reuse the completed audit.

# Three Moves

First, revalidate only the bounded decision guards and present the three
options. Second, wait for and record one explicit user selection, including
its accepted findings, limits, and expiry triggers; do not infer approval.
Third, answer this archive and prepare exactly one reciprocal NEXT mission
consistent with the choice. Do not execute the successor.

# Automatic Handoff

Only after an explicit choice, run applicable no-change lifecycle gates and
make the required local `docs: prepare next agent session` commit. Do not
implement an option, launch a successor, push, merge, publish, release, stash,
revert, bypass cleanup, remove the worktree, combine another dependency group,
or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
