package taskrun

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Roles of owner-submitted ordinary inputs. The caller chooses which file to
// submit; Ply binds the bytes it actually read by sha256 and validates their
// content. Reports are copied into Ply's private event store or report slot,
// and the acceptance entrypoint into the attempt snapshot it executes. Report
// evidence and the candidate review are not copied at submission: every later
// read re-reads the original path against the bound digest. Mode bits do not authenticate an owner or writer, so read access for
// others neither proves nor costs anything here, and Ply does not check the
// file's UID. Group/other write access is refused because it lets others
// change the file at any time. Native stores, controls, coordinator decisions
// and authority-bearing inputs keep the owner-only private file requirement.
const (
	inputDeliveryReport  = "delivery report"
	inputRoundReport     = "round report"
	inputReportEvidence  = "report evidence"
	inputAcceptance      = "acceptance entrypoint"
	inputCandidateReview = "candidate review"
)

// Private roles named in their retained owner-only refusals.
const (
	privateRunState          = "native workflow run state"
	privateRuntimeAcceptance = "runtime acceptance"
	privateNotificationRoute = "notification route"
	privateCoordinatorReview = "coordinator review"
)

// inputFileError keeps the actual errno (errors.Is) while naming the input.
// os.IsNotExist does not unwrap it; native callers that rely on it read
// through readFile, which returns raw errors.
type inputFileError struct {
	role, path string
	err        error
}

func (e *inputFileError) Error() string {
	return fmt.Sprintf("%s %s: %v; submit the absolute path of an existing regular file that you can read", e.role, e.path, e.err)
}
func (e *inputFileError) Unwrap() error { return e.err }

func readInput(path, role string, max int) ([]byte, error) {
	b, e := readFileAs(path, role, max, noSharedWriteAccess)
	if _, named := e.(*Error); e != nil && !named {
		return nil, &inputFileError{role, path, e}
	}
	return b, e
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
	logrus.WithFields(logrus.Fields{"role": role, "path": path, "mode": fmt.Sprintf("%04o", info.Mode().Perm()), "handling": "consumed by its sha256 binding and re-checked against that digest when read again; no chmod is needed"}).Warn("Input is readable by other local users")
}
