package gcpbilling

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

type reply struct {
	stdout string
	code   int
	err    error
}

// fakeCloud answers gcloud and bq in call order and records each call.
type fakeCloud struct {
	replies []reply
	calls   []string
	timeout []time.Duration
}

func (f *fakeCloud) Run(c process.Command) error {
	if c.Context == nil {
		return fmt.Errorf("unbounded command %s %q", c.Name, c.Args)
	}
	deadline, _ := c.Context.Deadline()
	f.calls = append(f.calls, c.Name+" "+strings.Join(c.Args, " "))
	f.timeout = append(f.timeout, time.Until(deadline).Round(time.Second))
	r := f.replies[0]
	f.replies = f.replies[1:]
	fmt.Fprint(c.Stdout, r.stdout)
	fmt.Fprint(c.Stderr, "STDERR-SECRET ya29.FAKETOKEN")
	if r.err != nil {
		return r.err
	}
	if r.code != 0 {
		return exec.Command("/bin/sh", "-c", fmt.Sprintf("exit %d", r.code)).Run()
	}
	return nil
}

var source = Source{Account: "019DC7-C01A85-D0D6A9", Table: "trip-1.billing_export.gcp_billing_export_v1_X", Location: "EU"}

const rowsBody = `[
 {"month":"202610","project_id":"trip-1","service":"Secret Manager","cost_type":"regular","currency":"NOK","cost":"153.02","credits":"0.0","last_export":"2026-11-01T15:08:12Z"},
 {"month":"202610","project_id":"trip-1","service":"Cloud Run","cost_type":"regular","currency":"NOK","cost":"2.86","credits":"-2.85","last_export":"2026-11-01T14:00:00Z"},
 {"month":"202610","service":"Invoice","cost_type":"tax","currency":"NOK","cost":"38.3","credits":"0.0","last_export":"2026-11-01T21:42:47Z"}]`

const budgetsBody = `[
 {"displayName":"Trip","amount":{"specifiedAmount":{"currencyCode":"NOK","units":"55","nanos":500000000}},
  "budgetFilter":{"calendarPeriod":"MONTH","creditTypesTreatment":"INCLUDE_ALL_CREDITS","projects":["projects/111"]}},
 {"displayName":"Guardrail","amount":{"specifiedAmount":{"currencyCode":"NOK","units":"75"}},
  "budgetFilter":{"calendarPeriod":"MONTH","creditTypesTreatment":"INCLUDE_ALL_CREDITS"}},
 {"displayName":"Services","amount":{"lastPeriodAmount":{}},
  "budgetFilter":{"creditTypesTreatment":"EXCLUDE_ALL_CREDITS","services":["services/x"]}},
 {"displayName":"Unknown project","amount":{"specifiedAmount":{"currencyCode":"NOK","units":"9"}},
  "budgetFilter":{"calendarPeriod":"MONTH","creditTypesTreatment":"INCLUDE_ALL_CREDITS","projects":["projects/999"]}}]`

const projectsBody = `[{"projectId":"trip-1","projectNumber":"111"},{"projectId":"gemini-2","projectNumber":"222"}]`

func ok(stdout string) reply { return reply{stdout: stdout} }

func TestReadSumsTheExportAndReadsBudgets(t *testing.T) {
	cloud := &fakeCloud{replies: []reply{ok("ya29.token\n"), ok(rowsBody), ok(budgetsBody), ok(projectsBody)}}
	got, err := Read(context.Background(), cloud, source)
	if err != nil {
		t.Fatal(err)
	}
	at := func(s string) time.Time { v, _ := time.Parse(time.RFC3339, s); return v }
	amount := func(v float64) *float64 { return &v }
	want := Snapshot{Currency: "NOK", Rows: []Row{
		{Month: "202610", ProjectID: "trip-1", Service: "Secret Manager", CostType: "regular", Cost: 153.02, LastExport: at("2026-11-01T15:08:12Z")},
		{Month: "202610", ProjectID: "trip-1", Service: "Cloud Run", CostType: "regular", Cost: 2.86, Credits: -2.85, LastExport: at("2026-11-01T14:00:00Z")},
		{Month: "202610", Service: "Invoice", CostType: "tax", Cost: 38.3, LastExport: at("2026-11-01T21:42:47Z")},
	}, Budgets: []Budget{
		{Name: "Trip", Amount: amount(55.5), Currency: "NOK", CalendarPeriod: "MONTH", CreditTypesTreatment: "INCLUDE_ALL_CREDITS", Projects: []string{"trip-1"}},
		{Name: "Guardrail", Amount: amount(75), Currency: "NOK", CalendarPeriod: "MONTH", CreditTypesTreatment: "INCLUDE_ALL_CREDITS", Projects: []string{}},
		{Name: "Services", CreditTypesTreatment: "EXCLUDE_ALL_CREDITS", Projects: []string{}, OtherFilters: true},
		{Name: "Unknown project", Amount: amount(9), Currency: "NOK", CalendarPeriod: "MONTH", CreditTypesTreatment: "INCLUDE_ALL_CREDITS", Projects: []string{}, OtherFilters: true},
	}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("snapshot\n%+v\nwant\n%+v", got, want)
	}
	for i, call := range []string{
		"gcloud auth print-access-token --quiet",
		"bq --headless --project_id=trip-1 query --location=EU --use_legacy_sql=false --format=json --max_rows=1000000 ",
		"gcloud billing budgets list --billing-account=019DC7-C01A85-D0D6A9 --billing-project=trip-1 --format=json --quiet",
		"gcloud projects list --format=json --quiet",
	} {
		if !strings.HasPrefix(cloud.calls[i], call) {
			t.Fatalf("call %d %q, want prefix %q", i, cloud.calls[i], call)
		}
	}
	if !strings.Contains(cloud.calls[1], "FROM `trip-1.billing_export.gcp_billing_export_v1_X`") {
		t.Fatalf("query does not read the table: %s", cloud.calls[1])
	}
	for i, d := range cloud.timeout {
		if d != CallTimeout {
			t.Fatalf("call %d timeout %v, want %v", i, d, CallTimeout)
		}
	}
}

