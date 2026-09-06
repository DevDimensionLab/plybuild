#!/usr/bin/env python3
"""Private validation, Docker proxy, and runtime bridge for image acceptance."""

import argparse
import datetime
import hashlib
import json
import os
import re
import stat
import subprocess
import sys
import tempfile
import time
from pathlib import Path


INSPECTION_SCRIPT = (
    "path=/bin/ply; "
    "test -e \"$path\" || exit 41; "
    "test ! -L \"$path\" || exit 42; "
    "test -f \"$path\" || exit 43; "
    "test -x \"$path\" || exit 44; "
    "printf 'path=/bin/ply\\ntype=regular\\n'; "
    "stat -c 'mode=%a' \"$path\"; "
    "stat -c 'size=%s' \"$path\"; "
    "sha256sum \"$path\" | awk '{print \"sha256=\" $1}'"
)
SHA256_RE = re.compile(r"^sha256:[0-9a-f]{64}$")
FILE_SHA_RE = re.compile(r"^[0-9a-f]{64}$")


class DockerAcceptanceError(Exception):
    pass


def sha256_file(path):
    digest = hashlib.sha256()
    with Path(path).open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def regular_file(path, executable=False):
    value_path = Path(path)
    try:
        value = value_path.lstat()
    except OSError as error:
        raise DockerAcceptanceError("missing file {}: {}".format(value_path, error))
    if not stat.S_ISREG(value.st_mode) or value_path.is_symlink():
        raise DockerAcceptanceError("path is not a non-symlink regular file: {}".format(value_path))
    if executable and not value.st_mode & 0o111:
        raise DockerAcceptanceError("path is not executable: {}".format(value_path))
    return value


def regular_directory(path):
    value_path = Path(path)
    value = value_path.lstat()
    if not stat.S_ISDIR(value.st_mode) or value_path.is_symlink():
        raise DockerAcceptanceError("path is not a non-symlink directory: {}".format(value_path))
    return value_path.resolve()


def outside_repository(path, repo, label):
    resolved = Path(path).resolve()
    repository = Path(repo).resolve()
    try:
        inside = os.path.commonpath((str(resolved), str(repository))) == str(repository)
    except ValueError:
        inside = False
    if inside:
        raise DockerAcceptanceError("{} is repository-local: {}".format(label, resolved))


def write_json(path, value, exclusive=False):
    destination = Path(path)
    destination.parent.mkdir(parents=True, exist_ok=True)
    flags = os.O_WRONLY | os.O_CREAT | (os.O_EXCL if exclusive else os.O_TRUNC)
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
        raise DockerAcceptanceError("invalid JSON {}: {}".format(path, error))


def read_text(path):
    regular_file(path)
    return Path(path).read_text(encoding="utf-8")


def required_environment(name):
    value = os.environ.get(name, "")
    if not value:
        raise DockerAcceptanceError("required environment is missing: " + name)
    return value


def append_call(record, argv):
    destination = Path(record)
    if destination.exists() or destination.is_symlink():
        regular_file(destination)
    else:
        regular_directory(destination.parent)
    line = json.dumps(
        {"argv": argv, "invoked_at_ns": time.time_ns()},
        sort_keys=True,
        separators=(",", ":"),
    ) + "\n"
    descriptor = os.open(destination, os.O_WRONLY | os.O_CREAT | os.O_APPEND, 0o600)
    with os.fdopen(descriptor, "a", encoding="utf-8") as stream:
        stream.write(line)
        stream.flush()
        os.fsync(stream.fileno())


