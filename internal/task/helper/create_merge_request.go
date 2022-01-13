package helper

import (
	"dependabot/internal/config"

	"github.com/gopaytech/go-commons/pkg/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

func CreateMergeRequest(projectId int, branchSourceName string, branchTargetName string, mergeRequestTitle string) (*gl.MergeRequest, error) {
	cfg := config.ProvideConfig()
	client, _ := config.ProvideGitlabClient(cfg)
	mergeRequestClient := gitlab.NewMergeRequest(client)

	id := gitlab.NameOrId{ID: projectId}
	mergeRequest, err := mergeRequestClient.Create(id, branchSourceName, branchTargetName, mergeRequestTitle)
	if err != nil {
		return nil, err
	}

	return mergeRequest, nil
}
