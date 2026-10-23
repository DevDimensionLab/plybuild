import json
import os
from pathlib import Path
import subprocess
import shlex
import sys
import tempfile
import unittest


SCRIPT = Path(__file__).with_name('coverage-preflight.py')


class CoveragePreflightTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name) / 'source'
        self.root.mkdir()
        self.out = Path(self.tmp.name) / 'before.json'
        self.git('init', '-q')
        (self.root / 'go.mod').write_text('module fixture\n')
        self.git('add', '.')
        self.git('-c', 'user.name=Fixture', '-c', 'user.email=fixture@example.invalid', 'commit', '-qm', 'fixture')
        self.oid = self.git('rev-parse', 'HEAD').strip()
        self.tree = self.git('rev-parse', 'HEAD^{tree}').strip()
        self.candidate = Path(self.tmp.name) / 'ply'
        self.candidate.write_text('#!/bin/sh\nprintf \'%s\\n\' \'%s\'\n' % ('%s', json.dumps({'source_commit': self.oid, 'dirty': False})))
        self.candidate.chmod(0o755)
        self.go = Path(self.tmp.name) / 'go'
        self.go.write_text('''#!/bin/sh
if [ "$1" = env ]; then printf '{"GOOS":"darwin","GOARCH":"arm64","CGO_ENABLED":"1","GOEXPERIMENT":"%s"}\\n' "$GOEXPERIMENT";
elif [ "$1" = list ]; then printf '{"Dir":"%s","GoFiles":["go.mod"]}\\n' "$PWD";
elif [ "$2" = -m ]; then printf 'fixture: go1.26.7\\n\\tbuild\\tGOOS=darwin\\n\\tbuild\\tGOARCH=arm64\\n\\tbuild\\tCGO_ENABLED=1\\n';
else echo 'go version go1.26.7 darwin/arm64'; fi
''')
        self.go.chmod(0o755)
        self.env = dict(os.environ, PLY_CANDIDATE_BINARY=str(self.candidate), PLY_CANDIDATE_OID=self.oid,
                        PLY_CANDIDATE_TREE=self.tree, PLY_REGRESSION_GO=str(self.go),
                        PLY_MR05_SCHEMA_PYTHON=sys.executable, GOOS='darwin', GOARCH='arm64', CGO_ENABLED='1',
                        GOENV='off', GOWORK='off', GOFLAGS='-mod=readonly', GOPROXY='off', GOSUMDB='off', GOTOOLCHAIN='local')
        for name in ('PLY_DELIVERY_TEST_BINARY', 'PLY_TASK_RUN_TEST_BINARY', 'PLY_SL01_SCHEMA_PYTHON'):
            self.env.pop(name, None)

    def git(self, *args):
        return subprocess.check_output(['git', '-C', str(self.root), *args], text=True)

    def invoke(self, *args, success=True):
        result = subprocess.run([sys.executable, '-I', '-B', str(SCRIPT), *args], env=self.env,
                                cwd=self.root, text=True, capture_output=True)
        self.assertEqual(result.returncode == 0, success, result.stdout + result.stderr)
        return result

    def begin(self):
        return self.invoke('--output', str(self.out), '--', str(self.go), 'test', '-json', '-count=1', '-timeout=30m',
                           '-coverprofile=' + str(Path(self.tmp.name) / 'coverage.out'), './...')

    def test_missing_relative_and_mismatched_inputs_fail_before_output(self):
        cases = [('PLY_CANDIDATE_BINARY', '', 'PLY_CANDIDATE_BINARY'),
                 ('PLY_CANDIDATE_BINARY', 'ply', 'absolute executable'),
                 ('PLY_CANDIDATE_OID', '0' * 40, 'OID/tree'),
                 ('PLY_MR05_SCHEMA_PYTHON', '/missing/python', 'executable'),
                 ('PLY_TASK_RUN_TEST_BINARY', '/other/ply', 'alias'),
                 ('GOENV', '', 'GOENV'), ('GOARCH', 'amd64', 'build setting'),
                 ('GOFLAGS', '-mod=readonly -tags=custom', 'tags')]
        for key, value, message in cases:
            with self.subTest(key=key, value=value):
                old = self.env.get(key)
                self.env[key] = value
                result = self.invoke('--output', str(self.out), success=False)
                self.assertIn(message, result.stderr)
                self.assertFalse(self.out.exists())
                if old is None:
                    del self.env[key]
                else:
                    self.env[key] = old

    def test_isolated_schema_import_rejects_pythonpath_only_dependency(self):
        python = Path(self.tmp.name) / 'python'
        python.write_text('#!/bin/sh\nexec ' + shlex.quote(sys.executable) + ' -S "$@"\n')
        python.chmod(0o755)
        self.env['PLY_MR05_SCHEMA_PYTHON'] = str(python)
        self.env['PYTHONPATH'] = str(Path(sys.executable).parent)
        self.assertIn('schema Python', self.invoke('--output', str(self.out), success=False).stderr)

    def test_success_records_context_hashes_and_refuses_overwrite(self):
        self.begin()
        before = json.loads(self.out.read_text())
        self.assertEqual(before['source']['oid'], self.oid)
        self.assertIn('go.mod', before['source']['files'])
        self.assertIn(str((self.root / 'go.mod').resolve()), before['build_inputs'])
        self.assertEqual(before['schema']['jsonschema_version'], '4.26.0')
        self.assertEqual(before['argv'][-1], './...')
        self.assertEqual(before['environment']['GOPROXY'], 'off')
        self.assertIn('sha256', before['candidate'])
        self.assertIn('already exists', self.invoke('--output', str(self.out), success=False).stderr)

    def finish(self, test_exit=0, success=True):
        log = Path(self.tmp.name) / 'test.jsonl'
        log.write_text(json.dumps({'Action': 'pass', 'Package': 'fixture', 'Test': 'TestOne', 'Elapsed': 2.5}) + '\n' +
                       json.dumps({'Action': 'pass', 'Package': 'fixture', 'Elapsed': 3.0}) + '\n')
        profile = Path(self.tmp.name) / 'coverage.out'
        profile.write_text('mode: set\nfixture/a.go:1.1,2.1 1 1\n')
        return self.invoke('--after', str(self.out), '--output', str(Path(self.tmp.name) / 'after.json'),
                           '--test-exit', str(test_exit), success=success)

    def test_finish_records_profile_and_package_times_without_qualifying_failure(self):
        self.begin()
        self.finish(1, success=False)
        after = json.loads((Path(self.tmp.name) / 'after.json').read_text())
        self.assertFalse(after['endpoint_checks_passed'])
        self.assertEqual(after['test_exit'], 1)
        self.assertEqual(after['packages'][0]['Elapsed'], 3.0)
        self.assertEqual(after['tests'][0]['Elapsed'], 2.5)
        self.assertIn('sha256', after['profile'])

    def test_finish_rejects_endpoint_source_and_binary_drift(self):
        for target in ('source', 'candidate'):
            with self.subTest(target=target):
                self.begin()
                path = self.root / 'go.mod' if target == 'source' else self.candidate
                original = path.read_bytes()
                path.write_bytes(original + b'\n')
                self.finish(success=False)
                after = json.loads((Path(self.tmp.name) / 'after.json').read_text())
                self.assertFalse(after['endpoint_checks_passed'])
                self.assertIn(target, ' '.join(after['gaps']))
                path.write_bytes(original)
                self.out.unlink()
                (Path(self.tmp.name) / 'after.json').unlink()

    def test_subdirectory_invocation_rejects_before_partial_source_hash_or_inputs(self):
        self.root = self.root / 'subdirectory'
        self.root.mkdir()
        self.env['PLY_CANDIDATE_BINARY'] = ''
        result = self.invoke('--output', str(self.out), success=False)
        self.assertIn('repository/module root', result.stderr)
        self.assertFalse(self.out.exists())

    def test_finish_rejects_changed_effective_go_build_context(self):
        self.env['GOEXPERIMENT'] = ''
        self.begin()
        self.env['GOEXPERIMENT'] = 'newinliner'
        self.finish(success=False)
        after = json.loads((Path(self.tmp.name) / 'after.json').read_text())
        self.assertFalse(after['endpoint_checks_passed'])
        self.assertIn('effective Go build context', ' '.join(after['gaps']))

    def test_finish_success_is_only_endpoint_and_exit_evidence(self):
        self.begin()
        self.finish()
        after = json.loads((Path(self.tmp.name) / 'after.json').read_text())
        self.assertTrue(after['endpoint_checks_passed'])
        self.assertIn('endpoint', after['limitation'])
        self.assertFalse(after['baseline_qualified'])


if __name__ == '__main__':
    unittest.main()
