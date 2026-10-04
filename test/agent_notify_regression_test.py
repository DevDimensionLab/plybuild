#!/usr/bin/env python3
"""Focused regressions found in local review; no network or installed skill state."""
import argparse
import errno
import json
import os
from pathlib import Path
import sys
import tempfile
import unittest
from unittest.mock import patch
sys.dont_write_bytecode=True
sys.path.insert(0,str(Path(__file__).resolve().parents[1]/'skills/ply-agent-notify/scripts'))
import notify
from common import encode


class HelperRegression(unittest.TestCase):
    def fixture(self):
        root=Path(tempfile.mkdtemp(prefix='agent-notify-regression-')).resolve()
        def put(name,value):
            p=root/name;p.write_bytes(encode(value));p.chmod(0o600);return str(p)
        handoff=put('handoff.json',{'kind':'FixtureHandoff@1'})
        route=put('route.json',{'name':'ply-log','webhook_env':'PLY_SLACK_WEBHOOK_URL','state_root':str(root/'state')})
        event=put('event.json',{'kind':'PlyAgentNotificationEvent@1','activity':'test','run':'run-1','event_id':'event-1','event_type':'agent_stopped','phase':'before_start','occurred_at':'2026-10-04T14:00:00Z','public':{'task_title':'Task','summary':'Stopped before start.','next_action':'Review the local return.','next_actor':'coordinator'}})
        task=put('task.json',{'kind':'PlyAgentNotificationTask@1','activity':'test','run':'run-1','worktree':str(root),'handoff':handoff,'route':route,'event_root':str(root/'events'),'actor_claim':'test','notification_authorized':True})
        return argparse.Namespace(task=task,event=event,trigger='event',credential_file=str(root/'absent'),report=None,start_receipt=None,ply_bin='not-invoked')

    def test_concurrent_lock_creation_uncertainty_is_not_not_attempted(self):
        args=self.fixture();original=os.open
        def open_with_race(path,flags,*a,**kw):
            if path=='.lock':raise FileNotFoundError(errno.ENOENT,'Injected concurrent create boundary')
            return original(path,flags,*a,**kw)
        with patch.object(notify.os,'open',open_with_race), patch.object(notify,'call') as cli:
            result=notify.deliver(args)
        self.assertEqual(result['state'],'unknown')
        cli.assert_not_called()


if __name__=='__main__':unittest.main()
