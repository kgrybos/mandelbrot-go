package color

import "math"

// HSV struct represents a color as HSV number triple.
type HSV struct {
	H uint16
	S float64
	V float64
}

func (color HSV) RGBA() (r, g, b, a uint32) {
	c := color.V * color.S
	x := c * (1 - math.Abs((math.Mod(float64(color.H)/60, 2) - 1)))
	m := color.V - c

	var rf, gf, bf float64
	switch {
	case color.H < 60:
		rf = c
		gf = x
		bf = 0
	case color.H < 120:
		rf = x
		gf = c
		bf = 0
	case color.H < 180:
		rf = 0
		gf = c
		bf = x
	case color.H < 240:
		rf = 0
		gf = x
		bf = c
	case color.H < 300:
		rf = x
		gf = 0
		bf = c
	case color.H < 360:
		rf = c
		gf = 0
		bf = x
	}

	r = uint32((rf + m) * 0xffff)
	g = uint32((gf + m) * 0xffff)
	b = uint32((bf + m) * 0xffff)
	a = 0xffff
	return r, g, b, a
}

type ColorCycle []HSV

func (cycle ColorCycle) Get(n int) HSV {
	return cycle[n%len(cycle)]
}

// makeColorPalette creates array of HSV encoded colors
// that are ordered by hue.
func MakeColorPalette(number int) ColorCycle {
	colors := make([]HSV, number)
	for i := 0; i < number; i++ {
		color := 30 + uint16(210./float64(number)*float64(i))
		colors[i] = HSV{color, 1, 1}
	}
	return colors
}
