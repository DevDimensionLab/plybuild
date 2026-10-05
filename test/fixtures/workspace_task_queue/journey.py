"""Isolated two-Task CLI journey. Every human claim is synthetic fixture data."""
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys

root = Path(sys.argv[1])
workspace = root / 'workspace'
repo = workspace / 'ply/main'
epic = workspace / 'ply/epic'
repo.mkdir(parents=True)
binary = root / 'ply'
helper = root / 'lifecycle-helper'
sequence = 0

def command(argv, cwd=workspace, expected=0, env=None):
    global sequence
    sequence += 1
    result = subprocess.run([str(v) for v in argv], cwd=cwd, env=env, capture_output=True, text=True)
    label = root / f'{sequence:03d}'
    label.with_suffix('.stdout').write_text(result.stdout)
    label.with_suffix('.stderr').write_text(result.stderr)
    label.with_suffix('.json').write_text(json.dumps({'argv': [str(v) for v in argv], 'cwd': str(cwd), 'exit': result.returncode}))
    assert (result.returncode == 0) == (expected == 0), (argv, result.returncode, result.stdout, result.stderr, root)
    return result.stdout

def ply(*args, cwd=workspace, expected=0):
    return command([binary, *args], cwd, expected)

def js(*args, cwd=workspace):
    return json.loads(ply(*args, '--format', 'json', cwd=cwd))

def git(path, *args):
    return command(['git', '-C', path, '-c', 'user.name=Ply tests', '-c', 'user.email=tests@example.invalid', *args]).strip()

def write(name, value):
    p = workspace / name
    p.write_text(json.dumps(value))
    return p

def generate(mode, name, *args):
    p = workspace / name
    command([helper, mode, p, *args])
    return p

def read(p):
    return json.loads(Path(p).read_text())

def revision(out):
    r = out['outcome_ref']
    return {k: r[k] for k in ['revision', 'manifest_sha256']}

def decision(out):
    r = out['outcome_ref']
    return {k: r[k] for k in ['id', 'manifest_sha256']}

def state_bytes():
    files = {str(p.relative_to(workspace)): hashlib.sha256(p.read_bytes()).hexdigest()
             for p in (workspace / '.ply').rglob('*') if p.is_file()}
    return files, git(repo, 'show-ref'), git(repo, 'worktree', 'list', '--porcelain')

def contents():
    return {str(p): hashlib.sha256(p.read_bytes()).hexdigest()
            for p in (workspace / '.ply/task-content/v1').rglob('*') if p.is_file()}

ply('workspace', 'init')
git(repo, 'init', '-b', 'main')
git(repo, 'commit', '--allow-empty', '-m', 'base')
b1 = git(repo, 'rev-parse', 'HEAD')
git(repo, 'worktree', 'add', '-b', 'epic', epic, b1)
ply('workspace', 'project', 'add', 'ply', '--name', 'Ply', '--wrapper', workspace / 'ply', '--repo', f'ply={repo}')
ply('workspace', 'epic', 'adopt', 'epic', '--title', 'Epic', '--project', 'ply', '--repo', 'ply', '--worktree', epic, '--ref', 'refs/heads/epic', '--expected-oid', b1)
assert git(epic, 'reflog', 'show', '--format=%H', 'refs/heads/epic').splitlines() == [b1]
for task_id in ['task', 'task-b']:
    ply('workspace', 'task', 'create', task_id, '--title', 'Same title', '--description', 'Synthetic H2 task', '--epic', 'epic', '--project', 'ply', '--repo', 'ply')
target = ['--project', 'ply', '--repo', 'ply', '--epic', 'epic']
old = state_bytes()
empty = js('workspace', 'task', 'queue', 'list', *target)
assert empty['revision'] == 0 and empty['pending'] == []
assert old == state_bytes(), 'empty queue read wrote state'

# Build drafts from exact public readback. The helper's transient input is not a stored record.
def select(task_id, rev, old_selection=None, old_spec=None):
    show = js('workspace', 'task', 'show', task_id)
    base = js('workspace', 'epic', 'base', 'show', 'epic', '--repo', 'ply')['versions'][-1]
    source = json.loads(json.dumps(show))
    source['resource']['persisted']['worktree'] = {
        'parent_worktree_id': base['worktree_id'], 'parent_ref': base['ref'],
        'parent_oid': base['oid'], 'parent_tree': base['tree']}
    source['resource']['observed']['target'] = {'oid': base['oid'], 'tree': base['tree']}
    show_path = write(f'{task_id}-source.json', source)
    path = generate('spec', f'{task_id}-spec-{rev}.json', show_path, workspace, rev, '-')
    draft = read(path)
    draft['task_id'] = task_id
    draft['publication_key'] = f'{task_id}/spec/{rev}'
    draft['expected_previous'] = old_spec
    if old_spec:
        for doc in draft['documents']:
            doc['source'] = {'kind': 'snapshot', 'manifest_sha256': old_spec['manifest_sha256'], 'document_id': doc['id']}
    path = write(path.name, draft)
    spec = js('workspace', 'task', 'spec', 'record', task_id, '--file', path)
    spec_path = write(f'{task_id}-record.json', spec)
    path = generate('assessment', f'{task_id}-assessment.json', spec_path)
    draft = read(path); draft['task_id'] = task_id; draft['publication_key'] = f'{task_id}/assessment/{rev}'
    path = write(path.name, draft)
    assessment = js('workspace', 'task', 'spec', 'assess', task_id, '--file', path)
    assessment_path = write(f'{task_id}-assess-result.json', assessment)
    path = generate('selection', f'{task_id}-selection.json', spec_path, assessment_path, show_path)
    draft = read(path); draft['task_id'] = task_id; draft['publication_key'] = f'{task_id}/selection/{rev}'; draft['expected_previous_selection'] = old_selection
    path = write(path.name, draft)
    selected = js('workspace', 'task', 'spec', 'select', task_id, '--file', path)
    return revision(spec), decision(selected)

