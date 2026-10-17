# FL01 startup qualification

The SL02 incident exposed two gaps in otherwise existing startup support:
recipient acceptance required handwritten protocol files and hashes, and a
coordinator could not record a sourced callback denial without impersonating
the owner. Native launch identity, preserved control binaries, same-tab resume,
actual permission checks, writer reservations and idempotent start already
existed. The historical SL02 attempt was not replayed or modified.

`execute begin` now owns acceptance bookkeeping. `observe-stop` preserves the
coordinator's actual denial separately; `resolve-stop` records the same owner's
cause repair before continuation. Readback distinguishes terminal creation,
Task prompt, observed owner receipt, actual permission and target changes.
Notification binding is optional native bookkeeping against an already
authorized route, with no transport or global setup effect.

The remaining inputs have specific reasons:

| Input or action | Why it remains |
| --- | --- |
| Human goal and delivery choice | Defines the requested result and exact allowed effects; reused afterward. |
| Actual policy, provider/session and filesystem roots | Ply cannot infer runtime authority from the requested launch policy. |
| Existing route authority and recipient timezone | A route's existence alone does not authorize sending or select the recipient's zone. |
| Sourced denial and actual repair | Missing callbacks alone prove neither a stop nor permission to retry. |
| Native folder trust after configuration-lock failure | Existing provider onboarding fallback; separately owned follow-up, outside this repair. |

## Real Claude qualification, 17 October 2026

The isolated Task `claude-native-begin` ran through ordinary installed Ply and
the installed `ply-delivery` skill in visible Herdr, using Claude Code 2.1.285,
observed model `claude-opus-5-5`, native auto mode, and session
`dc8e7d09-5660-4c0f-a1b5-2ade37d8a701`. Its Run is
`wfr_e4df4a156e3d22d7af007c1689cb49ea15ce443ced2e3d77577003b1a6241d67`.

Native `begin` was accepted by the provider's auto classifier and by Ply. The
owner supplied observation flags; it wrote no acceptance JSON, mandate or
startup hash/binding shell. No new global allow rule or permission mode was
introduced. The first relevant title-normalization regression ran at
16:27:07 UTC and failed with exit 1. The focused fix passed at 16:27:12 UTC.
The owner's first observed timestamp to red test was 45 seconds; exact elapsed
time from initial launch to that test was not measured.

After Task input, the owner asked zero extra human questions. Before Task input,
the configuration lock failed and Claude displayed folder trust. One unnecessary
clarification popup was raised; the coordinator handled one native trust action
under the existing exact-folder grant. The same native startup resumed. This is
a qualification limitation, not an end-to-end zero-interruption result and not
new human runtime acceptance. The lock repair belongs to a separate backlog
item. A transient local socket-token print reported by the delegate is retained
only as an incident note; no credential value belongs in qualification evidence.

The isolated candidate `c73f3128aee89612f5b67c4e13302dda726fa295` passed independent
parent review and native automatic verification (`verify-00000002`), then native
local Epic delivery and closure. No HumanQA record or delegate notification was
created. The delivery receipt digest is
`sha256:e6944288b1e7c53b61b3805672a51e7d8c03b152a0a63b9ebe4054c3e111fd43`.

Deterministic native and candidate-CLI tests separately exercise rejected
acceptance and rejected report observations, wrong/unknown authority, wrong
session, insufficient scope, lost-response replay, unchanged notification
identity, retained writer reservations, explicit same-owner continuation and
existing Codex startup. Synthetic provider tests are not the real Claude result.
