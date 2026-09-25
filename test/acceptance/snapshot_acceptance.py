#!/usr/bin/env python3
"""Private validation and invocation support for snapshot acceptance."""

import argparse
import hashlib
import json
import os
import re
import stat
import subprocess
import sys
import time
from pathlib import Path, PurePosixPath


EXPECTED_ARGS = ("release", "--snapshot", "--clean", "--skip=publish")
CREDENTIALS = (
    "GITHUB_TOKEN",
    "GITLAB_TOKEN",
    "GITEA_TOKEN",
    "HOMEBREW_TAP_GITHUB_TOKEN",
    "SNAPCRAFT_STORE_CREDENTIALS",
)


class SnapshotError(Exception):
    pass


def sha256_file(path):
    digest = hashlib.sha256()
    with Path(path).open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def regular_file(path, executable=False):
    path = Path(path)
    try:
        value = path.lstat()
    except OSError as error:
        raise SnapshotError("missing file {}: {}".format(path, error))
    if not stat.S_ISREG(value.st_mode) or path.is_symlink():
        raise SnapshotError("path is not a non-symlink regular file: {}".format(path))
    if executable and not value.st_mode & 0o111:
        raise SnapshotError("path is not executable: {}".format(path))
    return value


def outside_repository(path, repo, label):
    resolved = Path(path).resolve()
    repository = Path(repo).resolve()
    try:
        inside = os.path.commonpath((str(resolved), str(repository))) == str(repository)
    except ValueError:
        inside = False
    if inside:
        raise SnapshotError("{} is repository-local: {}".format(label, resolved))


def write_json(path, value, exclusive=False):
    destination = Path(path)
    destination.parent.mkdir(parents=True, exist_ok=True)
    flags = os.O_WRONLY | os.O_CREAT
    flags |= os.O_EXCL if exclusive else os.O_TRUNC
    descriptor = os.open(destination, flags, 0o600)
    with os.fdopen(descriptor, "w", encoding="utf-8") as stream:
        json.dump(value, stream, indent=2, sort_keys=True)
        stream.write("\n")
        stream.flush()
        os.fsync(stream.fileno())


def read_json(path):
    regular_file(path)
    try:
        return json.loads(Path(path).read_text(encoding="utf-8"))
    except (OSError, UnicodeDecodeError, json.JSONDecodeError) as error:
        raise SnapshotError("invalid JSON {}: {}".format(path, error))


def parse_config(config):
    regular_file(config)
    text = Path(config).read_text(encoding="utf-8")
    versions = re.findall(r"goreleaser/v(\d+\.\d+\.\d+)/", text)
    binaries = re.findall(r"^[ \t]*-[ \t]+binary:[ \t]*([A-Za-z0-9._-]+)[ \t]*$", text, re.MULTILINE)
    if len(versions) != 1:
        raise SnapshotError("GoReleaser config must declare exactly one schema version")
    if len(binaries) != 1:
        raise SnapshotError("GoReleaser config must declare exactly one build binary")
    return versions[0], binaries[0]


def run_checked(arguments):
    run = subprocess.run(
        arguments,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        check=False,
    )
    if run.returncode != 0:
        raise SnapshotError(
            "command failed ({}): {}".format(run.returncode, " ".join(arguments))
        )
    return run.stdout, run.stderr


def parse_go_metadata(text):
    result = {"build": {}}
    for raw_line in text.splitlines():
        fields = raw_line.strip().split()
        if not fields:
            continue
        if fields[0] == "path" and len(fields) >= 2:
            result["path"] = fields[1]
        elif fields[0] == "mod" and len(fields) >= 3:
            result["module"] = fields[1]
            result["version"] = fields[2]
        elif fields[0] == "build" and len(fields) >= 2 and "=" in fields[1]:
            key, value = fields[1].split("=", 1)
            result["build"][key] = value
    return result


