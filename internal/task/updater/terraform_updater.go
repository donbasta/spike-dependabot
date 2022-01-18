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

	gitTransport "github.com/go-git/go-git/v5/plumbing/transport"
)

type TerraformUpdater struct {
	GitCloneFunc git.CloneFunc
	GitAuth      gitTransport.AuthMethod
	Parser       parser.DependencyParser
	Manager      packageManager.PackageManager
}

func CreateTerraformUpdater() *TerraformUpdater {
	mainCfg := config.ProvideConfig()
	authMethod := config.ProvideGitAuth(mainCfg)
	return &TerraformUpdater{
		GitCloneFunc: git.Clone,
		GitAuth:      authMethod,
		Parser:       parser.CreateTerraformParser(),
		Manager:      packageManager.CreateTerraformPackageManager(),
	}
}

func (tf *TerraformUpdater) IsPackageDependencyRequirementFile(filepath string) bool {
	return tf.Manager.IsPackageDependencyRequirementFile(filepath)
}

func (tf *TerraformUpdater) GetPackageManagerName() string {
	return tf.Manager.GetPackageManagerName()
}

func (tf *TerraformUpdater) ProcessUpdateProjectDependencies(c *types.ProjectDependencies, s *service.SlackNotificationService, m *service.MergeRequestService) error {
	dt := time.Now()
	gitWorkingBranchName := fmt.Sprintf("scp-dependency-manager-bump/%s/%s", tf.Manager.GetPackageManagerName(), dt.Format("01-02-2006"))
	mergeRequestTitle := fmt.Sprintf("fix: bump %s dependencies (%s)", tf.Manager.GetPackageManagerName(), dt.Format("01-02-2006"))

	err := tf.updateProjectDependencyAndCommitChanges(c, gitWorkingBranchName, mergeRequestTitle)
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
	err = (*s).NotifyMerge(mainCfg.Slack.ChannelId, mergeRequestURL, c.Project.Name, tf.Manager.GetPackageManagerName())
	if err != nil {
		return err
	}

	return nil
}

func (tf *TerraformUpdater) updateLine(line string, c types.Dependency) (string, error) {
	trimmedLine := strings.TrimLeft(line, "")
	url, err := tf.Parser.GetBaseUrlFromRawSource(trimmedLine)
	if err != nil {
		return line, err
	}

	tokens := strings.Split(url, "?")
	baseUrl := tokens[0]
	baseUrl = strings.TrimRight(baseUrl, ".git")

	if baseUrl != c.SourceBaseUrl {
		return line, nil
	}

	if len(tokens) == 1 {
		return line, nil
	}

	separator := " = "
	lineTokens := strings.Split(line, separator)
	source := lineTokens[1]
	tokens = strings.Split(source, "?")
	tokens[1] = "ref=" + c.Version.String() + "\""
	lineTokens[1] = strings.Join(tokens, "?")
	updatedLine := strings.Join(lineTokens, separator)

	return updatedLine, nil
}

func (tf *TerraformUpdater) updateContentWithNewDependency(fileContent string, dependency types.Dependency) string {
	fileContentLines := strings.Split(fileContent, "\n")
	isModule := false
	tmpNewContent := []string{}

	for _, line := range fileContentLines {
		updatedLine := line
		trimmedLine := strings.Trim(line, " ")
		if len(trimmedLine) == 0 {
			tmpNewContent = append(tmpNewContent, "")
			continue
		}

		tokens := strings.Fields(line)
		if tokens[0] == "module" {
			isModule = true
		}

		attr := tokens[0]
		if isModule && (attr == "source") {
			updatedLine, _ = tf.updateLine(line, dependency)
			isModule = false
		}

		tmpNewContent = append(tmpNewContent, updatedLine)
	}

	newFileContent := strings.Join(tmpNewContent, "\n")
	return newFileContent
}

func (tf *TerraformUpdater) updateProjectDependencyAndCommitChanges(c *types.ProjectDependencies, gitWorkingBranchName string, commitMessage string) error {
	countTerraformDependencyUpdates := 0
	for _, dependencyUpdate := range c.Dependencies {
		if dependencyUpdate.Type == tf.Manager.GetPackageManagerName() {
			countTerraformDependencyUpdates += 1
		}
	}

	if countTerraformDependencyUpdates == 0 {
		return errors.NewOperationError(nil, "No dependency update found")
	}

	err := cloneRepoAndCommitChanges(c.Project.HTTPURLToRepo, c.Project.ID, gitWorkingBranchName, commitMessage, tf, c.Dependencies)
	if err != nil {
		return err
	}

	return nil
}
