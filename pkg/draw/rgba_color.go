package draw

import "image/color"

func RGBAToColor(RGBA [4]uint8) color.Color {
	return color.RGBA{RGBA[0], RGBA[1], RGBA[2], RGBA[3]}
}

// Get (r, g, b, a) values from Color, values are 0-255
func ColorRGBA(color color.Color) (r, g, b, a uint32) {
	r, g, b, a = color.RGBA()
	// (r, g, b, a) vales in color.RGBA() are left-shifted by 8 bits
	// thus we right-shift them to get real (0-255) values
	r >>= 8
	g >>= 8
	b >>= 8
	a >>= 8
	return
}

// Get (r, g, b,a) values from Color, values are 0-1
func ColorNormalRGBA(color color.Color) (r, g, b, a float64) {
	ir, ig, ib, ia := ColorRGBA(color)
	r, g, b, a = float64(ir), float64(ig), float64(ib), float64(ia)
	// get values from 0-1
	r /= 255
	g /= 255
	b /= 255
	a /= 255
	return
}
