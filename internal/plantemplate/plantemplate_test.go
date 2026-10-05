package plantemplate

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"testing/fstest"
	"unicode/utf8"
)

var wantPaths = []string{
	".gitignore",
	"AGENTS.md",
	"README.md",
	"archive/README.md",
	"design/README.md",
	"goal.md",
	"koe/README.md",
	"log.md",
	"next-task.md",
	"process-checks.md",
	"session-end.md",
	"specs/README.md",
	"status.md",
}

var languages = []string{"nb", "en"}

const quotedGoal = "Ship \"the thing\" with `backticks`, {{.Goal}}, $(touch /tmp/pwned) and 'quotes'\nsecond line"

func testInput(language, goal string) Input {
	return Input{
		Language:     language,
		RootPath:     "/srv/work space/root",
		PlanningPath: "/srv/work space/root/planning repo",
		Goal:         goal,
		Repositories: []Repository{
			{Path: "/srv/work space/product one", GitCommonDir: "/srv/work space/product one/.git"},
			{Path: "/srv/other/product `two`", GitCommonDir: "/srv/other/common dir.git"},
		},
	}
}

func mustRender(t *testing.T, input Input) Bundle {
	t.Helper()
	bundle, err := Render(input)
	if err != nil {
		t.Fatalf("Render(%s) returned %v", input.Language, err)
	}
	return bundle
}

func (b Bundle) content(t *testing.T, path string) string {
	t.Helper()
	for _, file := range b.Files {
		if file.Path == path {
			return string(file.Content)
		}
	}
	t.Fatalf("bundle has no %s", path)
	return ""
}

func (b Bundle) paths() []string {
	var paths []string
	for _, file := range b.Files {
		paths = append(paths, file.Path)
	}
	return paths
}

func (b Bundle) all() string {
	var all strings.Builder
	for _, file := range b.Files {
		all.Write(file.Content)
	}
	return all.String()
}

func TestRenderBothLanguagesProduceSameRequiredRolesInStableOrder(t *testing.T) {
	var perLanguage [][]string
	for _, language := range languages {
		bundle := mustRender(t, testInput(language, "A goal"))

		if bundle.TemplateID != "agentic/planning" || bundle.TemplateVersion != "1" {
			t.Fatalf("%s template identity = %q/%q", language, bundle.TemplateID, bundle.TemplateVersion)
		}
		if !reflect.DeepEqual(bundle.paths(), wantPaths) {
			t.Fatalf("%s paths = %v, want %v", language, bundle.paths(), wantPaths)
		}
		for _, file := range bundle.Files {
			if len(file.Content) == 0 || file.Content[len(file.Content)-1] != '\n' {
				t.Errorf("%s %s does not end with a final newline", language, file.Path)
			}
			if !utf8.Valid(file.Content) || strings.ContainsRune(string(file.Content), 0) {
				t.Errorf("%s %s is not valid text", language, file.Path)
			}
		}
		perLanguage = append(perLanguage, bundle.paths())
	}
	if !reflect.DeepEqual(perLanguage[0], perLanguage[1]) {
		t.Fatalf("languages disagree on paths: %v vs %v", perLanguage[0], perLanguage[1])
	}
}

func TestRenderProseFollowsSelectedLanguage(t *testing.T) {
	nb := mustRender(t, testInput("nb", "A goal"))
	en := mustRender(t, testInput("en", "A goal"))

	if !strings.Contains(nb.content(t, "AGENTS.md"), "Stående regler") || strings.Contains(nb.content(t, "AGENTS.md"), "Standing rules") {
		t.Error("nb AGENTS.md is not Norwegian Bokmål prose")
	}
	if !strings.Contains(en.content(t, "AGENTS.md"), "Standing rules") || strings.Contains(en.content(t, "AGENTS.md"), "Stående regler") {
		t.Error("en AGENTS.md is not English prose")
	}
	if nb.content(t, "README.md") == en.content(t, "README.md") {
		t.Error("README.md is identical in both languages")
	}
}

