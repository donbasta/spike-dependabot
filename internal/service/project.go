package service

import (
	"dependabot/internal/service/parser"

	gl "github.com/xanzy/go-gitlab"
)

type Project struct {
	project      *gl.Project
	dependencies []parser.Dependency
}
