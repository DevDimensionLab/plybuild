#!/usr/bin/env bash
# TEMPLATE-VERSION: 1
# Deliberately dishonest fixtures prove that the wrapper can reject false PASSes.
set -uo pipefail
export LC_ALL=C LANG=C

HERE="$(cd "$(dirname "$0")" && pwd)"
WRAPPER="$HERE/quality-audit.sh"
PARSER="$HERE/scorecard.py"
VENDOR="$HERE/vendor/quality-audit.sh"
REPORT_TEMPLATE="$HERE/vendor/audit-report.template.md"
LOCK="$HERE/../methodology.lock"
TMP_BASE="${TMPDIR:-/tmp}"
WORK="$(mktemp -d "${TMP_BASE%/}/test-ply-quality-audit.XXXXXX")"
if [ "${QUALITY_TEST_KEEP_WORK:-no}" = "yes" ]; then
  printf 'quality meta-test work: %s\n' "$WORK"
else
  trap 'rm -rf "$WORK"' EXIT
fi

fail() { printf 'FAIL: %s\n' "$1" >&2; exit 1; }
ok() { printf '  %s\n' "$1"; }

assert_no_generated_outputs() {
  directory="$1"
  label="$2"
  for name in scorecard.md scorecard.json raw-upstream-scorecard.md \
    .scorecard.json.pending .raw-upstream-scorecard.md.pending; do
    [ ! -e "$directory/$name" ] && [ ! -L "$directory/$name" ] \
      || fail "$label left generated artifact $name"
  done
}

printf 'T1 scripts parse and the vendored source matches its pin\n'
bash -n "$WRAPPER" || fail "wrapper does not parse"
bash -n "$VENDOR" || fail "vendored audit does not parse"
PYTHONPYCACHEPREFIX="$WORK/pycache" python3 -m py_compile "$PARSER" \
  || fail "scorecard parser does not compile"
python3 - "$PARSER" "$WORK/identity-probe" <<'PY' \
  || fail "selected-file post-scan identity guard accepted changed bytes"
from pathlib import Path
import sys

source = Path(sys.argv[1])
namespace = {"__name__": "quality_scorecard_meta", "__file__": str(source)}
exec(compile(source.read_text(encoding="utf-8"), str(source), "exec"), namespace)
repository = Path(sys.argv[2])
repository.mkdir()
(repository / "value.go").write_text("package fixture\n", encoding="utf-8")
context = {
    "manifest": {"production": ["value.go"], "tests": []},
    "metadata": {
        "production_files": {"count": 1, "sha256": namespace["selected_file_digest"](repository, ["value.go"])},
        "test_files": {"count": 0, "sha256": namespace["selected_file_digest"](repository, [])},
    },
}
(repository / "value.go").write_text("package fixture\n// changed\n", encoding="utf-8")
try:
    namespace["verify_selected_file_identity"](repository, context)
except namespace["AuditBroken"]:
    pass
else:
    raise SystemExit(1)
PY
python3 - "$PARSER" "$WORK/selection-probe" <<'PY' \
  || fail "post-scan build-selection guard accepted a new selected file"
from pathlib import Path
import sys

source = Path(sys.argv[1])
namespace = {"__name__": "quality_scorecard_selection_meta", "__file__": str(source)}
exec(compile(source.read_text(encoding="utf-8"), str(source), "exec"), namespace)
repository = Path(sys.argv[2])
repository.mkdir()
(repository / "go.mod").write_text("module example.invalid/selection\n\ngo 1.18\n", encoding="utf-8")
(repository / "value.go").write_text("package selection\n", encoding="utf-8")
context = namespace["go_build_context"](repository)
(repository / "added.go").write_text("package selection\n", encoding="utf-8")
try:
    namespace["verify_build_context_stable"](repository, context)
except namespace["AuditBroken"]:
    pass
else:
    raise SystemExit(1)
PY
python3 - "$PARSER" "$HERE" "$WORK/instrument-probe" <<'PY' \
  || fail "instrument identity guard accepted scanner bytes changed during audit"
from pathlib import Path
import shutil
import sys
from types import SimpleNamespace

source = Path(sys.argv[1])
tool_source = Path(sys.argv[2])
probe = Path(sys.argv[3])
probe.mkdir()
vendor = probe / "vendor"
vendor.mkdir()
for name in ("scorecard.py", "quality-audit.sh", "go-callscan.go.src", "q06-contract_test.go.src"):
    shutil.copyfile(tool_source / name, probe / name)
shutil.copyfile(tool_source / "vendor" / "quality-audit.sh", vendor / "quality-audit.sh")
shutil.copyfile(tool_source / "vendor" / "audit-report.template.md", vendor / "audit-report.template.md")
namespace = {"__name__": "quality_scorecard_instrument_meta", "__file__": str(probe / "scorecard.py")}
exec(compile((probe / "scorecard.py").read_text(encoding="utf-8"), str(probe / "scorecard.py"), "exec"), namespace)
args = SimpleNamespace(upstream=str(vendor / "quality-audit.sh"), partial=False, upstream_exit=1)
identity = namespace["instrument_identity"](args)
(probe / "go-callscan.go.src").write_text("changed\n", encoding="utf-8")
try:
    namespace["tool_metadata"](
        args, {"upstream_version": "1", "mode": "full mode"}, {"metadata": {}}, identity,
    )
except namespace["AuditBroken"]:
    pass
else:
    raise SystemExit(1)
PY
python3 - "$PARSER" "$WORK/inventory-probe" <<'PY' \
  || fail "inventory identity guard accepted a changed or newly created inventory"
from pathlib import Path
import sys

source = Path(sys.argv[1])
namespace = {"__name__": "quality_scorecard_inventory_meta", "__file__": str(source)}
exec(compile(source.read_text(encoding="utf-8"), str(source), "exec"), namespace)
repository = Path(sys.argv[2])
(repository / ".quality").mkdir(parents=True)
inventory_path = repository / ".quality" / "inventory"
inventory_path.write_text("[adapters]\ninternal/adapter/process\n", encoding="utf-8")
_, identity = namespace["inventory_snapshot"](repository, None)
inventory_path.write_text("[adapters]\ninternal/adapter/filesystem\n", encoding="utf-8")
try:
    namespace["verify_inventory_identity"](repository, identity)
except namespace["AuditBroken"]:
    pass
else:
    raise SystemExit(1)

inventory_path.unlink()
_, absent_identity = namespace["inventory_snapshot"](repository, None)
inventory_path.write_text("[adapters]\ninternal/adapter/process\n", encoding="utf-8")
try:
    namespace["verify_inventory_identity"](repository, absent_identity)
except namespace["AuditBroken"]:
    pass
else:
    raise SystemExit(1)
PY
python3 - "$PARSER" "$WORK/mixed-ratchet-probe.json" <<'PY' \
  || fail "mixed ratchet summary hid selected criteria that were not comparable"
import json
from pathlib import Path
import sys

source = Path(sys.argv[1])
namespace = {"__name__": "quality_scorecard_ratchet_meta", "__file__": str(source)}
exec(compile(source.read_text(encoding="utf-8"), str(source), "exec"), namespace)
build = {field: "fixture" for field in namespace["BASELINE_BUILD_FIELDS"]}
tool = {field: "fixture" for field in namespace["BASELINE_TOOL_FIELDS"]}
tool.update({"partial": False, "go_build": build})
inventory = {"present": True, "sha256": "a" * 64}
baseline_metrics = {
    "Q0.8": {"name": "scripts_without_meta_test", "value": 1, "direction": "max", "population": 1},
    "Q1.2": {"name": "process_exiting_calls", "value": 127, "direction": "max", "population": 1},
}
criteria = []
for level, count in ((0, 8), (1, 9), (2, 10), (3, 9)):
    for number in range(1, count + 1):
        criterion_id = "Q{}.{}".format(level, number)
        item = {"id": criterion_id}
        if criterion_id in baseline_metrics:
            item["metric"] = baseline_metrics[criterion_id]
        criteria.append(item)
baseline = {
    "schema_version": namespace["SCHEMA_VERSION"],
    "repository": {"module": "example.invalid/fixture", "commit": "baseline",
                   "tree": {"measurement_clean": True}},
    "inventory": inventory,
    "tool": tool,
    "criteria": criteria,
}
baseline_path = Path(sys.argv[2])
baseline_path.write_text(json.dumps(baseline), encoding="utf-8")
parsed = {"criteria": [
    {"id": "Q0.8", "verdict": "UNMEASURABLE", "measured": "empty",
     "metric": {"name": "scripts_without_meta_test", "value": 0,
                "direction": "max", "population": 0}},
    {"id": "Q1.2", "verdict": "FAIL", "measured": "held",
     "metric": {"name": "process_exiting_calls", "value": 127,
                "direction": "max", "population": 1}},
]}
result = namespace["apply_baseline"](
    parsed, baseline_path, "example.invalid/fixture", tool, inventory,
)
if not (
    result["status"] == "compared" and result["selected"] == 2 and
    result["compared"] == 1 and result["current_not_comparable"] == 1
):
    raise SystemExit(1)
PY
[ -f "$HERE/go-callscan.go.src" ] && [ ! -e "$HERE/go-callscan.go" ] \
  || fail "Go scanner source would be counted as product Go code"
grep -q 'tail -n +2' "$HERE/../baseline/README.md" \
  || fail "baseline recipe does not normalize the checkout-specific raw-report heading"
grep -q 'tool.go_build' "$HERE/../baseline/README.md" \
  || fail "baseline recipe does not bind the Go build context"
grep -q 'export LC_ALL=C LANG=C' "$WRAPPER" || fail "wrapper does not pin a locale-safe hash environment"
grep -q 'go env GOMODCACHE' "$WRAPPER" || fail "wrapper does not keep module cache outside isolated HOME"
expected="$(sed -n 's/^QUALITY_AUDIT_SHA256=//p' "$LOCK")"
actual="$(LC_ALL=C shasum -a 256 "$VENDOR" | awk '{print $1}')"
[ -n "$expected" ] && [ "$actual" = "$expected" ] \
  || fail "vendored audit drifted from methodology.lock"
expected="$(sed -n 's/^AUDIT_REPORT_TEMPLATE_SHA256=//p' "$LOCK")"
actual="$(LC_ALL=C shasum -a 256 "$REPORT_TEMPLATE" | awk '{print $1}')"
[ -n "$expected" ] && [ "$actual" = "$expected" ] \
  || fail "vendored report template drifted from methodology.lock"
expected_inventory="$(python3 - "$HERE/../baseline/scorecard.json" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as stream:
    print(json.load(stream)["inventory"]["sha256"])
PY
)"
actual_inventory="$(LC_ALL=C shasum -a 256 "$HERE/../baseline/inventory" | awk '{print $1}')"
[ -n "$expected_inventory" ] && [ "$actual_inventory" = "$expected_inventory" ] \
  || fail "preserved baseline inventory does not match scorecard metadata"
json_assert_build_context="$(python3 - "$HERE/../baseline/scorecard.json" <<'PY'
import json
import sys

value = json.load(open(sys.argv[1], encoding="utf-8"))["tool"]["go_build"]
required = (
    "version", "goos", "goarch", "cgo_enabled", "goenv", "goflags", "gowork", "production_files", "test_files",
    "go386", "goamd64", "goarm", "goarm64", "gomips", "gomips64", "goppc64", "goriscv64",
    "gowasm", "goexperiment",
)
print("yes" if all(key in value for key in required) and value["production_files"]["count"] > 0 else "no")
PY
)"
[ "$json_assert_build_context" = "yes" ] \
  || fail "baseline lacks a non-empty labeled Go build population"
grep -q 'CGO_ENABLED=0 GOENV=off GOWORK=off' "$WRAPPER" \
  || fail "wrapper does not align upstream and structured Go build context"
bash "$VENDOR" --self-test >/dev/null || fail "vendored audit's negative probes failed"
ok "shell/Python parse, checksum pin, and upstream negative probes passed"

mkdir -p "$WORK/repo/.quality" "$WORK/repo/internal/adapter/process" \
  "$WORK/repo/internal/testutil"
printf 'module example.invalid/qualityfixture\n\ngo 1.18\n' > "$WORK/repo/go.mod"
printf 'package qualityfixture\n\nfunc Value() int { return 1 }\n' > "$WORK/repo/value.go"
cat > "$WORK/repo/.quality/inventory" <<'EOF'
[subjects]
fixture = . : scripts/mutate-fixture

[adapters]
internal/adapter/process
EOF
cat > "$WORK/repo/internal/adapter/process/run.go" <<'EOF'
package process

import "os/exec"

func Command(name string, arguments ...string) *exec.Cmd {
	return exec.Command(name, arguments...)
}
EOF
cat > "$WORK/repo/internal/testutil/path.go" <<'EOF'
package testutil

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func WouldLeakIntoWorkingTree(path string) bool {
	if path == "" { return true }
	cwd, err := os.Getwd()
	if err != nil { return true }
	absolute, err := filepath.Abs(path)
	if err != nil { return true }
	relative, err := filepath.Rel(cwd, absolute)
	if err != nil { return true }
	return relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator)))
}
func ensureOutsideWorkingTree(path string) error {
	if WouldLeakIntoWorkingTree(path) { return fs.ErrInvalid }
	return nil
}
func WriteFileOutsideWorkingTree(path string, value []byte, mode fs.FileMode) error {
	if err := ensureOutsideWorkingTree(path); err != nil { return err }
	return os.WriteFile(path, value, mode)
}
func CopyFSOutsideWorkingTree(path string, source fs.FS) error {
	_ = source
	if err := ensureOutsideWorkingTree(path); err != nil { return err }
	target := filepath.Join(path, "copied")
	if err := ensureOutsideWorkingTree(target); err != nil { return err }
	return os.MkdirAll(target, 0o700)
}
EOF
cat > "$WORK/repo/internal/testutil/path_test.go" <<'EOF'
package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafeFixtureWritersRefuseRepositoryPaths(t *testing.T) {
	if err := WriteFileOutsideWorkingTree("", nil, 0o600); err == nil { t.Fatal("write was accepted") }
	if err := CopyFSOutsideWorkingTree("", os.DirFS(t.TempDir())); err == nil { t.Fatal("copy was accepted") }
}

func TestWouldLeakIntoRepositoryRejectsWorkingTreeAndAcceptsTempDir(t *testing.T) {
	if !WouldLeakIntoWorkingTree("") { t.Fatal("guard accepted empty path") }
	repositoryPath := filepath.Join(".", ".quality-control-cleanup")
	t.Cleanup(func() { _ = os.Remove(repositoryPath) })
	repositoryRoot, err := os.Getwd()
	if err != nil { t.Fatal(err) }
	if err := os.Symlink(repositoryRoot, filepath.Join(t.TempDir(), "repository-link")); err != nil { t.Fatal(err) }
}
EOF
cat > "$WORK/repo/value_test.go" <<'EOF'
package qualityfixture

import (
	"path/filepath"
	"testing"

	"example.invalid/qualityfixture/internal/testutil"
)

func TestValue(t *testing.T) {
	if Value() != 1 {
		t.Fatal("wrong value")
	}
	if err := testutil.WriteFileOutsideWorkingTree(filepath.Join(t.TempDir(), "fixture"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := testutil.CopyFSOutsideWorkingTree(t.TempDir(), nil); err != nil {
		t.Fatal(err)
	}
}
EOF
git -C "$WORK/repo" init -q
git -C "$WORK/repo" -c user.email=t@example.invalid -c user.name=t add .
git -C "$WORK/repo" -c user.email=t@example.invalid -c user.name=t commit -q -m init
commit="$(git -C "$WORK/repo" rev-parse HEAD)"
short="$(git -C "$WORK/repo" rev-parse --short HEAD)"

printf 'T1b Go environment cannot bypass or desynchronize audit execution\n'
goflags_repo="$WORK/goflags-repo"
mkdir -p "$goflags_repo"
cat > "$goflags_repo/go.mod" <<'EOF'
module example.invalid/goflags
go 1.18
EOF
printf 'package goflags\nfunc Value() int { return 1 }\n' > "$goflags_repo/value.go"
cat > "$goflags_repo/value_test.go" <<'EOF'
package goflags
import "testing"
func TestMustRunAndFail(t *testing.T) { t.Fatal("must run") }
EOF
git -C "$goflags_repo" init -q
git -C "$goflags_repo" -c user.email=t@example.invalid -c user.name=t add .
git -C "$goflags_repo" -c user.email=t@example.invalid -c user.name=t commit -q -m init
GOFLAGS='-run=^$' bash "$WRAPPER" "$goflags_repo" --only Q0.2 \
  --out "$WORK/goflags-bypass" >/dev/null 2>&1
rc=$?
[ "$rc" -eq 2 ] || fail "test-selecting GOFLAGS returned $rc, expected audit-broken exit 2"
assert_no_generated_outputs "$WORK/goflags-bypass" "test-selecting GOFLAGS"

host_goos="$(go env GOHOSTOS)"
case "$host_goos" in
  linux) cross_goos=darwin ;;
  *) cross_goos=linux ;;
esac
GOOS="$cross_goos" bash "$WRAPPER" "$goflags_repo" --only Q0.2 \
  --out "$WORK/cross-target" >/dev/null 2>&1
