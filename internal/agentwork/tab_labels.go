package agentwork

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

const (
	tabLabelWorkspaceLimit = 32
	tabLabelParallelLimit  = 4
	tabLabelByteLimit      = 256 << 10
	tabLabelTabLimit       = 4096
	tabLabelRuneLimit      = 256
	tabLabelDeadline       = time.Second
	// CommandContext kills the direct child at the shared deadline. WaitDelay
	// closes inherited pipes after at most another 50ms, even if descendants
	// retain them. All workers share that deadline; cleanup is concurrent.
	tabLabelCleanup = 50 * time.Millisecond
)

type tabKey struct{ workspace, tab string }

// These presentation reads deliberately do not belong to Runtime.List/Get or
// Identity. Overview alone invokes them, after constructing the base result.
func (h HerdrRuntime) tabLabels(agents []Agent) map[tabKey]string {
	ctx, cancel := context.WithTimeout(context.Background(), tabLabelDeadline)
	defer cancel()
	labels := map[tabKey]string{}
	workspaces := map[string]bool{}
	for _, a := range agents {
		workspaces[a.Workspace] = true
	}
	ids := make([]string, 0, len(workspaces))
	for id := range workspaces {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	if len(ids) > tabLabelWorkspaceLimit {
		ids = ids[:tabLabelWorkspaceLimit]
	}
	if len(ids) == 0 {
		return labels
	}
	jobs := make(chan string, len(ids))
	for _, id := range ids {
		jobs <- id
	}
	close(jobs)
	var mu sync.Mutex
	var workers sync.WaitGroup
	for worker := 0; worker < tabLabelParallelLimit && worker < len(ids); worker++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for id := range jobs {
				if ctx.Err() != nil {
					return
				}
				values := h.workspaceTabLabels(ctx, id)
				mu.Lock()
				for key, value := range values {
					labels[key] = value
				}
				mu.Unlock()
			}
		}()
	}
	workers.Wait()
	return labels
}

type tabLabelOutput struct {
	data []byte
}

func (b *tabLabelOutput) Write(p []byte) (int, error) {
	left := tabLabelByteLimit - len(b.data)
	if len(p) > left {
		b.data = append(b.data, p[:left]...)
		return left, fmt.Errorf("tab_metadata_limit")
	}
	b.data = append(b.data, p...)
	return len(p), nil
}

func (h HerdrRuntime) workspaceTabLabels(ctx context.Context, workspace string) map[tabKey]string {
	c := exec.CommandContext(ctx, h.Executable, "tab", "list", "--workspace", workspace)
	c.WaitDelay = tabLabelCleanup
	var out, errout tabLabelOutput
	c.Stdout, c.Stderr = &out, &errout
	if c.Run() != nil || ctx.Err() != nil || !utf8.Valid(out.data) {
		return nil
	}
	labels := decodeTabLabels(ctx, out.data, workspace)
	if ctx.Err() != nil {
		return nil
	}
	return labels
}

func decodeTabLabels(ctx context.Context, data []byte, workspace string) map[tabKey]string {
	var env struct {
		Error  json.RawMessage `json:"error"`
		Result *struct {
			Tabs *[]json.RawMessage `json:"tabs"`
		} `json:"result"`
	}
	if json.Unmarshal(data, &env) != nil || env.Result == nil || env.Result.Tabs == nil ||
		len(env.Error) > 0 && string(env.Error) != "null" || len(*env.Result.Tabs) > tabLabelTabLimit {
		return nil
	}
	labels := map[tabKey]string{}
	seen := map[tabKey]bool{}
	for _, raw := range *env.Result.Tabs {
		if ctx.Err() != nil {
			return nil
		}
		// Decode the join fields separately so a malformed label on a duplicate
		// row still invalidates that tab's metadata, regardless of row order.
		var row struct {
			Workspace string          `json:"workspace_id"`
			Tab       string          `json:"tab_id"`
			Label     json.RawMessage `json:"label"`
		}
		if json.Unmarshal(raw, &row) != nil || row.Workspace != workspace || row.Tab == "" {
			continue
		}
		key := tabKey{workspace, row.Tab}
		if seen[key] {
			delete(labels, key)
			continue
		}
		seen[key] = true
		// The shared strict decoder rejects invalid Unicode rather than
		// repairing lone surrogate escapes into an invented replacement name.
		value, err := canonicaljson.DecodeStrict(row.Label)
		label, ok := value.(string)
		if err == nil && ok && validTabLabel(label) {
			labels[key] = label
		}
	}
	return labels
}

func validTabLabel(label string) bool {
	if !utf8.ValidString(label) || strings.TrimSpace(label) == "" || utf8.RuneCountInString(label) > tabLabelRuneLimit {
		return false
	}
	for _, r := range label {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
