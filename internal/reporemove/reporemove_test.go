package reporemove

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

const (
	ws     = "/ws"
	digest = "sha256:plan"
)

type fake struct {
	target  Target
	plan    Plan
	errs    map[string]error
	changed map[string]bool
	calls   []string
	got     []Resolved
}

func newFake() *fake {
	return &fake{
		target:  Target{ProjectRoot: "/ws/ply", Remote: "upstream", Path: "ply/map/main", WrapperPresent: true},
		plan:    Plan{Digest: digest, Verdict: VerdictUnborn},
		errs:    map[string]error{},
		changed: map[string]bool{"metadata": true, "define": true, "wrapper": true},
	}
}

func (f *fake) deps() Dependencies {
	effect := func(name string) func(context.Context, Resolved, Plan) (bool, error) {
		return func(_ context.Context, r Resolved, p Plan) (bool, error) {
			f.calls = append(f.calls, name)
			f.got = append(f.got, r)
			if p.Digest != f.plan.Digest {
				return false, errors.New("effect got another plan")
			}
			return f.changed[name], f.errs[name]
		}
	}
	return Dependencies{
		Locate: func(_ context.Context, req Request) (Target, error) {
			f.calls = append(f.calls, "locate")
			return f.target, f.errs["locate"]
		},
		Inspect: func(_ context.Context, r Resolved) (Plan, error) {
			f.calls = append(f.calls, "inspect")
			f.got = append(f.got, r)
			return f.plan, f.errs["inspect"]
		},
		ClearWrapper:  effect("metadata"),
		Undefine:      effect("define"),
		RemoveWrapper: effect("wrapper"),
	}
}

func preview() Request {
	return Request{Workspace: ws, Project: "ply", RepoName: "map"}
}

func apply() Request {
	req := preview()
	req.Apply, req.Expect, req.Actor = true, digest, "ply-epic"
	return req
}

func requireStop(t *testing.T, err error, step Step, kind error) {
	t.Helper()
	var stopped *Error
	if !errors.As(err, &stopped) || stopped.Step != step || !errors.Is(err, kind) {
		t.Fatalf("error %v, want %s at %s", err, kind, step)
	}
}

func TestApplyRunsStepsInOrderWithResolvedTarget(t *testing.T) {
	f := newFake()
	f.changed["metadata"] = false

	result, err := Run(context.Background(), apply(), f.deps())

	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"locate", "inspect", "metadata", "define", "wrapper"}; !reflect.DeepEqual(f.calls, want) {
		t.Fatalf("calls %v", f.calls)
	}
	want := Resolved{Workspace: ws, Project: "ply", ProjectRoot: "/ws/ply", RepoName: "map", Wrapper: "/ws/ply/map", RepoDir: "/ws/ply/map/main",
		Leaf: "main", Path: "ply/map/main", Remote: "upstream", Actor: "ply-epic"}
	if result.Resolved != want || f.got[0] != want {
		t.Fatalf("resolved %+v, inspected %+v", result.Resolved, f.got[0])
	}
	if steps := []StepResult{{StepMetadata, false}, {StepDefine, true}, {StepWrapper, true}}; !reflect.DeepEqual(result.Steps, steps) {
		t.Fatalf("steps %+v", result.Steps)
	}
}

func TestPreviewHasNoEffectAndRemoteOverrideWins(t *testing.T) {
	f := newFake()
	req := preview()
	req.Remote = "mine"

	result, err := Run(context.Background(), req, f.deps())

	if err != nil || result.Plan.Digest != digest || result.Resolved.Remote != "mine" || len(result.Steps) != 0 {
		t.Fatalf("preview %+v %v", result, err)
	}
	if want := []string{"locate", "inspect"}; !reflect.DeepEqual(f.calls, want) {
		t.Fatalf("calls %v", f.calls)
	}
	f = newFake()
	f.target.Remote = ""
	if result, _ := Run(context.Background(), preview(), f.deps()); result.Resolved.Remote != DefaultRemote {
		t.Fatalf("default remote %q", result.Resolved.Remote)
	}
}