func TestRenderBindsContextAsDataAndNeverIntoPaths(t *testing.T) {
	baseline := mustRender(t, testInput("en", "plain")).paths()

	for _, language := range languages {
		input := testInput(language, quotedGoal)
		bundle := mustRender(t, input)

		goal := bundle.content(t, "goal.md")
		for _, line := range strings.Split(quotedGoal, "\n") {
			if !strings.Contains(goal, "> "+line) {
				t.Errorf("%s goal.md lost goal line %q:\n%s", language, line, goal)
			}
		}
		for _, name := range []string{"README.md", "AGENTS.md"} {
			content := bundle.content(t, name)
			for _, bound := range []string{
				input.RootPath, input.PlanningPath,
				"/srv/work space/product one", "/srv/work space/product one/.git",
				"/srv/other/product `two`", "/srv/other/common dir.git",
			} {
				if !strings.Contains(content, bound) {
					t.Errorf("%s %s does not bind %q", language, name, bound)
				}
			}
		}
		if !reflect.DeepEqual(bundle.paths(), baseline) {
			t.Errorf("%s paths changed with goal/repository text: %v", language, bundle.paths())
		}
		if strings.Contains(bundle.all(), "<no value>") {
			t.Errorf("%s bundle has an unfilled template value", language)
		}
	}
}

func TestRenderWithoutGoalSaysChoosingTheGoalIsNext(t *testing.T) {
	wantNext := map[string]string{"nb": "velge eller avklare målet", "en": "choose or clarify the goal"}
	for _, language := range languages {
		bundle := mustRender(t, testInput(language, "  \n"))

		for _, name := range []string{"goal.md", "next-task.md"} {
			if !strings.Contains(strings.ToLower(bundle.content(t, name)), wantNext[language]) {
				t.Errorf("%s %s does not say that choosing the goal is next:\n%s", language, name, bundle.content(t, name))
			}
		}
		if strings.Contains(bundle.content(t, "goal.md"), "> ") {
			t.Errorf("%s goal.md quotes a goal that was not supplied", language)
		}
	}
}

func TestRenderStatusDoesNotInventProgress(t *testing.T) {
	wantStatus := map[string][]string{
		"nb": {"Ingen analyse er gjort", "Ingen Epic er valgt", "Ingen tester er kjørt", "Ingen agent kjører", "Ingen menneskelig vurdering"},
		"en": {"No analysis has been done", "No Epic has been selected", "No tests have been run", "No agent is running", "No human judgment"},
	}
	for _, language := range languages {
		for _, goal := range []string{"", "Build a thing"} {
			status := mustRender(t, testInput(language, goal)).content(t, "status.md")
			for _, phrase := range wantStatus[language] {
				if !strings.Contains(status, phrase) {
					t.Errorf("%s status.md (goal %q) lacks honest line %q", language, goal, phrase)
				}
			}
		}
	}
}

func TestRenderLogOnlyDescribesActualSetup(t *testing.T) {
	for _, language := range languages {
		log := mustRender(t, testInput(language, "")).content(t, "log.md")

		if !strings.Contains(log, "agentic/planning") || !strings.Contains(log, "`"+language+"`") {
			t.Errorf("%s log.md does not record the template and language used:\n%s", language, log)
		}
		if got := strings.Count(log, "\n- "); got != 1 {
			t.Errorf("%s log.md has %d entries, want exactly the setup entry", language, got)
		}
	}
}

func TestRenderZeroRepositoriesStaysHonest(t *testing.T) {
	wantUnbound := map[string]string{"nb": "Ingen produktrepo er bundet ennå", "en": "No product repository is bound yet"}
	for _, language := range languages {
		input := testInput(language, "")
		input.Repositories = nil

		agents := mustRender(t, input).content(t, "AGENTS.md")

		if !strings.Contains(agents, wantUnbound[language]) {
			t.Errorf("%s AGENTS.md does not say that no repository is bound", language)
		}
	}
}

