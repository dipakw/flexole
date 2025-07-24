package main

import (
	"embed"
	"flexole/app"
)

//go:embed server.yml client.yml
var samples embed.FS

// Application version.
var version = "dev"

func main() {
	app.Run(&app.Config{
		Version: version,
		Samples: &samples,
	})
}
