package taskjournal

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLongDurationSnapshotRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name, start, finish string
		seconds             float64
	}{
		{"four_centuries", "1600-01-01T00:00:00Z", "2000-01-01T00:00:00Z", 12622780800},
		{"supported_extremes", "0000-01-01T00:00:00Z", "9999-12-31T23:59:59.999999999Z", 315569520000},
		{"second_boundary", "2026-10-05T21:15:51.999999999Z", "2026-10-05T21:15:52Z", 1e-9},
		{"ordinary_fraction", "2026-10-05T21:15:51.163625Z", "2026-10-05T21:20:10.178909Z", 259.015284},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, s := closedDurationFixture(t, tc.start, tc.finish)
			if got := s.Steps[0].DurationSeconds; got == nil || *got != tc.seconds {
				t.Fatalf("duration = %v, want %.17g", s.Steps[0], tc.seconds)
			}
			path := f.file(t, s)
			before, _ := os.ReadFile(path)
			sources := journalEvidenceFingerprint(t, f.root)
			offline, err := ReadSnapshot(path, "task", Options{})
			if err != nil || offline.SnapshotID != s.SnapshotID {
				t.Fatalf("roundtrip: %v", err)
			}
			after, _ := os.ReadFile(path)
			if !bytes.Equal(before, after) || journalEvidenceFingerprint(t, f.root) != sources {
				t.Fatal("read changed evidence")
			}
		})
	}
}

func TestOfflineRejectsSaturatedLongDuration(t *testing.T) {
	f, s := closedDurationFixture(t, "1600-01-01T00:00:00Z", "2000-01-01T00:00:00Z")
	for _, seconds := range []float64{time.Duration(1<<63 - 1).Seconds(), 12622780799} {
		s.Steps[0].DurationSeconds = ptr(seconds)
		s.SnapshotID = snapshotDigest(s)
		path := f.file(t, s)
		before, _ := os.ReadFile(path)
		sources := journalEvidenceFingerprint(t, f.root)
		_, err := ReadSnapshot(path, "task", Options{})
		if err == nil || !strings.Contains(err.Error(), "duration") {
			t.Fatalf("false duration %.17g: %v", seconds, err)
		}
		after, _ := os.ReadFile(path)
		if !bytes.Equal(before, after) || journalEvidenceFingerprint(t, f.root) != sources {
			t.Fatal("rejection changed evidence")
		}
	}
}

func TestLongTimelineProportionsAndKnownZero(t *testing.T) {
	f := fixture(t)
	in := f.input
	for i, bounds := range [][2]string{
		{"0001-01-01T00:00:00Z", "0401-01-01T00:00:00Z"},
		{"0401-01-01T00:00:00Z", "0801-01-01T00:00:00Z"},
		{"0801-01-01T00:00:00Z", "1201-01-01T00:00:00Z"},
	} {
		in.Step = &StepRef{[]string{"first", "middle", "last"}[i], "implementation", nil}
		in.PublicationKey, in.Type, in.OccurredAt = in.Step.ID+":start", "step_started", ptr(bounds[0])
		f.append(t, in)
		in.PublicationKey, in.Type, in.OccurredAt = in.Step.ID+":end", "step_finished", ptr(bounds[1])
		f.append(t, in)
	}
	s := f.show(t)
	text := Text(s, "timeline")
	if !strings.Contains(text, "Intervals on one UTC axis: 0001-01-01T00:00:00Z .. 1201-01-01T00:00:00Z") {
		t.Fatal(text)
	}
	for _, bar := range []string{
		strings.Repeat("=", 16) + strings.Repeat(" ", 32),
		strings.Repeat(" ", 15) + strings.Repeat("=", 17) + strings.Repeat(" ", 16),
		strings.Repeat(" ", 31) + strings.Repeat("=", 17),
	} {
		if !strings.Contains(text, "|"+bar+"|") {
			t.Fatalf("missing bar %q:\n%s", bar, text)
		}
	}
}