spec_a, selection_a = select('task', 1)
spec_b, selection_b = select('task-b', 1)

def set_queue(entries, key):
    q = js('workspace', 'task', 'queue', 'list', *target)
    draft = {'kind': 'WorkspaceTaskQueueDraft@1', 'schema_version': 1, 'publication_key': key,
             'project_id': 'ply', 'repo_id': 'ply', 'epic_id': 'epic', 'expected_revision': q['revision'],
             'entries': entries, 'registry_upgrade': q['registry']['required_upgrade'],
             'human_decision': {'actor_claim': 'synthetic fixture, not human QA', 'decided_at_utc': '2026-09-30T00:00:00Z',
                                'source': 'explicit_human_instruction', 'statement': 'Synthetic queue priority.'}}
    return js('workspace', 'task', 'queue', 'set', '--file', write('queue.json', draft))

q = set_queue([{'task_id': 'task', 'selection': selection_a}, {'task_id': 'task-b', 'selection': selection_b}], 'queue/initial')
assert [v['state'] for v in q['pending']] == ['ready', 'ready']
preview = js('workspace', 'task', 'prepare', '--next', cwd=epic)
assert preview['plan']['task_id'] == 'task'
old = state_bytes()
assert js('workspace', 'task', 'prepare', '--next', '--check', cwd=epic)['confirmation'] == preview['confirmation']
assert old == state_bytes(), 'preview wrote state'
prepared_a = js('workspace', 'task', 'prepare', '--next', *target, '--apply', '--confirm', preview['confirmation'])
assert prepared_a['state'] == 'prepared'
task = Path(prepared_a['plan']['worktree_path'])
task_ref = 'refs/heads/' + prepared_a['plan']['branch']
assert git(task, 'rev-parse', 'HEAD') == b1 and git(epic, 'rev-parse', 'HEAD') == b1
for item in prepared_a['plan']['required_inputs']:
    data = Path(item['locator']).read_bytes()
    assert item['sha256'] == 'sha256:' + hashlib.sha256(data).hexdigest() and item['size_bytes'] == len(data)
