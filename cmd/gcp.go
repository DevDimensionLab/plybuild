package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/gcpbilling"
	"github.com/spf13/cobra"
)

// gcpConfigFile holds the billing export and the Project links below the user's home.
const gcpConfigFile = ".config/ply/gcp.json"

// gcpBillingStateFile is the snapshot below $XDG_STATE_HOME, by default ~/.local/state.
const gcpBillingStateFile = "ply/gcp-billing.json"

// gcpBillingStaleAfter is the age of the last sync after which values are stale.
const gcpBillingStaleAfter = 24 * time.Hour

// gcpBillingZone is the time zone of Google's invoice month and budget period.
const gcpBillingZone = "America/Los_Angeles"

var errGCPBillingFailed = errors.New("reading Google Cloud billing failed")

// gcpBillingEnv is the export reader, clock, home directory and environment of ply gcp billing.
type gcpBillingEnv struct {
	read   func(ctx context.Context, src gcpbilling.Source) (gcpbilling.Snapshot, error)
	now    func() time.Time
	home   func() (string, error)
	getenv func(string) string
}

func systemGCPBillingEnv(runner process.Runner) gcpBillingEnv {
	return gcpBillingEnv{
		read: func(ctx context.Context, src gcpbilling.Source) (gcpbilling.Snapshot, error) {
			return gcpbilling.Read(ctx, runner, src)
		},
		now:    time.Now,
		home:   os.UserHomeDir,
		getenv: os.Getenv,
	}
}

type gcpBillingError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type gcpBillingSync struct {
	Kind          string           `json:"kind"`
	SchemaVersion int              `json:"schema_version"`
	Source        string           `json:"source"`
	Status        string           `json:"status"`
	Error         *gcpBillingError `json:"error,omitempty"`
	SyncedAt      string           `json:"synced_at,omitempty"`
	Months        *[]string        `json:"months,omitempty"`
	Currency      string           `json:"currency,omitempty"`
	Budgets       *int             `json:"budgets,omitempty"`
}

type gcpBillingQuota struct {
	Kind          string                `json:"kind"`
	SchemaVersion int                   `json:"schema_version"`
	Source        string                `json:"source"`
	Status        string                `json:"status"`
	Error         *gcpBillingError      `json:"error,omitempty"`
	SyncedAt      string                `json:"synced_at,omitempty"`
	Stale         *bool                 `json:"stale,omitempty"`
	Resources     *[]gcpBillingResource `json:"resources,omitempty"`
}

type gcpBillingResource struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	ProjectIDs       []string `json:"project_ids"`
	ObservedAt       string   `json:"observed_at,omitempty"`
	Unit             string   `json:"unit"`
	Currency         string   `json:"currency"`
	Window           string   `json:"window"`
	Used             float64  `json:"used"`
	Limit            float64  `json:"limit"`
	Remaining        float64  `json:"remaining"`
	RemainingPercent int      `json:"remaining_percent"`
	ResetsAt         string   `json:"resets_at"`
	Level            string   `json:"level"`
}

type gcpBillingUsage struct {
	Kind          string               `json:"kind"`
	SchemaVersion int                  `json:"schema_version"`
	Source        string               `json:"source"`
	Status        string               `json:"status"`
	Error         *gcpBillingError     `json:"error,omitempty"`
	SyncedAt      string               `json:"synced_at,omitempty"`
	Stale         *bool                `json:"stale,omitempty"`
	Month         string               `json:"month,omitempty"`
	ObservedAt    string               `json:"observed_at,omitempty"`
	Unit          string               `json:"unit,omitempty"`
	Currency      string               `json:"currency,omitempty"`
	Amount        *float64             `json:"amount,omitempty"`
	Tax           *float64             `json:"tax,omitempty"`
	Projects      *[]gcpBillingProject `json:"projects,omitempty"`
}

type gcpBillingProject struct {
	Key        string              `json:"key,omitempty"`
	ProjectIDs []string            `json:"project_ids"`
	Amount     float64             `json:"amount"`
	Services   []gcpBillingService `json:"services"`
}

type gcpBillingService struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

// gcpConfig is ~/.config/ply/gcp.json.
type gcpConfig struct {
	Billing struct {
		Account  string `json:"account"`
		Table    string `json:"table"`
		Location string `json:"location"`
	} `json:"billing"`
	Projects map[string][]string `json:"projects"`
}

func (c gcpConfig) source() gcpbilling.Source {
	return gcpbilling.Source{Account: c.Billing.Account, Table: c.Billing.Table, Location: c.Billing.Location}
}

