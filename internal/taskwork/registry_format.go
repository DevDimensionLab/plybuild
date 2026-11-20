package taskwork

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// RegistryFormat is the only work-item registry format this reader accepts.
// It is the physical envelope that keeps integration and publication evidence
// in task-content objects. There is no upgrade or migration from other formats.
const RegistryFormat = 11

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
		next = fmt.Sprintf("This release reads only work-item registry format %d and does not upgrade or migrate registries. Use a release that reads this format (docs/cli-installation.md); do not edit the registry.", RegistryFormat)
	}
	return &RegistryDiagnostic{Kind: "WorkspaceRegistryDiagnostic@1", SchemaVersion: 1, Code: code, Detail: detail, FormatVersion: version, SupportedFormats: []int{RegistryFormat}, NextAction: next}
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

// Inspect the already validated syntax tree before typed decoding. Any other
// declared format is rejected before referenced evidence objects are read.
func registryDeclaredFormat(root *yaml.Node) error {
	n := registryNodeField(root, "format_version")
	if n == nil {
		return registryDiagnostic("registry_invalid_content", "format_version is required", 0)
	}
	v, err := strconv.Atoi(n.Value)
	if err != nil || n.Kind != yaml.ScalarNode || n.Tag != "!!int" || v < 1 {
		return registryDiagnostic("registry_invalid_content", "format_version must be a positive integer", 0)
	}
	if v != RegistryFormat {
		return registryDiagnostic("registry_unsupported_format", fmt.Sprintf("declared format %d is not format %d", v, RegistryFormat), v)
	}
	return nil
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
var yamlUnmarshalerType = reflect.TypeOf((*yaml.Unmarshaler)(nil)).Elem()

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
		// Types that decode themselves own their members, including those of
		// historical delivery records; the shape validator checks current ones.
		if reflect.PointerTo(typ).Implements(yamlUnmarshalerType) {
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
				// Syntax validation guarantees !!str keys, but an explicit tag
				// can also occur on a collection. Match yaml's string decoding
				// before field lookup; a failed key skips its value entirely.
				name := key.Value
				if key.Kind != yaml.ScalarNode {
					if err := key.Decode(&name); err != nil {
						var typed *yaml.TypeError
						if checkTypes && errors.As(err, &typed) {
							messages = append(messages, typed.Errors...)
						}
						continue
					}
				}
				found := false
				for _, field := range fields {
					if field.name == name {
						walk(value, field.typeOf, checkTypes)
						found = true
						break
					}
				}
				if !found {
					messages = append(messages, fmt.Sprintf("line %d: field %s not found in type %s", key.Line, name, typ))
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
