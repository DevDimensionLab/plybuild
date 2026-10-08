package taskrun

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

var digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
var runPattern = regexp.MustCompile(`^trn_[0-9a-f]{64}$`)
var keyPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._/-]{0,127}$`)

func hash(b []byte) string { x := sha256.Sum256(b); return "sha256:" + hex.EncodeToString(x[:]) }
func Canonical(v any) ([]byte, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	x, e := canonicaljson.DecodeStrict(b)
	if e != nil {
		return nil, e
	}
	return canonicaljson.Marshal(x)
}
func digest(v any) string {
	b, e := Canonical(v)
	if e != nil {
		panic(e)
	}
	return hash(b)
}
func plain(s string, min, max int) bool {
	if !utf8.ValidString(s) || strings.TrimSpace(s) != s || utf8.RuneCountInString(s) < min || utf8.RuneCountInString(s) > max {
		return false
	}
	for _, c := range s {
		if unicode.IsControl(c) {
			return false
		}
	}
	return true
}
func key(s string) bool {
	if !keyPattern.MatchString(s) || strings.Contains(s, "//") || strings.HasSuffix(s, "/") {
		return false
	}
	for _, p := range strings.Split(s, "/") {
		if p == ".." {
			return false
		}
	}
	return true
}

// exactShape requires all members (including explicit nullable members), at every
// typed level. Raw WF objects are checked by the authoritative WF validators.
func exactShape(v any, t reflect.Type) error {
	if t == reflect.TypeOf(json.RawMessage{}) {
		return nil
	}
	if t.Kind() == reflect.Pointer {
		if v == nil {
			return nil
		}
		return exactShape(v, t.Elem())
	}
	if v == nil {
		return fmt.Errorf("unexpected null for %s", t)
	}
	switch t.Kind() {
	case reflect.Struct:
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object")
		}
		fields := map[string]reflect.Type{}
		optional := map[string]bool{}
		var add func(reflect.Type)
		add = func(t reflect.Type) {
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.Anonymous {
					add(f.Type)
					continue
				}
				name := strings.Split(f.Tag.Get("json"), ",")[0]
				if name != "-" {
					fields[name] = f.Type
					optional[name] = strings.Contains(f.Tag.Get("json"), ",omitempty")
				}
			}
		}
		add(t)
		if len(m) > len(fields) {
			return fmt.Errorf("missing or unknown fields for %s", t)
		}
		for n, ft := range fields {
			x, ok := m[n]
			if !ok {
				if optional[n] {
					continue
				}
				return fmt.Errorf("missing %s", n)
			}
			if e := exactShape(x, ft); e != nil {
				return fmt.Errorf("%s: %w", n, e)
			}
		}
	case reflect.Slice:
		a, ok := v.([]any)
		if !ok {
			return fmt.Errorf("expected array")
		}
		for _, x := range a {
			if e := exactShape(x, t.Elem()); e != nil {
				return e
			}
		}
	}
	return nil
}
func decode(b []byte, max int, out any) error {
	if len(b) > max {
		return invalid("document exceeds size limit")
	}
	if _, e := canonicaljson.DecodeStrict(b); e != nil {
		return invalid(e.Error())
	}
	var shape any
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	if e := dec.Decode(&shape); e != nil {
		return invalid(e.Error())
	}
	if e := exactShape(shape, reflect.TypeOf(out).Elem()); e != nil {
		return invalid(e.Error())
	}
	dec = json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if e := dec.Decode(out); e != nil {
		return invalid(e.Error())
	}
	return nil
}
func checkEnvelope(e Envelope, kind string) error {
	if e != env(kind) {
		return invalid("unsupported kind or schema_version")
	}
	return nil
}
func physical(path string, missing bool) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path || strings.ContainsRune(path, 0) {
		return invalid("path must be absolute and clean")
	}
	for p := path; ; p = filepath.Dir(p) {
		info, e := os.Lstat(p)
		if e != nil {
			if missing && os.IsNotExist(e) {
			} else {
				return e
			}
		} else if info.Mode()&os.ModeSymlink != 0 {
			return conflict("symlink path rejected")
		}
		if p == filepath.Dir(p) {
			break
		}
	}
	return nil
}
func readFile(path string, max int, private bool) ([]byte, error) {
	if e := physical(path, false); e != nil {
		return nil, e
	}
	before, e := os.Lstat(path)
	if e != nil {
		return nil, e
	}
	if !before.Mode().IsRegular() || before.Size() > int64(max) || (private && before.Mode().Perm()&0077 != 0) {
		return nil, invalid("expected bounded private regular file")
	}
	f, e := openRead(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	actual, e := f.Stat()
	if e != nil || !os.SameFile(before, actual) {
		return nil, conflict("file changed while opening")
	}
	b, e := io.ReadAll(io.LimitReader(f, int64(max)+1))
	if e != nil {
		return nil, e
	}
	after, e := os.Lstat(path)
	if e != nil || !os.SameFile(actual, after) || len(b) > max {
		return nil, conflict("file changed while reading")
	}
	return b, nil
}
func ReadRequest(path string) (Request, error) {
	var r Request
	b, e := readFile(path, 1<<20, false)
	if e != nil {
		return r, e
	}
	return parseRequest(b)
}
func parseRequest(b []byte) (Request, error) {
	return parseRequestForSurface(b, false)
}

// Herdr shares the native Task contract, but owns its provider-specific launch.
// Ordinary terminal and factory requests must remain Codex-only.
func parseRequestForSurface(b []byte, herdr bool) (Request, error) {
	var r Request
	if e := decode(b, 1<<20, &r); e != nil {
		return r, e
	}
	if r.Kind != env("request").Kind || (r.SchemaVersion != 1 && r.SchemaVersion != 2) {
		return r, invalid("unsupported request version")
	}
	var fields map[string]json.RawMessage
	_ = json.Unmarshal(b, &fields)
	_, hasFactory := fields["factory_test"]
	if (r.SchemaVersion == 1 && hasFactory) || (r.SchemaVersion == 2 && r.FactoryTest == nil) {
		return r, invalid("factory_test is required only for request schema 2")
	}
	if !key(r.RequestKey) || !regexp.MustCompile(`^pre_[0-9a-f]{64}$`).MatchString(r.PreparationID) || !digestPattern.MatchString(r.PreparationSHA256) {
		return r, invalid("invalid request key or preparation binding")
	}
	if e := physical(r.WorkspaceRoot, false); e != nil {
		return r, e
	}
	rt := r.Runtime
	p := rt.PermissionBinding
	providerOK := rt.Provider == "codex" || herdr && r.SchemaVersion == 1 && rt.Provider == "claude" && rt.ConfigProfile == nil && p.ProfileID == "manual"
	if !providerOK || (r.SchemaVersion == 1 && rt.Mode != "interactive" || r.SchemaVersion == 2 && (rt.Mode != "exec" || rt.ConfigProfile != nil)) || !plain(rt.Model, 1, 256) || rt.ConfigProfile != nil && !plain(*rt.ConfigProfile, 1, 256) {
		return r, invalid("invalid provider runtime; Herdr supports codex or claude (null config_profile, manual permission mode); native starts require Codex")
	}
	if p.AuthorityKind != "reported_contract_with_effective_policy" || !plain(p.ProfileID, 1, 256) || !digestPattern.MatchString(p.EffectivePolicySHA256) || len(p.Evidence) == 0 {
		return r, invalid("effective policy evidence is required")
	}
	last := ""
	for _, x := range p.Evidence {
		if x.Locator <= last || !digestPattern.MatchString(x.SHA256) || (x.Role != "effective_policy" && x.Role != "permission_proof" && x.Role != "runtime_contract") {
			return r, invalid("policy evidence must be locator sorted, unique and typed")
		}
		last = x.Locator
	}
	if r.Agreement != (Agreement{"A", 3, 5400, 2}) || r.ReturnPolicy != (ReturnPolicy{true, true, "separate", "separate"}) || !r.HumanAuthority.Authorized || (r.SchemaVersion == 1 && r.HumanAuthority.StartSurface != "human_ordinary_terminal" || r.SchemaVersion == 2 && r.HumanAuthority.StartSurface != "human_started_factory_test") || !plain(r.HumanAuthority.ActorClaim, 1, 256) {
		return r, invalid("agreement A, separate human gates and human start authority are required")
	}
	if r.SchemaVersion == 2 {
		if p.ProfileID != "factory-test" {
			return r, invalid("exec requires the factory-test managed profile")
		}
		f := r.FactoryTest
		if !digestPattern.MatchString(f.AuthorizationSHA256) || !filepath.IsAbs(f.AuthorizationPath) || f.Iteration < 1 || f.Iteration > 2 || f.TaskSlot < 1 || f.TaskSlot > 2 || f.ReasoningEffort != "low" || f.TimeoutSeconds < 1 || f.TimeoutSeconds > 240 {
			return r, invalid("invalid factory slot or budget")
		}
	}
	for _, x := range []Executable{rt.Executable, rt.PlyExecutable} {
		if !digestPattern.MatchString(x.SHA256) || !filepath.IsAbs(x.Path) || filepath.Clean(x.Path) != x.Path {
			return r, invalid("invalid executable binding")
		}
	}
	if _, e := workflowhandoff.ValidateTaskRunDraft(r.HandoffDraft); e != nil {
		return r, invalid(e.Error())
	}
	return r, nil
}
func runtimeBindings(r Runtime) ([]FileBinding, error) {
	return runtimeOperationBindings(r, []Executable{r.Executable, r.PlyExecutable})
}

// Launching a provider requires its executable. An accepted owner callback
// instead validates its control and policy; the launch identity stays in the
// frozen request, even when a package upgrade removes or replaces that file.
func runtimeOperationBindings(r Runtime, executables []Executable) ([]FileBinding, error) {
	out := []FileBinding{}
	for _, x := range executables {
		if e := verifyExecutable(x); e != nil {
			return nil, e
		}
		out = append(out, FileBinding{x.Path, x.SHA256})
	}
	for _, x := range r.PermissionBinding.Evidence {
		b, e := readFile(x.Locator, 1<<20, true)
		if e != nil {
			return nil, e
		}
		if hash(b) != x.SHA256 {
			return nil, conflict("policy evidence changed")
		}
		out = append(out, FileBinding{x.Locator, x.SHA256})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Locator < out[j].Locator })
	for i := 1; i < len(out); i++ {
		if out[i].Locator == out[i-1].Locator {
			return nil, invalid("runtime binding locators must be unique")
		}
	}
	return out, nil
}
func verifyExecutable(x Executable) error {
	if e := physical(x.Path, false); e != nil {
		return e
	}
	f, e := openRead(x.Path)
	if e != nil {
		return e
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		return e
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		return invalid("bound executable is not a regular executable")
	}
	h := sha256.New()
	if _, e = io.Copy(h, f); e != nil {
		return e
	}
	if "sha256:"+hex.EncodeToString(h.Sum(nil)) != x.SHA256 {
		return conflict("executable bytes changed")
	}
	return nil
}
func RunID(key string) string { return "trn_" + strings.TrimPrefix(hash([]byte(key)), "sha256:") }
func runPaths(r Request) Paths {
	root := filepath.Join(r.WorkspaceRoot, ".ply", "task-runs", "v1", "runs", RunID(r.RequestKey))
	return Paths{root, filepath.Join(root, "tmp"), filepath.Join(root, "reports", "report.json")}
}
func confirmation(v any) string {
	b, _ := Canonical(v)
	var m map[string]json.RawMessage
	_ = json.Unmarshal(b, &m)
	delete(m, "confirmation")
	delete(m, "next_argv")
	return digest(m)
}
func sortedReasons(r []Reason) []Reason {
	sort.Slice(r, func(i, j int) bool {
		if r[i].Code == r[j].Code {
			return r[i].Detail < r[j].Detail
		}
		return r[i].Code < r[j].Code
	})
	out := []Reason{}
	for _, x := range r {
		if len(out) == 0 || out[len(out)-1] != x {
			out = append(out, x)
		}
	}
	return out
}
func budgetValid(b *BudgetUsage) *bool {
	if b == nil {
		return nil
	}
	ok := b.InitialExecutionStarted && b.CorrectionRounds <= 3 && b.ActiveSeconds <= 5400 && b.EnvironmentMeasures <= 2
	return &ok
}
func validateReport(b []byte) (Report, error) {
	var r Report
	if e := decode(b, 4<<20, &r); e != nil {
		return r, e
	}
	if e := checkEnvelope(r.Envelope, "report"); e != nil {
		return r, e
	}
	if e := workflowhandoff.ValidateTaskRunSemantics(b); e != nil {
		return r, invalid(e.Error())
	}
	if r.BudgetUsage != nil {
		x := r.BudgetUsage
		if x.CorrectionRounds < 0 || x.CorrectionRounds > 99 || x.ActiveSeconds < 0 || x.ActiveSeconds > 86400 || x.EnvironmentMeasures < 0 || x.EnvironmentMeasures > 99 || !x.InitialExecutionStarted && (x.CorrectionRounds > 0 || x.EnvironmentMeasures > 0) || r.Outcome == "complete" && !x.InitialExecutionStarted {
			return r, invalid("contradictory or invalid budget usage")
		}
	}
	for _, x := range []*string{r.ProcessObservation, r.PreventiveFollowup} {
		if x != nil && !plain(*x, 1, 1000) {
			return r, invalid("invalid process observation or preventive followup")
		}
	}
	technical, _ := Canonical(r.TechnicalAssessment)
	if e := workspace.ValidateTaskRunTechnicalAssessment(technical); e != nil {
		return r, invalid(e.Error())
	}
	return r, nil
}
