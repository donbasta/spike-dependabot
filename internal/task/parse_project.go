package task

import (
	parser "dependabot/internal/file_parser"
	"encoding/base64"
	"log"

	gl "github.com/xanzy/go-gitlab"
)

func ansibleDependencyFile(file *gl.TreeNode) bool {
	fileName := file.Name
	return fileName == "requirements.yml" || fileName == "playbooks.yml" || fileName == "playbooks.yml.tmpl"
}

func terraformDependencyFile(file *gl.TreeNode) bool {
	fileName := file.Name
	filePath := file.Path
	return (fileName == "main.tf" || fileName == "main.tf.tmpl") && (filePath != "examples/main.tf")
}

func ParseProject(client *gl.Client, project *gl.Project) ([]parser.Dependency, error) {
	an := parser.AnsibleParser{}
	tf := parser.TerraformParser{}

	listOpt := gl.ListOptions{
		Page:    0,
		PerPage: 100,
	}

	ptrBool := func(b bool) *bool {
		return &b
	}
	ptrString := func(s string) *string {
		return &s
	}

	opt := &gl.ListTreeOptions{ListOptions: listOpt, Recursive: ptrBool(true)}
	id := project.ID
	tree, _, _ := client.Repositories.ListTree(id, opt)

	log.Println(project.Name)
	dependencies := []parser.Dependency{}
	for _, child := range tree {
		fileOptions := &gl.GetFileOptions{
			Ref: ptrString("master"),
		}
		file, _, _ := client.RepositoryFiles.GetFile(id, child.Path, fileOptions)
		if ansibleDependencyFile(child) {
			content, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				log.Fatal(err)
			}
			dep, _ := an.Parse(string(content))
			dependencies = append(dependencies, dep...)
		}
		if terraformDependencyFile(child) {
			content, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				log.Fatal(err)
			}
			dep, _ := tf.Parse(string(content))
			dependencies = append(dependencies, dep...)
		}
	}

	return dependencies, nil
}
