package maven

import (
	"encoding/xml"
	"errors"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLatestReleaseUP01MissingOrEmptyReleaseUsesVersionList(t *testing.T) {
	for _, test := range []struct {
		name    string
		release string
	}{
		{name: "missing"},
		{name: "self-closing", release: `<release/>`},
		{name: "empty", release: `<release></release>`},
		{name: "whitespace", release: "<release> \t\n\r </release>"},
	} {
		t.Run(test.name, func(t *testing.T) {
			hook := captureMavenPartialFailureLogs(t)
			metadata := releaseMetadataFromXML(t, test.release+`
				<latest>5.12.0</latest>
				<versions><version>5.11.4</version><version>5.12.0</version></versions>
				<lastUpdated>20261007080000</lastUpdated>`)

			assertLatestRelease(t, metadata, "5.12.0")
			for _, entry := range hook.entries {
				if entry.Level <= logrus.WarnLevel {
					t.Errorf("successful fallback produced %s diagnostic: %s", entry.Level, entry.Message)
				}
			}
		})
	}
}

func TestLatestReleaseUP02IgnoresLatestAndVersionListOrder(t *testing.T) {
	orders := [][]string{
		{"5.9.0", "5.12.0", "5.11.4", "5.13.0-SNAPSHOT", "5.14.0-RC1"},
		{"5.14.0-RC1", "5.13.0-SNAPSHOT", "5.11.4", "5.12.0", "5.9.0"},
		{"5.12.0", "5.13.0-SNAPSHOT", "5.9.0", "5.14.0-RC1", "5.11.4"},
	}
	for _, latest := range []string{"", "5.9.0", "5.13.0-SNAPSHOT"} {
		for _, versions := range orders {
			t.Run(latest+"/"+strings.Join(versions, ","), func(t *testing.T) {
				latestXML := ""
				if latest != "" {
					latestXML = "<latest>" + latest + "</latest>"
				}
				metadata := releaseMetadataFromXML(t, latestXML+releaseVersionsXML(versions...))
				assertLatestRelease(t, metadata, "5.12.0")
			})
		}
	}
}

func TestLatestReleaseUP03ExplicitReleaseHasPriority(t *testing.T) {
	for _, test := range []struct {
		name     string
		versions string
	}{
		{name: "higher-list", versions: releaseVersionsXML("5.12.0", "5.20.0")},
		{name: "missing-list"},
		{name: "invalid-unused-list", versions: releaseVersionsXML("broken")},
	} {
		t.Run(test.name, func(t *testing.T) {
			metadata := releaseMetadataFromXML(t, `<release>5.11.4</release><latest>5.20.0</latest>`+test.versions)
			assertLatestRelease(t, metadata, "5.11.4")
		})
	}
}

func TestLatestReleaseUP03NonReleaseValueUsesFallback(t *testing.T) {
	for _, version := range []string{"5.13.0-SNAPSHOT", "5.13.0-RC1", "5.13.0-m3"} {
		t.Run(version, func(t *testing.T) {
			metadata := releaseMetadataFromXML(t, "<release>"+version+"</release>"+
				releaseVersionsXML("5.11.4", "5.13.0-SNAPSHOT", "5.12.0", "5.14.0-RC1"))
			assertLatestRelease(t, metadata, "5.12.0")
		})
	}
}

func TestLatestReleaseUP03PreservesExistingReleaseSuffixes(t *testing.T) {
	for _, version := range []string{
		"5.11.4", "5.11.4.RELEASE", "5.11.4-release", "5.11.4.Final",
		"5.11.4-FINAL", "5.11.4-942da656", "5.11.4-abcdefgh",
	} {
		t.Run(version, func(t *testing.T) {
			t.Run("explicit-release", func(t *testing.T) {
				metadata := releaseMetadataFromXML(t, "<release>"+version+"</release>"+
					releaseVersionsXML("5.12.0"))
				assertLatestRelease(t, metadata, version)
			})
			t.Run("fallback", func(t *testing.T) {
				metadata := releaseMetadataFromXML(t, releaseVersionsXML("5.10.0", version,
					"5.13.0-SNAPSHOT", "5.14.0-RC1", "5.15.0-m3", "5.16.0-942da656-SNAPSHOT"))
				assertLatestRelease(t, metadata, version)
			})
		})
	}
}

func TestLatestReleaseUP04RejectsMissingSuitableRelease(t *testing.T) {
	for _, test := range []struct {
		name     string
		versions string
	}{
		{name: "missing-list"},
		{name: "empty-list", versions: `<versions/>`},
		{name: "non-releases-only", versions: releaseVersionsXML("5.13.0-SNAPSHOT", "5.14.0-RC1", "5.15.0-m3")},
	} {
		for _, latest := range []string{"", "<latest>5.20.0</latest>"} {
			t.Run(test.name+"/"+latest, func(t *testing.T) {
				metadata := releaseMetadataFromXML(t, latest+test.versions)
				version, err := metadata.LatestRelease()
				if !errors.Is(err, errNoSuitableRelease) {
					t.Fatalf("LatestRelease() = %q, %v; want no suitable release error", version.ToString(), err)
				}
				if version.ToString() != "" {
					t.Errorf("failed selection returned candidate %q", version.ToString())
				}
			})
		}
	}
}

func TestLatestReleaseUP05RejectsInvalidReleaseOrFallbackEntry(t *testing.T) {
	for _, test := range []struct {
		name       string
		versioning string
		invalid    string
	}{
		{name: "invalid-release", versioning: `<release>broken</release>` + releaseVersionsXML("5.12.0"), invalid: "broken"},
		{name: "invalid-release-with-whitespace", versioning: `<release> broken </release>` + releaseVersionsXML("5.12.0"), invalid: " broken "},
		{name: "invalid-first-entry", versioning: releaseVersionsXML("broken", "5.12.0"), invalid: "broken"},
		{name: "invalid-last-entry", versioning: releaseVersionsXML("5.12.0", "broken"), invalid: "broken"},
		{name: "empty-entry", versioning: releaseVersionsXML("5.12.0", ""), invalid: ""},
		{name: "snapshot-release-invalid-list", versioning: `<release>5.13.0-SNAPSHOT</release>` + releaseVersionsXML("5.12.0", "broken"), invalid: "broken"},
	} {
		t.Run(test.name, func(t *testing.T) {
			metadata := releaseMetadataFromXML(t, test.versioning)
			version, err := metadata.LatestRelease()
			if err == nil || !strings.Contains(err.Error(), "unable to parse version:"+test.invalid+" due to ") {
				t.Fatalf("LatestRelease() = %q, %v; want version parse error", version.ToString(), err)
			}
			if version.ToString() != "" {
				t.Errorf("invalid metadata returned guessed candidate %q", version.ToString())
			}
		})
	}
}

func releaseMetadataFromXML(t *testing.T, versioning string) RepositoryMetadata {
	t.Helper()
	var metadata RepositoryMetadata
	if err := xml.Unmarshal([]byte("<metadata><versioning>"+versioning+"</versioning></metadata>"), &metadata); err != nil {
		t.Fatalf("parse test metadata XML: %v", err)
	}
	return metadata
}

func releaseVersionsXML(versions ...string) string {
	return "<versions><version>" + strings.Join(versions, "</version><version>") + "</version></versions>"
}

func assertLatestRelease(t *testing.T, metadata RepositoryMetadata, want string) {
	t.Helper()
	version, err := metadata.LatestRelease()
	if err != nil {
		t.Fatalf("LatestRelease() failed: %v", err)
	}
	if version.ToString() != want {
		t.Errorf("LatestRelease() = %q, want %q", version.ToString(), want)
	}
}
