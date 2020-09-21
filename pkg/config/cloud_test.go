package config

import "testing"

func newMockCloudConfig() (cfg GitCloudConfig) {
	cfg.Impl.Path = "test/cloud-config"
	return
}

func TestGitCloudConfig_Services(t *testing.T) {
	cfg := newMockCloudConfig()

	services, err := cfg.Services()()
	if err != nil {
		t.Errorf("%v\n", err)
	}

	expected := "services"
	if services.Type != expected {
		t.Errorf("expected services type %s, got %s\n", expected, services.Type)
	}
}

func TestGitCloudConfig_LinkFromService(t *testing.T) {
	cfg := newMockCloudConfig()

	link, err := cfg.LinkFromService(cfg.Services(), "com.example", "flyway-demo", "info")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	expected := "http://localhost:8080/actuator/info"
	if link != expected {
		t.Errorf("expected link %s, got %s\n", expected, link)
	}
}

func TestGitCloudConfig_DefaultServiceEnvironmentUrl(t *testing.T) {
	cfg := newMockCloudConfig()

	services, _ := cfg.Services()()
	key := "info"
	defaultUrl, err := cfg.DefaultServiceEnvironmentUrl(services.Data[0], key)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	expected := "http://localhost:8080/actuator/info"
	if defaultUrl != expected {
		t.Errorf("expected default-url %s, got %s\n", expected, defaultUrl)
	}
}
