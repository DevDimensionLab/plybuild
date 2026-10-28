package agentwork

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

type callbackQualification struct {
	Path             string `json:"path"`
	SHA256           string `json:"sha256"`
	QualificationRef string `json:"qualification_ref"`
}

func (s *Service) callbackExecutor() (string, error) {
	invalid := fmt.Errorf("followup_qualified_callback_unavailable")
	raw, err := taskrun.PrivateRead(filepath.Join(s.Root, "ply-callback.local.json"), 16<<10)
	if err != nil {
		return "", invalid
	}
	var q callbackQualification
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&q) != nil || d.Decode(new(any)) != io.EOF || !filepath.IsAbs(q.Path) || filepath.Clean(q.Path) != q.Path || strings.ContainsAny(q.Path, "`\r\n") || strings.ContainsFunc(q.Path, unicode.IsControl) || !reference(q.QualificationRef) {
		return "", invalid
	}
	if taskrun.PhysicalPath(q.Path) != nil {
		return "", invalid
	}
	f, err := os.Open(q.Path)
	if err != nil {
		return "", invalid
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		return "", invalid
	}
	digest := sha256.New()
	if _, err = io.Copy(digest, f); err != nil || q.SHA256 != fmt.Sprintf("sha256:%x", digest.Sum(nil)) {
		return "", invalid
	}
	return "'" + strings.ReplaceAll(q.Path, "'", "'\\''") + "'", nil
}
