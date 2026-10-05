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


class CompactHelperRegression(HelperRegression):
    def compact(self):
        args=self.fixture()
        task=json.loads(Path(args.task).read_bytes())
        task.update(kind='PlyAgentNotificationTask@2',provider='codex',origin_cwd=str(Path.cwd().resolve()),context_root=str(Path.cwd().resolve().parent),timezone='Europe/Oslo')
        Path(args.task).write_bytes(encode(task))
        event=json.loads(Path(args.event).read_bytes());event['kind']='PlyAgentNotificationEvent@2';event['public']['status']='stopped'
        Path(args.event).write_bytes(encode(event))
        return args,task,event

    def test_v2_task_freezes_local_time_and_retains_target(self):
        args,task,event=self.compact()
        with patch.dict(os.environ,{'PLY_SLACK_WEBHOOK_URL':'invalid'}),patch.object(notify,'call') as cli:
            result=notify.deliver(args)
        cli.assert_not_called()
        request=json.loads(Path(result['request']).read_bytes())
        self.assertEqual(request['schema_version'],3)
        self.assertEqual(request['source']['worktree'],task['worktree'])
        self.assertEqual(request['presentation'],{'provider':'codex','origin_cwd':task['origin_cwd'],'context_root':task['context_root'],'context':Path(task['origin_cwd']).name,'timezone':'Europe/Oslo','local_occurred_at':'2026-10-04T16:00:00+02:00'})
        self.assertNotEqual(request['presentation']['origin_cwd'],request['source']['worktree'])

    def test_wrong_actual_cwd_fails_without_cli_or_event_reservation(self):
        args,task,event=self.compact();task['origin_cwd']=task['worktree'];task['context_root']=task['worktree'];Path(args.task).write_bytes(encode(task))
        with patch.object(notify,'call') as cli,self.assertRaises(notify.Failure): notify.deliver(args)
        cli.assert_not_called();self.assertFalse(Path(task['event_root']).exists())

    def test_bad_new_pairs_and_fields_fail_before_reservation(self):
        for variant in ['mixed','unknown-event','unknown-public','bad-status','unknown-zone','unknown-provider','surrogate','mention','long','line-separator']:
            with self.subTest(variant=variant):
                args,task,event=self.compact()
                if variant=='mixed':event['kind']='PlyAgentNotificationEvent@1'
                if variant=='unknown-event':event['extra']=True
                if variant=='unknown-public':event['public']['extra']=True
                if variant=='bad-status':event['public']['status']='done'
                if variant=='unknown-zone':task['timezone']='No/Such_Zone'
                if variant=='unknown-provider':task['provider']='model-name'
                if variant=='surrogate':event['public']['summary']='bad'+chr(0xd800)
                if variant=='mention':event['public']['summary']='Hi @channel'
                if variant=='long':event['public']['next_action']='ø'*141
                if variant=='line-separator':event['public']['summary']='one\u2028two'
                Path(args.task).write_bytes(encode(task));Path(args.event).write_text(json.dumps(event))
                with patch.object(notify,'call') as cli,self.assertRaises(notify.Failure):notify.deliver(args)
                cli.assert_not_called();self.assertFalse(Path(task['event_root']).exists())

    def test_damaged_preserved_request_retains_attempt_uncertainty(self):
        args,task,event=self.compact()
        with patch.dict(os.environ,{'PLY_SLACK_WEBHOOK_URL':'invalid'}):
            first=notify.deliver(args)
        request=Path(first['request'])
        (request.parent/'apply-started.json').write_bytes(encode({'notification_id':'ntf_'+'a'*64,'confirmation':'sha256:'+'b'*64}))
        request.write_bytes(b'{')
        with patch.object(notify,'call') as cli:
            result=notify.deliver(args)
        self.assertEqual(result['state'],'unknown')
        self.assertEqual(result['reason'],'attempt_preserved_readback_unavailable')
        cli.assert_not_called()

    def test_busy_event_never_reads_a_partially_written_request(self):
        args,task,event=self.compact()
        with patch.dict(os.environ,{'PLY_SLACK_WEBHOOK_URL':'invalid'}):first=notify.deliver(args)
        request=Path(first['request']);request.write_bytes(b'{')
        with patch.object(notify.fcntl,'flock',side_effect=BlockingIOError),patch.object(notify,'call') as cli:
            result=notify.deliver(args)
        self.assertEqual(result['state'],'unknown');self.assertEqual(result['reason'],'event_busy_preserve_local_return')
        cli.assert_not_called()


class UpgradeRegression(unittest.TestCase):
    def fixture(self):
        import install
        root=Path(tempfile.mkdtemp(prefix='notify-upgrade-')).resolve()
        old=root/'old';new=root/'new';dest=root/'installed';backup=root/'backup'
        for folder in [old,new]:
            folder.mkdir(mode=0o700);(folder/'scripts').mkdir(mode=0o700);(folder/'references').mkdir(mode=0o700)
            for name in install.FILES:(folder/name).write_bytes(('old '+name).encode())
        (new/'SKILL.md').write_bytes(b'new skill')
        install.install(old,dest)
        return install,old,new,dest,backup

    def test_exact_upgrade_preserves_private_backup(self):
        install,old,new,dest,backup=self.fixture()
        install.install(new,dest,upgrade_from=old,backup=backup)
        for name,mode in install.FILES.items():
            self.assertEqual((dest/name).read_bytes(),(new/name).read_bytes())
            self.assertEqual((backup/name).read_bytes(),(old/name).read_bytes())
            self.assertEqual((backup/name).stat().st_mode&0o777,mode)
            self.assertEqual((dest/name).stat().st_mode&0o777,mode)
        install.install(new,dest)  # Normal repeat remains idempotent.

    def test_customization_or_mode_conflict_preserves_all_bytes(self):
        for change in ['bytes','mode','extra']:
            with self.subTest(change=change):
                install,old,new,dest,backup=self.fixture()
                if change=='bytes':(dest/'SKILL.md').write_bytes(b'custom')
                if change=='mode':(dest/'scripts/notify.py').chmod(0o755)
                if change=='extra':(dest/'custom.txt').write_bytes(b'custom')
                before={str(p.relative_to(dest)):(p.read_bytes(),p.stat().st_mode) for p in dest.rglob('*') if p.is_file()}
                with self.assertRaises(notify.Failure):install.install(new,dest,upgrade_from=old,backup=backup)
                after={str(p.relative_to(dest)):(p.read_bytes(),p.stat().st_mode) for p in dest.rglob('*') if p.is_file()}
                self.assertEqual(before,after);self.assertFalse(backup.exists())


if __name__=='__main__':unittest.main()
