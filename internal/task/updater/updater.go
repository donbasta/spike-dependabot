package updater

import (
	"dependabot/internal/config"
	"dependabot/internal/db"
	"dependabot/internal/db/repository"
	"dependabot/internal/errors"
	"dependabot/internal/service"
	"dependabot/internal/task/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gopaytech/go-commons/pkg/git"
	"github.com/gopaytech/go-commons/pkg/gitlab"
	"github.com/gopaytech/go-commons/pkg/zlog"
)

func UpdateProjectDependencies(dependencyUpdates []types.ProjectDependencies) {
	mainConfig := config.ProvideConfig()

	gormDB, err := db.ProvideDB(mainConfig)
	if err != nil {
		log.Println(err)
	}

	packageUpdaters := []PackageUpdater{
		CreateAnsibleUpdater(),
		CreateTerraformUpdater(),
	}

	mergeRequestRepository := repository.NewMergeRequestRepository()
	mergeRequestService := service.NewMergeRequestService(gormDB, mergeRequestRepository)

	slackClient := config.ProvideSlackClient(mainConfig)
	slackNotificationService := service.NewSlackNotificationService(mainConfig, slackClient)

	client, _ := config.ProvideGitlabClient(mainConfig)
	mergeRequestClient := gitlab.NewMergeRequest(client)

	for _, dependencyUpdate := range dependencyUpdates {
		repositoryId := strconv.Itoa(dependencyUpdate.Project.ID)
		mergeRequestRecords, err := mergeRequestService.GetAllByRepositoryId(repositoryId)
		if err != nil {
			continue
		}

		if len(*mergeRequestRecords) > 0 {
			for i := range *mergeRequestRecords {
				id := gitlab.NameOrId{ID: dependencyUpdate.Project.ID}
				mergeRequestIID, _ := strconv.Atoi((*mergeRequestRecords)[i].MergeRequestIID)

				mergeRequest, _ := mergeRequestClient.Get(id, mergeRequestIID)

				branchName := mergeRequest.SourceBranch
				_, deleteBranchError := client.Branches.DeleteBranch(id.Get(), branchName)
				if deleteBranchError != nil {
					zlog.Info("Branch cannot be deleted, aborted")
					continue
				}

				mergeRequestService.Delete(&(*mergeRequestRecords)[i])
			}
		}

		for i := range packageUpdaters {
			err = packageUpdaters[i].ProcessUpdateProjectDependencies(&dependencyUpdate, &slackNotificationService, &mergeRequestService)
			if err != nil {
				log.Printf("[%s] Update error status: %s\n", dependencyUpdate.Project.Name, err)
			}
		}
	}

	zlog.Info("Completed updating all the dependencies in the repo")
}

func cloneRepoAndCommitChanges(
	repoUrl string, projectId int, branchName string, commitMessage string, p PackageUpdater, dependencies []types.Dependency,
) error {
	repoCloneTarget, err := ioutil.TempDir("", "*")
	if err != nil {
		zlog.Error(err, "error while creating temporary directory for repository")
		return err
	}

	defer os.RemoveAll(repoCloneTarget)

	mainCfg := config.ProvideConfig()
	authMethod := config.ProvideGitAuth(mainCfg)
	gitRepository, err := git.Clone(repoUrl, repoCloneTarget, authMethod)
	if err != nil {
		zlog.Error(err, "clone template repository %s from master to target %s", repoUrl, repoCloneTarget)
		return err
	}

	err = gitRepository.CreateBranch(branchName)
	if err != nil {
		return errors.NewOperationError(err, "create local repository branch %s ", branchName)
	}

	err = gitRepository.CheckoutBranch(branchName)
	if err != nil {
		return errors.NewOperationError(err, "checkout local repository branch %s ", branchName)
	}

	err = filepath.Walk(repoCloneTarget,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if p.IsPackageDependencyRequirementFile(path) {
				rawContent, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				content := string(rawContent)
				for _, dependencyUpdate := range dependencies {
					content = p.updateContentWithNewDependency(content, dependencyUpdate)
				}

				err = ioutil.WriteFile(path, []byte(content), 0644)
				if err != nil {
					return err
				}
			}

			return nil
		})

	if err != nil {
		return errors.NewOperationError(err, "traversing the repo %s ", branchName)
	}

	signature := object.Signature{
		Name:  "Automated Commit from Dependabot",
		Email: "cloud-automation@gopay.co.id",
		When:  time.Now(),
	}

	_, err = gitRepository.AddAllAndCommit(commitMessage, &signature)
	if err != nil {
		return errors.NewOperationError(err, "commit  [%s]", commitMessage)
	}

	client, _ := config.ProvideGitlabClient(config.ProvideConfig())
	id := gitlab.NameOrId{ID: projectId}
	_, _ = client.Branches.DeleteBranch(id.Get(), branchName)

	err = gitRepository.PushDefault()
	if err != nil {
		return errors.NewOperationError(err, "push [%s]", commitMessage)
	}

	return nil
}
