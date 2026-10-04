#!/usr/bin/env python3
"""Supplement TJ-01 acceptance against an isolated, native fixture manifest.

Uses the actual CLI. Only the supplied fixture is perturbed, with original bytes
restored in finally blocks. No runner, provider, or installed binary is started.
"""
import argparse
import copy
import hashlib
import json
from pathlib import Path
import subprocess
import unittest


class NativeBoundaries(unittest.TestCase):
    serial = 0

    def call(self, *args, ok=True):
        NativeBoundaries.serial += 1
        command = [CONFIG.ply, 'workspace', 'task', 'journal', *args]
        cp = subprocess.run(command, cwd=ROOT, capture_output=True, timeout=30)
        label = CONFIG.output / f'{NativeBoundaries.serial:03d}'
        label.with_suffix('.stdout').write_bytes(cp.stdout)
        label.with_suffix('.stderr').write_bytes(cp.stderr)
        label.with_suffix('.json').write_text(json.dumps(
            dict(command=command, cwd=str(ROOT), exit=cp.returncode), indent=2))
        self.assertEqual(cp.returncode == 0, ok, cp.stderr.decode())
        return cp

    def show(self):
        return json.loads(self.call('show', TASK, '--format', 'json').stdout)

    def assert_independent(self, snap):
        self.assertEqual(snap['coverage']['state'], 'partial')
        self.assertTrue(any(e['type'] == 'task_state' for e in snap['events']))
        self.assertTrue(any(e['data'].get('axis') == 'technical'
                            for e in snap['events']))

    def test_dirty_and_missing_worktree(self):
        worktree = Path(self.show()['task']['worktree']['locator'])
        dirty = worktree / 'boundary-dirty.txt'
        self.assertFalse(dirty.exists())
        dirty.write_text('Synthetic dirty candidate\n')
        try:
            snap = self.show()
            self.assertIsNone(snap['current']['candidate'])
            self.assertTrue(all(a['value'] == 'unknown'
                                for a in snap['current']['axes'].values()))
        finally:
            dirty.unlink()
        saved = worktree.with_name(worktree.name + '.boundary-saved')
        self.assertFalse(saved.exists())
        worktree.rename(saved)
        try:
            snap = self.show()
            self.assertIsNone(snap['current']['candidate'])
            self.assert_independent(snap)
        finally:
            saved.rename(worktree)

    def test_missing_content_document(self):
        src = next(s for s in self.show()['sources']
                   if s['kind'] == 'task_content_document')
        path = Path(src['locator'])
        saved = path.with_name(path.name + '.boundary-saved')
        self.assertFalse(saved.exists())
        path.rename(saved)
        try:
            snap = self.show()
            self.assert_independent(snap)
            self.assertTrue(any(s['locator'] == str(path) and s['status'] == 'missing'
                                for s in snap['sources']))
        finally:
            saved.rename(path)

    def test_native_sequence_clock_conflict(self):
        path = sorted((RUN / 'events').glob('*.json'))[-1]
        before = path.read_bytes()
        modified = json.loads(before)
        modified['recorded_at_utc'] = '2000-01-01T00:00:00Z'
        try:
            path.write_text(json.dumps(modified, sort_keys=True, separators=(',', ':')))
            snap = self.show()
            self.assertTrue(snap['coverage']['time_conflict_event_ids'])
            self.assertTrue(any(r['code'] == 'time_conflict'
                                for r in snap['coverage']['reasons']))
            self.assert_independent(snap)
        finally:
            path.write_bytes(before)

    def test_corrupt_optional_run_cache(self):
        path = RUN / 'result.json'
        before = path.read_bytes()
        try:
            path.write_text('{"invalid": true}')
            snap = self.show()
            self.assert_independent(snap)
            self.assertTrue(any(e['run_id'] == FIXTURE['runs'][0]['run_id']
                                and e['type'] == 'run_reserved'
                                for e in snap['events']))
        finally:
            path.write_bytes(before)

    def test_wrong_native_binding_and_append(self):
        path = RUN / 'request.json'
        before = path.read_bytes()
        modified = json.loads(before)
        modified['preparation_sha256'] = 'sha256:' + '0' * 64
        try:
            path.write_text(json.dumps(modified, sort_keys=True, separators=(',', ':')))
            snap = self.show()
            self.assertTrue(any(r['code'] == 'run_unbound'
                                for r in snap['coverage']['reasons']))
            self.assertFalse(any(e['run_id'] == FIXTURE['runs'][0]['run_id']
                                 and e['type'].startswith('run_')
                                 for e in snap['events']))
            self.assert_independent(snap)
        finally:
            path.write_bytes(before)
        source = CONFIG.output / 'evidence.txt'
        source.write_text('Synthetic binding evidence\n')
        data = dict(kind='ply.workspace.task-journal-event-input', schema_version=1,
                    publication_key='boundary-invalid-run', task_id=TASK,
                    actor=dict(id='observer', role='observer', session_id=None),
                    activity_id=None, run_binding=copy.deepcopy(FIXTURE['runs'][0]),
                    occurred_at=None,
                    time_basis=dict(kind='unknown', clock=None, precision='unknown',
                                    uncertainty=None),
                    sources=[dict(locator=str(source),
                                  sha256=hashlib.sha256(source.read_bytes()).hexdigest())],
                    type='note', title='Invalid binding', detail='', step=None,
                    outcome=None, candidate=None, relations=[], waiting=None)
        journal = ROOT / '.ply/task-process'
        self.assertFalse(journal.exists())
        for field in ('run_id', 'request_sha256', 'preparation_id', 'preparation_sha256'):
            bad = copy.deepcopy(data)
            bad['run_binding'][field] = 'sha256:' + '0' * 64 if field.endswith('sha256') else 'foreign'
            input_path = CONFIG.output / ('invalid-' + field + '.json')
            input_path.write_text(json.dumps(bad))
            self.call('record', TASK, '--file', str(input_path), '--format', 'json', ok=False)
            self.assertFalse(journal.exists())

    def test_argument_failures_without_write(self):
        self.call('show', 'nonexistent-task', ok=False)
        self.call('show', TASK, 'extra', ok=False)
        self.call('show', TASK, '--view', 'invalid', ok=False)
        self.call('record', TASK, '--format', 'json', ok=False)
        self.assertFalse((ROOT / '.ply/task-process').exists())


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--ply', required=True)
    parser.add_argument('--fixture', type=Path, required=True)
    parser.add_argument('--output', type=Path, required=True)
    CONFIG = parser.parse_args()
    CONFIG.output.mkdir(mode=0o700, parents=True, exist_ok=False)
    FIXTURE = json.loads(CONFIG.fixture.read_text())
    ROOT = Path(FIXTURE['workspace'])
    TASK = FIXTURE['task_id']
    RUN = ROOT / '.ply/task-runs/v1/runs' / FIXTURE['runs'][0]['run_id']
    unittest.main(argv=['task-journal-native-boundaries'], verbosity=2)
