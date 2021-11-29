package main

import (
	"dependabot/cmd/dependabot"
	"dependabot/internal/config"
)

func main() {
	config.Config()

	dependabot.Run()
}
