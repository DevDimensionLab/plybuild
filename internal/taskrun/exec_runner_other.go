//go:build !darwin

package taskrun

func systemExecRunner() ProcessRunner { return unsupportedRunner{} }
