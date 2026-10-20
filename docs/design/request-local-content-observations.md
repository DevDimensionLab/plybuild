# Verified content within registry observations

`systemWorkItemStore` owns the physical content observation used by its content
and queue closure validators. Every `readRegistry` with closure validation gets
one read-only `TaskContentStorage`; both validators share successful verified
bytes indexed by the exact root, kind and digest. Proposed-registry closure
validation in `publish` gets a separate reader. Registry decoding, every closure
check, and scope interpretation still run as before.

The observation ends when those synchronous validators return. There is no cache
on `Dependencies`, the system store, the locked session or the process. Another
Snapshot, publication's old-registry reread, proposed-registry validation,
recovery Snapshot and a later lock all start with physical reads. This deliberately
leaves repeated work across phases for later investigation. An observed content
success cannot make a registry reference or Git/runtime observation fresh.

The underlying storage still verifies physical paths, directory/file modes,
bounds and digests. Errors are never retained: a missing object can be created
and read later. Successful bytes are copied on retention and cache hits, so no
caller can alter the retained value. The observation cannot publish. Scope is
copied from the physical reader, preserving shared registration and Task closure
rules. It adds no registry snapshot reuse or decoded-document reuse.

The existing `taskContentCapture` serves a different read-contract purpose: it
attributes projection failures and final physical source-change checks to owners.
Its final rereads and error attribution remain intact. It is not used for mutable
registry generations.

`TestWorkItemSnapshotReusesVerifiedContentWithinObservation` exercises the consumed
store boundary. The baseline with only the existing store fault hook forwarded
for counting made 14 physical reads against 7 sources; the optimized reader makes
7. A second Snapshot performs those reads again. Other regression tests cover
corruption between snapshots, publication and relock, missing-then-created content,
root isolation, owned byte slices, scope, recovery and existing queue behavior.

Performance evidence is maintained privately with a frozen synthetic six-goal
queue, separate timing/count test executables, source archive, fixture digest and
replay scripts for subsequent A–D Tasks. Timing excludes fixture setup/reset and
mutation preview; the mutation measurement includes actual fixture-only goal
materialization and worktree preparation, without launching an agent. It does
not represent the user's complete registry or prove real startup latency.