func TestRenderRejectsInvalidInputWithoutOutput(t *testing.T) {
	cases := map[string]func(*Input){
		"empty language":         func(in *Input) { in.Language = "" },
		"unknown language":       func(in *Input) { in.Language = "de" },
		"wrong case language":    func(in *Input) { in.Language = "NB" },
		"padded language":        func(in *Input) { in.Language = " en" },
		"missing root":           func(in *Input) { in.RootPath = "" },
		"missing planning path":  func(in *Input) { in.PlanningPath = "" },
		"relative root":          func(in *Input) { in.RootPath = "root" },
		"relative planning path": func(in *Input) { in.PlanningPath = "../planning" },
		"newline in root":        func(in *Input) { in.RootPath = "/srv/root\n## injected" },
		"control in planning":    func(in *Input) { in.PlanningPath = "/srv/plan\x00ning" },
		"missing repo path":      func(in *Input) { in.Repositories[0].Path = "" },
		"missing git common dir": func(in *Input) { in.Repositories[1].GitCommonDir = "" },
		"relative repo path":     func(in *Input) { in.Repositories[0].Path = "product" },
		"newline in repo path":   func(in *Input) { in.Repositories[0].Path = "/srv/a\n/srv/b" },
		"nul in goal":            func(in *Input) { in.Goal = "goal\x00" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			input := testInput("en", "goal")
			mutate(&input)

			bundle, err := Render(input)

			if err == nil {
				t.Fatalf("Render accepted invalid input %+v", input)
			}
			if bundle.Files != nil || bundle.TemplateID != "" || bundle.TemplateVersion != "" {
				t.Fatalf("Render returned output %+v together with error %v", bundle, err)
			}
		})
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	for _, language := range languages {
		first := mustRender(t, testInput(language, quotedGoal))
		second := mustRender(t, testInput(language, quotedGoal))

		if !reflect.DeepEqual(first, second) {
			t.Errorf("%s: identical input produced different bundles", language)
		}
	}
}

func TestRenderDoesNotMutateItsInput(t *testing.T) {
	input := testInput("en", quotedGoal)
	before := testInput("en", quotedGoal)

	mustRender(t, input)

	if !reflect.DeepEqual(input, before) {
		t.Fatalf("Render mutated its input: %+v", input)
	}
}

func TestRenderInstructionsPreserveAutonomousPlanningAndSpecFocus(t *testing.T) {
	want := map[string][]string{
		"nb": {
			"Planagenten", "høy autonomi", "Spec", "sub-agenter", "Herdr", "samlet ansvar",
			"ukjent oppstart", "faktisk sandbox", "Produktkode", "dekkende menneskelig fullmakt",
			"mal er aldri et nytt produktmandat", "liten direkte løsning", "ikke en menneskelig start",
		},
		"en": {
			"planning agent", "high autonomy", "Spec", "subagents", "Herdr", "stay responsible",
			"unknown start", "actual sandbox", "Product code", "human authority that actually covers",
			"never a new product mandate", "small direct solution", "no extra human start",
		},
	}
	for _, language := range languages {
		agents := mustRender(t, testInput(language, "goal")).content(t, "AGENTS.md")
		for _, phrase := range want[language] {
			if !strings.Contains(agents, phrase) {
				t.Errorf("%s AGENTS.md lacks %q", language, phrase)
			}
		}
	}
}

func TestRenderContainsNoSourceProjectHistoryPrivatePathsOrRuntimeSettings(t *testing.T) {
	forbidden := []string{
		"/Users/", "/home/", "/private/", "codex", "claude", "gpt", "sonnet", "opus", "gemini",
		"codex-dev-start", "observer-log", "handoff", "permission-mode", "permissions.json", "settings.json",
		"auto_review", "on-request", "bypass", "api_key", "secret", "token", "password",
		"ply workflow", "ply build", "pflyt", "økt 1", "session 1",
	}
	for _, language := range languages {
		bundle := mustRender(t, testInput(language, "A plain goal"))
		all := strings.ToLower(bundle.all())
		for _, word := range forbidden {
			if strings.Contains(all, strings.ToLower(word)) {
				t.Errorf("%s bundle contains forbidden %q", language, word)
			}
		}
		for _, script := range []string{".sh", "./scripts", "scripts/"} {
			if strings.Contains(all, script) {
				t.Errorf("%s bundle references a startup/observer script (%q) that is not in the bundle", language, script)
			}
		}
	}
}

func TestRenderEmptyRoleDirectoriesHaveBriefGuidance(t *testing.T) {
	for _, language := range languages {
		bundle := mustRender(t, testInput(language, ""))
		for _, path := range []string{"design/README.md", "specs/README.md", "koe/README.md", "archive/README.md"} {
			content := bundle.content(t, path)
			if len(content) < 80 || len(content) > 1500 {
				t.Errorf("%s %s has %d bytes of guidance, want a brief but real note", language, path, len(content))
			}
		}
	}
}

