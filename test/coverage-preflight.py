import argparse
from datetime import datetime, timezone
import hashlib
import json
import os
from pathlib import Path
import shlex
import subprocess
import sys


REQUIRED_ENV = {'GOENV': 'off', 'GOWORK': 'off', 'GOTOOLCHAIN': 'local',
                'GOPROXY': 'off', 'GOSUMDB': 'off'}
CONTEXT_ENV = (*REQUIRED_ENV, 'GOFLAGS', 'GOOS', 'GOARCH', 'CGO_ENABLED', 'PATH',
               'HOME', 'TMPDIR', 'GOTMPDIR', 'GOCACHE', 'GOMODCACHE', 'CC', 'CXX',
               'CGO_CFLAGS', 'CGO_CPPFLAGS', 'CGO_CXXFLAGS', 'CGO_LDFLAGS',
               'PLY_CANDIDATE_BINARY', 'PLY_CANDIDATE_OID', 'PLY_CANDIDATE_TREE',
               'PLY_REGRESSION_GO', 'PLY_MR05_SCHEMA_PYTHON', 'PLY_SL01_SCHEMA_PYTHON',
               'PLY_DELIVERY_TEST_BINARY', 'PLY_TASK_RUN_TEST_BINARY')


def run(argv):
    result = subprocess.run(argv, text=True, capture_output=True, timeout=30)
    if result.returncode:
        raise ValueError(f'{shlex.join(argv)} failed: {result.stderr.strip()}')
    return result.stdout


def digest(path):
    with Path(path).open('rb') as stream:
        return hashlib.file_digest(stream, 'sha256').hexdigest()


def identity(path):
    path = Path(path)
    if not path.is_absolute() or not path.is_file() or not os.access(path, os.X_OK):
        raise ValueError(f'absolute executable required: {path}')
    return {'path': str(path), 'resolved': str(path.resolve()), 'sha256': digest(path)}


def required(name):
    value = os.environ.get(name)
    if not value:
        raise ValueError(f'{name} is required; no build, install or fallback')
    return value


def source():
    root = run(['git', 'rev-parse', '--show-toplevel']).strip()
    if Path(root).resolve() != Path.cwd().resolve() or not Path('go.mod').is_file():
        raise ValueError('run preflight from the repository/module root')
    files = run(['git', 'ls-files', '-z']).split('\0')
    return {'root': root, 'oid': run(['git', 'rev-parse', 'HEAD']).strip(),
            'tree': run(['git', 'rev-parse', 'HEAD^{tree}']).strip(),
            'status': run(['git', 'status', '--porcelain', '--untracked-files=all']),
            'files': {name: digest(name) if Path(name).is_file() else None for name in files if name}}


def build_inputs(go):
    raw = run([go, 'list', '-deps', '-test', '-json', './...'])
    decoder, files = json.JSONDecoder(), set()
    while raw.strip():
        package, end = decoder.raw_decode(raw.lstrip())
        raw = raw.lstrip()[end:]
        directory = Path(package.get('Dir', '.'))
        for key in ('GoFiles', 'CgoFiles', 'CFiles', 'CXXFiles', 'MFiles', 'HFiles', 'FFiles',
                    'SFiles', 'SwigFiles', 'SwigCXXFiles', 'SysoFiles', 'TestGoFiles', 'XTestGoFiles',
                    'EmbedFiles', 'TestEmbedFiles', 'XTestEmbedFiles'):
            for name in package.get(key, []):
                path = directory / name
                if path.is_file():
                    files.add(str(path.resolve()))
    return {path: digest(path) for path in sorted(files)}


