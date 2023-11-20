package main

import (
	"image"
	"math/cmplx"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

// check_mandelbrot checks how many iterations are needed
// to establish if a point is a part of the mandelbrot set.
func check_mandelbrot(c complex128, iterations int) int {
	z := 0 + 0i
	for i := 0; i < iterations; i++ {
		z = z*z + c
		if cmplx.Abs(z) >= 2 {
			return i
		}
	}
	return iterations
}

// make_color_palette create array of HSV encoded colors
// that are ordered by hue.
func make_color_palette(number int) []HSV {
	colors := make([]HSV, number)
	for i := 0; i < number; i++ {
		color := 30 + uint16(210./float64(number)*float64(i))
		colors[i] = HSV{color, 1, 1}
	}
	return colors
}

// map_range maps value x from range [inlow, inhigh] to range [outlow, outhigh]
func map_range(x, inlow, inhigh, outlow, outhigh float64) float64 {
	return (x-inlow)/(inhigh-inlow)*(outhigh-outlow) + outlow
}

func generateImage(width, height, iterations int) *image.RGBA {
	color_pallete := make_color_palette(iterations + 1)

	image := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			real := map_range(float64(x), 0, float64(width), -2, 2)
			im := map_range(float64(y), 0, float64(height), -2, 2)
			c := complex(real, im)
			chosen_color := color_pallete[check_mandelbrot(c, iterations)]
			image.Set(x, y, chosen_color)
		}
	}

	return image
}

func main() {
	if len(os.Args) != 4 {
		panic("Oczekiwano 3 argumentów: <szerokość> <wysokość> <liczba iteracji>")
	}

	width, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	height, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	iterations, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}

	a := app.New()
	w := a.NewWindow("Mandelbrot set")

	img := canvas.NewImageFromImage(generateImage(width, height, iterations))
	w.SetContent(img)
	w.Resize(fyne.NewSize(float32(width), float32(height)))

	w.ShowAndRun()
}
