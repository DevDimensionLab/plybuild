package workspace

import (
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"gopkg.in/yaml.v3"
)

// HistoricalRecord keeps a record of the removed delivery process exactly as it
// was stored. Ply neither reads nor validates its content; it only stays in
// place so a write never deletes frozen history.
type HistoricalRecord struct{ node yaml.Node }

func (h HistoricalRecord) MarshalYAML() (any, error) { return &h.node, nil }

func (h *HistoricalRecord) UnmarshalYAML(n *yaml.Node) error {
	h.node = *n
	return nil
}

// historicalAuthorityNode is decided by the recorded plan itself: only a plan of
// another kind than the supported ones was written by the removed process. Extra
// members never make a current authority historical.
func historicalAuthorityNode(n *yaml.Node) bool {
	plan := registryNodeField(n, "plan")
	if plan == nil {
		return false
	}
	kind := registryNodeField(plan, "kind")
	return kind != nil && kind.Value != currentIntegrationPlanKind && kind.Value != refreshIntegrationPlanKind
}

// historicalIntegrationKeys are the members only the removed process wrote on
// integration intents. An intent is accepted as historical only together with a
// historical authority.
var historicalIntegrationKeys = map[string]bool{"mode": true, "task_result_id": true, "human_qa_record_id": true}

func historicalIntegrationNode(n *yaml.Node) bool {
	if n.Kind != yaml.MappingNode {
		return false
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if historicalIntegrationKeys[n.Content[i].Value] {
			return true
		}
	}
	return false
}

// historicalGoalDelivery accepts the delivery member that earlier Goals carry.
// A Goal no longer selects delivery or acceptance; the member means nothing.
func historicalGoalDelivery(canonicaljson.Value) error { return nil }

// historyCopier is a registry value holding opaque history in an unexported
// field. A registry copy calls it with the original value of the same type.
type historyCopier interface{ copyHistory(from any) }

func (r *HistoricalRecord) copyHistory(from any) {
	r.node = *cloneYAMLNode(&[]yaml.Node{from.(HistoricalRecord).node}[0])
}
func (a *IntegrationAuthority) copyHistory(from any) {
	a.historical = cloneYAMLNode(from.(IntegrationAuthority).historical)
}
func (i *IntegrationIntent) copyHistory(from any) {
	i.historical = cloneYAMLNode(from.(IntegrationIntent).historical)
}

func cloneYAMLNode(n *yaml.Node) *yaml.Node {
	if n == nil {
		return nil
	}
	c := *n
	c.Alias = cloneYAMLNode(n.Alias)
	if n.Content != nil {
		c.Content = make([]*yaml.Node, len(n.Content))
		for i, child := range n.Content {
			c.Content[i] = cloneYAMLNode(child)
		}
	}
	return &c
}
