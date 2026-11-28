package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/ghapi"
	"github.com/devdimensionlab/plybuild/internal/githubbilling"
	"github.com/devdimensionlab/plybuild/internal/githubprs"
	"github.com/devdimensionlab/plybuild/internal/githubrepo"
	"github.com/devdimensionlab/plybuild/internal/quotahistory"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
	"github.com/spf13/cobra"
)

// githubConfigFile holds the included Actions minutes below the user's home.
const githubConfigFile = ".config/ply/github.json"

// quotaLowPercent is the remaining share below which a quota's level is low.
const quotaLowPercent = 20

var errGitHubActionsFailed = errors.New("reading GitHub Actions usage failed")

// githubActionsEnv is the billing reader, clock and home directory of ply github actions,
// and the pull request reader of ply github prs.
type githubActionsEnv struct {
	read    func(ctx context.Context, year int, month time.Month) (githubbilling.Usage, error)
	prs     func(ctx context.Context) (githubprs.Result, error)
	list    func(ctx context.Context) (githubprs.ListResult, error)
	details func(ctx context.Context) ([]githubprs.Detailed, error)
	now     func() time.Time
	home    func() (string, error)
	// record appends a point to the quota history; nil records nothing.
	record func(series string, p quotahistory.Point) error
}

func systemGitHubActionsEnv(runner process.Runner) githubActionsEnv {
	return githubActionsEnv{
		read: func(ctx context.Context, year int, month time.Month) (githubbilling.Usage, error) {
			return githubbilling.Read(ctx, runner, year, month)
		},
		prs:     func(ctx context.Context) (githubprs.Result, error) { return githubprs.Read(ctx, runner) },
		list:    func(ctx context.Context) (githubprs.ListResult, error) { return githubprs.List(ctx, runner) },
		details: func(ctx context.Context) ([]githubprs.Detailed, error) { return githubprs.Details(ctx, runner) },
		now:     time.Now,
		home:    os.UserHomeDir,
		record:  recordQuotaHistory,
	}
}

type githubActionsError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type githubActionsQuota struct {
	Kind             string              `json:"kind"`
	SchemaVersion    int                 `json:"schema_version"`
	ID               string              `json:"id"`
	Source           string              `json:"source"`
	Status           string              `json:"status"`
	Error            *githubActionsError `json:"error,omitempty"`
	ObservedAt       string              `json:"observed_at,omitempty"`
	Unit             string              `json:"unit,omitempty"`
	Used             *int                `json:"used,omitempty"`
	Limit            *int                `json:"limit,omitempty"`
	Remaining        *int                `json:"remaining,omitempty"`
	RemainingPercent *int                `json:"remaining_percent,omitempty"`
	ResetsAt         string              `json:"resets_at,omitempty"`
	Level            string              `json:"level,omitempty"`
}

type githubActionsUsage struct {
	Kind          string               `json:"kind"`
	SchemaVersion int                  `json:"schema_version"`
	Source        string               `json:"source"`
	Status        string               `json:"status"`
	Error         *githubActionsError  `json:"error,omitempty"`
	ObservedAt    string               `json:"observed_at,omitempty"`
	PeriodStart   string               `json:"period_start,omitempty"`
	Unit          string               `json:"unit,omitempty"`
	Repos         *[]githubActionsRepo `json:"repos,omitempty"`
}

type githubActionsRepo struct {
	Key          string    `json:"key"`
	ProjectIDs   *[]string `json:"project_ids,omitempty"`
	Minutes      int       `json:"minutes"`
	MacOSMinutes int       `json:"macos_minutes"`
}

