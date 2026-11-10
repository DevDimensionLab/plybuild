#!/usr/bin/env python3
"""Preserve an authorized agent event, preview once, apply once and read back through Ply."""
import argparse
import fcntl
import hashlib
import json
import os
import re
import shutil
import subprocess
import sys
import time
import unicodedata
from datetime import datetime, timezone
from pathlib import Path
from zoneinfo import ZoneInfo, ZoneInfoNotFoundError
sys.dont_write_bytecode = True
from common import CONFIG, ENV, Failure, credential, decode, digest, directory, encode, locator, physical, private, read_file, safe_public, write_once

STATES = {'transport_acknowledged', 'rejected', 'not_sent', 'rate_limited', 'unknown'}
ID = re.compile(r'ntf_[0-9a-f]{64}\Z')
DIGEST = re.compile(r'sha256:[0-9a-f]{64}\Z')
ZONE = re.compile(r'[A-Za-z0-9_+-]+(?:/[A-Za-z0-9_+-]+)*\Z')
STATUS = {'agent_finished': {'ready_for_review', 'ready_for_your_check', 'done'},
          'agent_stopped': {'stopped'}, 'feedback_required': {'needs_answer'}}
# Herdr origin lookup: public read-only calls, bounded per call and in total.
HERDR_LABEL = 32
HERDR_CALL_TIMEOUT = 1.0
HERDR_DEADLINE = 2.5
HERDR_OUTPUT_LIMIT = 64 << 10
HERDR_ANCESTORS = 64
HERDR_CANDIDATES = 8
HERDR_UNSAFE = {'Cc', 'Cf', 'Cs', 'Zl', 'Zp'}
# The owner's optional app link beside the route; same grammar as Ply's request-v6.
APP_LINK_FILE = 'app-link.json'
APP_LINK = re.compile(r'https://[a-z0-9](?:[a-z0-9-]*[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)+(?:/[A-Za-z0-9._~-]+)*/?\Z')
EVENT_TYPE = {status: kind for kind, statuses in STATUS.items() for status in statuses}
UTC_FORMAT = '%Y-%m-%dT%H:%M:%SZ'
SHORT_REQUIRED = ['activity', 'run', 'event_id', 'status', 'title', 'summary', 'next_action', 'provider']


def short_binding(args):
    """The task and event values of one notification, from arguments only."""
    if not all(getattr(args, name) for name in SHORT_REQUIRED):
        raise Failure('short_form_fields_required')
    cwd = Path.cwd()
    task = {'activity': args.activity, 'run': args.run, 'worktree': str(cwd), 'route': args.route,
            'event_root': str(events_root(args.route)), 'actor_claim': args.provider, 'provider': args.provider,
            'origin_cwd': str(cwd), 'context_root': args.context_root or str(cwd.parent), 'timezone': args.timezone}
    event = {'event_id': args.event_id, 'event_type': EVENT_TYPE.get(args.status), 'phase': args.phase,
             'occurred_at': datetime.now(timezone.utc).strftime(UTC_FORMAT),
             'public': {'status': args.status, 'task_title': args.title, 'summary': args.summary,
                        'next_action': args.next_action, 'next_actor': args.next_actor}}
    return task, event


def events_root(route):
    # Beside the route file, not its state root: a changed state_root keeps this
    # folder and meets the frozen binding. Ply refuses requests under its state root.
    # A route file in another directory is another namespace; this guards one route.
    return physical(route).parent/'notification-events'


class HerdrUnproven(Exception):
    pass


def herdr_label(value, secret):
    """A public display label, or None when it is unsafe."""
    if not isinstance(value, str):
        return None
    text = unicodedata.normalize('NFC', value)
    # Ordinary whitespace collapses; every other control or separator rejects.
    if any(unicodedata.category(c) in HERDR_UNSAFE for c in text if c not in '\t\n\r\f\v'):
        return None
    text = ' '.join(text.split())
    lower = text.lower()
    if not text or any(c in text for c in '<>') or any(s in lower for s in ['@channel', '@here', '@everyone']) or not safe_public(text, secret):
        return None
    if len(text) > HERDR_LABEL:
        text = text[:HERDR_LABEL-1].rstrip()+'…'
    return text


