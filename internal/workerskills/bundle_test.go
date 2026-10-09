package workerskills

import (
	"io/fs"
	"os"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"
)

func TestPinnedProvenanceAndLicense(t *testing.T) {
	m, err := Manifest()
	if err != nil {
		t.Fatal(err)
	}
	if m.SourceCommit != "1401c8b8030e023baeebb31781a6653fe8e93026" || m.SourceRepository != "https://github.com/addyosmani/agent-skills" || m.License != "MIT" || m.LicensePath != "upstream/LICENSE" {
		t.Fatalf("unexpected provenance: %+v", m)
	}
	license, err := embedded.ReadFile(m.LicensePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"MIT License", "Copyright (c) 2025 Addy Osmani", "Permission is hereby granted, free of charge", "THE SOFTWARE IS PROVIDED \"AS IS\""} {
		if !strings.Contains(string(license), expected) {
			t.Fatalf("license is missing %q", expected)
		}
	}
	upstream := 0
	for _, file := range m.Files {
		if file.SourcePath != "" {
			upstream++
			if file.Path != "upstream/"+file.SourcePath {
				t.Fatalf("source path was relocated without provenance: %+v", file)
			}
		} else if !strings.HasPrefix(file.Path, "adapters/") {
			t.Fatalf("unattributed non-adapter material: %+v", file)
		}
		if strings.HasSuffix(file.Path, "AGENTS.md") || strings.HasSuffix(file.Path, "CLAUDE.md") {
			t.Fatalf("upstream repository instructions were imported: %s", file.Path)
		}
	}
	if upstream != 10 {
		t.Fatalf("want five skills, four references and license; got %d source files", upstream)
	}
	if err := Verify(); err != nil {
		t.Fatal(err)
	}
}

func TestRoleBundlesContainExactSelectedBytesAndReferences(t *testing.T) {
	wantSkills := map[string][]string{
		RolePlanning:   {"planning-and-task-breakdown"},
		RoleDeveloper:  {"planning-and-task-breakdown", "incremental-implementation", "test-driven-development"},
		RoleCorrection: {"planning-and-task-breakdown", "incremental-implementation", "test-driven-development", "debugging-and-error-recovery"},
		RoleReviewer:   {"code-review-and-quality"},
	}
	m, err := Manifest()
	if err != nil {
		t.Fatal(err)
	}
	refPattern := regexp.MustCompile("`(\\.\\./\\.\\./references/[^`]+\\.md)`")
	for role, skills := range wantSkills {
		t.Run(role, func(t *testing.T) {
			bundle, err := Bundle(role)
			if err != nil {
				t.Fatal(err)
			}
			again, err := Bundle(role)
			if err != nil || bundle != again {
				t.Fatalf("role package is not deterministic: %v", err)
			}
			if !strings.Contains(bundle, "Role: "+role+"\n") || !strings.Contains(bundle, SourceCommit) || !strings.Contains(bundle, ManifestSHA256) {
				t.Fatal("role or provenance absent from actual prompt material")
			}
			common, _ := embedded.ReadFile("adapters/common.md")
			adapter, _ := embedded.ReadFile("adapters/" + role + ".md")
			upstreamStart := strings.Index(bundle, "## Upstream material")
			for _, instruction := range []string{string(common), string(adapter)} {
				at := strings.Index(bundle, instruction)
				if at < 0 || at >= upstreamStart {
					t.Fatal("Ply adapter is missing or follows upstream guidance")
				}
			}
			selection := m.Roles[role]
			var gotSkills []string
			for _, name := range selection.Files {
				raw, err := embedded.ReadFile(name)
				if err != nil {
					t.Fatal(err)
				}
				section := "--- BEGIN EMBEDDED FILE: " + name + " ---\n" + string(raw) + "\n--- END EMBEDDED FILE: " + name + " ---"
				if strings.Count(bundle, section) != 1 {
					t.Fatalf("exact source bytes must occur once: %s", name)
				}
				if strings.HasSuffix(name, "/SKILL.md") {
					gotSkills = append(gotSkills, path.Base(path.Dir(name)))
				}
				for _, match := range refPattern.FindAllSubmatch(raw, -1) {
					reference := path.Clean(path.Join(path.Dir(name), string(match[1])))
					if !contains(selection.Files, reference) {
						t.Fatalf("%s requires missing reference %s", name, reference)
					}
				}
			}
			if !reflect.DeepEqual(gotSkills, skills) {
				t.Fatalf("wrong role selection: got %v want %v", gotSkills, skills)
			}
			for _, file := range m.Files {
				if file.SourcePath != "" && !contains(selection.Files, file.Path) && strings.Contains(bundle, "--- BEGIN EMBEDDED FILE: "+file.Path+" ---") {
					t.Fatalf("unselected upstream instructions injected: %s", file.Path)
				}
			}
			t.Logf("%s: %d bytes, SHA-256 %s", role, len(bundle), digest([]byte(bundle)))
		})
	}
}

