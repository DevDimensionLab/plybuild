package maven

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseVersion(version string) (JavaVersion, error) {

	parts := strings.Split(version, ".")

	patchIndex := 2
	minorIndex := 1
	majorIndex := 0
	switch len(parts) {
	case 1:
		majorIndex = -1
		minorIndex = -1
		patchIndex = 0
	case 2:
		majorIndex = -1
		minorIndex = 0
		patchIndex = 1
	}

	var patchPart = parts[patchIndex]
	var suffixSeparator = ""
	var suffix = ""
	if suffixParts := strings.Split(parts[patchIndex], "-"); len(suffixParts) > 1 {
		patchPart = suffixParts[0]
		suffix = suffixParts[1]
		suffixSeparator = "-"
	} else if len(parts) == 4 {
		suffixSeparator = "."
		suffix = parts[3]
	}

	var err error
	var major = 0
	var minor = 0
	var patch = 0

	if -1 < majorIndex {
		major, err = strconv.Atoi(parts[majorIndex])
		if err != nil {
			return JavaVersion{}, err
		}
	}

	if -1 < minorIndex {
		minor, err = strconv.Atoi(parts[minorIndex])
		if err != nil {
			return JavaVersion{}, err
		}
	}
	patch, err = strconv.Atoi(patchPart)
	if err != nil {
		return JavaVersion{}, err
	}

	return JavaVersion{
		Major:           major,
		Minor:           minor,
		Patch:           patch,
		Suffix:          suffix,
		SuffixSeparator: suffixSeparator,
	}, nil
}

func (a JavaVersion) IsReleaseVersion() bool {
	if a.Suffix == "" {
		return true
	}
	if strings.ToUpper(a.Suffix) == "RELEASE" {
		return true
	}
	if strings.ToUpper(a.Suffix) == "FINAL" {
		return true
	}
	return false
}

func IsMajorUpgrade(old JavaVersion, new JavaVersion) bool {
	if old.Major < new.Major {
		return true
	}

	return false
}

func (a JavaVersion) IsDifferentFrom(b JavaVersion) bool {
	if a.Major != b.Major {
		return true
	}
	if a.Minor != b.Minor {
		return true
	}
	if a.Patch != b.Patch {
		return true
	}
	if a.Suffix != b.Suffix {
		return true
	}
	return false
}

func (a JavaVersion) IsLessThan(b JavaVersion) bool {
	aString := fmt.Sprintf("%09d%09d%09d%s", a.Major, a.Minor, a.Patch, a.Suffix)
	bString := fmt.Sprintf("%09d%09d%09d%s", b.Major, b.Minor, b.Patch, b.Suffix)
	return aString < bString
}

func (a JavaVersion) ToString() string {
	var firstPart = ""

	if a.Major == 0 {
		if a.Minor == 0 {
			firstPart = fmt.Sprintf("%d", a.Patch)
		} else {
			firstPart = fmt.Sprintf("%d.%d", a.Minor, a.Patch)
		}
	} else {
		firstPart = fmt.Sprintf("%d.%d.%d", a.Major, a.Minor, a.Patch)
	}

	if a.Suffix != "" {
		return fmt.Sprintf("%s%s%s", firstPart, a.SuffixSeparator, a.Suffix)
	} else {
		return firstPart
	}
}

type VersionSort []JavaVersion

func (a VersionSort) Len() int      { return len(a) }
func (a VersionSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a VersionSort) Less(i, j int) bool {
	return a[i].IsLessThan(a[j])
}
