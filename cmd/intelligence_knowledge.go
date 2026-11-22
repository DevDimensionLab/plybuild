package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/knowledge"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
	"github.com/spf13/cobra"
)

const intelligenceKnowledgeKind = "PlyIntelligenceKnowledge@1"

// Statuses of a knowledge source. The set is open: readers treat an unknown
// status as a source without items.
const (
	knowledgeSourceOK            = "ok"
	knowledgeSourceNotConfigured = "not_configured"
	knowledgeSourceInvalid       = "invalid"
	knowledgeSourceMissing       = "missing"
	knowledgeSourceUnreadable    = "unreadable"
)

// knowledgeSource is the Project's knowledge pointer and what came of it.
// Diagnostic is text for people, empty when Status is ok.
type knowledgeSource struct {
	Repo       string `json:"repo"`
	Path       string `json:"path"`
	Status     string `json:"status"`
	Diagnostic string `json:"diagnostic"`
}

// resolveKnowledgeDirectory is the physically resolved archive directory of
// project's knowledge pointer, or "" with a source that says why there is none. The pointer is
// checked only here, so a bad one never stops reading the workspace.
func resolveKnowledgeDirectory(project workspace.ProjectRecord) (string, knowledgeSource) {
	pointer := project.Knowledge
	if pointer == nil {
		return "", knowledgeSource{Status: knowledgeSourceNotConfigured, Diagnostic: fmt.Sprintf("Project %s has no knowledge source in %s", project.ID, workspacedef.File)}
	}
	source := knowledgeSource{Repo: string(pointer.Repo), Path: pointer.Path, Status: knowledgeSourceInvalid}
	if pointer.Locator == "" {
		source.Diagnostic = fmt.Sprintf("knowledge repository %s is not a repository of Project %s", pointer.Repo, project.ID)
		return "", source
	}
	if !filepath.IsLocal(filepath.FromSlash(pointer.Path)) {
		source.Diagnostic = fmt.Sprintf("knowledge path %s must be a relative path inside repository %s", pointer.Path, pointer.Repo)
		return "", source
	}
	dir := filepath.Join(pointer.Locator, filepath.FromSlash(pointer.Path))
	info, err := os.Stat(dir)
	switch {
	case os.IsNotExist(err):
		source.Status, source.Diagnostic = knowledgeSourceMissing, fmt.Sprintf("knowledge directory %s does not exist in repository %s", pointer.Path, pointer.Repo)
		return "", source
	case err == nil && !info.IsDir():
		source.Status, source.Diagnostic = knowledgeSourceMissing, fmt.Sprintf("knowledge directory %s is not a directory in repository %s", pointer.Path, pointer.Repo)
		return "", source
	}
	resolved, inside, err := resolveInside(pointer.Locator, dir)
	switch {
	case err != nil:
		source.Status, source.Diagnostic = knowledgeSourceUnreadable, fmt.Sprintf("cannot resolve knowledge path %s in repository %s: %v", pointer.Path, pointer.Repo, err)
		return "", source
	case !inside:
		source.Diagnostic = fmt.Sprintf("knowledge path %s leaves repository %s through a symbolic link", pointer.Path, pointer.Repo)
		return "", source
	}
	source.Status = knowledgeSourceOK
	return resolved, source
}

// resolveInside resolves every symbolic link in dir and reports whether the
// result lies inside the resolved repository locator. The search walks the
// resolved directory, since a walk does not follow a root that is a link.
func resolveInside(locator, dir string) (string, bool, error) {
	repository, err := filepath.EvalSymlinks(locator)
	if err != nil {
		return "", false, err
	}
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", false, err
	}
	rel, err := filepath.Rel(repository, resolved)
	return resolved, err == nil && filepath.IsLocal(rel), nil
}

// searchKnowledge searches the archive of project's knowledge pointer. A
// source that cannot be used gives an empty result, never an error.
func searchKnowledge(project workspace.ProjectRecord, query knowledge.Query) (knowledge.Result, knowledgeSource) {
	empty := knowledge.Result{Items: []knowledge.Item{}, Invalid: []knowledge.InvalidFile{}}
	dir, source := resolveKnowledgeDirectory(project)
	if source.Status != knowledgeSourceOK {
		return empty, source
	}
	result, err := knowledge.Search(dir, query)
	if err != nil {
		source.Status, source.Diagnostic = knowledgeSourceUnreadable, fmt.Sprintf("cannot read knowledge directory %s in repository %s: %v", source.Path, source.Repo, err)
		return empty, source
	}
	return result, source
}

