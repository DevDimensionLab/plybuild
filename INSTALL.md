# Install Ply

## Build and install from source

Clone the repository and install with the Makefile. `PLY_GO` is the absolute
path of a Go binary that reports the go.mod `toolchain` (Go 1.26.10):

```sh
git clone https://github.com/DevDimensionLab/plybuild.git
cd plybuild
PLY_GO=/absolute/go1.26.10/bin/go make install
ply version --format json
```

## Install a prebuilt package

A release reads only work-item registry format 11 (`registry_formats: [11]`).
A Project's registry is `<wrapper>/data/work/work-items.yaml`; a registry in any
other format is rejected without writes. Installation never reads, upgrades or
migrates a workspace registry.

A release contains `ply`, `ply-install`, `manifest.json` and this guide in a
`ply-<version>-darwin-arm64.tar.gz` archive. The publisher supplies the archive
SHA-256 and the **manifest SHA-256** alongside the package through the agreed
delivery channel. Retain that trusted manifest digest; a digest obtained only
from an untrusted download is not an authenticity check.

The consumer needs macOS on Apple Silicon and an existing destination directory.
No Go installation, source checkout, working old Ply or workspace access is
needed. Keep the release in a permanent directory, for example
`~/.local/share/ply/releases/<version>/`, for retries and other consumers. Do
not install from a maintainer's temporary build.

After verifying the archive's SHA-256 and extracting it, run the bundled helper.
Replace the example release path, digest and destination with exact values:

```sh
/absolute/releases/ply-v1.1.0-cb01a.1-darwin-arm64/ply-install \
  --package /absolute/releases/ply-v1.1.0-cb01a.1-darwin-arm64 \
  --sha256 sha256:<manifest-sha256> \
  --destination /absolute/bin/ply

/absolute/bin/ply version --format json
command -v ply
ply version --format json
```

`version` returns `PlyExecutable@1`: version, full source commit, dirty status,
OS, architecture, executable SHA-256 and readable registry formats. The absolute
and PATH-resolved executables must identify the same release. If the PATH points
elsewhere, put the intended destination directory first and clear the shell
command cache (`hash -r` in Bash, `rehash` in Zsh). The installer does not edit
shell settings.

The helper checks the caller-pinned manifest, both delivered executable digests,
Mach-O platform and candidate runtime identity. It serializes installations at
the destination. The old executable is preserved in `ply.backups/<sha256>` and
the operation is recorded in the sibling `ply.install.json` before replacement.
The replacement is fully written and synced in the same directory, then renamed
atomically. Exit zero and `PlyInstallReceipt@1` with `state: installed` mean the
installed executable was verified. Errors never claim a successful installation.
Installer files are independent of `.ply`, Git, queues and agent status.

## Retry and binary rollback

After an error or interruption, retain the package and rerun the **same command**.
The receipt distinguishes the exact old and new bytes. Recovery finishes the
pending installation if either known state is observed; unrelated bytes or a
missing/corrupt backup stop it. A different package cannot replace an unresolved
pending installation. Temporary files left by an abrupt process kill are never
executed or mistaken for the installed binary; they may be removed after recovery.

Restore the preserved executable with its exact currently installed digest:

```sh
/absolute/releases/ply-v1.1.0-cb01a.1-darwin-arm64/ply-install \
  --rollback --sha256 sha256:<installed-ply-sha256> \
  --destination /absolute/bin/ply
```

The receipt returns `state: rolled_back`. Rollback changes only executable bytes.
It does not downgrade a registry, rewrite authorization, or make an old reader
understand newer fields. A previous binary that could not read the workspace
still cannot read it. Reinstall the compatible package to resume ordinary work.
An installation without a previous binary has no binary rollback target.

## Maintainer release

From a clean committed checkout, build the consumer package:

```sh
PLY_GO=/absolute/go1.26.10/bin/go python3 scripts/package-cli v1.1.0-cb01a.1 /absolute/new-release-directory
```

The output directory must be absent and outside the checkout. Both executables
are built from the observed full source commit with `CGO_ENABLED=0`, `-trimpath`
and VCS metadata. The source is checked again before publication. The release
directory contains the extracted package, archive and both digest files. This
command performs no installation, registry migration, Git integration or upload.

Qualify those exact package bytes, including their installer, before delivering:

```sh
test/cb01a_acceptance.sh /absolute/new-release-directory/ply-v1.1.0-cb01a.1-darwin-arm64 \
  sha256:<manifest-sha256>
```
