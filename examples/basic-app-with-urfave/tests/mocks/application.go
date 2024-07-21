package mocks

import (
	"github.com/componego/componego"

	"github.com/componego/urfave-cli-integration/examples/basic-app-with-urfave/internal/application"
	"github.com/componego/urfave-cli-integration/examples/basic-app-with-urfave/internal/config"
)

type ApplicationMock struct {
	*application.Application
}

func NewApplicationMock() *ApplicationMock {
	return &ApplicationMock{
		Application: application.New(),
	}
}

func (a *ApplicationMock) ApplicationConfigInit(_ componego.ApplicationMode, _ any) (map[string]any, error) {
	return config.Read("config.test.json")
}
