// Package gcpbilling reads a Cloud Billing export in BigQuery and the billing
// account's budgets with the user's own gcloud and bq login, and keeps the
// result as one private local snapshot. It never changes the login, a project
// or a permission.
package gcpbilling

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/atomicfile"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// CallTimeout bounds each gcloud and bq call.
const CallTimeout = 60 * time.Second

var (
	ErrMissing       = errors.New("gcloud_missing")
	ErrAuth          = errors.New("gcloud_auth")
	ErrTimeout       = errors.New("timeout")
	ErrFailed        = errors.New("gcloud_failed")
	ErrMixedCurrency = errors.New("mixed_currency")
	ErrState         = errors.New("state_invalid")
)

// Source is the billing account and its export table, project.dataset.table.
type Source struct{ Account, Table, Location string }

var (
	account  = regexp.MustCompile(`^[0-9A-F]{6}-[0-9A-F]{6}-[0-9A-F]{6}$`)
	table    = regexp.MustCompile(`^[a-z][a-z0-9-]{4,28}[a-z0-9]\.[A-Za-z0-9_]+\.[A-Za-z0-9_]+$`)
	location = regexp.MustCompile(`^[A-Za-z0-9-]+$`)
)

// Valid reports whether every part has its Google form, which also keeps the
// table name safe to put in the query.
func (s Source) Valid() error {
	if !account.MatchString(s.Account) || !table.MatchString(s.Table) || !location.MatchString(s.Location) {
		return errors.New("billing account, table or location is not in Google's form")
	}
	return nil
}

func (s Source) project() string { return s.Table[:strings.IndexByte(s.Table, '.')] }

// Row is the export summed per invoice month (YYYYMM), project, service, cost
// type and currency. Credits are negative. ProjectID is empty for rows
// without a project, such as tax.
type Row struct {
	Month      string    `json:"month"`
	ProjectID  string    `json:"project_id,omitempty"`
	Service    string    `json:"service"`
	CostType   string    `json:"cost_type"`
	Cost       float64   `json:"cost"`
	Credits    float64   `json:"credits"`
	LastExport time.Time `json:"last_export"`
}

// Budget is one budget of the billing account. Amount is nil unless it is a
// specified amount. Projects holds project IDs. OtherFilters is true when the
// budget filters on more than projects, credit treatment and calendar period,
// or names a project that gcloud cannot list.
type Budget struct {
	Name                 string   `json:"name"`
	Amount               *float64 `json:"amount,omitempty"`
	Currency             string   `json:"currency,omitempty"`
	CalendarPeriod       string   `json:"calendar_period,omitempty"`
	CreditTypesTreatment string   `json:"credit_types_treatment,omitempty"`
	Projects             []string `json:"projects"`
	OtherFilters         bool     `json:"other_filters,omitempty"`
}

// Snapshot is one complete read. Currency is empty without rows.
type Snapshot struct {
	SyncedAt time.Time `json:"synced_at"`
	Currency string    `json:"currency,omitempty"`
	Rows     []Row     `json:"rows"`
	Budgets  []Budget  `json:"budgets"`
}

const query = "SELECT invoice.month AS month, project.id AS project_id, service.description AS service, cost_type, currency, " +
	"SUM(cost) AS cost, SUM((SELECT IFNULL(SUM(c.amount), 0) FROM UNNEST(credits) c)) AS credits, " +
	"FORMAT_TIMESTAMP('%%Y-%%m-%%dT%%H:%%M:%%SZ', MAX(export_time)) AS last_export " +
	"FROM `%s` GROUP BY 1, 2, 3, 4, 5"

