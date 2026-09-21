# Agent Session: Decide Kisielk Errcheck Product Direction

Status: NEXT
Session ID: `2026-09-21T200144+0200-decide-kisielk-errcheck-product-direction`
Created: `2026-09-21T20:01:44+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `13bdf0c4add170fc6a4a44d97005d72da74acc5fd6d1b20a5c71de239cd0bd7f`
Previous: [2026-09-21T185250+0200-evaluate-kisielk-errcheck-dependency.md](2026-09-21T185250+0200-evaluate-kisielk-errcheck-dependency.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P7 only by making one bounded product decision for selected exact-path
`github.com/kisielk/errcheck v1.5.0`. The completed independent evaluation
found that no exact-path stable release qualifies: selected v1.5.0 fails
ordinary loader/tests and resource ownership, highest Go-1.18-compatible
v1.8.0 fails exact Go 1.26.7 and repeatability, and v1.9.0+ raise the declared
Go floor. Choose exactly one direction below, record its ownership and expiry
guards, and stop. Do not repeat completed fixtures, add a direct target root,
implement a parent change, evaluate another dependency group, reopen an
earlier decision, or begin P8.

# Defensive Scope

This is a documentation-only product decision for an ordinary dependency-
quality review. Reuse completed public metadata, static source/API inspection,
admissible upstream tests, bounded ordinary command fixtures, project graph/
build measurements, and public advisory evidence. Do not fuzz, stress, probe
resource exhaustion, generate oversized or deeply nested input, generate
adversarial malformed source, or perform security or exploitability analysis.
All completed loader, repeatability, resource-close, and command results are
final; do not reproduce them.

Every disposable archive, cache, tool, report, and fixture must remain beneath
`${CODEX_SESSION_SCRATCH_ROOT:?}`. Never write to `/private/tmp`, `/tmp`, a
sibling of the managed root, or another external root. Verify containment and
remove task-owned scratch evidence before handoff.

# Authorized Roadmap

P2A-P6 are complete. P7 is blocked only on this errcheck decision after exact
Go 1.26.7, every accepted dependency move through Google UUID v1.4.0,
qualified go-cleanhttp v0.5.2, and all target-specific retained-module
decisions through exact inherited unloaded httprouter v1.2.0. Every earlier
outcome is final under its own guards. P8 remains queued.

The evaluation left product source, `go.mod`, and `go.sum` unchanged. MVS
selects exact v1.5.0 only through
`gogo/protobuf@v1.3.2 -> errcheck@v1.5.0`, on the genuine route from direct,
imported, loaded Viper v1.15.0. Errcheck and gogo/protobuf have negative why or
zero repository imports, zero production/complete-test package loads, and no
runtime reachability. Physical selection and zero target loading are not
qualification or risk acceptance. No httprouter, jtolds/gls, go-junit-report,
json-iterator, clockwork, demangle, strcase, memberlist, or other exception
transfers.

# Completed Evaluation

Fresh exact-path evidence resolves the public active unarchived non-fork MIT
repository and 16 stable releases v1.0.0-v1.20.0 plus excluded
v1.5.0-alpha. Lightweight tags form one ancestry; proxy/Git/sumdb identities
agree. There is no redirect, retraction, deprecation, replacement, `/v2`
module, fork promotion, or later stable. Current master is one workflow-only
commit after v1.20.0.

V1.0.0-v1.2.0 have no Go directive, v1.3.0-v1.6.3 declare Go 1.14,
v1.7.0/v1.8.0 declare Go 1.18, v1.9.0/v1.10.0 declare Go 1.22.0, and v1.20.0
declares Go 1.25.0. Selected v1.5.0 resolves ten modules all at or below Go
1.14; v1.8.0 resolves six modules all at or below Go 1.18. Go 1.18.10 rejects
v1.9.0+ at module parsing, making v1.8.0 the highest floor-compatible release.

Selected v1.5.0 exports checker, exclusions, results/errors, no-files error,
and mutable defaults; v1.6.1+ also exports the go/analysis Analyzer and
ReadExcludes. The command covers ordinary error-result checking, assertions,
blank assignments, exclusion/ignore filters, test/generated files, tags,
module mode, path rendering, diagnostics, and 0/1/2 exits. Caller-owned
checker/exclusion collections and analyzer/default globals require external
synchronization. Result uniqueness copies before sorting.

Selected v1.5.0 fails its Go 1.18.10 library test because its historical
go/packages loader reports missing type data; exact Go 1.26.7 loading panics/
fails against current type data, including ordinary command analysis of its
own testdata. Static inspection also finds `readfile` never closes source
files; v1.6.1 adds the missing close. No exhaustion test was performed.
V1.6.1-v1.8.0 pass count one on Go 1.18.10, but the v1.6 line fails the
exact-Go old loader and v1.7/v1.8 fail exact-Go compilation in old x/tools.
V1.9.0+ pass exact-Go count one but are floor-ineligible.

Every v1.6.1-and-later serious candidate fails count-ten and race-count-ten
because an upstream test stores a `t.TempDir()` loader closure in package-
global state and never restores it; later iterations deterministically use a
deleted directory. The race detector reports no data race, but repeatability
fails through latest. Ordinary v1.8.0/Go1.18 and v1.9.0/exact-Go command
fixtures otherwise produce identical deterministic checking/filtering/
diagnostic behavior. All serious candidate commands cross-build for six
ordinary targets on compatible SDKs.

A direct selected-v1.5.0 get manufactures three roots, changes no selection,
and tidy restores the common baseline. A v1.8.0 root also moves six unrelated
x/* selections; tidy removes errcheck but retains x/net/x/text churn and does
not converge. A v1.9.0 root raises the project Go line to 1.22.0. None is an
authorized exact dependency-only move. All projections keep zero target load.

Exact/package OSV and GitHub target results are empty. Pinned govulncheck
v1.8.0 finds no errcheck package/symbol trace, but module scans find public
GO-2026-6180 and GO-2026-6179 in every serious candidate's x/mod closure
through latest v1.20.0. No exploitability work was performed. Base and direct
v1.5.0 project populations are identical at 30/22/20/20 without a target
trace; guard OSV retains only the recorded Gorilla/retryablehttp pairs. The Go
module index and memberlist CNA identities remain byte-exact.

The unchanged project retains 234 modules, 3,599 graph edges, 355 production
and 429 complete-test entries, 197 module-backed complete-test entries across
41 modules, 1,067 sums, and the 432-line tidy projection. Real `go.mod`/
`go.sum` hashes remain
`7255a37243bc8dd4601ef065e5985b9ffc0b7a43346035a869568f99679ec874` /
`87c9efb4baa70c3cc37ba91c8833d06d078c7e168b607eaa12564aa442b886a7`.
All 32 target-plus-earlier why results and guarded imports/loads are zero;
including errcheck gives 32 guarded selections/217 incoming edges at
`500c57a1b4ffbdf5bb1296ef4dcb1bd16dbae12e54b86664fd1f8a2481a6c0dd`.
Accepted quality remains 27/27 Q0-Q2 PASS at L2.

# Required Product Decision

Choose exactly one. Option 1 is recommended because the target and its owner
are completely unloaded, retaining the inherited selection preserves every
genuine request and earlier guard, and no stable exact-path release satisfies
the simultaneous Go-floor, exact-Go, repeatability, resource, closure-
advisory, and target-only-selection contracts:

1. Explicitly retain exact selected, inherited, unloaded v1.5.0 without
   metadata changes under an errcheck-specific, non-transferable exception.
   Accept only the completed selected loader/test and source-close failures,
   later-release floor/exact-Go/repeatability/closure-advisory blockers,
   repository/archive/release/module/API/command, ownership/global-state,
   MVS/loading, project, vulnerability, and related findings. Guard exact
   v1.5.0; the sole gogo/protobuf v1.3.2 request and genuine Viper owner route;
   negative why/import/load results; runtime unreachability; exact graph/
   module/tidy/Go-floor state; every earlier guard; and no new advisory,
   independent defect, release, owner, qualified release, supported tidy-
   stable owner, or compatible genuine route to a qualified release.
2. Authorize exactly one separate measurement-only gogo/protobuf owning-
   parent/request study. It may resolve whether genuine compatible Viper,
   etcd/api, or other already-recorded owner changes remove errcheck or select
   a qualified exact-path stable release, and measure every API, behavior,
   Go-floor, graph, guarded-selection, vulnerability, and project consequence.
   It may recommend a later route but may not implement one, combine parent
   upgrades, add a direct errcheck root, or alter a finalized guard without its
   fresh owning decision.
3. Stop P7 explicitly with selected v1.5.0 unresolved and unaccepted. Prepare
   no implementation or dependency-evaluation successor and do not begin P8.

If none is acceptable, choose option 3. Do not infer acceptance from zero
loading, describe any release as qualified, add a direct target root, promote
an alternate path/branch/pseudo-version, implement a parent move, or transfer
another dependency's exception.

# Required Reading And Moves

Read this archive and its answered evaluation, the answered httprouter,
jtolds/gls, go-junit-report, and gogo/protobuf records, relevant retained-
module/ownership records, rolling handover, roadmap, `go.mod`, and `go.sum`.
Reuse the completed evaluation; revalidate only exact decision guards and
fresh advisory identities needed for the choice.

Record exactly one numbered choice and why. If option 1 is selected, define
the exact target-specific acceptance and expiry guards, leave source and
dependency metadata unchanged, and prepare one bounded evaluation of the next
queued P7 exact-path group without executing it. If option 2 is selected,
prepare exactly that one measurement-only study as the successor. If option 3
is selected, record the unresolved P7 stop and prepare no implementation or
dependency-evaluation successor.

In every case update the roadmap and rolling handover, answer this archive,
prepare any single reciprocal successor required by the selected option,
verify scratch containment and launcher state, and make only the local
`docs: prepare next agent session` handoff commit.

# Automatic Handoff

Do not execute a successor, push, merge, publish, release, stash, revert,
bypass cleanup, remove the worktree, change product source or dependency
metadata, add a direct errcheck root, implement a parent move, repeat completed
fixtures, evaluate another dependency in this decision turn, reopen
httprouter/jtolds/gls/go-junit-report/json-iterator/clockwork/mvn-pom-mutator/
demangle/pprof/strcase/memberlist work, write outside the managed scratch root,
or begin P8.
<!-- CODEX_SESSION_PROMPT_END -->
