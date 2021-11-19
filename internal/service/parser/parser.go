package parser

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"

	gl "github.com/xanzy/go-gitlab"
)

type PackageType string

const (
	Ansible   PackageType = "ansible"
	Terraform PackageType = "terraform"
)

type Semver struct {
	major, minor, patch int
}

func (s Semver) String() string {
	return fmt.Sprintf("v%d.%d.%d", s.major, s.minor, s.patch)
}

type Dependency struct {
	Url     string
	Version Semver
}

type DependencyParser interface {
	Parse(fileContent string) ([]Dependency, error)
}

type dependencyParser struct {
	packageType PackageType
}

// func (d *dependencyParser) Parse(fileContent string) ([]Dependency, error) {
// 	return []Dependency{{Url: "github-test", Version: Semver{1, 2, 3}}, {Url: "bitbucket-test", Version: Semver{1, 0, 1}}}, nil
// }

func ParseProject(client *gl.Client, project *gl.Project) ([]Dependency, error) {
	an := dependencyParser{packageType: Ansible}
	tf := dependencyParser{packageType: Terraform}

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
	dependencies := []Dependency{}
	for j := 0; j < len(t); j++ {
		ptrString := func(s string) *string {
			return &s
		}
		fileOptions := &gl.GetFileOptions{
			Ref: ptrString("master"),
		}
		file, _, _ := client.RepositoryFiles.GetFile(id, url.QueryEscape(t[j].Path), fileOptions)
		if t[j].Name == "requirements.yml" {
			content, _ := base64.StdEncoding.DecodeString(file.Content)
			dep, _ := an.Parse(string(content))
			dependencies = append(dependencies, dep...)
		} else if t[j].Name == "main.tf" {
			content, _ := base64.StdEncoding.DecodeString(file.Content)
			dep, _ := tf.Parse(string(content))
			dependencies = append(dependencies, dep...)
		}
	}

	return dependencies, nil
}
