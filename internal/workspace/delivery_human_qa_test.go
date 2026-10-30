package workspace

import "testing"

func TestLatestTaskHumanQAUsesActualTimeAndRejectsAmbiguousLatestAnswers(t *testing.T) {
	qa := func(id, at, outcome string) TaskHumanQARecord {
		return TaskHumanQARecord{ID: HumanQARecordID(id), TaskResultID: "result", ResultOID: "oid", ResultTree: "tree", Outcome: outcome, Actor: TaskHumanActorRecord{CompletedAtUTC: at}}
	}
	older := qa("ffff", "2026-10-01T12:00:00Z", "pass")
	newer := qa("0000", "2026-10-01T12:00:00.1Z", "fail")
	equalPass := qa("ffff", newer.Actor.CompletedAtUTC, "pass")
	equalFail := qa("aaaa", newer.Actor.CompletedAtUTC, "fail")
	later := qa("0001", "2026-10-01T12:00:01Z", "pass")
	unrelated := qa("other", "2026-10-01T13:00:00Z", "fail")
	unrelated.ResultOID = "other"
	for _, tc := range []struct {
		name    string
		records []TaskHumanQARecord
		want    HumanQARecordID
		wantErr bool
	}{
		{"fractional time beats ID order", []TaskHumanQARecord{newer, older}, newer.ID, false},
		{"opposite input order", []TaskHumanQARecord{older, newer}, newer.ID, false},
		{"equal time conflict", []TaskHumanQARecord{newer, equalPass}, "", true},
		{"equal time conflict reversed", []TaskHumanQARecord{equalPass, newer}, "", true},
		{"same outcome stable ID", []TaskHumanQARecord{newer, equalFail}, equalFail.ID, false},
		{"same outcome reversed", []TaskHumanQARecord{equalFail, newer}, equalFail.ID, false},
		{"old ambiguity superseded", []TaskHumanQARecord{newer, equalPass, later}, later.ID, false},
		{"different candidate ignored", []TaskHumanQARecord{later, unrelated}, later.ID, false},
		{"no matching answer", []TaskHumanQARecord{unrelated}, "", false},
		{"invalid time", []TaskHumanQARecord{qa("bad", "invalid", "pass")}, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LatestTaskHumanQA(tc.records, "result", "oid", "tree")
			if (err != nil) != tc.wantErr {
				t.Fatalf("selection error = %v; want error %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}
			if tc.want == "" {
				if got != nil {
					t.Fatalf("unexpected answer: %#v", got)
				}
			} else if got == nil || got.ID != tc.want {
				t.Fatalf("answer = %#v; want %s", got, tc.want)
			}
		})
	}
}
