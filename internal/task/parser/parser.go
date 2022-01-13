package parser

import "dependabot/internal/task/types"

type DependencyParser interface {
	ParseRequirementFile(fileContent string) ([]types.Dependency, error)
	GetBaseUrlFromRawSource(raw string) (string, error)
}
