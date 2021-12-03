package updater

import (
	"dependabot/internal/config"
	"log"

	"github.com/gopaytech/go-commons/pkg/git"
)

func UpdateProjects(changes []Changes) {
	mainCfg := config.ProvideConfig()
	authMethod := config.ProvideGitAuth(mainCfg)
	ansibleUpdater := &AnsibleUpdater{
		GitCloneFunc: git.Clone,
		GitAuth:      authMethod,
	}
	terraformUpdater := &TerraformUpdater{
		GitCloneFunc: git.Clone,
		GitAuth:      authMethod,
	}

	for i := 0; i < len(changes); i++ {
		err := ansibleUpdater.UpdateDependency(&changes[i])
		if err != nil {
			log.Fatal(err)
		}
		err = terraformUpdater.UpdateDependency(&changes[i])
		if err != nil {
			log.Fatal(err)
		}
	}
}
