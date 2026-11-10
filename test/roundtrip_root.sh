roundtrip_root() {
	local root=${PLY_ROUNDTRIP_ROOT:-} dir
	unsuitable() { printf '%s: unsuitable PLY_ROUNDTRIP_ROOT %s: %s\n' "$1" "'$root'" "$2" >&2; exit 2; }
	case $root in /*) ;; *) unsuitable "$1" 'set it to an existing absolute directory outside Git and Ply workspaces' ;; esac
	{ [ -d "$root" ] && [ ! -L "${root%/}" ] && [ -w "$root" ]; } || unsuitable "$1" 'not a writable, non-symlinked directory'
	root=$(cd "$root" && pwd -P)
	! git -C "$root" rev-parse --git-dir >/dev/null 2>&1 || unsuitable "$1" 'inside a Git repository'
	dir=$root; while :; do [ ! -e "$dir/.ply" ] || unsuitable "$1" "inside Ply workspace $dir"; [ "$dir" != / ] || break; dir=$(dirname "$dir"); done
	run_root=$(mktemp -d "$root/$2.XXXXXX")
	trap 'chmod -R u+rwx "$run_root" 2>/dev/null || true; rm -rf -- "$run_root"' EXIT
	trap 'exit 129' HUP
	trap 'exit 130' INT
	trap 'exit 143' TERM
	mkdir "$run_root/tmp"
	export TMPDIR="$run_root/tmp/" GOTMPDIR="$run_root/tmp"
}
