package main

import (
	"log"

	"github.com/Elliot727/gocvkit"
)

func main() {
	app, err := gocvkit.NewApp("config.toml")
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}
	defer app.Close()

	// Optional: you can add per-frame custom logic here if needed
	// (e.g. custom overlays beyond what's in config)
	if err := app.Run(nil); err != nil {
		log.Fatal("Run failed:", err)
	}
}
