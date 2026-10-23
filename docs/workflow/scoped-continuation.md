# Scoped same-owner continuation (DEP02 / TABBB0A4A)

This operation records a decision to continue specific native work while the
original denial and unknown effects remain active. It is **not** effect
reconciliation, OS sandboxing, runtime approval, or evidence that old writers
are absent. It does not change the frozen Goal, request, context, session,
mandate, provider policy, registry format, or delivery target.

The producer uses disposable modeled consumers only. Only pwa-plan and the
original PWA owner may later apply a qualified version to T792F7F9B. Missing
historical invocation, structured command return, fixture handle and wait stay
missing. Existing receipt qualification remains insufficient for resolving the
original partial effect.

## Public operation and request

From the exact original Task worktree and owner session:

```sh
NEW_PLY workflow execute disposition RUN --context ORIGINAL_CONTEXT --file request.json --format json
NEW_PLY workflow execute disposition RUN --context ORIGINAL_CONTEXT --check-operation prepare --format json
NEW_PLY workflow execute show RUN --format json
```

Use actual values and SHA-256 bindings. The request is ordinary JSON; native
records are private, immutable, chained files under `delivery/dispositions`.
Identical concurrent calls and lost-reply retries return the original record;
changed content under the same event ID fails. A new source/candidate assessment
needs a new ID. A historical successful registration retry is a readback, never
renewed operation eligibility.

```json
{
  "kind": "ply.workflow.scoped-disposition",
  "schema_version": 1,
  "event_id": "preparation-1",
  "run_id": "ACTUAL_RUN",
  "request_sha256": "sha256:ACTUAL_REQUEST",
  "mandate_sha256": "sha256:ACTUAL_MANDATE",
  "native_session": "ORIGINAL_SESSION",
  "source_oid": "CURRENT_HEAD",
  "source_tree": "CURRENT_HEAD_TREE",
  "retained": [{
    "stop": {"locator": "ORIGINAL_STOP", "sha256": "sha256:ACTUAL_STOP"},
    "unknowns": [
      "child_wait_and_descendants", "temp_socket_and_cleanup",
      "runtime_input_bytes", "filesystem_writes", "pending_writes"
    ],
    "conflicts": [],
    "assessment": "Explain the retained effects and any concrete possible overlaps; do not assert absence.",
    "provenance": "owner_and_static_assessment",
    "evidence": [{"locator": "PRESERVED_SOURCE", "sha256": "sha256:ACTUAL_SOURCE"}]
  }],
  "operations": [{
    "operation": "prepare",
    "action": "Describe the new already-authorized operation, never the refused action.",
    "boundary": "fresh_source_snapshot",
    "reads": ["path:/EXACT_TASK"],
    "writes": ["path:/EXACT_TASK"],
    "children": ["process:bounded-file-validator"],
    "assessment": "Assess independence against concrete source/artifact/child boundaries under retained uncertainty.",
    "provenance": "owner_and_static_assessment",
    "inputs": [{"locator": "OPERATION_SOURCE", "sha256": "sha256:ACTUAL_INPUT"}]
  }]
}
```

`retained` must name **every** unresolved stop, in original order. All five
unknown dimensions are retained even if no concrete collision is identified.
`conflicts` identifies concrete possible path/process/network/native-effect
collisions; it is not an exhaustive inventory of historical effects. Empty
conflicts does not mean no old writes or descendants. Owner/static assessments
are judgments; referenced source is not upgraded to runtime byte/write evidence.

The native decision uses a fresh observed candidate-input snapshot and immutable
new operation inputs. An unlocated unknown historical write is not automatically
a universal veto; the new product disposition permits bounded work while keeping
that uncertainty. A concrete possible overlap remains a veto for the dependent
operation. Source/input drift, extra unresolved stops, changed owner or native
unknown attempts cannot be hidden behind resource labels or a newer stop.

## Six-step consumer journey

