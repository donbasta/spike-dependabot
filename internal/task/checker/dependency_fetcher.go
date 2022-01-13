package checker

type DependencyFetcher interface {
	GetDependencyNewVersion(dependencyUrl string) (string, error)
}
