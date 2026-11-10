`PLY_ROUNDTRIP_ROOT=<dir> bash test/cli_surface_evidence_test.sh` checks the
existing surface script through its `GO` override. Its small build fixture fails
with exit 37, creates a CLI that emits separate streams and fails the first help
check (also with a `find` that fails), creates a CLI that passes every help call
so that a later assertion fails on an older help stream, or hangs until the run
gets TERM; a last case runs without `PLY_ROUNDTRIP_ROOT`. It runs the cases in a
`PLY_ROUNDTRIP_ROOT` subdirectory whose name contains a space. The regression
checks the exits 37, 1, 1, 1, 143 and 2, that the failing build or CLI streams
and exactly the help stream the failed assertion read reach stderr, that the run
directory was one `ply-cli-surface.XXXXXX` under `PLY_ROUNDTRIP_ROOT`, that the
hung child was observed and is gone, and that nothing is left in
`.build/compatibility/`, `PLY_ROUNDTRIP_ROOT` or the caller's `TMPDIR`. Without
an argument its own fixtures live in one
`.build/compatibility/ply-cli-evidence.XXXXXX` directory, removed on exit. An
absolute evidence directory passed as the first argument is used as given and
kept with per-case argv/stdout/stderr/exit for inspection.

Normal `make test-cli-surface` needs `PLY_ROUNDTRIP_ROOT`, builds one fresh
absolute CLI binary and runs the existing scenarios in one owned directory under
that root; the README build section describes it. A failed build replays its
stderr and returns the build's exit status; CLI assertion failures return 1.
Before its run directory is removed, a failed assertion prints the first 80
lines of each file it checked. Nothing is retained.

`test/al24-acceptance.sh` is the historical acceptance for an earlier native
candidate and still expects the former retention behaviour. It runs the regression and the real successful surface
script against a clean native candidate, checking that successful cleanup leaves
an unrelated sibling alone. It also records the supplied candidate binary's
version, Go build metadata and hash. Raw acceptance evidence is retained in a
private `/private/tmp/ply-al24-acceptance.*` directory, outside the native
runner's automatically removed temporary parent. This checks the exercised local shell/CLI
paths; it makes no claim about other platforms, providers or all signal cases.