// Read checks the login, sums the whole export and reads the budgets. Errors
// carry one of the sentinels above and never gcloud's or bq's output.
func Read(ctx context.Context, runner process.Runner, src Source) (Snapshot, error) {
	if err := call(ctx, runner, "gcloud", []string{"auth", "print-access-token", "--quiet"}, nil); err != nil {
		if errors.Is(err, ErrFailed) && ctx.Err() == nil {
			return Snapshot{}, fmt.Errorf("%w: gcloud is not logged in", ErrAuth)
		}
		return Snapshot{}, err
	}
	var rows []struct {
		Month      string `json:"month"`
		ProjectID  string `json:"project_id"`
		Service    string `json:"service"`
		CostType   string `json:"cost_type"`
		Currency   string `json:"currency"`
		Cost       string `json:"cost"`
		Credits    string `json:"credits"`
		LastExport string `json:"last_export"`
	}
	args := []string{"--headless", "--project_id=" + src.project(), "query", "--location=" + src.Location,
		"--use_legacy_sql=false", "--format=json", "--max_rows=1000000", fmt.Sprintf(query, src.Table)}
	if err := call(ctx, runner, "bq", args, &rows); err != nil {
		return Snapshot{}, err
	}
	snapshot := Snapshot{Rows: []Row{}, Budgets: []Budget{}}
	for _, r := range rows {
		cost, costErr := strconv.ParseFloat(r.Cost, 64)
		credits, creditsErr := strconv.ParseFloat(r.Credits, 64)
		exported, exportErr := time.Parse(time.RFC3339, r.LastExport)
		if costErr != nil || creditsErr != nil || exportErr != nil || math.IsNaN(cost+credits) || math.IsInf(cost+credits, 0) {
			return Snapshot{}, fmt.Errorf("%w: unreadable export row", ErrFailed)
		}
		if snapshot.Currency != "" && r.Currency != snapshot.Currency {
			return Snapshot{}, fmt.Errorf("%w: the export has more than one currency", ErrMixedCurrency)
		}
		snapshot.Currency = r.Currency
		snapshot.Rows = append(snapshot.Rows, Row{Month: r.Month, ProjectID: r.ProjectID, Service: r.Service, CostType: r.CostType,
			Cost: cost, Credits: credits, LastExport: exported.UTC()})
	}
	var budgets []struct {
		DisplayName string `json:"displayName"`
		Amount      struct {
			Specified *struct {
				CurrencyCode string `json:"currencyCode"`
				Units        string `json:"units"`
				Nanos        int64  `json:"nanos"`
			} `json:"specifiedAmount"`
		} `json:"amount"`
		Filter map[string]json.RawMessage `json:"budgetFilter"`
	}
	if err := call(ctx, runner, "gcloud", []string{"billing", "budgets", "list", "--billing-account=" + src.Account,
		"--billing-project=" + src.project(), "--format=json", "--quiet"}, &budgets); err != nil {
		return Snapshot{}, err
	}
	var projects []struct {
		ID     string `json:"projectId"`
		Number string `json:"projectNumber"`
	}
	if err := call(ctx, runner, "gcloud", []string{"projects", "list", "--format=json", "--quiet"}, &projects); err != nil {
		return Snapshot{}, err
	}
	ids := map[string]string{}
	for _, p := range projects {
		ids["projects/"+p.Number] = p.ID
	}
	for _, b := range budgets {
		budget := Budget{Name: b.DisplayName, Projects: []string{}}
		if s := b.Amount.Specified; s != nil {
			units, err := strconv.ParseInt(s.Units, 10, 64)
			if err != nil && s.Units != "" {
				return Snapshot{}, fmt.Errorf("%w: unreadable budget amount", ErrFailed)
			}
			amount := float64(units) + float64(s.Nanos)/1e9
			budget.Amount, budget.Currency = &amount, s.CurrencyCode
		}
		for key, value := range b.Filter {
			var err error
			switch key {
			case "calendarPeriod":
				err = json.Unmarshal(value, &budget.CalendarPeriod)
			case "creditTypesTreatment":
				err = json.Unmarshal(value, &budget.CreditTypesTreatment)
			case "projects":
				var names []string
				err = json.Unmarshal(value, &names)
				for _, name := range names {
					if id, ok := ids[name]; ok {
						budget.Projects = append(budget.Projects, id)
					} else {
						budget.OtherFilters = true
					}
				}
			default:
				budget.OtherFilters = true
			}
			if err != nil {
				return Snapshot{}, fmt.Errorf("%w: unreadable budget filter", ErrFailed)
			}
		}
		if budget.OtherFilters {
			budget.Projects = []string{}
		}
		snapshot.Budgets = append(snapshot.Budgets, budget)
	}
	return snapshot, nil
}

// call runs one command and decodes its stdout into out unless out is nil.
func call(ctx context.Context, runner process.Runner, name string, args []string, out interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, CallTimeout)
	defer cancel()
	var stdout bytes.Buffer
	err := runner.Run(process.Command{Name: name, Args: args, Context: ctx, Stdout: &stdout, Stderr: io.Discard})
	switch {
	case errors.Is(err, exec.ErrNotFound):
		return fmt.Errorf("%w: %s is not installed", ErrMissing, name)
	case errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded):
		return fmt.Errorf("%w: %s did not answer within %s", ErrTimeout, name, CallTimeout)
	case err != nil:
		return fmt.Errorf("%w: %s failed", ErrFailed, name)
	}
	if out != nil && json.Unmarshal(stdout.Bytes(), out) != nil {
		return fmt.Errorf("%w: unreadable %s output", ErrFailed, name)
	}
	return nil
}

const stateKind = "PlyGcpBillingState@1"

type state struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
	Snapshot
}

// Lock takes the lock beside path that lets one sync at a time read and save,
// creating the private directory. It waits for another holder.
func Lock(path string) (*filelock.Lock, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, err
	}
	return filelock.Acquire(path+".lock", filelock.Options{})
}

// Save replaces the snapshot at path as a whole, readable only by the user.
func Save(path string, s Snapshot) error {
	data, err := json.Marshal(state{Kind: stateKind, SchemaVersion: 1, Snapshot: s})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return atomicfile.Replace(path, data, 0o600)
}

// Load reads the snapshot at path. A missing file is os.ErrNotExist; a file
// that is not this version's form is ErrState.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var s state
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&s) != nil || s.Kind != stateKind || s.SchemaVersion != 1 || s.Rows == nil || s.Budgets == nil {
		return Snapshot{}, fmt.Errorf("%w: %s is not a Ply GCP billing snapshot; run ply gcp billing sync", ErrState, path)
	}
	return s.Snapshot, nil
}
