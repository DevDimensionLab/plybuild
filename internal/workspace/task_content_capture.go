package workspace

import (
	"bytes"
	"sort"
)

type capturedTaskContentKey struct {
	work         ProjectWork
	kind, digest string
}
type capturedTaskContentValue struct {
	bytes  []byte
	err    error
	owners map[string]bool
}

// taskContentCapture lets the domain validators share immutable managed reads
// within one projection. It is private to that request and cannot publish. The
// original storage still owns physical-path, size, mode, and digest validation.
type taskContentCapture struct {
	source *TaskContentStorage
	store  *TaskContentStorage
	owner  string
	values map[capturedTaskContentKey]*capturedTaskContentValue
}

func captureTaskContent(source *TaskContentStorage) *taskContentCapture {
	c := &taskContentCapture{source: source, values: map[capturedTaskContentKey]*capturedTaskContentValue{}}
	c.store = &TaskContentStorage{readCaptured: c.read, fault: func(string) error {
		return contentError("task_content_read_only", "captured content is read-only", nil)
	}}
	return c
}

func (c *taskContentCapture) read(work ProjectWork, kind, digest string) ([]byte, error) {
	key := capturedTaskContentKey{work, kind, digest}
	v := c.values[key]
	if v == nil {
		v = &capturedTaskContentValue{owners: map[string]bool{}}
		v.bytes, v.err = c.source.Read(work, kind, digest)
		c.values[key] = v
	}
	v.owners[c.owner] = true
	return v.bytes, v.err
}

// changedOwners performs one final physical validation of each successfully
// captured source. Absence or corruption is never repaired by this reader.
func (c *taskContentCapture) changedOwners() map[string][]string {
	keys := make([]capturedTaskContentKey, 0, len(c.values))
	for key := range c.values {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, b := keys[i], keys[j]
		if a.work.Dir != b.work.Dir {
			return a.work.Dir < b.work.Dir
		}
		if a.kind != b.kind {
			return a.kind < b.kind
		}
		return a.digest < b.digest
	})
	changed := map[string][]string{}
	for _, key := range keys {
		v := c.values[key]
		if v.err != nil {
			continue
		}
		now, err := c.source.Read(key.work, key.kind, key.digest)
		if err != nil || !bytes.Equal(v.bytes, now) {
			for owner := range v.owners {
				changed[owner] = append(changed[owner], key.digest)
			}
		}
	}
	return changed
}
