package service

import (
	"dependabot/internal/cache"
	"log"
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

func CheckDependency(client *gl.Client, projects []Project) {
	c := cache.ProvideCache()

	for i := 0; i < len(projects); i++ {
		dep := projects[i].dependencies
		for j := 0; j < len(dep); j++ {
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
					log.Printf("got from cache for url %s: version = %s\n", dep[j].Url, val.(string))
					continue
				}

				opts := &gl.ListReleasesOptions{}
				hehe := &group{client: client}
				id := gitlab.NewNameWithBaseUrl(dep[j].Url, "source.golabs.io")
				releases, _, _ := hehe.client.Releases.ListReleases(id.Get(), opts)
				if len(releases) > 0 {
					c.Set(dep[j].Url, releases[0].TagName, 0)
				}
				for k := 0; k < len(releases); k++ {
					log.Printf("%s %s %s", releases[k].TagName, releases[k].Name, releases[k].ReleasedAt)
				}
			}
		}
	}
}
