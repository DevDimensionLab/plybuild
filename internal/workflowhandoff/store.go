package workflowhandoff

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const storeRelativeRoot = ".ply/workflow/handoffs/v1"

type storeOps struct {
	lstat         func(string) (fs.FileInfo, error)
	readFile      func(string) ([]byte, error)
	readDir       func(string) ([]fs.DirEntry, error)
	mkdirAll      func(string, fs.FileMode) error
	mkdirTemp     func(string, string) (string, error)
	createTemp    func(string, string) (*os.File, error)
	openFile      func(string, int, fs.FileMode) (*os.File, error)
	write         func(*os.File, []byte) (int, error)
	chmodFile     func(*os.File, fs.FileMode) error
	chmodPath     func(string, fs.FileMode) error
	syncFile      func(*os.File) error
	closeFile     func(*os.File) error
	statFile      func(*os.File) (fs.FileInfo, error)
	rename        func(string, string) error
	remove        func(string) error
	removeAll     func(string) error
	syncDirectory func(string) error
	lock          func(*os.File) error
	unlock        func(*os.File) error
	replace       func(string, string) error
}

type systemStore struct {
	files FileSystem
	clock Clock
	ops   storeOps
}

func newSystemStore(files FileSystem, clock Clock) Store {
	return &systemStore{files: files, clock: clock, ops: storeOps{
		lstat: os.Lstat, readFile: os.ReadFile, readDir: os.ReadDir, mkdirAll: os.MkdirAll,
		mkdirTemp: os.MkdirTemp, createTemp: os.CreateTemp, openFile: os.OpenFile,
		write:     func(file *os.File, contents []byte) (int, error) { return file.Write(contents) },
		chmodFile: func(file *os.File, mode fs.FileMode) error { return file.Chmod(mode) }, chmodPath: os.Chmod,
		syncFile: func(file *os.File) error { return file.Sync() }, closeFile: func(file *os.File) error { return file.Close() },
		statFile: func(file *os.File) (fs.FileInfo, error) { return file.Stat() },
		rename:   os.Rename, remove: os.Remove, removeAll: os.RemoveAll,
		syncDirectory: storeSyncDirectory, lock: storePlatformLock, unlock: storePlatformUnlock,
		replace: storePlatformReplace,
	}}
}

func (store *systemStore) Create(input CreateStoreInput) (CreateStoreResult, error) {
	if input.Handoff.Identity.PublicationKey != input.PublicationKey || input.Handoff.SourceDraftSHA256 != input.SourceDraftSHA256 {
		return CreateStoreResult{}, classified(ErrorSchemaInvalid, "store create input does not match handoff publication identity", nil)
	}
	root := storeRoot(input.WorkspaceRoot)
	if err := store.ensureDirectory(root); err != nil {
		return CreateStoreResult{}, err
	}
	for _, relative := range []string{"publications", "activities", "runs"} {
		if err := store.ensureDirectory(filepath.Join(root, relative)); err != nil {
			return CreateStoreResult{}, err
		}
	}
	var lock *os.File
	if !input.StoreAlreadyLocked {
		var err error
		lock, err = store.acquire(filepath.Join(root, "store.lock"))
		if err != nil {
			return CreateStoreResult{}, err
		}
	}
	var result CreateStoreResult
	operationErr := func() error {
		publicationPath := filepath.Join(root, "publications", publicationFilename(input.PublicationKey))
		if _, err := store.ops.lstat(publicationPath); err == nil {
			fields, err := store.readRef(publicationPath, "ply.workflow.publication-ref", "publication_key_sha256", "source_draft_sha256", "activity_id", "run_id", "handoff_id", "handoff_sha256", "locator")
			if err != nil {
				return err
			}
			draft, _ := stringField(fields, "source_draft_sha256", "publication ref")
			if draft != input.SourceDraftSHA256 {
				return classified(ErrorConflict, fmt.Sprintf("publication key %s already has different content", input.PublicationKey), nil)
			}
			if objectMapString(fields, "publication_key_sha256") != digestBytes([]byte(input.PublicationKey)) {
				return storeConflict(publicationPath, "publication key digest does not match", nil)
			}
			locator, _ := stringField(fields, "locator", "publication ref")
			resolved, err := resolveStoreLocator(root, locator)
			if err != nil {
				return storeConflict(publicationPath, "publication locator is invalid", err)
			}
			snapshot, err := store.ReadByLocator(resolved)
			if err != nil {
				return err
			}
			if objectMapString(fields, "activity_id") != snapshot.Handoff.Identity.ActivityID || objectMapString(fields, "run_id") != snapshot.Handoff.Identity.RunID || objectMapString(fields, "handoff_id") != snapshot.Handoff.Identity.HandoffID || objectMapString(fields, "handoff_sha256") != snapshot.Handoff.SHA256 {
				return storeConflict(publicationPath, "publication ref does not match the handoff", nil)
			}
			result = CreateStoreResult{Snapshot: snapshot, Created: false}
			return nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return storeConflict(publicationPath, "cannot inspect publication ref", err)
		}
		if input.Revalidate != nil {
			if err := input.Revalidate(); err != nil {
				return err
			}
		}
		orphan, found, err := store.findRunByPublication(root, input.PublicationKey, input.SourceDraftSHA256)
		if err != nil {
			return err
		}
		if found {
			if !input.DeferActivityPublication {
				activityDir := filepath.Join(root, "activities", orphan.Identity.ActivityID)
				if err := store.ensureDirectory(filepath.Join(activityDir, "events")); err != nil {
					return err
				}
				if err := store.ensureLockFile(filepath.Join(activityDir, "activity.lock")); err != nil {
					return err
				}
				events, err := store.ops.readDir(filepath.Join(activityDir, "events"))
				if err != nil {
					return storeIO("read activity events", activityDir, err)
				}
				if len(events) == 0 {
					document := canonicaljson.Object{{Name: "document_kind", Value: "ply.workflow.handoff"}, {Name: "document_id", Value: orphan.Identity.HandoffID}, {Name: "sha256", Value: orphan.SHA256}, {Name: "locator", Value: relativeLocator(root, orphan.Locator)}}
					if _, err := store.publishActivityEvent(root, orphan, "published", document, nil, nil, nil); err != nil {
						return err
					}
				}
			}
			publication := envelope("ply.workflow.publication-ref", canonicaljson.Member{Name: "publication_key_sha256", Value: digestBytes([]byte(input.PublicationKey))}, canonicaljson.Member{Name: "source_draft_sha256", Value: input.SourceDraftSHA256}, canonicaljson.Member{Name: "activity_id", Value: orphan.Identity.ActivityID}, canonicaljson.Member{Name: "run_id", Value: orphan.Identity.RunID}, canonicaljson.Member{Name: "handoff_id", Value: orphan.Identity.HandoffID}, canonicaljson.Member{Name: "handoff_sha256", Value: orphan.SHA256}, canonicaljson.Member{Name: "locator", Value: relativeLocator(root, orphan.Locator)})
			publicationBytes, _ := canonicaljson.Marshal(publication)
			if _, err := store.atomicWriteOnce(publicationPath, publicationBytes, 0o400); err != nil {
				return err
			}
			if !input.DeferActivityPublication {
				if err := store.writeHead(root, orphan, "ready"); err != nil {
					return err
				}
			}
			result = CreateStoreResult{Snapshot: Snapshot{WorkspaceRoot: input.WorkspaceRoot, Handoff: orphan, HeadState: "ready"}, Created: false}
			return nil
		}

		runID := input.Handoff.Identity.RunID
		finalRun := filepath.Join(root, "runs", runID)
		if _, err := store.ops.lstat(finalRun); err == nil {
			return classified(ErrorConflict, fmt.Sprintf("run %s already exists without a matching publication ref", runID), nil)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return storeIO("inspect run", finalRun, err)
		}
		temporary, err := store.ops.mkdirTemp(filepath.Join(root, "runs"), ".run-*")
		if err != nil {
			return storeIO("mkdir run temp", filepath.Join(root, "runs"), err)
		}
		published := false
		defer func() {
			if !published {
				_ = store.ops.removeAll(temporary)
			}
		}()
		if err := store.ops.chmodPath(temporary, 0o700); err != nil {
			return storeIO("chmod run temp", temporary, err)
		}
		reply := filepath.Join(temporary, "reply")
		for _, directory := range []string{
			reply, filepath.Join(reply, "staging"), filepath.Join(reply, "start", "objects"), filepath.Join(reply, "start", "rejected"), filepath.Join(reply, "start", "events"), filepath.Join(reply, "terminal", "objects"), filepath.Join(reply, "terminal", "rejected"), filepath.Join(reply, "terminal", "events"), filepath.Join(reply, "artifacts", "sha256"),
		} {
			if err := store.ensureDirectory(directory); err != nil {
				return err
			}
		}
		if err := store.ensureLockFile(filepath.Join(reply, "reply.lock")); err != nil {
			return err
		}
		handoffName := "handoff." + strings.TrimPrefix(input.Handoff.SHA256, "sha256:") + ".json"
		if _, err := store.atomicWriteOnce(filepath.Join(temporary, handoffName), input.Handoff.Bytes, 0o400); err != nil {
			return err
		}
		if err := store.ops.syncDirectory(temporary); err != nil {
			return storeIO("sync run temp", temporary, err)
		}
		if err := store.ops.rename(temporary, finalRun); err != nil {
			return storeIO("publish run", finalRun, err)
		}
		published = true
		if err := store.ops.syncDirectory(filepath.Join(root, "runs")); err != nil {
			return storeIO("sync runs", filepath.Join(root, "runs"), err)
		}

		handoffLocator := filepath.Join(finalRun, handoffName)
		input.Handoff.Locator = handoffLocator
		input.Handoff.ReplyRoot = filepath.Join(finalRun, "reply")
		activityDir := filepath.Join(root, "activities", input.Handoff.Identity.ActivityID)
		if err := store.ensureDirectory(activityDir); err != nil {
			return err
		}
		if err := store.ensureDirectory(filepath.Join(activityDir, "events")); err != nil {
			return err
		}
		if err := store.ensureLockFile(filepath.Join(activityDir, "activity.lock")); err != nil {
			return err
		}
		if !input.DeferActivityPublication {
			document := canonicaljson.Object{{Name: "document_kind", Value: "ply.workflow.handoff"}, {Name: "document_id", Value: input.Handoff.Identity.HandoffID}, {Name: "sha256", Value: input.Handoff.SHA256}, {Name: "locator", Value: relativeLocator(root, handoffLocator)}}
			if _, err := store.publishActivityEvent(root, input.Handoff, "published", document, nil, nil, nil); err != nil {
				return err
			}
		}
		publication := envelope("ply.workflow.publication-ref",
			canonicaljson.Member{Name: "publication_key_sha256", Value: digestBytes([]byte(input.PublicationKey))},
			canonicaljson.Member{Name: "source_draft_sha256", Value: input.SourceDraftSHA256},
			canonicaljson.Member{Name: "activity_id", Value: input.Handoff.Identity.ActivityID}, canonicaljson.Member{Name: "run_id", Value: runID}, canonicaljson.Member{Name: "handoff_id", Value: input.Handoff.Identity.HandoffID}, canonicaljson.Member{Name: "handoff_sha256", Value: input.Handoff.SHA256}, canonicaljson.Member{Name: "locator", Value: relativeLocator(root, handoffLocator)})
		publicationBytes, _ := canonicaljson.Marshal(publication)
		if _, err := store.atomicWriteOnce(publicationPath, publicationBytes, 0o400); err != nil {
			return err
		}
		if !input.DeferActivityPublication {
			if err := store.writeHead(root, input.Handoff, "ready"); err != nil {
				return err
			}
		}
		result = CreateStoreResult{Snapshot: Snapshot{WorkspaceRoot: input.WorkspaceRoot, Handoff: input.Handoff, HeadState: "ready"}, Created: true}
		return nil
	}()
	var releaseErr error
	if lock != nil {
		releaseErr = store.release(lock)
	}
	if operationErr != nil {
		return CreateStoreResult{}, combineStoreErrors(operationErr, releaseErr)
	}
	if releaseErr != nil {
		return CreateStoreResult{}, releaseErr
	}
	return result, nil
}

