package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestWorkflowExecuteVerifyHelpManagedRequirements(t *testing.T) {
	c := newWorkflowExecuteCommand(taskrun.Dependencies{})
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	c.SetArgs([]string{"verify", "--help"})
	if err := c.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Native verification assembles and preserves task-requirements",
		"frozen Spec requirements and this exact candidate's acceptance and review evidence",
		"Do not hand-author a second requirements/status artifact",
		"cite candidate-review for observations preserved in the review",
		"a progress report or guide does not qualify a candidate",
		"Without --review",
		"--review supplies an actual candidate-bound DeliveryCandidateReview@1 record",
		"Missing, failed or unknown tests cannot produce a pass",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("verify help lacks %q", want)
		}
	}
}
