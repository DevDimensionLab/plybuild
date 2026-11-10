#!/usr/bin/env bash
# TEMPLATE-VERSION: 1
#
# Scores a repository against ACCEPTANCE-CRITERIA.md, in the directory beside this one.
#
# READ-ONLY with respect to the audited repository. Writes only to <repo>/target/quality-audit/.
#
#   quality-audit.sh [<repo>] [--quick] [--self-test] [--corpus <dir>]
#                    [--baseline <scorecard.json>] [--only Q1.2,Q2.*] [--out <dir>]
#
# --out writes the report somewhere other than <repo>/target/quality-audit -- necessary when
# the repository is read-only, and preferable when it must not be touched at all.
#
# Design rules this script must obey, because an audit that lies is worse than none:
#   - Four verdicts, never three: PASS / FAIL / UNMEASURABLE / N-A. A check whose population
#     is empty, or which could not run, is UNMEASURABLE -- never PASS.
#   - Every row prints WHAT WAS MEASURED, not just a verdict.
#   - Exit 1 means the project failed. Exit 2 means THIS SCRIPT is broken (self-test failed,
#     tree modified, timeout). Conflating them is the same error as reading a build failure
#     as a killed mutant.
#   - Negative probes: a check that has never been seen to fail is not evidence.
#
# CAUTION: this runs the audited project's own test suite, so it inherits whatever that suite
# does -- including reaching the network or writing files. That is not a flaw in the audit; it
# is what criterion Q0.4 measures. Run it somewhere you are willing to have those tests run.
set -uo pipefail

VERSION="1"
QUICK="no"
SELF_TEST_ONLY="no"
CORPUS_DIR=""
BASELINE=""
ONLY=""
OUT_OVERRIDE=""
REPO="."

while [ $# -gt 0 ]; do
  case "$1" in
    --quick)      QUICK="yes"; shift ;;
    --self-test)  SELF_TEST_ONLY="yes"; shift ;;
    --corpus)     CORPUS_DIR="${2:-}"; shift 2 ;;
    --baseline)   BASELINE="${2:-}"; shift 2 ;;
    --only)       ONLY="${2:-}"; shift 2 ;;
    --out)        OUT_OVERRIDE="${2:-}"; shift 2 ;;
    -h|--help)    sed -n '2,20p' "$0"; exit 0 ;;
    -*)           printf 'unknown option: %s\n' "$1" >&2; exit 2 ;;
    *)            REPO="$1"; shift ;;
  esac
done

TMP_BASE="${TMPDIR:-/tmp}"
WORK="$(mktemp -d "${TMP_BASE%/}/quality-audit.XXXXXX")"
trap 'rm -rf "$WORK"' EXIT

# ---------------------------------------------------------------------------------------------
# Result recording
# ---------------------------------------------------------------------------------------------
PASS_N=0; FAIL_N=0; UNMEAS_N=0; NA_N=0
ROWS=""            # id|level|verdict|title|measured|lift
PROBED=""          # ids that have a negative probe
NO_PROBE=""

record() {                       # record <id> <level> <verdict> <title> <measured> <lift>
  local id="$1" level="$2" verdict="$3" title="$4" measured="$5" lift="${6:-}"
  case "$verdict" in
    PASS)         PASS_N=$((PASS_N + 1)) ;;
    FAIL)         FAIL_N=$((FAIL_N + 1)) ;;
    UNMEASURABLE) UNMEAS_N=$((UNMEAS_N + 1)) ;;
    N-A)          NA_N=$((NA_N + 1)) ;;
  esac
  ROWS="${ROWS}${id}|${level}|${verdict}|${title}|${measured}|${lift}
"
}

wanted() {                       # honour --only
  [ -z "$ONLY" ] && return 0
  local id="$1" pattern
  local IFS=,
  for pattern in $ONLY; do
    case "$id" in $pattern) return 0 ;; esac
  done
  return 1
}

# ---------------------------------------------------------------------------------------------
# Predicates. Each takes a directory and prints a number or a string, so a NEGATIVE PROBE can
# call it against a synthetic fixture and require the answer to change.
# ---------------------------------------------------------------------------------------------

# Comment-stripped view of a shell file. An absence grep that skips this reads its own
# justification as a violation -- measured, twice.
code_only() { grep -v '^[[:space:]]*#' "$1" 2>/dev/null; }

p_fatal_exits() {                # process-exiting calls outside main()
  local dir="$1"
  grep -rn --include='*.go' 'log\.Fatal\|os\.Exit(' "$dir" 2>/dev/null \
    | grep -v '_test\.go' | grep -vc '/main\.go:' || true
}

p_direct_external() {            # external calls made outside the declared adapter packages
  local dir="$1" adapters="$2" out
  out=$(grep -rln --include='*.go' \
        'http\.Get(\|http\.Post(\|http\.NewRequest(\|exec\.Command(\|time\.Now()' \
        "$dir" 2>/dev/null | grep -v '_test\.go' || true)
  if [ -n "$adapters" ]; then
    out=$(printf '%s\n' "$out" | grep -v -F -f "$adapters" || true)
  fi
  printf '%s\n' "$out" | grep -c . || true
}

p_skipped_tests() {
  grep -rn --include='*_test.go' 't\.Skip\|t\.Skipf' "$1" 2>/dev/null | grep -c . || true
}

p_test_funcs() {
  grep -rh --include='*_test.go' '^func Test' "$1" 2>/dev/null | grep -c . || true
}

p_table_tests() {
  grep -rn --include='*_test.go' 't\.Run(' "$1" 2>/dev/null | grep -c . || true
}

