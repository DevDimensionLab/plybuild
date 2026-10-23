package workspace

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Compatibility is a finite read-only native projection. These facts confer no
// foreign-owner authority and are never accepted as caller JSON by the CLI.
type ParallelConsumerCompatibility struct {
	ID               string   `yaml:"id" json:"id"`
	Class            string   `yaml:"class" json:"class"`
	RemainingEffects []string `yaml:"remaining_effects" json:"remaining_effects"`
	ControlSHA256    string   `yaml:"control_sha256" json:"control_sha256"`
	EvidenceSHA256   string   `yaml:"evidence_sha256" json:"evidence_sha256"`
	Compatible       bool     `yaml:"compatible" json:"compatible"`
	Reasons          []string `yaml:"reasons" json:"reasons"`
}
type ParallelActivationEvidence struct {
	Target                  QueueTarget                     `yaml:"target" json:"target"`
	ReaderVersion           string                          `yaml:"reader_version" json:"reader_version"`
	ReaderCommit            string                          `yaml:"reader_commit" json:"reader_commit"`
	ReaderSHA256            string                          `yaml:"reader_sha256" json:"reader_sha256"`
	QualifiedManifestSHA256 string                          `yaml:"qualified_manifest_sha256" json:"qualified_manifest_sha256"`
	InventorySHA256         string                          `yaml:"inventory_sha256" json:"inventory_sha256"`
	Consumers               []ParallelConsumerCompatibility `yaml:"consumers" json:"consumers"`
}
type ParallelActivationPlan struct {
	Kind           string                     `yaml:"kind" json:"kind"`
	SchemaVersion  int                        `yaml:"schema_version" json:"schema_version"`
	Workspace      string                     `yaml:"workspace" json:"workspace"`
	FromVersion    int                        `yaml:"from_version" json:"from_version"`
	RegistrySHA256 string                     `yaml:"registry_sha256" json:"registry_sha256"`
	Evidence       ParallelActivationEvidence `yaml:"evidence" json:"evidence"`
}
type ParallelActivationReadback struct {
	State        string                 `json:"state"`
	Plan         ParallelActivationPlan `json:"plan"`
	Confirmation string                 `json:"confirmation"`
	Blockers     []string               `json:"blockers"`
	Backup       string                 `json:"backup"`
}
type parallelActivationPermit struct{ Plan ParallelActivationPlan }

