package config

import (
	gitTransport "github.com/go-git/go-git/v5/plumbing/transport"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
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

func ProvideGitAuth(config *Main) gitTransport.AuthMethod {
	return &gitHttp.BasicAuth{
		Username: config.Git.Username,
		Password: config.Git.Token,
	}
}
