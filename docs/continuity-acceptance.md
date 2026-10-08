# Try continuity after a provider upgrade

Install the reviewed candidate using the normal local Ply install procedure.
From this checkout, prepare an isolated exercise with that installed executable:

```shell
./scripts/continuity-journey prepare
# Or select the installed path explicitly:
./scripts/continuity-journey prepare --ply /absolute/installed/ply
```

Preparation starts a disposable native delivery with the historical reader from
commit `f2de0c638f1d2b1c38417aff70b79ca7bb4be31a`, qualifies its first candidate,
records an explicitly synthetic failed product check, and uses the installed
candidate to publish another Task's newer valid `delivery` metadata. It then
commits a correction and executes the old native verifier. The acceptance exits
0; qualification fails because that reader cannot understand the newer metadata.
The earlier candidate and its judgment remain separate from the correction.

Preparation then removes the accepted synthetic provider launcher and points the
fixture's ordinary `bin/codex` lookup at a distinct synthetic installation. The
accepted request still identifies the original path and digest. It executes the
pre-upgrade continuation reader from commit
`6472bb87e11de6512adcde65c03b1be41ce46646` and the original bound report callback;
both actually fail on the missing launcher. Their exits and outputs are retained,
and readback confirms that neither recorded a transition or report. The original
bound control predates `execute continue`, so that failing continuation uses a
separate preserved control; neither original control is overwritten.

The preserved historical source emits `task_content_invalid_input: object has
missing or unknown fields`. The live incident's unidentified old build wrapped
its schema failure in a publication-chain error. The regression reproduces the
same compatibility boundary and successful-verification/failed-qualification
split; it does not claim the historical binary is the incident binary.

The command prints an absolute `journey` path. Use its short commands in the
same terminal:

```shell
/printed/path/journey inspect
/printed/path/journey report-incomplete
/printed/path/journey continue
/printed/path/journey verify
/printed/path/journey inspect
```

Inspect shows the removed historical launcher, distinct current installation,
unchanged native session, preserved failures, successful verification attempt,
older qualified candidate and its prior judgment separately. Report-incomplete
uses the installed native `execute report --incomplete` path to durably record
the original owner's stopped outcome. It grants no verification, QA, control
replacement or integration authority. Repeating it records one report.

Continue supplies a fresh, explicitly synthetic runtime observation and
preserves a compatible control with durable proof for the same delivery. The
native checks compare its context, acceptance, current session, actual policy
and allowed effects with the accepted run and live transport observation. The
current provider installation supplies none of that authority. Verify reuses
the matching successful receipt and qualifies the correction without running
its acceptance again. The counter must remain two and the provider-start and
tab counts must remain one. Unrelated publications, the original request,
acceptance, controls and context remain intact. Repeating continue or verify
preserves the same transition and candidate.

Optionally exercise the separate human gate and local return in this disposable
workspace after examining those observations:

```shell
/printed/path/journey qa pass
/printed/path/journey integrate
/printed/path/journey inspect
```

`qa fail` and `qa blocked` are also explicit choices; `--answer 'your exact words'`
preserves additional feedback in the native observation beside that explicit
outcome. Multiline feedback is quoted reversibly within the native plain-text
limit. A failure does not integrate. This answer judges
the disposable candidate only. Tell the real delivery owner your actual result
for the installed implementation separately; automated fixtures never supply
that answer.

No command opens a real provider or terminal, sends a notification, touches a
live Task, changes global configuration, or publishes remotely. The provider,
accepted and current authority observations, and first QA record are labeled
synthetic. All native records and actual command outputs remain under the printed
private temporary directory. The script creates runtime observations only for
this isolated fixture; it cannot establish a real owner's current permissions.

## Genuine unsupported cases

A real owner whose launcher was removed or changed must provide a current
`--runtime-evidence FILE` observation for the same accepted native session,
policy and exact effects. The original acceptance or a newer executable on disk
is insufficient. This is the owner's runtime observation, not renewed human
permission. The native continuation rejects missing or stale observations,
another session/provider/Task, changed authority, altered required artifacts and
unsupported contracts. A future provider launch retains its own executable and
permission checks.

The bounded incomplete-report command can record `stopped` or `needs_input`
without establishing current authority for Task effects. It still requires the
exact original context, immutable inputs and live owner/session identity. If
those cannot be verified, its rejection remains a rejection; a local JSON file
alone is not a durable native report. Preserve the actual output and resolve
the reported dependency through the original owner's supported workflow.

Receipt reuse also requires the same source, acceptance script, review and
required evidence. If any changed, fresh verification is required. Continuation
and requalification never provide the exact candidate's actual human judgment.

## Regression prerequisites

The required native tests build both pinned historical sources using
`git archive`. Both commits must exist locally, so a shallow CI checkout must
fetch them or use a full checkout before acceptance. No network fetch happens
inside the tests. A runner may set `PLY_CONTINUITY_OLD_BINARY` and
`PLY_CONTINUITY_BEFORE_UPGRADE_BINARY` to physical executables already built from
their respective pinned sources to reuse those builds. Missing historical source
is a failing precondition, never a skipped compatibility test. The native suite
also tests a replacement at the old path and relocation of the old executable;
all scenarios preserve the original binding and install different current bytes.

`test/delivery_acceptance.sh` discovers the native continuity tests with the full
suite. The journey exporter is opt-in through `PLY_CONTINUITY_JOURNEY_ROOT`; the
full-suite runner rejects that variable so it cannot accidentally export another
workspace. `scripts/continuity-journey` sets it only for the exact exporter test.
