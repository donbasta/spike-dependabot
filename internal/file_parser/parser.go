package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Semver struct {
	major, minor, patch int
}

func MakeVersion(semver string) Semver {
	semver = strings.Trim(semver, " ")
	semverTrimV := semver[1:]
	nums := strings.Split(semverTrimV, ".")
	major, _ := strconv.Atoi(nums[0])
	minor, _ := strconv.Atoi(nums[1])
	patch, _ := strconv.Atoi(nums[2])
	return Semver{major, minor, patch}
}

func (s Semver) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.major, s.minor, s.patch)
}

type Dependency struct {
	Url     string
	Version Semver
}

type DependencyParser interface {
	Parse(fileContent string) ([]Dependency, error)
}
