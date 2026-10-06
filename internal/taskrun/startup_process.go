package taskrun

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Read local identities only: foreground-job inspection alone misses stopped
// and background providers that still belong to the preserved terminal shell.
func systemStartupProcessTree(shellPID int) ([]StartupProcess, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "/bin/ps", "-axo", "pid=,ppid=,pgid=")
	command.WaitDelay = time.Second
	var output workflowLimitedOutput
	command.Stdout = &output
	command.Stderr = io.Discard
	if err := command.Run(); err != nil {
		return nil, err
	}
	if output.overflow {
		return nil, errors.New("local process identities exceeded the inspection limit")
	}
	return startupProcessTreeFromTable(shellPID, string(output.data))
}

func startupProcessTreeFromTable(shellPID int, table string) ([]StartupProcess, error) {
	if shellPID <= 0 {
		return nil, errors.New("preserved shell identity is unavailable")
	}
	processes := map[int]StartupProcess{}
	children := map[int][]int{}
	for _, line := range strings.Split(table, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if len(fields) != 3 {
			return nil, errors.New("local process inspection returned invalid identities")
		}
		var ids [3]int
		for index, field := range fields {
			number, err := strconv.Atoi(field)
			if err != nil || number < 0 {
				return nil, errors.New("local process inspection returned invalid identities")
			}
			ids[index] = number
		}
		if ids[0] == 0 || ids[0] == ids[1] {
			return nil, errors.New("local process inspection returned ambiguous identities")
		}
		if _, exists := processes[ids[0]]; exists {
			return nil, errors.New("local process inspection returned duplicate identities")
		}
		processes[ids[0]] = StartupProcess{PID: ids[0], ParentPID: ids[1], ProcessGroupID: ids[2]}
		children[ids[1]] = append(children[ids[1]], ids[0])
	}
	if _, exists := processes[shellPID]; !exists {
		return nil, errors.New("preserved shell was not observed in the local process inspection")
	}
	queue, seen := []int{shellPID}, map[int]bool{}
	var tree []StartupProcess
	for len(queue) > 0 {
		pid := queue[0]
		queue = queue[1:]
		if seen[pid] {
			return nil, errors.New("local process inspection returned an ambiguous shell tree")
		}
		seen[pid] = true
		tree = append(tree, processes[pid])
		queue = append(queue, children[pid]...)
	}
	sort.Slice(tree, func(i, j int) bool { return tree[i].PID < tree[j].PID })
	return tree, nil
}

func deliveryRecoveryProcessTree(d Dependencies, shellPID, groupID int) ([]StartupProcess, error) {
	observe := d.StartupProcessTree
	if observe == nil {
		observe = systemStartupProcessTree
	}
	tree, err := observe(shellPID)
	if err != nil {
		return nil, workflowError(4, "local process inspection unavailable: "+err.Error())
	}
	if len(tree) != 1 || tree[0].PID != shellPID || tree[0].ParentPID < 0 || tree[0].ParentPID == shellPID || tree[0].ProcessGroupID != groupID {
		return nil, workflowError(4, "provider exit requires the observed preserved shell with no remaining child processes, including background or suspended jobs")
	}
	return append([]StartupProcess(nil), tree...), nil
}