func TestNoEffectBeforeAnyStepFails(t *testing.T) {
	for name, tc := range map[string]struct {
		req   func() Request
		setup func(*fake)
		step  Step
		kind  error
		calls []string
	}{
		"apply without expect": {req: func() Request { r := apply(); r.Expect = ""; return r }, step: StepPrepare, kind: ErrInvalid},
		"apply without actor":  {req: func() Request { r := apply(); r.Actor = ""; return r }, step: StepPrepare, kind: ErrInvalid},
		"repository path":      {req: func() Request { r := apply(); r.RepoName = "../x"; return r }, step: StepPrepare, kind: ErrInvalid},
		"relative workspace":   {req: func() Request { r := apply(); r.Workspace = "ws"; return r }, step: StepPrepare, kind: ErrInvalid},
		"locate fails": {setup: func(f *fake) { f.errs["locate"] = errors.New("unreadable") }, step: StepPrepare, kind: ErrInvalid,
			calls: []string{"locate"}},
		"unknown from a reader is no effect": {setup: func(f *fake) { f.errs["inspect"] = ErrUnknown }, step: StepCheck, kind: ErrInvalid,
			calls: []string{"locate", "inspect"}},
		"inspect conflict": {setup: func(f *fake) { f.errs["inspect"] = ErrConflict }, step: StepCheck, kind: ErrConflict,
			calls: []string{"locate", "inspect"}},
		"digest mismatch": {setup: func(f *fake) { f.plan.Digest = "sha256:changed" }, step: StepCheck, kind: ErrConflict,
			calls: []string{"locate", "inspect"}},
		"leaf of another layout": {setup: func(f *fake) { f.target.Path = "elsewhere/main" }, step: StepPrepare, kind: ErrInvalid,
			calls: []string{"locate"}},
		"leaf of several elements": {setup: func(f *fake) { f.target.Path = "ply/map/a/b" }, step: StepPrepare, kind: ErrInvalid,
			calls: []string{"locate"}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newFake()
			req := apply()
			if tc.req != nil {
				req = tc.req()
			}
			if tc.setup != nil {
				tc.setup(f)
			}

			result, err := Run(context.Background(), req, f.deps())

			requireStop(t, err, tc.step, tc.kind)
			if !reflect.DeepEqual(f.calls, tc.calls) || len(result.Steps) != 0 {
				t.Fatalf("calls %v, steps %+v", f.calls, result.Steps)
			}
		})
	}
}

func TestEffectFailuresKeepClassAndStop(t *testing.T) {
	for name, tc := range map[string]struct {
		failing string
		err     error
		step    Step
		kind    error
		calls   []string
		steps   []StepResult
	}{
		"metadata conflict": {"metadata", ErrConflict, StepMetadata, ErrConflict,
			[]string{"locate", "inspect", "metadata"}, nil},
		"define unclassified": {"define", errors.New("publish failed"), StepDefine, ErrUnknown,
			[]string{"locate", "inspect", "metadata", "define"}, []StepResult{{StepMetadata, true}}},
		"wrapper unknown": {"wrapper", ErrUnknown, StepWrapper, ErrUnknown,
			[]string{"locate", "inspect", "metadata", "define", "wrapper"}, []StepResult{{StepMetadata, true}, {StepDefine, true}}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newFake()
			f.errs[tc.failing] = tc.err

			result, err := Run(context.Background(), apply(), f.deps())

			requireStop(t, err, tc.step, tc.kind)
			if !reflect.DeepEqual(f.calls, tc.calls) || !reflect.DeepEqual(result.Steps, tc.steps) {
				t.Fatalf("calls %v, steps %+v", f.calls, result.Steps)
			}
		})
	}
}

func TestRerunStates(t *testing.T) {
	for name, tc := range map[string]struct {
		target  Target
		nothing bool
		kind    error
		reason  string
	}{
		"not defined and absent": {target: Target{ProjectRoot: "/ws/ply"}, nothing: true},
		"not defined and present": {target: Target{ProjectRoot: "/ws/ply", WrapperPresent: true}, kind: ErrConflict,
			reason: "ply repo create map"},
		"earlier quarantine": {target: Target{ProjectRoot: "/ws/ply", Path: "ply/map/main", WrapperPresent: true, Quarantines: []string{"/ws/ply/.ply-repo-remove-map-x"}},
			kind: ErrConflict, reason: "/ws/ply/.ply-repo-remove-map-x"},
		"quarantine without definition": {target: Target{ProjectRoot: "/ws/ply", Quarantines: []string{"/ws/ply/.ply-repo-remove-map-x"}},
			kind: ErrConflict, reason: "inspect it"},
	} {
		t.Run(name, func(t *testing.T) {
			f := newFake()
			f.target = tc.target

			result, err := Run(context.Background(), apply(), f.deps())

			if want := []string{"locate"}; !reflect.DeepEqual(f.calls, want) {
				t.Fatalf("calls %v", f.calls)
			}
			if tc.nothing {
				if err != nil || !result.Nothing || result.Resolved.Wrapper != "/ws/ply/map" {
					t.Fatalf("result %+v %v", result, err)
				}
				return
			}
			requireStop(t, err, StepCheck, tc.kind)
			if !strings.Contains(err.Error(), tc.reason) || result.Nothing {
				t.Fatalf("error %v, result %+v", err, result)
			}
		})
	}
}
