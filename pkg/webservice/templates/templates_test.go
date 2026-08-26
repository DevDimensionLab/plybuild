package templates

import (
	"crypto/sha256"
	"encoding/hex"
	htmltemplate "html/template"
	"strings"
	"testing"
	texttemplate "text/template"
)

type exactTemplateString string

const (
	_ exactTemplateString = Generate
	_ exactTemplateString = Upgrade
)

func TestTemplatesPreserveExactCompleteBytesCompositionAndBoundaries(t *testing.T) {
	fragments := []struct {
		name       string
		contents   string
		wantLength int
		wantSHA256 string
		wantPrefix string
		wantSuffix string
	}{
		{
			name: "shared header", contents: header,
			wantLength: 462, wantSHA256: "02a2bec76d462a8e2eba236f94e3afaaa764a259d74b399709b04e2d55c08327",
			wantPrefix: "\n<html>\n<head>\n", wantSuffix: "</head>\n<body>\n",
		},
		{
			name: "generate body", contents: generate,
			wantLength: 2420, wantSHA256: "5be666b2ee4a971ed7dfe6fb06c08dd4522589aefd617e8ea34e5993a5aba851",
			wantPrefix: "\n<div class=\"container\">\n", wantSuffix: "</div>\n",
		},
		{
			name: "upgrade body", contents: upgrade,
			wantLength: 309, wantSHA256: "c096570c516c365f8e3f41119d7034e62244e61e1fa8a07c3ccf4e998952c664",
			wantPrefix: "\n<div class=\"container\">\n", wantSuffix: "</div>\n",
		},
		{
			name: "shared footer", contents: footer,
			wantLength: 17, wantSHA256: "1f3a0cf746539d789ae2368f24bc649b2041b6298ca83996840be70a367add4d",
			wantPrefix: "\n</body>\n", wantSuffix: "</html>\n",
		},
	}
	if len(fragments) == 0 {
		t.Fatal("private template fragment characterization population is empty")
	}

	for _, fragment := range fragments {
		t.Run(fragment.name, func(t *testing.T) {
			assertExactTemplateBytes(t, fragment.contents, fragment.wantLength, fragment.wantSHA256)
			if !strings.HasPrefix(fragment.contents, fragment.wantPrefix) ||
				!strings.HasSuffix(fragment.contents, fragment.wantSuffix) {
				t.Fatalf("%s no longer has exact leading %q and trailing %q bytes",
					fragment.name, fragment.wantPrefix, fragment.wantSuffix)
			}
		})
	}

	publicTemplates := []struct {
		name       string
		contents   string
		body       string
		wantLength int
		wantSHA256 string
	}{
		{
			name: "Generate", contents: Generate, body: generate,
			wantLength: 2899, wantSHA256: "d3507d7c18f5a6d3fa0bcee1ef4465cfa8af2e59c5ccccb45b6174be774a23a8",
		},
		{
			name: "Upgrade", contents: Upgrade, body: upgrade,
			wantLength: 788, wantSHA256: "ae9194d9a6734d1fe568692e11327fc1b31b3f0d797fd558fd6ceba3ffa58ecc",
		},
	}
	if len(publicTemplates) == 0 {
		t.Fatal("public template characterization population is empty")
	}

	for _, template := range publicTemplates {
		t.Run(template.name, func(t *testing.T) {
			wantComposition := header + template.body + footer
			if template.contents != wantComposition {
				t.Fatalf("%s is not the exact shared header + private body + shared footer composition", template.name)
			}
			if !strings.HasPrefix(template.contents, "\n<html>\n<head>\n") ||
				!strings.HasSuffix(template.contents, "\n</body>\n</html>\n") {
				t.Fatalf("%s public leading or trailing bytes changed", template.name)
			}
			assertExactTemplateBytes(t, template.contents, template.wantLength, template.wantSHA256)
		})
	}
}

