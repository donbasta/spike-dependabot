package updater

import (
	"dependabot/internal/config"
	"dependabot/internal/db/entity"
	"dependabot/internal/errors"
	"dependabot/internal/service"
	"dependabot/internal/task/helper"
	"dependabot/internal/task/parser"
	"dependabot/internal/task/types"
	packageManager "dependabot/internal/task/types/package_manager"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gopaytech/go-commons/pkg/git"
	"gopkg.in/yaml.v2"

	gitTransport "github.com/go-git/go-git/v5/plumbing/transport"
)

type AnsibleUpdater struct {
	GitCloneFunc git.CloneFunc
	GitAuth      gitTransport.AuthMethod
	Parser       parser.DependencyParser
	Manager      packageManager.PackageManager
}

func CreateAnsibleUpdater() *AnsibleUpdater {
	mainCfg := config.ProvideConfig()
	authMethod := config.ProvideGitAuth(mainCfg)
	return &AnsibleUpdater{
		GitCloneFunc: git.Clone,
		GitAuth:      authMethod,
		Parser:       parser.CreateAnsibleParser(),
		Manager:      packageManager.CreateAnsiblePackageManager(),
	}
}

func (a *AnsibleUpdater) IsPackageDependencyRequirementFile(filepath string) bool {
	return a.Manager.IsPackageDependencyRequirementFile(filepath)
}

func (a *AnsibleUpdater) GetPackageManagerName() string {
	return a.Manager.GetPackageManagerName()
}

func (a *AnsibleUpdater) ProcessUpdateProjectDependencies(c *types.ProjectDependencies, s *service.SlackNotificationService, m *service.MergeRequestService) error {
	dt := time.Now()
	gitWorkingBranchName := fmt.Sprintf("scp-dependency-manager-bump/%s/%s", a.Manager.GetPackageManagerName(), dt.Format("01-02-2006"))
	mergeRequestTitle := fmt.Sprintf("fix: bump %s dependencies (%s)", a.Manager.GetPackageManagerName(), dt.Format("01-02-2006"))

	err := a.updateProjectDependencyAndCommitChanges(c, gitWorkingBranchName, mergeRequestTitle)
	if err != nil {
		return err
	}

	mergeRequest, err := helper.CreateMergeRequest(c.Project.ID, gitWorkingBranchName, "master", mergeRequestTitle)
	if err != nil {
		return err
	}
	if mergeRequest == nil {
		return err
	}

	mergeRequestURL := mergeRequest.WebURL
	if mergeRequestURL == "" {
		return err
	}

	repositoryId := strconv.Itoa(c.Project.ID)
	mergeRequestEntity := &entity.MergeRequest{
		RepositoryID:    repositoryId,
		MergeRequestIID: strconv.Itoa(mergeRequest.IID),
		RepositoryURL:   c.Project.HTTPURLToRepo,
	}

	_, err = (*m).Create(mergeRequestEntity)
	if err != nil {
		return err
	}

	mainCfg := config.ProvideConfig()
	err = (*s).NotifyMerge(mainCfg.Slack.ChannelId, mergeRequestURL, c.Project.Name, a.Manager.GetPackageManagerName())
	if err != nil {
		return err
	}

	return nil
}

func (a *AnsibleUpdater) updateContentWithNewDependency(fileContent string, dependency types.Dependency) string {
	byteContent := []byte(fileContent)
	var ansibleDependencies []packageManager.AnsibleDependency
	yaml.Unmarshal(byteContent, &ansibleDependencies)

	for i := range ansibleDependencies {
		if ansibleDependencies[i].Src == dependency.SourceRaw {
			ansibleDependencies[i].Version = dependency.Version.String()
		}
	}

	appendDashes := false
	var attrOrders [][]string
	buffer := []string{}
	linesOriginalContent := strings.Split(fileContent, "\n")
	for _, line := range linesOriginalContent {
		if line == "---" {
			appendDashes = true
			continue
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "- ") && len(buffer) > 0 {
			attrOrders = append(attrOrders, buffer)
			buffer = []string{}
		}

		line = strings.TrimLeft(line, "- ")
		tokens := strings.Split(line, ":")
		if len(tokens) == 0 {
			continue
		}
		buffer = append(buffer, tokens[0])
	}
	attrOrders = append(attrOrders, buffer)

	updatedByteContent, _ := yaml.Marshal(ansibleDependencies)
	updatedContent := string(updatedByteContent)
	lines := strings.Split(updatedContent, "\n")
	linesUpdatedContentNew := []string{}

	dependencyIdx := 0
	for i, line := range lines {
		if strings.HasPrefix(line, "- ") {
			buffer = []string{}
			orderedBuffer := []string{}
			buffer = append(buffer, strings.TrimLeft(line, "- "))
			j := i + 1
			for j < len(lines) && !strings.HasPrefix(lines[j], "- ") {
				buffer = append(buffer, strings.TrimSpace(lines[j]))
				j += 1
			}
			for _, attr := range attrOrders[dependencyIdx] {
				for _, b := range buffer {
					bb := strings.TrimLeft(b, "- ")
					bb = strings.TrimSpace(bb)
					if strings.HasPrefix(bb, attr) {
						indent := "  "
						if len(orderedBuffer) == 0 {
							indent = "- "
						}
						orderedBuffer = append(orderedBuffer, indent+bb)
						break
					}
				}
			}
			orderedBuffer = append(orderedBuffer, "")
			linesUpdatedContentNew = append(linesUpdatedContentNew, orderedBuffer...)
			dependencyIdx += 1
		}
	}

	if appendDashes {
		linesUpdatedContentNew = append([]string{"---"}, linesUpdatedContentNew...)
	}
	updatedContent = strings.Join(linesUpdatedContentNew, "\n")
	return updatedContent
}

func (a *AnsibleUpdater) updateProjectDependencyAndCommitChanges(c *types.ProjectDependencies, gitWorkingBranchName string, commitMessage string) error {
	countAnsibleDependencyUpdates := 0
	for _, dependencyUpdate := range c.Dependencies {
		if dependencyUpdate.Type == a.Manager.GetPackageManagerName() {
			countAnsibleDependencyUpdates += 1
		}
	}
	if countAnsibleDependencyUpdates == 0 {
		return errors.NewOperationError(nil, "No dependency update found")
	}

	err := cloneRepoAndCommitChanges(c.Project.HTTPURLToRepo, c.Project.ID, gitWorkingBranchName, commitMessage, a, c.Dependencies)
	if err != nil {
		return err
	}

	return nil
}
