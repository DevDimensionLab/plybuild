# Agent Session: Evaluate Gomarkdown Compatible Fallback

Status: NEXT
Session ID: `2026-09-04T171914+0200-evaluate-gomarkdown-compatible-fallback`
Created: `2026-09-04T17:19:14+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8b536f26e14e1584a27636e31460a0b9cb5942b94e4ad99884cdea4305020081`
Previous: [2026-09-04T145138+0200-upgrade-gomarkdown-dependency.md](2026-09-04T145138+0200-upgrade-gomarkdown-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 by independently evaluating the immediate pre-breaking gomarkdown
commit `45814d58469f9a462c899674ee58ca440b432dbf` as the highest known fallback
below the rejected latest pseudoversion. Resolve its canonical Go
pseudoversion from primary evidence and, only if it preserves the loaded
go-term-markdown behavior, security improvement, retained Go 1.18 floor, and
every quality contract, implement it as one exact gomarkdown selection.

# Authorized Roadmap

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and the
completed dependency groups through regexp2 v1.12.0. Viper v1.16.0, Emoji
v2.2.14, ini v1.67.3, fatih/color v1.19.0, and fsnotify v1.10.1 remain rejected
for exceeding the Go 1.18 floor. Latest gomarkdown
`v0.0.0-20260824154242-13c5cf49db8d` is now rejected for loaded rendering
incompatibility. Do not revisit those decisions or combine another group. P8
remains queued.

The latest gomarkdown replay reproduced the expected exact one-selection
closure and 22/33/22 -> 20/30/20 security improvement, but five of 30
previously passing go-term-markdown v0.1.4 renderer cases panic with
`Unknown node type *ast.ReferenceDefinition`. Primary repository evidence
attributes the new AST nodes to unsigned commit
`d18ffa2a47071de47ab9ad4dd61f0a01158ff97a`, whose immediate parent is unsigned
commit `45814d58469f9a462c899674ee58ca440b432dbf` at
2026-08-24T15:13:36Z. Treat the parent only as a candidate until proxy,
checksum, source, graph, behavior, and vulnerability evidence independently
confirms it.

This session may change only the exact canonical gomarkdown fallback pin's
required go.mod/go.sum metadata and the roadmap/handoff record. The focused
executable contract already exists. Do not change production Go, the consumer
renderer, another dependency, the language/toolchain declaration, quality
apparatus, Docker/release inputs, packaging, publishers, or P8 code.

# Measurements At Start

Focused contract commit is `8a6ef5fecd87cc773aed7934c9bff03fb229a8e8`,
exact parent `1cf74aae6c7d678b0829c0522857ddbc05e4e7b6`, clean tree
`b8d46e8eeb37162c467c5f9df7be861e6ead40e9`. After this documentation handoff,
continuity HEAD must have exact parent `8a6ef5f`, and ordinary and ignored
status must be empty. The retained dependency remains
`github.com/gomarkdown/markdown v0.0.0-20221013030248-663e2500819c`.

Use exact Go 1.26.7 at
`/private/tmp/ply-p7-toolchain-go1.26.7.GGMf8j/sdk/go/bin/go`, SHA-256
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`.
Put its directory first in PATH, keep GOENV=off, GOWORK=off,
GOTOOLCHAIN=local, and inject no ambient GOFLAGS. Retain the already verified
golangci-lint 2.12.2, GoReleaser 2.17.1, apidiff, and govulncheck v1.7.0
binaries and hashes recorded in the handover.

The rejected-latest replay is
`/private/tmp/ply-p7-gomarkdown-replay.1cf74aa.DA6bP5`, with verified
51,314-entry manifest SHA-256
`a1258b38e5ce00799d716b369b43e7d64f02d9b0562dba587e15b3662f7c02d6`
and selection-summary SHA-256
`ad6ab0f5d50768ba35d10df50b91e90062bb43241c34e85afedc74800967d12a`.
It contains the exact five-case panic evidence and primary proxy, repository,
checksum, and Go vulnerability records.

Commit-bound schema-2 review evidence is
`/private/tmp/ply-p7-gomarkdown-quality-review.8a6ef5f.SEwCxd`, with verified
2,306-entry manifest SHA-256
`dab3e71829e19cba1aaccca447e704052814da73623e5698f3c18e7ce591ad53`
and manual-evidence SHA-256
`0ebed9b0878ef52e8cb09792cdea008299e63ada4cd0b30a6ced2e4c3868c4c3`.

The exact accepted quality root is
`/private/tmp/ply-p7-gomarkdown-quality-parent.8a6ef5f.8w8yye`; its verified
252,604-entry manifest SHA-256 is
`2d72db874d9102a4e69d284b702c80d47cdb46b37a357774ed69c437243a2bfc`.
Exact make quality passes all 27 Q0-Q2 rows at L2, 80/80 mutations, 8/8
mutation and 4/4 acceptance populations, with held, regressed,
not-comparable, and dirty counts zero. Q0-Q2 scorecard SHA-256 is
`d2c2c6e4201060e5eb037f6a10e2d73ea0df88413609f5c7d4a6a04298a8905c`.

The independent final regression root is
`/private/tmp/ply-p7-gomarkdown-regression-gate.8a6ef5f.lRiNgs`; its verified
463-entry manifest SHA-256 is
`a000e23811cb85a9eff379537556af494940dd89c4142ca69d2971766ee0e682`.
It binds 234 selected modules, 3,556 graph edges, 429 packages, four loaded
gomarkdown packages, the exact 282-line tidy projection, 22/33/22 retained
vulnerability counts, and the expected full-audit L3-only exit.

# Role And Boundaries

From fresh external archives and caches, resolve the full fallback commit to
its canonical proxy pseudoversion, time, Go declaration, checksum pair, and
signature status. Prove that it is an ancestor of the rejected latest, excludes
the `ReferenceDefinition` AST commit, and includes or excludes each of
GO-2023-2074, GO-2024-3205, and GO-2026-5208 using primary Go vulnerability
fixed-version data and repository history.

Measure old versus fallback selected modules, graph edges, complete package
population, checksums, loaded packages/path, explicit exact-get diff, and
`go mod tidy -diff`. Require exactly one changed module selection and no
unexplained checksum, graph, or package change. Run all 34 focused
go-term-markdown renderer cases individually and require every case that passes
on the retained selection to pass on the fallback, including the five
reference-definition cases. Also require the repository's focused executable
contract, gomarkdown self-tests, build, complete tests/race/vet, pinned lint,
byte-identical public help, and identical API/CLI reports.

Stop and record rejection without editing dependency metadata if canonical
resolution, Go-floor compatibility, exact closure, rendering parity, or the
expected security delta fails. Security improvement is required; determine the
exact fallback population rather than assuming it equals latest.

# Required Reading

Confirm branch, exact ancestry, empty ordinary and ignored status, reciprocal
archive links, launcher `--check`, and the P7/P8 checkpoint before editing.
Verify every accepted manifest. Read this archive, the rolling handover, the P7
roadmap, go.mod/go.sum, the rejected-latest graph/renderer/vulnerability
evidence, the focused contract, and the toolchain, compatibility,
snapshot/Docker, quality, baseline-reproduction, and audit contracts.

# Three Moves

Only if all decision evidence passes, use exact Go 1.26.7 and exact
`go get github.com/gomarkdown/markdown@<canonical-fallback-pseudoversion>` for
one dependency-only commit. Do not hand-edit module metadata and do not use
tidy as the implementation command. Preserve all retained dependency versions,
including regexp2 v1.12.0, Cobra v1.10.2, YAML v3.0.5, go-md2man v2.0.7,
Blackfriday v2.1.0, pflag v1.0.10, Uniseg v0.4.7, Colorable v0.1.15, and
go-isatty v0.0.20.

Re-run the complete P7 dependency gate: focused Markdown behavior, graph/path,
tests/race/vet, pinned lint, help/API/CLI, launcher and Make contracts,
preflight, host plus fresh snapshot/Docker meta and acceptance, audit meta,
focused and exact Q0-Q2 audits, separate full audit, vulnerability comparison,
empty-HOME count-2, and final ordinary/ignored cleanliness. Exact make quality
must exit 0 with all 27 rows at L2 and zero held, regressed, not-comparable, or
dirty counts. The full audit may exit 1 only for the established queued L3
rows, never 2.

Keep all caches, projections, reports, generated artifacts, build contexts,
schema-2 evidence, and audit output outside the worktree. Warm caches from a
separate external Git archive; never run `go mod download all` inside a
measured tree. Never create `.agent-task/current.md` or
`.quality/manual-evidence.json`.

# Automatic Handoff

After the decision, rewrite the rolling handover and roadmap, answer this
archive, create exactly one reciprocal NEXT archive for the next measured P7
group, replace only launcher mutable regions, run launcher/handoff contracts,
and make the normal `docs: prepare next agent session` commit. Do not implement
that next group, launch a successor, push, merge, publish, release, stash,
revert, delete retained evidence/images, or remove the worktree.
<!-- CODEX_SESSION_PROMPT_END -->
