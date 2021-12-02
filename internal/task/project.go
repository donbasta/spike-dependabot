package task

import (
	parser "dependabot/internal/file_parser"

	gl "github.com/xanzy/go-gitlab"
)

type Project struct {
	project      *gl.Project
	dependencies []parser.Dependency
}
