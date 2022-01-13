package packageManager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformDependencyFilePath(t *testing.T) {
	manager := CreateTerraformPackageManager()

	filepath := "/example/main.tf"
	shouldFileBeParsed := manager.IsPackageDependencyRequirementFile(filepath)
	assert.False(t, shouldFileBeParsed)

	filepath = "/var/tmp/tf-module/terraform/module1/main.tf"
	shouldFileBeParsed = manager.IsPackageDependencyRequirementFile(filepath)
	assert.True(t, shouldFileBeParsed)
}

func TestTerraformPackageName(t *testing.T) {
	manager := CreateTerraformPackageManager()
	assert.Equal(t, "terraform", manager.GetPackageName())
}
