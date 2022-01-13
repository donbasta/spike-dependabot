package updater

import (
	"dependabot/internal/task/types"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformUpdaterNoDependencyUpdate(t *testing.T) {
	updater := CreateTerraformUpdater()
	projectUpdate := types.ProjectDependencies{
		Project:      nil,
		Dependencies: []types.Dependency{},
	}
	err := updater.updateProjectDependencyAndCommitChanges(&projectUpdate, "", "")
	assert.Error(t, err)
}

func TestTerraformUpdateFile(t *testing.T) {
	updater := CreateTerraformUpdater()

	expectedNewFileRaw, _ := ioutil.ReadFile("./test/terraform_new.tf")
	expectedNewFile := string(expectedNewFileRaw)

	oldFileRaw, _ := ioutil.ReadFile("./test/terraform_old.tf")
	oldFile := string(oldFileRaw)

	v1, _ := types.MakeVersion("v5.0.20")
	v2, _ := types.MakeVersion("v4.1.2")
	dependencies := []types.Dependency{
		{
			SourceBaseUrl: "source.golabs.io/gopay_infra/terraform/aws-basic-instance",
			Version:       *v1,
			Type:          "terraform",
		},
		{
			SourceBaseUrl: "source.golabs.io/gopay_infra/terraform/aws-redis",
			Version:       *v2,
			Type:          "terraform",
		},
	}
	updatedOldFile := updater.updateContentWithNewDependency(oldFile, dependencies[0])
	updatedOldFile = updater.updateContentWithNewDependency(updatedOldFile, dependencies[1])

	assert.Equal(t, expectedNewFile, updatedOldFile)
}

func TestTerraformPackageName(t *testing.T) {
	updater := CreateTerraformUpdater()

	path := "/example/main.tf"
	assert.False(t, updater.IsPackageDependencyRequirementFile(path))

	path = "/proj/dummy/template-scp/terraform/main.tf"
	assert.True(t, updater.IsPackageDependencyRequirementFile(path))
}
