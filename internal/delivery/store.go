package delivery

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"golang.org/x/sys/unix"
)

var idPattern = regexp.MustCompile(`^dlv_[a-f0-9]{64}$`)

func hash(data []byte) string { return fmt.Sprintf("sha256:%x", sha256.Sum256(data)) }
func canonical(v any) ([]byte, error) {
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	value, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		return nil, e
	}
	return canonicaljson.Marshal(value)
}
func digest(v any) string {
	raw, e := canonical(v)
	if e != nil {
		panic(e)
	}
	return hash(raw)
}
func shortHash(v string) string { return strings.TrimPrefix(hash([]byte(v)), "sha256:") }

func readJSON(path string, out any) error {
	raw, e := readBounded(path, 16<<20)
	if e != nil {
		return e
	}
	if _, e = canonicaljson.DecodeStrict(raw); e != nil {
		return fmt.Errorf("invalid JSON %s: %w", path, e)
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if e = d.Decode(out); e != nil {
		return e
	}
	if e = d.Decode(new(any)); e != io.EOF {
		return fmt.Errorf("unexpected trailing JSON")
	}
	return nil
}

func readBounded(path string, limit int64) ([]byte, error) {
	abs, e := filepath.Abs(path)
	if e != nil {
		return nil, e
	}
	actual, e := filepath.EvalSymlinks(abs)
	if e != nil {
		return nil, e
	}
	if actual != abs {
		return nil, fmt.Errorf("delivery path contains a symlink: %s", path)
	}
	before, e := os.Lstat(abs)
	if e != nil {
		return nil, e
	}
	if !before.Mode().IsRegular() {
		return nil, fmt.Errorf("delivery path is not a regular file: %s", path)
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
		return nil, fmt.Errorf("delivery file changed before read")
	}
	raw, e := io.ReadAll(io.LimitReader(f, limit+1))
	if e != nil {
		return nil, e
	}
	if int64(len(raw)) > limit {
		return nil, fmt.Errorf("delivery file exceeds %d bytes", limit)
	}
	after, e := f.Stat()
	if e != nil {
		return nil, e
	}
	current, e := os.Lstat(abs)
	if e != nil {
		return nil, e
	}
	if !os.SameFile(opened, current) || opened.ModTime() != after.ModTime() || opened.Size() != int64(len(raw)) || after.Size() != opened.Size() {
		return nil, fmt.Errorf("delivery file changed during read")
	}
	return raw, nil
}

func storeRoot(root string) string      { return filepath.Join(root, ".ply", "deliveries", "v1") }
func recordPath(root, id string) string { return filepath.Join(storeRoot(root), "records", id+".json") }

func physicalDirectories(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("delivery storage path must be absolute and clean")
	}
	current := string(filepath.Separator)
	for _, part := range strings.Split(strings.TrimPrefix(path, current), string(filepath.Separator)) {
		if part == "" {
			continue
		}
		current = filepath.Join(current, part)
		info, e := os.Lstat(current)
		if os.IsNotExist(e) {
			if e = os.Mkdir(current, 0700); e != nil && !os.IsExist(e) {
				return e
			}
			info, e = os.Lstat(current)
		}
		if e != nil {
			return e
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("delivery storage contains a symlink or non-directory: %s", current)
		}
	}
	return nil
}

