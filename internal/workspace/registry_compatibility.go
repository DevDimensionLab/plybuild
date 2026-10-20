package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/clirelease"
	"gopkg.in/yaml.v3"
)

const CurrentRegistryFormat = 5

// RegistryFormatBoundary records the immutable, explicitly authorized transition
// and the already available reader. It does not change any historical authority.
type RegistryFormatBoundary struct {
	FromVersion    int    `yaml:"from_version" json:"from_version"`
	FromSHA256     string `yaml:"from_sha256" json:"from_sha256"`
	ManifestSHA256 string `yaml:"manifest_sha256" json:"manifest_sha256"`
	ReaderVersion  string `yaml:"reader_version" json:"reader_version"`
	ReaderCommit   string `yaml:"reader_commit" json:"reader_commit"`
	ReaderSHA256   string `yaml:"reader_sha256" json:"reader_sha256"`
}

type registryUpgradePermit struct{ Package, ManifestSHA256 string }

type RegistryDiagnostic struct {
	Kind             string `json:"kind"`
	SchemaVersion    int    `json:"schema_version"`
	Code             string `json:"code"`
	Detail           string `json:"detail"`
	Artifact         string `json:"artifact"`
	FormatVersion    int    `json:"format_version"`
	SupportedFormats []int  `json:"supported_formats"`
	NextAction       string `json:"next_action"`
}

func (e *RegistryDiagnostic) Error() string { return e.Code + ": " + e.Detail + "; " + e.NextAction }

func registryDiagnostic(code, detail string, version int) *RegistryDiagnostic {
	next := "Run ply version --format json and retain this diagnostic with the unchanged registry and backups for producer repair before retrying. Do not strip fields or edit authorization history."
	if code == "registry_unsupported_format" {
		next = "Install a release supporting this format with its bundled ply-install and pinned manifest digest (docs/cli-installation.md); do not downgrade the registry."
	}
	if code == "registry_upgrade_required" || code == "registry_upgrade_conflict" {
		next = "Run ply workspace registry upgrade --package <release-directory> --sha256 <manifest-digest> --check --format json, inspect the fresh plan, then apply its exact confirmation."
	}
	return &RegistryDiagnostic{Kind: "WorkspaceRegistryDiagnostic@1", SchemaVersion: 1, Code: code, Detail: detail, FormatVersion: version, SupportedFormats: []int{1, 2, 3, 4, 5}, NextAction: next}
}

func RegistryError(err error) *RegistryDiagnostic {
	var d *RegistryDiagnostic
	if errors.As(err, &d) {
		return d
	}
	return nil
}

// Inspect actual syntax and the declared format before strict typed decoding.
// Unknown fields in a supported format are invalid input, never guessed versions.
func parseRegistryYAML(b []byte) (*yaml.Node, error) {
	d := yaml.NewDecoder(bytes.NewReader(b))
	var doc yaml.Node
	if err := d.Decode(&doc); err != nil {
		return nil, registryDiagnostic("registry_malformed_yaml", err.Error(), 0)
	}
	var extra yaml.Node
	if err := d.Decode(&extra); err != io.EOF {
		return nil, registryDiagnostic("registry_malformed_yaml", "expected exactly one YAML document", 0)
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, registryDiagnostic("registry_invalid_content", "registry must be a mapping", 0)
	}
	root := doc.Content[0]
	if err := validateQueueYAMLMap(root); err != nil {
		return nil, registryDiagnostic("registry_invalid_content", err.Error(), 0)
	}
	return root, nil
}

// Inspect the already validated syntax tree before typed decoding.
func registryDeclaredFormat(root *yaml.Node) (int, error) {
	for i := 0; i < len(root.Content); i += 2 {
		if root.Content[i].Value != "format_version" {
			continue
		}
		n := root.Content[i+1]
		v, err := strconv.Atoi(n.Value)
		if err != nil || n.Kind != yaml.ScalarNode || n.Tag != "!!int" || v < 1 {
			return 0, registryDiagnostic("registry_invalid_content", "format_version must be a positive integer", 0)
		}
		if v > CurrentRegistryFormat {
			return v, registryDiagnostic("registry_unsupported_format", fmt.Sprintf("declared format %d exceeds this reader", v), v)
		}
		return v, nil
	}
	return 0, registryDiagnostic("registry_invalid_content", "format_version is required", 0)
}

func classifyRegistryDecode(err error, version int) error {
	if err == nil || RegistryError(err) != nil {
		return err
	}
	code := "registry_invalid_content"
	var fields *yaml.TypeError
	if errors.As(err, &fields) {
		for _, message := range fields.Errors {
			if strings.Contains(message, "field ") && strings.Contains(message, " not found in type ") {
				code = "registry_unknown_field"
				break
			}
		}
	}
	return registryDiagnostic(code, err.Error(), version)
}

