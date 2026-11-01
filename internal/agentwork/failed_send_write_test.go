package agentwork

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFollowupFailedSendKeepsClaimWithoutAnotherSnapshot(t *testing.T) {
	s, d, calls := fixture(t)
	d.NextCheckAtUTC = s.now()
	f, err := s.Register(d)
	if err != nil {
		t.Fatal(err)
	}
	snapshots := func() map[string]string {
		t.Helper()
		paths, err := filepath.Glob(filepath.Join(s.dir(), "revisions", "*.json"))
		if err != nil {
			t.Fatal(err)
		}
		result := map[string]string{}
		for _, path := range paths {
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			result[filepath.Base(path)] = string(b)
		}
		return result
	}
	before := snapshots()
	var claim Followup
	var atSend map[string]string
	s.Send = func(Identity, string) error {
		calls.Add(1)
		var err error
		claim, err = s.Show(f.ID)
		if err != nil || claim.Transport.State != "unknown" {
			t.Fatalf("claim not durable before transport: %v %v", claim.Transport, err)
		}
		atSend = snapshots()
		return errors.New("input may have arrived but acknowledgment was lost")
	}
	got, err := s.Dispatch(f.ID)
	if err != nil || got.Transport.State != "unknown" || calls.Load() != 1 {
		t.Fatalf("unknown outcome lost: %v %v %d", got.Transport, err, calls.Load())
	}
	after := snapshots()
	addedBytes := 0
	for name, b := range after {
		if _, ok := before[name]; !ok {
			addedBytes += len(b)
		}
	}
	t.Logf("dispatch published snapshots=%d bytes=%d", len(after)-len(before), addedBytes)
	if !reflect.DeepEqual(atSend, after) {
		t.Fatal("failed send wrote another full snapshot after the durable unknown claim")
	}
	if got.Revision != claim.Revision {
		t.Fatal("returned revision differs from durable claim")
	}
	restarted := *s
	reread, err := restarted.Show(f.ID)
	if err != nil || reread.Transport.State != "unknown" || reread.Revision != claim.Revision {
		t.Fatalf("readback lost claim: %v %v", reread, err)
	}
	if _, err := restarted.Dispatch(f.ID); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 1 || !reflect.DeepEqual(after, snapshots()) {
		t.Fatal("retry repeated transport or storage")
	}
}
