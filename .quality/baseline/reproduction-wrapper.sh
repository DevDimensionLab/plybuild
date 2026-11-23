#!/usr/bin/env bash
# TEMPLATE-VERSION: 1
# Project wrapper around the pinned methodology audit. The upstream tool performs
# measurements; scorecard.py supplies structured output and the ratchet it documents
# but does not implement.
set -uo pipefail
export LC_ALL=C LANG=C CGO_ENABLED=0 GOENV=off GOWORK=off

HERE="$(cd "$(dirname "$0")" && pwd)"
UPSTREAM="${QUALITY_AUDIT_UPSTREAM:-$HERE/vendor/quality-audit.sh}"
PARSER="${QUALITY_AUDIT_PARSER:-$HERE/scorecard.py}"
REPO="."
OUT=""
BASELINE=""
MANUAL=""
MANUAL_EXPLICIT="no"
PASSTHROUGH="no"
PARTIAL="no"
FORWARD=()

while [ $# -gt 0 ]; do
  case "$1" in
    --baseline)
      [ $# -ge 2 ] || { printf '%s\n' 'AUDIT BROKEN: --baseline needs a path' >&2; exit 2; }
      BASELINE="$2"; shift 2 ;;
    --manual-evidence)
      [ $# -ge 2 ] || { printf '%s\n' 'AUDIT BROKEN: --manual-evidence needs a path' >&2; exit 2; }
      MANUAL="$2"; MANUAL_EXPLICIT="yes"; shift 2 ;;
    --out)
      [ $# -ge 2 ] || { printf '%s\n' 'AUDIT BROKEN: --out needs a path' >&2; exit 2; }
      OUT="$2"; FORWARD+=("$1" "$2"); shift 2 ;;
    --only)
      [ $# -ge 2 ] || { printf '%s\n' 'AUDIT BROKEN: --only needs a pattern' >&2; exit 2; }
      PARTIAL="yes"; FORWARD+=("$1" "$2"); shift 2 ;;
    -h|--help)
      PASSTHROUGH="yes"; FORWARD+=("$1"); shift ;;
    --self-test|--corpus)
      PASSTHROUGH="yes"; FORWARD+=("$1")
      if [ "$1" = "--corpus" ]; then
        [ $# -ge 2 ] || { printf '%s\n' 'AUDIT BROKEN: --corpus needs a path' >&2; exit 2; }
        FORWARD+=("$2"); shift 2
      else
        shift
      fi ;;
    -*) FORWARD+=("$1"); shift ;;
    *)  REPO="$1"; FORWARD+=("$1"); shift ;;
  esac
done

if [ "$PASSTHROUGH" = "yes" ]; then
  [ -f "$UPSTREAM" ] \
    || { printf 'AUDIT BROKEN: missing upstream audit %s\n' "$UPSTREAM" >&2; exit 2; }
  exec bash "$UPSTREAM" "${FORWARD[@]}"
fi

REPO_ABS="$(cd "$REPO" 2>/dev/null && pwd)" \
  || { printf 'AUDIT BROKEN: no such repository %s\n' "$REPO" >&2; exit 2; }
[ -n "$OUT" ] || OUT="$REPO_ABS/target/quality-audit"
[ -n "$MANUAL" ] || MANUAL="$REPO_ABS/.quality/manual-evidence.json"
REPORT="$OUT/scorecard.md"
JSON="$OUT/scorecard.json"
RAW_REPORT="$OUT/raw-upstream-scorecard.md"
PENDING_JSON="$OUT/.scorecard.json.pending"
PENDING_RAW="$OUT/.raw-upstream-scorecard.md.pending"

invalidate_artifact() {
  local artifact="$1"
  if [ -d "$artifact" ] && [ ! -L "$artifact" ]; then
    rmdir "$artifact" 2>/dev/null \
      || { printf 'AUDIT BROKEN: cannot invalidate generated artifact directory %s\n' "$artifact" >&2; return 1; }
  elif [ -e "$artifact" ] || [ -L "$artifact" ]; then
    unlink "$artifact" 2>/dev/null \
      || { printf 'AUDIT BROKEN: cannot invalidate generated artifact %s\n' "$artifact" >&2; return 1; }
  fi
}

invalidate_generated_outputs() {
  local status=0 artifact
  for artifact in "$JSON" "$RAW_REPORT" "$REPORT" "$PENDING_JSON" "$PENDING_RAW"; do
    invalidate_artifact "$artifact" || status=1
  done
  return "$status"
}