rc=$?
[ "$rc" -eq 2 ] || fail "cross-target Go context returned $rc, expected audit-broken exit 2"
assert_no_generated_outputs "$WORK/cross-target" "cross-target Go context"

workspace="$WORK/go-workspace"
workspace_repo="$workspace/repository"
workspace_dependency="$workspace/dependency"
mkdir -p "$workspace_repo" "$workspace_dependency"
cat > "$workspace_repo/go.mod" <<'EOF'
module example.invalid/workspace-repository
go 1.18
require example.invalid/workspace-dependency v0.0.0
EOF
cat > "$workspace_repo/value.go" <<'EOF'
package workspace_repository
import "example.invalid/workspace-dependency"
func Value() int { return workspace_dependency.Value() }
EOF
cat > "$workspace_dependency/go.mod" <<'EOF'
module example.invalid/workspace-dependency
go 1.18
EOF
cat > "$workspace_dependency/value.go" <<'EOF'
package workspace_dependency
func Value() int { return 1 }
EOF
cat > "$workspace/go.work" <<'EOF'
go 1.18
use (
	./repository
	./dependency
)
EOF
git -C "$workspace_repo" init -q
git -C "$workspace_repo" -c user.email=t@example.invalid -c user.name=t add .
git -C "$workspace_repo" -c user.email=t@example.invalid -c user.name=t commit -q -m init
GOPROXY=off GOWORK="$workspace/go.work" bash "$WRAPPER" "$workspace_repo" --only Q0.1 \
  --out "$WORK/workspace-context" >/dev/null 2>&1
rc=$?
[ "$rc" -eq 1 ] || fail "workspace isolation probe returned $rc, expected project-finding exit 1"
python3 - "$WORK/workspace-context/scorecard.json" <<'PY' \
  || fail "workspace build was not measured under the fixed module-only context"
import json
import sys

value = json.load(open(sys.argv[1], encoding="utf-8"))
criterion = value["criteria"][0]
build = value["tool"]["go_build"]
if not (
    criterion["id"] == "Q0.1" and criterion["verdict"] == "FAIL" and
    build["cgo_enabled"] == "0" and build["goenv"] == "off" and build["gowork"] == "off"
):
    raise SystemExit(1)
PY
ok "test selection/cross-target execution are rejected and Go workspace/CGO/GOENV context is aligned"

cat > "$WORK/fake-upstream.sh" <<'EOF'
#!/usr/bin/env bash
set -u
repo="."
out=""
while [ $# -gt 0 ]; do
  case "$1" in
    --out) out="$2"; shift 2 ;;
    --only) shift 2 ;;
    -*) shift ;;
    *) repo="$1"; shift ;;
  esac
done
[ -n "$out" ] || out="$repo/target/quality-audit"
mkdir -p "$out"
cp "$FAKE_REPORT" "$out/scorecard.md"
if [ "${FAKE_SEED_OUTPUTS:-no}" = "yes" ]; then
  printf 'stale upstream json\n' > "$out/scorecard.json"
  printf 'stale upstream raw\n' > "$out/raw-upstream-scorecard.md"
  printf 'stale upstream pending\n' > "$out/.scorecard.json.pending"
fi
if [ "${FAKE_EMPTY_REPORT:-no}" = "yes" ]; then
  : > "$out/scorecard.md"
fi
if [ "${FAKE_BLOCK_RAW:-no}" = "yes" ]; then
  (
    while [ ! -e "$out/.raw-upstream-scorecard.md.pending" ]; do sleep 0.01; done
    mkdir "$out/raw-upstream-scorecard.md"
  ) &
fi
if [ "${FAKE_BLOCK_JSON:-no}" = "yes" ]; then
  (
    while [ ! -e "$out/.scorecard.json.pending" ]; do sleep 0.01; done
    mkdir "$out/scorecard.json"
  ) &
fi
if [ "${FAKE_HARDLINK_JSON:-no}" = "yes" ]; then
  (
    while [ ! -e "$out/.scorecard.json.pending" ]; do sleep 0.01; done
    ln "$out/.scorecard.json.pending" "$out/scorecard.json"
  ) &
fi
if [ "${FAKE_HARDLINK_RAW:-no}" = "yes" ]; then
  (
    while [ ! -e "$out/.raw-upstream-scorecard.md.pending" ]; do sleep 0.01; done
    ln "$out/.raw-upstream-scorecard.md.pending" "$out/raw-upstream-scorecard.md"
  ) &
fi
if [ "${FAKE_DIRECTORY_REPORT:-no}" = "yes" ]; then
  unlink "$out/scorecard.md"
  mkdir "$out/scorecard.md"
fi
cat "$out/scorecard.md"
exit "${FAKE_EXIT:-0}"
EOF
chmod +x "$WORK/fake-upstream.sh"

write_report() {
  criterion="$1"
  verdict="$2"
  measured="$3"
  REPORT_CRITERION="$criterion"
  cat > "$WORK/report.md" <<EOF
Quality audit - fixture
commit ${REPORT_SHORT:-$short} . tool version 1 . quick mode

Denominators
  packages ${REPORT_PACKAGES:-1} . with tests ${REPORT_TESTED_PACKAGES:-1} . go files ${REPORT_GO_FILES:-1} . test functions ${REPORT_TEST_FUNCTIONS:-1} (${REPORT_SKIPPED:-0} skipped)
  table-driven tests 0 . scripts ${REPORT_SCRIPTS:-0} . mutation harnesses 0 . acceptance scripts 0
  enabled CI workflows 0 . inventory absent

Findings
  $verdict         $criterion fixture criterion
               measured: $measured

ATTAINED LEVEL: none
EOF
}

run_wrapper() {
  out="$1"
  shift
  if [ "${QUALITY_TEST_VERBOSE:-no}" = "yes" ]; then
    FAKE_REPORT="$WORK/report.md" QUALITY_AUDIT_UPSTREAM="$WORK/fake-upstream.sh" \
      bash "$WRAPPER" "${RUN_REPO:-$WORK/repo}" --out "$out" --only "$REPORT_CRITERION" "$@"
  else
    FAKE_REPORT="$WORK/report.md" QUALITY_AUDIT_UPSTREAM="$WORK/fake-upstream.sh" \
      bash "$WRAPPER" "${RUN_REPO:-$WORK/repo}" --out "$out" --only "$REPORT_CRITERION" "$@" >/dev/null 2>&1
  fi
}

json_assert() {
  file="$1"
  expression="$2"
  message="$3"
  if ! python3 - "$file" "$expression" <<'PY'
import json
import sys

value = json.load(open(sys.argv[1], encoding="utf-8"))
safe = {"value": value, "all": all, "any": any, "len": len}
if not eval(sys.argv[2], {"__builtins__": {}}, safe):
    raise SystemExit(1)
PY
  then
    fail "$message"
  fi
}

printf 'T2 ratchets distinguish improved, held, and regressed measurements\n'
write_report Q1.2 FAIL "127 process-exiting calls outside main() (RATCHET)"
run_wrapper "$WORK/baseline-seed"
rc=$?
[ "$rc" -eq 1 ] || fail "baseline seed returned $rc, expected 1"
python3 - "$WORK/baseline-seed/scorecard.json" "$WORK/baseline.json" <<'PY'
import json
import sys

value = json.load(open(sys.argv[1], encoding="utf-8"))
metrics = {
    "Q0.6": {"name": "skipped_tests", "value": 2, "direction": "max", "population": 1,
              "precondition_met": False},
    "Q0.8": {"name": "scripts_without_meta_test", "value": 1, "direction": "max", "population": 1},
    "Q1.1": {"name": "packages_without_tests", "value": 0, "direction": "max", "population": 1},
    "Q1.2": {"name": "process_exiting_calls", "value": 127, "direction": "max", "population": 1},
    "Q1.3": {"name": "direct_external_call_sites", "value": 9, "direction": "max", "population": 10,
              "declared_adapters": 1},
    "Q1.4": {"name": "covered_seams", "numerator": 0, "denominator": 1, "direction": "min_ratio"},
    "Q2.1": {"name": "valid_subject_harness_bindings", "numerator": 0, "denominator": 1,
              "direction": "min_ratio", "source_precondition_met": True},
    "Q3.4": {"name": "state_claim_phrases", "value": 0, "direction": "max", "population": 1},
}
criteria = []
for level, count in ((0, 8), (1, 9), (2, 10), (3, 9)):
    for number in range(1, count + 1):
        criterion_id = "Q{}.{}".format(level, number)
        item = {"id": criterion_id, "level": "L{}".format(level), "verdict": "FAIL",
                "title": "meta-test baseline", "measured": "fixture"}
        if criterion_id in metrics:
            item["metric"] = metrics[criterion_id]
        criteria.append(item)
value["criteria"] = criteria
value["tool"]["partial"] = False
json.dump(value, open(sys.argv[2], "w", encoding="utf-8"), indent=2, sort_keys=True)
PY
write_report Q1.2 PASS "128 process-exiting calls outside main() (RATCHET)"
run_wrapper "$WORK/ratchet" --baseline "$WORK/baseline.json"
rc=$?
[ "$rc" -eq 1 ] || fail "worsened ratchet returned $rc, expected 1"
json_assert "$WORK/ratchet/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['ratchet']['regressed'] == 1" \
  "worsened ratchet remained passing"
write_report Q1.2 FAIL "127 process-exiting calls outside main() (RATCHET)"
run_wrapper "$WORK/ratchet-held" --baseline "$WORK/baseline.json"
rc=$?
[ "$rc" -eq 0 ] || fail "held ratchet returned $rc, expected 0"
json_assert "$WORK/ratchet-held/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'PASS' and value['criteria'][0]['ratchet']['outcome'] == 'held'" \
  "held ratchet did not pass against the stored baseline"
write_report Q1.2 FAIL "126 process-exiting calls outside main() (RATCHET)"
run_wrapper "$WORK/ratchet-improved" --baseline "$WORK/baseline.json"
rc=$?
[ "$rc" -eq 0 ] || fail "improved ratchet returned $rc, expected 0"
json_assert "$WORK/ratchet-improved/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'PASS' and value['criteria'][0]['ratchet']['outcome'] == 'improved' and value['ratchet']['improved'] == 1 and value['ratchet']['baseline_sha256'] == '$(LC_ALL=C shasum -a 256 "$WORK/baseline.json" | awk '{print $1}')'" \
  "improved ratchet was not distinguished from a held metric"
python3 - "$WORK/baseline.json" "$WORK" <<'PY'
import copy
import json
import pathlib
import sys