def proxy():
    real = os.environ.get("PLY_SNAPSHOT_REAL_GORELEASER", "")
    record = os.environ.get("PLY_SNAPSHOT_CALL_RECORD", "")
    if tuple(sys.argv[1:]) != EXPECTED_ARGS:
        raise SnapshotError("GoReleaser argv is not the exact non-publishing snapshot command")
    leaked = [name for name in CREDENTIALS if os.environ.get(name, "")]
    if leaked:
        raise SnapshotError("release credentials reached GoReleaser: " + ", ".join(leaked))
    if not real or not record:
        raise SnapshotError("GoReleaser proxy controls are incomplete")
    regular_file(real, executable=True)
    write_json(
        record,
        {
            "argv": list(EXPECTED_ARGS),
            "credentials": {name: "cleared" for name in CREDENTIALS},
            "invoked_at_ns": time.time_ns(),
            "tool": str(Path(real).resolve()),
        },
        exclusive=True,
    )
    os.execv(real, [real, *EXPECTED_ARGS])


def validate_tool(args):
    version, binary = parse_config(args.config)
    tool_state = regular_file(args.tool, executable=True)
    outside_repository(args.tool, args.repo, "GoReleaser tool")
    version_stdout, version_stderr = run_checked([args.tool, "--version"])
    metadata_stdout, metadata_stderr = run_checked([args.go, "version", "-m", args.tool])
    Path(args.version_output).write_text(version_stdout + version_stderr, encoding="utf-8")
    Path(args.metadata_output).write_text(metadata_stdout + metadata_stderr, encoding="utf-8")
    found_version = re.findall(r"^GitVersion:[ \t]+(\S+)$", version_stdout, re.MULTILINE)
    found_platform = re.findall(r"^Platform:[ \t]+(\S+)$", version_stdout, re.MULTILINE)
    if found_version != [version]:
        raise SnapshotError("GoReleaser version does not match accepted config schema")
    if found_platform != [args.host_goos + "/" + args.host_goarch]:
        raise SnapshotError("GoReleaser tool platform does not match the host")
    metadata = parse_go_metadata(metadata_stdout)
    if metadata.get("path") != "github.com/goreleaser/goreleaser/v2":
        raise SnapshotError("GoReleaser build metadata has the wrong module path")
    if metadata.get("module") != "github.com/goreleaser/goreleaser/v2" or metadata.get("version") != "v" + version:
        raise SnapshotError("GoReleaser build metadata has the wrong module version")
    if metadata["build"].get("GOOS") != args.host_goos or metadata["build"].get("GOARCH") != args.host_goarch:
        raise SnapshotError("GoReleaser build metadata has the wrong platform")
    write_json(
        args.output,
        {
            "accepted_config_version": version,
            "configured_binary": binary,
            "go_build_metadata_sha256": sha256_file(args.metadata_output),
            "host_goarch": args.host_goarch,
            "host_goos": args.host_goos,
            "tool_path": str(Path(args.tool).resolve()),
            "tool_sha256": sha256_file(args.tool),
            "tool_size": tool_state.st_size,
            "version_output_sha256": sha256_file(args.version_output),
        },
    )


def validate_call(args):
    record = read_json(args.record)
    required = {"argv", "credentials", "invoked_at_ns", "tool"}
    if set(record) != required:
        raise SnapshotError("GoReleaser call record has unexpected fields")
    if record["argv"] != list(EXPECTED_ARGS):
        raise SnapshotError("GoReleaser call record has the wrong argv")
    if record["credentials"] != {name: "cleared" for name in CREDENTIALS}:
        raise SnapshotError("GoReleaser call record does not prove credential clearing")
    if record["tool"] != str(Path(args.tool).resolve()):
        raise SnapshotError("GoReleaser call record substituted the tool")
    if not isinstance(record["invoked_at_ns"], int) or record["invoked_at_ns"] < args.started * 1_000_000_000:
        raise SnapshotError("GoReleaser call record predates this run")


def no_symlink_components(root, relative):
    current = Path(root)
    for part in PurePosixPath(relative).parts:
        current = current / part
        try:
            value = current.lstat()
        except OSError as error:
            raise SnapshotError("artifact path is missing: {}".format(error))
        if stat.S_ISLNK(value.st_mode):
            raise SnapshotError("artifact path contains a symlink: {}".format(current))
    return current.lstat()


def fresh_regular(path, started, label):
    value = regular_file(path)
    if value.st_mtime < started:
        raise SnapshotError("{} predates this snapshot run".format(label))
    return value


