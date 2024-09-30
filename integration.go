package urfave_cli_integration

import (
	"context"
	"fmt"
	"os"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/componego/componego/impl/driver"
	"github.com/componego/componego/impl/runner/unhandled-errors"
	"github.com/componego/urfave-cli-integration/internal"
	"github.com/urfave/cli/v2"
)

const (
	appModeKey  = "componego:app:mode"
	exitCodeKey = "componego:app:exitCode"
)

// noinspection ALL
var (
	ToAction      = toAction
	ToApplication = toApplication
	ToCommand     = toCommand

	cliArgs = os.Args
)

type ApplicationCLI interface {
	componego.Application
	ApplicationCLI(cmd *cli.Command)
}

// RunWithContext runs the application with context and returns the exit code.
func RunWithContext(ctx context.Context, app componego.Application, appMode componego.ApplicationMode) int {
	cliApp := ToApplication(app)
	internal.SetMetadata(cliApp, appModeKey, appMode)
	err := cliApp.RunContext(ctx, cliArgs)
	if err == nil {
		return internal.GetMetadata[int](cliApp, exitCodeKey, componego.SuccessExitCode)
	}
	_, err = fmt.Fprintln(cliApp.ErrWriter, "\nERROR:", err.Error())
	if err != nil {
		panic(err)
	}
	exitCode := internal.GetMetadata[int](cliApp, exitCodeKey, componego.ErrorExitCode)
	if exitCode == componego.SuccessExitCode {
		return componego.ErrorExitCode
	}
	return exitCode
}

// Run runs the application and returns the exit code.
func Run(app componego.Application, appMode componego.ApplicationMode) int {
	return RunWithContext(context.Background(), app, appMode)
}

// RunAndExit runs the application and exits the program after stopping the application.
func RunAndExit(app componego.Application, appMode componego.ApplicationMode) {
	cli.OsExiter(Run(app, appMode))
}

func toAction(app componego.Application) cli.ActionFunc {
	return func(cliCtx *cli.Context) error {
		cliApp := cliCtx.App
		appMode := internal.GetMetadata[componego.ApplicationMode](cliApp, appModeKey, componego.ProductionMode)
		exitCode, err := driver.New(&driver.Options{
			AppIO:      application.NewIO(cliApp.Reader, cliApp.Writer, cliApp.ErrWriter),
			Additional: cliCtx,
		}).RunApplication(cliCtx.Context, app, appMode)
		if err != nil {
			errMessage := unhandled_errors.ToString(err, appMode, unhandled_errors.GetHandlers())
			_, err = fmt.Fprintln(cliApp.ErrWriter, errMessage)
		}
		internal.SetMetadata(cliApp, exitCodeKey, exitCode)
		return err
	}
}

func toApplication(app componego.Application) *cli.App {
	cliApp := cli.NewApp()
	if app, ok := app.(ApplicationCLI); ok {
		command := &cli.Command{}
		app.ApplicationCLI(command)
		cliApp.Name = command.Name
		cliApp.Usage = command.Usage
		cliApp.UsageText = command.UsageText
		cliApp.Description = command.Description
		cliApp.ArgsUsage = command.ArgsUsage
		cliApp.BashComplete = command.BashComplete
		cliApp.Before = command.Before
		cliApp.After = command.After
		cliApp.Action = command.Action
		cliApp.OnUsageError = command.OnUsageError
		cliApp.Commands = command.Subcommands
		cliApp.Flags = command.Flags
		cliApp.HideHelp = command.HideHelp
		cliApp.HideHelpCommand = command.HideHelpCommand
		cliApp.UseShortOptionHandling = command.UseShortOptionHandling
		cliApp.HelpName = command.HelpName
		cliApp.CustomAppHelpTemplate = command.CustomHelpTemplate
		cliApp.SkipFlagParsing = command.SkipFlagParsing
	}
	if cliApp.Usage == "" {
		cliApp.Usage = app.ApplicationName()
	}
	onUsageErrorOrig := cliApp.OnUsageError
	exitErrHandlerOrig := cliApp.ExitErrHandler
	cliApp.OnUsageError = func(cliCtx *cli.Context, err error, isSubcommand bool) error {
		internal.SetMetadata(cliCtx.App, exitCodeKey, componego.ErrorExitCode)
		if onUsageErrorOrig == nil {
			return onUsageError(cliCtx, err, isSubcommand)
		}
		return onUsageErrorOrig(cliCtx, err, isSubcommand)
	}
	cliApp.ExitErrHandler = func(cliCtx *cli.Context, err error) {
		exitCode := internal.GetMetadata[int](cliCtx.App, exitCodeKey, componego.SuccessExitCode)
		if exitCode == componego.SuccessExitCode && err != nil {
			internal.SetMetadata(cliCtx.App, exitCodeKey, componego.ErrorExitCode)
		}
		if exitErrHandlerOrig != nil {
			exitErrHandlerOrig(cliCtx, err)
		}
	}
	cliApp.Action = ToAction(app)
	return cliApp
}

func toCommand(commandName string, app componego.Application) *cli.Command {
	command := &cli.Command{
		Name:        commandName,
		Description: app.ApplicationName(),
		Usage:       app.ApplicationName(),
	}
	if app, ok := app.(ApplicationCLI); ok {
		app.ApplicationCLI(command)
	}
	command.Action = ToAction(app)
	return command
}

func onUsageError(cliCtx *cli.Context, err error, isSubcommand bool) error {
	_, err = fmt.Fprintln(cliCtx.App.ErrWriter, "Incorrect usage:", err.Error())
	if cliCtx.App.HideHelp || err != nil {
		return err
	}
	if _, err = fmt.Fprintln(cliCtx.App.ErrWriter); err != nil {
		return err
	}
	if isSubcommand {
		return cli.ShowCommandHelp(cliCtx, cliCtx.App.Name)
	}
	return cli.ShowAppHelp(cliCtx)
}