def herdr_run(argv, timeout):
    # argv only and no shell. The exit status is ignored: the JSON error object decides.
    done = subprocess.run(argv, stdin=subprocess.DEVNULL, stdout=subprocess.PIPE, stderr=subprocess.DEVNULL, timeout=timeout)
    if len(done.stdout) > HERDR_OUTPUT_LIMIT:
        raise HerdrUnproven('output_limit')
    return done.stdout.decode('utf-8')


def herdr_ancestors(text):
    """This process and its live ancestors, nearest first."""
    parent = {}
    for line in text.splitlines():
        fields = line.split()
        if len(fields) == 2:
            parent[int(fields[0])] = int(fields[1])
    chain = [os.getpid()]
    while len(chain) < HERDR_ANCESTORS and parent.get(chain[-1], 0) > 1:
        chain.append(parent[chain[-1]])
    return chain


def herdr_origin(provider, secret):
    """{'workspace', 'tab'} for a proven origin, otherwise None. Never raises."""
    herdr = os.environ.get('HERDR_BIN_PATH') or shutil.which('herdr')
    if os.environ.get('HERDR_ENV') != '1' or os.name != 'posix' or not herdr or not os.access(herdr, os.X_OK):
        return None
    end = time.monotonic()+HERDR_DEADLINE
    try:
        return herdr_lookup(provider, secret, herdr, end)
    except Exception:
        # Best effort: any failure sends exactly the notification sent without Herdr.
        return None


def herdr_lookup(provider, secret, herdr, end):
    def call(argv):
        left = end-time.monotonic()
        if left <= 0:
            raise HerdrUnproven('deadline')
        return herdr_run(argv, min(HERDR_CALL_TIMEOUT, left))

    def result(argv, key):
        reply = json.loads(call([herdr, *argv]))
        if not isinstance(reply, dict) or 'error' in reply:
            raise HerdrUnproven('herdr_error')
        return reply['result'][key]

    chain = herdr_ancestors(call(['ps', '-A', '-o', 'pid=,ppid=']))

    def distance(pane_id):
        # Only shell_pid is read; process-info also holds argv, which is never kept.
        pid = result(['pane', 'process-info', '--pane', pane_id], 'process_info')['shell_pid']
        return chain.index(pid) if isinstance(pid, int) and pid > 1 and pid in chain else None

    pane = None
    try:
        current = result(['pane', 'current', '--current'], 'pane')
        if distance(current['pane_id']) is not None:
            pane = current
    except Exception:
        # A stale inherited pane ID, or one that fails the proof, moves to the agent list.
        pane = None
    if pane is None:
        proven = {}
        agents = [a for a in result(['agent', 'list'], 'agents') if a.get('agent') == provider][:HERDR_CANDIDATES]
        for agent in agents:
            d = distance(agent['pane_id'])
            if d is not None:
                proven.setdefault(d, []).append(agent)
        if not proven or len(proven[min(proven)]) != 1:
            return None
        pane = proven[min(proven)][0]
    tab = herdr_label(result(['tab', 'get', pane['tab_id']], 'tab')['label'], secret)
    workspace = herdr_label(result(['workspace', 'get', pane['workspace_id']], 'workspace')['label'], secret)
    # One label without the other would be ambiguous.
    return {'workspace': workspace, 'tab': tab} if tab and workspace else None


def owner_app_link(route):
    """The owner-selected link, or None without the file. Never agent input."""
    path = physical(route).parent/APP_LINK_FILE
    if not os.path.lexists(path):
        return None
    value = decode(read_file(path, 4096, 0o600))
    link = value.get('app_link')
    if set(value) != {'app_link'} or not isinstance(link, str) or len(link) > 200 or not APP_LINK.fullmatch(link):
        raise Failure('invalid_app_link_owner_file')
    return link


def compact_text(value, limit):
    if not isinstance(value, str) or not value or value.strip() != value or len(value) > limit:
        raise Failure('invalid_compact_text')
    if any(unicodedata.category(c) in {'Cc', 'Cf', 'Cs', 'Zl', 'Zp'} for c in value) or any(c in value for c in '<>') or any(s in value.lower() for s in ['@channel', '@here', '@everyone']):
        raise Failure('unsafe_compact_text')


