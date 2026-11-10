package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

func digestTaskBytes(b []byte) string {
	d := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(d[:])
}
func validAbsoluteCleanPath(v string) bool {
	return filepath.IsAbs(v) && filepath.Clean(v) == v && !strings.ContainsRune(v, 0)
}
func validTaskText(v string, min, max int) bool {
	if !utf8.ValidString(v) || v != strings.TrimSpace(v) {
		return false
	}
	count := 0
	for _, r := range v {
		if unicode.IsControl(r) {
			return false
		}
		count++
	}
	return count >= min && count <= max
}
func validTaskUTC(v string) bool {
	_, err := time.Parse(time.RFC3339Nano, v)
	return err == nil && strings.HasSuffix(v, "Z")
}
func setString(v ...string) map[string]bool {
	m := map[string]bool{}
	for _, s := range v {
		m[s] = true
	}
	return m
}

func pathWithin(parent, path string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
