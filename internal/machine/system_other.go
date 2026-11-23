//go:build !darwin

package machine

import "runtime"

type system struct{}

// System reads nothing outside macOS; Read reports ErrUnsupportedOS.
func System() Reader { return system{} }

func (system) Cores() int                           { return runtime.NumCPU() }
func (system) LoadAverage() ([3]float64, error)     { return [3]float64{}, ErrUnsupportedOS }
func (system) MemoryPressureLevel() (uint32, error) { return 0, ErrUnsupportedOS }
func (system) Processes() ([]Process, error)        { return nil, ErrUnsupportedOS }
func (system) MemoryTotal() (uint64, error)         { return 0, ErrUnsupportedOS }
func (system) VMStat() ([]byte, error)              { return nil, ErrUnsupportedOS }
func (system) Swap() (uint64, uint64, error)        { return 0, 0, ErrUnsupportedOS }
