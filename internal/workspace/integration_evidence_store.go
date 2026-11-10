package workspace

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

// The physical envelope is deliberately newer than every inline reader. Its
// nested registry retains its own lifecycle format and validators. Only fresh
// stores use this representation; existing inline registries are never migrated.
const integrationEvidenceEnvelopeFormat = 9
const publicationEvidenceEnvelopeFormat = 10

// IntegrationEvidenceSourceKind identifies workspace-managed journal sources.
const IntegrationEvidenceSourceKind = "integration-evidence"
const integrationEvidenceKind = IntegrationEvidenceSourceKind
const integrationEvidenceLimit = 64 << 20

var integrationEvidenceCollections = []string{
	"integration_authorities", "integration_intents", "integration_attempts", "integration_results",
}

// Format 9 keeps publication records inline. Only format 10 includes them in
// the existing evidence store; physical representation never changes on write.
func evidenceCollections(publications bool) []string {
	if publications {
		return append(append([]string(nil), integrationEvidenceCollections...), "task_content_publications")
	}
	return integrationEvidenceCollections
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

func (store *systemWorkItemStore) encodeStoredRegistry(root string, registry *WorkItemRegistry, old WorkItemRegistry) ([]byte, error) {
	// Validate the complete logical state before publishing even an orphan object.
	external := old.integrationEvidence || old.RawSHA256 == ""
	raw, err := encodeWorkItemRegistry(*registry)
	if err != nil || !external {
		return raw, err
	}
	node, err := parseRegistryYAML(raw)
	if err != nil {
		return nil, err
	}
	storage := evidenceStorage(store)
	retained := retainedIntegrationEvidence(old, *registry)
	registry.integrationEvidenceRefs = nil
	registry.publicationEvidence = old.publicationEvidence || old.RawSHA256 == ""
	for _, name := range evidenceCollections(registry.publicationEvidence) {
		collection := registryNodeField(node, name)
		if collection == nil {
			continue
		} // Logical format 1 has no lifecycle evidence.
		for i, record := range collection.Content {
			identity, _ := evidenceIdentity(name)
			id := registryNodeField(record, identity)
			if id == nil {
				return nil, fmt.Errorf("%s evidence record has no %s", name, identity)
			}
			if digest := retained[name+"/"+id.Value]; digest != "" {
				// Retain the exact historical bytes and reference, even if another
				// compatible writer used a different valid YAML spelling.
				if _, err := storage.Read(root, integrationEvidenceKind, digest); err != nil {
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
			digest, err := storage.Publish(root, integrationEvidenceKind, payload)
			if err != nil {
				return nil, err
			}
			registry.integrationEvidenceRefs = append(registry.integrationEvidenceRefs, digest)
			collection.Content[i] = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: digest}
		}
	}
	format := integrationEvidenceEnvelopeFormat
	if registry.publicationEvidence {
		format = publicationEvidenceEnvelopeFormat
	}
	envelope := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Tag: "!!str", Value: "format_version"},
		{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.Itoa(format)},
		{Kind: yaml.ScalarNode, Tag: "!!str", Value: "registry"}, node,
	}}
	return marshalRegistryNode(envelope)
}

func (store *systemWorkItemStore) decodeStoredRegistry(root string, raw []byte) (WorkItemRegistry, error) {
	node, err := parseRegistryYAML(raw)
	if err != nil {
		return WorkItemRegistry{}, err
	}
	version := registryNodeField(node, "format_version")
	if version == nil || version.Kind != yaml.ScalarNode || version.Tag != "!!int" || (version.Value != strconv.Itoa(integrationEvidenceEnvelopeFormat) && version.Value != strconv.Itoa(publicationEvidenceEnvelopeFormat)) {
		return decodeWorkItemRegistry(raw)
	}
	logical := registryNodeField(node, "registry")
	if len(node.Content) != 4 || logical == nil || logical.Kind != yaml.MappingNode {
		return WorkItemRegistry{}, fmt.Errorf("invalid integration evidence envelope")
	}
	// Reject unsupported lifecycle formats before touching any referenced objects.
	if _, err := registryDeclaredFormat(logical); err != nil {
		return WorkItemRegistry{}, err
	}
	storage := evidenceStorage(store)
	var refs []string
	publications := version.Value == strconv.Itoa(publicationEvidenceEnvelopeFormat)
	for _, name := range evidenceCollections(publications) {
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
			payload, err := storage.Read(root, integrationEvidenceKind, ref.Value)
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
	registry.integrationEvidence = err == nil
	registry.publicationEvidence = publications && err == nil
	registry.integrationEvidenceRefs = refs
	return registry, err
}

// IntegrationEvidenceSources binds read-view and journal freshness to the actual
// immutable bytes consumed by the store, not a reserialization of their values.
func IntegrationEvidenceSources(root string, r WorkItemRegistry) []TaskJournalSource {
	out := []TaskJournalSource{}
	if !r.integrationEvidence {
		return out
	}
	for _, digest := range r.integrationEvidenceRefs {
		out = append(out, TaskJournalSource{integrationEvidenceKind, taskContentPath(root, integrationEvidenceKind, digest), digest})
	}
	return out
}

// Confirmed in-memory values may be reused only while all their external bytes
// still match. Immutable records are read and hashed on every locked snapshot.
func (store *systemWorkItemStore) verifyIntegrationEvidence(root string, r WorkItemRegistry) error {
	if !r.integrationEvidence {
		return nil
	}
	storage := evidenceStorage(store)
	for _, digest := range r.integrationEvidenceRefs {
		if _, err := storage.Read(root, integrationEvidenceKind, digest); err != nil {
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
	if !old.integrationEvidence {
		return out
	}
	collections := func(r WorkItemRegistry) []reflect.Value {
		values := []reflect.Value{reflect.ValueOf(r.IntegrationAuthorities), reflect.ValueOf(r.IntegrationIntents), reflect.ValueOf(r.IntegrationAttempts), reflect.ValueOf(r.IntegrationResults)}
		if old.publicationEvidence {
			values = append(values, reflect.ValueOf(r.TaskContentPublications))
		}
		return values
	}
	before, after := collections(old), collections(next)
	offset := 0
	for c, name := range evidenceCollections(old.publicationEvidence) {
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
// authoritative store. A byte hash alone cannot establish managed validity.
func IntegrationEvidenceSourceDigest(root, locator string) (string, error) {
	digest := "sha256:" + filepath.Base(locator)
	if !digestPattern.MatchString(digest) || locator != taskContentPath(root, integrationEvidenceKind, digest) {
		return "", contentError("task_content_integrity_conflict", "integration evidence source is outside its managed binding", nil)
	}
	if _, err := evidenceStorage(newSystemWorkItemStore()).Read(root, integrationEvidenceKind, digest); err != nil {
		return "", err
	}
	return digest, nil
}
