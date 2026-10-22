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

After the execution intent is durable, a source/base change still blocks the
old frozen request and metadata release, even before preparation or a Run exists.
Keep that intent, Goal, branch, worktree and refusal history. Do not restore old
source bytes or invent a Run to make `workflow execute backup-goal` eligible.
This physical binding remains retained; it does not prevent useful committed
work from continuing through a separately selected fresh Task.

## Continue committed work through a fresh Task

Use a new Task/Spec identity and the current registered Epic base. Export the
useful committed delta from its **original** base `B` to useful source `S2`;
diffing directly from a newer Epic base can accidentally delete its new files.
These are read-only Git operations on the old source:

```sh
git -C /physical/old-source diff --binary --full-index B S2 > /private/draft.patch
git -C /physical/old-source rev-parse S2 'S2^{tree}'
shasum -a 256 /private/draft.patch
```

Supply `/private/source-facts.md` with the old Task and exact Goal reference,
adoption ID, original B, frozen S1, useful S2 commit/tree, patch digest, and
selected current Epic base. State that the patch is an untrusted work draft
requiring a fresh owner's assessment. Uncommitted/dirty bytes are excluded and
remain in the old source. Git patches retain additions, deletions and executable
modes; the existing document publisher accepts UTF-8 inputs up to its document
limit. Inspect binary/large exports for suitability before registration.

Run from the registered Epic checkout. The following uses the existing first
Goal helper, including its extra-document inputs:

```sh
ply workspace task goal register FRESH_TASK --spec FRESH_SPEC \
  --document /private/fresh-spec.md \
  --extra-document draft=/private/draft.patch \
  --extra-document source-facts=/private/source-facts.md \
  --title 'Continue useful committed work' \
  --objective 'Assess and deliver the useful draft from the current Epic base' \
  --criterion 'The assessed useful behavior passes fresh acceptance and review' \
  --project PROJECT --repo REPO --epic EPIC \
  --delivery local_epic_integration --acceptance automatic \
  --responsible-actor 'Fresh delivery owner' --acceptance-timeout 1800 \
  --format json
# Repeat identical flags with --apply --confirmation sha256:PREVIEW_DIGEST.

ply workspace task show FRESH_TASK --format json
ply workspace task spec show FRESH_TASK --spec FRESH_SPEC --revision 1 --format json
ply workflow execute --project PROJECT --repo REPO --epic EPIC \
  --spec FRESH_SPEC --check --format json
# After inspecting readiness, omit --check to start this selected fresh Goal.
```

Registration preserves the supplied documents in native content storage and
returns their hashes and `contract.required_inputs` locators. It only appends
the new Goal; it neither reorders old/foreign pending entries nor applies a
patch, cherry-picks a commit, creates a worktree, starts an agent or accepts work.
Changed patch bytes invalidate the preview confirmation. Select the fresh Spec
explicitly; do not rely on `--next` to choose it ahead of retained entries.

Preparation creates a different source at the current Epic base. Start that
fresh bound Run and receive its mandate **before** changing this source. The
fresh owner reads and checks the preserved draft and source-facts inputs from
that mandate (including their digests), assesses their suitability against the
new base, then explicitly applies the chosen work. For example, inside the new
Task, with `PRESERVED_PATCH` set to its native draft input locator:

```sh
shasum -a 256 "$PRESERVED_PATCH"
git apply --check "$PRESERVED_PATCH"
git apply --index "$PRESERVED_PATCH"
git diff --cached
# After assessing the staged bytes and modes:
git commit -m 'Apply assessed useful draft in fresh Task'
```

A patch conflict is a new owner's implementation decision, not permission to
restore the old source or change its intent. The fresh candidate needs its own
meaningful tests, separate review of the entire current-base-to-candidate delta,
mandate/runtime decisions and candidate-bound acceptance. No old QA, Run,
permission or qualification is inherited. Native verification/integration/follow
then deliver and close only the fresh Task under its selected agreement; the old
physical source, dirty files, refs and native history stay retained.

`TestAdoptionFreshTaskAfterUnstartedIntent` exercises this journey through actual
native Goal/input/queue/preparation services and Git: an injected registry writer
failure after durable intent, useful S2 bytes, deletion of original-base README,
new executable mode, an advanced Epic base, changed-patch rejection, explicit
fresh-owner application from immutable storage after losing the export, complete
new review scope, separate review, verification, local integration and closeout.
It checks old Task-owned history/source and foreign queue preservation. Its
Herdr/provider/process/writer observations are injected; it proves modeled native
software behavior, not live provider acceptance or a live old-source recovery.
The opt-in first-Goal CLI regressions also use the exact make-built candidate for
public registration/extra-document/readiness commands. No live old Task/Run is
mutated by these tests.

Acceptance uses temporary Git repositories, injected writer observations and the
existing modeled Herdr/provider transport. It proves software behavior, including
separate review, verification and cleanup; it makes no claim about actual OS
writer absence, live provider acceptance, human product judgment or installation.
