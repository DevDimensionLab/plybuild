#!/bin/sh
set -eu
if [ "$#" -ne 2 ]; then
  echo 'Usage: PLY_DEP02_OLD_BINARY=/preserved/d58/ply-control demo-scoped-continuation.sh /exact/make/ply /private/evidence-directory' >&2
  exit 2
fi
candidate=$1
evidence=$2
repo=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd -P)
case "$candidate:$evidence" in /*:/*) ;; *) echo 'Candidate and evidence directory must be absolute.' >&2; exit 2;; esac
case "$evidence/" in "$repo/"*) echo 'Keep demo evidence outside the clean Task source.' >&2; exit 2;; esac
: "${PLY_DEP02_OLD_BINARY:?Select the preserved original d58 control; no provider launcher is used.}"
umask 077
mkdir -p "$evidence"
"$candidate" version --format json > "$evidence/candidate-version.json"
python3 - "$evidence/candidate-version.json" "$PLY_DEP02_OLD_BINARY" "$repo" <<'PY'
import hashlib,json,pathlib,subprocess,sys
v=json.loads(pathlib.Path(sys.argv[1]).read_text())
repo=sys.argv[3]
head=subprocess.check_output(['git','-C',repo,'rev-parse','HEAD'],text=True).strip()
assert v['source_commit']==head and v['dirty'] is False, 'Demo requires exact clean make binary'
assert not subprocess.check_output(['git','-C',repo,'status','--porcelain']), 'Demo requires clean candidate source'
old=pathlib.Path(sys.argv[2])
assert hashlib.sha256(old.read_bytes()).hexdigest()=='7044a9996acd208a2bacf9bc849e3f486b72a765aec5f4b1feb888458934cb32', 'Use the exact preserved original PWA d58 control'
PY
export PLY_DELIVERY_TEST_BINARY="$candidate"
export PLY_DEP02_EVIDENCE_DIR="$evidence"
cd "$repo"
go test -mod=readonly ./internal/taskrun -run '^TestDEP02' -count=1 -v