p_enabled_ci() {                 # workflows that are ENABLED -- .deactivated does not count
  local dir="$1"
  ls "$dir"/.github/workflows/*.yml "$dir"/.github/workflows/*.yaml 2>/dev/null | grep -c . || true
}

p_claim_rot() {                  # state claims in docs without an adjacent measurement command
  local dir="$1"
  grep -rniE 'currently|at the moment|is not (yet )?(built|installed|implemented)|commits ahead|needs (to be )?install' \
    --include='*.md' "$dir" 2>/dev/null | grep -c . || true
}

p_scripts_without_metatest() {
  local dir="$1" s base n=0
  [ -d "$dir/scripts" ] || { printf '0\n'; return; }
  for s in "$dir"/scripts/*; do
    [ -f "$s" ] || continue
    base="$(basename "$s")"
    case "$base" in test-*) continue ;; esac
    [ -f "$dir/scripts/test-$base" ] || n=$((n + 1))
  done
  printf '%s\n' "$n"
}

# ---------------------------------------------------------------------------------------------
# Negative probes. A check that has never been seen to fail is not evidence.
# Each builds a synthetic fixture in $WORK on which the predicate MUST report a violation.
# Probes here cover the text-based predicates; the toolchain-based checks are covered by
# tools/test-quality-audit.sh, which runs the whole audit against a deliberately bad repo.
# ---------------------------------------------------------------------------------------------
PROBE_FAILURES=""

probe() {                        # probe <id> <expected-nonzero-expression...>
  local id="$1"; shift
  if "$@"; then
    PROBED="$PROBED $id"
  else
    PROBE_FAILURES="$PROBE_FAILURES $id"
  fi
}

probe_fatal_exits() {
  local d="$WORK/probe-fatal"; mkdir -p "$d"
  printf 'package x\nimport "log"\nfunc Run() { log.Fatal("boom") }\n' > "$d/x.go"
  [ "$(p_fatal_exits "$d")" -gt 0 ]
}
probe_skipped_tests() {
  local d="$WORK/probe-skip"; mkdir -p "$d"
  printf 'package x\nimport "testing"\nfunc TestA(t *testing.T) { t.Skip("no") }\n' > "$d/x_test.go"
  [ "$(p_skipped_tests "$d")" -gt 0 ]
}
probe_direct_external() {
  local d="$WORK/probe-ext"; mkdir -p "$d"
  printf 'package x\nimport "os/exec"\nfunc Run() { _ = exec.Command("ls") }\n' > "$d/x.go"
  [ "$(p_direct_external "$d" "")" -gt 0 ]
}
probe_enabled_ci() {             # a .deactivated workflow must NOT count as CI
  local d="$WORK/probe-ci"; mkdir -p "$d/.github/workflows"
  printf 'on: push\n' > "$d/.github/workflows/lint.yaml.deactivated"
  [ "$(p_enabled_ci "$d")" -eq 0 ]
}
probe_scripts_without_metatest() {
  local d="$WORK/probe-scripts"; mkdir -p "$d/scripts"
  printf '#!/bin/sh\n' > "$d/scripts/lonely"
  [ "$(p_scripts_without_metatest "$d")" -gt 0 ]
}
probe_code_only() {              # an absence grep must not read its own comment as a hit
  local f="$WORK/probe-comment.sh"
  printf '#!/usr/bin/env bash\n# we deliberately do not use --since here\necho hi\n' > "$f"
  [ "$(code_only "$f" | grep -c -- '--since')" -eq 0 ]
}
probe_claim_rot() {
  local d="$WORK/probe-rot"; mkdir -p "$d"
  printf 'The binary is not yet installed.\n' > "$d/note.md"
  [ "$(p_claim_rot "$d")" -gt 0 ]
}

run_probes() {
  probe Q1.2  probe_fatal_exits
  probe Q0.6  probe_skipped_tests
  probe Q1.3  probe_direct_external
  probe L0-CI probe_enabled_ci
  probe Q0.8  probe_scripts_without_metatest
  probe Q2.10 probe_code_only
  probe Q3.4  probe_claim_rot
}

# ---------------------------------------------------------------------------------------------
# --corpus: audit this methodology library against its own invariants
# ---------------------------------------------------------------------------------------------
if [ -n "$CORPUS_DIR" ]; then
  # Resolve to an absolute path. With a relative one, `dirname` produces "." and the fallback
  # term-list lookup silently finds nothing -- so the neutrality check reports UNMEASURABLE
  # while looking like it ran. Measured while probing this very check.
  CORPUS_DIR="$( cd "$CORPUS_DIR" 2>/dev/null && pwd )" \
    || { printf 'no such directory for --corpus\n' >&2; exit 2; }
  printf 'Corpus audit of %s\n\n' "$CORPUS_DIR"
  corpus_fail=0
  # C1: no document may name a specific product, organisation, repository or person.
  # The term list lives OUTSIDE the corpus by default: a shareable methodology directory
  # must not itself contain the list of names it is scrubbed of.
  TERMS="$CORPUS_DIR/.forbidden-terms"
  [ -f "$TERMS" ] || TERMS="$(dirname "$CORPUS_DIR")/.quality-forbidden-terms"
  if [ -f "$TERMS" ]; then
    while IFS= read -r term; do
      [ -z "$term" ] && continue
      case "$term" in \#*) continue ;; esac
      # -w so a short name does not match inside an ordinary word.
      hits=$(grep -rniwF --include='*.md' --include='*.sh' -- "$term" "$CORPUS_DIR" 2>/dev/null | grep -c . || true)
      if [ "$hits" -ne 0 ]; then
        printf 'FAIL  C1 neutrality: "%s" appears %s times\n' "$term" "$hits"
        grep -rniwF --include='*.md' --include='*.sh' -- "$term" "$CORPUS_DIR" | head -5
        corpus_fail=1
      fi
    done < "$TERMS"
    [ "$corpus_fail" -eq 0 ] && printf 'PASS  C1 neutrality: no forbidden term appears\n'
  else
    printf 'UNMEASURABLE  C1 neutrality: no term list found, so nothing was checked.\n'
    printf '              Put one at %s\n' "$(dirname "$CORPUS_DIR")/.quality-forbidden-terms"
    corpus_fail=1
  fi
  # C2: every script parses.
  c2=0
  for s in "$CORPUS_DIR"/tools/*.sh; do
    bash -n "$s" || { printf 'FAIL  C2 %s does not parse\n' "$s"; c2=1; corpus_fail=1; }
  done
  [ "$c2" -eq 0 ] && printf 'PASS  C2 every script parses\n'
  # C3: every document is referenced from the index.
  c3=0
  for d in "$CORPUS_DIR"/*.md; do
    b="$(basename "$d")"
    [ "$b" = "README.md" ] && continue
    grep -q "$b" "$CORPUS_DIR/README.md" || { printf 'FAIL  C3 %s is not in the index\n' "$b"; c3=1; corpus_fail=1; }
  done
  [ "$c3" -eq 0 ] && printf 'PASS  C3 every document is in the index\n'
  # C4: every cross-reference resolves.
  c4=0
  for d in "$CORPUS_DIR"/*.md; do
    for ref in $(grep -o '](\([0-9A-Za-z._/-]*\.md\))' "$d" | sed 's/](\(.*\))/\1/' | sort -u); do
      [ -f "$CORPUS_DIR/$ref" ] || { printf 'FAIL  C4 %s links to missing %s\n' "$(basename "$d")" "$ref"; c4=1; corpus_fail=1; }
    done
  done
  [ "$c4" -eq 0 ] && printf 'PASS  C4 every cross-reference resolves\n'
  # C5: the corpus must state no volatile project state.
  # Anchored to the start of the line: an unanchored match also hits the CHECK that forbids
  # the pattern, which is the same shape as a static check reading its own justification.
  if grep -rnE '^[[:space:]]*exec > >\(tee' "$CORPUS_DIR"/tools/*.sh > /dev/null 2>&1; then
    printf 'FAIL  C5 a template uses exec > >(tee ...), which can leave the log empty\n'; corpus_fail=1
  else
    printf 'PASS  C5 no template uses exec > >(tee ...)\n'
  fi
  # C6: the JVM document claims one row per criterion. Check it rather than trust it.
  if [ -f "$CORPUS_DIR/ACCEPTANCE-CRITERIA.md" ] && [ -f "$CORPUS_DIR/08-jvm-adaptation.md" ]; then
    c6=0
    for cid in $( grep -oE '^### (Q[0-9]\.[0-9]+)' "$CORPUS_DIR/ACCEPTANCE-CRITERIA.md" | awk '{print $2}' ); do
      grep -q "$cid" "$CORPUS_DIR/08-jvm-adaptation.md" \
        || { printf 'FAIL  C6 %s has no row in the JVM translation table\n' "$cid"; c6=1; corpus_fail=1; }
    done
    [ "$c6" -eq 0 ] && printf 'PASS  C6 every criterion has a row in the JVM translation table\n'
  fi
  printf '\n'
  [ "$corpus_fail" -eq 0 ] && { printf 'corpus OK\n'; exit 0; }
  printf 'corpus has findings\n'; exit 1
fi

# ---------------------------------------------------------------------------------------------
# Self-test
# ---------------------------------------------------------------------------------------------
run_probes
if [ -n "$PROBE_FAILURES" ]; then
  printf 'AUDIT BROKEN: these checks did not fail on a fixture built to break them:%s\n' "$PROBE_FAILURES" >&2
  printf 'A check that has never been seen to fail is not evidence.\n' >&2
  exit 2
fi
if [ "$SELF_TEST_ONLY" = "yes" ]; then
  printf 'self-test OK: %s negative probes fired\n' "$(printf '%s' "$PROBED" | wc -w | tr -d ' ')"
  exit 0
fi

# ---------------------------------------------------------------------------------------------
# Preflight
# ---------------------------------------------------------------------------------------------
[ -d "$REPO" ] || { printf 'no such directory: %s\n' "$REPO" >&2; exit 2; }
REPO="$(cd "$REPO" && pwd)"

OUT_DIR="${OUT_OVERRIDE:-$REPO/target/quality-audit}"
# LOAD-BEARING: a verdict that does not reach disk is an audit failure, not a project failure.
# Measured while dry-running this script: the target repository was read-only, the report was
# never written, and the audit still printed a complete score and exited 1 -- i.e. exactly the
# failure 03-acceptance-scripts.md warns about, in the tool that teaches it.
if ! mkdir -p "$OUT_DIR" 2>/dev/null || ! : > "$OUT_DIR/.write-probe" 2>/dev/null; then
  printf 'AUDIT BROKEN: cannot write the report to %s\n' "$OUT_DIR" >&2
  printf 'The verdict must reach disk. Use --out <writable-dir>.\n' >&2
  exit 2
fi
rm -f "$OUT_DIR/.write-probe"
REPORT="$OUT_DIR/scorecard.md"

if [ ! -f "$REPO/go.mod" ]; then
  printf 'This audit measures Go projects; %s has no go.mod.\n' "$REPO" >&2
  printf 'For a JVM project see 08-jvm-adaptation.md; the criteria are the same, the commands differ.\n' >&2
  exit 2
fi

# The audit writes its own report under target/, so that path is excluded from the
# comparison -- otherwise the audit would report ITSELF as having modified the tree.
tree_hash() {
  ( cd "$REPO" && git status --porcelain --untracked-files=all 2>/dev/null \
      | grep -v ' target/\|^?? target/' | shasum | awk '{print $1}' )
}
TREE_BEFORE=$( tree_hash )

# ---------------------------------------------------------------------------------------------
# Denominators (Q0.7) -- every number below is meaningless without these
# ---------------------------------------------------------------------------------------------
export GOCACHE="${TMP_BASE%/}/gocache-audit" GOTMPDIR="${TMP_BASE%/}" CGO_ENABLED=0

PKGS=$( cd "$REPO" && go list ./... 2>/dev/null | grep -c . || true )
PKGS_WITH_TESTS=$( cd "$REPO" && go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./... 2>/dev/null | grep -c . || true )
TESTFUNCS=$( p_test_funcs "$REPO" )
SKIPPED=$( p_skipped_tests "$REPO" )
# Count by BASENAME. Filtering the full path once excluded every script in a repository whose
# parent directory happened to contain "/test-" -- found by the audit's own test.
count_scripts() {
  local n=0 f base
  [ -d "$REPO/scripts" ] || { printf '0\n'; return; }
  for f in "$REPO"/scripts/*; do
    [ -f "$f" ] || continue
    base="$(basename "$f")"
    case "$base" in test-*) continue ;; esac
    n=$((n + 1))
  done
  printf '%s\n' "$n"
}
SCRIPTS=$( count_scripts )
HARNESSES=$( ls "$REPO"/scripts/mutate-* 2>/dev/null | grep -c . || true )
VERIFIERS=$( ls "$REPO"/scripts/verify-* 2>/dev/null | grep -c . || true )
GOFILES=$( find "$REPO" -name '*.go' -not -path '*/vendor/*' 2>/dev/null | grep -c . || true )

INVENTORY="$REPO/.quality/inventory"
adapters_file=""
if [ -f "$INVENTORY" ]; then
  adapters_file="$WORK/adapters"
  sed -n '/^\[adapters\]/,/^\[/p' "$INVENTORY" | grep -v '^\[' | grep -v '^[[:space:]]*$' > "$adapters_file" || true
  [ -s "$adapters_file" ] || adapters_file=""
fi
inv_section() { [ -f "$INVENTORY" ] || return 0; sed -n "/^\[$1\]/,/^\[/p" "$INVENTORY" | grep -v '^\[' | grep -v '^[[:space:]]*$' || true; }

# ---------------------------------------------------------------------------------------------
# L0
# ---------------------------------------------------------------------------------------------
if wanted Q0.1; then
  build_out=$( cd "$REPO" && go build ./... 2>&1 ); build_rc=$?
  if [ "$build_rc" -eq 0 ]; then
    record Q0.1 L0 PASS "one command builds every package" "go build ./... clean over $PKGS packages"
  else
    record Q0.1 L0 FAIL "one command builds every package" "go build ./... exited $build_rc: $(printf '%s' "$build_out" | head -1)"
  fi
fi

TEST_RC=127; TEST_RAN=0
if wanted Q0.2; then
  test_out=$( cd "$REPO" && go test ./... -count=1 2>&1 ); TEST_RC=$?
  TEST_RAN=$( printf '%s' "$test_out" | grep -c '^ok\|^--- PASS' || true )
  if printf '%s' "$test_out" | grep -q 'no test files' && [ "$TESTFUNCS" -eq 0 ]; then
    record Q0.2 L0 UNMEASURABLE "one command runs every test, uncached" "0 test functions in the whole repository - a green run proves nothing"
  elif [ "$TEST_RC" -eq 0 ]; then
    record Q0.2 L0 PASS "one command runs every test, uncached" "$TESTFUNCS test functions across $PKGS_WITH_TESTS of $PKGS packages"
  else
    first_fail=$( printf '%s' "$test_out" | grep -m1 -- '--- FAIL' | sed 's/^[[:space:]]*//' )
    first_pkg=$( printf '%s' "$test_out" | grep -m1 '^FAIL[[:space:]]' | awk '{print $2}' )
    record Q0.2 L0 FAIL "one command runs every test, uncached" \
      "go test ./... -count=1 exited $TEST_RC; first failure: ${first_fail:-none reported} in ${first_pkg:-unknown package}"
  fi
fi

if wanted Q0.3; then
  vet_out=$( cd "$REPO" && go vet ./... 2>&1 ); vet_rc=$?
  lint_cfg="none"
  for c in .golangci.yml .golangci.yaml .golangci.toml; do [ -f "$REPO/$c" ] && lint_cfg="$c"; done
  if [ "$vet_rc" -eq 0 ] && [ "$lint_cfg" != "none" ]; then
    record Q0.3 L0 PASS "the static gate is clean" "go vet clean; linter config $lint_cfg"
  elif [ "$vet_rc" -eq 0 ]; then
    record Q0.3 L0 FAIL "the static gate is clean" "go vet clean, but there is no checked-in linter config" "ACCEPTANCE-CRITERIA.md Q0.3"
  else
    record Q0.3 L0 FAIL "the static gate is clean" "go vet exited $vet_rc: $(printf '%s' "$vet_out" | head -1)"
  fi
fi

if wanted Q0.4; then
  if [ "$QUICK" = "yes" ]; then
    record Q0.4 L0 UNMEASURABLE "tests are hermetic" "not run (--quick); rerun without --quick"
  elif [ "$TEST_RC" -ne 0 ]; then
    record Q0.4 L0 UNMEASURABLE "tests are hermetic" "the ordinary suite is not green, so an isolated run proves nothing"
  else
    herm_home="$WORK/home"; mkdir -p "$herm_home"
    herm_out=$( cd "$REPO" && env HOME="$herm_home" go test ./... -count=1 2>&1 ); herm_rc=$?
    if [ "$herm_rc" -eq 0 ]; then
      record Q0.4 L0 PASS "tests are hermetic" "suite green with HOME pointed at an empty directory"
    else
      record Q0.4 L0 FAIL "tests are hermetic" "suite fails with an isolated HOME: $(printf '%s' "$herm_out" | grep -m1 FAIL)"  "04-seams-and-test-doubles.md §5"
    fi
  fi
fi

if wanted Q0.5; then
  if [ "$TEST_RC" -ne 0 ] || [ "$TESTFUNCS" -eq 0 ]; then
    record Q0.5 L0 UNMEASURABLE "the suite leaves the tree byte-identical" "the suite did not run green, so an unchanged tree proves nothing"
  else
    tree_now=$( tree_hash )
    if [ "$tree_now" = "$TREE_BEFORE" ]; then
      record Q0.5 L0 PASS "the suite leaves the tree byte-identical" "checksum unchanged (${TREE_BEFORE:0:12})"
    else
      record Q0.5 L0 FAIL "the suite leaves the tree byte-identical" "${TREE_BEFORE:0:12} -> ${tree_now:0:12}; run git status" "04-seams-and-test-doubles.md §6"
    fi
  fi
fi

if wanted Q0.6; then
  guard_defs=$( grep -rn --include='*_test.go' 'func wouldLeakIntoRepository' "$REPO" 2>/dev/null | grep -c . || true )
  guard_uses=$( grep -rn --include='*_test.go' 'wouldLeakIntoRepository(' "$REPO" 2>/dev/null | grep -c . || true )
  guard_uses=$(( guard_uses - guard_defs ))
  if [ "$guard_defs" -gt 0 ] && [ "$guard_uses" -gt 0 ] && [ "$SKIPPED" -eq 0 ]; then
    record Q0.6 L0 PASS "no test writes into the repository; no permanently skipped tests" "leak guard used at $guard_uses call sites; 0 skipped tests"
  else
    record Q0.6 L0 FAIL "no test writes into the repository; no permanently skipped tests" \
      "leak guard defined $guard_defs times, called $guard_uses times; $SKIPPED skipped tests of $TESTFUNCS" \
      "04-seams-and-test-doubles.md §6"
  fi
fi

if wanted Q0.7; then
  record Q0.7 L0 PASS "the denominators are recorded" \
    "$PKGS packages, $GOFILES go files, $TESTFUNCS test functions, $SCRIPTS scripts, $HARNESSES harnesses, $VERIFIERS acceptance scripts"
fi

if wanted Q0.8; then
  missing=$( p_scripts_without_metatest "$REPO" )
  if [ "$SCRIPTS" -eq 0 ]; then
    record Q0.8 L0 UNMEASURABLE "every script has a meta-test and parses" "no scripts/ directory - nothing to measure"
  elif [ "$missing" -eq 0 ]; then
    record Q0.8 L0 PASS "every script has a meta-test and parses" "$SCRIPTS scripts, all with a test- counterpart"
  else
    record Q0.8 L0 FAIL "every script has a meta-test and parses" "$missing of $SCRIPTS scripts have no scripts/test-<name>" "01-mutation-harness.md §11"
  fi
fi

# ---------------------------------------------------------------------------------------------
# L1
# ---------------------------------------------------------------------------------------------
if wanted Q1.1; then
  no_tests=$(( PKGS - PKGS_WITH_TESTS ))
  if [ "$PKGS" -eq 0 ]; then
    record Q1.1 L1 UNMEASURABLE "no package with logic has zero tests" "go list found 0 packages"
  elif [ "$no_tests" -eq 0 ]; then
    record Q1.1 L1 PASS "no package with logic has zero tests" "0 of $PKGS packages without tests"
  else
    record Q1.1 L1 FAIL "no package with logic has zero tests" "$no_tests of $PKGS packages have no test files (RATCHET)" "04-seams-and-test-doubles.md §10"
  fi
fi

if wanted Q1.2; then
  fatals=$( p_fatal_exits "$REPO" )
  if [ "$fatals" -eq 0 ]; then
    record Q1.2 L1 PASS "the entry layer is testable" "0 process-exiting calls outside main()"
  else
    record Q1.2 L1 FAIL "the entry layer is testable" "$fatals process-exiting calls outside main() - the layer cannot be driven from a test at all (RATCHET)" "04-seams-and-test-doubles.md §10"
  fi
fi

if wanted Q1.3; then
  if [ -z "$adapters_file" ]; then
    direct=$( p_direct_external "$REPO" "" )
    record Q1.3 L1 UNMEASURABLE "external effects go through an injection point" \
      "$direct files make direct http/exec/time calls, but .quality/inventory declares no [adapters], so this cannot distinguish an adapter from a violation" \
      "04-seams-and-test-doubles.md §2"
  else
    direct=$( p_direct_external "$REPO" "$adapters_file" )
    if [ "$direct" -eq 0 ]; then
      record Q1.3 L1 PASS "external effects go through an injection point" "0 direct calls outside $(grep -c . "$adapters_file") declared adapters"
    else
      record Q1.3 L1 FAIL "external effects go through an injection point" "$direct files make direct calls outside the declared adapters (RATCHET)" "04-seams-and-test-doubles.md §2"
    fi
  fi
fi

if wanted Q1.4; then
  seams=$( inv_section seams | grep -c . || true )
  if [ "$seams" -eq 0 ]; then
    record Q1.4 L1 UNMEASURABLE "every seam has an argument-swap test" \
      "no [seams] declared in .quality/inventory - an empty inventory cannot fail this check" "04-seams-and-test-doubles.md §1"
  else
    covered=0
    while IFS= read -r line; do
      [ -z "$line" ] && continue
      harness=$( printf '%s' "$line" | awk -F: '{print $2}' | tr -d ' ' )
      label=$( printf '%s' "$line" | sed 's/.*"\(.*\)".*/\1/' )
      [ -f "$REPO/$harness" ] && grep -qF -- "$label" "$REPO/$harness" && covered=$((covered + 1))
    done <<< "$( inv_section seams )"
    if [ "$covered" -eq "$seams" ]; then
      record Q1.4 L1 PASS "every seam has an argument-swap test" "$covered of $seams declared seams covered by a named mutation"
    else
      record Q1.4 L1 FAIL "every seam has an argument-swap test" "$covered of $seams declared seams covered" "04-seams-and-test-doubles.md §1"
    fi
  fi
fi

if wanted Q1.5; then
  nows=$( grep -rn --include='*.go' 'time\.Now()' "$REPO" 2>/dev/null | grep -v '_test\.go' | grep -c . || true )
  clocks=$( grep -rn --include='*.go' 'Clock ' "$REPO" 2>/dev/null | grep -c . || true )
  if [ "$clocks" -gt 0 ] && [ "$nows" -le "$clocks" ]; then
    record Q1.5 L1 PASS "time is injected" "$nows direct time.Now() calls against $clocks clock declarations"
  else
    record Q1.5 L1 FAIL "time is injected" "$nows direct time.Now() calls, $clocks clock declarations" "04-seams-and-test-doubles.md §7"
  fi
fi

for id in Q1.6 Q1.7 Q1.9; do
  wanted "$id" || continue
  case "$id" in
    Q1.6) record Q1.6 L1 UNMEASURABLE "doubles record arguments; defaults for every side effect" "MANUAL: read the test harness. Does every dependency get a default double, and is the struct passed on whole rather than rebuilt field by field?" "04-seams-and-test-doubles.md §4" ;;
    Q1.7) record Q1.7 L1 UNMEASURABLE "partial-failure commands are asserted on content" "MANUAL: does any test assert only the exit code of a command that reports per-item failures?" "04-seams-and-test-doubles.md §9" ;;
    Q1.9) record Q1.9 L1 UNMEASURABLE "tests refuse to pass on an empty population" "MANUAL: read the tests that iterate a collection. Is there an assertion - not a comment - that fails when the collection is empty?" "02-reading-a-surviving-mutant.md §2" ;;
  esac