func (store *systemStore) findRunByPublication(root, publicationKey, sourceDraftSHA256 string) (handoffDocument, bool, error) {
	runs := filepath.Join(root, "runs")
	entries, err := store.ops.readDir(runs)
	if err != nil {
		return handoffDocument{}, false, storeIO("read runs", runs, err)
	}
	var found *handoffDocument
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		files, err := store.ops.readDir(filepath.Join(runs, entry.Name()))
		if err != nil {
			return handoffDocument{}, false, storeIO("read run", entry.Name(), err)
		}
		for _, file := range files {
			if !strings.HasPrefix(file.Name(), "handoff.") || !strings.HasSuffix(file.Name(), ".json") {
				continue
			}
			document, err := store.readHandoff(filepath.Join(runs, entry.Name(), file.Name()))
			if err != nil {
				return handoffDocument{}, false, err
			}
			if document.Identity.PublicationKey != publicationKey {
				continue
			}
			if document.SourceDraftSHA256 != sourceDraftSHA256 {
				return handoffDocument{}, false, classified(ErrorConflict, fmt.Sprintf("publication key %s already has different content", publicationKey), nil)
			}
			if found != nil {
				return handoffDocument{}, false, storeConflict(runs, "multiple complete runs have the same publication key", nil)
			}
			copy := document
			found = &copy
		}
	}
	if found == nil {
		return handoffDocument{}, false, nil
	}
	return *found, true, nil
}

func (store *systemStore) SubmitStart(input SubmitStoreInput) (SubmitStoreResult, error) {
	return store.submit(input, "start")
}
func (store *systemStore) SubmitResult(input SubmitStoreInput) (SubmitStoreResult, error) {
	return store.submit(input, "terminal")
}

