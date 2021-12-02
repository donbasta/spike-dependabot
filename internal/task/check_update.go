package task

import (
	"dependabot/internal/cache"
	"strings"

	gitlab "github.com/gopaytech/go-commons/pkg/gitlab"
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
	c := cache.ProvideCache()
	ret := []Changes{}

	for i := 0; i < len(projects); i++ {
		dep := projects[i].dependencies
		tmp := Changes{project: projects[i].project, DepChanges: []DependencyChange{}}
		for j := 0; j < len(dep); j++ {
			newVersion := ""
			source := getSourceTypeFromURL(dep[j].Url)
			if source == "" {
				continue
			}

			if source == "github" {
				//TODO: fetch version from github
				continue
			}

			if source == "golabs" {
				if val, found := c.Get(dep[j].Url); found {
					newVersion = val.(string)
				} else {
					opts := &gl.ListReleasesOptions{}
					hehe := &group{client: client}
					id := gitlab.NewNameWithBaseUrl(dep[j].Url, "source.golabs.io")
					releases, _, _ := hehe.client.Releases.ListReleases(id.Get(), opts)
					if len(releases) > 0 {
						newVersion = releases[0].TagName
						c.Set(dep[j].Url, newVersion, 0)
					}
				}
			}

			if newVersion != "" && newVersion != dep[j].Version.String() {
				tmp.DepChanges = append(tmp.DepChanges, DependencyChange{source: dep[j].Url, old: dep[j].Version.String(), new: newVersion})
			}
		}
		ret = append(ret, tmp)
	}

	return ret
}
