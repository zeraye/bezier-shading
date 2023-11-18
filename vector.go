package bezier_shading

import (
	"math"

	"github.com/zeraye/bezier-shading/pkg/geom"
	"gonum.org/v1/gonum/stat/combin"
)

type Vec struct {
	x, y, z float64
}

const eps = 10e-12

// return vec from p0 to p1 (vec=p1-p0)
func vecFromPoints(p0, p1 geom.Point) Vec {
	return Vec{p1.X - p0.X, p1.Y - p0.Y, 0}
}

func dotProduct(vec0, vec1 Vec) float64 {
	product := 0.0
	product += vec0.x * vec1.x
	product += vec0.y * vec1.y
	product += vec0.z * vec1.z
	return product
}

func crossProduct(vec0, vec1 Vec) Vec {
	return Vec{
		vec0.y*vec1.z - vec0.z*vec1.y,
		vec0.z*vec1.x - vec0.x*vec1.z,
		vec0.x*vec1.y - vec0.y*vec1.x,
	}
}

func magnitude(vec Vec) float64 {
	return math.Sqrt(vec.x*vec.x + vec.y*vec.y + vec.z*vec.z)
}

func normalize(vec Vec) Vec {
	mag := magnitude(vec)
	if mag == 0 {
		return Vec{0, 0, 0}
	}
	return Vec{vec.x / mag, vec.y / mag, vec.z / mag}
}

func mult(scalar float64, vec Vec) Vec {
	return Vec{vec.x * scalar, vec.y * scalar, vec.z * scalar}
}

func add(vec0, vec1 Vec) Vec {
	return Vec{
		vec0.x + vec1.x,
		vec0.y + vec1.y,
		vec0.z + vec1.z,
	}
}

func add3(vec0, vec1, vec2 Vec) Vec {
	return Vec{
		vec0.x + vec1.x + vec2.x,
		vec0.y + vec1.y + vec2.y,
		vec0.z + vec1.z + vec2.z,
	}
}

func minus(vec0, vec1 Vec) Vec {
	return Vec{
		vec0.x - vec1.x,
		vec0.y - vec1.y,
		vec0.z - vec1.z,
	}
}

func b(i, n int, t float64) float64 {
	return float64(combin.Binomial(n, i)) * math.Pow(t, float64(i)) * math.Pow(1-t, float64(n-i))
}

func bezier(u, v float64, pointsHeight [][]float64) Vec {
	z := 0.0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			z += (pointsHeight[i][j] / 100) * b(i, 3, u) * b(j, 3, v)
		}
	}
	vec := Vec{u, v, z}
	return vec
}

func bezierDU(u, v float64, pointsHeight [][]float64) Vec {
	z := 0.0
	for i := 0; i <= 2; i++ {
		for j := 0; j <= 3; j++ {
			z += (pointsHeight[i+1][j]/100 - pointsHeight[i][j]/100) * b(i, 2, u) * b(j, 3, v)
		}
	}
	z *= 3
	vec := Vec{1, 0, z}
	return vec
}

func bezierDV(u, v float64, pointsHeight [][]float64) Vec {
	z := 0.0
	for i := 0; i <= 3; i++ {
		for j := 0; j <= 2; j++ {
			z += (pointsHeight[i][j+1]/100 - pointsHeight[i][j]/100) * b(i, 3, u) * b(j, 2, v)
		}
	}
	z *= 3
	vec := Vec{0, 1, z}
	return vec
}
