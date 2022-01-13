package types

import (
	gl "github.com/xanzy/go-gitlab"
)

type ProjectDependencies struct {
	Project      *gl.Project
	Dependencies []Dependency
}

type Dependency struct {
	SourceRaw     string
	SourceBaseUrl string
	Version       SemanticVersion
	Type          string
}
