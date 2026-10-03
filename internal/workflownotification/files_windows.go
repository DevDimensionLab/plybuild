//go:build windows

package workflownotification

import "os"

func unsupported() error {
	return fail(1, "unsupported_platform", "Notification storage requires a supported Unix filesystem.")
}
func openDirectory(string) (*os.File, error)               { return nil, unsupported() }
func readFile(string, int64, bool) ([]byte, error)         { return nil, unsupported() }
func readAt(*os.File, string, int64, bool) ([]byte, error) { return nil, unsupported() }
func stateDirectory(string, bool) (*os.File, error)        { return nil, unsupported() }
func lockState(*os.File) (func(), error)                   { return nil, unsupported() }
func atomicState(*os.File, []byte) error                   { return unsupported() }

func bindRoot(*os.File, string, string) error { return unsupported() }