func newGitHubCommand(env githubActionsEnv) *cobra.Command {
	github := &cobra.Command{
		Use:               "github",
		Short:             "Read what Ply needs from GitHub",
		Args:              cobra.NoArgs,
		RunE:              func(c *cobra.Command, args []string) error { return c.Help() },
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	}
	actions := &cobra.Command{
		Use:   "actions",
		Short: "Read the GitHub Actions minutes used this month",
		Long: "Read the GitHub Actions minutes of the user gh is logged in as, from GitHub's billing usage (gh scope user).\n" +
			"Each command reads GitHub when it runs; quota also adds the month's minutes to this machine's quota history (ply quota history). Minutes count against the included quota as\n" +
			"Linux ×1, Windows ×2 and macOS ×10; storage and other runners are not counted. The included minutes per\n" +
			"month are read from ~/" + githubConfigFile + `: {"actions": {"included_minutes": 3000}}.`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error { return c.Help() },
	}
	var format string
	add := func(c *cobra.Command) {
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		c.Args = func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("expected no positional arguments")
			}
			if format != "text" && format != "json" {
				return errors.New("--format must be text or json")
			}
			return nil
		}
		actions.AddCommand(c)
	}
	add(&cobra.Command{
		Use:     "quota",
		Short:   "Show the included Actions minutes used and left this month, with a level",
		Long:    "Show the Actions minutes used this month against the included minutes, what is left and the level:\nexhausted at 0% left, low below 20%, else ok. Without the included minutes, limit, remaining and level are left out.\nThe exit status is 1 when GitHub could not be read.",
		Example: "  ply github actions quota\n  ply github actions quota --format json",
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, err := githubActionsLimit(env)
			if err != nil {
				return err
			}
			report, usage, failed := githubActionsRead(cmd.Context(), env)
			quota := githubActionsQuota{Kind: "PlyGitHubActionsQuota@1", SchemaVersion: 1, ID: "github.actions", Source: "github", Status: report.Status, Error: report.Error}
			if failed == nil {
				quota.ObservedAt, quota.Unit, quota.ResetsAt = report.ObservedAt, report.Unit, githubActionsTime(report.period.AddDate(0, 1, 0))
				used := 0
				for _, r := range usage {
					used += r.Minutes
				}
				quota.Used = &used
				if limit > 0 {
					remaining := max(0, limit-used)
					percent := remaining * 100 / limit
					quota.Limit, quota.Remaining, quota.RemainingPercent, quota.Level = &limit, &remaining, &percent, quotaLevel(remaining <= 0, percent)
				}
				recordActionsHistory(env.record, quota)
			}
			if err := githubActionsWrite(cmd.OutOrStdout(), format, quota, githubActionsQuotaText(quota)); err != nil {
				return err
			}
			return failed
		},
	})
	add(&cobra.Command{
		Use:     "usage",
		Short:   "List the Actions minutes each repository used this month, with its Projects",
		Long:    "List the repositories with Actions minutes this month, most first, with the IDs of the Projects in ply-workspace.json\nwhose repository has that GitHub remote. Outside a workspace project_ids is left out. The minutes of all rows add up to\nused in ply github actions quota. The exit status is 1 when GitHub could not be read.",
		Example: "  ply github actions usage\n  ply github actions usage --format json",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, repos, failed := githubActionsRead(cmd.Context(), env)
			if failed == nil {
				projects, err := githubProjectsByRemote()
				if err != nil {
					return err
				}
				for i := range repos {
					if projects != nil {
						ids := append([]string{}, projects[strings.ToLower(repos[i].Key)]...)
						repos[i].ProjectIDs = &ids
					}
				}
				report.PeriodStart, report.Repos = githubActionsTime(report.period), &repos
			}
			if err := githubActionsWrite(cmd.OutOrStdout(), format, report.githubActionsUsage, githubActionsUsageText(report.githubActionsUsage)); err != nil {
				return err
			}
			return failed
		},
	})
	github.AddCommand(actions)
	github.AddCommand(newGitHubPRsCommand(env))
	return github
}

// githubError is the code and description of a failed ply github read.
func githubError(err error) *githubActionsError {
	code := ghapi.ErrFailed
	for _, known := range []error{ghapi.ErrMissing, ghapi.ErrAuth, ghapi.ErrScope, ghapi.ErrTimeout} {
		if errors.Is(err, known) {
			code = known
		}
	}
	return &githubActionsError{Code: code.Error(), Description: strings.TrimPrefix(err.Error(), code.Error()+": ")}
}

type githubActionsReport struct {
	githubActionsUsage
	period time.Time
}

// githubActionsRead reads this month's billing usage and sums the counted
// minutes per repository, most first; a failure is reported in the result.
func githubActionsRead(ctx context.Context, env githubActionsEnv) (githubActionsReport, []githubActionsRepo, error) {
	now := env.now().UTC()
	period := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	report := githubActionsReport{githubActionsUsage: githubActionsUsage{Kind: "PlyGitHubActionsUsage@1", SchemaVersion: 1, Source: "github"}, period: period}
	usage, err := env.read(ctx, period.Year(), period.Month())
	if err != nil {
		report.Status, report.Error = "error", githubError(err)
		return report, nil, errGitHubActionsFailed
	}
	report.Status, report.ObservedAt, report.Unit = "ok", githubActionsTime(now), "minutes"
	type sum struct{ minutes, macOS float64 }
	sums := map[string]*sum{}
	for _, item := range usage.Items {
		multiplier := githubActionsMultiplier(item)
		if multiplier == 0 {
			continue
		}
		key := usage.Login + "/" + item.Repository
		if sums[key] == nil {
			sums[key] = &sum{}
		}
		sums[key].minutes += item.Quantity * multiplier
		if multiplier == 10 {
			sums[key].macOS += item.Quantity * multiplier
		}
	}
	repos := []githubActionsRepo{}
	for key, s := range sums {
		if minutes := int(math.Round(s.minutes)); minutes > 0 {
			repos = append(repos, githubActionsRepo{Key: key, Minutes: minutes, MacOSMinutes: int(math.Round(s.macOS))})
		}
	}
	sort.Slice(repos, func(i, j int) bool {
		if repos[i].Minutes != repos[j].Minutes {
			return repos[i].Minutes > repos[j].Minutes
		}
		return repos[i].Key < repos[j].Key
	})
	return report, repos, nil
}

