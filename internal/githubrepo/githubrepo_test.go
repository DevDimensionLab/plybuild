package githubrepo

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

const (
	notFound   = `{"data":{"repository":null},"errors":[{"type":"NOT_FOUND","path":["repository"],"locations":[{"line":1,"column":40}],"message":"PAYLOAD-SECRET"}]}`
	queryField = "query=query($owner:String!,$name:String!){repository(owner:$owner,name:$name){nameWithOwner visibility isEmpty sshUrl}}"
)

func repository(fields string) string { return `{"data":{"repository":{` + fields + `}}}` + "\n" }

func observed(nameWithOwner, visibility, isEmpty, sshURL string) string {
	return repository(fmt.Sprintf(`"nameWithOwner":%q,"visibility":%q,"isEmpty":%s,"sshUrl":%q`, nameWithOwner, visibility, isEmpty, sshURL))
}

var privateEmpty = observed("Org/svc", "PRIVATE", "true", "git@github.com:Org/svc.git")

type reply struct {
	stdout string
	code   int
}

// call records the context state at run time, since EnsurePrivate cancels its readback context on return.
type call struct {
	args, env   []string
	ctxErr      error
	deadline    time.Time
	hasDeadline bool
}

type fakeGH struct {
	replies []reply
	calls   []call
	cancel  context.CancelFunc
}

func (f *fakeGH) Run(c process.Command) error {
	if c.Name != "gh" || c.Dir != "" || c.Start || c.Context == nil || c.Env == nil {
		return fmt.Errorf("unexpected command %s %q: want gh without Dir or Start, with Context and Env", c.Name, c.Args)
	}
	deadline, hasDeadline := c.Context.Deadline()
	f.calls = append(f.calls, call{c.Args, c.Env, c.Context.Err(), deadline, hasDeadline})
	r := f.replies[0]
	f.replies = f.replies[1:]
	fmt.Fprint(c.Stdout, r.stdout)
	fmt.Fprint(c.Stderr, "STDERR-SECRET gho_FAKETOKEN")
	if c.Args[0] == "repo" && f.cancel != nil {
		f.cancel()
		return context.Canceled
	}
	if r.code != 0 {
		return exec.Command("/bin/sh", "-c", fmt.Sprintf("exit %d", r.code)).Run()
	}
	return nil
}

func (f *fakeGH) creates() int {
	n := 0
	for _, c := range f.calls {
		if c.args[0] == "repo" {
			n++
		}
	}
	return n
}

func readArgs(owner, name string) []string {
	return []string{"api", "graphql", "--hostname", "github.com", "-f", queryField, "-f", "owner=" + owner, "-f", "name=" + name}
}

func assertSecretFree(t *testing.T, err error) {
	t.Helper()
	for _, secret := range []string{"SECRET", "gho_", "git@", "https://"} {
		if err != nil && strings.Contains(err.Error(), secret) {
			t.Fatalf("error %q contains %q", err, secret)
		}
	}
}

func TestEnsurePrivateAcceptsVisiblePrivateEmptyWithoutCreate(t *testing.T) {
	long := strings.Repeat("a", 19) + "-" + strings.Repeat("b", 19)
	for _, test := range []struct{ owner, name, stdout, url string }{
		{"Org", "svc", privateEmpty, "git@github.com:Org/svc.git"},
		{"org", "SVC", privateEmpty, "git@github.com:Org/svc.git"},
		{long, "_a.b-c", observed(long+"/_a.b-c", "PRIVATE", "true", "git@github.com:"+long+"/_a.b-c.git"), "git@github.com:" + long + "/_a.b-c.git"},
	} {
		gh := &fakeGH{replies: []reply{{test.stdout, 0}}}

		repo, created, err := EnsurePrivate(context.Background(), gh, test.owner, test.name)

		if err != nil || created || repo.SSHURL != test.url || len(gh.calls) != 1 {
			t.Fatalf("%s/%s = %#v, %v, %v; %d calls", test.owner, test.name, repo, created, err, len(gh.calls))
		}
		if args := gh.calls[0].args; !reflect.DeepEqual(args, readArgs(test.owner, test.name)) {
			t.Fatalf("read args %q", args)
		}
	}
}

