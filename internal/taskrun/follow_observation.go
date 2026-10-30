package taskrun

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

const workflowFollowInterval = 250 * time.Millisecond
const workflowHerdrReadTimeout = 15 * time.Second

var workflowObservationWindowExpired = workflowError(5, "Run remains bound to the same session; observation timed out without stopping or restarting it; current provider liveness is not established")

// One read continuation belongs to this invocation, never to provider effects.
// The caller releases the store lock, waits its normal cadence and rechecks
// current authority/identity on the next iteration of the existing follow.
type followObservation struct{ timeouts int }

func (f *followObservation) failure(d Dependencies, id string, deadline time.Time, err error) (error, bool) {
	var call *HerdrCallError
	if !errors.As(err, &call) || !call.ReadOnly || call.Phase != "agent get" || call.Class != "timeout" {
		return err, false
	}
	if call.ObserverDeadline {
		return workflowObservationWindowExpired, false
	}
	f.timeouts++
	resume := f.timeouts == 1 && time.Until(deadline) > workflowFollowInterval
	limit := d.HerdrTimeout
	if limit <= 0 {
		limit = workflowHerdrReadTimeout
	}
	logrus.WithFields(logrus.Fields{"run_id": id, "source": "Herdr agent get", "class": call.Class, "request_timeout": limit, "remaining": time.Until(deadline), "read_continuation": resume}).Warn("Follow observation request timed out; provider state unknown")
	if f.timeouts > 1 {
		err = workflowError(5, "Transport observation unavailable after bounded read continuation: Herdr agent get class=timeout; provider state unknown; no provider input was sent")
	}
	return err, resume
}

func workflowFollowReadback(d Dependencies, root, id string, err error) (WorkflowRun, error) {
	o, showErr := WorkflowShow(d, root, id)
	if showErr != nil {
		return o, showErr
	}
	o.NextAction = WorkflowAction{"coordinator", "Continue observation with ply workflow run follow " + id + " --timeout 45 --format json; this reads the same Run/session without sending provider input."}
	return o, err
}