// Node.Decode does not expose KnownFields. Collect the same unknown-field
// diagnostics against the typed schema, retaining the decoder's type errors.
// Maps intentionally remain open (queue requests have their own semantic rules).
func decodeRegistryNode(root *yaml.Node, registry *WorkItemRegistry) error {
	var messages []string
	var walk func(*yaml.Node, reflect.Type, bool)
	walk = func(n *yaml.Node, typ reflect.Type, checkTypes bool) {
		for typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		if n.Tag == "!!null" {
			return
		}
		switch typ.Kind() {
		case reflect.Struct:
			if n.Kind != yaml.MappingNode {
				break
			}
			fields := yamlStructShapeFields(typ)
			for i := 0; i < len(n.Content); i += 2 {
				key, value := n.Content[i], n.Content[i+1]
				found := false
				for _, field := range fields {
					if field.name == key.Value {
						walk(value, field.typeOf, checkTypes)
						found = true
						break
					}
				}
				if !found {
					messages = append(messages, fmt.Sprintf("line %d: field %s not found in type %s", key.Line, key.Value, typ))
				}
			}
			return
		case reflect.Slice, reflect.Array:
			if n.Kind == yaml.SequenceNode {
				for _, child := range n.Content {
					walk(child, typ.Elem(), checkTypes)
				}
				return
			}
		}
		if checkTypes {
			if err := n.Decode(reflect.New(typ).Interface()); err != nil {
				var typed *yaml.TypeError
				if errors.As(err, &typed) {
					messages = append(messages, typed.Errors...)
				}
			}
		}
	}
	walk(root, reflect.TypeOf(*registry), false)
	if err := root.Decode(registry); err != nil {
		var typed *yaml.TypeError
		if !errors.As(err, &typed) {
			return err
		}
		// Revisit only invalid input to interleave unknown-field and type errors
		// in traversal order, including flow mappings on the same source line.
		messages = nil
		walk(root, reflect.TypeOf(*registry), true)
	}
	if len(messages) == 0 {
		return nil
	}
	return &yaml.TypeError{Errors: messages}
}

func validateRegistryBoundary(r WorkItemRegistry) error {
	if r.FormatVersion < 5 {
		if r.FormatBoundary != nil {
			return fmt.Errorf("legacy format cannot contain a format boundary")
		}
		return nil
	}
	b := r.FormatBoundary
	if b == nil || b.FromVersion != 4 || !digestPattern.MatchString(b.FromSHA256) || !digestPattern.MatchString(b.ManifestSHA256) || !digestPattern.MatchString(b.ReaderSHA256) || !validOIDText(b.ReaderCommit) || !strings.HasPrefix(b.ReaderVersion, "v") {
		return fmt.Errorf("format 5 requires an identified format-4 transition and delivered reader")
	}
	return nil
}

func boundaryFor(old WorkItemRegistry, m clirelease.Manifest, pin string) *RegistryFormatBoundary {
	return &RegistryFormatBoundary{FromVersion: old.FormatVersion, FromSHA256: old.RawSHA256, ManifestSHA256: pin, ReaderVersion: m.Version, ReaderCommit: m.SourceCommit, ReaderSHA256: m.SHA256}
}

// Called under the work-item lock at publication, not just during preview.
func validateRegistryFormatTransition(old, next WorkItemRegistry) error {
	if next.FormatVersion < old.FormatVersion {
		return registryDiagnostic("registry_upgrade_conflict", "registry downgrade is forbidden", old.FormatVersion)
	}
	if old.FormatVersion == 5 {
		if !reflect.DeepEqual(old.FormatBoundary, next.FormatBoundary) {
			return registryDiagnostic("registry_upgrade_conflict", "immutable format boundary changed", old.FormatVersion)
		}
		return nil
	}
	if next.FormatVersion < 5 {
		return nil
	} // Existing 1->2/3/4 authorization remains in its native command.
	p := next.upgradePermit
	if p == nil {
		return registryDiagnostic("registry_upgrade_required", "format 5 requires a confirmed native upgrade and available reader", old.FormatVersion)
	}
	if old.FormatVersion != 4 || next.FormatBoundary == nil || next.FormatBoundary.FromSHA256 != old.RawSHA256 {
		return registryDiagnostic("registry_upgrade_conflict", "registry changed since confirmation", old.FormatVersion)
	}
	m, _, err := clirelease.Load(p.Package, p.ManifestSHA256)
	if err != nil {
		return registryDiagnostic("registry_upgrade_required", "compatible reader unavailable: "+err.Error(), old.FormatVersion)
	}
	if !reflect.DeepEqual(next.FormatBoundary, boundaryFor(old, m, p.ManifestSHA256)) {
		return registryDiagnostic("registry_upgrade_conflict", "reader or upgrade basis changed", old.FormatVersion)
	}
	before, err := encodeWorkItemRegistry(old)
	if err != nil {
		return err
	}
	copy := next
	copy.FormatVersion, copy.FormatBoundary, copy.upgradePermit = old.FormatVersion, nil, nil
	after, err := encodeWorkItemRegistry(copy)
	if err != nil || !bytes.Equal(before, after) {
		return registryDiagnostic("registry_upgrade_conflict", "format transition must preserve all existing records", old.FormatVersion)
	}
	return nil
}

