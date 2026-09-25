# Agent Session: Upgrade Pflag Dependency

Status: ANSWERED - HISTORY
Session ID: `2026-09-03T233249+0200-upgrade-pflag-dependency`
Created: `2026-09-03T23:32:49+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `75e9108b605624f468394d3041fb6db36a8a4dec8b11ed7077298d23286b996f`
Previous: [2026-09-03T212310+0200-upgrade-cobra-dependency.md](2026-09-03T212310+0200-upgrade-cobra-dependency.md)
Next: [2026-09-04T013713+0200-upgrade-viper-closure.md](2026-09-04T013713+0200-upgrade-viper-closure.md)
Outcome: `github.com/spf13/pflag` moved from v1.0.9 to v1.0.10 as the sole selected-module change at `360b2f3`; fresh primary release, proxy, tag, checksum, graph, behavior, and vulnerability evidence confirmed 233 modules, 3,551 edges, the exact 429-package population, and 22/33/22 vulnerability parity. Exact `make quality` exits 0 with all 27 Q0-Q2 rows passing at clean L2, and the independent 40-stage regression passes with no production, API, CLI, toolchain, artifact, quality, baseline, or unrelated dependency drift.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its next smallest dependency group: upgrade existing indirect
`github.com/spf13/pflag v1.0.9` to latest `v1.0.10` as a one-module move.
Independently reverify the release decision from primary evidence, preserve
behavior and every quality contract, and finish with no unrelated drift.

# Authorized Roadmap

P2A-P6 are complete. The Go 1.26.7 toolchain, terminal pair, logrus group, and
Cobra closure are complete; P7 remains active for dependency groups and P8
remains queued. This session may change only pflag's required `go.mod`/`go.sum`
metadata, a focused executable contract if one is needed, and the roadmap
record.

The external probe selected v1.0.10 as the latest pflag release. It declares
Go 1.12 and changes only the existing pflag selection: 233 selected modules,
3,551 graph edges, and the exact 429-package test population are retained. The
historical `go mod tidy -diff` projection grows from 238 to 240 lines only for
the two v1.0.10 checksums. Do not silently widen this group; stop and record the
decision if independent replay changes any other selected module or checksum.

Viper is deliberately not combined with this move. The retained wider probe
found that even its smallest next release, v1.16.0, changes 31 selected modules,
adds `github.com/google/s2a-go`, and raises selected `testify` to a module that
declares Go 1.20. That requires a separate scope and compatibility decision.

# Measurements At Start

The Cobra implementation commit is
`e14ed5ecd0856893c282f37995077542845fb563`, exact parent
`b10076effe353edfae1020e0585175332cc2ea5d`, clean tree
`f17f7cf92e3c52be38f1a5180e283023a76aa75a`. After the handoff, the continuity
HEAD must have exact parent `e14ed5e`, and ordinary and ignored status must be
empty.

The declared and verified toolchain remains Go 1.26.7. The retained official
executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; audit
reproduction invokes literal `go`. Do not inject ambient `GOFLAGS`.

The independently replayed Cobra selection root is
`/private/tmp/ply-p7-cobra-selection.b10076e.94XcUY`, with verified
43,656-entry manifest SHA-256
`5e8a8cce3fe2495620987bce5bc9ebb98ac6db871aa4d7e14875a00f45f05e8f`.
The next-group pflag selection root is
`/private/tmp/ply-p7-next-pflag-selection.e14ed5e-final.IJI34Q`, with verified
38,164-entry manifest SHA-256
`16b28f1dd4fd329b7404d16cb00367a8610f9dd354b3bfed87c9ac4eebf0c5d4`.

The exact Cobra quality root is
`/private/tmp/ply-p7-cobra-quality-gate.e14ed5e.1788465147.56571`, with verified
234,708-entry manifest SHA-256
`a59de1e99bdb27b1bcd1e364f2a5feb1d9ba7c1dfec1a1fc381fd93327749e71`.
Its Q0-Q2 scorecard SHA-256 is
`3756e807d0fb56aee5f49b4a43d60e80526c6ab7f1b8456b388cccedfb20b9b8`:
all 27 criteria pass at L2, populations are 8/8 and 4/4, and held, regressed,
not-comparable, and dirty counts are zero.

The independent Cobra regression root is
`/private/tmp/ply-p7-cobra-regression-gate.e14ed5e.H7N4M5`, with verified
97,882-entry manifest SHA-256
`ec18888b35837e6353825cecec7edeb4bd1daa25a7cc428c4ace567130fb80f1`.
All 40 stages pass. Its full report exits 1 only for queued L3 rows Q3.1,
Q3.3, Q3.4, and Q3.7; it is not a P7 dependency-group exit gate.

# Role And Boundaries

Verify time-sensitive pflag release metadata, tag identity, module `go`
requirement, checksums, release changes, and relevant vulnerability information
from primary Go module, repository, and Go vulnerability sources. Record why
v1.0.10 is compatible with pinned Go 1.26.7, golangci-lint 2.12.2, GoReleaser
2.17.1, Cobra v1.10.1, and the retained `go 1.18` language floor.

Do not change production Go behavior, public Go API, CLI semantics, toolchain
declarations, Docker or release inputs, any dependency other than pflag,
quality-tool versions, quality thresholds, baseline numeric debt,
compatibility allowlists, acceptance or mutation populations, publishers,
registries, credentials, inactive packaging, or P8 domain code. Do not combine
Viper, a source fix, or another dependency with this move; stop on an
incompatible behavior change.

Keep module/build caches, graph projections, vulnerability results, reports,
generated artifacts, build contexts, schema-2 evidence, and audit output
outside the worktree. Warm fresh module caches from external Git archives;
never use a real-tree `go mod download all` as a tidy surrogate. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint block before
editing. Verify the accepted manifests. Read the rolling handover, this
archive, the P7 roadmap, `go.mod`/`go.sum`, pflag/Cobra callers and tests,
toolchain and baseline-reproduction contracts, compatibility, snapshot/Docker,
quality, and audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, `go mod tidy -diff` projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, and a
vulnerability baseline. Reverify v1.0.10 against the current Go proxy and
repository release/tag evidence. Compare old/new graphs exactly and require the
sole changed selection and new checksum pair to be pflag v1.0.10.

Make one focused dependency-only implementation commit using the exact Go
tool, without hand-editing dependency metadata. Re-run focused pflag/Cobra
callers and tests, pinned lint, complete tests/race/vet, API/CLI and
entry/subprocess compatibility, launcher and Make contracts, complete
preflight, host acceptance, fresh snapshot/Docker meta and acceptance, audit
meta, focused and exact Q0-Q2 audits, the separate full audit, vulnerability
comparison, and empty-HOME count-2. Refresh external schema-2 evidence when
commit binding requires it.

Require exact `make quality` exit 0 at L2 with all 27 Q0-Q2 rows present,
8/8 mutation and 4/4 acceptance populations, and zero held, regressed,
not-comparable, or dirty counts. The full report may exit 1 only for the four
queued L3 rows, never 2. Retain immutable external evidence and verified
manifests for graph selection, quality, regression, and artifacts.

# Automatic Handoff

After this group succeeds, keep P7 active and P8 queued. Rewrite the rolling
handover and roadmap, answer this archive, create exactly one reciprocal NEXT
archive for the next evidence-selected dependency group, replace only the
launcher's mutable regions, run launcher and handoff contracts, and make the
normal `docs: prepare next agent session` commit. Do not implement the next
group in this session.

Do not launch a successor, push, merge, publish, release, delete retained
evidence or local images, stash, revert, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
