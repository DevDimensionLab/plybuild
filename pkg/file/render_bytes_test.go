package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderBytesRendersTemplateInMemory(t *testing.T) {
	got, err := RenderBytes("greeting", []byte("hello {{.Name}}\n"), struct{ Name string }{Name: "world `x` \"y\""})

	if err != nil {
		t.Fatalf("RenderBytes returned %v, want nil", err)
	}
	if want := "hello world `x` \"y\"\n"; string(got) != want {
		t.Fatalf("RenderBytes = %q, want %q", got, want)
	}
}

func TestRenderBytesDoesNotEvaluateDataAsTemplate(t *testing.T) {
	got, err := RenderBytes("data", []byte("{{.Value}}"), struct{ Value string }{Value: "{{.Missing}} $(id)"})

	if err != nil {
		t.Fatalf("RenderBytes returned %v, want nil", err)
	}
	if string(got) != "{{.Missing}} $(id)" {
		t.Fatalf("RenderBytes evaluated data: %q", got)
	}
}

func TestRenderBytesIsDeterministic(t *testing.T) {
	source := []byte("{{range .Items}}[{{.}}]{{end}}")
	data := struct{ Items []string }{Items: []string{"b", "a", "c"}}

	first, err := RenderBytes("list", source, data)
	if err != nil {
		t.Fatal(err)
	}
	second, err := RenderBytes("list", source, data)
	if err != nil {
		t.Fatal(err)
	}

	if string(first) != "[b][a][c]" || string(first) != string(second) {
		t.Fatalf("RenderBytes not deterministic: %q then %q", first, second)
	}
}

func TestRenderBytesReturnsParseErrorInsteadOfPanicking(t *testing.T) {
	var got []byte
	var err error
	var recovered interface{}

	func() {
		defer func() { recovered = recover() }()
		got, err = RenderBytes("broken.tmpl", []byte("before {{"), struct{}{})
	}()

	if recovered != nil {
		t.Fatalf("RenderBytes panicked: %v", recovered)
	}
	if err == nil || !strings.Contains(err.Error(), "broken.tmpl") {
		t.Fatalf("RenderBytes parse error = %v, want error naming the template", err)
	}
	if got != nil {
		t.Fatalf("RenderBytes returned partial output %q on parse error", got)
	}
}

func TestRenderBytesReturnsExecuteErrorWithoutPartialOutput(t *testing.T) {
	got, err := RenderBytes("exec.tmpl", []byte("written before {{.Nope}}"), struct{ Value string }{})

	if err == nil {
		t.Fatal("RenderBytes accepted a template referencing a missing field")
	}
	if got != nil {
		t.Fatalf("RenderBytes returned partial output %q on execute error", got)
	}
}

func TestRenderBytesRejectsMissingMapKeys(t *testing.T) {
	got, err := RenderBytes("map.tmpl", []byte("{{.absent}}"), map[string]string{})

	if err == nil {
		t.Fatalf("RenderBytes rendered a missing map key as %q, want an error", got)
	}
}

func TestRenderStillWritesValidTemplateOutput(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "input path with spaces.tmpl")
	output := filepath.Join(dir, "output.txt")
	if err := os.WriteFile(input, []byte("artifact={{.ArtifactID}}\ngroup={{.GroupID}}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := Render(input, output, struct {
		ArtifactID string
		GroupID    string
	}{ArtifactID: "demo", GroupID: "com.example"})

	if err != nil {
		t.Fatalf("Render returned %v, want nil", err)
	}
	written, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if want := "artifact=demo\ngroup=com.example\n"; string(written) != want {
		t.Fatalf("Render wrote %q, want %q", written, want)
	}
}
