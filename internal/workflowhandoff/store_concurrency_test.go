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

func TestConcurrentDifferentCreatesWithSamePublicationKeyHaveOneAuthority(t *testing.T) {
	root := physicalTempDir(t)
	store := newSystemStore(systemFileSystem{}, systemClock{})
	first := testStoreHandoff(t, root)
	second := first
	second.Identity.RunID = "run_cccccccccccccccccccccccccccccccc"
	second.Identity.HandoffID = "hnd_dddddddddddddddddddddddddddddddd"
	second.Identity.StartReceiptID = "rcp_eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	second.Identity.TerminalResultID = "res_ffffffffffffffffffffffffffffffff"
	second.SourceDraftSHA256 = digestBytes([]byte("different draft"))
	second.Value = replaceObjectMember(second.Value, "source_draft_sha256", second.SourceDraftSHA256)
	identityValue, _ := objectMember(second.Value, "identity")
	identity := identityValue.(canonicaljson.Object)
	identity = replaceObjectMember(identity, "run_id", second.Identity.RunID)
	identity = replaceObjectMember(identity, "handoff_id", second.Identity.HandoffID)
	identity = replaceObjectMember(identity, "start_receipt_id", second.Identity.StartReceiptID)
	identity = replaceObjectMember(identity, "terminal_result_id", second.Identity.TerminalResultID)
	second.Value = replaceObjectMember(second.Value, "identity", identity)
	second.Bytes, _ = canonicaljson.Marshal(second.Value)
	second.SHA256 = digestBytes(second.Bytes)
	inputs := []CreateStoreInput{
		{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: first.SourceDraftSHA256, Handoff: first},
		{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: second.SourceDraftSHA256, Handoff: second},
	}
	start := make(chan struct{})
	errors := make(chan error, 2)
	var group sync.WaitGroup
	for _, input := range inputs {
		input := input
		group.Add(1)
		go func() {
			defer group.Done()
			<-start
			_, err := store.Create(input)
			errors <- err
		}()
	}
	close(start)
	group.Wait()
	close(errors)
	succeeded, conflicted := 0, 0
	for err := range errors {
		switch {
		case err == nil:
			succeeded++
		case IsClass(err, ErrorConflict):
			conflicted++
		default:
			t.Fatalf("create error = %v", err)
		}
	}
	if succeeded != 1 || conflicted != 1 {
		t.Fatalf("succeeded=%d conflicted=%d", succeeded, conflicted)
	}
	entries, err := os.ReadDir(filepath.Join(storeRoot(root), "runs"))
	if err != nil {
		t.Fatal(err)
	}
	runs := 0
	for _, entry := range entries {
		if entry.IsDir() && entry.Name()[0] != '.' {
			runs++
		}
	}
	if runs != 1 {
		t.Fatalf("published runs = %d", runs)
	}
}

func TestConcurrentTerminalWritersAreIdempotentOrConflict(t *testing.T) {
	for _, different := range []bool{false, true} {
		name := "identical"
		if different {
			name = "different"
		}
		t.Run(name, func(t *testing.T) {
			dependencies, created, snapshot, principal := createTestHandoff(t)
			startPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "terminal-race-start.json", minimalStartDraftValue(t, snapshot, principal))
			startResult, err := SubmitStart(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: startPath})
			if err != nil {
				t.Fatal(err)
			}
			first := minimalTerminalDraftValue(snapshot, principal, startResult)
			second := first
			if different {
				second = replaceObjectMember(first, "summary", "A competing bounded result.")
			}
			paths := []string{
				writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "terminal-race-first.json", first),
				writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "terminal-race-second.json", second),
			}
			start := make(chan struct{})
			results := make(chan SubmitResult, 2)
			errors := make(chan error, 2)
			var group sync.WaitGroup
			for _, path := range paths {
				path := path
				group.Add(1)
				go func() {
					defer group.Done()
					<-start
					result, err := SubmitResultDocument(dependencies, SubmitInput{HandoffLocator: created.Locator, DraftPath: path})
					results <- result
					errors <- err
				}()
			}
			close(start)
			group.Wait()
			close(results)
			close(errors)
			succeeded, conflicted := 0, 0
			for err := range errors {
				if err == nil {
					succeeded++
				} else if IsClass(err, ErrorConflict) {
					conflicted++
				} else {
					t.Fatalf("terminal writer error = %v", err)
				}
			}
			if succeeded != 2 || (!different && conflicted != 0) {
				if different && succeeded == 1 && conflicted == 1 {
					// Expected competing payload outcome.
				} else {
					t.Fatalf("succeeded=%d conflicted=%d", succeeded, conflicted)
				}
			}
			createdCount := 0
			for result := range results {
				if result.Created {
					createdCount++
				}
			}
			if createdCount != 1 {
				t.Fatalf("created terminal results = %d", createdCount)
			}
			read, err := dependencies.Store.ReadByLocator(created.Locator)
			if err != nil || read.Terminal == nil {
				t.Fatalf("terminal snapshot = %#v, %v", read, err)
			}
			entries, err := os.ReadDir(filepath.Join(read.Handoff.ReplyRoot, "terminal", "events"))
			if err != nil || len(entries) != 2 {
				t.Fatalf("terminal attempts = %d, %v", len(entries), err)
			}
		})
	}
}

