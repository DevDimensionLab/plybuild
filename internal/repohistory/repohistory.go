// Package repohistory holds the rule that turns recorded repository bindings
// into one historical Git common directory per repository ID.
package repohistory

import "fmt"

// Binding is one recorded repository identity. An empty GitCommonDir records
// no identity.
type Binding struct{ RepoID, GitCommonDir string }

// ConflictError reports a repository ID bound to two Git common directories.
type ConflictError struct{ RepoID string }

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflicting historical repository identity for %s", e.RepoID)
}

// CommonDirs maps each repository ID to its historical Git common directory.
// Bindings without a common directory are skipped before their ID is used.
func CommonDirs(bs []Binding) (map[string]string, error) {
	identities := map[string]string{}
	for _, b := range bs {
		if b.GitCommonDir == "" {
			continue
		}
		if old := identities[b.RepoID]; old != "" && old != b.GitCommonDir {
			return nil, &ConflictError{RepoID: b.RepoID}
		}
		identities[b.RepoID] = b.GitCommonDir
	}
	return identities, nil
}