func TestTemplatesPreserveExactFormsInputsActionsAndCDNReferences(t *testing.T) {
	tokens := []struct {
		name      string
		contents  string
		token     string
		wantCount int
	}{
		{name: "generate form action", contents: Generate, token: `<form class="form-inline" action="/api/generate" method="POST">`, wantCount: 1},
		{name: "upgrade form action", contents: Upgrade, token: `<form class="form-inline" action="/api/upgrade" method="POST">`, wantCount: 1},
		{name: "groupId input name", contents: Generate, token: `name="groupId"`, wantCount: 1},
		{name: "artifactId input name", contents: Generate, token: `name="artifactId"`, wantCount: 1},
		{name: "package input name", contents: Generate, token: `name="package"`, wantCount: 1},
		{name: "name input name", contents: Generate, token: `name="name"`, wantCount: 1},
		{name: "description input name", contents: Generate, token: `name="description"`, wantCount: 1},
		{name: "language input name", contents: Generate, token: `name="language"`, wantCount: 2},
		{name: "templates input name", contents: Generate, token: `name="templates"`, wantCount: 1},
		{name: "dependencies input name", contents: Generate, token: `name="dependencies"`, wantCount: 1},
		{name: "project group action", contents: Generate, token: `{{.ProjectConfig.GroupId}}`, wantCount: 1},
		{name: "project artifact action", contents: Generate, token: `{{.ProjectConfig.ArtifactId}}`, wantCount: 1},
		{name: "project package action", contents: Generate, token: `{{.ProjectConfig.Package}}`, wantCount: 1},
		{name: "project name action", contents: Generate, token: `{{.ProjectConfig.Name}}`, wantCount: 1},
		{name: "project description action", contents: Generate, token: `{{.ProjectConfig.Description}}`, wantCount: 1},
		{name: "cloud templates range action", contents: Generate, token: `{{range .CloudConfig.Templates }}`, wantCount: 1},
		{name: "dependency groups range action", contents: Generate, token: `{{range .IoResponse.Dependencies.Values }}`, wantCount: 1},
		{name: "dependency values range action", contents: Generate, token: `{{range .Values }}`, wantCount: 1},
		{name: "name actions", contents: Generate, token: `{{ .Name }}`, wantCount: 3},
		{name: "dependency ID action", contents: Generate, token: `{{ .Id }}`, wantCount: 1},
		{name: "range end actions", contents: Generate, token: `{{end}}`, wantCount: 3},
		{name: "upgrade template actions", contents: Upgrade, token: `{{`, wantCount: 0},
		{
			name: "generate Bootstrap stylesheet CDN", contents: Generate,
			token: `https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css`, wantCount: 1,
		},
		{
			name: "generate Bootstrap script CDN", contents: Generate,
			token: `https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/js/bootstrap.bundle.min.js`, wantCount: 1,
		},
		{
			name: "upgrade Bootstrap stylesheet CDN", contents: Upgrade,
			token: `https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css`, wantCount: 1,
		},
		{
			name: "upgrade Bootstrap script CDN", contents: Upgrade,
			token: `https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/js/bootstrap.bundle.min.js`, wantCount: 1,
		},
	}
	if len(tokens) == 0 {
		t.Fatal("template token characterization population is empty")
	}

	for _, token := range tokens {
		t.Run(token.name, func(t *testing.T) {
			if got := strings.Count(token.contents, token.token); got != token.wantCount {
				t.Fatalf("template token %q occurs %d times, want exactly %d", token.token, got, token.wantCount)
			}
		})
	}
}

func TestTemplatesParseWithExactCallerSelectedEngines(t *testing.T) {
	parsers := []struct {
		name       string
		contents   string
		parse      func(string) (string, error)
		wantParsed string
	}{
		{
			name: "Generate through text/template", contents: Generate, wantParsed: "generateTemplate",
			parse: func(contents string) (string, error) {
				parsed, err := texttemplate.New("generateTemplate").Parse(contents)
				if err != nil {
					return "", err
				}
				return parsed.Name(), nil
			},
		},
		{
			name: "Upgrade through html/template", contents: Upgrade, wantParsed: "upgradeTemplate",
			parse: func(contents string) (string, error) {
				parsed, err := htmltemplate.New("upgradeTemplate").Parse(contents)
				if err != nil {
					return "", err
				}
				return parsed.Name(), nil
			},
		},
	}
	if len(parsers) == 0 {
		t.Fatal("caller-selected template parser characterization population is empty")
	}

	for _, parser := range parsers {
		t.Run(parser.name, func(t *testing.T) {
			parsedName, err := parser.parse(parser.contents)
			if err != nil {
				t.Fatalf("caller-selected parser rejected existing template bytes: %v", err)
			}
			if parsedName != parser.wantParsed {
				t.Fatalf("parsed template name is %q, want exact caller name %q", parsedName, parser.wantParsed)
			}
		})
	}
}

func assertExactTemplateBytes(t *testing.T, contents string, wantLength int, wantSHA256 string) {
	t.Helper()
	digest := sha256.Sum256([]byte(contents))
	gotSHA256 := hex.EncodeToString(digest[:])
	if len(contents) != wantLength || gotSHA256 != wantSHA256 {
		t.Fatalf("complete template bytes have length %d and SHA-256 %s, want length %d and SHA-256 %s",
			len(contents), gotSHA256, wantLength, wantSHA256)
	}
}
