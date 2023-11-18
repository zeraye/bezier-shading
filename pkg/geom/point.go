package geom

import "math"

type Point struct {
	X, Y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func Dist(p0, p1 *Point) float64 {
	return math.Sqrt(math.Pow(p0.X-p1.X, 2) + math.Pow(p0.Y-p1.Y, 2))
}
