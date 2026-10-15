"""Candidate-bound, private/offline MR-05 regressions; not full product delivery."""
import hashlib
import json
import os
from pathlib import Path
import re
import signal
import stat
import subprocess
import sys
import tempfile
import time

# Source-derived expected ref supplied by root. Actual public native preview
# must match it before start; mismatch requires reviewed Goal revision. No override.
EXPECTED_REF = 'refs/heads/ply-task/mr-05-bound-runtime-control-replacement-ae7cec4dc362'
REPORT_KIND = 'MR05ReplacementFocusedAcceptance@1'
EVIDENCE_SCOPE = 'candidate-bound deterministic regressions; implementation stage only; AC01 through AC10 remain root-owned product obligations'
GO = Path('/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go')
GO_SHA = 'sha256:9da68c657a8344623d37fc9dc048d845011736409249bc924dd9af47a61594e6'
MODULE_CACHE = Path('/Users/perottochristensen/go/pkg/mod')
SCHEMA_PYTHON = Path('/Users/perottochristensen/.local/share/ply/gs01-20261009/schema-venv/bin/python')
SCHEMA_PYTHON_SHA = 'sha256:01328ad249b1c29f731960419d4145c427b4e3456378d0a53645c4586ae081f9'
TOTAL_SECONDS = 780
MODULE = 'github.com/devdimensionlab/plybuild'
REQUIRED = {'./internal/taskrun': ('TestMR05ControlDiscoveryAndOwnership', 'TestMR05ControlPreviewStatusReadOnlyStrictRequest', 'TestMR05ControlStaleBasisAndABA', 'TestMR05ControlIdempotentConcurrentApply', 'TestMR05ControlFaultsPreserveUnknown', 'TestMR05ControlAttemptBoundariesAndSameRunResume', 'TestMR05ControlLegacyAndNativeMutationRaces'), './internal/weekrun': ('TestMR05WeekControlScopeAndOptIn', 'TestMR05WeekPauseDispatchReservationRace', 'TestMR05WeekControlResumePreservesSelectionAndDeadline', 'TestMR05WeekControlDeadlineSupersedesPendingResume', 'TestMR05WeekControlChildAndLegacyAuthority'), './cmd': ('TestMR05ControlCandidateCLIAndSchema', 'TestMR05ControlCandidateCLIReadOnlyErrors')}


def require(ok, detail):
    if not ok:
        raise RuntimeError(detail)


def sha(path):
    h = hashlib.sha256()
    with Path(path).open('rb') as stream:
        for block in iter(lambda: stream.read(1024 * 1024), b''):
            h.update(block)
    return 'sha256:' + h.hexdigest()


def physical_env(name):
    value = os.environ.get(name, '')
    p = Path(value)
    require(value and p.is_absolute() and p.resolve(strict=True) == p,
            name + ' must be an existing physical absolute path')
    return p


def group_alive(pid):
    try:
        os.killpg(pid, 0)
        return True
    except ProcessLookupError:
        return False