| Step | Native boundary and actual effects | Retained uncertainty and guards |
| --- | --- | --- |
| 1. Prepare | `prepare` / `fresh_source_snapshot`: original owner writes acceptance, source, validator and docs in its Task; declares build/cache/temp writes and bounded children. Snapshot may start dirty. | Does not prove which bytes the old fixture opened or whether its writes ended. Concrete overlapping input/write/child effects block. Check before working; new source requires a fresh record. |
| 2. Candidate and tests | Original owner runs authorized non-browser checks/build/artifact validation, waits its own new children, commits a clean candidate, then registers its exact OID/tree and operation inputs. | Actual check/build/precache effects are writes and children, not read-only. No Playwright, browser, process probe or old-command replay. Historical runtime filesystem and pending writes remain unknown. |
| 3. Review | `review` / `exact_clean_candidate`: coordinator starts the ordinary separate native review; reviewer returns from its own context/session. | Reads exact candidate; ordinary authorized test cache/temp effects must be declared. No implementation writer. Dirty/drifted source, findings, changed basis or absent native return cannot pass. |
| 4. Verify | `verify` / `candidate_and_isolated_verifier`: original owner uses the current bound control and unchanged original context, omits `--review`, supplies the bound candidate executable. | Native acceptance/binary snapshots, isolated HOME/TMPDIR, bounded execution, actual exit and before/after candidate checks remain. Bind acceptance and executable bytes in `inputs`. No hand-authored receipt or review bypass. |
| 5. Deliver | `integrate` / `qualified_candidate_and_exact_target`: ordinary automatic local integration under the frozen agreement. | Exact qualified/current candidate, fresh review, inputs_match, exit 0, acceptance, target/base and ordinary writer/attempt guards remain. Registration itself neither integrates nor supplies QA. |
| 6. Retain source | Delivery readback reports `delivered` only after the native local receipt and queue/base closure. | `retained_effects_cleanup_pending` keeps original owner and source. No cleanup, retire, writer transfer, follow/probe, terminal-close fabrication or replay follows from this disposition. |

Resource declarations describe concrete operation effects. Native source,
acceptance, run-evidence, verifier scratch and frozen integration-target boundaries
are also checked; omission or renaming cannot remove those dependencies. The
candidate snapshot covers tracked and nonignored untracked source, not an
exhaustive runtime filesystem audit. Ignored build/cache outputs are declared
operation effects; important external inputs (including acceptance and binary)
are bound separately. Native operations recheck candidate and input stability.

## Frozen controls

The actual PWA control is d58dbd4 (preserved SHA-256
7044a9996acd208a2bacf9bc849e3f486b72a765aec5f4b1feb888458934cb32).
It keeps its original active-stop veto. The installed producer base 68ea053
is not the original consumer control. The newer CLI
registers the disposition without rewriting that executable or its request or
context. Use the existing supported compatibility transition when necessary:

```sh
NEW_PLY workflow execute continue RUN --context ORIGINAL_CONTEXT --check --format json
NEW_PLY workflow execute continue RUN --context ORIGINAL_CONTEXT --format json
# Use preview.control_executable.path returned by continue:
COMPATIBLE_CONTROL workflow execute verify RUN --context ORIGINAL_CONTEXT --candidate-binary CANDIDATE --format json
COMPATIBLE_CONTROL workflow execute integrate RUN --context ORIGINAL_CONTEXT --format json
```

The original control is retained and remains readable; only the native compatible
control becomes the callback executable. Continue starts no owner and replays no
input. The original Codex 0.163 launcher has been legitimately removed. Use the
existing `--runtime-evidence` protocol with actual same-session observations
from the original owner; the executable test models this missing-launcher case.
Without that observation, continuation remains blocked. No fabricated policy
or reinstall is supplied by this feature.

`delivery_status.scoped_continuation` shows the immutable record, input and source
bindings, retained original effects, individual operation eligibility/block
reasons, runtime per-action status and cleanup blocker. Eligibility is not an
executed action or successful review/test/delivery. Old stop, ordinary resolution
and retirement operations retain their original semantics.
