# Integration evidence storage

New work-item stores publish integration evidence once, outside the mutable
registry. Existing inline stores retain their representation. There is no
migration command, automatic conversion, dual write, garbage collection or
change to integration approval/retry rules in this implementation.

## Physical and logical formats

A fresh `.ply/work-items.yaml` uses physical envelope format 9:

```yaml
format_version: 9
registry:
  format_version: 3
  # Other ordinary registry fields retain their existing logical format.
  integration_authorities:
    - sha256:<digest of one authority's YAML bytes>
  integration_intents: []
  integration_attempts: []
  integration_results: []
```

The four integration arrays contain ordered digest references. Each referenced
object contains exactly one existing typed record, stored at
`.ply/task-content/v1/integration-evidence/sha256/<hex digest>`. Managed directories
are physical 0700 directories; objects are physical 0400 files, bounded to 64 MiB
per record and checked against their digest. This reuses the existing workspace
content publisher and work-item lock. No extra writer, store lock or process is
introduced. Logical lifecycle formats 1–8 and their existing transition permits
keep their meanings inside the envelope; format 1 has no integration arrays.

`systemWorkItemStore` is the authoritative boundary. It validates the complete
logical registry before publication, publishes/syncs the objects, then atomically
publishes/syncs the index. Readback verifies referenced bytes and runs the same
strict record, reference, plan-digest and exact-effect/reflog validators used for
inline history. Missing, corrupt, wrongly typed or misbound evidence stops the
read/write; there is no empty-history fallback or inline repair path.

An attempt's prepared and result-recorded states are separate immutable object
versions. Only the current version is referenced; the prepared object remains
retained. Interrupted publication can leave unreferenced objects. Retry verifies
and reuses these objects without rewriting them. Readers see the old complete
index or the new complete index, never a partially hydrated successful state.
A lost reply/directory-sync error retains existing readback/recovery semantics;
it creates no permission for another Git effect.

A lock-period cache still verifies every referenced object's bytes before reusing
its typed registry. Read-view freshness and Task journal provenance include the
actual referenced objects, including physical-path and immutable-mode checks. Unchanged
records retain their exact object references even when their valid YAML spelling
differs from the current serializer. Normal integration writers, automatic delivery,
closeout, repository locking and recovery use the same hydrated registry rather
than separate copies of history.

## Compatibility and later delivery

Existing inline format 1–8 files continue to be read and written inline, including
their explicit native lifecycle upgrades. They receive no write reduction until
a separately authorized migration is designed and qualified. No live workspace
or installed executable is changed by this Task.

Older readers reject the new physical envelope as unsupported format 9. The new
executable advertises format 9; its internal logical registry decoder deliberately
continues to reject an envelope used as a logical registry. Consumers must use
the workspace store, including the updated direct Task journal reader.

The release package declares formats 1–9, matching the executable. The
manifest/executable equality guard rejects mismatched packages. Every consumer
of a newly created format-9 store must use a compatible reader. Existing inline
stores remain inline; installing this release does not migrate shared data.

## Measured publication

`TestIntegrationEvidenceTaskCreationDoesNotRewriteHistory` builds eight actual
validated no-effect/retry histories in private Git/workspace fixtures (eight each
of authorities, intents, attempts and results), seeds equivalent inline and fresh
stores, and publishes one new unbound Task through the real store in each.
It counts actual new/replaced/changed payload files by inode and native file
state, then sums their published sizes. All 32 prior evidence objects must retain
their inode/state. Locks, directories and the unchanged workspace marker are not
payloads. This measures published bytes, not kernel/block-device write traffic,
read traffic, fsync count or temporary-file overhead.

One isolated run on 2026-11-01 measured:

| Measure | Inline | Fresh object-backed store |
| --- | ---: | ---: |
| Published payload bytes for one new Task | 138,581 | 15,934 |
| Published payload files | 1 | 1 |
| Historical objects rewritten | inline history in the registry | 0 of 32 |

The measured reduction was **88.50%**. Retained payload stock in the fresh seeded
store was 130,691 bytes; stock is distinct from publication I/O. Temporary fixture
path lengths can slightly change absolute byte counts. The immutable evidence
still occupies space and is still validated/read; this change does not claim
zero history read/serialization cost or constant-size reference indices.

Regression coverage includes legacy reader rejection/preservation, historical
lookups, missing/corrupt objects, correct-hash invalid records, wrong-category
references, cached-session corruption, object/index interruption and lost replies.
Existing integration journeys exercise concurrency, source/target drift,
unknown/partial results and reflog recovery through the new fresh-store boundary.
Unreferenced objects and old attempt versions are retained intentionally; nothing
in this Task qualifies them for deletion. The duplicate inline payload is absent
from fresh indices; legacy format readers and validators remain necessary.

## Corrected readback and journal boundary

Persisted integration readback uses the observed physical `RawSHA256` in both
normal and nullable paths. Successful publication returns the store's observed
generation, including its new physical digest. A logical registry without an
observed digest reports unknown persisted hash/freshness; serializing an inline
model cannot establish physical persistence. Regression checks cover format 9,
legacy, lost Git replies, directory-sync reply loss, replay and Task show.

Journal initial reads, final `before_freshness` checks and batch rechecks validate
integration sources through workspace's managed evidence reader. Changes to 0400
object mode or 0700 directory mode invalidate the affected source and integration
axis even when bytes are unchanged. Batch source caching retains the stronger
managed validation when the same path also appears as an ordinary artifact.
End-to-end journal regressions use a restored synthetic private Git/workspace
fixture and exercise both early corruption and the final freshness checkpoint.
