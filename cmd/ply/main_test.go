package main

import "testing"

func TestMainDelegatesToCommand(t *testing.T) {
	original := execute
	t.Cleanup(func() {
		execute = original
	})

	calls := 0
	execute = func() {
		calls++
	}

	main()

	if calls != 1 {
		t.Fatalf("command entrypoint was called %d times, want 1", calls)
	}
}
