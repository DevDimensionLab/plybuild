package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTaskCloseoutSharesImmutableContentWithoutCopying(t *testing.T) {
	f, in := closeoutFixture(t, true, false)
	p, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !p.Ready {
		t.Fatalf("preview: %+v %v", p, err)
	}
	receipt, err := ApplyTaskCloseout(f.dependencies, in, p.PlanSHA256)
	if err != nil {
		t.Fatal(err)
	}
	var shared int
	var sharedBytes int64
	for _, e := range receipt.Retention.Entries {
		managed := false
		for _, kind := range []string{"objects", "manifests", "requests"} {
			managed = managed || e.OriginLocator == taskContentPath(f.root, kind, e.SHA256)
		}
		if !managed {
			continue
		}
		origin, err := os.Stat(e.OriginLocator)
		if err != nil {
			t.Fatal(err)
		}
		retained, err := os.Stat(e.RetainedLocator)
		if err != nil {
			t.Fatal(err)
		}
		if !os.SameFile(origin, retained) {
			t.Fatalf("closeout wrote another copy of immutable content: %s (%d bytes)", e.OriginLocator, e.SizeBytes)
		}
		shared++
		sharedBytes += e.SizeBytes
	}
	if shared == 0 {
		t.Fatal("fixture did not exercise immutable input retention")
	}
	t.Logf("closeout retained %d immutable references without copying %d bytes", shared, sharedBytes)
	if _, err := os.Stat(f.taskPath); !os.IsNotExist(err) {
		t.Fatal("completed source worktree was not removed")
	}
	if err := verifyCloseoutRetention(f.dependencies, f.root, receipt.Plan, *receipt.Retention); err != nil {
		t.Fatal(err)
	}
	for _, e := range receipt.Retention.Entries {
		if e.Role == "spec_manifest" {
			if err := os.Remove(e.OriginLocator); err != nil {
				t.Fatal(err)
			}
			b, err := ReadTaskRetainedEvidence(f.root, in.TaskID, e.OriginLocator, e.SHA256)
			if err != nil || int64(len(b)) != e.SizeBytes {
				t.Fatalf("retention lost bytes after original name was removed: %v", err)
			}
			return
		}
	}
	t.Fatal("fixture did not exercise manifest resolution")
}

func TestCloseoutSharingPreservesCopiesAndRejectsDamagedContent(t *testing.T) {
	for _, scenario := range []string{"objects", "manifests", "requests", "mutable", "foreign-workspace", "unsupported-link", "existing", "wrong-target", "writable", "symlink", "tampered"} {
		t.Run(scenario, func(t *testing.T) {
			root, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(filepath.Join(root, MarkerDirectory), 0700); err != nil {
				t.Fatal(err)
			}
			s := &TaskContentStorage{}
			b := []byte(strings.Repeat("immutable fixture content\n", 4096))
			kind := "objects"
			if scenario == "manifests" || scenario == "requests" {
				kind = scenario
			}
			sha, err := s.Publish(root, kind, b)
			if err != nil {
				t.Fatal(err)
			}
			origin := taskContentPath(root, kind, sha)
			dest := taskContentPath(root, "closeout", sha)
			wantShared, wantError := true, false
			switch scenario {
			case "mutable":
				origin = filepath.Join(root, "source.txt")
				if err := os.WriteFile(origin, b, 0600); err != nil {
					t.Fatal(err)
				}
				wantShared = false
			case "foreign-workspace":
				other, err := filepath.EvalSymlinks(t.TempDir())
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Mkdir(filepath.Join(other, MarkerDirectory), 0700); err != nil {
					t.Fatal(err)
				}
				if _, err := s.Publish(other, kind, b); err != nil {
					t.Fatal(err)
				}
				origin, wantShared = taskContentPath(other, kind, sha), false
			case "unsupported-link":
				s.fault = func(stage string) error {
					if stage == "object-link" {
						return errors.New("fixture filesystem has no hardlinks")
					}
					return nil
				}
				wantShared = false
			case "existing", "wrong-target":
				if _, err := s.Publish(root, "closeout", b); err != nil {
					t.Fatal(err)
				}
				wantShared = false
				if scenario == "wrong-target" {
					if err := os.Chmod(dest, 0600); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(dest, []byte("wrong"), 0400); err != nil {
						t.Fatal(err)
					}
					if err := os.Chmod(dest, 0400); err != nil {
						t.Fatal(err)
					}
					wantError = true
				}
			case "writable", "tampered":
				if err := os.Chmod(origin, 0600); err != nil {
					t.Fatal(err)
				}
				if scenario == "tampered" {
					if err := os.WriteFile(origin, []byte("changed"), 0400); err != nil {
						t.Fatal(err)
					}
					if err := os.Chmod(origin, 0400); err != nil {
						t.Fatal(err)
					}
				}
				wantError = true
			case "symlink":
				if err := os.Rename(origin, origin+".real"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(origin+".real", origin); err != nil {
					t.Fatal(err)
				}
				wantError = true
			}
			got, err := s.retainCloseout(root, origin, b)
			if wantError {
				if err == nil {
					t.Fatal("damaged or mutable managed content was shared")
				}
				if scenario == "wrong-target" {
					raw, err := os.ReadFile(dest)
					if err != nil || string(raw) != "wrong" {
						t.Fatal("wrong existing target was overwritten")
					}
				}
				return
			}
			if err != nil || got != sha {
				t.Fatalf("retention failed: %s %v", got, err)
			}
			a, err := os.Stat(origin)
			if err != nil {
				t.Fatal(err)
			}
			c, err := os.Stat(dest)
			if err != nil || os.SameFile(a, c) != wantShared {
				t.Fatalf("wrong storage relationship, shared=%v err=%v", wantShared, err)
			}
			if c.Mode().Perm() != 0400 {
				t.Fatal("retained object is not immutable")
			}
			if _, err := s.retainCloseout(root, origin, b); err != nil {
				t.Fatalf("retry failed: %v", err)
			}
		})
	}
}

func TestCloseoutSharingSyncFailureKeepsSourceAndResumes(t *testing.T) {
	f, in := closeoutFixture(t, true, false)
	p, err := PreviewTaskCloseout(f.dependencies, in)
	if err != nil || !p.Ready {
		t.Fatalf("preview: %+v %v", p, err)
	}
	syncs := 0
	f.dependencies.TaskContent.fault = func(stage string) error {
		if stage == "recovery-sync" {
			syncs++
			if syncs == 2 {
				return errors.New("fixture target sync failed after link")
			}
		}
		return nil
	}
	if _, err := ApplyTaskCloseout(f.dependencies, in, p.PlanSHA256); err == nil {
		t.Fatal("incomplete sync was reported as complete retention")
	}
	if _, err := os.Stat(f.taskPath); err != nil {
		t.Fatal("source was removed before durable retention")
	}
	f.dependencies.TaskContent.fault = nil
	receipt, err := ApplyTaskCloseout(f.dependencies, in, p.PlanSHA256)
	if err != nil || receipt.State != "complete" {
		t.Fatalf("retry did not recover the same closeout: %+v %v", receipt, err)
	}
}
