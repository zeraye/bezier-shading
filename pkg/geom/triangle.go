package geom

type Triangle struct {
	P0, P1, P2 *Point
}

func NewTriangle(p0, p1, p2 *Point) *Triangle {
	return &Triangle{p0, p1, p2}
}
