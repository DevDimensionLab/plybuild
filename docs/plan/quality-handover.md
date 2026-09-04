# Quality Upgrade Handover

Generated: 2026-09-04T17:19:14+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository And Continuity

- Worktree `/Users/perottochristensen/github/ply/upgrade-quality`, branch
  `codex/upgrade-quality`, base `master` at `5635d50`.
- Latest implementation is focused Markdown contract commit
  `8a6ef5fecd87cc773aed7934c9bff03fb229a8e8`, exact parent
  `1cf74aae6c7d678b0829c0522857ddbc05e4e7b6`, clean tree
  `b8d46e8eeb37162c467c5f9df7be861e6ead40e9`.
- The next implementation, if accepted, must follow the documentation handoff
  whose exact parent is `8a6ef5f`; start with empty ordinary and ignored
  status.
- No dependency metadata, production Go behavior, API/CLI, toolchain,
  Docker/release input, quality population, publisher, packaging, or P8 code
  changed in the rejected-latest session.
- No push, merge, publication, release, stash, successor launch, retained
  evidence/image deletion, or worktree removal occurred.

P2A-P6 are complete. P7 remains active after the Go 1.26.7 toolchain and the
completed dependency groups through regexp2 v1.12.0. Viper v1.16.0, Emoji
v2.2.14, ini v1.67.3, fatih/color v1.19.0, and fsnotify v1.10.1 remain rejected
for exceeding the retained Go 1.18 floor. Latest gomarkdown is now additionally
rejected for loaded renderer incompatibility. P8 remains queued.

The answered latest-gomarkdown archive links reciprocally to exactly one NEXT
archive for a separately measured pre-breaking gomarkdown fallback. Only the
launcher's mutable header and prompt regions may change. Never create
`.agent-task/current.md` or `.quality/manual-evidence.json`; keep all caches,
projections, reports, artifacts, build contexts, and audit evidence external.

## Latest Gomarkdown Decision

Rejected indirect `github.com/gomarkdown/markdown`
`v0.0.0-20221013030248-663e2500819c` ->
`v0.0.0-20260824154242-13c5cf49db8d`. Primary proxy evidence again reports no
tagged versions, resolves latest at 2026-08-24T15:42:42Z to unsigned commit
`13c5cf49db8d0bd189d6d04337e618ae01a83c8f`, declares Go 1.12, and yields
module/go.mod checksums `h1:8VtgBGEPLZ2Yn0Fuh6Pwmy3qF6indeaqy8mrBMbUKRQ=` and
`h1:JDGcbDT52eL4fju3sZ4TeHGsQwhG9nbDV21aMyhwPoA=`.

The exact external get changes one selection and the main-module edge only.
Both states retain 234 modules, 3,556 graph edges, 429 test packages, and four
loaded packages through
`plybuild/cmd -> go-term-markdown -> gomarkdown/markdown`. Exact get changes
only the indirect requirement and adds the new checksum pair while retaining
the old pair. Tidy projects 282 -> 285 lines. Main Go remains 1.18 and the
candidate declares Go 1.12.

General build, main tests, lint, help, API/CLI, gomarkdown self-tests, and
artifact meta-contracts pass. Vulnerability scans reproduce the intended
22/33/22 -> 20/30/20 improvement with no additions: GO-2023-2074 and
GO-2024-3205 leave both symbol scans, and those plus GO-2026-5208 leave the
module scan. The improvement is measured but not retained because the
candidate is behaviorally incompatible.

The decisive consumer replay ran all 34 go-term-markdown v0.1.4 `TestRender`
subtests independently. Retained gomarkdown passes 30; latest passes 25. Five
old-pass/candidate-fail cases—`Amps_and_angle_encoding`,
`Links,_reference_style`, `Links,_shortcut_references`,
`Literal_quotes_in_titles`, and `Markdown_Documentation_-_Basics`—panic with
`Unknown node type *ast.ReferenceDefinition`.

Latest commit `d18ffa2a47071de47ab9ad4dd61f0a01158ff97a` deliberately keeps reference
definitions as AST nodes for round trip. go-term-markdown v0.1.4's renderer
panics on unknown nodes. This loaded user-visible path makes latest
incompatible despite the general green checks.

