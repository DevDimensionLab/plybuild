# Quality audit — <project> — <date>

**Commit:** <sha>  **Auditor:** <who>  **Tool version:** <n>

> Produced with `tools/quality-audit.sh`, plus the four manual findings no script can do.
> Criteria: ACCEPTANCE-CRITERIA.md

## Denominators

<Paste the denominator block from the scorecard. Every number below is meaningless without it.>

## Score

| Level | Passed | Failed | Unmeasurable |
|---|---|---|---|
| L0 Measurable | | | |
| L1 Tested | | | |
| L2 Measured | | | |
| L3 Self-defending | | | |

**Attained level:** <L0 / L1 / L2 / L3 / none>

`UNMEASURABLE` counts as not-passed. It is not a softer FAIL — it means the check could not have
failed, so it is not evidence either way.

## Findings

<One row per non-PASS criterion, with what was measured. Keep the audit's own wording; it states
the measurement rather than a claim.>

| Id | Verdict | Measured | Lift |
|---|---|---|---|
| | | | |

## Manual findings

<Required for L3, and worth doing at every level. Numbers, not adjectives.>

**1. Seam test.** <Which calls, what was swapped, whether the suite went red.>

**2. Real output, read as a user.** <What was run, and what the output actually says. Does it
answer the question the feature exists for?>

**3. Three survivors, classified.** <Missing test / redundant code / fixture that never gets
there — which, and why.>

**4. One premise verified.** <Which claim in the project's own documentation, and whether the code
agrees.>

## Lift plan

<Ordered by level, then by cost. For each move: the change, and the measurement that will show it
landed. Name what each failing criterion blocks.>

1. **<move>** — blocks <criteria>. Done when: <measurement>.
2. **<move>** — Done when: <measurement>.
3. **<move>** — Done when: <measurement>.

## Baseline

<Store the scorecard as the ratchet baseline, and say where:>
`target/quality-audit/scorecard.md` at commit <sha>.
