//go:build !windows

package taskrun

import (
	"os"
	"testing"
)

func TestPlyUpgradeWorkerControl(t *testing.T) {
	f, _ := mr05Worker(t)
	f.own()
	p := mr05Preview(t, f, "pause", "tc0e6e1fe-worker-upgrade")
	if err := os.Remove(f.s.Request.Runtime.PlyExecutable.Path); err != nil {
		t.Fatal(err)
	}
	if _, err := WorkerControls(f.d).Apply(p); err != nil {
		t.Fatalf("Ply upgrade invalidated worker control: %v", err)
	}
	mr05Consume(t, f)
	if f.s.Result.Worker.Phase != "paused" {
		t.Fatalf("pause was not consumed: %+v", f.s.Result.Worker)
	}
}
