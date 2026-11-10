# WI01 compatibility expectations

These fixed bytes were captured and verified against the pre-extraction source
at `450702c8a486c037bcf0f3d3c69646d41135fbf1`, before moving production types.
There is no runtime fixture updater: tests read the committed expectations.

`wi01Models` supplies explicit inputs covering every moved field, nil/bound Task
worktrees and nil/empty repo slices. models.json/yaml retain wire names, order
and null/array behavior. Non-empty format-1/2 registries preserve reencoding;
intent.txt preserves the canonical create-operation intent digest. raw-v*.txt
bind original registry bytes prefixed with the fixed comment in the test, so
normalizing/reencoding before hashing would fail.

Newer history/adoption uses the existing TestEpicBaseB1B2AndHistoricalSpec and
TestAdoptionExistingSourceJourney fixtures. Full/scoped closure and read-count
regressions remain in the workspace/workspaceview suites.
