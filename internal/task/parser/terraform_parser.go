package parser

import (
	"dependabot/internal/errors"
	"dependabot/internal/task/types"
	"strings"
)

type terraformParser struct {
}

func CreateTerraformParser() *terraformParser {
	return &terraformParser{}
}

func (tf *terraformParser) GetBaseUrlFromRawSource(raw string) (string, error) {
	remotes := []string{"github.com", "source.golabs.io"}

	for _, remote := range remotes {
		if idx := strings.Index(raw, remote); idx != -1 {
			return raw[idx:], nil
		}
	}

	return "", errors.NewInvalidSourceError(nil, "Invalid or not supported source")
}

func (tf *terraformParser) convertSourceToDependency(source string) (types.Dependency, error) {
	url, err := tf.GetBaseUrlFromRawSource(source)
	if err != nil {
		return types.Dependency{}, err
	}

	tokens := strings.Split(url, "?")
	baseUrl := tokens[0]
	baseUrl = strings.TrimSuffix(baseUrl, ".git")

	dependency := types.Dependency{
		SourceRaw:     source,
		SourceBaseUrl: baseUrl,
		Type:          "terraform",
	}

	var version *types.SemanticVersion
	if len(tokens) == 1 {
		version, _ = types.MakeVersion("latest")
	} else {
		version, _ = types.MakeVersion(strings.TrimLeft(tokens[1], "ref="))
	}
	dependency.Version = *version

	return dependency, nil
}

func (tf *terraformParser) ParseRequirementFile(fileContent string) ([]types.Dependency, error) {
	lines := strings.Split(fileContent, "\n")
	dependencies := []types.Dependency{}
	isModule := false

	for _, line := range lines {
		line := strings.Trim(line, " ")
		if len(line) == 0 {
			continue
		}

		tokens := strings.Fields(line)
		if tokens[0] == "module" {
			isModule = true
		}

		attr := tokens[0]
		if isModule && (attr == "source") && len(tokens) >= 3 {
			rawSource := strings.Trim(tokens[2], "\"")
			buffer, err := tf.convertSourceToDependency(rawSource)
			if err != nil {
				return []types.Dependency{}, err
			}

			dependencies = append(dependencies, buffer)
			isModule = false
		}
	}

	return dependencies, nil
}
