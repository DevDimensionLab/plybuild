package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/gcpbilling"
)

var gcpBillingNow = time.Date(2026, time.November, 18, 21, 0, 0, 0, time.UTC)

const gcpBillingConfig = `{"billing": {"account": "019DC7-C01A85-D0D6A9",
  "table": "trip-1.billing_export.gcp_billing_export_v1_X", "location": "EU"},
 "projects": {"trip-1": ["trip"], "gemini-2": ["trip"]}}`

func gcpBillingSnapshot() gcpbilling.Snapshot {
	at := func(s string) time.Time { v, _ := time.Parse(time.RFC3339, s); return v }
	amount := func(v float64) *float64 { return &v }
	row := func(month, project, service, costType string, cost, credits float64, exported string) gcpbilling.Row {
		return gcpbilling.Row{Month: month, ProjectID: project, Service: service, CostType: costType, Cost: cost, Credits: credits, LastExport: at(exported)}
	}
	budget := func(name string, limit float64, projects ...string) gcpbilling.Budget {
		return gcpbilling.Budget{Name: name, Amount: amount(limit), Currency: "NOK", CalendarPeriod: "MONTH",
			CreditTypesTreatment: "INCLUDE_ALL_CREDITS", Projects: append([]string{}, projects...)}
	}
	unsupported := budget("Services", 5)
	unsupported.OtherFilters = true
	lastPeriod := budget("Last period", 1, "trip-1")
	lastPeriod.Amount = nil
	excluding := budget("Without credits", 1, "trip-1")
	excluding.CreditTypesTreatment = "EXCLUDE_ALL_CREDITS"
	return gcpbilling.Snapshot{Currency: "NOK", Rows: []gcpbilling.Row{
		row("202610", "trip-1", "Secret Manager", "regular", 153.02, 0, "2026-11-01T15:08:12Z"),
		row("202610", "trip-1", "Cloud Run", "regular", 2.86, -2.85, "2026-11-01T14:00:00Z"),
		row("202610", "trip-1", "Certificate Manager", "regular", 0, 0, "2026-11-01T14:00:00Z"),
		row("202610", "", "Invoice", "tax", 38.3, 0, "2026-11-01T21:42:47Z"),
		row("202610", "", "Invoice", "rounding_error", -0.0001, 0, "2026-11-01T17:56:29Z"),
		row("202610", "", "Support", "regular", 1, 0, "2026-11-01T10:00:00Z"),
		row("202610", "other-3", "Cloud Storage", "regular", 0.10, 0, "2026-11-01T10:00:00Z"),
		row("202611", "trip-1", "Secret Manager", "regular", 3.58, 0, "2026-11-02T23:54:10Z"),
		row("202611", "trip-1", "Cloud Run", "regular", 0.42, -0.42, "2026-11-02T20:00:00Z"),
		row("202611", "trip-1", "Invoice", "tax", 9, 0, "2026-11-03T00:00:00Z"),
	}, Budgets: []gcpbilling.Budget{
		budget("Trip", 55, "trip-1"), budget("Gemini", 20, "gemini-2"), budget("Guardrail", 75),
		budget("Trip strict", 60, "trip-1"), unsupported, lastPeriod, excluding,
	}}
}

// gcpBillingHome writes the configuration unless config is empty; the state
// directory is XDG_STATE_HOME under the same temporary directory.
func gcpBillingHome(t *testing.T, config string) (home, state string) {
	t.Helper()
	home = t.TempDir()
	if config != "" {
		dir := filepath.Join(home, ".config", "ply")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "gcp.json"), []byte(config), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return home, filepath.Join(home, "xdg-state")
}

func gcpBillingTestEnv(home, state string, read func(context.Context, gcpbilling.Source) (gcpbilling.Snapshot, error)) gcpBillingEnv {
	return gcpBillingEnv{
		read: read,
		now:  func() time.Time { return gcpBillingNow },
		home: func() (string, error) { return home, nil },
		getenv: func(name string) string {
			if name == "XDG_STATE_HOME" {
				return state
			}
			return ""
		},
	}
}

func gcpBillingReadFails(t *testing.T) func(context.Context, gcpbilling.Source) (gcpbilling.Snapshot, error) {
	return func(context.Context, gcpbilling.Source) (gcpbilling.Snapshot, error) {
		t.Fatal("a read command contacted Google")
		return gcpbilling.Snapshot{}, nil
	}
}

func executeGCP(t *testing.T, env gcpBillingEnv, args ...string) (string, error) {
	t.Helper()
	command := newGCPCommand(env)
	command.SilenceUsage = true // as ExecuteE does for RootCmd
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs(args)
	err := command.ExecuteContext(context.Background())
	return stdout.String(), err
}

