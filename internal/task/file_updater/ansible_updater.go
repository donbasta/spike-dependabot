package updater

import (
	"io/ioutil"
	"os"

	"github.com/gopaytech/go-commons/pkg/git"
	"github.com/gopaytech/go-commons/pkg/zlog"

	gitTransport "github.com/go-git/go-git/v5/plumbing/transport"
)

type AnsibleUpdater struct {
	GitCloneFunc git.CloneFunc
	GitAuth      gitTransport.AuthMethod
}

func (a *AnsibleUpdater) Update(c *Changes) error {
	zlog.Info("create temp folder as Template git clone target")
	repoCloneTarget, err := ioutil.TempDir("", "*")
	if err != nil {
		zlog.Error(err, "create temp folder as Template git clone target")
		return err
	}

	defer os.RemoveAll(repoCloneTarget)
	repoUrl := c.Project.HTTPURLToRepo
	zlog.Info("clone template repository %s master to target %s", repoUrl, repoCloneTarget)
	_, err = a.GitCloneFunc(repoUrl, repoCloneTarget, a.GitAuth)
	if err != nil {
		zlog.Error(err, "clone template repository %s from master to target %s", repoUrl, repoCloneTarget)
		return err
	}

	return nil
}