baseline = json.load(open(sys.argv[1], encoding="utf-8"))
directory = pathlib.Path(sys.argv[2])
variants = {}
variants["baseline-goos"] = copy.deepcopy(baseline)
variants["baseline-goos"]["tool"]["go_build"]["goos"] = "not-the-current-goos"
variants["baseline-goflags"] = copy.deepcopy(baseline)
variants["baseline-goflags"]["tool"]["go_build"]["goflags"] = "-tags=forged"
variants["baseline-goarm64"] = copy.deepcopy(baseline)
variants["baseline-goarm64"]["tool"]["go_build"]["goarm64"] = "forged"
variants["baseline-scanner"] = copy.deepcopy(baseline)
variants["baseline-scanner"]["tool"]["call_scanner_sha256"] = "0" * 64
variants["baseline-inventory"] = copy.deepcopy(baseline)
variants["baseline-inventory"]["inventory"]["sha256"] = "0" * 64
variants["baseline-partial"] = copy.deepcopy(baseline)
variants["baseline-partial"]["tool"]["partial"] = True
variants["baseline-duplicate"] = copy.deepcopy(baseline)
variants["baseline-duplicate"]["criteria"].append(copy.deepcopy(variants["baseline-duplicate"]["criteria"][0]))
variants["baseline-dirty"] = copy.deepcopy(baseline)
variants["baseline-dirty"]["repository"]["tree"]["measurement_clean"] = False
for name, value in variants.items():
    (directory / (name + ".json")).write_text(
        json.dumps(value, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
PY
for invalid_baseline in baseline-goos baseline-goflags baseline-goarm64 baseline-scanner baseline-inventory baseline-partial baseline-duplicate baseline-dirty; do
  run_wrapper "$WORK/$invalid_baseline-out" --baseline "$WORK/$invalid_baseline.json"
  rc=$?
  [ "$rc" -eq 2 ] || fail "$invalid_baseline returned $rc, expected audit-broken exit 2"
  [ ! -e "$WORK/$invalid_baseline-out/scorecard.json" ] \
    || fail "$invalid_baseline published authoritative JSON"
done
python3 - "$WORK/baseline.json" "$WORK/empty-ratchet-baseline.json" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as stream:
    value = json.load(stream)
next(item for item in value["criteria"] if item["id"] == "Q1.2")["metric"]["population"] = 0
with open(sys.argv[2], "w", encoding="utf-8") as stream:
    json.dump(value, stream)
PY
write_report Q1.2 FAIL "127 process-exiting calls outside main() (RATCHET)"
run_wrapper "$WORK/empty-ratchet" --baseline "$WORK/empty-ratchet-baseline.json"
rc=$?
[ "$rc" -eq 2 ] || fail "empty ratchet baseline returned $rc, expected audit-broken exit 2"
[ ! -e "$WORK/empty-ratchet/scorecard.json" ] \
  || fail "empty ratchet baseline published authoritative JSON"
write_report Q0.8 PASS "0 scripts, all with a test- counterpart"
run_wrapper "$WORK/empty-current-ratchet" --baseline "$WORK/baseline.json"
rc=$?
[ "$rc" -eq 1 ] || fail "empty current ratchet returned $rc, expected finding exit 1"
json_assert "$WORK/empty-current-ratchet/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'UNMEASURABLE' and value['ratchet']['status'] == 'not_comparable' and value['ratchet']['compared'] == 0" \
  "empty current ratchet did not publish an unmeasurable structured result"
write_report Q0.6 PASS "leak guard used at 2 call sites; 0 skipped tests"
run_wrapper "$WORK/q06-pass" --baseline "$WORK/baseline.json"
rc=$?
[ "$rc" -eq 0 ] || fail "Q0.6 PASS-format ratchet returned $rc, expected 0"
json_assert "$WORK/q06-pass/scorecard.json" \
  "value['criteria'][0]['metric']['value'] == 0 and value['criteria'][0]['metric']['precondition_met'] and value['ratchet']['compared'] == 1" \
  "Q0.6 PASS text did not produce a comparable metric"
REPORT_SCRIPTS=1
write_report Q0.8 PASS "1 scripts, all with a test- counterpart"
run_wrapper "$WORK/q08-pass" --baseline "$WORK/baseline.json"
rc=$?
unset REPORT_SCRIPTS
[ "$rc" -eq 0 ] || fail "Q0.8 PASS-format ratchet returned $rc, expected 0"
json_assert "$WORK/q08-pass/scorecard.json" \
  "value['criteria'][0]['metric']['value'] == 0 and value['criteria'][0]['metric']['population'] == 1 and value['ratchet']['compared'] == 1" \
  "Q0.8 PASS text did not produce a comparable metric"
write_report Q1.3 PASS "0 direct calls outside 1 declared adapters"
run_wrapper "$WORK/q13-pass" --baseline "$WORK/baseline.json"
rc=$?
[ "$rc" -eq 0 ] || fail "Q1.3 PASS-format ratchet returned $rc, expected 0"
json_assert "$WORK/q13-pass/scorecard.json" \
  "value['criteria'][0]['metric']['value'] == 0 and value['criteria'][0]['metric']['population'] == 1 and value['criteria'][0]['metric']['declared_adapters'] == 1 and value['ratchet']['compared'] == 1" \
  "local Q1.3 measurement did not produce a comparable metric"
ok "regressed/held/improved semantics and PASS formats are comparable"

printf 'T3 a PASS over an empty population becomes UNMEASURABLE\n'
write_report Q1.4 PASS "0 of 0 declared seams covered by a named mutation"
run_wrapper "$WORK/empty"
rc=$?
[ "$rc" -eq 1 ] || fail "empty population returned $rc, expected 1"
json_assert "$WORK/empty/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'UNMEASURABLE'" \
  "empty declared seam population remained PASS"
ok "0 declared seams cannot produce PASS"

printf 'T4 Q0.1 cannot pass with zero packages or Go files\n'
REPORT_PACKAGES=0 REPORT_TESTED_PACKAGES=0 REPORT_GO_FILES=0 REPORT_TEST_FUNCTIONS=0 \
  write_report Q0.1 PASS "go build ./... clean over 0 packages"
run_wrapper "$WORK/empty-build-population"
rc=$?
[ "$rc" -eq 1 ] || fail "empty build population returned $rc, expected 1"
json_assert "$WORK/empty-build-population/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'UNMEASURABLE'" \
  "Q0.1 passed with no packages and no Go files"
run_wrapper "$WORK/non-ratchet-with-baseline" --baseline "$WORK/baseline.json"
rc=$?
[ "$rc" -eq 1 ] || fail "non-ratchet --only with baseline returned $rc, expected finding exit 1"
json_assert "$WORK/non-ratchet-with-baseline/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'UNMEASURABLE' and value['ratchet']['status'] == 'not_applicable' and value['ratchet']['compared'] == 0" \
  "valid baseline required a globally comparable ratchet for non-ratchet --only"
ok "empty and non-ratchet selected populations remain structured with a baseline"

printf 'T5 an empty tree hash makes the audit broken, not passing\n'
write_report Q0.5 PASS "checksum unchanged ()"
mkdir -p "$WORK/empty-hash"
printf 'STALE AUTHORITATIVE JSON\n' > "$WORK/empty-hash/scorecard.json"
printf 'STALE RAW REPORT\n' > "$WORK/empty-hash/raw-upstream-scorecard.md"
printf 'STALE TRANSIENT REPORT\n' > "$WORK/empty-hash/scorecard.md"
export FAKE_SEED_OUTPUTS=yes
run_wrapper "$WORK/empty-hash"
rc=$?
unset FAKE_SEED_OUTPUTS
[ "$rc" -eq 2 ] || fail "empty checksum returned $rc, expected audit-broken exit 2"
[ ! -f "$WORK/empty-hash/scorecard.json" ] || fail "empty checksum still produced structured evidence"
[ -s "$WORK/empty-hash/raw-upstream-scorecard.md" ] \
  || fail "current raw report was not preserved after parser failure"
grep -q 'checksum unchanged ()' "$WORK/empty-hash/raw-upstream-scorecard.md" \
  || fail "stale raw report survived the failed audit"
[ ! -e "$WORK/empty-hash/scorecard.md" ] || fail "transient report was not moved to the raw path"
[ ! -e "$WORK/empty-hash/.scorecard.json.pending" ] \
  || fail "upstream-created pending JSON survived parser failure"
[ ! -e "$WORK/empty-hash/.raw-upstream-scorecard.md.pending" ] \
  || fail "private raw-report staging survived parser failure"
ok "exit 2 invalidates stale JSON/raw and preserves only the current raw report"

printf 'T6 a baseline that aliases authoritative output is preserved and rejected\n'
mkdir -p "$WORK/baseline-alias"
cp "$WORK/baseline.json" "$WORK/baseline-alias/scorecard.json"
cp "$WORK/baseline-alias/scorecard.json" "$WORK/baseline-alias.before"
write_report Q1.2 FAIL "127 process-exiting calls outside main() (RATCHET)"
run_wrapper "$WORK/baseline-alias" --baseline "$WORK/baseline-alias/scorecard.json"
rc=$?
[ "$rc" -eq 2 ] || fail "baseline/output alias returned $rc, expected 2"
cmp -s "$WORK/baseline-alias.before" "$WORK/baseline-alias/scorecard.json" \
  || fail "baseline/output alias was deleted or changed"
mkdir -p "$WORK/baseline-hardlink-alias"
cp "$WORK/baseline.json" "$WORK/baseline-hardlink-input.json"
ln "$WORK/baseline-hardlink-input.json" "$WORK/baseline-hardlink-alias/scorecard.json"
run_wrapper "$WORK/baseline-hardlink-alias" --baseline "$WORK/baseline-hardlink-input.json"
rc=$?
[ "$rc" -eq 2 ] || fail "hard-linked baseline/output alias returned $rc, expected 2"
[ -f "$WORK/baseline-hardlink-input.json" ] && [ -f "$WORK/baseline-hardlink-alias/scorecard.json" ] \
  || fail "hard-linked baseline/output alias was invalidated before rejection"
ok "baseline path and inode aliases fail before invalidation"

printf 'T6b manual evidence that aliases output is preserved and rejected\n'
mkdir -p "$WORK/manual-alias"
printf 'MANUAL EVIDENCE MUST SURVIVE\n' > "$WORK/manual-alias/scorecard.json"
cp "$WORK/manual-alias/scorecard.json" "$WORK/manual-alias.before"
run_wrapper "$WORK/manual-alias" --manual-evidence "$WORK/manual-alias/scorecard.json"
rc=$?
[ "$rc" -eq 2 ] || fail "manual/output alias returned $rc, expected 2"
cmp -s "$WORK/manual-alias.before" "$WORK/manual-alias/scorecard.json" \
  || fail "manual evidence/output alias was deleted or changed"
mkdir -p "$WORK/manual-hardlink-alias"
printf 'MANUAL HARDLINK MUST SURVIVE\n' > "$WORK/manual-hardlink-input.json"
ln "$WORK/manual-hardlink-input.json" "$WORK/manual-hardlink-alias/scorecard.json"
run_wrapper "$WORK/manual-hardlink-alias" --manual-evidence "$WORK/manual-hardlink-input.json"
rc=$?
[ "$rc" -eq 2 ] || fail "hard-linked manual/output alias returned $rc, expected 2"
[ -f "$WORK/manual-hardlink-input.json" ] && [ -f "$WORK/manual-hardlink-alias/scorecard.json" ] \
  || fail "hard-linked manual/output alias was invalidated before rejection"
ok "manual path and inode aliases fail before invalidation"

printf 'T7 manual evidence from another commit cannot pass\n'
write_report Q3.9 PASS "score reproduced twice and 4 manual findings were read"
sed 's#github.com/devdimensionlab/plybuild#example.invalid/qualityfixture#' \
  "$HERE/../baseline/manual-evidence.json" > "$WORK/stale-evidence.json"
run_wrapper "$WORK/stale" --manual-evidence "$WORK/stale-evidence.json"
rc=$?
[ "$rc" -eq 1 ] || fail "stale evidence returned $rc, expected 1"
json_assert "$WORK/stale/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['manual_evidence']['status'] == 'stale'" \
  "stale manual evidence remained passing"
run_wrapper "$WORK/missing-evidence" --manual-evidence "$WORK/does-not-exist.json"
rc=$?
[ "$rc" -eq 2 ] || fail "explicit missing evidence returned $rc, expected audit-broken exit 2"
printf '{}\n' > "$WORK/empty-evidence.json"
run_wrapper "$WORK/empty-evidence" --manual-evidence "$WORK/empty-evidence.json"
rc=$?
[ "$rc" -eq 1 ] || fail "empty evidence returned $rc, expected project-finding exit 1"
json_assert "$WORK/empty-evidence/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['manual_evidence']['status'] == 'invalid'" \
  "empty manual evidence remained passing"
ok "stale and empty evidence fail; an explicit missing path is exit 2"

printf 'T8 manual evidence is unique, digest-bound, and valid only for a clean measured tree\n'
tree="$(git -C "$WORK/repo" rev-parse 'HEAD^{tree}')"
python3 - "$HERE/../baseline/manual-evidence.json" "$WORK/valid-evidence.json" "$commit" "$tree" \
  "$WORK/repo/.quality/inventory" "$PARSER" "$WRAPPER" "$WORK/fake-upstream.sh" \
  "$HERE/go-callscan.go.src" "$HERE/q06-contract_test.go.src" <<'PY'
import hashlib
import json
import sys

with open(sys.argv[1], encoding="utf-8") as stream:
    value = json.load(stream)
value["repository"] = {
    "module": "example.invalid/qualityfixture",
    "commit": sys.argv[3],
    "tree": {
        "commit_tree": sys.argv[4],
        "status_sha256": hashlib.sha256(b"\0").hexdigest(),
        "inventory_overlay_sha256": None,
    },
}
inventory_sha256 = hashlib.sha256(open(sys.argv[5], "rb").read()).hexdigest()
value["inventory"] = {
    "path": ".quality/inventory", "present": True,
    "sha256": inventory_sha256, "overlay": False,
}

def digest(path):
    return hashlib.sha256(open(path, "rb").read()).hexdigest()

value["instrument"] = {
    "wrapper_version": 1,
    "upstream_version": "1",
    "upstream_sha256": digest(sys.argv[8]),
    "parser_sha256": digest(sys.argv[6]),
    "wrapper_sha256": digest(sys.argv[7]),
    "call_scanner_sha256": digest(sys.argv[9]),
    "q06_contract_sha256": digest(sys.argv[10]),
    "report_template_sha256": None,
}
with open(sys.argv[2], "w", encoding="utf-8") as stream:
    json.dump(value, stream, indent=2, sort_keys=True)
    stream.write("\n")
PY
write_report Q3.9 UNMEASURABLE "MANUAL evidence required"
run_wrapper "$WORK/valid-manual" --manual-evidence "$WORK/valid-evidence.json"
rc=$?
[ "$rc" -eq 0 ] || fail "valid manual evidence returned $rc, expected 0"
evidence_digest="$(LC_ALL=C shasum -a 256 "$WORK/valid-evidence.json" | awk '{print $1}')"
json_assert "$WORK/valid-manual/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'PASS' and value['manual_evidence']['status'] == 'valid' and value['manual_evidence']['sha256'] == '$evidence_digest' and value['repository']['tree']['measurement_clean']" \
  "valid manual evidence was not digest-bound to the clean tree"
write_report Q3.9 FAIL "2 runs differed even though 4 manual findings were read"
run_wrapper "$WORK/determinism-failed-manual" --manual-evidence "$WORK/valid-evidence.json"
rc=$?
[ "$rc" -eq 1 ] || fail "failed determinism with valid manual evidence returned $rc, expected 1"
json_assert "$WORK/determinism-failed-manual/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['manual_evidence']['status'] == 'valid'" \
  "manual evidence overrode the audit's failed deterministic rerun"
run_wrapper "$WORK/determinism-failed-no-manual"
rc=$?
[ "$rc" -eq 1 ] || fail "failed determinism without manual evidence returned $rc, expected 1"
json_assert "$WORK/determinism-failed-no-manual/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['manual_evidence']['status'] == 'not_supplied'" \
  "absence of manual evidence downgraded the audit's failed deterministic rerun"
write_report Q3.9 UNMEASURABLE "MANUAL evidence required"
cp "$WORK/repo/value.go" "$WORK/value.go.clean"
printf '\n// dirty measurement\n' >> "$WORK/repo/value.go"
run_wrapper "$WORK/dirty-manual" --manual-evidence "$WORK/valid-evidence.json"
rc=$?
[ "$rc" -eq 1 ] || fail "dirty-tree manual evidence returned $rc, expected 1"
json_assert "$WORK/dirty-manual/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['manual_evidence']['status'] == 'stale' and not value['repository']['tree']['measurement_clean']" \
  "manual evidence remained valid on a dirty measured tree"
cp "$WORK/value.go.clean" "$WORK/repo/value.go"
python3 - "$WORK/valid-evidence.json" "$WORK/duplicate-evidence.json" "$WORK/class-evidence.json" <<'PY'
import copy
import json
import sys

with open(sys.argv[1], encoding="utf-8") as stream:
    original = json.load(stream)
duplicate = copy.deepcopy(original)
duplicate["findings"]["seam_test"] = [duplicate["findings"]["seam_test"][0]] * 3
with open(sys.argv[2], "w", encoding="utf-8") as stream:
    json.dump(duplicate, stream)
invalid_class = copy.deepcopy(original)
invalid_class["findings"]["survivors"][0]["classification"] = "probably fine"
with open(sys.argv[3], "w", encoding="utf-8") as stream:
    json.dump(invalid_class, stream)
PY
run_wrapper "$WORK/duplicate-manual" --manual-evidence "$WORK/duplicate-evidence.json"
rc=$?
[ "$rc" -eq 1 ] || fail "duplicate manual receipts returned $rc, expected 1"
json_assert "$WORK/duplicate-manual/scorecard.json" \
  "value['manual_evidence']['status'] == 'invalid'" \
  "duplicate manual receipts were accepted"
run_wrapper "$WORK/class-manual" --manual-evidence "$WORK/class-evidence.json"
rc=$?
[ "$rc" -eq 1 ] || fail "invalid survivor class returned $rc, expected 1"
json_assert "$WORK/class-manual/scorecard.json" \
  "value['manual_evidence']['status'] == 'invalid'" \
  "survivor classification outside the documented vocabulary was accepted"
ok "manual receipts are tree/digest-bound, unique, and vocabulary-checked"

printf 'T8c six manual L1/L2 rows require bound, truthful criterion receipts\n'
python3 "$HERE/test-manual-evidence.py" "$PARSER" \
  || fail "criterion-bound manual evidence meta-test failed"
ok "all six rows are reachable and malformed or false-PASS receipts fail closed"

printf 'T8b measurement identity includes ignored and index-hidden source changes\n'
ignored_repo="$WORK/ignored-source-repo"
git clone -q "$WORK/repo" "$ignored_repo"
printf 'hidden*.go\n' > "$ignored_repo/.gitignore"
git -C "$ignored_repo" -c user.email=t@example.invalid -c user.name=t add .gitignore
git -C "$ignored_repo" -c user.email=t@example.invalid -c user.name=t commit -q -m ignore-fixture
cat > "$ignored_repo/hidden_effect.go" <<'EOF'
package qualityfixture
import "os"
func HiddenEffect(path string) error { return os.WriteFile(path, nil, 0o600) }
EOF
export RUN_REPO="$ignored_repo" REPORT_SHORT="$(git -C "$ignored_repo" rev-parse --short HEAD)"
write_report Q3.9 PASS "score reproduced twice and 4 manual findings were read"
run_wrapper "$WORK/ignored-source"
rc=$?
[ "$rc" -eq 1 ] || fail "ignored source identity probe returned $rc, expected 1"
json_assert "$WORK/ignored-source/scorecard.json" \
  "not value['repository']['tree']['measurement_clean'] and 'hidden_effect.go' in value['repository']['tree']['dirty_paths']" \
  "ignored Go source was accepted as part of a clean measured tree"
assume_repo="$WORK/assume-unchanged-repo"
git clone -q "$WORK/repo" "$assume_repo"
git -C "$assume_repo" update-index --assume-unchanged value.go
python3 - "$assume_repo/value.go" <<'PY'
from pathlib import Path
import sys

path = Path(sys.argv[1])
path.write_text(path.read_text(encoding="utf-8").replace("return 1", "return 2"), encoding="utf-8")
PY
export RUN_REPO="$assume_repo" REPORT_SHORT="$(git -C "$assume_repo" rev-parse --short HEAD)"
write_report Q3.9 PASS "score reproduced twice and 4 manual findings were read"
run_wrapper "$WORK/assume-unchanged-source"
rc=$?
[ "$rc" -eq 1 ] || fail "assume-unchanged source identity probe returned $rc, expected 1"
json_assert "$WORK/assume-unchanged-source/scorecard.json" \
  "not value['repository']['tree']['measurement_clean'] and 'value.go' in value['repository']['tree']['dirty_paths']" \
  "assume-unchanged source edit was accepted as part of a clean measured tree"
mode_repo="$WORK/mode-source-repo"
git clone -q "$WORK/repo" "$mode_repo"
chmod 0645 "$mode_repo/value.go"
export RUN_REPO="$mode_repo" REPORT_SHORT="$(git -C "$mode_repo" rev-parse --short HEAD)"
write_report Q3.9 PASS "score reproduced twice and 4 manual findings were read"
run_wrapper "$WORK/other-executable-bit-source"
rc=$?
[ "$rc" -eq 1 ] || fail "non-owner executable-bit identity probe returned $rc, expected 1"
json_assert "$WORK/other-executable-bit-source/scorecard.json" \
  "not value['repository']['tree']['measurement_clean'] and 'value.go' in value['repository']['tree']['dirty_paths']" \
  "a Git-canonical executable mode change outside the owner bit was accepted as clean"
vendor_repo="$WORK/vendor-source-repo"
git clone -q "$WORK/repo" "$vendor_repo"
mkdir -p "$vendor_repo/vendor/example.invalid/dependency"
printf 'package dependency\nfunc Value() int { return 1 }\n' > \
  "$vendor_repo/vendor/example.invalid/dependency/dependency.go"
git -C "$vendor_repo" -c user.email=t@example.invalid -c user.name=t add vendor
git -C "$vendor_repo" -c user.email=t@example.invalid -c user.name=t commit -q -m vendor-fixture
printf 'package dependency\nfunc Value() int { return 2 }\n' > \
  "$vendor_repo/vendor/example.invalid/dependency/dependency.go"
export RUN_REPO="$vendor_repo" REPORT_SHORT="$(git -C "$vendor_repo" rev-parse --short HEAD)"
write_report Q3.9 PASS "score reproduced twice and 4 manual findings were read"
run_wrapper "$WORK/vendor-source"
rc=$?
[ "$rc" -eq 1 ] || fail "vendored source identity probe returned $rc, expected 1"
json_assert "$WORK/vendor-source/scorecard.json" \
  "not value['repository']['tree']['measurement_clean'] and 'vendor/example.invalid/dependency/dependency.go' in value['repository']['tree']['dirty_paths']" \
  "modified tracked vendored source was accepted as part of a clean measured tree"
unset RUN_REPO REPORT_SHORT
ok "filesystem-vs-HEAD identity catches ignored, mode, index-hidden, and vendored changes"

printf 'T9 local Q0.6 requires the central guard contract and safe consumers\n'
q06_bad="$WORK/q06-bad-repo"
git clone -q "$WORK/repo" "$q06_bad"
python3 - "$q06_bad/internal/testutil/path.go" <<'PY'
from pathlib import Path
import sys

path = Path(sys.argv[1])
value = path.read_text(encoding="utf-8")
path.write_text(value.replace(
    "\treturn os.WriteFile(path, value, mode)",
    "\t_ = os.WriteFile(\"LEAK-IN-REPOSITORY\", nil, 0o600)\n\treturn os.WriteFile(path, value, mode)",
), encoding="utf-8")
PY
cat >> "$q06_bad/internal/testutil/path_test.go" <<'EOF'

func TestUnsafeSupportWriterIsStillMeasured(t *testing.T) {
	if err := os.WriteFile("LEAK", nil, 0o600); err != nil { t.Fatal(err) }
	_ = os.Remove("LEAK-CLEANUP")
}
EOF
cat > "$q06_bad/internal/testutil/leak.go" <<'EOF'
package testutil
import "os"
func UnsafeHelper() error { return os.WriteFile("LEAK", nil, 0o600) }
EOF
cat > "$q06_bad/helper_test.go" <<'EOF'
package qualityfixture
import "example.invalid/qualityfixture/internal/testutil"
func unsafeProductionTestHelperCall() { _ = testutil.UnsafeHelper() }
EOF
cat > "$q06_bad/unsafe_test.go" <<'EOF'
package qualityfixture
import (
	"bytes"
	"encoding/json"
	"io"
	filesystem "os"
	"path/filepath"
	"syscall"
	"testing"
)
func unsafeFixtureWrite(path string) error { return filesystem.WriteFile(path, nil, 0o600) }
func unsafeFixtureMkdir(path string) error { return filesystem.MkdirAll(path, 0o700) }
func unsafeFileWrite(file *filesystem.File) { _, _ = file.Write(nil) }
func unsafeCreateTemp() { _, _ = filesystem.CreateTemp(".", "fixture-") }
func unsafeMkdirTemp() { _, _ = filesystem.MkdirTemp("repo-relative", "fixture-") }
func unsafeFunctionValue(path string) error {
	write := filesystem.WriteFile
	return write(path, nil, 0o600)
}
func unsafeMethodValue(file *filesystem.File) { write := file.Write; _, _ = write(nil) }
func unsafeCopyFS(path string) error { return filesystem.CopyFS(path, nil) }
func unsafeRootMethod(root *filesystem.Root, path string) error {
	write := root.WriteFile
	return write(path, nil, 0o600)
}
func unsafeWriterFactory() func(string, []byte, filesystem.FileMode) error {
	return filesystem.WriteFile
}
type unsafeOperations struct { write func(string, []byte, filesystem.FileMode) error }
func unsafeFactoryWrite(path string) error { write := unsafeWriterFactory(); return write(path, nil, 0o600) }
func unsafeFieldWrite(path string) error {
	return (unsafeOperations{write: filesystem.WriteFile}).write(path, nil, 0o600)
}
var unsafeInterfaceOutput io.Writer = filesystem.Stdout
var safeInterfaceOutput io.Writer = &bytes.Buffer{}
func unsafeInterfaceWrite() { _, _ = unsafeInterfaceOutput.Write(nil) }
func safeInterfaceWrite() { _, _ = safeInterfaceOutput.Write(nil) }
func unsafeConvertedInterfaceWrite() {
	output := io.Writer(filesystem.Stdout)
	_, _ = output.Write(nil)
}
func safeConvertedInterfaceWrite() {
	output := io.Writer(&bytes.Buffer{})
	_, _ = output.Write(nil)
}
func injectedInterfaceWriter(output io.Writer) { _, _ = output.Write(nil) }
func interfaceWriterFactory() io.Writer { return filesystem.Stdout }
func unsafeFactoryInterfaceWrite() { _, _ = interfaceWriterFactory().Write(nil) }
type interfaceDependencies struct { output io.Writer }
func unsafeFieldInterfaceWrite() {
	dependencies := interfaceDependencies{output: filesystem.Stdout}
	_, _ = dependencies.output.Write(nil)
}
func multiInterfaceWriter() (io.Writer, error) { return filesystem.Stdout, nil }
func unsafeMultiInterfaceWrite() { output, _ := multiInterfaceWriter(); _, _ = output.Write(nil) }
func namedInterfaceWriter() (output io.Writer) { output = filesystem.Stdout; return }
func unsafeNamedInterfaceWrite() { _, _ = namedInterfaceWriter().Write(nil) }
func unsafeAssignedFieldInterfaceWrite() {
	var dependencies interfaceDependencies
	dependencies.output = filesystem.Stdout
	_, _ = dependencies.output.Write(nil)
}
func unsafeIndexedInterfaceWrite() {
	outputs := []io.Writer{filesystem.Stdout}
	_, _ = outputs[0].Write(nil)
}
func unsafeRangedInterfaceWrite() {
	outputs := []io.Writer{filesystem.Stdout}
	for _, output := range outputs { _, _ = output.Write(nil) }
}
func unsafeAssertedInterfaceWrite() {
	output := any(filesystem.Stdout).(io.Writer)
	_, _ = output.Write(nil)
}
func unsafeCommaOKAssertedInterfaceWrite() {
	output, ok := any(filesystem.Stdout).(io.Writer)
	if ok { _, _ = output.Write(nil) }
}
func unsafeInterfaceCopy(source io.Reader) {
	var output io.Writer = filesystem.Stdout
	_, _ = io.Copy(output, source)
}
func injectedInterfaceCopy(destination io.Writer, source io.Reader) {
	_, _ = io.Copy(destination, source)
}
func unsafeClosureInterfaceWrite() {
	factory := func() io.Writer { return filesystem.Stdout }
	_, _ = factory().Write(nil)
}
func unsafeRawSyscall(path string) {
	fd, _ := syscall.Open(path, syscall.O_WRONLY|syscall.O_CREAT, 0o600)
	_, _ = syscall.Write(fd, nil)
}
func unsafeAsymmetricRenames(t *testing.T) {
	_ = syscall.Rename(filepath.Join(t.TempDir(), "source"), "LEAK-IN-REPOSITORY")
	_ = syscall.Rename("LEAK-IN-REPOSITORY", filepath.Join(t.TempDir(), "destination"))
}
func mutateTemporaryPath(path *string) { *path = "." }
func unsafeMutatedTemporaryPath(t *testing.T) {
	path := t.TempDir()
	mutateTemporaryPath(&path)
	_ = filesystem.WriteFile(filepath.Join(path, "LEAK"), nil, 0o600)
}
func unsafeSymlinkEscape(t *testing.T) {
	_ = filesystem.Symlink("..", filepath.Join(t.TempDir(), "repository"))
}
func unsafeEscapedTemporaryPath(t *testing.T) {
	path := t.TempDir()
	_ = json.Unmarshal([]byte(`"."`), &path)
	_ = filesystem.WriteFile(filepath.Join(path, "LEAK"), nil, 0o600)
}
func unsafeEnvironmentTemporaryPath(t *testing.T) {
	t.Setenv("TMPDIR", ".")
	_ = filesystem.WriteFile(filepath.Join(filesystem.TempDir(), "LEAK"), nil, 0o600)
}
func unsafeEnvironmentTestingPath(t *testing.T) {
	t.Setenv("TMPDIR", ".")
	_ = filesystem.WriteFile(filepath.Join(t.TempDir(), "LEAK"), nil, 0o600)
}
func poisonTemporaryEnvironment(t *testing.T) { t.Setenv("TMPDIR", ".") }
func unsafeHelperEnvironmentTestingPath(t *testing.T) {
	poisonTemporaryEnvironment(t)
	_ = filesystem.WriteFile(filepath.Join(t.TempDir(), "LEAK"), nil, 0o600)
}
func unsafeSyscallEnvironmentTestingPath(t *testing.T) {
	_ = syscall.Setenv("TMPDIR", ".")
	_ = filesystem.WriteFile(filepath.Join(t.TempDir(), "LEAK"), nil, 0o600)
}
EOF
cat >> "$q06_bad/go.mod" <<'EOF'

require golang.org/x/sys v0.0.0
replace golang.org/x/sys => ./external/xsys
EOF
mkdir -p "$q06_bad/external/xsys/unix"
cat > "$q06_bad/external/xsys/go.mod" <<'EOF'
module golang.org/x/sys
go 1.18
EOF
cat > "$q06_bad/external/xsys/unix/unix.go" <<'EOF'
package unix
const O_WRONLY = 1
const O_CREAT = 64
func ByteSliceFromString(value string) ([]byte, error) { return append([]byte(value), 0), nil }
func Open(path string, flags int, mode uint32) (int, error) { return -1, nil }
func Setenv(key, value string) error { return nil }
EOF
cat > "$q06_bad/xsys_test.go" <<'EOF'
package qualityfixture
import (
	"os"
	"path/filepath"
	"testing"
	"golang.org/x/sys/unix"
)
func pureXSysHelper() { _, _ = unix.ByteSliceFromString("pure") }
func unsafeXSysOpen() { _, _ = unix.Open("LEAK", unix.O_WRONLY|unix.O_CREAT, 0o600) }
func unsafeXSysEnvironment(t *testing.T) {
	_ = unix.Setenv("TMPDIR", ".")
	_ = os.WriteFile(filepath.Join(t.TempDir(), "LEAK"), nil, 0o600)
}
EOF
cat > "$q06_bad/file_type_test.go" <<'EOF'
package qualityfixture
import "os"
type auditOutputFile = os.File
EOF
cat > "$q06_bad/file_write_test.go" <<'EOF'
package qualityfixture
func unsafeCrossFileWrite(file *auditOutputFile) { _, _ = file.Write(nil) }
EOF
mkdir -p "$q06_bad/internal/filetype"
cat > "$q06_bad/internal/filetype/file.go" <<'EOF'
package filetype
import "os"
type OutputFile = os.File
EOF
cat > "$q06_bad/cross_package_write_test.go" <<'EOF'
package qualityfixture
import "example.invalid/qualityfixture/internal/filetype"
func unsafeCrossPackageWrite(file *filetype.OutputFile) { _, _ = file.Write(nil) }
EOF
mkdir -p "$q06_bad/internal/resource"
cat > "$q06_bad/internal/resource/output.go" <<'EOF'
package resource
import (
	"bytes"
	"io"
	"net"
	"os"
)
var Output io.Writer = os.Stdout
func OutputFactory() io.Writer { return os.Stdout }
EOF
cat > "$q06_bad/cross_package_interface_write_test.go" <<'EOF'
package qualityfixture
import "example.invalid/qualityfixture/internal/resource"
func unsafeImportedInterfaceWrite() { _, _ = resource.Output.Write(nil) }
func unsafeImportedFactoryWrite() { _, _ = resource.OutputFactory().Write(nil) }
EOF
cat > "$q06_bad/copy_test.go" <<'EOF'
package qualityfixture
import (
	"io"
	"os"
)
func unsafeCopy(file *os.File, source io.Reader) { _, _ = io.Copy(file, source) }
EOF
cat > "$q06_bad/buffer_test.go" <<'EOF'
package qualityfixture
import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)
func unrelatedBufferWrite(buffer *bytes.Buffer) {
	_, _ = buffer.Write(nil)
	_, _ = io.Copy(buffer, strings.NewReader("memory only"))
	_ = os.ErrNotExist
}
func safeTemporaryDirectories(t *testing.T) {
	_, _ = os.CreateTemp(t.TempDir(), "fixture-")
	_, _ = os.MkdirTemp("", "system-temp-")
	directory := t.TempDir()
	file, _ := os.CreateTemp(directory, "aliased-")
	if file != nil {
		_, _ = file.Write(nil)
		alias := file
		_, _ = alias.Write(nil)
		var output io.Writer = alias
		_, _ = io.Copy(output, strings.NewReader("temp only"))
		_ = file.Close()
	}
	_ = os.WriteFile(filepath.Join(directory, "fixture"), nil, 0o600)
	root, _ := os.OpenRoot(directory)
	if root != nil {
		_ = root.WriteFile("root-fixture", nil, 0o600)
		_ = root.Close()
	}
	readOnly, _ := os.OpenFile("fixture.txt", os.O_RDONLY, 0)
	if readOnly != nil { _ = readOnly.Close() }
	fd, _ := syscall.Open("fixture.txt", syscall.O_RDONLY, 0)
	if fd >= 0 { _ = syscall.Close(fd) }
	made, _ := os.MkdirTemp("", "made-")
	if made != "" {
		defer os.RemoveAll(made)
		_ = os.WriteFile(filepath.Join(made, "fixture"), nil, 0o600)
	}
	_ = os.WriteFile(t.TempDir()+"/concatenated", nil, 0o600)
	_ = os.Symlink("fixture-target", filepath.Join(t.TempDir(), "link"))
}
type temporaryDependencies struct { directory string }
func temporaryFactory(t *testing.T) *temporaryDependencies {
	return &temporaryDependencies{directory: t.TempDir()}
}
func safeTemporaryFactoryField(t *testing.T) {
	_ = os.WriteFile(filepath.Join(temporaryFactory(t).directory, "fixture"), nil, 0o600)
}
func safeTestingTB(test testing.TB) {
	_ = os.WriteFile(filepath.Join(test.TempDir(), "fixture"), nil, 0o600)
}
func safeTestingB(benchmark *testing.B) {
	_ = os.WriteFile(filepath.Join(benchmark.TempDir(), "fixture"), nil, 0o600)
}
func safeTestingF(fuzz *testing.F) {
	_ = os.WriteFile(filepath.Join(fuzz.TempDir(), "fixture"), nil, 0o600)
}
EOF
export RUN_REPO="$q06_bad" REPORT_SHORT="$(git -C "$q06_bad" rev-parse --short HEAD)"
write_report Q0.6 PASS "leak guard used at 3 call sites; 0 skipped tests"
run_wrapper "$WORK/q06-unsafe"
rc=$?
[ "$rc" -eq 1 ] || fail "unsafe direct test write returned $rc, expected 1"
json_assert "$WORK/q06-unsafe/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['unsafe_direct_test_writes'] == 46 and len([v for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites'] if v['path'] == 'unsafe_test.go']) == 35 and len([v for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites'] if v['path'] == 'cross_package_interface_write_test.go']) == 2 and len([v for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites'] if v['path'] == 'xsys_test.go']) == 2 and len([v for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites'] if v['path'] == 'internal/testutil/leak.go']) == 1 and len([v for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites'] if v['path'] == 'internal/testutil/path.go']) == 1 and len([v for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites'] if v['path'] == 'internal/testutil/path_test.go']) == 2 and any(v['path'] == 'file_write_test.go' for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites']) and any(v['path'] == 'cross_package_write_test.go' for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites']) and any(v['path'] == 'copy_test.go' for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites']) and all(v['path'] != 'buffer_test.go' for v in value['criteria'][0]['metric']['unsafe_direct_test_write_sites']) and not value['criteria'][0]['metric']['precondition_met']" \
  "Q0.6 did not resolve direct/escaping/interface/raw-syscall effects or incorrectly counted memory-only writes"
rm "$q06_bad/unsafe_test.go" "$q06_bad/file_type_test.go" "$q06_bad/file_write_test.go" \
  "$q06_bad/cross_package_write_test.go" "$q06_bad/cross_package_interface_write_test.go" \
  "$q06_bad/copy_test.go" "$q06_bad/buffer_test.go" "$q06_bad/xsys_test.go" \
  "$q06_bad/helper_test.go" "$q06_bad/internal/testutil/leak.go" \
  "$q06_bad/internal/filetype/file.go" \
  "$q06_bad/internal/resource/output.go"
rmdir "$q06_bad/internal/filetype" "$q06_bad/internal/resource"
python3 - "$q06_bad/internal/testutil/path.go" <<'PY'
import pathlib
import sys

path = pathlib.Path(sys.argv[1])
value = path.read_text(encoding="utf-8")
path.write_text(value.replace("if WouldLeakIntoWorkingTree(path)", "if false"), encoding="utf-8")
PY
run_wrapper "$WORK/q06-bypassed-guard"
rc=$?
[ "$rc" -eq 1 ] || fail "bypassed central guard returned $rc, expected 1"
json_assert "$WORK/q06-bypassed-guard/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['central_guard_call_sites'] == 0 and not value['criteria'][0]['metric']['precondition_met']" \
  "Q0.6 accepted a bypassed central guard"
q06_receiver_controls="$WORK/q06-receiver-controls-repo"
git clone -q "$WORK/repo" "$q06_receiver_controls"
cat > "$q06_receiver_controls/internal/testutil/path_test.go" <<'EOF'
package testutil
import "testing"
type falseControl struct{}
func (falseControl) TestWouldLeakIntoRepositoryRejectsWorkingTreeAndAcceptsTempDir(t *testing.T) {}
func (falseControl) TestSafeFixtureWritersRefuseRepositoryPaths(t *testing.T) {}
EOF
export RUN_REPO="$q06_receiver_controls" REPORT_SHORT="$(git -C "$q06_receiver_controls" rev-parse --short HEAD)"
run_wrapper "$WORK/q06-receiver-controls"
rc=$?
[ "$rc" -eq 1 ] || fail "receiver control methods returned $rc, expected 1"
json_assert "$WORK/q06-receiver-controls/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['repository_control_tests_passed'] == 0 and not value['criteria'][0]['metric']['precondition_met']" \
  "receiver methods with control-test names were treated as executed package tests"
q06_tagged_controls="$WORK/q06-tagged-controls-repo"
git clone -q "$WORK/repo" "$q06_tagged_controls"
python3 - "$q06_tagged_controls/internal/testutil/path_test.go" <<'PY'
import pathlib
import sys

path = pathlib.Path(sys.argv[1])
path.write_text("//go:build quality_never\n\n" + path.read_text(encoding="utf-8"), encoding="utf-8")
PY
export RUN_REPO="$q06_tagged_controls" REPORT_SHORT="$(git -C "$q06_tagged_controls" rev-parse --short HEAD)"
run_wrapper "$WORK/q06-tagged-controls"
rc=$?
[ "$rc" -eq 1 ] || fail "inactive tagged controls returned $rc, expected 1"
json_assert "$WORK/q06-tagged-controls/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['repository_control_tests_passed'] == 0 and not value['criteria'][0]['metric']['precondition_met']" \
  "exit 0 with no selected control tests was treated as two passing controls"
q06_dishonest="$WORK/q06-dishonest-repo"
git clone -q "$WORK/repo" "$q06_dishonest"
cat > "$q06_dishonest/internal/testutil/path.go" <<'EOF'
package testutil
import "io/fs"
func WouldLeakIntoWorkingTree(path string) bool { return false }
func ensureOutsideWorkingTree(path string) error {
	// WouldLeakIntoWorkingTree(path)
	return nil
}
func WriteFileOutsideWorkingTree(path string, value []byte, mode fs.FileMode) error {
	// ensureOutsideWorkingTree(path)
	return nil
}
func CopyFSOutsideWorkingTree(path string, source fs.FS) error {
	// ensureOutsideWorkingTree(path)
	return nil
}
EOF
cat > "$q06_dishonest/internal/testutil/path_test.go" <<'EOF'
package testutil
import "testing"
func TestWouldLeakIntoRepositoryRejectsWorkingTreeAndAcceptsTempDir(t *testing.T) {}
func TestSafeFixtureWritersRefuseRepositoryPaths(t *testing.T) {}
EOF
export RUN_REPO="$q06_dishonest" REPORT_SHORT="$(git -C "$q06_dishonest" rev-parse --short HEAD)"
run_wrapper "$WORK/q06-dishonest"
rc=$?
[ "$rc" -eq 1 ] || fail "dishonest Q0.6 contract returned $rc, expected 1"
json_assert "$WORK/q06-dishonest/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['contract_tests_passed'] == 0 and not value['criteria'][0]['metric']['precondition_met']" \
  "Q0.6 accepted empty controls and an ineffective guard"
q06_shadow="$WORK/q06-shadow-repo"
git clone -q "$WORK/repo" "$q06_shadow"
cat > "$q06_shadow/dot_shadow_test.go" <<'EOF'
package qualityfixture
import (
	. "example.invalid/qualityfixture/internal/testutil"
	"os"
)
var _ = CopyFSOutsideWorkingTree
func unsafeDotShadow(path string) error {
	WriteFileOutsideWorkingTree := os.WriteFile
	return WriteFileOutsideWorkingTree(path, nil, 0o600)
}
EOF
cat > "$q06_shadow/selector_shadow_test.go" <<'EOF'
package qualityfixture
import (
	"io/fs"
	tu "example.invalid/qualityfixture/internal/testutil"
)
var _ = tu.WriteFileOutsideWorkingTree
type fakeSafeWriters struct {
	WriteFileOutsideWorkingTree func(string, []byte, fs.FileMode) error
	CopyFSOutsideWorkingTree func(string, fs.FS) error
}
func fakeSafeWriterNames(path string) {
	tu := fakeSafeWriters{
		WriteFileOutsideWorkingTree: func(string, []byte, fs.FileMode) error { return nil },
		CopyFSOutsideWorkingTree: func(string, fs.FS) error { return nil },
	}
	_ = tu.WriteFileOutsideWorkingTree(path, nil, 0o600)
	_ = tu.CopyFSOutsideWorkingTree(path, nil)
}
EOF
export RUN_REPO="$q06_shadow" REPORT_SHORT="$(git -C "$q06_shadow" rev-parse --short HEAD)"
write_report Q0.6 PASS "leak guard used at 4 call sites; 0 skipped tests"
run_wrapper "$WORK/q06-shadow"
rc=$?
[ "$rc" -eq 1 ] || fail "shadowed safe-writer names returned $rc, expected 1"
json_assert "$WORK/q06-shadow/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['unsafe_direct_test_writes'] == 1 and value['criteria'][0]['metric']['write_safe_call_sites'] == 2 and value['criteria'][0]['metric']['copy_safe_call_sites'] == 2" \
  "Q0.6 trusted shadowed safe-writer spellings instead of resolved function identities"
unset RUN_REPO REPORT_SHORT
ok "Q0.6 rejects unsafe writes, bypassed guards, shadowed names, and empty dishonest controls"

printf 'T10 local Q1.3 and Q2.1 measurements enforce exact inventory contracts\n'
missing_adapter_repo="$WORK/missing-adapter-repo"
git clone -q "$WORK/repo" "$missing_adapter_repo"
printf '\ninternal/adapter/missing\n' >> "$missing_adapter_repo/.quality/inventory"
export RUN_REPO="$missing_adapter_repo" REPORT_SHORT="$(git -C "$missing_adapter_repo" rev-parse --short HEAD)"
write_report Q1.3 PASS "0 direct calls outside 2 declared adapters"
run_wrapper "$WORK/missing-adapter"
rc=$?
[ "$rc" -eq 1 ] || fail "missing declared adapter returned $rc, expected 1"
json_assert "$WORK/missing-adapter/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['value'] == 0 and value['criteria'][0]['metric']['missing_adapter_paths'] == ['internal/adapter/missing']" \
  "Q1.3 passed while a declared adapter path did not exist"
python3 - "$WORK/baseline.json" "$WORK/missing-adapter-baseline.json" \
  "$missing_adapter_repo/.quality/inventory" <<'PY'
import hashlib
import json
import sys

value = json.load(open(sys.argv[1], encoding="utf-8"))
metric = next(item["metric"] for item in value["criteria"] if item["id"] == "Q1.3")
metric["invalid_adapter_paths"] = [{"path": "internal/adapter/missing", "status": "missing"}]
value["inventory"]["sha256"] = hashlib.sha256(open(sys.argv[3], "rb").read()).hexdigest()
with open(sys.argv[2], "w", encoding="utf-8") as stream:
    json.dump(value, stream)
PY
run_wrapper "$WORK/missing-adapter-ratchet" --baseline "$WORK/missing-adapter-baseline.json"
rc=$?
[ "$rc" -eq 1 ] || fail "stored Q1.3 debt returned $rc, expected finding exit 1"
json_assert "$WORK/missing-adapter-ratchet/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['ratchet']['status'] == 'not_comparable' and value['ratchet']['current_not_comparable'] == 1" \
  "stored missing-adapter debt broke the audit or falsely passed the current invalid precondition"
unset RUN_REPO REPORT_SHORT
broad_adapter_repo="$WORK/broad-adapter-repo"
git clone -q "$WORK/repo" "$broad_adapter_repo"
python3 - "$broad_adapter_repo/.quality/inventory" <<'PY'
from pathlib import Path
import sys

path = Path(sys.argv[1])
path.write_text(path.read_text(encoding="utf-8").replace("internal/adapter/process", "internal"), encoding="utf-8")
PY
export RUN_REPO="$broad_adapter_repo" REPORT_SHORT="$(git -C "$broad_adapter_repo" rev-parse --short HEAD)"
write_report Q1.3 PASS "0 direct calls outside 1 declared adapters"
run_wrapper "$WORK/broad-adapter"
rc=$?
[ "$rc" -eq 1 ] || fail "broad ancestor adapter returned $rc, expected 1"
json_assert "$WORK/broad-adapter/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['value'] == 1 and value['criteria'][0]['metric']['invalid_adapter_paths'][0]['status'] == 'no_selected_production_package'" \
  "Q1.3 accepted a broad ancestor directory as an exact production adapter package"
unset RUN_REPO REPORT_SHORT
invalid_subject_repo="$WORK/invalid-subject-ratchet-repo"
git clone -q "$WORK/repo" "$invalid_subject_repo"
python3 - "$invalid_subject_repo/.quality/inventory" <<'PY'
from pathlib import Path
import sys

path = Path(sys.argv[1])
path.write_text(path.read_text(encoding="utf-8").replace("fixture = . :", "fixture = .quality :"), encoding="utf-8")
PY
python3 - "$WORK/baseline.json" "$WORK/invalid-subject-baseline.json" \
  "$invalid_subject_repo/.quality/inventory" <<'PY'
import hashlib
import json
import sys

value = json.load(open(sys.argv[1], encoding="utf-8"))
metric = next(item["metric"] for item in value["criteria"] if item["id"] == "Q2.1")
metric.update({"numerator": 0, "denominator": 1, "source_precondition_met": True})
value["inventory"]["sha256"] = hashlib.sha256(open(sys.argv[3], "rb").read()).hexdigest()
with open(sys.argv[2], "w", encoding="utf-8") as stream:
    json.dump(value, stream)
PY
export RUN_REPO="$invalid_subject_repo" REPORT_SHORT="$(git -C "$invalid_subject_repo" rev-parse --short HEAD)"
write_report Q2.1 PASS "0 harnesses for 1 declared subjects of 1 packages"
run_wrapper "$WORK/invalid-subject-ratchet" --baseline "$WORK/invalid-subject-baseline.json"
rc=$?
[ "$rc" -eq 1 ] || fail "invalid subject source with baseline returned $rc, expected finding exit 1"
json_assert "$WORK/invalid-subject-ratchet/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and not value['criteria'][0]['metric']['source_precondition_met'] and value['ratchet']['status'] == 'not_comparable' and value['ratchet']['current_not_comparable'] == 1" \
  "Q2.1 ratchet upgraded an invalid current subject source to PASS"
unset RUN_REPO REPORT_SHORT
scan_repo="$WORK/scan-repo"
mkdir -p "$scan_repo/.quality" "$scan_repo/internal/adapter/process" \
  "$scan_repo/internal/filetype" "$scan_repo/internal/testutil" "$scan_repo/internal/testutilish" \
  "$scan_repo/internal/pkg" "$scan_repo/internal/testonly" "$scan_repo/scripts" \
  "$scan_repo/target" "$scan_repo/vendor"
cat > "$scan_repo/go.mod" <<'EOF'
module example.invalid/scanfixture

go 1.18

require (
	example.invalid/prompt v0.0.0
	example.invalid/scanfixture/plugin v0.0.0
	golang.org/x/sys v0.0.0
)
replace example.invalid/prompt => ./external/prompt
replace example.invalid/scanfixture/plugin => ./external/plugin
replace golang.org/x/sys => ./external/xsys
EOF
mkdir -p "$scan_repo/external/prompt"
cat > "$scan_repo/external/prompt/go.mod" <<'EOF'
module example.invalid/prompt
go 1.18
EOF
cat > "$scan_repo/external/prompt/prompt.go" <<'EOF'
package prompt
type Prompt struct{}
func (*Prompt) Run() (string, error) { return "", nil }
type Worker struct{}
func Command(string) *Worker { return &Worker{} }
func Now() int { return 0 }
func (*Worker) Write([]byte) (int, error) { return 0, nil }
func (*Worker) Do() {}
func (*Worker) Get() {}
EOF
mkdir -p "$scan_repo/external/plugin"
cat > "$scan_repo/external/plugin/go.mod" <<'EOF'
module example.invalid/scanfixture/plugin
go 1.18
EOF
cat > "$scan_repo/external/plugin/plugin.go" <<'EOF'
package plugin
type Buffer struct{}
func New() *Buffer { return &Buffer{} }
func (*Buffer) Write([]byte) (int, error) { return 0, nil }
EOF
mkdir -p "$scan_repo/external/xsys/unix"
cat > "$scan_repo/external/xsys/go.mod" <<'EOF'
module golang.org/x/sys
go 1.18
EOF
cat > "$scan_repo/external/xsys/unix/unix.go" <<'EOF'
package unix
const O_WRONLY = 1
const O_CREAT = 64
func ByteSliceFromString(value string) ([]byte, error) { return append([]byte(value), 0), nil }
func Major(value uint64) uint32 { return uint32(value) }
func Open(path string, flags int, mode uint32) (int, error) { return -1, nil }
EOF
cat > "$scan_repo/.quality/inventory" <<'EOF'
[subjects]
ready   = internal/pkg : scripts/mutate-ready
not-executable = internal/pkg : scripts/mutate-not-executable
missing = internal/pkg : scripts/mutate-missing
missing-source = does/not/exist : scripts/mutate-ready
test-only-source = internal/testonly : scripts/mutate-ready
outside-source = ../../outside : scripts/mutate-ready
quality-source = .quality/fake.go : scripts/mutate-ready
target-source = target/fake.go : scripts/mutate-ready
vendor-source = vendor/fake.go : scripts/mutate-ready

[adapters]
internal/adapter/process
EOF
printf 'package qualityfake\nfunc Value() {}\n' > "$scan_repo/.quality/fake.go"
printf 'package targetfake\nfunc Value() {}\n' > "$scan_repo/target/fake.go"
printf 'package vendorfake\nfunc Value() {}\n' > "$scan_repo/vendor/fake.go"
cat > "$scan_repo/internal/testonly/only_test.go" <<'EOF'
package testonly
func testOnlyHelper() {}
EOF
cat > "$scan_repo/internal/adapter/process/run.go" <<'EOF'
package process
import "os/exec"
func Run() { _ = exec.Command("true") }
EOF
cat > "$scan_repo/internal/testutil/write.go" <<'EOF'
package testutil
import "os"
func Write(path string) error { return os.WriteFile(path, nil, 0o600) }
EOF
cat > "$scan_repo/internal/pkg/write.go" <<'EOF'
package pkg
import "os"
func Write(path string) error { return os.WriteFile(path, nil, 0o600) }
func WriteValue(path string) error { write := os.WriteFile; return write(path, nil, 0o600) }
EOF
cat > "$scan_repo/internal/pkg/write_test.go" <<'EOF'
package pkg
import "os"
func ignoredTestWrite(path string) error { return os.WriteFile(path, nil, 0o600) }
EOF
cat > "$scan_repo/internal/testutilish/write.go" <<'EOF'
package testutilish
import "os"
func Write(path string) error { return os.WriteFile(path, nil, 0o600) }
EOF
cat > "$scan_repo/internal/pkg/client.go" <<'EOF'
package pkg
func Do(client *HTTPClient) { _, _ = client.Do(nil) }
func Get(client *HTTPClient) { _, _ = client.Get("https://example.invalid") }
func GetValue(client *HTTPClient) { get := client.Get; _, _ = get("https://example.invalid") }
EOF
cat > "$scan_repo/internal/pkg/types.go" <<'EOF'
package pkg
import "net/http"
type HTTPClient = http.Client
EOF
cat > "$scan_repo/internal/pkg/file_types.go" <<'EOF'
package pkg
import "os"
type OutputFile = os.File
EOF
cat > "$scan_repo/internal/pkg/file_write.go" <<'EOF'
package pkg
func FileWrite(file *OutputFile) { _, _ = file.Write(nil) }
func FileWriteValue(file *OutputFile) { write := file.Write; _, _ = write(nil) }
EOF
cat > "$scan_repo/internal/filetype/file.go" <<'EOF'
package filetype
import "os"
type OutputFile = os.File
EOF
cat > "$scan_repo/internal/pkg/cross_package_write.go" <<'EOF'
package pkg
import "example.invalid/scanfixture/internal/filetype"
func CrossPackageWrite(file *filetype.OutputFile) { _, _ = file.Write(nil) }
EOF
cat > "$scan_repo/internal/pkg/copy.go" <<'EOF'
package pkg
import (
	"io"
	"os"
)
func Copy(file *os.File, source io.Reader) { _, _ = io.Copy(file, source) }
func CopyValue(file *os.File, source io.Reader) { copyToFile := io.Copy; _, _ = copyToFile(file, source) }
EOF
cat > "$scan_repo/internal/pkg/modern_os.go" <<'EOF'
package pkg
import "os"
func CopyFS(path string) error { return os.CopyFS(path, nil) }
func OpenRoot(path string) (*os.Root, error) { return os.OpenRoot(path) }
func RootWrite(root *os.Root, path string) error { return root.WriteFile(path, nil, 0o600) }
EOF
cat > "$scan_repo/internal/pkg/higher_order.go" <<'EOF'
package pkg
import "os"
var ExportedWriter = os.WriteFile
func WriterFactory() func(string, []byte, os.FileMode) error { return os.WriteFile }
type operations struct { write func(string, []byte, os.FileMode) error }
func FactoryWrite(path string) error { write := WriterFactory(); return write(path, nil, 0o600) }
func FieldWrite(path string) error { return (operations{write: os.WriteFile}).write(path, nil, 0o600) }
EOF
cat > "$scan_repo/internal/pkg/interface_write.go" <<'EOF'
package pkg
import (
	"bytes"
	"io"
	"net/http"
	"os"
)
var interfaceFileOutput io.Writer = os.Stdout
var interfaceMemoryOutput io.Writer = &bytes.Buffer{}
func InterfaceFileWrite() { _, _ = interfaceFileOutput.Write(nil) }
func InterfaceMemoryWrite() { _, _ = interfaceMemoryOutput.Write(nil) }
func InterfaceConvertedFileWrite() {
	output := io.Writer(os.Stdout)
	_, _ = output.Write(nil)
}
func InterfaceConvertedMemoryWrite() {
	output := io.Writer(&bytes.Buffer{})
	_, _ = output.Write(nil)
}
func InterfaceWriterSeam(output io.Writer) { _, _ = output.Write(nil) }
func InterfaceReaderSeam(input io.Reader) { _, _ = input.Read(nil) }
func interfaceWriterFactory() io.Writer { return os.Stdout }
func InterfaceFactoryWrite() { _, _ = interfaceWriterFactory().Write(nil) }
type interfaceDependencies struct { output io.Writer }
func InterfaceFieldWrite() {
	dependencies := interfaceDependencies{output: os.Stdout}
	_, _ = dependencies.output.Write(nil)
}
type httpDoer interface { Do(*http.Request) (*http.Response, error) }
var interfaceHTTPClient httpDoer = http.DefaultClient
var interfaceHTTPTransport http.RoundTripper = http.DefaultTransport
func InterfaceHTTPDo(request *http.Request) { _, _ = interfaceHTTPClient.Do(request) }
func InterfaceHTTPRoundTrip(request *http.Request) { _, _ = interfaceHTTPTransport.RoundTrip(request) }
func InterfaceHTTPSeam(doer httpDoer, request *http.Request) { _, _ = doer.Do(request) }
EOF
cat > "$scan_repo/internal/pkg/raw_syscall.go" <<'EOF'
package pkg
import "syscall"
func RawSyscallWrite(path string) {
	fd, _ := syscall.Open(path, syscall.O_WRONLY|syscall.O_CREAT, 0o600)
	_, _ = syscall.Write(fd, nil)
}
func RawSyscallPure() int {
	_ = syscall.O_RDONLY
	return syscall.Getpid()
}
EOF
cat > "$scan_repo/internal/pkg/origin_variants.go" <<'EOF'
package pkg
import (
	"bytes"
	"io"
	"os"
)
func multiWriter() (io.Writer, error) { return os.Stdout, nil }
func MultiWriterUse() { output, _ := multiWriter(); _, _ = output.Write(nil) }
func namedWriter() (output io.Writer) { output = os.Stdout; return }
func NamedWriterUse() { _, _ = namedWriter().Write(nil) }
type originDependencies struct { output io.Writer }
func AssignedFieldWriterUse() {
	var dependencies originDependencies
	dependencies.output = os.Stdout
	_, _ = dependencies.output.Write(nil)
}
func IndexedWriterUse() {
	outputs := []io.Writer{os.Stdout}
	_, _ = outputs[0].Write(nil)
}
func RangedWriterUse() {
	outputs := []io.Writer{os.Stdout}
	for _, output := range outputs { _, _ = output.Write(nil) }
}
func AssertedWriterUse() {
	output := any(os.Stdout).(io.Writer)
	_, _ = output.Write(nil)
}
func CommaOKAssertedWriterUse() {
	output, ok := any(os.Stdout).(io.Writer)
	if ok { _, _ = output.Write(nil) }
}
func InterfaceCopyUse(source io.Reader) {
	var output io.Writer = os.Stdout
	_, _ = io.Copy(output, source)
}
func MemoryInterfaceCopyUse(source io.Reader) {
	var output io.Writer = &bytes.Buffer{}
	_, _ = io.Copy(output, source)
}
func ClosureWriterUse() {
	factory := func() io.Writer { return os.Stdout }
	_, _ = factory().Write(nil)
}
func compositeWriterFactory() *originDependencies {
	return &originDependencies{output: os.Stdout}
}
func FactoryFieldWriterUse() { _, _ = compositeWriterFactory().output.Write(nil) }
func writerSliceFactory() []io.Writer { return []io.Writer{os.Stdout} }
func FactoryElementWriterUse() { _, _ = writerSliceFactory()[0].Write(nil) }
func AppendedWriterUse() {
	outputs := []io.Writer{&bytes.Buffer{}}
	outputs = append(outputs, os.Stdout)
	_, _ = outputs[1].Write(nil)
}
func MutatedCollectionWriterUse() {
	outputs := map[string]io.Writer{}
	outputs["real"] = os.Stdout
	_, _ = outputs["real"].Write(nil)
}
func PointerFieldMutationWriterUse() {
	dependencies := originDependencies{output: &bytes.Buffer{}}
	pointer := &dependencies
	pointer.output = os.Stdout
	_, _ = dependencies.output.Write(nil)
}
func mutateOriginDependencies(dependencies *originDependencies) { dependencies.output = os.Stdout }
func CalledFieldMutationWriterUse() {
	dependencies := originDependencies{output: &bytes.Buffer{}}
	mutateOriginDependencies(&dependencies)
	_, _ = dependencies.output.Write(nil)
}
func MemoryMutationControls() {
	dependencies := originDependencies{output: &bytes.Buffer{}}
	pointer := &dependencies
	pointer.output = &bytes.Buffer{}
	_, _ = dependencies.output.Write(nil)
	outputs := []io.Writer{&bytes.Buffer{}}
	outputs = append(outputs, &bytes.Buffer{})
	_, _ = outputs[1].Write(nil)
}
type deterministicDependencies struct { output io.Writer }
func DeterministicOriginKinds() {
	dependencies := deterministicDependencies{output: &bytes.Buffer{}}
	left := &dependencies
	right := &dependencies
	left.output = os.Stdout
	connection := &net.TCPConn{}
	right.output = connection
	_, _ = dependencies.output.Write(nil)
	outputs := []io.Writer{&bytes.Buffer{}}
	leftOutputs := outputs
	rightOutputs := outputs
	leftOutputs[0] = os.Stdout
	rightOutputs[0] = connection
	_, _ = outputs[0].Write(nil)
}
EOF
cat > "$scan_repo/internal/pkg/catalog_effects.go" <<'EOF'
package pkg
import (
	"context"
	"crypto/tls"
	"database/sql"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"testing/fstest"
	"time"
)
func CatalogEffects(database *sql.DB, dialer *tls.Dialer, resolver *net.Resolver) {
	_ = (&exec.Cmd{Path: "/bin/true"}).Run()
	_, _ = net.LookupHost("example.invalid")
	_, _ = tls.Dial("tcp", "example.invalid:443", nil)
	_, _ = dialer.Dial("tcp", "example.invalid:443")
	_, _ = dialer.DialContext(context.Background(), "tcp", "example.invalid:443")
	_, _ = resolver.LookupIPAddr(context.Background(), "example.invalid")
	_, _ = fs.ReadFile(os.DirFS("."), "fixture")
	time.AfterFunc(time.Second, func() {})
	_, _ = database.Exec("select 1")
}
func InMemoryFSIsNotExternal() {
	_, _ = fs.ReadFile(fstest.MapFS{"fixture": &fstest.MapFile{Data: []byte("ok")}}, "fixture")
}
EOF
mkdir -p "$scan_repo/internal/resource" "$scan_repo/internal/resourceuse"
cat > "$scan_repo/internal/resource/output.go" <<'EOF'
package resource
import (
	"bytes"
	"io"
	"os"
)
var Output io.Writer = os.Stdout
func OutputFactory() io.Writer { return os.Stdout }
type Dependencies struct { Output io.Writer }
var Default = Dependencies{Output: os.Stdout}
var Outputs = []io.Writer{os.Stdout}
var MemoryDefault = Dependencies{Output: &bytes.Buffer{}}
var MemoryOutputs = []io.Writer{&bytes.Buffer{}}
EOF
cat > "$scan_repo/internal/resourceuse/write.go" <<'EOF'
package resourceuse
import "example.invalid/scanfixture/internal/resource"
func ImportedInterfaceWrite() { _, _ = resource.Output.Write(nil) }
func ImportedFactoryWrite() { _, _ = resource.OutputFactory().Write(nil) }
func ImportedCompositeWrite() { _, _ = resource.Default.Output.Write(nil) }
func ImportedElementWrite() { _, _ = resource.Outputs[0].Write(nil) }
func ImportedMemoryControls() {
	_, _ = resource.MemoryDefault.Output.Write(nil)
	_, _ = resource.MemoryOutputs[0].Write(nil)
}
EOF
cat > "$scan_repo/internal/pkg/cgo_effect.go" <<'EOF'
package pkg
/*
#include <stdio.h>
static void quality_open(void) {
	FILE *file = fopen("unsafe.txt", "w");
	if (file != NULL) { fclose(file); }
}
*/
import "C"
func CgoWrite() { C.quality_open() }
EOF
cat > "$scan_repo/internal/pkg/unrelated.go" <<'EOF'
package pkg
import (
	"bytes"
	"io"
	"strings"
)
type worker struct{}
func (worker) Do() {}
func IgnoreWorkerDo(value worker) { value.Do() }
func IgnoreBufferWrite(value *bytes.Buffer) { _, _ = value.Write(nil) }
func IgnoreMemoryCopy(output *bytes.Buffer) { _, _ = io.Copy(output, strings.NewReader("memory only")) }
func IgnoreInjectedCopy(output io.Writer, source io.Reader) { _, _ = io.Copy(output, source) }
EOF
cat > "$scan_repo/internal/pkg/argument_origins.go" <<'EOF'
package pkg
import (
	"bytes"
	"io"
	"net/http"
	"os"
)
type argumentDoer interface { Do(*http.Request) (*http.Response, error) }
func writeArgument(output io.Writer) { _, _ = output.Write(nil) }
func requestArgument(client argumentDoer, request *http.Request) { _, _ = client.Do(request) }
func KnownArguments(request *http.Request) {
	writeArgument(os.Stdout)
	requestArgument(http.DefaultClient, request)
	writeArgument(&bytes.Buffer{})
}
func ForwardInjectedArguments(output io.Writer, client argumentDoer, request *http.Request) {
	writeArgument(output)
	requestArgument(client, request)
}
func ConcreteInjectedClient(client *http.Client, request *http.Request) { _, _ = client.Do(request) }
func KnownConcreteClient(request *http.Request) { ConcreteInjectedClient(http.DefaultClient, request) }
func LiteralArguments(request *http.Request) {
	func(client argumentDoer) { _, _ = client.Do(request) }(http.DefaultClient)
	call := func(client argumentDoer) { _, _ = client.Do(request) }
	call(http.DefaultClient)
}
type argumentDependencies struct { client argumentDoer; output io.Writer }
func useArgumentDependencies(dependencies argumentDependencies, request *http.Request) {
	_, _ = dependencies.client.Do(request)
	_, _ = dependencies.output.Write(nil)
}
func CompositeArguments(request *http.Request) {
	dependencies := argumentDependencies{client: http.DefaultClient, output: &bytes.Buffer{}}
	useArgumentDependencies(dependencies, request)
}
type methodDependencies struct { output io.Writer }
func (dependencies methodDependencies) run() { _, _ = dependencies.output.Write(nil) }
func usePointerDependencies(dependencies *methodDependencies) { _, _ = dependencies.output.Write(nil) }
func MethodArgumentOrigins() {
	(methodDependencies{output: os.Stdout}).run()
	dependencies := methodDependencies{output: os.Stdout}
	dependencies.run()
	run := dependencies.run
	run()
	usePointerDependencies(&methodDependencies{output: os.Stdout})
	(methodDependencies{output: &bytes.Buffer{}}).run()
}
type concreteFileDependencies struct { output *os.File }
func injectedFileField(dependencies concreteFileDependencies) { _, _ = dependencies.output.Write(nil) }
func identityFile(output *os.File) *os.File { return output }
func InjectedFileControls(dependencies concreteFileDependencies, output *os.File) {
	injectedFileField(dependencies)
	_, _ = identityFile(output).Write(nil)
	_, _ = identityFile(os.Stdout).Write(nil)
}
EOF
cat > "$scan_repo/internal/pkg/xsys.go" <<'EOF'
package pkg
import "golang.org/x/sys/unix"
func XSysEffects(path string) {
	_, _ = unix.ByteSliceFromString("pure")
	_ = unix.Major(0)
	_, _ = unix.Open(path, unix.O_WRONLY|unix.O_CREAT, 0o600)
}
EOF
cat > "$scan_repo/internal/pkg/third_party.go" <<'EOF'
package pkg
import "example.invalid/prompt"
func ResolvedThirdPartyTuple(instance *prompt.Prompt) string {
	value, _ := instance.Run()
	worker := prompt.Command("safe")
	_, _ = worker.Write(nil)
	worker.Do()
	worker.Get()
	_ = prompt.Now()
	return value
}
EOF
cat > "$scan_repo/internal/pkg/nested_third_party.go" <<'EOF'
package pkg
import "example.invalid/scanfixture/plugin"
func ResolvedNestedThirdParty() { buffer := plugin.New(); _, _ = buffer.Write(nil) }
EOF
cat > "$scan_repo/internal/pkg/server.go" <<'EOF'
package pkg
import "net/http"
func Serve(server *http.Server) { _ = server.ListenAndServe() }
EOF
cat > "$scan_repo/internal/pkg/zip.go" <<'EOF'
package pkg
import "archive/zip"
func Open(path string) { _, _ = zip.OpenReader(path) }
EOF
cat > "$scan_repo/internal/pkg/dot.go" <<'EOF'
package pkg
import . "os"
func DotWrite(path string) error { return WriteFile(path, nil, 0o600) }
EOF
cat > "$scan_repo/internal/pkg/symlink.go" <<'EOF'
package pkg
import "os"
func Symlink(oldPath, newPath string) error { return os.Symlink(oldPath, newPath) }
EOF
cat > "$scan_repo/internal/pkg/inactive.go" <<'EOF'
//go:build quality_never

package pkg
import "os"
func TaggedWrite(path string) error { return os.WriteFile(path, nil, 0o600) }
EOF
mkdir -p "$scan_repo/internal/notadapter/processish"
cat > "$scan_repo/internal/notadapter/processish/run.go" <<'EOF'
package processish
import "os/exec"
func Run() { _ = exec.Command("true") }
EOF
printf '#!/usr/bin/env bash\nexit 1\n' > "$scan_repo/scripts/mutate-ready"
printf '#!/usr/bin/env bash\nexit 1\n' > "$scan_repo/scripts/mutate-not-executable"
printf '#!/usr/bin/env bash\nexit 1\n' > "$scan_repo/scripts/mutate-unrelated"
chmod +x "$scan_repo/scripts/mutate-ready"
chmod +x "$scan_repo/scripts/mutate-unrelated"
git -C "$scan_repo" init -q
git -C "$scan_repo" -c user.email=t@example.invalid -c user.name=t add .
git -C "$scan_repo" -c user.email=t@example.invalid -c user.name=t commit -q -m init
chmod 0400 "$scan_repo/scripts/mutate-not-executable"
scan_commit="$(git -C "$scan_repo" rev-parse HEAD)"
export RUN_REPO="$scan_repo" REPORT_SHORT="$(git -C "$scan_repo" rev-parse --short HEAD)"
export GOFLAGS=-mod=mod
write_report Q1.3 PASS "0 direct calls outside 1 declared adapters"
run_wrapper "$WORK/direct-effects"
rc=$?
[ "$rc" -eq 1 ] || fail "direct-effects probe returned $rc, expected 1"
json_assert "$WORK/direct-effects/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['value'] == 59 and value['criteria'][0]['metric']['population'] == 60 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/pkg/catalog_effects.go']) == 5 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/pkg/argument_origins.go']) == 11 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/pkg/write.go']) == 2 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/pkg/modern_os.go']) == 2 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/pkg/higher_order.go']) == 4 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/pkg/interface_write.go']) == 6 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/pkg/origin_variants.go' and v['kind'] == 'filesystem']) == 17 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/pkg/raw_syscall.go' and v['kind'] == 'filesystem']) == 2 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/resourceuse/write.go' and v['kind'] == 'filesystem']) == 4 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/pkg/xsys.go' and v['kind'] == 'filesystem']) == 1 and any(v['path'] == 'internal/pkg/symlink.go' for v in value['criteria'][0]['metric']['violations']) and any(v['path'] == 'internal/pkg/zip.go' for v in value['criteria'][0]['metric']['violations']) and any(v['path'] == 'internal/pkg/dot.go' for v in value['criteria'][0]['metric']['violations']) and any(v['path'] == 'internal/testutilish/write.go' for v in value['criteria'][0]['metric']['violations']) and all(v['path'] not in ('internal/testutil/write.go', 'internal/pkg/write_test.go', 'internal/pkg/unrelated.go', 'internal/pkg/third_party.go', 'internal/pkg/nested_third_party.go', 'internal/pkg/client.go', 'internal/pkg/file_write.go', 'internal/pkg/copy.go', 'internal/pkg/cross_package_write.go', 'internal/pkg/cgo_effect.go', 'internal/resource/output.go') for v in value['criteria'][0]['metric']['violations'])" \
  "Q1.3 did not resolve function/interface/syscall/cgo effects, ignore memory-only methods, or enforce exact adapters"
for iteration in 1 2 3 4 5; do
  run_wrapper "$WORK/direct-effects-repeat-$iteration"
  rc=$?
  [ "$rc" -eq 1 ] || fail "deterministic direct-effects repeat $iteration returned $rc"
  cmp -s "$WORK/direct-effects/scorecard.json" "$WORK/direct-effects-repeat-$iteration/scorecard.json" \
    || fail "multi-origin Q1.3 output varied in process $iteration"
done
stale_gopath="$WORK/stale-gopath"
mkdir -p "$stale_gopath/src/example.invalid/scanfixture/internal/resource"
cat > "$stale_gopath/src/example.invalid/scanfixture/internal/resource/output.go" <<'EOF'
package resource
type Writer interface { Write([]byte) (int, error) }
var Output Writer
func OutputFactory() Writer { return nil }
EOF
GO111MODULE=off GOPATH="$stale_gopath" go install example.invalid/scanfixture/internal/resource \
  || fail "could not build the stale local-package archive probe"
GOPATH="$stale_gopath" run_wrapper "$WORK/stale-local-archive"
rc=$?
[ "$rc" -eq 1 ] || fail "stale local-package archive probe returned $rc, expected 1"
json_assert "$WORK/stale-local-archive/scorecard.json" \
  "value['criteria'][0]['metric']['value'] == 59 and len([v for v in value['criteria'][0]['metric']['violations'] if v['path'] == 'internal/resourceuse/write.go' and v['kind'] == 'filesystem']) == 4" \
  "GOPATH export data shadowed selected local-module source origins"
default_production_count="$(python3 - "$WORK/direct-effects/scorecard.json" <<'PY'
import json
import sys
print(json.load(open(sys.argv[1], encoding="utf-8"))["tool"]["go_build"]["production_files"]["count"])
PY
)"
GOFLAGS='-mod=mod -tags=quality_never' run_wrapper "$WORK/tagged-direct-effects"
rc=$?
[ "$rc" -eq 1 ] || fail "tag-enabled direct-effects probe returned $rc, expected 1"
json_assert "$WORK/tagged-direct-effects/scorecard.json" \
  "value['criteria'][0]['metric']['value'] == 60 and value['criteria'][0]['metric']['population'] == 61 and any(v['path'] == 'internal/pkg/inactive.go' for v in value['criteria'][0]['metric']['violations']) and value['tool']['go_build']['goflags'] == '-mod=mod -tags=quality_never' and value['tool']['go_build']['production_files']['count'] == $default_production_count + 1" \
  "build-tag activation did not change the labeled selected population by one file/effect"
