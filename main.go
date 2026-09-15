package main

import (
	"embed"
	"io/fs"
	"os"

	"wombat/internal/app"
)

//go:embed all:frontend/public
var embeddedAssets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	assets, err := fs.Sub(embeddedAssets, "frontend/public")
	if err != nil {
		panic(err)
	}
	os.Exit(app.Run(assets, icon))
}
