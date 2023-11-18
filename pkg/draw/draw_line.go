package draw

import (
	"image"
	"image/color"
	"math"

	"github.com/zeraye/bezier-shading/pkg/geom"
)

func DrawLine(p0, p1 geom.Point, color color.Color, img *image.RGBA) {
	BresenhamDrawLine(p0, p1, color, img)
}

// Bresenham's line algorithm: https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func BresenhamDrawLine(p0, p1 geom.Point, color color.Color, img *image.RGBA) {
	if math.Abs(p1.Y-p0.Y) < math.Abs(p1.X-p0.X) {
		if p0.X > p1.X {
			DrawLineLow(p1, p0, color, img)
		} else {
			DrawLineLow(p0, p1, color, img)
		}
	} else {
		if p0.Y > p1.Y {
			DrawLineHigh(p1, p0, color, img)
		} else {
			DrawLineHigh(p0, p1, color, img)
		}
	}
}

func DrawLineLow(p0, p1 geom.Point, color color.Color, img *image.RGBA) {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	yi := 1.0
	if dy < 0 {
		yi = -1
		dy = -dy
	}
	D := 2*dy - dx
	y := p0.Y

	for x := p0.X; x < p1.X; x++ {
		img.Set(int(x), int(y), color)
		if D > 0 {
			y += yi
			D += 2 * (dy - dx)
		} else {
			D += 2 * dy
		}
	}
}

func DrawLineHigh(p0, p1 geom.Point, color color.Color, img *image.RGBA) {
	dx := p1.X - p0.X
	dy := p1.Y - p0.Y
	xi := 1.0
	if dx < 0 {
		xi = -1
		dx = -dx
	}
	D := 2*dx - dy
	x := p0.X

	for y := p0.Y; y < p1.Y; y++ {
		img.Set(int(x), int(y), color)
		if D > 0 {
			x += xi
			D += 2 * (dx - dy)
		} else {
			D += 2 * dx
		}
	}
}