def validate_artifact(args):
    _, configured_binary = parse_config(args.config)
    outside_repository(args.source, args.repo, "snapshot source")
    source = Path(args.source).resolve()
    dist = source / "dist"
    if dist.is_symlink() or not dist.is_dir():
        raise SnapshotError("snapshot dist is not a regular directory")
    artifacts_path = dist / "artifacts.json"
    metadata_path = dist / "metadata.json"
    fresh_regular(artifacts_path, args.started, "artifacts metadata")
    fresh_regular(metadata_path, args.started, "release metadata")
    artifacts = read_json(artifacts_path)
    metadata = read_json(metadata_path)
    if not isinstance(artifacts, list) or not artifacts:
        raise SnapshotError("artifacts metadata has an empty population")
    if metadata.get("commit") != args.commit:
        raise SnapshotError("release metadata commit does not match the commanded source")
    if metadata.get("runtime") != {"goos": args.host_goos, "goarch": args.host_goarch}:
        raise SnapshotError("release metadata runtime does not match the host")
    binaries = [value for value in artifacts if isinstance(value, dict) and value.get("type") == "Binary"]
    candidates = [
        value for value in binaries
        if value.get("goos") == args.host_goos and value.get("goarch") == args.host_goarch
    ]
    if len(candidates) != 1:
        raise SnapshotError("host binary population is {}, expected exactly 1".format(len(candidates)))
    candidate = candidates[0]
    extra = candidate.get("extra")
    if not isinstance(extra, dict) or extra.get("Binary") != configured_binary:
        raise SnapshotError("host binary metadata does not match the configured binary")
    build_id = extra.get("ID")
    target = candidate.get("target")
    name = candidate.get("name")
    if not all(isinstance(value, str) and value for value in (build_id, target, name)):
        raise SnapshotError("host binary metadata is incomplete")
    target_prefix = args.host_goos + "_" + args.host_goarch
    if target != target_prefix and not target.startswith(target_prefix + "_"):
        raise SnapshotError("host binary target has the wrong platform path")
    expected_name = configured_binary + (".exe" if args.host_goos == "windows" else "")
    if name != expected_name:
        raise SnapshotError("host binary has the wrong artifact name")
    expected_relative = "dist/{0}_{1}/{2}".format(build_id, target, name)
    relative = candidate.get("path")
    if not isinstance(relative, str) or PurePosixPath(relative).as_posix() != relative:
        raise SnapshotError("host binary path is not normalized")
    if relative != expected_relative:
        raise SnapshotError("host binary path does not match its metadata")
    artifact = source / PurePosixPath(relative)
    artifact_state = no_symlink_components(source, relative)
    if not stat.S_ISREG(artifact_state.st_mode):
        raise SnapshotError("host artifact is not a regular file")
    if not artifact_state.st_mode & 0o111:
        raise SnapshotError("host artifact is not executable")
    if artifact_state.st_size <= 0:
        raise SnapshotError("host artifact is empty")
    if artifact_state.st_mtime < args.started:
        raise SnapshotError("host artifact predates this snapshot run")
    build_stdout, build_stderr = run_checked([args.go, "version", "-m", str(artifact)])
    if build_stderr:
        raise SnapshotError("go version -m wrote stderr for the host artifact")
    build = parse_go_metadata(build_stdout)["build"]
    if build.get("GOOS") != args.host_goos or build.get("GOARCH") != args.host_goarch:
        raise SnapshotError("host artifact build metadata has the wrong platform")
    if build.get("vcs.revision") != args.commit or build.get("vcs.modified") != "false":
        raise SnapshotError("host artifact is not bound to the clean commanded commit")
    write_json(
        args.output,
        {
            "artifact": {
                "build_id": build_id,
                "path": str(artifact),
                "relative_path": relative,
                "sha256": sha256_file(artifact),
                "size": artifact_state.st_size,
                "target": target,
            },
            "binary_records": len(binaries),
            "host_candidate_records": len(candidates),
            "metadata_sha256": sha256_file(metadata_path),
            "source_commit": args.commit,
        },
    )


def file_identity(path):
    value = regular_file(path, executable=True)
    return {
        "device": value.st_dev,
        "inode": value.st_ino,
        "mode": stat.S_IMODE(value.st_mode),
        "path": str(Path(path).resolve()),
        "sha256": sha256_file(path),
        "size": value.st_size,
    }


def identity(args):
    write_json(args.output, file_identity(args.path))


