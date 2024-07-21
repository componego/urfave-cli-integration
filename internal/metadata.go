/*
Copyright 2024 Volodymyr Konstanchuk and the Componego Framework contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internal

import (
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
	return defaultValue
}

func SetMetadata(app *cli.App, key string, value any) {
	if app.Metadata == nil {
		app.Metadata = make(map[string]any)
	}
	app.Metadata[key] = value
}
