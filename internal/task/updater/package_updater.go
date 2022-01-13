package updater

import (
	"dependabot/internal/service"
	"dependabot/internal/task/types"
)

type PackageUpdater interface {
	IsPackageDependencyRequirementFile(filepath string) bool
	ProcessUpdateProjectDependencies(c *types.ProjectDependencies, s *service.SlackNotificationService, m *service.MergeRequestService) error
	updateContentWithNewDependency(fileContent string, dependency types.Dependency) string
	updateProjectDependencyAndCommitChanges(c *types.ProjectDependencies, gitWorkingBranchName string, commitMessage string) error
}
