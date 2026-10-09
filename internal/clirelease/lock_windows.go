//go:build windows

package clirelease

import "fmt"

func installationLock(string) (func(), error) {
	return nil, fmt.Errorf("release_platform_mismatch: this installer supports macOS arm64")
}
