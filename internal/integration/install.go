package integration

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// CommandSpec contains argv, not shell source. The absolute executable and
// physical working directory are visible parts of the frozen installation plan.
type CommandSpec struct {
	Executable string   `json:"executable"`
	Arguments  []string `json:"arguments"`
	CWD        string   `json:"cwd"`
}

type InstallProfile struct {
	Kind           string      `json:"kind"`
	SchemaVersion  int         `json:"schema_version"`
	Name           string      `json:"name"`
	CandidateOID   string      `json:"candidate_oid"`
	ArtifactPath   string      `json:"artifact_path"`
	ArtifactSHA256 string      `json:"artifact_sha256"`
	Install        CommandSpec `json:"install"`
	Verify         CommandSpec `json:"verify"`
}

type InstallObservation struct {
	State          string `json:"state"`
	Reason         string `json:"reason,omitempty"`
	NextAction     string `json:"next_action"`
	ProfileSHA256  string `json:"profile_sha256"`
	CandidateOID   string `json:"candidate_oid"`
	ArtifactSHA256 string `json:"artifact_sha256,omitempty"`
	ObservedAtUTC  string `json:"observed_at_utc"`
}

type InstallAdapter interface {
	Observe(InstallProfile) (InstallObservation, error)
	Install(InstallProfile) (InstallObservation, error)
}

type LocalInstallAdapter struct {
	Run func(cwd, program string, args ...string) ([]byte, error)
	Now func() time.Time
}

var installDigestPattern = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)
var installNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)

func cleanAbsolute(path string) bool {
	return filepath.IsAbs(path) && filepath.Clean(path) == path && !strings.ContainsAny(path, "\x00\r\n")
}

func validateCommand(c CommandSpec) error {
	if !cleanAbsolute(c.Executable) || !cleanAbsolute(c.CWD) || len(c.Arguments) > 256 {
		return fmt.Errorf("install commands require an absolute executable, physical cwd and explicit argument list")
	}
	for _, arg := range c.Arguments {
		if len(arg) > 16<<10 || strings.ContainsRune(arg, '\x00') {
			return fmt.Errorf("invalid install command argument")
		}
	}
	return nil
}

func ValidateInstallProfile(p InstallProfile, candidateOID string) error {
	if p.Kind != "ply.integration.install-profile" || p.SchemaVersion != 1 || !installNamePattern.MatchString(p.Name) ||
		!mergeOIDPattern.MatchString(p.CandidateOID) || p.CandidateOID != candidateOID ||
		!cleanAbsolute(p.ArtifactPath) || !installDigestPattern.MatchString(p.ArtifactSHA256) {
		return fmt.Errorf("install profile must bind the exact candidate and installed artifact SHA-256")
	}
	if err := validateCommand(p.Install); err != nil {
		return err
	}
	return validateCommand(p.Verify)
}

