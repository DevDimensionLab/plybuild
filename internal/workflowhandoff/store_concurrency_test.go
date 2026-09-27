package workflowhandoff

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func TestConcurrentIdenticalCreateAndStart(t *testing.T) {
	root := physicalTempDir(t)
	store := newSystemStore(systemFileSystem{}, systemClock{})
	input := CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)}

	start := make(chan struct{})
	results := make(chan CreateStoreResult, 2)
	errors := make(chan error, 2)
	var group sync.WaitGroup
	for index := 0; index < 2; index++ {
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			result, err := store.Create(input)
			results <- result
			errors <- err
		}()
	}
	close(start)
	group.Wait()
	close(results)
	close(errors)
	created := 0
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	var snapshot Snapshot
	for result := range results {
		if result.Created {
			created++
		}
		snapshot = result.Snapshot
	}
	if created != 1 {
		t.Fatalf("created count = %d", created)
	}

	document := testAcceptedDocument(t, "rcp_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "started")
	start = make(chan struct{})
	submitResults := make(chan SubmitStoreResult, 2)
	errors = make(chan error, 2)
	group = sync.WaitGroup{}
	for index := 0; index < 2; index++ {
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			result, err := store.SubmitStart(SubmitStoreInput{Snapshot: snapshot, Phase: "start", Document: document})
			submitResults <- result
			errors <- err
		}()
	}
	close(start)
	group.Wait()
	close(submitResults)
	close(errors)
	created = 0
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	for result := range submitResults {
		if result.Created {
			created++
		}
	}
	if created != 1 {
		t.Fatalf("accepted start count = %d", created)
	}
}

func TestConcurrentSupersedeAndSubmitStartHaveOneWinner(t *testing.T) {
	dependencies, created, snapshot, principal := createTestHandoff(t)
	startPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "race-start.json", minimalStartDraftValue(t, snapshot, principal))
	replacement := replaceObjectMember(minimalHandoffDraftValue(snapshot.Handoff.Target.Worktree, snapshot.Handoff.Target.Ref, snapshot.Handoff.Target.OID), "publication_key", "wf01/race-replacement")
	replacementPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "race-replacement.json", replacement)

	start := make(chan struct{})
	errors := make(chan error, 2)
	var group sync.WaitGroup
	group.Add(2)
	go func() {
		defer group.Done()
		<-start
		_, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
		errors <- err
	}()
	go func() {
		defer group.Done()
		<-start
		_, err := Supersede(dependencies, SupersedeInput{HandoffID: created.HandoffID, DraftPath: replacementPath, Reason: "Correct the handoff before start."})
		errors <- err
	}()
	close(start)
	group.Wait()
	close(errors)
	succeeded := 0
	for err := range errors {
		if err == nil {
			succeeded++
		} else if !IsClass(err, ErrorStale) {
			t.Fatalf("loser error = %v", err)
		}
	}
	if succeeded != 1 {
		t.Fatalf("successful transitions = %d", succeeded)
	}

	old, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if (old.Start != nil) == (old.HeadState == "superseded") {
		t.Fatalf("old snapshot must be either started or superseded: %#v", old)
	}
}

func TestConcurrentConflictingStartWritersKeepOneAcceptedDocument(t *testing.T) {
	dependencies, created, snapshot, principal := createTestHandoff(t)
	first := minimalStartDraftValue(t, snapshot, principal)
	otherPrincipal := append(canonicaljson.Object{}, principal...)
	otherPrincipal = replaceObjectMember(otherPrincipal, "session_id", "competing-session")
	second := replaceObjectMember(first, "principal", otherPrincipal)
	firstPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "first-start.json", first)
	secondPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "second-start.json", second)
	start := make(chan struct{})
	results := make(chan SubmitResult, 2)
	errors := make(chan error, 2)
	var group sync.WaitGroup
	for _, path := range []string{firstPath, secondPath} {
		path := path
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			result, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: path})
			results <- result
			errors <- err
		}()
	}
	close(start)
	group.Wait()
	close(results)
	close(errors)
	succeeded := 0
	for err := range errors {
		if err == nil {
			succeeded++
		} else if !IsClass(err, ErrorConflict) {
			t.Fatalf("loser error = %v", err)
		}
	}
	if succeeded != 1 {
		t.Fatalf("successful starts = %d", succeeded)
	}
	read, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if read.Start == nil {
		t.Fatal("accepted start is missing")
	}
	entries, err := os.ReadDir(filepath.Join(read.Handoff.ReplyRoot, "start", "events"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("start attempt events = %d", len(entries))
	}
}