atomic_replace_regular() {
  python3 - "$1" "$2" <<'PY'
import os
import stat
import sys

source, destination = sys.argv[1:]
try:
    before = os.lstat(source)
    if not stat.S_ISREG(before.st_mode) or before.st_size == 0:
        raise OSError("source is not a non-empty regular file")
    try:
        current = os.lstat(destination)
    except FileNotFoundError:
        current = None
    if current is not None and not stat.S_ISREG(current.st_mode):
        raise OSError("destination exists and is not a regular file")
    os.replace(source, destination)
    after = os.lstat(destination)
    if (after.st_dev, after.st_ino) != (before.st_dev, before.st_ino):
        raise OSError("source changed during atomic publication")
    if not stat.S_ISREG(after.st_mode) or after.st_size == 0:
        raise OSError("published artifact is not a non-empty regular file")
    try:
        os.lstat(source)
    except FileNotFoundError:
        pass
    else:
        raise OSError("source name survived atomic publication")
except OSError as error:
    print("{} -> {}: {}".format(source, destination, error), file=sys.stderr)
    raise SystemExit(1)
PY
}

reject_generated_alias() {
  local label="$1" input="$2" generated alias
  [ -n "$input" ] || return 0
  for generated in "$JSON" "$RAW_REPORT" "$REPORT" "$PENDING_JSON" "$PENDING_RAW"; do
    alias="$(python3 - "$input" "$generated" <<'PY'
import os
import sys

source, generated = (os.path.abspath(value) for value in sys.argv[1:])
same_path = os.path.realpath(source) == os.path.realpath(generated)
same_file = False
if os.path.exists(source) and os.path.exists(generated):
    same_file = os.path.samefile(source, generated)
print("yes" if same_path or same_file else "no")
PY
)" || { printf 'AUDIT BROKEN: cannot compare %s with generated path %s\n' "$label" "$generated" >&2; return 1; }
    if [ "$alias" = "yes" ]; then
      printf 'AUDIT BROKEN: %s aliases generated output %s\n' "$label" "$generated" >&2
      return 1
    fi
  done
}

if ! command -v python3 >/dev/null 2>&1; then
  # Without Python canonical alias checks are unavailable. It is still safe to
  # invalidate generated names when no baseline or manual-evidence input exists.
  if [ -z "$BASELINE" ] && [ "$MANUAL_EXPLICIT" = "no" ] && \
      [ ! -e "$MANUAL" ] && [ ! -L "$MANUAL" ]; then
    invalidate_generated_outputs || true
  fi
  printf '%s\n' 'AUDIT BROKEN: python3 is required for structured output' >&2
  exit 2
fi
reject_generated_alias --baseline "$BASELINE" || exit 2
reject_generated_alias --manual-evidence "$MANUAL" || exit 2

invalidate_generated_outputs || exit 2
[ -f "$UPSTREAM" ] || { printf 'AUDIT BROKEN: missing upstream audit %s\n' "$UPSTREAM" >&2; exit 2; }
[ -f "$PARSER" ] || { printf 'AUDIT BROKEN: missing scorecard parser %s\n' "$PARSER" >&2; exit 2; }
[ "$MANUAL_EXPLICIT" = "no" ] || [ -f "$MANUAL" ] \
  || { printf 'AUDIT BROKEN: explicit manual evidence does not exist: %s\n' "$MANUAL" >&2; exit 2; }

AUDIT_ENV_ROOT="$(mktemp -d /tmp/ply-quality-audit-env.XXXXXX 2>/dev/null)" \
  || { printf '%s\n' 'AUDIT BROKEN: cannot create isolated audit environment' >&2; exit 2; }
cleanup_audit_environment() {
  [ -n "${AUDIT_ENV_ROOT:-}" ] || return 0
  python3 - "$AUDIT_ENV_ROOT" <<'PY' >/dev/null 2>&1 || true
import os
import shutil
import stat
import sys

root = sys.argv[1]
for directory, directories, files in os.walk(root, topdown=False):
    for name in files + directories:
        path = os.path.join(directory, name)
        try:
            os.chmod(path, os.lstat(path).st_mode | stat.S_IRUSR | stat.S_IWUSR | stat.S_IXUSR)
        except OSError:
            pass
shutil.rmtree(root, ignore_errors=True)
PY
}
trap cleanup_audit_environment EXIT HUP INT TERM
mkdir "$AUDIT_ENV_ROOT/tmp" "$AUDIT_ENV_ROOT/go-tmp" \
  || { printf '%s\n' 'AUDIT BROKEN: cannot initialize isolated audit environment' >&2; exit 2; }