def validate_docker_argv(argv):
    if not argv:
        raise DockerAcceptanceError("Docker argv is empty")
    if argv[0] in {"login", "logout", "push", "tag"} or "--push" in argv:
        raise DockerAcceptanceError("Docker publication operation is forbidden")

    tag = required_environment("PLY_DOCKER_EXPECTED_TAG")
    context = required_environment("PLY_DOCKER_EXPECTED_CONTEXT")
    run_id = required_environment("PLY_DOCKER_ACCEPTANCE_RUN_ID")
    image_id = os.environ.get("PLY_DOCKER_EXPECTED_IMAGE_ID", "")
    selected_context = required_environment("PLY_DOCKER_EXPECTED_CONTEXT_NAME")
    exact_queries = (
        ["version", "--format", "{{json .}}"],
        ["info", "--format", "{{json .}}"],
        ["context", "show"],
        ["context", "inspect", selected_context],
        ["buildx", "version"],
        ["buildx", "inspect"],
    )
    if argv in exact_queries:
        return
    if argv == [
        "build", "--no-cache", "--force-rm", "--tag", tag,
        "--build-arg", "PLY_ACCEPTANCE_RUN_ID=" + run_id, context,
    ]:
        return
    if argv == ["image", "inspect", tag]:
        return
    if image_id and argv == ["image", "inspect", image_id]:
        return
    if argv == ["image", "ls", "--no-trunc", "--quiet", tag]:
        return
    if image_id and argv == ["history", "--no-trunc", "--format", "{{json .}}", image_id]:
        return
    if argv[0] == "run":
        expected = json.loads(required_environment("PLY_DOCKER_EXPECTED_RUN_ARGV"))
        if argv != expected:
            raise DockerAcceptanceError("Docker run argv is not the exact accepted invocation")
        return
    if argv[:2] == ["container", "inspect"]:
        cid = required_environment("PLY_DOCKER_EXPECTED_CONTAINER_ID")
        if argv == ["container", "inspect", cid, "--format", "{{.Image}}"]:
            return
    if argv[:2] == ["container", "rm"]:
        cid = required_environment("PLY_DOCKER_EXPECTED_CONTAINER_ID")
        if argv == ["container", "rm", cid]:
            return
    raise DockerAcceptanceError("Docker argv is outside the acceptance allowlist: " + repr(argv))


def docker_proxy(argv):
    real = required_environment("PLY_DOCKER_REAL_BIN")
    record = required_environment("PLY_DOCKER_CALL_RECORD")
    regular_file(real, executable=True)
    validate_docker_argv(argv)
    append_call(record, argv)
    os.execve(real, [real, *argv], os.environ.copy())


def proxy_run(argv, environment=None, capture=False):
    helper = str(Path(__file__).resolve())
    run_environment = os.environ.copy()
    if environment:
        run_environment.update(environment)
    return subprocess.run(
        [helper, "docker", *argv],
        env=run_environment,
        stdout=subprocess.PIPE if capture else None,
        stderr=subprocess.PIPE if capture else None,
        text=True,
        check=False,
    )


def runtime_bridge(command_args):
    repo = required_environment("PLY_DOCKER_REPO_ROOT")
    runtime_root = regular_directory(required_environment("PLY_DOCKER_RUNTIME_ROOT"))
    outside_repository(runtime_root, repo, "runtime receipt root")
    image_id = required_environment("PLY_DOCKER_EXPECTED_IMAGE_ID")
    if not SHA256_RE.fullmatch(image_id):
        raise DockerAcceptanceError("runtime image ID is malformed")
    flow = required_environment("PLY_DOCKER_EXPECTED_FLOW")
    if not command_args or command_args[0] != flow:
        raise DockerAcceptanceError("runtime command does not match the selected verifier")

    home_text = required_environment("HOME")
    home = Path(home_text).resolve()
    mount_root = home.parent
    regular_directory(mount_root)
    outside_repository(mount_root, repo, "verifier mount root")
    try:
        home.relative_to(mount_root)
    except ValueError as error:
        raise DockerAcceptanceError("HOME escapes the verifier mount root") from error

    descriptor, cid_name = tempfile.mkstemp(prefix="container-", suffix=".cid", dir=runtime_root)
    os.close(descriptor)
    os.unlink(cid_name)
    run_argv = [
        "run",
        "--cidfile",
        cid_name,
        "--add-host",
        "host.docker.internal:host-gateway",
        "--volume",
        "{0}:{0}".format(mount_root),
        "--env",
        "HOME={}".format(home),
        "--workdir",
        str(mount_root),
        image_id,
        *command_args,
    ]
    run = proxy_run(
        run_argv,
        {"PLY_DOCKER_EXPECTED_RUN_ARGV": json.dumps(run_argv, separators=(",", ":"))},
    )
    if not Path(cid_name).is_file():
        raise DockerAcceptanceError("Docker run produced no container ID receipt")
    container_id = Path(cid_name).read_text(encoding="utf-8").strip()
    if not re.fullmatch(r"[0-9a-f]{64}", container_id):
        raise DockerAcceptanceError("Docker returned a malformed container ID")

    container_environment = {"PLY_DOCKER_EXPECTED_CONTAINER_ID": container_id}
    inspect = proxy_run(
        ["container", "inspect", container_id, "--format", "{{.Image}}"],
        container_environment,
        capture=True,
    )
    actual_image = inspect.stdout.strip() if inspect.returncode == 0 else ""
    removed = proxy_run(["container", "rm", container_id], container_environment, capture=True)
    support_valid = inspect.returncode == 0 and removed.returncode == 0 and actual_image == image_id

    receipt_fd, receipt_name = tempfile.mkstemp(prefix="receipt-", suffix=".json", dir=runtime_root)
    os.close(receipt_fd)
    write_json(
        receipt_name,
        {
            "command_args": command_args,
            "container_id": container_id,
            "container_image": actual_image,
            "home": str(home),
            "image_id": image_id,
            "mount_root": str(mount_root),
            "run_argv": run_argv,
            "run_exit": run.returncode,
            "support_valid": support_valid,
        },
    )
    if not support_valid:
        raise DockerAcceptanceError("runtime container did not use the immutable image ID")
    return run.returncode


