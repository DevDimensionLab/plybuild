package workspace

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestHistoricalRegistryPoliciesRoundtrip(t *testing.T) {
	for _, policies := range []string{"[]", "[{task_id: historical-task, mode: required, activation_publication_key: null, legacy_result_ids: []}]"} {
		r := emptyWorkItemRegistry()
		b, err := encodeWorkItemRegistry(r)
		if err != nil {
			t.Fatal(err)
		}
		b = bytes.Replace(b, []byte("task_problem_revisions:"), []byte("task_spec_policies: "+policies+"\ntask_problem_revisions:"), 1)
		restored, err := decodeWorkItemRegistry(b)
		if err != nil {
			t.Fatal(err)
		}
		out, err := encodeWorkItemRegistry(restored)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(restored.TaskSpecPolicies, func() *[]TaskSpecPolicy {
			again, e := decodeWorkItemRegistry(out)
			if e != nil {
				t.Fatal(e)
			}
			return again.TaskSpecPolicies
		}()) {
			t.Fatal("policy values changed")
		}
		if !bytes.Contains(out, []byte("task_spec_policies:")) {
			t.Fatal("historical policies dropped")
		}
		again, err := decodeWorkItemRegistry(out)
		if err != nil {
			t.Fatal(err)
		}
		stable, err := encodeWorkItemRegistry(again)
		if err != nil || !bytes.Equal(out, stable) {
			t.Fatal("history changed on roundtrip", err)
		}
		bad := strings.Replace(string(out), "task_spec_policies:", "unknown_policies:", 1)
		if _, err := decodeWorkItemRegistry([]byte(bad)); RegistryError(err) == nil || RegistryError(err).Code != "registry_unknown_field" {
			t.Fatal("unknown field accepted", err)
		}
	}
}

func TestConsumerRegistryHistoricalSnapshot(t *testing.T) {
	path := os.Getenv("PLY_CONSUMER_REGISTRY_FIXTURE")
	if path == "" {
		t.Skip("private consumer snapshot not supplied")
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	r, err := decodeWorkItemRegistry(before)
	if err != nil {
		t.Fatal(err)
	}
	out, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := decodeWorkItemRegistry(out)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r, restored) {
		t.Fatal("historical registry content changed")
	}
	for i, a := range r.IntegrationAuthorities {
		if integrationPlanDigest(a.Plan) != integrationPlanDigest(restored.IntegrationAuthorities[i].Plan) {
			t.Fatal("historical plan digest changed")
		}
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("read changed snapshot")
	}
}

func TestHistoricalGuardTamperingIsRejected(t *testing.T) {
	path := os.Getenv("PLY_CONSUMER_REGISTRY_FIXTURE")
	if path == "" {
		t.Skip("private consumer snapshot not supplied")
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	tampered := bytes.Replace(before, []byte("relevance: current"), []byte("relevance: stale"), 1)
	if bytes.Equal(before, tampered) {
		t.Fatal("fixture does not contain expected historical relevance")
	}
	if _, err := decodeWorkItemRegistry(tampered); err == nil {
		t.Fatal("changed history with original digest accepted")
	}
}
