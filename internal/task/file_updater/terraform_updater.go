package updater

import (
	gl "github.com/xanzy/go-gitlab"
)

func terraformDependencyFile(file *gl.TreeNode) bool {
	fileName := file.Name
	filePath := file.Path
	return (fileName == "main.tf" || fileName == "main.tf.tmpl") && (filePath != "examples/main.tf")
}

type TerraformUpdater struct {
}

func (t *TerraformUpdater) Update(fileChanges *Changes) error {
	return nil
}