func (store *systemStore) submit(input SubmitStoreInput, phase string) (SubmitStoreResult, error) {
	if input.Phase != "" && input.Phase != phase {
		return SubmitStoreResult{}, classified(ErrorSchemaInvalid, "submit phase mismatch", nil)
	}
	snapshot, err := store.ReadByLocator(input.Snapshot.Handoff.Locator)
	if err != nil {
		return SubmitStoreResult{}, err
	}
	if snapshot.HeadState != "ready" {
		lock, lockErr := store.acquire(filepath.Join(snapshot.Handoff.ReplyRoot, "reply.lock"))
		if lockErr == nil {
			objectPath := filepath.Join(snapshot.Handoff.ReplyRoot, phase, "objects", strings.TrimPrefix(input.Document.SHA256, "sha256:")+".json")
			_, _ = store.atomicWriteOnce(objectPath, input.Document.Bytes, 0o400)
			input.Document.Locator = objectPath
			_, _ = store.publishReplyEvent(snapshot, phase, "late", input.Document, ErrorStale, "handoff lifecycle is closed")
			_ = store.release(lock)
		}
		return SubmitStoreResult{}, classified(ErrorStale, fmt.Sprintf("handoff %s is %s", snapshot.Handoff.Identity.HandoffID, snapshot.HeadState), nil)
	}
	if phase == "terminal" && snapshot.Start == nil {
		return SubmitStoreResult{}, classified(ErrorStartRequired, "a terminal result requires an accepted start receipt", nil)
	}
	reply := snapshot.Handoff.ReplyRoot
	lock, err := store.acquire(filepath.Join(reply, "reply.lock"))
	if err != nil {
		return SubmitStoreResult{}, err
	}
	var result SubmitStoreResult
	operationErr := func() error {
		current, err := store.ReadByLocator(input.Snapshot.Handoff.Locator)
		if err != nil {
			return err
		}
		snapshot = current
		if snapshot.HeadState != "ready" {
			objectPath := filepath.Join(reply, phase, "objects", strings.TrimPrefix(input.Document.SHA256, "sha256:")+".json")
			_, _ = store.atomicWriteOnce(objectPath, input.Document.Bytes, 0o400)
			input.Document.Locator = objectPath
			_, _ = store.publishReplyEvent(snapshot, phase, "late", input.Document, ErrorStale, "handoff lifecycle is closed")
			return classified(ErrorStale, fmt.Sprintf("handoff %s is %s", snapshot.Handoff.Identity.HandoffID, snapshot.HeadState), nil)
		}
		if phase == "terminal" && snapshot.Start == nil {
			return classified(ErrorStartRequired, "a terminal result requires an accepted start receipt", nil)
		}
		acceptedPath := filepath.Join(reply, phase, "accepted.ref")
		if _, err := store.ops.lstat(acceptedPath); err == nil {
			fields, err := store.readRef(acceptedPath, "ply.workflow.accepted-ref", "activity_id", "run_id", "handoff_id", "phase", "document_id", "document_sha256", "locator", "attempt_event_sha256")
			if err != nil {
				return err
			}
			id, _ := stringField(fields, "document_id", "accepted ref")
			digest, _ := stringField(fields, "document_sha256", "accepted ref")
			locator, _ := stringField(fields, "locator", "accepted ref")
			if id == input.Document.DocumentID && digest == input.Document.SHA256 {
				resolved, err := resolveStoreLocator(storeRoot(snapshot.WorkspaceRoot), locator)
				if err != nil {
					return storeConflict(acceptedPath, "accepted ref locator is invalid", err)
				}
				existing, err := store.readAccepted(resolved, snapshot.Handoff, phase)
				if err != nil {
					return err
				}
				result = SubmitStoreResult{Document: *existing, Created: false}
				return nil
			}
			objectPath := filepath.Join(reply, phase, "objects", strings.TrimPrefix(input.Document.SHA256, "sha256:")+".json")
			_, _ = store.atomicWriteOnce(objectPath, input.Document.Bytes, 0o400)
			_, _ = store.publishReplyEvent(snapshot, phase, "conflict", input.Document, ErrorConflict, "accepted slot already contains different bytes")
			return classified(ErrorConflict, fmt.Sprintf("%s slot already contains document %s", phase, id), nil)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return storeConflict(acceptedPath, "cannot inspect accepted ref", err)
		}
		if input.Revalidate != nil {
			if err := input.Revalidate(); err != nil {
				objectPath := filepath.Join(reply, phase, "objects", strings.TrimPrefix(input.Document.SHA256, "sha256:")+".json")
				_, _ = store.atomicWriteOnce(objectPath, input.Document.Bytes, 0o400)
				input.Document.Locator = objectPath
				_, _ = store.publishReplyEvent(snapshot, phase, "conflict", input.Document, classOf(err), err.Error())
				return err
			}
		}
		if phase == "terminal" {
			for _, artifact := range input.Artifacts {
				destination := filepath.Join(reply, "artifacts", "sha256", strings.TrimPrefix(artifact.SHA256, "sha256:"))
				if _, err := store.atomicWriteOnce(destination, artifact.Bytes, 0o400); err != nil {
					return err
				}
			}
		}
		objectPath := filepath.Join(reply, phase, "objects", strings.TrimPrefix(input.Document.SHA256, "sha256:")+".json")
		if _, err := store.atomicWriteOnce(objectPath, input.Document.Bytes, 0o400); err != nil {
			return err
		}
		input.Document.Locator = objectPath
		eventDigest, err := store.publishReplyEvent(snapshot, phase, "accepted", input.Document, "", "")
		if err != nil {
			return err
		}
		ref := envelope("ply.workflow.accepted-ref", canonicaljson.Member{Name: "activity_id", Value: snapshot.Handoff.Identity.ActivityID}, canonicaljson.Member{Name: "run_id", Value: snapshot.Handoff.Identity.RunID}, canonicaljson.Member{Name: "handoff_id", Value: snapshot.Handoff.Identity.HandoffID}, canonicaljson.Member{Name: "phase", Value: phase}, canonicaljson.Member{Name: "document_id", Value: input.Document.DocumentID}, canonicaljson.Member{Name: "document_sha256", Value: input.Document.SHA256}, canonicaljson.Member{Name: "locator", Value: relativeLocator(storeRoot(snapshot.WorkspaceRoot), objectPath)}, canonicaljson.Member{Name: "attempt_event_sha256", Value: eventDigest})
		refBytes, _ := canonicaljson.Marshal(ref)
		if _, err := store.atomicWriteOnce(acceptedPath, refBytes, 0o400); err != nil {
			return err
		}
		result = SubmitStoreResult{Document: input.Document, Created: true}
		return nil
	}()
	releaseErr := store.release(lock)
	if operationErr != nil {
		return SubmitStoreResult{}, combineStoreErrors(operationErr, releaseErr)
	}
	if releaseErr != nil {
		return SubmitStoreResult{}, releaseErr
	}
	return result, nil
}

func (store *systemStore) Cancel(input ControlStoreInput) (ControlStoreResult, error) {
	input.State = "cancelled"
	return store.control(input, false)
}
func (store *systemStore) Abandon(input ControlStoreInput) (ControlStoreResult, error) {
	input.State = "abandoned_unknown"
	return store.control(input, true)
}
func (store *systemStore) control(input ControlStoreInput, requiresStart bool) (ControlStoreResult, error) {
	snapshot, err := store.ReadByID(input.WorkspaceRoot, input.HandoffID)
	if err != nil {
		return ControlStoreResult{}, err
	}
	if snapshot.HeadState != "ready" {
		return ControlStoreResult{}, classified(ErrorStale, fmt.Sprintf("handoff %s is %s", input.HandoffID, snapshot.HeadState), nil)
	}
	if requiresStart {
		if snapshot.Start == nil || snapshot.Start.Outcome != "started" || snapshot.Terminal != nil {
			return ControlStoreResult{}, classified(ErrorStale, "handoff cannot be abandoned in its current state", nil)
		}
	} else if snapshot.Start != nil {
		return ControlStoreResult{}, classified(ErrorStale, "started handoff cannot be cancelled", nil)
	}
	activityLock, err := store.acquire(filepath.Join(storeRoot(input.WorkspaceRoot), "activities", snapshot.Handoff.Identity.ActivityID, "activity.lock"))
	if err != nil {
		return ControlStoreResult{}, err
	}
	replyLock, err := store.acquire(filepath.Join(snapshot.Handoff.ReplyRoot, "reply.lock"))
	if err != nil {
		_ = store.release(activityLock)
		return ControlStoreResult{}, err
	}
	current, err := store.ReadByID(input.WorkspaceRoot, input.HandoffID)
	if err != nil {
		_ = store.release(replyLock)
		_ = store.release(activityLock)
		return ControlStoreResult{}, err
	}
	if current.HeadState != "ready" || (requiresStart && (current.Start == nil || current.Start.Outcome != "started" || current.Terminal != nil)) || (!requiresStart && current.Start != nil) {
		_ = store.release(replyLock)
		_ = store.release(activityLock)
		return ControlStoreResult{}, classified(ErrorStale, "handoff cannot be controlled in its current state", nil)
	}
	snapshot = current
	eventType := "cancelled"
	if requiresStart {
		eventType = "abandoned_unknown"
	}
	_, operationErr := store.publishActivityEvent(storeRoot(input.WorkspaceRoot), snapshot.Handoff, eventType, nil, input.Reason, nil, nil)
	if operationErr == nil {
		operationErr = store.writeHead(storeRoot(input.WorkspaceRoot), snapshot.Handoff, input.State)
	}
	replyRelease := store.release(replyLock)
	activityRelease := store.release(activityLock)
	if operationErr != nil {
		return ControlStoreResult{}, combineStoreErrors(operationErr, combineStoreErrors(replyRelease, activityRelease))
	}
	if replyRelease != nil {
		return ControlStoreResult{}, replyRelease
	}
	if activityRelease != nil {
		return ControlStoreResult{}, activityRelease
	}
	snapshot.HeadState = input.State
	return ControlStoreResult{Snapshot: snapshot}, nil
}

