package taskrun

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Roles of owner-submitted delivery inputs. Each is bound by sha256 and Ply
// keeps its own private copy, so read access for other local users grants no
// authority. Write access for them is still refused: they could choose the
// bytes before binding. Native stores, controls and authority-bearing inputs,
// such as runtime acceptance, keep the owner-only private file requirement.
const (
	inputDeliveryReport  = "delivery report"
	inputReportEvidence  = "report evidence"
	inputAcceptance      = "acceptance entrypoint"
	inputCandidateReview = "candidate review"
)

// Private roles named in their retained owner-only refusals.
const (
	privateRunState          = "native workflow run state"
	privateRuntimeAcceptance = "runtime acceptance"
	privateNotificationRoute = "notification route"
)

func readInput(path, role string, max int) ([]byte, error) {
	return readFileAs(path, role, max, ownerWritableAccess)
}

func readInputBinding(path, role string, max int) (FileBinding, []byte, error) {
	raw, e := readInput(path, role, max)
	return FileBinding{path, hash(raw)}, raw, e
}

// inputBound re-reads an ordinary input by its bound digest.
func inputBound(b FileBinding, role string, max int) ([]byte, error) {
	raw, e := readInput(b.Locator, role, max)
	if e != nil {
		return nil, e
	}
	if hash(raw) != b.SHA256 {
		return nil, workflowError(4, role+" "+b.Locator+": bytes changed since they were bound; restore the bound bytes or submit a new input")
	}
	return raw, nil
}

// warnSharedInput explains, once at submission, that a mode which used to
// stop the callback is now consumed.
func warnSharedInput(role, path string) {
	info, e := os.Lstat(path)
	if e != nil || info.Mode().Perm()&0077 == 0 {
		return
	}
	logrus.WithFields(logrus.Fields{"role": role, "path": path, "mode": fmt.Sprintf("%04o", info.Mode().Perm()), "handling": "consumed by its sha256 binding; Ply keeps its own private copy and needs no chmod"}).Warn("Input is readable by other local users")
}
