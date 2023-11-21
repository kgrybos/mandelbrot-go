package main

import (
	"fmt"
	"image"
	"os"
	"runtime/pprof"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/deeean/go-vector/vector2"
)

func doIter(iterations int, z, point complex128) (int, complex128) {
	for i := 0; i < iterations; i++ {
		z = z*z + point
		if real(z)*real(z)+imag(z)*imag(z) >= 4 {
			return i, z
		}
	}
	return -1, z
}

// CheckMandelbrot checks how many iterations are needed
// to establish if a point is a part of the mandelbrot set.
func CheckMandelbrot(points []complex128) []int {
	currentIteration := 0
	iterations := 128
	zs := make([]complex128, len(points))
	mandelbrotIters := make([]int, len(points))
	for i := range mandelbrotIters {
		mandelbrotIters[i] = -1
	}

	totalPoints := 0
	for {
		pointsChanged := 0
		// fmt.Println(iterations)
		for pointI, point := range points {
			if mandelbrotIters[pointI] == -1 {
				var i int
				i, zs[pointI] = doIter(iterations-currentIteration, zs[pointI], point)
				if i != -1 {
					mandelbrotIters[pointI] = currentIteration + i
					pointsChanged++
					totalPoints++
				}
			}
		}

		currentIteration = iterations
		iterations *= 2

		// if every point has diverged break
		if totalPoints == len(points) {
			break
		}

		// at least some point should've diverged
		if float64(totalPoints)/float64(len(points)) > 0.05 {
			// very few new divergent points since last check
			if float64(pointsChanged)/float64(len(points)-totalPoints) < 0.01 {
				break
			}
		}
	}

	return mandelbrotIters
}

// makeColorPalette create array of HSV encoded colors
// that are ordered by hue.
func makeColorPalette(number int) []HSV {
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
	imageSize int,
	colorPallete []HSV,
) *image.RGBA {
	points := make([]complex128, 0, imageSize*imageSize)
	for x := 0; x < imageSize; x++ {
		for y := 0; y < imageSize; y++ {
			real := center.X - size/2 + size*float64(x)/float64(imageSize)
			im := center.Y - size/2 + size*float64(y)/float64(imageSize)
			points = append(points, complex(real, im))
		}
	}

	mandelbrotIters := CheckMandelbrot(points)

	image := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))
	for x := 0; x < imageSize; x++ {
		for y := 0; y < imageSize; y++ {
			iterations := mandelbrotIters[x*imageSize+y]

			color := HSV{0, 0, 0}
			if iterations >= 0 {
				color = colorPallete[iterations%colorNumber]
			}
			image.Set(x, y, color)
		}
	}

	return image
}

const windowSize = 700
const colorNumber = 200

func main() {
	cpufile, err := os.Create("cpu.pprof")
	if err != nil {
		panic(err)
	}
	err = pprof.StartCPUProfile(cpufile)
	if err != nil {
		panic(err)
	}
	defer cpufile.Close()
	defer pprof.StopCPUProfile()

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
	w.Resize(fyne.NewSize(windowSize, windowSize))

	go func() {
		size := float64(2)
		colorPallete := makeColorPalette(colorNumber)

		for {
			frame := generateFrame(
				center,
				size,
				windowSize,
				colorPallete,
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
