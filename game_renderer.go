package bezier_shading

import (
	"image"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/zeraye/bezier-shading/pkg/draw"
	"github.com/zeraye/bezier-shading/pkg/geom"
)

type gameRenderer struct {
	raster  *canvas.Raster
	objects []fyne.CanvasObject
	game    *Game
}

func (gr *gameRenderer) Destroy() {
}

func (gr *gameRenderer) Layout(size fyne.Size) {
	gr.raster.Resize(size)
}

func (gr *gameRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(gr.game.config.UI.RasterWidth), float32(gr.game.config.UI.RasterHeight))
}

func (gr *gameRenderer) Objects() []fyne.CanvasObject {
	return gr.objects
}

func (gr *gameRenderer) Refresh() {
	canvas.Refresh(gr.raster)
}

// Draw game raster (canvas, not menu)
func (gr *gameRenderer) Draw(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(gr.game.config.UI.RasterWidth), int(gr.game.config.UI.RasterHeight)))

	// draw raster background
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, draw.RGBAToColor(gr.game.config.UI.BackgroundColorRGBA))
		}
	}

	blueColor := draw.RGBAToColor([4]uint8{0, 0, 255, 255})
	whiteColor := draw.RGBAToColor([4]uint8{255, 255, 255, 255})
	yellowColor := draw.RGBAToColor([4]uint8{255, 255, 0, 255})

	var wg sync.WaitGroup
	wg.Add(len(gr.game.triangles))
	for _, tri := range gr.game.triangles {
		go FillPolygon([]*geom.Point{tri.P0, tri.P1, tri.P2}, gr.game.backgroundSolidColor, img, gr.game, &wg)
	}
	wg.Wait()

	if gr.game.showMesh {
		wg.Add(len(gr.game.triangles))
		for _, tri := range gr.game.triangles {
			go OutlineTriangle(tri, blueColor, img, &wg)
		}
		wg.Wait()
	}

	for points_row_index := range gr.game.points {
		for _, point := range gr.game.points[points_row_index] {
			if point == gr.game.pointHeight {
				draw.DrawCircle(*point, 8, blueColor, true, img)
			} else {
				draw.DrawCircle(*point, 8, whiteColor, true, img)
			}
		}
	}

	draw.DrawCircle(*gr.game.LightPoint, 8, yellowColor, true, img)

	// draw raster border
	for x := 0; x < img.Bounds().Dx(); x++ {
		img.Set(x, 0, draw.RGBAToColor(gr.game.config.UI.RasterBorderColorRGBA))
		img.Set(x, img.Bounds().Dy()-1, draw.RGBAToColor(gr.game.config.UI.RasterBorderColorRGBA))
	}
	for y := 0; y < img.Bounds().Dx(); y++ {
		img.Set(0, y, draw.RGBAToColor(gr.game.config.UI.RasterBorderColorRGBA))
		img.Set(img.Bounds().Dx()-1, y, draw.RGBAToColor(gr.game.config.UI.RasterBorderColorRGBA))
	}

	gr.game.Busy = false

	return img
}
