package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Window   WindowConfig
	UI       UIConfig
	Defaults DefaultsConfig
	Light    LightConfig
}

type WindowConfig struct {
	Name      string
	Width     int
	Height    int
	FixedSize bool
	Theme     string
}

type UIConfig struct {
	BackgroundColorRGBA          [4]uint8
	SecondaryBackgroundColorRGBA [4]uint8
	RasterBorderColorRGBA        [4]uint8
	RasterWidth                  int
	RasterHeight                 int
}

type DefaultsConfig struct {
	Kd                              float64 // coefficient describing the impact of a given component on the result (0-1)
	Ks                              float64 // coefficient describing the impact of a given component on the result (0-1)
	M                               float64 // coefficient describing how much a given triangle is changed (1-100)
	LightColorRGBA                  [4]uint8
	LightAnimation                  bool
	LightHeight                     float64
	DefaultBackgroundSolidColorRGBA [4]uint8
	Triangulation                   int // number of triangles at the side of square
	InterpolationPointsPerSide      int
}

type LightConfig struct {
	SpiralMinRadius         float64
	SpiralRadiusDelta       float64
	SpiralAngleDelta        float64
	SpiralUpdateMiliseconds int64
}

func Load(r io.Reader) (*Config, error) {
	var data Config
	_, err := toml.NewDecoder(r).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func LoadStandard(dir string, filename string) (*Config, error) {
	path := filepath.Join(dir, filename)
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer r.Close()
	return Load(r)
}
