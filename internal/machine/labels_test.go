package machine

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func labelsPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "ply", "process-labels.json")
}

func TestLabelsSetListRemove(t *testing.T) {
	path := labelsPath(t)
	security := "security"
	got, err := SetLabel(path, Label{Name: "Google Chrome Helper (Renderer)", DisplayName: "Chrome tab", Category: &security})
	if err != nil || got.Name != "Google Chrome Helper (Renderer)" || *got.Category != "security" {
		t.Fatalf("set: %+v %v", got, err)
	}
	if _, err := SetLabel(path, Label{Name: "node", DisplayName: "Node.js"}); err != nil {
		t.Fatal(err)
	}
	if _, err := SetLabel(path, Label{Name: "node", DisplayName: "Node"}); err != nil {
		t.Fatal(err)
	}
	list, err := ListLabels(path)
	b, _ := json.Marshal(list)
	if err != nil || string(b) != `[{"name":"Google Chrome Helper (Renderer)","display_name":"Chrome tab","category":"security"},{"name":"node","display_name":"Node","category":null}]` {
		t.Fatalf("list: %s %v", b, err)
	}
	removed, err := RemoveLabel(path, "node")
	if err != nil || !removed {
		t.Fatalf("remove: %v %v", removed, err)
	}
	removed, err = RemoveLabel(path, "node")
	if err != nil || removed {
		t.Fatalf("remove again: %v %v", removed, err)
	}
	if _, err := os.Stat(path + ".lock"); err == nil {
		t.Fatal("lock beside the file has the wrong name")
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(path), "process-labels.lock")); err != nil {
		t.Fatal(err)
	}
}

func TestLabelsMissingFileIsEmpty(t *testing.T) {
	path := labelsPath(t)
	labels, err := ReadLabels(path)
	if err != nil || len(labels) != 0 {
		t.Fatalf("got %v %v", labels, err)
	}
	if _, err := os.Stat(filepath.Dir(path)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("reading created the directory")
	}
	list, err := ListLabels(path)
	if err != nil || list == nil || len(list) != 0 {
		t.Fatalf("list: %v %v", list, err)
	}
}

func TestLabelsValidation(t *testing.T) {
	path := labelsPath(t)
	unknown := "antivirus"
	for _, l := range []Label{
		{Name: "", DisplayName: "x"},
		{Name: "a\nb", DisplayName: "x"},
		{Name: strings.Repeat("n", 201), DisplayName: "x"},
		{Name: "x", DisplayName: ""},
		{Name: "x", DisplayName: "two\nlines"},
		{Name: "x", DisplayName: strings.Repeat("é", 61)},
		{Name: "x", DisplayName: "x", Category: &unknown},
	} {
		if _, err := SetLabel(path, l); err == nil || !strings.HasPrefix(err.Error(), "invalid_arguments: ") {
			t.Fatalf("%+v: %v", l, err)
		}
	}
	if _, err := SetLabel(path, Label{Name: strings.Repeat("n", 200), DisplayName: strings.Repeat("é", 60)}); err != nil {
		t.Fatalf("limits count characters: %v", err)
	}
}

func TestLabelsUnreadableFileIsKeptAndReported(t *testing.T) {
	path := labelsPath(t)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadLabels(path); !errors.Is(err, ErrLabelsFormat) {
		t.Fatalf("read: %v", err)
	}
	_, err := SetLabel(path, Label{Name: "x", DisplayName: "x"})
	if !errors.Is(err, ErrLabelsFormat) || !strings.HasPrefix(err.Error(), "unknown_format: ") || !strings.Contains(err.Error(), path) {
		t.Fatalf("set: %v", err)
	}
	if _, err := RemoveLabel(path, "x"); !errors.Is(err, ErrLabelsFormat) {
		t.Fatalf("remove: %v", err)
	}
	if b, _ := os.ReadFile(path); string(b) != "not json" {
		t.Fatalf("file changed: %q", b)
	}
	if err := os.WriteFile(path, nil, 0o000); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o000); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadLabels(path); !errors.Is(err, ErrLabelsRead) || !strings.HasPrefix(err.Error(), "read_failed: ") {
		t.Fatalf("unreadable: %v", err)
	}
}

