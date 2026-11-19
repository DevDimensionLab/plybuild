// Package taskactor is Ply's one rule for who acted on a Task when nobody
// named an actor, and the role a Task's agent works in. Recording uses it for
// new records; plan trace uses it for records stored with actor unknown.
package taskactor

import (
	"encoding/json"
	"os"

	"github.com/devdimensionlab/plybuild/internal/agentstart"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Unknown is the actor when the rule finds none.
const Unknown = "unknown"

// Roles a Task can be registered with; developer is the default.
const (
	RoleDeveloper = "developer"
	RoleReviewer  = "reviewer"
)

// Agent is the actor goal-start gives the agent it starts for task.
func Agent(provider, task string) string {
	return provider + ":" + task
}

// StartedAgent is the agent goal-start started for task according to its
// receipt, or "" without one.
func StartedAgent(work workspace.ProjectWork, task string) string {
	b, err := os.ReadFile(workspace.GoalStartReceiptPath(work, workspace.TaskID(task)))
	if err != nil {
		return ""
	}
	var r agentstart.Receipt
	if json.Unmarshal(b, &r) != nil || r.Provider == "" || r.Key != task {
		return ""
	}
	return Agent(r.Provider, task)
}

// Of is the actor of a record nobody named an actor for: Ply's own
// operations get the owner of the Task's plan, the agent's records get the
// started agent. An empty candidate gives Unknown.
func Of(plyOperation bool, agent, planOwner string) string {
	actor := agent
	if plyOperation {
		actor = planOwner
	}
	if actor == "" {
		return Unknown
	}
	return actor
}

// Role is the role the latest goal register gave the Task, else developer.
func Role(records []taskjournal.Contribution) string {
	role := RoleDeveloper
	for _, r := range records {
		if r.Event != nil && r.Event.Ply != nil && r.Event.Ply.Operation == "goal_register" && r.Event.Ply.Refs["role"] != "" {
			role = r.Event.Ply.Refs["role"]
		}
	}
	return role
}
