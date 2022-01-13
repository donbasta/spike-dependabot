package checker

import (
	"dependabot/internal/errors"
	"dependabot/internal/task/types"
	"log"
	"strings"

	gl "github.com/xanzy/go-gitlab"
)

func getSourceTypeFromURL(url string) (string, error) {
	if strings.Contains(url, "source.golabs.io") {
		return "golabs", nil
	}

	if strings.Contains(url, "github.com") {
		return "github", nil
	}

	return "", errors.NewInvalidSourceError(nil, "Invalid or not supported source")
}

func CheckMultipleProjectsDependencyUpdate(client *gl.Client, projects []types.ProjectDependencies) []types.ProjectDependencies {
	allDependencyUpdates := []types.ProjectDependencies{}

	dependencyFetcher := map[string]DependencyFetcher{
		"golabs": GolabsDependencyVersionFetcher{Client: client},
		"github": GithubDependencyVersionFetcher{},
	}

	for _, project := range projects {
		projectUpdates := types.ProjectDependencies{
			Project:      project.Project,
			Dependencies: []types.Dependency{},
		}

		for _, dependency := range project.Dependencies {
			var currentDependencyVersion string
			source, err := getSourceTypeFromURL(dependency.SourceBaseUrl)

			if err != nil {
				continue
			}

			currentDependencyVersion, err = dependencyFetcher[source].GetDependencyNewVersion(dependency.SourceBaseUrl)

			if err != nil {
				log.Fatalln(err)
				continue
			}

			projectDependencyVersion := dependency.Version.String()
			if currentDependencyVersion != projectDependencyVersion && projectDependencyVersion != "latest" {
				v, _ := types.MakeVersion(currentDependencyVersion)
				dependencyUpdate := types.Dependency{
					SourceRaw:     dependency.SourceRaw,
					SourceBaseUrl: dependency.SourceBaseUrl,
					Version:       *v,
					Type:          dependency.Type,
				}

				projectUpdates.Dependencies = append(projectUpdates.Dependencies, dependencyUpdate)
			}
		}

		allDependencyUpdates = append(allDependencyUpdates, projectUpdates)
	}

	return allDependencyUpdates
}
