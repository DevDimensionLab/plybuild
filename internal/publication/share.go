package publication

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// ShareKind names share.json: the Claude artifact of a path, one fixed link
// that each share of a newer run updates.
const ShareKind = "PlyPublicationShare@1"

// Share states. Pending and unknown hold the reservation of the path until a
// result is recorded; shared and failed end the attempt.
const (
	ShareStatePending = "pending"
	ShareStateShared  = "shared"
	ShareStateFailed  = "failed"
	ShareStateUnknown = "unknown"
)

// Share error classes.
const (
	ClassSharePending      = "publication_share_pending"
	ClassShareConflict     = "publication_share_conflict"
	ClassShareLaunchFailed = "publication_share_launch_failed"
)

const (
	shareName = "share.json"
	maxError  = 200
)

var artifactPattern = regexp.MustCompile(`^https://claude\.ai/artifact/[A-Za-z0-9_-]{1,64}$`)

// Path is the identity of a publication path.
type Path struct {
	Repo, Module, Category, Agent string
}

// ShareSession is the Herdr tab and agent session of the latest attempt, for
// talking to or resuming the agent; Ply does neither.
type ShareSession struct {
	Provider string `json:"provider"`
	Pane     string `json:"pane,omitempty"`
	Tab      string `json:"tab,omitempty"`
	ID       string `json:"id,omitempty"`
}

// ShareRecord is the share of a path as readbacks show it.
type ShareRecord struct {
	Project    string        `json:"project"`
	Repo       string        `json:"repo,omitempty"`
	Module     string        `json:"module,omitempty"`
	Category   string        `json:"category"`
	Agent      string        `json:"agent,omitempty"`
	State      string        `json:"state"`
	Run        string        `json:"run"`
	URL        *string       `json:"url"`
	SharedRun  *string       `json:"shared_run"`
	Error      *string       `json:"error"`
	StartedAt  string        `json:"started_at"`
	FinishedAt *string       `json:"finished_at"`
	Actor      string        `json:"actor,omitempty"`
	Session    *ShareSession `json:"session"`
}

// Share is the content of share.json.
type Share struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
	ShareRecord
}

func (s Share) open() bool { return s.State == ShareStatePending || s.State == ShareStateUnknown }

func shareConflict(detail string) error { return &Error{Class: ClassShareConflict, Detail: detail} }

// ReserveShare starts an attempt to share run (the newest visible run when
// empty) of path p, before any agent is started. A path whose shared run is
// run is returned unchanged with started false; an open attempt is refused.
func (s Store) ReserveShare(project string, p Path, run, actor string) (Share, bool, error) {
	if err := checkLine("--actor", actor, maxTitle, false); err != nil {
		return Share{}, false, err
	}
	v, err := s.Show(p.Repo, p.Module, p.Category, p.Agent, run)
	if err != nil {
		return Share{}, false, err
	}
	var result Share
	started := false
	err = s.locked(func() error {
		directory := filepath.Dir(v.Path)
		previous, found, err := readShare(directory)
		if err != nil {
			return err
		}
		switch {
		case found && previous.open():
			return &Error{Class: ClassSharePending, Detail: "run " + previous.Run + " is being shared (" + previous.State + ")"}
		case found && previous.State == ShareStateShared && previous.SharedRun != nil && *previous.SharedRun == v.Run:
			result = previous
			return nil
		}
		result = Share{Kind: ShareKind, SchemaVersion: SchemaVersion, ShareRecord: ShareRecord{Project: project, Repo: p.Repo, Module: storedModule(p.Module),
			Category: p.Category, Agent: p.Agent, State: ShareStatePending, Run: v.Run, StartedAt: s.now().UTC().Format(time.RFC3339), Actor: actor}}
		if found {
			result.URL, result.SharedRun = previous.URL, previous.SharedRun
		}
		started = true
		return writeShare(directory, result)
	})
	return result, started && err == nil, err
}

// FinishShareLaunch records how starting the agent for run ended: pending
// when it started, failed when it did not, unknown when the effect is unknown.
func (s Store) FinishShareLaunch(p Path, run string, session *ShareSession, state, detail string) (Share, error) {
	if state != ShareStatePending && state != ShareStateFailed && state != ShareStateUnknown {
		return Share{}, invalid("launch state must be pending, failed or unknown")
	}
	return s.updateShare(p, run, func(share *Share) {
		share.State, share.Session = state, session
		if state != ShareStatePending {
			share.Error = oneLine(detail)
		}
		if state == ShareStateFailed {
			share.FinishedAt = s.stamp()
		}
	}, func(share Share) bool { return share.State == ShareStatePending && share.Session == nil })
}

