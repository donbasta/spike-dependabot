package helper

import (
	"dependabot/internal/config"
	"log"
	"strconv"
)

func GetDefaultGroupList() []int {
	mainCfg := config.ProvideConfig()
	listID := mainCfg.Groups.ListID

	ret := []int{}
	for _, id := range listID {
		i, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			continue
		}

		ret = append(ret, i)
	}

	return ret
}
