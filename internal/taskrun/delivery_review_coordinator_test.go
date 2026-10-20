package taskrun

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func reviewRoundRequestFile(t *testing.T, round DeliveryReviewRoundStatus) deliveryReviewRound {
	t.Helper()
	var r deliveryReviewRound
	b, err := os.ReadFile(round.Request.Locator)
	if err == nil {
		err = json.Unmarshal(b, &r)
	}
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// RV01 AC01: an inherited HERDR_PANE_ID can name a pane that no longer holds
// the coordinator. Ply binds the pane whose shell actually holds this process,
// never the focused pane, and the owner cannot hide behind a stale identity.
func TestRV01ReviewBindsActualCoordinatorDespiteStaleInheritedPane(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	oid, _ := reviewRoundCandidate(t, f, "fixed attempt")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w1:p1")
	workflowTestModel(t, f.workflowFixture, map[string]any{"pane_shells": map[string]any{"w-fixture:root": os.Getppid(), "w-fixture:user": 99999}, "focused_pane": "w-fixture:user"})
	requested, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
	if err != nil {
		t.Fatal(err)
	}
	round := reviewRoundLatest(t, requested)
	if r := reviewRoundRequestFile(t, round); r.CoordinatorPane != "w-fixture:root" || round.CandidateOID != oid || round.ReviewerStart != "started" {
		t.Fatalf("review did not bind the actual coordinator pane: %+v %+v", r, round)
	}
	reviewRoundIn(t, round.ReviewerPane)
	workflowTestModel(t, f.workflowFixture, map[string]any{"pane_shells": nil})
	passed := reviewRoundReturnFile(t, f, "passed.json", `{"decision":"passed","summary":"Looks right","findings":[],"evidence_gaps":[]}`)
	if o, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, passed); err != nil {
		t.Fatal(err)
	}
	if r := reviewRoundLatest(t, o); r.Decision != "passed" || o.DeliveryStatus.Review.CurrentCandidate != "review_passed" {
		t.Fatalf("the reviewer's return did not complete the actual coordinator's round: %+v", o.DeliveryStatus.Review)
	}

	// The owner's own pane holding this process is refused despite a stale
	// inherited identity that names some other pane.
	reviewRoundCandidate(t, f, "fixed attempt, owner edit")
	reviewRoundIn(t, "w1:p1")
	workflowTestModel(t, f.workflowFixture, map[string]any{"pane_shells": map[string]any{o.Transport.PaneID: os.Getppid()}})
	if _, err = WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"}); err == nil || !strings.Contains(err.Error(), "cannot coordinate its own review") {
		t.Fatalf("the owner coordinated its own review through a stale pane identity: %v", err)
	}
}

// RV01 AC02: an unknown or ambiguous coordinator stops before any reviewer
// effect or round file, with a concrete next step. A focused user pane or an
// existing pane named by a stale identity is never bound.
func TestRV01UnknownOrAmbiguousCoordinatorStopsBeforeReviewerLaunch(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed attempt")
	coordinator := reviewRoundAsCoordinator(f)
	cases := []struct {
		name, env string
		shells    map[string]any
		want      string
	}{
		{"stale identity, focused user pane", "w1:p1", map[string]any{"w-fixture:user": 99999}, "unknown"},
		{"existing pane not holding the caller", "w-fixture:user", map[string]any{"w-fixture:user": 99999}, "unknown"},
		{"outside Herdr", "", map[string]any{"w-fixture:user": 99999}, "unknown"},
		{"two panes hold the caller", "w1:p1", map[string]any{"w-fixture:a": os.Getppid(), "w-fixture:b": os.Getppid()}, "ambiguous"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reviewRoundIn(t, c.env)
			workflowTestModel(t, f.workflowFixture, map[string]any{"pane_shells": c.shells, "focused_pane": "w-fixture:user"})
			_, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
			if err == nil || !strings.Contains(err.Error(), c.want) || !strings.Contains(err.Error(), "from the coordinator's own Herdr pane") {
				t.Fatalf("coordinator identity was not refused with a next step: %v", err)
			}
			if n := len(reviewRoundArgv(t, f, "tab create")); n != 1 {
				t.Fatalf("a reviewer tab was created for an unknown coordinator (%d tab creates)", n)
			}
			if _, err = os.Stat(filepath.Join(o.Paths.RunRoot, "review")); !os.IsNotExist(err) {
				t.Fatalf("a review round was written for an unknown coordinator: %v", err)
			}
		})
	}
}

