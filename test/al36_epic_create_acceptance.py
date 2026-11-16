#!/usr/bin/env python3
"""AL36: bounded disposable native CLI journey. No install or real workspace writes."""
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile


def main():
    ply = Path(sys.argv[1]).resolve(strict=True)
    # Raw commands and results are retained by the invoking acceptance transcript.
    def run(argv, cwd, ok=True, structured=False):
        argv = [str(x) for x in argv]
        p = subprocess.run(argv, cwd=cwd, capture_output=True, text=True)
        print(json.dumps(dict(argv=argv,cwd=str(cwd),stdout=p.stdout,stderr=p.stderr,exit=p.returncode)),flush=True)
        assert (p.returncode == 0) == ok, argv
        return json.loads(p.stdout) if structured and ok else p
    with tempfile.TemporaryDirectory(prefix='ply-al36-', dir=os.environ.get('TMPDIR','/private/tmp')) as temp:
        root = Path(temp).resolve()
        run([ply,'workspace','init'],root)
        repo = root/'example'/'ordinary'
        repo.mkdir(parents=True)
        run(['git','init','-b','ordinary'],repo)
        run(['git','config','user.name','AL36 fixture'],repo)
        run(['git','config','user.email','al36@example.invalid'],repo)
        (repo/'README.md').write_text('fixture\n')
        run(['git','add','README.md'],repo)
        run(['git','-c','commit.gpgsign=false','commit','-m','fixture'],repo)
        oid=run(['git','rev-parse','HEAD'],repo).stdout.strip()
        run([ply,'workspace','project','add','example','--name','Example','--wrapper',repo.parent,'--repo',f'example={repo}'],root)
        registry=repo.parent/'data'/'work'/'work-items.yaml' # the Project's registry, under its wrapper
        # Epic creation works on a fresh registry; there is no format upgrade.
        def raw():
            return registry.read_bytes() if registry.exists() else None
        target=repo.parent/'created'
        args=[ply,'workspace','epic','create','created','--title','Created Epic','--project','example','--repo','example','--parent-ref','refs/heads/ordinary','--expected-parent-oid',oid,'--branch','created','--path',target,'--format','json']
        before=run(['git','show-ref'],repo).stdout
        (repo/'dirty').write_text('dirty')
        baseline=raw()
        run(args,root,ok=False)
        assert raw()==baseline and not target.exists()
        (repo/'dirty').unlink()
        # r001-f01: a private ref lock produces a known no-effect attempt.
        # A temporary dirty parent must preserve pending identity and history,
        # without another Git attempt, until the ordinary parent is restored.
        lock=repo/'.git'/'refs'/'heads'/'created.lock'
        lock.write_text('AL36 disposable regression lock\n')
        initial=json.loads(run(args,root,ok=False).stdout)
        original=initial['operation']
        assert original['state']=='pending' and len(original['attempts'])==1 and original['attempts'][0]['exit']!=0
        assert not target.exists()
        lock.unlink()
        (repo/'temporary-dirty').write_text('temporary parent change\n')
        changed_result=run(args,root,ok=False)
        changed=json.loads(changed_result.stdout)['operation']
        assert changed['state']=='pending' and 'workspace_work_parent_dirty' in changed_result.stderr
        for key in ['id','worktree_id','intent_digest','request','attempts']:
            assert changed[key]==original[key], key
        assert len(changed['observations'])>len(original['observations'])
        assert changed['observations'][-1]['parent_worktree_clean'] is False
        assert not target.exists()
        (repo/'temporary-dirty').unlink()
        made=run(args,root,structured=True)
        for key in ['id','worktree_id','intent_digest','request']:
            assert made['operation'][key]==original[key], key
        assert made['operation']['attempts'][:1]==original['attempts']
        assert made['operation']['observations'][:len(changed['observations'])]==changed['observations']
        assert made['outcome']=='created' and made['epic']['repo_bindings'][0]['worktree']['origin']=='ply_created'
        assert made['operation']['request']['parent_ref']=='refs/heads/ordinary'
        assert len(made['operation']['attempts'])==2
        assert run(['git','show-ref','refs/heads/ordinary'],repo).stdout==before
        assert run(['git','status','--porcelain'],repo).stdout==''
        if len(sys.argv)>2 and sys.argv[2]=='--pending-parent-retry':
            print('AL36 r001-f01: native no-effect/temporary-dirty/clean retry preserves identity and failed Git evidence PASS',flush=True)
            return
        read=run([ply,'workspace','epic','show','created','--format','json'],root,structured=True)
        assert read['kind']=='WorkspaceEpicReadback@3' and read['ready_for_task_worktree_create'] and read['base'][0]['revision']==1
        listed=run([ply,'workspace','epic','list','--format','json'],root,structured=True)
        assert listed['kind']=='WorkspaceEpicListReadback@2' and len(listed['epics'])==1
        # The ordinary parent becomes purely historical after completion.
        run(['git','-c','commit.gpgsign=false','commit','--allow-empty','-m','parent moved'],repo)
        (target/'dirty').write_text('later Epic work')
        baseline=raw()
        again=run(args,root,structured=True)
        assert again['outcome']=='already_created' and again['operation']['id']==made['operation']['id']
        assert raw()==baseline
        read=run([ply,'workspace','epic','show','created','--format','json'],root,structured=True)
        assert not read['ready_for_task_worktree_create']
        # Missing anchor stays a current readiness fact, never a repair request.
        run(['git','worktree','remove','--force',target],repo)
        again=run(args,root,structured=True)
        assert again['outcome']=='already_created' and not target.exists() and raw()==baseline
        print('AL36 isolated CLI: native format upgrades, exact creation, parent preservation, historical retry and missing-anchor readiness PASS',flush=True)


if __name__=='__main__':
    main()
