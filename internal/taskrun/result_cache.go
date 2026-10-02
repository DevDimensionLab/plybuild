package taskrun

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// Each cache version is validated against its own exact schema and immutable
// chain prefix. Projection is in memory only; historical bytes are never upgraded.
func decodeResultCache(raw []byte) (Result, error) {
	var header Envelope
	if e := json.Unmarshal(raw, &header); e != nil {
		return Result{}, integrity(e.Error())
	}
	if header.Kind != "ply.workspace.task-run-result" || header.SchemaVersion != 1 && header.SchemaVersion != 2 && header.SchemaVersion != 3 {
		return Result{}, integrity("unsupported result cache version")
	}
	if header.SchemaVersion == 1 {
		var m map[string]json.RawMessage
		if e := json.Unmarshal(raw, &m); e != nil {
			return Result{}, e
		}
		if _, ok := m["runtime_facts"]; ok {
			return Result{}, integrity("unexpected v2 cache field")
		}
		if _, ok := m["task_execution"]; ok {
			return Result{}, integrity("unexpected v2 cache field")
		}
		// Check the original bytes for duplicate keys before rebuilding the shape.
		if _, e := canonicaljson.DecodeStrict(raw); e != nil {
			return Result{}, e
		}
		m["runtime_facts"], _ = Canonical(RuntimeFacts{})
		m["task_execution"], _ = Canonical(TaskExecution{})
		raw, _ = json.Marshal(m)
	}
	var r Result
	if e := decode(raw, 4<<20, &r); e != nil {
		return r, integrity(e.Error())
	}
	return r, nil
}
func resultCacheValue(r Result, version int) any {
	r.SchemaVersion = version
	if version == 2 || version == 3 {
		return r
	}
	raw, _ := Canonical(r)
	var m map[string]any
	json.Unmarshal(raw, &m)
	delete(m, "runtime_facts")
	delete(m, "task_execution")
	return m
}