// RV01 AC03: the reviewer runs with process-local readonly Go modules on top
// of the ambient flags, and its brief demands a clean exact candidate after
// tests. A dirty worktree cannot return passed; the actual negative judgment
// is sealed and the drift stays for inspection, never reset.
func TestRV01ReviewerKeepsCandidateAndCannotPassADirtyWorktree(t *testing.T) {
	t.Setenv("GOFLAGS", "-mod=mod -tags=fixture")
	f, o, _ := reviewRoundFixture(t)
	oid, _ := reviewRoundCandidate(t, f, "fixed attempt")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator"})
	if err != nil {
		t.Fatal(err)
	}
	tabs := reviewRoundArgv(t, f, "tab create")
	if tab := "\x00" + strings.Join(tabs[len(tabs)-1], "\x00") + "\x00"; !strings.Contains(tab, "\x00--env\x00GOFLAGS=-mod=mod -tags=fixture -mod=readonly\x00") {
		t.Fatalf("reviewer tab lacks process-local readonly Go modules: %q", tabs[len(tabs)-1])
	}
	round := reviewRoundLatest(t, requested)
	brief, err := os.ReadFile(round.Brief.Locator)
	if err != nil {
		t.Fatal(err)
	}
	target := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	for _, want := range []string{"-mod=readonly", "git -C " + ShellQuote(target) + " status --porcelain", "rev-parse HEAD", oid, "Never reset"} {
		if !strings.Contains(string(brief), want) {
			t.Fatalf("reviewer brief lacks %q:\n%s", want, brief)
		}
	}

	if err = os.WriteFile(filepath.Join(target, "go.mod"), []byte("module drifted\n"), 0600); err != nil {
		t.Fatal(err)
	}
	reviewRoundIn(t, round.ReviewerPane)
	passed := reviewRoundReturnFile(t, f, "passed.json", `{"decision":"passed","summary":"Tests pass","findings":[],"evidence_gaps":[]}`)
	if _, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, passed); err == nil || !strings.Contains(err.Error(), "go.mod") || !strings.Contains(err.Error(), "blocked or changes_requested") {
		t.Fatalf("a dirty Task worktree returned passed: %v", err)
	}
	blocked := reviewRoundReturnFile(t, f, "blocked.json", `{"decision":"blocked","summary":"Tests changed go.mod","findings":[],"evidence_gaps":["go.mod drifted from the candidate during tests: ?? go.mod"]}`)
	if o, err = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, blocked); err != nil {
		t.Fatal(err)
	}
	if r := reviewRoundLatest(t, o); r.Decision != "blocked" || len(r.EvidenceGaps) != 1 {
		t.Fatalf("the actual negative review was not preserved: %+v", r)
	}
	if b, err := os.ReadFile(filepath.Join(target, "go.mod")); err != nil || string(b) != "module drifted\n" {
		t.Fatalf("observed drift was not preserved: %q %v", b, err)
	}
}

