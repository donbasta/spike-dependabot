package task

import (
	checker "dependabot/internal/dependency_checker"
	updater "dependabot/internal/task/file_updater"
	"strings"

	gl "github.com/xanzy/go-gitlab"
)

func getSourceTypeFromURL(url string) string {
	if strings.Contains(url, "source.golabs.io") {
		return "golabs"
	}
	if strings.Contains(url, "github.com") {
		return "github"
	}
	return ""
}

func CheckDependency(client *gl.Client, projects []Project) []updater.Changes {
	ret := []updater.Changes{}

	gitlabClient := &checker.GitlabDependencyChecker{Client: client}
	githubClient := &checker.GithubDependencyChecker{}

	for _, p := range projects {
		tmp := updater.Changes{Project: p.project, DepChanges: []updater.DependencyChange{}}
		for _, dependency := range p.dependencies {
			newVersion := ""
			source := getSourceTypeFromURL(dependency.Url)
			if source == "" {
				continue
			}

			if source == "github" {
				newVersion = githubClient.Check(&dependency)
			}

			if source == "golabs" {
				newVersion = gitlabClient.Check(&dependency)
			}

			if newVersion != "" && newVersion != dependency.Version.String() {
				tmp.DepChanges = append(tmp.DepChanges, updater.DependencyChange{Source: dependency.Url, Old: dependency.Version.String(), New: newVersion})
			}
		}

		ret = append(ret, tmp)
	}

	return ret
}
