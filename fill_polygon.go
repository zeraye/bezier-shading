package bezier_shading

import (
	"cmp"
	"image"
	"image/color"
	"math"
	"slices"
	"sync"

	"github.com/zeraye/bezier-shading/pkg/draw"
	"github.com/zeraye/bezier-shading/pkg/geom"
)

func mod(value, mod int) int {
	if value < 0 {
		return value + mod
	}
	return value
}

func FillPolygon(points []*geom.Point, color color.Color, img *image.RGBA, g *Game, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(points) < 3 {
		return
	}

	ind := make([]int, len(points))
	for i := 0; i < len(points); i++ {
		ind[i] = i
	}
	slices.SortFunc(ind, func(i, j int) int {
		return cmp.Compare(points[i].Y, points[j].Y)
	})

	n_arr := []Vec{}
	z_arr := []float64{}
	for i := 0; i < len(points); i++ {
		ndu := bezierDU(points[i].X/float64(g.config.UI.RasterWidth), points[i].Y/float64(g.config.UI.RasterHeight), g.pointsHeight)
		ndv := bezierDV(points[i].X/float64(g.config.UI.RasterWidth), points[i].Y/float64(g.config.UI.RasterHeight), g.pointsHeight)
		n_arr = append(n_arr, normalize(crossProduct(ndu, ndv)))
		z_arr = append(z_arr, bezier(points[i].X/float64(g.config.UI.RasterWidth), points[i].Y/float64(g.config.UI.RasterHeight), g.pointsHeight).z)
	}

	ymin := points[ind[0]].Y
	ymax := points[ind[len(points)-1]].Y
	aet := []*geom.Segment{}

	for y := ymin; y < ymax; y++ {
		for k := range ind[:len(ind)-1] {
			if points[ind[k]].Y == y {
				curr := points[ind[k]]
				prev := points[mod(ind[k]-1, len(points))]
				next := points[mod(ind[k]+1, len(points))]
				if prev.Y >= curr.Y {
					aet = append(aet, geom.NewSegment(prev, curr))
				} else {
					// remove (prev, curr from aet)
					for i, seg := range aet {
						if *seg == *geom.NewSegment(prev, curr) {
							aet[i] = aet[len(aet)-1]
							aet = aet[:len(aet)-1]
							break
						}
					}
				}
				if next.Y >= curr.Y {
					aet = append(aet, geom.NewSegment(next, curr))
				} else {
					// remove (next, curr from aet)
					for i, seg := range aet {
						if *seg == *geom.NewSegment(next, curr) {
							aet[i] = aet[len(aet)-1]
							aet = aet[:len(aet)-1]
							break
						}
					}
				}
			}
		}
		slices.SortFunc(aet, func(s0, s1 *geom.Segment) int {
			return cmp.Compare(s0.P0.X, s1.P0.X)
		})

		for i := 0; i < len(aet)/2; i++ {
			x0 := getX(y, *aet[i])
			x1 := getX(y, *aet[i+1])
			for x := x0; x < x1; x++ {
				var normalmapVec *Vec = nil

				if !g.isBackgroundSolidColor && g.backgroundImage != nil {
					color = g.backgroundImage.At(int(x), int(y))
				}
				if g.normalMap != nil {
					normalmapVec = getNormalVecFromColor(g.normalMap.At(int(x), int(y)))
				}

				img.Set(int(x), int(y), calcColor(color, g, x, y, n_arr, z_arr, points, normalmapVec))
			}
		}
	}
}

func calcColor(c color.Color, g *Game, x, y float64, n_arr []Vec, z_arr []float64, points []*geom.Point, normalmapVec *Vec) color.Color {
	kd := g.menu.kdSlider.Value
	ks := g.menu.ksSlider.Value
	ILr, ILg, ILb, _ := draw.ColorNormalRGBA(g.lightColor)
	IOr, IOg, IOb, _ := draw.ColorNormalRGBA(c)
	m := g.menu.mSlider.Value

	p := geom.NewPoint(x, y)
	v0 := vecFromPoints(*points[0], *points[1])
	v1 := vecFromPoints(*points[0], *points[2])
	v2 := vecFromPoints(*points[0], *p)
	d00 := dotProduct(v0, v0)
	d01 := dotProduct(v0, v1)
	d11 := dotProduct(v1, v1)
	d20 := dotProduct(v2, v0)
	d21 := dotProduct(v2, v1)
	denom := d00*d11 - d01*d01
	wv := (d11*d20 - d01*d21) / denom
	ww := (d00*d21 - d01*d20) / denom
	wu := 1 - wv - ww
	weight := Vec{wu, wv, ww}

	n := add3(mult(weight.x, n_arr[0]), mult(weight.y, n_arr[1]), mult(weight.z, n_arr[2]))
	n = normalize(n)

	z := z_arr[0]*weight.x + z_arr[1]*weight.y + z_arr[2]*weight.z

	if normalmapVec != nil {
		binorm := crossProduct(n, Vec{0, 0, 1})
		binorm = normalize(binorm)
		if (n == Vec{0, 0, 1}) {
			binorm = Vec{0, 1, 0}
		}

		normalmapVec.y = -normalmapVec.y

		tangent := crossProduct(binorm, n)
		tangent = normalize(tangent)

		n.x = normalmapVec.x*tangent.x + normalmapVec.y*tangent.y + normalmapVec.z*tangent.z
		n.y = normalmapVec.x*binorm.x + normalmapVec.y*binorm.y + normalmapVec.z*binorm.z
		n.z = normalmapVec.x*n.x + normalmapVec.y*n.y + normalmapVec.z*n.z

		n = normalize(n)

		z += normalmapVec.z
	}

	z *= 100

	l := Vec{(g.LightPoint.X - x), (g.LightPoint.Y - y), g.lightHeight - z}
	l = normalize(l)

	v := Vec{0, 0, 1}

	r := minus(mult(2*dotProduct(n, l), n), l)
	r = normalize(r)

	cosNL := dotProduct(n, l)
	if cosNL < 0 {
		cosNL = 0
	}
	cosVR := dotProduct(v, r)
	if cosVR < 0 {
		cosVR = 0
	}
	cosmVR := math.Pow(cosVR, m)

	kdcosNL := kd * cosNL * 255
	kscosmVR := ks * cosmVR * 255

	Ir := math.Min(ILr*IOr*kdcosNL+ILr*IOr*kscosmVR, 255)
	Ig := math.Min(ILg*IOg*kdcosNL+ILg*IOg*kscosmVR, 255)
	Ib := math.Min(ILb*IOb*kdcosNL+ILb*IOb*kscosmVR, 255)

	return color.RGBA{uint8(Ir), uint8(Ig), uint8(Ib), 255}
}

func getX(y float64, s geom.Segment) float64 {
	return s.P0.X + (y-s.P0.Y)*growthRate(s)
}

func growthRate(s geom.Segment) float64 {
	if s.P1.Y-s.P0.Y == 0 {
		return 0
	}
	return (s.P1.X - s.P0.X) / (s.P1.Y - s.P0.Y)
}

func getNormalVecFromColor(c color.Color) *Vec {
	r, g, b, _ := draw.ColorNormalRGBA(c)
	return &Vec{(r - 0.5) * 2, (g - 0.5) * 2, b}
}

func OutlineTriangle(tri *geom.Triangle, color color.Color, img *image.RGBA, wg *sync.WaitGroup) {
	defer wg.Done()

	draw.DrawLine(*tri.P0, *tri.P1, color, img)
	draw.DrawLine(*tri.P1, *tri.P2, color, img)
	draw.DrawLine(*tri.P2, *tri.P0, color, img)
}