done

if wanted Q1.8; then
  if [ "$QUICK" = "yes" ]; then
    record Q1.8 L1 UNMEASURABLE "no state leaks between runs" "not run (--quick)"
  elif [ "$TEST_RC" -ne 0 ] || [ "$TESTFUNCS" -eq 0 ]; then
    record Q1.8 L1 UNMEASURABLE "no state leaks between runs" "the ordinary suite is not green, so a second run proves nothing"
  else
    ( cd "$REPO" && go test ./... -count=2 > /dev/null 2>&1 ); c2_rc=$?
    if [ "$c2_rc" -eq 0 ]; then
      record Q1.8 L1 PASS "no state leaks between runs" "go test ./... -count=2 green"
    else
      record Q1.8 L1 FAIL "no state leaks between runs" "green with -count=1, exit $c2_rc with -count=2" "04-seams-and-test-doubles.md §5"
    fi
  fi
fi

# ---------------------------------------------------------------------------------------------
# L2
# ---------------------------------------------------------------------------------------------
if wanted Q2.1; then
  subjects=$( inv_section subjects | grep -c . || true )
  if [ "$subjects" -eq 0 ] && [ "$HARNESSES" -eq 0 ]; then
    record Q2.1 L2 FAIL "every declared subject has a mutation harness" "0 harnesses in a repository of $PKGS packages" "01-mutation-harness.md"
  elif [ "$subjects" -eq 0 ]; then
    record Q2.1 L2 UNMEASURABLE "every declared subject has a mutation harness" "$HARNESSES harnesses exist, but .quality/inventory declares no [subjects], so coverage cannot be judged" "01-mutation-harness.md"
  else
    record Q2.1 L2 PASS "every declared subject has a mutation harness" "$HARNESSES harnesses for $subjects declared subjects of $PKGS packages"
  fi
