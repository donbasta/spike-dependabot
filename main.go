package main

import (
	"dependabot/cmd/dependabot"
	"dependabot/internal/config"
	"log"
)

func main() {
	config.Config()

	err := dependabot.NewRootCommand().Execute()
	if err != nil {
		log.Fatal(err)
	}
}
