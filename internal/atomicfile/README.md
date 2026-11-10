# File replacement

`Replace(destination, contents, mode) error` owns the shared temp/write/chmod/
file-sync/close/replace/directory-sync/cleanup sequence. It needs an existing
parent directory. Callers own encoding, validation, locks, domain errors and
reconciliation. This package knows no Task, Run, registry or workspace schema.
It does not implement write-once, no-clobber or physical-path validation.

A temporary file named `.atomicfile-*.tmp` is created in the destination's
parent directory. Before replacement, failures trigger close and cleanup;
secondary failures are included in the error text without losing the primary
cause (`errors.Is`). Cleanup can itself fail and leave a temporary file.

All failures are `*Error`, with operation, path, underlying cause and `Replaced`.
`Replaced=true` means the replacement primitive returned success, followed by a
directory-sync failure: complete new bytes were published, but directory
durability was not confirmed. `Replaced=false` means no successful replacement
was observed. In particular, an error from rename or MoveFileEx is not proof
that no external effect occurred. Inspect and reconcile the destination before
retrying after uncertain effect; the boolean is not a retry policy. Concurrent
external writers can change the destination independently of this call.

On non-Windows systems replacement uses `os.Rename`; directory sync opens,
syncs and closes the directory. On Windows it uses `MoveFileEx` with
`MOVEFILE_REPLACE_EXISTING | MOVEFILE_WRITE_THROUGH`; directory sync is a no-op,
matching the existing consumers. No Windows runtime or power-loss guarantee is
implied. File modes retain the OS's existing semantics. A successful return
reports completion of these primitives, not proof against SSD or power failure.

The exported surface is one operation and one error type. Fault operations are
private, per call and have no global mutable hooks. Tests share the same internal
sequence used by Replace, and also exercise the exported operation against real
files. Consumers have a narrow private replacement function for boundary faults;
they translate error context and retain the wrapped error for `errors.As`.
Workflowhandoff's write-once operations and other workspace stores keep their
existing mechanics because their additional semantics are outside this change.
