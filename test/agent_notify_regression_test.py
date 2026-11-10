#!/usr/bin/env python3
"""Focused regressions found in local review; no network or installed skill state."""
import errno
import json
import os
from pathlib import Path
import shutil
import sys
import tempfile
import time
import unittest
from unittest.mock import patch
sys.dont_write_bytecode=True
sys.path.insert(0,str(Path(__file__).resolve().parents[1]/'skills/ply-agent-notify/scripts'))
import notify
from common import encode


class ShortHelperRegression(unittest.TestCase):
    ID='ntf_'+'a'*64
    CONFIRM='sha256:'+'b'*64

    def setUp(self):
        self.root=Path(tempfile.mkdtemp(prefix='agent-notify-short-')).resolve()
        self.state=self.root/'state'
        self.cwd=self.root/'job';self.cwd.mkdir()
        self.route=self.root/'route.json'
        self.route.write_bytes(encode({'name':'ply-log','webhook_env':'PLY_SLACK_WEBHOOK_URL','state_root':str(self.state)}))
        previous=Path.cwd();os.chdir(self.cwd);self.addCleanup(os.chdir,previous)
        env=patch.dict(os.environ,{'PLY_SLACK_WEBHOOK_URL':'https://hooks.slack.com/services/T0/B0/fixture'});env.start();self.addCleanup(env.stop)
        # A real Herdr session must never reach these fake-only tests.
        os.environ.pop('HERDR_ENV',None)
        self.sends=[];self.state_value='transport_acknowledged'

    def args(self,**change):
        values={'activity':'job-1','run':'job-1','event-id':'finished','status':'done','title':'ℹ️ Jobb','summary':'Ferdig.','next-action':'No action required.','provider':'claude',**change}
        argv=['--route',str(self.route),'--ply-bin','stub']
        for name,value in values.items():
            if value is not None:argv+=['--'+name,value]
        return notify.parser().parse_args(argv)

    def fake(self,ply,args,secret=None):
        if args[0]=='show':return 0,{'kind':'ply.workflow.notification','notification_id':self.ID,'state':self.state_value}
        if '--check' in args:
            self.requests=getattr(self,'requests',[])+[json.loads(Path(args[2]).read_bytes())]
            return 0,{'kind':'ply.workflow.notification-preview','allowed':True,'notification_id':self.ID,'confirmation':self.CONFIRM,'existing_state':None}
        self.sends.append(args)
        return 5 if self.state_value=='unknown' else 0,{'kind':'ply.workflow.notification','notification_id':self.ID,'state':self.state_value}

    def deliver(self,args=None,**change):
        with patch.object(notify,'call',self.fake):
            return notify.deliver(args or self.args(**change))

    def stored(self):
        return sorted(p.name for p in (self.root/'notification-events').rglob('*') if p.is_file())

    def test_normal_finish_needs_no_documents_and_stores_no_copies(self):
        result=self.deliver()
        self.assertEqual(result['state'],'transport_acknowledged');self.assertEqual(len(self.sends),1)
        request=self.requests[0]
        self.assertEqual(request['schema_version'],4)
        self.assertEqual(request['source'],{'kind':'external','activity':'job-1','run':'job-1','worktree':str(self.cwd)})
        self.assertEqual(set(request['event']),{'id','type','phase','occurred_at'})
        self.assertEqual((request['event']['id'],request['event']['type'],request['event']['phase']),('finished','agent_finished','after_start'))
        self.assertEqual(request['public'],{'status':'done','task_title':'ℹ️ Jobb','summary':'Ferdig.','next_action':'No action required.','next_actor':'user'})
        self.assertEqual(request['presentation']['context'],'job');self.assertEqual(request['presentation']['timezone'],'Europe/Oslo')
        self.assertEqual(self.stored(),['.lock','apply-started.json','binding.json','request.json'])

    def test_status_selects_event_type(self):
        for status,summary,kind in [('stopped','Stoppet.','agent_stopped'),('needs_answer','Hvilket mål?','feedback_required')]:
            self.deliver(status=status,summary=summary,**{'event-id':status})
            self.assertEqual(self.requests[-1]['event']['type'],kind)

    def test_same_id_after_lost_response_only_reads_and_keeps_time(self):
        self.state_value='unknown'
        first=self.deliver();self.assertEqual(first['state'],'unknown')
        before=Path(first['request']).read_bytes()
        with patch.object(notify,'datetime',wraps=notify.datetime) as clock:
            clock.now.return_value=notify.datetime(2030,1,1,tzinfo=notify.timezone.utc)
            second=self.deliver()
        self.assertEqual(second['state'],'unknown');self.assertEqual(second['reason'],'preserved_attempt_readback')
        self.assertEqual(len(self.sends),1);self.assertEqual(Path(first['request']).read_bytes(),before)

    def test_changed_text_on_used_id_is_conflict_without_send(self):
        self.deliver()
        result=self.deliver(summary='Endret.')
        self.assertEqual(result['state'],'unknown');self.assertEqual(result['reason'],'binding_conflict_preserved_attempt')
        self.assertEqual(len(self.sends),1)

    def test_unsent_changed_text_stays_not_attempted(self):
        with patch.dict(os.environ,{'PLY_SLACK_WEBHOOK_URL':'invalid'}):
            self.deliver()
            with self.assertRaises(notify.Failure) as err:self.deliver(summary='Endret.')
        self.assertEqual(err.exception.reason,'existing_bytes_conflict');self.assertEqual(self.sends,[])

    def test_route_change_on_used_id_is_conflict(self):
        self.deliver()
        moved=self.root/'moved.json'
        moved.write_bytes(encode({'name':'ply-log','webhook_env':'PLY_SLACK_WEBHOOK_URL','state_root':str(self.state)+'-other'}))
        args=self.args();args.route=str(moved)
        result=self.deliver(args)
        self.assertEqual(result['state'],'unknown');self.assertEqual(result['reason'],'binding_conflict_preserved_attempt')
        self.assertEqual(len(self.sends),1)

    def test_mixed_or_incomplete_input_fails_before_any_state(self):
        for change in [{'status':'bogus'},{'activity':None},{'run':None},{'event-id':None},{'title':None},{'provider':None},{'provider':'gpt'},{'summary':'Hi @channel'},
                       {'timezone':'No/Such_Zone'},{'summary':'bad'+chr(0xd800)},{'next-action':'ø'*141},{'summary':'one\u2028two'}]:
            with self.subTest(change=change),self.assertRaises(notify.Failure):self.deliver(**change)
        self.assertFalse((self.root/'notification-events').exists());self.assertEqual(self.sends,[])

    def test_removed_bound_task_documents_are_not_arguments(self):
        # The bound Task entrance and its source documents belonged to the removed delivery process.
        for flag in ['--task','--event','--report','--start-receipt']:
            with self.subTest(flag=flag),patch('sys.stderr'),self.assertRaises(SystemExit):
                notify.parser().parse_args(['--route',str(self.route),flag,str(self.root/'x.json')])

    def test_concurrent_lock_creation_uncertainty_is_not_not_attempted(self):
        original=os.open
        def open_with_race(path,flags,*a,**kw):
            if path=='.lock':raise FileNotFoundError(errno.ENOENT,'Injected concurrent create boundary')
            return original(path,flags,*a,**kw)
        with patch.object(notify.os,'open',open_with_race):
            result=self.deliver()
        self.assertEqual(result['state'],'unknown');self.assertEqual(self.sends,[])

    def test_damaged_preserved_request_retains_attempt_uncertainty(self):
        with patch.dict(os.environ,{'PLY_SLACK_WEBHOOK_URL':'invalid'}):
            first=self.deliver()
        request=Path(first['request'])
        (request.parent/'apply-started.json').write_bytes(encode({'notification_id':self.ID,'confirmation':self.CONFIRM}))
        request.write_bytes(b'{')
        with patch.object(notify,'call') as cli:
            result=notify.deliver(self.args())
        self.assertEqual(result['state'],'unknown')
        self.assertEqual(result['reason'],'attempt_preserved_readback_unavailable')
        cli.assert_not_called()

    def test_busy_event_never_reads_a_partially_written_request(self):
        with patch.dict(os.environ,{'PLY_SLACK_WEBHOOK_URL':'invalid'}):first=self.deliver()
        request=Path(first['request']);request.write_bytes(b'{')
        with patch.object(notify.fcntl,'flock',side_effect=BlockingIOError),patch.object(notify,'call') as cli:
            result=notify.deliver(self.args())
        self.assertEqual(result['state'],'unknown');self.assertEqual(result['reason'],'event_busy_preserve_local_return')
        cli.assert_not_called()

    def test_marker_of_an_earlier_entrance_blocks_short_send(self):
        # A removed bound-Task entrance left its shared marker beside the route; it stays history.
        key=notify.attempt_markers({'activity':'job-1','run':'job-1','route':str(self.route),'event_root':str(self.root/'notification-events')},{'event_id':'finished'})[0]
        folder=self.root/'notification-events'/key
        for path in [folder.parent,folder]:path.mkdir(mode=0o700,exist_ok=True)
        marker=folder/'apply-started.json';marker.write_bytes(encode({'notification_id':self.ID,'confirmation':self.CONFIRM}));marker.chmod(0o600)
        self.state_value='unknown'
        result=self.deliver()
        self.assertEqual(result['state'],'unknown');self.assertEqual(result['reason'],'preserved_attempt_readback')
        self.assertEqual(self.sends,[])

    def test_invalid_replay_of_used_event_is_never_not_attempted(self):
        outside=str(self.root.parent)
        other=tempfile.mkdtemp(prefix='agent-notify-outside-',dir=outside)
        self.deliver()
        for change in [{'summary':'bad <text>'},{'status':'needs_answer'},{'context-root':other},{'status':'bogus'},{'next-actor':'nobody'},{'phase':'later'},{'title':'x'*61},{'summary':None},{'title':None}]:
            with self.subTest(change=change):
                result=self.deliver(**change)
                self.assertEqual(result['state'],'unknown');self.assertEqual(result['reason'],'invalid_input_preserved_attempt')
                self.assertEqual(len(self.sends),1)

    def test_short_field_errors_on_used_event_stay_unknown(self):
        self.deliver()
        args=self.args();args.summary=''
        result=self.deliver(args)
        self.assertEqual(result['state'],'unknown');self.assertEqual(result['reason'],'invalid_input_preserved_attempt')
        self.assertEqual(len(self.sends),1)

    def test_same_argv_from_deleted_cwd_on_used_event_stays_unknown(self):
        self.deliver()
        self.cwd.rmdir()
        result=self.deliver()
        self.assertEqual(result['state'],'unknown');self.assertEqual(result['reason'],'invalid_input_preserved_attempt')
        self.assertEqual(len(self.sends),1)

    def test_private_directory_name_is_not_public_text_in_short_form(self):
        for name in ['a<b>','x'*97]:
            with self.subTest(name=name):
                cwd=self.cwd/name;cwd.mkdir();os.chdir(cwd)
                result=self.deliver(**{'event-id':name[:5]})
                self.assertEqual(result['state'],'transport_acknowledged')
                self.assertEqual(self.requests[-1]['presentation']['context'],name)

    def test_outside_context_root_still_fails_without_state(self):
        other=tempfile.mkdtemp(prefix='agent-notify-outside-',dir=str(self.root))
        with self.assertRaises(notify.Failure) as err:self.deliver(**{'context-root':other})
        self.assertEqual(err.exception.reason,'origin_outside_context_root');self.assertEqual(self.sends,[])

    def test_usage_names_only_the_short_form(self):
        skill=Path(notify.__file__).parent.parent
        usage=(skill/'references/usage.md').read_text();short=(skill/'SKILL.md').read_text()
        for flag in ['--activity','--run','--event-id','--status','--provider']:self.assertIn(flag,short)
        for flag in ['--task','--event ','--report','--start-receipt']:self.assertNotIn(flag,usage+short)
        for status,kind in notify.EVENT_TYPE.items():self.assertRegex(usage,r'`%s`[^\n]*`%s`'%(kind,status))


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

    def test_upgrade_from_the_v4_source_installs_a_helper_that_starts(self):
        import install,io,subprocess,tarfile
        root=Path(tempfile.mkdtemp(prefix='notify-v4-upgrade-')).resolve();self.addCleanup(shutil.rmtree,root)
        repo=Path(__file__).resolve().parents[1]
        archive=subprocess.run(['git','archive','8a400fe','skills/ply-agent-notify'],cwd=repo,capture_output=True,check=True).stdout
        with tarfile.open(fileobj=io.BytesIO(archive)) as t:t.extractall(root,filter='data')
        old,dest=root/'skills/ply-agent-notify',root/'installed'
        install.install(old,dest)
        install.install(repo/'skills/ply-agent-notify',dest,upgrade_from=old,backup=root/'backup')
        started=subprocess.run([sys.executable,dest/'scripts/notify.py','--help'],capture_output=True)
        self.assertEqual(started.returncode,0,started.stderr)

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


