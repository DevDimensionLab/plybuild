package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newMockLocalConfig() (cfg LocalConfigDir) {
	cfg.impl.Path = "test/local-config"
	return
}

func TestLocalConfig_Config(t *testing.T) {
	cfg := newMockLocalConfig()

	localCfg, err := cfg.Config()
	if err != nil {
		t.Fatalf("load local config fixture: %v", err)
	}

	expected := "git@github.com:devdimensionlab/ply-config.git"
	if localCfg.CloudConfig.Git.Url != expected {
		t.Errorf("expected cloudConfig.git.url %s, got %s\n", expected, localCfg.CloudConfig.Git.Url)
	}
}

func TestLocalConfigRewriteDropsLegacySourceProviderAndKeepsOtherFields(t *testing.T) {
	const legacyToken = "legacy-access-token"
	const nexusURL = "https://nexus.example.invalid"
	cfg := OpenLocalConfig(t.TempDir())
	content := "sourceProvider:\n  host: legacy-host\n  access_token: " + legacyToken +
		"\nnexus:\n  url: " + nexusURL + "\n"
	if err := os.WriteFile(filepath.Join(cfg.impl.Path, localConfigFileName), []byte(content), 0644); err != nil {
		t.Fatalf("write legacy local config: %v", err)
	}

	localCfg, err := cfg.Config()
	if err != nil {
		t.Fatalf("legacy sourceProvider block must not fail config load: %v", err)
	}
	if err := cfg.UpdateLocalConfig(localCfg); err != nil {
		t.Fatalf("rewrite local config: %v", err)
	}

	if localCfg.Nexus.Url != nexusURL {
		t.Errorf("nexus.url was %q, want %q", localCfg.Nexus.Url, nexusURL)
	}
	rewritten, err := os.ReadFile(cfg.FilePath())
	if err != nil {
		t.Fatalf("read rewritten local config: %v", err)
	}
	text := string(rewritten)
	if !strings.Contains(text, nexusURL) {
		t.Errorf("rewritten config lost nexus.url:\n%s", text)
	}
	if strings.Contains(text, "sourceProvider") || strings.Contains(text, legacyToken) {
		t.Errorf("rewritten config kept the legacy sourceProvider block or token:\n%s", text)
	}
}
