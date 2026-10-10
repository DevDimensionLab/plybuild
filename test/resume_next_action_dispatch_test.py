"""Offline routing/binding tests; these are not product acceptance receipts.

Git observations are simulated. No checkout/ref is created or changed, no Go
test process is started, and the ordinary fullsuite is never run here.
"""

import contextlib
import hashlib
import io
import json
import os
from pathlib import Path
import shlex
import subprocess
import sys
import tempfile
import unittest
from unittest import mock


REPO = Path(__file__).resolve().parent.parent
REF = "refs/heads/ply-task/aw-01-current-resume-guidance-038968865f8f"
DRIVER = REPO / "test/resume-next-action-acceptance.sh"
DRIVER_SHA = "a5b3bb5e6f6d1062851f4c926fb5939ab096f89a315506ef06a4f5482321feb4"
DEFAULT_SHA = "61abee5a1bab34699d07fcc4883c248cbc019597bc4540561917638efe6d68ac"


class DispatchTest(unittest.TestCase):
    def setUp(self):
        self.scratch = tempfile.TemporaryDirectory(prefix="aw01-dispatch-", dir=os.environ["TMPDIR"])
        self.addCleanup(self.scratch.cleanup)
        self.root = Path(self.scratch.name).resolve()

    def script(self, path, body):
        path.write_text(body)
        path.chmod(0o700)

    def test_only_exact_task_ref_dispatches_before_default_setup(self):
        repo = self.root / "candidate source"
        test = repo / "test"
        test.mkdir(parents=True)
        dispatcher = test / "delivery_acceptance.sh"
        dispatcher.write_bytes((REPO / "test/delivery_acceptance.sh").read_bytes())
        self.script(test / "resume-next-action-acceptance.sh", "#!/bin/sh\npwd -P\nexit 0\n")
        bindir = self.root / "bin"
        bindir.mkdir()
        self.script(bindir / "git", "#!" + sys.executable + "\n" +
                    "import os,sys\n"
                    "expected = [\"-c\", \"core.fsmonitor=false\", \"-c\", \"core.hooksPath=/dev/null\", \"-C\", os.environ[\"ROUTE_ROOT\"], \"symbolic-ref\", \"--quiet\", \"HEAD\"]\n"
                    "assert sys.argv[1:] == expected, sys.argv\n"
                    "assert os.environ[\"GIT_OPTIONAL_LOCKS\"] == \"0\"\n"
                    "value = os.environ[\"ROUTE_REF\"]\n"
                    "print(value)\n"
                    "sys.exit(0 if value else 1)\n")
        self.script(bindir / "mktemp", "#!/bin/sh\nprintf 'ordinary-default\\n'\nexit 73\n")
        for ref in (REF, REF + "-nearby", "refs/heads/ply-task/another-task", "refs/heads/ply_absence_20261010", "refs/heads/main", ""):
            with self.subTest(ref=ref or "detached HEAD"):
                env = dict(os.environ, PATH=str(bindir) + os.pathsep + os.environ["PATH"], ROUTE_REF=ref, ROUTE_ROOT=str(repo))
                result = subprocess.run(["/bin/sh", str(dispatcher)], cwd=self.root, env=env, capture_output=True, text=True, timeout=10)
                self.assertEqual(result.returncode, 0 if ref == REF else 73, result.stderr)
                if ref == REF:
                    self.assertEqual(result.stdout.strip(), str(repo))

    def test_reviewed_driver_and_ordinary_default_remain_byte_identical(self):
        self.assertEqual(hashlib.sha256(DRIVER.read_bytes()).hexdigest(), DRIVER_SHA)
        source = (REPO / "test/delivery_acceptance.sh").read_bytes()
        default = source[source.index(b"# Native automatic verification"):]
        self.assertEqual(hashlib.sha256(default).hexdigest(), DEFAULT_SHA)

    def test_driver_requires_matching_native_candidate_and_clean_source(self):
        # Load the exact embedded Python, replacing only its external Git/test
        # process adapters. Deliberately never manufacture a Go pass stream.
        source = DRIVER.read_text()
        python = shlex.split(source[source.index("exec python3"):])[5]
        python = python[:python.rindex("\ntry:\n    raise SystemExit(main())")]
        namespace = {}
        exec(compile(python, str(DRIVER), "exec"), namespace)
        repo = self.root / "repo"
        (repo / "internal/taskrun").mkdir(parents=True)
        (repo / "go.mod").write_text("module fixture\n")
        for name in ("home", "tmp"):
            (self.root / name).mkdir()
        candidate = self.root / "candidate"
        self.script(candidate, "#!/bin/sh\nexit 99\n")
        oid, tree = "a" * 40, "b" * 40
        base = dict(os.environ, HOME=str(self.root / "home"), TMPDIR=str(self.root / "tmp"), PLY_CANDIDATE_BINARY=str(candidate), PLY_CANDIDATE_OID=oid, PLY_CANDIDATE_TREE=tree)
        previous = Path.cwd()
        self.addCleanup(os.chdir, previous)
        os.chdir(repo)
        cases = (
            ("matching", {}, REF, "", None),
            ("missing oid", {"PLY_CANDIDATE_OID": ""}, REF, "", "OID/tree are required"),
            ("wrong oid", {"PLY_CANDIDATE_OID": "c" * 40}, REF, "", "do not match"),
            ("missing tree", {"PLY_CANDIDATE_TREE": ""}, REF, "", "OID/tree are required"),
            ("wrong tree", {"PLY_CANDIDATE_TREE": "c" * 40}, REF, "", "do not match"),
            ("dirty", {}, REF, " M source.go", "must remain clean"),
            ("untracked", {}, REF, "?? unrelated", "must remain clean"),
            ("other ref", {}, REF + "-nearby", "", "exact selected Task ref"),
            ("detached", {}, "", "", "Git observation failed"),
            ("missing binary", {"PLY_CANDIDATE_BINARY": ""}, REF, "", "PLY_CANDIDATE_BINARY"),
        )
        for name, changes, ref, status, error in cases:
            with self.subTest(case=name):
                def git(argv, **kwargs):
                    self.assertEqual(argv[:5], ["/usr/bin/git", "-c", "core.fsmonitor=false", "-c", "core.hooksPath=/dev/null"])
                    self.assertEqual(kwargs["cwd"], repo)
                    self.assertEqual(kwargs["env"]["GIT_OPTIONAL_LOCKS"], "0")
                    command = tuple(argv[5:])
                    outputs = {("symbolic-ref", "--quiet", "HEAD"): ref, ("rev-parse", "HEAD"): oid, ("rev-parse", "HEAD^{tree}"): tree, ("status", "--porcelain=v1", "--untracked-files=all"): status}
                    return subprocess.CompletedProcess(argv, 1 if command[0] == "symbolic-ref" and not ref else 0, outputs[command], "")

                output = io.StringIO()
                with mock.patch.dict(os.environ, dict(base, **changes), clear=True), mock.patch.object(sys, "argv", ["driver", str(DRIVER)]), mock.patch("subprocess.run", side_effect=git), mock.patch("subprocess.Popen", side_effect=RuntimeError("routing test: Go deliberately not started")) as start, contextlib.redirect_stdout(output):
                    if error:
                        with self.assertRaisesRegex(RuntimeError, error):
                            namespace["main"]()
                        start.assert_not_called()
                    else:
                        self.assertEqual(namespace["main"](), 1)
                        start.assert_called_once()
                        result = json.loads(output.getvalue())
                        self.assertEqual(result["outcome"], "fail")
                        self.assertIn("Go deliberately not started", result["runner_error"])
                        self.assertTrue(all(row == {"run": 0, "pass": 0} for row in result["required_tests"].values()))


if __name__ == "__main__":
    unittest.main(verbosity=2)
