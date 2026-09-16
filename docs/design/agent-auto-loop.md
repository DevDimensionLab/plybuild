# Automatic Agent Loop Controller

Status: accepted for the `codex/upgrade-quality` worktree.

## Purpose

`codex-dev-start.sh` owns one validated Codex lifecycle. It intentionally stops
when the roadmap is complete, when a real failure occurs, or when a bounded
product choice requires authorization. `codex-dev-auto.sh` is the outer
controller that distinguishes those cases and continues only after a separate
decision agent has recorded a safe choice.

Start the automatic controller from the worktree root:

```sh
./codex-dev-auto.sh
```

It runs the ordinary launcher unchanged. After every launcher exit it starts a
fresh ephemeral Codex decision turn with structured JSON output. The decision
turn inspects the tracked lifecycle, repository state, and retained launcher
logs, then returns exactly one action:

- `restart`: record one explicit bounded authorization in one local commit,
  leave a valid `NEXT` session, and run the launcher again;
- `complete`: confirm both the lifecycle and authorized roadmap are complete;
- `blocked`: preserve the evidence and stop for human review.

The controller does not blindly choose option 1. It prefers a recommended
option 1 only when the active session explicitly offers it, the current
selection and dependency metadata stay unchanged, the target remains unloaded
and runtime-unreachable, no new finding exists, and exact expiry guards are
recorded. Architecture, parent, source, dependency, Go-floor, expired-guard,
or ambiguous choices remain human decisions.

## Verification Boundary

Decision-agent prose is not authority. Before a restart, the controller
requires all of the following:

1. exactly one new commit whose parent is the pre-decision HEAD;
2. the reported commit resolves to the new HEAD;
3. changed paths are limited to `codex-dev-start.sh`, the rolling plan and
   handover, the session archive directory, and the continuity design;
4. a clean worktree and a successful launcher check in `NEXT` state;
5. a non-empty recorded decision and matching structured `NEXT` result.

Completion is accepted only when the decision agent leaves HEAD unchanged and
both its structured result and the launcher report `COMPLETE`. A blocked result
must also leave HEAD unchanged. The controller never pushes, merges, releases,
stashes, reverts, removes the worktree, or implements the product choice.

## Logs And Cleanup

Every invocation prints its external run directory. By default the directory
is below `${TMPDIR:-/tmp}/ply-codex-dev-auto-logs`; it is rejected if it is
inside the worktree, is a symlink, is a non-empty unowned directory, or has an
unknown ownership marker.

Each cycle retains:

- a timestamped high-level `run.log`;
- complete live launcher console output;
- the decision prompt, schema, final JSON, bounded JSONL, and bounded stderr;
- bounded tails of launcher final messages, JSONL, and stderr.

Raw launcher supervisor trees and decision scratch are removed after every
cycle, including interrupt cleanup. Ten run directories are retained by
default, and pruning removes only recognized `run-*` directories under the
owned log root. An atomic active-run lock prevents two outer controllers from
starting against the same log root. Inspect the latest retained state with:

```sh
./codex-dev-auto.sh --status
```

Validate configuration without starting either agent with:

```sh
./codex-dev-auto.sh --check
```

`CODEX_DEV_AUTO_LOG_ROOT`, `CODEX_DEV_AUTO_LOG_KEEP_RUNS`, and
`CODEX_DEV_AUTO_MAX_CYCLES` override the bounded defaults. Exit status 0 means
verified completion, 2 means a retained blocker or rejected decision, and 130
means interruption.

## Runtime Contract

The decision turn uses fresh `codex exec --ephemeral`, JSONL events, an output
schema, and a final-message file. These are the documented non-interactive
automation surfaces in the
[official Codex documentation](https://learn.chatgpt.com/docs/non-interactive-mode).
The contract test substitutes recording launchers and a recording Codex
executable, so it never invokes a real agent or network service.

Recheck with:

```sh
make test-agent-auto
```