// projectIDs are the Ply Projects linked to a GCP project; empty when none.
func (c gcpConfig) projectIDs(gcpProject string) []string {
	return append([]string{}, c.Projects[gcpProject]...)
}

func newGCPCommand(env gcpBillingEnv) *cobra.Command {
	gcp := &cobra.Command{
		Use:               "gcp",
		Short:             "Read what Ply needs from Google Cloud",
		Args:              cobra.NoArgs,
		RunE:              func(c *cobra.Command, args []string) error { return c.Help() },
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	}
	billing := &cobra.Command{
		Use:   "billing",
		Short: "Sync and read the Google Cloud cost and budgets of one billing account",
		Long: "Sync and read the Cloud Billing export in BigQuery and the billing account's budgets. Only sync contacts Google,\n" +
			"with the user's own gcloud and bq login; it replaces the local snapshot in $XDG_STATE_HOME/" + gcpBillingStateFile + "\n" +
			"(default ~/.local/state). quota and usage read only that snapshot; values are stale 24 hours after the last sync.\n" +
			"Ply never syncs by itself. Amounts are cost plus credits in Google's invoice month (Los Angeles time), without tax.\n" +
			"The export and the links from GCP projects to Ply Projects are read from ~/" + gcpConfigFile + ":\n" +
			`{"billing": {"account": "XXXXXX-XXXXXX-XXXXXX", "table": "project.dataset.table", "location": "EU"},` + "\n" +
			` "projects": {"<gcp project id>": ["<ply project id>"]}}`,
		Args: cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error { return c.Help() },
	}
	var format, month string
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
		billing.AddCommand(c)
	}
	add(&cobra.Command{
		Use:     "sync",
		Short:   "Read the billing export and budgets from Google and replace the local snapshot",
		Long:    "Read the whole billing export summed per invoice month, project, service and cost type, the budgets and the\nproject numbers, and replace the local snapshot. On any failure the old snapshot stays. The exit status is 1 when\nGoogle could not be read or the configuration is invalid.",
		Example: "  ply gcp billing sync\n  ply gcp billing sync --format json",
		RunE: func(cmd *cobra.Command, args []string) error {
			report := gcpBillingSync{Kind: "PlyGcpBillingSync@1", SchemaVersion: 1, Source: "gcp"}
			config, path, failed, err := gcpBillingOpen(env)
			if err == nil && failed == nil && config != nil {
				failed, err = gcpBillingSyncSnapshot(cmd.Context(), env, *config, path, &report)
			}
			if err != nil {
				return err
			}
			gcpBillingStatus(&report.Status, &report.Error, config, failed)
			if err := githubActionsWrite(cmd.OutOrStdout(), format, report, gcpBillingSyncText(report)); err != nil {
				return err
			}
			return gcpBillingExit(failed)
		},
	})
	add(&cobra.Command{
		Use:     "quota",
		Short:   "Show what is used and left of each budget this month, with a level",
		Long:    "Show each budget for one project or the whole account against this month's amount from the local snapshot:\nexhausted when nothing is left, low below 20% left, else ok. Budgets with other filters, other credit treatment,\nanother period or last period's amount are left out. Of two budgets with the same id the lowest is used.",
		Example: "  ply gcp billing quota\n  ply gcp billing quota --format json",
		RunE: func(cmd *cobra.Command, args []string) error {
			report := gcpBillingQuota{Kind: "PlyGcpBillingQuota@1", SchemaVersion: 1, Source: "gcp"}
			config, snapshot, failed, err := gcpBillingLoad(env)
			if err != nil {
				return err
			}
			if failed == nil && snapshot != nil {
				period, err := gcpBillingCurrentPeriod(env)
				if err != nil {
					return err
				}
				resources := gcpBillingResources(*config, *snapshot, period)
				report.SyncedAt, report.Stale, report.Resources = githubActionsTime(snapshot.SyncedAt), gcpBillingStale(env, *snapshot), &resources
			}
			gcpBillingStatus(&report.Status, &report.Error, config, failed)
			if report.Status == "ok" && snapshot == nil {
				report.Status = "pending"
			}
			if err := githubActionsWrite(cmd.OutOrStdout(), format, report, gcpBillingQuotaText(report)); err != nil {
				return err
			}
			return gcpBillingExit(failed)
		},
	})
	usage := &cobra.Command{
		Use:     "usage",
		Short:   "List the month's amount per GCP project and service, with its Projects",
		Long:    "List the amount of each GCP project and its services in one invoice month from the local snapshot, most first,\nwith the Ply Projects linked in ~/" + gcpConfigFile + ". Tax and rounding are summed apart as tax. The amount\nis the sum of the projects. --month defaults to the current invoice month.",
		Example: "  ply gcp billing usage\n  ply gcp billing usage --month 2026-10 --format json",
		RunE: func(cmd *cobra.Command, args []string) error {
			period, err := gcpBillingCurrentPeriod(env)
			if err != nil {
				return err
			}
			if month != "" {
				if period, err = time.ParseInLocation("2006-01", month, period.Location()); err != nil {
					return errors.New("--month must be YYYY-MM")
				}
			}
			report := gcpBillingUsage{Kind: "PlyGcpBillingUsage@1", SchemaVersion: 1, Source: "gcp"}
			config, snapshot, failed, err := gcpBillingLoad(env)
			if err != nil {
				return err
			}
			if failed == nil && snapshot != nil {
				report.SyncedAt, report.Stale, report.Month = githubActionsTime(snapshot.SyncedAt), gcpBillingStale(env, *snapshot), period.Format("2006-01")
				gcpBillingFillUsage(&report, *config, *snapshot, period)
			}
			gcpBillingStatus(&report.Status, &report.Error, config, failed)
			if report.Status == "ok" && snapshot == nil {
				report.Status = "pending"
			}
			if err := githubActionsWrite(cmd.OutOrStdout(), format, report, gcpBillingUsageText(report)); err != nil {
				return err
			}
			return gcpBillingExit(failed)
		},
	}
	usage.Flags().StringVar(&month, "month", "", "invoice month YYYY-MM (default: the current one)")
	add(usage)
	gcp.AddCommand(billing)
	return gcp
}

