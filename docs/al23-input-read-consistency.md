# AL23 input reads during in-place writes

The existing reader can consume mixed bytes during a same-inode, same-size
rewrite. It does not establish an atomic snapshot. The existing digest binding
remains truthful: it hashes the consumed bytes, and a later `inputBound` read
rejects different bytes. This investigation found no false final-file binding
and adds no consistency guard or new gate.

`TestAL23SameInodeReadBindsConsumedBytes` reads four bytes from `AAAAaaaa`,
writes `BBBBbbbb` through a real `WriteAt` on the same file, then reads the
remaining bytes. The inode and size stay unchanged. The consumed `AAAAbbbb`
hashes to `18907f150da56804603eca81221ccfb54b47f6bf085eac866eeaa2b48ed746e5`;
the final `BBBBbbbb` hashes to
`3c949f46c3369edb1e5f0bfb9880e71e910bc05a745117540a1b13957e1af6cb`.
The test uses the existing `inputBound` consumer to reject the final file under
the consumed-byte digest, accept the exact consumed bytes, then bind and accept
the completed writer's stable bytes. It follows `readInputBinding`'s existing
`FileBinding{path, hash(raw)}` rule. Metadata identity alone proves no content
consistency; this test makes no arbitrary-writer snapshot guarantee.

The only production change is a local read-function argument, supplied with
`io.ReadAll` by `readFileAs`. Checks, byte limits and binding rules are unchanged.
The initial undefined-seam compile failure is seam TDD, not evidence of a
production behavioral defect. A baseline run instrumented the original
`796383021488` reader solely with this read seam through a Go overlay; it passed
the same deterministic mutation and instability controls. No failed behavioral
assertion or production fix is claimed.

`TestAL23ReadInstabilityControls` retains rejection of inode replacement and
growth beyond the bounded read, with role/path diagnostics and stable retry.
Existing named AL05 tests preserve ordinary 0644 input, role protections,
acceptance/review qualification and execution of the exact acceptance snapshot
when the source changes after binding. That after-binding fault is separate
from AL23's during-read mutation.

`TestAL23CandidateCLIRejectsFinalBytesForMixedDigest` uses the actual supplied
make-built CLI and the existing isolated native-workspace fixture. Provider and
owner bootstrap are modeled; no live provider is called. The CLI rejects report
evidence carrying the mixed-byte digest when the path contains final bytes,
leaves report history empty, then accepts a fresh stable binding and an
idempotent retry without qualifying a candidate merely from progress. The
under-read write is injected only in the unit test, not in the CLI subprocess.
Each CLI invocation logs argv, cwd, exit, stdout and stderr.

The retained acceptance entrypoint is `test/al23-acceptance.sh`. Native verify
runs it once against the reviewed clean commit and exact candidate binary,
preserving raw output, binary hash/build metadata, named focused tests, vet and
diffcheck. Development baseline and seam-TDD evidence is kept outside the Task
at `/private/tmp/ply-al23-evidence`; native verification and review are separate
facts from local integration and coordinator installation/closeout.
