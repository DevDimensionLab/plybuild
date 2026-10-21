# AL13 review return prose verification

Review evidence gaps and finding evidence now share the existing 256 KiB
return-file budget. Their substantive Unicode text, internal paragraphs and
tabs are preserved, with the existing outer-whitespace trimming. Summary,
identity, verdict, session and sealed-history rules remain in place. The parser
also refuses trailing data after the single review JSON object.

Development used cached Go 1.26.7, `GOTOOLCHAIN=local`,
`GOFLAGS=-mod=readonly`, existing caches and physical private scratch directories.
Raw evidence is retained at `/private/tmp/ply-al13-evidence`.

Before production changes, a CLI built from base
`c91e188da203eda4adb8307d49c68bb7bb95a1fa` reproduced:

- `red-prose.stdout/.stderr/.exit`: blocked multiline evidence refused with
  CLI exit 2, Go test exit 1.
- `red-finding-evidence.stdout/.stderr/.exit`: adjacent finding evidence refused
  with CLI exit 2, Go test exit 1.
- `red-trailing.stdout/.stderr/.exit`: second JSON object actually accepted
  with CLI exit 0; the behavioral regression exited 1. This was executed proof,
  rather than only source inspection.
- `red-parser.stdout/.stderr/.exit`: the focused parser unit regression exited 1
  on ordinary long single-line explanation before the fix.

`base-ply` retains that old CLI. `red-cli.stdout/.stderr/.exit` and `red-cli/`
repeat the behavioral failures against those same old product bytes using the
current CLI regressions and separate raw process files. These are fresh modeled
fixtures, not replay of any real runtime denial.

`green.stdout/.stderr/.exit` records Go test exit 0 for this anchored selection:

```text
^(TestAL13ReviewReturnProseParser|TestAL13ReviewReturnProseThroughCandidateCLI|TestAL13ReviewReturnRejectsInvalidCLIInput)$
```

The command was `go test ./internal/taskrun -run <selection> -count=1 -v`, with
`PLY_DELIVERY_TEST_BINARY=/private/tmp/ply-al13-evidence/candidate-ply` and
`PLY_AL13_EVIDENCE_DIR=/private/tmp/ply-al13-evidence/green-cli`.
`green-cli/` contains separate actual JSON stdout, stderr and exit files.
The initial green attempt (`green-initial.*`) had passing CLI cases but a parser
test expectation that overlooked existing outer-whitespace trimming. Correcting
that expectation required no further production change.

The three named tests cover long single-line and multiline Unicode prose;
adjacent finding evidence; all three verdicts; passed-with-findings and
changes-without-findings refusal; meaningful blocked evidence; control text;
schema/structure errors; malformed/trailing JSON; the unchanged total bound;
wrong pane/context/session; a short existing return; identical callback after
modeled response loss; refusal to change a sealed judgment; and byte-for-byte
preservation of the sealed input, assessment, context, request and Run state.
The CLI is a real executable. Herdr and provider sessions are the existing
synthetic fixtures, not a live provider pilot or human product acceptance.

`vet.stdout/.stderr/.exit` records exit 0 for
`go vet ./internal/taskrun ./cmd`. Shell syntax checks passed for the committed
`test/al13-acceptance.sh` and the frozen external acceptance entrypoint.

The acceptance entrypoint calls `test/al13-acceptance.sh`, which consumes
`PLY_CANDIDATE_BINARY` and verifies the exact clean candidate OID/tree. Its
five explicitly named tests add the existing short-return CLI case and
same-candidate evidence reassessment with immutable earlier history. It checks
that every named test actually ran, then runs vet and diff checks. The full
native acceptance is reserved for after independent exact-candidate review.
This document records development evidence, not a passed independent review,
qualified TaskResult, human QA or delivered integration.
