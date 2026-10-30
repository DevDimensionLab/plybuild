# PWA consumer binding and executable demo (TABBB0A4A)

> Historisk dokument etter fjerningen av den kompilerte prosesskjeden.
> Execute-, worker/week-, delivery- og integrasjonsorkestreringen som omtales
> her er ikke en kjørbar oppskrift i dagens kilde. Fjernede kildefiler er
> bevart som historiske stier. Se [gjeldende direkte operasjoner](../../README.md#human-owned-process).

This completes the [early six-step contract](scoped-continuation.md) with
consumer-01–03. It supplies no actual PWA disposition or installed qualification.
The producer's tests use disposable Runs, fake Herdr/provider identities and
modeled reviewer judgments. Native records, CLI callbacks, Git, npm/Node children,
file writes, waits, verification and local branch integration are real.

## Original context and runtime evidence

The preserved PWA executable is d58dbd4def24bee385307db7dc18fe1583958533,
SHA-256 7044a9996acd208a2bacf9bc849e3f486b72a765aec5f4b1feb888458934cb32.
Its original launcher was Codex 0.163 and has legitimately been removed. Do not
restore or invoke it. The old d58 active-stop gate is expected to keep refusing;
it never acquires the new disposition command.

From the original Task cwd/session, the **qualified new CLI** first registers
and reads the disposition using the **original context**, before switching the
callback executable. Save each chosen disposition request privately (0600),
using a new file/ID for changed bytes. This registration uses existing accepted authority and
current source facts, without requiring the removed launcher. Native observed
owner identity must still match. No command is sent to an updater or shell.

```sh
NEW_PLY workflow execute disposition RUN --context ORIGINAL_CONTEXT --file preparation-1.json --format json
NEW_PLY workflow execute disposition RUN --context ORIGINAL_CONTEXT --check-operation prepare --format json
```

The existing runtime-observation contract and trust boundary are documented in
[active delivery continuity](../design/active-delivery-continuity.md#provider-upgrades-and-current-runtime-observations).
The exact struct is `DeliveryRuntimeObservation` in
`internal/taskrun/delivery_runtime_observation.go`; it accepts these fields only:

```json
{
  "kind": "ply.workflow.delivery-runtime-observation",
  "schema_version": 1,
  "run_id": "ACTUAL_ORIGINAL_RUN",
  "request_sha256": "sha256:ORIGINAL_REQUEST",
  "session_id": "ply:ACTUAL_ORIGINAL_RUN",
  "context": {"locator": "ORIGINAL_CONTEXT", "sha256": "sha256:ORIGINAL_CONTEXT_BYTES"},
  "acceptance": {"locator": "ORIGINAL_ACCEPTANCE", "sha256": "sha256:ORIGINAL_ACCEPTANCE_BYTES"},
  "runtime_claim": {
    "runtime_id": "codex",
    "model_id": null,
    "profile_id": null,
    "effective_policy_sha256": "sha256:PRESERVED_UNCHANGED_RUNTIME_OBSERVATION",
    "native_session_id": "ACTUAL_ORIGINAL_NATIVE_SESSION"
  },
  "delivery_permission": {
    "launch_contract_sha256": "sha256:ORIGINAL_BOUND_LAUNCH_CONTRACT",
    "permission_confirmed": false,
    "startup_protocol": "mandate_receipt_v1",
    "mandate_sha256": "sha256:ORIGINAL_FROZEN_MANDATE",
    "actual_policy_evidence": [{"locator": "PRESERVED_RECEIPT_RUNTIME_OBSERVATION", "sha256": "sha256:ACTUAL_OBSERVATION_BYTES", "role": "runtime_observation"}],
    "delivery_agreement_sha256": "sha256:ORIGINAL_AGREEMENT",
    "allowed_effects": ["COPY_THE_EXACT_ALREADY_ACCEPTED_EFFECTS"]
  },
  "observed_at_utc": "ACTUAL_RFC3339_UTC"
}
```

This is the actual PWA consumer's `mandate_receipt_v1` protocol:
`runtime_permission=runtime_per_action`, `permission_confirmed=false`, and the
original frozen mandate digest. The field named `effective_policy_sha256` binds
the preserved recipient runtime observation; it does **not** assert OS policy
coverage or blanket permission. Keep the exact already-recorded protocol,
receipt observations, agreement/effects and locator-sorted evidence. Never derive
permission from root/helper metadata or change these values to match an example.

Historical legacy Runs retain their separate accepted-authority protocol and
may have `permission_confirmed=true` with `effective_policy` evidence. That is
historical compatibility only, not the PWA protocol and not an instruction to
make a new blanket permission claim.

Save the actual chosen observation privately (0600). Its time must satisfy the
existing native freshness guard: ten minutes, at most thirty seconds clock skew.
This is a technical observation freshness check, not a deadline for a person.

The original owner observes its actual current provider/session and policy from
its running agent's supplied runtime context, with supported read-only Herdr
owner identity/cwd/status observations. Herdr does not establish effective policy.
The owner must substantiate its current runtime context and the unchanged
accepted observation/mandate bindings; copying launch-request settings or locating a new executable does not
suffice. Unknown/changed authority stops the transition. No ps, browser,
launcher/version/quota probe, global policy change or private transcript search
is necessary or authorized by this route.

Use the **same bytes and path** for preview and apply:

```sh
NEW_PLY workflow execute continue RUN --context ORIGINAL_CONTEXT --runtime-evidence CURRENT_OBSERVATION --check --format json
NEW_PLY workflow execute continue RUN --context ORIGINAL_CONTEXT --runtime-evidence CURRENT_OBSERVATION --format json
```

Read `preview.control_executable.path` from the apply response. Keep the old
control/context/request bytes. Use that returned compatible control for owner
callbacks; no restart, new session or rewritten original control is needed.

## Refresh points and actual child effects

The executable journey begins with dirty README, package and validator source.
It retains unknown historical input bytes, exhaustive filesystem/pending writes,
child/wait/descendants and temp/socket/cleanup throughout. It does not reduce the
old effects to listener labels or present a static source as runtime proof.

1. Register `dirty-preparation` against current dirty inputs and all original
   stops. Check `prepare` before source/acceptance/validator work.
2. Intentional writes change the source and acceptance inputs. Keep that old
   record; create `prepared-node-inputs` with a new ID, fresh native snapshot,
   current input hashes and assessments **before** running the changed checks.
3. Execute actual offline npm `check`, `build` and `test:build` scripts. Their
   modeled Node programs spawn and wait a Node child, write `dist` and
   `dist-server`, validate built JSON artifacts/precache, and write cache/temp.
   These are actual non-browser subprocesses and writes, not read-only tests.
4. Preserve the actual preparation output; commit the exact complete source.
   Register `exact-candidate` with the new clean OID/tree, acceptance, executable
   and preparation-output bindings before dependent review/verify/integrate.
5. Corrections, validator/acceptance changes, another commit or changed important
   input require another immutable assessment/event with current source and
   inputs. A new candidate requires fresh independent review and verification.

Ignored build/cache outputs are explicit effects rather than candidate source
inputs. Important consumed external inputs must be bound separately. Concrete
possible overlaps block their dependent operation, including implicit native
source/input/scratch/target dependencies even if a declaration omits them.
Unlocated old writes remain unknown; current snapshots do not prove their absence.
Historical registration retry returns its original observation and never renews
operation eligibility after drift. `show` remains readable when work is blocked.

## Review, verify, local delivery and retained cleanup

The coordinator starts a normal native review from its separate pane:

```sh
NEW_PLY workflow execute review RUN --coordinator ACTUAL_COORDINATOR --format json
# Reviewer, from the separately started native reviewer pane and context:
REVIEW_CONTROL workflow execute review-return RUN --context REVIEW_CONTEXT --file ACTUAL_REVIEW_RETURN --format json
# Original owner, original context, returned compatible owner control:
COMPATIBLE_CONTROL workflow execute verify RUN --context ORIGINAL_CONTEXT --candidate-binary EXACT_CANDIDATE --format json
COMPATIBLE_CONTROL workflow execute integrate RUN --context ORIGINAL_CONTEXT --format json
NEW_PLY workflow execute show RUN --format json
NEW_PLY workspace task queue list --project PROJECT --repo REPO --epic EPIC --schema-version 2 --format json
NEW_PLY workspace epic base show EPIC --repo REPO --format json
```

The review must return `passed` for the current clean candidate. Owner verification
omits `--review`; a hand-authored review cannot replace the separate round. Check
`verification.current`, `inputs_match`, qualified outcome and actual exit 0, and
current automatic acceptance. Review evidence gaps preserve uncertainty without
inventing browser QA. The integration receipt binds the exact candidate/target,
queue closure and base update; `final_delivery.state=delivered` is distinct from
cleanup.

This **automatic** route requires **no explicit release**. The original owner
remains reserved while historical effects are retained. Ordinary release semantics
are preserved: a live stop cannot acquire writer-transfer authority through the
new integrate scope. `release`, `retire`, cleanup, policy expansion, browser
replay and foreign writer operations are not supported disposition operations.

The CLI journey demonstrates a native integration reservation with a deliberately
lost pre-effect reply. Report, new verification and another registration remain
blocked while it is pending; only the same exact native integration ledger may
resume. The resumed CLI integrates once. It returns exit 5 and
`retained_effects_cleanup_pending`: actual delivery is observed, cleanup is not
complete. Original stop/unknowns, source and owner remain. No follow/ps inspection,
provider/pane close or source deletion is used to manufacture completion.

Concrete blockers remain separate: changed source/acceptance/binary, relevant
possible writes/children, added or older omitted stops, unknown native attempts,
missing/negative/stale review, noncurrent qualification/acceptance, changed frozen
target/base or ordinary writer guards. Native checks never transform a technical
pass into human QA. Root alone later supplies the exact installed qualification;
pwa-plan still owns every actual consumer action.

## Run the exact candidate demo

Build a clean committed candidate with `make build PLY_BUILD_OUTPUT=/private/ply`.
Use the preserved original control path supplied by the coordinator, then:

```sh
PLY_DEP02_OLD_BINARY=/absolute/preserved/d58/ply-control scripts/demo-scoped-continuation.sh /private/ply /private/dep02-evidence
```

The script checks clean source/make identity and the exact original-control SHA
before running focused public-CLI fixtures. It exports original runtime-observation
and continuation bytes, request, review-return, acceptance text, actual Node output,
verification receipt, qualified/delivered readbacks and queue/base observations.
The disposable paths in these files are historical fixture paths; the exported
bytes are the evidence. These are modeled consumer facts, not actual PWA delivery
or live reviewer acceptance. The separate producer review/verify and root's later
installed qualification remain required.
