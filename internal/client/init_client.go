package client

import (
	gitlab "github.com/gopaytech/go-commons/pkg/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

func CreateClient(token string) *gl.Client {
	client, _ := gitlab.NewClient("https://source.golabs.io", token)
	return client
}