def main():
    require(isinstance(EXPECTED_REF, str) and EXPECTED_REF.startswith('refs/heads/')
            and 'PLACEHOLDER' not in EXPECTED_REF, 'route is unbound; root must bind the fresh native Task ref')
    repo = Path.cwd().resolve(strict=True)
    require((repo / 'go.mod').is_file(), 'run from the bound candidate repository')
    home, tmp, candidate = (physical_env(name) for name in ('HOME', 'TMPDIR', 'PLY_CANDIDATE_BINARY'))
    require(home.is_dir() and tmp.is_dir() and home != tmp, 'distinct native private HOME/TMPDIR are required')
    require(all(stat.S_IMODE(p.stat().st_mode) & 0o077 == 0 for p in (home, tmp)), 'HOME/TMPDIR must be private')
    require(stat.S_ISREG(candidate.stat().st_mode) and os.access(candidate, os.X_OK), 'candidate must be a regular executable')
    oid, tree = (os.environ.get(key, '') for key in ('PLY_CANDIDATE_OID', 'PLY_CANDIDATE_TREE'))
    require(all(re.fullmatch(r'(?:[0-9a-f]{40}|[0-9a-f]{64})', value) for value in (oid, tree)), 'native candidate OID/tree are required')
    require(GO.is_file() and os.access(GO, os.X_OK) and MODULE_CACHE.is_dir() and sha(GO) == GO_SHA,
            'pinned existing Go/toolchain cache is unavailable; downloading is forbidden')
    require(SCHEMA_PYTHON.is_file() and os.access(SCHEMA_PYTHON, os.X_OK) and sha(SCHEMA_PYTHON) == SCHEMA_PYTHON_SHA,
            'pinned existing schema interpreter is unavailable; installation is forbidden')
    env = {'HOME': str(home), 'TMPDIR': str(tmp), 'PATH': os.environ.get('PATH', '/usr/bin:/bin'),
           'LC_ALL': 'C', 'GIT_OPTIONAL_LOCKS': '0', 'GIT_CONFIG_NOSYSTEM': '1'}

    def git(*args):
        p = subprocess.run(['/usr/bin/git', '-c', 'core.fsmonitor=false', '-c', 'core.hooksPath=/dev/null', *args],
                           cwd=repo, env=env, capture_output=True, text=True, timeout=10)
        require(p.returncode == 0, 'read-only Git observation failed: ' + p.stderr.strip())
        return p.stdout.strip()

    def binding():
        require(git('symbolic-ref', '--quiet', 'HEAD') == EXPECTED_REF, 'wrong Task ref')
        require(git('rev-parse', 'HEAD') == oid and git('rev-parse', 'HEAD^{tree}') == tree, 'native candidate OID/tree mismatch')
        require(git('status', '--porcelain=v1', '--untracked-files=all') == '', 'candidate source must remain clean')

    binding()
    candidate_sha = sha(candidate)
    scratch = Path(tempfile.mkdtemp(prefix='mr05-replacement-acceptance-', dir=tmp))
    os.chmod(scratch, 0o700)
    for name in ('tmp', 'go-cache', 'go-tmp'):
        (scratch / name).mkdir(mode=0o700)
    env.update({'TMPDIR': str(scratch / 'tmp'), 'GOCACHE': str(scratch / 'go-cache'),
                'GOTMPDIR': str(scratch / 'go-tmp'), 'GOMODCACHE': str(MODULE_CACHE),
                'GOROOT': str(GO.parent.parent), 'GOTOOLCHAIN': 'local', 'GOENV': 'off',
                'GOWORK': 'off', 'GOPROXY': 'off', 'GOSUMDB': 'off',
                'PLY_CANDIDATE_BINARY': str(candidate), 'PLY_DELIVERY_TEST_BINARY': str(candidate),
                'PLY_CANDIDATE_OID': oid, 'PLY_CANDIDATE_TREE': tree,
                'PLY_MR05_SCHEMA_PYTHON': str(SCHEMA_PYTHON)})
    identity = {'source_ref': EXPECTED_REF, 'candidate_oid': oid, 'candidate_tree': tree,
                'candidate_binary': {'path': str(candidate), 'sha256': candidate_sha},
                'go': {'path': str(GO), 'sha256': GO_SHA},
                'schema_python': {'path': str(SCHEMA_PYTHON), 'sha256': SCHEMA_PYTHON_SHA},
                'driver_sha256': sha(Path(__file__).resolve(strict=True)),
                'required_tests': REQUIRED, 'artifact_directory': str(scratch), 'scope': EVIDENCE_SCOPE}
    (scratch / 'identity.json').write_text(json.dumps(identity, indent=2) + '\n')
    start = time.monotonic()
    interrupted = {'signal': None}
    commands, errors, coverage = [], [], {}

    def signal_request(number, frame):
        if interrupted['signal'] is None:
            interrupted['signal'] = number

    handlers = {n: signal.signal(n, signal_request) for n in (signal.SIGTERM, signal.SIGHUP, signal.SIGINT)}

    def run(label, argv):
        child, code, timed_out, problem = None, None, False, None
        cleanup = []
        stdout_path, stderr_path = scratch / (label + '.stdout'), scratch / (label + '.stderr')
        with stdout_path.open('wb') as stdout, stderr_path.open('wb') as stderr:
            try:
                if interrupted['signal'] is None:
                    if TOTAL_SECONDS - (time.monotonic() - start) <= 0:
                        timed_out = True
                    else:
                        child = subprocess.Popen(argv, cwd=repo, env=env, stdout=stdout, stderr=stderr, start_new_session=True)
                while child is not None and interrupted['signal'] is None:
                    remaining = TOTAL_SECONDS - (time.monotonic() - start)
                    if remaining <= 0:
                        timed_out = True
                        break
                    try:
                        code = child.wait(timeout=min(0.25, remaining))
                        break
                    except subprocess.TimeoutExpired:
                        pass
            except BaseException as error:
                problem = type(error).__name__ + ': ' + str(error)
            finally:
                if child is not None:
                    alive = group_alive(child.pid)
                    if code == 0 and alive:
                        cleanup.append('a descendant outlived the test runner')
                    if alive and child.returncode is None:
                        # This sole waiter has not reaped the child: its PID cannot
                        # be reused. Never signal a group after leader reaping.
                        try:
                            require(os.getpgid(child.pid) == child.pid, 'owned child process group changed')
                            os.killpg(child.pid, signal.SIGKILL)
                        except ProcessLookupError:
                            pass
                        except BaseException as error:
                            cleanup.append('cannot positively stop own unreaped test process group: ' + str(error))
                    elif alive:
                        cleanup.append('group remains after leader reaping; no late signal is authorized')
                    try:
                        code = child.wait(timeout=10)
                    except BaseException as error:
                        cleanup.append('test child exit unresolved: ' + str(error))
                    until = time.monotonic() + 3
                    while group_alive(child.pid) and time.monotonic() < until:
                        time.sleep(0.05)
                    if group_alive(child.pid):
                        cleanup.append('test process group still observed after cleanup')
        result = {'label': label, 'argv': argv, 'exit_code': code, 'timed_out': timed_out,
                  'interrupted_signal': interrupted['signal'], 'runner_error': problem, 'cleanup_errors': cleanup,
                  'stdout': {'path': str(stdout_path), 'sha256': sha(stdout_path)},
                  'stderr': {'path': str(stderr_path), 'sha256': sha(stderr_path)}}
        commands.append(result)
        require(code == 0 and not timed_out and interrupted['signal'] is None and problem is None and not cleanup,
                'command did not complete cleanly: ' + label)
        return stdout_path

    try:
        schema_out = run('schema-prerequisite', [str(SCHEMA_PYTHON), '-I', '-B', '-c',
            'import importlib.metadata; assert importlib.metadata.version("jsonschema") == "4.26.0"; print("jsonschema 4.26.0")'])
        require(schema_out.read_text().strip() == 'jsonschema 4.26.0', 'schema prerequisite was not positively observed')
        for index, (package, names) in enumerate(REQUIRED.items(), 1):
            output = run('go-' + str(index), [str(GO), 'test', '-json', '-mod=readonly', '-count=1', '-p=1',
                '-parallel=1', '-timeout=240s', '-run', '^(' + '|'.join(names) + ')$', package])
            observed = {name: {'run': 0, 'pass': 0} for name in names}
            package_pass, bad = 0, []
            expected_package = MODULE + '/' + package.removeprefix('./')
            for line in output.read_text().splitlines():
                event = json.loads(line)
                require(event.get('Package') == expected_package, 'unexpected Go package output')
                action, name = event.get('Action'), event.get('Test')
                if action in ('fail', 'skip'):
                    bad.append({'test': name, 'action': action})
                if name in observed and action in ('run', 'pass'):
                    observed[name][action] += 1
                if name is None and action == 'pass':
                    package_pass += 1
            coverage[package] = {'required': observed, 'package_pass': package_pass, 'bad_events': bad}
            require(not bad and package_pass == 1 and all(value == {'run': 1, 'pass': 1} for value in observed.values()),
                    'mandatory MR-05 regression missing, skipped, failed or repeated: ' + package)
    except Exception as error:
        errors.append(str(error))
    finally:
        for n, handler in handlers.items():
            signal.signal(n, handler)
    try:
        binding()
        require(sha(candidate) == candidate_sha, 'candidate CLI bytes changed during acceptance')
    except Exception as error:
        errors.append(str(error))
    if interrupted['signal'] is not None:
        errors.append('interrupted by signal ' + str(interrupted['signal']))
    if set(coverage) != set(REQUIRED):
        errors.append('all mandatory packages were not observed')
    result = {'kind': REPORT_KIND, 'outcome': 'pass' if not errors else 'fail', 'identity': identity,
              'coverage': coverage, 'commands': commands, 'errors': errors,
              'interrupted_signal': interrupted['signal'], 'elapsed_seconds': time.monotonic() - start,
              'scope': EVIDENCE_SCOPE}
    (scratch / 'summary.json').write_text(json.dumps(result, indent=2) + '\n')
    print(json.dumps(result, sort_keys=True))
    return 0 if not errors else 1


if __name__ == '__main__':
    try:
        raise SystemExit(main())
    except Exception as error:
        print(json.dumps({'kind': REPORT_KIND, 'outcome': 'fail', 'error': str(error), 'scope': EVIDENCE_SCOPE}))
        raise SystemExit(1)
