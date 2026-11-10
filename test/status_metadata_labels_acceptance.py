#!/usr/bin/env python3
"""T03079C2A: isolated staged-CLI semantics, monotonic timings and child cleanup.

Run with the test-only jsonschema dependency and the exact private candidate.
This does not measure or qualify the downstream HTTP surface or claim general
standalone descendant cleanup. External-group samples explicitly supply caller
killpg; r001-f01 remains a separately reproduced open baseline limitation.
"""
import concurrent.futures
import json
import os
from pathlib import Path
import signal
import subprocess
import sys
import tempfile
import time
import unicodedata

from jsonschema import Draft202012Validator

BINARY = str(Path(sys.argv[1]).resolve(strict=True))
SCHEMA_ROOT = Path(__file__).resolve().parents[1] / "schemas/agent-status-metadata-v1"
schema = json.loads((SCHEMA_ROOT / "schema.json").read_text())
Draft202012Validator.check_schema(schema)
validator = Draft202012Validator(schema)
samples = []


def validate(value):
    validator.validate(value)
    for agent in value["agents"]:
        for text in agent.values():
            assert len(text.encode("utf-16-le")) // 2 <= 256
        label = agent.get("tab_label")
        if label is not None:
            assert label and label.strip() == label
            assert not any(unicodedata.category(c) in ("Cc", "Cf", "Cs") or c in "\u2028\u2029\ufffd/\\@" for c in label)
            assert not any(m in label.lower() for m in ("session_id", "terminal_id", "cwd=", "token=", "secret=", "password=", "://"))


for path in sorted((SCHEMA_ROOT / "examples").glob("*.json")):
    if not path.name.startswith("source-"):
        validate(json.loads(path.read_text()))
for path in sorted((SCHEMA_ROOT / "negative").glob("*.json")):
    try:
        validate(json.loads(path.read_text()))
    except Exception:
        pass
    else:
        raise AssertionError(f"negative fixture accepted: {path.name}")

SHIM = r'''import json, os, pathlib, subprocess, sys, time
root = pathlib.Path(__file__).parent
args = sys.argv[1:]
with (root/'calls').open('a') as f: f.write(json.dumps(args)+'\n')
if args == ['agent', 'list']: name='agents'
elif len(args)==4 and args[:3]==['tab','list','--workspace']: name=args[3]
else: sys.exit(91)
(root/('process-'+str(os.getpid()))).write_text(json.dumps(dict(pid=os.getpid(), pgid=os.getpgrp(), name=name)))
config=json.loads((root/'config').read_text())
mode=config.get('list' if name=='agents' else 'label','ok')
if name=='agents': time.sleep(config.get('list_delay',0))
if mode=='hang': time.sleep(20)
if mode=='exit': print('synthetic-private-diagnostic',file=sys.stderr); sys.exit(7)
if mode=='malformed': print('{'); sys.exit(0)
if mode=='oversize': print(' ' * ((8<<20)+1 if name=='agents' else (256<<10)+1)); sys.exit(0)
if mode=='stderr': print('x' * ((8<<20)+1 if name=='agents' else (256<<10)+1),file=sys.stderr)
if mode in ('pipes','tree-hang'):
    p=subprocess.Popen(['/bin/sleep','20'],stdout=sys.stdout,stderr=sys.stderr)
    (root/('child-'+str(p.pid))).write_text(str(p.pid))
    if mode=='tree-hang': time.sleep(20)
sys.stdout.write((root/name).read_text()); sys.stdout.flush()
if mode=='slow-exit': time.sleep(20)
'''

