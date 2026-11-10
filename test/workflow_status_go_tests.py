#!/usr/bin/env python3
"""Run the workflow-status regression groups and preserve the complete test inventory."""
import json
from pathlib import Path
import subprocess
import sys
import time

repo = Path(__file__).resolve().parents[1]
evidence = Path(sys.argv[1]).resolve()
go = sys.argv[2] if len(sys.argv) > 2 else 'go'
inventory = subprocess.run([go, 'test', '-json', '-list', '^Test', './...'], cwd=repo, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
(evidence / 'go-inventory.jsonl').write_text(inventory.stdout)
(evidence / 'go-inventory.stderr').write_text(inventory.stderr)
if inventory.returncode:
    sys.stderr.write(inventory.stderr)
    raise SystemExit(inventory.returncode)
registered = set()
for line in inventory.stdout.splitlines():
    event = json.loads(line)
    name = event.get('Output', '').strip()
    if name.startswith('Test') and ' ' not in name:
        registered.add((event['Package'], name))

groups = [
    ('cli-and-projection', ['./cmd', './internal/workspaceview'], '.'),
    ('registered-queue-and-content', ['./internal/workspace'], '^Test(RegisteredQueue|TaskContent|TaskGoal|TaskQueue)'),
    ('preserved-run-and-report', ['./internal/taskrun'], '^Test(ReadInventory|InventoryDelivery|InventoryCaptured|DeliveryOwnerActualAcceptance)'),
]
observed = {}
started_tests = set()
results = []
failed = False
for name, packages, pattern in groups:
    argv = [go, 'test', '-json', '-count=1', '-timeout=5m', *packages, '-run', pattern]
    print(f'Running {name}: {" ".join(argv)}', flush=True)
    started = time.monotonic()
    with (evidence / f'{name}.jsonl').open('w') as output:
        result = subprocess.run(argv, cwd=repo, stdout=output, stderr=subprocess.STDOUT)
    for line in (evidence / f'{name}.jsonl').read_text().splitlines():
        try:
            event = json.loads(line)
        except ValueError:
            continue
        test = event.get('Test')
        if test and '/' not in test:
            key = (event['Package'], test)
            if event.get('Action') == 'run':
                started_tests.add(key)
            elif event.get('Action') in ('pass', 'fail', 'skip'):
                observed[key] = event['Action']
    results.append(dict(id=name, argv=argv, cwd=str(repo), exit=result.returncode, elapsed_seconds=time.monotonic()-started))
    print(f'{name}: exit {result.returncode}, {results[-1]["elapsed_seconds"]:.1f}s; log {evidence / (name + ".jsonl")}', flush=True)
    failed |= result.returncode != 0

report = dict(kind='WorkflowStatusGoTestCoverage@1', groups=results, tests=[dict(package=p, test=t, outcome=observed.get((p,t), 'incomplete' if (p,t) in started_tests else 'not_run')) for p,t in sorted(registered)])
(evidence / 'go-test-coverage.json').write_text(json.dumps(report, indent=2)+'\n')
counts = {state: sum(test['outcome']==state for test in report['tests']) for state in ('pass','fail','skip','incomplete','not_run')}
print(f'Test inventory: {counts}. Tests outside the affected groups were not run; see go-test-coverage.json.', flush=True)
raise SystemExit(1 if failed or counts['fail'] or counts['incomplete'] else 0)
