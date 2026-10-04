#!/usr/bin/env python3
"""Install only ply-agent-notify; preserve any existing different destination."""
import argparse
import json
import os
import sys
from pathlib import Path
sys.dont_write_bytecode = True
from common import DEST, Failure, directory, physical, read_file, write_once

FILES = {'SKILL.md': 0o600, 'references/task-template.json': 0o600,
         'references/usage.md': 0o600, 'scripts/common.py': 0o600,
         'scripts/notify.py': 0o700, 'scripts/bootstrap.py': 0o700, 'scripts/install.py': 0o700}


def install(source, destination):
    source, destination = physical(source), physical(destination)
    data = {name: read_file(source/name) for name in FILES}
    if os.path.lexists(destination):
        with directory(destination, secure=True):
            pass
        actual = set()
        for parent, dirs, files in os.walk(destination, followlinks=False):
            with directory(parent, secure=True):
                pass
            for name in dirs:
                with directory(Path(parent)/name, secure=True):
                    pass
            actual.update(str((Path(parent)/f).relative_to(destination)) for f in files)
        if actual != set(FILES):
            raise Failure('existing_install_conflict')
        for name, mode in FILES.items():
            if read_file(destination/name, mode=mode) != data[name]:
                raise Failure('existing_install_conflict')
        return
    with directory(destination, create=True, secure=True):
        pass
    for child in ['scripts', 'references']:
        with directory(destination/child, create=True, secure=True):
            pass
    for name, mode in FILES.items():
        write_once(destination/name, data[name], mode)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--destination', default=str(DEST), help='Exact skill directory; parent must exist')
    args = parser.parse_args()
    try:
        install(Path(__file__).resolve().parent.parent, args.destination)
        print(json.dumps({'kind':'PlyAgentNotificationInstall@1', 'state':'installed', 'destination':args.destination}))
        return 0
    except (Failure, OSError, ValueError):
        print(json.dumps({'kind':'PlyAgentNotificationInstall@1', 'state':'not_attempted', 'reason':'Destination is unsafe or conflicts with the skill; existing bytes were preserved.'}))
        return 2


if __name__ == '__main__':
    sys.exit(main())