fi

if wanted Q2.2; then
  if [ "$HARNESSES" -eq 0 ]; then
    record Q2.2 L2 UNMEASURABLE "each harness declares enough mutations and can fail" "no harnesses to measure" "01-mutation-harness.md §2"
  else
    bad=""; total=0
    for h in "$REPO"/scripts/mutate-*; do
      n=$( grep -cE '^mutate2? "' "$h" || true ); total=$((total + n))
      [ "$n" -ge 8 ] || bad="$bad $(basename "$h"):${n}mut"
      grep -q '^\[ "\$survived" -eq 0 \]' "$h" || bad="$bad $(basename "$h"):nogate"
    done
    if [ -z "$bad" ]; then
      record Q2.2 L2 PASS "each harness declares enough mutations and can fail" "$total mutations across $HARNESSES harnesses, all with a failing gate"
    else
      record Q2.2 L2 FAIL "each harness declares enough mutations and can fail" "problems:$bad" "01-mutation-harness.md §2"
    fi
  fi
fi

if wanted Q2.3; then
  if [ "$HARNESSES" -eq 0 ]; then
    record Q2.3 L2 UNMEASURABLE "each harness has a T1-T10 meta-test" "no harnesses to measure" "01-mutation-harness.md §8"
  else
    bad=""
    for h in "$REPO"/scripts/mutate-*; do
      t="$REPO/scripts/test-$(basename "$h")"
      if [ ! -f "$t" ]; then bad="$bad $(basename "$h"):nometatest"; continue; fi
      ctrl=$( code_only "$t" | grep -cE '(^|[^A-Za-z])T([1-9]|10)[^0-9]' || true )
      [ "$ctrl" -ge 10 ] || bad="$bad $(basename "$t"):${ctrl}controls"
    done
    if [ -z "$bad" ]; then
      record Q2.3 L2 PASS "each harness has a T1-T10 meta-test" "$HARNESSES meta-tests, each with at least 10 controls in code"
    else
      record Q2.3 L2 FAIL "each harness has a T1-T10 meta-test" "problems:$bad" "01-mutation-harness.md §8"
    fi
  fi
