# Agent Session: Characterize Compatibility

Status: ANSWERED - HISTORY
Session ID: `2026-08-24T103558+0200-characterize-compatibility`
Created: `2026-08-24T10:35:58+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8a9b5bab9598d2b76683a3d55d3415bc4ef62c13591a5156ccabf9c6c46bc971`
Previous: [2026-08-24T085458+0200-make-manual-evidence-reachable.md](2026-08-24T085458+0200-make-manual-evidence-reachable.md)
Next: [2026-08-24T110354+0200-make-distribution-local.md](2026-08-24T110354+0200-make-distribution-local.md)
Outcome: P2A completed at 14764fd with pinned API and normalized Cobra compatibility gates, four fresh-host acceptance verifiers and 26 meta-controls, a clean full gate, and zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete checkpoint P2A in exactly three measured moves: pin and gate the
public Go API against v1.0.1, capture an order-independent Cobra CLI contract,
and add four falsifiable host-binary acceptance verifiers. Preserve public API
and CLI behavior and leave zero comparable ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P2A is active and
P2B-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md. A queued checkpoint is approved work, so the
launcher must remain NEXT until all authorized checkpoints are complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P2A
plan, the continuity design, and the current quality report. Regenerate reports
rather than relying on ignored output. The full audit may exit 1 for documented
findings; exit 2 invalidates the checkpoint.

# Role And Boundaries

Work autonomously in this worktree on codex/upgrade-quality. Make one focused
commit per measured quality move. Do not push, merge, release, publish, remove
the worktree, stash inherited changes, revert user work, or run destructive Git
commands. Never run make release or make release-brew. Stop if compatibility
cannot be established, the audit exits 2, or a comparable ratchet regresses.

Keep Go 1.18 for the module. Pin external compatibility tooling without adding
its modern toolchain requirements to go.mod. Use temporary writable caches
under /private/tmp and keep tests independent of developer HOME, public
services, and repository writes.

# Required Reading

Read docs/plan/quality-handover.md, the P2A section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, .quality/inventory, and the existing
Makefile/script contract tests before changing behavior.

# Three Moves

1. Pin golang.org/x/exp/cmd/apidiff at
   v0.0.0-20260709172345-9ea1abe57597 outside the project module. Add a
   deterministic v1.0.1 comparison and explicit machine-readable allowlist.
   The current compatible cmd/ply package addition is allowed; incompatible
   changes are not.
2. Export the Cobra command and flag tree to sorted JSON, compare it with a
   v1.0.1 baseline and explicit delta allowlist, and bind help/output/exit
   behavior for host install, status, upgrade, and build.
3. Add scripts/verify-install, verify-status, verify-upgrade, and verify-build
   with matching negative meta-tests. Use fresh host artifacts, isolated HOME,
   copied fixtures, and loopback inputs. Each verifier must fail for missing
   artifacts, false assertions, bad exit behavior, or project-state leakage.

After each move, run focused tests and measure from a clean commit. At P2A exit,
run the full gate and hermetic count-2 test.

# Automatic Handoff

Before this agent session ends, finish and commit a coherent move or record an
exact resumable state. Rewrite the rolling handover, update checkpoint statuses,
answer this archive, create one linked NEXT archive, and replace the launcher's
mutable session regions. No separate agent-restart message is required.

If P2A is complete, activate P2B. If P2A is unfinished, keep P2A active and
write a concrete resume mission. Stage only handoff files in the handoff commit
named docs: prepare next agent session. Do not launch the next session
yourself. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
