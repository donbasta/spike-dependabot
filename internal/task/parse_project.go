package task

import (
	parser "dependabot/internal/file_parser"
	"encoding/base64"
	"log"

	gl "github.com/xanzy/go-gitlab"
)

func ParseProject(client *gl.Client, project *gl.Project) ([]parser.Dependency, error) {
	an := parser.AnsibleParser{}
	tf := parser.TerraformParser{}

	listOpt := gl.ListOptions{
		Page:    0,
		PerPage: 100,
	}

	newTrue := func(b bool) *bool {
		return &b
	}

	opt := &gl.ListTreeOptions{ListOptions: listOpt, Recursive: newTrue(true)}
	id := project.ID
	tree, _, _ := client.Repositories.ListTree(id, opt)

	log.Println(project.Name)
	dependencies := []parser.Dependency{}
	for j := 0; j < len(tree); j++ {
		ptrString := func(s string) *string {
			return &s
		}
		fileOptions := &gl.GetFileOptions{
			Ref: ptrString("master"),
		}
		file, _, _ := client.RepositoryFiles.GetFile(id, tree[j].Path, fileOptions)
		fileName := tree[j].Name
		filePath := tree[j].Path
		if fileName == "requirements.yml" || fileName == "playbooks.yml" || fileName == "playbooks.yml.tmpl" {
			content, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				log.Fatal(err)
			}
			dep, _ := an.Parse(string(content))
			dependencies = append(dependencies, dep...)
		}
		if (fileName == "main.tf" || fileName == "main.tf.tmpl") && (filePath != "examples/main.tf") {
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
