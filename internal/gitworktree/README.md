# Git worktree inventory decoding

`ParseInventory(raw []byte) ([]Entry, error)` decodes the fields of
`git worktree list --porcelain -z`. `Entry` contains only Locator, Ref, OID,
Locked, Prunable, Bare and Detached. The package uses only the standard library.

NUL separates fields; an empty field finishes a record. As in the original
Ply parsers, a new `worktree` field also finishes the preceding record, and EOF
finishes the last record without requiring a trailing NUL. Leading/repeated empty
fields are ignored. Values retain spaces and newlines. Lock/prune reasons are
accepted but only their presence is represented. Unknown fields, repeated fields
inside a record, fields outside a record and empty locators return an error and
no partial result. A subsequent `worktree` begins a new record, not a duplicate
field in the preceding record.

The decoder preserves wire order and does not inspect the filesystem, validate
OID/ref syntax or require HEAD/branch fields. In particular, bare records without
HEAD are representable. It does not change or retain the input bytes; each call
returns independent entries. Empty input returns a non-nil empty slice (`[]`).

Ply's workspace type aliases Entry. Its adapters retain different policies:

| Adapter | Validation and ordering |
| --- | --- |
| `parseWorktreeInventory` (ListWorktrees and integration) | Requires nonempty OID, and nonempty ref unless bare/detached. Resolves every locator to an existing physical directory, preserving filesystem error causes. Sorts by the resolved locator. |
| `parseReadInventory` (inventory/status and closeout) | Requires an absolute clean locator, valid OID unless bare, and full branch ref unless bare/detached. Does not require a checkout to exist. Sorts by the reported locator. |

Format errors now come from this one decoder; their private wording may differ
from the old read parser. Adapter validation messages and filesystem causes stay
with the adapters. Decoding precedes adapter validation, so when several things
are invalid, a format error can be reported before an earlier record's value or
filesystem error. Rejection is preserved; no readiness or identity is attested.
Git execution, before/after inventory comparison, common-directory/ref/OID
verification, cleanliness and effect guards remain in their existing consumers.

This extraction preserves master policy: even a missing prunable checkout fails
the physical adapter. When incorporating it into the workflow Epic, retain that
Epic's existing narrowly scoped missing+prunable exception in the adapter,
including its stored-path validation and permission-error rejection. That policy
does not belong in this decoder or a parser flag.
