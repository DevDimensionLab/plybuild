package workspace

import (
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"strings"
	"testing"
)

func TestAL10RequirementEvidenceOutcomesByReportCompleteness(t *testing.T) {
	oid, tree := strings.Repeat("a", 40), strings.Repeat("b", 40)
	basis := TaskSpecBasis{TaskID: "task", SpecID: "spec", Spec: TaskRevisionRef{Revision: 1, ManifestSHA256: "sha256:" + strings.Repeat("c", 64)}}
	spec := canonicaljson.Object{{Name: "requirements", Value: []canonicaljson.Value{canonicaljson.Object{{Name: "id", Value: "ac-01"}}, canonicaljson.Object{{Name: "id", Value: "ac-02"}}}}}
	verifiers := []TaskVerifierResultRecord{{VerifierID: "test", Outcome: "failed"}}
	evidence := func(second string, mutate func(canonicaljson.Object) canonicaljson.Object) []byte {
		entry := func(id, outcome string) canonicaljson.Value {
			return canonicaljson.Object{{Name: "id", Value: id}, {Name: "outcome", Value: outcome}, {Name: "verifier_ids", Value: []canonicaljson.Value{"test"}}, {Name: "artifact_ids", Value: []canonicaljson.Value{}}, {Name: "reason", Value: "Actual observation."}}
		}
		v := canonicaljson.Object{{Name: "kind", Value: "WorkspaceTaskRequirementEvidence@1"}, {Name: "schema_version", Value: int64(1)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"}, {Name: "task_id", Value: "task"}, {Name: "spec_id", Value: "spec"}, {Name: "spec", Value: canonicaljson.Object{{Name: "revision", Value: int64(1)}, {Name: "manifest_sha256", Value: basis.Spec.ManifestSHA256}}}, {Name: "result_oid", Value: oid}, {Name: "result_tree", Value: tree}, {Name: "requirements", Value: []canonicaljson.Value{entry("ac-01", "passed"), entry("ac-02", second)}}}
		if mutate != nil {
			v = mutate(v)
		}
		b, err := canonicaljson.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	for _, outcome := range []string{"failed", "not_run", "unknown"} {
		raw := evidence(outcome, nil)
		if err := ValidateTaskRequirementEvidence(raw, basis, spec, oid, tree, verifiers, nil); err == nil || !strings.Contains(err.Error(), "requirement ac-02 has outcome "+outcome) {
			t.Fatalf("complete evidence accepted %s: %v", outcome, err)
		}
		if err := ValidateIncompleteTaskRequirementEvidence(raw, basis, spec, oid, tree, verifiers, nil); err != nil {
			t.Fatalf("incomplete report refused mixed passed/%s evidence: %v", outcome, err)
		}
	}
	replace := func(name string, value canonicaljson.Value) func(canonicaljson.Object) canonicaljson.Object {
		return func(v canonicaljson.Object) canonicaljson.Object {
			for i := range v {
				if v[i].Name == name {
					v[i].Value = value
				}
			}
			return v
		}
	}
	for name, raw := range map[string][]byte{
		"candidate": evidence("not_run", replace("result_oid", strings.Repeat("0", 40))),
		"spec":      evidence("not_run", replace("spec_id", "other")),
		"id set":    evidence("not_run", replace("requirements", []canonicaljson.Value{})),
	} {
		if err := ValidateIncompleteTaskRequirementEvidence(raw, basis, spec, oid, tree, verifiers, nil); err == nil {
			t.Fatalf("incomplete evidence with wrong %s binding accepted", name)
		}
	}
	if err := ValidateIncompleteTaskRequirementEvidence(evidence("not_run", nil), basis, spec, oid, tree, nil, nil); err == nil {
		t.Fatal("incomplete evidence citing an absent verifier accepted")
	}
}
