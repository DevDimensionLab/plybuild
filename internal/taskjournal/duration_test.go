package taskjournal

import (
	"bytes"
	"encoding/json"
	"math"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func closedDurationFixture(t *testing.T, start, finish string) (testFixture, Snapshot) {
	t.Helper()
	f := fixture(t)
	in := f.input
	in.PublicationKey, in.Type = "duration:start", "step_started"
	in.Step = &StepRef{"measured", "implementation", nil}
	in.OccurredAt = ptr(start)
	in.TimeBasis.Precision = "nanosecond"
	f.append(t, in)
	in.PublicationKey, in.Type = "duration:finish", "step_finished"
	in.OccurredAt = ptr(finish)
	f.append(t, in)
	return f, f.show(t)
}

func TestFractionalDurationSnapshotRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name, start, finish, number, textSeconds string
		seconds                                  float64
	}{
		{"native_wait", "2026-10-05T21:15:51.163625Z", "2026-10-05T21:20:10.178909Z", "259.015284", "259.015284", 259.015284},
		{"nanosecond", "2026-10-05T21:15:51Z", "2026-10-05T21:15:51.000000001Z", "1e-9", "0.000000001", 1e-9},
		{"full_fraction", "2026-10-05T21:15:51Z", "2026-10-05T21:15:52.123456789Z", "1.123456789", "1.123456789", 1.123456789},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, s := closedDurationFixture(t, tc.start, tc.finish)
			before := nativeFingerprint(t, f.root)
			if len(s.Steps) != 1 || s.Steps[0].DurationSeconds == nil || *s.Steps[0].DurationSeconds != tc.seconds {
				t.Fatalf("measured fraction was lost: %+v", s.Steps)
			}
			b, err := Canonical(s)
			if err != nil {
				t.Fatalf("installed JSON show cannot encode the measured duration: %v", err)
			}
			if !bytes.Contains(b, []byte(`"duration_seconds":`+tc.number)) {
				t.Fatalf("fraction changed in serialized snapshot: want %s", tc.number)
			}
			offline, err := ReadSnapshot(f.file(t, s), "task", Options{Order: "recorded"})
			if err != nil || offline.SnapshotID != s.SnapshotID || offline.Steps[0].DurationSeconds == nil || *offline.Steps[0].DurationSeconds != tc.seconds {
				t.Fatalf("offline fractional snapshot rejected or changed: %+v %v", offline.Steps, err)
			}
			if !strings.Contains(Text(offline, "timeline"), tc.textSeconds+" seconds") {
				t.Fatalf("text view rounded the preserved duration: %s", Text(offline, "timeline"))
			}
			if nativeFingerprint(t, f.root) != before {
				t.Fatal("read-only roundtrip changed native sources")
			}
		})
	}
}

func TestSnapshotFractionalDecodingKeepsStrictBoundaries(t *testing.T) {
	for _, number := range []string{"259.0152840000000000000001", "1e-400", "-0", "-0.0", "-0.125", "1e309", "9007199254740993.5"} {
		if _, err := decodeJournalJSON([]byte(`{"kind":"ply.workspace.task-journal-snapshot","steps":[{"duration_seconds":` + number + `}]}`)); err == nil {
			t.Errorf("unsafe or precision-losing duration %s accepted", number)
		}
	}
	for _, document := range []string{
		`{"kind":"ply.workspace.task-journal-snapshot","steps":[{"duration_seconds":0.125,"duration_seconds":0.125}]}`,
		`{"kind":"ply.workspace.task-journal-snapshot","steps":[{"duration_seconds":0.125}],"schema_version":1.0}`,
		`{"kind":"ply.workspace.task-journal-snapshot","steps":{"0":{"duration_seconds":0.125}}}`,
		`{"kind":"ply.workspace.task-journal-snapshot","steps":[{"duration_seconds":0.125}],"unknown":"\uD800"}`,
		`{"kind":"ply.workspace.task-journal-snapshot","steps":[{"duration_seconds":0.125}]} {}`,
	} {
		if _, err := decodeJournalJSON([]byte(document)); err == nil {
			t.Errorf("malformed or out-of-scope fractional snapshot accepted: %s", document)
		}
	}
	for _, number := range []string{"1.25e-1", "0.1250"} {
		value, err := decodeJournalJSON([]byte(`{"kind":"ply.workspace.task-journal-snapshot","steps":[{"duration_seconds":` + number + `}]}`))
		if err != nil {
			t.Fatal(err)
		}
		b, err := canonicaljson.Marshal(value)
		if err != nil || !bytes.Contains(b, []byte(`"duration_seconds":0.125`)) {
			t.Fatalf("equivalent decimal spelling changed value: %s %v", b, err)
		}
	}
	if _, err := decodeJournalJSON([]byte(`{"kind":"ply.workspace.task-journal-snapshot","steps":[{"duration_seconds":0e999999999999999999}]} `)); err != nil {
		t.Fatalf("exact zero exponent was not handled safely: %v", err)
	}
}