selected = js('workspace', 'task', 'spec', 'show', 'task', '--spec', 'solution', '--revision', '1')
selected_path = write('selected.json', selected)
(workspace / 'input.txt').write_text('immutable fixture input\n')
handoff_draft = generate('handoff', 'handoff-draft.json', task, task_ref, b1, workspace / 'input.txt', selected_path)
text = ply('workflow', 'handoff', 'create', '--file', handoff_draft)
handoff = next(line.removeprefix('Handoff: ') for line in text.splitlines() if line.startswith('Handoff: '))
start_draft = generate('start', 'start-draft.json', handoff, root)
text = ply('workflow', 'handoff', 'submit-start', '--handoff', handoff, '--file', start_draft, cwd=task)
start = next(line.removeprefix('Start receipt: ') for line in text.splitlines() if line.startswith('Start receipt: '))
(task / 'delivery.txt').write_text('delivered\n')
git(task, 'add', 'delivery.txt'); git(task, 'commit', '-m', 'H2 synthetic delivery')
b2 = git(task, 'rev-parse', 'HEAD'); tree2 = git(task, 'rev-parse', 'HEAD^{tree}')
os.environ['WS06_RESULT_OID'] = b2; os.environ['WS06_RESULT_TREE'] = tree2
reply = Path(read(handoff)['reply_capability']['reply_root'])
stdout = reply / 'staging/verifier.stdout'; stdout.write_text(git(task, 'status', '--short'))
requirements = reply / 'staging/task-requirements.json'
command([helper, 'requirements', requirements, handoff])
terminal_draft = generate('terminal', 'terminal-draft.json', handoff, start, stdout, requirements)
text = ply('workflow', 'handoff', 'submit-result', '--handoff', handoff, '--file', terminal_draft, cwd=task)
terminal = next(line.removeprefix('Terminal result: ') for line in text.splitlines() if line.startswith('Terminal result: '))
inspection = js('workflow', 'handoff', 'inspect', '--handoff', handoff)
inspection_path = write('inspection.json', inspection)
show_a = js('workspace', 'task', 'show', 'task'); show_path = write('show-a.json', show_a)
result_draft = generate('task-result', 'result.json', show_path, handoff, start, terminal, inspection_path)
result = js('workspace', 'task', 'result', 'record', 'task', '--file', result_draft)
result_path = write('result-readback.json', result)
(workspace / 'qa.txt').write_text('Synthetic QA fixture only; no actual human QA.\n')
qa_draft = generate('qa', 'qa.json', result_path, workspace / 'qa.txt')
qa = js('workspace', 'task', 'qa', 'record', 'task', '--file', qa_draft)
args = ['workspace', 'task', 'integrate', 'task', '--result', result['record']['id'], '--qa', qa['record']['id'], '--expected-result-oid', b2, '--expected-parent-oid', b1]
before_check = state_bytes()
plan = js(*args, '--check')
assert state_bytes() == before_check, 'integration check changed state'
integrated = js(*args, '--apply', '--confirm', plan['plan']['sha256'])
assert integrated['classification'] == 'exact_effect' and git(epic, 'rev-parse', 'HEAD') == b2
after_apply = state_bytes()
assert js(*args, '--apply', '--confirm', plan['plan']['sha256'])['classification'] == 'exact_effect'
assert state_bytes() == after_apply, 'identical integration retry changed state'
assert git(epic, 'reflog', 'show', '--format=%H', 'refs/heads/epic').splitlines() == [b2, b1]
preserved = contents()
persisted_a = js('workspace', 'task', 'show', 'task')['resource']['persisted']
q = js('workspace', 'task', 'queue', 'list', *target)
js('workspace', 'task', 'queue', 'advance', *target, '--preparation', prepared_a['preparation']['id'], '--expected-revision', str(q['revision']), '--reason', 'Synthetic explicit continuation')
assert js('workspace', 'task', 'prepare', '--next', *target)['confirmation'] is None
base_plan = js('workspace', 'epic', 'base', 'update', 'epic', '--repo', 'ply', '--check')
base = js('workspace', 'epic', 'base', 'update', 'epic', '--repo', 'ply', '--apply', '--confirm', base_plan['confirmation'])
assert base['state'] == 'committed' and base['versions'][-1]['revision'] == 2
historical_a = js('workspace', 'task', 'show', 'task')
assert historical_a['freshness']['selection'] == 'current'
assert historical_a['freshness']['target'] == 'stale'
assert historical_a['results'][0]['basis_status'] == 'bound'
assert historical_a['results'][0]['relevance'] == 'stale'
assert js('workspace', 'task', 'prepare', '--next', *target)['confirmation'] is None
assert js('workspace', 'task', 'show', 'task')['resource']['persisted'] == persisted_a
assert js('workspace', 'task', 'result', 'record', 'task', '--file', result_draft)['record'] == result['record']
assert js('workspace', 'task', 'qa', 'record', 'task', '--file', qa_draft)['record'] == qa['record']
assert js(*args, '--apply', '--confirm', plan['plan']['sha256'])['classification'] == 'exact_effect'
assert js(*args, '--check')['task_spec_relevance']['relevance'] == 'stale'
assert js('workspace', 'task', 'spec', 'record', 'task', '--file', workspace / 'task-spec-1.json')['outcome_ref']['manifest_sha256'] == spec_a['manifest_sha256']
_, new_b = select('task-b', 2, selection_b, spec_b)
set_queue([{'task_id': 'task-b', 'selection': new_b}], 'queue/new-base')
preview_b = js('workspace', 'task', 'prepare', 'task-b', *target)
prepared_b = js('workspace', 'task', 'prepare', 'task-b', *target, '--apply', '--confirm', preview_b['confirmation'])
assert prepared_b['state'] == 'prepared' and prepared_b['plan']['parent_oid'] == b2
assert prepared_b['plan']['worktree_path'] != str(task)
assert git(task, 'rev-parse', 'HEAD') == b2, 'A worktree was moved by base registration'
(Path(prepared_b['plan']['worktree_path']) / 'work.txt').write_text('normal ongoing work\n')
assert js('workspace', 'task', 'preparation', 'show', prepared_b['preparation']['id'])['state'] == 'prepared'
assert js('workspace', 'task', 'preparation', 'show', prepared_a['preparation']['id'])['plan']['parent_oid'] == b1
assert js('workspace', 'task', 'show', 'task')['resource']['persisted'] == persisted_a
assert all(Path(p).exists() and hashlib.sha256(Path(p).read_bytes()).hexdigest() == h for p, h in preserved.items()), 'historical content bytes changed'
assert git(repo, 'rev-parse', 'refs/heads/main') == b1
(root / 'H2-result.json').write_text(json.dumps({'scenario': 'H2', 'status': 'passed', 'base_1': b1, 'base_2': b2, 'preparation_a': prepared_a, 'preparation_b': prepared_b, 'preserved_content': preserved}, indent=2))
print('H2: passed; all human claims were synthetic fixture data')