fi

if wanted Q2.4; then
  record Q2.4 L2 UNMEASURABLE "declared equals killed, measured by running it" \
    "not run by this audit: a harness run mutates source and can take many minutes. Run each scripts/test-mutate-* yourself and record the numbers." \
    "01-mutation-harness.md §8"
fi

if wanted Q2.5; then
  features=$( inv_section features | grep -c . || true )
  if [ "$features" -eq 0 ] && [ "$VERIFIERS" -eq 0 ]; then
    record Q2.5 L2 FAIL "each user-facing feature has an acceptance script" "0 acceptance scripts" "03-acceptance-scripts.md"
  elif [ "$features" -eq 0 ]; then
    record Q2.5 L2 UNMEASURABLE "each user-facing feature has an acceptance script" "$VERIFIERS acceptance scripts exist, but no [features] declared" "03-acceptance-scripts.md"
  else
    missing=0
    for v in $( inv_section features | awk -F= '{print $2}' | tr -d ' ' ); do
      [ -f "$REPO/$v" ] || missing=$((missing + 1))
      [ -f "$REPO/$(dirname "$v")/test-$(basename "$v")" ] || missing=$((missing + 1))
    done
    if [ "$missing" -eq 0 ]; then
      record Q2.5 L2 PASS "each user-facing feature has an acceptance script" "$features features, each with a verify- and a test-verify- script"
    else
      record Q2.5 L2 FAIL "each user-facing feature has an acceptance script" "$missing missing scripts across $features features" "03-acceptance-scripts.md"
    fi
  fi
