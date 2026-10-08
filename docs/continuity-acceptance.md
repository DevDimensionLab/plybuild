# Try active-delivery continuity

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

The preserved historical source emits `task_content_invalid_input: object has
missing or unknown fields`. The live incident's unidentified old build wrapped
its schema failure in a publication-chain error. The regression reproduces the
same compatibility boundary and successful-verification/failed-qualification
split; it does not claim the historical binary is the incident binary.

The command prints an absolute `journey` path. Use its short commands in the
same terminal:

```shell
/printed/path/journey inspect
/printed/path/journey continue
/printed/path/journey verify
/printed/path/journey inspect
```

Inspect shows the current source, successful verification attempt, older
qualified candidate and its prior judgment separately. Continue preserves a
compatible control with durable proof for the same delivery. Verify reuses the
matching successful receipt and qualifies the correction without running its
acceptance again. The counter must remain two and the provider-start and tab
counts must remain one. Unrelated publications, the original control and the
original context remain intact. Repeating continue or verify preserves the same
transition and candidate.

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
accepted authority and first QA record are labeled synthetic. All native records
and actual command outputs remain under the printed private temporary directory.

## Regression prerequisites

The required native tests build the actual historical source using `git archive`.
The pinned commit must exist locally, so a shallow CI checkout must fetch that
commit or use a full checkout before acceptance. No network fetch happens inside
the tests. A runner may set `PLY_CONTINUITY_OLD_BINARY` to a physical executable
already built from the pinned source to reuse that build. Missing historical
source is a failing precondition, never a skipped compatibility test.

`test/delivery_acceptance.sh` discovers the native continuity tests with the full
suite. The journey exporter is opt-in through `PLY_CONTINUITY_JOURNEY_ROOT`; the
full-suite runner rejects that variable so it cannot accidentally export another
workspace. `scripts/continuity-journey` sets it only for the exact exporter test.
