package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func claudeConfigFixture(t *testing.T, data string, mode os.FileMode) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".claude.json")
	if err = os.WriteFile(path, []byte(data), mode); err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
	return path
}

func claudeConfigRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func claudeConfigValue(t *testing.T, b []byte) any {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestClaudeConfigPreservesValuesAndCreatesPrivateBackup(t *testing.T) {
	key := `/physical/repo with "quotes" and \slash`
	for _, input := range []string{
		`{}`, `{"projects":{}}`,
		`{"projects":{"/other":{"hasTrustDialogAccepted":false,"unknown":[null,true,"keep"]}},"decimal":0.1234567890123456789,"large":9007199254740993123456789,"exponent":1e9999,"secret":"never put me in a receipt"}`,
	} {
		t.Run(input, func(t *testing.T) {
			path := claudeConfigFixture(t, input, 0640)
			before := claudeConfigRead(t, path)
			want := claudeConfigValue(t, before).(map[string]any)
			projects, ok := want["projects"].(map[string]any)
			if !ok {
				projects = map[string]any{}
				want["projects"] = projects
			}
			projects[key] = map[string]any{"hasTrustDialogAccepted": true}
			effect := workflowApplyClaudeTrust(Dependencies{}, path, key)
			if effect.State != "written" || effect.ConfigPath != path || effect.ProjectKey != key {
				t.Fatalf("trust write failed: %+v", effect)
			}
			after := claudeConfigRead(t, path)
			if !reflect.DeepEqual(want, claudeConfigValue(t, after)) {
				t.Fatal("unrelated JSON values changed")
			}
			if effect.BeforeSHA256 != hash(before) || effect.AfterSHA256 != hash(after) {
				t.Fatal("effect did not bind actual before and after bytes")
			}
			if effect.BackupPath == "" || !bytes.Equal(claudeConfigRead(t, effect.BackupPath), before) {
				t.Fatal("exact original backup missing")
			}
			info, _ := os.Stat(path)
			backup, _ := os.Stat(effect.BackupPath)
			if info.Mode().Perm() != 0640 || backup.Mode().Perm() != 0600 {
				t.Fatalf("configuration/backup modes differ: %v %v", info.Mode(), backup.Mode())
			}
			receipt, _ := json.Marshal(effect)
			if bytes.Contains(receipt, []byte("never put me")) {
				t.Fatal("receipt leaked configuration content")
			}
		})
	}
}

func TestClaudeConfigPreservesExistingProjectAndNoOp(t *testing.T) {
	key := "/physical/repo"
	path := claudeConfigFixture(t, `{"projects":{"/physical/repo":{"hasTrustDialogAccepted":false,"other":{"n":3.25}}}}`, 0600)
	effect := workflowApplyClaudeTrust(Dependencies{}, path, key)
	if effect.State != "written" {
		t.Fatalf("false flag not updated: %+v", effect)
	}
	want := claudeConfigValue(t, []byte(`{"projects":{"/physical/repo":{"hasTrustDialogAccepted":true,"other":{"n":3.25}}}}`))
	if !reflect.DeepEqual(want, claudeConfigValue(t, claudeConfigRead(t, path))) {
		t.Fatal("existing project fields changed")
	}
	before := claudeConfigRead(t, path)
	info, _ := os.Stat(path)
	entries, _ := os.ReadDir(filepath.Dir(path))
	effect = workflowApplyClaudeTrust(Dependencies{}, path, key)
	afterInfo, _ := os.Stat(path)
	afterEntries, _ := os.ReadDir(filepath.Dir(path))
	if effect.State != "already" || effect.BeforeSHA256 != hash(before) || effect.AfterSHA256 != hash(before) || effect.BackupPath != "" {
		t.Fatalf("already-trusted effect differs: %+v", effect)
	}
	if !bytes.Equal(before, claudeConfigRead(t, path)) || !os.SameFile(info, afterInfo) || !info.ModTime().Equal(afterInfo.ModTime()) || info.Mode() != afterInfo.Mode() || len(entries) != len(afterEntries) {
		t.Fatal("already-trusted configuration was touched")
	}
}

func TestClaudeConfigMissingAndInvalidStayUntouched(t *testing.T) {
	key := "/physical/repo"
	for _, input := range []string{
		`not json`, `null`, `[]`, `{"projects":null}`, `{"projects":[]}`,
		`{"projects":{"/physical/repo":null}}`, `{"projects":{"/physical/repo":[]}}`,
		`{"projects":{"/physical/repo":{"hasTrustDialogAccepted":"true"}}}`,
		`{"projects":{"/physical/repo":{"hasTrustDialogAccepted":null}}}`,
		`{"other":{"x":1,"x":2}}`, `{"projects":{},"projects":{}}`,
		`{"other":[{"x":1,"\u0078":2}]}`, `{"x":1} {"y":2}`,
		`{"\ud800":1}`, `{"unknown":"\udfff"}`, `{"unknown":"\ud800x"}`,
	} {
		t.Run(input, func(t *testing.T) {
			path := claudeConfigFixture(t, input, 0600)
			before, _ := os.Stat(path)
			effect := workflowApplyClaudeTrust(Dependencies{}, path, key)
			after, _ := os.Stat(path)
			entries, _ := os.ReadDir(filepath.Dir(path))
			if effect.State != "failed" || effect.Reason == "" || string(claudeConfigRead(t, path)) != input || !os.SameFile(before, after) || len(entries) != 1 {
				t.Fatalf("invalid configuration was changed: %+v", effect)
			}
		})
	}
	path := claudeConfigFixture(t, `{}`, 0600)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	effect := workflowApplyClaudeTrust(Dependencies{}, path, key)
	entries, _ := os.ReadDir(filepath.Dir(path))
	if effect.State != "skipped" || len(entries) != 0 || effect.BackupPath != "" {
		t.Fatalf("missing configuration created artifacts: %+v", effect)
	}
}

func TestClaudeConfigPreservesUnicodeAndReportsPostPublicationDrift(t *testing.T) {
	path := claudeConfigFixture(t, `{"\ud83d\ude00":"\ud83d\ude01","literal":"\\ud800"}`, 0600)
	effect := workflowApplyClaudeTrust(Dependencies{}, path, "/repo")
	if effect.State != "written" {
		t.Fatalf("valid escaped Unicode rejected: %+v", effect)
	}
	updated := []byte(`{"external":"after-publication"}`)
	d := Dependencies{Fault: func(point string) error {
		if point == "claude_trust_after_publish" {
			return os.WriteFile(path, updated, 0600)
		}
		return nil
	}}
	effect = workflowApplyClaudeTrust(d, path, "/another-repo")
	if effect.State != "unknown" || !bytes.Equal(updated, claudeConfigRead(t, path)) || effect.AfterSHA256 != hash(updated) {
		t.Fatalf("post-publication drift hidden or rolled back: %+v", effect)
	}
}

func TestClaudeConfigDoesNotPublishAnUnreadableOversizeResult(t *testing.T) {
	// JSON's HTML escaping expands each '<' to six bytes. The original fits
	// the read bound; its serialized replacement must also remain readable.
	input := `{"payload":"` + strings.Repeat("<", workflowClaudeConfigLimit/6) + `"}`
	path := claudeConfigFixture(t, input, 0600)
	effect := workflowApplyClaudeTrust(Dependencies{}, path, "/repo")
	entries, _ := os.ReadDir(filepath.Dir(path))
	if effect.State != "failed" || string(claudeConfigRead(t, path)) != input || len(entries) != 1 {
		t.Fatalf("oversize replacement was published: %+v", effect)
	}
}

func TestClaudeConfigRejectsSymlinkAndSpecialFile(t *testing.T) {
	for _, kind := range []string{"file_symlink", "parent_symlink", "directory", "lock_symlink"} {
		t.Run(kind, func(t *testing.T) {
			path := claudeConfigFixture(t, `{}`, 0600)
			original := path
			switch kind {
			case "file_symlink":
				path += ".link"
				if err := os.Symlink(original, path); err != nil {
					t.Fatal(err)
				}
			case "parent_symlink":
				link := filepath.Dir(path) + "-alias"
				if err := os.Symlink(filepath.Dir(path), link); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = os.Remove(link) })
				path = filepath.Join(link, filepath.Base(path))
			case "directory":
				path += ".dir"
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
			case "lock_symlink":
				if err := os.Symlink(original, path+".lock"); err != nil {
					t.Fatal(err)
				}
			}
			effect := workflowApplyClaudeTrust(Dependencies{}, path, "/physical/repo")
			if effect.State != "failed" || string(claudeConfigRead(t, original)) != `{}` {
				t.Fatalf("unsafe file accepted or original changed: %+v", effect)
			}
		})
	}
}

