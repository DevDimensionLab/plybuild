package taskexecute

import (
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestProjectDataRootPublicCLIStartupAndReadback(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		binary = filepath.Join(physicalTemp(t), "ply")
		build := exec.Command("go", "build", "-o", binary, "./cmd/ply")
		build.Dir = "../.."
		if output, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build CLI: %s %v", output, err)
		}
	}
	for _, name := range []string{"data-one", "data-two"} {
		t.Run(name, func(t *testing.T) {
			profile := projectDataProfile(t)
			_, in, root, epic := launcherFixtureWithCLI(t, "codex", binary)
			data := filepath.Join(root, name)
			setProjectDataRoot(t, profile, data)
			if err := os.WriteFile(in.Runtime.HerdrExecutable, []byte(recoveryLauncherHerdr), 0700); err != nil {
				t.Fatal(err)
			}
			call := func(cwd string, args ...string) []byte {
				t.Helper()
				command := exec.Command(binary, args...)
				command.Dir = cwd
				var stderr strings.Builder
				command.Stderr = &stderr
				output, err := command.Output()
				if err != nil {
					t.Fatalf("CLI %v: %v\n%s\n%s", args, err, output, stderr.String())
				}
				return output
			}
			var preview, first, repeated Result
			args := []string{"workflow", "execute", "--spec", "observable-goal", "--format", "json"}
			if err := json.Unmarshal(call(epic, append(args, "--check")...), &preview); err != nil || preview.Goal == nil {
				t.Fatalf("preview failed: %v", err)
			}
			if !strings.HasPrefix(preview.Goal.AcceptancePath, data+"/task-executions/") {
				t.Fatalf("configured root ignored: %s", preview.Goal.AcceptancePath)
			}
			if _, err := os.Stat(data); !os.IsNotExist(err) {
				t.Fatal("CLI check created data root")
			}
			if err := json.Unmarshal(call(epic, args...), &first); err != nil || first.Run == nil || first.State != "started" {
				t.Fatalf("startup failed: %+v %v", first, err)
			}
			setProjectDataRoot(t, profile, filepath.Join(root, "changed-root"))
			if err := json.Unmarshal(call(epic, args...), &repeated); err != nil || repeated.Run == nil || repeated.State != "existing" || repeated.Run.RunID != first.Run.RunID || repeated.RequestPath != first.RequestPath {
				t.Fatalf("repeat changed attempt: %+v %v", repeated, err)
			}
			var shown struct {
				RunID string `json:"run_id"`
			}
			if err := json.Unmarshal(call(root, "workflow", "execute", "show", first.Run.RunID, "--format", "json"), &shown); err != nil || shown.RunID != first.Run.RunID {
				t.Fatalf("readback from another cwd failed: %+v %v", shown, err)
			}
			var model map[string]any
			raw, err := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
			if err != nil || json.Unmarshal(raw, &model) != nil || model["tabs"] != float64(1) || model["starts"] != float64(1) || model["prompts"] != float64(1) {
				t.Fatalf("duplicate transport: %s %v", raw, err)
			}
			for _, absent := range []string{filepath.Join(root, "changed-root"), filepath.Join(root, ".ply", "task-executions")} {
				if _, err := os.Stat(absent); !os.IsNotExist(err) {
					t.Fatalf("unexpected attempt files: %s", absent)
				}
			}
			counts := map[string]int64{}
			for _, area := range []string{filepath.Join(root, ".ply"), data} {
				err := filepath.WalkDir(area, func(path string, entry fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if !entry.Type().IsRegular() {
						return nil
					}
					b, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					if utf8.Valid(b) && !strings.ContainsRune(string(b), 0) {
						counts["text_files"]++
						counts["text_bytes"] += int64(len(b))
						if strings.HasPrefix(path, filepath.Dir(preview.Goal.AcceptancePath)+string(filepath.Separator)) {
							counts["attempt_text_files"]++
							counts["attempt_text_bytes"] += int64(len(b))
						}
					}
					return nil
				})
				if err != nil {
					t.Fatal(err)
				}
			}
			t.Logf("completed startup/repeat/readback text: %v", counts)
		})
	}
}