def compact_presentation(task, event):
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
    if event['phase'] not in {'before_start', 'after_start'} or event['phase'] == 'before_start' and event['event_type'] != 'agent_stopped':
        raise Failure('invalid_event_phase')
    if event['event_type'] == 'feedback_required' and (public['next_actor'] != 'user' or '?' not in public['summary']):
        raise Failure('necessary_question_required')
    try:
        instant = datetime.strptime(event['occurred_at'], UTC_FORMAT).replace(tzinfo=timezone.utc)
        if instant.strftime(UTC_FORMAT) != event['occurred_at']:
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
    if not all([args.activity, args.run, args.event_id]):
        raise Failure('short_form_fields_required')
    known = {'activity': args.activity, 'run': args.run, 'route': args.route, 'event_root': str(events_root(args.route))}
    def bind_and_deliver():
        task, event = short_binding(args)
        return deliver_bound(args, task, event)
    return preserving_attempt(attempt_markers(known, {'event_id': args.event_id})[1], bind_and_deliver)


def attempt_markers(task, event):
    key = hashlib.sha256(encode({'activity':task['activity'], 'run':task['run'], 'event':event.get('event_id')})).hexdigest()
    return key, [physical(task['event_root'])/key/'apply-started.json', events_root(task['route'])/key/'apply-started.json']


def deliver_bound(args, task, supplied_event):
    return preserving_attempt(attempt_markers(task, supplied_event)[1],
                              lambda: deliver_event(args, task, supplied_event))


def preserving_attempt(markers, deliver_call):
    try:
        return deliver_call()
    except (Failure, OSError, ValueError, KeyError, TypeError):
        # Invalid input cannot erase a known attempt on this event.
        if any(os.path.lexists(p) for p in markers):
            return outcome('unknown', 'invalid_input_preserved_attempt')
        raise


def deliver_event(args, task, supplied_event):
    expected = {'kind': 'PlyAgentNotificationEvent@2', 'activity': task['activity'], 'run': task['run']}
    if any(key in supplied_event and supplied_event[key] != want for key, want in expected.items()):
        raise Failure('event_identity_mismatch')
    event = {**expected, **supplied_event}
    presentation, instant = compact_presentation(task, event)
    root = physical(task['event_root'])
    key, markers = attempt_markers(task, event)
    folder = root/key
    route = decode(read_file(task['route'], 16 << 10))
    if route.get('webhook_env') != ENV:
        raise Failure('named_credential_required')
    source = {k: task[k] for k in ['activity', 'run', 'worktree']}
    source['kind'] = 'external'
    schema = 4
    request = {'kind':'ply.workflow.notification-request', 'schema_version':schema,
               'route':route['name'], 'source':source,
               'event':{'id':event['event_id'], 'type':event['event_type'], 'phase':event['phase'],
                        'occurred_at':event['occurred_at']},
               'sender':{'actor_claim':task['actor_claim']}, 'public':event['public'],
               'presentation':presentation}
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
    origin_lookup = 'none'
    # Only a fresh event looks up its origin; a used one keeps its first bytes.
    if not os.path.lexists(request_path) and not any(os.path.lexists(p) for p in markers):
        herdr = herdr_origin(task['provider'], secret or os.environ.get(ENV, ''))
        if herdr:
            schema = request['schema_version'] = 5
            presentation['herdr'] = herdr
            origin_lookup = 'resolved'
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
            attempt_path, shared_path = markers
            if os.path.lexists(request_path):
                prior = decode(read_file(request_path, mode=0o600))
                stored = prior.get('presentation')
                if isinstance(stored, dict) and prior.get('schema_version') in (4, 5, 6):
                    # The first attempt's origin and link win over this call.
                    schema = request['schema_version'] = prior['schema_version']
                    for name in ['herdr', 'app_link']:
                        presentation.pop(name, None)
                        if name in stored:
                            presentation[name] = stored[name]
                    origin_lookup = 'resolved' if 'herdr' in stored else 'none'
                if prior.get('schema_version') != schema or not isinstance(stored, dict):
                    raise Failure('existing_bytes_conflict')
                frozen = stored.get('local_occurred_at')
                try:
                    local = datetime.fromisoformat(frozen.replace('Z', '+00:00'))
                    # The first call froze the instant; a replay supplies none of its own.
                    instant = local
                    request['event']['occurred_at'] = local.astimezone(timezone.utc).strftime(UTC_FORMAT)
                    if local.utcoffset() is None or local != instant or local.isoformat(timespec='seconds').replace('+00:00', 'Z') != frozen:
                        raise ValueError()
                except (AttributeError, TypeError, ValueError):
                    raise Failure('invalid_preserved_local_time') from None
                presentation['local_occurred_at'] = frozen
            elif 'local_occurred_at' not in presentation:
                presentation['local_occurred_at'] = local_stamp(instant, presentation['timezone'])
            if not os.path.lexists(request_path) and not any(os.path.lexists(p) for p in markers):
                # Only a new request reads the owner's choice; its bytes then freeze it.
                link = owner_app_link(task['route'])
                if link:
                    schema = request['schema_version'] = 6
                    presentation['app_link'] = link
            if not safe_public(request, secret or os.environ.get(ENV, '')):
                raise Failure('unsafe_notification_text')
            # Freeze route locator as well as request bytes, before any send.
            write_once(folder/'binding.json', encode({'route':task['route'], 'state_root':route['state_root'], 'route_sha256':locator(task['route'])['sha256']}))
            write_once(request_path, encode(request))
            details = {'event_id':event['event_id'], 'request':str(request_path), 'origin_lookup':origin_lookup}
            marker = next((p for p in markers if os.path.lexists(p)), None)
            if marker:
                prior = decode(read_file(marker, mode=0o600))
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
                    started = encode({'notification_id':notification_id, 'confirmation':confirmation})
                    # Both entrances to one event see this marker, whichever folder they own.
                    for path in [shared_path.parent.parent, shared_path.parent]:
                        with directory(path, create=True, secure=True):
                            pass
                    write_once(shared_path, started)
                    write_once(attempt_path, started)
                    # A helper death after this marker never permits a second apply,
                    # even if the original child has not yet reserved Ply state.
                    apply_code, applied = call(args.ply_bin, command+['--apply', '--confirm', confirmation], secret)
                    show_code, state = show(args.ply_bin, task['route'], notification_id)
                    if applied.get('kind') != 'ply.workflow.notification' or applied.get('notification_id') != notification_id or applied.get('state') != state:
                        state = 'unknown'
                    result = outcome(state, 'ply_readback', notification_id=notification_id, preview_exit=code,
                                     apply_exit=apply_code, show_exit=show_code, **details)
            return result
        except (OSError, subprocess.SubprocessError, Failure) as err:
            # Never claim not-attempted after the helper's durable apply boundary.
            if any(os.path.lexists(p) for p in markers):
                reason = 'binding_conflict_preserved_attempt' if isinstance(err, Failure) and err.reason == 'existing_bytes_conflict' else 'attempt_preserved_readback_unavailable'
                return outcome('unknown', reason)
            if isinstance(err, Failure):
                raise
            raise Failure('local_io_or_cli_unavailable') from None
        finally:
            os.close(lock)


