package parser

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type AnsibleParser struct {
}

type AnsibleDependency struct {
	Name    string `yaml:"name,omitempty"`
	Src     string `yaml:"src,omitempty"`
	Version string `yaml:"version,omitempty"`
	Scm     string `yaml:"scm,omitempty"`
}

func (an *AnsibleParser) Parse(fileContent string) ([]Dependency, error) {
	byteContent := []byte(fileContent)

	var d []AnsibleDependency
	yaml.Unmarshal(byteContent, &d)

	deps := []Dependency{}

	reVersion, _ := regexp.Compile(`^v\d*\.\d*\.\d*$`)

	for i := 0; i < len(d); i++ {
		version := d[i].Version
		matched := reVersion.MatchString(version)
		if !matched {
			//TODO: should change to latest or default handling when src is absent
			version = "v0.0.0"
		}
		url := d[i].Src
		if len(url) > 4 && url[:4] == "git@" {
			url = url[4:]
			url = strings.Replace(url, ":", "/", 1)
		}
		if len(url) > 4 && url[len(url)-4:] == ".git" {
			url = url[:len(url)-4]
		}
		deps = append(deps, Dependency{Url: url, Version: MakeVersion(version)})
	}

	return deps, nil
}
