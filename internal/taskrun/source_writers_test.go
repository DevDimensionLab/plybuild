package taskrun

import (
	"encoding/json"
	"os"
	"testing"
)

func TestAdoptionSourceWriterObservations(t *testing.T) {
	for _, tc := range []struct {
		name  string
		table []LocalProcess
		err   error
		state string
	}{
		{"idle", []LocalProcess{{PID: os.Getpid(), CWD: "/return"}}, nil, "absent"},
		{"writer", []LocalProcess{{PID: os.Getpid()}, {PID: 4242, CWD: "/source/src"}}, nil, "present"},
		{"refused", nil, os.ErrPermission, "unknown"},
		{"incomplete", []LocalProcess{{PID: 4242}}, nil, "unknown"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := Dependencies{LocalProcesses: func(session, transcript string) ([]LocalProcess, error) {
				if session != "" || transcript != "" {
					t.Fatal("invented unregistered session")
				}
				return tc.table, tc.err
			}}
			inventory := func(kind string) (json.RawMessage, error) {
				if kind == "pane" {
					return json.RawMessage(`{"panes":[]}`), nil
				}
				return json.RawMessage(`{"agents":[]}`), nil
			}
			got, e := ObserveSourceWriters(d, "/source", inventory)
			if got.State != tc.state || (e == nil) != (tc.state == "absent") {
				t.Fatalf("got %+v err=%v", got, e)
			}
		})
	}
	d := Dependencies{LocalProcesses: func(string, string) ([]LocalProcess, error) { t.Fatal("managed writer excluded"); return nil, nil }}
	got, e := ObserveSourceWriters(d, "/source", func(string) (json.RawMessage, error) {
		return json.RawMessage(`{"panes":[{"workspace_id":"w","tab_id":"t","pane_id":"p","terminal_id":"terminal","cwd":"/source"}]}`), nil
	})
	if e == nil || got.State != "present" {
		t.Fatal("known managed writer excluded")
	}
}

// r001-f02: valid managed IDs do not locate a pane or pending agent. A local
// idle table cannot replace missing managed location evidence.
func TestAdoptionManagedInventoryRequiresKnownLocation(t *testing.T) {
	for _, kind := range []string{"pane", "agent"} {
		for _, tc := range []struct {
			name, cwd, foreground, state string
			localCalls                   int
		}{
			{"complete-cwd", "/return", "", "absent", 1},
			{"complete-foreground", "", "/return", "absent", 1},
			{"source-cwd", "/source/src", "", "present", 0},
			{"source-foreground", "", "/source/src", "present", 0},
			{"missing-both-locations", "", "", "unknown", 0},
		} {
			t.Run(kind+"/"+tc.name, func(t *testing.T) {
				calls := 0
				d := Dependencies{LocalProcesses: func(session, transcript string) ([]LocalProcess, error) {
					calls++
					if session != "" || transcript != "" {
						t.Fatal("invented source session")
					}
					return []LocalProcess{{PID: os.Getpid(), CWD: "/return"}}, nil
				}}
				inventory := func(requested string) (json.RawMessage, error) {
					rows := []map[string]string{}
					if requested == kind {
						rows = append(rows, map[string]string{"workspace_id": "w", "tab_id": "t", "pane_id": "p", "terminal_id": "term", "cwd": tc.cwd, "foreground_cwd": tc.foreground})
					}
					b, e := json.Marshal(map[string]any{requested + "s": rows})
					return json.RawMessage(b), e
				}
				got, e := ObserveSourceWriters(d, "/source", inventory)
				if got.State != tc.state || (e == nil) != (tc.state == "absent") {
					t.Fatalf("managed location evidence classified as %+v err=%v; want %s", got, e, tc.state)
				}
				if calls != tc.localCalls {
					t.Fatalf("local inspection calls=%d want=%d", calls, tc.localCalls)
				}
			})
		}
	}
}
