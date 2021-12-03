package task

import (
	checker "dependabot/internal/dependency_checker"
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

type DependencyChange struct {
	source   string
	old, new string
}

type Changes struct {
	project    *gl.Project
	DepChanges []DependencyChange
}

func CheckDependency(client *gl.Client, projects []Project) []Changes {
	ret := []Changes{}

	gitlabClient := &checker.GitlabDependencyChecker{Client: client}
	githubClient := &checker.GithubDependencyChecker{}

	for _, p := range projects {
		tmp := Changes{project: p.project, DepChanges: []DependencyChange{}}
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
				tmp.DepChanges = append(tmp.DepChanges, DependencyChange{source: dependency.Url, old: dependency.Version.String(), new: newVersion})
			}
		}

		ret = append(ret, tmp)
	}

	return ret
}
