// Package machine reads how hard this machine is working: load average per
// logical core and memory pressure. Nothing is stored; every read is fresh.
package machine

import (
	"errors"
	"fmt"
	"math"
	"time"
)

const StatusKind = "MachineStatusReadback@1"

// Reasons for a null value.
const (
	ReasonUnknownLevel = "unknown_level"
	ReasonReadFailed   = "read_failed"
	// ReasonUnknownFormat: vm_stat output lacks a line memory_used needs.
	ReasonUnknownFormat = "unknown_format"
)

var ErrUnsupportedOS = errors.New("unsupported_os: machine status reads macOS only")

// Reader is the system seam: the darwin reader and fakes in tests.
type Reader interface {
	Cores() int
	LoadAverage() ([3]float64, error)
	MemoryPressureLevel() (uint32, error)
	MemoryTotal() (uint64, error)
	VMStat() ([]byte, error)
	Swap() (used, total uint64, err error)
}

type Windows struct {
	M1  float64 `json:"1m"`
	M5  float64 `json:"5m"`
	M15 float64 `json:"15m"`
}

type Status struct {
	Kind                 string   `json:"kind"`
	SchemaVersion        int      `json:"schema_version"`
	ObservedAt           string   `json:"observed_at"`
	Cores                int      `json:"cores"`
	Load                 Windows  `json:"load"`
	LoadPercent          Windows  `json:"load_percent"`
	MemoryPressure       *string  `json:"memory_pressure"`
	MemoryPressureReason *string  `json:"memory_pressure_reason"`
	MemoryTotalBytes     uint64   `json:"memory_total_bytes"`
	MemoryUsedBytes      *uint64  `json:"memory_used_bytes"`
	MemoryUsedPercent    *float64 `json:"memory_used_percent"`
	MemoryUsedReason     *string  `json:"memory_used_reason"`
	SwapUsedBytes        *uint64  `json:"swap_used_bytes"`
	SwapTotalBytes       *uint64  `json:"swap_total_bytes"`
	SwapReason           *string  `json:"swap_reason"`
}

var pressureNames = map[uint32]string{1: "normal", 2: "warn", 4: "critical"}

// Read returns the status at now. A load average or physical memory size that
// cannot be read fails the call; any other value that cannot be read is null
// with a reason.
func Read(r Reader, now time.Time) (Status, error) {
	load, err := r.LoadAverage()
	if errors.Is(err, ErrUnsupportedOS) {
		return Status{}, err
	}
	if err != nil {
		return Status{}, fmt.Errorf("%s: load average: %w", ReasonReadFailed, err)
	}
	total, err := r.MemoryTotal()
	if err != nil {
		return Status{}, fmt.Errorf("%s: hw.memsize: %w", ReasonReadFailed, err)
	}
	cores := r.Cores()
	percent := func(l float64) float64 { return math.Round(l*1000/float64(cores)) / 10 }
	round2 := func(l float64) float64 { return math.Round(l*100) / 100 }
	s := Status{
		Kind: StatusKind, SchemaVersion: 1, ObservedAt: now.UTC().Format(time.RFC3339), Cores: cores,
		Load:             Windows{round2(load[0]), round2(load[1]), round2(load[2])},
		LoadPercent:      Windows{percent(load[0]), percent(load[1]), percent(load[2])},
		MemoryTotalBytes: total,
	}
	s.readMemoryUsed(r)
	s.readSwap(r)
	level, err := r.MemoryPressureLevel()
	name, known := pressureNames[level]
	switch {
	case err != nil:
		s.MemoryPressureReason = ptr(ReasonReadFailed)
	case !known:
		s.MemoryPressureReason = ptr(ReasonUnknownLevel)
	default:
		s.MemoryPressure = &name
	}
	return s, nil
}

// Text is the one-line human form.
func (s Status) Text() string {
	pressure := unknown(s.MemoryPressureReason)
	if s.MemoryPressure != nil {
		pressure = *s.MemoryPressure
	}
	p := s.LoadPercent
	return fmt.Sprintf("Load %.0f%% (1m) %.0f%% (5m) %.0f%% (15m) of %d cores; %s, %s; memory pressure %s", p.M1, p.M5, p.M15, s.Cores, s.memoryText(), s.swapText(), pressure)
}

const gib = 1 << 30

func (s Status) memoryText() string {
	if s.MemoryUsedPercent == nil {
		return "memory " + unknown(s.MemoryUsedReason)
	}
	return fmt.Sprintf("memory %.0f%% used (%.1f of %.1f GiB)", *s.MemoryUsedPercent, float64(*s.MemoryUsedBytes)/gib, float64(s.MemoryTotalBytes)/gib)
}

func (s Status) swapText() string {
	if s.SwapUsedBytes == nil {
		return "swap " + unknown(s.SwapReason)
	}
	return fmt.Sprintf("swap %.1f of %.1f GiB", float64(*s.SwapUsedBytes)/gib, float64(*s.SwapTotalBytes)/gib)
}

func unknown(reason *string) string {
	if reason == nil {
		return "unknown"
	}
	return "unknown (" + *reason + ")"
}

func ptr(s string) *string { return &s }

func (s *Status) readMemoryUsed(r Reader) {
	out, err := r.VMStat()
	if err != nil {
		s.MemoryUsedReason = ptr(ReasonReadFailed)
		return
	}
	used, ok := memoryUsed(out)
	if !ok {
		s.MemoryUsedReason = ptr(ReasonUnknownFormat)
		return
	}
	percent := math.Round(float64(used)*1000/float64(s.MemoryTotalBytes)) / 10
	s.MemoryUsedBytes, s.MemoryUsedPercent = &used, &percent
}

func (s *Status) readSwap(r Reader) {
	used, total, err := r.Swap()
	if err != nil {
		s.SwapReason = ptr(ReasonReadFailed)
		return
	}
	s.SwapUsedBytes, s.SwapTotalBytes = &used, &total
}
