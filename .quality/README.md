# Quality apparatus

This directory applies Engagement B (Lift) from the pinned methodology. Run the
project wrapper, not the vendored upstream script directly:

```sh
bash .quality/tools/quality-audit.sh . \
  --baseline .quality/baseline/scorecard.json \
  --manual-evidence .quality/manual-evidence.json
```

The wrapper writes deterministic `scorecard.json` and
`raw-upstream-scorecard.md` under `target/quality-audit/`. The raw Markdown is
preserved for traceability but can contain upstream verdicts corrected by the
wrapper. The JSON is authoritative. Exit `0` means every selected criterion
passes, exit `1` means measured project findings, and exit `2` means the audit
or its evidence is broken. A missing population remains `UNMEASURABLE`; a
baseline regression is `FAIL`. A held or improved ratchet passes only when its
measurement preconditions and populations are non-empty.

The local structured pass corrects known gaps in the pinned upstream Markdown:
Q0.6 validates the central safe-writer apparatus and actually executed control
tests, Q1.3 uses a type-resolved Go scan with exact adapter boundaries, and Q2.1
validates every subject's production roots and named executable harness.
`internal/testutil` and `*_test.go` are explicit test-support
exclusions from the production-side-effect population; an ordinary production
package is never excluded by substring. Local packages are loaded from source;
function and method values are traced to their bound standard-library effects,
including across local package boundaries. Interface origins are followed
through factories, named and multiple returns, fields, indexes, ranges, type
assertions, closures, and exported package bindings. Module export data is used
to distinguish resolved third-party methods from known effect candidates; an
unresolved potentially effectful call makes the audit broken. In-memory
`io.Copy` destinations are not counted as filesystem effects. The catalog
includes process execution, network/DNS, HTTP and TLS, database, clock,
filesystem, `syscall`, and cgo boundaries, including `os.CopyFS`, `os.OpenRoot`,
and effectful `*os.Root` methods.

The scanner consumes the file population selected by `go list`, so native-platform
GOOS/GOARCH filename constraints, CGO, build tags, and `GOFLAGS` are honored.
Cross-target `GOOS`/`GOARCH` contexts fail closed because this initial scanner is
compiled and executed on the audit host. JSON tool
metadata binds the result to `go version`, `GOOS`, `GOARCH`, `CGO_ENABLED`,
`GOENV`, `GOWORK`, `GOFLAGS`, `GOEXPERIMENT`, the architecture-specific Go
selectors, and
path-and-content digests for the selected production/test files. Path-sensitive
flags such as `-overlay`, `-modfile`, `-toolexec`, and `-C` are rejected because
their effective source population cannot be bound by this initial scanner.
Selected files, build context, inventory bytes, and every audit instrument are
revalidated before publication. Measurement cleanliness compares repository
bytes directly with `HEAD`, including ignored files, tracked vendor files, and
index-hidden changes.

A supplied baseline is accepted only when it is a complete clean measurement
with the same module, inventory checksum, audit instruments, Go toolchain,
build selectors, and selected source population. Empty current populations
remain structured `UNMEASURABLE` results; an empty or incomplete stored
ratchet population is broken evidence.
The wrapper fixes `CGO_ENABLED=0`, `GOENV=off`, and `GOWORK=off` to match the
pinned upstream audit and ignore machine-local persisted Go settings. `GOFLAGS`
that can skip tests, short-circuit execution, or change the source path outside
the bound manifest are rejected.

Copy `.quality/baseline/manual-evidence.json` to `.quality/manual-evidence.json`
only as a schema example. Replace every receipt, commit tree, and optional
inventory-overlay hash with measurements from the commit being audited. Manual
evidence is accepted only for a clean measured tree. Evidence from another
commit, a dirty tree, duplicate receipts, or an undocumented survivor
classification is rejected.

Do not use the current output file as its own baseline. A `--baseline` path that
aliases any generated artifact is rejected before outputs are invalidated.

Run the audit meta-test with:

```sh
bash .quality/tools/test-quality-audit.sh
```

`methodology.lock` records the source revision and content checksums for each
vendored file. Re-vendoring is a deliberate update: refresh the files, the lock,
and the meta-test in the same change.
