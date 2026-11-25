package machine

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestParseProcesses(t *testing.T) {
	out := "  82334 111.0 404416 /opt/homebrew/bin/node\n" +
		"   375   2.1 250304 /Applications/Google Chrome.app/Contents/Frameworks/Google Chrome Helper (Renderer)\n" +
		"     9   0.0     16 launchd\n" +
		"garbage line\n" +
		"   777  50.0    100 /bin/ps\n"
	got := parseProcesses([]byte(out), 777)
	want := []Process{
		{PID: 82334, Name: "node", CPUPercent: 111, RSSBytes: 404416 * 1024},
		{PID: 375, Name: "Google Chrome Helper (Renderer)", CPUPercent: 2.1, RSSBytes: 250304 * 1024},
		{PID: 9, Name: "launchd", CPUPercent: 0, RSSBytes: 16 * 1024},
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("got  %+v\nwant %+v", got, want)
	}
}

func TestParseProcessesKeepsTopTen(t *testing.T) {
	var b strings.Builder
	for pid := 1; pid <= 12; pid++ {
		cpu := float64(pid)
		if pid == 11 || pid == 12 {
			cpu = 10 // ties with pid 10, ordered by rss then pid
		}
		fmt.Fprintf(&b, "%d %.1f %d p%d\n", pid, cpu, 100-pid, pid)
	}
	got := parseProcesses([]byte(b.String()), 0)
	var pids []int
	for _, p := range got {
		pids = append(pids, p.PID)
	}
	if fmt.Sprint(pids) != "[10 11 12 9 8 7 6 5 4 3]" {
		t.Fatalf("got %v", pids)
	}
	if parseProcesses(nil, 0) == nil {
		t.Fatal("no processes must be an empty list, not null")
	}
}

func TestWithProcessesJSON(t *testing.T) {
	s, _ := Read(fakeReader{cores: 8, load: [3]float64{20, 1, 0.5}, level: 1}, now)
	base, _ := json.Marshal(s)
	b, _ := json.Marshal(WithProcesses(s, []Process{{PID: 1, Name: "x", CPUPercent: 1.5, RSSBytes: 2048}}, nil))
	want := strings.TrimSuffix(string(base), "}") + `,"processes":[{"pid":1,"name":"x","cpu_percent":1.5,"rss_bytes":2048,"display_name":null,"category":null}],"processes_reason":null,"process_labels_reason":null}`
	if string(b) != want {
		t.Fatalf("got  %s\nwant %s", b, want)
	}
	b, _ = json.Marshal(WithProcesses(s, nil, errors.New("boom")))
	if !strings.HasSuffix(string(b), `"processes":null,"processes_reason":"read_failed","process_labels_reason":null}`) || !strings.Contains(string(b), `"load_percent"`) {
		t.Fatalf("got %s", b)
	}
}

func TestWithProcessesText(t *testing.T) {
	s, _ := Read(fakeReader{cores: 16, load: [3]float64{4.91, 9.34, 9.97}, level: 1}, now)
	got := WithProcesses(s, []Process{{PID: 82334, Name: "node", CPUPercent: 111, RSSBytes: 404416 * 1024}}, nil).Text()
	if got != s.Text()+"\n  111.0%  395 MiB  82334 node" {
		t.Fatalf("got %q", got)
	}
	got = WithProcesses(s, nil, errors.New("boom")).Text()
	if got != s.Text()+"\n  processes unknown (read_failed)" {
		t.Fatalf("got %q", got)
	}
}

func TestReadProcessesOnThisMachine(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS only")
	}
	got, err := System().(ProcessReader).Processes()
	if err != nil || len(got) == 0 || len(got) > 10 {
		t.Fatalf("got %d processes, %v", len(got), err)
	}
	for _, p := range got {
		if p.PID <= 0 || strings.Contains(p.Name, "/") || p.Name == "" {
			t.Fatalf("bad row %+v", p)
		}
	}
}