fi

if [ "$VERIFIERS" -eq 0 ]; then
  for id in Q2.6 Q2.7 Q2.8 Q2.9 Q2.10; do
    wanted "$id" || continue
    record "$id" L2 UNMEASURABLE "acceptance-script properties" "no acceptance scripts to inspect" "03-acceptance-scripts.md"
  done
else
  v_over=0; v_probe=0; v_guard=0; v_tee=0; v_exec=0; v_tmp=0; v_gate=0; v_badinput=0; v_heredoc=0
  for v in "$REPO"/scripts/verify-*; do
    C=$( code_only "$v" )
    printf '%s' "$C" | grep -q ':-' && v_over=$((v_over + 1))
    printf '%s' "$C" | grep -q -- '--help' && v_probe=$((v_probe + 1))
    printf '%s' "$C" | grep -qi 'proves nothing' && v_guard=$((v_guard + 1))
    [ "$( printf '%s' "$C" | grep -c 'tee -a' )" -ge 3 ] && v_tee=$((v_tee + 1))
    printf '%s' "$C" | grep -q 'exec > >(tee' && v_exec=$((v_exec + 1))
    printf '%s' "$C" | grep -q '%/}' && v_tmp=$((v_tmp + 1))
    printf '%s' "$C" | grep -qE '^\[ "\$(fail|failures)" -eq 0 \]' && v_gate=$((v_gate + 1))
    printf '%s' "$C" | grep -q 'python3 - <<' && v_heredoc=$((v_heredoc + 1))
  done
  wanted Q2.6 && { [ "$v_over" -eq "$VERIFIERS" ] && [ "$v_probe" -eq "$VERIFIERS" ] && [ "$v_tee" -eq "$VERIFIERS" ] && [ "$v_exec" -eq 0 ] && [ "$v_tmp" -eq "$VERIFIERS" ] && [ "$v_gate" -eq "$VERIFIERS" ] \
    && record Q2.6 L2 PASS "acceptance scripts measure the right artefact and log the verdict" "$VERIFIERS scripts: all have an override, a help probe, >=3 tee calls, a stripped TMPDIR and a failing gate" \
    || record Q2.6 L2 FAIL "acceptance scripts measure the right artefact and log the verdict" "of $VERIFIERS: override $v_over, help probe $v_probe, tee>=3 $v_tee, exec-tee $v_exec, TMPDIR stripped $v_tmp, gate $v_gate" "03-acceptance-scripts.md §2"; }
  wanted Q2.7 && { [ "$v_guard" -eq "$VERIFIERS" ] \
    && record Q2.7 L2 PASS "each acceptance script has a falsifiability guard" "$v_guard of $VERIFIERS refuse a population that cannot falsify" \
    || record Q2.7 L2 FAIL "each acceptance script has a falsifiability guard" "$v_guard of $VERIFIERS have one" "03-acceptance-scripts.md §4"; }
  wanted Q2.8 && record Q2.8 L2 UNMEASURABLE "a control proves the input drives the behaviour" \
    "MANUAL: read the control. It must compare a MAGNITUDE (an excess), not a COUNT - a count can be identical whether or not the limit works, so grepping for it would pass for the wrong reason." "03-acceptance-scripts.md §5"
  wanted Q2.9 && { [ "$v_badinput" -eq 0 ] \
    && record Q2.9 L2 UNMEASURABLE "bad input and read-only behaviour are controlled" "MANUAL: does a control require BOTH a non-zero exit AND zero produced artefacts?" "03-acceptance-scripts.md §6" \
    || record Q2.9 L2 PASS "bad input and read-only behaviour are controlled" "present"; }
  wanted Q2.10 && { [ "$v_heredoc" -eq 0 ] \
    && record Q2.10 L2 PASS "static checks strip comments; verdict logic is not in a stdin heredoc" "0 of $VERIFIERS pipe a program into python3 on stdin" \
    || record Q2.10 L2 FAIL "static checks strip comments; verdict logic is not in a stdin heredoc" "$v_heredoc of $VERIFIERS read the PROGRAM from stdin, so a program that also reads data there gets empty input" "03-acceptance-scripts.md §8"; }
