#!/usr/bin/env python3
"""Exercise exact PM01 candidate in an isolated workspace; no provider or delivery."""
import hashlib
import json
import os
from pathlib import Path
import subprocess
import tempfile

ply = Path(os.environ['PLY_CANDIDATE_BINARY']).resolve(strict=True)
# Native verification's scratch is removed; retain raw command evidence separately.
evidence = Path(tempfile.mkdtemp(prefix='ply-pm01-evidence-', dir='/private/tmp')).resolve()
print(f'PM01 raw evidence: {evidence}', flush=True)
sequence = 0


def run(argv, cwd, code=0):
    global sequence
    sequence += 1
    argv = [str(a) for a in argv]
    result = subprocess.run(argv, cwd=cwd, capture_output=True, text=True)
    record = dict(argv=argv, cwd=str(cwd), stdout=result.stdout,
                  stderr=result.stderr, exit=result.returncode)
    (evidence / f'{sequence:03}.json').write_text(json.dumps(record, indent=2))
    print(json.dumps(record), flush=True)
    assert result.returncode == code, record
    return result.stdout


def cli(root, *args, code=0):
    return run([ply, 'workspace', *args], root, code)


def git(repo, *args):
    return run(['git', '-C', repo, *args], repo)


def repo_state(repo):
    return (git(repo, 'show-ref'), git(repo, 'status', '--porcelain=v1'),
            git(repo, 'worktree', 'list', '--porcelain'))


with tempfile.TemporaryDirectory(prefix='ply-pm01-fixture-') as temporary:
    base = Path(temporary).resolve()
    root = base / 'workspace'
    wrapper = root / 'shared'
    first = wrapper / 'main'
    extra = base / 'external/main'
    for repo in (first, extra):
        repo.mkdir(parents=True)
        git(repo, 'init', '-q', '-b', 'main')
        (repo / 'source').write_text('baseline\n')
        git(repo, 'add', 'source')
        git(repo, '-c', 'user.name=PM01 fixture', '-c', 'user.email=fixture@example.invalid',
            '-c', 'commit.gpgsign=false', 'commit', '-qm', 'Fixture baseline')
    cli(root, 'init')
    cli(root, 'project', 'add', 'shared', '--name', 'Shared', '--wrapper', wrapper,
        '--repo', f'first={first}')
    # Preserve actual native Epic and Task history, without starting a Run.
    cli(root, 'epic', 'adopt', 'history', '--title', 'Existing history', '--project', 'shared',
        '--repo', 'first', '--worktree', first, '--ref', 'refs/heads/main',
        '--expected-oid', git(first, 'rev-parse', 'HEAD').strip())
    cli(root, 'task', 'create', 'history-task', '--title', 'Existing task', '--description',
        'Preserved binding', '--project', 'shared', '--repo', 'first', '--epic', 'history', '--format', 'json')
    linked = base / 'linked'
    git(extra, 'worktree', 'add', '-q', '-b', 'linked', linked)
    alias = base / 'alias'
    alias.symlink_to(extra, target_is_directory=True)
    for repo in (first, extra):
        (repo / 'source').write_text('dirty tracked source\n')
        (repo / 'untracked').write_text('dirty untracked source\n')
    states = {repo: repo_state(repo) for repo in (first, extra)}
    protected = {p.relative_to(root): p.read_bytes() for p in (root / '.ply').rglob('*')
                 if p.is_file() and p.name not in ('projects.yaml', 'projects.lock')}
    before = json.loads(cli(root, 'project', 'show', 'shared', '--format', 'json'))
    out = cli(root, 'project', 'add-member', 'shared', '--repo', f'extra={extra}')
    assert 'Added repository member extra' in out
    registry = root / '.ply/projects.yaml'
    registered = registry.read_bytes()
    for path in (extra, alias, linked):
        assert 'Already registered repository member extra' in cli(
            root, 'project', 'add-member', 'shared', '--repo', f'extra={path}')
        assert registry.read_bytes() == registered
    after = json.loads(cli(root, 'project', 'show', 'shared', '--format', 'json'))
    assert after['project'] == dict(before['project'], repo_count=2)
    assert after['repositories'] == [dict(repo_id='extra', locator=str(extra), git_common_dir=str(extra / '.git'), wrapper=None),
                                     before['repositories'][0]]
    listed = json.loads(cli(root, 'project', 'list', '--format', 'json'))
    assert len(listed['projects']) == 1 and listed['projects'][0]['repo_count'] == 2
    for project, member, error in (
            ('shared', f'alias={linked}', 'workspace_project_conflict'),
            ('shared', f'extra={first}', 'workspace_project_conflict'),
            ('shared', f'missing={base / "missing"}', 'workspace_project_path_error'),
            ('absent', f'extra={extra}', 'workspace_project_not_found'),
            ('shared', f'notgit={base}', 'workspace_project_repo_invalid')):
        sequence_before = sequence
        cli(root, 'project', 'add-member', project, '--repo', member, code=1)
        record = json.loads((evidence / f'{sequence_before + 1:03}.json').read_text())
        assert error in record['stderr']
        assert registry.read_bytes() == registered
    # Known ID can be shared by another Project; no locator rebinding.
    other = root / 'other'
    other.mkdir()
    cli(root, 'project', 'add', 'other', '--name', 'Other', '--wrapper', other, '--repo', f'first={first}')
    cli(root, 'project', 'add-member', 'other', '--repo', f'extra={linked}')
    reused = json.loads(cli(root, 'project', 'show', 'other', '--format', 'json'))
    assert reused['repositories'] == after['repositories']
    assert {p.relative_to(root): p.read_bytes() for p in (root / '.ply').rglob('*')
            if p.is_file() and p.name not in ('projects.yaml', 'projects.lock')} == protected
    for repo in (first, extra):
        assert repo_state(repo) == states[repo]
        assert (repo / 'source').read_text() == 'dirty tracked source\n'
        assert (repo / 'untracked').read_text() == 'dirty untracked source\n'
    task = json.loads(cli(root, 'task', 'show', 'history-task', '--format', 'json'))
    assert task
    cli(root, 'project', 'add-member', '--help')
    print('PM01 PASS: add/readback/retry/reuse/aliases/conflicts/missing/dirty/history preservation', flush=True)
