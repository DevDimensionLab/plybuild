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
