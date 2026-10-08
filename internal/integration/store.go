package integration

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"golang.org/x/sys/unix"
)

var integrationID = regexp.MustCompile(`^int_[a-f0-9]{64}$`)
var sha256Pattern = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)

func byteHash(b []byte) string { return fmt.Sprintf("sha256:%x", sha256.Sum256(b)) }
func digest(v any) string {
	b, e := taskrun.Canonical(v)
	if e != nil {
		panic(e)
	}
	return byteHash(b)
}
func operationID(hash string) string { return "int_" + strings.TrimPrefix(hash, "sha256:") }
func storeRoot(root string) string   { return filepath.Join(root, ".ply", "integrations", "v1") }
func receiptPath(root, id string) string {
	return filepath.Join(storeRoot(root), "records", id+".json")
}

// All reads, including previews, avoid mkdir, locks and access-time based state.
func readBounded(path string, limit int64) ([]byte, error) {
	abs, e := filepath.Abs(path)
	if e != nil {
		return nil, e
	}
	physical, e := filepath.EvalSymlinks(abs)
	if e != nil {
		return nil, e
	}
	if abs != physical {
		return nil, fmt.Errorf("integration path contains a symlink: %s", path)
	}
	before, e := os.Lstat(abs)
	if e != nil {
		return nil, e
	}
	if !before.Mode().IsRegular() {
		return nil, fmt.Errorf("integration evidence must be a regular file: %s", path)
	}
	f, e := os.Open(abs)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	opened, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if !os.SameFile(before, opened) {
		return nil, fmt.Errorf("integration evidence changed before read")
	}
	b, e := io.ReadAll(io.LimitReader(f, limit+1))
	if e != nil {
		return nil, e
	}
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("integration evidence exceeds limit")
	}
	after, e := f.Stat()
	if e != nil {
		return nil, e
	}
	current, e := os.Lstat(abs)
	if e != nil {
		return nil, e
	}
	if !os.SameFile(opened, current) || opened.ModTime() != after.ModTime() || opened.Size() != after.Size() || int64(len(b)) != opened.Size() {
		return nil, fmt.Errorf("integration evidence changed during read")
	}
	return b, nil
}
func readJSON(path string, out any) error {
	b, e := readBounded(path, 32<<20)
	if e != nil {
		return e
	}
	if _, e = canonicaljson.DecodeStrict(b); e != nil {
		return e
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if e = d.Decode(out); e != nil {
		return e
	}
	if e = d.Decode(new(any)); e != io.EOF {
		return fmt.Errorf("trailing JSON")
	}
	return nil
}
func ensureDirectories(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("nonphysical integration store")
	}
	current := string(filepath.Separator)
	for _, part := range strings.Split(strings.TrimPrefix(path, current), string(filepath.Separator)) {
		if part == "" {
			continue
		}
		current = filepath.Join(current, part)
		i, e := os.Lstat(current)
		if os.IsNotExist(e) {
			if e = os.Mkdir(current, 0700); e != nil && !os.IsExist(e) {
				return e
			}
			i, e = os.Lstat(current)
		}
		if e != nil {
			return e
		}
		if !i.IsDir() || i.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("integration storage contains symlink or nondirectory: %s", current)
		}
	}
	return nil
}
func atomicJSON(path string, v any, immutable bool) error {
	b, e := taskrun.Canonical(v)
	if e != nil {
		return e
	}
	if e = ensureDirectories(filepath.Dir(path)); e != nil {
		return e
	}
	if old, err := readBounded(path, 32<<20); err == nil {
		if bytes.Equal(old, b) {
			return nil
		}
		if immutable {
			return fmt.Errorf("immutable integration evidence conflict: %s", path)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	f, e := os.CreateTemp(filepath.Dir(path), ".publish-")
	if e != nil {
		return e
	}
	defer os.Remove(f.Name())
	if e = f.Chmod(0600); e == nil {
		_, e = f.Write(b)
	}
	if e == nil {
		e = f.Sync()
	}
	ce := f.Close()
	if e == nil {
		e = ce
	}
	if e != nil {
		return e
	}
	if immutable {
		if e = os.Chmod(f.Name(), 0400); e != nil {
			return e
		}
	}
	if e = os.Rename(f.Name(), path); e != nil {
		return e
	}
	d, e := os.Open(filepath.Dir(path))
	if e != nil {
		return e
	}
	defer d.Close()
	return d.Sync()
}
func withLock(root string, fn func() error) error {
	dir := storeRoot(root)
	if e := ensureDirectories(dir); e != nil {
		return e
	}
	fd, e := unix.Open(filepath.Join(dir, "integration.lock"), unix.O_CREAT|unix.O_RDWR|unix.O_NOFOLLOW, 0600)
	if e != nil {
		return e
	}
	f := os.NewFile(uintptr(fd), "integration.lock")
	defer f.Close()
	if e = unix.Flock(fd, unix.LOCK_EX); e != nil {
		return e
	}
	defer unix.Flock(fd, unix.LOCK_UN)
	return fn()
}
func readReceipt(root, id string) (Receipt, error) {
	var r Receipt
	if !integrationID.MatchString(id) {
		return r, fmt.Errorf("invalid Integration ID")
	}
	if e := readJSON(receiptPath(root, id), &r); e != nil {
		return r, e
	}
	if r.Kind != "ply.integration.receipt" || r.SchemaVersion != 1 || r.ID != id || r.Plan.Workspace != root || digest(r.Plan) != r.PlanSHA256 || operationID(r.PlanSHA256) != id || r.Decision.PlanSHA256 != r.PlanSHA256 || r.Decision.IntegrationID != id {
		return r, fmt.Errorf("integration plan or decision integrity conflict")
	}
	b, e := readBounded(r.DecisionFile.Locator, 1<<20)
	if e != nil {
		return r, e
	}
	if byteHash(b) != r.DecisionFile.SHA256 || digest(r.Decision) != r.DecisionFile.SHA256 {
		return r, fmt.Errorf("integration decision bytes changed")
	}
	previous := ""
	for _, event := range r.Events {
		claimed := event.SHA256
		event.SHA256 = ""
		if event.PreviousSHA256 != previous || digest(event) != claimed || event.IntegrationID != id || event.CandidateOID != r.Plan.TaskResult.ResultOID {
			return r, fmt.Errorf("integration history integrity conflict")
		}
		previous = claimed
	}
	if e := validateReconsideration(r, r.State == "cancelled"); e != nil {
		return r, e
	}
	if r.State == "reconsidering" && r.Reconsideration == nil {
		return r, fmt.Errorf("reconsideration is missing its preserved negative decision")
	}
	return r, nil
}
func (s *Service) save(r *Receipt) error {
	r.UpdatedAtUTC = s.now()
	return atomicJSON(receiptPath(r.Plan.Workspace, r.ID), r, false)
}
func (s *Service) event(r *Receipt, kind string) {
	id := "ine_" + strings.TrimPrefix(byteHash([]byte(r.ID+"/"+kind)), "sha256:")
	for _, v := range r.Events {
		if v.ID == id {
			return
		}
	}
	e := Event{ID: id, Kind: kind, IntegrationID: r.ID, DeliveryID: r.Plan.DeliveryID, TaskID: r.Plan.TaskResult.TaskID, CandidateOID: r.Plan.TaskResult.ResultOID, TargetRef: r.Plan.Agreement.TargetRef, IntegratedOID: r.IntegratedOID, AtUTC: s.now()}
	if r.PR != nil {
		e.PRURL = r.PR.URL
	}
	if len(r.Events) > 0 {
		e.PreviousSHA256 = r.Events[len(r.Events)-1].SHA256
	}
	e.SHA256 = digest(e)
	r.Events = append(r.Events, e)
}

type effectReservation struct {
	Kind          string `json:"kind"`
	Key           string `json:"key"`
	IntegrationID string `json:"integration_id"`
}

func reserve(r Receipt) error {
	if !sha256Pattern.MatchString(r.Plan.EffectKey) {
		return fmt.Errorf("invalid integration effect key")
	}
	path := filepath.Join(storeRoot(r.Plan.Workspace), "effects", strings.TrimPrefix(r.Plan.EffectKey, "sha256:")+".json")
	var old effectReservation
	if e := readJSON(path, &old); e == nil {
		return reserveSupersession(r, old)
	} else if !os.IsNotExist(e) {
		return e
	}
	return atomicJSON(path, effectReservation{"ply.integration.reservation", r.Plan.EffectKey, r.ID}, true)
}
