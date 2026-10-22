package workspace

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// AL38 red-before-code regression: external ownership crosses Epic/Task boundaries.
func TestExternalIssuesWorkspaceOwnershipAndPrimary(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	before, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	issue := ExternalIssue{Provider: ExternalJira, Key: "PLY-123", URL: "https://jira.example.invalid/browse/PLY-123"}
	in := ExternalIssueInput{SubjectKind: "epic", SubjectID: "epic", Operation: ExternalIssueAdd, Issue: &issue, ActorClaim: "test"}
	added, err := UpdateExternalIssues(f.dependencies, in)
	if err != nil || !added.Changed || len(added.Refs) != 1 || added.Primary != nil {
		t.Fatalf("add: %#v %v", added, err)
	}
	retry, err := UpdateExternalIssues(f.dependencies, in)
	if err != nil || retry.Changed || retry.Revision != 1 {
		t.Fatalf("retry: %#v %v", retry, err)
	}
	in.SubjectKind, in.SubjectID = "task", "task"
	if _, err := UpdateExternalIssues(f.dependencies, in); err == nil {
		t.Fatal("same issue linked to distinct Task")
	}
	rev := added.Revision
	selected, err := UpdateExternalIssues(f.dependencies, ExternalIssueInput{SubjectKind: "epic", SubjectID: "epic", Operation: ExternalIssuePrimary, RefID: added.Refs[0].ID, ExpectedRevision: &rev, ActorClaim: "test"})
	if err != nil || selected.Primary == nil {
		t.Fatalf("primary: %#v %v", selected, err)
	}
	task, err := ShowExternalIssues(f.dependencies, "task", "task")
	if err != nil || task.Effective == nil || task.Effective.CommitKey != "PLY-123" || task.Effective.SubjectKind != "epic" {
		t.Fatalf("fallback: %#v %v", task, err)
	}
	after, err := os.ReadFile(workItemsPath(root))
	if err != nil || string(before) != string(after) {
		t.Fatal("binding changed identity registry", err)
	}
}

