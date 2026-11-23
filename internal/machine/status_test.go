package machine

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeReader struct {
	cores               int
	load                [3]float64
	loadErr             error
	level               uint32
	levelErr            error
	memTotal            uint64
	memTotalErr         error
	vmstat              string
	vmstatErr           error
	swapUsed, swapTotal uint64
	swapErr             error
}

// vmstat16k is vm_stat output with 16 KiB pages: used = (1000 - 100 + 300 + 200) pages.
const vmstat16k = `Mach Virtual Memory Statistics: (page size of 16384 bytes)
Pages free:                              235923.
Pages wired down:                           300.
Pages purgeable:                            100.
Anonymous pages:                           1000.
Pages occupied by compressor:               200.
`

func (f fakeReader) Cores() int                           { return f.cores }
func (f fakeReader) LoadAverage() ([3]float64, error)     { return f.load, f.loadErr }
func (f fakeReader) MemoryPressureLevel() (uint32, error) { return f.level, f.levelErr }
func (f fakeReader) MemoryTotal() (uint64, error) {
	if f.memTotal == 0 && f.memTotalErr == nil {
		return 1 << 34, nil
	}
	return f.memTotal, f.memTotalErr
}
func (f fakeReader) VMStat() ([]byte, error)       { return []byte(f.vmstat), f.vmstatErr }
func (f fakeReader) Swap() (uint64, uint64, error) { return f.swapUsed, f.swapTotal, f.swapErr }

var now = time.Date(2026, 11, 23, 21, 4, 5, 123, time.FixedZone("CET", 3600))

func TestLoadPercent(t *testing.T) {
	s, err := Read(fakeReader{cores: 16, load: [3]float64{4.91, 9.34, 9.97}, level: 1}, now)
	if err != nil {
		t.Fatal(err)
	}
	if s.LoadPercent != (Windows{30.7, 58.4, 62.3}) || s.Load != (Windows{4.91, 9.34, 9.97}) || s.Cores != 16 {
		t.Fatalf("got %+v", s)
	}
	s, _ = Read(fakeReader{cores: 8, load: [3]float64{20, 1.23456, 0}, level: 1}, now)
	if s.LoadPercent.M1 != 250 || s.Load.M5 != 1.23 {
		t.Fatalf("got %+v", s)
	}
}

func TestMemoryPressure(t *testing.T) {
	for _, c := range []struct {
		level         uint32
		err           error
		value, reason string
	}{
		{1, nil, "normal", ""}, {2, nil, "warn", ""}, {4, nil, "critical", ""},
		{3, nil, "", ReasonUnknownLevel}, {0, errors.New("boom"), "", ReasonReadFailed},
	} {
		s, err := Read(fakeReader{cores: 1, level: c.level, levelErr: c.err}, now)
		if err != nil {
			t.Fatal(err)
		}
		if deref(s.MemoryPressure) != c.value || deref(s.MemoryPressureReason) != c.reason {
			t.Fatalf("level %d: got %v %v", c.level, deref(s.MemoryPressure), deref(s.MemoryPressureReason))
		}
	}
}

func TestLoadReadFails(t *testing.T) {
	if _, err := Read(fakeReader{cores: 1, loadErr: errors.New("boom")}, now); err == nil || !strings.HasPrefix(err.Error(), "read_failed") {
		t.Fatalf("got %v", err)
	}
	if _, err := Read(fakeReader{cores: 1, loadErr: ErrUnsupportedOS}, now); !errors.Is(err, ErrUnsupportedOS) || !strings.HasPrefix(err.Error(), "unsupported_os") {
		t.Fatalf("got %v", err)
	}
}

