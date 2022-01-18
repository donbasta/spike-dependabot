package config

import (
	"context"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/gopaytech/go-commons/pkg/atom"
	"github.com/gopaytech/go-commons/pkg/zlog"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Main struct {
	Database Database
	Git      Git
	Slack    Slack
	Groups   Groups
}

type Git struct {
	Token    string `env:"GIT_TOKEN,required"`
	URL      string `env:"GIT_URL,required"`
	Username string `env:"GIT_USERNAME,required"`
	Email    string `env:"GIT_EMAIL,required"`
}

type Database struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
}

type Slack struct {
	Token     string `env:"SLACK_TOKEN"`
	ChannelId string `env:"SLACK_CHANNEL_ID"`
}

type Groups struct {
	ListID []string `env:"LIST_GROUP_ID"`
}

var (
	once       sync.Once
	mainConfig *Main
)

const projectDirName = "scp-dependency-manager"

func loadFromEnvFile() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		zlog.Info("Error while loading .env file")
	}
}

func load() (*Main, error) {
	main := Main{}
	ctx := context.Background()

	loadFromEnvFile()

	err := envconfig.Process(ctx, &main)
	return &main, err
}

var reset = new(atom.Bool)

func init() {
	reset.Set(false)
}

func Reset() {
	reset.Set(true)
}

func Config() *Main {
	once.Do(func() {
		config, err := load()
		if err != nil {
			log.Fatalf("Processing config failed %s", err)
		}

		mainConfig = config
	})

	if reset.Get() {
		config, err := load()
		if err != nil {
			log.Fatalf("Processing config failed %s", err)
		}

		mainConfig = config
		reset.Set(false)
	}

	return mainConfig
}
