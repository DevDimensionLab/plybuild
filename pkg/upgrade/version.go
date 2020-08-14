package upgrade

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseVersion(version string) (JavaVersion, error) {

	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 4 {
		return JavaVersion{}, errors.New("could not parse version: " + version)
	}

	var patchPart = parts[2]
	var suffixSeparator = ""
	var suffix = ""
	if suffixParts := strings.Split(parts[2], "-"); len(suffixParts) > 1 {
		patchPart = suffixParts[0]
		suffix = suffixParts[1]
		suffixSeparator = "-"
	} else if len(parts) == 4 {
		suffixSeparator = "."
		suffix = parts[3]
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return JavaVersion{}, err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return JavaVersion{}, err
	}
	patch, err := strconv.Atoi(patchPart)
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

	if a.Suffix == "RELEASE" {
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

func (a JavaVersion) ToString() string {
	firstPart := fmt.Sprintf("%d.%d.%d", a.Major, a.Minor, a.Patch)

	if a.Suffix != "" {
		return fmt.Sprintf("%s%s%s", firstPart, a.SuffixSeparator, a.Suffix)
	} else {
		return firstPart
	}
}
