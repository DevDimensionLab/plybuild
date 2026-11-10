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
        ("envelope-truncation", dict(envelope([]), truncated=True), "partial", "partial", 0, None),
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
    # r001-f01: a malformed duplicate still invalidates the valid pane claim.
    malformed = {}
    for name, field in (("invalid-status", "agent_status"), ("null-tab", "tab_id"), ("invalid-provider", "agent")):
        bad = row()
        bad[field] = None if name == "null-tab" else 42
        malformed[name] = json.dumps(bad)
    malformed["duplicate-status-before-pane"] = '{"agent":"codex","agent_status":"idle","agent_status":42,"pane_id":"w1:p1","tab_id":"t","workspace_id":"w"}'
    malformed["duplicate-status-after-pane"] = '{"pane_id":"w1:p1","agent":"codex","agent_status":"idle","agent_status":42,"tab_id":"t","workspace_id":"w"}'
    malformed["duplicate-provider-before-pane"] = '{"agent":"codex","agent":42,"pane_id":"w1:p1","tab_id":"t","workspace_id":"w","agent_status":"idle"}'
    malformed["duplicate-pane-same-claim"] = '{"agent":"codex","pane_id":"w1:p1","pane_id":"w1:p1","tab_id":"t","workspace_id":"w","agent_status":"idle"}'
    malformed["duplicate-pane-different-claims"] = '{"agent":"codex","pane_id":"w1:p1","pane_id":"w1:p2","tab_id":"t","workspace_id":"w","agent_status":"idle"}'
    for name, bad in malformed.items():
        for reverse in (False, True):
            rows = [json.dumps(row()), bad]
            if reverse:
                rows.reverse()
            rows.append(json.dumps(row("w1:p-safe")))
            raw = '{"result":{"type":"agent_list","agents":[' + ','.join(rows) + ']}}'
            modes.append((f"malformed-duplicate-{name}-reverse-{reverse}", raw, "partial", "partial", 1, None))
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
        assert elapsed < 2.6, (name, elapsed)
        v = json.loads(proc.stdout)
        assert set(v) == {"kind", "schema_version", "observed_at_utc", "sources", "agents"}, v.keys()
        assert v["kind"] == "AgentStatusMetadataReadback@1" and v["schema_version"] == 1
        assert len(v["sources"]) == 1
        source = v["sources"][0]
        assert (source["status"], source["coverage"], len(v["agents"])) == (status, coverage, count), name
        assert source["source"] == "herdr.agent.list"
        assert source["source_observed_at_utc"] == (None if status == "unavailable" else v["observed_at_utc"])
        if name.startswith("malformed-duplicate-"):
            assert [a["pane"] for a in v["agents"]] == ["w1:p-safe"], name
        for agent in v["agents"]:
            assert set(agent) == {"pane", "provider", "status", "tab", "workspace"}
            assert agent["pane"]
            assert agent["provider"] in ("codex", "claude")
            assert agent["status"] in ("working", "idle", "done", "blocked", "unknown")
            for text in agent.values():
                assert isinstance(text, str) and len(text.encode("utf-16-le")) // 2 <= 256
                assert "never-export" not in text
        print(f"SM01 {name}: exit=0 source={status} coverage={coverage} agents={count} elapsed={elapsed:.3f}s margin_to_2600ms={2.6-elapsed:.3f}s", flush=True)
    observed = calls.read_text().splitlines()
    assert observed.count("agent list") == len(modes), "extra list or cache"
    assert all(call in ("agent list", "tab list --workspace w1") for call in observed), "unexpected source operation"
    assert not (root / ".ply").exists(), "readback created native state"
print("SM01 CLI conformance passed")