fi

# ---------------------------------------------------------------------------------------------
# L3
# ---------------------------------------------------------------------------------------------
RULES=""
for f in AGENTS.md CLAUDE.md CONTRIBUTING.md docs/AI_WORK_RULES.md; do
  [ -f "$REPO/$f" ] && RULES="$RULES $f"
done

if wanted Q3.1; then
  if [ -z "$RULES" ]; then
    record Q3.1 L3 FAIL "the rules file is numbered and names its checks" "no AGENTS.md, CLAUDE.md, CONTRIBUTING.md or rules document found" "06-documentation-as-apparatus.md §5"
  else
    numbered=0
    for f in $RULES; do numbered=$(( numbered + $( grep -cE '^[0-9]+\. ' "$REPO/$f" || true ) )); done
    if [ "$numbered" -ge 10 ]; then
      record Q3.1 L3 PASS "the rules file is numbered and names its checks" "$numbered numbered rules in$RULES"
    else
      record Q3.1 L3 FAIL "the rules file is numbered and names its checks" "only $numbered numbered rules in$RULES" "06-documentation-as-apparatus.md §5"
    fi
  fi
fi

if wanted Q3.2; then
  doctest=$( grep -rln --include='*_test.go' 'README' "$REPO" 2>/dev/null | grep -c . || true )
  if [ "$doctest" -gt 0 ]; then
    record Q3.2 L3 PASS "the overview is tested against the real command tree" "$doctest test file(s) reference the README"
  else
    record Q3.2 L3 FAIL "the overview is tested against the real command tree" "no test references the README" "06-documentation-as-apparatus.md §4"
  fi
fi

wanted Q3.3 && record Q3.3 L3 UNMEASURABLE "entry documents are under test, order-independently" \
  "MANUAL: is there a test on the entry/boot document, and are its assertions independent of incidental order? An assertion that fails early makes every assertion below it unreachable." "06-documentation-as-apparatus.md §4"

if wanted Q3.4; then
  rot=$( p_claim_rot "$REPO" )
  record Q3.4 L3 UNMEASURABLE "claim rot is bounded" \
    "$rot state-claim phrases in markdown (RATCHET, indicator only - the phrase list is a heuristic and proves nothing on its own)" \
    "06-documentation-as-apparatus.md §3"
fi