func TestConcurrentCancelAndSubmitStartHaveOneWinner(t *testing.T) {
	dependencies, created, snapshot, principal := createTestHandoff(t)
	startPath := writeCanonicalTestFile(t, snapshot.WorkspaceRoot, "cancel-race-start.json", minimalStartDraftValue(t, snapshot, principal))
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
		_, err := Cancel(dependencies, ControlInput{HandoffID: created.HandoffID, Reason: "Cancel concurrently before start."})
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
	read, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if !((read.Start != nil && read.HeadState == "ready") || (read.Start == nil && read.HeadState == "cancelled")) {
		t.Fatalf("incoherent cancel/start snapshot = %#v", read)
	}
}

func TestReadbackRemainsCoherentWhenHeadOrAcceptedRefsChangeMidRead(t *testing.T) {
	t.Run("head changes to cancelled after old head bytes", func(t *testing.T) {
		root := physicalTempDir(t)
		base := newSystemStore(systemFileSystem{}, systemClock{})
		created, err := base.Create(CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)})
		if err != nil {
			t.Fatal(err)
		}
		reader := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
		original := reader.ops.readFile
		triggered := false
		reader.ops.readFile = func(path string) ([]byte, error) {
			contents, err := original(path)
			if err == nil && filepath.Base(path) == "head.ref" && !triggered {
				triggered = true
				if _, cancelErr := base.Cancel(ControlStoreInput{WorkspaceRoot: root, HandoffID: HandoffID(created.Snapshot.Handoff.Identity.HandoffID), Reason: "Concurrent cancellation."}); cancelErr != nil {
					t.Fatalf("cancel during read: %v", cancelErr)
				}
			}
			return contents, err
		}
		observed, err := reader.ReadByLocator(created.Snapshot.Handoff.Locator)
		if err != nil {
			t.Fatal(err)
		}
		if observed.HeadState != "ready" || observed.Start != nil || observed.Terminal != nil {
			t.Fatalf("mixed old snapshot = %#v", observed)
		}
		current, err := base.ReadByLocator(created.Snapshot.Handoff.Locator)
		if err != nil || current.HeadState != "cancelled" {
			t.Fatalf("current snapshot = %#v, %v", current, err)
		}
	})

	t.Run("terminal ref appears after start ref bytes", func(t *testing.T) {
		root := physicalTempDir(t)
		base := newSystemStore(systemFileSystem{}, systemClock{})
		created, err := base.Create(CreateStoreInput{WorkspaceRoot: root, PublicationKey: "wf01/test", SourceDraftSHA256: digestBytes([]byte("draft")), Handoff: testStoreHandoff(t, root)})
		if err != nil {
			t.Fatal(err)
		}
		start := testAcceptedDocument(t, created.Snapshot.Handoff.Identity.StartReceiptID, "started")
		if _, err := base.SubmitStart(SubmitStoreInput{Snapshot: created.Snapshot, Phase: "start", Document: start}); err != nil {
			t.Fatal(err)
		}
		started, err := base.ReadByLocator(created.Snapshot.Handoff.Locator)
		if err != nil {
			t.Fatal(err)
		}
		terminal := testAcceptedDocument(t, started.Handoff.Identity.TerminalResultID, "complete")
		reader := newSystemStore(systemFileSystem{}, systemClock{}).(*systemStore)
		original := reader.ops.readFile
		triggered := false
		reader.ops.readFile = func(path string) ([]byte, error) {
			contents, err := original(path)
			if err == nil && filepath.Base(path) == "accepted.ref" && filepath.Base(filepath.Dir(path)) == "start" && !triggered {
				triggered = true
				if _, submitErr := base.SubmitResult(SubmitStoreInput{Snapshot: started, Phase: "terminal", Document: terminal}); submitErr != nil {
					t.Fatalf("terminal during read: %v", submitErr)
				}
			}
			return contents, err
		}
		observed, err := reader.ReadByLocator(created.Snapshot.Handoff.Locator)
		if err != nil {
			t.Fatal(err)
		}
		if observed.Start == nil || observed.Terminal == nil || observed.Start.DocumentID != start.DocumentID || observed.Terminal.DocumentID != terminal.DocumentID {
			t.Fatalf("mixed terminal snapshot = %#v", observed)
		}
	})
}
