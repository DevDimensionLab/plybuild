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
BASE = 'fd864af8b2bbe09726f85297b169f43dca958599'
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

    def notify(self, folder, extra=(), expected=(0,), env=None, **kwargs):
        c = self.command(self.notify_args(folder, extra, **kwargs),env=env,expected=expected)
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

    def v1_state(self):
        # Materialize and compile the exact prior candidate as an isolated fixture.
        baseline=self.root/'baseline-source';baseline.mkdir(mode=0o700)
        c=subprocess.run(['git','archive',BASE],cwd=REPO,env=self.env,capture_output=True,check=True)
        with tarfile.open(fileobj=io.BytesIO(c.stdout)) as archive:
            archive.extractall(baseline,filter='data')
        old=self.root/'baseline-driver'
        self.command(['go','test','-c','-o',old,'./test/notification-driver'],cwd=baseline)
        for response in [200,403]:
            f=self.fixture('legacy-'+str(response),response={'kind':'http','status':response,'body':'ok' if response==200 else 'action_prohibited'})
            source={'kind':'external','activity':'fixture/activity','run':'run-1','worktree':str(f/'work'), 'handoff':loc(f/'handoff.md'),'start_receipt':loc(f/'start.json',True),'report':loc(f/'report.json',True)}
            req={'kind':'ply.workflow.notification-request','schema_version':1,'route':'ply-log','source':source,
                 'gate':{'id':'gate-1','revision':1,'opened_at':'2026-10-04T14:00:00Z','reason':'result_control','state':'waiting_for_human'},
                 'sender':{'actor_claim':'legacy-fixture'},'public':{'task_title':'Legacy report','next_action':'Start result control.'}}
            put(f/'v1.json',req)
            route=f/'config/route.json';baseargs=['--config',f/'driver.json','--','workflow','notification']
            send=['send','--file',f/'v1.json','--route',route,'--format','json']
            p=json.loads(self.command([old,*baseargs,*send,'--check'],env=self.secret_env()).stdout)
            self.command([old,*baseargs,*send,'--apply','--confirm',p['confirmation']],env=self.secret_env(),expected=(0,5))
            state=(f/'state/state.json').read_bytes()
            r=json.loads(self.command([self.driver,*baseargs,'show',p['notification_id'],'--route',route,'--format','json']).stdout)
            assert r['state']==('transport_acknowledged' if response==200 else 'rejected')
            self.command([self.driver,*baseargs,*send,'--apply','--confirm',p['confirmation']],env=self.secret_env())
            assert len(self.posts(f))==1 and (f/'state/state.json').read_bytes()==state
            if response==403:
                retry=['retry',p['notification_id'],'--route',route,'--format','json']
                p2=json.loads(self.command([self.driver,*baseargs,*retry,'--check'],env=self.secret_env()).stdout)
                self.command([self.driver,*baseargs,*retry,'--apply','--confirm',p2['confirmation']],env=self.secret_env(),expected=(5,))
                assert len(self.posts(f))==2
                self.command([self.driver,*baseargs,*retry,'--apply','--confirm',p2['confirmation']],env=self.secret_env())
                assert len(self.posts(f))==2
        self.checks.append('actual v1 baseline state reads byte-for-byte, replays and explicitly retries under the new binary')

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
        self.setup();self.journeys();self.failures();self.incomplete_attempt_record();self.concurrency();self.unsafe_setup();self.v1_state();self.privacy()
        # Production executable help and preview are exercised without transport.
        for parts in [[],['workflow'],['workflow','notification'],['workflow','notification','send'],['workflow','notification','show'],['workflow','notification','retry']]:
            self.command([self.ply,*parts,'--help'])
        f=self.fixture('production-preview')
        # Create/preserve a request through the real helper, then production check.
        r=self.notify(f,['--credential-file',str(f/'config/missing')],expected=(2,))
        c=self.command([self.ply,'workflow','notification','send','--file',r['request'],'--route',f/'config/route.json','--check','--format','json'],env=self.secret_env())
        assert json.loads(c.stdout)['allowed'] is True and not self.posts(f) and not (f/'state/state.json').exists()
        self.privacy()
        put(self.root/'result.json',{'kind':'PlyAgentNotificationAcceptance@1','outcome':'passed','checks':self.checks,'commands':'commands.json','root':str(self.root),'live_network':False,'v1_baseline':BASE})
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