// RecordShareResult ends the attempt for run with the artifact url or an
// error. The url must be the path's fixed link once it has one. Recording the
// same result again returns the share unchanged.
func (s Store) RecordShareResult(p Path, run, url, detail string) (Share, error) {
	if (url == "") == (detail == "") {
		return Share{}, invalid("exactly one of --url and --error is required")
	}
	if url != "" && !artifactPattern.MatchString(url) {
		return Share{}, invalid("--url must be https://claude.ai/artifact/<id>")
	}
	if err := checkLine("--error", detail, maxError, false); err != nil {
		return Share{}, err
	}
	var replay *Share
	share, err := s.updateShare(p, run, func(share *Share) {
		share.FinishedAt = s.stamp()
		if url != "" {
			shared := share.Run
			share.State, share.URL, share.SharedRun, share.Error = ShareStateShared, &url, &shared, nil
			return
		}
		share.State, share.Error = ShareStateFailed, &detail
	}, func(share Share) bool {
		if share.open() {
			return url == "" || share.URL == nil || *share.URL == url
		}
		if sameResult(share, url, detail) {
			replay = &share
		}
		return false
	})
	if replay != nil {
		return *replay, nil
	}
	return share, err
}

func sameResult(share Share, url, detail string) bool {
	if url != "" {
		return share.State == ShareStateShared && share.URL != nil && *share.URL == url && share.SharedRun != nil && *share.SharedRun == share.Run
	}
	return share.State == ShareStateFailed && share.Error != nil && *share.Error == detail
}

// updateShare changes the share of p for run under the lock when allowed says so.
func (s Store) updateShare(p Path, run string, change func(*Share), allowed func(Share) bool) (Share, error) {
	relative, err := location(p.Repo, p.Module, p.Category, p.Agent)
	if err != nil {
		return Share{}, err
	}
	if !runPattern.MatchString(run) {
		return Share{}, invalid("--run is not a run ID")
	}
	directory := filepath.Join(s.root, relative)
	var result Share
	err = s.locked(func() error {
		share, found, err := readShare(directory)
		if err != nil {
			return err
		}
		if !found || share.Run != run {
			return shareConflict("no attempt to share run " + run)
		}
		if !allowed(share) {
			return shareConflict("the attempt to share run " + run + " is " + share.State)
		}
		change(&share)
		result = share
		return writeShare(directory, share)
	})
	return result, err
}

func (s Store) stamp() *string {
	at := s.now().UTC().Format(time.RFC3339)
	return &at
}

func (s Store) locked(f func() error) error {
	if err := os.MkdirAll(s.root, 0o755); err != nil {
		return ioError("create publication directory", err)
	}
	lock, err := filelock.Acquire(filepath.Join(s.root, lockName), filelock.Options{})
	if err != nil {
		return ioError("lock publications", err)
	}
	err = f()
	if releaseErr := lock.Release(); releaseErr != nil && err == nil {
		err = &Error{Class: ClassIOError, Detail: "the share was stored; releasing the lock failed"}
	}
	return err
}

// oneLine keeps a launch error to one line of at most maxError characters.
func oneLine(detail string) *string {
	runes := []rune(detail)
	for i, r := range runes {
		if r < ' ' || r == 0x7f {
			runes[i] = ' '
		}
	}
	if len(runes) > maxError {
		runes = runes[:maxError]
	}
	line := string(runes)
	return &line
}

func writeShare(directory string, share Share) error {
	temp := filepath.Join(directory, "."+shareName+tempSuffix)
	if err := writeJSON(temp, shareName, share); err != nil {
		_ = os.Remove(temp)
		return err
	}
	if err := os.Rename(temp, filepath.Join(directory, shareName)); err != nil {
		_ = os.Remove(temp)
		return ioError("write "+shareName, err)
	}
	return nil
}

// readShare returns the share of directory and whether one exists.
func readShare(directory string) (Share, bool, error) {
	contents, err := os.ReadFile(filepath.Join(directory, shareName))
	if errors.Is(err, fs.ErrNotExist) {
		return Share{}, false, nil
	}
	if err != nil {
		return Share{}, false, ioError("read "+shareName, err)
	}
	var share Share
	if err := json.Unmarshal(contents, &share); err != nil || share.Kind != ShareKind || !runPattern.MatchString(share.Run) {
		return Share{}, false, &Error{Class: ClassIOError, Detail: shareName + " is not a valid " + ShareKind}
	}
	return share, true, nil
}

// shareOf is the share readback of directory: nil without a readable share.
func shareOf(directory string) (*ShareRecord, error) {
	share, found, err := readShare(directory)
	if err != nil || !found {
		return nil, err
	}
	return &share.ShareRecord, nil
}
