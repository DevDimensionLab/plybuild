# Agent Session: Evaluate Inconshreveable Mousetrap Dependency

Status: NEXT
Session ID: `2026-09-05T013847+0200-evaluate-inconshreveable-mousetrap-dependency`
Created: `2026-09-05T01:38:47+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `998c2ecf3f1fe6fcd163e05e4b374173f4cfed4be2d819c7ed9bc3a8a33b6a49`
Previous: [2026-09-05T011325+0200-evaluate-hashicorp-hcl-dependency.md](2026-09-05T011325+0200-evaluate-hashicorp-hcl-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating selected indirect
`github.com/inconshreveable/mousetrap v1.1.0` as one bounded dependency group.
Resolve the canonical latest release and highest floor-compatible candidate
from primary evidence. Implement one exact changed selection only if it
preserves the retained Go 1.18 floor, has an explained minimal closure, and
passes every quality contract; do not manufacture a dependency commit when
the exact selected release is already canonical latest.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain,
completed dependency groups through accepted Google go-cmp v0.6.0, and the
rejected/no-change Hashicorp HCL evaluation. Go-cmp v0.7.0, Viper v1.16.0,
Emoji v2.2.14, ini v1.67.3, fatih/color v1.19.0, fsnotify v1.10.1, latest
gomarkdown `v0.0.0-20260824154242-13c5cf49db8d`, and an HCL selection change
remain rejected for their recorded floor, loaded-behavior,
release-qualification, or self-test failures. Do not revisit those decisions
or combine another module group. P8 remains queued.

The current build list selects inconshreveable/mousetrap v1.1.0. Treat the
target version, latest release, Go declaration, closure, loaded package
population, source history, and test status as unknown until independently
resolved. This session may change only mousetrap's exact required
go.mod/go.sum metadata and the roadmap/handoff record. Do not change
production Go, another dependency, language/toolchain declarations, quality
apparatus, Docker/release inputs, packaging, publishers, or P8 code.

# Measurements At Start

Latest implementation remains dependency-only go-cmp commit
`c314bcb440b5f94871d71a249bca7ae7f87d5543`, exact parent
`cbdb0a915d7d30bb6633cfbb2ccf0323aea1b77a`, clean tree
`9a212b377cb08e7d6e3044fb6dae071058b6242b`. The HCL no-change handoff commit
must have exact parent `99f508d927047f3803e96d082efc15e2dd2422b7`; go.mod/go.sum remain
byte-identical to c314bcb, and ordinary and ignored status must be empty.

Current dependency measurements are 234 selected modules, 3,557 graph edges,
429 complete test packages, a 308-line tidy projection, and exact
Darwin-symbol/Darwin-module/Windows-symbol vulnerability populations 20/30/20.
The retained main module declares Go 1.18 and prefers toolchain Go 1.26.7.

Use exact Go 1.26.7 at
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in PATH, keep GOENV=off, GOWORK=off,
GOTOOLCHAIN=local, and inject no ambient GOFLAGS. Recovery evidence at
`/private/tmp/ply-p7-go1.26.7-recovery.qvk4zs` has a verified two-entry
manifest SHA-256
`1b30193f4f4811f2515223c53c004603afa0a4812b8a571a24c49b252122e83b`.
Retain the verified golangci-lint 2.12.2, GoReleaser 2.17.1, apidiff, and
govulncheck v1.7.0 binaries and hashes recorded in the handover.

Rejected/no-change HCL evidence is
`/private/tmp/ply-p7-hcl-selection.99f508d.CFvyqD`, with verified 21,427-entry
manifest SHA-256
`ad2954c3ebcf8bfe40e5e1684fcfcf662610555e6ed113515750db1a91c1acff`
and selection-summary SHA-256
`c20b0d6354c35adf0e05916610d03a93c0210bb613ae889250316bffc34132d9`.
It proves canonical HCL v1 latest is already v1.0.0, exact get is a no-op,
targeted Vault/Nomad prereleases are not stable candidates, v2 is a different
Go-1.23 module, closure and vulnerabilities remain exact, loaded consumers
pass, and required default HCL module tests/race fail on an upstream test-only
vet diagnostic. No HCL implementation or post-implementation quality run
occurred.

Accepted go-cmp selection, schema-2 review, exact-quality, and regression
evidence remain respectively
`/private/tmp/ply-p7-go-cmp-selection.cbdb0a9.kbzc4C`,
`/private/tmp/ply-p7-go-cmp-quality-review.c314bcb.XeWg4T`,
`/private/tmp/ply-p7-go-cmp-quality-parent.c314bcb.zxFKyw`, and
`/private/tmp/ply-p7-go-cmp-regression-gate.c314bcb-final.rMW0l5`, with
verified manifest populations/digests recorded in the handover. Exact accepted
quality passes all 21 stages and 27 Q0-Q2 rows at L2, 80/80 mutations, all
8/8 mutation and 4/4 acceptance populations, and zero held, regressed,
not-comparable, or dirty counts.

# Role And Boundaries

From fresh external archives and caches, resolve mousetrap releases through
the Go proxy, checksum database, and upstream repository. Record exact tag
commits/times, module Go declarations, checksum pairs, source identity, tag
and commit signature status, relevant release history, and primary Go
vulnerability data. Prove which stable release is canonical latest and which
release is highest compatible with Go 1.18. Distinguish stable releases,
prereleases, retractions, and major-path changes explicitly; do not infer
compatibility from a successful modern-toolchain build.

Measure old versus candidate selected modules, graph edges, complete package
population, checksums, loaded packages/path, explicit exact-get diff, and
`go mod tidy -diff`. Require an explained minimal selection/edge/checksum
closure. If mousetrap is loaded, run focused consumer behavior in addition to
the candidate module's complete tests, repository build, complete
tests/race/vet, pinned lint, byte-identical public help, and identical API/CLI
reports.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, source identity, module
self-tests, loaded behavior, or any repository quality contract fails. If the
selected release is already canonical latest and exact get is a no-op, record
that no-change decision without forcing an implementation commit. Measure the
exact vulnerability population; do not assume parity.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, P7 in
the roadmap, go.mod/go.sum, rejected HCL evidence, accepted go-cmp evidence,
and the toolchain, compatibility, snapshot/Docker, quality,
baseline-reproduction, and audit contracts. Preserve the documented btree
regression-manifest correction.

# Three Moves

Only if all decision evidence passes and the selection changes, use exact Go
1.26.7 and exact
`go get github.com/inconshreveable/mousetrap@<selected-version>` for one
dependency-only commit. Do not hand-edit module metadata and do not use tidy
as implementation. Preserve every retained selection, including HCL v1.0.0,
go-cmp v0.6.0, btree v1.1.3, gomarkdown fallback, regexp2 v1.12.0, Cobra
v1.10.2, YAML v3.0.5, go-md2man v2.0.7, Blackfriday v2.1.0, pflag v1.0.10,
Uniseg v0.4.7, Colorable v0.1.15, and go-isatty v0.0.20.

After a changed selection, re-run the complete P7 dependency gate: focused
behavior when loaded, graph/path, tests/race/vet, pinned lint, help/API/CLI,
launcher and Make contracts, preflight, host plus fresh snapshot/Docker meta
and acceptance, audit meta, focused and exact Q0-Q2 audits, separate full
audit, vulnerability comparison, empty-HOME count-2, and final
ordinary/ignored cleanliness. Exact make quality must exit 0 with all 27 rows
at L2 and zero held, regressed, not-comparable, or dirty counts. The full
audit may exit 1 only for established queued L3 rows, never 2.

Keep all caches, projections, reports, generated artifacts, build contexts,
schema-2 evidence, and audit output outside the worktree. Warm caches from a
separate external Git archive; never run `go mod download all` inside a
measured tree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next measured P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not
implement that next group, launch a successor, push, merge, publish, release,
stash, revert, delete retained evidence/images, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
