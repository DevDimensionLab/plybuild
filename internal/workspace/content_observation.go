package workspace

import "bytes"

// observeTaskContent owns verified bytes for one closure-validation phase. Both
// content and queue validators consume the same reader. Each registry snapshot
// and each proposed publication creates a new phase, including recovery reads;
// no cached observation survives a registry reread, lock acquisition or publish.
// Only successful physical reads are retained. This reader cannot publish, and
// callers receive owned copies so parsing or editing cannot poison later reads.
func observeTaskContent(source *TaskContentStorage) *TaskContentStorage {
	verified := map[capturedTaskContentKey][]byte{}
	return &TaskContentStorage{scope: source.scope, readCaptured: func(root, kind, digest string) ([]byte, error) {
		key := capturedTaskContentKey{root, kind, digest}
		if b, ok := verified[key]; ok {
			return bytes.Clone(b), nil
		}
		b, err := source.Read(root, kind, digest)
		if err != nil {
			return nil, err
		}
		verified[key] = bytes.Clone(b)
		return b, nil
	}}
}
