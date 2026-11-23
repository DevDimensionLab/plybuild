package machine

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var vmStatPageSize = regexp.MustCompile(`\(page size of (\d+) bytes\)`)

// memoryUsed is Activity Monitor's "Memory Used" from vm_stat output: app
// memory (anonymous minus purgeable pages), wired pages and the pages the
// compressor occupies, times the page size. ok is false when a line is missing.
func memoryUsed(out []byte) (used uint64, ok bool) {
	size := vmStatPageSize.FindSubmatch(out)
	if size == nil {
		return 0, false
	}
	pageSize, err := strconv.ParseUint(string(size[1]), 10, 64)
	if err != nil || pageSize == 0 {
		return 0, false
	}
	pages := map[string]uint64{}
	lines := bufio.NewScanner(bytes.NewReader(out))
	for lines.Scan() {
		label, value, found := strings.Cut(lines.Text(), ":")
		count, err := strconv.ParseUint(strings.TrimSuffix(strings.TrimSpace(value), "."), 10, 64)
		if found && err == nil {
			pages[label] = count
		}
	}
	var sum int64
	for label, sign := range map[string]int64{"Anonymous pages": 1, "Pages purgeable": -1, "Pages wired down": 1, "Pages occupied by compressor": 1} {
		count, found := pages[label]
		if !found {
			return 0, false
		}
		sum += sign * int64(count)
	}
	if sum < 0 {
		return 0, false
	}
	return uint64(sum) * pageSize, true
}

// decodeSwap reads struct xsw_usage: total, avail and used as uint64.
func decodeSwap(raw []byte) (used, total uint64, err error) {
	if len(raw) < 24 {
		return 0, 0, fmt.Errorf("vm.swapusage has %d bytes, want at least 24", len(raw))
	}
	return binary.LittleEndian.Uint64(raw[16:24]), binary.LittleEndian.Uint64(raw[0:8]), nil
}
