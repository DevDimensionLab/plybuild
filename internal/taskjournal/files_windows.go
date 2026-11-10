//go:build windows

package taskjournal

import "os"

func createStore(string, string) (*os.File, error) { return nil, unsupported() }

// Fail closed until descriptor-relative no-reparse-point storage is available.
func unsupported() error                       { return invalid("physical journal storage is unavailable on Windows") }
func openDir(string, bool) (*os.File, error)   { return nil, unsupported() }
func readFile(string, int) ([]byte, error)     { return nil, unsupported() }
func sourceHash(string) (string, error)        { return "", unsupported() }
func lock(*os.File) (*os.File, error)          { return nil, unsupported() }
func unlock(*os.File) error                    { return unsupported() }
func linkAt(*os.File, string, string) error    { return unsupported() }
func removeAt(*os.File, string) error          { return unsupported() }
func writeTemp(*os.File, string, []byte) error { return unsupported() }
func regularMode(string) error                 { return unsupported() }
