package urfave_cli_integration

import (
	"testing"

	"github.com/componego/componego"
	"github.com/componego/componego/impl/application"
	"github.com/urfave/cli/v2"
)

func TestRunApplication(t *testing.T) {
	currentArgs := cliArgs
	t.Cleanup(func() {
		cliArgs = currentArgs
	})
	cliArgs = []string{"application"}

	t.Run("basic test", func(t *testing.T) {
		expectedExitCode := 123
		called := false
		appFactory := application.NewFactory("Test Application")
		appFactory.SetApplicationAction(func(_ componego.Environment, options any) (int, error) {
			if _, ok := options.(*cli.Context); !ok {
				t.Fatal("context not found")
			}
			called = true
			return expectedExitCode, nil
		})
		actualExitCode := Run(appFactory.Build(), componego.TestMode)
		if expectedExitCode != actualExitCode {
			t.Fatal("wrong exit code")
		}
		if !called {
			t.Fatal("the required function was not called")
		}
	})

	t.Run("command line configuration", func(t *testing.T) {
		called := false
		app := &testApplication{
			applicationCLIFunc: func(cmd *cli.Command) {
				cmd.Usage = "usage text"
			},
			applicationAction: func(_ componego.Environment, options any) (int, error) {
				ctxCli, ok := options.(*cli.Context)
				if !ok {
					t.Fatal("context not found")
				}
				if ctxCli.Command.Usage != "usage text" {
					t.Fatal("data doesn't match")
				}
				called = true
				return componego.SuccessExitCode, nil
			},
		}
		if Run(app, componego.TestMode) != componego.SuccessExitCode {
			t.Fatal("wrong exit code")
		}
		if !called {
			t.Fatal("the required function was not called")
		}
	})
}

type testApplication struct {
	applicationCLIFunc func(cmd *cli.Command)
	applicationAction  func(env componego.Environment, options any) (int, error)
}

// ApplicationName belongs to interface componego.Application.
func (t *testApplication) ApplicationName() string {
	return "Test App v0.0.1"
}

// ApplicationCLI belongs to interface ApplicationCLI.
func (t *testApplication) ApplicationCLI(cmd *cli.Command) {
	t.applicationCLIFunc(cmd)
}

// ApplicationAction belongs to interface componego.Application.
func (t *testApplication) ApplicationAction(env componego.Environment, options any) (int, error) {
	return t.applicationAction(env, options)
}

var (
	_ componego.Application = (*testApplication)(nil)
	_ ApplicationCLI        = (*testApplication)(nil)
)