python3 - "$REPO_ABS" "$AUDIT_ENV_ROOT" <<'PY' || exit 2
import os
import sys

repo, temporary = (os.path.realpath(value) for value in sys.argv[1:])
try:
    inside = os.path.commonpath((repo, temporary)) == repo
except ValueError:
    inside = False
if inside:
    print("AUDIT BROKEN: isolated audit environment is inside the measured repository", file=sys.stderr)
    raise SystemExit(1)
PY
export TMPDIR="$AUDIT_ENV_ROOT/tmp" GOTMPDIR="$AUDIT_ENV_ROOT/go-tmp"

effective_goflags="$(go env GOFLAGS 2>/dev/null)" \
  || { printf '%s\n' 'AUDIT BROKEN: go env GOFLAGS failed' >&2; exit 2; }
python3 - "$effective_goflags" <<'PY' || exit 2
import shlex
import sys

try:
    tokens = shlex.split(sys.argv[1])
except ValueError as error:
    print("AUDIT BROKEN: GOFLAGS is malformed: " + str(error), file=sys.stderr)
    raise SystemExit(1)

for token in tokens:
    normalized = "-" + token[2:] if token.startswith("--") else token
    name = normalized.split("=", 1)[0]
    if name in {
        "-C", "-modfile", "-overlay", "-toolexec",
        "-args", "-count", "-exec", "-failfast", "-list", "-n", "-run",
        "-short", "-shuffle", "-skip",
    } or name.startswith("-test."):
        print("AUDIT BROKEN: GOFLAGS may bypass or desynchronize the audit: " + token, file=sys.stderr)
        raise SystemExit(1)
PY

target_goos="$(go env GOOS 2>/dev/null)" || target_goos=""
target_goarch="$(go env GOARCH 2>/dev/null)" || target_goarch=""
host_goos="$(go env GOHOSTOS 2>/dev/null)" || host_goos=""
host_goarch="$(go env GOHOSTARCH 2>/dev/null)" || host_goarch=""
[ -n "$target_goos" ] && [ -n "$target_goarch" ] && [ -n "$host_goos" ] && [ -n "$host_goarch" ] \
  || { printf '%s\n' 'AUDIT BROKEN: cannot determine host and target Go platforms' >&2; exit 2; }
[ "$target_goos/$target_goarch" = "$host_goos/$host_goarch" ] \
  || { printf 'AUDIT BROKEN: cross-target Go context %s/%s cannot execute the native audit scanner (%s/%s)\n' \
       "$target_goos" "$target_goarch" "$host_goos" "$host_goarch" >&2
       exit 2; }

hash_probe="$(printf '' | shasum 2>/dev/null | awk '{print $1}')"
case "$hash_probe" in
  *[!0-9a-f]*|'') printf '%s\n' 'AUDIT BROKEN: shasum did not produce a hash' >&2; exit 2 ;;
esac
[ "${#hash_probe}" -eq 40 ] \
  || { printf '%s\n' 'AUDIT BROKEN: shasum produced a malformed hash' >&2; exit 2; }

instrument_sha256() {
  local path="$1" digest
  [ -f "$path" ] && [ ! -L "$path" ] \
    || { printf 'AUDIT BROKEN: missing regular quality instrument %s\n' "$path" >&2; return 1; }
  digest="$(LC_ALL=C shasum -a 256 "$path" 2>/dev/null | awk '{print $1}')"
  case "$digest" in
    *[!0-9a-f]*|'') printf 'AUDIT BROKEN: cannot hash quality instrument %s\n' "$path" >&2; return 1 ;;
  esac
  [ "${#digest}" -eq 64 ] \
    || { printf 'AUDIT BROKEN: malformed quality instrument hash for %s\n' "$path" >&2; return 1; }
  printf '%s\n' "$digest"
}

export_instrument_hash() {
  local name="$1" path="$2" digest
  digest="$(instrument_sha256 "$path")" || return 1
  printf -v "$name" '%s' "$digest"
  export "$name"
}

