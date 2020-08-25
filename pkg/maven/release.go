package maven

import (
	"errors"
	"sort"
)

func (meta Metadata) LatestRelease() (JavaVersion, error) {
	version, err := ParseVersion(meta.Versioning.Release)
	if err != nil {
		return JavaVersion{}, err
	}

	if !version.IsReleaseVersion() {
		return getLatestRelease(meta.Versioning.Versions.Version)
	}

	return version, nil
}

func getLatestRelease(versions []string) (JavaVersion, error) {
	var javaVersions []JavaVersion
	for _, versionString := range versions {
		version, err := ParseVersion(versionString)
		if err != nil {
			return JavaVersion{}, err
		}
		javaVersions = append(javaVersions, version)
	}

	sort.Sort(VersionSort(javaVersions))

	for i := len(javaVersions) - 1; i >= 0; i-- {
		if javaVersions[i].IsReleaseVersion() {
			return javaVersions[i], nil
		}
	}

	return JavaVersion{}, errors.New("could not find a suitable release version")
}