// RV01 AC04: an actual go test in a disposable module, whose ambient flags
// select -mod=mod, keeps go.mod unchanged under the reviewer's process-local
// GOFLAGS. The same invocation with only the ambient flags rewrites it.
func TestRV01ReviewerGoTestPreservesModuleMetadata(t *testing.T) {
	goTool, err := exec.LookPath("go")
	if err != nil {
		t.Skip("no Go tool on PATH")
	}
	module := t.TempDir()
	files := map[string]string{
		"go.mod":      "module example.com/fixture\n\ngo 1.18\n\nreplace example.com/dep => ./dep\n",
		"dep/go.mod":  "module example.com/dep\n\ngo 1.18\n",
		"dep/dep.go":  "package dep\n\nfunc V() int { return 1 }\n",
		"fix.go":      "package fixture\n\nimport \"example.com/dep\"\n\nfunc F() int { return dep.V() }\n",
		"fix_test.go": "package fixture\n\nimport \"testing\"\n\nfunc TestF(t *testing.T) {\n\tif F() != 1 {\n\t\tt.Fatal()\n\t}\n}\n",
	}
	for name, body := range files {
		if err = os.MkdirAll(filepath.Dir(filepath.Join(module, name)), 0700); err == nil {
			err = os.WriteFile(filepath.Join(module, name), []byte(body), 0600)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	goTest := func(goflags string) string {
		c := exec.Command(goTool, "test", "./...")
		c.Dir = module
		c.Env = append(os.Environ(), "GOFLAGS="+goflags, "GOTOOLCHAIN=local", "GOPROXY=off")
		_, _ = c.CombinedOutput()
		b, err := os.ReadFile(filepath.Join(module, "go.mod"))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	t.Setenv("GOFLAGS", "-mod=mod")
	if after := goTest(deliveryReviewGoFlags()); after != files["go.mod"] {
		t.Fatalf("the reviewer's go test changed go.mod:\n%s", after)
	}
	if after := goTest(os.Getenv("GOFLAGS")); after == files["go.mod"] {
		t.Fatal("the fixture does not exercise -mod=mod: go.mod stayed unchanged")
	}
}

// RV01 actual root topology: the coordinator's tools run under a provider
// daemon, so no pane shell holds the caller and HERDR_PANE_ID stays stale.
// Ply does not guess; the coordinator names its verified pane, which must be
// live in fresh Herdr state, agree with any observed ancestry and not be the
// owner's pane.
func TestRV01DetachedCoordinatorNamesItsVerifiedLivePane(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	oid, _ := reviewRoundCandidate(t, f, "fixed attempt")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w1:p1")
	detached := map[string]any{"pane_shells": map[string]any{"w-fixture:root": 99998, "w-fixture:user": 99999}, "focused_pane": "w-fixture:user"}
	workflowTestModel(t, f.workflowFixture, detached)
	request := func(pane string) error {
		_, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator", CoordinatorPane: pane})
		return err
	}
	if err := request(""); err == nil || !strings.Contains(err.Error(), "unknown") || !strings.Contains(err.Error(), "--coordinator-pane") {
		t.Fatalf("a detached caller was bound without a verified pane: %v", err)
	}
	if err := request("w1:p1"); err == nil || !strings.Contains(err.Error(), "could not be verified as a live Herdr pane") {
		t.Fatalf("a stale pane was accepted as coordinator: %v", err)
	}
	if err := request(o.Transport.PaneID); err == nil || !strings.Contains(err.Error(), "cannot coordinate its own review") {
		t.Fatalf("the owner's pane was accepted as coordinator: %v", err)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"pane_shells": map[string]any{"w-fixture:root": os.Getppid(), "w-fixture:user": 99999}})
	if err := request("w-fixture:user"); err == nil || !strings.Contains(err.Error(), "contradicts pane w-fixture:root") {
		t.Fatalf("a named pane contradicting the observed ancestry was accepted: %v", err)
	}
	if n := len(reviewRoundArgv(t, f, "tab create")); n != 1 {
		t.Fatalf("a refused coordinator created a reviewer tab (%d tab creates)", n)
	}

	workflowTestModel(t, f.workflowFixture, detached)
	out, err := reviewRoundCLI(t, deliveryCLIBinary(t), f.Parent, "w1:p1", "review", o.RunID, "--coordinator", "fixture root coordinator", "--coordinator-pane", "w-fixture:root", "--format", "json")
	if err != nil {
		t.Fatalf("execute review --coordinator-pane: %v\n%s", err, out)
	}
	var run WorkflowRun
	if err = json.Unmarshal([]byte(out), &run); err != nil {
		t.Fatalf("execute review output: %v\n%s", err, out)
	}
	round := reviewRoundLatest(t, run)
	if r := reviewRoundRequestFile(t, round); r.CoordinatorPane != "w-fixture:root" || round.CandidateOID != oid || round.ReviewerStart != "started" {
		t.Fatalf("the verified coordinator pane was not bound: %+v %+v", r, round)
	}
}

// Review r001-f01: when several pane shells hold the caller, a named pane
// must be one of them; an unrelated live pane is never bound.
func TestRV01ExplicitPaneMustBeAmongObservedAncestryPanes(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "fixed attempt")
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w1:p1")
	workflowTestModel(t, f.workflowFixture, map[string]any{"pane_shells": map[string]any{"w-fixture:a": os.Getppid(), "w-fixture:b": os.Getppid(), "w-fixture:user": 99999}, "focused_pane": "w-fixture:user"})
	_, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator", CoordinatorPane: "w-fixture:user"})
	if err == nil || !strings.Contains(err.Error(), "w-fixture:a, w-fixture:b") {
		t.Fatalf("an unrelated live pane was bound while ancestry names other panes: %v", err)
	}
	if n := len(reviewRoundArgv(t, f, "tab create")); n != 1 {
		t.Fatalf("a refused coordinator created a reviewer tab (%d tab creates)", n)
	}
	requested, err := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root coordinator", CoordinatorPane: "w-fixture:b"})
	if err != nil {
		t.Fatal(err)
	}
	if r := reviewRoundRequestFile(t, reviewRoundLatest(t, requested)); r.CoordinatorPane != "w-fixture:b" {
		t.Fatalf("the named ancestry pane was not bound: %+v", r)
	}
}
