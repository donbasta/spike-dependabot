package updater

import (
	gitTransport "github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/gopaytech/go-commons/pkg/git"
)

type TerraformUpdater struct {
	GitCloneFunc git.CloneFunc
	GitAuth      gitTransport.AuthMethod
}

func (t *TerraformUpdater) UpdateDependency(c *Changes) error {
	return nil
}
