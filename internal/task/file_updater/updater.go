package updater

import (
	gl "github.com/xanzy/go-gitlab"
)

type Updater interface {
	UpdateDependency(fileChanges *Changes) error
	SubmitMergeRequest() error
}

type DependencyChange struct {
	Source   string
	Old, New string
}

type Changes struct {
	Project    *gl.Project
	DepChanges []DependencyChange
}
