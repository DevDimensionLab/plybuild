package agentwork

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

const maxStoreBytes = 8 << 20

type operation struct {
	Digest string   `json:"digest"`
	Result Followup `json:"result"`
}
type record struct {
	Value      Followup             `json:"value"`
	TargetRefs []string             `json:"target_refs"`
	Target     Identity             `json:"target"`
	Draft      Draft                `json:"draft"`
	Requests   map[string]operation `json:"requests"`
}
type database struct {
	Sequence int               `json:"sequence"`
	Records  map[string]record `json:"records"`
	Keys     map[string]string `json:"keys"`
}

func (s *Service) dir() string { return filepath.Join(s.Root, ".ply", "agent-work", "v1") }
func (s *Service) read() (database, error) {
	db := database{Records: map[string]record{}, Keys: map[string]string{}}
	dir := filepath.Join(s.dir(), "revisions")
	if e := taskrun.PhysicalPath(dir); e != nil {
		if os.IsNotExist(e) {
			return db, nil
		}
		return db, fmt.Errorf("followup_store_unavailable")
	}
	entries, e := os.ReadDir(dir)
	if e != nil {
		return db, fmt.Errorf("followup_store_unavailable")
	}
	names := []string{}
	for _, v := range entries {
		if len(v.Name()) > 0 && v.Name()[0] == '.' {
			continue
		}
		if !regexp.MustCompile(`^[0-9]{12}\.json$`).MatchString(v.Name()) || !v.Type().IsRegular() {
			return db, fmt.Errorf("followup_store_invalid")
		}
		names = append(names, v.Name())
	}
	if len(names) == 0 {
		return db, nil
	}
	sort.Strings(names)
	name := names[len(names)-1]
	b, e := taskrun.PrivateRead(filepath.Join(dir, name), maxStoreBytes)
	if e != nil {
		return db, fmt.Errorf("followup_store_unavailable")
	}
	if e = json.Unmarshal(b, &db); e != nil || db.Records == nil || db.Keys == nil || fmt.Sprintf("%012d.json", db.Sequence) != name {
		return db, fmt.Errorf("followup_store_invalid")
	}
	return db, nil
}
func (s *Service) save(db *database) error {
	db.Sequence++
	b, e := json.Marshal(db)
	if e != nil {
		return e
	}
	if len(b) > maxStoreBytes {
		return fmt.Errorf("followup_store_limit")
	}
	return taskrun.PrivateWriteOnce(filepath.Join(s.dir(), "revisions", fmt.Sprintf("%012d.json", db.Sequence)), b)
}
func (s *Service) locked(f func(*database) error) error {
	return taskrun.WithPrivateLock(s.dir(), func() error {
		db, e := s.read()
		if e != nil {
			return e
		}
		return f(&db)
	})
}
func validID(id string) bool { return regexp.MustCompile(`^fup_[a-f0-9]{64}$`).MatchString(id) }
func (s *Service) Show(id string) (Followup, error) {
	if !validID(id) {
		return Followup{}, fmt.Errorf("followup_invalid_id")
	}
	db, e := s.read()
	if e != nil {
		return Followup{}, e
	}
	r, ok := db.Records[id]
	if !ok {
		return Followup{}, fmt.Errorf("followup_not_found")
	}
	return s.public(r.Value), nil
}
func (s *Service) List(due bool) (List, error) {
	o := List{ListKind, 1, s.now(), []Followup{}}
	db, e := s.read()
	if e != nil {
		return o, e
	}
	ids := []string{}
	for id := range db.Records {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	n := s.native()
	for _, id := range ids {
		f := db.Records[id].Value
		if !due || f.State == "waiting" && f.Transport.State == "not_sent" && s.isDue(f) {
			o.Items = append(o.Items, s.publicWithNative(f, n))
		}
	}
	return o, nil
}
func cycleRef(id string, cycle int) string { return id + "/cycle/" + strconv.Itoa(cycle) }

// RequestStatus preserves the original mutation result after a lost reply. Show
// separately returns current state; neither operation retries an effect.
func (s *Service) RequestStatus(id, requestID string) (Followup, error) {
	if !validID(id) || !reference(requestID) {
		return Followup{}, fmt.Errorf("followup_invalid_request_identity")
	}
	db, e := s.read()
	if e != nil {
		return Followup{}, e
	}
	r, ok := db.Records[id]
	if !ok {
		return Followup{}, fmt.Errorf("followup_not_found")
	}
	p, ok := r.Requests[requestID]
	if !ok {
		return Followup{}, fmt.Errorf("followup_request_not_recorded")
	}
	return s.public(p.Result), nil
}
