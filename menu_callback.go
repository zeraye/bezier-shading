package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/zeraye/bezier-shading/pkg/draw"
)

func lightColorPickerCallback(g *Game, lightColorLabel *widget.Label) func(color.Color) {
	return func(c color.Color) {
		g.lightColor = c
		r, g, b, _ := draw.ColorRGBA(g.lightColor)
		lightColorLabel.Text = fmt.Sprintf("color: (%d, %d, %d)", r, g, b)
		lightColorLabel.Refresh()
	}
}

func lightColorButtonTapped(g *Game, lightColorLabel *widget.Label) func() {
	return func() {
		dialog.ShowColorPicker("Color picker", "light color", lightColorPickerCallback(g, lightColorLabel), g.window)
	}
}

func animationButtonTapped(g *Game) func() {
	return func() {
		if g.LightAnimation {
			g.menu.lightAnimationButton.Text = "Start animation"
		} else {
			g.menu.lightAnimationButton.Text = "Stop animation"
		}
		g.menu.lightAnimationButton.Refresh()
		g.LightAnimation = !g.LightAnimation
	}
}

func surfaceButtonTapped(g *Game) func() {
	return func() {
		if g.surface == "bezier" {
			g.menu.surfaceButton.Text = "Hemisphere (currently)"
		} else {
			g.menu.surfaceButton.Text = "Bezier (currently)"
		}
		g.menu.surfaceButton.Refresh()
		if g.surface == "bezier" {
			g.surface = "hemisphere"
		} else {
			g.surface = "bezier"
		}
	}
}

func normalMapfileOpenCallback(g *Game, normalMapLabel *widget.Label) func(fyne.URIReadCloser, error) {
	return func(urc fyne.URIReadCloser, err error) {
		if err != nil {
			panic(err)
		}
		if urc == nil {
			return
		}
		g.normalMap, err = getImageFromFilePath(urc.URI().Path())
		if err != nil {
			panic(err)
		}
		normalMapLabel.Text = "file: " + urc.URI().Name()
		normalMapLabel.Refresh()
	}
}

func normalMapButtonTapped(g *Game, normalMapLabel *widget.Label) func() {
	return func() {
		dialog.ShowFileOpen(normalMapfileOpenCallback(g, normalMapLabel), g.window)
	}
}

func backgroundSolidColorPickerCallback(g *Game, backgroundSolidColorLabel *widget.Label) func(color.Color) {
	return func(c color.Color) {
		g.backgroundSolidColor = c
		r, g, b, _ := draw.ColorRGBA(g.backgroundSolidColor)
		backgroundSolidColorLabel.Text = fmt.Sprintf("color: (%d, %d, %d)", r, g, b)
		backgroundSolidColorLabel.Refresh()
	}
}

func backgroundSolidColorButtonTapped(g *Game, backgroundSolidColorLabel *widget.Label) func() {
	return func() {
		dialog.ShowColorPicker("Color picker", "background solid color", backgroundSolidColorPickerCallback(g, backgroundSolidColorLabel), g.window)
	}
}

func backgroundImagefileOpenCallback(g *Game, backgroundImageLabel *widget.Label) func(fyne.URIReadCloser, error) {
	return func(urc fyne.URIReadCloser, err error) {
		if err != nil {
			panic(err)
		}
		if urc == nil {
			return
		}
		g.backgroundImage, err = getImageFromFilePath(urc.URI().Path())
		if err != nil {
			panic(err)
		}
		backgroundImageLabel.Text = "file: " + urc.URI().Name()
		backgroundImageLabel.Refresh()
	}
}

func backgroundImageButtonTapped(g *Game, backgroundImageLabel *widget.Label) func() {
	return func() {
		dialog.ShowFileOpen(backgroundImagefileOpenCallback(g, backgroundImageLabel), g.window)
	}
}

func backgroundRadioButtonChanged(g *Game, backgroundRadioButton *widget.RadioGroup) func(string) {
	return func(option string) {
		backgroundRadioButton.SetSelected(option)
		backgroundRadioButton.Refresh()
		if option == "Solid color" {
			g.menu.backgroundSolidColorLabel.Show()
			g.menu.backgroundSolidColorButton.Show()
			g.menu.backgroundImageLabel.Hide()
			g.menu.backgroundImageButton.Hide()
			g.isBackgroundSolidColor = true
		} else if option == "Image" {
			g.menu.backgroundSolidColorLabel.Hide()
			g.menu.backgroundSolidColorButton.Hide()
			g.menu.backgroundImageLabel.Show()
			g.menu.backgroundImageButton.Show()
			g.isBackgroundSolidColor = false
		} else {
			panic("Invalid entry for background radio button")
		}
	}
}

func triangulationSliderChanged(g *Game, triangulationSlider *widget.Slider) func(float64) {
	return func(value float64) {
		triangulationSlider.Value = value
		g.triangulation = int(value)
		g.triangles = makeTriangles(g.config, g.points, g.triangulation)
		triangulationSlider.Refresh()
		g.Refresh()
	}
}

func alphaSliderChanged(g *Game, alphaSlider *widget.Slider) func(float64) {
	return func(value float64) {
		alphaSlider.Value = value
		g.alpha = value
		alphaSlider.Refresh()
		g.Refresh()
	}
}

func betaSliderChanged(g *Game, betaSlider *widget.Slider) func(float64) {
	return func(value float64) {
		betaSlider.Value = value
		g.beta = value
		betaSlider.Refresh()
		g.Refresh()
	}
}

func pointsHeightSliderChanged(g *Game, pointsHeightSlider *widget.Slider) func(float64) {
	return func(value float64) {
		if g.pointHeight != nil {
			for points_row_index := range g.points {
				for point_index, point := range g.points[points_row_index] {
					if point == g.pointHeight {
						g.pointsHeight[points_row_index][point_index] = value
					}
				}
			}
			pointsHeightSlider.Value = value
		} else {
			pointsHeightSlider.Value = 0
		}
		pointsHeightSlider.Refresh()
		g.Refresh()
	}
}

func triangulationCheckChanged(g *Game) func(bool) {
	return func(value bool) {
		g.showMesh = value
	}
}