GOFLAGS="-mod=mod -overlay=$WORK/quality-overlay.json" run_wrapper "$WORK/overlay-build-context"
rc=$?
[ "$rc" -eq 2 ] || fail "path-sensitive GOFLAGS overlay returned $rc, expected audit-broken exit 2"
[ ! -e "$WORK/overlay-build-context/scorecard.json" ] \
  || fail "overlay build context published authoritative JSON without scanning overlaid contents"
GOFLAGS="-mod=mod --overlay=$WORK/quality-overlay.json" run_wrapper "$WORK/double-dash-overlay-build-context"
rc=$?
[ "$rc" -eq 2 ] || fail "double-dash GOFLAGS overlay returned $rc, expected audit-broken exit 2"
[ ! -e "$WORK/double-dash-overlay-build-context/scorecard.json" ] \
  || fail "double-dash overlay published authoritative JSON without scanning overlaid contents"
cat > "$scan_repo/internal/pkg/unresolved.go" <<'EOF'
package pkg
import "example.invalid/unavailable"
func UnresolvedWrite(output *unavailable.Output) { _, _ = output.Write(nil) }
func UnresolvedCommand() { _ = unavailable.Command("run") }
func UnresolvedCommandAlias() { command := unavailable.Command; _ = command("run") }
EOF
run_wrapper "$WORK/unresolved-effect"
rc=$?
[ "$rc" -eq 2 ] || fail "unresolved effect candidate returned $rc, expected audit-broken exit 2"
[ ! -e "$WORK/unresolved-effect/scorecard.json" ] \
  || fail "unresolved effect candidate left authoritative JSON"