type RegistryUpgradePlan struct {
	Kind           string                  `json:"kind"`
	SchemaVersion  int                     `json:"schema_version"`
	Workspace      string                  `json:"workspace"`
	FromVersion    int                     `json:"from_version"`
	ToVersion      int                     `json:"to_version"`
	RegistrySHA256 string                  `json:"registry_sha256"`
	Package        string                  `json:"package"`
	ManifestSHA256 string                  `json:"manifest_sha256"`
	Boundary       *RegistryFormatBoundary `json:"boundary"`
	Backup         string                  `json:"backup"`
	Confirmation   string                  `json:"confirmation"`
	State          string                  `json:"state"`
}

func registryUpgradePlan(root string, r WorkItemRegistry, directory, pin string, m clirelease.Manifest) RegistryUpgradePlan {
	b := boundaryFor(r, m, pin)
	if r.FormatVersion == 5 {
		b = r.FormatBoundary
	}
	p := RegistryUpgradePlan{Kind: "WorkspaceRegistryUpgrade@1", SchemaVersion: 1, Workspace: root, FromVersion: 4, ToVersion: 5, RegistrySHA256: b.FromSHA256, Package: directory, ManifestSHA256: pin, Boundary: b, Backup: taskContentPath(root, "backups", b.FromSHA256), State: "ready"}
	p.Confirmation = queueDigest(p)
	return p
}

func UpgradeRegistry(d Dependencies, directory, pin, confirmation string, apply bool) (RegistryUpgradePlan, error) {
	var out RegistryUpgradePlan
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	physical, err := filepath.EvalSymlinks(directory)
	if err != nil || !filepath.IsAbs(directory) || physical != directory {
		return out, WorkInvalidArguments("release directory must be absolute and physical")
	}
	m, _, err := clirelease.Load(directory, pin)
	if err != nil {
		return out, registryDiagnostic("registry_upgrade_required", err.Error(), 0)
	}
	inspect := func(r WorkItemRegistry) error {
		if r.FormatVersion != 4 && r.FormatVersion != 5 {
			return registryDiagnostic("registry_upgrade_required", "first use the existing native Task/queue upgrade from formats 1-3 to format 4", r.FormatVersion)
		}
		out = registryUpgradePlan(root, r, directory, pin, m)
		if r.FormatVersion == 5 {
			if !reflect.DeepEqual(r.FormatBoundary, boundaryFor(WorkItemRegistry{FormatVersion: 4, RawSHA256: r.FormatBoundary.FromSHA256}, m, pin)) {
				return registryDiagnostic("registry_upgrade_conflict", "registry already has a different format boundary", 5)
			}
			out.State = "already_upgraded"
		}
		if apply && confirmation != out.Confirmation {
			return registryDiagnostic("registry_upgrade_conflict", "exact fresh preview confirmation is required", r.FormatVersion)
		}
		return nil
	}
	if !apply {
		r, err := d.WorkItems.Snapshot(root)
		if err != nil {
			return out, err
		}
		return outAfterInspect(&out, inspect(r))
	}
	err = d.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		if err = inspect(r); err != nil {
			return err
		}
		if r.FormatVersion == 5 {
			return platformSyncDirectory(filepath.Dir(workItemsPath(root)))
		}
		b, err := os.ReadFile(workItemsPath(root))
		if err != nil {
			return err
		}
		if digestTaskBytes(b) != r.RawSHA256 {
			return registryDiagnostic("registry_upgrade_conflict", "registry changed during upgrade", r.FormatVersion)
		}
		if _, err = d.TaskContent.Publish(root, "backups", b); err != nil {
			return err
		}
		r.FormatVersion, r.FormatBoundary = 5, out.Boundary
		r.upgradePermit = &registryUpgradePermit{Package: directory, ManifestSHA256: pin}
		if err = s.Publish(r); err != nil {
			return err
		}
		out.State = "upgraded"
		return nil
	})
	return out, err
}

func outAfterInspect(out *RegistryUpgradePlan, err error) (RegistryUpgradePlan, error) {
	return *out, err
}
