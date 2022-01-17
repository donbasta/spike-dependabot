package checker

import (
	"dependabot/internal/cache"
	"strings"

	"github.com/gopaytech/go-commons/pkg/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

type GolabsDependencyVersionFetcher struct {
	Client *gl.Client
}

func (g *GolabsDependencyVersionFetcher) getBaseUrl() string {
	return "source.golabs.io"
}

func (g *GolabsDependencyVersionFetcher) getProjectLatestRelease(projectUrl string) (*gl.Release, error) {
	options := &gl.ListReleasesOptions{}
	projectUrl = strings.TrimPrefix(projectUrl, "https://")
	id := gitlab.NewNameWithBaseUrl(projectUrl, g.getBaseUrl())
	releases, _, err := g.Client.Releases.ListReleases(id.Get(), options)

	if err != nil {
		return nil, err
	}

	if len(releases) > 0 {
		return releases[0], nil
	}

	return nil, nil
}

func (g GolabsDependencyVersionFetcher) GetDependencyNewVersion(dependencyUrl string) (string, error) {
	var newVersion string

	c := cache.ProvideCache()

	if val, found := c.Get(dependencyUrl); found {
		newVersion = val.(string)
	} else {
		latestRelease, err := g.getProjectLatestRelease(dependencyUrl)
		if err != nil {
			return "", err
		}

		if latestRelease != nil {
			newVersion = latestRelease.TagName
			c.Set(dependencyUrl, newVersion, 0)
		}
	}

	return newVersion, nil
}
