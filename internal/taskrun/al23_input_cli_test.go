//go:build darwin

package taskrun

import ()

// Provider/owner bootstrap is modeled by the existing fixture. All report
// consumption below runs the actual supplied candidate CLI. The deterministic
// under-read write is covered separately, not injected into this subprocess.