class AppLinkOwnerFileRegression(unittest.TestCase):
    """The optional owner file beside the route selects request-v6 for new events only."""
    ID,CONFIRM=ShortHelperRegression.ID,ShortHelperRegression.CONFIRM
    args,fake,deliver=ShortHelperRegression.args,ShortHelperRegression.fake,ShortHelperRegression.deliver
    LINK='https://plyapp-pwa.web.app'

    def setUp(self):
        ShortHelperRegression.setUp(self)
        self.owner=self.root/'app-link.json'

    def own(self,value,mode=0o600):
        self.owner.write_bytes(encode(value) if isinstance(value,dict) else value);self.owner.chmod(mode)

    def not_read(self):
        return patch.object(notify,'owner_app_link',side_effect=AssertionError('owner file read'))

    def test_missing_file_keeps_v4_and_valid_file_selects_v6(self):
        self.deliver(**{'event-id':'without'})
        self.assertEqual(self.requests[-1]['schema_version'],4);self.assertNotIn('app_link',self.requests[-1]['presentation'])
        self.own({'app_link':self.LINK})
        self.deliver(**{'event-id':'with'})
        self.assertEqual(self.requests[-1]['schema_version'],6);self.assertEqual(self.requests[-1]['presentation']['app_link'],self.LINK)
        self.assertEqual(len(self.sends),2)

    def test_v6_replay_keeps_the_first_link_without_reading_the_owner_file(self):
        self.own({'app_link':self.LINK})
        first=self.deliver();before=Path(first['request']).read_bytes()
        for change in [{'app_link':'https://other.example'},b'not json',None]:
            if change is None:self.owner.unlink()
            else:self.own(change)
            with self.subTest(change=change),self.not_read():
                again=self.deliver()
                self.assertEqual(again['reason'],'preserved_attempt_readback')
                self.assertEqual(Path(first['request']).read_bytes(),before)
        self.assertEqual(json.loads(before)['presentation']['app_link'],self.LINK);self.assertEqual(len(self.sends),1)

    def test_enabling_the_file_never_recasts_an_existing_v4_event(self):
        first=self.deliver();before=Path(first['request']).read_bytes()
        self.own({'app_link':self.LINK})
        with self.not_read():again=self.deliver()
        self.assertEqual(again['reason'],'preserved_attempt_readback');self.assertEqual(len(self.sends),1)
        self.assertEqual(Path(first['request']).read_bytes(),before);self.assertEqual(json.loads(before)['schema_version'],4)

    def test_attempt_marker_without_request_reads_back_without_the_owner_file(self):
        self.own({'app_link':self.LINK})
        key,markers=notify.attempt_markers({'activity':'job-1','run':'job-1','event_root':str(self.root/'notification-events'),'route':str(self.route)},{'event_id':'finished'})
        markers[0].parent.mkdir(parents=True,mode=0o700);markers[0].parent.parent.chmod(0o700)
        markers[0].write_bytes(encode({'notification_id':self.ID}));markers[0].chmod(0o600)
        with self.not_read():result=self.deliver()
        self.assertEqual(result['reason'],'preserved_attempt_readback');self.assertEqual(self.sends,[])

    def test_invalid_owner_file_fails_before_any_request(self):
        for name,value,mode in [('extra key',{'app_link':self.LINK,'label':'x'},0o600),('http',{'app_link':'http://plyapp-pwa.web.app'},0o600),
                                ('markup',{'app_link':'https://plyapp-pwa.web.app/|x'},0o600),('not string',{'app_link':[self.LINK]},0o600),
                                ('query',{'app_link':self.LINK+'?x=1'},0o600),('not json',b'{',0o600),('readable',{'app_link':self.LINK},0o644)]:
            with self.subTest(name):
                self.own(value,mode)
                with self.assertRaises(notify.Failure):self.deliver(**{'event-id':name.replace(' ','-')})
                self.assertFalse(list((self.root/'notification-events').rglob('request.json')))
        self.assertEqual(self.sends,[])


