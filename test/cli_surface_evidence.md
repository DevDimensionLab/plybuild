`bash test/cli_surface_evidence_test.sh` checks failure retention through the
existing surface script's `GO` override. Its small build fixture fails with exit
37 or creates a CLI that emits separate streams and fails the first help check.
The regression checks the script exit, retained raw build/CLI streams, stderr
locator, private directory mode, and an unrelated sibling. It cleans its own
fixtures on success. Pass an absolute evidence directory as the first argument
to retain those fixtures and per-case argv/stdout/stderr/exit for inspection.

Normal `make test-cli-surface` still builds one fresh absolute CLI binary and
runs the existing scenarios. A successful run removes its temporary area. A
failed run retains that private area and reports its path on stderr. Build
streams are saved as `build.stdout` and `build.stderr`; a failed build also
replays its stderr and returns the build's exit status. CLI assertion failures
return 1 and retain the existing per-scenario streams. Remove a retained area
manually once its evidence is no longer needed.

`test/al24-acceptance.sh` runs the regression and the real successful surface
script against a clean native candidate, checking that successful cleanup leaves
an unrelated sibling alone. It also records the supplied candidate binary's
version, Go build metadata and hash. Raw acceptance evidence is retained in a
private `/private/tmp/ply-al24-acceptance.*` directory, outside the native
runner's automatically removed temporary parent. This checks the exercised local shell/CLI
paths; it makes no claim about other platforms, providers or all signal cases.