def validate_daemon(args):
    version = read_json(args.version)
    info = read_json(args.info)
    context_name = read_text(args.context_show).strip()
    contexts = read_json(args.context_inspect)
    buildx_version = read_text(args.buildx_version).strip()
    buildx_inspect = read_text(args.buildx_inspect)
    if not isinstance(version.get("Client"), dict) or not isinstance(version.get("Server"), dict):
        raise DockerAcceptanceError("Docker version lacks real Client and Server identities")
    if not version["Client"].get("Version") or not version["Server"].get("Version"):
        raise DockerAcceptanceError("Docker Client or Server version is empty")
    if not isinstance(info, dict) or not info.get("ID") or not info.get("ServerVersion"):
        raise DockerAcceptanceError("Docker info lacks a server identity")
    if info.get("OSType") != "linux":
        raise DockerAcceptanceError("Docker daemon OSType is not linux")
    architecture = {"aarch64": "arm64", "x86_64": "amd64"}.get(
        info.get("Architecture"), info.get("Architecture")
    )
    if architecture not in {"arm64", "amd64"}:
        raise DockerAcceptanceError("Docker daemon architecture is unsupported")
    if not context_name or not isinstance(contexts, list) or len(contexts) != 1:
        raise DockerAcceptanceError("selected Docker context identity is ambiguous")
    context = contexts[0]
    if context.get("Name") != context_name:
        raise DockerAcceptanceError("Docker context inspection does not match the selected context")
    endpoint = context.get("Endpoints", {}).get("docker", {}).get("Host", "")
    if not endpoint.startswith("unix://"):
        raise DockerAcceptanceError("Docker context endpoint is not a local Unix socket")
    if not buildx_version or not re.search(r"^Driver:\s+docker\s*$", buildx_inspect, re.MULTILINE):
        raise DockerAcceptanceError("Docker builder identity is missing or not daemon-backed")
    platforms = re.findall(r"^Platforms:\s+(.+)$", buildx_inspect, re.MULTILINE)
    if len(platforms) != 1 or "linux/" + architecture not in platforms[0].split(", "):
        raise DockerAcceptanceError("Docker builder does not support the daemon platform")
    write_json(
        args.output,
        {
            "architecture": architecture,
            "buildx_version": buildx_version,
            "builder_driver": "docker",
            "builder_platforms": platforms[0].split(", "),
            "client_version": version["Client"]["Version"],
            "context": context_name,
            "endpoint": endpoint,
            "host_id": info["ID"],
            "host_name": info.get("Name", ""),
            "os": "linux",
            "server_version": version["Server"]["Version"],
        },
    )


def parse_created(value):
    if value.endswith("Z"):
        value = value[:-1] + "+00:00"
    return datetime.datetime.fromisoformat(value).timestamp()