func TestEnsurePrivateRejectsVisibleMismatchAsConflictWithoutCreate(t *testing.T) {
	for _, test := range []struct{ stdout, reason string }{
		{observed("Org/svc", "PUBLIC", "true", "git@github.com:Org/svc.git"), "is PUBLIC"},
		{observed("Org/svc", "INTERNAL", "true", "git@github.com:Org/svc.git"), "is INTERNAL"},
		{observed("Org/svc", "PRIVATE", "false", "git@github.com:Org/svc.git"), "is not empty"},
		{observed("Org/other", "PRIVATE", "true", "git@github.com:Org/other.git"), "identity differs"},
		{observed("Org/svc", "PRIVATE", "true", "git@github.com:org/svc.git"), "SSH URL differs"},
		{observed("Org/svc", "PRIVATE", "true", "https://github.com/Org/svc.git"), "SSH URL differs"},
	} {
		gh := &fakeGH{replies: []reply{{test.stdout, 0}}}

		_, created, err := EnsurePrivate(context.Background(), gh, "Org", "svc")

		if !errors.Is(err, ErrConflict) || created || len(gh.calls) != 1 || err.Error() != "conflict: GitHub repository Org/svc "+test.reason {
			t.Fatalf("%s: %v, created %v, calls %d", test.stdout, err, created, len(gh.calls))
		}
	}
}

func TestEnsurePrivateTreatsUnexpectedReadAsUnknownWithoutCreate(t *testing.T) {
	forbidden := strings.Replace(notFound, "NOT_FOUND", "FORBIDDEN", 1)
	for _, test := range []reply{
		{repository(`"nameWithOwner":"Org/svc","visibility":"PRIVATE","isEmpty":"true","sshUrl":"git@github.com:Org/svc.git"`), 0},
		{repository(`"nameWithOwner":"Org/svc","visibility":"PRIVATE","isEmpty":null,"sshUrl":"git@github.com:Org/svc.git"`), 0},
		{repository(`"nameWithOwner":"Org/svc","visibility":"PRIVATE","sshUrl":"git@github.com:Org/svc.git"`), 0},
		{repository(`"nameWithOwner":"Org/svc","visibility":"PRIVATE","isEmpty":true,"sshUrl":"git@github.com:Org/svc.git","url":"x"`), 0},
		{repository(`"nameWithOwner":"","visibility":"PRIVATE","isEmpty":true,"sshUrl":"git@github.com:/.git"`), 0},
		{repository(`"nameWithOwner":"Org/svc","visibility":"SECRET","isEmpty":true,"sshUrl":"git@github.com:Org/svc.git"`), 0},
		{`{"data":{}}`, 0},
		{`{"data":null}`, 0},
		{`{"data":{"repository":null}}`, 0},
		{strings.TrimSuffix(privateEmpty, "\n") + `,"extensions":{}}`, 0},
		{privateEmpty + privateEmpty, 0},
		{privateEmpty + "trailing", 0},
		{privateEmpty, 1},
		{privateEmpty, 2},
		{privateEmpty, 4},
		{notFound, 0},
		{notFound, 4},
		{forbidden, 1},
		{strings.Replace(notFound, `}]}`, `},{"type":"NOT_FOUND","path":["repository"]}]}`, 1), 1},
		{strings.Replace(notFound, `["repository"]`, `["repository","owner"]`, 1), 1},
		{strings.Replace(notFound, `"message"`, `"x":true,"message"`, 1), 1},
		{strings.Replace(notFound, `"type"`, `"Type"`, 1), 1},
		{strings.TrimSuffix(privateEmpty, "}\n") + `,"errors":[]}`, 0},
		{strings.Replace(notFound, `"repository":null`, `"repository":null,"viewer":{}`, 1), 1},
		{strings.Replace(notFound, `{"data":{"repository":null},`, `{"data":{"repository":null},"x":1,`, 1), 1},
		{"", 0},
		{"", 1},
		{"not json", 1},
		{"null", 0},
	} {
		gh := &fakeGH{replies: []reply{test}}

		_, created, err := EnsurePrivate(context.Background(), gh, "Org", "svc")

		if !errors.Is(err, ErrUnknown) || created || len(gh.calls) != 1 {
			t.Fatalf("%q exit %d: %v, created %v, calls %d", test.stdout, test.code, err, created, len(gh.calls))
		}
		assertSecretFree(t, err)
	}
}

