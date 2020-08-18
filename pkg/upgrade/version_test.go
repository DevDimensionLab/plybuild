package upgrade

import "testing"

func TestParseVersion(t *testing.T) {

	_, err := ParseVersion("foobar")
	if err == nil {
		t.Errorf("foobar should not be pased correctly as a version")
	}

	version, err := ParseVersion("1.2.3")
	if err != nil {
		t.Errorf("err")
		panic(err)
	}

	if version.Major != 1 {
		t.Errorf("Major should be 1")
	}
	if version.Minor != 2 {
		t.Errorf("Minor should be 2")
	}
	if version.Patch != 3 {
		t.Errorf("Patch should be 3")
	}
}

func TestParseVersionWithSuffix(t *testing.T) {
	version, err := ParseVersion("1.2.3-rc")
	if err != nil {
		t.Errorf("err")
		panic(err)
	}

	if version.Suffix != "rc" {
		t.Errorf("Suffix should be rc")
	}
}

func TestIsReleaseVersion(t *testing.T) {
	notRelease, _ := ParseVersion("1.2.3-rc")
	if notRelease.IsReleaseVersion() {
		t.Errorf("1.2.3-rc is not a release version")
	}

	notRelease2, _ := ParseVersion("1.2.3-m3")
	if notRelease2.IsReleaseVersion() {
		t.Errorf("1.2.3-m3 is not a release version")
	}

	release1, _ := ParseVersion("1.2.3")
	if !release1.IsReleaseVersion() {
		t.Errorf("1.2.3 is a release version")
	}

	release2, _ := ParseVersion("1.2.3.RELEASE")
	if !release2.IsReleaseVersion() {
		t.Errorf("1.2.3.RELEASE is a release version")
	}
}

func TestIsDifferentVersion(t *testing.T) {
	sameVersionA, _ := ParseVersion("1.2.3")
	sameVersionB, _ := ParseVersion("1.2.3")
	if sameVersionA.IsDifferentFrom(sameVersionB) {
		t.Errorf("sameVersionA and sameVersionB is not the same")
	}

	differentVersionA, _ := ParseVersion("1.2.3")
	differentVersionB, _ := ParseVersion("1.2.4")

	if !differentVersionA.IsDifferentFrom(differentVersionB) {
		t.Errorf("differentVersionA and differentVersionB is not the same version")
	}
}

func TestVersionToString(t *testing.T) {
	versionA := "1.2.3"
	parsedVersionA, _ := ParseVersion(versionA)

	if parsedVersionA.ToString() != versionA {
		t.Errorf("VersionToString of parsedVersionA should be: " + versionA)
	}

	versionB := "1.2.3-rc"
	parsedVersionB, _ := ParseVersion(versionB)

	if parsedVersionB.ToString() != versionB {
		t.Errorf("VersionToString of parsedVersionB should be: " + versionB)
	}

	versionC := "1.2.3.RELEASE"
	parsedVersionC, _ := ParseVersion(versionC)

	if parsedVersionC.ToString() != versionC {
		t.Errorf("VersionToString of parsedVersionC should be: " + versionC)
	}
}

func TestIsMajorUpgrade(t *testing.T) {
	mainVersion, _ := ParseVersion("1.2.3")
	newMinorVersion, _ := ParseVersion("1.3.0")
	newMajorVersion, _ := ParseVersion("2.2.3")

	if IsMajorUpgrade(mainVersion, newMinorVersion) {
		t.Errorf("%s is not a major upgrade from %s", newMinorVersion.ToString(), mainVersion.ToString())
	}

	if !IsMajorUpgrade(mainVersion, newMajorVersion) {
		t.Errorf("%s should be a major upgrade from %s", newMajorVersion.ToString(), mainVersion.ToString())
	}
}

func TestJustOneDigitVersionAsPatchVersion(t *testing.T) {
	version := "1"
	parsedVersion, err := ParseVersion(version)

	if err != nil {
		t.Errorf("should accept version=%s, got: %s", version, err.Error() )
	}

	if parsedVersion.Major != 0 {
		t.Errorf("expected major to be 0 for one digit version, not %d", parsedVersion.Major )
	}

	if parsedVersion.Minor != 0 {
		t.Errorf("expected minor to be 0 for one digit version, not %d", parsedVersion.Minor )
	}

	if parsedVersion.Patch != 1 {
		t.Errorf("expected patch to be %s for one digit version, not %d", version, parsedVersion.Major )
	}

}

func TestJustTwoDigitVersionAsPatchVersion(t *testing.T) {
	version := "4.2"
	parsedVersion, err := ParseVersion(version)

	if err != nil {
		t.Errorf("should accept version=%s, got: %s", version, err.Error() )
	}

	if parsedVersion.Major != 0 {
		t.Errorf("expected major to be 0 for two digit version, not %d", parsedVersion.Major )
	}

	if parsedVersion.Minor != 4 {
		t.Errorf("expected minor to be 4 for two digit version, not %d", parsedVersion.Minor )
	}

	if parsedVersion.Patch != 2 {
		t.Errorf("expected patch to be 2 for one digit version, not %d",  parsedVersion.Patch )
	}
}


