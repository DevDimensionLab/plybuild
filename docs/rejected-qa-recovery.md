# Recover a rejected human QA reservation

`workflow execute recover-qa` repairs one narrowly identified reservation: the
original human QA attempt rejected the known `+00:00` timestamp form before
publishing a human QA record. It preserves the rejected attestation, attempt,
original callback control and native history. Recovery does not record a human
verdict or authorize integration.

Run the reviewed repair executable from the original Task worktree, in the same
live owner session. Supply the exact original `qa-...` attempt ID, private
context and a current runtime observation. The observation must come from the
owner's actual current provider, session and effective authority; requested
permissions or a newer executable on disk are insufficient. Its contract is
documented in [active-delivery-continuity.md](design/active-delivery-continuity.md#provider-upgrades-and-current-runtime-observations).

First inspect the read-only preview:

```shell
/absolute/repair/ply workflow execute recover-qa RUN_ID \
  --attempt 'qa-<64-hex-attempt-digest>' \
  --context /absolute/original/context.json \
  --runtime-evidence /absolute/current-runtime.json \
  --check --format json
```

Replace the example attempt ID with the preserved rejected attempt. Inspect its
run, attempt, state and next action. Apply only the exact confirmation digest
returned by that preview:

```shell
/absolute/repair/ply workflow execute recover-qa RUN_ID \
  --attempt 'qa-<64-hex-attempt-digest>' \
  --context /absolute/original/context.json \
  --runtime-evidence /absolute/current-runtime.json \
  --confirm sha256:EXACT_PREVIEW_DIGEST --format json
```

Applying without `--confirm` is rejected. `--check` and `--confirm` are mutually
exclusive. Native checks must still establish the same owner, accepted authority,
unchanged required artifacts, known timestamp rejection and absence of QA
publication. A different attempt, another failure cause, an unknown publication
result or changed authority must not be cleared by this operation. A retry keeps
the original attempt ID, so it cannot clear a later QA reservation. Inspect any
uncertain result before retrying.

The recovery receipt establishes only that the rejected reservation was safely
recovered. Prepare a new attestation file with supported UTC timestamps, preserving
the human's actual answer and exact candidate binding. Keep the rejected file
unchanged. Then use the original bound callback control's ordinary QA command:

```shell
/absolute/original/ply-control workflow execute qa RUN_ID OUTCOME \
  --context /absolute/original/context.json \
  --evidence /absolute/corrected-human-attestation.json --format json
```

`OUTCOME` is the actual human's `pass`, `fail` or `blocked` judgment. Recovery,
technical tests and a successful command do not supply that judgment. Only an
actual matching candidate pass and the original delivery authority permit the
separate integration operation. The repair command does not replace the original
control, restart a provider or submit a replacement human answer.

Repair-tool tests use isolated native fixtures with explicitly synthetic provider,
runtime and human data. Their results validate the repair mechanism; they are not
a human pass for this implementation or for the blocked delivery. Testing must
not edit a live reservation or reuse its private evidence as a fixture.
