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
HELPER_BASELINE = 'c8f9dd00ee36d3f57191355c278d982808d421a8'
PRE_APP_LINK = '69e902e513fd44002657cbe0a48e7a77e2241aef'
# The source before DP01 removed requests v1-v3 and the helper's bound-Task entrance.
PRE_DP01 = '20b294055569d09e4ca42b5f48a5b8f6731b0b57'
CASE = {'agent_finished':('done','The local return is preserved.','Review the preserved local return.','coordinator'),
        'agent_stopped':('stopped','The local return is preserved.','Review the preserved local return.','coordinator'),
        'feedback_required':('needs_answer','Which agreed target should I use?','Answer in the active agent task.','user')}
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
        # No real Herdr lookup in this fake-only run; short requests stay v4.
        self.env.pop('HERDR_ENV', None)
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
        self.cases = {}

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
        put(folder/'work'/'return.json', {'private':'PRIVATE_SOURCE_SENTINEL'})
        status, summary, action, actor = CASE[kind]
        self.cases[folder] = {'event-id':name,'status':status,'phase':phase,'title':'A useful agent task','summary':summary,
                              'next-action':action,'next-actor':actor,'provider':'codex','timezone':'Europe/Oslo'}
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

    def notify_args(self, folder, extra=()):
        args = [sys.executable,self.skill/'scripts/notify.py','--activity','fixture/activity','--run','run-1',
                '--route',folder/'config/route.json','--ply-bin',folder/'ply-fixture','--credential-file',folder/'config/webhook']
        for name,value in self.cases[folder].items():
            args += ['--'+name,value]
        return [str(x) for x in [*args,*extra]]

    def notify(self, folder, extra=(), expected=(0,), env=None, cwd=None):
        c = self.command(self.notify_args(folder, extra),env=env,expected=expected,cwd=cwd or folder/'work')
        return json.loads(c.stdout)

    def posts(self, folder):
        path = folder/'posts.jsonl'
        return [json.loads(line) for line in path.read_text().splitlines()] if path.exists() else []

    def journeys(self):
        for kind,label in [('agent_finished','Done'),('feedback_required','Needs answer'),('agent_stopped','Stopped')]:
            before_start = kind=='agent_stopped'
            f=self.fixture(kind,kind,'before_start' if before_start else 'after_start')
            local=(f/'work/return.json').read_bytes()
            r=self.notify(f)
            assert r['state']=='transport_acknowledged' and r['product_outcome']=='unchanged'
            posts=self.posts(f);assert len(posts)==1
            text=posts[0]['payload']['text']
            assert ('\n'+label+': A useful agent task — ') in text
            assert 'PRIVATE_SOURCE_SENTINEL' not in text and str(f) not in text
            request=json.loads(Path(r['request']).read_bytes())
            assert request['schema_version']==4 and set(request['source'])=={'kind','activity','run','worktree'}
            calls=[json.loads(x) for x in (f/'calls.jsonl').read_text().splitlines()]
            assert [x['args'][2] for x in calls]==['send','send','show']
            assert [x['credential_present'] for x in calls]==[True,True,False]
            again=self.notify(f)
            assert again['notification_id']==r['notification_id'] and len(self.posts(f))==1
            assert (f/'work/return.json').read_bytes()==local
            if kind=='feedback_required':
                self.notify(f,['--trigger','waiting_for_feedback'])
                assert len(self.posts(f))==1
        for trigger in ['optional_feedback','internal_correction']:
            f=self.fixture(trigger)
            r=self.notify(f,['--trigger',trigger]);assert r['state']=='suppressed' and len(self.posts(f))==0
            assert not (f/'config/notification-events').exists() and not (f/'calls.jsonl').exists()
        self.checks.append('finished, necessary feedback, before-start stop, replay, waiting and suppression use the actual installed helper')

    def failures(self):
        f=self.fixture('unsafe-text')
        self.cases[f]['summary']='Hi <!here>'
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
            before=(f/'work/return.json').read_bytes()
            r=self.notify(f,expected=(5,));assert r['state']==want and len(self.posts(f))==1
            if name=='zero-exit-rejection': assert r['apply_exit']==0
            again=self.notify(f,['--port-revision','2'],expected=(5,));assert again['state']==want and len(self.posts(f))==1
            assert (f/'work/return.json').read_bytes()==before
            self.cases[f]['summary']='Different bytes for the same event.'
            self.notify(f,expected=(5,));assert len(self.posts(f))==1 and (f/'work/return.json').read_bytes()==before
        self.checks.append('credential absence/priority, unsafe text, rejected versus lost response, exit-zero rejection, changed event and port-revision replay')

    def incomplete_attempt_record(self):
        f=self.fixture('missing-attempt-id')
        first=self.notify(f)
        assert first['state']=='transport_acknowledged' and len(self.posts(f))==1
        request=Path(first['request'])
        preserved={p:p.read_bytes() for p in [f/'work/return.json',request,request.parent/'binding.json']}
        marker=request.parent/'apply-started.json'
        attempt=json.loads(marker.read_bytes())
        del attempt['notification_id']
        put(marker,attempt)
        marker_bytes=marker.read_bytes()
        calls=(f/'calls.jsonl').read_bytes()
        replay=self.command(self.notify_args(f),expected=(0,2,5),cwd=f/'work')
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
        def run():return subprocess.run(args,env=self.env,capture_output=True,cwd=f/'work')
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

    def stored_history_and_helper_upgrade(self):
        """Records written before DP01 by its CLI and helper are history: read, never sent again."""
        source=self.root/'pre-dp01-source';source.mkdir(mode=0o700)
        c=subprocess.run(['git','archive',PRE_DP01],cwd=REPO,env=self.env,capture_output=True,check=True)
        with tarfile.open(fileobj=io.BytesIO(c.stdout)) as archive:
            archive.extractall(source,filter='data')
        old=self.root/'pre-dp01-driver'
        self.command(['go','test','-c','-o',old,'./test/notification-driver'],cwd=source)
        f=self.fixture('request-versions')
        route=f/'config/route.json';baseargs=['--config',f/'driver.json','--','workflow','notification']
        handoff=put(f/'handoff.md',b'Offline task mandate.\n')
        report=put(f/'report.json',{'kind':'FixtureReport@1','activity':'fixture/activity','run':'run-1'})
        start=put(f/'start.json',{'kind':'FixtureStart@1','activity':'fixture/activity','run':'run-1'})
        origin=f/'ply'/'planning';origin.mkdir(mode=0o700,parents=True)
        records=[]
        for version in [1,2,3]:
            for state,transport in [('transport_acknowledged',{'kind':'http','status':200,'body':'ok'}),('rejected',{'kind':'http','status':403,'body':'action_prohibited'}),('unknown',{'kind':'unknown'})]:
                name='v%d-%s'%(version,state)
                cfg=json.loads((f/'driver.json').read_bytes());cfg['transport']=transport;put(f/'driver.json',cfg)
                source_value={'kind':'external','activity':'fixture/activity','run':'run-1','worktree':str(f/'work'),'handoff':loc(handoff),'start_receipt':loc(start,True),'report':loc(report,True)}
                req={'kind':'ply.workflow.notification-request','schema_version':version,'route':'ply-log','source':source_value,'sender':{'actor_claim':'request-version-fixture'}}
                if version==1:
                    req['gate']={'id':name,'revision':1,'opened_at':'2026-10-04T14:00:00Z','reason':'result_control','state':'waiting_for_human'}
                    req['public']={'task_title':'Version 1 report','next_action':'Start result control.'}
                else:
                    public={'task_title':'Varsler','summary':'Ferdig.','next_action':'No action required.','next_actor':'coordinator'}
                    if version==3:public['status']='done'
                    kind='PlyAgentNotificationEvent@%d'%(version-1)
                    record=put(f/(name+'-event.json'),{'kind':kind,'activity':'fixture/activity','run':'run-1','event_id':name,'event_type':'agent_finished','phase':'after_start','occurred_at':'2026-10-04T14:00:00Z','public':public})
                    req['public']=public
                    req['event']={'id':name,'type':'agent_finished','phase':'after_start','occurred_at':'2026-10-04T14:00:00Z','record':loc(record,True)}
                    if version==3:
                        req['presentation']={'provider':'codex','origin_cwd':str(origin),'context_root':str(f),'context':'ply/planning','timezone':'Europe/Oslo','local_occurred_at':'2026-10-04T16:00:00+02:00'}
                path=put(f/(name+'.json'),req)
                send=['send','--file',path,'--route',route,'--format','json']
                p=json.loads(self.command([old,*baseargs,*send,'--check'],env=self.secret_env()).stdout)
                result=json.loads(self.command([old,*baseargs,*send,'--apply','--confirm',p['confirmation']],env=self.secret_env(),expected=(0,5)).stdout)
                records.append({'version':version,'state':state,'path':str(path),'preview':p,'result':result})
        original=(f/'state/state.json').read_bytes()
        put(f/'request-version-provenance.json',{'producer':PRE_DP01,'state_sha256':sha(f/'state/state.json'),'records':records,'network':False})
        cfg=json.loads((f/'driver.json').read_bytes());cfg['transport']={'kind':'http','status':200,'body':'ok'};put(f/'driver.json',cfg)
        for record in records:
            p=record['preview'];send=['send','--file',record['path'],'--route',route,'--format','json']
            r=json.loads(self.command([self.driver,*baseargs,'show',p['notification_id'],'--route',route,'--format','json']).stdout)
            assert (r['notification_id'],r['state'],r['request_sha256'])==(p['notification_id'],record['result']['state'],p['request_sha256']),(record['version'],record['state'],r)
            self.command([self.driver,*baseargs,*send,'--check'],env=self.secret_env(),expected=(2,))
            self.command([self.driver,*baseargs,*send,'--apply','--confirm',p['confirmation']],env=self.secret_env(),expected=(2,))
            self.command([self.driver,*baseargs,'retry',p['notification_id'],'--route',route,'--format','json','--check'],env=self.secret_env(),expected=(4,))
        assert len(self.posts(f))==9 and (f/'state/state.json').read_bytes()==original
        # A v4 short event with an identity a stored v2/v3 record already used conflicts and only reads.
        for record in records:
            if record['version']==1:continue
            self.cases[f]['event-id']=json.loads(Path(record['path']).read_bytes())['event']['id']
            result=self.notify(f,expected=(0,5))
            assert result['reason']=='preview_rejected' and result['preview_exit']==4 and len(self.posts(f))==9
        self.cases[f]['event-id']='mixed-new'
        assert self.notify(f)['state']=='transport_acknowledged' and len(self.posts(f))==10
        self.checks.append('pre-DP01 CLI v1/v2/v3 acknowledged/rejected/unknown records read back as history; resend, retry and a v4 recast of the same identity never post')
        # The pre-DP01 helper's bound Task wrote its shared marker beside the route; the
        # short form on that identity only reads it.
        legacy=self.fixture('legacy-helper',response={'kind':'unknown'})
        task={'kind':'PlyAgentNotificationTask@1','activity':'fixture/activity','run':'run-1','worktree':str(legacy/'work'),'handoff':str(put(legacy/'handoff.md',b'Mandate.\n')),
              'route':str(legacy/'config/route.json'),'event_root':str(legacy/'events'),'actor_claim':'offline-agent','notification_authorized':True}
        put(legacy/'task.json',task)
        put(legacy/'event.json',{'kind':'PlyAgentNotificationEvent@1','activity':'fixture/activity','run':'run-1','event_id':'legacy-helper','event_type':'agent_stopped','phase':'before_start',
                                 'occurred_at':'2026-10-04T14:00:00Z','public':{'task_title':'A useful agent task','summary':'Stopped.','next_action':'Review the local return.','next_actor':'coordinator'}})
        (legacy/'ply-old').write_text((legacy/'ply-fixture').read_text().replace(repr(str(self.driver)),repr(str(old))));(legacy/'ply-old').chmod(0o700)
        argv=[sys.executable,source/'skills/ply-agent-notify/scripts/notify.py','--task',legacy/'task.json','--event',legacy/'event.json','--ply-bin',legacy/'ply-old','--credential-file',legacy/'config/webhook']
        first=json.loads(self.command(argv,expected=(5,)).stdout);assert first['state']=='unknown' and len(self.posts(legacy))==1
        self.cases[legacy].update({'status':'stopped','phase':'before_start','summary':'Stopped.','next-action':'Review the local return.'})
        repeated=self.notify(legacy,expected=(5,))
        assert repeated['reason']=='preserved_attempt_readback' and repeated['state']=='unknown' and len(self.posts(legacy))==1
        self.checks.append('an event the pre-DP01 helper started through its bound Task is only read back by the short form')
        baseline=self.root/'helper-baseline-source';baseline.mkdir(mode=0o700)
        c=subprocess.run(['git','archive',HELPER_BASELINE,'skills/ply-agent-notify'],cwd=REPO,env=self.env,capture_output=True,check=True)
        with tarfile.open(fileobj=io.BytesIO(c.stdout)) as archive:
            archive.extractall(baseline,filter='data')
        # Exercise the actual explicit upgrade command only in this private fixture.
        installed=self.root/'installed'/'old-notify';backup=self.root/'installed'/'old-notify-backup'
        self.command([sys.executable,baseline/'skills/ply-agent-notify/scripts/install.py','--destination',installed])
        # The owner removes the retired template before the exact upgrade check.
        for folder in [baseline/'skills/ply-agent-notify',installed]:(folder/'references/task-template.json').unlink()
        self.command([sys.executable,REPO/'skills/ply-agent-notify/scripts/install.py','--destination',installed,'--upgrade-from',baseline/'skills/ply-agent-notify','--backup',backup])
        for path in self.skill.rglob('*'):
            if not path.is_file():continue
            name=path.relative_to(self.skill)
            assert (installed/name).read_bytes()==path.read_bytes()
            assert (backup/name).read_bytes()==(baseline/'skills/ply-agent-notify'/name).read_bytes()
            assert stat.S_IMODE((installed/name).stat().st_mode)==stat.S_IMODE(path.stat().st_mode)
            assert stat.S_IMODE((backup/name).stat().st_mode)==stat.S_IMODE(path.stat().st_mode)
        self.checks.append('explicit fixture upgrade from helper baseline c8f9dd0 retains private backup and permissions')

    def compact_journeys(self):
        for provider in ['codex','claude']:
            for status,label,event_type in [('ready_for_review','Ready for review','agent_finished'),('ready_for_your_check','Ready for your check','agent_finished'),('done','Done','agent_finished'),('stopped','Stopped','agent_stopped'),('needs_answer','Needs answer','feedback_required')]:
                f=self.fixture('compact-'+provider+'-'+status,event_type)
                case=self.cases[f];case.update(provider=provider,status=status)
                result=self.notify(f)
                payload=self.posts(f)[0]['payload']
                clock=json.loads(Path(result['request']).read_bytes())['presentation']['local_occurred_at'][11:19]
                assert payload['text']==clock+' · '+provider+'\n'+label+': '+case['title']+' — '+case['summary']+'\nNext: '+case['next-action'],payload
                request=Path(result['request']);before=request.read_bytes()
                self.notify(f,env={**self.env,'TZ':'Pacific/Honolulu'})
                assert len(self.posts(f))==1 and request.read_bytes()==before
        f=self.fixture('compact-outside-root')
        other=f/'other';other.mkdir(mode=0o700)
        assert self.notify(f,['--context-root',str(other)],expected=(2,))['reason']=='origin_outside_context_root'
        assert not self.posts(f)
        self.checks.append('real helper/command payloads for both providers and five statuses, TZ-independent replay and origin outside the context root rejected')

    def frozen_timezone(self):
        import zipfile
        f=self.fixture('frozen-zone',response={'kind':'http','status':403,'body':'action_prohibited'})
        # Both runtimes see a private named zone for the first attempt. Removing
        # this zone simulates a tzdata update without touching host configuration.
        tz=f/'tz';(tz/'Fixture').mkdir(parents=True,mode=0o700)
        goroot=self.command(['go','env','GOROOT']).stdout.decode().strip()
        with zipfile.ZipFile(Path(goroot)/'lib/time/zoneinfo.zip') as z:
            (tz/'Fixture/Frozen').write_bytes(z.read('Europe/Oslo'))
        self.cases[f]['timezone']='Fixture/Frozen'
        env={**self.env,'ZONEINFO':str(tz),'PYTHONTZPATH':str(tz)}
        result=self.notify(f,env=env,expected=(5,))
        assert result['state']=='rejected'
        req=Path(result['request']);request_bytes=req.read_bytes();state=(f/'state/state.json').read_bytes()
        (tz/'Fixture/Frozen').rename(tz/'Fixture/Retired')
        repeat=self.notify(f,env=env,expected=(5,))
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
        self.cases[f]['event-id']='new-retired-zone-event'
        assert self.notify(f,env=env,expected=(2,))['reason']=='unknown_timezone_or_unrepresentable_time'
        self.checks.append('removed IANA zone leaves historical database, helper replay, explicit retry and frozen payload intact; fresh input rejects the unavailable zone')

    def privacy(self):
        for directory in [self.root/'commands',self.root/'installed',*[p for p in self.root.iterdir() if p.is_dir() and (p/'config').exists()]]:
            for p in directory.rglob('*'):
                if p.is_symlink() or not p.is_file() or p in self.credentials:
                    continue
                # Skip files reached through intentionally unsafe fixture directory links.
                if any(q.is_symlink() for q in p.parents):continue
                b=p.read_bytes()
                assert self.secret.encode() not in b and self.other.encode() not in b, 'secret leaked into '+str(p.relative_to(self.root))
        self.checks.append('canary credentials occur only in the explicitly private bootstrap credential files, never output/payload/receipts')

    def short_args(self, f, event_id, summary='Ferdig.'):
        return [str(x) for x in [sys.executable,self.skill/'scripts/notify.py','--activity','fixture/activity','--run','run-1','--event-id',event_id,'--status','done',
                '--title','ℹ️ Jobb','--summary',summary,'--next-action','No action required.','--provider','claude','--route',f/'config/route.json',
                '--ply-bin',f/'ply-fixture','--credential-file',f/'config/webhook']]

    def short_journeys(self):
        f=self.fixture('agent_finished_short')
        short=lambda *a,**k:json.loads(self.command(self.short_args(f,'agent_finished_short',*a),cwd=f/'work',**k).stdout)
        r=short();assert r['state']=='transport_acknowledged' and 'receipt' not in r
        posts=self.posts(f);assert len(posts)==1
        text=posts[0]['payload']['text'];assert text.endswith('Done: ℹ️ Jobb — Ferdig.\nNext: No action required.') and ' · claude\n' in text and 'work' not in text
        request=Path(r['request']);before=request.read_bytes();parsed=json.loads(before)
        assert parsed['schema_version']==4 and 'handoff' not in parsed['source'] and 'record' not in parsed['event']
        assert sorted(p.name for p in request.parent.iterdir())==['.lock','apply-started.json','binding.json','request.json']
        again=short();assert again['notification_id']==r['notification_id'] and again['state']=='transport_acknowledged'
        assert request.read_bytes()==before and len(self.posts(f))==1
        changed=short('Endret.',expected=(5,));assert changed['state']=='unknown' and len(self.posts(f))==1
        f=self.fixture('unknown_short',response={'kind':'unknown'})
        args=self.short_args(f,'unknown_short')
        r=json.loads(self.command(args,cwd=f/'work',expected=(5,)).stdout);assert r['state']=='unknown' and len(self.posts(f))==1
        self.checks.append('short form sends once without documents; replay, changed text and unknown only read')

    def app_link_journey(self):
        f=self.fixture('app-link');link='https://plyapp-pwa.web.app';owner=f/'config/app-link.json'
        short=lambda event_id,**k:json.loads(self.command(self.short_args(f,event_id),cwd=f/'work',**k).stdout)
        first=short('before-owner');assert json.loads(Path(first['request']).read_bytes())['schema_version']==4
        put(owner,{'app_link':link})
        assert short('before-owner')['reason']=='preserved_attempt_readback' and len(self.posts(f))==1
        r=short('with-owner');posts=self.posts(f);payload=posts[-1]['payload']
        assert len(posts)==2 and 'blocks' not in posts[0]['payload']
        assert {k:v for k,v in payload.items() if k!='blocks'}=={'text':payload['text'],'mrkdwn':False,'parse':'none','unfurl_links':False,'unfurl_media':False}
        assert payload['blocks']==[{'type':'section','text':{'type':'plain_text','text':payload['text'],'emoji':False}},
                                   {'type':'context','elements':[{'type':'mrkdwn','text':'<'+link+'|Åpne plyapp>','verbatim':True}]}]
        request=Path(r['request']).read_bytes()
        for change in [{'app_link':'https://other.example'},b'{',None]:
            put(owner,change) if change is not None else owner.unlink()
            assert short('with-owner')['reason']=='preserved_attempt_readback' and len(self.posts(f))==2 and Path(r['request']).read_bytes()==request
        # Rolling back the CLI after a v6 record fails closed on the whole route; it never sends.
        source=self.root/'pre-app-link-source';source.mkdir(mode=0o700)
        c=subprocess.run(['git','archive',PRE_APP_LINK],cwd=REPO,env=self.env,capture_output=True,check=True)
        with tarfile.open(fileobj=io.BytesIO(c.stdout)) as archive:
            archive.extractall(source,filter='data')
        old=self.root/'pre-app-link-driver'
        self.command(['go','test','-c','-o',old,'./test/notification-driver'],cwd=source)
        show=['--config',f/'driver.json','--','workflow','notification','show',first['notification_id'],'--route',f/'config/route.json','--format','json']
        rolled=json.loads(self.command([old,*show],expected=(1,)).stdout)
        assert rolled['state']=='unknown' and rolled['reasons'][0]['code']=='state_unavailable'
        assert json.loads(self.command([self.driver,*show]).stdout)['state']=='transport_acknowledged' and len(self.posts(f))==2
        self.checks.append('owner app-link file selects v6 for new short events only; plain_text body plus Ply-built link block; v4 and v6 replays ignore later owner-file changes; pre-v6 CLI fails closed on v6 state without sending')

    def run(self):
        self.command([sys.executable,REPO/'test/agent_notify_regression_test.py'])
        self.setup();self.journeys();self.failures();self.incomplete_attempt_record();self.concurrency();self.unsafe_setup();self.compact_journeys();self.short_journeys();self.frozen_timezone();self.stored_history_and_helper_upgrade();self.app_link_journey();self.privacy()
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
        put(self.root/'result.json',{'kind':'PlyAgentNotificationAcceptance@1','outcome':'passed','checks':self.checks,'commands':'commands.json','root':str(self.root),'live_network':False,'helper_baseline':HELPER_BASELINE})
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