def snapshot():
    state = source()
    candidate = identity(required('PLY_CANDIDATE_BINARY'))
    go = identity(required('PLY_REGRESSION_GO'))
    python = identity(required('PLY_MR05_SCHEMA_PYTHON'))
    for alias, expected in [('PLY_DELIVERY_TEST_BINARY', candidate['path']),
                            ('PLY_TASK_RUN_TEST_BINARY', candidate['path']),
                            ('PLY_SL01_SCHEMA_PYTHON', python['path'])]:
        if os.environ.get(alias) not in (None, '', expected):
            raise ValueError(f'conflicting input alias {alias}; expected {expected}')
    for key, value in REQUIRED_ENV.items():
        if os.environ.get(key) != value:
            raise ValueError(f'{key} must explicitly be {value}')
    flags = shlex.split(required('GOFLAGS'))
    if '-mod=readonly' not in flags:
        raise ValueError('GOFLAGS must preserve -mod=readonly')
    for flag in flags:
        if flag != '-mod=readonly' and not flag.startswith('-tags='):
            raise ValueError(f'unsupported GOFLAGS for this coverage setup: {flag}')
    tags = [flag.split('=', 1)[1] for flag in flags if flag.startswith('-tags=')]
    if len(tags) > 1:
        raise ValueError('GOFLAGS must have at most one -tags= value')
    if state['oid'] != required('PLY_CANDIDATE_OID') or state['tree'] != required('PLY_CANDIDATE_TREE'):
        raise ValueError('source OID/tree do not match intended candidate')
    if state['status']:
        raise ValueError('source must be clean before coverage')
    version = json.loads(run([candidate['path'], 'version', '--format', 'json']))
    if version.get('source_commit') != state['oid'] or version.get('dirty') is not False:
        raise ValueError('packaged candidate source_commit/dirty does not match clean source')
    candidate['version'] = version
    candidate['build_metadata'] = run([go['path'], 'version', '-m', candidate['path']])
    settings = {}
    for line in candidate['build_metadata'].splitlines():
        if line.startswith('\tbuild\t') and '=' in line:
            key, value = line.removeprefix('\tbuild\t').split('=', 1)
            settings[key] = value
    for key in ('GOOS', 'GOARCH', 'CGO_ENABLED'):
        if settings.get(key) != required(key):
            raise ValueError(f'candidate build setting {key} does not match explicit environment')
    for key in ('CGO_CFLAGS', 'CGO_CPPFLAGS', 'CGO_CXXFLAGS', 'CGO_LDFLAGS'):
        if os.environ.get(key) and settings.get(key) != os.environ[key]:
            raise ValueError(f'candidate build setting {key} does not match environment')
    if settings.get('-tags', '') != (tags[0] if tags else ''):
        raise ValueError('candidate build tags do not match GOFLAGS tags')
    go['version'] = run([go['path'], 'version']).strip()
    go['environment'] = json.loads(run([go['path'], 'env', '-json']))
    if go['version'].split()[2] not in candidate['build_metadata'].splitlines()[0].split():
        raise ValueError('candidate Go toolchain version does not match test toolchain')
    for key in ('GOEXPERIMENT', 'GOARM64', 'GOAMD64', 'GOARM', 'GO386', 'GOMIPS', 'GOMIPS64', 'GOPPC64', 'GORISCV64', 'GOWASM'):
        if key in settings and settings[key] != go['environment'].get(key, ''):
            raise ValueError(f'candidate build setting {key} does not match effective Go context')
    try:
        schema = json.loads(run([python['path'], '-I', '-B', '-c', '''import hashlib, importlib.metadata, json, jsonschema, pathlib, sys
version = importlib.metadata.version('jsonschema')
assert version == '4.26.0', 'expected jsonschema 4.26.0, got ' + version
jsonschema.Draft202012Validator.check_schema({'type': 'object'})
packages = {}
for name in ('jsonschema', 'jsonschema-specifications', 'referencing', 'rpds-py', 'attrs'):
    dist = importlib.metadata.distribution(name)
    packages[name] = {'version': dist.version, 'files': {str(f): hashlib.sha256(pathlib.Path(dist.locate_file(f)).read_bytes()).hexdigest() for f in dist.files or [] if pathlib.Path(dist.locate_file(f)).is_file() and '__pycache__' not in str(f)}}
print(json.dumps({'python_version': sys.version, 'jsonschema_version': version, 'module': jsonschema.__file__, 'packages': packages}))''']))
    except (ValueError, json.JSONDecodeError) as error:
        raise ValueError(f'schema Python must support pinned jsonschema under -I -B: {error}') from error
    schema['executable'] = python
    return {'source': state, 'candidate': candidate, 'go': go, 'schema': schema,
            'build_inputs': build_inputs(go['path']),
            'environment': {key: os.environ.get(key) for key in CONTEXT_ENV}}


def write(path, value):
    with Path(path).open('x') as output:
        json.dump(value, output, indent=2, sort_keys=True)
        output.write('\n')


