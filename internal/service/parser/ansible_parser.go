package parser

import (
	"log"
	"regexp"

	"gopkg.in/yaml.v2"
)

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
		log.Println(d[i])
		version := d[i].Version
		matched := reVersion.MatchString("aaxbb")
		if !matched {
			//TODO: should change to latest
			version = "v0.0.0"
		}
		deps = append(deps, Dependency{Url: d[i].Src, Version: MakeVersion(version)})
	}

	return deps, nil
}