with tempfile.TemporaryDirectory(prefix="T03079C2A-cli-") as directory:
    root = Path(directory)
    shim = root / "herdr"
    shim.write_text('#!' + sys.executable + '\n' + SHIM)
    shim.chmod(0o700)
    env = dict(os.environ, PATH=str(root) + os.pathsep + os.environ["PATH"], HOME=str(root), TMPDIR=str(root))
    argv = [BINARY, 'workflow', 'agents', '--metadata-only', '--format', 'json', '--workspace', str(root/'absent')]
    expected = []

    def setup(count, workspaces):
        global expected
        rows, labels, expected = [], {}, []
        for n in range(count):
            w = f'w{n % workspaces:02d}'
            tab, pane, label = f't{n}', f'p{n:04d}', f'Name {n} 🌱'
            provider, status = ('codex', 'idle') if n % 2 else ('claude', 'working')
            rows.append(dict(agent=provider, agent_status=status, pane_id=pane, tab_id=tab, workspace_id=w, agent_session=False, terminal_id={'private': 'never-export'}))
            labels.setdefault(w, []).append(dict(workspace_id=w, tab_id=tab, label=label))
            expected.append(dict(pane=pane, tab=tab, workspace=w, status=status, provider=provider, tab_label=label))
        (root/'agents').write_text(json.dumps(dict(result=dict(type='agent_list', agents=rows))))
        for w, tabs in labels.items():
            (root/w).write_text(json.dumps(dict(result=dict(tabs=tabs)), ensure_ascii=False, separators=(',', ':')))
        (root/'config').write_text('{}')

    def read(name, healthy=True, list_failure=False, cancel=False):
        start = time.monotonic()
        proc = subprocess.Popen(argv, cwd=root, env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        if cancel:
            # Wait only for our isolated fixture's startup, never inspect real agents.
            until = time.monotonic() + 1
            while not (root/'calls').exists() and time.monotonic() < until:
                time.sleep(.005)
            proc.send_signal(signal.SIGTERM)
        stdout, stderr = proc.communicate(timeout=3)
        assert proc.returncode == 0 and not stderr, (name, proc.returncode, stderr)
        value = json.loads(stdout)
        validate(value)
        source = value['sources'][0]
        if healthy:
            assert source['status'] == ('available' if expected else 'empty') and source['coverage'] == 'complete', (name, source)
            assert value['agents'] == expected, (name, len(value['agents']), len(expected))
        elif list_failure or cancel:
            assert source['status'] == 'unavailable' and not value['agents'], (name, value)
        else:
            without_labels = [{k: v for k, v in a.items() if k != 'tab_label'} for a in expected]
            assert source['coverage'] == 'complete' and value['agents'] == without_labels, (name, value)
        # Only isolated children created by this test are observed. Cleanup must
        # have completed when the full response and semantic validation return.
        for path in root.glob('child-*'):
            pid = int(path.read_text())
            try:
                os.kill(pid, 0)
            except ProcessLookupError:
                pass
            else:
                raise AssertionError((name, 'owned child still alive', pid))
        elapsed = time.monotonic() - start
        samples.append(dict(name=name, elapsed_ms=elapsed*1000, agents=len(value['agents']), expected_labels=len(expected) if healthy else 0, source=source['status']))
        print(json.dumps(samples[-1]), flush=True)
        assert elapsed <= 2.6, (name, elapsed)
        return value

    for count, workspaces in [(9, 3), (4096, 4), (4096, 32)]:
        setup(count, workspaces)
        for iteration in range(3):
            read(f'healthy-{count}-{workspaces}-sample-{iteration}')
        with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
            list(pool.map(lambda n: read(f'four-way-{count}-{workspaces}-{n}'), range(4)))
    setup(9, 3)
    # Beyond the four supported reads: no new process-wide scheduler/cache is
    # introduced. Observe bounded per-request worker recovery at eight callers.
    with concurrent.futures.ThreadPoolExecutor(max_workers=8) as pool:
        list(pool.map(lambda n: read(f'eight-way-stress-{n}'), range(8)))
    for target in ('list', 'label'):
        for mode in ('exit', 'hang', 'slow-exit', 'pipes', 'malformed', 'oversize', 'stderr'):
            for path in root.glob('child-*'): path.unlink()
            (root/'config').write_text(json.dumps({target: mode}))
            read(f'{target}-{mode}', healthy=False, list_failure=target=='list')
            (root/'config').write_text('{}')
            read(f'recovery-{target}-{mode}')
    (root/'config').write_text(json.dumps(dict(list_delay=1.3, label='hang')))
    read('shared-list-label-budget', healthy=False)
    (root/'calls').unlink()
    (root/'config').write_text(json.dumps(dict(list='hang')))
    read('cancel-list', healthy=False, cancel=True)
    (root/'config').write_text('{}')
    read('recovery-cancel')
    # Model an explicit external cleanup owner, not an inferred one. This
    # fixture supplies killpg after exit; start_new_session alone guarantees no
    # cleanup. The separate r001-f01 baseline reproducer omits caller killpg and
    # remains red for unmanaged group-leader timeout. No standalone pass here.
    for target in ('list', 'label'):
        for stop in ('normal-timeout', 'term', 'hard-kill'):
            for pattern in ('process-*', 'child-*'):
                for path in root.glob(pattern): path.unlink()
            (root/'config').write_text(json.dumps({target: 'tree-hang'}))
            start = time.monotonic()
            proc = subprocess.Popen(argv, cwd=root, env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE, start_new_session=True)
            deadline = time.monotonic()+1.5
            owned = []
            while time.monotonic()<deadline:
                owned = [json.loads(p.read_text()) for p in root.glob('process-*') if p.stat().st_size]
                active = [p for p in owned if (p['name']=='agents') == (target=='list')]
                if active and list(root.glob('child-*')): break
                time.sleep(.005)
            assert active, ('fixture never started', target, stop)
            owned += [dict(pid=int(p.read_text()), pgid=os.getpgid(int(p.read_text())), name='grandchild') for p in root.glob('child-*')]
            assert all(p['pgid']==proc.pid for p in owned), ('child escaped external group', owned, proc.pid)
            if stop=='term': proc.send_signal(signal.SIGTERM)
            if stop=='hard-kill': proc.kill()  # Deliberately kill only Ply first.
            stdout, stderr = proc.communicate(timeout=2.6)
            try: os.killpg(proc.pid, signal.SIGKILL)
            except ProcessLookupError: pass
            # Include workers that started between the startup observation and
            # cancellation, plus every fixture grandchild, before checking cleanup.
            owned = [json.loads(p.read_text()) for p in root.glob('process-*') if p.stat().st_size]
            owned += [dict(pid=int(p.read_text())) for p in root.glob('child-*')]
            # Reap/observe only the known isolated fixture processes.
            until=time.monotonic()+.2
            alive=[]
            while time.monotonic()<until:
                alive=[]
                for child in owned:
                    try: os.kill(child['pid'],0); alive.append(child['pid'])
                    except ProcessLookupError: pass
                if not alive: break
                time.sleep(.005)
            assert not alive, ('external cleanup leaked children', alive)
            elapsed=time.monotonic()-start
            assert elapsed<=2.6, ('external cleanup bound',elapsed)
            if stop!='hard-kill':
                assert proc.returncode==0 and not stderr
                validate(json.loads(stdout))
            sample=dict(name=f'external-group-{target}-{stop}',elapsed_ms=elapsed*1000,agents=None,expected_labels=0,source='cancelled' if stop!='normal-timeout' else 'timeout')
            samples.append(sample);print(json.dumps(sample),flush=True)
            (root/'config').write_text('{}')
            read(f'external-recovery-{target}-{stop}')
    setup(0, 1)
    read('genuine-empty')
    for line in (root/'calls').read_text().splitlines():
        args = json.loads(line)
        assert args == ['agent', 'list'] or (len(args)==4 and args[:3]==['tab','list','--workspace'])
    assert not (root/'.ply').exists(), 'native state created'

print(json.dumps(dict(task_code='T03079C2A', samples=samples, max_elapsed_ms=max(s['elapsed_ms'] for s in samples), whole_http_3s='not_run; root/server owner responsibility', cleanup_scope='isolated caller-owned group cleanup plus nonleader cleanup; no standalone group-leader pass', open_baseline_limitation='T03079C2A/r001-f01: unmanaged group-leader descendants may survive normal timeout; tracked by Ply-plan')), flush=True)
