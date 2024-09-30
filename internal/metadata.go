package internal

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func GetMetadata[T any](app *cli.App, key string, defaultValue T) T {
	if app.Metadata == nil {
		return defaultValue
	}
	value, ok := app.Metadata[key]
	if !ok {
		return defaultValue
	}
	if castedValue, ok := value.(T); ok {
		return castedValue
	}
	panic(fmt.Sprintf("could not cast the metadata value to type %T", defaultValue))
}

func SetMetadata(app *cli.App, key string, value any) {
	if app.Metadata == nil {
		app.Metadata = make(map[string]any)
	}
	app.Metadata[key] = value
}
