package packageManager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnsibleDependencyFilePath(t *testing.T) {
	manager := CreateAnsiblePackageManager()

	path := "/var/ansible-role/requirements.yml"
	shouldFileBeParsed := manager.IsPackageDependencyRequirementFile(path)
	assert.True(t, shouldFileBeParsed)

	path = "/var/ansible-role/requirements.yaml"
	shouldFileBeParsed = manager.IsPackageDependencyRequirementFile(path)
	assert.True(t, shouldFileBeParsed)

	path = "/var/tmp/ansible-playbook/ansible/module1/test.yaml"
	shouldFileBeParsed = manager.IsPackageDependencyRequirementFile(path)
	assert.False(t, shouldFileBeParsed)
}

func TestAnsiblePackageName(t *testing.T) {
	manager := CreateAnsiblePackageManager()
	assert.Equal(t, "ansible", manager.GetPackageName())
}