export_instrument_hash QUALITY_AUDIT_EXPECTED_UPSTREAM_SHA256 "$UPSTREAM" || exit 2
export_instrument_hash QUALITY_AUDIT_EXPECTED_PARSER_SHA256 "$PARSER" || exit 2
export_instrument_hash QUALITY_AUDIT_EXPECTED_WRAPPER_SHA256 "$HERE/quality-audit.sh" || exit 2
export_instrument_hash QUALITY_AUDIT_EXPECTED_CALL_SCANNER_SHA256 "$HERE/go-callscan.go.src" || exit 2
export_instrument_hash QUALITY_AUDIT_EXPECTED_Q06_CONTRACT_SHA256 "$HERE/q06-contract_test.go.src" || exit 2
if [ -f "$(dirname "$UPSTREAM")/audit-report.template.md" ]; then
  export_instrument_hash QUALITY_AUDIT_EXPECTED_REPORT_TEMPLATE_SHA256 \
    "$(dirname "$UPSTREAM")/audit-report.template.md" || exit 2
else
  unset QUALITY_AUDIT_EXPECTED_REPORT_TEMPLATE_SHA256
fi

if [ -z "${GOMODCACHE:-}" ]; then
  GOMODCACHE="$(go env GOMODCACHE 2>/dev/null)"
  [ -n "$GOMODCACHE" ] \
    || { printf '%s\n' 'AUDIT BROKEN: go env GOMODCACHE returned no path' >&2; exit 2; }
  export GOMODCACHE
fi

bash "$UPSTREAM" "${FORWARD[@]}"
upstream_rc=$?
case "$upstream_rc" in
  0|1) ;;
  2) invalidate_generated_outputs || true; exit 2 ;;
  *) printf 'AUDIT BROKEN: upstream returned undocumented exit %s\n' "$upstream_rc" >&2
     invalidate_generated_outputs || true
     exit 2 ;;
esac

[ -f "$REPORT" ] && [ ! -L "$REPORT" ] && [ -s "$REPORT" ] \
  || { printf 'AUDIT BROKEN: upstream wrote no regular report to %s\n' "$REPORT" >&2
       invalidate_generated_outputs || true
       exit 2; }

# The upstream process owns no generated artifact after this point. Move its
# report to a private path and remove any authoritative-looking side outputs it
# may have created before invoking the structured parser.
invalidate_artifact "$JSON" || { invalidate_generated_outputs || true; exit 2; }
invalidate_artifact "$RAW_REPORT" || { invalidate_generated_outputs || true; exit 2; }
invalidate_artifact "$PENDING_JSON" || { invalidate_generated_outputs || true; exit 2; }
invalidate_artifact "$PENDING_RAW" || { invalidate_generated_outputs || true; exit 2; }
if ! atomic_replace_regular "$REPORT" "$PENDING_RAW"; then
  printf 'AUDIT BROKEN: cannot isolate upstream report as %s\n' "$PENDING_RAW" >&2
  invalidate_generated_outputs || true
  exit 2
fi

commit="$(git -C "$REPO_ABS" rev-parse HEAD 2>/dev/null)" \
  || { invalidate_generated_outputs || true
       printf 'AUDIT BROKEN: %s is not a git repository\n' "$REPO_ABS" >&2
       exit 2; }

parser_args=(
  --report "$PENDING_RAW"
  --output "$PENDING_JSON"
  --display-output "$JSON"
  --repo "$REPO_ABS"
  --commit "$commit"
  --upstream "$UPSTREAM"
  --upstream-exit "$upstream_rc"
)
[ -z "$BASELINE" ] || parser_args+=(--baseline "$BASELINE")
[ ! -f "$MANUAL" ] || parser_args+=(--manual-evidence "$MANUAL")
[ "$PARTIAL" = "no" ] || parser_args+=(--partial)

python3 "$PARSER" "${parser_args[@]}"
parser_rc=$?
case "$parser_rc" in
  0|1) ;;
  *) invalidate_artifact "$JSON" || true
     invalidate_artifact "$PENDING_JSON" || true
     parser_rc=2 ;;
esac
if ! atomic_replace_regular "$PENDING_RAW" "$RAW_REPORT"; then
  invalidate_generated_outputs || true
  printf 'AUDIT BROKEN: cannot preserve raw report as %s\n' "$RAW_REPORT" >&2
  exit 2
fi
if [ "$parser_rc" -ne 2 ] && ! atomic_replace_regular "$PENDING_JSON" "$JSON"; then
  invalidate_generated_outputs || true
  printf 'AUDIT BROKEN: cannot publish structured scorecard as %s\n' "$JSON" >&2
  exit 2
fi
printf 'raw upstream scorecard: %s\n' "$RAW_REPORT"
exit "$parser_rc"
