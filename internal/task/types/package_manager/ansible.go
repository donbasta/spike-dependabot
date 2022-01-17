package packageManager

import "path/filepath"

type ansiblePackageManager struct {
}

func (a *ansiblePackageManager) IsPackageDependencyRequirementFile(path string) bool {
	fileName := filepath.Base(path)
	return fileName == "requirements.yml" || fileName == "playbooks.yml" || fileName == "playbooks.yml.tmpl" || fileName == "requirements.yaml" || fileName == "playbooks.yaml" || fileName == "playbooks.yaml.tmpl"
}

func (a *ansiblePackageManager) GetPackageName() string {
	return "ansible"
}

func CreateAnsiblePackageManager() *ansiblePackageManager {
	return &ansiblePackageManager{}
}