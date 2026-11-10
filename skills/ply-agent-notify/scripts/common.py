"""Private, physical files and strict JSON for the notification skill (POSIX)."""
import hashlib
import json
import os
import re
import stat
from contextlib import contextmanager
from pathlib import Path

CONFIG = Path('/Users/perottochristensen/.config/ply/notifications/ply-log')
STATE = Path('/Users/perottochristensen/.local/state/ply/agent-notifications/ply-log')
DEST = Path('/Users/perottochristensen/.agents/skills/ply-agent-notify')
ENV = 'PLY_SLACK_WEBHOOK_URL'
WEBHOOK = re.compile(r'https://hooks\.slack\.com/services/[A-Za-z0-9_-]+/[A-Za-z0-9_-]+/[A-Za-z0-9_-]+\Z')


class Failure(Exception):
    def __init__(self, reason, state='not_attempted'):
        self.reason, self.state = reason, state


def physical(path):
    path = os.fspath(path)
    if not path.startswith('/') or os.path.normpath(path) != path:
        raise Failure('physical_absolute_path_required')
    return Path(path)


def private(info, mode):
    if info.st_uid != os.getuid() or stat.S_IMODE(info.st_mode) != mode:
        raise Failure('unsafe_owner_or_mode')


@contextmanager
def directory(path, create=False, parents=False, secure=False):
    path = physical(path)
    fd = os.open('/', os.O_RDONLY | os.O_DIRECTORY)
    try:
        parts = path.parts[1:]
        for i, part in enumerate(parts):
            try:
                nxt = os.open(part, os.O_RDONLY | os.O_DIRECTORY | os.O_NOFOLLOW, dir_fd=fd)
            except FileNotFoundError:
                if not create or (i != len(parts)-1 and not parents):
                    raise
                try:
                    os.mkdir(part, 0o700, dir_fd=fd)
                    os.fsync(fd)
                except FileExistsError:
                    pass
                nxt = os.open(part, os.O_RDONLY | os.O_DIRECTORY | os.O_NOFOLLOW, dir_fd=fd)
            os.close(fd)
            fd = nxt
        if secure:
            private(os.fstat(fd), 0o700)
        yield fd
    finally:
        os.close(fd)


def read_file(path, limit=1 << 20, mode=None):
    path = physical(path)
    with directory(path.parent, secure=mode is not None) as parent:
        fd = os.open(path.name, os.O_RDONLY | os.O_NOFOLLOW | os.O_NONBLOCK, dir_fd=parent)
        try:
            before = os.fstat(fd)
            if not stat.S_ISREG(before.st_mode) or before.st_size > limit:
                raise Failure('invalid_file')
            if mode is not None:
                private(before, mode)
                if before.st_nlink != 1:
                    raise Failure('linked_private_file')
            chunks = []
            size = 0
            while size <= limit:
                chunk = os.read(fd, min(65536, limit+1-size))
                if not chunk:
                    break
                chunks.append(chunk)
                size += len(chunk)
            after = os.fstat(fd)
            leaf = os.stat(path.name, dir_fd=parent, follow_symlinks=False)
            if size > limit or (before.st_size, before.st_mtime_ns, before.st_ctime_ns) != (after.st_size, after.st_mtime_ns, after.st_ctime_ns) or (after.st_ino, after.st_dev) != (leaf.st_ino, leaf.st_dev):
                raise Failure('file_changed')
            with directory(path.parent) as current:
                a, b = os.fstat(parent), os.fstat(current)
                if (a.st_dev, a.st_ino) != (b.st_dev, b.st_ino):
                    raise Failure('parent_changed')
            return b''.join(chunks)
        finally:
            os.close(fd)


def encode(value):
    return (json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(',', ':'))+'\n').encode('utf-8')


def decode(raw):
    def pairs(items):
        result = {}
        for key, value in items:
            if key in result:
                raise Failure('duplicate_json_key')
            result[key] = value
        return result
    def bad_number(_):
        raise Failure('invalid_json_number')
    try:
        text = raw.decode('utf-8')
        value = json.loads(text, object_pairs_hook=pairs, parse_constant=bad_number)
        # Reject lone escaped surrogates too.
        json.dumps(value, ensure_ascii=False).encode('utf-8')
        if not isinstance(value, dict):
            raise Failure('json_object_required')
        return value
    except (ValueError, UnicodeError):
        raise Failure('invalid_json') from None


def write_once(path, data, mode=0o600):
    """Exclusive durable creation; existing different bytes are never replaced."""
    path = physical(path)
    with directory(path.parent, secure=True) as parent:
        try:
            fd = os.open(path.name, os.O_WRONLY | os.O_CREAT | os.O_EXCL | os.O_NOFOLLOW, mode, dir_fd=parent)
        except FileExistsError:
            if read_file(path, max(len(data), 1 << 20), mode) != data:
                raise Failure('existing_bytes_conflict')
            return False
        with os.fdopen(fd, 'wb') as out:
            out.write(data)
            out.flush()
            os.fsync(out.fileno())
        os.fsync(parent)
        return True


def digest(data):
    return 'sha256:'+hashlib.sha256(data).hexdigest()


def locator(path, metadata=False):
    raw = read_file(path, 16 << 20 if metadata else 1 << 20)
    out = {'path': str(physical(path)), 'sha256': digest(raw)}
    if metadata:
        obj = decode(raw)
        out['kind'] = obj['kind']
    return out


def credential(path):
    # Presence, even an empty or invalid value, forbids hidden file fallback.
    if ENV in os.environ:
        value = os.environ[ENV]
    else:
        value = read_file(path, 4096, 0o600).decode('utf-8')
    if not WEBHOOK.fullmatch(value):
        raise Failure('invalid_credential')
    return value


def safe_public(value, secret):
    if isinstance(value, str):
        return not re.search(r'https?://hooks\.(slack\.com|slack-gov\.com)\b', value, re.I) and (not secret or secret not in value)
    if isinstance(value, dict):
        return all(safe_public(k, secret) and safe_public(v, secret) for k, v in value.items())
    if isinstance(value, list):
        return all(safe_public(v, secret) for v in value)
    return True