rm "$scan_repo/internal/pkg/unresolved.go"
write_report Q2.1 PASS "3 harnesses for 3 declared subjects of 1 packages"
run_wrapper "$WORK/subject-harnesses"
rc=$?
[ "$rc" -eq 1 ] || fail "subject-harness probe returned $rc, expected 1"
json_assert "$WORK/subject-harnesses/scorecard.json" \
  "value['criteria'][0]['verdict'] == 'FAIL' and value['criteria'][0]['metric']['name'] == 'valid_subject_harness_bindings' and value['criteria'][0]['metric']['numerator'] == 1 and value['criteria'][0]['metric']['denominator'] == 9 and len(value['criteria'][0]['metric']['invalid_subject_bindings']) == 8 and len([issue for subject in value['criteria'][0]['metric']['invalid_subject_bindings'] for issue in subject['issues'] if issue['code'] == 'excluded_root']) == 3 and any(issue['code'] == 'missing' and issue['kind'] == 'source' for subject in value['criteria'][0]['metric']['invalid_subject_bindings'] for issue in subject['issues']) and any(issue['code'] == 'no_production_go' for subject in value['criteria'][0]['metric']['invalid_subject_bindings'] for issue in subject['issues']) and any(issue['code'] == 'invalid_path' for subject in value['criteria'][0]['metric']['invalid_subject_bindings'] for issue in subject['issues'])" \
  "Q2.1 did not bind each subject to real repository production roots and its executable harness"
