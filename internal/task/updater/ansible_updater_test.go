package updater

import (
	"dependabot/internal/task/types"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnsibleUpdaterNoDependencyUpdate(t *testing.T) {
	updater := CreateAnsibleUpdater()
	projectUpdate := types.ProjectDependencies{
		Project:      nil,
		Dependencies: []types.Dependency{},
	}
	err := updater.updateProjectDependencyAndCommitChanges(&projectUpdate, "", "")
	assert.Error(t, err)
}

func TestAnsibleUpdateFile(t *testing.T) {
	updater := CreateAnsibleUpdater()

	expectedNewFileRaw, _ := ioutil.ReadFile("./test/ansible_new.yaml")
	expectedNewFile := string(expectedNewFileRaw)

	oldFileRaw, _ := ioutil.ReadFile("./test/ansible_old.yaml")
	oldFile := string(oldFileRaw)

	v1, _ := types.MakeVersion("v2.0.0")
	dependencies := []types.Dependency{
		{
			SourceRaw:     "git@source.golabs.io:farras.f.interns/rabbitmq-role.git",
			SourceBaseUrl: "source.golabs.io/farras.f.interns/rabbitmq-role",
			Version:       *v1,
			Type:          "ansible",
		},
	}

	updatedOldFile := updater.updateContentWithNewDependency(oldFile, dependencies[0])

	assert.Equal(t, expectedNewFile, updatedOldFile)
}

func TestAnsibleUpdateFile2(t *testing.T) {
	updater := CreateAnsibleUpdater()

	expectedNewFileRaw, _ := ioutil.ReadFile("./test/ansible_new_2.yaml")
	expectedNewFile := string(expectedNewFileRaw)

	oldFileRaw, _ := ioutil.ReadFile("./test/ansible_old_2.yaml")
	oldFile := string(oldFileRaw)

	v1, _ := types.MakeVersion("v1.2.9")
	v2, _ := types.MakeVersion("v1.0.4")
	dependencies := []types.Dependency{
		{
			SourceRaw:     "https://source.golabs.io/gopay_infra/ansible/playbooks/redis-playbook.git",
			SourceBaseUrl: "source.golabs.io/gopay_infra/ansible/playbooks/redis-playbook",
			Version:       *v1,
			Type:          "ansible",
		},
		{
			SourceRaw:     "https://source.golabs.io/gopay_infra/ansible/playbooks/scp-playbook.git",
			SourceBaseUrl: "source.golabs.io/gopay_infra/ansible/playbooks/scp-playbook",
			Version:       *v2,
			Type:          "ansible",
		},
	}

	updatedOldFile := updater.updateContentWithNewDependency(oldFile, dependencies[0])
	updatedOldFile = updater.updateContentWithNewDependency(updatedOldFile, dependencies[1])

	assert.Equal(t, expectedNewFile, updatedOldFile)
}

func TestAnsiblePackageRequirementFileName(t *testing.T) {
	updater := CreateAnsibleUpdater()

	path := "/main/requirements.yml"
	assert.True(t, updater.IsPackageDependencyRequirementFile(path))

	path = "/proj/dummy/template-scp/ansible/inventories.yml"
	assert.False(t, updater.IsPackageDependencyRequirementFile(path))
}

func TestAnsiblePackageName(t *testing.T) {
	updater := CreateAnsibleUpdater()

	assert.Equal(t, "ansible", updater.GetPackageManagerName())
}
