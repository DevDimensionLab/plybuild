#!/usr/bin/env python3
"""AL38 offline fixture; run with a make-built or retained installed CLI path."""
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile

binary = str(Path(sys.argv[1] if len(sys.argv) > 1 else os.environ['PLY_CANDIDATE_BINARY']).resolve())
with tempfile.TemporaryDirectory(prefix='ply-external-issues-') as temporary:
    root = Path(temporary).resolve()
    workspace = root / 'workspace'
    main = workspace / 'fixture' / 'main'
    epic = workspace / 'fixture' / 'epic'
    task = workspace / 'fixture' / 'task'
    home = root / 'home'
    home.mkdir()
    main.mkdir(parents=True)
    env = dict(os.environ, HOME=str(home), GIT_CONFIG_NOSYSTEM='1', GIT_CONFIG_GLOBAL=os.devnull)

    def call(argv, cwd=workspace, ok=True):
        p = subprocess.run(argv, cwd=cwd, env=env, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, timeout=30)
        if (p.returncode == 0) != ok:
            raise AssertionError(f'{argv}: exit {p.returncode}\n{p.stdout}\n{p.stderr}')
        return p.stdout if ok else p.stderr

    def git(*args, cwd=main):
        return call(['git', '-c', 'user.name=Ply fixture', '-c', 'user.email=fixture@example.invalid', *args], cwd).strip()

    def ply(*args, ok=True):
        return call([binary, 'workspace', *args], ok=ok)

    def links(kind, op, *args, ok=True):
        value = ply(kind, 'external-issue', op, kind, *args, '--format', 'json', ok=ok)
        return json.loads(value) if ok else value

    ply('init')
    git('init', '-b', 'main')
    (main / 'README.md').write_text('fixture\n')
    git('add', 'README.md')
    git('commit', '-m', 'fixture')
    oid = git('rev-parse', 'HEAD')
    git('worktree', 'add', '-b', 'epic', str(epic), oid)
    ply('project', 'add', 'fixture', '--name', 'Fixture', '--wrapper', str(main.parent), '--repo', f'fixture={main}')
    ply('epic', 'adopt', 'epic', '--title', 'Epic', '--project', 'fixture', '--repo', 'fixture', '--worktree', str(epic), '--ref', 'refs/heads/epic', '--expected-oid', oid)
    ply('task', 'create', 'task', '--title', 'Task', '--description', 'Offline fixture', '--project', 'fixture', '--repo', 'fixture', '--epic', 'epic')
    ply('task', 'worktree', 'create', 'task', '--branch', 'task', '--path', str(task), '--expected-parent-oid', oid)
    (task / 'dirty.txt').write_text('ongoing task work\n')
    refs_before = git('show-ref')
    dirty_before = git('status', '--porcelain', cwd=task)
    marker = workspace / '.ply'

    def metadata():
        return {str(p.relative_to(marker)): hashlib.sha256(p.read_bytes()).hexdigest()
                for p in marker.rglob('*') if p.is_file() and p.name != 'external-issues.json'}

    before = metadata()
    empty = links('task', 'show')
    assert empty['refs'] == [] and empty['primary'] is None and empty['effective'] is None
    assert 'external_issues' not in json.loads(ply('task', 'show', 'task', '--format', 'json'))
    jira = ['--provider', 'jira', '--key', 'PLY-123', '--url', 'http://JIRA.example.invalid:80/jira/browse/PLY-123', '--actor', 'fixture', '--primary', '--expected-revision', '0']
    added = links('epic', 'add', *jira)
    assert added['changed'] and added['primary']['commit_key'] == 'PLY-123'
    assert added['refs'][0]['jira_context_path'] == '/jira'
    ref = added['refs'][0]['id']
    bytes_before = (marker / 'external-issues.json').read_bytes()
    retry = links('epic', 'add', *jira)
    assert not retry['changed'] and retry['revision'] == 1
    assert (marker / 'external-issues.json').read_bytes() == bytes_before
    assert links('task', 'show')['effective']['subject_kind'] == 'epic'
    failure = links('task', 'add', *jira, ok=False)
    assert 'already linked to epic/epic' in failure
    assert (marker / 'external-issues.json').read_bytes() == bytes_before
    gh = links('task', 'add', '--provider', 'github', '--host', 'git.example.invalid', '--owner', 'OWNER', '--repository', 'REPO', '--number', '7', '--url', 'https://GIT.example.invalid:443/owner/Repo/issues/007/', '--actor', 'fixture', '--primary', '--expected-revision', '1')
    assert gh['effective']['commit_key'] == 'owner/repo#7' and gh['effective']['subject_kind'] == 'task'
    task_ref = gh['refs'][0]['id']
    assert json.loads(ply('task', 'show', 'task', '--format', 'json'))['external_issues']['effective'] == gh['effective']
    assert json.loads(ply('epic', 'show', 'epic', '--format', 'json'))['external_issues']['primary']['ref_id'] == ref
    cleared = links('task', 'primary', '--clear', '--expected-revision', '2', '--actor', 'fixture')
    assert cleared['primary'] is None and cleared['effective']['subject_kind'] == 'epic'
    # The old primary-add request must not reverse the deliberate clear.
    links('task', 'add', '--provider', 'github', '--host', 'git.example.invalid', '--owner', 'owner', '--repository', 'repo', '--number', '7', '--url', 'https://git.example.invalid/owner/repo/issues/7', '--actor', 'fixture', '--primary', '--expected-revision', '1', ok=False)
    selected = links('task', 'primary', task_ref, '--expected-revision', '3', '--actor', 'fixture')
    assert selected['effective']['subject_kind'] == 'task'
    unlinked = links('task', 'unlink', task_ref, '--expected-revision', '4', '--actor', 'fixture')
    assert unlinked['refs'] == [] and unlinked['primary'] is None and unlinked['effective']['subject_kind'] == 'epic'
    assert len(unlinked['events']) == 4
    links('epic', 'unlink', ref, '--expected-revision', '5', '--actor', 'fixture')
    assert links('task', 'list')['effective'] is None
    human = ply('task', 'external-issue', 'show', 'task')
    assert 'offline_unverified' in human and 'Effective commit reference: none' in human
    assert metadata() == before, 'unrelated workspace metadata changed'
    assert git('show-ref') == refs_before, 'Git refs changed'
    assert git('status', '--porcelain', cwd=task) == dirty_before, 'dirty source changed'
    assert (task / 'dirty.txt').read_text() == 'ongoing task work\n'
    print('AL38 make/installed CLI fixture passed: typed add/retry/conflict/primary/fallback/unlink; native metadata and Git unchanged')
