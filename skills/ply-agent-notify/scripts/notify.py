#!/usr/bin/env python3
"""Preserve an authorized agent event, preview once, apply once and read back through Ply."""
import argparse
import fcntl
import hashlib
import json
import os
import re
import subprocess
import sys
import uuid
import unicodedata
from datetime import datetime, timezone
from pathlib import Path
from zoneinfo import ZoneInfo, ZoneInfoNotFoundError
sys.dont_write_bytecode = True
from common import CONFIG, ENV, Failure, credential, decode, directory, encode, locator, physical, private, read_file, safe_public, write_once

STATES = {'transport_acknowledged', 'rejected', 'not_sent', 'rate_limited', 'unknown'}
ID = re.compile(r'ntf_[0-9a-f]{64}\Z')
DIGEST = re.compile(r'sha256:[0-9a-f]{64}\Z')
ZONE = re.compile(r'[A-Za-z0-9_+-]+(?:/[A-Za-z0-9_+-]+)*\Z')
STATUS = {'agent_finished': {'ready_for_review', 'ready_for_your_check', 'done'},
          'agent_stopped': {'stopped'}, 'feedback_required': {'needs_answer'}}


def compact_text(value, limit):
    if not isinstance(value, str) or not value or value.strip() != value or len(value) > limit:
        raise Failure('invalid_compact_text')
    if any(unicodedata.category(c) in {'Cc', 'Cf', 'Cs', 'Zl', 'Zp'} for c in value) or any(c in value for c in '<>') or any(s in value.lower() for s in ['@channel', '@here', '@everyone']):
        raise Failure('unsafe_compact_text')


def compact_presentation(task, event, args):
    if task['provider'] not in {'codex', 'claude'}:
        raise Failure('explicit_provider_required')
    origin, root = physical(task['origin_cwd']), physical(task['context_root'])
    for path in [root, origin]:
        with directory(path):
            pass
    if Path.cwd() != origin:
        raise Failure('origin_cwd_mismatch')
    try:
        context = origin.relative_to(root).as_posix()
    except ValueError:
        raise Failure('origin_outside_context_root') from None
    compact_text(context, 96)
    zone = task['timezone']
    if not isinstance(zone, str) or not ZONE.fullmatch(zone) or zone == 'Local' or len(zone) > 128:
        raise Failure('invalid_timezone')
    if set(event) != {'kind', 'activity', 'run', 'event_id', 'event_type', 'phase', 'occurred_at', 'public'}:
        raise Failure('invalid_event_fields')
    public = event['public']
    if not isinstance(public, dict) or set(public) != {'task_title', 'summary', 'next_action', 'next_actor', 'status'}:
        raise Failure('invalid_public_fields')
    for name, limit in [('task_title', 60), ('summary', 100), ('next_action', 140)]:
        compact_text(public[name], limit)
    if public['status'] not in STATUS.get(event['event_type'], set()) or public['next_actor'] not in {'user', 'coordinator'}:
        raise Failure('invalid_event_status')
    if event['phase'] not in {'before_start', 'after_start'} or event['phase'] == 'before_start' and (event['event_type'] != 'agent_stopped' or args.start_receipt):
        raise Failure('invalid_event_phase')
    if event['event_type'] == 'agent_finished' and not args.report:
        raise Failure('report_required')
    if event['event_type'] == 'feedback_required' and (public['next_actor'] != 'user' or '?' not in public['summary']):
        raise Failure('necessary_question_required')
    try:
        instant = datetime.strptime(event['occurred_at'], '%Y-%m-%dT%H:%M:%SZ').replace(tzinfo=timezone.utc)
        if instant.strftime('%Y-%m-%dT%H:%M:%SZ') != event['occurred_at']:
            raise ValueError()
    except (ValueError, TypeError):
        raise Failure('invalid_event_time') from None
    return {'provider': task['provider'], 'origin_cwd': str(origin), 'context_root': str(root),
            'context': context, 'timezone': zone}, instant


def local_stamp(instant, zone):
    try:
        return instant.astimezone(ZoneInfo(zone)).isoformat(timespec='seconds').replace('+00:00', 'Z')
    except (ZoneInfoNotFoundError, ValueError, OverflowError):
        raise Failure('unknown_timezone_or_unrepresentable_time') from None


def outcome(state, reason, **fields):
    return {'kind': 'PlyAgentNotificationDelivery@1', 'state': state, 'reason': reason,
            'product_outcome': 'unchanged', **fields}


def call(ply, args, secret=None):
    env = os.environ.copy()
    env.pop(ENV, None)
    if secret is not None:
        env[ENV] = secret
    # No shell, eval, output forwarding or secret in argv. Show gets no credential.
    result = subprocess.run([ply, 'workflow', 'notification', *args, '--format', 'json'],
                            env=env, capture_output=True, timeout=30)
    try:
        value = decode(result.stdout)
    except Failure:
        value = {}
    return result.returncode, value