// A category Ply does not know any more is passed on as stored.
func TestLabelsUnknownStoredCategoryIsPassedOn(t *testing.T) {
	path := labelsPath(t)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"labels":[{"name":"x","display_name":"X","category":"antivirus"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	labels, err := ReadLabels(path)
	if err != nil || *labels["x"].Category != "antivirus" {
		t.Fatalf("got %v %v", labels, err)
	}
}

func TestWithProcessesNamesTheLabelledRows(t *testing.T) {
	s, _ := Read(fakeReader{cores: 8, load: [3]float64{1, 1, 1}, level: 1}, now)
	security := "security"
	labels := map[string]Label{
		"Google Chrome Helper (Renderer)": {Name: "Google Chrome Helper (Renderer)", DisplayName: "Chrome tab"},
		"wdavdaemon":                      {Name: "wdavdaemon", DisplayName: "Microsoft Defender", Category: &security},
	}
	processes := []Process{
		{PID: 1, Name: "wdavdaemon", CPUPercent: 84.2, RSSBytes: 298 << 20},
		{PID: 2, Name: "Google Chrome Helper (Renderer)", CPUPercent: 2, RSSBytes: 1 << 20},
		{PID: 3, Name: "node", CPUPercent: 1, RSSBytes: 1 << 20},
	}
	got := WithLabels(WithProcesses(s, processes, nil), labels, nil)
	b, _ := json.Marshal(got)
	for _, want := range []string{
		`{"pid":1,"name":"wdavdaemon","cpu_percent":84.2,"rss_bytes":312475648,"display_name":"Microsoft Defender","category":"security"}`,
		`{"pid":2,"name":"Google Chrome Helper (Renderer)","cpu_percent":2,"rss_bytes":1048576,"display_name":"Chrome tab","category":null}`,
		`{"pid":3,"name":"node","cpu_percent":1,"rss_bytes":1048576,"display_name":null,"category":null}`,
		`"process_labels_reason":null}`,
	} {
		if !strings.Contains(string(b), want) {
			t.Fatalf("missing %s in\n%s", want, b)
		}
	}
	text := got.Text()
	if !strings.Contains(text, "\n  84.2%  298 MiB  1 wdavdaemon (Microsoft Defender, security)") ||
		!strings.Contains(text, "2 Google Chrome Helper (Renderer) (Chrome tab)\n") || !strings.HasSuffix(text, "3 node") {
		t.Fatalf("text %q", text)
	}
}

func TestWithLabelsReportsAnUnreadableTable(t *testing.T) {
	s, _ := Read(fakeReader{cores: 8, load: [3]float64{1, 1, 1}, level: 1}, now)
	got := WithLabels(WithProcesses(s, []Process{{PID: 1, Name: "x"}}, nil), nil, ErrLabelsFormat)
	b, _ := json.Marshal(got)
	if !strings.Contains(string(b), `"display_name":null,"category":null}],"processes_reason":null,"process_labels_reason":"unknown_format"}`) {
		t.Fatalf("got %s", b)
	}
	if !strings.HasSuffix(got.Text(), "\n  process labels unknown (unknown_format)") {
		t.Fatalf("text %q", got.Text())
	}
	// No processes, no labels reason.
	got = WithLabels(WithProcesses(s, nil, errors.New("boom")), nil, ErrLabelsRead)
	b, _ = json.Marshal(got)
	if !strings.HasSuffix(string(b), `"processes":null,"processes_reason":"read_failed","process_labels_reason":null}`) {
		t.Fatalf("got %s", b)
	}
}
