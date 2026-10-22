package workspace

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"

	"github.com/devdimensionlab/plybuild/internal/clirelease"
)

func validateEpicRegistryUpgradeTransition(old, next WorkItemRegistry) error {
	if old.FormatVersion == 6 {
		if !reflect.DeepEqual(old.FormatBoundary, next.FormatBoundary) || !reflect.DeepEqual(old.EpicCreateBoundary, next.EpicCreateBoundary) {
			return fmt.Errorf("immutable registry boundaries changed")
		}
		return nil
	}
	if old.FormatVersion != 5 || next.upgradePermit == nil {
		return registryDiagnostic("registry_upgrade_required", "format 6 requires an explicit digest-bound native upgrade from format 5", old.FormatVersion)
	}
	p := next.upgradePermit
	m, _, e := clirelease.Load(p.Package, p.ManifestSHA256)
	if e != nil {
		return e
	}
	if !slices.Contains(m.RegistryFormats, 6) || !reflect.DeepEqual(next.EpicCreateBoundary, boundaryFor(old, m, p.ManifestSHA256)) {
		return registryDiagnostic("registry_upgrade_conflict", "format 6 reader or exact registry basis differs", old.FormatVersion)
	}
	copy := next
	copy.FormatVersion = 5
	copy.EpicCreateBoundary = nil
	copy.EpicCreateOperations = nil
	copy.upgradePermit = nil
	a, e := encodeWorkItemRegistry(old)
	if e != nil {
		return e
	}
	b, e := encodeWorkItemRegistry(copy)
	if e != nil || !bytes.Equal(a, b) {
		return fmt.Errorf("format 6 upgrade must preserve existing semantic records")
	}
	return nil
}

// UpgradeEpicCreateRegistry is the explicit 5 -> 6 release-pinned transition.
// The ordinary 4 -> 5 command and old-format flows retain their original contract.
func UpgradeEpicCreateRegistry(d Dependencies, directory, pin, confirmation string, apply bool) (RegistryUpgradePlan, error) {
	var out RegistryUpgradePlan
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return out, e
	}
	physical, e := filepath.EvalSymlinks(directory)
	if e != nil || !filepath.IsAbs(directory) || physical != directory {
		return out, WorkInvalidArguments("release directory must be absolute and physical")
	}
	m, _, e := clirelease.Load(directory, pin)
	if e != nil {
		return out, e
	}
	if !slices.Contains(m.RegistryFormats, 6) {
		return out, registryDiagnostic("registry_upgrade_required", "release does not declare a format 6 reader", 0)
	}
	inspect := func(r WorkItemRegistry) error {
		if r.FormatVersion != 5 && r.FormatVersion != 6 {
			return registryDiagnostic("registry_upgrade_required", "first complete the existing native transition to format 5", r.FormatVersion)
		}
		b := boundaryFor(r, m, pin)
		if r.FormatVersion == 6 {
			b = r.EpicCreateBoundary
			if !reflect.DeepEqual(b, boundaryFor(WorkItemRegistry{FormatVersion: 5, RawSHA256: b.FromSHA256}, m, pin)) {
				return registryDiagnostic("registry_upgrade_conflict", "format 6 already has a different boundary", 6)
			}
		}
		out = RegistryUpgradePlan{Kind: "WorkspaceRegistryUpgrade@1", SchemaVersion: 1, Workspace: root, FromVersion: 5, ToVersion: 6, RegistrySHA256: b.FromSHA256, Package: directory, ManifestSHA256: pin, Boundary: b, Backup: taskContentPath(root, "backups", b.FromSHA256), State: "ready"}
		out.Confirmation = queueDigest(out)
		if r.FormatVersion == 6 {
			out.State = "already_upgraded"
		}
		if apply && confirmation != out.Confirmation {
			return registryDiagnostic("registry_upgrade_conflict", "exact fresh preview confirmation required", r.FormatVersion)
		}
		return nil
	}
	if !apply {
		r, e := d.WorkItems.Snapshot(root)
		if e != nil {
			return out, e
		}
		return outAfterInspect(&out, inspect(r))
	}
	e = d.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
		r, e := s.Snapshot()
		if e != nil {
			return e
		}
		if e = inspect(r); e != nil {
			return e
		}
		if r.FormatVersion == 6 {
			return platformSyncDirectory(filepath.Dir(workItemsPath(root)))
		}
		b, e := os.ReadFile(workItemsPath(root))
		if e != nil {
			return e
		}
		if digestTaskBytes(b) != r.RawSHA256 {
			return fmt.Errorf("registry bytes changed during upgrade")
		}
		if _, e = d.TaskContent.Publish(root, "backups", b); e != nil {
			return e
		}
		r.FormatVersion = 6
		r.EpicCreateBoundary = out.Boundary
		r.EpicCreateOperations = []EpicCreateOperation{}
		r.upgradePermit = &registryUpgradePermit{Package: directory, ManifestSHA256: pin}
		if e = s.Publish(r); e != nil {
			return e
		}
		out.State = "upgraded"
		return nil
	})
	return out, e
}
