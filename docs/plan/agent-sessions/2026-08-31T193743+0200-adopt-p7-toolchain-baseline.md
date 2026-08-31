# Agent Session: Adopt P7 Toolchain Baseline

Status: NEXT
Session ID: `2026-08-31T193743+0200-adopt-p7-toolchain-baseline`
Created: `2026-08-31T19:37:43+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `ee0464ab62a1f4e35f9b0b0d452f8a7c87ed24dc5a9822a7fb6e7d85923285f5`
Previous: [2026-08-31T113850+0200-complete-p6-quality-exit-gate.md](2026-08-31T113850+0200-complete-p6-quality-exit-gate.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Begin P7 with one bounded toolchain-baseline move: select the maintained Go
1.26 or 1.27 line supported by the current pinned quality tools, migrate every
contractual Go toolchain declaration together, and reproduce old/new quality
measurements with both instrument identities. Keep the move toolchain-only and
finish only when the declared toolchain matches the verified toolchain with no
unexplained debt, API, CLI, acceptance, or artifact drift.

# Authorized Roadmap

P2A-P6 are complete. P7 is active and P8 remains queued. This session may
implement only the first P7 maintained-toolchain baseline move; dependency
version upgrades and P8 domain work remain out of scope.

Inventory every active and contractual Go identity before choosing a version,
including `go.mod`, Docker build inputs, active and deactivated CI/release
inputs, documentation, quality baseline instrument metadata, and executable
contracts. Verify time-sensitive support from primary Go, golangci-lint,
GoReleaser, and relevant action/image sources. Choose Go 1.26 or 1.27 from
evidence rather than availability alone, and record why the other line was not
selected.

# Measurements At Start

The P6 implementation commit is
`097a9f15782c773e47a61d903f5d09f5c96408f0`, exact parent
`b8b803c807166b00ed2cd257789b2337b41b9eab`, tree
`8f9616f7c4efd90a4d346adb0af4156566e1914d`. After the P6 handoff, the new
continuity HEAD must have exact parent `097a9f1`, and ordinary and ignored
status must be empty.

The module currently declares Go 1.18 while the verified P6 gates ran with Go
1.26.2. The Docker builder is a floating `golang:alpine`; the deactivated
release workflow still names Go 1.15. Do not assume either deactivated input or
floating image is irrelevant until its archive/build contracts are traced.

The current exact quality gate is
`/private/tmp/ply-p6-quality-gate.097a9f1-final.lOV5NN`, with verified
15,849-entry manifest SHA-256
`b89ac6645f814d40d2444ddca3a808919afacf40c873b576a170bc0e5a4091a6`.
Its Q0-Q2 scorecard SHA-256 is
`bfe32a5e7eb90e9c7e9cc756a47f38bfefd68866e754cbb030653d409fce5322`:
all 27 criteria pass, L2 is attained, populations are 8/4, and held,
regressed, not-comparable, and dirty counts are zero.

The independent regression root is
`/private/tmp/ply-p6-regression-gate.097a9f1.YGeuN9`, with verified
15,458-entry manifest SHA-256
`85fc11883c3ecd9ce444ae23777ef34c36944273d03fbc1f869b17fa477c7de1`.
The full report exits 1 only for queued L3 rows Q3.1, Q3.3, Q3.4, and Q3.7; it
is not a P7 exit gate.

# Role And Boundaries

Change only the smallest coherent set of toolchain declarations, their focused
contracts, exact-toolchain baseline metadata when reproduction proves it, and
the roadmap record. Do not upgrade direct or indirect dependency versions,
change production Go behavior or API/CLI semantics, reactivate CI, alter
acceptance or mutation populations, loosen quality thresholds, or modify
publishers, registries, credentials, or release behavior.

Do not edit numeric baseline debt merely to make the new toolchain pass. First
reproduce the current baseline with the old verified Go identity, then run the
same instrument with the candidate identity. Preserve every comparable numeric
debt value. If the candidate changes measurement semantics, dependency
selection, public output, generated artifacts, or compatibility in a way that
cannot be explained and contract-tested inside this move, stop and report the
decision instead of broadening scope.

Keep downloaded toolchains, caches, module caches, reports, generated files,
build contexts, schema-2 evidence, and audit output outside the worktree. Do
not change a global tool installation merely to make the candidate available;
use an explicit external regular tool path and record its version and checksum.
Never create `.agent-task/current.md` or `.quality/manual-evidence.json`.

# Required Reading

Confirm branch, exact ancestry, clean and ignored status, reciprocal archive
links, launcher `--check`, and the authorized P7 checkpoint before editing.
Read the rolling handover, this archive, the P7 roadmap, toolchain declarations,
baseline identity/reproduction contracts, Docker/snapshot inputs, compatibility
contracts, and pinned-tool support policies.

# Three Moves

Probe both the current and candidate Go executables explicitly. Record
`go version`, `go env` identity fields, executable SHA-256, module/toolchain
selection, and the support evidence used for the decision. Run old/new
baseline reproduction from clean external state before changing the checked-in
baseline. Make one focused implementation commit before creating any
commit-bound current evidence.

After implementation, refresh the external schema-2 document for the exact
clean commit if the audit requires commit-bound manual receipts. Run the
focused six-receipt audit, then the exact complete `make quality` apparatus
with fresh snapshot and fresh local Docker acceptance. Require exact
`--only Q0.*,Q1.*,Q2.*` audit exit 0 at L2 with every criterion present and no
held material debt. The separate full L3 report may exit 1, never 2.

Then rerun toolchain declaration contracts, API/CLI and entry/subprocess
compatibility, pinned lint, complete tests/race/vet, launcher and Make
contracts, complete preflight, host acceptance, snapshot/Docker meta and
acceptance, audit meta, focused and full audits, and empty-HOME count-2. Retain
immutable external evidence and a verified manifest for the complete move.

# Automatic Handoff

After this bounded toolchain move succeeds, keep P7 active and P8 queued.
Rewrite the rolling handover and roadmap, answer this archive, create exactly
one reciprocal NEXT archive for the first small dependency group, replace only
the launcher's mutable regions, run launcher and handoff contracts, and make
the normal `docs: prepare next agent session` commit. Do not implement the
dependency group in the toolchain session.

Do not launch a successor, push, merge, publish, release, delete retained P6
evidence or local images, stash, revert, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
