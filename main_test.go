package main

import (
	"testing"

	"github.com/deeean/go-vector/vector2"
	"github.com/kgrybos/mandelbrot-go/color"
)

func BenchmarkCheckMandelbrot(b *testing.B) {
	points := make([]complex128, 490000)
	center := vector2.Vector2{X: -1.162779, Y: 0.2713448}
	// center := vector2.Vector2{X: 100, Y: 100}
	size := 0.001
	for x := 0; x < 700; x++ {
		for y := 0; y < 700; y++ {
			real := center.X - size/2 + size*float64(x)/float64(700)
			im := center.Y - size/2 + size*float64(y)/float64(700)
			points = append(points, complex(real, im))
		}
	}

	for i := 0; i < b.N; i++ {
		CheckMandelbrot(points, PrecisionInfo{precision: 0.1, mindiv: 0.5, maxiter: -1})
	}
}

func BenchmarkGenerateFrame(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateFrame(
			vector2.Vector2{X: -1.162779, Y: 0.2713448},
			0.001,
			700,
			PrecisionInfo{precision: 0.1, mindiv: 0.5, maxiter: -1},
			color.MakeColorPalette(200),
			8,
		)
	}
}
