# Retire resources while retaining delivery effects

From the containing workspace outside the Task source, use:

```sh
ply workflow execute retire RUN_ID --retain-effects \
  --observer 'actual coordinator' --reason 'resource ownership only' --format json
```

This opt-in mode supports started, stopped, interactive local Epic deliveries whose owner and every known Task/session writer are already absent. It never closes a pane or sends provider input. Native attempts with unresolved outcomes, terminal history, foreign/corrupt bindings and unknown writer observations block the dependent operation.

Retirement preserves dirty and unique source at its original location, every stop and its evidence, provider annotations and all original effect claims. Partial and unknown effects remain partial and unknown. First recording emits a warning to stderr with the Run, original effects, retained record and source digest. JSON stdout stays parseable.

`execute show` exposes `delivery_status.resource_retirement`. `recorded_pending_release` means the retirement event is durable but exact preparation release remains outstanding. `retired` means that preparation is released or advanced. Neither state means delivered, effect-resolved, permission-confirmed or replayable. Earlier tests, reviews and human judgments remain historical; owner callbacks and continuation stay closed. Fresh work requires an independent Task mandate.

Repeat the same Run, flag, observer and reason after interruption. Pending release revalidates the native bindings and current source around fresh writer absence. Drift during that observation blocks the call. Later settled local changes may progress: an append-only retained observation binds the new source/effects to the original event and previous observation before release of the exact still-owned preparation. Original records are unchanged. Completed retries read history without fresh writer inspection or source hashing; they never advance a successor or mutate a foreign queue slot.

Absence observations cover local processes in the Task directory, holding the provider transcript or naming the session. Processes that left all those scopes are outside this observation method. Resource closure proves no external operation completed or failed. Independent development must retain any unknown external dependency until separately observed through a legitimate supported path.

Default delivery retirement and report retirement retain their existing contracts. The `workflow run retire` alias parses the explicit flag but refuses it for report Runs.