def validate_terminal(args):
    regular_file(args.stdout)
    lines = Path(args.stdout).read_text(encoding="utf-8").splitlines()
    terminal = "verify-{}: PASS".format(args.name)
    if lines.count(terminal) != 1 or not lines or lines[-1] != terminal:
        raise SnapshotError("{} verifier lacks one terminal PASS record".format(args.name))
    write_json(
        args.output,
        {
            "name": args.name,
            "stdout_sha256": sha256_file(args.stdout),
            "terminal": terminal,
        },
    )


def validate_verifier(args):
    validate_terminal(args)
    regular_file(args.trace)
    artifact = str(Path(args.artifact).resolve())
    trace_lines = Path(args.trace).read_text(encoding="utf-8", errors="replace").splitlines()
    traced = [line for line in trace_lines if line.startswith("TRACE:")]
    behavioral = [
        line for line in traced
        if artifact + " " + args.name + " " in line and "--target" in line and "--help" not in line
    ]
    if not behavioral:
        raise SnapshotError(
            "{} verifier trace has no non-help command against the snapshot artifact".format(args.name)
        )
    value = read_json(args.output)
    value.update(
        {
            "artifact_path": artifact,
            "artifact_sha256": sha256_file(artifact),
            "behavioral_trace_records": len(behavioral),
            "trace_sha256": sha256_file(args.trace),
        }
    )
    write_json(args.output, value)


def field(args):
    value = read_json(args.path)
    for name in args.field.split("."):
        if not isinstance(value, dict) or name not in value:
            raise SnapshotError("missing JSON field: " + args.field)
        value = value[name]
    if isinstance(value, (dict, list)):
        print(json.dumps(value, sort_keys=True, separators=(",", ":")))
    elif isinstance(value, bool):
        print("true" if value else "false")
    elif value is None:
        print("null")
    else:
        print(value)


def parser():
    value = argparse.ArgumentParser()
    commands = value.add_subparsers(dest="command", required=True)

    tool = commands.add_parser("validate-tool")
    tool.add_argument("--config", required=True)
    tool.add_argument("--tool", required=True)
    tool.add_argument("--go", required=True)
    tool.add_argument("--repo", required=True)
    tool.add_argument("--host-goos", required=True)
    tool.add_argument("--host-goarch", required=True)
    tool.add_argument("--version-output", required=True)
    tool.add_argument("--metadata-output", required=True)
    tool.add_argument("--output", required=True)

    call = commands.add_parser("validate-call")
    call.add_argument("--record", required=True)
    call.add_argument("--tool", required=True)
    call.add_argument("--started", required=True, type=int)

    artifact = commands.add_parser("validate-artifact")
    artifact.add_argument("--repo", required=True)
    artifact.add_argument("--source", required=True)
    artifact.add_argument("--config", required=True)
    artifact.add_argument("--go", required=True)
    artifact.add_argument("--host-goos", required=True)
    artifact.add_argument("--host-goarch", required=True)
    artifact.add_argument("--commit", required=True)
    artifact.add_argument("--started", required=True, type=int)
    artifact.add_argument("--output", required=True)

    identity_parser = commands.add_parser("identity")
    identity_parser.add_argument("--path", required=True)
    identity_parser.add_argument("--output", required=True)

    terminal = commands.add_parser("validate-terminal")
    terminal.add_argument("--name", required=True)
    terminal.add_argument("--stdout", required=True)
    terminal.add_argument("--output", required=True)

    verifier = commands.add_parser("validate-verifier")
    verifier.add_argument("--name", required=True)
    verifier.add_argument("--artifact", required=True)
    verifier.add_argument("--stdout", required=True)
    verifier.add_argument("--trace", required=True)
    verifier.add_argument("--output", required=True)

    field_parser = commands.add_parser("field")
    field_parser.add_argument("--path", required=True)
    field_parser.add_argument("--field", required=True)
    return value


def main():
    if len(sys.argv) > 1 and sys.argv[1] == "release":
        proxy()
        return 0
    args = parser().parse_args()
    if args.command == "validate-tool":
        validate_tool(args)
    elif args.command == "validate-call":
        validate_call(args)
    elif args.command == "validate-artifact":
        validate_artifact(args)
    elif args.command == "identity":
        identity(args)
    elif args.command == "validate-terminal":
        validate_terminal(args)
    elif args.command == "validate-verifier":
        validate_verifier(args)
    elif args.command == "field":
        field(args)
    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except Exception as error:
        print("snapshot acceptance support: {}".format(error), file=sys.stderr)
        sys.exit(1)
