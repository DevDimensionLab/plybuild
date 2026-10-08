# Human-started integration and Task closeout

`ply integration` lets a person review one qualified Delivery, confirm its exact
candidate and effects, and follow integration through Task closeout. Run it from
an existing workspace checkout, usually the registered return checkout. The
current directory never selects the target.

Technical verification, human judgment, integration, installation, closeout and
notification remain separate facts. A published PR is not a merged PR. An
accepted asynchronous request is not an observed merge. Failed installation
does not undo or hide an already observed integration.

## Prepare and hand off the candidate

Use the [Delivery workflow](workflow-delivery.md) to qualify and register the
TaskResult. New goals can explicitly select human integration ownership:

```json
{
  "schema_version": 2,
  "integration_owner": "human",
  "mode": "local_epic_integration",
  "project_id": "example",
  "repo_id": "product",
  "epic_id": "product-fixes",
  "target_ref": "refs/heads/product-fixes",
  "target_worktree": "/absolute/workspace/product/product-fixes"
}
```

For this local agreement the developer qualifies and records Delivery, then
releases writing ownership. Agent callbacks cannot integrate. For `pull_request`,
the developer still publishes the exact PR with required QA and metadata before
handover; merge requires the separate human plan. The current owner releases:

```shell
ply workflow execute release wfr_<run> --context /absolute/private/context.json \
  --reason 'The candidate is ready and this owner has stopped writing.'
```

This revocation-only command may use a newer installed CLI. It does not replace
the frozen control executable or change the original agreement. A person without
the owner's private context can use `ply integration release --delivery dlv_<id>`
only when Ply positively observes that exact native provider has ended. An idle
session or unknown process state does not release ownership.

When handing over from an older immutable controller, finish its verification,
report and required QA callbacks first. New release revokes that controller. A
later `fail` or `blocked` before integration keeps immutable answer/QA history and
clears only the effective takeover, allowing the original owner to correct it.

Schema 1 agreements retain their original meaning. A separately preserved human
takeover can use an existing qualified Delivery; it adds no PR merge or main/master
permission to historical agent mandates. After takeover, agent integration
entrypoints are blocked. A negative judgment returns correction ownership and
preserves the answer and prior evidence.

## Select, inspect and answer

```shell
ply integration list --project example
ply integration --delivery dlv_<id> --check
ply integration --delivery dlv_<id>
ply integration show int_<id> --format json
ply integration scan --project example
```

Multiple Deliveries require an explicit ID. Preview shows candidate commit/tree,
source, registered target, prior QA, merge method, optional install commands,
cleanup resources and notification destination. It rechecks this exact plan after
input; changed evidence or choices invalidate the answer.

Enter exactly `pass`, `fail` or `blocked` followed by Enter. `pass` records the
person's judgment of this candidate and confirmation of the displayed effects.
Existing valid QA is shown, so an already reviewed PR does not require repeating
the same product exercise. A later negative answer blocks an older pass.
`--actor` and `--observation` preserve the person's claim and findings as local
provenance, not identity authentication. EOF, an incomplete line, silence and
`--force` cannot replace an answer.

`list`, `show`, `scan` and `--check` perform no workspace, Git, registry, queue or
external mutations. PR preview reads current GitHub facts. `--format json`
returns envelopes with `kind` and `schema_version: 1`, exact plan and decision
bindings, observed receipts, reasons, history and next action. Read-only preview
can return a blocked state successfully; incomplete execution exits nonzero.

## PR merge and local return

Select a method for an already published, verified PR:

```shell
ply integration --delivery dlv_<pr> --method merge --check
ply integration --delivery dlv_<pr> --method squash
```

Ply checks repository, head, base, permissions, mergeability and permitted methods,
then sends the expected head to GitHub. It requests direct merge with rule bypass
disabled. It does not create an absent PR, force push, choose a fallback method,
request auto-merge or enroll a PR in a queue. The actual merge OID and URL are
preserved, including squash. The local base checkout stays unchanged.

Local Epic return requires no PR flags:

```shell
ply integration --delivery dlv_<epic-return>
```

An explicit `local_branch_integration` agreement also supports `refs/heads/main`
and `refs/heads/master` at the registered return checkout. It needs no remote,
GitHub tools or credentials; a configured remote adds no network effect. Changed
or dirty source/parent and incompatible history block return without rebase or
reset. Native integration updates the exact Epic base and closes only the correct
Task queue entry. PR closeout does not close a newer current Task after publication
already closed the original entry.

## Optional installation

Installation defaults to none. Place a private profile outside the Task source
and select it with `--install-profile /absolute/private/install.json`:

