package dependabot

import "github.com/google/wire"

func InitializeServerCommand() (ServerCommand, error) {
	wire.Build(
		NewServerCommand,
	)
	return nil, nil
}