func atomicJSON(path string, value any, immutable bool) error {
	raw, e := canonical(value)
	if e != nil {
		return e
	}
	if e = physicalDirectories(filepath.Dir(path)); e != nil {
		return e
	}
	if previous, err := readBounded(path, 16<<20); err == nil {
		if bytes.Equal(previous, raw) {
			return nil
		}
		if immutable {
			return fmt.Errorf("immutable delivery content conflict: %s", path)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	f, e := os.CreateTemp(filepath.Dir(path), ".publish-")
	if e != nil {
		return e
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if e = f.Chmod(0600); e == nil {
		_, e = f.Write(raw)
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
		if e = os.Chmod(tmp, 0400); e != nil {
			return e
		}
	}
	if e = os.Rename(tmp, path); e != nil {
		return e
	}
	d, e := os.Open(filepath.Dir(path))
	if e != nil {
		return e
	}
	defer d.Close()
	return d.Sync()
}

// A workspace-wide advisory lock protects effect reservations across IDs as
// well as concurrent processes. Atomic snapshots preserve every event already
// acknowledged; no command is launched before its intent is durable.
func withLock(root string, fn func() error) error {
	dir := storeRoot(root)
	if e := physicalDirectories(dir); e != nil {
		return e
	}
	path := filepath.Join(dir, "delivery.lock")
	fd, e := unix.Open(path, unix.O_CREAT|unix.O_RDWR|unix.O_NOFOLLOW, 0600)
	if e != nil {
		return e
	}
	f := os.NewFile(uintptr(fd), path)
	defer f.Close()
	if e = unix.Flock(fd, unix.LOCK_EX); e != nil {
		return e
	}
	defer unix.Flock(fd, unix.LOCK_UN)
	return fn()
}

func readReceipt(root, id string) (Receipt, error) {
	var r Receipt
	if !idPattern.MatchString(id) {
		return r, fmt.Errorf("invalid Delivery ID")
	}
	if e := readJSON(recordPath(root, id), &r); e != nil {
		return r, e
	}
	if r.Kind != "ply.delivery.receipt" || r.SchemaVersion != 1 || r.ID != id || r.Manifest.Workspace != root || r.RegistrationSHA256 != digest(r.Registration) || r.ManifestSHA256 != digest(r.Manifest) {
		return r, fmt.Errorf("delivery manifest integrity conflict")
	}
	previous := ""
	for _, ev := range r.Events {
		claimed := ev.SHA256
		ev.SHA256 = ""
		if ev.PreviousSHA256 != previous || claimed != digest(ev) || ev.DeliveryID != r.ID || ev.CandidateOID != r.Manifest.TaskResult.ResultOID {
			return r, fmt.Errorf("delivery history integrity conflict")
		}
		previous = claimed
	}
	return r, nil
}

func (s *Service) save(root string, r *Receipt) error {
	r.UpdatedAtUTC = s.now()
	return atomicJSON(recordPath(root, r.ID), r, false)
}
func (s *Service) event(r *Receipt, kind, key, detail, url string) {
	id := "dle_" + shortHash(r.ID+"/"+key)
	for _, e := range r.Events {
		if e.ID == id {
			return
		}
	}
	e := Event{ID: id, Kind: kind, AtUTC: s.now(), Actor: r.Manifest.OwnerClaim, Origin: "ply.delivery", DeliveryID: r.ID, CandidateOID: r.Manifest.TaskResult.ResultOID, AttemptID: r.AttemptID, Detail: detail, URL: url}
	if len(r.Events) > 0 {
		e.PreviousSHA256 = r.Events[len(r.Events)-1].SHA256
	}
	e.SHA256 = digest(e)
	r.Events = append(r.Events, e)
}

type effect struct {
	Kind            string        `json:"kind"`
	SchemaVersion   int           `json:"schema_version"`
	Key             string        `json:"key"`
	OwnerDeliveryID string        `json:"owner_delivery_id"`
	PushStarted     bool          `json:"push_started"`
	PushObserved    bool          `json:"push_observed"`
	CreateStarted   bool          `json:"create_started"`
	LocalStarted    bool          `json:"local_started"`
	PR              *PRReceipt    `json:"pull_request"`
	Local           *LocalReceipt `json:"local_integration"`
	Complete        bool          `json:"complete"`
}

func effectPath(root, key string) string {
	return filepath.Join(storeRoot(root), "effects", strings.TrimPrefix(key, "sha256:")+".json")
}
func readEffect(root, key string) (effect, error) {
	var e effect
	err := readJSON(effectPath(root, key), &e)
	if err == nil && (e.Key != key || e.Kind != "ply.delivery.effect" || e.SchemaVersion != 1 || !idPattern.MatchString(e.OwnerDeliveryID)) {
		err = fmt.Errorf("delivery effect identity conflict")
	}
	return e, err
}
func saveEffect(root string, e effect) error { return atomicJSON(effectPath(root, e.Key), e, false) }

func listReceipts(root string) ([]Receipt, error) {
	entries, e := os.ReadDir(filepath.Join(storeRoot(root), "records"))
	if os.IsNotExist(e) {
		return []Receipt{}, nil
	}
	if e != nil {
		return nil, e
	}
	result := []Receipt{}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".publish-") {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".json")
		r, e := readReceipt(root, id)
		if e != nil {
			return nil, e
		}
		result = append(result, r)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].CreatedAtUTC == result[j].CreatedAtUTC {
			return result[i].ID < result[j].ID
		}
		return result[i].CreatedAtUTC < result[j].CreatedAtUTC
	})
	return result, nil
}
