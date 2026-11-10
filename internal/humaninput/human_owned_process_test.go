package humaninput

import "testing"

func TestHumanOwnedAnswerUseNeedsNoReceiptRitual(t *testing.T) {
	f := newFixture(t)
	q := f.ask(f.choice("human-owned/direct-use"))
	answer := f.answer(q.QuestionID, "request-direct-use", "full")
	got, err := Continue(f.d, f.root, q.QuestionID, answer.RequestID, digestOf(answer), "Used the human's answer.")
	if err != nil {
		t.Fatal(err)
	}
	if got.State != StateContinued || got.Receipt != nil || got.Continuation == nil {
		t.Fatalf("unexpected observation: %+v", got)
	}
	shown, err := Show(f.d, f.root, q.QuestionID)
	if err != nil || shown.State != StateContinued || shown.Receipt != nil {
		t.Fatalf("direct use was not preserved honestly: %+v %v", shown, err)
	}
}
