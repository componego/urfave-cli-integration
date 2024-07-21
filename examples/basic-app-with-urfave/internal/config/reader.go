package config

func Read(filename string) (map[string]any, error) {
	_ = filename
	// This function can read data from a file, but we have hardcoded the values for this example.
	return map[string]any{
		"hello": map[string]any{
			"message": "Hello World!",
		},
	}, nil
}
