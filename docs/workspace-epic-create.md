# Create a repo-specific Epic

`adopt` registers an existing clean worktree. `create` makes one new branch and
worktree from an explicit ordinary parent; it does not start Tasks or agents.

```sh
ply workspace epic create feature --title 'Feature' --project example --repo example \
  --parent-ref refs/heads/develop --expected-parent-oid <full-commit-oid> \
  --branch epic/feature --path /physical/wrapper/feature --format json
```

The parent must be a clean, symbolic checkout in the registered repository,
without Git operations or Task/Epic ownership. The target's parent directory must
already exist. The leaf and branch must be vacant. Targets inside a worktree,
Git administration or `.ply` are rejected. Version one rejects trees that select external filters,
symlinks, gitlinks, missing objects, sparse checkout and worktree-specific config.
Git runs with hooks/fsmonitor disabled and without inherited Git redirection,
external configuration or transport permission. Unused configured filter drivers do
not block creation; their clean/smudge/process commands are disabled before
parent status inspection, including when working attributes differ from the commit.

Creation requires registry format 6. Ordinary existing operations keep their
current format. After a compatible release is delivered to workspace consumers,
preview `workspace registry upgrade --to-format 6 --package <release> --sha256
<manifest-digest> --check --format json`, then apply with its exact `--confirm`.
Format 4 first uses the existing upgrade to 5. Each transition preserves all
semantic records and the exact original bytes in a native backup. Creation
itself never upgrades a registry.

Before Git runs, one typed intent reserves the Epic ID, parent, target and branch.
`epic show <id> --format json` and `epic list` expose pending operations separately
from usable Epics. An identical retry observes the actual effect: no effect may
retry, branch-only adds the worktree without recreating the branch, and an exact
worktree finalizes metadata even after a lost/nonzero command response. Partial
or unknown effects preserve reconciliation state and evidence; there is no
supported automatic reconciliation or cleanup command in this version. Keep the
record and Git state for producer inspection. No reset, prune or name substitution
is performed. A changed request conflicts with the preserved intent.

Completion atomically records the Epic, `ply_created` anchor, initial base and
completed operation. An identical completed retry reports `already_created`
without requiring the historical parent to remain unchanged. Current Epic
readiness is a separate `epic show` observation: later work, base advancement or
a missing anchor never triggers recreation or rollback. Adopted Epics retain
`adopted` provenance. Detailed creation readback is `WorkspaceEpicReadback@3`;
format-6 list readback is `WorkspaceEpicListReadback@2`.
