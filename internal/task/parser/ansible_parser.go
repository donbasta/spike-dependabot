package parser

import (
	"dependabot/internal/task/types"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type ansibleParser struct {
}

type AnsibleDependency struct {
	Name    string `yaml:"name,omitempty"`
	Src     string `yaml:"src,omitempty"`
	Version string `yaml:"version,omitempty"`
	Scm     string `yaml:"scm,omitempty"`
}

func CreateAnsibleParser() *ansibleParser {
	return &ansibleParser{}
}

func (an *ansibleParser) GetBaseUrlFromRawSource(raw string) (string, error) {
	baseUrl := raw
	if strings.HasPrefix(baseUrl, "git@") {
		baseUrl = strings.TrimLeft(baseUrl, "git@")
		baseUrl = strings.Replace(baseUrl, ":", "/", 1)
	}

	baseUrl = strings.TrimPrefix(baseUrl, "https://")
	baseUrl = strings.TrimSuffix(baseUrl, ".git")
	return baseUrl, nil
}

func (an *ansibleParser) ParseRequirementFile(fileContent string) ([]types.Dependency, error) {
	byteContent := []byte(fileContent)

	var ansibleDependencies []AnsibleDependency
	yaml.Unmarshal(byteContent, &ansibleDependencies)

	dependencies := []types.Dependency{}

	regexVersion, _ := regexp.Compile(`^v\d*\.\d*\.\d*$`)

	for _, ansibleDependency := range ansibleDependencies {
		version := ansibleDependency.Version
		matched := regexVersion.MatchString(version)
		if !matched {
			version = "latest"
		}

		url, _ := an.GetBaseUrlFromRawSource(ansibleDependency.Src)

		v, _ := types.MakeVersion(version)
		dependencyFromFile := types.Dependency{
			SourceRaw:     ansibleDependency.Src,
			SourceBaseUrl: url,
			Version:       *v,
			Type:          "ansible",
		}
		dependencies = append(dependencies, dependencyFromFile)
	}

	return dependencies, nil
}