func (store *systemStore) Supersede(input SupersedeStoreInput) (CreateStoreResult, error) {
	root := storeRoot(input.WorkspaceRoot)
	storeLock, err := store.acquire(filepath.Join(root, "store.lock"))
	if err != nil {
		return CreateStoreResult{}, err
	}
	var activityLock, replyLock *os.File
	var created CreateStoreResult
	operationErr := func() error {
		old, err := store.ReadByID(input.WorkspaceRoot, input.OldHandoffID)
		if err != nil {
			return err
		}
		if old.HeadState != "ready" || old.Start != nil {
			return classified(ErrorStale, "only an unstarted ready handoff can be superseded", nil)
		}
		if input.Create.WorkspaceRoot != input.WorkspaceRoot || input.Create.Handoff.Identity.ActivityID != old.Handoff.Identity.ActivityID || !input.Create.DeferActivityPublication {
			return classified(ErrorConflict, "replacement is not bound to the existing activity", nil)
		}
		activityLock, err = store.acquire(filepath.Join(root, "activities", old.Handoff.Identity.ActivityID, "activity.lock"))
		if err != nil {
			return err
		}
		replyLock, err = store.acquire(filepath.Join(old.Handoff.ReplyRoot, "reply.lock"))
		if err != nil {
			return err
		}
		current, err := store.ReadByID(input.WorkspaceRoot, input.OldHandoffID)
		if err != nil {
			return err
		}
		if current.HeadState != "ready" || current.Start != nil {
			return classified(ErrorStale, "only an unstarted ready handoff can be superseded", nil)
		}
		input.Create.StoreAlreadyLocked = true
		created, err = store.Create(input.Create)
		if err != nil {
			return err
		}
		document := canonicaljson.Object{{Name: "document_kind", Value: "ply.workflow.handoff"}, {Name: "document_id", Value: created.Snapshot.Handoff.Identity.HandoffID}, {Name: "sha256", Value: created.Snapshot.Handoff.SHA256}, {Name: "locator", Value: relativeLocator(root, created.Snapshot.Handoff.Locator)}}
		if _, err = store.publishActivityEvent(root, current.Handoff, "superseded", document, input.Reason, created.Snapshot.Handoff.Identity.RunID, nil); err != nil {
			return err
		}
		return store.writeHead(root, created.Snapshot.Handoff, "ready")
	}()
	var replyRelease, activityRelease error
	if replyLock != nil {
		replyRelease = store.release(replyLock)
	}
	if activityLock != nil {
		activityRelease = store.release(activityLock)
	}
	storeRelease := store.release(storeLock)
	combinedRelease := combineStoreErrors(replyRelease, combineStoreErrors(activityRelease, storeRelease))
	if operationErr != nil {
		return CreateStoreResult{}, combineStoreErrors(operationErr, combinedRelease)
	}
	if combinedRelease != nil {
		return CreateStoreResult{}, combinedRelease
	}
	return created, nil
}

func (store *systemStore) ReadByID(workspaceRoot string, id HandoffID) (Snapshot, error) {
	runs := filepath.Join(storeRoot(workspaceRoot), "runs")
	entries, err := store.ops.readDir(runs)
	if err != nil {
		return Snapshot{}, storeIO("read runs", runs, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		files, err := store.ops.readDir(filepath.Join(runs, entry.Name()))
		if err != nil {
			return Snapshot{}, storeIO("read run", entry.Name(), err)
		}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "handoff.") && strings.HasSuffix(file.Name(), ".json") {
				locator := filepath.Join(runs, entry.Name(), file.Name())
				document, err := store.readHandoff(locator)
				if err != nil {
					return Snapshot{}, err
				}
				if document.Identity.HandoffID == string(id) {
					return store.readSnapshot(workspaceRoot, document)
				}
			}
		}
	}
	return Snapshot{}, classified(ErrorNotFound, fmt.Sprintf("handoff %s was not found", id), nil)
}

func (store *systemStore) ReadByLocator(locator string) (Snapshot, error) {
	if !filepath.IsAbs(locator) || filepath.Clean(locator) != locator {
		return Snapshot{}, classified(ErrorNotFound, "handoff locator must be absolute and clean", nil)
	}
	physical, err := store.files.EvalSymlinks(locator)
	if err != nil || physical != locator {
		return Snapshot{}, classified(ErrorNotFound, "handoff locator must be a physical path", err)
	}
	document, err := store.readHandoff(locator)
	if err != nil {
		return Snapshot{}, err
	}
	root, workspaceRoot, ok := rootsFromLocator(locator)
	if !ok {
		return Snapshot{}, classified(ErrorNotFound, fmt.Sprintf("%s is not a canonical handoff locator", locator), nil)
	}
	expected := filepath.Join(root, "runs", document.Identity.RunID, "handoff."+strings.TrimPrefix(document.SHA256, "sha256:")+".json")
	if locator != expected || document.Workspace.Root != workspaceRoot {
		return Snapshot{}, storeConflict(locator, "handoff locator or workspace binding is not canonical", nil)
	}
	return store.readSnapshot(workspaceRoot, document)
}

