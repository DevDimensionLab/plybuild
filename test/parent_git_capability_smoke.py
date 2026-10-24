#!/usr/bin/env python3
"""T0523E9D1: bounded CLI smoke; accepts a candidate or installed binary.

Usage: python3 test/parent_git_capability_smoke.py /absolute/ply /new/private/output
Creates a fresh physical /private/tmp workspace and retains command/effect evidence.
Setup writes finish before snapshots. Only format responses are modeled; all other
Git reads are allowlisted and delegated to pinned Homebrew Git. No provider starts.
"""
import hashlib
import json
import os
from pathlib import Path
import signal
import subprocess
import sys
import time


def main():
    binary = Path(sys.argv[1]).resolve(strict=True)
    output = Path(sys.argv[2]).absolute()
    output = output.resolve()
    assert output.is_relative_to(Path('/private/tmp')), 'use physical external scratch'
    output.mkdir(mode=0o700)  # Refuse an existing directory or an earlier attempt.
    workspace = output / 'workspace'
    repo, epic = workspace / 'ply/main', workspace / 'ply/epic'
    repo.mkdir(parents=True)
    home, temp = output / 'home', output / 'tmp'
    home.mkdir(); temp.mkdir()
    git = '/opt/homebrew/bin/git'
    env = {'PATH': '/opt/homebrew/bin:/usr/bin:/bin', 'HOME': str(home),
           'TMPDIR': str(temp), 'LC_ALL': 'C', 'LANG': 'C',
           'GIT_CONFIG_NOSYSTEM': '1', 'GIT_CONFIG_GLOBAL': '/dev/null',
           'GIT_TERMINAL_PROMPT': '0', 'GIT_OPTIONAL_LOCKS': '0'}
    deadline, sequence = time.monotonic() + 120, 0

    def command(argv, expected=0, cwd=workspace, child_env=None, required_error=None):
        nonlocal sequence
        sequence += 1
        remaining = deadline - time.monotonic()
        assert remaining > 0, '120-second smoke budget exhausted'
        argv = list(map(str, argv))
        started = time.monotonic()
        process = subprocess.Popen(argv, cwd=cwd, env=child_env or env,
                                   stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                                   text=True, start_new_session=True)
        try:
            stdout, stderr = process.communicate(timeout=min(30, remaining))
        except subprocess.TimeoutExpired:
            os.killpg(process.pid, signal.SIGKILL)
            process.communicate()
            raise
        result = subprocess.CompletedProcess(argv, process.returncode, stdout, stderr)
        (output / f'{sequence:03d}.json').write_text(json.dumps({
            'argv': argv, 'cwd': str(cwd), 'stdout': result.stdout,
            'stderr': result.stderr, 'exit': result.returncode,
            'elapsed_seconds': time.monotonic() - started}, indent=2) + '\n')
        assert (result.returncode == 0) == (expected == 0), (argv, result, output)
        if required_error:
            assert required_error in result.stderr, (argv, result, output)
        return result.stdout

    def ply(*args, expected=0, child_env=None, required_error=None):
        return command([binary, *args], expected, child_env=child_env, required_error=required_error)

    def js(*args):
        return json.loads(ply(*args, '--format', 'json'))

    def real_git(path, *args):
        return command([git, '-C', path, '-c', 'user.name=Ply fixture',
                        '-c', 'user.email=tests@example.invalid', *args]).strip()

    def write(name, value):
        path = workspace / name
        path.write_text(json.dumps(value))
        return path

    def draft(kind, key):
        return {'kind': kind, 'schema_version': 1, 'format': 'json', 'format_version': 1,
                'canonicalization': 'RFC8785', 'publication_key': key, 'task_id': 'task',
                'registry_upgrade': None, 'recorder': {
                    'actor_claim': 'synthetic fixture', 'control_surface': 'isolated CLI smoke',
                    'recorded_at_utc': '2026-10-24T00:00:00Z'}}

    def revision(value):
        return {k: value['outcome_ref'][k] for k in ('revision', 'manifest_sha256')}

    def decision(value):
        return {k: value['outcome_ref'][k] for k in ('id', 'manifest_sha256')}

    target = ['--project', 'ply', '--repo', 'ply', '--epic', 'epic']
    ply('workspace', 'init')
    real_git(repo, 'init', '-b', 'main')
    real_git(repo, 'commit', '--allow-empty', '-m', 'isolated base')
    oid = real_git(repo, 'rev-parse', 'HEAD')
    real_git(repo, 'worktree', 'add', '-b', 'epic', epic, oid)
    ply('workspace', 'project', 'add', 'ply', '--name', 'Fixture',
        '--wrapper', workspace / 'ply', '--repo', f'ply={repo}')
    ply('workspace', 'epic', 'adopt', 'epic', '--title', 'Fixture', '--project', 'ply',
        '--repo', 'ply', '--worktree', epic, '--ref', 'refs/heads/epic', '--expected-oid', oid)
    ply('workspace', 'task', 'create', 'task', '--title', 'Fixture', '--description',
        'Synthetic parent capability fixture', '--epic', 'epic', '--project', 'ply', '--repo', 'ply')
    show = js('workspace', 'task', 'show', 'task')
    base = js('workspace', 'epic', 'base', 'show', 'epic', '--repo', 'ply')['versions'][-1]
    spec = draft('WorkspaceTaskSpecDraft@1', 'fixture/spec')
    spec.update(spec_id='solution', expected_previous=None, problem=show['problem']['head'],
                title='Synthetic execution solution', parts={}, documents=[], supporting=[])
    refs = {}
    for name in ('abstract', 'functional', 'technical'):
        path = workspace / f'{name}.md'
        data = f'# Fixture {name}\n\nPreserve exact parent capability diagnostics.\n'.encode()
        path.write_bytes(data)
        spec['documents'].append({'id': name, 'source': {
            'kind': 'file', 'locator': str(path), 'sha256': 'sha256:' + hashlib.sha256(data).hexdigest(),
            'size_bytes': len(data), 'media_type': 'text/markdown', 'git_provenance': None}})
        refs[name] = [{'document_id': name, 'section': None}]
        spec['parts'][name] = {'state': 'present', 'reason': None, 'documents': refs[name]}
    spec.update(requirements=[{'id': 'f-01', 'functional_refs': refs['functional'],
                               'technical_refs': refs['technical'], 'verification_ids': ['tests'],
                               'acceptance': 'Preserve exact diagnostics without preview writes.'}],
                removed_requirement_ids=[], dependencies=[], change_reason='Initial fixture',
                phases=[{'id': 'verify', 'purpose': 'Verify fixture', 'requirement_ids': ['f-01'],
                         'entry_criteria': [], 'exit_criteria': ['Evidence retained'],
                         'verification_ids': ['tests']}], implementation_basis={
                    'project_id': 'ply', 'repo_id': 'ply', 'epic_id': 'epic',
                    'git_common_dir': str(repo / '.git'), 'parent_worktree_id': base['worktree_id'],
                    'parent_ref': base['ref'], 'parent_oid': base['oid'], 'parent_tree': base['tree'],
                    'start_oid': base['oid'], 'start_tree': base['tree']})
    recorded = js('workspace', 'task', 'spec', 'record', 'task', '--file', write('spec.json', spec))
    assessment = draft('WorkspaceTaskSpecAssessmentDraft@1', 'fixture/assessment')
    assessment.update(spec_id='solution', spec=revision(recorded), expected_previous_assessment=None,
                      outcome='ready', reason='Synthetic fixture only', open_questions=[], documents=[],
                      checks=[{'id': k, 'outcome': 'pass', 'reason': 'Synthetic fixture only',
                               'evidence_document_ids': []} for k in (
                          'acceptance_coverage', 'implementation_basis', 'problem_coverage',
                          'recovery', 'scope_and_phases', 'three_parts')])
    assessed = js('workspace', 'task', 'spec', 'assess', 'task', '--file', write('assessment.json', assessment))
    human = {'actor_claim': 'synthetic fixture, not human QA',
             'decided_at_utc': '2026-10-24T00:00:00Z', 'source': 'explicit_human_instruction',
             'statement': 'Synthetic local fixture selection only.'}
    selection = draft('WorkspaceTaskSolutionSelectionDraft@1', 'fixture/selection')
    selection.update(expected_previous_selection=None, action='select', reason='Synthetic fixture',
                     human_decision=human, solution={'spec_id': 'solution', 'spec': revision(recorded),
                     'problem': show['problem']['head'], 'assessment': decision(assessed)})
    selected = js('workspace', 'task', 'spec', 'select', 'task', '--file', write('selection.json', selection))
    queue = js('workspace', 'task', 'queue', 'list', *target)
    queued = {'kind': 'WorkspaceTaskQueueDraft@1', 'schema_version': 1, 'publication_key': 'fixture/queue',
              'project_id': 'ply', 'repo_id': 'ply', 'epic_id': 'epic',
              'expected_revision': queue['revision'], 'entries': [{'task_id': 'task', 'selection': decision(selected)}],
              'registry_upgrade': queue['registry']['required_upgrade'], 'human_decision': human}
    js('workspace', 'task', 'queue', 'set', '--file', write('queue.json', queued))
    goal_args = ['workspace', 'task', 'goal', 'register', 'goal-task', '--spec', 'fixture-goal',
                 '--document', workspace / 'technical.md', '--title', 'Fixture goal',
                 '--objective', 'Preserve parent diagnostics', '--criterion', 'Preview writes no state',
                 *target, '--delivery', 'local_epic_integration', '--acceptance', 'human']
    first = js(*goal_args)
    js(*goal_args, '--apply', '--confirmation', first['confirmation'])
    # Supported public preparation provides a genuine pre-change confirmation.
    fresh = js('workspace', 'task', 'prepare', 'task', *target)
    assert fresh['confirmation'] and fresh['plan']['observed_parent']['ref_format'] == 'files'

    fixture_bin = output / 'fixture-bin'
    fixture_bin.mkdir()
    mode, git_log = output / 'format.json', output / 'git-reads.jsonl'
    wrapper = fixture_bin / 'git'
    wrapper.write_text(f'#!{sys.executable} -S\n' + '''import json, os, subprocess, sys
from pathlib import Path
args = sys.argv[1:]
original = list(args)
while len(args) >= 2 and args[0] == '-c': args = args[2:]
if args and args[0] == 'check-ref-format': pass
else:
    assert len(args) >= 3 and args[0] == '-C'
    assert Path(args[1]).resolve().is_relative_to(Path(ROOT))
    args = args[2:]
assert args and (args[0] in ('rev-parse', 'status', 'show-ref', 'check-ref-format') or
    args[:3] == ['symbolic-ref', '-q', 'HEAD'] or args[:2] == ['worktree', 'list'] or
    args[:2] == ['merge-base', '--is-ancestor'] or
    args[:5] == ['diff-tree', '-r', '-z', '--no-commit-id', '--name-only'])
mode = json.loads(Path(MODE).read_text())
if args == ['rev-parse', mode['option']]:
    out, err, code = mode['output'], '', 0
else:
    result = subprocess.run([GIT, *original], capture_output=True, timeout=5)
    out, err, code = result.stdout.decode(), result.stderr.decode(), result.returncode
with open(LOG, 'a') as log:
    log.write(json.dumps({'argv': original, 'cwd': os.getcwd(), 'stdout': out, 'stderr': err, 'exit': code}) + '\\n')
sys.stdout.write(out); sys.stderr.write(err); sys.exit(code)
'''.replace('ROOT', repr(str(workspace / 'ply'))).replace('MODE', repr(str(mode)))
        .replace('GIT', repr(git)).replace('LOG', repr(str(git_log))))
    wrapper.chmod(0o700)
    fixture_env = dict(env, PATH=str(fixture_bin) + ':' + env['PATH'])

    def snapshot():
        files = {str(p.relative_to(workspace)): hashlib.sha256(p.read_bytes()).hexdigest()
                 for p in (workspace / '.ply').rglob('*') if p.is_file()}
        return {'registry_and_content': files,
                'directories': sorted(str(p.relative_to(workspace)) for p in workspace.rglob('*') if p.is_dir()),
                'refs': real_git(repo, 'show-ref'), 'inventory': real_git(repo, 'worktree', 'list', '--porcelain')}

    old = snapshot()
    (output / 'before.json').write_text(json.dumps(old, indent=2) + '\n')
    for option, response in (('--show-ref-format', '--show-ref-format\n'),
                             ('--show-ref-format', '\n'), ('--show-ref-format', 'reftable\n'),
                             ('--show-object-format=storage', '--show-object-format=storage\n')):
        mode.write_text(json.dumps({'option': option, 'output': response}))
        for selector in ('task', '--next'):
            # Exercise both renderers on echo; the other format cases retain
            # JSON coverage within the same 120-second smoke bound.
            for fmt in (('json', 'text') if response == '--show-ref-format\n' else ('json',)):
                out = ply('workspace', 'task', 'prepare', selector, *target, '--format', fmt, child_env=fixture_env)
                assert 'task_queue_git_capability_unsupported' in out and option in out
                if fmt == 'json':
                    value = json.loads(out)
                    assert value['state'] == 'blocked' and value['freshness'] == 'unknown'
                    assert value['confirmation'] is None and value['plan'] is None
                    assert json.dumps(response.removesuffix('\n')) in value['reasons'][0]['message']
                else:
                    assert '--apply' not in out and 'Freshness: unknown' in out
        # Only the invalid goal preview is invoked: it fails before runtime/provider observation.
        out = ply('workflow', 'execute', '--spec', 'fixture-goal', *target, '--check',
                  '--format', 'json', expected=1, child_env=fixture_env,
                  required_error='task_queue_git_capability_unsupported')
        assert json.loads(out).get('goal') is None
        # Apply an earlier supported digest after only the observed format changed.
        out = ply('workspace', 'task', 'prepare', 'task', *target, '--apply',
                  '--confirm', fresh['confirmation'], '--format', 'json', expected=1, child_env=fixture_env)
        value = json.loads(out)
        assert value['confirmation'] is None and value['preparation'] is None
        assert value['reasons'][0]['code'] == 'task_queue_git_capability_unsupported'
        assert snapshot() == old, 'preview/refused apply changed registry/content/refs/inventory'
    (output / 'after.json').write_text(json.dumps(snapshot(), indent=2) + '\n')
    reads = [json.loads(line) for line in git_log.read_text().splitlines()]
    assert sum('--show-ref-format' in row['argv'] for row in reads) >= 3
    assert sum('--show-object-format=storage' in row['argv'] for row in reads) >= 1
    print(f'T0523E9D1 CLI smoke PASS; {sequence} commands; evidence {output}')


if __name__ == '__main__':
    main()
