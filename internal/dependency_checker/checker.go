package checker

import (
	"dependabot/internal/cache"
	parser "dependabot/internal/file_parser"

	gl "github.com/xanzy/go-gitlab"

	"github.com/gopaytech/go-commons/pkg/gitlab"
)

type DependencyChecker interface {
	Check(d *parser.Dependency)
}

type GithubDependencyChecker struct {
}

type GitlabDependencyChecker struct {
	Client *gl.Client
}

func (g *GitlabDependencyChecker) Check(d *parser.Dependency) string {
	newVersion := ""
	c := cache.ProvideCache()
	if val, found := c.Get(d.Url); found {
		newVersion = val.(string)
	} else {
		opts := &gl.ListReleasesOptions{}
		id := gitlab.NewNameWithBaseUrl(d.Url, "source.golabs.io")
		releases, _, _ := g.Client.Releases.ListReleases(id.Get(), opts)
		if len(releases) > 0 {
			newVersion = releases[0].TagName
			c.Set(d.Url, newVersion, 0)
		}
	}
	return newVersion
}

func (g *GithubDependencyChecker) Check(d *parser.Dependency) string {
	// TODO implement dependency checker with github client
	newVersion := ""
	return newVersion
}