func (store *systemStore) readSnapshot(workspaceRoot string, handoff handoffDocument) (Snapshot, error) {
	snapshot := Snapshot{WorkspaceRoot: workspaceRoot, Handoff: handoff, HeadState: "ready"}
	root := storeRoot(workspaceRoot)
	headPath := filepath.Join(root, "activities", handoff.Identity.ActivityID, "head.ref")
	if _, err := store.ops.lstat(headPath); err == nil {
		fields, err := store.readRef(headPath, "ply.workflow.activity-head-ref", "activity_id", "run_id", "handoff_id", "state", "event_ordinal", "event_sha256", "handoff_sha256", "locator")
		if err != nil {
			return Snapshot{}, err
		}
		headID, _ := stringField(fields, "handoff_id", "head ref")
		if objectMapString(fields, "activity_id") != handoff.Identity.ActivityID {
			return Snapshot{}, storeConflict(headPath, "head activity does not match the handoff", nil)
		}
		state := objectMapString(fields, "state")
		if !setOf("ready", "cancelled", "superseded", "abandoned_unknown")[state] {
			return Snapshot{}, storeConflict(headPath, "head state is invalid", nil)
		}
		ordinal, ok := fields["event_ordinal"].(int64)
		eventDigest := objectMapString(fields, "event_sha256")
		if !ok || ordinal < 1 || !validateDigest(eventDigest) {
			return Snapshot{}, storeConflict(headPath, "head event binding is invalid", nil)
		}
		eventName := fmt.Sprintf("%020d-%s.json", ordinal, strings.TrimPrefix(eventDigest, "sha256:"))
		eventBytes, err := store.ops.readFile(filepath.Join(root, "activities", handoff.Identity.ActivityID, "events", eventName))
		if err != nil || digestBytes(eventBytes) != eventDigest {
			return Snapshot{}, storeConflict(headPath, "head event is missing or has different bytes", err)
		}
		headLocator, err := resolveStoreLocator(root, objectMapString(fields, "locator"))
		if err != nil {
			return Snapshot{}, storeConflict(headPath, "head locator is invalid", err)
		}
		headHandoff, err := store.readHandoff(headLocator)
		if err != nil {
			return Snapshot{}, err
		}
		if headHandoff.Identity.HandoffID != headID || headHandoff.Identity.RunID != objectMapString(fields, "run_id") || headHandoff.SHA256 != objectMapString(fields, "handoff_sha256") {
			return Snapshot{}, storeConflict(headPath, "head ref does not match its handoff", nil)
		}
		if headID != handoff.Identity.HandoffID {
			snapshot.HeadState = "superseded"
			snapshot.ReplacementID = headID
		} else {
			snapshot.HeadState = state
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return Snapshot{}, storeConflict(headPath, "cannot inspect head", err)
	}
	for _, phase := range []string{"start", "terminal"} {
		acceptedPath := filepath.Join(handoff.ReplyRoot, phase, "accepted.ref")
		if _, err := store.ops.lstat(acceptedPath); err == nil {
			fields, err := store.readRef(acceptedPath, "ply.workflow.accepted-ref", "activity_id", "run_id", "handoff_id", "phase", "document_id", "document_sha256", "locator", "attempt_event_sha256")
			if err != nil {
				return Snapshot{}, err
			}
			if objectMapString(fields, "activity_id") != handoff.Identity.ActivityID || objectMapString(fields, "run_id") != handoff.Identity.RunID || objectMapString(fields, "handoff_id") != handoff.Identity.HandoffID || objectMapString(fields, "phase") != phase {
				return Snapshot{}, storeConflict(acceptedPath, "accepted ref binding does not match the handoff", nil)
			}
			relative, _ := stringField(fields, "locator", "accepted ref")
			resolved, err := resolveStoreLocator(root, relative)
			if err != nil {
				return Snapshot{}, storeConflict(acceptedPath, "accepted ref locator is not a clean relative path", err)
			}
			document, err := store.readAccepted(resolved, handoff, phase)
			if err != nil {
				return Snapshot{}, err
			}
			if objectMapString(fields, "document_id") != document.DocumentID || objectMapString(fields, "document_sha256") != document.SHA256 {
				return Snapshot{}, storeConflict(acceptedPath, "accepted ref document binding does not match stored bytes", nil)
			}
			if phase == "start" {
				snapshot.Start = document
			} else {
				if err := store.validateStoredArtifacts(handoff, document); err != nil {
					return Snapshot{}, err
				}
				snapshot.Terminal = document
			}
		} else if !errors.Is(err, fs.ErrNotExist) {
			return Snapshot{}, storeConflict(acceptedPath, "cannot inspect accepted ref", err)
		}
		attempts, err := store.readAttempts(handoff, phase)
		if err != nil {
			return Snapshot{}, err
		}
		snapshot.Attempts = append(snapshot.Attempts, attempts...)
	}
	return snapshot, nil
}

func (store *systemStore) validateStoredArtifacts(handoff handoffDocument, terminal *acceptedDocument) error {
	value, found := objectMember(terminal.Value, "artifacts")
	if !found {
		return nil
	}
	artifacts, ok := value.([]canonicaljson.Value)
	if !ok {
		return storeConflict(terminal.Locator, "terminal artifacts are invalid", nil)
	}
	for _, raw := range artifacts {
		artifact := raw.(canonicaljson.Object)
		if objectString(artifact, "kind") != "managed" {
			continue
		}
		digest := objectString(artifact, "sha256")
		locator := objectString(artifact, "locator")
		expected := filepath.Join(handoff.ReplyRoot, "artifacts", "sha256", strings.TrimPrefix(digest, "sha256:"))
		if locator != expected {
			return storeConflict(terminal.Locator, "managed artifact locator is not canonical", nil)
		}
		info, err := store.ops.lstat(locator)
		if err != nil || !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
			return storeConflict(locator, "managed artifact is missing or not regular", err)
		}
		contents, err := store.ops.readFile(locator)
		if err != nil {
			return storeIO("read managed artifact", locator, err)
		}
		if digestBytes(contents) != digest || int64(len(contents)) != objectInt(artifact, "size_bytes") {
			return storeConflict(locator, "managed artifact bytes do not match the terminal result", nil)
		}
	}
	return nil
}

func (store *systemStore) readAttempts(handoff handoffDocument, phase string) ([]Attempt, error) {
	directory := filepath.Join(handoff.ReplyRoot, phase, "events")
	entries, err := store.ops.readDir(directory)
	if err != nil {
		return nil, storeIO("read reply events", directory, err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	result := make([]Attempt, 0, len(entries))
	for sequence, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			return nil, storeConflict(filepath.Join(directory, entry.Name()), "unexpected reply event entry", nil)
		}
		path := filepath.Join(directory, entry.Name())
		contents, err := store.ops.readFile(path)
		if err != nil {
			return nil, storeIO("read reply event", path, err)
		}
		value, err := canonicaljson.DecodeStrict(contents)
		if err != nil {
			return nil, storeConflict(path, "reply event is invalid", err)
		}
		fields, err := exactObject(value, "reply event", "kind", "schema_version", "format", "format_version", "canonicalization", "activity_id", "run_id", "handoff_id", "phase", "classification", "occurred_at_utc", "document_id", "digest", "size_bytes", "locator", "error_class", "detail")
		if err != nil {
			return nil, storeConflict(path, err.Error(), err)
		}
		if err := validateEnvelope(fields, "ply.workflow.reply-attempt-event", "reply event"); err != nil {
			return nil, storeConflict(path, err.Error(), err)
		}
		if objectMapString(fields, "activity_id") != handoff.Identity.ActivityID || objectMapString(fields, "run_id") != handoff.Identity.RunID || objectMapString(fields, "handoff_id") != handoff.Identity.HandoffID || objectMapString(fields, "phase") != phase {
			return nil, storeConflict(path, "reply event binding does not match the handoff", nil)
		}
		expectedName := fmt.Sprintf("%020d-%s.json", sequence+1, strings.TrimPrefix(digestBytes(contents), "sha256:"))
		if entry.Name() != expectedName || !setOf("accepted", "rejected", "conflict", "late")[objectMapString(fields, "classification")] || !validateUTC(objectMapString(fields, "occurred_at_utc")) {
			return nil, storeConflict(path, "reply event metadata or filename is invalid", nil)
		}
		size, ok := fields["size_bytes"].(int64)
		if !ok || size < 0 {
			return nil, storeConflict(path, "reply event size is invalid", nil)
		}
		if fields["digest"] != nil {
			digest, ok := fields["digest"].(string)
			if !ok || !validateDigest(digest) {
				return nil, storeConflict(path, "reply event digest is invalid", nil)
			}
		}
		attempt := Attempt{Phase: phase, Classification: objectMapString(fields, "classification"), ErrorClass: objectMapString(fields, "error_class"), Detail: objectMapString(fields, "detail"), Sequence: int64(sequence + 1)}
		attempt.SizeBytes = size
		if fields["document_id"] != nil {
			attempt.DocumentID, _ = fields["document_id"].(string)
		}
		if fields["digest"] != nil {
			attempt.Digest, _ = fields["digest"].(string)
		}
		if fields["locator"] != nil {
			attempt.Locator, _ = fields["locator"].(string)
		}
		result = append(result, attempt)
	}
	return result, nil
}

