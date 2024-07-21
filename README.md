# ComponeGo Framework + Urfave CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/componego/urfave-cli-integration)](https://goreportcard.com/report/github.com/componego/urfave-cli-integration)
[![Tests](https://github.com/componego/urfave-cli-integration/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/componego/urfave-cli-integration/actions/workflows/tests.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/componego/urfave-cli-integration.svg)](https://pkg.go.dev/github.com/componego/urfave-cli-integration)

This package provides easy integration of the [Componego Framework](https://github.com/componego/componego) with [Urfave CLI V2](https://github.com/urfave/cli).

All [Urfave CLI](https://github.com/urfave/cli) features are supported.
You can use a convenient solution for processing command line arguments with the [Componego Framework](https://github.com/componego/componego) together.

### Documentation

The application must be launched in a special way for this to work.

#### main.go
```go
package main

import (
    "github.com/componego/componego"
    "github.com/componego/urfave-cli-integration"

    "github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/internal/application"
)

func main() {
    urfave_cli_integration.RunAndExit(application.New(), componego.ProductionMode)
}
```
This is necessary because the current integration may launch different [componego applications](https://github.com/componego/componego) depending on the command line arguments.


Inside the application, you can access the [Urfave CLI](https://github.com/urfave/cli) context.
You can also modify any command options. For example, add a couple of required flags for the command line.

#### application.go
```go
package application

import (
    "github.com/componego/componego"
    "github.com/componego/urfave-cli-integration"
    "github.com/urfave/cli/v2"
)

type Application struct {
}

func New() *Application {
    return &Application{}
}

// ApplicationName belongs to interface componego.Application.
func (a *Application) ApplicationName() string {
    return "Application Name v0.0.1"
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
    // ...
}

// ApplicationConfigInit belongs to interface componego.ApplicationConfigInit.
func (a *Application) ApplicationConfigInit(_ componego.ApplicationMode, options any) (map[string]any, error) {
    cliCtx := options.(*cli.Context) // <----
    // ...
}

// ApplicationAction belongs to interface componego.Application.
func (a *Application) ApplicationAction(env componego.Environment, options any) (int, error) {
    cliCtx := options.(*cli.Context) // <----
    // ...
}

var (
    _ componego.Application                 = (*Application)(nil)
    _ componego.ApplicationConfigInit       = (*Application)(nil)
    _ urfave_cli_integration.ApplicationCLI = (*Application)(nil)
)
```

Subcommands can be added like this:

```go
// ApplicationCLI belongs to interface urfave_cli_integration.ApplicationCLI.
func (a *Application) ApplicationCLI(cmd *cli.Command) {
    cmd.Subcommands = []*cli.Command{
        urfave_cli_integration.ToCommand("subcommand", application.New()),
    }
    // ....
}
```

Examples available [here](./examples).

### Contributing

We are open to improvements and suggestions. Pull requests are welcome.

### License

The source code of the repository is licensed under the [MIT license](./LICENSE).
