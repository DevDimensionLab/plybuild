# Agent Session: Remove Process Exits

Status: NEXT
Session ID: `2026-08-24T113208+0200-remove-process-exits`
Created: `2026-08-24T11:32:08+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `31c6f9b006cf6efc0102f7370b2e9cf39beb242099a567631466e120bbfca86f`
Previous: [2026-08-24T110354+0200-make-distribution-local.md](2026-08-24T110354+0200-make-distribution-local.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the first P3 move in one measured change: remove process termination
from non-main code, make both executable entry points delegate to one
error-returning command path, and reduce Q1.2 from 127 process-exiting calls
outside `main` to zero. Preserve the P2A public API, CLI, subprocess, and host
acceptance contracts with zero comparable ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P2B is complete, P3 is
active, and P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md. A queued checkpoint is approved work, so the
launcher must remain NEXT until all authorized checkpoints are complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
both `main.go` entry points, `cmd.Execute()`, `cmd.RootCmd`, process-exit call
sites, the Q1.2 scanner, and the API/CLI/subprocess contracts. Regenerate
quality reports rather than relying on ignored output. The clean P2B audit at
312d168 exited 1 for 16 documented findings with two improved, five held, zero
regressed ratchets, and 127 process-exiting calls outside `main`. Exit 2
invalidates a checkpoint.

# Role And Boundaries

Work autonomously in this worktree on codex/upgrade-quality. Make one focused
implementation commit for this P3 move. Do not push, merge, publish, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands. Do not execute any publisher; `make release` is now a local snapshot
alias, but distribution execution is outside this move.

Keep Go 1.18 and preserve exported `cmd.Execute()` and `cmd.RootCmd` symbols and
signatures, the normalized Cobra tree, subprocess output/exit behavior, and all
four host acceptance flows. If evidence shows external callers require
`cmd.Execute()` itself to terminate the process, stop for an explicit migration
decision instead of hiding an exit behind a function variable or adapter.

Do not introduce subprocess, HTTP, or filesystem adapters until the process
boundary is green and measured. Do not broaden into Docker, dependencies,
cloud, Spring, distribution, or mutation harnesses.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`, both executable
entry points, the command execution path, and the P2A API/CLI/subprocess and
host acceptance contracts before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start with a failing entry-boundary contract. It must cover both `main.go`
   and `cmd/ply/main.go`, prove they delegate to the same error-returning command
   path, and reject process termination outside the two `main` functions.

2. Add the error-returning execution path, keep the exported compatibility
   surface intact, move exit-code selection and termination to the executable
   boundaries, and preserve root/status/upgrade/build plus unknown-command
   stdout, stderr, and exit behavior. Do not begin adapter work in this move.

3. Run the focused boundary tests, API/CLI compatibility, host acceptance, the
   full checkpoint gate, and the empty-HOME count-2 test from a clean commit.
   The full audit may exit 1 for known findings, but Q1.2 must reach zero and no
   comparable ratchet may regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent P3 move or record
an exact resumable state. Rewrite the rolling handover, update checkpoint
status, answer this archive, create one linked NEXT archive, and replace the
launcher's mutable session regions. No separate agent-restart message is
required.

Keep P3 active for the ordered adapter moves after the exit boundary is green.
Stage only handoff files in the handoff commit named
`docs: prepare next agent session`. Do not launch the next session yourself.
COMPLETE is valid only after every authorized checkpoint through P8 is complete.
<!-- CODEX_SESSION_PROMPT_END -->
