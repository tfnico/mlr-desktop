package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Initialize logger first thing
	if err := InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Add panic recovery for catastrophic failures
	defer func() {
		if r := recover(); r != nil {
			Log.WithFields(logrus.Fields{
				"panic_value": r,
			}).Fatal("Application crashed with unrecovered panic")
		}
	}()

	LogInfo("Application starting", logrus.Fields{
		"app_name": "mlr-desktop",
	})

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:         "mlr-desktop",
		DisableResize: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		LogError(err, "Wails runtime error", logrus.Fields{
			"app_name": "mlr-desktop",
		})
		println("Error:", err.Error())
		os.Exit(1)
	}

	LogInfo("Application shutting down normally", logrus.Fields{
		"app_name": "mlr-desktop",
	})
}

