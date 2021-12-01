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

				opts := &gl.ListTagsOptions{}
				hehe := &group{client: client}
				id := gitlab.NewNameWithBaseUrl(dep[j].Url, "source.golabs.io")
				tags, _, _ := hehe.client.Tags.ListTags(id.Get(), opts)
				//TODO: get latest tags from list of tags
				if len(tags) > 0 {
					c.Set(dep[j].Url, tags[0].Name, 0)
				}
			}
		}
	}
}
