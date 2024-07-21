package application

import (
	"fmt"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/urfave-cli-integration"
	"github.com/urfave/cli/v2"

	"github.com/componego/urfave-cli-integration/examples/basic-app-with-urfave/internal/config"
)

type Application struct{}

func New() *Application {
	return &Application{}
}

// ApplicationName belongs to interface componego.Application.
func (a *Application) ApplicationName() string {
	return "Urfave CLI Integration Example v0.0.1"
}

// ApplicationCLI belongs to interface urfave_cli_integration.ApplicationCLI.
func (a *Application) ApplicationCLI(cmd *cli.Command) {
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Usage:    "config filename",
			Required: true,
		},
	}
}

// ApplicationConfigInit belongs to interface componego.ApplicationConfigInit.
func (a *Application) ApplicationConfigInit(_ componego.ApplicationMode, options any) (map[string]any, error) {
	// Options are *cli.Context because the application is running inside urfave/cli.
	// See main function.
	cliCtx := options.(*cli.Context)
	return config.Read(cliCtx.String("config"))
}

// ApplicationAction belongs to interface componego.Application.
func (a *Application) ApplicationAction(env componego.Environment, _ any) (int, error) {
	configMessage, err := env.ConfigProvider().ConfigValue("hello.message", nil)
	if err == nil {
		_, err = fmt.Fprintln(env.ApplicationIO().OutputWriter(), configMessage)
	}
	return application.ExitWrapper(err)
}

var (
	_ componego.Application                 = (*Application)(nil)
	_ componego.ApplicationConfigInit       = (*Application)(nil)
	_ urfave_cli_integration.ApplicationCLI = (*Application)(nil)
)