unset RUN_REPO REPORT_SHORT GOFLAGS
ok "receiver types distinguish real effects; test support and adapters are exact"

printf 'T11 structured output is deterministic\n'
write_report Q1.4 PASS "0 of 0 declared seams covered by a named mutation"
run_wrapper "$WORK/repeat"
rc=$?
[ "$rc" -eq 1 ] || fail "repeat audit returned $rc, expected 1"
cmp -s "$WORK/empty/scorecard.json" "$WORK/repeat/scorecard.json" \
  || fail "identical input did not produce byte-identical JSON"
ok "two conversions of the same commit/report are byte-identical"

printf 'T12 invalid inherited locale still yields a real hash and cleans audit work\n'
mkdir -p "$WORK/audit-tmp"
LC_ALL=C.UTF-8 LANG=C.UTF-8 TMPDIR="$WORK/audit-tmp/" \
  bash "$WRAPPER" "$WORK/repo" --quick --only Q0.2,Q0.5,Q1.3 \
  --out "$WORK/real-audit" >/dev/null 2>&1
rc=$?
[ "$rc" -eq 0 ] || fail "focused real audit returned $rc, expected selected-criteria exit 0"
json_assert "$WORK/real-audit/scorecard.json" \
  "[item for item in value['criteria'] if item['id'] == 'Q0.5'][0]['verdict'] == 'PASS'" \
  "locale-safe real audit did not record a Q0.5 PASS"