func TestRenderStaysWithinSmallBundleLimits(t *testing.T) {
	for _, language := range languages {
		bundle := mustRender(t, testInput(language, quotedGoal))

		if size := len(bundle.content(t, "AGENTS.md")); size >= 8000 {
			t.Errorf("%s AGENTS.md is %d bytes, want < 8000", language, size)
		}
		if size := len(bundle.all()); size >= 40000 {
			t.Errorf("%s bundle is %d bytes, want < 40000", language, size)
		}
	}
}

func TestRenderFilesReturnsParseErrorWithoutPanicOrPartialOutput(t *testing.T) {
	src := fstest.MapFS{
		"en/agents.tmpl": {Data: []byte("fine {{.RootPath}}\n")},
		"en/readme.tmpl": {Data: []byte("broken {{")},
	}
	roles := []role{{Path: "AGENTS.md", Source: "agents"}, {Path: "README.md", Source: "readme"}}
	var files []File
	var err error
	var recovered interface{}

	func() {
		defer func() { recovered = recover() }()
		files, err = renderFiles(src, roles, "en", view{RootPath: "/srv"})
	}()

	if recovered != nil {
		t.Fatalf("renderFiles panicked: %v", recovered)
	}
	if err == nil || !strings.Contains(err.Error(), "en/readme.tmpl") {
		t.Fatalf("renderFiles error = %v, want a parse error naming en/readme.tmpl", err)
	}
	if files != nil {
		t.Fatalf("renderFiles returned partial output %+v", files)
	}
}

func TestRenderFilesReturnsExecuteAndMissingSourceErrors(t *testing.T) {
	roles := []role{{Path: "AGENTS.md", Source: "agents"}}

	_, executeErr := renderFiles(fstest.MapFS{"en/agents.tmpl": {Data: []byte("{{.NoSuchField}}\n")}}, roles, "en", view{})
	if executeErr == nil {
		t.Error("renderFiles accepted a template that references an unknown field")
	}
	_, missingErr := renderFiles(fstest.MapFS{}, roles, "en", view{})
	if missingErr == nil {
		t.Error("renderFiles accepted a missing template source")
	}
}

func TestRenderFilesReturnsFilesInStablePathOrder(t *testing.T) {
	src := fstest.MapFS{"en/x.tmpl": {Data: []byte("x\n")}}
	unsorted := []role{{Path: "z.md", Source: "x"}, {Path: "a/b.md", Source: "x"}, {Path: ".hidden", Source: "x"}}

	files, err := renderFiles(src, unsorted, "en", view{})

	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range files {
		got = append(got, f.Path)
	}
	if want := []string{".hidden", "a/b.md", "z.md"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("renderFiles order = %v, want %v", got, want)
	}
}

func TestRenderFilesRejectsUnsafeOutputPaths(t *testing.T) {
	src := fstest.MapFS{"en/x.tmpl": {Data: []byte("x\n")}}
	for _, path := range []string{"", "/abs.md", "../escape.md", "dir/../../escape.md", "a//b.md", "dir/", "back\\slash.md"} {
		if _, err := renderFiles(src, []role{{Path: path, Source: "x"}}, "en", view{}); err == nil {
			t.Errorf("renderFiles accepted unsafe output path %q", path)
		}
	}
	dup := []role{{Path: "a.md", Source: "x"}, {Path: "a.md", Source: "x"}}
	if _, err := renderFiles(src, dup, "en", view{}); err == nil {
		t.Error("renderFiles accepted duplicate output paths")
	}
}

