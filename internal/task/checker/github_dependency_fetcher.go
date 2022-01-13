package checker

type GithubDependencyVersionFetcher struct {
}

func (g GithubDependencyVersionFetcher) GetDependencyNewVersion(dependencyUrl string) (string, error) {
	var newVersion string
	return newVersion, nil
}