func ActivateParallelRegistry(d Dependencies, confirmation string, apply bool, observe func(WorkItemRegistry) (ParallelActivationEvidence, error)) (ParallelActivationReadback, error) {
	out := ParallelActivationReadback{Blockers: []string{}}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	if observe == nil {
		return out, fmt.Errorf("native compatibility observer is required")
	}
	inspect := func(r WorkItemRegistry) error {
		if r.FormatVersion == ParallelRegistryFormat {
			if r.ParallelActivation == nil {
				return fmt.Errorf("parallel activation has no native publication receipt")
			}
			out.Plan = *r.ParallelActivation
			out.Confirmation = queueDigest(out.Plan)
			out.State = "activated"
			out.Backup = taskContentPath(root, "backups", out.Plan.RegistrySHA256)
			if apply && confirmation != out.Confirmation {
				return registryDiagnostic("registry_upgrade_conflict", "confirmation belongs to another activation", r.FormatVersion)
			}
			return nil
		}
		if r.FormatVersion < 4 || r.FormatVersion > 7 {
			return registryDiagnostic("parallel_activation_predecessor", "first use the existing supported legacy content/queue upgrade", r.FormatVersion)
		}
		evidence, e := observe(r)
		if e != nil {
			return e
		}
		if validateQueueTarget(r, evidence.Target) != nil || !validOIDText(evidence.ReaderCommit) || evidence.ReaderVersion == "" || !digestPattern.MatchString(evidence.ReaderSHA256) || !digestPattern.MatchString(evidence.QualifiedManifestSHA256) || !digestPattern.MatchString(evidence.InventorySHA256) || evidence.Consumers == nil {
			return fmt.Errorf("activation requires exact native qualified reader and complete compatibility inventory")
		}
		previous := ""
		for _, consumer := range evidence.Consumers {
			if consumer.ID <= previous || consumer.RemainingEffects == nil || consumer.Reasons == nil || !digestPattern.MatchString(consumer.EvidenceSHA256) {
				return fmt.Errorf("invalid native consumer inventory")
			}
			previous = consumer.ID
			if !consumer.Compatible {
				out.Blockers = append(out.Blockers, consumer.ID+": "+strings.Join(consumer.Reasons, "; "))
			}
		}
		out.Plan = ParallelActivationPlan{Kind: "WorkspaceParallelActivationPlan@1", SchemaVersion: 1, Workspace: root, FromVersion: r.FormatVersion, RegistrySHA256: r.RawSHA256, Evidence: evidence}
		out.Confirmation = queueDigest(out.Plan)
		out.State = "preview"
		out.Backup = taskContentPath(root, "backups", r.RawSHA256)
		if len(out.Blockers) > 0 {
			out.State = "blocked"
		}
		return nil
	}
	if !apply {
		err = WithTaskSpecRegistrationSnapshot(d, root, func(s *TaskSpecSession) error { return inspect(s.Registry) })
		return out, err
	}
	err = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			if e = inspect(r); e != nil {
				return e
			}
			if out.State == "activated" {
				return nil
			}
			if confirmation != out.Confirmation {
				return registryDiagnostic("registry_upgrade_conflict", "registry, reader or native consumer facts changed since preview", r.FormatVersion)
			}
			if len(out.Blockers) > 0 {
				return registryDiagnostic("activation_consumer_compatibility_unknown", strings.Join(out.Blockers, " | "), r.FormatVersion)
			}
			raw, e := os.ReadFile(workItemsPath(root))
			if e != nil || digestTaskBytes(raw) != r.RawSHA256 {
				return fmt.Errorf("activation registry bytes changed")
			}
			if _, e = d.TaskContent.Publish(root, "backups", raw); e != nil {
				return e
			}
			plan := out.Plan
			r.FormatVersion = ParallelRegistryFormat
			r.ParallelActivation = &plan
			r.SourceReconciliations = []TaskSourceReconciliation{}
			proof := plan.Evidence
			r.ParallelBoundary = &RegistryFormatBoundary{FromVersion: plan.FromVersion, FromSHA256: plan.RegistrySHA256, ManifestSHA256: proof.QualifiedManifestSHA256, ReaderVersion: proof.ReaderVersion, ReaderCommit: proof.ReaderCommit, ReaderSHA256: proof.ReaderSHA256}
			r.parallelPermit = &parallelActivationPermit{Plan: plan}
			if _, e = publishWorkItemRegistryRecover(session, r); e != nil {
				return e
			}
			out.State = "activated"
			return nil
		})
	})
	return out, err
}

func validateParallelRegistryTransition(old, next WorkItemRegistry) error {
	if old.FormatVersion == ParallelRegistryFormat {
		if !reflect.DeepEqual(old.ParallelBoundary, next.ParallelBoundary) || !reflect.DeepEqual(old.ParallelActivation, next.ParallelActivation) || !reflect.DeepEqual(old.FormatBoundary, next.FormatBoundary) || !reflect.DeepEqual(old.EpicCreateBoundary, next.EpicCreateBoundary) || !reflect.DeepEqual(old.AdoptionBoundary, next.AdoptionBoundary) {
			return fmt.Errorf("immutable format boundaries changed")
		}
		return nil
	}
	permit := next.parallelPermit
	if permit == nil || old.FormatVersion < 4 || old.FormatVersion > 7 || permit.Plan.RegistrySHA256 != old.RawSHA256 || permit.Plan.FromVersion != old.FormatVersion || next.ParallelActivation == nil || !contentTypedEqual(*next.ParallelActivation, permit.Plan) || next.ParallelBoundary == nil || next.ParallelBoundary.FromSHA256 != old.RawSHA256 {
		return registryDiagnostic("parallel_activation_required", "use the explicit qualified parallel activation plan", old.FormatVersion)
	}
	copy := next
	copy.FormatVersion = old.FormatVersion
	copy.ParallelBoundary = nil
	copy.ParallelActivation = nil
	copy.SourceReconciliations = nil
	copy.parallelPermit = nil
	before, e := encodeWorkItemRegistry(old)
	if e != nil {
		return e
	}
	after, e := encodeWorkItemRegistry(copy)
	if e != nil || !bytes.Equal(before, after) {
		return fmt.Errorf("parallel activation must preserve every predecessor record")
	}
	return nil
}