// gcpBillingOpen reads the configuration and resolves the snapshot path. A
// nil configuration without failure means not configured.
func gcpBillingOpen(env gcpBillingEnv) (*gcpConfig, string, *gcpBillingError, error) {
	home, err := env.home()
	if err != nil {
		return nil, "", nil, fmt.Errorf("home directory for %s: %v", gcpConfigFile, err)
	}
	state := env.getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(state) {
		state = filepath.Join(home, ".local", "state")
	}
	path := filepath.Join(state, filepath.FromSlash(gcpBillingStateFile))
	data, err := os.ReadFile(filepath.Join(home, filepath.FromSlash(gcpConfigFile)))
	if errors.Is(err, os.ErrNotExist) {
		return nil, path, nil, nil
	}
	var config gcpConfig
	if err == nil {
		err = json.Unmarshal(data, &config)
	}
	if err == nil {
		err = config.source().Valid()
	}
	if err != nil {
		return nil, path, &gcpBillingError{Code: "config_invalid",
			Description: "~/" + gcpConfigFile + " does not give a valid billing account, export table and location"}, nil
	}
	return &config, path, nil, nil
}

// gcpBillingLoad reads the configuration and the snapshot; a nil snapshot
// without failure means not configured or not synced yet.
func gcpBillingLoad(env gcpBillingEnv) (*gcpConfig, *gcpbilling.Snapshot, *gcpBillingError, error) {
	config, path, failed, err := gcpBillingOpen(env)
	if err != nil || failed != nil || config == nil {
		return config, nil, failed, err
	}
	snapshot, err := gcpbilling.Load(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return config, nil, nil, nil
	case err != nil:
		return config, nil, &gcpBillingError{Code: gcpbilling.ErrState.Error(),
			Description: "the local snapshot is unreadable; run ply gcp billing sync"}, nil
	}
	return config, &snapshot, nil, nil
}

// gcpBillingSyncSnapshot reads Google and replaces the snapshot under the
// sync lock; a read failure is reported in the result.
func gcpBillingSyncSnapshot(ctx context.Context, env gcpBillingEnv, config gcpConfig, path string, report *gcpBillingSync) (failure *gcpBillingError, err error) {
	lock, err := gcpbilling.Lock(path)
	if err != nil {
		return nil, fmt.Errorf("lock the GCP billing snapshot: %v", err)
	}
	defer func() {
		if releaseErr := lock.Release(); releaseErr != nil && failure == nil && err == nil {
			err = releaseErr
		}
	}()
	snapshot, err := env.read(ctx, config.source())
	if err != nil {
		code := gcpbilling.ErrFailed
		for _, known := range []error{gcpbilling.ErrMissing, gcpbilling.ErrAuth, gcpbilling.ErrTimeout, gcpbilling.ErrMixedCurrency} {
			if errors.Is(err, known) {
				code = known
			}
		}
		return &gcpBillingError{Code: code.Error(), Description: strings.TrimPrefix(err.Error(), code.Error()+": ")}, nil
	}
	snapshot.SyncedAt = env.now().UTC()
	if err := gcpbilling.Save(path, snapshot); err != nil {
		return nil, fmt.Errorf("save the GCP billing snapshot: %v", err)
	}
	months := map[string]bool{}
	for _, r := range snapshot.Rows {
		if len(r.Month) == 6 {
			months[r.Month[:4]+"-"+r.Month[4:]] = true
		}
	}
	list := []string{}
	for m := range months {
		list = append(list, m)
	}
	sort.Strings(list)
	budgets := len(snapshot.Budgets)
	report.SyncedAt, report.Months, report.Currency, report.Budgets = githubActionsTime(snapshot.SyncedAt), &list, snapshot.Currency, &budgets
	return nil, nil
}