func TestStatusJSONShape(t *testing.T) {
	s, _ := Read(fakeReader{cores: 8, load: [3]float64{20, 1, 0.5}, level: 1, memTotal: 1 << 30, vmstat: vmstat16k, swapUsed: 1 << 29, swapTotal: 1 << 31}, now)
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"kind":"MachineStatusReadback@1","schema_version":1,"observed_at":"2026-11-23T20:04:05Z","cores":8,"load":{"1m":20,"5m":1,"15m":0.5},"load_percent":{"1m":250,"5m":12.5,"15m":6.3},"memory_pressure":"normal","memory_pressure_reason":null,"memory_total_bytes":1073741824,"memory_used_bytes":22937600,"memory_used_percent":2.1,"memory_used_reason":null,"swap_used_bytes":536870912,"swap_total_bytes":2147483648,"swap_reason":null}`
	if string(b) != want {
		t.Fatalf("got  %s\nwant %s", b, want)
	}
}

func TestText(t *testing.T) {
	s, _ := Read(fakeReader{cores: 16, load: [3]float64{4.91, 9.34, 9.97}, level: 1, memTotal: 48 << 30, vmstat: vmstat16k, swapUsed: 5299 << 20, swapTotal: 6 << 30}, now)
	if got := s.Text(); got != "Load 31% (1m) 58% (5m) 62% (15m) of 16 cores; memory 0% used (0.0 of 48.0 GiB), swap 5.2 of 6.0 GiB; memory pressure normal" {
		t.Fatalf("got %q", got)
	}
	s, _ = Read(fakeReader{cores: 16, level: 3}, now)
	if got := s.Text(); !strings.HasSuffix(got, "; memory unknown (unknown_format), swap 0.0 of 0.0 GiB; memory pressure unknown (unknown_level)") {
		t.Fatalf("got %q", got)
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func TestMemoryUsed(t *testing.T) {
	s, err := Read(fakeReader{cores: 1, level: 1, memTotal: 48 << 30, vmstat: vmstat16k}, now)
	if err != nil || *s.MemoryUsedBytes != 1400*16384 || s.MemoryTotalBytes != 48<<30 || *s.MemoryUsedPercent != 0 || s.MemoryUsedReason != nil {
		t.Fatalf("got %+v %v", s, err)
	}
	half := strings.Replace(vmstat16k, "Anonymous pages:                           1000.", "Anonymous pages:                        1572464.", 1)
	s, _ = Read(fakeReader{cores: 1, level: 1, memTotal: 48 << 30, vmstat: half}, now)
	if *s.MemoryUsedPercent != 50 { // (1572464-100+300+200) pages of 16 KiB = 24 GiB of 48
		t.Fatalf("got %v", *s.MemoryUsedPercent)
	}
	for _, c := range []struct {
		vmstat string
		err    error
		reason string
	}{
		{strings.Replace(vmstat16k, "Pages wired down:                           300.\n", "", 1), nil, ReasonUnknownFormat},
		{strings.Replace(vmstat16k, "(page size of 16384 bytes)", "", 1), nil, ReasonUnknownFormat},
		{strings.Replace(vmstat16k, "(page size of 16384 bytes)", "(page size of 0 bytes)", 1), nil, ReasonUnknownFormat},
		{"", errors.New("boom"), ReasonReadFailed},
	} {
		s, err := Read(fakeReader{cores: 1, level: 1, vmstat: c.vmstat, vmstatErr: c.err}, now)
		if err != nil || s.MemoryUsedBytes != nil || s.MemoryUsedPercent != nil || deref(s.MemoryUsedReason) != c.reason || s.LoadPercent.M1 != 0 {
			t.Fatalf("%q: got %+v %v", c.reason, s, err)
		}
	}
	if _, err := Read(fakeReader{cores: 1, memTotalErr: errors.New("boom")}, now); err == nil || !strings.HasPrefix(err.Error(), "read_failed") {
		t.Fatalf("hw.memsize failure: %v", err)
	}
}

func TestSwap(t *testing.T) {
	s, _ := Read(fakeReader{cores: 1, swapErr: errors.New("boom")}, now)
	if s.SwapUsedBytes != nil || s.SwapTotalBytes != nil || deref(s.SwapReason) != ReasonReadFailed {
		t.Fatalf("got %+v", s)
	}
	raw := make([]byte, 32)
	raw[0], raw[8], raw[16] = 6, 1, 5 // total 6, avail 1, used 5 (little endian)
	used, total, err := decodeSwap(raw)
	if used != 5 || total != 6 || err != nil {
		t.Fatalf("got %d %d %v", used, total, err)
	}
	if _, _, err := decodeSwap(raw[:20]); err == nil {
		t.Fatal("accepted a short xsw_usage")
	}
}
