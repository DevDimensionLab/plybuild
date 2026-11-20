package workspace

import (
	"errors"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// DefinitionsFile lives beside .ply, independently of disposable runtime data.
const DefinitionsFile = workspacedef.File

// ReadDefinitions reads the definitions through workspacedef and applies the
// workspace Project and repository rules.
func ReadDefinitions(root string) (*workspacedef.Definitions, string, error) {
	defs, digest, err := workspacedef.Read(root)
	var defErr *workspacedef.Error
	if errors.As(err, &defErr) {
		return nil, "", definitionError(defErr)
	}
	if err != nil || defs == nil {
		return nil, digest, err
	}
	if err = validateDefinitions(*defs); err != nil {
		return nil, "", err
	}
	return defs, digest, nil
}

type definitionRootFiles struct {
	FileSystem
	root string
}

func (f definitionRootFiles) Getwd() (string, error) { return f.root, nil }

// Config discovery can recreate only the workspace marker. Init retains its
// existing atomic publication, ancestor validation and concurrent-init guards.
func initDefinitionWorkspace(files FileSystem, root string) (bool, error) {
	defs, _, err := ReadDefinitions(root)
	if err != nil || defs == nil {
		return false, err
	}
	_, err = Init(Dependencies{Files: definitionRootFiles{files, root}})
	return err == nil, err
}