// gcpBillingStatus sets ok, not_configured or error from the reading.
func gcpBillingStatus(status *string, failure **gcpBillingError, config *gcpConfig, failed *gcpBillingError) {
	switch {
	case failed != nil:
		*status, *failure = "error", failed
	case config == nil:
		*status = "not_configured"
	default:
		*status = "ok"
	}
}

func gcpBillingExit(failed *gcpBillingError) error {
	if failed != nil {
		return errGCPBillingFailed
	}
	return nil
}

// gcpBillingCurrentPeriod is the start of the current invoice month.
func gcpBillingCurrentPeriod(env gcpBillingEnv) (time.Time, error) {
	zone, err := time.LoadLocation(gcpBillingZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time zone %s: %v", gcpBillingZone, err)
	}
	now := env.now().In(zone)
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, zone), nil
}

func gcpBillingStale(env gcpBillingEnv, s gcpbilling.Snapshot) *bool {
	stale := env.now().Sub(s.SyncedAt) > gcpBillingStaleAfter
	return &stale
}

// gcpBillingCounted is false for tax and rounding, which budgets leave out.
func gcpBillingCounted(r gcpbilling.Row) bool {
	return r.CostType != "tax" && r.CostType != "rounding_error"
}

func gcpBillingRound(v float64) float64 { return math.Round(v*100) / 100 }

// gcpBillingResources gives each supported budget its use in the period,
// sorted by id; of two budgets with the same id the lowest limit is kept.
func gcpBillingResources(config gcpConfig, snapshot gcpbilling.Snapshot, period time.Time) []gcpBillingResource {
	month := period.Format("200601")
	byID := map[string]gcpBillingResource{}
	for _, b := range snapshot.Budgets {
		if b.OtherFilters || b.Amount == nil || *b.Amount <= 0 || b.CalendarPeriod != "MONTH" ||
			b.CreditTypesTreatment != "INCLUDE_ALL_CREDITS" || len(b.Projects) > 1 {
			continue
		}
		id, project := "gcp.account", ""
		if len(b.Projects) == 1 {
			project = b.Projects[0]
			id = "gcp." + project
		}
		if old, ok := byID[id]; ok && old.Limit <= *b.Amount {
			continue
		}
		used, observed := 0.0, time.Time{}
		for _, r := range snapshot.Rows {
			if r.Month == month && gcpBillingCounted(r) && (project == "" || r.ProjectID == project) {
				used += r.Cost + r.Credits
				if r.LastExport.After(observed) {
					observed = r.LastExport
				}
			}
		}
		resource := gcpBillingResource{ID: id, Name: b.Name, ProjectIDs: config.projectIDs(project), Unit: "amount",
			Currency: b.Currency, Window: "month", Used: gcpBillingRound(used), Limit: *b.Amount, ResetsAt: githubActionsTime(period.AddDate(0, 1, 0))}
		if !observed.IsZero() {
			resource.ObservedAt = githubActionsTime(observed)
		}
		resource.Remaining = gcpBillingRound(max(0, resource.Limit-resource.Used))
		resource.RemainingPercent = int(math.Floor(resource.Remaining / resource.Limit * 100))
		resource.Level = quotaLevel(resource.Remaining <= 0, resource.RemainingPercent)
		byID[id] = resource
	}
	resources := []gcpBillingResource{}
	for _, r := range byID {
		resources = append(resources, r)
	}
	sort.Slice(resources, func(i, j int) bool { return resources[i].ID < resources[j].ID })
	return resources
}