func (store *systemStore) readHandoff(locator string) (handoffDocument, error) {
	info, err := store.ops.lstat(locator)
	if err != nil {
		return handoffDocument{}, classified(ErrorNotFound, fmt.Sprintf("handoff %s was not found", locator), err)
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return handoffDocument{}, storeConflict(locator, "handoff is not a regular file", nil)
	}
	bytes, err := store.ops.readFile(locator)
	if err != nil {
		return handoffDocument{}, storeIO("read handoff", locator, err)
	}
	digest := digestBytes(bytes)
	base := filepath.Base(locator)
	if base != "handoff."+strings.TrimPrefix(digest, "sha256:")+".json" {
		return handoffDocument{}, storeConflict(locator, "handoff filename digest does not match bytes", nil)
	}
	value, err := canonicaljson.DecodeStrict(bytes)
	if err != nil {
		return handoffDocument{}, storeConflict(locator, "handoff JSON is invalid", err)
	}
	object, ok := value.(canonicaljson.Object)
	if !ok {
		return handoffDocument{}, storeConflict(locator, "handoff JSON is not an object", nil)
	}
	if _, hasKind := objectMember(object, "kind"); hasKind {
		if err := validateStoredHandoff(object); err != nil {
			return handoffDocument{}, storeConflict(locator, "handoff schema is invalid", err)
		}
	}
	document, err := handoffCore(object)
	if err != nil {
		return handoffDocument{}, storeConflict(locator, err.Error(), err)
	}
	expectedReplyRoot := filepath.Join(filepath.Dir(locator), "reply")
	if document.ReplyRoot != "" && document.ReplyRoot != expectedReplyRoot {
		return handoffDocument{}, storeConflict(locator, "reply root does not match the canonical run path", nil)
	}
	document.Value = object
	document.Bytes = bytes
	document.SHA256 = digest
	document.Locator = locator
	document.ReplyRoot = expectedReplyRoot
	return document, nil
}

func handoffCore(value canonicaljson.Object) (handoffDocument, error) {
	identityValue, found := objectMember(value, "identity")
	if !found {
		return handoffDocument{}, fmt.Errorf("handoff identity is missing")
	}
	idFields, err := exactObject(identityValue, "identity", "activity_id", "run_id", "handoff_id", "start_receipt_id", "terminal_result_id", "publication_key", "activity_key", "created_at_utc")
	if err != nil {
		return handoffDocument{}, err
	}
	get := func(name string) string { value, _ := stringField(idFields, name, "identity"); return value }
	document := handoffDocument{Identity: identity{ActivityID: get("activity_id"), RunID: get("run_id"), HandoffID: get("handoff_id"), StartReceiptID: get("start_receipt_id"), TerminalResultID: get("terminal_result_id"), PublicationKey: get("publication_key"), ActivityKey: get("activity_key"), CreatedAtUTC: get("created_at_utc")}}
	if source, ok := objectMember(value, "source_draft_sha256"); ok {
		document.SourceDraftSHA256, _ = source.(string)
	}
	if goalValue, ok := objectMember(value, "goal"); ok {
		if goal, ok := goalValue.(canonicaljson.Object); ok {
			title, _ := objectMember(goal, "title")
			document.GoalTitle, _ = title.(string)
		}
	}
	if budgetValue, ok := objectMember(value, "budget"); ok {
		if budget, ok := budgetValue.(canonicaljson.Object); ok {
			rounds, _ := objectMember(budget, "max_rounds")
			document.MaxRounds, _ = rounds.(int64)
		}
	}
	if replyValue, ok := objectMember(value, "reply_capability"); ok {
		reply, ok := replyValue.(canonicaljson.Object)
		if ok {
			capability, _ := objectMember(reply, "capability_id")
			secret, _ := objectMember(reply, "secret")
			root, _ := objectMember(reply, "reply_root")
			document.ReplyCapabilityID, _ = capability.(string)
			document.ReplySecret, _ = secret.(string)
			document.ReplyRoot, _ = root.(string)
		}
	}
	if workspaceValue, ok := objectMember(value, "workspace_binding"); ok {
		binding, ok := workspaceValue.(canonicaljson.Object)
		if ok {
			document.Workspace.Root = objectString(binding, "root")
			if versionValue, found := objectMember(binding, "marker_format_version"); found {
				if version, ok := versionValue.(int64); ok {
					document.Workspace.MarkerFormatVersion = int(version)
				}
			}
			document.Workspace.MarkerSHA256 = objectString(binding, "marker_sha256")
		}
	}
	if projectValue, ok := objectMember(value, "project_binding"); ok {
		binding, ok := projectValue.(canonicaljson.Object)
		if ok {
			document.ProjectID = workspace.ProjectID(objectString(binding, "project_id"))
			document.RepoID = workspace.RepoID(objectString(binding, "repo_id"))
			document.RegisteredLocator = objectString(binding, "registered_locator")
			document.RegisteredGitCommonDir = objectString(binding, "registered_git_common_dir")
		}
	}
	if targetValue, ok := objectMember(value, "target_binding"); ok {
		target, ok := targetValue.(canonicaljson.Object)
		if ok {
			document.Target.Worktree = objectString(target, "worktree")
			document.Target.GitCommonDir = objectString(target, "git_common_dir")
			document.Target.Ref = objectString(target, "ref")
			document.Target.OID = objectString(target, "oid")
			document.Target.Tree = objectString(target, "tree")
			if statusValue, found := objectMember(target, "status_policy"); found {
				document.Target.Status, _ = validateStatusPolicy(statusValue)
			}
		}
	}
	return document, nil
}

func (store *systemStore) readAccepted(locator string, handoff handoffDocument, phase string) (*acceptedDocument, error) {
	bytes, err := store.ops.readFile(locator)
	if err != nil {
		return nil, storeIO("read accepted document", locator, err)
	}
	value, err := canonicaljson.DecodeStrict(bytes)
	if err != nil {
		return nil, storeConflict(locator, "accepted document is invalid", err)
	}
	object, ok := value.(canonicaljson.Object)
	if !ok {
		return nil, storeConflict(locator, "accepted document is not an object", nil)
	}
	if _, hasKind := objectMember(handoff.Value, "kind"); hasKind {
		if err := validateStoredAccepted(object, phase, handoff.ReplyCapabilityID, handoff.ReplySecret); err != nil {
			return nil, storeConflict(locator, "accepted document schema or proof is invalid", err)
		}
	}
	document := &acceptedDocument{Value: object, Bytes: bytes, SHA256: digestBytes(bytes), Locator: locator}
	if filepath.Base(locator) != strings.TrimPrefix(document.SHA256, "sha256:")+".json" {
		return nil, storeConflict(locator, "accepted document filename digest does not match bytes", nil)
	}
	for _, name := range []string{"receipt_id", "result_id", "id"} {
		if v, ok := objectMember(object, name); ok {
			document.DocumentID, _ = v.(string)
		}
	}
	for _, name := range []string{"acceptance", "reported_outcome", "outcome"} {
		if v, ok := objectMember(object, name); ok {
			document.Outcome, _ = v.(string)
		}
	}
	document.Summary = objectString(object, "summary")
	document.Meaning = objectString(object, "meaning")
	return document, nil
}