grep -qE 'checksum unchanged \([0-9a-f]{12}\)' "$WORK/real-audit/scorecard.json" \
  || fail "real audit did not persist a non-empty checksum"
leftovers="$(find "$WORK/audit-tmp" -mindepth 1 -maxdepth 1 ! -name 'gocache-audit' -print)"
[ -z "$leftovers" ] || fail "audit work directory was not cleaned: $leftovers"
ok "invalid locale was replaced and go-list/scanner work and private build caches were removed"

printf 'T13 upstream audit-broken exits leave no transient authoritative-looking report\n'
write_report Q1.2 FAIL "127 process-exiting calls outside main() (RATCHET)"
mkdir -p "$WORK/no-python-bin" "$WORK/missing-python-preflight"
ln -s "$(command -v dirname)" "$WORK/no-python-bin/dirname"
ln -s "$(command -v unlink)" "$WORK/no-python-bin/unlink"
ln -s "$(command -v rmdir)" "$WORK/no-python-bin/rmdir"
printf 'STALE\n' > "$WORK/missing-python-preflight/scorecard.json"
PATH="$WORK/no-python-bin" /bin/bash "$WRAPPER" "$WORK/repo" \
  --out "$WORK/missing-python-preflight" --only Q1.2 >/dev/null 2>&1
rc=$?
[ "$rc" -eq 2 ] || fail "missing Python preflight returned $rc, expected 2"
assert_no_generated_outputs "$WORK/missing-python-preflight" "missing Python preflight"
mkdir -p "$WORK/missing-upstream-preflight"
printf 'STALE\n' > "$WORK/missing-upstream-preflight/scorecard.json"
QUALITY_AUDIT_UPSTREAM="$WORK/does-not-exist-upstream" \
  bash "$WRAPPER" "$WORK/repo" --out "$WORK/missing-upstream-preflight" \
  --only Q1.2 >/dev/null 2>&1
rc=$?
[ "$rc" -eq 2 ] || fail "missing upstream preflight returned $rc, expected 2"
assert_no_generated_outputs "$WORK/missing-upstream-preflight" "missing upstream preflight"
mkdir -p "$WORK/missing-parser-preflight"
printf 'STALE\n' > "$WORK/missing-parser-preflight/scorecard.json"
QUALITY_AUDIT_PARSER="$WORK/does-not-exist-parser" \
  bash "$WRAPPER" "$WORK/repo" --out "$WORK/missing-parser-preflight" \
  --only Q1.2 >/dev/null 2>&1
rc=$?
[ "$rc" -eq 2 ] || fail "missing parser preflight returned $rc, expected 2"
assert_no_generated_outputs "$WORK/missing-parser-preflight" "missing parser preflight"
export FAKE_EXIT=9 FAKE_SEED_OUTPUTS=yes
run_wrapper "$WORK/invalid-upstream-exit"
rc=$?
unset FAKE_EXIT FAKE_SEED_OUTPUTS
[ "$rc" -eq 2 ] || fail "undocumented upstream exit returned $rc, expected audit-broken exit 2"
assert_no_generated_outputs "$WORK/invalid-upstream-exit" "undocumented upstream exit"
export FAKE_EXIT=2 FAKE_SEED_OUTPUTS=yes
run_wrapper "$WORK/upstream-exit-two"
rc=$?
unset FAKE_EXIT FAKE_SEED_OUTPUTS
[ "$rc" -eq 2 ] || fail "upstream exit 2 returned $rc, expected audit-broken exit 2"
assert_no_generated_outputs "$WORK/upstream-exit-two" "upstream exit 2"
export FAKE_EMPTY_REPORT=yes FAKE_SEED_OUTPUTS=yes
run_wrapper "$WORK/empty-upstream-report"
rc=$?
unset FAKE_EMPTY_REPORT FAKE_SEED_OUTPUTS
[ "$rc" -eq 2 ] || fail "empty upstream report returned $rc, expected audit-broken exit 2"
assert_no_generated_outputs "$WORK/empty-upstream-report" "empty upstream report"
export FAKE_DIRECTORY_REPORT=yes
run_wrapper "$WORK/directory-upstream-report"
rc=$?
unset FAKE_DIRECTORY_REPORT
[ "$rc" -eq 2 ] || fail "directory upstream report returned $rc, expected audit-broken exit 2"
assert_no_generated_outputs "$WORK/directory-upstream-report" "directory upstream report"
ok "dependency/upstream failures and invalid reports remove every generated artifact"

printf 'T14 a raw-report preservation failure invalidates every staged artifact\n'
write_report Q0.6 PASS "leak guard used at 2 call sites; 0 skipped tests"
export FAKE_BLOCK_RAW=yes
run_wrapper "$WORK/raw-preservation-failure"
rc=$?
unset FAKE_BLOCK_RAW
[ "$rc" -eq 2 ] || fail "raw-report preservation failure returned $rc, expected 2"
assert_no_generated_outputs "$WORK/raw-preservation-failure" "raw-report preservation failure"
export FAKE_BLOCK_JSON=yes
run_wrapper "$WORK/json-publication-failure"
rc=$?
unset FAKE_BLOCK_JSON
[ "$rc" -eq 2 ] || fail "JSON publication directory collision returned $rc, expected 2"
assert_no_generated_outputs "$WORK/json-publication-failure" "JSON publication directory collision"
export FAKE_HARDLINK_JSON=yes
run_wrapper "$WORK/json-hardlink-publication-failure"
rc=$?
unset FAKE_HARDLINK_JSON
[ "$rc" -eq 2 ] || fail "JSON publication hard-link collision returned $rc, expected 2"
assert_no_generated_outputs "$WORK/json-hardlink-publication-failure" "JSON publication hard-link collision"
export FAKE_HARDLINK_RAW=yes
run_wrapper "$WORK/raw-hardlink-publication-failure"
rc=$?
unset FAKE_HARDLINK_RAW
[ "$rc" -eq 2 ] || fail "raw publication hard-link collision returned $rc, expected 2"
assert_no_generated_outputs "$WORK/raw-hardlink-publication-failure" "raw publication hard-link collision"
ok "post-parser publication failures remove JSON, reports, pending files, blockers, and hard links"

printf 'T15 old and new baseline instruments reproduce with identical numeric debt\n'
old_effect_repo="$WORK/baseline-old-effect-repository"
old_reproduction_repo="$WORK/baseline-old-reproduction-repository"
old_structured_repo="$WORK/baseline-old-structured-repository"
new_effect_repo="$WORK/baseline-new-effect-repository"
new_reproduction_repo="$WORK/baseline-new-reproduction-repository"
new_structured_repo="$WORK/baseline-new-structured-repository"
baseline_drift_probe="$WORK/baseline-drift-probe"
old_instrument="$WORK/baseline-old-instrument"
old_effect_out="$WORK/baseline-old-effect-output"
old_out="$WORK/baseline-old-output"
new_effect_out="$WORK/baseline-new-effect-output"
new_out="$WORK/baseline-new-output"
old_effect_home="$WORK/baseline-old-effect-home"
new_effect_home="$WORK/baseline-new-effect-home"
old_reproduction_home="$WORK/baseline-old-reproduction-home"
new_reproduction_home="$WORK/baseline-new-reproduction-home"
reproduction_home_fixture="$HERE/../baseline/reproduction-home"
migration="$HERE/../baseline/instrument-migration.json"
repository_root="$(cd "$HERE/../.." && pwd)"
old_source="$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["old"]["instrument_source_commit"])' "$migration")"
new_source="$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["new"]["instrument_source_commit"])' "$migration")"
baseline_commit=5635d50bd161a9a5aa81fc4332cc0c9d68885d08
baseline_overlay_sha256=4cec690f46b9595d70bce9a08c164aa56d46ac5bb3d69bfe84f83d9006c06e8d
baseline_gomodcache="$(go env GOMODCACHE)"
[ -n "$baseline_gomodcache" ] && [ "${baseline_gomodcache#/}" != "$baseline_gomodcache" ] \
  || fail "could not resolve the shared Go module cache before isolating baseline HOME"

prepare_baseline_reproduction_repo() {
  destination="$1"
  git clone -q --shared --no-checkout "$repository_root" "$destination" || return 1
  git -C "$destination" checkout -q --detach "$baseline_commit" || return 1
  mkdir -p "$destination/.quality" || return 1
  cp "$HERE/../baseline/inventory" "$destination/.quality/inventory"
}

prepare_baseline_reproduction_home() {
  destination="$1"
  mkdir -p "$destination" || return 1
  cp -R "$reproduction_home_fixture/." "$destination" || return 1
  python3 - "$migration" "$destination" <<'PY'
import hashlib
import json
from pathlib import Path
import sys

metadata = json.load(open(sys.argv[1], encoding="utf-8"))["measured_source"]["reproduction_home"]
root = Path(sys.argv[2]).resolve()
relative = metadata["active_profile_path"]
if not relative or Path(relative).is_absolute() or Path(relative).as_posix() != relative or ".." in Path(relative).parts:
    raise SystemExit(1)
active_profile = root / relative
if metadata["active_profile_type"] != "directory" or not active_profile.is_dir() or active_profile.is_symlink():
    raise SystemExit(1)
marker = active_profile / ".fixture"
if not marker.is_file() or marker.is_symlink():
    raise SystemExit(1)
if hashlib.sha256(marker.read_bytes()).hexdigest() != metadata["marker_sha256"]:
    raise SystemExit(1)
PY
}

assert_baseline_reproduction_input() {
  identity_parser="$1"
  identity_repo="$2"
  python3 - "$identity_parser" "$identity_repo" "$baseline_commit" "$baseline_overlay_sha256" <<'PY'
from pathlib import Path
import sys

source, repository, commit, overlay = sys.argv[1:]
namespace = {"__name__": "quality_baseline_input_contract", "__file__": source}
exec(compile(Path(source).read_text(encoding="utf-8"), source, "exec"), namespace)
inventory, inventory_identity = namespace["inventory_snapshot"](repository, overlay)
tree = namespace["tree_identity"](repository, commit, overlay)
namespace["verify_inventory_identity"](repository, inventory_identity)
if not tree["measurement_clean"] or tree["dirty_paths"] != [".quality/inventory"]:
    raise SystemExit(1)
PY
}