func TestClaudeConfigPreservesExistingNativeLock(t *testing.T) {
	path := claudeConfigFixture(t, `{"keep":1}`, 0600)
	lockPath := path + ".lock"
	if err := os.Mkdir(lockPath, 0700); err != nil {
		t.Fatal(err)
	}
	before, _ := os.Stat(lockPath)
	effect := workflowApplyClaudeTrust(Dependencies{}, path, "/repo")
	after, err := os.Stat(lockPath)
	if effect.State != "failed" || string(claudeConfigRead(t, path)) != `{"keep":1}` || err != nil || !os.SameFile(before, after) || !before.ModTime().Equal(after.ModTime()) {
		t.Fatalf("existing native lock was bypassed or changed: %+v", effect)
	}
}

func TestClaudeConfigReleasesOnlyItsOwnNativeLock(t *testing.T) {
	for _, drift := range []string{"none", "replace", "mtime"} {
		t.Run(drift, func(t *testing.T) {
			path := claudeConfigFixture(t, `{"keep":1}`, 0600)
			lockPath := path + ".lock"
			var foreign os.FileInfo
			d := Dependencies{Fault: func(point string) error {
				if point != "claude_trust_before_publish" || drift == "none" {
					return nil
				}
				if drift == "replace" {
					if err := os.Rename(lockPath, lockPath+".preserved"); err != nil {
						return err
					}
					if err := os.Mkdir(lockPath, 0700); err != nil {
						return err
					}
				} else {
					changed := time.Now().Add(time.Second)
					if err := os.Chtimes(lockPath, changed, changed); err != nil {
						return err
					}
				}
				foreign, _ = os.Stat(lockPath)
				return nil
			}}
			effect := workflowApplyClaudeTrust(d, path, "/repo")
			after, err := os.Stat(lockPath)
			if drift == "none" {
				if effect.State != "written" || !os.IsNotExist(err) {
					t.Fatalf("own lock was not released: %+v, %v", effect, err)
				}
				return
			}
			if effect.State != "failed" || string(claudeConfigRead(t, path)) != `{"keep":1}` || foreign == nil || err != nil || !os.SameFile(foreign, after) || !foreign.ModTime().Equal(after.ModTime()) {
				t.Fatalf("lost lock was ignored or a foreign lock was removed: %+v", effect)
			}
		})
	}
}