func newIntelligenceKnowledgeCommand(d workspace.Dependencies) *cobra.Command {
	var project, format, repo, theme, tag, status string
	var limit int
	command := &cobra.Command{
		Use:   "knowledge [WORD...]",
		Short: "Find posts in the Project's curated knowledge archive",
		Long: "Reads the k-*.md posts under the knowledge directory that the Project's \"knowledge\" pointer in " + workspacedef.File + " names, and writes nothing.\n" +
			"A word hits a post when it is part of a word in its id, tags, title, situation line or body; case does not matter. Posts with more distinct words come first, then posts whose scope fits --repo, then hits in id, tags or title before the situation line before the body, then id.\n" +
			"Without words every post is listed by theme and id, and --limit does not apply. A post without status or project is not filtered out on them.\n" +
			"Without a usable pointer the answer is empty with a diagnostic and exit status 0. Files that cannot be read as posts are listed under invalid.",
		Example: "  ply intelligence knowledge headless skjermbilde\n" +
			"  ply intelligence knowledge refactor mutant --repo ply --format json\n" +
			"  ply intelligence knowledge --theme tooling",
		Args: func(cmd *cobra.Command, args []string) error {
			switch {
			case format != "text" && format != "json":
				return intelligenceInvalidArguments("--format must be text or json")
			case limit < 1:
				return intelligenceInvalidArguments("--limit must be at least 1")
			}
			_, err := knowledgeStatuses(status)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			resolved, err := resolveIntelligenceProject(d, project)
			if err != nil {
				return err
			}
			statuses, _ := knowledgeStatuses(status)
			query := knowledge.Query{Words: args, Project: string(resolved.ID), Repo: repo, Theme: theme, Tag: tag, Statuses: statuses, Limit: limit}
			result, source := searchKnowledge(resolved, query)
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string              `json:"kind"`
					SchemaVersion int                 `json:"schema_version"`
					Project       workspace.ProjectID `json:"project"`
					Source        knowledgeSource     `json:"source"`
					knowledge.Result
				}{intelligenceKnowledgeKind, intelligenceSchemaVersion, resolved.ID, source, result})
			}
			return writeKnowledgeText(cmd.OutOrStdout(), source, result, len(args) > 0)
		},
	}
	command.Flags().StringVar(&project, "project", "", "Project ID (default: the Project whose wrapper contains the current directory)")
	command.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	command.Flags().StringVar(&repo, "repo", "", "rank posts whose scope fits this repository first; filters nothing")
	command.Flags().StringVar(&theme, "theme", "", "only posts of this theme, the first directory under the archive")
	command.Flags().StringVar(&tag, "tag", "", "only posts with exactly this tag")
	command.Flags().StringVar(&status, "status", strings.Join(knowledge.DefaultStatuses, ","), "comma-separated statuses, or all")
	command.Flags().IntVar(&limit, "limit", 10, "at most this many matches")
	return command
}

// knowledgeStatuses parses --status: nil for all.
func knowledgeStatuses(value string) ([]string, error) {
	if value == "all" {
		return nil, nil
	}
	statuses := strings.Split(value, ",")
	for _, s := range statuses {
		if s == "" {
			return nil, intelligenceInvalidArguments("--status must be a comma-separated list of statuses, or all")
		}
	}
	return statuses, nil
}

func writeKnowledgeText(out io.Writer, source knowledgeSource, result knowledge.Result, words bool) error {
	var b strings.Builder
	if source.Status != knowledgeSourceOK {
		fmt.Fprintf(&b, "knowledge source %s: %s\n", source.Status, source.Diagnostic)
	}
	if source.Status == knowledgeSourceOK && words && result.Matched == 0 {
		b.WriteString("no matches\n")
	}
	for _, item := range result.Items {
		writeKnowledgeItem(&b, "", item)
	}
	for _, invalid := range result.Invalid {
		fmt.Fprintf(&b, "invalid %s: %s\n", invalid.File, invalid.Error)
	}
	if shown := len(result.Items); shown < result.Matched {
		fmt.Fprintf(&b, "%d of %d matches shown; raise --limit for more\n", shown, result.Matched)
	}
	if result.Hidden > 0 {
		noun := "posts"
		if result.Hidden == 1 {
			noun = "post"
		}
		fmt.Fprintf(&b, "%d more %s hidden by --status, --theme, --tag or the Project scope\n", result.Hidden, noun)
	}
	_, err := io.WriteString(out, b.String())
	return err
}

// writeKnowledgeItem writes the two lines of one post, the second indented
// two more spaces than indent.
func writeKnowledgeItem(b *strings.Builder, indent string, item knowledge.Item) {
	repos := strings.Join(item.Scope.Repos, ",")
	if repos == "" {
		repos = "-"
	}
	var summary []string
	for _, part := range []string{item.Title, item.Situation} {
		if part != "" {
			summary = append(summary, part)
		}
	}
	fmt.Fprintf(b, "%s%s  %s/%s  %s\n%s  %s\n", indent, item.ID, item.Theme, item.Status, repos, indent, strings.Join(summary, " — "))
}

// sharedExperience is at most three posts of project's knowledge archive
// that text points at, ranked for repos (see knowledge.Relevant). A source
// that cannot be used gives none, so a command that shows them never fails
// or changes on account of the archive.
func sharedExperience(project workspace.ProjectRecord, text string, repos ...string) []knowledge.Item {
	dir, source := resolveKnowledgeDirectory(project)
	if source.Status != knowledgeSourceOK {
		return nil
	}
	items, err := knowledge.Relevant(dir, text, repos, string(project.ID))
	if err != nil {
		return nil
	}
	return items
}

// writeSharedExperience writes items under "Shared experience" with the
// path to open, and nothing without items.
func writeSharedExperience(b *strings.Builder, items []knowledge.Item) {
	if len(items) == 0 {
		return
	}
	b.WriteString("Shared experience (open the file at path):\n")
	for _, item := range items {
		writeKnowledgeItem(b, "  ", item)
		fmt.Fprintf(b, "    path: %s\n", item.Path)
	}
}