Focused contract `TestTerminalMarkdownReferenceDefinitionsRemainRenderable`
was added in commit `8a6ef5f`. It runs the public consumer renderer and checks
the resolved reference text and URL. It has no direct `fatih/color` import and
no loop, preserving dependency classification, the 168-subject manual
iteration population, and the retained 282-line tidy projection.

## Accepted Evidence

- Decision replay `/private/tmp/ply-p7-gomarkdown-replay.1cf74aa.DA6bP5`:
  verified 51,314-entry manifest SHA-256
  `a1258b38e5ce00799d716b369b43e7d64f02d9b0562dba587e15b3662f7c02d6`;
  selection-summary SHA-256
  `ad6ab0f5d50768ba35d10df50b91e90062bb43241c34e85afedc74800967d12a`.
- Commit review `/private/tmp/ply-p7-gomarkdown-quality-review.8a6ef5f.SEwCxd`:
  verified 2,306-entry manifest SHA-256
  `dab3e71829e19cba1aaccca447e704052814da73623e5698f3c18e7ce591ad53`;
  manual-evidence SHA-256
  `0ebed9b0878ef52e8cb09792cdea008299e63ada4cd0b30a6ced2e4c3868c4c3`.
- Quality root `/private/tmp/ply-p7-gomarkdown-quality-parent.8a6ef5f.8w8yye`:
  verified 252,604-entry manifest SHA-256
  `2d72db874d9102a4e69d284b702c80d47cdb46b37a357774ed69c437243a2bfc`;
  Q0-Q2 scorecard SHA-256
  `d2c2c6e4201060e5eb037f6a10e2d73ea0df88413609f5c7d4a6a04298a8905c`.
  Exact `make quality` passes all 27 rows at L2, all 80 mutations, 8/8
  mutation groups, 4/4 host features, snapshot, and Docker, with zero held,
  regressed, not-comparable, or dirty counts.
- Regression root
  `/private/tmp/ply-p7-gomarkdown-regression-gate.8a6ef5f.lRiNgs`: verified
  463-entry manifest SHA-256
  `a000e23811cb85a9eff379537556af494940dd89c4142ca69d2971766ee0e682`;
  regression-summary SHA-256
  `34366be3008763d6ded65dcce0f78bc0518e00fa7f75953e26e43b2ad3cf7031`;
  full scorecard SHA-256
  `e46eede53c6b646c7f8a4a9182c92fd2b5ef0e414d1f4ca1c48ea060c7445b4e`.
  It binds retained 234/3,556/429/4 graph/path, 282-line tidy, 22/33/22
  vulnerabilities, and full-audit exit 1 only for queued Q3.1, Q3.3, Q3.4,
  and Q3.7. Its pre-existing Q3.4 heuristic is the sole full-audit ratchet
  regression; the P7 Q0-Q2 exit gate has zero regressions.

Tool identities remain Go 1.26.7
`9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6`,
golangci-lint 2.12.2
`3ba856c13833c4cda2182eb71bdbc7f96ddad339a1cdda8c36f88bfb1e34bd6f`,
GoReleaser 2.17.1
`f5f08a777bc1b9321fdebae0a03f6f20af3713c17d6089bf28061d7791ca578c`,
apidiff `0c55d9e385a57d6af9b69a9abd3b71d986123dbe82e687fcaecb99f9a0acba20`,
and govulncheck v1.7.0
`0db06024e71e82bd3e4444d57b757abf86fd4b7d0b08c1c21ca811597c953cd5`.

## Next Objective

Independently evaluate immediate pre-breaking gomarkdown commit
`45814d58469f9a462c899674ee58ca440b432dbf` as the highest known fallback.
Resolve its canonical proxy pseudoversion/checksums; prove exact one-selection
closure, Go 1.18 floor compatibility, applicable security fixes, and exclusion
of the reference-definition AST commit.

Run all 34 consumer cases and require every retained pass, especially the five
reference cases. Require the focused contract, gomarkdown self-tests, exact
graph/path/tidy/checksum behavior, full Go/quality/artifact gates, API/CLI/help,
and vulnerability comparison. Implement with exact `go get` only if every
decision gate passes; otherwise record rejection without dependency changes.

Warm caches from separate Git archives. Never use `go mod download all` in a
measured tree or tidy as implementation. Stop before another dependency,
consumer/source fix, language/toolchain or quality-policy change, packaging,
publication, or P8. Do not push, merge, publish, release, delete evidence,
stash, revert, launch a successor, or remove the worktree.