func TestEnsurePrivateCreatesOnceWhenNotFound(t *testing.T) {
	gh := &fakeGH{replies: []reply{{notFound, 1}, {"https://github.com/Org/svc\n", 0}, {privateEmpty, 0}}}

	repo, created, err := EnsurePrivate(context.Background(), gh, "Org", "svc")

	if err != nil || !created || repo.SSHURL != "git@github.com:Org/svc.git" || len(gh.calls) != 3 {
		t.Fatalf("EnsurePrivate = %#v, %v, %v; calls %d", repo, created, err, len(gh.calls))
	}
	want := [][]string{readArgs("Org", "svc"), {"repo", "create", "Org/svc", "--private"}, readArgs("Org", "svc")}
	for i, c := range gh.calls {
		if !reflect.DeepEqual(c.args, want[i]) {
			t.Fatalf("call %d args %q, want %q", i, c.args, want[i])
		}
	}
}

func TestEnsurePrivateReportsCreateAndReadbackFailuresAsUnknown(t *testing.T) {
	public := observed("Org/svc", "PUBLIC", "true", "git@github.com:Org/svc.git")
	for _, test := range []struct {
		create   int
		readback reply
		message  string
	}{
		{1, reply{notFound, 1}, "unknown: create Org/svc: exit 1; readback: not visible"},
		{1, reply{privateEmpty, 0}, "unknown: create Org/svc: exit 1; readback: private empty"},
		{1, reply{public, 0}, "unknown: create Org/svc: exit 1; readback: is PUBLIC"},
		{4, reply{"garbage", 0}, "unknown: create Org/svc: exit 4; readback: unknown"},
		{0, reply{notFound, 1}, "unknown: readback Org/svc: not visible"},
		{0, reply{public, 0}, "unknown: readback Org/svc: is PUBLIC"},
	} {
		gh := &fakeGH{replies: []reply{{notFound, 1}, {"", test.create}, test.readback}}

		_, created, err := EnsurePrivate(context.Background(), gh, "Org", "svc")

		if !errors.Is(err, ErrUnknown) || created || gh.creates() != 1 || len(gh.calls) != 3 || err.Error() != test.message {
			t.Fatalf("%v, created %v, creates %d; want %q", err, created, gh.creates(), test.message)
		}
	}
}

func TestEnsurePrivateReadsBackWithFreshDeadlineAfterCanceledCreate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gh := &fakeGH{replies: []reply{{notFound, 1}, {"", 0}, {privateEmpty, 0}}, cancel: cancel}

	_, created, err := EnsurePrivate(ctx, gh, "Org", "svc")

	if !errors.Is(err, ErrUnknown) || created || gh.creates() != 1 || err.Error() != "unknown: create Org/svc: canceled; readback: private empty" {
		t.Fatalf("EnsurePrivate = %v, %v; creates %d", created, err, gh.creates())
	}
	readback := gh.calls[2]
	if readback.ctxErr != nil || !readback.hasDeadline || time.Until(readback.deadline) > readbackTimeout {
		t.Fatalf("readback context err %v, deadline %v %v", readback.ctxErr, readback.deadline, readback.hasDeadline)
	}
}

