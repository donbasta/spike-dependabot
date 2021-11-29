package config

import (
	"github.com/gopaytech/go-commons/pkg/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

func ProvideConfig() *Main {
	main := Config()
	return main
}

func ProvideGitlabClient(config *Main) (*gl.Client, error) {
	return gitlab.NewClient(config.Git.URL, config.Git.Token)
}
