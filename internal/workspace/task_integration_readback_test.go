package workspace

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func TestTaskIntegrationReadbackIsCanonicalAndHasOneAction(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	result, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := MarshalTaskIntegrationReadback(result.Readback)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := canonicaljson.DecodeStrict(encoded); err != nil || !bytes.Contains(encoded, []byte(`"kind":"WorkspaceTaskIntegrationReadback@1"`)) || bytes.HasSuffix(encoded, []byte("\n")) {
		t.Fatalf("readback = %s, error = %v", encoded, err)
	}
	text := RenderTaskIntegrationText(result.Readback)
	if strings.Count(text, "\n") != 6 || !strings.Contains(text, "Git changed: no") || !strings.Contains(text, "--confirm sha256:") {
		t.Fatalf("text = %q", text)
	}
}

func TestFormatTwoTaskShowUsesNullableIntegrationReadback(t *testing.T) {
	fixture := newWorkItemJourneyFixture(t)
	result, err := ShowTask(fixture.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	if result.Integration == nil || result.Integration.NextAction.Kind != "record_task_result" {
		t.Fatalf("integration = %#v", result.Integration)
	}
	encoded, err := MarshalTaskReadback(result)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(encoded, []byte(`"task_result":null`)) || !bytes.Contains(encoded, []byte(`"plan":null`)) {
		t.Fatalf("readback = %s", encoded)
	}
}

func TestTaskShowPreservesDurableIntegrationWhenLiveObservationFails(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	if _, err := ApplyTaskIntegration(fixture.dependencies, input); err != nil {
		t.Fatal(err)
	}
	fixture.dependencies.IntegrationGit = unsupportedIntegrationGit{TaskIntegrationGit: fixture.dependencies.IntegrationGit}

	shown, err := ShowTask(fixture.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	if shown.Integration == nil || shown.Integration.Classification != "exact_effect" || shown.Integration.RecoveryStatus != "complete" || shown.Integration.NextAction.Kind != "none" {
		t.Fatalf("integration = %#v", shown.Integration)
	}
	encoded, err := MarshalTaskReadback(shown)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"authority":{`, `"intent":{`, `"attempt":{`, `"integration_result":{`, `"plan":null`} {
		if !bytes.Contains(encoded, []byte(field)) {
			t.Fatalf("readback does not preserve %s: %s", field, encoded)
		}
	}
}

func TestTaskShowKeepsDurableOutcomeSeparateFromFreshSourceObservation(t *testing.T) {
	fixture := newIntegrationJourneyFixture(t)
	input := TaskIntegrationInput{TaskID: "task", TaskResultID: fixture.resultID, HumanQARecordID: fixture.qaID, ExpectedResultOID: fixture.resultOID, ExpectedParentOID: fixture.oid}
	checked, err := CheckTaskIntegration(fixture.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	if _, err := ApplyTaskIntegration(fixture.dependencies, input); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, fixture.taskPath, "commit", "--allow-empty", "-m", "source moved after durable result")
	currentSourceOID := runLocalGit(t, fixture.taskPath, "rev-parse", "HEAD")

	shown, err := ShowTask(fixture.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	if shown.Integration == nil || shown.Integration.Classification != "exact_effect" || shown.Integration.RecoveryStatus != "complete" || shown.Integration.NextAction.Kind != "none" {
		t.Fatalf("integration = %#v", shown.Integration)
	}
	observedSource := integrationReadbackObject(t, integrationReadbackMember(t, shown.Integration.Value, "observed_source"), "observed_source")
	values := integrationReadbackObject(t, integrationReadbackMember(t, observedSource, "values"), "observed_source.values")
	freshness := integrationReadbackObject(t, integrationReadbackMember(t, observedSource, "freshness"), "observed_source.freshness")
	if got := integrationReadbackMember(t, values, "oid"); got != currentSourceOID {
		t.Fatalf("fresh observed source OID = %#v, want %s", got, currentSourceOID)
	}
	if got := integrationReadbackMember(t, freshness, "oid"); got != "fresh" {
		t.Fatalf("fresh observed source OID freshness = %#v", got)
	}
	authority := integrationReadbackObject(t, integrationReadbackMember(t, shown.Integration.Value, "authority"), "authority")
	plan := integrationReadbackObject(t, integrationReadbackMember(t, authority, "plan"), "authority.plan")
	authoritySource := integrationReadbackObject(t, integrationReadbackMember(t, plan, "observed_source"), "authority.plan.observed_source")
	if got := integrationReadbackMember(t, authoritySource, "oid"); got != fixture.resultOID {
		t.Fatalf("durable authority source OID = %#v, want %s", got, fixture.resultOID)
	}
	result := integrationReadbackObject(t, integrationReadbackMember(t, shown.Integration.Value, "integration_result"), "integration_result")
	if got := integrationReadbackMember(t, result, "outcome"); got != "exact_effect" {
		t.Fatalf("durable integration outcome = %#v", got)
	}
}