func TestUnknownRoleFailsWithoutPartialPrompt(t *testing.T) {
	for _, role := range []string{"", "Developer", "developer\nreviewer", "../../adapters/common.md", "shipper"} {
		if text, err := Bundle(role); err == nil || text != "" {
			t.Errorf("unknown role %q produced instructions: error=%v", role, err)
		}
	}
}

func TestBrokenAssetsFailBeforeProducingInstructions(t *testing.T) {
	checks := []struct {
		name string
		edit func(fstest.MapFS)
	}{
		{"source bytes", func(files fstest.MapFS) { files["upstream/skills/test-driven-development/SKILL.md"].Data[0] ^= 1 }},
		{"license bytes", func(files fstest.MapFS) { files["upstream/LICENSE"].Data[0] ^= 1 }},
		{"adapter bytes", func(files fstest.MapFS) { files["adapters/common.md"].Data[0] ^= 1 }},
		{"manifest role or hash change", func(files fstest.MapFS) { files["manifest.json"].Data = []byte(`{"roles":{"developer":{"files":[]}}}`) }},
		{"missing reference", func(files fstest.MapFS) { delete(files, "upstream/references/testing-patterns.md") }},
		{"missing manifest", func(files fstest.MapFS) { delete(files, "manifest.json") }},
		{"extra repository instructions", func(files fstest.MapFS) {
			files["upstream/AGENTS.md"] = &fstest.MapFile{Data: []byte("unreviewed instructions")}
		}},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			files := copyAssets(t)
			check.edit(files)
			if text, err := bundle(files, RoleDeveloper); err == nil || text != "" {
				t.Fatalf("broken package produced prompt: error=%v", err)
			}
		})
	}
}

func TestManifestCallerCannotChangeLaterBundles(t *testing.T) {
	original, err := Bundle(RoleDeveloper)
	if err != nil {
		t.Fatal(err)
	}
	m, err := Manifest()
	if err != nil {
		t.Fatal(err)
	}
	m.Files[0].SHA256 = "changed"
	m.Roles[RoleDeveloper].Files[0] = "unreviewed.md"
	delete(m.Roles, RoleReviewer)
	after, err := Bundle(RoleDeveloper)
	if err != nil || original != after {
		t.Fatalf("caller changed shared manifest state: %v", err)
	}
	if _, err := Bundle(RoleReviewer); err != nil {
		t.Fatal(err)
	}
}

// Re-executing the compiled tests from an empty directory checks that assets
// are in the executable, independent of the caller's working directory/config.
func TestBundleFromUnrelatedWorkingDirectory(t *testing.T) {
	if os.Getenv("PLY_SKILLS_OFFLINE_HELPER") == "1" {
		for _, role := range []string{RolePlanning, RoleDeveloper, RoleCorrection, RoleReviewer} {
			if _, err := Bundle(role); err != nil {
				t.Fatal(err)
			}
		}
		return
	}
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(binary, "-test.run=^TestBundleFromUnrelatedWorkingDirectory$")
	cmd.Dir = t.TempDir()
	cmd.Env = []string{"PLY_SKILLS_OFFLINE_HELPER=1"}
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("standalone embedded package failed: %v\n%s", err, output)
	}
}

func copyAssets(t *testing.T) fstest.MapFS {
	t.Helper()
	files := fstest.MapFS{}
	if err := fs.WalkDir(embedded, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		raw, err := embedded.ReadFile(name)
		if err == nil {
			files[name] = &fstest.MapFile{Data: append([]byte{}, raw...), Mode: 0444}
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
	return files
}

func contains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}
