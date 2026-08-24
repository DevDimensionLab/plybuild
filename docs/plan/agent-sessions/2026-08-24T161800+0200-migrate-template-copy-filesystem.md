# Agent Session: Migrate Template Copy Filesystem

Status: NEXT
Session ID: `2026-08-24T161800+0200-migrate-template-copy-filesystem`
Created: `2026-08-24T16:18:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `18f2a2a2f5fb20952c574df30d1edaa01e23e6e719606bea56b177ed45bd0133`
Previous: [2026-08-24T151703+0200-migrate-maven-metadata-http.md](2026-08-24T151703+0200-migrate-maven-metadata-http.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the next P3 move in one measured change: introduce the declared
`internal/adapter/filesystem` boundary through the coherent non-existing-target
template copy flow and kill the `template-copy` source/destination argument
swap. Keep Q1.2 at zero, reduce only migrated Q1.3 sites, and preserve every
P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, `pkg/template/template.go` and its tests/callers,
`pkg/file/file.go` copy/merge helpers and every caller/test, relevant project and
template types, the existing process/HTTP/cloud seam patterns, the
Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 scanners, and the P2A compatibility, subprocess, and
host acceptance contracts. Regenerate ignored reports outside the measured
tree or remove them before the clean audit. Implementation commit `ffc4e77`
has 65 tests across 13 of 24 packages, Q0.6 has zero unsafe direct test writes,
Q1.2 is 0, Q1.3 is 74 violations of 79 production effect sites with three
adapter paths absent, and Q1.4 is 5 of 8. Exact Q2.1 remains 0 of 8 executable
harnesses; the upstream filename-only denominator sees three non-executable
`mutate-*` paths. The clean gate passed, the full audit exited 1 for 16
documented findings and never 2, and comparable ratchets were five improved,
two held, and zero regressed.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the template-copy filesystem move. Do not push,
merge, publish, remove the worktree, stash inherited changes, revert user work,
or run destructive Git commands. Do not invoke a publisher or distribution
command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Preserve exported
`template.MergeTemplate`, `file.CopyOrMerge`, and `file.CopyFile` signatures and
every caller. Production must retain exact template filtering and target-path
derivation, existing-target merge selection, non-existing-target copy behavior,
source bytes and mode, missing destination-directory creation, logging, error
order and semantics, search/replace, render, and Maven merge follow-ups.

Migrate only the filesystem effects needed by the coherent non-existing-target
template copy flow. Use resolvable interfaces and complete dependency values,
not stored function dependencies. Do not create a second generic filesystem
implementation or broaden into existing-target merge internals, unrelated
write/delete/rename/grep/render calls, Git hooks, cloud, Spring/download, HTTP,
clock, server, Docker, dependencies, distribution, or formal mutation-harness
scope.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`,
`pkg/template/template.go`, its tests and callers, `pkg/file/file.go` copy/merge
helpers and every caller/test, relevant project/template types, the existing
process/HTTP/cloud interfaces and recording tests, the
Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 implementation in `.quality/tools`, and the P2A
compatibility and host acceptance contracts before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with recording template copy and filesystem contracts. Prove the
   complete source path stays before the complete resolved target path, a
   missing target selects copy while an existing target retains merge, source
   bytes and mode reach the destination, a missing destination directory is
   created, dependency errors retain their order and behavior, safe defaults
   perform no mutation, the complete dependency reaches the flow, and an empty
   recorded population fails.

2. Introduce the thinnest zero-value-safe `internal/adapter/filesystem`
   interface boundary and route only the migrated copy effects through it.
   Preserve source read, destination-directory existence/creation, source mode
   lookup, logging, destination write, and error ordering without duplicating
   `CopyOrMerge` or merge behavior. Bind the immutable `template-copy` label
   through a new non-executable `scripts/mutate-template` seam driver and its
   Q0.8 meta-test, without claiming the later P5 harness.

3. Run focused template/file/adapter/argument-swap tests,
   Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 measurements, API/CLI compatibility, all four host
   acceptance flows, the full checkpoint gate, and empty-HOME count-2 from a
   clean commit. Expect Q1.2 to stay zero, `template-copy` to raise Q1.4 from 5
   of 8 to 6 of 8, and exact Q2.1 to remain 0 of 8. Accept only regenerated
   Q1.3 values. The full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent template-copy
filesystem move or record an exact resumable state. Rewrite the rolling
handover, update the P3 measurements, answer this archive, create one linked
NEXT archive, replace the launcher's mutable regions, run the launcher
contract, and make the separate handoff-only commit
`docs: prepare next agent session`.

Keep P3 active for later filesystem adapter moves. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