def show(ply, route, notification_id):
    code, value = call(ply, ['show', notification_id, '--route', route])
    state = value.get('state')
    if code != 0 or value.get('kind') != 'ply.workflow.notification' or value.get('notification_id') != notification_id or state not in STATES:
        return code, 'unknown'
    return code, state


def deliver(args):
    if args.trigger != 'event':
        return outcome('suppressed', args.trigger)
    task = decode(read_file(args.task, 64 << 10))
    required = {'kind', 'activity', 'run', 'worktree', 'handoff', 'route', 'event_root', 'actor_claim', 'notification_authorized'}
    compact = task.get('kind') == 'PlyAgentNotificationTask@2'
    if compact:
        required |= {'provider', 'origin_cwd', 'context_root', 'timezone'}
    if set(task) != required or task['kind'] not in {'PlyAgentNotificationTask@1', 'PlyAgentNotificationTask@2'} or task['notification_authorized'] is not True:
        raise Failure('task_route_authority_required')
    event = decode(read_file(args.event))
    event_kind = 'PlyAgentNotificationEvent@2' if compact else 'PlyAgentNotificationEvent@1'
    if event.get('kind') != event_kind or event.get('activity') != task['activity'] or event.get('run') != task['run']:
        raise Failure('event_identity_mismatch')
    if compact:
        presentation, instant = compact_presentation(task, event, args)
    route = decode(read_file(task['route'], 16 << 10))
    if route.get('webhook_env') != ENV:
        raise Failure('named_credential_required')
    record = locator(args.event)
    record['kind'] = event_kind
    source = {k: task[k] for k in ['activity', 'run', 'worktree']}
    source.update(kind='external', handoff=locator(task['handoff']))
    for name in ['start_receipt', 'report']:
        path = getattr(args, name)
        if path:
            source[name] = locator(path, metadata=True)
    request = {'kind':'ply.workflow.notification-request', 'schema_version':3 if compact else 2,
               'route':route['name'], 'source':source,
               'event':{'id':event['event_id'], 'type':event['event_type'], 'phase':event['phase'],
                        'occurred_at':event['occurred_at'], 'record':record},
               'sender':{'actor_claim':task['actor_claim']}, 'public':event['public']}
    root = physical(task['event_root'])
    # The version-independent key is also used by older Task@1 callers.
    key = hashlib.sha256(encode({'activity':task['activity'], 'run':task['run'], 'event':event['event_id']})).hexdigest()
    folder = root/key
    if compact:
        request['presentation'] = presentation
        request_path = folder/'request.json'
        if not os.path.lexists(request_path):
            # Reject an unknown zone on a fresh event before creating its folder.
            # An existing request is read only after taking the event lock below.
            presentation['local_occurred_at'] = local_stamp(instant, presentation['timezone'])
    secret, credential_error = None, None
    try:
        secret = credential(args.credential_file)
    except (Failure, OSError, UnicodeError):
        credential_error = 'credential_missing_invalid_or_unsafe'
    if not all(safe_public(value, secret or os.environ.get(ENV, '')) for value in [request, task, route]):
        raise Failure('unsafe_notification_text')
    # Route and port revision deliberately do not change this identity. Keep this
    # event root fixed across callers and restarts, just like the route state root.
    with directory(root, create=True, secure=True):
        pass
    with directory(folder, create=True, secure=True) as fd:
        try:
            lock = os.open('.lock', os.O_RDWR | os.O_CREAT | os.O_NOFOLLOW | os.O_NONBLOCK, 0o600, dir_fd=fd)
        except FileNotFoundError:
            # A concurrent O_CREAT/O_NOFOLLOW can report ENOENT on macOS.
            # Treat path replacement or another caller conservatively; no retry.
            return outcome('unknown', 'event_changed_or_busy_preserve_local_return')
        try:
            import stat
            info = os.fstat(lock)
            private(info, 0o600)
            if not stat.S_ISREG(info.st_mode) or info.st_nlink != 1:
                raise Failure('unsafe_event_lock')
            try:
                fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
            except BlockingIOError:
                return outcome('unknown', 'event_busy_preserve_local_return')
            request_path = folder/'request.json'
            attempt_path = folder/'apply-started.json'
            if compact:
                if os.path.lexists(request_path):
                    prior = decode(read_file(request_path, mode=0o600))
                    stored = prior.get('presentation')
                    if prior.get('schema_version') != 3 or not isinstance(stored, dict):
                        raise Failure('existing_bytes_conflict')
                    frozen = stored.get('local_occurred_at')
                    try:
                        local = datetime.fromisoformat(frozen.replace('Z', '+00:00'))
                        if local.utcoffset() is None or local != instant or local.isoformat(timespec='seconds').replace('+00:00', 'Z') != frozen:
                            raise ValueError()
                    except (AttributeError, TypeError, ValueError):
                        raise Failure('invalid_preserved_local_time') from None
                    presentation['local_occurred_at'] = frozen
                elif 'local_occurred_at' not in presentation:
                    presentation['local_occurred_at'] = local_stamp(instant, presentation['timezone'])
                if not safe_public(request, secret or os.environ.get(ENV, '')):
                    raise Failure('unsafe_notification_text')
            # Freeze route locator as well as request bytes, before any send.
            write_once(folder/'binding.json', encode({'route':task['route'], 'state_root':route['state_root'], 'route_sha256':locator(task['route'])['sha256']}))
            write_once(request_path, encode(request))
            details = {'event_id':event['event_id'], 'request':str(request_path)}
            if os.path.lexists(attempt_path):
                prior = decode(read_file(attempt_path, mode=0o600))
                notification_id = prior.get('notification_id')
                if not isinstance(notification_id, str) or not ID.fullmatch(notification_id):
                    raise Failure('invalid_attempt_record', 'unknown')
                code, state = show(args.ply_bin, task['route'], notification_id)
                result = outcome(state, 'preserved_attempt_readback', notification_id=notification_id, show_exit=code, **details)
            elif credential_error:
                result = outcome('not_attempted', credential_error, **details)
            else:
                command = ['send', '--file', str(request_path), '--route', task['route']]
                code, preview = call(args.ply_bin, command+['--check'], secret)
                notification_id = preview.get('notification_id')
                confirmation = preview.get('confirmation')
                if code != 0 or preview.get('kind') != 'ply.workflow.notification-preview' or preview.get('allowed') is not True or not isinstance(notification_id, str) or not ID.fullmatch(notification_id) or not isinstance(confirmation, str) or not DIGEST.fullmatch(confirmation):
                    state = preview.get('state', 'not_attempted')
                    result = outcome(state if state in STATES else 'not_attempted', 'preview_rejected', preview_exit=code, **details)
                elif preview.get('existing_state') is not None:
                    show_code, state = show(args.ply_bin, task['route'], notification_id)
                    result = outcome(state, 'preserved_ply_readback', notification_id=notification_id, preview_exit=code, show_exit=show_code, **details)
                else:
                    write_once(attempt_path, encode({'notification_id':notification_id, 'confirmation':confirmation}))
                    # A helper death after this marker never permits a second apply,
                    # even if the original child has not yet reserved Ply state.
                    apply_code, applied = call(args.ply_bin, command+['--apply', '--confirm', confirmation], secret)
                    show_code, state = show(args.ply_bin, task['route'], notification_id)
                    if applied.get('kind') != 'ply.workflow.notification' or applied.get('notification_id') != notification_id or applied.get('state') != state:
                        state = 'unknown'
                    result = outcome(state, 'ply_readback', notification_id=notification_id, preview_exit=code,
                                     apply_exit=apply_code, show_exit=show_code, **details)
            receipt = folder/('delivery-'+uuid.uuid4().hex+'.json')
            write_once(receipt, encode(result))
            return {**result, 'receipt':str(receipt)}
        except (OSError, subprocess.SubprocessError, Failure) as err:
            # Never claim not-attempted after the helper's durable apply boundary.
            if os.path.lexists(folder/'apply-started.json'):
                reason = 'binding_conflict_preserved_attempt' if isinstance(err, Failure) and err.reason == 'existing_bytes_conflict' else 'attempt_preserved_readback_unavailable'
                return outcome('unknown', reason)
            if isinstance(err, Failure):
                raise
            raise Failure('local_io_or_cli_unavailable') from None
        finally:
            os.close(lock)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--task', required=True, help='Authorized task context JSON')
    parser.add_argument('--event', help='Already preserved strict event-record JSON')
    parser.add_argument('--start-receipt', help='Actual start receipt JSON, when available')
    parser.add_argument('--report', help='Actual report JSON; required for agent_finished')
    parser.add_argument('--ply-bin', default='ply', help='Ply executable (one argv entry)')
    parser.add_argument('--credential-file', default=str(CONFIG/'webhook'), help='Exact private credential file, used only when the environment variable is absent')
    parser.add_argument('--trigger', choices=['event', 'optional_feedback', 'internal_correction', 'waiting_for_feedback'], default='event')
    parser.add_argument('--port-revision', type=int, help='Caller annotation only; never changes event identity or permits resend')
    args = parser.parse_args()
    if args.trigger == 'event' and not args.event:
        parser.error('--event is required for an event notification')
    try:
        result = deliver(args)
    except Failure as err:
        result = outcome(err.state, err.reason)
    except (OSError, ValueError, KeyError, TypeError, subprocess.SubprocessError):
        result = outcome('not_attempted', 'invalid_local_input_or_cli')
    print(json.dumps(result, ensure_ascii=False))
    if result['state'] in ['transport_acknowledged', 'suppressed']:
        return 0
    return 2 if result['state'] == 'not_attempted' else 5


if __name__ == '__main__':
    sys.exit(main())
