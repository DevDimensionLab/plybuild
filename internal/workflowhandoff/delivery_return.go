package workflowhandoff

type DeliveryHumanAttestation struct {
	Kind           string `json:"kind"`
	SchemaVersion  int    `json:"schema_version"`
	TaskID         string `json:"task_id"`
	TaskResultID   string `json:"task_result_id"`
	ResultOID      string `json:"result_oid"`
	ResultTree     string `json:"result_tree"`
	Outcome        string `json:"outcome"`
	ActorClaim     string `json:"actor_claim"`
	StartSurface   string `json:"start_surface"`
	StartedAtUTC   string `json:"started_at_utc"`
	CompletedAtUTC string `json:"completed_at_utc"`
	Answer         string `json:"answer"`
	Observation    string `json:"observation"`
}
