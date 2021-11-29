package service

import (
	"dependabot/internal/service/parser"
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
	t, _, _ := client.Repositories.ListTree(id, opt)

	log.Println(project.Name)
	dependencies := []parser.Dependency{}
	for j := 0; j < len(t); j++ {
		ptrString := func(s string) *string {
			return &s
		}
		fileOptions := &gl.GetFileOptions{
			Ref: ptrString("master"),
		}
		file, _, _ := client.RepositoryFiles.GetFile(id, t[j].Path, fileOptions)
		if t[j].Name == "requirements.yml" && t[j].Path == "requirements.yml" {
			content, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				log.Fatal(err)
			}
			dep, _ := an.Parse(string(content))
			dependencies = append(dependencies, dep...)
		}
		if t[j].Name == "main.tf" && t[j].Path == "main.tf" {
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
