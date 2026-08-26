package config

import (
	"path/filepath"
	"testing"
)

func newMockCloudConfig() (cfg GitCloudConfig) {
	cfg.Impl.Path = "test/cloud-config"
	return
}

func TestGitCloudConfig_Services(t *testing.T) {
	cfg := newMockCloudConfig()

	services, err := cfg.Services()()
	if err != nil {
		t.Fatalf("load services fixture: %v", err)
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
		t.Fatalf("resolve service link: %v", err)
	}

	expected := "http://localhost:8080/actuator/info"
	if link != expected {
		t.Errorf("expected link %s, got %s\n", expected, link)
	}
}

func TestGitCloudConfig_DefaultServiceEnvironmentUrl(t *testing.T) {
	cfg := newMockCloudConfig()

	services, err := cfg.Services()()
	if err != nil {
		t.Fatalf("load services fixture: %v", err)
	}
	key := "info"
	defaultUrl, err := cfg.DefaultServiceEnvironmentUrl(services.Data[0], key)
	if err != nil {
		t.Fatalf("resolve default environment URL: %v", err)
	}

	expected := "http://localhost:8080/actuator/info"
	if defaultUrl != expected {
		t.Errorf("expected default-url %s, got %s\n", expected, defaultUrl)
	}
}

func TestGitCloudConfigValidTemplatesFromReturnsExactOrderedPartialResultWithError(t *testing.T) {
	isolateConfigTestHome(t)
	cfg := newMockCloudConfig()

	templates, err := cfg.ValidTemplatesFrom([]string{"test-template", "test-template", "missing-template"})

	wantError := "could not find any valid templates with name: missing-template"
	if err == nil || err.Error() != wantError {
		t.Fatalf("ValidTemplatesFrom error was %v, want exact error %q", err, wantError)
	}
	if len(templates) != 1 {
		t.Fatalf("ValidTemplatesFrom returned %#v, want one populated partial template", templates)
	}
	if templates[0].Name != "test-template" ||
		templates[0].Project.Path != filepath.Join("test", "cloud-config", "templates", "test-template") ||
		templates[0].Project.Config.Name != "test-template" {
		t.Fatalf("ValidTemplatesFrom partial result was %#v, want exact first unique template", templates[0])
	}
}
