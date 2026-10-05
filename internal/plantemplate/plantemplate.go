// Package plantemplate renders the bundled agentic/planning template. It is pure:
// it reads only embedded sources and performs no filesystem, Git, network or
// agent effects. Creating the repository from the bundle is the caller's job.
package plantemplate

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/devdimensionlab/plybuild/pkg/file"
)

const (
	templateID      = "agentic/planning"
	templateVersion = "1"

	languageNorwegian = "nb"
	languageEnglish   = "en"
)

//go:embed templates
var embedded embed.FS

// Repository is one explicitly bound product repository.
type Repository struct{ Path, GitCommonDir string }

// Input is the context bound into the rendered documents.
type Input struct {
	Language, RootPath, PlanningPath, Goal string
	Repositories                           []Repository
}

// File is one rendered text file at a fixed, relative, slash-separated path.
type File struct {
	Path    string
	Content []byte
}

// Bundle is the rendered template with Files in byte order of Path.
type Bundle struct {
	TemplateID, TemplateVersion string
	Files                       []File
}

type role struct{ Path, Source string }

var roles = []role{
	{".gitignore", "gitignore"},
	{"AGENTS.md", "agents"},
	{"README.md", "readme"},
	{"archive/README.md", "archive"},
	{"design/README.md", "design"},
	{"goal.md", "goal"},
	{"koe/README.md", "koe"},
	{"log.md", "log"},
	{"next-task.md", "next-task"},
	{"process-checks.md", "process-checks"},
	{"session-end.md", "session-end"},
	{"specs/README.md", "specs"},
	{"status.md", "status"},
}

type repositoryView struct{ Path, PathCode, GitCommonDir, CommonDirCode string }

type view struct {
	TemplateID, TemplateVersion, Language string
	RootPath, RootCode                    string
	PlanningPath, PlanningCode            string
	HasGoal                               bool
	GoalQuote                             string
	HasRepositories                       bool
	Repositories                          []repositoryView
}

// Render validates the input and renders the bundle. On any error it returns an
// empty Bundle.
func Render(input Input) (Bundle, error) {
	if input.Language != languageNorwegian && input.Language != languageEnglish {
		return Bundle{}, fmt.Errorf("unsupported language %q: want %s or %s", input.Language, languageNorwegian, languageEnglish)
	}
	if err := validatePath("root path", input.RootPath); err != nil {
		return Bundle{}, err
	}
	if err := validatePath("planning path", input.PlanningPath); err != nil {
		return Bundle{}, err
	}
	for i, repository := range input.Repositories {
		if err := validatePath(fmt.Sprintf("repository %d path", i+1), repository.Path); err != nil {
			return Bundle{}, err
		}
		if err := validatePath(fmt.Sprintf("repository %d Git common dir", i+1), repository.GitCommonDir); err != nil {
			return Bundle{}, err
		}
	}
	goal := strings.TrimSpace(input.Goal)
	if err := validateGoal(goal); err != nil {
		return Bundle{}, err
	}

	v := view{
		TemplateID:      templateID,
		TemplateVersion: templateVersion,
		Language:        input.Language,
		RootPath:        input.RootPath,
		RootCode:        codeSpan(input.RootPath),
		PlanningPath:    input.PlanningPath,
		PlanningCode:    codeSpan(input.PlanningPath),
		HasGoal:         goal != "",
		HasRepositories: len(input.Repositories) > 0,
	}
	if v.HasGoal {
		v.GoalQuote = quoteBlock(goal)
	}
	for _, repository := range input.Repositories {
		v.Repositories = append(v.Repositories, repositoryView{
			Path:          repository.Path,
			PathCode:      codeSpan(repository.Path),
			GitCommonDir:  repository.GitCommonDir,
			CommonDirCode: codeSpan(repository.GitCommonDir),
		})
	}

	sources, err := fs.Sub(embedded, "templates")
	if err != nil {
		return Bundle{}, err
	}
	files, err := renderFiles(sources, roles, input.Language, v)
	if err != nil {
		return Bundle{}, err
	}
	return Bundle{TemplateID: templateID, TemplateVersion: templateVersion, Files: files}, nil
}

func renderFiles(sources fs.FS, roles []role, language string, v view) ([]File, error) {
	seen := map[string]bool{}
	var files []File
	for _, r := range roles {
		if err := validateOutputPath(r.Path); err != nil {
			return nil, err
		}
		if seen[r.Path] {
			return nil, fmt.Errorf("duplicate output path %q", r.Path)
		}
		seen[r.Path] = true

		source := path.Join(language, r.Source+".tmpl")
		raw, err := fs.ReadFile(sources, source)
		if err != nil {
			return nil, fmt.Errorf("read template %s: %w", source, err)
		}
		content, err := file.RenderBytes(source, raw, v)
		if err != nil {
			return nil, err
		}
		files = append(files, File{Path: r.Path, Content: content})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func validateOutputPath(p string) error {
	if p == "" || strings.ContainsAny(p, "\\\x00") || strings.HasPrefix(p, "/") || path.Clean(p) != p || p == ".." || strings.HasPrefix(p, "../") {
		return fmt.Errorf("unsafe output path %q", p)
	}
	return nil
}

func validatePath(name, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	if !filepath.IsAbs(value) {
		return fmt.Errorf("%s must be absolute", name)
	}
	if strings.IndexFunc(value, unicode.IsControl) >= 0 {
		return fmt.Errorf("%s must not contain control characters", name)
	}
	return nil
}

func validateGoal(goal string) error {
	for _, r := range goal {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return fmt.Errorf("goal must not contain control characters")
		}
	}
	return nil
}

// codeSpan wraps s in a Markdown code span whose delimiter is longer than any
// backtick run in s, so s survives as data.
func codeSpan(s string) string {
	longest, run := 0, 0
	for _, r := range s {
		if r == '`' {
			run++
			if run > longest {
				longest = run
			}
		} else {
			run = 0
		}
	}
	delimiter := strings.Repeat("`", longest+1)
	if strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`") {
		return delimiter + " " + s + " " + delimiter
	}
	return delimiter + s + delimiter
}

// quoteBlock renders text as a Markdown block quote, line by line.
func quoteBlock(text string) string {
	text = strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line == "" {
			lines[i] = ">"
		} else {
			lines[i] = "> " + line
		}
	}
	return strings.Join(lines, "\n")
}
