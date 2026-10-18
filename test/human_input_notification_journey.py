#!/usr/bin/env python3
"""HI03: model the installed skills with real CLI data and offline notification I/O.

This is a cooperative agent/app-role fixture, not a live phone or passkey test.
All writes are under a retained physical /private/tmp directory. No live transport,
user credential lookup, app code, global installation or existing Q is used.
"""
import argparse
import copy
from datetime import datetime
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import tempfile

sys.dont_write_bytecode = True
REPO = Path(__file__).resolve().parents[1]
ENV = 'PLY_SLACK_WEBHOOK_URL'
BASE = 'd294c6b352d75748d5d24dd82f97c94bbcb1259c'


def put(path, value):
    path.write_bytes(value if isinstance(value, bytes) else
                     (json.dumps(value, ensure_ascii=False, sort_keys=True)+'\n').encode())
    path.chmod(0o600)
    return path


def examples(skill):
    return [json.loads(s) for s in re.findall(r'```json\n(.*?)\n\s*```', skill.read_text(), re.S)]


class Journey:
    def __init__(self, ply, driver):
        self.root = Path(tempfile.mkdtemp(prefix='hi03-journey-', dir='/private/tmp')).resolve()
        self.ply, self.driver = str(ply), str(driver)
        self.env = dict(os.environ)
        self.env.pop(ENV, None)
        for key in ['HOME', 'TMPDIR']:
            path = self.root/key.lower(); path.mkdir(mode=0o700)
            self.env[key] = str(path)
        self.env['PYTHONDONTWRITEBYTECODE'] = '1'
        self.repo = self.root/'ws/demo/main'; self.repo.mkdir(parents=True)
        self.ws = self.root/'ws'
        self.skill = self.root/'ply-agent-notify'
        self.human = self.root/'ply-human-input'
        self.checks = []

    def command(self, args, cwd=None, env=None, codes=(0,), structured=True):
        r = subprocess.run([str(a) for a in args], cwd=cwd or self.repo,
                           env=env or self.env, capture_output=True, timeout=60)
        assert r.returncode in codes, (args, r.returncode, r.stdout.decode(), r.stderr.decode())
        return json.loads(r.stdout) if structured else r.stdout

    def cli(self, *args):
        return self.command([self.ply, *args, '--format', 'json'])

    def input(self, *args):
        return self.cli('workflow', 'input', *args, '--workspace', self.ws)

    def setup(self):
        self.command(['git', 'init', '-q', '-b', 'main'], structured=False)
        self.command(['git', '-c', 'user.name=HI03', '-c', 'user.email=hi03@example.invalid',
                      'commit', '-q', '--allow-empty', '-m', 'base'], structured=False)
        self.command([self.ply, 'workspace', 'init'], cwd=self.ws, structured=False)
        self.command([self.ply, 'workspace', 'project', 'add', 'demo', '--name', 'Demo',
                      '--wrapper', 'demo', '--repo', 'demo=demo/main'], cwd=self.ws, structured=False)
        self.command([sys.executable, REPO/'skills/ply-agent-notify/scripts/install.py',
                      '--destination', self.skill])
        self.command([REPO/'scripts/install-human-input-skill', '--destination', self.human], structured=False)
        self.feedback = next(x for x in examples(self.human/'SKILL.md') if x.get('event_type') == 'feedback_required')
        self.done = next(x for x in examples(self.skill/'SKILL.md') if x.get('event_type') == 'agent_finished')
        for folder, source in [(self.skill, REPO/'skills/ply-agent-notify'), (self.human, REPO/'skills/ply-human-input')]:
            for p in folder.rglob('*'):
                if p.is_file(): assert p.read_bytes() == (source/p.relative_to(folder)).read_bytes()

    def question(self, topic):
        draft = examples(self.human/'SKILL.md')[0]
        draft.update(publication_key='hi03/'+topic, project_id='demo', repo_id='demo',
                     activity_id='hi03', provider='codex', origin_cwd=str(self.repo),
                     text='PRIVATE_QUESTION '+('Details. '*30)+'Which note should be written?',
                     reason='PRIVATE_REASON: the scope remains undecided.')
        path = put(self.root/(topic+'-draft.json'), draft)
        return self.input('ask', '--file', path), path

    def fixture(self, q, name, unknown=False):
        f = self.root/name; f.mkdir(mode=0o700)
        (f/'state').mkdir(mode=0o700)
        route = put(f/'route.json', {'kind':'ply.workflow.notification-route','schema_version':1,
                    'name':'ply-log','transport':'slack_incoming_webhook','channel_label':'ply-log',
                    'webhook_env':ENV,'state_root':str(f/'state')})
        handoff = put(f/'handoff.md', b'Offline authorized fixture mandate. PRIVATE_MANDATE\n')
        task = {'kind':'PlyAgentNotificationTask@2','activity':'notification/job', 'run':'run-1',
                'worktree':str(self.repo),'handoff':str(handoff),'route':str(route),
                'event_root':str(f/'events'),'actor_claim':'codex-fixture','notification_authorized':True,
                'provider':'codex','origin_cwd':str(self.repo),'context_root':str(self.ws),'timezone':'Europe/Oslo'}
        put(f/'task.json', task)
        event = copy.deepcopy(self.feedback)
        event.update(event_id='human-input/'+q['question_id'], occurred_at=datetime.fromisoformat(q['created_at_utc'].replace('Z','+00:00')).strftime('%Y-%m-%dT%H:%M:%SZ'))
        put(f/'event.json', event)
        put(f/'driver.json', {'kind':'PlyNotificationAcceptanceCase@1','artifacts_root':str(self.root),
             'now':q['created_at_utc'],'post_log':str(f/'posts.jsonl'),
             'transport':{'kind':'unknown'} if unknown else {'kind':'http','status':200,'body':'ok'},'fault':None})
        # The only notification executable is the existing network-free production graph.
        wrapper = ('#!'+sys.executable+'\nimport os,sys\n'
                   'os.execv('+repr(self.driver)+', ['+repr(self.driver)+', "--config", '+
                   repr(str(f/'driver.json'))+', "--", *sys.argv[1:]])\n')
        (f/'ply-offline').write_text(wrapper); (f/'ply-offline').chmod(0o700)
        return f

    def notify(self, f, trigger=None, skill=None, report=None, authorized=True):
        if not authorized: return {'state':'not_attempted'}  # Skill's missing-binding branch.
        args = [sys.executable, (skill or self.skill)/'scripts/notify.py', '--task', f/'task.json',
                '--event', f/'event.json', '--ply-bin', f/'ply-offline', '--credential-file', f/'absent']
        if trigger: args += ['--trigger', trigger]
        if report: args += ['--report', report]
        # Synthetic fixture string only; this executable cannot perform a real POST.
        env = {**self.env, ENV:'https://hooks.slack.com/services/Tfixture/Bfixture/offlineonly'}
        return self.command(args, env=env, codes=(0, 2, 5))

    def pending_notify(self, q, f, **kwargs):
        # The modeled agent follows step Notify: only the exact still-waiting Q alerts.
        current = self.input('show', q['question_id'])
        if current['state'] != 'waiting_for_answer': return {'state':'not_needed'}
        return self.notify(f, **kwargs)

    def posts(self, f):
        p = f/'posts.jsonl'
        return p.read_bytes().splitlines() if p.exists() else []

    def answer(self, q):
        # Simulate app-role immutable answer APIs; passkey/mobile UI is owned by plyapp.
        preview = self.input('answer', 'preview', q['question_id'], '--request-id', 'app-'+q['question_id'], '--option-id', 'full')
        path = put(self.root/(q['question_id']+'-answer.json'), preview['request'])
        assert self.input('answer', 'apply', '--file', path)['state'] == 'recorded'

    def journey(self):
        q, draft = self.question('pending'); f = self.fixture(q, 'pending')
        assert self.input('ask', '--file', draft)['question_id'] == q['question_id']
        first = self.pending_notify(q, f)
        assert first['state'] == 'transport_acknowledged', first
        request = Path(first['request']); before = request.read_bytes(); event = (f/'event.json').read_bytes()
        posts = self.posts(f); assert len(posts) == 1
        text = json.loads(posts[0])['payload']['text']
        assert '⚠️' in text and '/plyapp' in text and 'Questions' in text and 'passkey' in text
        assert len(text.splitlines()) == 3 and len(text) <= 480
        assert all(s not in text for s in ['PRIVATE_', q['text'], q['reason'], 'Short usage note', 'Full reference'])
        assert json.loads(before)['source']['activity'] == 'notification/job' != q['scope']['activity_id']
        assert self.input('wait', q['question_id'], '--timeout-seconds', '1')['outcome'] == 'still_waiting'
        assert self.pending_notify(q, f, trigger='waiting_for_feedback')['state'] == 'suppressed'
        assert self.pending_notify(q, f)['state'] == 'transport_acknowledged'  # Fresh helper process after lost reply.
        assert self.posts(f) == posts and request.read_bytes() == before and (f/'event.json').read_bytes() == event
        self.answer(q)
        assert self.pending_notify(q, f)['state'] == 'not_needed'
        answer = self.input('wait', q['question_id'], '--timeout-seconds', '1')['question']['answer']
        flags = ['--request-id', answer['request_id'], '--request-sha256', answer['request_sha256']]
        assert self.input('acknowledge', q['question_id'], *flags)['state'] == 'agent_received'
        put(self.repo/'decision.txt', b'Full note.\n')  # Observable work selected by the stored answer.
        assert answer['option_id'] == 'full'
        assert self.input('continue', q['question_id'], *flags, '--summary', 'Wrote the full note.')['state'] == 'agent_reported_continued'
        assert self.input('ask', '--file', draft)['answer'] == answer
        assert self.posts(f) == posts
        self.checks.append('installed Q publish → one warning alert → restart/readback → same app-role answer/acknowledge/use')

        for state in ['answered','withdrawn']:
            q, _ = self.question(state); f = self.fixture(q, state)
            if state == 'answered': self.answer(q)
            else: self.input('withdraw', q['question_id'], '--reason', 'Decision became unnecessary.')
            assert self.pending_notify(q, f)['state'] == 'not_needed'
            assert not (f/'events').exists() and not self.posts(f)
        q, _ = self.question('unknown'); f = self.fixture(q, 'unknown', unknown=True)
        r = self.pending_notify(q, f); assert r['state'] == 'unknown', r
        request = Path(r['request']); before = request.read_bytes(); posts = self.posts(f)
        assert len(posts) == 1 and self.pending_notify(q, f)['state'] == 'unknown'
        assert self.posts(f) == posts and request.read_bytes() == before
        self.checks.append('answered/withdrawn skip alert; unknown reuses attempt and request bytes with no second POST')

        q, _ = self.question('unauthorized'); f = self.fixture(q, 'unauthorized')
        assert self.pending_notify(q, f, authorized=False)['state'] == 'not_attempted'
        assert not (f/'events').exists() and self.input('show', q['question_id'])['state'] == 'waiting_for_answer'
        task = json.loads((f/'task.json').read_bytes()); task['notification_authorized'] = False; put(f/'task.json', task)
        assert self.notify(f)['state'] == 'not_attempted' and not self.posts(f)
        self.checks.append('missing authority remains not_attempted with the local Q intact')

        q, _ = self.question('done'); f = self.fixture(q, 'done')
        event = copy.deepcopy(self.done); event.update(event_id='done-1', occurred_at=json.loads((f/'event.json').read_bytes())['occurred_at'])
        put(f/'event.json', event)
        report = put(f/'report.json', {'kind':'ply.workflow.delivery-report','summary':'PRIVATE_REPORT'})
        assert self.notify(f, report=report)['state'] == 'transport_acknowledged'
        text = json.loads(self.posts(f)[0])['payload']['text']
        assert 'ℹ️' in text and '⚠' not in text and 'PRIVATE_REPORT' not in text
        self.checks.append('installed done event has info indicator; native original source stays private')

    def legacy(self):
        old = self.root/'old-notify'; old.mkdir(mode=0o700)
        for name in ['SKILL.md','references/task-template.json','references/usage.md','scripts/common.py',
                     'scripts/notify.py','scripts/bootstrap.py','scripts/install.py']:
            p = old/name; p.parent.mkdir(mode=0o700, exist_ok=True)
            put(p, self.command(['git', 'show', BASE+':skills/ply-agent-notify/'+name], cwd=REPO, structured=False))
        for unknown in [False, True]:
            dest = self.root/('upgrade-notify-'+str(unknown))
            self.command([sys.executable, old/'scripts/install.py', '--destination', dest])
            topic = 'legacy-unknown' if unknown else 'legacy-acknowledged'
            q, _ = self.question(topic); f = self.fixture(q, topic, unknown)
            event = json.loads((f/'event.json').read_bytes())
            event.update(kind='PlyAgentNotificationEvent@2',activity='notification/job',run='run-1')
            event['public']['task_title'] = 'Original title without indicator'
            put(f/'event.json', event)
            r = self.pending_notify(q, f, skill=dest)
            assert r['state'] == ('unknown' if unknown else 'transport_acknowledged'), r
            request = Path(r['request']); before = request.read_bytes(); posts = self.posts(f)
            state_files = {p:p.read_bytes() for p in (f/'state').rglob('*') if p.is_file()}
            self.command([sys.executable, self.skill/'scripts/install.py', '--destination', dest,
                          '--upgrade-from', old,
                          '--backup', self.root/('backup-'+str(unknown))])
            assert self.pending_notify(q, f, skill=dest)['state'] == r['state']
            assert request.read_bytes() == before and self.posts(f) == posts
            assert all(p.read_bytes() == raw for p,raw in state_files.items())
        self.checks.append('v3 requests, old titles, payloads and acknowledged/unknown state bytes survive installed helper upgrade')

    def run(self):
        self.setup()
        before = {p:p.read_bytes() for p in (self.ws/'.ply').glob('*.yaml')}
        self.journey(); self.legacy()
        assert all(p.read_bytes() == raw for p,raw in before.items()), 'notification journey created Task/Run or altered registry'
        result = {'kind':'HI03OfflineJourney@1','state':'passed','root':str(self.root),
                  'live_transport':False,'physical_phone_or_passkey_evidence':False,'checks':self.checks}
        put(self.root/'result.json', result)
        print(json.dumps(result, ensure_ascii=False))


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--ply', required=True, type=Path)
    parser.add_argument('--driver', required=True, type=Path)
    args = parser.parse_args()
    Journey(args.ply.resolve(), args.driver.resolve()).run()