if wanted Q3.5; then
  notes=$( ls "$REPO"/docs/design/*.md 2>/dev/null | grep -c . || true )
  if [ "$notes" -eq 0 ]; then
    record Q3.5 L3 FAIL "features have a design note with a premises table" "no docs/design/ notes" "06-documentation-as-apparatus.md §1"
  else
    withprem=$( grep -rlniE 'premise|measured, not assumed|premisser' "$REPO"/docs/design/*.md 2>/dev/null | grep -c . || true )
    if [ "$withprem" -eq "$notes" ]; then
      record Q3.5 L3 PASS "features have a design note with a premises table" "$notes notes, all with a premises section"
    else
      record Q3.5 L3 FAIL "features have a design note with a premises table" "$withprem of $notes notes have a premises section" "06-documentation-as-apparatus.md §1"
    fi
  fi
fi

if wanted Q3.6; then
  ho=$( ls "$REPO"/docs/plan/*handover* "$REPO"/docs/*handover* 2>/dev/null | grep -c . || true )
  if [ "$ho" -gt 0 ]; then
    record Q3.6 L3 PASS "a rolling handover exists" "$ho handover document(s) - whether it is rewritten rather than appended to needs a human read"
  else
    record Q3.6 L3 FAIL "a rolling handover exists" "none found" "06-documentation-as-apparatus.md §2"
  fi
fi

wanted Q3.7 && record Q3.7 L3 UNMEASURABLE "no drift between repository and installed instruction copies" \
  "MANUAL: cmp -s each instruction file in the repository against its installed copy, BOTH directions, and check what the install step actually iterates over." "07-agent-harness.md §3"

if wanted Q3.8; then
  if [ -f "$REPO/Makefile" ] && grep -qE '^(quality|verify|check):' "$REPO/Makefile"; then
    record Q3.8 L3 PASS "one entry point runs the apparatus" "Makefile has a quality/verify/check target"
  else
    record Q3.8 L3 FAIL "one entry point runs the apparatus" "no single make target runs tests, meta-tests and acceptance scripts" "01-mutation-harness.md §11"
  fi
fi

wanted Q3.9 && record Q3.9 L3 UNMEASURABLE "the score is reproducible and a human has read the output" \
  "MANUAL, and deliberately so: the top level cannot be reached by a script alone. Record the four manual findings from ACCEPTANCE-CRITERIA.md 'How to score', with numbers." "05-evidence-discipline.md §4"

# ---------------------------------------------------------------------------------------------
# Report
# ---------------------------------------------------------------------------------------------
level_summary() {
  local lvl="$1" p=0 f=0 u=0
  while IFS='|' read -r id level verdict title measured lift; do
    [ "$level" = "$lvl" ] || continue
    case "$verdict" in PASS) p=$((p+1));; FAIL) f=$((f+1));; UNMEASURABLE) u=$((u+1));; esac
  done <<< "$ROWS"
  printf '%s %s/%s   PASS %s  FAIL %s  UNMEASURABLE %s\n' "$lvl" "$p" "$((p+f+u))" "$p" "$f" "$u"
}

level_complete() {
  local lvl="$1"
  while IFS='|' read -r id level verdict title measured lift; do
    [ "$level" = "$lvl" ] || continue
    [ "$verdict" = "PASS" ] || return 1
  done <<< "$ROWS"
  return 0
}

{
  printf 'Quality audit — %s\n' "$REPO"
  printf 'commit %s · tool version %s · %s\n\n' \
    "$( cd "$REPO" && git rev-parse --short HEAD 2>/dev/null || echo 'not a git repo' )" \
    "$VERSION" "$( [ "$QUICK" = yes ] && echo 'quick mode' || echo 'full mode' )"

  printf 'Denominators\n'
  printf '  packages %s · with tests %s · go files %s · test functions %s (%s skipped)\n' \
    "$PKGS" "$PKGS_WITH_TESTS" "$GOFILES" "$TESTFUNCS" "$SKIPPED"
  printf '  table-driven tests %s · scripts %s · mutation harnesses %s · acceptance scripts %s\n' \
    "$( p_table_tests "$REPO" )" "$SCRIPTS" "$HARNESSES" "$VERIFIERS"
  printf '  enabled CI workflows %s · inventory %s\n\n' \
    "$( p_enabled_ci "$REPO" )" "$( [ -f "$INVENTORY" ] && echo present || echo 'absent (several checks are therefore UNMEASURABLE)' )"

  printf 'Findings\n'
  while IFS='|' read -r id level verdict title measured lift; do
    [ -z "$id" ] && continue
    printf '  %-12s %-5s %s\n' "$verdict" "$id" "$title"
    printf '               measured: %s\n' "$measured"
    [ -n "$lift" ] && printf '               lift:     %s\n' "$lift"
    case " $PROBED " in *" $id "*) ;; *) printf '               (no negative probe for this check)\n' ;; esac
  done <<< "$ROWS"

  printf '\nLevels\n'
  for lvl in L0 L1 L2 L3; do printf '  %s\n' "$( level_summary "$lvl" )"; done

  attained="none"
  for lvl in L0 L1 L2 L3; do
    if level_complete "$lvl"; then attained="$lvl"; else break; fi
  done
  printf '\nATTAINED LEVEL: %s\n' "$attained"

  printf '\nStill required by hand (no script can do these — ACCEPTANCE-CRITERIA.md “How to score”)\n'
  printf '  1. seam test on the three highest-value external calls\n'
  printf '  2. read one real human output, with real numbers, as a user\n'
  printf '  3. classify three survivors (or three untested branches) by the three causes\n'
  printf '  4. verify one premise the project documents, against the code\n'
} 2>&1 | tee "$REPORT"

# The verdict must reach disk. An ordinary pipe: `exec > >(tee ...)` can fail and leave the
# log empty while everything on screen looks normal.

TREE_AFTER=$( tree_hash )
if [ "$TREE_BEFORE" != "$TREE_AFTER" ]; then
  printf '\nAUDIT BROKEN: this audit changed the repository (%s -> %s). The score is void.\n' \
    "${TREE_BEFORE:0:12}" "${TREE_AFTER:0:12}" | tee -a "$REPORT" >&2
  exit 2
fi

printf '\nreport: %s\n' "$REPORT"
[ "$FAIL_N" -eq 0 ] && [ "$UNMEAS_N" -eq 0 ]
