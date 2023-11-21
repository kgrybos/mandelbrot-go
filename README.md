# Mandelbort set zoom

Program napisany w `go` generujący animację przybliżania w zbiorze Mandelbrota.

<img src="mandelbrot-zoom.gif" alt="Screenshot symulacji" width=320>

## Kompilacja

```sh
go build
```

## Użycie

```
usage: ./mandelbrot-go [flags] [x] [y]
  -colors int
    	number of colors used to graph Mandelbroot set (default 200)
  -maxiter int
    	maximum number of iteration per frame (default -1)
  -mindiv float
    	minimal percent of found divergent points to continue to next frame (range from 0 to 1) (default 0.05)
  -precision float
    	precision of Mandelbrot set computation (range from 0 to 1) (default 0.01)
  -prof
    	collect profiling data
  -size int
    	window size (default 700)
  -workers int
    	number of threads (default 8)
  -zoom float
    	zoom amount per frame (range from 0 to 1) (default 0.9)
```