write_baseline_reproduction_effect_manifest() {
  identity_parser="$1"
  identity_repo="$2"
  manifest_output="$3"
  python3 - "$identity_parser" "$identity_repo" "$manifest_output" <<'PY'
import hashlib
import json
import os
from pathlib import Path
import stat
import sys

source, repository_value, output = sys.argv[1:]
namespace = {"__name__": "quality_baseline_effect_contract", "__file__": source}
exec(compile(Path(source).read_text(encoding="utf-8"), source, "exec"), namespace)
repository = Path(repository_value).resolve()
effects = []
for relative in namespace["measurement_dirty_paths"](repository):
    path = repository / relative
    if path.is_symlink():
        target = os.fsencode(os.readlink(path))
        effect = {
            "path": relative,
            "type": "symlink",
            "target_sha256": hashlib.sha256(target).hexdigest(),
            "relocatable_target_sha256": hashlib.sha256(
                target.replace(os.fsencode(str(repository)), b"<BASELINE_REPOSITORY>")
            ).hexdigest(),
        }
    elif path.is_file():
        contents = path.read_bytes()
        effect = {
            "path": relative,
            "type": "file",
            "mode": stat.S_IMODE(path.stat().st_mode),
            "sha256": hashlib.sha256(contents).hexdigest(),
            "relocatable_sha256": hashlib.sha256(
                contents.replace(os.fsencode(str(repository)), b"<BASELINE_REPOSITORY>")
            ).hexdigest(),
        }
    elif path.exists():
        effect = {"path": relative, "type": "other", "mode": stat.S_IMODE(path.stat().st_mode)}
    else:
        effect = {"path": relative, "type": "missing"}
    effects.append(effect)
Path(output).write_text(
    json.dumps(effects, ensure_ascii=False, separators=(",", ":"), sort_keys=True) + "\n",
    encoding="utf-8",
)
PY
}

compare_baseline_reproduction_effect_manifests() {
  old_manifest="$1"
  new_manifest="$2"
  python3 - "$old_manifest" "$new_manifest" <<'PY'
import json
import sys

old = json.load(open(sys.argv[1], encoding="utf-8"))
new = json.load(open(sys.argv[2], encoding="utf-8"))
if len(old) != len(new):
    raise SystemExit(1)
for old_effect, new_effect in zip(old, new):
    for key in ("path", "type", "mode"):
        if old_effect.get(key) != new_effect.get(key):
            raise SystemExit(1)
    if old_effect["type"] == "file":
        if old_effect["relocatable_sha256"] != new_effect["relocatable_sha256"]:
            raise SystemExit(1)
    elif old_effect["type"] == "symlink":
        if old_effect["relocatable_target_sha256"] != new_effect["relocatable_target_sha256"]:
            raise SystemExit(1)
PY
}

prepare_baseline_reproduction_repo "$old_effect_repo" \
  || fail "could not prepare exact old baseline effect input"
prepare_baseline_reproduction_repo "$old_reproduction_repo" \
  || fail "could not prepare exact old baseline reproduction input"
prepare_baseline_reproduction_repo "$old_structured_repo" \
  || fail "could not prepare exact old baseline structured input"
prepare_baseline_reproduction_repo "$new_effect_repo" \
  || fail "could not prepare exact new baseline effect input"
prepare_baseline_reproduction_repo "$new_reproduction_repo" \
  || fail "could not prepare exact new baseline reproduction input"
prepare_baseline_reproduction_repo "$new_structured_repo" \
  || fail "could not prepare exact new baseline structured input"
prepare_baseline_reproduction_repo "$baseline_drift_probe" \
  || fail "could not prepare baseline source-drift probe"
prepare_baseline_reproduction_home "$old_reproduction_home" \
  || fail "could not prepare pinned old baseline HOME fixture"
prepare_baseline_reproduction_home "$new_reproduction_home" \
  || fail "could not prepare pinned new baseline HOME fixture"
mkdir -p "$old_instrument" "$old_effect_out" "$old_out" "$new_effect_out" "$new_out" \
  "$old_effect_home" "$new_effect_home" \
  "$WORK/baseline-old-effect-gocache" "$WORK/baseline-old-gocache" \
  "$WORK/baseline-new-effect-gocache" "$WORK/baseline-new-gocache"
git -C "$repository_root" archive "$old_source" -- \
  .quality/tools .quality/baseline/manual-evidence.json .quality/baseline/scorecard.json \
  | tar -x -C "$old_instrument"

assert_baseline_reproduction_input \
  "$old_instrument/.quality/tools/scorecard.py" "$old_effect_repo" \
  || fail "exact old baseline effect input was not accepted before audit execution"
assert_baseline_reproduction_input \
  "$old_instrument/.quality/tools/scorecard.py" "$old_reproduction_repo" \
  || fail "exact old baseline reproduction input was not accepted before audit execution"
assert_baseline_reproduction_input \
  "$old_instrument/.quality/tools/scorecard.py" "$old_structured_repo" \
  || fail "exact old baseline structured input was not accepted before measurement"
assert_baseline_reproduction_input "$PARSER" "$new_effect_repo" \
  || fail "exact new baseline effect input was not accepted before audit execution"
assert_baseline_reproduction_input "$PARSER" "$new_reproduction_repo" \
  || fail "exact new baseline reproduction input was not accepted before audit execution"
assert_baseline_reproduction_input "$PARSER" "$new_structured_repo" \
  || fail "exact new baseline structured input was not accepted before measurement"
printf '\nT15 pre-existing baseline source drift\n' >> "$baseline_drift_probe/README.md"
if assert_baseline_reproduction_input \
  "$old_instrument/.quality/tools/scorecard.py" "$baseline_drift_probe" >/dev/null 2>&1; then
  fail "pre-existing baseline source drift was accepted as test-generated output"
fi
write_baseline_reproduction_effect_manifest \
  "$old_instrument/.quality/tools/scorecard.py" "$old_effect_repo" "$old_effect_out/effects-before.json" \
  || fail "could not record old baseline effects before audit execution"
write_baseline_reproduction_effect_manifest \
  "$PARSER" "$new_effect_repo" "$new_effect_out/effects-before.json" \
  || fail "could not record new baseline effects before audit execution"
compare_baseline_reproduction_effect_manifests \
  "$old_effect_out/effects-before.json" "$new_effect_out/effects-before.json" \
  || fail "old and new baseline execution inputs differ before audit execution"

LC_ALL=C LANG=C CGO_ENABLED=0 GOENV=off GOWORK=off HOME="$old_effect_home" \
  GOMODCACHE="$baseline_gomodcache" \
  GOCACHE="$WORK/baseline-old-effect-gocache" \
  bash "$old_instrument/.quality/tools/vendor/quality-audit.sh" \
  "$old_effect_repo" --out "$old_effect_out" >/dev/null
old_effect_rc=$?
[ "$old_effect_rc" -eq 1 ] \
  || fail "old baseline effect probe returned $old_effect_rc, expected 1"
write_baseline_reproduction_effect_manifest \
  "$old_instrument/.quality/tools/scorecard.py" "$old_effect_repo" "$old_effect_out/effects-after.json" \
  || fail "could not record old test-generated baseline effects"
cmp -s "$old_effect_out/effects-before.json" "$old_effect_out/effects-after.json" \
  && fail "historical baseline audit produced no test-generated output for the focused contract"

LC_ALL=C LANG=C CGO_ENABLED=0 GOENV=off GOWORK=off HOME="$new_effect_home" \
  GOMODCACHE="$baseline_gomodcache" \
  GOCACHE="$WORK/baseline-new-effect-gocache" \
  bash "$VENDOR" "$new_effect_repo" --out "$new_effect_out" >/dev/null
new_effect_rc=$?
[ "$new_effect_rc" -eq 1 ] \
  || fail "new baseline effect probe returned $new_effect_rc, expected 1"
write_baseline_reproduction_effect_manifest \
  "$PARSER" "$new_effect_repo" "$new_effect_out/effects-after.json" \
  || fail "could not record new test-generated baseline effects"
cmp -s "$new_effect_out/effects-before.json" "$new_effect_out/effects-after.json" \
  && fail "new historical baseline audit produced no test-generated output for the focused contract"
compare_baseline_reproduction_effect_manifests \
  "$old_effect_out/effects-after.json" "$new_effect_out/effects-after.json" \
  || fail "old and new baseline effect probes changed different paths or relocation-independent bytes"
tail -n +2 "$old_effect_out/scorecard.md" > "$old_effect_out/reproduced-raw.body"
tail -n +2 "$new_effect_out/scorecard.md" > "$new_effect_out/reproduced-raw.body"
cmp -s "$old_effect_out/reproduced-raw.body" "$new_effect_out/reproduced-raw.body" \
  || fail "old and new baseline effect-probe reports differ"

write_baseline_reproduction_effect_manifest \
  "$old_instrument/.quality/tools/scorecard.py" "$old_reproduction_repo" "$old_out/effects-before.json" \
  || fail "could not record old reproduction effects before audit execution"
write_baseline_reproduction_effect_manifest \
  "$PARSER" "$new_reproduction_repo" "$new_out/effects-before.json" \
  || fail "could not record new reproduction effects before audit execution"
compare_baseline_reproduction_effect_manifests \
  "$old_out/effects-before.json" "$new_out/effects-before.json" \
  || fail "old and new baseline reproduction inputs differ before audit execution"

LC_ALL=C LANG=C CGO_ENABLED=0 GOENV=off GOWORK=off HOME="$old_reproduction_home" \
  GOMODCACHE="$baseline_gomodcache" \
  GOCACHE="$WORK/baseline-old-gocache" \
  bash "$old_instrument/.quality/tools/vendor/quality-audit.sh" \
  "$old_reproduction_repo" --out "$old_out" >/dev/null
old_upstream_rc=$?
[ "$old_upstream_rc" -eq 1 ] \
  || fail "old upstream baseline returned $old_upstream_rc, expected 1"
write_baseline_reproduction_effect_manifest \
  "$old_instrument/.quality/tools/scorecard.py" "$old_reproduction_repo" "$old_out/effects-after.json" \
  || fail "could not record old reproduction effects after audit execution"
LC_ALL=C LANG=C CGO_ENABLED=0 GOENV=off GOWORK=off HOME="$old_reproduction_home" \
  GOMODCACHE="$baseline_gomodcache" \
  GOCACHE="$WORK/baseline-old-gocache" \
  python3 "$old_instrument/.quality/tools/scorecard.py" \
  --report "$old_out/scorecard.md" \
  --output "$old_out/scorecard.json" \
  --repo "$old_structured_repo" \
  --commit "$baseline_commit" \
  --upstream "$old_instrument/.quality/tools/vendor/quality-audit.sh" \
  --upstream-exit "$old_upstream_rc" \
  --manual-evidence "$old_instrument/.quality/baseline/manual-evidence.json" \
  --inventory-overlay-sha256 "$baseline_overlay_sha256" \
  >/dev/null
old_structured_rc=$?
[ "$old_structured_rc" -eq 1 ] \
  || fail "old structured baseline returned $old_structured_rc, expected 1"
assert_baseline_reproduction_input \
  "$old_instrument/.quality/tools/scorecard.py" "$old_structured_repo" \
  || fail "old structured measurement changed its exact baseline input"

LC_ALL=C LANG=C CGO_ENABLED=0 GOENV=off GOWORK=off HOME="$new_reproduction_home" \
  GOMODCACHE="$baseline_gomodcache" \
  GOCACHE="$WORK/baseline-new-gocache" \
  bash "$VENDOR" "$new_reproduction_repo" --out "$new_out" >/dev/null
new_upstream_rc=$?
[ "$new_upstream_rc" -eq 1 ] \
  || fail "new upstream baseline returned $new_upstream_rc, expected 1"
write_baseline_reproduction_effect_manifest \
  "$PARSER" "$new_reproduction_repo" "$new_out/effects-after.json" \
  || fail "could not record new reproduction effects after audit execution"
compare_baseline_reproduction_effect_manifests \
  "$old_out/effects-after.json" "$new_out/effects-after.json" \
  || fail "old and new baseline reproductions changed different paths or relocation-independent bytes"
LC_ALL=C LANG=C CGO_ENABLED=0 GOENV=off GOWORK=off HOME="$new_reproduction_home" \
  GOMODCACHE="$baseline_gomodcache" \
  GOCACHE="$WORK/baseline-new-gocache" python3 "$PARSER" \
  --report "$new_out/scorecard.md" \
  --output "$new_out/scorecard.json" \
  --repo "$new_structured_repo" \
  --commit "$baseline_commit" \
  --upstream "$VENDOR" \
  --upstream-exit "$new_upstream_rc" \
  --manual-evidence "$HERE/../baseline/manual-evidence.json" \
  --inventory-overlay-sha256 "$baseline_overlay_sha256" \
  >/dev/null
new_structured_rc=$?
[ "$new_structured_rc" -eq 1 ] \
  || fail "new structured baseline returned $new_structured_rc, expected 1"
assert_baseline_reproduction_input "$PARSER" "$new_structured_repo" \
  || fail "new structured measurement changed its exact baseline input"

tail -n +2 "$old_out/scorecard.md" > "$old_out/reproduced-raw.body"
tail -n +2 "$new_out/scorecard.md" > "$new_out/reproduced-raw.body"
tail -n +2 "$HERE/../baseline/raw-upstream-scorecard.md" > "$new_out/stored-raw.body"
cmp -s "$old_out/reproduced-raw.body" "$new_out/reproduced-raw.body" \
  || fail "old and new raw baseline bodies differ"
cmp -s "$new_out/reproduced-raw.body" "$new_out/stored-raw.body" \
  || fail "baseline raw report body did not reproduce"
cmp -s "$old_out/scorecard.json" "$old_instrument/.quality/baseline/scorecard.json" \
  || fail "old authoritative baseline JSON did not reproduce byte-for-byte"
cmp -s "$new_out/scorecard.json" "$HERE/../baseline/scorecard.json" \
  || fail "new authoritative baseline JSON did not reproduce byte-for-byte"

python3 - "$migration" "$old_out/scorecard.json" "$new_out/scorecard.json" \
  "$old_instrument/.quality/baseline/manual-evidence.json" \
  "$HERE/../baseline/manual-evidence.json" "$new_out/reproduced-raw.body" \
  "$repository_root" "$old_source" "$new_source" <<'PY' \
  || fail "baseline instrument migration record or debt comparison is invalid"
import hashlib
import json
import subprocess
import sys

manifest_path, old_path, new_path, old_manual, new_manual, raw_body, repository, old_source, new_source = sys.argv[1:]
manifest = json.load(open(manifest_path, encoding="utf-8"))
old = json.load(open(old_path, encoding="utf-8"))
new = json.load(open(new_path, encoding="utf-8"))

def file_digest(path):
    return hashlib.sha256(open(path, "rb").read()).hexdigest()

def numeric_leaves(value, path=""):
    if isinstance(value, bool):
        return {}
    if isinstance(value, (int, float)):
        return {path: value}
    if isinstance(value, dict):
        result = {}
        for key, child in sorted(value.items()):
            result.update(numeric_leaves(child, path + "/" + str(key)))
        return result
    if isinstance(value, list):
        result = {}
        for index, child in enumerate(value):
            result.update(numeric_leaves(child, path + "/" + str(index)))
        return result
    return {}

scope_old = {"denominators": old["denominators"], "criteria": old["criteria"]}
scope_new = {"denominators": new["denominators"], "criteria": new["criteria"]}
old_numeric = numeric_leaves(scope_old)
new_numeric = numeric_leaves(scope_new)
comparison = manifest["comparison"]
measured_source = manifest["measured_source"]
assert old["criteria"] == new["criteria"] and comparison["criteria_equal"] is True
assert old["denominators"] == new["denominators"] and comparison["denominators_equal"] is True
assert old_numeric == new_numeric
assert len(old_numeric) == comparison["numeric_debt_leaves"]
assert next(value for value in old["criteria"] if value["id"] == "Q3.9")["verdict"] == comparison["q3_9_old_verdict"] == "PASS"
assert next(value for value in new["criteria"] if value["id"] == "Q3.9")["verdict"] == comparison["q3_9_new_verdict"] == "PASS"
assert file_digest(old_path) == manifest["old"]["scorecard_sha256"]
assert file_digest(new_path) == manifest["new"]["scorecard_sha256"]
assert file_digest(old_manual) == manifest["old"]["manual_evidence_sha256"]
assert file_digest(new_manual) == manifest["new"]["manual_evidence_sha256"]
assert file_digest(raw_body) == comparison["raw_report_body_sha256"]
for label, scorecard, source in (("old", old, old_source), ("new", new, new_source)):
    assert source == manifest[label]["instrument_source_commit"]
    assert scorecard["repository"]["commit"] == measured_source["commit"]
    assert scorecard["repository"]["tree"]["commit_tree"] == measured_source["commit_tree"]
    assert scorecard["repository"]["tree"]["inventory_overlay_sha256"] == measured_source["inventory_overlay_sha256"]
    assert all(scorecard["tool"].get(key) == value for key, value in manifest[label]["tool"].items())
    parser_bytes = subprocess.check_output([
        "git", "-C", repository, "show", source + ":.quality/tools/scorecard.py",
    ])
    assert hashlib.sha256(parser_bytes).hexdigest() == manifest[label]["tool"]["parser_sha256"]
PY
ok "old/new instruments, raw body, Q3.9, and all 228 numeric debt leaves reproduced"

printf 'OK: 15 quality-audit controls passed\n'