// The acceptance entrypoint supplies the exact make-built binary. All commands
// below read only an exported snapshot in a private fixture, without discovery.
func TestLongDurationOfflineCLI(t *testing.T) {
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		t.Skip("exact candidate binary supplied by acceptance")
	}
	f, s := closedDurationFixture(t, "0001-01-01T00:00:00Z", "1201-01-01T00:00:00Z")
	path := f.file(t, s)
	before, _ := os.ReadFile(path)
	sources := journalEvidenceFingerprint(t, f.root)
	for _, format := range []string{"json", "text"} {
		cmd := exec.Command(binary, "workspace", "task", "journal", "show", "task", "--snapshot", path, "--format", format, "--view", "timeline")
		cmd.Dir = t.TempDir()
		output, err := cmd.CombinedOutput()
		t.Logf("offline %s command %q; result: %v; output:\n%s", format, cmd.Args, err, output)
		if err != nil {
			t.Fatalf("offline %s: %v: %s", format, err, output)
		}
		if format == "json" {
			out := filepath.Join(f.root, "cli-snapshot.json")
			write(t, out, output)
			got, err := ReadSnapshot(out, "task", Options{})
			if err != nil || got.SnapshotID != s.SnapshotID {
				t.Fatalf("CLI roundtrip: %v", err)
			}
		} else if !bytes.Contains(output, []byte("|"+strings.Repeat("=", 48)+"|")) {
			t.Fatalf("CLI timeline: %s", output)
		}
	}

	s.Steps[0].DurationSeconds = ptr(time.Duration(1<<63 - 1).Seconds())
	s.SnapshotID = snapshotDigest(s)
	rejected := f.file(t, s)
	rejectedBefore, _ := os.ReadFile(rejected)
	cmd := exec.Command(binary, "workspace", "task", "journal", "show", "task", "--snapshot", rejected, "--format", "json")
	cmd.Dir = t.TempDir()
	output, err := cmd.CombinedOutput()
	t.Logf("rejected saturated snapshot command %q; result: %v; output:\n%s", cmd.Args, err, output)
	if err == nil || !bytes.Contains(output, []byte("duration")) {
		t.Fatalf("CLI accepted false duration: %v: %s", err, output)
	}
	rejectedAfter, _ := os.ReadFile(rejected)
	if !bytes.Equal(rejectedBefore, rejectedAfter) {
		t.Fatal("CLI rewrote rejected historical snapshot")
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) || journalEvidenceFingerprint(t, f.root) != sources {
		t.Fatal("CLI changed evidence")
	}
}

// Include immutable contribution records, not only workspace/native sources.
func journalEvidenceFingerprint(t *testing.T, root string) string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(filepath.Join(root, ".ply"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files[path] = hash(b)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return digest(files)
}

func TestElapsedSecondsPreservesDurationValues(t *testing.T) {
	start, _ := parseTime("2026-10-05T21:15:51.999999999Z")
	for _, d := range []time.Duration{0, 1, 2, 999999999, 1000000000, 1123456789, 259015284000, 1<<53 - 1, 1<<63 - 2, 1<<63 - 1} {
		for _, d := range []time.Duration{d, -d} {
			finish := start.Add(d)
			if got, want := elapsedSeconds(start, finish), d.Seconds(); got != want {
				t.Fatalf("%s: %.17g != %.17g", d, got, want)
			}
		}
	}
}

func TestLongIntervalUnknownStatesStayNull(t *testing.T) {
	for _, mode := range []string{"missing_start", "missing_finish", "clock", "uncertainty", "reversed", "no_end"} {
		t.Run(mode, func(t *testing.T) {
			f := fixture(t)
			start := f.input
			start.Type, start.PublicationKey = "step_started", "long:start"
			start.Step = &StepRef{"long", "implementation", nil}
			start.OccurredAt = ptr("0001-01-01T00:00:00Z")
			finish := start
			finish.Type, finish.PublicationKey, finish.OccurredAt = "step_finished", "long:finish", ptr("1201-01-01T00:00:00Z")
			state := "time_unknown"
			switch mode {
			case "missing_start":
				start.OccurredAt, start.TimeBasis = nil, unknownTime()
			case "missing_finish":
				finish.OccurredAt, finish.TimeBasis = nil, unknownTime()
			case "clock":
				finish.TimeBasis.Clock = ptr("other")
			case "uncertainty":
				start.TimeBasis.Uncertainty = &Uncertainty{"0001-01-01T00:00:00Z", "0001-01-01T00:00:01Z"}
			case "reversed":
				start.OccurredAt, finish.OccurredAt = finish.OccurredAt, start.OccurredAt
				state = "time_conflict"
			case "no_end":
				state = "no_end_recorded"
			}
			f.append(t, start)
			if mode != "no_end" {
				f.append(t, finish)
			}
			s := f.show(t)
			if p := s.Steps[0]; p.State != state || p.DurationSeconds != nil {
				t.Fatalf("changed state: %+v", p)
			}
			path := f.file(t, s)
			before := journalEvidenceFingerprint(t, f.root)
			offline, err := ReadSnapshot(path, "task", Options{})
			if err != nil || !equal(offline.Coverage, s.Coverage) {
				t.Fatalf("offline state/coverage: %v", err)
			}
			text := Text(offline, "timeline")
			if !strings.Contains(text, "|"+strings.Repeat(" ", 48)+"|") {
				t.Fatal(text)
			}
			if before != journalEvidenceFingerprint(t, f.root) {
				t.Fatal("state roundtrip changed evidence")
			}
		})
	}
}
