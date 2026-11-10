package workspace

import "fmt"

// This nested contract marker is emitted only for newly completed native creates.
// Older strict readers reject the unknown field before accepting the registry.
const nativeReceiptVersion = 1

func validateNativeReceipt(r WorkItemRegistry, t TaskRecord, hasOperation bool) error {
	w := t.Worktree
	if w.NativeReceiptVersion != nativeReceiptVersion || hasOperation || (t.WorktreeState != WorkItemReady && t.WorktreeState != WorkItemRetired) || w.Origin != "ply_created" || w.OwnerKind != "task" || w.OwnerID != t.ID || w.GitCommonDir != t.GitCommonDir || w.ParentEpicID != t.ParentEpicID || !worktreeIDPattern.MatchString(string(w.ID)) || !operationIDPattern.MatchString(string(w.OperationID)) || !validStoredPath(w.Locator) || !validFullBranchRef(w.Ref) || w.OID != w.ParentOID || w.Tree != w.ParentTree || !validOIDText(w.OID) || !validOIDText(w.Tree) || !historicalEpicBaseMatches(r, t.ParentEpicID, t.RepoID, w.ParentWorktreeID, w.ParentRef, w.ParentOID, w.ParentTree) {
		return fmt.Errorf("Task %s has invalid native receipt", t.ID)
	}
	// Only the original intent's identity fields participate in this hash. No
	// operation/observation is reconstituted from the completion receipt.
	digest, err := nativeIntentDigest(t.GitCommonDir, t.ParentEpicID, w.ParentOID, w.ParentRef, w.ParentTree, w.ParentWorktreeID, t.ProjectID, t.RepoID, w.Ref, w.Locator, t.ID)
	if err != nil || digest != w.IntentDigest {
		return fmt.Errorf("Task %s native receipt digest differs", t.ID)
	}
	return nil
}

// Completed bindings keep reserving their ref and path, including retired Tasks.
func nativeReceiptReservation(r WorkItemRegistry, common, ref, path string) error {
	for _, t := range r.Tasks {
		if w := t.Worktree; w != nil && w.Origin == "ply_created" && (pathsOverlap(path, w.Locator) || common == w.GitCommonDir && ref == w.Ref) {
			return workError(ErrorWorkIdentityConflict, "Git resource reserved by Task "+string(t.ID), nil)
		}
	}
	return nil
}
