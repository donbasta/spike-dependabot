package crawler

import (
	"dependabot/internal/task/parser"
	"dependabot/internal/task/types"
	packageManager "dependabot/internal/task/types/package_manager"
	"encoding/base64"

	gitlab "github.com/gopaytech/go-commons/pkg/gitlab"
	gl "github.com/xanzy/go-gitlab"
)

func createDefaultListOptions() gl.ListOptions {
	return gl.ListOptions{
		Page:    0,
		PerPage: 100,
	}
}

func CrawlMultipleGroups(client *gl.Client, groupIds []int) ([]types.ProjectDependencies, error) {
	crawledProjects := []types.ProjectDependencies{}

	for _, groupId := range groupIds {
		projects, err := CrawlGroupProjects(client, groupId)
		if err != nil {
			return nil, err
		}

		ch := make(chan struct{}, len(projects))

		for _, project := range projects {
			go func(p *gl.Project) {
				dependencies, _ := CrawlProjectFilesAndGetDependencies(client, p)
				crawledProjects = append(crawledProjects, types.ProjectDependencies{
					Project:      p,
					Dependencies: dependencies,
				})
				ch <- struct{}{}
			}(project)
		}

		for range projects {
			<-ch
		}
	}

	return crawledProjects, nil
}

func CrawlGroupProjects(client *gl.Client, groupId int) ([]*gl.Project, error) {
	id := gitlab.NameOrId{ID: groupId}
	listOptions := createDefaultListOptions()
	listGroupProjectOptions := &gl.ListGroupProjectsOptions{
		IncludeSubgroups: gl.Bool(true),
		ListOptions:      listOptions,
	}
	result, _, _ := client.Groups.ListGroupProjects(id.Get(), listGroupProjectOptions)

	return result, nil
}

func CrawlProjectFilesAndGetDependencies(client *gl.Client, project *gl.Project) ([]types.Dependency, error) {
	ansiblePackageManager := packageManager.CreateAnsiblePackageManager()
	ansibleParser := parser.CreateAnsibleParser()

	terraformPackageManager := packageManager.CreateTerraformPackageManager()
	terraformParser := parser.CreateTerraformParser()

	opt := &gl.ListTreeOptions{
		ListOptions: createDefaultListOptions(),
		Recursive:   gl.Bool(true),
	}
	projectId := project.ID
	repositoryTree, _, _ := client.Repositories.ListTree(projectId, opt)

	dependencies := []types.Dependency{}
	for _, child := range repositoryTree {
		fileOptions := &gl.GetFileOptions{
			Ref: gl.String("master"),
		}

		path := child.Path
		file, _, _ := client.RepositoryFiles.GetFile(projectId, path, fileOptions)

		if ansiblePackageManager.IsPackageDependencyRequirementFile(path) {
			content, _ := base64.StdEncoding.DecodeString(file.Content)
			dep, _ := ansibleParser.ParseRequirementFile(string(content))
			dependencies = append(dependencies, dep...)
		}

		if terraformPackageManager.IsPackageDependencyRequirementFile(path) {
			content, _ := base64.StdEncoding.DecodeString(file.Content)
			dep, _ := terraformParser.ParseRequirementFile(string(content))
			dependencies = append(dependencies, dep...)
		}
	}

	return dependencies, nil
}
