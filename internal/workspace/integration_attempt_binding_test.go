package workspace

import (
	"reflect"
	"strings"
	"testing"
)

func TestIntegrationAttemptBindingPreservesValidHistory(t *testing.T) {
	r := evidenceHistoryFixture(t)
	raw, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	got, err := decodeWorkItemRegistry(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.IntegrationAuthorities, r.IntegrationAuthorities) ||
		!reflect.DeepEqual(got.IntegrationIntents, r.IntegrationIntents) ||
		!reflect.DeepEqual(got.IntegrationAttempts, r.IntegrationAttempts) ||
		!reflect.DeepEqual(got.IntegrationResults, r.IntegrationResults) {
		t.Fatal("valid attempt/result history changed on roundtrip")
	}
}

func TestIntegrationObservationRejectCrossAuthorityAttempt(t *testing.T) {
	r := evidenceHistoryFixture(t)
	first, second := &r.IntegrationResults[0], &r.IntegrationResults[1]
	if first.AuthorityID == second.AuthorityID || first.IntentID == second.IntentID ||
		!reflect.DeepEqual(first.BeforeObservation, second.BeforeObservation) {
		t.Fatal("fixture needs distinct authorities and intents with identical observations")
	}
	first.AttemptID, second.AttemptID = second.AttemptID, first.AttemptID
	for i := range r.IntegrationAttempts {
		a := &r.IntegrationAttempts[i]
		if a.ID == *first.AttemptID {
			id := first.ID
			a.ResultID = &id
		}
		if a.ID == *second.AttemptID {
			id := second.ID
			a.ResultID = &id
		}
	}
	if _, err := encodeWorkItemRegistry(r); err == nil {
		t.Fatal("cross-authority attempt/result pairs accepted despite identical observations")
	} else if !strings.Contains(err.Error(), "has invalid attempt") {
		t.Fatalf("expected attempt binding rejection, got: %v", err)
	}
}
