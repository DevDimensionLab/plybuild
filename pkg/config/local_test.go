package config

import "testing"

func newMockLocalConfig() (cfg LocalConfig) {
	cfg.impl.Path = "test/local-config"
	return
}

func TestLocalConfig_Config(t *testing.T) {
	cfg := newMockLocalConfig()

	localCfg, err := cfg.Config()
	if err != nil {
		t.Errorf("%v\n", err)
	}

	expected := "test-source-provider"
	if localCfg.SourceProvider.Host != expected {
		t.Errorf("expected sourceProvider.Host %s, got %s\n", expected, localCfg.SourceProvider.Host)
	}
}