func (store *systemStore) publishActivityEvent(root string, handoff handoffDocument, eventType string, document canonicaljson.Value, reason any, related any, eventError any) (string, error) {
	eventsDir := filepath.Join(root, "activities", handoff.Identity.ActivityID, "events")
	entries, err := store.ops.readDir(eventsDir)
	if err != nil {
		return "", storeIO("read activity events", eventsDir, err)
	}
	ordinal := int64(len(entries) + 1)
	if document == nil {
		document = nil
	}
	if reason == nil {
		reason = nil
	}
	if related == nil {
		related = nil
	}
	if eventError == nil {
		eventError = nil
	}
	event := envelope("ply.workflow.activity-event", canonicaljson.Member{Name: "activity_id", Value: handoff.Identity.ActivityID}, canonicaljson.Member{Name: "ordinal", Value: ordinal}, canonicaljson.Member{Name: "event_type", Value: eventType}, canonicaljson.Member{Name: "occurred_at_utc", Value: store.now()}, canonicaljson.Member{Name: "run_id", Value: handoff.Identity.RunID}, canonicaljson.Member{Name: "handoff_id", Value: handoff.Identity.HandoffID}, canonicaljson.Member{Name: "document", Value: document}, canonicaljson.Member{Name: "reason", Value: reason}, canonicaljson.Member{Name: "related_run_id", Value: related}, canonicaljson.Member{Name: "error", Value: eventError})
	bytes, _ := canonicaljson.Marshal(event)
	digest := digestBytes(bytes)
	name := fmt.Sprintf("%020d-%s.json", ordinal, strings.TrimPrefix(digest, "sha256:"))
	_, err = store.atomicWriteOnce(filepath.Join(eventsDir, name), bytes, 0o400)
	return digest, err
}

func (store *systemStore) publishReplyEvent(snapshot Snapshot, phase, classification string, document acceptedDocument, errorClass ErrorClass, detail string) (string, error) {
	eventsDir := filepath.Join(snapshot.Handoff.ReplyRoot, phase, "events")
	entries, err := store.ops.readDir(eventsDir)
	if err != nil {
		return "", storeIO("read reply events", eventsDir, err)
	}
	sequence := int64(len(entries) + 1)
	var id, digest, locator canonicaljson.Value
	if document.DocumentID != "" {
		id = document.DocumentID
	}
	if document.SHA256 != "" {
		digest = document.SHA256
	}
	if document.Locator != "" {
		locator = document.Locator
	}
	event := envelope("ply.workflow.reply-attempt-event", canonicaljson.Member{Name: "activity_id", Value: snapshot.Handoff.Identity.ActivityID}, canonicaljson.Member{Name: "run_id", Value: snapshot.Handoff.Identity.RunID}, canonicaljson.Member{Name: "handoff_id", Value: snapshot.Handoff.Identity.HandoffID}, canonicaljson.Member{Name: "phase", Value: phase}, canonicaljson.Member{Name: "classification", Value: classification}, canonicaljson.Member{Name: "occurred_at_utc", Value: store.now()}, canonicaljson.Member{Name: "document_id", Value: id}, canonicaljson.Member{Name: "digest", Value: digest}, canonicaljson.Member{Name: "size_bytes", Value: int64(len(document.Bytes))}, canonicaljson.Member{Name: "locator", Value: locator}, canonicaljson.Member{Name: "error_class", Value: string(errorClass)}, canonicaljson.Member{Name: "detail", Value: detail})
	bytes, _ := canonicaljson.Marshal(event)
	eventDigest := digestBytes(bytes)
	name := fmt.Sprintf("%020d-%s.json", sequence, strings.TrimPrefix(eventDigest, "sha256:"))
	_, err = store.atomicWriteOnce(filepath.Join(eventsDir, name), bytes, 0o400)
	return eventDigest, err
}

