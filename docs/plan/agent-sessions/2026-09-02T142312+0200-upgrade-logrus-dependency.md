# Agent Session: Upgrade Logrus Dependency

Status: NEXT
Session ID: `2026-09-02T142312+0200-upgrade-logrus-dependency`
Created: `2026-09-02T14:23:12+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `aca406817d3a846547d13d63a91def327089ed97b586fcb8b93fa38dcb43ddaf`
Previous: [2026-09-01T154646+0200-upgrade-terminal-dependencies.md](2026-09-01T154646+0200-upgrade-terminal-dependencies.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its next smallest dependency group: upgrade only direct
`github.com/sirupsen/logrus v1.9.0` to `v1.9.3`. Independently reverify the
release decision from primary evidence, preserve behavior and every quality
contract, and finish only with no unrelated module, API, CLI, acceptance,
artifact, vulnerability, or baseline drift.

# Authorized Roadmap

P2A-P6 are complete. The Go 1.26.7 toolchain move and terminal-pair group are
complete; P7 remains active for dependency groups and P8 remains queued. This
session may change only logrus's required `go.mod`/`go.sum` metadata, a focused
executable contract if one is needed, and the roadmap record.

The external probe selected v1.9.3 because it is the highest logrus release
compatible with the retained `go 1.18` language floor. Current v1.10.x releases
require Go 1.23. The v1.9.3 probe changed only logrus: no coupled direct or
indirect version moved. Stop and record the decision if independent replay
requires a wider minimum-version-selection closure.

# Measurements At Start

The terminal-pair implementation commit is
`0f93a527d1f95cc043a695101746d359c74c9c54`, exact parent
`e25a98546c20f8707ffe4f71b1d544fc8fb4944f`, clean tree
`ebbc91e08907b37175faf45ddc3283e8bd0ed771`. After the handoff, the continuity
HEAD must have exact parent `0f93a52`, and ordinary and ignored status must be
empty.

The declared and verified toolchain remains Go 1.26.7. The retained official
executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; audit
reproduction invokes literal `go`. Do not inject ambient `GOFLAGS`.

The exact terminal selection root is
`/private/tmp/ply-p7-terminal-selection.ctFhJk`, with verified 76,252-entry
manifest SHA-256
`c5cb5e2706ecccfc8eea7c80cadaa747c458e605b2919ab80cd378b263f84aca`.
The next-group selection root is
`/private/tmp/ply-p7-next-group-selection.0f93a52.Nl74oE`, with verified
31,623-entry manifest SHA-256
`e91b0260398550ed95d15063dbfbb96e01c2467b9f878ea70a0e459f25ec9b01`.

The exact quality root is
`/private/tmp/ply-p7-terminal-quality-gate.0f93a52.V0TBoc`, with verified
268,572-entry manifest SHA-256
`e8027bcb9e2a30ff2fbbf3e3bb8c1a1072abe98a7f55ff6a1fb313c4200d91df`.
Its Q0-Q2 scorecard SHA-256 is
`1ec4203b099e49017f67a906f19a2ceeb69f1e86b7ca33bd6f8cd1bf23a27cc5`:
all 27 criteria pass at L2, populations are 8/8 and 4/4, and held, regressed,
not-comparable, and dirty counts are zero.

The independent regression root is
`/private/tmp/ply-p7-terminal-regression-gate.0f93a52-final.hIYQFq`, with
verified 68,961-entry manifest SHA-256
`1a4414491a83262c4e7684dcc50d7c8372310a66d41e67cf0b9b9b30aca4964a`.
All 40 stages pass. Its full report exits 1 only for queued L3 rows Q3.1,
Q3.3, Q3.4, and Q3.7; it is not a P7 dependency-group exit gate.

# Role And Boundaries

Verify time-sensitive logrus releases, module `go` requirements, release
metadata, and relevant vulnerability information from primary Go module,
repository, and Go vulnerability sources. Record why the selected version is
compatible with pinned Go 1.26.7, golangci-lint 2.12.2, GoReleaser 2.17.1, and
the retained `go 1.18` language floor.

Do not change production Go behavior, public Go API, CLI semantics, toolchain
declarations, Docker or release inputs, another dependency, quality-tool
versions, quality thresholds, baseline numeric debt, compatibility allowlists,
acceptance or mutation populations, publishers, registries, credentials,
inactive packaging, or P8 domain code. Do not combine a source fix with this
dependency move; stop on an incompatible behavior change.

Keep module/build caches, graph projections, vulnerability results, reports,
generated artifacts, build contexts, schema-2 evidence, and audit output
outside the worktree. Warm fresh module caches from external Git archives;
never use a real-tree `go mod download all` as a tidy surrogate. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint block before
editing. Verify the accepted manifests. Read the rolling handover, this
archive, the P7 roadmap, `go.mod`/`go.sum`, logrus callers and tests, toolchain
and baseline-reproduction contracts, compatibility, snapshot/Docker, quality,
and audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, `go mod tidy -diff` projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, and a
vulnerability baseline. Probe v1.9.3 in an external copy first. Compare old/new
graphs exactly and require every changed module and checksum to be logrus.

Make one focused dependency-only implementation commit using the exact Go
tool, without hand-editing dependency metadata. Re-run focused logrus callers
and tests, pinned lint, complete tests/race/vet, API/CLI and entry/subprocess
compatibility, launcher and Make contracts, complete preflight, host
acceptance, fresh snapshot/Docker meta and acceptance, audit meta, focused and
exact Q0-Q2 audits, the separate full audit, vulnerability comparison, and
empty-HOME count-2. Refresh external schema-2 evidence when commit binding
requires it.

Require exact `make quality` exit 0 at L2 with all 27 Q0-Q2 rows present,
8/8 mutation and 4/4 acceptance populations, and zero held, regressed,
not-comparable, or dirty counts. The full report may exit 1 only for the four
queued L3 rows, never 2. Retain immutable external evidence and verified
manifests for graph selection, quality, regression, and artifacts.

# Automatic Handoff

After this group succeeds, keep P7 active and P8 queued. Rewrite the rolling
handover and roadmap, answer this archive, create exactly one reciprocal NEXT
archive for the next evidence-selected small dependency group, replace only
the launcher's mutable regions, run launcher and handoff contracts, and make
the normal `docs: prepare next agent session` commit. Do not implement the next
group in this session.

Do not launch a successor, push, merge, publish, release, delete retained
evidence or local images, stash, revert, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
