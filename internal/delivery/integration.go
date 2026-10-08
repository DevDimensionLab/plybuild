package delivery

import (
	"fmt"
	"reflect"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

// ValidateForIntegration revalidates the existing Delivery's native technical
// evidence and authority, including completed PR deliveries. It neither records
// QA nor performs a delivery effect. The human integration service owns its
// separate decision and must not reinterpret PR publication as merge authority.
func (s *Service) ValidateForIntegration(cwd, id string) (Receipt, taskrun.DeliveryAuthority, error) {
	var empty taskrun.DeliveryAuthority
	d, root, err := s.context(cwd)
	if err != nil {
		return Receipt{}, empty, err
	}
	r, err := readReceipt(root, id)
	if err != nil {
		return r, empty, err
	}
	registry, _, result, basis, err := registryCandidate(d, root, r.Manifest.TaskID, r.Manifest.TaskResult.ID)
	if err != nil {
		return r, empty, err
	}
	if digest(result) != r.Manifest.TaskResultSHA256 || !reflect.DeepEqual(basis, r.Manifest.Spec) {
		return r, empty, fmt.Errorf("registered native TaskResult or Spec changed")
	}
	if err = checkEvidence(d, result, basis); err != nil {
		return r, empty, err
	}
	if err = sourceUnchanged(d, result); err != nil {
		return r, empty, err
	}
	auth, err := taskrun.ValidateDeliveryAuthority(d, root, r.Manifest.WorkflowRunID, result, r.Manifest.Agreement)
	if err != nil {
		return r, empty, err
	}
	if auth.RequestSHA256 != r.Manifest.RequestSHA256 || auth.MandateSHA256 != r.Manifest.MandateSHA256 || auth.PreparationID != r.Manifest.PreparationID || auth.ExpectedParentOID != r.Manifest.ExpectedParentOID {
		return r, empty, fmt.Errorf("preserved native delivery authority changed")
	}
	r.HumanQA, err = taskrun.LatestDeliveryHumanQA(registry.HumanQARecords, result.ID, result.ResultOID, result.ResultTree)
	if err != nil {
		return r, empty, err
	}
	if r.HumanQA != nil {
		for _, a := range r.HumanQA.Evidence {
			if err = rehash(a.Locator, a.SHA256, a.SizeBytes); err != nil {
				return r, empty, fmt.Errorf("human QA evidence: %w", err)
			}
		}
	}
	return r, auth, nil
}
