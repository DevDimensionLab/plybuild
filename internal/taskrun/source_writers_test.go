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
