#!/bin/sh
# Run from the exact native candidate Task cwd; no heredoc or global temp file.
set -eu
exec python3 -I -B ./test/mr05-projection-acceptance.py
