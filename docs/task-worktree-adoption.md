# Adopt an existing Task source

This implements sections 1–7 of [the adoption contract](task-worktree-adoption-spec.md).
The contract's section 8 remains the historical AL44 planning boundary.

Use a previously unused active Task with its first Goal already registered and
pending. The Goal must select automatic local Epic delivery. Run from the
registered Epic return checkout. The source must be an unclaimed physical linked
worktree, on a full local branch, clean and idle, with useful committed work
strictly descended from the current Epic base. Clean status alone does not
establish writer absence: native managed inventories and bounded local cwd
observations must succeed. Writers outside that bounded observation can remain
unobservable; the receipt is not an ongoing external-writer lease.

Adoption requires registry format 7. Upgrade explicitly through the existing
release-pinned `workspace registry upgrade --to-format 7` preview/apply after all
workspace consumers have a compatible delivered reader. Existing formats and
creation records retain their original semantics; no historical records migrate
implicitly during adoption. Keep the native upgrade's byte-exact backup.

```sh
ply workspace task worktree adopt TASK_ID \
  --project PROJECT --repo REPO --epic EPIC \
  --path /physical/source --ref refs/heads/useful \
  --expected-oid SOURCE_OID --format json
# Repeat the same inputs with --apply --confirmation sha256:PREVIEW_DIGEST.

ply workspace task show TASK_ID --format json
ply workflow execute --spec GOAL_SPEC_ID --check
```

Preview writes nothing. Apply reserves native metadata and binds the existing
source; it moves no Git ref, creates no worktree, starts no agent and inherits no
qualification. Its receipt records source/start `S`, chosen Epic parent `B`, the
exact Goal and provenance. Ordinary Goal execution uses this path and branch,
starts a fresh mandate at `S`, and reviews/tests the entire `B..C` delta. The
review brief names `B` and `S` separately. Delivery and later cleanup use the
existing native integration, retention and owner-handoff/follow lifecycle.
Execution authorizes ordinary source cleanup, including this adopted branch and
worktree; adoption alone does not authorize deleting either.

If source, Goal or Epic base changes before execution materialization:

```sh
ply workspace task worktree release-adoption ADOPTION_ID \
  --reason 'Source or chosen basis changed before execution' --format json
# Repeat with --apply --confirmation sha256:RELEASE_PREVIEW_DIGEST.
```

Release removes only the exact still-owned metadata claim. It preserves Git,
files, Goal/Problem, queue order and all adoption/release receipts. Source
availability, cleanliness and writer absence are not release prerequisites. The
pending Task reports `adoption_released` and cannot fall through to creation.
Freshly preview adoption with the current source OID and valid current basis;
new identity follows the causal release chain even when bytes are unchanged.
Old adopt/release/execution tokens cannot revive or clear the new generation.

Execution and release compete under the same registry lock. Execution persists
its request before content publication; a lost publication reply recovers that
same request and preparation. A materialization intent, preparation, Run or
possible start fences release. Unknown metadata publication must be resolved by
native readback; it is never treated as writer absence or no effect.

A bounded remaining recovery limit: after the execution intent is durable, a
source/base change blocks the frozen request and metadata release. This feature
adds no generic cancellation or permission to restore old bytes. Preserve the
intent and actual observations for supported recovery or separate coordinator
selection. The isolated regression records this boundary before any Run exists.

Acceptance uses temporary Git repositories, injected writer observations and the
existing modeled Herdr/provider transport. It proves software behavior, including
separate review, verification and cleanup; it makes no claim about actual OS
writer absence, live provider acceptance, human product judgment or installation.
