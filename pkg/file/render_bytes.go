package file

import (
	"bytes"
	"text/template"
)

// RenderBytes renders a Go template held in memory. It performs no filesystem,
// Git or network access, returns parse and execute errors instead of panicking,
// and returns no output when rendering fails.
func RenderBytes(name string, source []byte, data interface{}) ([]byte, error) {
	t, err := template.New(name).Option("missingkey=error").Parse(string(source))
	if err != nil {
		return nil, err
	}

	var rendered bytes.Buffer
	if err := t.Execute(&rendered, data); err != nil {
		return nil, err
	}

	return rendered.Bytes(), nil
}