func InstallProfileDigest(p InstallProfile) string {
	raw, _ := json.Marshal(p)
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func LoadInstallProfile(path string) (InstallProfile, error) {
	var p InstallProfile
	if !cleanAbsolute(path) {
		return p, fmt.Errorf("install profile path must be absolute")
	}
	physical, err := filepath.EvalSymlinks(path)
	if err != nil || physical != path {
		return p, fmt.Errorf("install profile requires its physical file path")
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > 64<<10 {
		return p, fmt.Errorf("install profile must be a regular local file smaller than 64 KiB")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}
	if _, err = canonicaljson.DecodeStrict(raw); err != nil {
		return p, err
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&p); err != nil {
		return p, err
	}
	return p, ValidateInstallProfile(p, p.CandidateOID)
}

// ValidateInstallController prevents self-removal or self-replacement. The
// service calls this with its actual executable and every retained control input
// before it lets cleanup retire source. Installing an ordinary ply executable
// is safe when the running controller has a separate preserved path.
func ValidateInstallController(p InstallProfile, controller, source string, inputs ...string) error {
	if !cleanAbsolute(controller) || !cleanAbsolute(source) {
		return fmt.Errorf("controller and cleanup source must be physical absolute paths")
	}
	physicalController, err := filepath.EvalSymlinks(controller)
	if err != nil {
		return err
	}
	physicalSource, err := filepath.EvalSymlinks(source)
	if err != nil {
		return err
	}
	inside := func(path string) bool {
		rel, err := filepath.Rel(physicalSource, path)
		return err == nil && (rel == "." || rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
	}
	if inside(physicalController) || physicalController == p.ArtifactPath {
		return fmt.Errorf("preserve the integration controller outside the cleanup source and installed artifact")
	}
	controllerInfo, err := os.Stat(physicalController)
	if err != nil {
		return err
	}
	if artifactInfo, err := os.Stat(p.ArtifactPath); err == nil && os.SameFile(controllerInfo, artifactInfo) {
		return fmt.Errorf("installed artifact aliases the running integration controller")
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, path := range inputs {
		physical, err := filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}
		if inside(physical) {
			return fmt.Errorf("preserve integration control inputs outside the cleanup source")
		}
	}
	return nil
}

func checkCommandLocation(c CommandSpec) error {
	physical, err := filepath.EvalSymlinks(c.CWD)
	if err != nil || physical != c.CWD {
		return fmt.Errorf("install working directory is missing or has changed physical identity")
	}
	info, err := os.Stat(c.CWD)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("install working directory is unavailable")
	}
	executable, err := os.Stat(c.Executable)
	if err != nil || !executable.Mode().IsRegular() || executable.Mode().Perm()&0111 == 0 {
		return fmt.Errorf("install executable is unavailable or not executable")
	}
	return nil
}

func (a *LocalInstallAdapter) run(c CommandSpec) ([]byte, error) {
	if err := checkCommandLocation(c); err != nil {
		return nil, err
	}
	if a.Run != nil {
		return a.Run(c.CWD, c.Executable, c.Arguments...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, c.Executable, c.Arguments...)
	cmd.Dir = c.CWD
	stdout, stderr := adapterBuffer{limit: 2 << 20}, adapterBuffer{limit: 8 << 10}
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.Bytes(), err
	}
	if stdout.over {
		return nil, fmt.Errorf("install observation exceeded output limit")
	}
	return stdout.Bytes(), nil
}

func installedDigest(path string) (string, error) {
	parent, err := filepath.EvalSymlinks(filepath.Dir(path))
	if err != nil || parent != filepath.Dir(path) {
		return "", fmt.Errorf("installed artifact parent changed physical identity")
	}
	before, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if !before.Mode().IsRegular() {
		return "", fmt.Errorf("installed artifact is not a regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(before, opened) {
		return "", fmt.Errorf("installed artifact changed while opening")
	}
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	after, err := os.Lstat(path)
	if err != nil || !os.SameFile(before, after) || before.Size() != after.Size() || !before.ModTime().Equal(after.ModTime()) {
		return "", fmt.Errorf("installed artifact changed during verification")
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}

func (a *LocalInstallAdapter) observation(p InstallProfile) InstallObservation {
	now := time.Now
	if a.Now != nil {
		now = a.Now
	}
	return InstallObservation{State: "effect_unknown", Reason: "effect_unknown", NextAction: "Observe the installed artifact or obtain explicit human clarification before retrying installation; preserve the Task worktree.", ProfileSHA256: InstallProfileDigest(p), CandidateOID: p.CandidateOID, ObservedAtUTC: now().UTC().Format(time.RFC3339Nano)}
}

func (a *LocalInstallAdapter) Observe(p InstallProfile) (InstallObservation, error) {
	o := a.observation(p)
	if err := ValidateInstallProfile(p, p.CandidateOID); err != nil {
		return o, err
	}
	digest, err := installedDigest(p.ArtifactPath)
	if err != nil {
		return o, err
	}
	o.ArtifactSHA256 = digest
	if digest != p.ArtifactSHA256 {
		// Wrong bytes do not prove that installation never ran or had no other
		// effect. Recovery must not interpret this as safe automatic retry.
		return o, nil
	}
	if _, err = a.run(p.Verify); err != nil {
		o.State, o.Reason, o.NextAction = "install_failed", "install_verification_failed", "The artifact digest matches but its verification failed; inspect the verification result before retrying the installation step."
		return o, err
	}
	// Verification itself must not silently replace the artifact being attested.
	after, err := installedDigest(p.ArtifactPath)
	if err != nil || after != digest {
		return o, fmt.Errorf("installed artifact changed during its verification command")
	}
	o.State, o.Reason, o.NextAction = "verified", "", "Continue the remaining closeout steps."
	return o, nil
}

func (a *LocalInstallAdapter) Install(p InstallProfile) (InstallObservation, error) {
	o := a.observation(p)
	if err := ValidateInstallProfile(p, p.CandidateOID); err != nil {
		return o, err
	}
	if err := checkCommandLocation(p.Install); err != nil {
		o.State, o.Reason, o.NextAction = "install_failed", "install_not_started", "Restore the preserved install executable and working directory, then explicitly resume the install step."
		return o, err
	}
	artifactParent, err := filepath.EvalSymlinks(filepath.Dir(p.ArtifactPath))
	if err != nil || artifactParent != filepath.Dir(p.ArtifactPath) {
		o.State, o.Reason = "install_failed", "install_not_started"
		return o, fmt.Errorf("installed artifact parent is missing or has changed physical identity")
	}
	if info, err := os.Lstat(p.ArtifactPath); err == nil && !info.Mode().IsRegular() {
		o.State, o.Reason = "install_failed", "install_not_started"
		return o, fmt.Errorf("installed artifact path is a symlink or nonregular resource")
	} else if err != nil && !os.IsNotExist(err) {
		o.State, o.Reason = "install_failed", "install_not_started"
		return o, err
	}
	_, runErr := a.run(p.Install)
	observed, observeErr := a.Observe(p)
	if observed.State == "verified" {
		return observed, nil
	}
	if runErr != nil {
		var exit *exec.ExitError
		if errors.As(runErr, &exit) && exit.ProcessState != nil && exit.ProcessState.Exited() {
			observed.State, observed.Reason, observed.NextAction = "install_failed", "install_failed", "The install command exited unsuccessfully; inspect its effects before explicitly retrying the install step."
		}
		return observed, runErr
	}
	if observeErr != nil {
		return observed, observeErr
	}
	// A successful process exit still needs the artifact and verify evidence.
	observed.State, observed.Reason, observed.NextAction = "install_failed", "installed_artifact_mismatch", "The install command finished but the selected artifact was not verified; inspect the installed output before an explicit retry."
	return observed, nil
}
