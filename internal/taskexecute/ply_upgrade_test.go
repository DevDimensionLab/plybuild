package taskexecute

import (
	"bytes"
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"os"
	"path/filepath"
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

// Old preview@1 emits only a request digest. No saved basis or historical Ply
// identity exists in that request. A mismatch cannot distinguish an upgrade
// from changed policy/mandate; accepting it would bypass the approved inputs.
func TestPlyUpgradeLegacyOpaquePreviewLimit(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	p := boundPreview(t, d, in)
	p.Request.Binding = legacyBoundBinding(p.Request, p.Basis)
	if err := os.WriteFile(filepath.Join(root, "bin", "ply"), []byte("#!/bin/sh\n# upgraded\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	out, err := ApplyBoundStart(d, p.Request)
	if err == nil || out.ExecutionEffect != "none" {
		t.Fatalf("unverifiable legacy digest accepted: %+v %v", out, err)
	}
	if _, err := os.Stat(boundIntentPath(root, p.Request.RequestID)); !os.IsNotExist(err) {
		t.Fatalf("unverifiable request created intent: %v", err)
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