func TestEnsurePrivateRejectsInvalidInputWithoutCalls(t *testing.T) {
	gh := &fakeGH{}
	for _, test := range []struct{ owner, name string }{
		{"", "svc"}, {"-x", "svc"}, {"a--b", "svc"}, {"a-", "svc"}, {"a_b", "svc"}, {"github.com/Org", "svc"},
		{strings.Repeat("a", 40), "svc"}, {strings.Repeat("a", 20) + "-" + strings.Repeat("b", 19), "svc"},
		{"Org", ""}, {"Org", "."}, {"Org", ".."}, {"Org", "x.git"}, {"Org", "x.GIT"}, {"Org", "a/b"}, {"Org", "-x"},
		{"Org", strings.Repeat("a", 101)}, {"Org", "æ"},
	} {
		if _, _, err := EnsurePrivate(context.Background(), gh, test.owner, test.name); !errors.Is(err, ErrInvalid) {
			t.Fatalf("%q/%q: %v, want ErrInvalid", test.owner, test.name, err)
		}
	}
	var noContext context.Context
	if _, _, err := EnsurePrivate(noContext, gh, "Org", "svc"); !errors.Is(err, ErrInvalid) {
		t.Fatalf("nil context: %v", err)
	}
	if _, _, err := EnsurePrivate(context.Background(), nil, "Org", "svc"); !errors.Is(err, ErrInvalid) {
		t.Fatalf("nil runner: %v", err)
	}
	if len(gh.calls) != 0 {
		t.Fatalf("%d calls for invalid input", len(gh.calls))
	}
}

func TestEnsurePrivateFiltersEnvironmentForEveryCall(t *testing.T) {
	for _, name := range []string{"GH_REPO", "GH_DEBUG", "DEBUG", "GH_FORCE_TTY", "CLICOLOR_FORCE", "GH_PAGER", "PAGER", "GIT_DIR", "GIT_WORK_TREE", "GIT_CONFIG_GLOBAL", "GIT_CONFIG_COUNT", "NO_COLOR"} {
		t.Setenv(name, "inherited")
	}
	t.Setenv("GH_HOST", "example.com")
	t.Setenv("GH_CONFIG_DIR", "/private/gh-config")
	gh := &fakeGH{replies: []reply{{notFound, 1}, {"", 0}, {privateEmpty, 0}}}

	if _, _, err := EnsurePrivate(context.Background(), gh, "Org", "svc"); err != nil {
		t.Fatal(err)
	}
	for i, c := range gh.calls {
		env := "\n" + strings.Join(c.env, "\n") + "\n"
		if strings.Contains(env, "inherited") || strings.Contains(env, "example.com") {
			t.Fatalf("call %d env keeps an inherited value that must be dropped", i)
		}
		if !strings.Contains(env, "\nGH_CONFIG_DIR=/private/gh-config\n") {
			t.Fatalf("call %d env lacks the inherited GH_CONFIG_DIR", i)
		}
		for _, entry := range []string{"GH_HOST=github.com", "GH_PROMPT_DISABLED=1", "NO_COLOR=1", "GH_NO_UPDATE_NOTIFIER=1", "GH_NO_EXTENSION_UPDATE_NOTIFIER=1"} {
			if strings.Count(env, "\n"+entry+"\n") != 1 {
				t.Fatalf("call %d env lacks exactly one %s", i, entry)
			}
		}
	}
}

func TestValidMatchesEnsurePrivateInputRule(t *testing.T) {
	for _, test := range []struct {
		owner, name string
		want        bool
	}{
		{"Org", "svc", true}, {"a-b", "x.y_z-1", true}, {"Org", ".hidden", true},
		{"", "svc", false}, {"-x", "svc", false}, {strings.Repeat("a", 40), "svc", false},
		{"Org", "", false}, {"Org", ".", false}, {"Org", "..", false}, {"Org", "x.GIT", false}, {"Org", "a/b", false},
	} {
		if got := Valid(test.owner, test.name); got != test.want {
			t.Fatalf("Valid(%q, %q) = %v", test.owner, test.name, got)
		}
	}
}
