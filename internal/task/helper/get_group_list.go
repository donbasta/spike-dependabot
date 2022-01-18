package helper

import (
	"dependabot/internal/config"
	"strconv"

	"github.com/gopaytech/go-commons/pkg/zlog"
)

func GetDefaultGroupList() []int {
	mainCfg := config.ProvideConfig()
	listID := mainCfg.Groups.ListID

	listIntID := []int{}
	for _, id := range listID {
		intId, err := strconv.Atoi(id)
		if err != nil {
			zlog.Info(err.Error())
			continue
		}

		listIntID = append(listIntID, intId)
	}

	return listIntID
}
