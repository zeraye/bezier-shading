package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/zeraye/bezier-shading/pkg/config"
	"github.com/zeraye/bezier-shading/pkg/draw"
	"github.com/zeraye/bezier-shading/pkg/geom"
)

type Game struct {
	Busy           bool
	LightPoint     *geom.Point
	LightAnimation bool

	widget.BaseWidget

	config *config.Config
	window fyne.Window
	menu   *Menu

	lightColor             color.Color
	lightHeight            float64
	backgroundSolidColor   color.Color
	backgroundImage        image.Image
	normalMap              image.Image
	isBackgroundSolidColor bool
	points                 [][]*geom.Point
	pointsHeight           [][]float64
	triangulation          int
	triangles              []*geom.Triangle
	pointHeight            *geom.Point
	showMesh               bool
	surface                string
	alpha                  float64
	beta                   float64
}

func NewGame(config *config.Config, window fyne.Window) *Game {
	menu := NewMenu(config)
	lightColor := draw.RGBAToColor(config.Defaults.LightColorRGBA)
	lightAnimation := config.Defaults.LightAnimation
	lightHeight := config.Defaults.LightHeight
	lightPoint := geom.NewPoint(float64(config.UI.RasterWidth)/2, float64(config.UI.RasterHeight)/2)
	triangulation := config.Defaults.Triangulation
	backgroundSolidColor := draw.RGBAToColor(config.Defaults.DefaultBackgroundSolidColorRGBA)

	size := config.Defaults.InterpolationPointsPerSide
	points := make([][]*geom.Point, size)
	for i := 0; i < size; i++ {
		points[i] = make([]*geom.Point, size)
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			points[i][j] = geom.NewPoint(float64(config.UI.RasterWidth*i)/float64(size-1), float64(config.UI.RasterHeight*j)/float64(size-1))
		}
	}
	triangles := makeTriangles(config, points, triangulation)

	pointsHeight := make([][]float64, size)
	for i := 0; i < size; i++ {
		pointsHeight[i] = make([]float64, size)
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			pointsHeight[i][j] = 0
		}
	}

	var backgroundImage image.Image = nil
	var normalMap image.Image = nil

	game := &Game{
		config:                 config,
		menu:                   menu,
		window:                 window,
		lightColor:             lightColor,
		LightAnimation:         lightAnimation,
		lightHeight:            lightHeight,
		LightPoint:             lightPoint,
		backgroundSolidColor:   backgroundSolidColor,
		backgroundImage:        backgroundImage,
		normalMap:              normalMap,
		isBackgroundSolidColor: true,
		points:                 points,
		pointsHeight:           pointsHeight,
		triangulation:          triangulation,
		triangles:              triangles,
		Busy:                   true,
		showMesh:               false,
		surface:                "bezier",
		alpha:                  0,
		beta:                   0,
	}

	game.ExtendBaseWidget(game)

	return game
}

func (g *Game) BuildUI() fyne.CanvasObject {
	return container.NewBorder(nil, nil, g.menu.BuildUI(g), g)
}

func (g *Game) CreateRenderer() fyne.WidgetRenderer {
	renderer := &gameRenderer{game: g}
	raster := canvas.NewRaster(renderer.Draw)
	renderer.raster = raster
	renderer.objects = []fyne.CanvasObject{raster}

	return renderer
}

func (g *Game) Tapped(ev *fyne.PointEvent) {
	mouse_pos := geom.NewPoint(float64(ev.Position.X), float64(ev.Position.Y))

	for points_row_index := range g.points {
		for point_index, point := range g.points[points_row_index] {
			if geom.Dist(point, mouse_pos) <= 8 {
				g.pointHeight = nil
				g.menu.pointsHeightSlider.SetValue(g.pointsHeight[points_row_index][point_index])
				g.pointHeight = point
			}
		}
	}
}

func (g *Game) TappedSecondary(ev *fyne.PointEvent) {
}

func (g *Game) Dragged(ev *fyne.DragEvent) {
	mouse_pos := geom.NewPoint(float64(ev.Position.X), float64(ev.Position.Y))
	g.LightPoint = mouse_pos
}

func (g *Game) DragEnd() {
}

func makeTriangles(config *config.Config, points [][]*geom.Point, triangulation int) []*geom.Triangle {
	size := config.Defaults.InterpolationPointsPerSide
	triangles := []*geom.Triangle{}
	for i := 0; i < size-1; i++ {
		for j := 0; j < size-1; j++ {
			for m := 0; m < triangulation; m++ {
				for n := 0; n < triangulation; n++ {
					sideLength := float64(config.UI.RasterWidth) / float64((size-1)*triangulation)
					triangles = append(triangles,
						geom.NewTriangle(
							geom.NewPoint(points[i][j].X+float64(m)*sideLength, points[i][j].Y+float64(n)*sideLength),
							geom.NewPoint(points[i][j].X+float64(m+1)*sideLength, points[i][j].Y+float64(n)*sideLength),
							geom.NewPoint(points[i][j].X+float64(m)*sideLength, points[i][j].Y+float64(n+1)*sideLength),
						),
						geom.NewTriangle(
							geom.NewPoint(points[i][j].X+float64(m+1)*sideLength, points[i][j].Y+float64(n)*sideLength),
							geom.NewPoint(points[i][j].X+float64(m+1)*sideLength, points[i][j].Y+float64(n+1)*sideLength),
							geom.NewPoint(points[i][j].X+float64(m)*sideLength, points[i][j].Y+float64(n+1)*sideLength),
						),
					)
				}
			}
		}
	}
	return triangles
}
