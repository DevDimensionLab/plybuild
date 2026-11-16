package workspace

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

// IntegrationEvidenceSourceKind identifies workspace-managed journal sources.
const IntegrationEvidenceSourceKind = "integration-evidence"
const integrationEvidenceKind = IntegrationEvidenceSourceKind
const integrationEvidenceLimit = 64 << 20

// The physical registry keeps these collections as references to immutable
// evidence objects, identified per record.
var integrationEvidenceCollections = []string{
	"integration_authorities", "integration_intents", "integration_attempts", "integration_results", "task_content_publications",
}

func evidenceIdentity(name string) (yamlField, goField string) {
	if name == "task_content_publications" {
		return "publication_key", "PublicationKey"
	}
	return "id", "ID"
}

func registryNodeField(root *yaml.Node, key string) *yaml.Node {
	if root.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == key {
			return root.Content[i+1]
		}
	}
	return nil
}

func marshalRegistryNode(root *yaml.Node) ([]byte, error) {
	var out bytes.Buffer
	e := yaml.NewEncoder(&out)
	e.SetIndent(2)
	if err := e.Encode(root); err != nil {
		return nil, err
	}
	if err := e.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func evidenceStorage(store *systemWorkItemStore) *TaskContentStorage {
	return &TaskContentStorage{fault: func(stage string) error { return store.fail("integration-evidence-" + stage) }}
}

func (store *systemWorkItemStore) encodeStoredRegistry(work ProjectWork, registry *WorkItemRegistry, old WorkItemRegistry) ([]byte, error) {
	// Validate the complete logical state before publishing even an orphan object.
	raw, err := encodeWorkItemRegistry(*registry)
	if err != nil {
		return nil, err
	}
	node, err := parseRegistryYAML(raw)
	if err != nil {
		return nil, err
	}
	storage := evidenceStorage(store)
	retained := retainedIntegrationEvidence(old, *registry)
	registry.integrationEvidenceRefs = nil
	for _, name := range integrationEvidenceCollections {
		collection := registryNodeField(node, name)
		for i, record := range collection.Content {
			identity, _ := evidenceIdentity(name)
			id := registryNodeField(record, identity)
			if id == nil {
				return nil, fmt.Errorf("%s evidence record has no %s", name, identity)
			}
			if digest := retained[name+"/"+id.Value]; digest != "" {
				// Retain the exact historical bytes and reference, even if another
				// compatible writer used a different valid YAML spelling.
				if _, err := storage.Read(work, integrationEvidenceKind, digest); err != nil {
					return nil, err
				}
				registry.integrationEvidenceRefs = append(registry.integrationEvidenceRefs, digest)
				collection.Content[i] = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: digest}
				continue
			}
			payload, err := marshalRegistryNode(record)
			if err != nil {
				return nil, err
			}
			if len(payload) > integrationEvidenceLimit {
				return nil, fmt.Errorf("integration evidence exceeds %d bytes", integrationEvidenceLimit)
			}
			digest, err := storage.Publish(work, integrationEvidenceKind, payload)
			if err != nil {
				return nil, err
			}
			registry.integrationEvidenceRefs = append(registry.integrationEvidenceRefs, digest)
			collection.Content[i] = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: digest}
		}
	}
	envelope := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: "!!str", Value: "format_version"},
		{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.Itoa(RegistryFormat)},
		{Kind: yaml.ScalarNode, Tag: "!!str", Value: "registry"}, node,
	}}
	return marshalRegistryNode(envelope)
}

func (store *systemWorkItemStore) decodeStoredRegistry(work ProjectWork, raw []byte) (WorkItemRegistry, error) {
	node, err := parseRegistryYAML(raw)
	if err != nil {
		return WorkItemRegistry{}, err
	}
	// Reject any other format before touching referenced objects.
	if err := registryDeclaredFormat(node); err != nil {
		return WorkItemRegistry{}, err
	}
	logical := registryNodeField(node, "registry")
	if len(node.Content) != 4 || node.Content[0].Value != "format_version" || logical == nil || logical.Kind != yaml.MappingNode {
		return WorkItemRegistry{}, registryDiagnostic("registry_invalid_content", "registry must contain exactly format_version and a registry mapping", RegistryFormat)
	}
	storage := evidenceStorage(store)
	var refs []string
	for _, name := range integrationEvidenceCollections {
		collection := registryNodeField(logical, name)
		if collection == nil {
			continue
		} // Normal shape validation rejects missing lists.
		if collection.Kind != yaml.SequenceNode || collection.Tag != "!!seq" {
			return WorkItemRegistry{}, fmt.Errorf("%s must be an evidence reference array", name)
		}
		for i, ref := range collection.Content {
			if ref.Kind != yaml.ScalarNode || ref.Tag != "!!str" || !digestPattern.MatchString(ref.Value) {
				return WorkItemRegistry{}, fmt.Errorf("invalid %s evidence reference", name)
			}
			refs = append(refs, ref.Value)
			payload, err := storage.Read(work, integrationEvidenceKind, ref.Value)
			if err != nil {
				return WorkItemRegistry{}, err
			}
			record, err := parseRegistryYAML(payload)
			if err != nil {
				return WorkItemRegistry{}, err
			}
			collection.Content[i] = record
		}
	}
	hydrated, err := marshalRegistryNode(logical)
	if err != nil {
		return WorkItemRegistry{}, err
	}
	registry, err := decodeWorkItemRegistry(hydrated)
	registry.integrationEvidenceRefs = refs
	return registry, err
}

