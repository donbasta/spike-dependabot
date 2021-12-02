package task

import (
	"dependabot/internal/task/handler"

	gitlab "github.com/gopaytech/go-commons/pkg/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

type group struct {
	client *gl.Client
}

func CrawlGroups(client *gl.Client) ([]Project, error) {
	watchedProjects := []Project{}
	groupIDs := handler.GetGroupList()
	for _, groupID := range groupIDs {
		projects, err := crawlGroup(client, groupID)
		if err != nil {
			return nil, err
		}

		sem := make(chan struct{}, len(projects))
		for _, project := range projects {
			go func(p *gl.Project) {
				deps, _ := ParseProject(client, p)
				watchedProjects = append(watchedProjects, Project{p, deps})
				sem <- struct{}{}
			}(project)
		}
		for j := 0; j < len(projects); j++ {
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
