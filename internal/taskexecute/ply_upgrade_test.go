package taskexecute

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPlyUpgradeBoundPreviewApply(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	p := boundPreview(t, d, in)
	if err := os.WriteFile(filepath.Join(root, "bin", "ply"), []byte("#!/bin/sh\n# TC0E6E1FE upgraded Ply\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	out, err := ApplyBoundStart(d, p.Request)
	if err != nil {
		t.Fatalf("upgrade invalidated approved request: %+v %v", out, err)
	}
	again, err := ApplyBoundStart(d, p.Request)
	if err != nil || again.ExpectedRunID != out.ExpectedRunID {
		t.Fatalf("replay: %+v %v", again, err)
	}
	var model struct{ Tabs, Starts, Prompts int }
	raw, err := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
	if err != nil || json.Unmarshal(raw, &model) != nil || model.Tabs != 1 || model.Starts != 1 || model.Prompts != 1 {
		t.Fatalf("duplicate effects: %s %v", raw, err)
	}
}

// A legacy request contains only its opaque checksum. Startup uses the current
// named Goal and authority without recovering a prior basis or old executable.
func TestPlyUpgradeLegacyOpaqueRequestWithoutBasisOrHistory(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	p := boundPreview(t, d, in)
	p.Request.Binding = legacyBoundBinding(p.Request, p.Basis)
	requestPath := filepath.Join(root, "legacy-request.json")
	original, _ := taskrun.Canonical(p.Request)
	if err := os.WriteFile(requestPath, original, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "ply"), []byte("#!/bin/sh\n# upgraded\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(boundIntentPath(root, p.Request.RequestID)); !os.IsNotExist(err) {
		t.Fatal("fixture retained a basis")
	}
	request, err := ReadBoundStartRequest(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ApplyBoundStart(d, request)
	if err != nil {
		t.Fatalf("old opaque request required history: %+v %v", out, err)
	}
	again, err := ApplyBoundStart(d, request)
	if err != nil || again.ExpectedRunID != out.ExpectedRunID {
		t.Fatalf("legacy replay: %+v %v", again, err)
	}
	var model struct{ Starts, Prompts int }
	raw, err := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
	if err != nil || json.Unmarshal(raw, &model) != nil || model.Starts != 1 || model.Prompts != 1 {
		t.Fatalf("legacy request duplicated input: %s %v", raw, err)
	}
	after, err := os.ReadFile(requestPath)
	if err != nil || !bytes.Equal(original, after) {
		t.Fatalf("old request rewritten: %v", err)
	}
}

func TestPlyUpgradeInformationalBindingKeepsRequestIdentity(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	p := boundPreview(t, d, in)
	out, err := ApplyBoundStart(d, p.Request)
	if err != nil {
		t.Fatal(err)
	}
	path := boundIntentPath(root, p.Request.RequestID)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, binding := range []string{"sha256:" + strings.Repeat("0", 64), ""} {
		request := p.Request
		request.Binding = binding
		again, err := ApplyBoundStart(d, request)
		if err != nil || again.ExpectedRunID != out.ExpectedRunID {
			t.Fatalf("informational checksum replaced identity: %+v %v", again, err)
		}
	}
	changed := p.Request
	changed.Input.SpecID = "another-goal"
	if _, err := ApplyBoundStart(d, changed); err == nil {
		t.Fatal("same identity accepted a different named Goal")
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("replay rewrote original request")
	}
	var model struct{ Starts, Prompts int }
	raw, _ := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
	if json.Unmarshal(raw, &model) != nil || model.Starts != 1 || model.Prompts != 1 {
		t.Fatalf("checksum replay duplicated input: %s", raw)
	}
}

func TestPlyUpgradeStartUsesCurrentNamedGoalAndPolicy(t *testing.T) {
	for _, change := range []string{"goal", "policy"} {
		t.Run(change, func(t *testing.T) {
			d, in, root, _ := boundFixture(t)
			p := boundPreview(t, d, in)
			if change == "goal" {
				boundChangeGoal(t, d, root, "goal")
			} else {
				path := filepath.Join(root, "codex-home", "config.toml")
				raw, _ := os.ReadFile(path)
				if err := os.WriteFile(path, bytes.Replace(raw, []byte("network = false"), []byte("network = true"), 1), 0600); err != nil {
					t.Fatal(err)
				}
			}
			current := boundPreview(t, d, in)
			if _, err := ApplyBoundStart(d, p.Request); err != nil {
				t.Fatalf("historical preview overrode current named Goal/policy: %v", err)
			}
			intent, err := readBoundIntent(root, p.Request.RequestID)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(intent.Basis, current.Basis) {
				t.Fatal("startup did not capture the current registered Goal/policy")
			}
			if intent.Request.Binding != p.Request.Binding {
				t.Fatal("startup rewrote historical checksum")
			}
		})
	}
}

func TestPlyUpgradeLegacyRetainedIntentKeepsIdentity(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	p := boundPreview(t, d, in)
	p.Request.Binding = legacyBoundBinding(p.Request, p.Basis)
	intent := boundStartIntent{Kind: "ply.workflow.bound-start-intent", SchemaVersion: 1, Request: p.Request, Basis: p.Basis, RunID: boundRunID(p.Request), RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
	raw, err := taskrun.Canonical(intent)
	if err != nil {
		t.Fatal(err)
	}
	path := boundIntentPath(root, p.Request.RequestID)
	if err = writeOnce(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(p.Basis.Runtime.Runtime.PlyExecutable.Path); err != nil {
		t.Fatal(err)
	}
	out, err := ApplyBoundStart(d, p.Request)
	if err != nil || out.ExpectedRunID != intent.RunID || out.State != "retained_attempt" {
		t.Fatalf("legacy intent lost identity after upgrade: %+v %v", out, err)
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(raw, after) {
		t.Fatalf("legacy intent was rewritten: %v", err)
	}
}

func legacyBoundBinding(r BoundStartRequest, b BoundStartBasis) string {
	r.Binding = ""
	if b.DependencySHA256 != "" {
		b.RegistrySHA256 = ""
	}
	return boundHash(struct {
		Request BoundStartRequest `json:"request"`
		Basis   BoundStartBasis   `json:"basis"`
	}{r, b})
}

func TestPlyUpgradeBindingOptionalAtRequestBoundary(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	p := boundPreview(t, d, in)
	body, _ := json.Marshal(p.Request)
	var object map[string]any
	if err := json.Unmarshal(body, &object); err != nil {
		t.Fatal(err)
	}
	delete(object, "binding")
	body, _ = json.Marshal(object)
	path := filepath.Join(root, "request-without-checksum.json")
	if err := os.WriteFile(path, body, 0600); err != nil {
		t.Fatal(err)
	}
	request, err := ReadBoundStartRequest(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ApplyBoundStart(d, request); err != nil {
		t.Fatalf("checksum omission blocked valid current startup: %v", err)
	}
}

func TestPlyUpgradeCurrentSnapshotStartsAfterEarlierQueueChange(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	p := boundPreview(t, d, in)
	d.Fault = func(point string) error {
		if point == "bound_start_before_lock" {
			boundChangeQueue(t, d, root)
		}
		return nil
	}
	out, err := ApplyBoundStart(d, p.Request)
	if err != nil {
		t.Fatalf("change before current snapshot required another preview: %+v %v", out, err)
	}
}

func TestPlyUpgradeUnknownAttemptIgnoresChecksumWithoutReplay(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	p := boundPreview(t, d, in)
	d.Fault = func(point string) error {
		if point == "bound_start_intent_published" {
			return errors.New("lost response after durable intent")
		}
		return nil
	}
	first, err := ApplyBoundStart(d, p.Request)
	if err == nil || first.ExecutionEffect != "unknown" {
		t.Fatalf("unknown attempt not preserved: %+v %v", first, err)
	}
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	d.Fault = nil
	request := p.Request
	request.Binding = ""
	again, err := ApplyBoundStart(d, request)
	if err != nil || again.State != "retained_attempt" || again.ExecutionEffect != "unknown" || again.ExpectedRunID != first.ExpectedRunID {
		t.Fatalf("checksum replay changed unknown attempt: %+v %v", again, err)
	}
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("unknown attempt replay wrote native state")
	}
	if _, err := os.Stat(filepath.Join(root, "bin", "recovery-model.json")); !os.IsNotExist(err) {
		t.Fatal("unknown attempt replay launched provider or input")
	}
}