func TestExternalIssuesJiraContextsAreDistinctIdentities(t *testing.T) {
	first, err := NormalizeExternalIssue(ExternalIssue{Provider: ExternalJira, Key: "PLY-123", URL: "https://example.test/jira-a/browse/PLY-123"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := NormalizeExternalIssue(ExternalIssue{Provider: ExternalJira, Key: "PLY-123", URL: "https://example.test/jira-b/browse/PLY-123"})
	if err != nil {
		t.Fatal(err)
	}
	if first.ID == second.ID {
		t.Fatal("distinct Jira base contexts collapsed to the same identity")
	}
	alias, err := NormalizeExternalIssue(ExternalIssue{Provider: ExternalJira, Key: "ply-123", URL: "https://EXAMPLE.test:443/jira-%61/browse/%50LY-123/"})
	if err != nil || alias != first {
		t.Fatalf("same instance alias did not normalize: %#v %v", alias, err)
	}
}

func TestExternalIssuesNormalizationRejectsAmbiguousMetadata(t *testing.T) {
	jira := ExternalIssue{Provider: ExternalJira, Key: "PLY-123", URL: "http://JIRA.example.invalid:80/jira/browse/ply-123/"}
	normalized, err := NormalizeExternalIssue(jira)
	if err != nil || normalized.URL != "http://jira.example.invalid/jira/browse/PLY-123" || normalized.JiraContextPath != "/jira" {
		t.Fatalf("HTTP Jira: %#v %v", normalized, err)
	}
	github := ExternalIssue{Provider: ExternalGitHub, Host: "git.example.invalid", Owner: "Owner", Repository: "Repo", Number: 123, URL: "https://GIT.example.invalid:443/OWNER/Repo/issues/00123/"}
	gh, err := NormalizeExternalIssue(github)
	if err != nil || gh.CommitKey() != "owner/repo#123" || gh.URL != "https://git.example.invalid/owner/repo/issues/123" {
		t.Fatalf("GitHub alias: %#v %v", gh, err)
	}
	for _, raw := range []string{
		"ftp://jira.example.invalid/browse/PLY-123", "https://user:secret@jira.example.invalid/browse/PLY-123", "https://jira.example.invalid/browse/PLY-123?x=y", "https://jira.example.invalid/browse/PLY-123?", "https://jira.example.invalid/browse/PLY-123#", "https://jira.example.invalid/browse/PLY-123#fragment", "https://jira.example.invalid/browse/OTHER-1", "https://jira.example.invalid/browse/PLY-123\n", "https://jira.example.invalid/browse/%0APLY-123", "https://jira.example.invalid/a%2fb/browse/PLY-123", "https://jira.example.invalid/a%5cb/browse/PLY-123", "https://jira.example.invalid/../browse/PLY-123", "https://jira.example.invalid//browse/PLY-123", "https://jira.example.invalid/unsupported/PLY-123", "https://jira.example.invalid:/browse/PLY-123", "https://jira.example.invalid:65536/browse/PLY-123", "https://jira.example.invalid/browse/PLY-123//", "https://[::1]/browse/PLY-123", "https://jirá.invalid/browse/PLY-123", "https://jira.example.invalid/browse/%2550LY-123",
	} {
		t.Run(raw, func(t *testing.T) {
			in := jira
			in.URL = raw
			if _, err := NormalizeExternalIssue(in); err == nil {
				t.Fatalf("accepted %q", raw)
			}
		})
	}
	bad := []ExternalIssue{jira, jira, jira, github, github, github, github, github, github, github, github}
	bad[0].Provider = "unknown"
	bad[1].Key = "PLY-0"
	bad[2].Number = 1
	bad[3].Number = 0
	bad[4].Number = -1
	bad[5].Number = 124
	bad[6].Host = "other.invalid"
	bad[7].Owner = "other"
	bad[8].Repository = "other"
	bad[9].URL = "https://git.example.invalid/Owner/Repo/pull/123"
	bad[10].Host = ""
	for i, in := range bad {
		if _, err := NormalizeExternalIssue(in); err == nil {
			t.Fatalf("invalid identity %d accepted", i)
		}
	}
}

func TestExternalIssuesPrimaryCASUnlinkAndHistory(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	epicIssue := ExternalIssue{Provider: ExternalJira, Key: "PLY-1", URL: "https://jira.example.invalid/browse/PLY-1"}
	gh := ExternalIssue{Provider: ExternalGitHub, Host: "git.example.invalid", Owner: "Owner", Repository: "Repo", Number: 2, URL: "https://git.example.invalid/owner/repo/issues/2"}
	update := func(kind, id, op, ref string, issue *ExternalIssue, primary bool, rev *int) ExternalIssuesMutation {
		t.Helper()
		out, err := UpdateExternalIssues(f.dependencies, ExternalIssueInput{SubjectKind: kind, SubjectID: id, Operation: op, RefID: ref, Issue: issue, Primary: primary, ExpectedRevision: rev, ActorClaim: "test"})
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	zero := 0
	epic := update("epic", "epic", ExternalIssueAdd, "", &epicIssue, true, &zero)
	rev := epic.Revision
	task := update("task", "task", ExternalIssueAdd, "", &gh, true, &rev)
	if task.Effective == nil || task.Effective.CommitKey != "owner/repo#2" || task.Effective.SubjectKind != "task" {
		t.Fatalf("Task precedence: %#v", task)
	}
	// Stale writes fail without publication, including stale add-with-primary.
	stale := ExternalIssueInput{SubjectKind: "task", SubjectID: "task", Operation: ExternalIssuePrimary, ExpectedRevision: &rev, ActorClaim: "test"}
	path := task.Source.Locator
	before, _ := os.ReadFile(path)
	if _, err := UpdateExternalIssues(f.dependencies, stale); err == nil {
		t.Fatal("stale primary accepted")
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(before) {
		t.Fatal("stale primary published")
	}
	rev = task.Revision
	cleared := update("task", "task", ExternalIssuePrimary, "", nil, false, &rev)
	if cleared.Primary != nil || cleared.Effective == nil || cleared.Effective.SubjectKind != "epic" {
		t.Fatalf("clear fallback: %#v", cleared)
	}
	rev = cleared.Revision
	selected := update("task", "task", ExternalIssuePrimary, task.Refs[0].ID, nil, false, &rev)
	rev = selected.Revision
	unlinked := update("task", "task", ExternalIssueUnlink, task.Refs[0].ID, nil, false, &rev)
	if unlinked.Primary != nil || len(unlinked.Refs) != 0 || unlinked.Effective == nil || unlinked.Effective.RefID != epic.Refs[0].ID || len(unlinked.Events) != 4 {
		t.Fatalf("unlink/history: %#v", unlinked)
	}
	rev = unlinked.Revision
	epic = update("epic", "epic", ExternalIssueUnlink, epic.Refs[0].ID, nil, false, &rev)
	read, err := ShowExternalIssues(f.dependencies, "task", "task")
	if err != nil || read.Effective != nil {
		t.Fatalf("cleared Epic fallback: %#v %v", read, err)
	}
	// Explicit unlink releases workspace identity without erasing history.
	moved := update("task", "task", ExternalIssueAdd, "", &epicIssue, false, nil)
	if moved.Primary != nil || len(moved.Refs) != 1 {
		t.Fatalf("moved: %#v", moved)
	}
	// The same normalized identity with a different canonical scheme conflicts.
	changed := epicIssue
	changed.URL = "http://jira.example.invalid/browse/PLY-1"
	before, _ = os.ReadFile(path)
	if _, err := UpdateExternalIssues(f.dependencies, ExternalIssueInput{SubjectKind: "task", SubjectID: "task", Operation: ExternalIssueAdd, Issue: &changed, ActorClaim: "test"}); err == nil {
		t.Fatal("conflicting URL accepted")
	}
	after, _ = os.ReadFile(path)
	if string(after) != string(before) {
		t.Fatal("conflict published")
	}
}

func TestExternalIssuesLegacyAndUnavailableSourceDoNotChangeOtherMetadata(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		t.Run(fmt.Sprint(legacy), func(t *testing.T) {
			f := newWorkItemJourneyFixtureVersion(t, legacy)
			root := filepath.Dir(f.wrapper)
			// Worktree/source observation adapters are deliberately unavailable.
			f.dependencies.WorkGit = nil
			// Dirty adopted worktree is irrelevant to metadata mutation and readback.
			if err := os.WriteFile(filepath.Join(f.wrapper, "epic", "dirty.txt"), []byte("live work"), 0600); err != nil {
				t.Fatal(err)
			}
			files := map[string][]byte{}
			err := filepath.WalkDir(filepath.Join(root, MarkerDirectory), func(path string, entry fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if entry.Type().IsRegular() {
					b, e := os.ReadFile(path)
					if e != nil {
						return e
					}
					files[path] = b
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			refsBefore := runLocalGit(t, filepath.Join(f.wrapper, "main"), "show-ref")
			issue := ExternalIssue{Provider: ExternalJira, Key: "PLY-7", URL: "http://offline.invalid/context/browse/PLY-7"}
			out, err := UpdateExternalIssues(f.dependencies, ExternalIssueInput{SubjectKind: "task", SubjectID: "task", Operation: ExternalIssueAdd, Issue: &issue, ActorClaim: "test"})
			if err != nil || len(out.Refs) != 1 {
				t.Fatalf("legacy metadata %#v %v", out, err)
			}
			for path, before := range files {
				after, err := os.ReadFile(path)
				if err != nil || string(before) != string(after) {
					t.Fatalf("changed unrelated metadata %s: %v", path, err)
				}
			}
			if after := runLocalGit(t, filepath.Join(f.wrapper, "main"), "show-ref"); after != refsBefore {
				t.Fatal("Git refs changed")
			}
			if err := os.Rename(filepath.Join(f.wrapper, "epic"), filepath.Join(f.wrapper, "epic-unavailable")); err != nil {
				t.Fatal(err)
			}
			read, err := ShowExternalIssues(f.dependencies, "task", "task")
			if err != nil || len(read.Refs) != 1 || read.ProviderStatus != "offline_unverified" {
				t.Fatalf("offline read: %#v %v", read, err)
			}
		})
	}
}

type externalAttemptStore struct {
	WorkItemStore
	attempted chan struct{}
}

func (s externalAttemptStore) WithLock(root string, fn func(WorkItemStoreSession) error) error {
	close(s.attempted)
	return s.WorkItemStore.WithLock(root, fn)
}
func TestExternalIssuesConcurrentCrossKindPublishesOneWinner(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	entered, release, attempted := make(chan struct{}), make(chan struct{}), make(chan struct{})
	results := make(chan error, 2)
	issue := ExternalIssue{Provider: ExternalJira, Key: "PLY-8", URL: "https://offline.invalid/browse/PLY-8"}
	in := ExternalIssueInput{SubjectKind: "epic", SubjectID: "epic", Operation: ExternalIssueAdd, Issue: &issue, ActorClaim: "test"}
	go func() {
		_, err := updateExternalIssues(f.dependencies, in, workspaceMetadataPublisher{fault: func(stage string) error {
			if stage == "replace" {
				close(entered)
				select {
				case <-release:
				case <-time.After(5 * time.Second):
					return fmt.Errorf("publication barrier timeout")
				}
			}
			return nil
		}})
		results <- err
	}()
	select {
	case <-entered:
	case err := <-results:
		t.Fatalf("writer failed before barrier: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("writer did not reach barrier")
	}
	other := f.dependencies
	other.WorkItems = externalAttemptStore{other.WorkItems, attempted}
	in.SubjectKind, in.SubjectID = "task", "task"
	go func() { _, err := UpdateExternalIssues(other, in); results <- err }()
	select {
	case <-attempted:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("contender did not reach lock")
	}
	close(release)
	successes, conflicts := 0, 0
	for i := 0; i < 2; i++ {
		select {
		case err := <-results:
			if err == nil {
				successes++
			} else if hasWorkErrorClass(err, ErrorWorkIdentityConflict) {
				conflicts++
			} else {
				t.Fatal(err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("writer did not return")
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("wins=%d conflicts=%d", successes, conflicts)
	}
	s, err := readExternalIssues(filepath.Dir(f.wrapper))
	if err != nil || len(s.store.Events) != 1 || len(s.owners) != 1 {
		t.Fatalf("publication %#v %v", s, err)
	}
}

func TestExternalIssuesPrimaryAddRetryIsNoEffectButCannotRestoreStalePrimary(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	zero := 0
	issue := ExternalIssue{Provider: ExternalJira, Key: "PLY-9", URL: "https://jira.invalid/browse/PLY-9"}
	in := ExternalIssueInput{SubjectKind: "task", SubjectID: "task", Operation: ExternalIssueAdd, Issue: &issue, Primary: true, ExpectedRevision: &zero, ActorClaim: "test"}
	first, err := UpdateExternalIssues(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(first.Source.Locator)
	if err != nil {
		t.Fatal(err)
	}
	retry, err := UpdateExternalIssues(f.dependencies, in)
	if err != nil || retry.Changed || retry.Revision != 1 {
		t.Fatalf("exact add-primary retry refused: %#v %v", retry, err)
	}
	after, _ := os.ReadFile(first.Source.Locator)
	if string(before) != string(after) {
		t.Fatal("retry rewrote bytes")
	}
	rev := first.Revision
	clear := ExternalIssueInput{SubjectKind: "task", SubjectID: "task", Operation: ExternalIssuePrimary, ExpectedRevision: &rev, ActorClaim: "test"}
	cleared, err := UpdateExternalIssues(f.dependencies, clear)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := UpdateExternalIssues(f.dependencies, in); err == nil {
		t.Fatal("stale add restored deliberately cleared primary")
	}
	repeated, err := UpdateExternalIssues(f.dependencies, clear)
	if err != nil || repeated.Changed || repeated.Revision != cleared.Revision {
		t.Fatalf("no-effect clear retry %#v %v", repeated, err)
	}
}

func TestExternalIssuesCrossProjectUniquenessAndDistinctJiraContexts(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	if _, _, _, err := f.dependencies.Projects.Add(root, ProjectRecord{ID: "other", Name: "Other", Wrapper: "/missing/other", RepoIDs: []RepoID{"other"}}, []RepoRecord{{ID: "other", Locator: "/missing/other/main", GitCommonDir: "/missing/other/main/.git"}}); err != nil {
		t.Fatal(err)
	}
	err := f.dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		e := r.Epics[0]
		e.ID = "other"
		e.ProjectID = "other"
		e.RepoBindings = append([]EpicRepoBinding{}, e.RepoBindings...)
		b := &e.RepoBindings[0]
		b.RepoID = "other"
		b.GitCommonDir = "/missing/other/main/.git"
		b.Worktree.ID = "wt_11111111111111111111111111111111"
		b.Worktree.OwnerID = "other"
		b.Worktree.Locator = "/missing/other/epic"
		b.Worktree.GitCommonDir = b.GitCommonDir
		r.Epics = append(r.Epics, e)
		return session.Publish(r)
	})
	if err != nil {
		t.Fatal(err)
	}
	issue := ExternalIssue{Provider: ExternalJira, Key: "PLY-123", URL: "https://example.test/jira-a/browse/PLY-123"}
	in := ExternalIssueInput{SubjectKind: "task", SubjectID: "task", Operation: ExternalIssueAdd, Issue: &issue, ActorClaim: "test"}
	first, err := UpdateExternalIssues(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	in.SubjectKind, in.SubjectID = "epic", "other"
	issue.URL = "https://EXAMPLE.test:443/jira-%61/browse/PLY-123/"
	if _, err := UpdateExternalIssues(f.dependencies, in); !hasWorkErrorClass(err, ErrorWorkIdentityConflict) {
		t.Fatalf("cross-project duplicate not refused: %v", err)
	}
	issue.URL = "https://example.test/jira-b/browse/PLY-123"
	other, err := UpdateExternalIssues(f.dependencies, in)
	if err != nil || len(other.Refs) != 1 || other.Refs[0].ID == first.Refs[0].ID {
		t.Fatalf("distinct instance: %#v %v", other, err)
	}
}

func TestExternalIssuesStrictStoreAndPublicationReadback(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	issue := ExternalIssue{Provider: ExternalJira, Key: "PLY-5", URL: "https://offline.invalid/browse/PLY-5"}
	in := ExternalIssueInput{SubjectKind: "task", SubjectID: "task", Operation: ExternalIssueAdd, Issue: &issue, ActorClaim: "test"}
	out, err := updateExternalIssues(f.dependencies, in, workspaceMetadataPublisher{fault: func(stage string) error {
		if stage == "directory-sync" {
			return fmt.Errorf("injected post-replace uncertainty")
		}
		return nil
	}})
	if err == nil || out.Changed {
		t.Fatalf("unknown durability presented as success: %#v %v", out, err)
	}
	read, err := ShowExternalIssues(f.dependencies, "task", "task")
	if err != nil || read.Revision != 1 || len(read.Refs) != 1 {
		t.Fatalf("publication observation: %#v %v", read, err)
	}
	retry, err := UpdateExternalIssues(f.dependencies, in)
	if err != nil || retry.Changed || retry.Revision != 1 {
		t.Fatalf("identical recovery %#v %v", retry, err)
	}
	good, err := os.ReadFile(read.Source.Locator)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{
		strings.Replace(string(good), `"schema_version": 1`, `"schema_version": 2`, 1),
		strings.Replace(string(good), `"schema_version": 1`, `"schema_version": 1, "extra": true`, 1),
		strings.Replace(string(good), `"schema_version": 1`, `"schema_version": 1, "schema_version": 1`, 1),
		strings.Replace(string(good), `"sequence": 1`, `"sequence": 2`, 1),
		strings.Replace(string(good), `"provider": "jira"`, `"provider": "other"`, 1),
	} {
		if err := os.WriteFile(read.Source.Locator, []byte(bad), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := readExternalIssues(root); err == nil {
			t.Fatal("corrupt metadata accepted")
		}
	}
	if err := os.WriteFile(read.Source.Locator, good, 0600); err != nil {
		t.Fatal(err)
	}
	task, err := ShowTask(f.dependencies, "task")
	if err != nil || task.ExternalIssues == nil {
		t.Fatalf("Task projection: %#v %v", task.ExternalIssues, err)
	}
	b, err := MarshalTaskReadback(task)
	if err != nil || !strings.Contains(string(b), `"external_issues"`) || !strings.Contains(string(b), `"offline_unverified"`) {
		t.Fatalf("Task JSON: %s %v", b, err)
	}
}

func TestExternalIssuesPreserveQueuedGoalAndHistoricalSpec(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	firstGoal := goalFixture(t, f, "task", "external-issue-fixture-goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &firstGoal}})
	goalFixture(t, f, "task", "external-issue-fixture-goal-revised")
	before, err := os.ReadFile(workItemsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.TaskSpecRevisions) != 2 {
		t.Fatalf("expected historical and current Specs, got %d", len(r.TaskSpecRevisions))
	}
	issue := ExternalIssue{Provider: ExternalGitHub, Host: "git.offline.invalid", Owner: "fixture", Repository: "repo", Number: 42, URL: "https://git.offline.invalid/fixture/repo/issues/42"}
	out, err := UpdateExternalIssues(f.dependencies, ExternalIssueInput{SubjectKind: "task", SubjectID: "task", Operation: ExternalIssueAdd, Issue: &issue, ActorClaim: "test"})
	if err != nil || !out.Changed {
		t.Fatalf("queued goal metadata %#v %v", out, err)
	}
	after, err := os.ReadFile(workItemsPath(root))
	if err != nil || string(before) != string(after) {
		t.Fatalf("metadata changed queued goal or historical Spec bindings: %v", err)
	}
	if _, err := f.dependencies.WorkItems.Snapshot(root); err != nil {
		t.Fatalf("immutable content closure changed: %v", err)
	}
}
