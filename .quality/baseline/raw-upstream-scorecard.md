Quality audit — /private/tmp/ply-baseline-git.PNojvw/repo
commit 5635d50 · tool version 1 · full mode

Denominators
  packages 20 · with tests 5 · go files 100 · test functions 33 (2 skipped)
  table-driven tests 0 · scripts 1 · mutation harnesses 0 · acceptance scripts 0
  enabled CI workflows 0 · inventory present

Findings
  PASS         Q0.1  one command builds every package
               measured: go build ./... clean over 20 packages
               (no negative probe for this check)
  FAIL         Q0.2  one command runs every test, uncached
               measured: go test ./... -count=1 exited 1; first failure: --- FAIL: TestSortAndWritePom_sort_enabled_by_default (0.00s) in github.com/devdimensionlab/plybuild/pkg/config
               (no negative probe for this check)
  FAIL         Q0.3  the static gate is clean
               measured: go vet clean, but there is no checked-in linter config
               lift:     ACCEPTANCE-CRITERIA.md Q0.3
               (no negative probe for this check)
  UNMEASURABLE Q0.4  tests are hermetic
               measured: the ordinary suite is not green, so an isolated run proves nothing
               (no negative probe for this check)
  UNMEASURABLE Q0.5  the suite leaves the tree byte-identical
               measured: the suite did not run green, so an unchanged tree proves nothing
               (no negative probe for this check)
  FAIL         Q0.6  no test writes into the repository; no permanently skipped tests
               measured: leak guard defined 0 times, called 0 times; 2 skipped tests of 33
               lift:     04-seams-and-test-doubles.md §6
  PASS         Q0.7  the denominators are recorded
               measured: 20 packages, 100 go files, 33 test functions, 1 scripts, 0 harnesses, 0 acceptance scripts
               (no negative probe for this check)
  FAIL         Q0.8  every script has a meta-test and parses
               measured: 1 of 1 scripts have no scripts/test-<name>
               lift:     01-mutation-harness.md §11
  FAIL         Q1.1  no package with logic has zero tests
               measured: 15 of 20 packages have no test files (RATCHET)
               lift:     04-seams-and-test-doubles.md §10
               (no negative probe for this check)
  FAIL         Q1.2  the entry layer is testable
               measured: 127 process-exiting calls outside main() - the layer cannot be driven from a test at all (RATCHET)
               lift:     04-seams-and-test-doubles.md §10
  FAIL         Q1.3  external effects go through an injection point
               measured: 9 files make direct calls outside the declared adapters (RATCHET)
               lift:     04-seams-and-test-doubles.md §2
  FAIL         Q1.4  every seam has an argument-swap test
               measured: 0 of 8 declared seams covered
               lift:     04-seams-and-test-doubles.md §1
               (no negative probe for this check)
  FAIL         Q1.5  time is injected
               measured: 1 direct time.Now() calls, 0 clock declarations
               lift:     04-seams-and-test-doubles.md §7
               (no negative probe for this check)
  UNMEASURABLE Q1.6  doubles record arguments; defaults for every side effect
               measured: MANUAL: read the test harness. Does every dependency get a default double, and is the struct passed on whole rather than rebuilt field by field?
               lift:     04-seams-and-test-doubles.md §4
               (no negative probe for this check)
  UNMEASURABLE Q1.7  partial-failure commands are asserted on content
               measured: MANUAL: does any test assert only the exit code of a command that reports per-item failures?
               lift:     04-seams-and-test-doubles.md §9
               (no negative probe for this check)
  UNMEASURABLE Q1.9  tests refuse to pass on an empty population
               measured: MANUAL: read the tests that iterate a collection. Is there an assertion - not a comment - that fails when the collection is empty?
               lift:     02-reading-a-surviving-mutant.md §2
               (no negative probe for this check)
  UNMEASURABLE Q1.8  no state leaks between runs
               measured: the ordinary suite is not green, so a second run proves nothing
               (no negative probe for this check)
  PASS         Q2.1  every declared subject has a mutation harness
               measured: 0 harnesses for 8 declared subjects of 20 packages
               (no negative probe for this check)
  UNMEASURABLE Q2.2  each harness declares enough mutations and can fail
               measured: no harnesses to measure
               lift:     01-mutation-harness.md §2
               (no negative probe for this check)
  UNMEASURABLE Q2.3  each harness has a T1-T10 meta-test
               measured: no harnesses to measure
               lift:     01-mutation-harness.md §8
               (no negative probe for this check)
  UNMEASURABLE Q2.4  declared equals killed, measured by running it
               measured: not run by this audit: a harness run mutates source and can take many minutes. Run each scripts/test-mutate-* yourself and record the numbers.
               lift:     01-mutation-harness.md §8
               (no negative probe for this check)
  FAIL         Q2.5  each user-facing feature has an acceptance script
               measured: 8 missing scripts across 4 features
               lift:     03-acceptance-scripts.md
               (no negative probe for this check)
  UNMEASURABLE Q2.6  acceptance-script properties
               measured: no acceptance scripts to inspect
               lift:     03-acceptance-scripts.md
               (no negative probe for this check)
  UNMEASURABLE Q2.7  acceptance-script properties
               measured: no acceptance scripts to inspect
               lift:     03-acceptance-scripts.md
               (no negative probe for this check)
  UNMEASURABLE Q2.8  acceptance-script properties
               measured: no acceptance scripts to inspect
               lift:     03-acceptance-scripts.md
               (no negative probe for this check)
  UNMEASURABLE Q2.9  acceptance-script properties
               measured: no acceptance scripts to inspect
               lift:     03-acceptance-scripts.md
               (no negative probe for this check)
  UNMEASURABLE Q2.10 acceptance-script properties
               measured: no acceptance scripts to inspect
               lift:     03-acceptance-scripts.md
  FAIL         Q3.1  the rules file is numbered and names its checks
               measured: no AGENTS.md, CLAUDE.md, CONTRIBUTING.md or rules document found
               lift:     06-documentation-as-apparatus.md §5
               (no negative probe for this check)
  FAIL         Q3.2  the overview is tested against the real command tree
               measured: no test references the README
               lift:     06-documentation-as-apparatus.md §4
               (no negative probe for this check)
  UNMEASURABLE Q3.3  entry documents are under test, order-independently
               measured: MANUAL: is there a test on the entry/boot document, and are its assertions independent of incidental order? An assertion that fails early makes every assertion below it unreachable.
               lift:     06-documentation-as-apparatus.md §4
               (no negative probe for this check)
  UNMEASURABLE Q3.4  claim rot is bounded
               measured: 0 state-claim phrases in markdown (RATCHET, indicator only - the phrase list is a heuristic and proves nothing on its own)
               lift:     06-documentation-as-apparatus.md §3
  FAIL         Q3.5  features have a design note with a premises table
               measured: no docs/design/ notes
               lift:     06-documentation-as-apparatus.md §1
               (no negative probe for this check)
  FAIL         Q3.6  a rolling handover exists
               measured: none found
               lift:     06-documentation-as-apparatus.md §2
               (no negative probe for this check)
  UNMEASURABLE Q3.7  no drift between repository and installed instruction copies
               measured: MANUAL: cmp -s each instruction file in the repository against its installed copy, BOTH directions, and check what the install step actually iterates over.
               lift:     07-agent-harness.md §3
               (no negative probe for this check)
  FAIL         Q3.8  one entry point runs the apparatus
               measured: no single make target runs tests, meta-tests and acceptance scripts
               lift:     01-mutation-harness.md §11
               (no negative probe for this check)
  UNMEASURABLE Q3.9  the score is reproducible and a human has read the output
               measured: MANUAL, and deliberately so: the top level cannot be reached by a script alone. Record the four manual findings from ACCEPTANCE-CRITERIA.md 'How to score', with numbers.
               lift:     05-evidence-discipline.md §4
               (no negative probe for this check)

Levels
  L0 2/8   PASS 2  FAIL 4  UNMEASURABLE 2
  L1 0/9   PASS 0  FAIL 5  UNMEASURABLE 4
  L2 1/10   PASS 1  FAIL 1  UNMEASURABLE 8
  L3 0/9   PASS 0  FAIL 5  UNMEASURABLE 4

ATTAINED LEVEL: none

Still required by hand (no script can do these — ACCEPTANCE-CRITERIA.md “How to score”)
  1. seam test on the three highest-value external calls
  2. read one real human output, with real numbers, as a user
  3. classify three survivors (or three untested branches) by the three causes
  4. verify one premise the project documents, against the code