func TestOfflineRejectsFalseFractionalDuration(t *testing.T) {
	f, s := closedDurationFixture(t, "2026-10-05T21:15:51.163625Z", "2026-10-05T21:20:10.178909Z")
	for _, value := range []float64{259, 259.015285, math.Nextafter(259.015284, math.Inf(1))} {
		s.Steps[0].DurationSeconds = ptr(value)
		s.SnapshotID = snapshotDigest(s)
		if _, err := ReadSnapshot(f.file(t, s), "task", Options{}); err == nil {
			t.Fatalf("digest-valid false duration %g accepted", value)
		}
	}
}

func TestFractionalNumbersRemainScopedToSnapshotDurations(t *testing.T) {
	valid := map[string]any{"kind": "ply.workspace.task-journal-snapshot", "schema_version": 1, "steps": []any{map[string]any{"duration_seconds": 0.125}}}
	if _, err := Canonical(valid); err != nil {
		t.Fatalf("snapshot duration rejected: %v", err)
	}
	for name, document := range map[string]any{
		"other_kind":       map[string]any{"kind": "ply.workspace.task-journal-event-input", "steps": []any{map[string]any{"duration_seconds": 0.125}}},
		"version":          map[string]any{"kind": "ply.workspace.task-journal-snapshot", "schema_version": 1.25},
		"nested_payload":   map[string]any{"kind": "ply.workspace.task-journal-snapshot", "events": []any{map[string]any{"data": map[string]any{"duration_seconds": 0.125}}}},
		"wrong_step_field": map[string]any{"kind": "ply.workspace.task-journal-snapshot", "steps": []any{map[string]any{"sequence": 0.125}}},
		"scalar":           0.125,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Canonical(document); err == nil {
				t.Fatal("unrelated journal contract acquired fractional numbers")
			}
		})
	}
	if _, err := canonicaljson.DecodeStrict([]byte(`{"duration_seconds":0.125}`)); err == nil {
		t.Fatal("global strict integer parser was weakened")
	}
}

func TestWholeSecondSnapshotKeepsLegacyDigest(t *testing.T) {
	f, s := closedDurationFixture(t, "2026-10-05T21:15:51Z", "2026-10-05T21:20:10Z")
	// Reproduce the previous integer-only snapshot digest on the same bytes.
	legacyCanonical := func(v any) []byte {
		t.Helper()
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		x, err := canonicaljson.DecodeStrict(b)
		if err != nil {
			t.Fatal(err)
		}
		b, err = canonicaljson.Marshal(x)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	oldBytes := legacyCanonical(s)
	var fields map[string]any
	if err := json.Unmarshal(oldBytes, &fields); err != nil {
		t.Fatal(err)
	}
	delete(fields, "snapshot_id")
	delete(fields, "selection")
	want := "sha256:" + hash(legacyCanonical(fields))
	if s.SnapshotID != want {
		t.Fatalf("legacy integer snapshot digest changed: got %s want %s", s.SnapshotID, want)
	}
	path := f.file(t, s)
	write(t, path, oldBytes)
	if _, err := ReadSnapshot(path, "task", Options{}); err != nil {
		t.Fatalf("legacy snapshot became unreadable: %v", err)
	}
}
