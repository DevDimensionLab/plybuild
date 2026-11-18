package agentstatus

import (
	"os"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/testutil"
)

func TestMain(m *testing.M) { os.Exit(testutil.Isolate(m)) }
