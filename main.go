package main

import (
	"embed"
	"flexole/app"
)

//go:embed server.yml client.yml
var samples embed.FS

func main() {
	app.Run(&samples)
}
