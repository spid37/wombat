package app

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/wailsapp/wails/v2"
	wailsoptions "github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

var (
	appname = "Wombat"
	semver  = "0.0.0-dev"
)

// Run is the main function to run the application
func Run(assets fs.FS, icon []byte) int {
	appData, err := appDataLocation(appname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open app data directory: %v\n", err)
		return 1
	}
	defer crashlog(appData)

	api := NewAPI(appData)

	err = wails.Run(&wailsoptions.App{
		Title:            appname,
		Width:            1200,
		Height:           820,
		MinWidth:         800,
		MinHeight:        600,
		BackgroundColour: &wailsoptions.RGBA{R: 46, G: 52, B: 64, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  api.startup,
		OnDomReady: api.domReady,
		OnShutdown: api.shutdown,
		Bind: []interface{}{
			api,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarDefault(),
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   appname,
				Message: "Cross platform gRPC client",
				Icon:    icon,
			},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "app: error running app: %v\n", err)
		return 1
	}
	return 0
}

func crashlog(appData string) {
	if isDevProcess() {
		return
	}
	if r := recover(); r != nil {
		if _, err := os.Stat(appData); os.IsNotExist(err) {
			_ = os.MkdirAll(appData, 0700)
		}
		var b bytes.Buffer
		b.WriteString(fmt.Sprintf("%+v\n\n", r))
		buf := make([]byte, 1<<20)
		s := runtime.Stack(buf, true)
		b.Write(buf[0:s])
		_ = os.WriteFile(filepath.Join(appData, "crash.log"), b.Bytes(), 0644)
	}
}

// isDevProcess reports whether the process was launched by `wails dev`.
func isDevProcess() bool {
	_, ok := os.LookupEnv("devserver")
	return ok
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
