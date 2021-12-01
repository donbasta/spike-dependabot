package service

import (
	"dependabot/internal/service/helper"

	gitlab "github.com/gopaytech/go-commons/pkg/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

type group struct {
	client *gl.Client
}

type empty struct{}

func CrawlGroups(client *gl.Client) ([]Project, error) {
	watchedProjects := []Project{}
	groupIDs := helper.GetGroupList()
	for i := 0; i < len(groupIDs); i++ {
		groupID := groupIDs[i]
		projects, err := crawlGroup(client, groupID)
		if err != nil {
			return nil, err
		}
		size := len(projects)

		sem := make(chan empty, size)
		for j := 0; j < size; j++ {
			go func(p *gl.Project) {
				deps, _ := ParseProject(client, p)
				watchedProjects = append(watchedProjects, Project{p, deps})
				sem <- empty{}
			}(projects[j])
		}
		for j := 0; j < size; j++ {
			<-sem
		}
	}
	return watchedProjects, nil
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