// gcpBillingFillUsage sums the period per project and service, most first.
func gcpBillingFillUsage(report *gcpBillingUsage, config gcpConfig, snapshot gcpbilling.Snapshot, period time.Time) {
	month := period.Format("200601")
	services := map[string]map[string]float64{}
	tax, observed := 0.0, time.Time{}
	for _, r := range snapshot.Rows {
		if r.Month != month {
			continue
		}
		if r.LastExport.After(observed) {
			observed = r.LastExport
		}
		if !gcpBillingCounted(r) {
			tax += r.Cost + r.Credits
			continue
		}
		if services[r.ProjectID] == nil {
			services[r.ProjectID] = map[string]float64{}
		}
		services[r.ProjectID][r.Service] += r.Cost + r.Credits
	}
	projects, total := []gcpBillingProject{}, 0.0
	for key, sums := range services {
		project := gcpBillingProject{Key: key, ProjectIDs: config.projectIDs(key), Services: []gcpBillingService{}}
		raw := 0.0
		for name, amount := range sums {
			raw += amount
			if rounded := gcpBillingRound(amount); rounded != 0 {
				project.Services = append(project.Services, gcpBillingService{Name: name, Amount: rounded})
			}
		}
		if project.Amount = gcpBillingRound(raw); project.Amount == 0 {
			continue
		}
		sort.Slice(project.Services, func(i, j int) bool {
			a, b := project.Services[i], project.Services[j]
			return a.Amount > b.Amount || a.Amount == b.Amount && a.Name < b.Name
		})
		total += project.Amount
		projects = append(projects, project)
	}
	sort.Slice(projects, func(i, j int) bool {
		a, b := projects[i], projects[j]
		return a.Amount > b.Amount || a.Amount == b.Amount && a.Key < b.Key
	})
	amount, roundedTax := gcpBillingRound(total), gcpBillingRound(tax)
	if !observed.IsZero() {
		report.ObservedAt = githubActionsTime(observed)
	}
	report.Unit, report.Currency, report.Amount, report.Tax, report.Projects = "amount", snapshot.Currency, &amount, &roundedTax, &projects
}

// gcpBillingStateText is the line for a reading without values, or "".
func gcpBillingStateText(status string, failure *gcpBillingError, syncedAt string, stale *bool) string {
	switch {
	case failure != nil:
		return fmt.Sprintf("Google Cloud: %s (%s)\n", failure.Description, failure.Code)
	case status == "not_configured":
		return "Google Cloud: not configured; see ply gcp billing --help\n"
	case status == "pending":
		return "Google Cloud: not synced yet; run ply gcp billing sync\n"
	case stale != nil && *stale:
		return fmt.Sprintf("Google Cloud: last synced %s, more than 24 hours ago\n", syncedAt)
	}
	return ""
}

func gcpBillingSyncText(s gcpBillingSync) string {
	if text := gcpBillingStateText(s.Status, s.Error, "", nil); text != "" {
		return text
	}
	return fmt.Sprintf("Google Cloud billing synced: %d months, %d budgets, %s\n", len(*s.Months), *s.Budgets, s.Currency)
}

func gcpBillingQuotaText(q gcpBillingQuota) string {
	text := gcpBillingStateText(q.Status, q.Error, q.SyncedAt, q.Stale)
	if q.Resources == nil {
		return text
	}
	var b strings.Builder
	b.WriteString(text)
	for _, r := range *q.Resources {
		resets, _ := time.Parse(time.RFC3339, r.ResetsAt)
		zone, _ := time.LoadLocation(gcpBillingZone)
		fmt.Fprintf(&b, "Google Cloud %s (%s): %.2f of %.2f %s used, %.2f left (%d%%, %s), resets %s\n",
			r.ID, r.Name, r.Used, r.Limit, r.Currency, r.Remaining, r.RemainingPercent, r.Level, resets.In(zone).Format(time.DateOnly))
	}
	return b.String()
}

func gcpBillingUsageText(u gcpBillingUsage) string {
	text := gcpBillingStateText(u.Status, u.Error, u.SyncedAt, u.Stale)
	if u.Projects == nil {
		return text
	}
	var b strings.Builder
	b.WriteString(text)
	fmt.Fprintf(&b, "Google Cloud %s: %.2f %s, tax %.2f %s\n", u.Month, *u.Amount, u.Currency, *u.Tax, u.Currency)
	for _, p := range *u.Projects {
		key := p.Key
		if key == "" {
			key = "(no project)"
		}
		fmt.Fprintf(&b, "%s  %.2f\n", key, p.Amount)
		for _, s := range p.Services {
			fmt.Fprintf(&b, "  %s  %.2f\n", s.Name, s.Amount)
		}
	}
	return b.String()
}

func init() {
	RootCmd.AddCommand(newGCPCommand(systemGCPBillingEnv(process.SystemRunner().Runner)))
}