def finish(args):
    before = json.loads(Path(args.after).read_text())
    gaps = []
    current_source = source()
    current_inputs = {}
    current_go_environment = None
    current_build_inputs = None
    if current_source != before['source']:
        gaps.append('source endpoint changed')
    for name, old in [('candidate', before['candidate']), ('go', before['go']),
                      ('schema Python', before['schema']['executable'])]:
        current_inputs[name] = identity(old['path'])
        if any(current_inputs[name][key] != old[key] for key in ('path', 'resolved', 'sha256')):
            gaps.append(f'{name} endpoint changed')
    for key, value in before['environment'].items():
        if os.environ.get(key) != value:
            gaps.append(f'environment endpoint changed: {key}')
    try:
        current = snapshot()
        current_go_environment = current['go']['environment']
        current_build_inputs = current['build_inputs']
        if current_build_inputs != before['build_inputs']:
            gaps.append('selected Go source/test/embed build inputs endpoint changed')
        if {key: value for key, value in current['go']['environment'].items() if key != 'GOGCCFLAGS'} != {key: value for key, value in before['go']['environment'].items() if key != 'GOGCCFLAGS'}:
            gaps.append('effective Go build context endpoint changed')
        if current['schema'] != before['schema']:
            gaps.append('schema packages endpoint changed')
    except ValueError as error:
        gaps.append(str(error))
    log = Path(args.after).with_name('test.jsonl')
    packages, tests = [], []
    if log.is_file():
        for line in log.read_text().splitlines():
            try:
                event = json.loads(line)
            except json.JSONDecodeError:
                gaps.append('invalid test.jsonl event')
                continue
            if event.get('Action') in ('pass', 'fail', 'skip'):
                (tests if event.get('Test') else packages).append(event)
        if not packages or any(e['Action'] == 'fail' for e in packages + tests):
            gaps.append('missing package completion or failed test event')
    else:
        gaps.append('missing test.jsonl')
    profile_args = [arg.split('=', 1)[1] for arg in before['argv'] if arg.startswith('-coverprofile=')]
    profile = Path(profile_args[0]) if len(profile_args) == 1 else None
    if profile is None or not profile.is_file() or not profile.stat().st_size:
        gaps.append('missing completed coverage profile')
    if args.test_exit != 0:
        gaps.append(f'go test exited {args.test_exit}')
    result = {'before': {'path': str(Path(args.after).resolve()), 'sha256': digest(args.after)},
              'finished_at': datetime.now(timezone.utc).isoformat(), 'source': current_source,
              'inputs': current_inputs, 'effective_go_environment': current_go_environment,
              'build_inputs': current_build_inputs, 'test_exit': args.test_exit, 'packages': packages, 'tests': tests,
              'profile': {'path': str(profile), 'sha256': digest(profile)} if profile and profile.is_file() else None,
              'test_log_sha256': digest(log) if log.is_file() else None,
              'gaps': gaps, 'endpoint_checks_passed': not gaps, 'baseline_qualified': False,
              'limitation': 'Hash equality is endpoint evidence, not continuous monitoring or proof of executed argv. Consumer must confirm full-suite exit and immutable build context before analysis.'}
    write(args.output, result)
    if gaps:
        raise ValueError('; '.join(gaps))


def main():
    parser = argparse.ArgumentParser(description='Validate explicit coverage inputs and record endpoint evidence; never builds, installs or runs tests.')
    parser.add_argument('--output', required=True)
    parser.add_argument('--after')
    parser.add_argument('--test-exit', type=int)
    parser.add_argument('argv', nargs=argparse.REMAINDER)
    args = parser.parse_args()
    if Path(args.output).exists():
        raise ValueError(f'output already exists: {args.output}')
    if args.after:
        if args.test_exit is None:
            raise ValueError('--after requires actual --test-exit')
        finish(args)
    else:
        result = snapshot()
        result['recorded_at'] = datetime.now(timezone.utc).isoformat()
        result['argv'] = args.argv[1:] if args.argv[:1] == ['--'] else args.argv
        result['argv_status'] = 'declared, not executed by preflight'
        write(args.output, result)
    print(f'coverage preflight evidence: {args.output}')


if __name__ == '__main__':
    try:
        main()
    except (ValueError, OSError, subprocess.SubprocessError) as error:
        print(f'coverage preflight: {error}', file=sys.stderr)
        sys.exit(1)
