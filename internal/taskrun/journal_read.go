package taskrun

// JournalReadback exposes the validated immutable run chain to the process
// journal. No process, provider, cache write, lock or recovery is involved.
type JournalReadback struct {
	Request Request
	Binding *Binding
	Events  []Event
	Hashes  []string
	Result  Result
}

func ReadProcessJournal(d Dependencies, root, id string) (JournalReadback, error) {
	j, err := readJournal(root, id)
	out := JournalReadback{j.Request, j.Binding, j.Events, j.Hashes, j.Result}
	if err != nil {
		return out, err
	}
	if err = validateRunEvidence(d, j); err != nil {
		return out, err
	}
	return out, nil
}
