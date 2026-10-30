package taskrun

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
)

// Presentation comes from frozen pre-start Task data, never the Run ID (which
// includes the request) or a tab's current name/focus. No identity is rebound.
var taskTitleCode = regexp.MustCompile(`^([A-Z]{2,8}[0-9]{1,8}[A-Z]{0,4})(?:[:\s]|$)`)

func (g deliveryReceiptGoal) presentation() (code, title string) {
	title = strings.Join(strings.Fields(g.Title), " ")
	if match := taskTitleCode.FindStringSubmatch(title); match != nil {
		return match[1], strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(title, match[1]), ":"))
	}
	sum := sha256.Sum256([]byte(g.TaskID))
	return fmt.Sprintf("T%X", sum[:4]), title
}

func deliveryTaskPresentation(s workflowState) (deliveryReceiptGoal, bool) {
	if s.Request.Delivery == nil || s.Request.Delivery.StartupProtocol != MandateReceiptV1 {
		return deliveryReceiptGoal{}, false
	}
	g, err := readDeliveryReceiptGoal(s.Request.Delivery.Goal)
	return g, err == nil && g.TaskID != ""
}

func workflowTaskTabLabel(s workflowState, role string) string {
	if g, ok := deliveryTaskPresentation(s); ok {
		code, title := g.presentation()
		prefix := code + " " + role
		const limit = 40
		room := limit - len([]rune(prefix)) - 1
		if room <= 0 || title == "" {
			return prefix
		}
		runes := []rune(title)
		if len(runes) > room {
			title = string(runes[:room-1]) + "…"
		}
		return prefix + " " + title
	}
	// Historical/report-only transports keep their original labels and contracts.
	return role + " " + s.Request.Herdr.TabLabel
}

func deliveryTaskCodeInstructions(g deliveryReceiptGoal) string {
	code, _ := g.presentation()
	return fmt.Sprintf("Use task code %s in references and returns. This is a presentation label; native Task/Run/session identity and authority remain unchanged. Full task title: %s. Task: %s.\n\n", code, g.Title, g.TaskID)
}

func workflowTaskOwnerTabLabel(s workflowState) string {
	if _, ok := deliveryTaskPresentation(s); ok {
		return workflowTaskTabLabel(s, "owner")
	}
	return "run " + s.Request.Herdr.TabLabel
}