```json
{
  "kind": "ply.integration.install-profile",
  "schema_version": 1,
  "name": "local-product",
  "candidate_oid": "<exact candidate commit>",
  "artifact_path": "/absolute/bin/product",
  "artifact_sha256": "sha256:<expected installed file digest>",
  "install": {
    "executable": "/usr/bin/install",
    "arguments": ["-m", "0755", "/absolute/staging/product", "/absolute/bin/product"],
    "cwd": "/absolute/staging"
  },
  "verify": {
    "executable": "/absolute/bin/product",
    "arguments": ["--version"],
    "cwd": "/absolute/staging"
  }
}
```

Supply real values for the selected candidate and built artifact. Commands use
argument arrays. Both the installed file hash and verification command must pass.
Integration precedes installation; cleanup waits for verification. Failed or
unknown installation keeps the source and remains visible followup.

Ply preserves a controller outside the source and installed artifact. If the
running executable would be replaced or removed, `controller_required` includes
the exact saved `resume` command. Invoke that command explicitly; Ply does not
silently restart or replay input.

## Closeout and recovery

Closeout preserves result, Spec, verification, review, QA and integration evidence
outside the source and retains the candidate commit under an archive ref. It then
removes only the exact clean Task worktree and expected local branch and completes
its lifecycle. Git GC cannot remove the archived candidate, including after
squash. Historical manifests are never rewritten.

Use `--keep` in the original plan to finish while retaining source and branch.
Otherwise tracked, untracked or ignored changes, hidden index flags, nested
repositories, submodules, worktree locks, changed identities, symlinks and active
or unknown writers block cleanup. Running inside the source being removed is
also blocked. There is no force cleanup, broad prune or remote branch deletion.

```shell
ply integration --delivery dlv_<id> --keep
ply integration resume int_<id> --check
ply integration resume int_<id>
ply integration resume int_<id> --retry-install
ply integration resume int_<id> --retry-merge
```

Resume observes the same operation and completes remaining steps. Unknown effects
are never blindly repeated. `--retry-install` allows a specifically observed
failed installation; restore approved profile bytes first. A pending or unknown
PR is observed by the preserved remote operation ID. Concurrent calls and
multiple Delivery IDs for the same candidate/target share effect ownership.

`--retry-merge` requires positive evidence that the previous request had no merge
effect, a fresh exact PR preflight and the current plan's human pass. An open PR
alone is insufficient. Previous requests and failed read observations remain in
history; uncertainty about a newer attempt cannot reuse an older rejection.

If an unexecuted plan must change, inspect and explicitly reconsider it:

```shell
ply integration reconsider int_<id> --check
ply integration reconsider int_<id> --observation 'The selected method is unavailable.'
```

The second command asks for an actual `fail` or `blocked` answer. It preserves the
original decision and separately records negative native QA and revocation.
Source ownership must be released again before a new plan is confirmed. Immutable
reservation history connects the replacement plan to the revoked one. Pending,
unknown, integrated, installed or closed effects cannot be cancelled this way.

Required followup stays separate from observed integration. Completed Tasks leave
active lists; history reports retired or deliberately kept resources. Notification
status cannot undo integration or make incomplete closeout appear complete.

## Legacy reconciliation

Scan distinguishes active work, integrated work with outstanding closeout,
retained/retired resources, missing sources and unowned worktrees. Age, directory
name and physical absence never prove integration.

```shell
ply integration scan --project example
ply integration release --task legacy-task
ply integration --task legacy-task --reconcile --check
ply integration --task legacy-task --reconcile
```

Legacy release requires the person's exact `released` answer and checks known
native owners. Reconciliation requires one unambiguous exact native integration
and source identity. It closes only that historical Task without another merge,
invented Delivery, replacement historical QA or installation. An unfinished
Integration must be resumed, including when required installation failed;
reconciliation cannot bypass it.

## Optional integration notification

Notifications default to off. Select an existing local route explicitly with
`--notification-route /absolute/private/route.json`. Route and state must survive
Task cleanup. No personal route, recipient or credential is embedded in versioned
files, and global configuration is not changed.

Ply records one stable integration event only after observing the actual effect,
with Task/Delivery/Integration identities, candidate, target and actual OID/PR.
Native transport preserves attempts and deduplicates by that event. Transport
acknowledgment differs from human receipt. Failed and unknown sends stay visible;
ordinary resume observes a previous attempt without blindly resending.

For a positively proven undelivered attempt, native transport offers a separate
explicit retry. Inspect its exact event, route and confirmation digest first:

```shell
ply integration retry-notification int_<id> --check
ply integration retry-notification int_<id> --apply --confirm sha256:<shown-digest>
```

Unknown and acknowledged sends cannot be retried. This command runs no merge,
installation or cleanup; repeating an already used confirmation observes its
recorded attempt instead of sending again.

## Verification

`test/delivery_acceptance.sh` builds the actual CLI and runs all Go packages,
isolated native taskrun fixtures, static checks and whitespace checks. Integration
tests use disposable Git repositories and controlled GitHub/notification adapters.
Synthetic provider and human records are labelled as test data. They never merge
a real PR, clean up an existing Task or supply human approval of the implementation.
