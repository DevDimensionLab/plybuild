#!/usr/bin/env python3
"""SM01 hermetic public CLI conformance and bounded source timing."""
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import time

binary = str(Path(sys.argv[1]).resolve(strict=True))

def row(pane="w1:p1", status="idle"):
    return dict(agent="codex", agent_status=status, pane_id=pane,
                tab_id="w1:t1", workspace_id="w1", agent_session=False,
                terminal_id={"private": "never-export"}, title="never-export")

def envelope(rows, **extra):
    return dict(result=dict(type="agent_list", agents=rows, **extra))

with tempfile.TemporaryDirectory(prefix="sm01-cli-") as tmp:
    root = Path(tmp)
    fixture = root / "response.json"
    calls = root / "calls"
    shim = root / "herdr"
    env = dict(os.environ, PATH=str(root) + os.pathsep + os.environ["PATH"])
    modes = [
        ("sessionless", envelope([row(status="future")]), "available", "complete", 1, None),
        ("empty", envelope([]), "empty", "complete", 0, None),
        ("partial", envelope([row(), 42]), "partial", "partial", 1, None),
        ("duplicate", envelope([row(), row()]), "partial", "partial", 0, None),
        ("unknown-type", dict(result=dict(type="future", agents=[])), "unavailable", "unknown", 0, None),
        ("truncation-signal", envelope([], truncated=True), "partial", "partial", 0, None),
        ("utf16-256", envelope([row("😀" * 128)]), "available", "complete", 1, None),
        ("utf16-257", envelope([row("😀" * 128 + "a")]), "partial", "partial", 0, None),
        ("control", envelope([row("bad\ntext")]), "partial", "partial", 0, None),
        ("4096", envelope([row(f"p{n:04d}") for n in range(4096)]), "available", "complete", 4096, None),
        ("4097", envelope([row(f"p{n:04d}") for n in range(4097)]), "partial", "partial", 4096, None),
        ("malformed", "{truncated", "unavailable", "unknown", 0, None),
        ("oversize", " " * (8 << 20) + json.dumps(envelope([])), "unavailable", "unknown", 0, None),
        ("unavailable", None, "unavailable", "unknown", 0, "exit 7"),
        ("timeout", None, "unavailable", "unknown", 0, "sleep 20"),
        ("inherited-pipes", None, "unavailable", "unknown", 0, "sleep 20 & wait"),
    ]
    for name, payload, status, coverage, count, command in modes:
        fixture.write_text(payload if isinstance(payload, str) else json.dumps(payload), encoding="utf-8")
        # No provider execution; this executable accepts exactly the public List operation.
        shim.write_text("#!/bin/sh\nprintf '%s\\n' \"$*\" >> '" + str(calls) + "'\n"
                        "[ \"$*\" = 'agent list' ] || exit 91\n" +
                        (command or "cat '" + str(fixture) + "'") + "\n")
        shim.chmod(0o700)
        start = time.monotonic()
        proc = subprocess.run([binary, "workflow", "agents", "--metadata-only", "--format", "json",
                               "--workspace", str(root / "absent-workspace")],
                              cwd=root, env=env, capture_output=True, timeout=4.5)
        elapsed = time.monotonic() - start
        assert proc.returncode == 0, (name, proc.returncode, proc.stderr)
        assert not proc.stderr, (name, proc.stderr)
        assert elapsed < 4.0, (name, elapsed)
        v = json.loads(proc.stdout)
        assert set(v) == {"kind", "schema_version", "observed_at_utc", "sources", "agents"}, v.keys()
        assert v["kind"] == "AgentStatusMetadataReadback@1" and v["schema_version"] == 1
        assert len(v["sources"]) == 1
        source = v["sources"][0]
        assert (source["status"], source["coverage"], len(v["agents"])) == (status, coverage, count), name
        assert source["source"] == "herdr.agent.list"
        assert source["source_observed_at_utc"] == (None if status == "unavailable" else v["observed_at_utc"])
        for agent in v["agents"]:
            assert set(agent) == {"pane", "provider", "status", "tab", "workspace"}
            assert agent["pane"]
            assert agent["provider"] in ("codex", "claude")
            assert agent["status"] in ("working", "idle", "done", "blocked", "unknown")
            for text in agent.values():
                assert isinstance(text, str) and len(text.encode("utf-16-le")) // 2 <= 256
                assert "never-export" not in text
        print(f"SM01 {name}: exit=0 source={status} coverage={coverage} agents={count} elapsed={elapsed:.3f}s margin_to_5s={5-elapsed:.3f}s", flush=True)
    assert calls.read_text().splitlines() == ["agent list"] * len(modes), "extra source operation or cache"
    assert not (root / ".ply").exists(), "readback created native state"
print("SM01 CLI conformance passed")
