package taskrun

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unicode/utf8"
)

// ClaudeTrustEffect reports only configuration identity and the bounded trust
// effect. It never includes configuration contents or a native readiness claim.
type ClaudeTrustEffect struct {
	State        string `json:"state"`
	ConfigPath   string `json:"config_path"`
	ProjectKey   string `json:"project_key"`
	BeforeSHA256 string `json:"before_sha256,omitempty"`
	AfterSHA256  string `json:"after_sha256,omitempty"`
	BackupPath   string `json:"backup_path,omitempty"`
	Reason       string `json:"reason,omitempty"`
}

func workflowApplyClaudeTrust(d Dependencies, configPath, projectKey string) (effect ClaudeTrustEffect) {
	effect = ClaudeTrustEffect{State: "failed", ConfigPath: configPath, ProjectKey: projectKey}
	if !filepath.IsAbs(configPath) || filepath.Clean(configPath) != configPath || !plain(projectKey, 1, 4096) || !filepath.IsAbs(projectKey) || filepath.Clean(projectKey) != projectKey {
		effect.Reason = "configuration path or project key is not absolute and clean"
		return effect
	}
	root, err := physicalRoot(filepath.Dir(configPath), false)
	if err != nil {
		if os.IsNotExist(err) {
			effect.State, effect.Reason = "skipped", "configuration file is absent; native onboarding is unchanged"
		} else {
			effect.Reason = "configuration parent is not an accessible physical directory"
		}
		return effect
	}
	defer root.Close()
	name := filepath.Base(configPath)
	before, err := workflowReadClaudeConfig(root, name)
	if err != nil {
		if os.IsNotExist(err) {
			effect.State, effect.Reason = "skipped", "configuration file is absent; native onboarding is unchanged"
		} else {
			effect.Reason = "configuration is not an accessible physical regular file"
		}
		return effect
	}
	effect.BeforeSHA256, effect.AfterSHA256 = hash(before.data), hash(before.data)
	_, already, err := workflowMergeClaudeTrust(before.data, projectKey)
	if err != nil {
		effect.Reason = "configuration JSON or trust structure is unsupported"
		return effect
	}
	if already {
		faultErr := d.fault("claude_trust_before_already")
		effect = workflowConfirmClaudeAlready(root, effect, before)
		if faultErr != nil {
			effect.State = "failed"
			effect.Reason = "existing configuration trust could not be confirmed"
		}
		return effect
	}
	// The initial read keeps missing, invalid and already-trusted files free of
	// side effects. Re-read under this lock so simultaneous Ply launches merge
	// into the latest file instead of replacing each other's project entries.
	lock, err := workflowLockClaudeConfig(root, name+".lock")
	if err != nil {
		effect.AfterSHA256 = ""
		if actual, readErr := workflowReadClaudeConfig(root, name); readErr == nil && workflowClaudeParentCurrent(root, configPath) {
			effect.AfterSHA256 = hash(actual.data)
		}
		effect.Reason = "configuration lock is unavailable"
		return effect
	}
	defer func() {
		if lock.release() != nil {
			if effect.State == "written" {
				effect.State, effect.Reason = "unknown", "configuration was replaced but lock release is uncertain; no rollback attempted"
			} else if effect.Reason == "" {
				effect.Reason = "configuration lock release could not be confirmed"
			}
		}
	}()
	effect = workflowWriteClaudeTrust(d, root, lock, effect)
	return effect
}

const workflowClaudeConfigLimit = 8 << 20

type workflowClaudeConfigSnapshot struct {
	data []byte
	info os.FileInfo
}

type workflowClaudeConfigLock struct {
	root     *os.Root
	name     string
	info     os.FileInfo
	acquired time.Time
}

