#!/usr/bin/env python3
"""Scope regression for the installed-view acceptance harness itself."""
import os
from pathlib import Path
import subprocess
import tempfile
from types import SimpleNamespace
import unittest
from unittest.mock import patch

from workspace_views_acceptance import Journey


class AcceptanceScope(unittest.TestCase):
    def test_inherited_git_repository_overrides_cannot_redirect_fixture_writes(self):
        with tempfile.TemporaryDirectory(prefix="ply-acceptance-scope-") as temporary:
            root = Path(temporary).resolve()
            canary = root / "unrelated-canary"
            canary.mkdir()
            clean = {key: value for key, value in os.environ.items() if not key.startswith("GIT_")}
            clean.update(GIT_CONFIG_NOSYSTEM="1", GIT_CONFIG_GLOBAL="/dev/null")

            def git(*args):
                command = ["git", "-c", "user.name=Fixture", "-c", "user.email=fixture@invalid",
                           "-c", "commit.gpgsign=false", "-c", "core.hooksPath=/dev/null", *args]
                return subprocess.check_output(command, cwd=canary, env=clean, text=True).strip()

            git("init", "--initial-branch=main")
            git("commit", "--allow-empty", "-m", "Preserve canary")
            original = git("rev-parse", "HEAD")
            overrides = dict(GIT_DIR=str(canary / ".git"), GIT_WORK_TREE=str(canary),
                             GIT_INDEX_FILE=str(canary / "unexpected-index"),
                             GIT_OBJECT_DIRECTORY=str(canary / ".git" / "objects"))
            with patch.dict(os.environ, overrides):
                journey = Journey(SimpleNamespace(ply=Path("/usr/bin/true"), baseline=None, root=root / "journey"))
            main = journey.workspace / "main"
            main.mkdir()
            journey.git(main, "init", "--initial-branch=main")
            journey.git(main, "commit", "--allow-empty", "-m", "Fixture only")
            self.assertEqual(git("rev-parse", "HEAD"), original, "fixture commit escaped through caller Git environment")
            self.assertTrue((main / ".git").is_dir())
            self.assertFalse((canary / "unexpected-index").exists())


if __name__ == "__main__":
    unittest.main()
