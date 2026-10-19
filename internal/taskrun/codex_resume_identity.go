package taskrun

// CodexResumeIdentity is a local runtime observation, never Herdr hook metadata.
// Codex 0.163 can defer SessionStart until the first turn after thread/resume.
// Its exclusive thread writer lock exists before that turn. The observed pane,
// process executable and exclusive lock together identify the live conversation.
type CodexResumeIdentity struct {
	Method        string     `json:"method"`
	NativeSession string     `json:"native_session"`
	PaneID        string     `json:"pane_id"`
	PID           int        `json:"pid"`
	Executable    Executable `json:"executable"`
	LockPath      string     `json:"lock_path"`
	ObservedAtUTC string     `json:"observed_at_utc"`
}
