# Agent Session: Upgrade Terminal Dependencies

Status: ANSWERED - HISTORY
Session ID: `2026-09-01T154646+0200-upgrade-terminal-dependencies`
Created: `2026-09-01T15:46:46+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `3b941599d6de87621bcdfe4551818606034b787d49ed11feaf8a7ddff53d94ef`
Previous: [2026-08-31T193743+0200-adopt-p7-toolchain-baseline.md](2026-08-31T193743+0200-adopt-p7-toolchain-baseline.md)
Next: [2026-09-02T142312+0200-upgrade-logrus-dependency.md](2026-09-02T142312+0200-upgrade-logrus-dependency.md)
Outcome: `golang.org/x/term` moved from v0.5.0 to v0.29.0 and its sole MVS coupling `golang.org/x/sys` moved from v0.5.0 to v0.30.0 at `0f93a52`; the retained Go 1.18 floor selects the highest compatible terminal pair, exact Q0-Q2 exits 0 at clean L2, and the complete independent regression has no unrelated module, API, CLI, acceptance, artifact, or reachable-vulnerability drift.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 with its first small dependency group: upgrade only direct
`golang.org/x/term v0.5.0` and the coupled `golang.org/x/sys v0.5.0` selection
required by it. Choose current compatible versions from primary evidence,
preserve behavior and all quality contracts, and finish only with no unrelated
module, API, CLI, acceptance, artifact, or baseline drift.

# Authorized Roadmap

P2A-P6 are complete. The bounded Go 1.26.7 toolchain baseline move is complete;
P7 remains active for dependency groups and P8 remains queued. This session may
change only the terminal pair's required `go.mod`/`go.sum` dependency metadata,
their focused executable contract if one is needed, and the roadmap record.

Do not update another direct or indirect module merely because `go get` or
`go mod tidy` proposes it. First determine the terminal pair's exact minimum
version-selection closure. If it requires a wider group, stop and record the
decision instead of broadening scope.

# Measurements At Start

The P7 toolchain implementation commit is
`16ecb67eebb6450226d3264f30a64a14391b5245`, exact parent
`9852ed15b6e22ec0263bc4f81a973f4e088d7cf5`, clean tree
`94262fbfd426692fdc44f9eb21f85065ed40a6e2`. After the handoff, the continuity
HEAD must have exact parent `16ecb67`, and ordinary and ignored status must be
empty.

The declared and verified toolchain is Go 1.26.7. The retained official
executable is
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in `PATH` as well as passing exact `GO`; audit
reproduction invokes literal `go`. Do not inject ambient `GOFLAGS`.

The exact quality root is
`/private/tmp/ply-p7-quality-gate.16ecb67-final.hdvr9q`, with verified
246,560-entry manifest SHA-256
`15136fa28368eb4ada0b81dfddc4f36daa19f17b7ab00d3a4ea85d496066cb8e`.
Its Q0-Q2 scorecard SHA-256 is
`0cf6f16797ec73ed8d51944a77dc0cde4de85e174526b915afb58542e98d4cec`:
all 27 criteria pass at L2, populations are 8/8 and 4/4, and held, regressed,
not-comparable, and dirty counts are zero.

The independent regression root is
`/private/tmp/ply-p7-regression-gate.16ecb67-final.4nMFvg`, with verified
67,726-entry manifest SHA-256
`48be0c3b828c60702ff4f06c513b8720dba4259c9df39d99cf23b0f655dffded`.
All 36 stages pass. Its full report exits 1 only for queued L3 rows Q3.1,
Q3.3, Q3.4, and Q3.7; it is not a P7 dependency-group exit gate.

# Role And Boundaries

Verify time-sensitive terminal-pair releases, module `go` requirements,
release notes, and relevant vulnerability information from primary Go module,
repository, and Go vulnerability sources. Record why the selected versions are
compatible with the pinned Go 1.26.7, golangci-lint 2.12.2, GoReleaser 2.17.1,
and the retained `go 1.18` language floor.

Do not change production Go behavior, public Go API, CLI semantics, toolchain
declarations, Docker or release inputs, quality-tool versions, quality
thresholds, baseline numeric debt, compatibility allowlists, acceptance or
mutation populations, publishers, registries, credentials, inactive packaging,
or P8 domain code. Do not combine a source fix with this dependency move; stop
on an incompatible behavior change.

Keep module/build caches, graph projections, vulnerability results, reports,
generated artifacts, build contexts, schema-2 evidence, and audit output
outside the worktree. Warm a fresh module cache from an external Git archive;
running `go mod download all` in the real tree adds 208 historical checksum
lines and is not an authorized tidy. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint block before
editing. Verify the accepted manifests. Read the rolling handover, this
archive, the P7 roadmap, `go.mod`/`go.sum`, terminal-module callers and tests,
toolchain declaration and baseline-reproduction contracts, compatibility,
snapshot/Docker, quality, and audit contracts.

# Three Moves

From clean external state, record the old module graph, selected versions,
package population, checksums, `go mod tidy -diff` projection, build/test/lint
result, public help, API/CLI reports, generated artifact contracts, and a
vulnerability baseline. Probe the candidate pair in an external copy first.
Compare old/new graphs exactly and require every changed module and checksum to
belong to the authorized terminal pair closure.

Make one focused dependency-only implementation commit. Do not hand-edit
dependency metadata to force the desired graph. Re-run the focused terminal
callers/tests, pinned lint, complete tests/race/vet, API/CLI and
entry/subprocess compatibility, launcher and Make contracts, complete
preflight, host acceptance, fresh snapshot/Docker meta and acceptance, audit
meta, focused and exact Q0-Q2 audits, the separate full audit, vulnerability
comparison, and empty-HOME count-2. Refresh external schema-2 evidence only if
the exact commit binding requires it.

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