// IntegrationEvidenceSources binds read-view and journal freshness to the actual
// immutable bytes consumed by the store, not a reserialization of their values.
func IntegrationEvidenceSources(work ProjectWork, r WorkItemRegistry) []TaskJournalSource {
	out := []TaskJournalSource{}
	for _, digest := range r.integrationEvidenceRefs {
		out = append(out, TaskJournalSource{integrationEvidenceKind, taskContentPath(work, integrationEvidenceKind, digest), digest})
	}
	return out
}

// Confirmed in-memory values may be reused only while all their external bytes
// still match. Immutable records are read and hashed on every locked snapshot.
func (store *systemWorkItemStore) verifyIntegrationEvidence(work ProjectWork, r WorkItemRegistry) error {
	storage := evidenceStorage(store)
	for _, digest := range r.integrationEvidenceRefs {
		if _, err := storage.Read(work, integrationEvidenceKind, digest); err != nil {
			return err
		}
	}
	return nil
}

// Arrays are ordered and validated by the logical registry. Match by record
// identity, since an insertion can change an existing record's array position.
// A changed attempt gets a new object; an unchanged record retains its bytes.
func retainedIntegrationEvidence(old, next WorkItemRegistry) map[string]string {
	out := map[string]string{}
	if old.RawSHA256 == "" {
		return out
	} // A registry that was never stored has no evidence objects to retain.
	collections := func(r WorkItemRegistry) []reflect.Value {
		return []reflect.Value{reflect.ValueOf(r.IntegrationAuthorities), reflect.ValueOf(r.IntegrationIntents), reflect.ValueOf(r.IntegrationAttempts), reflect.ValueOf(r.IntegrationResults), reflect.ValueOf(r.TaskContentPublications)}
	}
	before, after := collections(old), collections(next)
	offset := 0
	for c, name := range integrationEvidenceCollections {
		_, field := evidenceIdentity(name)
		type prior struct {
			value  any
			digest string
		}
		byID := map[string]prior{}
		for i := 0; i < before[c].Len(); i++ {
			record := before[c].Index(i)
			byID[record.FieldByName(field).String()] = prior{record.Interface(), old.integrationEvidenceRefs[offset]}
			offset++
		}
		for i := 0; i < after[c].Len(); i++ {
			record := after[c].Index(i)
			id := record.FieldByName(field).String()
			previous, found := byID[id]
			if found && reflect.DeepEqual(previous.value, record.Interface()) {
				out[name+"/"+id] = previous.digest
			}
		}
	}
	return out
}

// IntegrationEvidenceSourceDigest revalidates a captured source through the same
// physical-path, directory-mode, immutable-mode, bound and digest checks as the
// authoritative store. A byte hash alone cannot establish managed validity. The
// source must be managed by the work of a Project of the workspace at root.
func IntegrationEvidenceSourceDigest(root, locator string) (string, error) {
	digest := "sha256:" + filepath.Base(locator)
	if !digestPattern.MatchString(digest) {
		return "", contentError("task_content_integrity_conflict", "integration evidence source is outside its managed binding", nil)
	}
	works, err := DefinedWork(root)
	if err != nil {
		return "", err
	}
	for _, work := range works {
		if locator != taskContentPath(work, integrationEvidenceKind, digest) {
			continue
		}
		if _, err := evidenceStorage(newSystemWorkItemStore()).Read(work, integrationEvidenceKind, digest); err != nil {
			return "", err
		}
		return digest, nil
	}
	return "", contentError("task_content_integrity_conflict", "integration evidence source is outside its managed binding", nil)
}
