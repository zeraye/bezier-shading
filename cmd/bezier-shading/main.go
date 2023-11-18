package main

import (
	"log"
	"math"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	bezier_shading "github.com/zeraye/bezier-shading"
	"github.com/zeraye/bezier-shading/pkg/config"
)

func main() {
	config, err := config.LoadStandard("config", "config.toml")
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("FYNE_THEME", config.Window.Theme)

	app := app.NewWithID(config.Window.Name)
	window := app.NewWindow(config.Window.Name)
	game := bezier_shading.NewGame(config, window)

	window.SetContent(game.BuildUI())
	window.Resize(fyne.NewSize(float32(config.Window.Width), float32(config.Window.Height)))
	window.SetFixedSize(config.Window.FixedSize)

	midX := float64(config.UI.RasterWidth) / 2
	midY := float64(config.UI.RasterHeight) / 2
	Rmax := float64(config.UI.RasterWidth) / 2
	Rmin := config.Light.SpiralMinRadius
	R := Rmin
	angle := 0.0
	incremental := true

	go updateLightPoint(game, &midX, &midY, &R, &angle)
	go updateRadiusAngle(config, &R, &Rmin, &Rmax, &angle, &incremental)

	window.ShowAndRun()
}

func updateLightPoint(game *bezier_shading.Game, midX, midY, R, angle *float64) {
	for range time.Tick(time.Duration(time.Millisecond)) {
		if !game.Busy {
			game.Busy = true
			if game.LightAnimation {
				game.LightPoint.X = *midX + *R*math.Sin(*angle)
				game.LightPoint.Y = *midY + *R*math.Cos(*angle)
			}
			game.Refresh()
		}
	}
}

func updateRadiusAngle(config *config.Config, R, Rmin, Rmax, angle *float64, incremental *bool) {
	duration := int64(time.Millisecond) * config.Light.SpiralUpdateMiliseconds
	for range time.Tick(time.Duration(duration)) {
		if *incremental {
			*R += config.Light.SpiralRadiusDelta
		} else {
			*R -= config.Light.SpiralRadiusDelta
		}
		if *R <= *Rmin || *R >= *Rmax {
			*incremental = !*incremental
		}
		*angle += config.Light.SpiralAngleDelta
	}
}
