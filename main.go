package main

import (
	"flag"
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
func CheckMandelbrot(points []complex128, precisionInfo PrecisionInfo) []int {
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

		if precisionInfo.maxiter >= 0 && iterations >= precisionInfo.maxiter {
			break
		}

		// if every point has diverged break
		if totalPoints == len(points) {
			break
		}

		// at least some point should've diverged
		if float64(totalPoints)/float64(len(points)) > precisionInfo.mindiv {
			// very few new divergent points since last check
			if float64(pointsChanged)/float64(len(points)-totalPoints) < precisionInfo.precision {
				break
			}
		}
	}

	return mandelbrotIters
}

// makeColorPalette creates array of HSV encoded colors
// that are ordered by hue.
func makeColorPalette(number int) []HSV {
	colors := make([]HSV, number)
	for i := 0; i < number; i++ {
		color := 30 + uint16(210./float64(number)*float64(i))
		colors[i] = HSV{color, 1, 1}
	}
	return colors
}

// worker processes part of Mandelbrot image
func worker(
	upperLeft vector2.Vector2,
	step float64,
	precisionInfo PrecisionInfo,
	colorPallete []HSV,
	colorNumber int,
	image *image.RGBA,
	region image.Rectangle,
	completionChannel chan struct{},
) {
	width := region.Dx()
	height := region.Dy()
	points := make([]complex128, 0, width*height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			real := upperLeft.X + float64(x)*step
			im := upperLeft.Y + float64(y)*step
			points = append(points, complex(real, im))
		}
	}

	mandelbrotIters := CheckMandelbrot(points, precisionInfo)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			iterations := mandelbrotIters[x*height+y]

			color := HSV{0, 0, 0}
			if iterations >= 0 {
				color = colorPallete[iterations%colorNumber]
			}
			image.Set(region.Min.X+x, region.Min.Y+y, color)
		}
	}

	completionChannel <- struct{}{}
}

// generateFrame creates new Mandelbrot set image
func generateFrame(
	center vector2.Vector2,
	size float64,
	imageSize int,
	precisionInfo PrecisionInfo,
	colorPallete []HSV,
	colorNumber int,
	numberWorkers int,
) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))
	completionChannel := make(chan struct{})

	upperLeft := *center.SubScalar(size / 2)
	rowsPerWorker := imageSize / numberWorkers
	sizePerWorker := size / float64(numberWorkers)
	step := size / float64(imageSize)
	for i := 0; i < numberWorkers; i++ {
		lastWorker := 0
		if i == numberWorkers-1 {
			lastWorker = imageSize - rowsPerWorker*numberWorkers
		}

		region := image.Rect(
			0,
			i*rowsPerWorker,
			imageSize,
			i*rowsPerWorker+rowsPerWorker+lastWorker,
		)

		go worker(
			upperLeft,
			step,
			precisionInfo,
			colorPallete,
			colorNumber,
			img,
			region,
			completionChannel,
		)
		upperLeft.Y += sizePerWorker
	}

	for i := 0; i < numberWorkers; i++ {
		<-completionChannel
	}

	return img
}

type PrecisionInfo struct {
	precision float64
	mindiv    float64
	maxiter   int
}

func main() {
	profFlag := flag.Bool("prof", false, "collect profiling data")
	numberWorkers := flag.Int("workers", 8, "number of threads")
	zoom := flag.Float64("zoom", 0.9, "zoom amount per frame (range from 0 to 1)")
	precision := flag.Float64("precision", 0.01, "precision of Mandelbrot set computation (range from 0 to 1)")
	minDiv := flag.Float64("mindiv", 0.05, "minimal percent of found divergent points to continue to next frame (range from 0 to 1)")
	maxiter := flag.Int("maxiter", -1, "maximum number of iteration per frame")
	numberColors := flag.Int("colors", 200, "number of colors used to graph Mandelbroot set")
	windowSize := flag.Int("size", 700, "window size")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] [x] [y]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *profFlag {
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
	}

	x := -1.16278126
	y := 0.27134518
	if len(flag.Args()) == 1 || len(flag.Args()) > 2 {
		flag.Usage()
		panic("Wrong number of arguments")
	} else if len(flag.Args()) == 2 {
		var err error
		x, err = strconv.ParseFloat(flag.Arg(0), 64)
		if err != nil {
			flag.Usage()
			panic(err)
		}

		y, err = strconv.ParseFloat(flag.Arg(1), 64)
		if err != nil {
			flag.Usage()
			panic(err)
		}
	}

	fmt.Printf("Point: %.10f %.10f\n", x, y)

	center := vector2.Vector2{X: x, Y: y}

	app := app.New()
	w := app.NewWindow("Mandelbrot set")
	w.Resize(fyne.NewSize(float32(*windowSize), float32(*windowSize)))

	img := image.NewRGBA(image.Rect(0, 0, *windowSize, *windowSize))
	canvasImage := canvas.NewImageFromImage(img)
	w.SetContent(canvasImage)

	size := float64(2)
	colorPallete := makeColorPalette(*numberColors)
	go func() {
		for {
			frame := generateFrame(
				center,
				size,
				*windowSize,
				PrecisionInfo{
					precision: *precision,
					mindiv:    *minDiv,
					maxiter:   *maxiter,
				},
				colorPallete,
				*numberColors,
				*numberWorkers,
			)
			canvasImage.Image = frame
			canvasImage.Refresh()
			size *= *zoom
		}
	}()

	w.ShowAndRun()
}
