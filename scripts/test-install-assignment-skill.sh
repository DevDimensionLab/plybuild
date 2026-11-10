#!/usr/bin/env bash
# Meta-test of install-assignment-skill.sh. make preflight requires every
# production file in scripts/ to have a scripts/test-<name> meta-test, so this
# file only runs test/assignment_skill_install_test.sh in a fresh root.
# TMPDIR must be set explicitly to an absolute directory outside /tmp (the
# owner's .build or external synthetic root); there is no fallback.

set -euo pipefail

fail() {
	printf 'assignment skill meta-test: %s\n' "$*" >&2
	exit 1
}

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
[[ -n "${TMPDIR:-}" ]] || fail 'TMPDIR must be set explicitly'
[[ "$TMPDIR" == /* && -d "$TMPDIR" && ! -L "$TMPDIR" ]] || fail 'TMPDIR must be an existing absolute directory'
physical=$(cd "$TMPDIR" && pwd -P)
case $physical in
/tmp | /tmp/* | /private/tmp | /private/tmp/*) fail 'TMPDIR must not be under /tmp' ;;
esac
temp_root=$(mktemp -d "$physical/ply-assignment-skill-test.XXXXXX")
cleanup_temp_root() {
	find "$temp_root" -type d -exec chmod u+rwx {} + 2>/dev/null || true
	rm -rf -- "$temp_root"
}
trap cleanup_temp_root EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

bash "$repo_root/test/assignment_skill_install_test.sh" "$temp_root"