func TestCodeSpanSurvivesBackticksAndEdgeSpaces(t *testing.T) {
	cases := map[string]string{
		"/srv/plain":    "`/srv/plain`",
		"/srv/a`b":      "``/srv/a`b``",
		"/srv/a``b`c":   "```/srv/a``b`c```",
		"`/srv/start":   "`` `/srv/start ``",
		"/srv/end`":     "`` /srv/end` ``",
		"/srv/sp ace":   "`/srv/sp ace`",
		"/srv/tick`end": "``/srv/tick`end``",
	}
	for in, want := range cases {
		if got := codeSpan(in); got != want {
			t.Errorf("codeSpan(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestQuoteBlockPrefixesEveryLineAndNormalizesLineEndings(t *testing.T) {
	if got, want := quoteBlock("one\r\n\r\ntwo\rthree"), "> one\n>\n> two\n> three"; got != want {
		t.Errorf("quoteBlock = %q, want %q", got, want)
	}
}

func bundleBytes(bundle Bundle) int {
	total := 0
	for _, file := range bundle.Files {
		total += len(file.Content)
	}
	return total
}

func manyRepositories(count int) []Repository {
	var repositories []Repository
	for i := 0; i < count; i++ {
		repositories = append(repositories, Repository{
			Path:         fmt.Sprintf("/work/root/product-service-%03d/main", i),
			GitCommonDir: fmt.Sprintf("/work/root/product-service-%03d/main/.git", i),
		})
	}
	return repositories
}

func assertRejectedWithEmptyBundle(t *testing.T, input Input, what string) {
	t.Helper()
	bundle, err := Render(input)
	if err == nil {
		t.Fatalf("Render accepted %s (bundle %d bytes, AGENTS.md %d bytes)", what, bundleBytes(bundle), len(bundle.content(t, "AGENTS.md")))
	}
	if bundle.Files != nil || bundle.TemplateID != "" || bundle.TemplateVersion != "" {
		t.Fatalf("Render returned %+v together with error %v", bundle, err)
	}
}

func TestRenderRejectsGoalThatPushesBundleOverLimit(t *testing.T) {
	for _, language := range languages {
		input := testInput(language, strings.Repeat("x", 40000))

		assertRejectedWithEmptyBundle(t, input, "a 40000-byte goal")
	}
}

func TestRenderBundleLimitIsMeasuredInUTF8BytesNotCharacters(t *testing.T) {
	for _, language := range languages {
		// 25000 characters of 2 bytes each: under 40000 characters, over 40000 bytes.
		input := testInput(language, strings.Repeat("å", 25000))

		assertRejectedWithEmptyBundle(t, input, "a goal of 25000 two-byte characters")
	}
}

func TestRenderBundleLimitIsExactlyBelow40000Bytes(t *testing.T) {
	for _, language := range languages {
		for _, fill := range []string{"x", "å"} {
			accepted := sort.Search(40000, func(n int) bool {
				_, err := Render(testInput(language, strings.Repeat(fill, n+1)))
				return err != nil
			})
			// accepted is the largest goal length that renders.
			bundle := mustRender(t, testInput(language, strings.Repeat(fill, accepted)))
			if size := bundleBytes(bundle); size >= 40000 || size < 39900 {
				t.Errorf("%s fill %q: largest accepted bundle is %d bytes, want just below 40000", language, fill, size)
			}
			if _, err := Render(testInput(language, strings.Repeat(fill, accepted+1))); err == nil {
				t.Errorf("%s fill %q: Render accepted one more character than the limit allows", language, fill)
			}
		}
	}
}

func TestRenderRejectsManyRepositoriesThatPushAgentsOverLimit(t *testing.T) {
	for _, language := range languages {
		input := testInput(language, "")
		input.Repositories = manyRepositories(64)

		assertRejectedWithEmptyBundle(t, input, "64 repositories")
	}
}

func TestRenderAgentsLimitIsMeasuredInUTF8BytesNotCharacters(t *testing.T) {
	for _, language := range languages {
		input := testInput(language, "")
		// About 3000 extra characters (6000 bytes) in AGENTS.md: below 8000 characters, above 8000 bytes.
		input.Repositories = []Repository{{
			Path:         "/work/" + strings.Repeat("å", 2000),
			GitCommonDir: "/work/" + strings.Repeat("å", 1000),
		}}

		assertRejectedWithEmptyBundle(t, input, "AGENTS.md of fewer than 8000 characters but more than 8000 bytes")
	}
}

func TestRenderAgentsLimitIsExactlyBelow8000Bytes(t *testing.T) {
	for _, language := range languages {
		accepted := sort.Search(500, func(n int) bool {
			input := testInput(language, "")
			input.Repositories = manyRepositories(n + 1)
			_, err := Render(input)
			return err != nil
		})
		input := testInput(language, "")
		input.Repositories = manyRepositories(accepted)
		bundle := mustRender(t, input)
		if size := len(bundle.content(t, "AGENTS.md")); size >= 8000 || size < 7800 {
			t.Errorf("%s: largest accepted AGENTS.md is %d bytes with %d repositories, want just below 8000", language, size, accepted)
		}
		input.Repositories = manyRepositories(accepted + 1)
		if _, err := Render(input); err == nil {
			t.Errorf("%s: Render accepted %d repositories beyond the AGENTS.md limit", language, accepted+1)
		}
	}
}
