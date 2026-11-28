package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/quotahistory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// quotaHistoryPeriod is the period ply quota history reads without --since.
const quotaHistoryPeriod = 8 * 24 * time.Hour

// quotaEnv is the clock of ply quota history; the directory comes from the environment.
type quotaEnv struct {
	now func() time.Time
}

// recordQuotaHistory appends p to series in this machine's quota history.
func recordQuotaHistory(series string, p quotahistory.Point) error {
	dir, err := quotahistory.DefaultDir(os.Getenv, os.UserHomeDir)
	if err == nil {
		_, err = quotahistory.Open(dir).Record(series, p)
	}
	return err
}

// logHistoryError keeps a failed history write out of the command's output and exit status.
func logHistoryError(err error) {
	if err != nil {
		logrus.Debugf("quota history not written: %v", err)
	}
}

// recordActionsHistory adds the used minutes of a read GitHub Actions quota to the history.
func recordActionsHistory(record func(string, quotahistory.Point) error, q githubActionsQuota) {
	if record == nil || q.Used == nil {
		return
	}
	at, err := time.Parse(time.RFC3339, q.ObservedAt)
	if err != nil {
		logHistoryError(err)
		return
	}
	p := quotahistory.Point{At: at, Used: float64(*q.Used)}
	if q.Limit != nil {
		limit := float64(*q.Limit)
		p.Limit = &limit
	}
	if reset, err := time.Parse(time.RFC3339, q.ResetsAt); err == nil {
		p.ResetsAt = &reset
	}
	logHistoryError(record(q.ID, p))
}

func quotaInvalid(detail string) error { return fmt.Errorf("quota_invalid_arguments: %s", detail) }

func newQuotaCommand(env quotaEnv) *cobra.Command {
	command := &cobra.Command{
		Use:               "quota",
		Short:             "Read the quota history of this machine",
		Args:              cobra.NoArgs,
		RunE:              func(c *cobra.Command, args []string) error { return c.Help() },
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	}
	var series []string
	var since, until, format string
	history := &cobra.Command{
		Use:   "history",
		Short: "Show the quota points ply agent and ply github actions quota have recorded on this machine",
		Long: "Read the quota history in $XDG_STATE_HOME/" + quotahistory.StateFile + " (by default below ~/.local/state). ply agent and\n" +
			"ply github actions quota add a point when a value changed or the last point is ten minutes older; this command only reads,\n" +
			"without network or workspace. Series: claude.five_hour, claude.weekly, codex.five_hour and codex.weekly (used percent) and\n" +
			"github.actions (used minutes this month, with the included minutes as limit). Without --since the last 8 days are read;\n" +
			"both bounds are inclusive RFC 3339 times. Lines that cannot be read are listed under invalid.",
		Example: "  ply quota history --format json\n  ply quota history --series claude.weekly --since 2026-11-01T00:00:00Z --format json",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return quotaInvalid("expected no positional arguments")
			}
			if format != "text" && format != "json" {
				return quotaInvalid("--format must be text or json")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			from, to, err := quotaPeriod(env.now(), since, until)
			if err != nil {
				return err
			}
			dir, err := quotahistory.DefaultDir(os.Getenv, os.UserHomeDir)
			if err != nil {
				return err
			}
			h, err := quotahistory.Open(dir).Read(series, from, to)
			if err != nil {
				if errors.Is(err, quotahistory.ErrUnknownSeries) {
					return quotaInvalid(err.Error())
				}
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string `json:"kind"`
					SchemaVersion int    `json:"schema_version"`
					Since         string `json:"since"`
					Until         string `json:"until"`
					quotahistory.History
				}{quotahistory.Kind, 1, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339), h})
			}
			return writeQuotaHistoryText(cmd.OutOrStdout(), h)
		},
	}
	history.Flags().StringArrayVar(&series, "series", nil, "series ID; repeat for more (default: every series)")
	history.Flags().StringVar(&since, "since", "", "first time, RFC 3339 (default: 8 days ago)")
	history.Flags().StringVar(&until, "until", "", "last time, RFC 3339 (default: now)")
	history.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	command.AddCommand(history)
	return command
}

// quotaPeriod is the period to read: --since and --until, by default the last eight days.
func quotaPeriod(now time.Time, since, until string) (time.Time, time.Time, error) {
	to := now.UTC().Truncate(time.Second)
	from := to.Add(-quotaHistoryPeriod)
	var err error
	if since != "" {
		if from, err = time.Parse(time.RFC3339, since); err != nil {
			return from, to, quotaInvalid("--since must be an RFC 3339 time")
		}
	}
	if until != "" {
		if to, err = time.Parse(time.RFC3339, until); err != nil {
			return from, to, quotaInvalid("--until must be an RFC 3339 time")
		}
	}
	return from, to, nil
}

func writeQuotaHistoryText(w io.Writer, h quotahistory.History) error {
	var out strings.Builder
	for _, s := range h.Series {
		fmt.Fprintf(&out, "%s (%s): %d points", s.ID, s.Unit, len(s.Points))
		if n := len(s.Points); n > 0 {
			last := s.Points[n-1]
			fmt.Fprintf(&out, ", last %g at %s", last.Used, last.At.UTC().Format(time.RFC3339))
		}
		out.WriteString("\n")
	}
	for _, invalid := range h.Invalid {
		fmt.Fprintf(&out, "warning: %s line %d: %s\n", invalid.File, invalid.Line, invalid.Error)
	}
	_, err := io.WriteString(w, out.String())
	return err
}

func init() {
	RootCmd.AddCommand(newQuotaCommand(quotaEnv{now: time.Now}))
}