// Static binding: Claude Code 2.1.285, /opt/homebrew/Caskroom/claude-code/2.1.285/claude,
// SHA-256 51f09bd1e021d9fa8a1864c179799bd37cb39962a937935c5cf6823398e86db4,
// ei at byte 178148756 uses proper-lockfile's `${configPath}.lock` directory.
// Its normal writer reads after locking; native fallback and exit writers can
// still write unlocked. Ply never takes over an existing lock or assumes full
// native exclusion. A two-second lease stays below native's ten-second stale
// threshold without duplicating its heartbeat or stale-recovery machinery.
func workflowLockClaudeConfig(root *os.Root, name string) (*workflowClaudeConfigLock, error) {
	deadline := time.Now().Add(2 * time.Second)
	for {
		acquired := time.Now()
		err := root.Mkdir(name, 0700)
		if err == nil {
			info, statErr := root.Lstat(name)
			if statErr != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return nil, errors.New("configuration lock identity unavailable")
			}
			return &workflowClaudeConfigLock{root, name, info, acquired}, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}
		info, err := root.Lstat(name)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
			return nil, errors.New("configuration lock is not a physical directory")
		}
		if !time.Now().Before(deadline) {
			return nil, errors.New("configuration lock is busy")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (lock *workflowClaudeConfigLock) owned() bool {
	current, err := lock.root.Lstat(lock.name)
	return err == nil && current.IsDir() && current.Mode() == lock.info.Mode() && os.SameFile(current, lock.info) && current.ModTime().Equal(lock.info.ModTime())
}

func (lock *workflowClaudeConfigLock) fresh() bool {
	now := time.Now()
	return lock.owned() && now.Sub(lock.acquired) < 2*time.Second && now.Sub(lock.info.ModTime()) < 2*time.Second && !lock.info.ModTime().After(now.Add(time.Second))
}

func (lock *workflowClaudeConfigLock) release() error {
	if !lock.owned() {
		return errors.New("configuration lock ownership changed")
	}
	return lock.root.Remove(lock.name)
}

func workflowReadClaudeConfig(root *os.Root, name string) (workflowClaudeConfigSnapshot, error) {
	var result workflowClaudeConfigSnapshot
	before, err := root.Lstat(name)
	if err != nil {
		return result, err
	}
	if !before.Mode().IsRegular() || before.Size() > workflowClaudeConfigLimit {
		return result, errors.New("unsupported configuration file")
	}
	f, err := workflowOpenClaudeConfig(root, name)
	if err != nil {
		return result, err
	}
	defer f.Close()
	actual, err := f.Stat()
	if err != nil || !workflowSameClaudeFile(before, actual) {
		return result, errors.New("configuration changed while opening")
	}
	b, err := io.ReadAll(io.LimitReader(f, workflowClaudeConfigLimit+1))
	if err != nil {
		return result, err
	}
	after, err := root.Lstat(name)
	if err != nil || !workflowSameClaudeFile(actual, after) || len(b) > workflowClaudeConfigLimit {
		return result, errors.New("configuration changed while reading")
	}
	return workflowClaudeConfigSnapshot{b, after}, nil
}

func workflowSameClaudeFile(a, b os.FileInfo) bool {
	return a != nil && b != nil && a.Mode().IsRegular() && b.Mode().IsRegular() && os.SameFile(a, b) && a.Size() == b.Size() && a.Mode() == b.Mode() && a.ModTime().Equal(b.ModTime())
}

func workflowClaudeParentCurrent(root *os.Root, path string) bool {
	parent := filepath.Dir(path)
	if physical(parent, false) != nil {
		return false
	}
	current, err := os.Lstat(parent)
	if err != nil {
		return false
	}
	bound, err := root.Stat(".")
	return err == nil && current.IsDir() && os.SameFile(current, bound)
}

// A no-op still asserts trust at the named configuration path. Check the
// pinned parent and exact snapshot again without acquiring a lock or writing.
func workflowConfirmClaudeAlready(root *os.Root, effect ClaudeTrustEffect, before workflowClaudeConfigSnapshot) ClaudeTrustEffect {
	effect.State, effect.AfterSHA256 = "failed", ""
	effect.Reason = "configuration parent or snapshot changed before trust confirmation"
	if !workflowClaudeParentCurrent(root, effect.ConfigPath) {
		return effect
	}
	current, err := workflowReadClaudeConfig(root, filepath.Base(effect.ConfigPath))
	if err != nil || !workflowClaudeParentCurrent(root, effect.ConfigPath) {
		return effect
	}
	effect.AfterSHA256 = hash(current.data)
	if workflowSameClaudeFile(before.info, current.info) && bytes.Equal(before.data, current.data) {
		effect.State, effect.Reason = "already", ""
	}
	return effect
}

func workflowWriteClaudeTrust(d Dependencies, root *os.Root, lock *workflowClaudeConfigLock, effect ClaudeTrustEffect) (out ClaudeTrustEffect) {
	out = effect
	name := filepath.Base(out.ConfigPath)
	before, err := workflowReadClaudeConfig(root, name)
	if err != nil {
		out.BeforeSHA256, out.AfterSHA256 = "", ""
		out.Reason = "configuration could not be re-read under the lock"
		return out
	}
	out.BeforeSHA256 = hash(before.data)
	var intended []byte
	replaceAttempted := false
	// Inspect while still holding the lock. A native writer can ignore this
	// lock; freshness checks detect observed drift but are not a native CAS.
	defer func() {
		if out.State == "already" {
			out = workflowConfirmClaudeAlready(root, out, before)
			return
		}
		out.AfterSHA256 = ""
		after, readErr := workflowReadClaudeConfig(root, name)
		if readErr == nil && workflowClaudeParentCurrent(root, out.ConfigPath) {
			out.AfterSHA256 = hash(after.data)
			if replaceAttempted && out.State == "written" && !bytes.Equal(after.data, intended) {
				out.State, out.Reason = "unknown", "configuration changed after publication; no rollback attempted"
			} else if replaceAttempted && out.State == "failed" && (!workflowSameClaudeFile(before.info, after.info) || !bytes.Equal(before.data, after.data)) {
				out.State, out.Reason = "unknown", "configuration replacement outcome is uncertain; no rollback attempted"
			}
		} else if replaceAttempted {
			out.State, out.Reason = "unknown", "configuration poststate is unavailable; no rollback attempted"
		}
	}()
	var already bool
	intended, already, err = workflowMergeClaudeTrust(before.data, out.ProjectKey)
	if err != nil {
		out.Reason = "configuration JSON or trust structure cannot be updated safely"
		return out
	}
	if already {
		if d.fault("claude_trust_before_already") != nil {
			out.Reason = "existing configuration trust could not be confirmed"
			return out
		}
		out.State, out.Reason = "already", ""
		return out
	}
	if !workflowClaudeParentCurrent(root, out.ConfigPath) {
		out.Reason = "configuration parent changed before the write"
		return out
	}
	if !lock.fresh() {
		out.Reason = "configuration lock ownership or freshness changed before the write"
		return out
	}
	if d.fault("claude_trust_before_backup") != nil {
		out.Reason = "configuration backup was not created"
		return out
	}
	nonce := make([]byte, 16)
	if _, err = rand.Read(nonce); err != nil {
		out.Reason = "private configuration filenames could not be generated"
		return out
	}
	suffix := hex.EncodeToString(nonce)
	backup := name + ".ply-trust-backup-" + suffix
	if _, err = workflowCreateClaudeFile(root, backup, before.data, 0600); err != nil {
		out.Reason = "configuration backup could not be preserved"
		return out
	}
	out.BackupPath = filepath.Join(filepath.Dir(out.ConfigPath), backup)
	if workflowSyncClaudeDirectory(root) != nil {
		out.Reason = "configuration backup durability could not be confirmed"
		return out
	}
	if d.fault("claude_trust_before_temp") != nil {
		out.Reason = "configuration replacement was not prepared"
		return out
	}
	temp := name + ".ply-trust-tmp-" + suffix
	tempInfo, err := workflowCreateClaudeFile(root, temp, intended, before.info.Mode().Perm())
	if err != nil {
		out.Reason = "configuration replacement could not be prepared"
		return out
	}
	defer workflowRemoveOwnedClaudeFile(root, temp, tempInfo)
	if d.fault("claude_trust_before_publish") != nil {
		out.Reason = "configuration replacement was not published"
		return out
	}
	fresh, err := workflowReadClaudeConfig(root, name)
	if err != nil || !workflowSameClaudeFile(before.info, fresh.info) || !bytes.Equal(before.data, fresh.data) || !workflowClaudeParentCurrent(root, out.ConfigPath) {
		out.Reason = "configuration changed before publication; observed changes were preserved"
		return out
	}
	if !lock.fresh() {
		out.Reason = "configuration lock ownership or freshness changed before publication"
		return out
	}
	replaceAttempted = true
	if root.Rename(temp, name) != nil {
		out.Reason = "configuration replacement failed"
		return out
	}
	out.State, out.Reason = "unknown", "configuration was replaced but completion is uncertain; no rollback attempted"
	if d.fault("claude_trust_after_publish") != nil || d.fault("claude_trust_before_directory_sync") != nil || workflowSyncClaudeDirectory(root) != nil {
		return out
	}
	out.State, out.Reason = "written", ""
	return out
}

// Create privately before restoring the original configuration mode. Only
// completed, synced bytes may be renamed over the configuration. Backups remain
// 0600 even when the original configuration is readable by another user.
func workflowCreateClaudeFile(root *os.Root, name string, b []byte, mode os.FileMode) (info os.FileInfo, err error) {
	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return nil, err
	}
	owned, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	defer func() {
		if err != nil {
			workflowRemoveOwnedClaudeFile(root, name, owned)
		}
	}()
	if _, err = f.Write(b); err == nil {
		err = f.Chmod(mode)
	}
	if err == nil {
		err = f.Sync()
	}
	if err == nil {
		info, err = f.Stat()
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	return info, err
}

func workflowRemoveOwnedClaudeFile(root *os.Root, name string, owned os.FileInfo) {
	actual, err := root.Lstat(name)
	if err == nil && owned != nil && actual.Mode().IsRegular() && os.SameFile(owned, actual) {
		_ = root.Remove(name)
	}
}

func workflowSyncClaudeDirectory(root *os.Root) error {
	f, err := root.Open(".")
	if err != nil {
		return err
	}
	defer f.Close()
	return f.Sync()
}

func workflowMergeClaudeTrust(b []byte, key string) ([]byte, bool, error) {
	if !utf8.Valid(b) || !json.Valid(b) || !workflowClaudeUnicodeValid(b) {
		return nil, false, errors.New("invalid JSON")
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	if err := workflowCheckClaudeJSON(dec, 0); err != nil {
		return nil, false, err
	}
	object := func(raw []byte) (map[string]json.RawMessage, error) {
		var v map[string]json.RawMessage
		if err := json.Unmarshal(raw, &v); err != nil || v == nil {
			return nil, errors.New("expected JSON object")
		}
		return v, nil
	}
	config, err := object(b)
	if err != nil {
		return nil, false, err
	}
	projects := map[string]json.RawMessage{}
	if raw, ok := config["projects"]; ok {
		projects, err = object(raw)
		if err != nil {
			return nil, false, err
		}
	}
	project := map[string]json.RawMessage{}
	if raw, ok := projects[key]; ok {
		project, err = object(raw)
		if err != nil {
			return nil, false, err
		}
	}
	if raw, ok := project["hasTrustDialogAccepted"]; ok {
		switch string(bytes.TrimSpace(raw)) {
		case "true":
			return b, true, nil
		case "false":
		default:
			return nil, false, errors.New("expected trust boolean")
		}
	}
	project["hasTrustDialogAccepted"] = json.RawMessage("true")
	projects[key], err = json.Marshal(project)
	if err != nil {
		return nil, false, err
	}
	config["projects"], err = json.Marshal(projects)
	if err != nil {
		return nil, false, err
	}
	// Compact output avoids unbounded indentation amplification for nested
	// unknown values. The replacement must fit the same readback bound.
	after, err := json.Marshal(config)
	if len(after)+1 > workflowClaudeConfigLimit {
		return nil, false, errors.New("configuration replacement exceeds size limit")
	}
	return append(after, '\n'), false, err
}

// encoding/json replaces unpaired UTF-16 surrogates with U+FFFD when decoding
// object keys. Reject them instead of silently renaming unrelated config keys.
// JSON syntax has already been validated, so every escape has its full width.
func workflowClaudeUnicodeValid(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] != '\\' {
			continue
		}
		i++
		if b[i] != 'u' {
			continue
		}
		value, _ := strconv.ParseUint(string(b[i+1:i+5]), 16, 16)
		i += 4
		if value >= 0xdc00 && value <= 0xdfff {
			return false
		}
		if value >= 0xd800 && value <= 0xdbff {
			if i+6 >= len(b) || b[i+1] != '\\' || b[i+2] != 'u' {
				return false
			}
			low, err := strconv.ParseUint(string(b[i+3:i+7]), 16, 16)
			if err != nil || low < 0xdc00 || low > 0xdfff {
				return false
			}
			i += 6
		}
	}
	return true
}

// Token inspection rejects duplicate keys at every depth without interpreting
// unknown numbers as float64 or imposing Ply's integer-only JSON contract.
func workflowCheckClaudeJSON(dec *json.Decoder, depth int) error {
	if depth > 256 {
		return errors.New("configuration nesting exceeds limit")
	}
	token, err := dec.Token()
	if err != nil {
		return err
	}
	delim, nested := token.(json.Delim)
	if !nested {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for dec.More() {
			token, err = dec.Token()
			if err != nil {
				return err
			}
			key, ok := token.(string)
			if !ok || seen[key] {
				return errors.New("duplicate or invalid JSON object key")
			}
			seen[key] = true
			if err = workflowCheckClaudeJSON(dec, depth+1); err != nil {
				return err
			}
		}
	case '[':
		for dec.More() {
			if err = workflowCheckClaudeJSON(dec, depth+1); err != nil {
				return err
			}
		}
	default:
		return errors.New("invalid JSON delimiter")
	}
	_, err = dec.Token()
	return err
}
