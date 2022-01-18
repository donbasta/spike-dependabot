package packageManager

import "path/filepath"

type terraformPackageManager struct {
}

func (t *terraformPackageManager) IsPackageDependencyRequirementFile(path string) bool {
	fileName := filepath.Base(path)
	fileDir := filepath.Base(filepath.Dir(path))
	return (fileName == "main.tf" || fileName == "main.tf.tmpl") && (fileDir != "examples") && (fileDir != "example")
}

func (t *terraformPackageManager) GetPackageManagerName() string {
	return "terraform"
}

func CreateTerraformPackageManager() *terraformPackageManager {
	return &terraformPackageManager{}
}
