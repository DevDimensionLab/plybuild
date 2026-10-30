package taskrun

import (
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type reversedDeliveryQAIDs struct {
	workspace.TaskLifecycleIDSource
	count int
}

func (s *reversedDeliveryQAIDs) NewHumanQARecordID() (workspace.HumanQARecordID, error) {
	s.count++
	if s.count == 1 {
		return workspace.HumanQARecordID("hqa_" + strings.Repeat("f", 32)), nil
	}
	return workspace.HumanQARecordID("hqa_" + strings.Repeat("0", 32)), nil
}