func (store *systemStore) writeHead(root string, handoff handoffDocument, state string) error {
	eventsDir := filepath.Join(root, "activities", handoff.Identity.ActivityID, "events")
	entries, err := store.ops.readDir(eventsDir)
	if err != nil || len(entries) == 0 {
		return storeIO("read head event", eventsDir, err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	eventBytes, err := store.ops.readFile(filepath.Join(eventsDir, entries[len(entries)-1].Name()))
	if err != nil {
		return storeIO("read head event", eventsDir, err)
	}
	ordinal := int64(len(entries))
	head := envelope("ply.workflow.activity-head-ref", canonicaljson.Member{Name: "activity_id", Value: handoff.Identity.ActivityID}, canonicaljson.Member{Name: "run_id", Value: handoff.Identity.RunID}, canonicaljson.Member{Name: "handoff_id", Value: handoff.Identity.HandoffID}, canonicaljson.Member{Name: "state", Value: state}, canonicaljson.Member{Name: "event_ordinal", Value: ordinal}, canonicaljson.Member{Name: "event_sha256", Value: digestBytes(eventBytes)}, canonicaljson.Member{Name: "handoff_sha256", Value: handoff.SHA256}, canonicaljson.Member{Name: "locator", Value: relativeLocator(root, handoff.Locator)})
	bytes, _ := canonicaljson.Marshal(head)
	return store.atomicReplace(filepath.Join(root, "activities", handoff.Identity.ActivityID, "head.ref"), bytes, 0o600)
}

func (store *systemStore) publishReplyRejection(snapshot Snapshot, phase, classification string, raw []byte, class ErrorClass, detail string) error {
	lock, err := store.acquire(filepath.Join(snapshot.Handoff.ReplyRoot, "reply.lock"))
	if err != nil {
		return err
	}
	operationErr := func() error {
		current, err := store.ReadByLocator(snapshot.Handoff.Locator)
		if err != nil {
			return err
		}
		if current.HeadState != "ready" || current.Terminal != nil {
			classification = "late"
			class = ErrorStale
			detail = "handoff lifecycle is closed"
		} else if phase == "start" && current.Start != nil {
			classification = "conflict"
			class = ErrorConflict
			detail = "start slot already contains an accepted receipt"
		}
		digest := ""
		if raw != nil {
			digest = digestBytes(raw)
		}
		limit := maxResultBytes
		if phase == "start" {
			limit = maxStartBytes
		}
		if raw != nil && len(raw) <= limit {
			_, _ = store.atomicWriteOnce(filepath.Join(current.Handoff.ReplyRoot, phase, "rejected", strings.TrimPrefix(digest, "sha256:")+".bin"), raw, 0o400)
		}
		document := acceptedDocument{Bytes: raw, SHA256: digest}
		_, err = store.publishReplyEvent(current, phase, classification, document, class, detail)
		return err
	}()
	releaseErr := store.release(lock)
	return combineStoreErrors(operationErr, releaseErr)
}

func (store *systemStore) ensureDirectory(path string) error {
	if err := store.ops.mkdirAll(path, 0o700); err != nil {
		return storeIO("mkdir", path, err)
	}
	info, err := store.ops.lstat(path)
	if err != nil {
		return storeIO("inspect directory", path, err)
	}
	if !info.IsDir() || info.Mode()&fs.ModeSymlink != 0 {
		return storeConflict(path, "store path is not a regular directory", nil)
	}
	if err := store.ops.chmodPath(path, 0o700); err != nil {
		return storeIO("chmod directory", path, err)
	}
	return nil
}
func (store *systemStore) ensureLockFile(path string) error {
	file, err := store.ops.openFile(path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return storeIO("open lock", path, err)
	}
	if err := store.ops.chmodFile(file, 0o600); err != nil {
		_ = store.ops.closeFile(file)
		return storeIO("chmod lock", path, err)
	}
	return store.ops.closeFile(file)
}
func (store *systemStore) acquire(path string) (*os.File, error) {
	if err := store.ensureLockFile(path); err != nil {
		return nil, err
	}
	file, err := store.ops.openFile(path, os.O_RDWR, 0o600)
	if err != nil {
		return nil, storeIO("open lock", path, err)
	}
	opened, statErr := store.ops.statFile(file)
	actual, pathErr := store.ops.lstat(path)
	if statErr != nil || pathErr != nil || !opened.Mode().IsRegular() || !actual.Mode().IsRegular() || actual.Mode()&fs.ModeSymlink != 0 || !os.SameFile(opened, actual) {
		_ = store.ops.closeFile(file)
		return nil, storeConflict(path, "lock identity changed", firstError(statErr, pathErr))
	}
	if err := store.ops.lock(file); err != nil {
		_ = store.ops.closeFile(file)
		return nil, storeIO("lock", path, err)
	}
	return file, nil
}
func (store *systemStore) release(file *os.File) error {
	unlockErr := store.ops.unlock(file)
	closeErr := store.ops.closeFile(file)
	if unlockErr != nil {
		return storeIO("unlock", file.Name(), unlockErr)
	}
	if closeErr != nil {
		return storeIO("close lock", file.Name(), closeErr)
	}
	return nil
}
func (store *systemStore) atomicWriteOnce(destination string, contents []byte, mode fs.FileMode) (bool, error) {
	if existing, err := store.ops.readFile(destination); err == nil {
		if string(existing) == string(contents) {
			return false, nil
		}
		return false, classified(ErrorConflict, fmt.Sprintf("immutable destination %s already contains different bytes", destination), nil)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return false, storeIO("read destination", destination, err)
	}
	directory := filepath.Dir(destination)
	temporary, err := store.ops.createTemp(directory, ".ply-*.tmp")
	if err != nil {
		return false, storeIO("open temp", directory, err)
	}
	temporaryPath := temporary.Name()
	renamed := false
	abort := func(operation string, primary error) (bool, error) {
		_ = store.ops.closeFile(temporary)
		if !renamed {
			_ = store.ops.remove(temporaryPath)
		}
		return false, storeIO(operation, temporaryPath, primary)
	}
	written, err := store.ops.write(temporary, contents)
	if err != nil {
		return abort("write", err)
	}
	if written != len(contents) {
		return abort("write", io.ErrShortWrite)
	}
	if err := store.ops.chmodFile(temporary, mode); err != nil {
		return abort("chmod", err)
	}
	if err := store.ops.syncFile(temporary); err != nil {
		return abort("sync", err)
	}
	if err := store.ops.closeFile(temporary); err != nil {
		return abort("close", err)
	}
	if _, err := store.ops.lstat(destination); err == nil {
		return false, classified(ErrorConflict, fmt.Sprintf("immutable destination %s appeared during publish", destination), nil)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return false, storeIO("inspect destination", destination, err)
	}
	if err := store.ops.rename(temporaryPath, destination); err != nil {
		return abort("rename", err)
	}
	renamed = true
	if err := store.ops.syncDirectory(directory); err != nil {
		return false, storeIO("sync directory", directory, err)
	}
	return true, nil
}
func (store *systemStore) atomicReplace(destination string, contents []byte, mode fs.FileMode) error {
	directory := filepath.Dir(destination)
	temporary, err := store.ops.createTemp(directory, ".ply-*.tmp")
	if err != nil {
		return storeIO("open temp", directory, err)
	}
	path := temporary.Name()
	replaced := false
	defer func() {
		if !replaced {
			_ = store.ops.remove(path)
		}
	}()
	written, err := store.ops.write(temporary, contents)
	if err != nil || written != len(contents) {
		_ = store.ops.closeFile(temporary)
		if err == nil {
			err = io.ErrShortWrite
		}
		return storeIO("write", path, err)
	}
	if err := store.ops.chmodFile(temporary, mode); err != nil {
		_ = store.ops.closeFile(temporary)
		return storeIO("chmod", path, err)
	}
	if err := store.ops.syncFile(temporary); err != nil {
		_ = store.ops.closeFile(temporary)
		return storeIO("sync", path, err)
	}
	if err := store.ops.closeFile(temporary); err != nil {
		return storeIO("close", path, err)
	}
	if err := store.ops.replace(path, destination); err != nil {
		return storeIO("replace", destination, err)
	}
	replaced = true
	if err := store.ops.syncDirectory(directory); err != nil {
		return storeIO("sync directory", directory, err)
	}
	return nil
}

func (store *systemStore) readRef(path, kind string, names ...string) (map[string]canonicaljson.Value, error) {
	contents, err := store.ops.readFile(path)
	if err != nil {
		return nil, storeIO("read ref", path, err)
	}
	value, err := canonicaljson.DecodeStrict(contents)
	if err != nil {
		return nil, storeConflict(path, "ref is not canonical JSON", err)
	}
	fields, err := exactObject(value, "ref", append([]string{"kind", "schema_version", "format", "format_version", "canonicalization"}, names...)...)
	if err != nil {
		return nil, storeConflict(path, err.Error(), err)
	}
	if err := validateEnvelope(fields, kind, "ref"); err != nil {
		return nil, storeConflict(path, err.Error(), err)
	}
	return fields, nil
}
func (store *systemStore) now() string {
	clock := store.clock
	if clock == nil {
		return time.Now().UTC().Format(time.RFC3339Nano)
	}
	return clock.Now().UTC().Format(time.RFC3339Nano)
}
func storeRoot(workspaceRoot string) string {
	return filepath.Join(workspaceRoot, filepath.FromSlash(storeRelativeRoot))
}
func rootsFromLocator(locator string) (string, string, bool) {
	marker := string(filepath.Separator) + filepath.FromSlash(storeRelativeRoot) + string(filepath.Separator) + "runs" + string(filepath.Separator)
	index := strings.Index(locator, marker)
	if index < 0 {
		return "", "", false
	}
	workspaceRoot := locator[:index]
	return storeRoot(workspaceRoot), workspaceRoot, true
}
func publicationFilename(key string) string {
	sum := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x.ref", sum)
}
func relativeLocator(root, locator string) string {
	relative, _ := filepath.Rel(root, locator)
	return filepath.ToSlash(relative)
}
func resolveStoreLocator(root, relative string) (string, error) {
	if relative == "" || path.IsAbs(relative) || path.Clean(relative) != relative || relative == ".." || strings.HasPrefix(relative, "../") || filepath.ToSlash(filepath.FromSlash(relative)) != relative {
		return "", fmt.Errorf("locator must be a clean slash-separated relative path")
	}
	resolved := filepath.Join(root, filepath.FromSlash(relative))
	if resolved == root || !strings.HasPrefix(resolved, root+string(filepath.Separator)) {
		return "", fmt.Errorf("locator escapes store root")
	}
	return resolved, nil
}
func envelope(kind string, members ...canonicaljson.Member) canonicaljson.Object {
	result := canonicaljson.Object{{Name: "kind", Value: kind}, {Name: "schema_version", Value: int64(1)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"}}
	return append(result, members...)
}
func objectString(object canonicaljson.Object, name string) string {
	value, _ := objectMember(object, name)
	text, _ := value.(string)
	return text
}
func storeConflict(path, detail string, err error) error {
	return classified(ErrorWorkspaceConflict, fmt.Sprintf("%s: %s", path, detail), err)
}
func storeIO(operation, path string, err error) error {
	if err == nil {
		err = errors.New("unknown store error")
	}
	return classified(ErrorIO, fmt.Sprintf("%s %s: %v", operation, path, err), err)
}
func combineStoreErrors(primary, secondary error) error {
	if primary == nil {
		return secondary
	}
	if secondary == nil {
		return primary
	}
	return fmt.Errorf("%w; release: %v", primary, secondary)
}
