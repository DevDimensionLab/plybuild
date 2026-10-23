package taskrun

import (
	"os"
	"os/exec"
	"testing"
)

func TestDeliveryCLIBinaryBuildIgnoresFixtureCWD(t *testing.T) {
	t.Setenv("PLY_DELIVERY_TEST_BINARY", "")
	t.Chdir(t.TempDir())
	before, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	binary := deliveryCLIBinary(t)
	if output, err := exec.Command(binary, "version", "--format", "json").CombinedOutput(); err != nil {
		t.Fatalf("built CLI is unavailable: %v\n%s", err, output)
	}
	if after, err := os.Getwd(); err != nil || after != before {
		t.Fatalf("CLI build changed fixture cwd: %q -> %q, %v", before, after, err)
	}
}