func TestReadClassifiesFailuresWithoutCloudOutput(t *testing.T) {
	missing := reply{err: &exec.Error{Name: "gcloud", Err: exec.ErrNotFound}}
	for name, test := range map[string]struct {
		replies []reply
		want    error
	}{
		"gcloud missing":  {[]reply{missing}, ErrMissing},
		"bq missing":      {[]reply{ok("t"), {err: &exec.Error{Name: "bq", Err: exec.ErrNotFound}}}, ErrMissing},
		"not logged in":   {[]reply{{code: 1}}, ErrAuth},
		"timeout":         {[]reply{ok("t"), {err: context.DeadlineExceeded}}, ErrTimeout},
		"query failed":    {[]reply{ok("t"), {code: 2}}, ErrFailed},
		"budgets failed":  {[]reply{ok("t"), ok(rowsBody), {code: 1}}, ErrFailed},
		"projects failed": {[]reply{ok("t"), ok(rowsBody), ok(budgetsBody), {code: 1}}, ErrFailed},
		"unreadable rows": {[]reply{ok("t"), ok("Waiting on job")}, ErrFailed},
		"bad cost":        {[]reply{ok("t"), ok(`[{"month":"202610","cost":"x","credits":"0","currency":"NOK","last_export":"2026-11-01T00:00:00Z"}]`)}, ErrFailed},
		"mixed currency": {[]reply{ok("t"), ok(`[{"month":"202610","cost":"1","credits":"0","currency":"NOK","last_export":"2026-11-01T00:00:00Z"},` +
			`{"month":"202611","cost":"1","credits":"0","currency":"EUR","last_export":"2026-11-02T00:00:00Z"}]`)}, ErrMixedCurrency},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := Read(context.Background(), &fakeCloud{replies: test.replies}, source)
			if !errors.Is(err, test.want) {
				t.Fatalf("error %v, want %v", err, test.want)
			}
			if strings.Contains(err.Error(), "SECRET") || strings.Contains(err.Error(), "ya29") {
				t.Fatalf("error exposes cloud output: %v", err)
			}
		})
	}
}

func TestReadCanceledByTheCallerIsNotATimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Read(ctx, &fakeCloud{replies: []reply{{err: context.Canceled}}}, source)
	if !errors.Is(err, ErrFailed) {
		t.Fatalf("error %v, want %v", err, ErrFailed)
	}
}

func TestSourceValid(t *testing.T) {
	if err := source.Valid(); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []Source{
		{Account: "x", Table: source.Table, Location: "EU"},
		{Account: source.Account, Table: "trip-1.billing_export", Location: "EU"},
		{Account: source.Account, Table: "trip-1.d.t` WHERE 1=1 --", Location: "EU"},
		{Account: source.Account, Table: source.Table, Location: ""},
	} {
		if bad.Valid() == nil {
			t.Fatalf("%+v is valid", bad)
		}
	}
}

func TestSaveAndLoadKeepTheSnapshotPrivate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "ply", "gcp-billing.json")
	if _, err := Load(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("load before save: %v", err)
	}
	amount := 55.0
	snapshot := Snapshot{SyncedAt: time.Date(2026, 11, 18, 21, 0, 0, 0, time.UTC), Currency: "NOK",
		Rows:    []Row{{Month: "202611", ProjectID: "trip-1", Service: "Cloud Run", CostType: "regular", Cost: 0.42, Credits: -0.42, LastExport: time.Date(2026, 11, 2, 23, 54, 10, 0, time.UTC)}},
		Budgets: []Budget{{Name: "Trip", Amount: &amount, Currency: "NOK", CalendarPeriod: "MONTH", CreditTypesTreatment: "INCLUDE_ALL_CREDITS", Projects: []string{"trip-1"}}}}
	if err := Save(path, snapshot); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil || !reflect.DeepEqual(got, snapshot) {
		t.Fatalf("load %+v, %v", got, err)
	}
	for p, mode := range map[string]os.FileMode{path: 0o600, filepath.Dir(path): 0o700} {
		if info, err := os.Stat(p); err != nil || info.Mode().Perm() != mode {
			t.Fatalf("%s mode %v, %v", p, info.Mode().Perm(), err)
		}
	}
	for _, body := range []string{`{}`, `not json`, `{"kind":"PlyGcpBillingState@1","schema_version":2}`} {
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(path); !errors.Is(err, ErrState) {
			t.Fatalf("load %s: %v", body, err)
		}
	}
}

func TestLockCreatesThePrivateDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ply", "gcp-billing.json")
	lock, err := Lock(path)
	if err != nil {
		t.Fatal(err)
	}
	defer lock.Release()
	if lock.Path() != path+".lock" {
		t.Fatalf("lock path %s", lock.Path())
	}
	if info, err := os.Stat(filepath.Dir(path)); err != nil || info.Mode().Perm() != 0o700 {
		t.Fatalf("directory %v, %v", info, err)
	}
}