// gcpBillingSynced runs sync with the snapshot and returns an env whose read fails.
func gcpBillingSynced(t *testing.T) gcpBillingEnv {
	t.Helper()
	home, state := gcpBillingHome(t, gcpBillingConfig)
	env := gcpBillingTestEnv(home, state, func(_ context.Context, src gcpbilling.Source) (gcpbilling.Snapshot, error) {
		if src != (gcpbilling.Source{Account: "019DC7-C01A85-D0D6A9", Table: "trip-1.billing_export.gcp_billing_export_v1_X", Location: "EU"}) {
			return gcpbilling.Snapshot{}, fmt.Errorf("source %+v", src)
		}
		return gcpBillingSnapshot(), nil
	})
	if _, err := executeGCP(t, env, "billing", "sync"); err != nil {
		t.Fatal(err)
	}
	env.read = gcpBillingReadFails(t)
	return env
}

func TestGCPBillingSyncSavesTheSnapshotUnderXDGStateHome(t *testing.T) {
	home, state := gcpBillingHome(t, gcpBillingConfig)
	env := gcpBillingTestEnv(home, state, func(context.Context, gcpbilling.Source) (gcpbilling.Snapshot, error) {
		return gcpBillingSnapshot(), nil
	})
	stdout, err := executeGCP(t, env, "billing", "sync", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	want := `{"kind":"PlyGcpBillingSync@1","schema_version":1,"source":"gcp","status":"ok","synced_at":"2026-11-18T21:00:00Z",` +
		`"months":["2026-10","2026-11"],"currency":"NOK","budgets":7}` + "\n"
	if stdout != want {
		t.Fatalf("stdout\n%s\nwant\n%s", stdout, want)
	}
	saved, err := gcpbilling.Load(filepath.Join(state, "ply", "gcp-billing.json"))
	if err != nil || !saved.SyncedAt.Equal(gcpBillingNow) || len(saved.Rows) != 10 {
		t.Fatalf("saved %+v, %v", saved, err)
	}
}

func TestGCPBillingSyncWithoutXDGStateHomeUsesLocalState(t *testing.T) {
	home, _ := gcpBillingHome(t, gcpBillingConfig)
	env := gcpBillingTestEnv(home, "", func(context.Context, gcpbilling.Source) (gcpbilling.Snapshot, error) {
		return gcpBillingSnapshot(), nil
	})
	if _, err := executeGCP(t, env, "billing", "sync"); err != nil {
		t.Fatal(err)
	}
	if _, err := gcpbilling.Load(filepath.Join(home, ".local", "state", "ply", "gcp-billing.json")); err != nil {
		t.Fatal(err)
	}
}

func TestGCPBillingSyncFailureKeepsTheOldSnapshot(t *testing.T) {
	env := gcpBillingSynced(t)
	env.now = func() time.Time { return gcpBillingNow.Add(time.Hour) }
	env.read = func(context.Context, gcpbilling.Source) (gcpbilling.Snapshot, error) {
		return gcpbilling.Snapshot{}, fmt.Errorf("%w: gcloud is not logged in", gcpbilling.ErrAuth)
	}
	stdout, err := executeGCP(t, env, "billing", "sync", "--format", "json")
	want := `{"kind":"PlyGcpBillingSync@1","schema_version":1,"source":"gcp","status":"error",` +
		`"error":{"code":"gcloud_auth","description":"gcloud is not logged in"}}` + "\n"
	if err == nil || stdout != want {
		t.Fatalf("stdout %s, error %v", stdout, err)
	}
	env.read = gcpBillingReadFails(t)
	quota, err := executeGCP(t, env, "billing", "quota", "--format", "json")
	if err != nil || !strings.Contains(quota, `"synced_at":"2026-11-18T21:00:00Z"`) {
		t.Fatalf("quota %s, error %v", quota, err)
	}
}

func TestGCPBillingQuotaGivesEachSupportedBudgetALevel(t *testing.T) {
	stdout, err := executeGCP(t, gcpBillingSynced(t), "billing", "quota", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	// November in Los Angeles: trip-1 used 3.58 net of credits; tax is left out.
	// Trip strict (60) loses to Trip (55); Services, Last period and Without credits are not supported.
	want := `{"kind":"PlyGcpBillingQuota@1","schema_version":1,"source":"gcp","status":"ok","synced_at":"2026-11-18T21:00:00Z","stale":false,"resources":[` +
		`{"id":"gcp.account","name":"Guardrail","project_ids":[],"observed_at":"2026-11-02T23:54:10Z","unit":"amount","currency":"NOK","window":"month",` +
		`"used":3.58,"limit":75,"remaining":71.42,"remaining_percent":95,"resets_at":"2026-12-01T08:00:00Z","level":"ok"},` +
		`{"id":"gcp.gemini-2","name":"Gemini","project_ids":["trip"],"unit":"amount","currency":"NOK","window":"month",` +
		`"used":0,"limit":20,"remaining":20,"remaining_percent":100,"resets_at":"2026-12-01T08:00:00Z","level":"ok"},` +
		`{"id":"gcp.trip-1","name":"Trip","project_ids":["trip"],"observed_at":"2026-11-02T23:54:10Z","unit":"amount","currency":"NOK","window":"month",` +
		`"used":3.58,"limit":55,"remaining":51.42,"remaining_percent":93,"resets_at":"2026-12-01T08:00:00Z","level":"ok"}]}` + "\n"
	if stdout != want {
		t.Fatalf("stdout\n%s\nwant\n%s", stdout, want)
	}
}

func TestGCPBillingMonthFollowsLosAngeles(t *testing.T) {
	env := gcpBillingSynced(t)
	// 1 December 07:59 UTC is still 30 November in Los Angeles.
	env.now = func() time.Time { return time.Date(2026, time.December, 1, 7, 59, 0, 0, time.UTC) }
	stdout, err := executeGCP(t, env, "billing", "usage", "--format", "json")
	if err != nil || !strings.Contains(stdout, `"month":"2026-11"`) {
		t.Fatalf("stdout %s, error %v", stdout, err)
	}
	env.now = func() time.Time { return time.Date(2026, time.December, 1, 8, 0, 0, 0, time.UTC) }
	stdout, err = executeGCP(t, env, "billing", "quota", "--format", "json")
	if err != nil || !strings.Contains(stdout, `"used":0,"limit":55`) || !strings.Contains(stdout, `"resets_at":"2027-01-01T08:00:00Z"`) {
		t.Fatalf("stdout %s, error %v", stdout, err)
	}
}

func TestGCPBillingStaleAfterTwentyFourHours(t *testing.T) {
	env := gcpBillingSynced(t)
	for offset, want := range map[time.Duration]string{24 * time.Hour: `"stale":false`, 24*time.Hour + time.Second: `"stale":true`} {
		env.now = func() time.Time { return gcpBillingNow.Add(offset) }
		for _, command := range []string{"quota", "usage"} {
			stdout, err := executeGCP(t, env, "billing", command, "--format", "json")
			if err != nil || !strings.Contains(stdout, want) {
				t.Fatalf("%s after %v: stdout %s, error %v", command, offset, stdout, err)
			}
		}
	}
}

func TestQuotaLevel(t *testing.T) {
	for _, test := range []struct {
		none    bool
		percent int
		want    string
	}{
		{true, 0, "exhausted"}, {false, 0, "low"}, {false, 19, "low"}, {false, 20, "ok"}, {false, 100, "ok"},
	} {
		if got := quotaLevel(test.none, test.percent); got != test.want {
			t.Fatalf("level(%v, %d%%) = %q, want %q", test.none, test.percent, got, test.want)
		}
	}
}

func TestGCPBillingQuotaExhaustedWhenUsedReachesTheBudget(t *testing.T) {
	home, state := gcpBillingHome(t, gcpBillingConfig)
	env := gcpBillingTestEnv(home, state, func(context.Context, gcpbilling.Source) (gcpbilling.Snapshot, error) {
		s := gcpBillingSnapshot()
		s.Rows = append(s.Rows, gcpbilling.Row{Month: "202611", ProjectID: "gemini-2", Service: "Gemini", CostType: "regular", Cost: 25, LastExport: gcpBillingNow})
		return s, nil
	})
	if _, err := executeGCP(t, env, "billing", "sync"); err != nil {
		t.Fatal(err)
	}
	stdout, err := executeGCP(t, env, "billing", "quota", "--format", "json")
	if err != nil || !strings.Contains(stdout, `"used":25,"limit":20,"remaining":0,"remaining_percent":0,"resets_at":"2026-12-01T08:00:00Z","level":"exhausted"`) {
		t.Fatalf("stdout %s, error %v", stdout, err)
	}
}

func TestGCPBillingUsageSplitsTheMonthByProjectAndService(t *testing.T) {
	stdout, err := executeGCP(t, gcpBillingSynced(t), "billing", "usage", "--month", "2026-10", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	want := `{"kind":"PlyGcpBillingUsage@1","schema_version":1,"source":"gcp","status":"ok","synced_at":"2026-11-18T21:00:00Z","stale":false,` +
		`"month":"2026-10","observed_at":"2026-11-01T21:42:47Z","unit":"amount","currency":"NOK","amount":154.13,"tax":38.3,"projects":[` +
		`{"key":"trip-1","project_ids":["trip"],"amount":153.03,"services":[{"name":"Secret Manager","amount":153.02},{"name":"Cloud Run","amount":0.01}]},` +
		`{"project_ids":[],"amount":1,"services":[{"name":"Support","amount":1}]},` +
		`{"key":"other-3","project_ids":[],"amount":0.1,"services":[{"name":"Cloud Storage","amount":0.1}]}]}` + "\n"
	if stdout != want {
		t.Fatalf("stdout\n%s\nwant\n%s", stdout, want)
	}
}

func TestGCPBillingUsageOfAMonthWithoutData(t *testing.T) {
	stdout, err := executeGCP(t, gcpBillingSynced(t), "billing", "usage", "--month", "2026-08", "--format", "json")
	if err != nil || !strings.HasSuffix(stdout, `"month":"2026-08","unit":"amount","currency":"NOK","amount":0,"tax":0,"projects":[]}`+"\n") {
		t.Fatalf("stdout %s, error %v", stdout, err)
	}
	if _, err := executeGCP(t, gcpBillingSynced(t), "billing", "usage", "--month", "2026-8"); err == nil {
		t.Fatal("--month 2026-8 accepted")
	}
}

func TestGCPBillingReadsWithoutConfigOrSnapshot(t *testing.T) {
	for name, test := range map[string]struct {
		config, state, status string
		exit                  bool
	}{
		"not configured": {"", "", `"status":"not_configured"}`, false},
		"pending":        {gcpBillingConfig, "", `"status":"pending"}`, false},
		"config invalid": {`{"billing": {"account": "x"}}`, "", `"status":"error","error":{"code":"config_invalid",`, true},
		"state invalid":  {gcpBillingConfig, "{}", `"status":"error","error":{"code":"state_invalid",`, true},
	} {
		t.Run(name, func(t *testing.T) {
			home, state := gcpBillingHome(t, test.config)
			if test.state != "" {
				path := filepath.Join(state, "ply", "gcp-billing.json")
				if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(test.state), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			for _, command := range []string{"quota", "usage"} {
				stdout, err := executeGCP(t, gcpBillingTestEnv(home, state, gcpBillingReadFails(t)), "billing", command, "--format", "json")
				if (err != nil) != test.exit || !strings.Contains(stdout, `"source":"gcp",`+test.status) {
					t.Fatalf("%s: stdout %s, error %v", command, stdout, err)
				}
			}
		})
	}
}

func TestGCPBillingSyncWithoutConfig(t *testing.T) {
	home, state := gcpBillingHome(t, "")
	stdout, err := executeGCP(t, gcpBillingTestEnv(home, state, gcpBillingReadFails(t)), "billing", "sync", "--format", "json")
	if err != nil || stdout != `{"kind":"PlyGcpBillingSync@1","schema_version":1,"source":"gcp","status":"not_configured"}`+"\n" {
		t.Fatalf("stdout %s, error %v", stdout, err)
	}
	if _, err := os.Stat(state); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("sync without config wrote state: %v", err)
	}
}

func TestGCPBillingText(t *testing.T) {
	env := gcpBillingSynced(t)
	quota, err := executeGCP(t, env, "billing", "quota")
	want := "Google Cloud gcp.account (Guardrail): 3.58 of 75.00 NOK used, 71.42 left (95%, ok), resets 2026-12-01\n" +
		"Google Cloud gcp.gemini-2 (Gemini): 0.00 of 20.00 NOK used, 20.00 left (100%, ok), resets 2026-12-01\n" +
		"Google Cloud gcp.trip-1 (Trip): 3.58 of 55.00 NOK used, 51.42 left (93%, ok), resets 2026-12-01\n"
	if err != nil || quota != want {
		t.Fatalf("quota %q, error %v", quota, err)
	}
	usage, err := executeGCP(t, env, "billing", "usage", "--month", "2026-10")
	want = "Google Cloud 2026-10: 154.13 NOK, tax 38.30 NOK\n" +
		"trip-1  153.03\n  Secret Manager  153.02\n  Cloud Run  0.01\n" +
		"(no project)  1.00\n  Support  1.00\n" +
		"other-3  0.10\n  Cloud Storage  0.10\n"
	if err != nil || usage != want {
		t.Fatalf("usage %q, error %v", usage, err)
	}
	env.now = func() time.Time { return gcpBillingNow.Add(25 * time.Hour) }
	if quota, _ := executeGCP(t, env, "billing", "quota"); !strings.HasPrefix(quota, "Google Cloud: last synced 2026-11-18T21:00:00Z, more than 24 hours ago\n") {
		t.Fatalf("stale quota %q", quota)
	}
	home, state := gcpBillingHome(t, gcpBillingConfig)
	if pending, _ := executeGCP(t, gcpBillingTestEnv(home, state, nil), "billing", "usage"); pending != "Google Cloud: not synced yet; run ply gcp billing sync\n" {
		t.Fatalf("pending %q", pending)
	}
}
