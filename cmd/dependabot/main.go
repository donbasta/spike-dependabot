package dependabot

import (
	"dependabot/internal/client"
	"dependabot/internal/service"
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

	service.CrawlGroups(client)
}