func TestClaudeConfigPreservesConcurrentPlyProjects(t *testing.T) {
	path := claudeConfigFixture(t, `{"keep":1.5}`, 0600)
	const count = 12
	results := make(chan ClaudeTrustEffect, count)
	var group sync.WaitGroup
	for i := 0; i < count; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			results <- workflowApplyClaudeTrust(Dependencies{}, path, "/repo/"+strings.Repeat("x", i+1))
		}(i)
	}
	group.Wait()
	close(results)
	for effect := range results {
		if effect.State != "written" {
			t.Fatalf("concurrent writer failed: %+v", effect)
		}
	}
	actual := claudeConfigValue(t, claudeConfigRead(t, path)).(map[string]any)
	if len(actual["projects"].(map[string]any)) != count || actual["keep"] != json.Number("1.5") {
		t.Fatal("concurrent writer lost a project or existing value")
	}
}

func TestClaudeConfigDetectsExternalDriftBeforePublish(t *testing.T) {
	path := claudeConfigFixture(t, `{"keep":1}`, 0600)
	updated := []byte(`{"keep":2,"external":"preserve"}`)
	d := Dependencies{Fault: func(point string) error {
		if point == "claude_trust_before_publish" {
			return os.WriteFile(path, updated, 0600)
		}
		return nil
	}}
	effect := workflowApplyClaudeTrust(d, path, "/physical/repo")
	if effect.State != "failed" || !bytes.Equal(updated, claudeConfigRead(t, path)) || effect.AfterSHA256 != hash(updated) {
		t.Fatalf("observed external change overwritten or not reported: %+v", effect)
	}
}

