package main

import (
	"github.com/componego/componego"
	"github.com/componego/urfave-cli-integration"

	"github.com/componego/urfave-cli-integration/examples/basic-app-with-urfave/internal/application"
)

func main() {
	urfave_cli_integration.RunAndExit(application.New(), componego.DeveloperMode)
}