def validate_image(args):
    values = read_json(args.inspect)
    if not isinstance(values, list) or len(values) != 1:
        raise DockerAcceptanceError("image inspection population is not exactly one")
    value = values[0]
    image_id = value.get("Id", "")
    if not SHA256_RE.fullmatch(image_id):
        raise DockerAcceptanceError("image ID is malformed")
    tag_ids = [line.strip() for line in read_text(args.tag_ids).splitlines() if line.strip()]
    if tag_ids != [image_id]:
        raise DockerAcceptanceError("tag resolution is missing, ambiguous, or stale")
    if value.get("RepoTags") != [args.tag]:
        raise DockerAcceptanceError("image has an unexpected tag population")
    if value.get("Os") != args.os or value.get("Architecture") != args.architecture:
        raise DockerAcceptanceError("image platform does not match the selected daemon platform")
    config = value.get("Config")
    if not isinstance(config, dict) or config.get("Entrypoint") != ["/bin/ply"]:
        raise DockerAcceptanceError("image entrypoint is not exactly /bin/ply")
    if config.get("Cmd") not in (None, []):
        raise DockerAcceptanceError("image carries an unexpected default command")
    if config.get("User") not in (None, "", "0") or config.get("WorkingDir") not in (None, "", "/"):
        raise DockerAcceptanceError("image user or working-directory configuration is unexpected")
    environment = config.get("Env", []) or []
    if any(item.startswith(("GOOS=", "GOARCH=", "CGO_ENABLED=")) for item in environment):
        raise DockerAcceptanceError("builder-only Go configuration leaked into the final image")
    if not isinstance(value.get("Size"), int) or value["Size"] <= 0:
        raise DockerAcceptanceError("image size is empty")
    if not isinstance(value.get("RootFS", {}).get("Layers"), list) or not value["RootFS"]["Layers"]:
        raise DockerAcceptanceError("image layer population is empty")
    created = value.get("Created", "")
    try:
        created_epoch = parse_created(created)
    except (TypeError, ValueError) as error:
        raise DockerAcceptanceError("image creation time is malformed") from error
    if created_epoch + 1 < args.started:
        raise DockerAcceptanceError("image predates the commanded build")
    write_json(
        args.output,
        {
            "architecture": value["Architecture"],
            "created": created,
            "descriptor_digest": value.get("Descriptor", {}).get("digest", ""),
            "entrypoint": config["Entrypoint"],
            "environment": environment,
            "image_id": image_id,
            "os": value["Os"],
            "repo_digests": value.get("RepoDigests", []),
            "size": value["Size"],
            "tag": args.tag,
        },
    )


def validate_image_file(args):
    fields = {}
    for line in read_text(args.path).splitlines():
        if "=" not in line:
            raise DockerAcceptanceError("image file inspection contains a malformed line")
        key, value = line.split("=", 1)
        if key in fields:
            raise DockerAcceptanceError("image file inspection contains duplicate fields")
        fields[key] = value
    if set(fields) != {"path", "type", "mode", "size", "sha256"}:
        raise DockerAcceptanceError("image file inspection fields are incomplete")
    if fields["path"] != "/bin/ply" or fields["type"] != "regular":
        raise DockerAcceptanceError("/bin/ply is missing, symlinked, or not regular")
    try:
        mode = int(fields["mode"], 8)
        size = int(fields["size"])
    except ValueError as error:
        raise DockerAcceptanceError("/bin/ply mode or size is malformed") from error
    if not mode & 0o111 or size <= 0 or not FILE_SHA_RE.fullmatch(fields["sha256"]):
        raise DockerAcceptanceError("/bin/ply is not a non-empty executable with a digest")
    write_json(args.output, {"mode": fields["mode"], "path": "/bin/ply", "sha256": fields["sha256"], "size": size, "type": "regular"})


def file_identity(args):
    value = regular_file(args.path, executable=True)
    write_json(
        args.output,
        {
            "mode": stat.S_IMODE(value.st_mode),
            "path": str(Path(args.path).resolve()),
            "sha256": sha256_file(args.path),
            "size": value.st_size,
        },
    )


def validate_terminal(args):
    lines = read_text(args.stdout).splitlines()
    terminal = "verify-{}: PASS".format(args.name)
    if lines.count(terminal) != 1 or not lines or lines[-1] != terminal:
        raise DockerAcceptanceError("{} verifier lacks one terminal PASS record".format(args.name))
    write_json(args.output, {"name": args.name, "stdout_sha256": sha256_file(args.stdout), "terminal": terminal})


def validate_runtime(args):
    root = regular_directory(args.root)
    receipts = sorted(root.glob("receipt-*.json"))
    if len(receipts) != 3:
        raise DockerAcceptanceError("{} runtime receipt population is {}, expected 3".format(args.flow, len(receipts)))
    values = [read_json(path) for path in receipts]
    container_ids = []
    commands = []
    for value in values:
        if value.get("support_valid") is not True:
            raise DockerAcceptanceError("runtime receipt did not validate its support path")
        if value.get("image_id") != args.image_id or value.get("container_image") != args.image_id:
            raise DockerAcceptanceError("runtime receipt substituted the immutable image")
        run_argv = value.get("run_argv")
        if not isinstance(run_argv, list) or args.image_id not in run_argv or args.tag in run_argv:
            raise DockerAcceptanceError("runtime argv used a tag or omitted the immutable image ID")
        commands.append(value.get("command_args"))
        container_ids.append(value.get("container_id"))
    if len(set(container_ids)) != 3:
        raise DockerAcceptanceError("runtime container identity population is duplicated")
    help_commands = [value for value in commands if value == [args.flow, "--help"]]
    behavior = [value for value in commands if isinstance(value, list) and value and value[0] == args.flow and "--target" in value and "--help" not in value]
    if len(help_commands) != 1 or len(behavior) != 2:
        raise DockerAcceptanceError("verifier lacks the exact help and non-help runtime population")
    write_json(
        args.output,
        {
            "behavioral_runtime_receipts": len(behavior),
            "container_ids": container_ids,
            "flow": args.flow,
            "image_id": args.image_id,
            "receipt_count": len(receipts),
            "receipts_sha256": [sha256_file(path) for path in receipts],
        },
    )


