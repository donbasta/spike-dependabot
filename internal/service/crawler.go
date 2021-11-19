package service

import (
	"log"

	gitlab "github.com/gopaytech/go-commons/pkg/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

type group struct {
	client *gl.Client
}

func CrawlGroup(client *gl.Client) ([]*gl.Project, error) {
	// groupID := 3663 //roles
	groupID := 3725 //playbooks
	id := gitlab.NameOrId{ID: groupID}

	listOpts := gl.ListOptions{
		Page:    0,
		PerPage: 100,
	}

	opts := &gl.ListGroupProjectsOptions{IncludeSubgroups: gl.Bool(true), ListOptions: listOpts}
	hehe := &group{client: client}
	result, response, _ := hehe.client.Groups.ListGroupProjects(id.Get(), opts)
	log.Println(response)
	return result, nil
}
