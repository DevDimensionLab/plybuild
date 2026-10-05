#!/usr/bin/env python3
"""Retained offline fixtures exercising the installed skill and actual Ply CLI graph."""
import argparse
import concurrent.futures
import hashlib
import io
import json
import os
from pathlib import Path
import secrets
import stat
import subprocess
import sys
import tarfile
import tempfile
import time

REPO = Path(__file__).resolve().parents[1]
BASE = 'c8f9dd00ee36d3f57191355c278d982808d421a8'
ENV = 'PLY_SLACK_WEBHOOK_URL'


def encoded(value):
    return (json.dumps(value, sort_keys=True)+'\n').encode()


def put(path, value):
    path.write_bytes(encoded(value) if not isinstance(value, bytes) else value)
    path.chmod(0o600)
    return path


def sha(path):
    return 'sha256:'+hashlib.sha256(path.read_bytes()).hexdigest()


def loc(path, metadata=False):
    value = {'path':str(path), 'sha256':sha(path)}
    if metadata:
        value['kind'] = json.loads(path.read_bytes())['kind']
    return value


class Harness:
    def __init__(self, root):
        self.root = root
        self.secret = 'https://hooks.slack.com/services/Tfixture/Bfixture/'+secrets.token_hex(24)
        self.other = 'https://hooks.slack.com/services/Tfixture/Bfixture/'+secrets.token_hex(24)
        self.env = dict(os.environ)
        self.env.pop(ENV, None)
        self.env['PYTHONDONTWRITEBYTECODE'] = '1'
        self.env['TMPDIR'] = str(root/'tmp')
        (root/'tmp').mkdir(mode=0o700)
        (root/'commands').mkdir(mode=0o700)
        (root/'installed').mkdir(mode=0o700)
        self.commands = []
        self.checks = []
        self.credentials = set()
        self.skill = root/'installed'/'ply-agent-notify'
        self.ply = root/'ply'
        self.driver = root/'notification-driver'
        self.counter = 0

    def command(self, argv, env=None, expected=(0,), cwd=REPO):
        self.counter += 1
        name = '%03d'%self.counter
        c = subprocess.run([str(x) for x in argv], cwd=cwd, env=self.env if env is None else env, capture_output=True)
        for data in [c.stdout, c.stderr]:
            assert self.secret.encode() not in data and self.other.encode() not in data, 'credential leaked in subprocess output'
        out = self.root/'commands'/(name+'.stdout')
        err = self.root/'commands'/(name+'.stderr')
        out.write_bytes(c.stdout); err.write_bytes(c.stderr)
        self.commands.append({'argv':[str(x) for x in argv], 'cwd':str(cwd), 'exit':c.returncode,
                              'stdout':str(out.relative_to(self.root)), 'stderr':str(err.relative_to(self.root))})
        put(self.root/'commands.json', self.commands)
        assert c.returncode in expected, 'Unexpected exit %s for %s; see %s'%(c.returncode, argv[0], out)
        return c

    def secret_env(self, value=None):
        return {**self.env, ENV:self.secret if value is None else value}

    def setup(self):
        self.command(['go','build','-o',self.ply,'./cmd/ply'])
        self.command(['go','test','-c','-o',self.driver,'./test/notification-driver'])
        self.command(['make','install-agent-notify','AGENT_NOTIFY_DEST='+str(self.skill)])
        self.command(['make','install-agent-notify','AGENT_NOTIFY_DEST='+str(self.skill)])
        for p in self.skill.rglob('*'):
            if p.is_file():
                want = 0o700 if p.name in ['notify.py','bootstrap.py','install.py'] else 0o600
                assert stat.S_IMODE(p.stat().st_mode)==want
            else:
                assert stat.S_IMODE(p.stat().st_mode)==0o700
        self.checks.append('fixture installation is idempotent and private')

    def fixture(self, name, kind='agent_finished', phase='after_start', response=None):
        folder = self.root/name; folder.mkdir(mode=0o700)
        (folder/'work').mkdir(mode=0o700)
        config = folder/'config'
        self.command([sys.executable,self.skill/'scripts/bootstrap.py','--config-dir',config,'--state-root',folder/'state'], env=self.secret_env())
        self.credentials.add(config/'webhook')
        self.command([sys.executable,self.skill/'scripts/bootstrap.py','--config-dir',config,'--state-root',folder/'state'])
        assert stat.S_IMODE((config/'webhook').stat().st_mode)==0o600
        assert stat.S_IMODE(config.stat().st_mode)==0o700
        handoff = put(folder/'handoff.md', b'Offline task mandate. PRIVATE_SOURCE_SENTINEL\n')
        report = put(folder/'report.json', {'kind':'FixtureReport@1','activity':'fixture/activity','run':'run-1','outcome':'reported','private':'PRIVATE_SOURCE_SENTINEL'})
        start = put(folder/'start.json', {'kind':'FixtureStart@1','activity':'fixture/activity','run':'run-1'})
        event = {'kind':'PlyAgentNotificationEvent@1','activity':'fixture/activity','run':'run-1','event_id':name,'event_type':kind,'phase':phase,'occurred_at':'2026-10-04T14:00:00Z',
                 'public':{'task_title':'A useful agent task', 'summary':'Which agreed target should I use?' if kind=='feedback_required' else 'The local return is preserved.',
                           'next_action':'Answer in the active agent task.' if kind=='feedback_required' else 'Review the preserved local return.',
                           'next_actor':'user' if kind=='feedback_required' else 'coordinator'}}
        put(folder/'event.json',event)
        task = {'kind':'PlyAgentNotificationTask@1','activity':'fixture/activity','run':'run-1','worktree':str(folder/'work'), 'handoff':str(handoff),'route':str(config/'route.json'),
                'event_root':str(folder/'events'), 'actor_claim':'offline-agent','notification_authorized':True}
        put(folder/'task.json',task)
        cfg = {'kind':'PlyNotificationAcceptanceCase@1','artifacts_root':str(self.root),'now':'2026-10-04T14:00:00Z','post_log':str(folder/'posts.jsonl'),
               'transport':response or {'kind':'http','status':200,'body':'ok'},'fault':None}
        put(folder/'driver.json',cfg)
        wrapper = '''#!/usr/bin/env python3
import hashlib,json,os,subprocess,sys
from pathlib import Path
folder=Path(__file__).parent
cfg=json.loads((folder/'driver.json').read_bytes())
args=sys.argv[1:]
present='PLY_SLACK_WEBHOOK_URL' in os.environ
with (folder/'calls.jsonl').open('a') as f: f.write(json.dumps({'args':args,'credential_present':present})+'\\n')
if args[2]=='show' and present: sys.exit(87)
if args[2]=='send' and cfg.get('expected_credential_sha256') and hashlib.sha256(os.environ.get('PLY_SLACK_WEBHOOK_URL','').encode()).hexdigest()!=cfg['expected_credential_sha256']: sys.exit(88)
c=subprocess.run([DRIVER,'--config',str(folder/'driver.json'),'--',*args],capture_output=True)
sys.stdout.buffer.write(c.stdout);sys.stderr.buffer.write(c.stderr)
sys.exit(0 if cfg.get('force_zero_exit') else c.returncode)
'''.replace('DRIVER',repr(str(self.driver)))
        (folder/'ply-fixture').write_text(wrapper); (folder/'ply-fixture').chmod(0o700)
        return folder

    def notify_args(self, folder, extra=(), report=True, receipt=True):
        args = [sys.executable,self.skill/'scripts/notify.py','--task',folder/'task.json','--event',folder/'event.json', '--ply-bin',folder/'ply-fixture','--credential-file',folder/'config/webhook']
        if report: args += ['--report',folder/'report.json']
        if receipt: args += ['--start-receipt',folder/'start.json']
        return [str(x) for x in [*args,*extra]]

    def notify(self, folder, extra=(), expected=(0,), env=None, cwd=REPO, **kwargs):
        c = self.command(self.notify_args(folder, extra, **kwargs),env=env,expected=expected,cwd=cwd)
        return json.loads(c.stdout)

    def posts(self, folder):
        path = folder/'posts.jsonl'
        return [json.loads(line) for line in path.read_text().splitlines()] if path.exists() else []

    def journeys(self):
        for kind in ['agent_finished','feedback_required','agent_stopped']:
            before_start = kind=='agent_stopped'
            f=self.fixture(kind,kind,'before_start' if before_start else 'after_start')
            before={p:p.read_bytes() for p in [f/'event.json', f/'report.json', f/'start.json']}
            r=self.notify(f,report=not before_start,receipt=not before_start)
            assert r['state']=='transport_acknowledged' and r['product_outcome']=='unchanged'
            posts=self.posts(f);assert len(posts)==1
            text=posts[0]['payload']['text']
            assert ('Your answer is needed' if kind=='feedback_required' else 'Agent stopped' if before_start else 'Agent finished — reported, not yet controlled') in text
            assert ('User continues in the active task' if kind=='feedback_required' else 'Coordinator continues') in text
            assert 'PRIVATE_SOURCE_SENTINEL' not in text and str(f) not in text
            request=json.loads(Path(r['request']).read_bytes())
            if before_start: assert 'report' not in request['source'] and 'start_receipt' not in request['source']
            calls=[json.loads(x) for x in (f/'calls.jsonl').read_text().splitlines()]
            assert [x['args'][2] for x in calls]==['send','send','show']
            assert [x['credential_present'] for x in calls]==[True,True,False]
            again=self.notify(f,report=not before_start,receipt=not before_start)
            assert again['notification_id']==r['notification_id'] and len(self.posts(f))==1
            for p,b in before.items(): assert p.read_bytes()==b
            if kind=='feedback_required':
                self.notify(f,['--trigger','waiting_for_feedback'])
                assert len(self.posts(f))==1
        for trigger in ['optional_feedback','internal_correction']:
            f=self.fixture(trigger)
            r=self.notify(f,['--trigger',trigger]);assert r['state']=='suppressed' and len(self.posts(f))==0
            assert not (f/'events').exists() and not (f/'calls.jsonl').exists()
        self.checks.append('finished, necessary feedback, before-start stop, replay, waiting and suppression use the actual installed helper')

    def failures(self):
        f=self.fixture('missing-report')
        assert self.notify(f,report=False,expected=(2,))['state']=='not_attempted' and not self.posts(f)
        f=self.fixture('identity-mismatch')
        event=json.loads((f/'event.json').read_bytes());event['run']='other';put(f/'event.json',event)
        assert self.notify(f,expected=(2,))['state']=='not_attempted' and not self.posts(f)
        f=self.fixture('missing-credential')
        r=self.notify(f,['--credential-file',str(f/'config/missing')],expected=(2,))
        assert r['state']=='not_attempted' and not self.posts(f)
        f=self.fixture('invalid-environment')
        assert self.notify(f,env=self.secret_env('invalid'),expected=(2,))['state']=='not_attempted' and not self.posts(f)
        f=self.fixture('environment-priority')
        cfg=json.loads((f/'driver.json').read_bytes());cfg['expected_credential_sha256']=hashlib.sha256(self.other.encode()).hexdigest();put(f/'driver.json',cfg)
        assert self.notify(f,env=self.secret_env(self.other))['state']=='transport_acknowledged'
        for name,transport,want in [('rejected',{'kind':'http','status':403,'body':'action_prohibited'},'rejected'),('lost-response',{'kind':'unknown'},'unknown'),('zero-exit-rejection',{'kind':'http','status':403,'body':'action_prohibited'},'rejected')]:
            f=self.fixture(name,response=transport)
            if name=='zero-exit-rejection':
                cfg=json.loads((f/'driver.json').read_bytes());cfg['force_zero_exit']=True;put(f/'driver.json',cfg)
            before=(f/'report.json').read_bytes()
            r=self.notify(f,expected=(5,));assert r['state']==want and len(self.posts(f))==1
            if name=='zero-exit-rejection': assert r['apply_exit']==0
            again=self.notify(f,['--port-revision','2'],expected=(5,));assert again['state']==want and len(self.posts(f))==1
            assert (f/'report.json').read_bytes()==before
            event=json.loads((f/'event.json').read_bytes());event['public']['summary']='Different bytes for the same event.';put(f/'event.json',event)
            self.notify(f,expected=(5,));assert len(self.posts(f))==1 and (f/'report.json').read_bytes()==before
        self.checks.append('credential absence/priority, required sources, identity mismatch, rejected versus lost response, exit-zero rejection, changed event and port-revision replay')

    def incomplete_attempt_record(self):
        f=self.fixture('missing-attempt-id')
        first=self.notify(f)
        assert first['state']=='transport_acknowledged' and len(self.posts(f))==1
        request=Path(first['request'])
        preserved={p:p.read_bytes() for p in [f/'handoff.md',f/'report.json',f/'start.json',
                   f/'event.json',request,request.parent/'binding.json',Path(first['receipt'])]}
        marker=request.parent/'apply-started.json'
        attempt=json.loads(marker.read_bytes())
        del attempt['notification_id']
        put(marker,attempt)
        marker_bytes=marker.read_bytes()
        calls=(f/'calls.jsonl').read_bytes()
        replay=self.command(self.notify_args(f),expected=(0,2,5))
        result=json.loads(replay.stdout)
        observed={'state':result['state'],'exit':replay.returncode,'posts':len(self.posts(f)),
                  'local_return_preserved':all(p.read_bytes()==b for p,b in preserved.items()),
                  'marker_preserved':marker.read_bytes()==marker_bytes,
                  'no_new_cli_calls':(f/'calls.jsonl').read_bytes()==calls,
                  'product_outcome':result['product_outcome']}
        put(f/'observed.json',observed)
        assert observed['posts']==1 and observed['no_new_cli_calls']
        assert observed['local_return_preserved'] and observed['marker_preserved']
        assert observed['product_outcome']=='unchanged'
        assert observed['state']=='unknown', observed
        assert observed['exit']==5, observed
        self.checks.append('an incomplete preserved attempt remains unknown with exit 5, one POST and unchanged local return')

    def concurrency(self):
        f=self.fixture('concurrent')
        args=self.notify_args(f)
        def run():return subprocess.run(args,env=self.env,capture_output=True)
        with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
            results=list(pool.map(lambda _:run(),range(2)))
        for i,c in enumerate(results):
            assert self.secret.encode() not in c.stdout+c.stderr
            put(f/('process-%d.stdout'%i),c.stdout);put(f/('process-%d.stderr'%i),c.stderr)
            assert c.returncode in [0,5], c.stdout.decode()
        assert len(self.posts(f))==1
        assert self.notify(f)['state']=='transport_acknowledged' and len(self.posts(f))==1
        self.checks.append('concurrent helper processes and a later process reserve at most one POST')

    def unsafe_setup(self):
        f=self.fixture('unsafe-setup')
        script=self.skill/'scripts/bootstrap.py'
        args=[sys.executable,script,'--config-dir',f/'config','--state-root',f/'state']
        originals={p:p.read_bytes() for p in [f/'config/route.json', f/'config/webhook']}
        self.command(args,env=self.secret_env(self.other),expected=(2,))
        for p,b in originals.items(): assert p.read_bytes()==b
        (f/'config/webhook').chmod(0o644)
        self.command(args,expected=(2,))
        assert stat.S_IMODE((f/'config/webhook').stat().st_mode)==0o644
        self.notify(f,expected=(2,))
        (f/'config/webhook').chmod(0o600) # owned fixture test restoration only
        link=f/'config-link';link.symlink_to(f/'config',target_is_directory=True)
        self.command([sys.executable,script,'--config-dir',link,'--state-root',f/'state'],env=self.secret_env(),expected=(2,))
        self.notify(f,['--credential-file',str(link/'webhook')],expected=(2,))
        leaf=f/'config/webhook-link';leaf.symlink_to(f/'config/webhook')
        self.notify(f,['--credential-file',str(leaf)],expected=(2,))
        (f/'config').chmod(0o755)
        self.command(args,expected=(2,));self.notify(f,expected=(2,))
        (f/'config').chmod(0o700)
        # A foreign uid cannot be created by an unprivileged harness. Inject the
        # caller uid at the real file validation boundary, with no chown/escalation.
        snippet="import sys;sys.dont_write_bytecode=True;sys.path.insert(0,sys.argv[1]);import common;uid=common.os.getuid();common.os.getuid=lambda:uid+1;import bootstrap;sys.argv=['bootstrap','--config-dir',sys.argv[2],'--state-root',sys.argv[3]];sys.exit(bootstrap.main())"
        self.command([sys.executable,'-c',snippet,self.skill/'scripts',f/'config',f/'state'],expected=(2,))
        # An existing installation with different bytes is preserved as a whole.
        p=self.skill/'SKILL.md';original=p.read_bytes();p.write_bytes(original+b'\nLocal customization.\n')
        self.command(['make','install-agent-notify','AGENT_NOTIFY_DEST='+str(self.skill)],expected=(2,))
        assert p.read_bytes()==original+b'\nLocal customization.\n';p.write_bytes(original)
        install_link=self.root/'installed/skill-link';install_link.symlink_to(self.skill,target_is_directory=True)
        self.command(['make','install-agent-notify','AGENT_NOTIFY_DEST='+str(install_link)],expected=(2,))
        self.skill.chmod(0o755)
        self.command(['make','install-agent-notify','AGENT_NOTIFY_DEST='+str(self.skill)],expected=(2,));self.skill.chmod(0o700)
        snippet="import sys;sys.dont_write_bytecode=True;sys.path.insert(0,sys.argv[1]);import common;uid=common.os.getuid();common.os.getuid=lambda:uid+1;import install;sys.argv=['install','--destination',sys.argv[2]];sys.exit(install.main())"
        self.command([sys.executable,'-c',snippet,self.skill/'scripts',self.skill],expected=(2,))
        for p,b in originals.items():assert p.read_bytes()==b
        assert not self.posts(f)
        self.checks.append('bootstrap/install conflicts, directory/leaf symlinks, unsafe modes, and injected foreign-owner observation are rejected without repair')

    def legacy_state(self):
        # Build the exact pre-change command graph. Every old state byte below
        # comes from that binary, never a hand-crafted database or new serializer.
        baseline=self.root/'baseline-source';baseline.mkdir(mode=0o700)
        c=subprocess.run(['git','archive',BASE],cwd=REPO,env=self.env,capture_output=True,check=True)
        with tarfile.open(fileobj=io.BytesIO(c.stdout)) as archive:
            archive.extractall(baseline,filter='data')
        old=self.root/'baseline-driver'
        self.command(['go','test','-c','-o',old,'./test/notification-driver'],cwd=baseline)
        f=self.fixture('legacy-mixed')
        route=f/'config/route.json';baseargs=['--config',f/'driver.json','--','workflow','notification']
        records=[]
        for version in [1,2]:
            for state,transport in [('transport_acknowledged',{'kind':'http','status':200,'body':'ok'}),('rejected',{'kind':'http','status':403,'body':'action_prohibited'}),('unknown',{'kind':'unknown'})]:
                name='v%d-%s'%(version,state)
                cfg=json.loads((f/'driver.json').read_bytes());cfg['transport']=transport;put(f/'driver.json',cfg)
                source={'kind':'external','activity':'fixture/activity','run':'run-1','worktree':str(f/'work'),'handoff':loc(f/'handoff.md'),'start_receipt':loc(f/'start.json',True),'report':loc(f/'report.json',True)}
                req={'kind':'ply.workflow.notification-request','schema_version':version,'route':'ply-log','source':source,'sender':{'actor_claim':'legacy-fixture'}}
                if version==1:
                    req['gate']={'id':name,'revision':1,'opened_at':'2026-10-04T14:00:00Z','reason':'result_control','state':'waiting_for_human'}
                    req['public']={'task_title':'Legacy report','next_action':'Start result control.'}
                else:
                    event=json.loads((f/'event.json').read_bytes());event['event_id']=name
                    event_path=put(f/(name+'-event.json'),event)
                    req['public']=event['public']
                    req['event']={'id':name,'type':event['event_type'],'phase':event['phase'],'occurred_at':event['occurred_at'],'record':loc(event_path,True)}
                path=put(f/(name+'.json'),req)
                send=['send','--file',path,'--route',route,'--format','json']
                p=json.loads(self.command([old,*baseargs,*send,'--check'],env=self.secret_env()).stdout)
                old_result=self.command([old,*baseargs,*send,'--apply','--confirm',p['confirmation']],env=self.secret_env(),expected=(0,5)).stdout
                records.append({'version':version,'state':state,'path':str(path),'preview':p,'result':json.loads(old_result)})
        original=(f/'state/state.json').read_bytes()
        put(f/'baseline-state.json',original)
        put(f/'baseline-provenance.json',{'base_oid':BASE,'producer_sha256':sha(old),'state_sha256':sha(f/'baseline-state.json'),'records':records,'network':False})
        for record in records:
            p=record['preview'];send=['send','--file',record['path'],'--route',route,'--format','json']
            r=json.loads(self.command([self.driver,*baseargs,'show',p['notification_id'],'--route',route,'--format','json']).stdout)
            assert r==record['result'],(record['version'],record['state'],r)
            preview=json.loads(self.command([self.driver,*baseargs,*send,'--check'],env=self.secret_env()).stdout)
            assert preview['payload']==p['payload'] and preview['payload_sha256']==p['payload_sha256']
            self.command([self.driver,*baseargs,*send,'--apply','--confirm',p['confirmation']],env=self.secret_env())
            assert len(self.posts(f))==6 and (f/'state/state.json').read_bytes()==original
        for record in records:
            p=record['preview'];retry=['retry',p['notification_id'],'--route',route,'--format','json']
            if record['state']!='rejected':
                self.command([self.driver,*baseargs,*retry,'--check'],env=self.secret_env(),expected=(4,))
                continue
            cfg=json.loads((f/'driver.json').read_bytes());cfg['transport']={'kind':'http','status':200,'body':'ok'};put(f/'driver.json',cfg)
            p2=json.loads(self.command([self.driver,*baseargs,*retry,'--check'],env=self.secret_env()).stdout)
            self.command([self.driver,*baseargs,*retry,'--apply','--confirm',p2['confirmation']],env=self.secret_env())
            count=len(self.posts(f))
            assert self.posts(f)[-1]['payload']==p['payload']
            self.command([self.driver,*baseargs,*retry,'--apply','--confirm',p2['confirmation']],env=self.secret_env())
            assert len(self.posts(f))==count
        assert len(self.posts(f))==8
        # Add v3 to this same unchanged route/state; older records must still show.
        task,event=self.compact_task(f,'codex','done','mixed-new')
        result=self.notify(f,cwd=task['origin_cwd'])
        assert result['state']=='transport_acknowledged' and len(self.posts(f))==9
        for record in records:
            self.command([self.driver,*baseargs,'show',record['preview']['notification_id'],'--route',route,'--format','json'])
        # Recast the exact prior v2 identity, including unknown, as v3.
        for record in records:
            if record['version']!=2:continue
            task,event=self.compact_task(f,'codex','done',json.loads(Path(record['path']).read_bytes())['event']['id'])
            result=self.notify(f,cwd=task['origin_cwd'],expected=(0,5))
            assert result['reason']=='preview_rejected' and result['preview_exit']==4 and len(self.posts(f))==9
        # The old helper first writes an immutable v2 request; the new helper
        # must read that exact byte sequence without another apply.
        old_helper=self.fixture('legacy-helper')
        wrapper=old_helper/'ply-fixture'
        wrapper.write_text(wrapper.read_text().replace(str(self.driver),str(old)))
        argv=self.notify_args(old_helper);argv[1]=str(baseline/'skills/ply-agent-notify/scripts/notify.py')
        result=json.loads(self.command(argv).stdout)
        request=Path(result['request']);original_request=request.read_bytes()
        wrapper.write_text(wrapper.read_text().replace(str(old),str(self.driver)))
        repeated=self.notify(old_helper)
        assert repeated['reason']=='preserved_attempt_readback' and request.read_bytes()==original_request and len(self.posts(old_helper))==1
        # Exercise the actual explicit upgrade command only in this private fixture.
        installed=self.root/'installed'/'old-notify';backup=self.root/'installed'/'old-notify-backup'
        self.command([sys.executable,baseline/'skills/ply-agent-notify/scripts/install.py','--destination',installed])
        self.command([sys.executable,REPO/'skills/ply-agent-notify/scripts/install.py','--destination',installed,'--upgrade-from',baseline/'skills/ply-agent-notify','--backup',backup])
        for path in self.skill.rglob('*'):
            if not path.is_file():continue
            name=path.relative_to(self.skill)
            assert (installed/name).read_bytes()==path.read_bytes()
            assert (backup/name).read_bytes()==(baseline/'skills/ply-agent-notify'/name).read_bytes()
            assert stat.S_IMODE((installed/name).stat().st_mode)==stat.S_IMODE(path.stat().st_mode)
            assert stat.S_IMODE((backup/name).stat().st_mode)==stat.S_IMODE(path.stat().st_mode)
        self.checks.append('base c8f9dd0 binary produced mixed v1/v2 acknowledged/rejected/unknown state; exact readback/payload/replay, permitted retry, mixed v3 and version conflicts pass')

    def compact_task(self,f,provider,status,event_id=None):
        origin=f/'ply'/('planning' if provider=='codex' else 'varsler');origin.mkdir(mode=0o700,parents=True,exist_ok=True)
        task=json.loads((f/'task.json').read_bytes())
        task.update(kind='PlyAgentNotificationTask@2',provider=provider,origin_cwd=str(origin),context_root=str(f),timezone='Europe/Oslo');put(f/'task.json',task)
        event=json.loads((f/'event.json').read_bytes());event['kind']='PlyAgentNotificationEvent@2';event['public']['status']=status
        if event_id:event['event_id']=event_id
        put(f/'event.json',event)
        return task,event

    def compact_journeys(self):
        for provider in ['codex','claude']:
            for status,label,event_type in [('ready_for_review','Ready for review','agent_finished'),('ready_for_your_check','Ready for your check','agent_finished'),('done','Done','agent_finished'),('stopped','Stopped','agent_stopped'),('needs_answer','Needs answer','feedback_required')]:
                f=self.fixture('compact-'+provider+'-'+status,event_type)
                task,event=self.compact_task(f,provider,status)
                result=self.notify(f,cwd=task['origin_cwd'])
                payload=self.posts(f)[0]['payload']
                assert payload['text']=='16:00:00 · '+provider+' · '+('ply/planning' if provider=='codex' else 'ply/varsler')+'\n'+label+': '+event['public']['task_title']+' — '+event['public']['summary']+'\nNext: '+event['public']['next_action']
                request=Path(result['request']);before=request.read_bytes()
                assert json.loads(before)['source']['worktree']!=task['origin_cwd']
                self.notify(f,cwd=task['origin_cwd'],env={**self.env,'TZ':'Pacific/Honolulu'})
                assert len(self.posts(f))==1 and request.read_bytes()==before
                # The same helper event key prevents replay during a schema switch.
                event['kind']='PlyAgentNotificationEvent@1';del event['public']['status'];put(f/'event.json',event)
                task['kind']='PlyAgentNotificationTask@1'
                for key in ['provider','origin_cwd','context_root','timezone']:del task[key]
                put(f/'task.json',task)
                replay=self.notify(f,expected=(5,));assert replay['state']=='unknown' and len(self.posts(f))==1
        f=self.fixture('compact-wrong-cwd');task,event=self.compact_task(f,'codex','done')
        assert self.notify(f,expected=(2,))['reason']=='origin_cwd_mismatch'
        assert not (f/'events').exists() and not self.posts(f)
        self.checks.append('real helper/command payloads for both providers and five statuses, origin distinct from target, TZ-independent replay, cross-version helper conflict and wrong cwd rejection')

    def frozen_timezone(self):
        import zipfile
        f=self.fixture('frozen-zone',response={'kind':'http','status':403,'body':'action_prohibited'})
        task,event=self.compact_task(f,'codex','done')
        # Both runtimes see a private named zone for the first attempt. Removing
        # this zone simulates a tzdata update without touching host configuration.
        tz=f/'tz';(tz/'Fixture').mkdir(parents=True,mode=0o700)
        goroot=self.command(['go','env','GOROOT']).stdout.decode().strip()
        with zipfile.ZipFile(Path(goroot)/'lib/time/zoneinfo.zip') as z:
            (tz/'Fixture/Frozen').write_bytes(z.read('Europe/Oslo'))
        task['timezone']='Fixture/Frozen';put(f/'task.json',task)
        env={**self.env,'ZONEINFO':str(tz),'PYTHONTZPATH':str(tz)}
        result=self.notify(f,cwd=task['origin_cwd'],env=env,expected=(5,))
        assert result['state']=='rejected'
        req=Path(result['request']);request_bytes=req.read_bytes();state=(f/'state/state.json').read_bytes()
        (tz/'Fixture/Frozen').rename(tz/'Fixture/Retired')
        repeat=self.notify(f,cwd=task['origin_cwd'],env=env,expected=(5,))
        assert repeat['state']=='rejected' and len(self.posts(f))==1
        assert req.read_bytes()==request_bytes and (f/'state/state.json').read_bytes()==state
        cfg=json.loads((f/'driver.json').read_bytes());cfg['now']='2026-10-12T01:00:00Z';put(f/'driver.json',cfg)
        env['TZ']='Pacific/Honolulu'
        args=[self.driver,'--config',f/'driver.json','--','workflow','notification']
        route=f/'config/route.json'
        retry=['retry',result['notification_id'],'--route',route,'--format','json']
        p=json.loads(self.command([*args,*retry,'--check'],env={**env,ENV:self.secret}).stdout)
        self.command([*args,*retry,'--apply','--confirm',p['confirmation']],env={**env,ENV:self.secret},expected=(5,))
        assert len(self.posts(f))==2 and self.posts(f)[0]['payload']==self.posts(f)[1]['payload']
        # A genuinely fresh event with the retired zone must fail, not use UTC.
        event['event_id']='new-retired-zone-event';put(f/'event.json',event)
        assert self.notify(f,cwd=task['origin_cwd'],env=env,expected=(2,))['reason']=='unknown_timezone_or_unrepresentable_time'
        self.checks.append('removed IANA zone leaves historical database, helper replay, explicit retry and frozen payload intact; fresh input rejects the unavailable zone')

    def privacy(self):
        for directory in [self.root/'commands',self.root/'installed',*[p for p in self.root.iterdir() if p.is_dir() and (p/'task.json').exists()]]:
            for p in directory.rglob('*'):
                if p.is_symlink() or not p.is_file() or p in self.credentials:
                    continue
                # Skip files reached through intentionally unsafe fixture directory links.
                if any(q.is_symlink() for q in p.parents):continue
                b=p.read_bytes()
                assert self.secret.encode() not in b and self.other.encode() not in b, 'secret leaked into '+str(p.relative_to(self.root))
        self.checks.append('canary credentials occur only in the explicitly private bootstrap credential files, never output/payload/receipts')

    def run(self):
        self.command([sys.executable,REPO/'test/agent_notify_regression_test.py'])
        self.setup();self.journeys();self.failures();self.incomplete_attempt_record();self.concurrency();self.unsafe_setup();self.compact_journeys();self.frozen_timezone();self.legacy_state();self.privacy()
        # Production executable help and preview are exercised without transport.
        for parts in [[],['workflow'],['workflow','notification'],['workflow','notification','send'],['workflow','notification','show'],['workflow','notification','retry']]:
            self.command([self.ply,*parts,'--help'])
        for script in ['notify.py','install.py']:
            self.command([sys.executable,self.skill/'scripts'/script,'--help'])
        f=self.fixture('production-preview')
        # Create/preserve a request through the real helper, then production check.
        r=self.notify(f,['--credential-file',str(f/'config/missing')],expected=(2,))
        c=self.command([self.ply,'workflow','notification','send','--file',r['request'],'--route',f/'config/route.json','--check','--format','json'],env=self.secret_env())
        assert json.loads(c.stdout)['allowed'] is True and not self.posts(f) and not (f/'state/state.json').exists()
        self.privacy()
        put(self.root/'result.json',{'kind':'PlyAgentNotificationAcceptance@1','outcome':'passed','checks':self.checks,'commands':'commands.json','root':str(self.root),'live_network':False,'legacy_baseline':BASE})
        print(json.dumps({'outcome':'passed','checks':len(self.checks),'fixture':str(self.root),'commands':len(self.commands)}))


def main():
    parser=argparse.ArgumentParser(description=__doc__);parser.add_argument('--root');args=parser.parse_args()
    if args.root:
        root=Path(args.root)
        assert root.is_absolute() and root.parent.resolve()==root.parent and not os.path.lexists(root)
        root.mkdir(mode=0o700)
    else:
        root=Path(tempfile.mkdtemp(prefix='ply-agent-notify-')).resolve()
    print('Retained offline fixture: '+str(root),flush=True)
    h=Harness(root)
    try:h.run()
    except Exception:
        put(root/'result.json',{'kind':'PlyAgentNotificationAcceptance@1','outcome':'failed','checks':h.checks,'commands':'commands.json','root':str(root),'live_network':False})
        raise


if __name__=='__main__':main()
