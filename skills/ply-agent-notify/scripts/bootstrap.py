#!/usr/bin/env python3
"""One-time, explicitly authorized private route/credential setup; never sends."""
import argparse
import getpass
import json
import os
import sys
sys.dont_write_bytecode = True
from common import CONFIG, STATE, ENV, WEBHOOK, Failure, decode, directory, encode, physical, read_file, write_once


def bootstrap(config, state, value):
    config, state = physical(config), physical(state)
    if not WEBHOOK.fullmatch(value):
        raise Failure('invalid_credential')
    route = encode({'kind': 'ply.workflow.notification-route', 'schema_version': 1,
                    'name': 'ply-log', 'transport': 'slack_incoming_webhook',
                    'channel_label': 'ply-log', 'webhook_env': ENV, 'state_root': str(state)})
    # Inspect existing leaves before creating any new file; conflicts stay intact.
    for path, wanted in [(config/'webhook', value.encode()), (config/'route.json', route)]:
        try:
            current = read_file(path, 16384, 0o600)
        except FileNotFoundError:
            continue
        if current != wanted:
            raise Failure('existing_bytes_conflict')
    with directory(config, create=True, parents=True, secure=True):
        pass
    with directory(state, create=True, parents=True, secure=True):
        pass
    write_once(config/'webhook', value.encode())
    write_once(config/'route.json', route)
    return {'kind': 'PlyAgentNotificationBootstrap@1', 'state': 'ready',
            'route': str(config/'route.json'), 'state_root': str(state), 'sent': False}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--config-dir', default=str(CONFIG), help='Explicit private configuration directory (fixture paths supported)')
    parser.add_argument('--state-root', default=str(STATE), help='Fixed private notification state root')
    args = parser.parse_args()
    try:
        if ENV in os.environ:
            value = os.environ[ENV]
        elif os.path.lexists(physical(args.config_dir)/'webhook'):
            value = read_file(physical(args.config_dir)/'webhook', 4096, 0o600).decode('utf-8')
        else:
            # Fail rather than fall back to an echoed stdin prompt.
            if not sys.stdin.isatty():
                raise Failure('hidden_terminal_input_required')
            import warnings
            with warnings.catch_warnings():
                warnings.simplefilter('error', getpass.GetPassWarning)
                value = getpass.getpass('Existing Slack webhook (hidden): ')
        print(json.dumps(bootstrap(args.config_dir, args.state_root, value)))
        return 0
    except (Failure, OSError, ValueError, KeyError, getpass.GetPassWarning):
        print(json.dumps({'kind': 'PlyAgentNotificationBootstrap@1', 'state': 'not_attempted',
                          'reason': 'Setup rejected; check the credential, physical paths, owner, modes and existing bytes. No message was sent.'}))
        return 2


if __name__ == '__main__':
    sys.exit(main())
