package main

import (
	"fmt"
	"image"
	"math/cmplx"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/deeean/go-vector/vector2"
)

// check_mandelbrot checks how many iterations are needed
// to establish if a point is a part of the mandelbrot set.
func check_mandelbrot(points []complex128) []int {
	current_iteration := 0
	iterations := 128
	zs := make([]complex128, len(points))
	mandelbrot_iters := make([]int, len(points))
	for i := range mandelbrot_iters {
		mandelbrot_iters[i] = -1
	}

	total_points := 0
	for {
		points_changed := 0
		// fmt.Println(iterations)
		for point_i := range points {
			point := points[point_i]
			z := zs[point_i]

			if mandelbrot_iters[point_i] == -1 {
				for i := current_iteration; i < iterations; i++ {
					z = z*z + point
					if cmplx.Abs(z) >= 2 {
						mandelbrot_iters[point_i] = i
						points_changed++
						total_points++
						break
					}
				}
				zs[point_i] = z
			}
		}

		current_iteration = iterations
		iterations *= 2

		// if every point has diverged break
		if total_points == len(points) {
			break
		}

		// at least some point should've diverged
		if float64(total_points)/float64(len(points)) > 0.05 {
			// very few new divergent points since last check
			if float64(points_changed)/float64(len(points)-total_points) < 0.01 {
				break
			}
		}
	}

	return mandelbrot_iters
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

func generateFrame(
	center vector2.Vector2,
	size float64,
	image_size int,
	color_pallete []HSV,
) *image.RGBA {
	points := make([]complex128, 0, image_size*image_size)
	for x := 0; x < image_size; x++ {
		for y := 0; y < image_size; y++ {
			real := center.X - size/2 + size*float64(x)/float64(image_size)
			im := center.Y - size/2 + size*float64(y)/float64(image_size)
			points = append(points, complex(real, im))
		}
	}

	mandelbrot_iters := check_mandelbrot(points)

	image := image.NewRGBA(image.Rect(0, 0, image_size, image_size))
	for x := 0; x < image_size; x++ {
		for y := 0; y < image_size; y++ {
			iterations := mandelbrot_iters[x*image_size+y]

			color := HSV{0, 0, 0}
			if iterations >= 0 {
				color = color_pallete[iterations%color_number]
			}
			image.Set(x, y, color)
		}
	}

	return image
}

const window_size = 700
const color_number = 200

func main() {
	x := -1.162779
	y := 0.2713448
	if len(os.Args) == 3 {
		var err error
		x, err = strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			panic(err)
		}

		y, err = strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Użycie: ./mandelbrot-go <x> <y>")
		fmt.Printf("Domyślny punkt: %f %f", x, y)
	}

	center := vector2.Vector2{X: x, Y: y}

	app := app.New()
	w := app.NewWindow("Mandelbrot set")
	w.Resize(fyne.NewSize(window_size, window_size))

	go func() {
		size := float64(2)
		color_pallete := make_color_palette(color_number)

		for {
			frame := generateFrame(
				center,
				size,
				window_size,
				color_pallete,
			)
			img := canvas.NewImageFromImage(frame)
			w.SetContent(img)

			size *= 0.9
			// fmt.Println("==================")
			// fmt.Println(size)
		}
	}()

	w.ShowAndRun()
}
