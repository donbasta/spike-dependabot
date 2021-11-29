package service

import (
	"dependabot/internal/service/helper"
	"log"

	gitlab "github.com/gopaytech/go-commons/pkg/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

type group struct {
	client *gl.Client
}

func CrawlGroups(client *gl.Client) {
	groupIDs := helper.GetGroupList()
	for i := 0; i < len(groupIDs); i++ {
		groupID := groupIDs[i]
		projects, err := crawlGroup(client, groupID)
		if err != nil {
			return
		}
		size := len(projects)
		for i := 0; i < size; i++ {
			deps, _ := ParseProject(client, projects[i])
			for j := 0; j < len(deps); j++ {
				log.Println(deps[j].Url, " ", deps[j].Version)
			}
		}
	}
}

func crawlGroup(client *gl.Client, groupID int) ([]*gl.Project, error) {
	id := gitlab.NameOrId{ID: groupID}

	listOpts := gl.ListOptions{
		Page:    0,
		PerPage: 100,
	}

	opts := &gl.ListGroupProjectsOptions{IncludeSubgroups: gl.Bool(true), ListOptions: listOpts}
	hehe := &group{client: client}
	result, _, _ := hehe.client.Groups.ListGroupProjects(id.Get(), opts)
	return result, nil
}
