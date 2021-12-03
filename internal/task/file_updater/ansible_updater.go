package updater

import (
	"dependabot/internal/errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gopaytech/go-commons/pkg/git"
	"github.com/gopaytech/go-commons/pkg/zlog"

	"github.com/go-git/go-git/v5/plumbing/object"
	gitTransport "github.com/go-git/go-git/v5/plumbing/transport"
)

type AnsibleUpdater struct {
	GitCloneFunc git.CloneFunc
	GitAuth      gitTransport.AuthMethod
}

func (a *AnsibleUpdater) UpdateDependency(c *Changes) error {
	if len(c.DepChanges) == 0 {
		return nil
	}

	zlog.Info("create temp folder as Template git clone target")
	repoCloneTarget, err := ioutil.TempDir("", "*")
	if err != nil {
		zlog.Error(err, "error while cloning repository")
		return err
	}

	defer os.RemoveAll(repoCloneTarget)
	repoUrl := c.Project.HTTPURLToRepo
	zlog.Info("clone template repository %s master to target %s", repoUrl, repoCloneTarget)
	gitRepository, err := a.GitCloneFunc(repoUrl, repoCloneTarget, a.GitAuth)
	if err != nil {
		zlog.Error(err, "clone template repository %s from master to target %s", repoUrl, repoCloneTarget)
		return err
	}

	testDependencyName := "gcloud-rabbitmq"
	gitWorkingBranchName := fmt.Sprintf("dependabot/ansible-terraform/%s-%s", testDependencyName, c.DepChanges[0].New)

	zlog.Info("create local repository branch %s ", gitWorkingBranchName)
	err = gitRepository.CreateBranch(gitWorkingBranchName)
	if err != nil {
		return errors.NewOperationError(err, "create local repository branch %s ", gitWorkingBranchName)
	}

	zlog.Info("checkout local repository branch %s ", gitWorkingBranchName)
	err = gitRepository.CheckoutBranch(gitWorkingBranchName)
	if err != nil {
		return errors.NewOperationError(err, "checkout local repository branch %s ", gitWorkingBranchName)
	}

	title := fmt.Sprintf("Bump: %s to %s", testDependencyName, c.DepChanges[0].New)
	zlog.Info("commit [%s]", title)

	signature := object.Signature{
		Name:  "Automated Commit from Dependabot",
		Email: "cloud-automation@gopay.co.id",
		When:  time.Now(),
	}
	_, err = gitRepository.AddAllAndCommit(title, &signature)
	if err != nil {
		return errors.NewOperationError(err, "commit  [%s]", title)
	}

	zlog.Info("push [%s]", title)
	err = gitRepository.PushDefault()
	if err != nil {
		return errors.NewOperationError(err, "push [%s]", title)
	}

	return nil
}
