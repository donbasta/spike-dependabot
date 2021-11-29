package dependabot

import (
	"dependabot/internal/client"
	"dependabot/internal/service"
	"dependabot/internal/service/parser"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading environment file")
	}
	secretKey := os.Getenv("SECRET_KEY")

	client := client.CreateClient(secretKey)
	projects, err := service.CrawlGroup(client)
	if err != nil {
		log.Println("error")
		return
	}
	size := len(projects)
	for i := 0; i < size; i++ {
		if projects[i].Name != "aws-basic-instance" {
			continue
		}
		deps, _ := parser.ParseProject(client, projects[i])
		for j := 0; j < len(deps); j++ {
			log.Println(deps[j].Url, " ", deps[j].Version)
		}
	}
}
