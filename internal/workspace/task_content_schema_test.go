package workspace

import (
	"bytes"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"os"
	"path/filepath"
	"testing"
)

func TestTaskContentRejectsOversizedDraftAndSourceBeforePublication(t *testing.T) {
	for _, kind := range []string{"draft", "source"} {
		t.Run(kind, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			root := filepath.Dir(f.wrapper)
			in := contentFixtureSpec(t, f, "fixture/oversized", nil)
			before, err := os.ReadFile(workItemsPath(root))
			if err != nil {
				t.Fatal(err)
			}
			if kind == "draft" {
				raw, err := os.ReadFile(in.File)
				if err != nil {
					t.Fatal(err)
				}
				// Valid JSON padded just beyond the public 256 KiB file boundary.
				raw = append(raw, bytes.Repeat([]byte(" "), (256<<10)+1-len(raw))...)
				if err = os.WriteFile(in.File, raw, 0600); err != nil {
					t.Fatal(err)
				}
			} else {
				// The declaration remains valid, but the physical source exceeds 4 MiB.
				if err = os.WriteFile(filepath.Join(root, "solution.md"), bytes.Repeat([]byte("x"), (4<<20)+1), 0600); err != nil {
					t.Fatal(err)
				}
			}
			out, err := RecordTaskSpec(f.dependencies, in)
			var contentErr *TaskContentError
			if !errors.As(err, &contentErr) || contentErr.Code != "task_content_invalid_input" || out.Attempted.ObjectPublicationStarted {
				t.Fatalf("oversized %s crossed the publication boundary: %#v %v", kind, out, err)
			}
			after, err := os.ReadFile(workItemsPath(root))
			if err != nil || !bytes.Equal(before, after) {
				t.Fatal("oversized input changed the registry")
			}
		})
	}
}

func TestTaskContentDraftStrictShapeAndLimits(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	in := contentFixtureSpec(t, f, "fixture/spec", nil)
	raw, _ := os.ReadFile(in.File)
	for _, name := range []string{"unknown", "missing", "unresolved_ref", "technical_na", "requirement_order", "phase_ref", "dependencies"} {
		t.Run(name, func(t *testing.T) {
			v, err := canonicaljson.DecodeStrict(raw)
			if err != nil {
				t.Fatal(err)
			}
			m := contentFields(v)
			switch name {
			case "unknown":
				m["unknown"] = true
			case "missing":
				delete(m, "problem")
			case "unresolved_ref":
				m["documents"] = []canonicaljson.Value{}
			case "technical_na":
				parts := contentFields(m["parts"])
				parts["technical"] = contentObject(map[string]canonicaljson.Value{"state": "not_applicable", "reason": nil, "documents": []canonicaljson.Value{}})
				m["parts"] = contentObject(parts)
			case "requirement_order":
				a := contentArray(m, "requirements")
				m["requirements"] = append(a, a[0])
			case "phase_ref":
				a := contentArray(m, "phases")
				p := contentFields(a[0])
				p["requirement_ids"] = []canonicaljson.Value{"missing"}
				m["phases"] = []canonicaljson.Value{contentObject(p)}
			case "dependencies":
				m["dependencies"] = []canonicaljson.Value{"future-task"}
			}
			b, _ := canonicaljson.Marshal(contentObject(m))
			path := filepath.Join(root, name+".json")
			if err = os.WriteFile(path, b, 0600); err != nil {
				t.Fatal(err)
			}
			out, err := RecordTaskSpec(f.dependencies, TaskContentInput{TaskID: "task", File: path})
			if err == nil || out.Attempted.ObjectPublicationStarted {
				t.Fatalf("invalid draft accepted: %#v %v", out, err)
			}
		})
	}
}

func TestTaskContentVersionThreeEmptyRoundTrip(t *testing.T) {
	r := emptyWorkItemRegistry()
	b, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	got, err := decodeWorkItemRegistry(b)
	if err != nil || got.FormatVersion != 3 || got.TaskContentPublications == nil || got.TaskSpecPolicies == nil {
		t.Fatalf("v3 shape: %#v %v", got, err)
	}
	again, err := encodeWorkItemRegistry(got)
	if err != nil || string(b) != string(again) {
		t.Fatal("v3 did not round trip")
	}
}

func TestTaskSpecRequirementRemovalAndOrderedPhases(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	in := contentFixtureSpec(t, f, "fixture/r1", nil)
	raw, _ := os.ReadFile(in.File)
	v, _ := canonicaljson.DecodeStrict(raw)
	m := contentFields(v)
	phase := contentFields(contentArray(m, "phases")[0])
	phase["id"] = "z-first"
	second := contentFields(contentArray(m, "phases")[0])
	second["id"] = "a-second"
	m["phases"] = []canonicaljson.Value{contentObject(phase), contentObject(second)}
	raw, _ = canonicaljson.Marshal(contentObject(m))
	if err := os.WriteFile(in.File, raw, 0600); err != nil {
		t.Fatal(err)
	}
	first, err := RecordTaskSpec(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	saved, err := readContentManifest(f.dependencies.TaskContent, root, first.OutcomeRef.ManifestSHA256)
	if err != nil {
		t.Fatal(err)
	}
	phases := contentArray(contentFields(saved), "phases")
	if contentString(contentFields(phases[0]), "id") != "z-first" || contentString(contentFields(phases[1]), "id") != "a-second" {
		t.Fatal("ordered phases were silently sorted")
	}
	m["publication_key"] = "fixture/r2"
	m["expected_previous"] = contentRefValue(TaskRevisionRef{1, first.OutcomeRef.ManifestSHA256})
	m["requirements"] = []canonicaljson.Value{}
	m["phases"] = []canonicaljson.Value{}
	raw, _ = canonicaljson.Marshal(contentObject(m))
	if err = os.WriteFile(in.File, raw, 0600); err != nil {
		t.Fatal(err)
	}
	rejected, err := RecordTaskSpec(f.dependencies, in)
	if err == nil || rejected.Attempted.ObjectPublicationStarted {
		t.Fatal("silent requirement removal was published")
	}
	m["removed_requirement_ids"] = []canonicaljson.Value{"f-01"}
	raw, _ = canonicaljson.Marshal(contentObject(m))
	if err = os.WriteFile(in.File, raw, 0600); err != nil {
		t.Fatal(err)
	}
	secondRevision, err := RecordTaskSpec(f.dependencies, in)
	if err != nil || *secondRevision.OutcomeRef.Revision != 2 {
		t.Fatalf("explicit draft removal failed: %#v %v", secondRevision, err)
	}
	saved, err = readContentManifest(f.dependencies.TaskContent, root, secondRevision.OutcomeRef.ManifestSHA256)
	if err != nil {
		t.Fatal(err)
	}
	if err = taskSpecStructurallyReady(saved); err == nil {
		t.Fatal("empty requirements became ready")
	}
}
