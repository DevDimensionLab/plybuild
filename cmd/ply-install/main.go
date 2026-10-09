// ply-install is shipped beside the prebuilt CLI, so an obsolete installed Ply
// and a machine without Go or a source checkout can still install a release.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/devdimensionlab/plybuild/internal/clirelease"
)

func run() int {
	flags := flag.NewFlagSet("ply-install", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	pack := flags.String("package", "", "absolute extracted release directory")
	digest := flags.String("sha256", "", "expected sha256: manifest digest (installed binary digest for rollback)")
	dest := flags.String("destination", "", "absolute ordinary Ply executable path; parent must exist")
	rollback := flags.Bool("rollback", false, "restore the recorded binary preimage only; never downgrade registry data")
	if err := flags.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 || *digest == "" || *dest == "" || !*rollback && *pack == "" || *rollback && *pack != "" {
		flags.Usage()
		return 2
	}
	i := clirelease.Installer{}
	var r clirelease.Receipt
	var err error
	if *rollback {
		r, err = i.Rollback(*dest, *digest)
	} else {
		r, err = i.Install(*pack, *digest, *dest)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err = json.NewEncoder(os.Stdout).Encode(r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func main() { os.Exit(run()) }
