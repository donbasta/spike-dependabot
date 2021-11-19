package dependabot

import (
	"dependabot/internal/client"
	"dependabot/internal/service"
	"dependabot/internal/service/parser"
	"log"
)

func Run() {
	client := client.CreateClient("_anDFoWTfDoysp8RUfap")
	projects, err := service.CrawlGroup(client)
	if err != nil {
		log.Println("error")
		return
	}
	size := len(projects)
	for i := 0; i < size; i++ {
		if projects[i].Name != "rabbitmq-playbook" {
			continue
		}
		deps, _ := parser.ParseProject(client, projects[i])
		for j := 0; j < len(deps); j++ {
			log.Println(deps[j].Url, " ", deps[j].Version)
		}
	}
}
