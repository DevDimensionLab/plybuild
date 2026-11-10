compat_source_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)
compat_root="$compat_source_root/.build/compatibility"
compat_guard_dir() {
	if [[ -L $2 || (-e $2 && ! -d $2) ]]; then
		printf '%s: refusing symlinked or non-directory %s\n' "$1" "$2" >&2
		exit 2
	fi
	mkdir -p -- "$2" || exit 2
}
compat_run_root() {
	compat_guard_dir "$1" "$compat_source_root/.build"
	compat_guard_dir "$1" "$compat_root"
	run_root=$(mktemp -d "$compat_root/$2.XXXXXX") || exit 2
	trap compat_cleanup EXIT
	trap 'exit 129' HUP
	trap 'exit 130' INT
	trap 'exit 143' TERM
	mkdir "$run_root/tmp"
	export TMPDIR="$run_root/tmp/" GOTMPDIR="$run_root/tmp"
}
compat_cleanup() {
	find "$run_root" -type d -exec chmod u+rwx {} + 2>/dev/null || true
	rm -rf -- "$run_root"
}
