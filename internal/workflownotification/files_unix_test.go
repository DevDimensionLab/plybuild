//go:build !windows

package workflownotification

import (
	"os"
	"syscall"
	"testing"
)

type changedOwnerInfo struct {
	os.FileInfo
	stat syscall.Stat_t
}

func (i changedOwnerInfo) Sys() any { return &i.stat }
func TestPrivateStateRejectsDifferentOwner(t *testing.T) {
	f := newFixture(t)
	info, e := os.Stat(f.root)
	if e != nil {
		t.Fatal(e)
	}
	stat := *info.Sys().(*syscall.Stat_t)
	stat.Uid++
	if privateInfo(changedOwnerInfo{info, stat}) {
		t.Fatal("foreign-owned state accepted")
	}
}