// githubActionsMultiplier is how many included minutes one minute of the row
// uses; 0 for rows that do not count.
func githubActionsMultiplier(item githubbilling.Item) float64 {
	if item.Product != "actions" || item.UnitType != "Minutes" {
		return 0
	}
	switch {
	case strings.HasPrefix(item.SKU, "Actions Linux"):
		return 1
	case strings.HasPrefix(item.SKU, "Actions Windows"):
		return 2
	case strings.HasPrefix(item.SKU, "Actions macOS"):
		return 10
	}
	return 0
}

// quotaLevel is the level of every Ply quota: exhausted only when nothing is
// left, even when the floored percent is already 0.
func quotaLevel(nothingLeft bool, remainingPercent int) string {
	switch {
	case nothingLeft:
		return "exhausted"
	case remainingPercent < quotaLowPercent:
		return "low"
	}
	return "ok"
}

// githubActionsLimit reads the included minutes; 0 when not configured.
func githubActionsLimit(env githubActionsEnv) (int, error) {
	home, err := env.home()
	if err != nil {
		return 0, fmt.Errorf("home directory for %s: %v", githubConfigFile, err)
	}
	path := filepath.Join(home, filepath.FromSlash(githubConfigFile))
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	var config struct {
		Actions struct {
			IncludedMinutes int `json:"included_minutes"`
		} `json:"actions"`
	}
	if err == nil {
		err = json.Unmarshal(data, &config)
	}
	if err != nil || config.Actions.IncludedMinutes < 0 {
		return 0, fmt.Errorf("%s is not readable as {\"actions\": {\"included_minutes\": <minutes>}}", path)
	}
	return config.Actions.IncludedMinutes, nil
}

// githubProjectsByRemote maps lower-case owner/repo to the IDs of the Projects
// that have it as a repository remote; nil outside a workspace.
func githubProjectsByRemote() (map[string][]string, error) {
	observed, err := workspace.ObserveContaining(workspace.SystemDependencies())
	var project *workspace.ProjectError
	if errors.As(err, &project) && project.Class == workspace.ErrorWorkspaceNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defs, _, err := workspacedef.Read(observed.Root)
	if err != nil || defs == nil {
		return map[string][]string{}, err
	}
	projects := map[string][]string{}
	for _, p := range defs.Projects {
		for _, r := range p.Repos {
			if owner, name, ok := githubrepo.FromRemote(r.Remote); ok {
				key := strings.ToLower(owner + "/" + name)
				if n := len(projects[key]); n == 0 || projects[key][n-1] != p.ID {
					projects[key] = append(projects[key], p.ID)
				}
			}
		}
	}
	return projects, nil
}

func githubActionsWrite(out io.Writer, format string, value interface{}, text string) error {
	if format == "json" {
		return json.NewEncoder(out).Encode(value)
	}
	_, err := io.WriteString(out, text)
	return err
}

func githubActionsQuotaText(q githubActionsQuota) string {
	switch {
	case q.Error != nil:
		return fmt.Sprintf("GitHub Actions: %s (%s)\n", q.Error.Description, q.Error.Code)
	case q.Limit == nil:
		return fmt.Sprintf("GitHub Actions: %d minutes used, no included minutes configured in ~/%s\n", *q.Used, githubConfigFile)
	}
	return fmt.Sprintf("GitHub Actions: %d of %d minutes used, %d left (%d%%, %s), resets %s\n",
		*q.Used, *q.Limit, *q.Remaining, *q.RemainingPercent, q.Level, q.ResetsAt[:len(time.DateOnly)])
}

func githubActionsUsageText(u githubActionsUsage) string {
	if u.Error != nil {
		return fmt.Sprintf("GitHub Actions: %s (%s)\n", u.Error.Description, u.Error.Code)
	}
	var b strings.Builder
	for _, r := range *u.Repos {
		fmt.Fprintf(&b, "%s  %d", r.Key, r.Minutes)
		if r.MacOSMinutes > 0 {
			fmt.Fprintf(&b, " (macOS %d)", r.MacOSMinutes)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func githubActionsTime(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func init() {
	RootCmd.AddCommand(newGitHubCommand(systemGitHubActionsEnv(process.SystemRunner().Runner)))
}
