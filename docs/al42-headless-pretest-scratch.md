# AL42 headless pretest scratch

New headless pretest intents bind `pretest_scratch` to the private physical
allocation returned by the existing automatic verifier allocator. The same path
is frozen in TMPDIR and the existing verifier sandbox policy. Inherited TMPDIR
inside a workspace falls through to the allocator's existing external fallback.
HOME, acceptance snapshot, candidate identity, native logs and native reservation
ownership retain their existing contracts. The sandbox-exec argv is unchanged.

Only preparation without a published intent can remove a never-started allocation.
An existing or unobservable intent on a failed preparation preserves scratch and
warns with Run/job/path. After launch, cleanup requires a positively exited,
quiescent process, complete untruncated output, no timeout/error, zero exit and
successful durable publication of both logs and result. Reserved-result recovery
validates native result/log bindings before applying the same decision. Unknown
reservations are never replayed. Failed tests still enter bounded correction.
Legacy intents lack the optional field and confer no scratch cleanup authority;
legacy results lack the optional completion field and remain readable.

Failed, incomplete, interrupted or unknown attempts and publication errors retain
child evidence with a WARN locator. Cleanup errors warn without changing the
native test outcome. There is no background cleanup or discovery of old roots.

## Verification

The first regression was run before the production change:

```
go test ./internal/taskrun -run '^TestPretestScratchReservedWorkspaceFixture$' -count=1
```

It failed with `workspace_init_nested` when workspace.Init used the actual
reserved pretest TMPDIR under worker/attempts/pretest-0001/tmp. Raw red output is
preserved at `/private/tmp/ply-al42-evidence/red.log` for this execution.

`TestPretestScratchReservedWorkspaceFixture` now creates that workspace and checks
physical/private allocation, Run/job/candidate binding, bounded policy and the
unchanged sandbox wrapper. Its dedicated runner is explicitly `/usr/bin/true`:
reservation is real, while provider and OS sandbox execution are not claimed.

Focused tests cover allocation failure, acceptance snapshot conflict, reservation
failure, an error after durable intent publication, completion conditions, legacy
TMPDIR preservation, native stdout/stderr/result publication failure, retained
child evidence, unknown reservation recovery without replacement, and bounded
failed-test correction. Process completion inputs in those tests are modeled.

`TestPretestScratchOrdinaryCLIWorkspace` separately invokes an ordinary unsandboxed
/bin/sh with the reserved environment and exact make-built PLY_CANDIDATE_BINARY.
The shell uses mktemp and `ply workspace init`, writes child evidence, exits, and
the test checks its surviving locator. This proves the real CLI can host a
workspace under the chosen scratch and child evidence survives shell exit; it
is not evidence of live OS sandbox enforcement or provider behavior.

The private native acceptance entrypoint runs only these named AL42 tests plus
the existing policy-construction test. It does not run a full suite or historical
refused sandbox/PS/provider journeys. Parent owns separate native review,
closeout, installation, installed smoke and the whole-job notification.
