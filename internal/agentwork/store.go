package agentwork

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/privatestore"
	"github.com/devdimensionlab/plybuild/internal/workspace"
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

// stores are the agent-work directories of every defined Project, sorted by
// Project. A follow-up lives in exactly one of them; its ID includes the
// Project, so IDs are unique across the workspace.
func (s *Service) stores() ([]workspace.ProjectWork, error) { return workspace.DefinedWork(s.Root) }

// store is the agent-work directory of project. The wrapper must exist.
func (s *Service) store(project string) (string, error) {
	works, e := s.stores()
	if e != nil {
		return "", e
	}
	for _, w := range works {
		if string(w.Project) == project {
			if e := w.Create(); e != nil {
				return "", e
			}
			return w.AgentWork(), nil
		}
	}
	return "", fmt.Errorf("followup_project_not_found")
}

// locate is the store holding follow-up id, searched in every Project.
func (s *Service) locate(id string) (string, record, error) {
	works, e := s.stores()
	if e != nil {
		return "", record{}, e
	}
	for _, w := range works {
		db, e := read(w.AgentWork())
		if e != nil {
			return "", record{}, e
		}
		if r, ok := db.Records[id]; ok {
			return w.AgentWork(), r, nil
		}
	}
	return "", record{}, fmt.Errorf("followup_not_found")
}
func read(store string) (database, error) {
	db := database{Records: map[string]record{}, Keys: map[string]string{}}
	dir := filepath.Join(store, "revisions")
	if e := privatestore.PhysicalPath(dir); e != nil {
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
	b, e := privatestore.PrivateRead(filepath.Join(dir, name), maxStoreBytes)
	if e != nil {
		return db, fmt.Errorf("followup_store_unavailable")
	}
	if e = json.Unmarshal(b, &db); e != nil || db.Records == nil || db.Keys == nil || fmt.Sprintf("%012d.json", db.Sequence) != name {
		return db, fmt.Errorf("followup_store_invalid")
	}
	return db, nil
}
func save(store string, db *database) error {
	db.Sequence++
	b, e := json.Marshal(db)
	if e != nil {
		return e
	}
	if len(b) > maxStoreBytes {
		return fmt.Errorf("followup_store_limit")
	}
	return privatestore.PrivateWriteOnce(filepath.Join(store, "revisions", fmt.Sprintf("%012d.json", db.Sequence)), b)
}
func locked(store string, f func(*database) error) error {
	return privatestore.WithPrivateLock(store, func() error {
		db, e := read(store)
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
	_, r, e := s.locate(id)
	if e != nil {
		return Followup{}, e
	}
	return s.public(r.Value), nil
}
func (s *Service) List(due bool) (List, error) {
	return s.list(due, s.native)
}

// list projects each row against one request-local native basis.
func (s *Service) list(due bool, native func() NativeSnapshot) (List, error) {
	o := List{ListKind, 1, s.now(), []Followup{}}
	works, e := s.stores()
	if e != nil {
		return o, e
	}
	records := map[string]record{}
	for _, w := range works {
		db, e := read(w.AgentWork())
		if e != nil {
			return o, e
		}
		for id, r := range db.Records {
			records[id] = r
		}
	}
	ids := []string{}
	for id := range records {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	n := native()
	for _, id := range ids {
		f := records[id].Value
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
	_, r, e := s.locate(id)
	if e != nil {
		return Followup{}, e
	}
	p, ok := r.Requests[requestID]
	if !ok {
		return Followup{}, fmt.Errorf("followup_request_not_recorded")
	}
	return s.public(p.Result), nil
}
