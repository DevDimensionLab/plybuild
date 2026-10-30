package agentwork

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestHumanFollowupIgnoresQuestionAndProcessNextActor(t *testing.T) {
	for _, phase := range []string{"needs_input", "completed", "human_next"} {
		t.Run(phase, func(t *testing.T) {
			d, h, root, repo := questionWorkspace(t)
			target := Identity{"claude", "human-chosen-session", "workspace", "tab", "pane", "terminal"}
			askQuestion(t, h, root, repo, "human-followup", &target.Session)
			before := treeSnapshot(t, root+"/.ply/human-input")
			binary, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			qualifyTestCallback(t, root, binary)
			now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
			sends := 0
			s := Service{Root: root, Now: func() time.Time { return now }, Runtime: &fakeRuntime{agents: []Observation{{Identity: target, Status: "idle", Ready: boolp(true)}}}, Native: func() NativeSnapshot {
				n := readNative(d, root)
				n.Bindings = append(n.Bindings, Binding{Identity: target, Work: Work{Phase: phase, NextActor: "human", Freshness: "fresh"}, Veto: "clear"})
				return n
			}, Send: func(got Identity, _ string) error {
				if got != target {
					t.Fatal("follow-up changed the selected transport identity")
				}
				sends++
				return nil
			}}
			f, err := s.Register(Draft{PublicationKey: "human-followup", OwnerRef: "owner", MandateRef: "mandate", Registrar: Registrar{Kind: HumanRequested, ActorRef: "human"}, TargetRef: target.Ref(), WaitKind: "status", WaitRef: target.Ref(), NextCheckAtUTC: s.now()})
			if err != nil {
				t.Fatal(err)
			}
			for i := 0; i < 2; i++ {
				f, err = s.Dispatch(f.ID)
				if err != nil {
					t.Fatal(err)
				}
			}
			if sends != 1 || f.Transport.State != "sent" {
				t.Fatalf("human follow-up vetoed by process state: sends=%d transport=%+v", sends, f.Transport)
			}
			if !reflect.DeepEqual(before, treeSnapshot(t, root+"/.ply/human-input")) {
				t.Fatal("follow-up changed the separate question history")
			}
		})
	}
}
