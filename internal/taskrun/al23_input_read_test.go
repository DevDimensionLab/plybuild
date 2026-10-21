package taskrun

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The injected boundary controls read order only. Writes use the real OS file,
// with no timing race, global hook, provider, or claim of an atomic snapshot.
func TestAL23SameInodeReadBindsConsumedBytes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "evidence.log")
	original, replacement := []byte("AAAAaaaa"), []byte("BBBBbbbb")
	if err := os.WriteFile(path, original, 0644); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := readFileAsWithRead(path, inputReportEvidence, 8, noSharedWriteAccess, func(r io.Reader) ([]byte, error) {
		first := make([]byte, 4)
		if _, err := io.ReadFull(r, first); err != nil {
			return nil, err
		}
		w, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
		_, writeErr := w.WriteAt(replacement, 0)
		closeErr := w.Close()
		if writeErr != nil {
			return nil, writeErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		rest, err := io.ReadAll(r)
		return append(first, rest...), err
	})
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(path)
	if err != nil || !os.SameFile(before, after) || before.Size() != after.Size() {
		t.Fatalf("not same inode/size: %v", err)
	}
	want := []byte("AAAAbbbb")
	if !bytes.Equal(raw, want) {
		t.Fatalf("consumed %q, want %q", raw, want)
	}
	// This is readInputBinding's exact binding rule. The existing consumer must
	// reject the final file even though inode and size match the original.
	binding := FileBinding{path, hash(raw)}
	if binding.SHA256 == hash(original) || binding.SHA256 == hash(replacement) {
		t.Fatal("mixed bytes falsely bound to a complete version")
	}
	if _, err := inputBound(binding, inputReportEvidence, 8); err == nil || !strings.Contains(err.Error(), "bytes changed since they were bound") || !strings.Contains(err.Error(), path) || !strings.Contains(err.Error(), inputReportEvidence) {
		t.Fatalf("final file accepted as mixed bound bytes: %v", err)
	}
	// Returning the actual consumed bytes succeeds, as does a fresh binding of
	// the completed writer's stable file. No human attestation is involved.
	if err := os.WriteFile(path, raw, 0644); err != nil {
		t.Fatal(err)
	}
	if got, err := inputBound(binding, inputReportEvidence, 8); err != nil || !bytes.Equal(got, want) {
		t.Fatalf("bound bytes: %q %v", got, err)
	}
	if err := os.WriteFile(path, replacement, 0644); err != nil {
		t.Fatal(err)
	}
	stable, got, err := readInputBinding(path, inputReportEvidence, 8)
	if err != nil || !bytes.Equal(got, replacement) || stable.SHA256 != hash(replacement) {
		t.Fatalf("stable binding: %+v %q %v", stable, got, err)
	}
	if _, err := inputBound(stable, inputReportEvidence, 8); err != nil {
		t.Fatal(err)
	}
	t.Logf("consumed=%q sha256=%s; final=%q sha256=%s; final rejected under consumed digest, stable resubmission passed", raw, binding.SHA256, got, stable.SHA256)
}

func TestAL23ReadInstabilityControls(t *testing.T) {
	for _, change := range []string{"replace inode", "grow beyond limit"} {
		t.Run(change, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "input")
			if err := os.WriteFile(path, []byte("12345678"), 0644); err != nil {
				t.Fatal(err)
			}
			_, err := readFileAsWithRead(path, inputAcceptance, 8, noSharedWriteAccess, func(r io.Reader) ([]byte, error) {
				if change == "replace inode" {
					next := path + ".next"
					if err := os.WriteFile(next, []byte("abcdefgh"), 0644); err != nil {
						return nil, err
					}
					if err := os.Rename(next, path); err != nil {
						return nil, err
					}
				} else {
					if err := os.WriteFile(path, []byte("123456789"), 0644); err != nil {
						return nil, err
					}
				}
				return io.ReadAll(r)
			})
			if err == nil || !strings.Contains(err.Error(), "file changed while reading") || !strings.Contains(err.Error(), path) || !strings.Contains(err.Error(), inputAcceptance) {
				t.Fatalf("instability: %v", err)
			}
			if _, err := readInput(path, inputAcceptance, 16); err != nil {
				t.Fatalf("stable retry: %v", err)
			}
		})
	}
}