def parser():
    # No prefix matching: the removed --event must never reach --event-id.
    parser = argparse.ArgumentParser(description=__doc__, allow_abbrev=False)
    parser.add_argument('--activity', help='Stable job identity')
    parser.add_argument('--run', help='Stable run identity; part of the event identity, never changed on continuation')
    parser.add_argument('--event-id', help='Stable ID of this logical event; a replay reads, never resends')
    parser.add_argument('--status', help='done, ready_for_review, ready_for_your_check, stopped or needs_answer')
    parser.add_argument('--title', help='Public title, at most 60 characters')
    parser.add_argument('--summary', help='Public summary, at most 100 characters')
    parser.add_argument('--next-action', help='Public next action, at most 140 characters')
    parser.add_argument('--next-actor', default='user', help='user or coordinator')
    parser.add_argument('--provider', help='Actual agent runtime: codex or claude')
    parser.add_argument('--phase', default='after_start', help='after_start, or before_start for a stop')
    parser.add_argument('--route', default=str(CONFIG/'route.json'), help='Agreed route file; a used event stays bound to it')
    parser.add_argument('--context-root', help='Physical context root; defaults to the parent of the working directory')
    parser.add_argument('--timezone', default='Europe/Oslo', help='Recipient IANA time zone')
    parser.add_argument('--ply-bin', default='ply', help='Ply executable (one argv entry)')
    parser.add_argument('--credential-file', default=str(CONFIG/'webhook'), help='Exact private credential file, used only when the environment variable is absent')
    parser.add_argument('--trigger', choices=['event', 'optional_feedback', 'internal_correction', 'waiting_for_feedback'], default='event')
    parser.add_argument('--port-revision', type=int, help='Caller annotation only; never changes event identity or permits resend')
    return parser


def main():
    args = parser().parse_args()
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