def validate_calls(args):
    regular_file(args.path)
    calls = []
    for line in Path(args.path).read_text(encoding="utf-8").splitlines():
        value = json.loads(line)
        if set(value) != {"argv", "invoked_at_ns"} or not isinstance(value["argv"], list):
            raise DockerAcceptanceError("Docker call record is malformed")
        calls.append(value)
    if not calls:
        raise DockerAcceptanceError("Docker call record population is empty")
    build = [value for value in calls if value["argv"][:1] == ["build"]]
    runs = [value for value in calls if value["argv"][:1] == ["run"]]
    expected_build = [
        "build", "--no-cache", "--force-rm", "--tag", args.tag,
        "--build-arg", "PLY_ACCEPTANCE_RUN_ID=" + args.run_id, args.context,
    ]
    if len(build) != 1 or build[0]["argv"] != expected_build:
        raise DockerAcceptanceError("Docker build did not run exactly once with accepted argv")
    if len(runs) != 10:
        raise DockerAcceptanceError("Docker run population is {}, expected 10".format(len(runs)))
    if any(args.image_id not in value["argv"] or args.tag in value["argv"] for value in runs):
        raise DockerAcceptanceError("a Docker run did not use only the immutable image ID")
    if any(value["argv"][0] in {"login", "push", "tag"} or "--push" in value["argv"] for value in calls):
        raise DockerAcceptanceError("Docker publication call reached the command record")
    write_json(args.output, {"build_calls": 1, "call_count": len(calls), "run_calls": len(runs), "sha256": sha256_file(args.path)})


def field(args):
    value = read_json(args.path)
    for name in args.field.split("."):
        if not isinstance(value, dict) or name not in value:
            raise DockerAcceptanceError("missing JSON field: " + args.field)
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

    daemon = commands.add_parser("validate-daemon")
    for name in ("version", "info", "context-show", "context-inspect", "buildx-version", "buildx-inspect", "output"):
        daemon.add_argument("--" + name, required=True)

    image = commands.add_parser("validate-image")
    for name in ("inspect", "tag-ids", "tag", "os", "architecture", "output"):
        image.add_argument("--" + name, required=True)
    image.add_argument("--started", required=True, type=int)

    image_file = commands.add_parser("validate-image-file")
    image_file.add_argument("--path", required=True)
    image_file.add_argument("--output", required=True)

    identity = commands.add_parser("identity")
    identity.add_argument("--path", required=True)
    identity.add_argument("--output", required=True)

    terminal = commands.add_parser("validate-terminal")
    terminal.add_argument("--name", required=True)
    terminal.add_argument("--stdout", required=True)
    terminal.add_argument("--output", required=True)

    runtime = commands.add_parser("validate-runtime")
    for name in ("root", "flow", "image-id", "tag", "output"):
        runtime.add_argument("--" + name, required=True)

    calls = commands.add_parser("validate-calls")
    for name in ("path", "tag", "context", "run-id", "image-id", "output"):
        calls.add_argument("--" + name, required=True)

    field_parser = commands.add_parser("field")
    field_parser.add_argument("--path", required=True)
    field_parser.add_argument("--field", required=True)
    return value


def main():
    if len(sys.argv) > 1 and sys.argv[1] == "docker":
        docker_proxy(sys.argv[2:])
        return 0
    support_commands = {
        "validate-daemon", "validate-image", "validate-image-file", "identity",
        "validate-terminal", "validate-runtime", "validate-calls", "field",
    }
    if len(sys.argv) > 1 and sys.argv[1] not in support_commands:
        return runtime_bridge(sys.argv[1:])
    args = parser().parse_args()
    if args.command == "validate-daemon":
        validate_daemon(args)
    elif args.command == "validate-image":
        validate_image(args)
    elif args.command == "validate-image-file":
        validate_image_file(args)
    elif args.command == "identity":
        file_identity(args)
    elif args.command == "validate-terminal":
        validate_terminal(args)
    elif args.command == "validate-runtime":
        validate_runtime(args)
    elif args.command == "validate-calls":
        validate_calls(args)
    elif args.command == "field":
        field(args)
    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except Exception as error:
        print("docker acceptance support: {}".format(error), file=sys.stderr)
        sys.exit(1)
