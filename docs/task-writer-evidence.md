# Task writer evidence

`internal/taskrun/delivery_task_writers.go` is the shared module for composing
Herdr identity and inventory observations with local Task/session process
evidence. Its action entrypoints return present, positively absent, or unknown
alongside the original evidence. Inspection failure, incomplete identities,
missing observer processes and duplicate process identities stay unknown.

The existing Herdr runner and `Dependencies.LocalProcesses` are the seam, with
production and fixture adapters. There is no new dependency interface or
configurable policy engine. Tests exercise composition at these same adapters;
caller tests exercise the actual effects and history.

Resume supplies its identified invoking pane. Only a plain pane and its exact
direct-parent shell without positive session evidence can be excluded; other
ancestors, siblings, Task cwd processes and session holders remain writers.
Retirement reconciles local processes for the current resumed owner and may
wait within its existing timeout for positively observed writers. Unknown
inspection returns immediately without replay. Cleanup checks local writers
for every owner, including a headless owner's actual session. Human recovery
continues to record Herdr absence with process liveness unknown; reading old
records grants no new writer-start or cleanup authority.

The callers still own their distinct policies: resume reserves a new generation,
retirement closes the selected pane and releases its queue only after absence,
and cleanup waits for owner exit before continuing the recorded closeout.
Follow observes a released owner until it settles; failed inspection never
becomes a reason to close another pane. A failed cleanup preserves the delivered
integration and resumes the same closeout, without another Git integration.

The module earns depth by removing inventory/local-evidence composition,
PID-presence interpretation, process-liveness promotion and error-string
classification from resume, retirement and automatic cleanup. Deleting it would
put that knowledge back into those callers. The composition and unknown/present
distinction now have locality in one implementation, while historical records
and action-specific effects retain their existing authority.

## Registering and optionally starting a fresh backup Goal

From the containing workspace or registered return checkout:

```sh
ply workflow execute backup-goal RUN_ID --task fresh-backup --spec fresh-backup-v1 --provider codex --format json
# Reuse the same arguments and returned confirmation for explicit publication:
ply workflow execute backup-goal RUN_ID --task fresh-backup --spec fresh-backup-v1 --provider codex --apply --confirmation sha256:... --format json
```

This G04 dependency requires an already natively retired interactive Goal
Spec@2 Run with an automatic local Epic agreement, released source preparation,
unchanged pre-work source and no candidate or review round. Preview reads only;
apply publishes a fresh Task/Problem/Goal and appends its exact reference while
retaining the current preparation and all pending references/order. It creates
no branch/worktree and starts no provider. Identical retries recover the same
native publication keys, including an interrupted partial publication. Inspect
Task publications and queue after an uncertain acknowledgement before repeating
identical input; changed identity/provider/settings conflict rather than replace.

Provider and fresh Task/Spec IDs are required. Model and effort are optional:
omission resolves through normal fresh native defaults (Codex's existing local
settings, or Claude's native defaults). Preview shows resolved values and source;
no old provider-specific model/effort is copied. Optional overrides use shared
launch validation. Registration does not require Herdr launch configuration or
provider permission acceptance. The returned `workflow execute --spec ... --check`
is a separate read-only preflight from the return checkout and may report missing
launch configuration; it is not a provider start.

To register and start together, run from the registered return checkout:

```sh
ply workflow execute backup-goal RUN_ID --task fresh-backup --spec fresh-backup-v1 --provider codex --start --format json
# Use this opt-in preview's confirmation, keeping the same runtime options:
ply workflow execute backup-goal RUN_ID --task fresh-backup --spec fresh-backup-v1 --provider codex --start --apply --confirmation sha256:... --format json
```

`--herdr-workspace` and `--permission-profile` use the existing native controls
and defaults. The read-only preview resolves local runtime and destination
configuration; live Herdr, base and resource checks remain with native execute.
Registration and execution have separate fields in the opt-in response. A later
startup failure preserves the registered refs and the native next action. Repeat
identical input for the same native Run/pane/session readback, including after
Task preparation. Once a native intent or request exists, explicit destination/profile
choices must match it before readiness or startup, even before Run publication.
Omitted choices retain that launch and preview shows its preserved effective
profile. An uncertain launch never automatically restarts, resumes or
switches providers. A new launch creates its own mandate and callback context;
startup is not an owner receipt, runtime permission observation or product work.
An unrelated current preparation blocks startup without changing its queue.

Original product requirements and immutable document descriptors remain binding.
A normative runtime identity addendum records the source Run/request/Goal/ref and
fresh executor; it overrides historical runtime prose without deleting arbitrary
product text. No old acceptance, permissions, tests, review, QA, candidate,
qualification, result or delivery transfers. Active retirement before this supported boundary,
post-work/source transfer, review quota, report/legacy, PR and cross-repository
cases remain later G04 work; this bounded registration/start feature does not complete G04.