class HerdrOriginRegression(unittest.TestCase):
    """A fake herdr executable; the proof uses this test process's real ancestry."""
    ID,CONFIRM=ShortHelperRegression.ID,ShortHelperRegression.CONFIRM
    args,fake,deliver=ShortHelperRegression.args,ShortHelperRegression.fake,ShortHelperRegression.deliver

    def setUp(self):
        ShortHelperRegression.setUp(self)
        self.herdr=self.root/'herdr'
        self.herdr.write_text('#!'+sys.executable+' -I\n'
            'import json,sys,time\nfrom pathlib import Path\n'
            'here=Path(__file__).parent;cfg=json.loads((here/"herdr.json").read_text())\n'
            'key=" ".join(sys.argv[1:])\n'
            'with open(here/"herdr-calls.log","a") as log:log.write(key+"\\n")\n'
            'time.sleep(cfg.get("sleep",0))\n'
            'reply=cfg["replies"].get(key,{"error":{"code":"not_found"}})\n'
            'print(json.dumps(reply));sys.exit(1 if "error" in reply else 0)\n')
        self.herdr.chmod(0o700)
        os.environ.update(HERDR_ENV='1',HERDR_BIN_PATH=str(self.herdr))
        self.parent=os.getppid()
        self.configure()

    def configure(self,env_pid=None,agents=(),tab='pn01-herdr-notification-spec',workspace='ply-master',stale=False,sleep=0):
        pane={'pane_id':'w9:pB','tab_id':'w9:tB','workspace_id':'w9','agent':'claude'}
        replies={'pane current --current':{'error':{'code':'pane_not_found'}} if stale else {'result':{'pane':pane}},
                 'pane process-info --pane w9:pB':{'result':{'process_info':{'pane_id':'w9:pB','shell_pid':self.parent if env_pid is None else env_pid}}},
                 'agent list':{'result':{'agents':[{'pane_id':p,'tab_id':p.replace(':p',':t'),'workspace_id':'w9','agent':a} for p,a,_ in agents]}},
                 'tab get w9:tB':{'result':{'tab':{'label':tab}}},'workspace get w9':{'result':{'workspace':{'label':workspace}}}}
        for p,a,pid in agents:
            replies['pane process-info --pane '+p]={'result':{'process_info':{'pane_id':p,'shell_pid':pid}}}
            replies['tab get '+p.replace(':p',':t')]={'result':{'tab':{'label':'moved-'+p[-1]}}}
        (self.root/'herdr.json').write_text(json.dumps({'replies':replies,'sleep':sleep}))

    def calls(self):
        log=self.root/'herdr-calls.log'
        return log.read_text().splitlines() if log.exists() else []

    def test_missing_herdr_env_or_binary_sends_v4_without_lookup(self):
        for name,change in [('env',{'HERDR_ENV':'0'}),('binary',{'HERDR_BIN_PATH':str(self.root/'absent')})]:
            with self.subTest(name),patch.dict(os.environ,change):
                result=self.deliver(**{'event-id':name})
                self.assertEqual(result['origin_lookup'],'none');self.assertEqual(self.requests[-1]['schema_version'],4)
        self.assertEqual(self.calls(),[])

    def test_proven_env_pane_names_its_workspace_and_tab(self):
        result=self.deliver()
        self.assertEqual(result['origin_lookup'],'resolved');self.assertEqual(len(self.sends),1)
        request=self.requests[0]
        self.assertEqual(request['schema_version'],5)
        self.assertEqual(request['presentation']['herdr'],{'workspace':'ply-master','tab':'pn01-herdr-notification-spec'})
        self.assertNotIn('agent list',self.calls())

    def test_owner_link_keeps_the_proven_origin_in_v6(self):
        (self.root/'app-link.json').write_bytes(encode({'app_link':AppLinkOwnerFileRegression.LINK}));(self.root/'app-link.json').chmod(0o600)
        result=self.deliver()
        self.assertEqual(result['origin_lookup'],'resolved');self.assertEqual(self.requests[0]['schema_version'],6)
        self.assertEqual(self.requests[0]['presentation']['herdr'],{'workspace':'ply-master','tab':'pn01-herdr-notification-spec'})
        self.assertEqual(self.requests[0]['presentation']['app_link'],AppLinkOwnerFileRegression.LINK)

    def test_stale_env_pane_falls_back_to_the_proven_moved_pane(self):
        self.configure(stale=True,agents=[('w9:pX','codex',self.parent),('w9:pY','claude',999999),('w9:pZ','claude',self.parent)])
        self.deliver()
        self.assertEqual(self.requests[0]['presentation']['herdr'],{'workspace':'ply-master','tab':'moved-Z'})

    def test_unproven_or_ambiguous_origin_sends_v4(self):
        for name,agents in [('unproven',[('w9:pX','codex',self.parent),('w9:pY','claude',999999)]),
                            ('ambiguous',[('w9:pY','claude',self.parent),('w9:pZ','claude',self.parent)])]:
            with self.subTest(name):
                self.configure(env_pid=999999,agents=agents)
                result=self.deliver(**{'event-id':name})
                self.assertEqual(result['origin_lookup'],'none');self.assertEqual(self.requests[-1]['schema_version'],4)
                self.assertNotIn('herdr',self.requests[-1]['presentation'])

    def test_hanging_herdr_ends_within_the_deadline_and_still_sends_once(self):
        self.configure(sleep=30)
        started=time.monotonic()
        result=self.deliver()
        self.assertLess(time.monotonic()-started,3.0)
        self.assertEqual(result['origin_lookup'],'none');self.assertEqual(self.requests[0]['schema_version'],4);self.assertEqual(len(self.sends),1)

    def test_rename_or_close_after_the_first_attempt_replays_frozen_bytes(self):
        for name,first_env in [('v5',{}),('v4',{'HERDR_ENV':'0'})]:
            with self.subTest(name):
                self.configure()
                with patch.dict(os.environ,first_env):
                    first=self.deliver(**{'event-id':name})
                before,calls=Path(first['request']).read_bytes(),len(self.calls())
                for change in [{'tab':'renamed'},{'stale':True}]:
                    self.configure(**change)
                    again=self.deliver(**{'event-id':name})
                    self.assertEqual(again['reason'],'preserved_attempt_readback')
                    self.assertEqual(again['origin_lookup'],first['origin_lookup'])
                    self.assertEqual(Path(first['request']).read_bytes(),before)
                    self.assertEqual(len(self.calls()),calls)
        self.assertEqual(len(self.sends),2)

    def test_unsafe_label_omits_both_and_long_label_is_truncated(self):
        self.configure(tab='<!here>')
        self.deliver(**{'event-id':'unsafe'})
        self.assertEqual(self.requests[-1]['schema_version'],4)
        self.configure(workspace='w'*40)
        self.deliver(**{'event-id':'long'})
        self.assertEqual(self.requests[-1]['presentation']['herdr']['workspace'],'w'*31+'…')

    def test_label_sanitizing(self):
        label=notify.herdr_label
        self.assertEqual(label('  ply \t  master ',''),'ply master')
        self.assertEqual(label('e\u0301','' ),'\u00e9')
        self.assertEqual(label('x'*31+' yz',''),'x'*31+'…')
        self.assertEqual(label('x'*30+' yz',''),'x'*30+'…')
        for bad in ['','  ','a\x07b','a\u202eb','one\u2028two','a<b','@Here now','team @channel','https://hooks.slack.com/services/T/B/x','has-secret-in-it',None,7]:
            self.assertIsNone(label(bad,'secret'),bad)


if __name__=='__main__':unittest.main()
