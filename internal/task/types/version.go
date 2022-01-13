package types

import (
	"dependabot/internal/errors"
	"fmt"
	"strconv"
	"strings"
)

type SemanticVersion struct {
	major, minor, patch int
}

func MakeVersion(semver string) (*SemanticVersion, error) {
	semver = strings.Trim(semver, " ")
	if semver == "latest" {
		return &SemanticVersion{0, 0, 0}, nil
	}

	semver = strings.TrimLeft(semver, "v")
	nums := strings.Split(semver, ".")
	if len(nums) != 3 {
		return nil, errors.NewInvalidVersionError(nil, "version invalid")
	}

	major, errMajor := strconv.Atoi(nums[0])
	minor, errMinor := strconv.Atoi(nums[1])
	patch, errPatch := strconv.Atoi(nums[2])
	if errMajor != nil || errMinor != nil || errPatch != nil {
		return nil, errors.NewInvalidVersionError(nil, "invalid version, not an integer")
	}

	return &SemanticVersion{major, minor, patch}, nil
}

func LatestVersion(semver *SemanticVersion) bool {
	return semver.major == 0 && semver.minor == 0 && semver.patch == 0
}

func (s SemanticVersion) String() string {
	if LatestVersion(&s) {
		return "latest"
	}

	return fmt.Sprintf("v%d.%d.%d", s.major, s.minor, s.patch)
}
