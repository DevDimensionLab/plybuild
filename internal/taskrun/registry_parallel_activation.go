package taskrun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/buildinfo"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"
)

// Parallel activation reuses the completed native qualified-CLI gate. It neither
// installs a CLI nor weakens closing/cleanup qualification. The retained exact
// candidate supplies the compatible coordinator, independently of old owners.
func ActivateParallelWorkspace(d Dependencies, root, qualifiedRun, coordinator, confirmation string, apply bool) (workspace.ParallelActivationReadback, error) {
	var out workspace.ParallelActivationReadback
	if err := containing(d, root); err != nil {
		return out, err
	}
	installation, binary, err := qualifiedCLIInstallation(d, root, qualifiedRun, filepath.Join(root, "ply-parallel-reader"), coordinator)
	if err != nil {
		return out, err
	}
	if hash(binary) != installation.Source.BinarySHA256 {
		return out, workflowError(4, "qualified reader bytes changed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, installation.Controller.Locator, "version", "--format", "json")
	command.Env = []string{"PATH=/usr/bin:/bin", "LC_ALL=C"}
	command.Dir = root
	raw, err := command.Output()
	if err != nil {
		return out, fmt.Errorf("qualified reader identity unavailable: %w", err)
	}
	var identity buildinfo.Identity
	if err = decode(raw, 1<<20, &identity); err != nil {
		return out, err
	}
	if identity.SHA256 != installation.Source.BinarySHA256 || identity.Dirty || identity.SourceCommit != installation.Source.OID || !slices.Contains(identity.RegistryFormats, workspace.ParallelRegistryFormat) {
		return out, workflowError(4, "exact qualified clean reader must identify its delivered commit and format8 capability")
	}
	observe := func(r workspace.WorkItemRegistry) (workspace.ParallelActivationEvidence, error) {
		s, e := workflowRead(root, qualifiedRun)
		if e != nil {
			return workspace.ParallelActivationEvidence{}, e
		}
		agreement := deliveryEffectiveAgreement(s)
		if agreement == nil {
			return workspace.ParallelActivationEvidence{}, workflowError(4, "qualified reader has no registered target")
		}
		var target workspace.QueueTarget
		for _, ep := range r.Epics {
			if ep.ID != agreement.EpicID {
				continue
			}
			for _, b := range ep.RepoBindings {
				if b.RepoID == agreement.RepoID {
					target = workspace.QueueTarget{ProjectID: ep.ProjectID, RepoID: b.RepoID, EpicID: ep.ID, GitCommonDir: b.GitCommonDir, ParentWorktreeID: b.Worktree.ID, ParentLocator: b.Worktree.Locator, ParentRef: b.Worktree.Ref}
				}
			}
		}
		consumers, inventorySHA, e := parallelActivationConsumers(d, root, r, installation.Source.BinarySHA256)
		return workspace.ParallelActivationEvidence{Target: target, ReaderVersion: identity.Version, ReaderCommit: identity.SourceCommit, ReaderSHA256: identity.SHA256, QualifiedManifestSHA256: digest(installation.Source), InventorySHA256: inventorySHA, Consumers: consumers}, e
	}
	run := func() error {
		var e error
		out, e = workspace.ActivateParallelRegistry(d.Workspace, confirmation, apply, observe)
		return e
	}
	// A preview observes existing stores without creating a lock or any directory.
	if !apply {
		err = run()
		return out, err
	}
	err = withStore(root, run)
	return out, err
}

func parallelActivationConsumers(d Dependencies, root string, r workspace.WorkItemRegistry, readerSHA string) ([]workspace.ParallelConsumerCompatibility, string, error) {
	inventory, err := ReadInventory(d, root)
	if err != nil {
		return nil, "", err
	}
	out := []workspace.ParallelConsumerCompatibility{}
	owned := map[string]bool{}
	known := map[string]bool{}
	addUnknown := func(id, reason string, evidence any) {
		out = append(out, workspace.ParallelConsumerCompatibility{ID: id, Class: "unknown", RemainingEffects: []string{"unknown_native_effect"}, EvidenceSHA256: digest(evidence), Reasons: []string{reason}})
	}
	if len(inventory.Reasons) > 0 {
		addUnknown("inventory/native", "activation_consumer_compatibility_unknown: incomplete native inventory", inventory.Reasons)
	}
	for _, row := range inventory.Runs {
		known[row.RunID] = true
		if row.Freshness != "fresh" {
			addUnknown("run/"+row.RunID, "activation_consumer_compatibility_unknown: damaged or incomplete native binding", row)
			continue
		}
		if row.Transport == nil || *row.Transport != "herdr" {
			addUnknown("run/"+row.RunID, "activation_consumer_compatibility_unknown: legacy native remaining effects require an identified compatible reader", row)
			continue
		}
		s, e := workflowRead(root, row.RunID)
		if e != nil {
			return nil, "", e
		}
		owned[s.Request.PreparationID] = true
		fact, e := parallelConsumerClass(d, s, readerSHA)
		if e != nil {
			addUnknown("run/"+row.RunID, e.Error(), row)
			continue
		}
		out = append(out, fact)
	}
	for _, prep := range r.TaskPreparations {
		active := false
		for _, event := range r.TaskQueueEvents {
			if event.QueueID != prep.Plan.QueueID {
				continue
			}
			id, _ := event.Request["preparation_id"].(string)
			if id == prep.ID {
				switch event.Kind {
				case "reserve", "reserve_parallel":
					active = true
				case "advance", "release":
					active = false
				}
			}
		}
		if active && !owned[prep.ID] {
			addUnknown("preparation/"+prep.ID, "activation_consumer_compatibility_unknown: preparation has no complete native Run/control binding", prep)
		}
	}
	// Bound-start intents without their already classified Run can still perform
	// first-start effects. Week state remains exclusive; an unproven controller
	// is a specific compatibility blocker, not a process/liveness prerequisite.
	for _, directory := range []string{"bound-goal-starts", "week-runs"} {
		base := filepath.Join(root, ".ply", directory)
		entries, e := os.ReadDir(base)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return nil, "", e
		}
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".lock") {
				continue
			}
			if !entry.IsDir() {
				addUnknown(directory+"/"+entry.Name(), "activation_consumer_compatibility_unknown: unrecognized native control entry", entry.Name())
				continue
			}
			path := filepath.Join(base, entry.Name(), "intent.json")
			if directory == "week-runs" {
				addUnknown(directory+"/"+entry.Name(), "activation_consumer_compatibility_unknown: exact remaining week controller/effect compatibility is unavailable", entry.Name())
				continue
			}
			raw, e := readFile(path, 1<<20, true)
			if e != nil {
				addUnknown(directory+"/"+entry.Name(), "activation_consumer_compatibility_unknown: incomplete bound start", path)
				continue
			}
			var intent struct {
				RunID string `json:"run_id"`
			}
			if e = json.Unmarshal(raw, &intent); e != nil || !known[intent.RunID] {
				addUnknown(directory+"/"+entry.Name(), "activation_consumer_compatibility_unknown: bound first-start has no classified Run", hash(raw))
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, digest([]any{inventory, out}), nil
}

func parallelConsumerClass(d Dependencies, s workflowState, readerSHA string) (workspace.ParallelConsumerCompatibility, error) {
	fact := workspace.ParallelConsumerCompatibility{ID: "run/" + s.Result.RunID, Class: "PAR-A01", RemainingEffects: []string{"source_work", "review_verify", "integration", "closeout"}, ControlSHA256: s.Request.Runtime.PlyExecutable.SHA256, EvidenceSHA256: digest(s), Reasons: []string{}}
	runtime, err := workflowEffectiveRuntime(s)
	if err != nil {
		return fact, err
	}
	fact.ControlSHA256 = runtime.PlyExecutable.SHA256
	ownerCompatible := false
	if bytes, e := workflowBound(FileBinding{runtime.PlyExecutable.Path, runtime.PlyExecutable.SHA256}, 128<<20); e == nil {
		ownerCompatible = hash(bytes) == readerSHA
	}
	if !deliveryRun(s.Request) || s.Result.Delivery == nil {
		return fact, workflowError(4, "activation_consumer_compatibility_unknown: report-only contract still has unclassified legitimate return effects")
	}
	ds := s.Result.Delivery
	if a := ds.Attempt; a != nil && a.State != "recorded" {
		return fact, workflowError(4, "activation_consumer_compatibility_unknown: unresolved native effect "+a.ID)
	}
	if retired, e := deliveryRetirement(s); e != nil {
		return fact, e
	} else if retired != nil {
		fact.Class = "PAR-A02"
		fact.RemainingEffects = []string{}
		fact.Compatible = true
		return fact, nil
	}
	if len(ds.Candidates) > 0 {
		c := ds.Candidates[len(ds.Candidates)-1]
		release, e := deliveryRelease(s, c)
		if e != nil {
			return fact, e
		}
		integrated, observed, e := automaticIntegratedResult(s, c.TaskResult.ID)
		if e != nil {
			return fact, e
		}
		closeout, e := workflowTaskCloseout(s.Request.WorkspaceRoot, s)
		if e != nil {
			return fact, e
		}
		if ds.Phase == "completed" && observed && integrated.Completed && closeout != nil && closeout.State == "complete" && closeout.LifecycleCompleted {
			fact.Class = "PAR-A02"
			fact.RemainingEffects = []string{}
			fact.Compatible = true
			return fact, nil
		}
		var releaseRecord deliveryOwnershipRelease
		if release != nil {
			raw, err := workflowBound(*release, 1<<20)
			if err != nil {
				return fact, err
			}
			if err = decode(raw, 1<<20, &releaseRecord); err != nil {
				return fact, err
			}
		}
		if release != nil && releaseRecord.Origin == "automatic_delivery_completed" && observed && integrated.Completed && ds.Phase == "closing" {
			fact.Class = "PAR-A03"
			fact.RemainingEffects = []string{"coordinator_original_closeout"}
			fact.Compatible = readerSHA != ""
			fact.EvidenceSHA256 = digest([]any{s, release, integrated, closeout, readerSHA})
			if !fact.Compatible {
				fact.Reasons = []string{"compatible qualified coordinator follow/resume path is unavailable"}
			}
			return fact, nil
		}
		if release != nil {
			fact.Class = "PAR-A04"
			fact.RemainingEffects = []string{"human_integration", "restorable_owner_correction", "coordinator_closeout"}
		}
	}
	fact.Compatible = ownerCompatible && readerSHA != ""
	if !fact.Compatible {
		fact.Reasons = []string{"exact frozen owner/correction control lacks proven format8 compatibility; finish through its existing route or wait"}
	}
	return fact, nil
}
