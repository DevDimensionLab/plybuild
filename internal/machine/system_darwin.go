package machine

import (
	"encoding/binary"
	"fmt"
	"os/exec"
	"runtime"

	"golang.org/x/sys/unix"
)

type system struct{}

// System reads this machine through sysctl.
func System() Reader { return system{} }

func (system) Cores() int { return runtime.NumCPU() }

// LoadAverage decodes struct loadavg: three uint32 fixed-point values,
// padding, then the int64 scale they are divided by.
func (system) LoadAverage() ([3]float64, error) {
	var load [3]float64
	raw, err := unix.SysctlRaw("vm.loadavg")
	if err != nil {
		return load, err
	}
	if len(raw) < 24 {
		return load, fmt.Errorf("vm.loadavg has %d bytes, want 24", len(raw))
	}
	scale := int64(binary.LittleEndian.Uint64(raw[16:24]))
	if scale <= 0 {
		return load, fmt.Errorf("vm.loadavg scale %d", scale)
	}
	for i := range load {
		load[i] = float64(binary.LittleEndian.Uint32(raw[4*i:])) / float64(scale)
	}
	return load, nil
}

func (system) MemoryPressureLevel() (uint32, error) {
	return unix.SysctlUint32("kern.memorystatus_vm_pressure_level")
}

// Processes runs ps for every user's processes; the ps process itself is left out.
func (system) Processes() ([]Process, error) {
	ps := exec.Command("/bin/ps", "-Ao", "pid=,pcpu=,rss=,comm=")
	out, err := ps.Output()
	if err != nil {
		return nil, err
	}
	return parseProcesses(out, ps.Process.Pid), nil
}

func (system) MemoryTotal() (uint64, error) { return unix.SysctlUint64("hw.memsize") }

func (system) VMStat() ([]byte, error) { return exec.Command("/usr/bin/vm_stat").Output() }

func (system) Swap() (used, total uint64, err error) {
	raw, err := unix.SysctlRaw("vm.swapusage")
	if err != nil {
		return 0, 0, err
	}
	return decodeSwap(raw)
}
