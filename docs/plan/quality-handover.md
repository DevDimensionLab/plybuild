# Quality Upgrade Handover

Generated: 2026-08-25T04:48:50+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 549685d618f8.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, and 549685d in roadmap order. The
separate operational continuity implementation is 1b85711 and changes no Go
quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T044850+0200-migrate-file-find-all.md. The
find-first predecessor is answered history, and the reciprocal archive graph
has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 24 Preserved

The existing zero-value-safe filesystem adapter now exposes only one additional
operation, Walk(root string, callback filepath.WalkFunc) error. Its system
implementation delegates the exact root and callback once to filepath.Walk,
and a zero dependency returns filesystem.ErrNoFilesystem before root access.
Every complete filesystem recording double was extended mechanically.

Exported file.FindFirst(fileSuffix string, dir string) (result string, err
error) retains its signature and observable behavior. It is the production
wrapper around one private complete filesystem dependency, and production
selects filesystem.System(). The private helper delegates the exact root and
callback through filesystem.Walk. Its unchanged callback ignores fi and errIn,
matches only strings.HasSuffix(path, fileSuffix), assigns the exact matching
path before returning io.EOF, stops on the first matching callback, and
normalizes only a final walk error equal to io.EOF. Every other exact error and
named result remains unchanged.

Eight top-level adapter and file recording contracts bring the suite to 202
tests. They prove complete dependency delivery, production system selection,
exact root and callback values, callback order, ignored metadata and errors,
case-sensitive exact suffix matching, first-match result and stop, no-match
empty result, exact non-EOF errors with the current partial result, final EOF
normalization, safe zero behavior without developer-path access, system walk
behavior, and non-empty walk and callback populations. The file contracts
perform no real filesystem mutation.

Only FindFirst's direct filepath.Walk moved. FindAll, GrepRecursive, IDE
cleanup, Render, ClearDir, Move, DeleteAll, DeleteSingleFile, append-open,
file-read, directory-create, file-create, file-overwrite, file-existence,
Bitbucket, Wpost, supervisor, inventory, public API, callers, and every P2A
contract remain unchanged.

## Measured Quality State

The clean full audit at 549685d618f8 reports:

- Absolute L0: 8 of 8.
- 202 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 20 guarded safe-writer sites, 15 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 42 direct external sites outside five declared adapters of 58
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims across 38 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: tree f483d9b2bf5e69dd55b094bc361c514d7efb80fc, status
  SHA-256 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d,
  and zero dirty paths.

The clean report is /private/tmp/ply-find-first-clean-audit-549685d and the
focused report is /private/tmp/ply-find-first-focused.3DyVdh. Compatibility,
preflight, Go-gate, acceptance, and hermetic outputs also stayed under
/private/tmp.

## Decisions And Learned Facts

1. The adapter Walk operation is a resolvable interface dependency, not a
   function-valued effect, and is safe at its zero value.
2. FindFirst's callback body stayed unchanged, including deliberately ignoring
   callback metadata and errors.
3. Q1.3 moved from 43/58 to 42/58 because the direct FindFirst walk disappeared
   and the already-counted system Walk boundary was reused.
4. FindAll is the next isolated flow. It can reuse the existing Walk operation,
   nominally moving Q1.3 to 41/57 without another adapter extension.
5. FindAll deliberately preserves traversal order, appends every path matching
   strings.HasSuffix unless SuffixIn reports an excluded substring, ignores fi
   and errIn, never stops on a match, and normalizes only final io.EOF.
6. .quality/inventory is baseline-checksum-bound. Do not relabel seams or claim
   P5 mutation coverage.
7. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight, and
   keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
8. The ignored .agent-task/current.md compatibility mirror is absent.
9. Supervisor success requires process, structured-stream, and committed
   repository evidence; final prose is observable only.
10. The partial-raw-log signal fixture can flake. The first complete make test
    run did; the standalone launcher, preflight, and immediate complete rerun
    all passed 62 controls. Never weaken or skip this contract.

## Next Objective

Move only exported file.FindAll(suffix string, excludes []string, dir string)
(result []string, err error)'s direct filepath.Walk(dir, callback) through the
existing zero-value-safe filesystem Walk operation.

Start red with file recording contracts for the complete dependency,
production system selection, exact root and callback delivery, callback order,
ignored metadata/errors, exact suffix and exclusion decisions, ordered
all-match accumulation, no-match nil slice and nil error, exact non-EOF walk
errors with the current partial result, final io.EOF normalization, safe zero
behavior without developer-path access, and non-empty walk/callback
populations. These contracts must not mutate the real filesystem.

Keep FindAll's signature and callback body unchanged. Production must select
filesystem.System(), and the private helper must carry one complete filesystem
dependency. Reuse adapter Walk without extending the adapter or changing
complete doubles. Do not move GrepRecursive, IDE cleanup, SuffixIn, callers, or
another operation, and do not enter HTTP/process, config, clock/server, P4,
P5, or later roadmap work.

Expected direction is nominally Q1.3 41/57 with Q1.2, Q1.4, and exact Q2.1
held. Regenerate the structured result and require zero comparable ratchet
regressions.

## Verification Notes

Completed from implementation commit 549685d618f8:

- Red evidence: focused compilation failed on absent adapter Walk and private
  find-first helper/dependency symbols.
- Focused filesystem, file, both caller, and adjacent complete-double package
  tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls.
- API, CLI, subprocess compatibility: PASS.
- make test: first run hit the known signal flake; immediate full rerun PASS.
- make test-install, uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable.
- Clean full audit: exit 1, 16 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable.
- Implementation commit: 549685d618f8d11ceda2f5d5964fc83ee7e538d6
  (quality: route find-first walk through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, all file code/tests and FindAll callers, the adapter and
complete doubles, and the named audits before editing. Confirm branch, HEAD,
status, reciprocal links, and ./codex-dev-start.sh --check. Begin red for only
FindAll's walk boundary, then finish with one implementation commit and one
separate handoff-only commit.

Stop before recursive grep, IDE cleanup, another adapter operation or family,
Q1.4 expansion, P4, mutation, Docker, cloud, distribution, or publication.
Stop on API/CLI change, comparable ratchet regression, full-audit exit 2, or
failure to isolate implementation from generated and handoff-only state.