func TestClaudeConfigFaultsPreserveActualOutcome(t *testing.T) {
	for _, point := range []string{"claude_trust_before_backup", "claude_trust_before_temp", "claude_trust_before_publish", "claude_trust_after_publish", "claude_trust_before_directory_sync"} {
		t.Run(point, func(t *testing.T) {
			path := claudeConfigFixture(t, `{"keep":3}`, 0600)
			d := Dependencies{Fault: func(p string) error {
				if p == point {
					return errors.New("secret error content")
				}
				return nil
			}}
			effect := workflowApplyClaudeTrust(d, path, "/physical/repo")
			after := claudeConfigRead(t, path)
			postPublish := point == "claude_trust_after_publish" || point == "claude_trust_before_directory_sync"
			if postPublish {
				if effect.State != "unknown" || !bytes.Contains(after, []byte("hasTrustDialogAccepted")) {
					t.Fatalf("published effect hidden: %+v", effect)
				}
			} else if effect.State != "failed" || string(after) != `{"keep":3}` {
				t.Fatalf("pre-publication failure changed configuration: %+v", effect)
			}
			if effect.AfterSHA256 != hash(after) || strings.Contains(effect.Reason, "secret error") {
				t.Fatalf("incorrect or sensitive receipt: %+v", effect)
			}
			matches, _ := filepath.Glob(path + ".ply-trust-tmp-*")
			if len(matches) != 0 {
				t.Fatal("owned temp file leaked")
			}
		})
	}
}

func TestClaudeConfigAlreadyTrustedDoesNotEvenCreateLock(t *testing.T) {
	path := claudeConfigFixture(t, "{\n  \"projects\": {\"/repo\": {\"hasTrustDialogAccepted\": true}}\n}\n", 0644)
	past := time.Unix(1234567890, 1234)
	if err := os.Chtimes(path, past, past); err != nil {
		t.Fatal(err)
	}
	before, _ := os.Stat(path)
	effect := workflowApplyClaudeTrust(Dependencies{}, path, "/repo")
	after, _ := os.Stat(path)
	entries, _ := os.ReadDir(filepath.Dir(path))
	if effect.State != "already" || len(entries) != 1 || !os.SameFile(before, after) || !after.ModTime().Equal(past) || after.Mode().Perm() != 0644 {
		t.Fatalf("already-trusted path was mutated: %+v", effect)
	}
}

func TestClaudeConfigAlreadyRequiresCurrentParentAndSnapshot(t *testing.T) {
	const original = `{"projects":{"/repo":{"hasTrustDialogAccepted":true}}}`
	const replacement = `{"projects":{"/repo":{"hasTrustDialogAccepted":false}},"external":"preserve"}`
	for _, branch := range []string{"early", "locked"} {
		for _, drift := range []string{"parent", "file", "content"} {
			t.Run(branch+"/"+drift, func(t *testing.T) {
				path := claudeConfigFixture(t, original, 0640)
				calls := 0
				d := Dependencies{Fault: func(point string) error {
					if point != "claude_trust_before_already" {
						return nil
					}
					calls++
					switch drift {
					case "parent":
						parent := filepath.Dir(path)
						if err := os.Rename(parent, parent+".moved"); err != nil {
							return err
						}
						if err := os.Mkdir(parent, 0700); err != nil {
							return err
						}
						return os.WriteFile(path, []byte(replacement), 0600)
					case "file":
						if err := os.WriteFile(path+".replacement", []byte(replacement), 0600); err != nil {
							return err
						}
						return os.Rename(path+".replacement", path)
					default:
						return os.WriteFile(path, []byte(replacement), 0640)
					}
				}}
				var effect ClaudeTrustEffect
				if branch == "early" {
					effect = workflowApplyClaudeTrust(d, path, "/repo")
				} else {
					root, err := physicalRoot(filepath.Dir(path), false)
					if err != nil {
						t.Fatal(err)
					}
					defer root.Close()
					lock, err := workflowLockClaudeConfig(root, filepath.Base(path)+".lock")
					if err != nil {
						t.Fatal(err)
					}
					effect = workflowWriteClaudeTrust(d, root, lock, ClaudeTrustEffect{State: "failed", ConfigPath: path, ProjectKey: "/repo"})
					if err = lock.release(); err != nil {
						t.Fatal(err)
					}
				}
				entries, _ := os.ReadDir(filepath.Dir(path))
				if calls != 1 || effect.State != "failed" && effect.State != "unknown" || string(claudeConfigRead(t, path)) != replacement || len(entries) != 1 || effect.BackupPath != "" {
					t.Fatalf("stale no-op accepted or replacement changed: calls=%d effect=%+v", calls, effect)
				}
				if effect.AfterSHA256 != "" && effect.AfterSHA256 != hash([]byte(replacement)) {
					t.Fatalf("stale bytes reported as current poststate: %+v", effect)
				}
			})
		}
	}
}
