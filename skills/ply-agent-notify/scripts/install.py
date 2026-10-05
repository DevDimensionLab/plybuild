#!/usr/bin/env python3
"""Install ply-agent-notify, or explicitly upgrade an exact owned copy with backup."""
import argparse
import json
import os
import sys
import uuid
from pathlib import Path
sys.dont_write_bytecode = True
from common import DEST, Failure, directory, physical, read_file, write_once

FILES = {'SKILL.md': 0o600, 'references/task-template.json': 0o600,
         'references/usage.md': 0o600, 'scripts/common.py': 0o600,
         'scripts/notify.py': 0o700, 'scripts/bootstrap.py': 0o700, 'scripts/install.py': 0o700}


def verify(destination, data):
    with directory(destination, secure=True):
        pass
    actual, actual_dirs = set(), set()
    for parent, dirs, files in os.walk(destination, followlinks=False):
        with directory(parent, secure=True):
            pass
        for name in dirs:
            path = Path(parent)/name
            with directory(path, secure=True):
                pass
            actual_dirs.add(str(path.relative_to(destination)))
        actual.update(str((Path(parent)/f).relative_to(destination)) for f in files)
    if actual != set(FILES) or actual_dirs != {'scripts', 'references'}:
        raise Failure('existing_install_conflict')
    for name, mode in FILES.items():
        if read_file(destination/name, mode=mode) != data[name]:
            raise Failure('existing_install_conflict')


def create(destination, data):
    # Exclusive root creation prevents a pre-existing backup being reused.
    with directory(destination.parent) as parent:
        os.mkdir(destination.name, 0o700, dir_fd=parent)
        os.fsync(parent)
    with directory(destination, create=True, secure=True):
        pass
    for child in ['scripts', 'references']:
        with directory(destination/child, create=True, secure=True):
            pass
    for name, mode in FILES.items():
        write_once(destination/name, data[name], mode)


def install(source, destination, upgrade_from=None, backup=None):
    source, destination = physical(source), physical(destination)
    data = {name: read_file(source/name) for name in FILES}
    if bool(upgrade_from) != bool(backup):
        raise Failure('upgrade_requires_source_and_backup')
    if not upgrade_from:
        if os.path.lexists(destination):
            verify(destination, data)
        else:
            create(destination, data)
        return
    previous, backup = physical(upgrade_from), physical(backup)
    if backup == destination or destination in backup.parents or backup in destination.parents or os.path.lexists(backup):
        raise Failure('new_separate_backup_required')
    old = {name: read_file(previous/name) for name in FILES}
    verify(destination, old)
    create(backup, old)
    verify(backup, old)
    verify(destination, old)
    # Only the enumerated owned files are replaced, with their exact expected
    # bytes/modes checked first. A failure preserves the complete prior backup.
    for name, mode in FILES.items():
        if data[name] == old[name]:
            continue
        path = destination/name
        if read_file(path, mode=mode) != old[name]:
            raise Failure('existing_install_conflict')
        with directory(path.parent, secure=True) as parent:
            pending = '.upgrade-'+uuid.uuid4().hex
            fd = os.open(pending, os.O_WRONLY | os.O_CREAT | os.O_EXCL | os.O_NOFOLLOW, mode, dir_fd=parent)
            with os.fdopen(fd, 'wb') as out:
                out.write(data[name]); out.flush(); os.fsync(out.fileno())
            os.replace(pending, path.name, src_dir_fd=parent, dst_dir_fd=parent)
            os.fsync(parent)
    verify(destination, data)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--destination', default=str(DEST), help='Exact skill directory; parent must exist')
    parser.add_argument('--upgrade-from', help='Reviewed previous skill source; installed bytes and modes must match')
    parser.add_argument('--backup', help='New separate private backup directory; required with --upgrade-from')
    args = parser.parse_args()
    try:
        install(Path(__file__).resolve().parent.parent, args.destination, args.upgrade_from, args.backup)
        print(json.dumps({'kind':'PlyAgentNotificationInstall@1', 'state':'installed', 'destination':args.destination}))
        return 0
    except (Failure, OSError, ValueError):
        print(json.dumps({'kind':'PlyAgentNotificationInstall@1', 'state':'unknown' if args.upgrade_from else 'not_attempted', 'reason':'Installation did not complete. Inspect the destination and any preserved backup before recovery; no automatic repair was attempted.'}))
        return 2


if __name__ == '__main__':
    sys.exit(main())
