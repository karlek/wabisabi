package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"time"

	rand7i "github.com/7i/rand"
	"github.com/Sirupsen/logrus"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/karlek/profile"
	"github.com/karlek/progress/barcli"
	"github.com/karlek/wabisabi/coloring"
	"github.com/karlek/wabisabi/fractal"
	"github.com/karlek/wabisabi/mandel"
	"github.com/karlek/wabisabi/plot"
	"github.com/karlek/wabisabi/render"
)

func main() {
	defer profile.Start(profile.CPUProfile).Stop()
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	parseAdvancedFlags()
	parseFunctionFlag()

	if factor == -1 {
		factor = 1 / tries
	}

	// Handle interrupts as fails, so we can chain with an image viewer.
	inter := make(chan os.Signal, 1)
	signal.Notify(inter, os.Interrupt)
	go func(inter chan os.Signal) {
		<-inter
		os.Exit(1)
	}(inter)

	if err := buddha(); err != nil {
		logrus.Fatalln(err)
	}
}

func buddha() (err error) {
	// Create coloring scheme for the buddhabrot rendering.
	var grad coloring.Gradient
	// grad.AddColor(colorful.Color{0.02, 0.01, 0.01})
	// grad.AddColor(colorful.Color{0.02, 0.01, 0.02})
	// grad.AddColor(colorful.Color{0.0, 0.0, 0.0})
	// grad.AddColor(colorful.Color{0.0, 0.0, 0.1})
	// grad.AddColor(colorful.Color{0.1, 0.1, 0.1})
	// grad.AddColor(colorful.Color{0.3, 0.3, 0.3})
	// grad.AddColor(colorful.Color{0.00, 0.00, 0.00})

	// grad.AddColor(colorful.Color{0, 0.0, 0})
	// grad.AddColor(colorful.Color{0.11, 0.0, 0.08})
	// grad.AddColor(colorful.Color{0, 0.5, 1})
	// grad.AddColor(colorful.Color{1, 0.5, 0})
	// grad.AddColor(colorful.Color{1, 1, 1})

	grad.AddColor(colorful.Color{0, 1, 0})
	grad.AddColor(colorful.Color{0, 0, 1})
	grad.AddColor(colorful.Color{1, 0, 0})
	// grad.AddColor(colorful.Color{0, 0, 0})
	// grad.AddColor(colorful.Color{0, 0, 0})
	// grad.AddColor(colorful.Color{.65, 1, 0})
	// grad.AddColor(colorful.Color{1, 1, 1})

	ranges := []float64{
		// 20.0 / float64(iterations),
		// 200.0 / float64(iterations),
		// 200.0 / float64(iterations),
		// 1000.0 / float64(iterations),
		// 2000.0 / float64(iterations),
		// 2000.0 / float64(iterations),
		// 20000.0 / float64(iterations),
		// 0.0000001,
		// 0.000001,
		0.00001,
		0.0001,
		0.001,
		// 0.01,
		// 0.1,
		// 0.5,
	}
	// xor thing
	// orbit gradient
	// function for iteration
	method := coloring.NewColoring(mode, grad, ranges)

	logrus.Println("[.] Initializing.")
	var frac *fractal.Fractal
	var ren *render.Render
	// Load previous histograms and render the image with, maybe, new options.

	settings := logrus.Fields{
		"factor":     factor,
		"f":          getFunctionName(f),
		"out":        out,
		"load":       load,
		"save":       save,
		"anti":       anti,
		"brot":       getFunctionName(brot),
		"palette":    palettePath,
		"tries":      tries,
		"bailout":    bailout,
		"offset":     offset,
		"exposure":   exposure,
		"width":      width,
		"height":     height,
		"iterations": iterations,
		"zoom":       zoom,
		"seed":       seed,
	}
	logrus.WithFields(settings).Println("Config")

	ren = render.New(width, height, f, factor, exposure)
	if load {
		logrus.Println("[-] Loading visits.")
		frac, ren, err = loadArt()
		if err != nil {
			return err
		}
		fmt.Println(frac, ren)
	} else {
		// Fill our histogram bins of the orbits.
		frac = fractal.New(width, height, iterations, method, coefficient, bailout, plane, zoom, offsetReal, offsetImag, seed, intermediaryPoints)
		fmt.Println(frac)
		fillHistograms(frac, runtime.NumCPU())
		if save {
			logrus.Println("[i] Saving r, g, b channels")
			if err := saveArt(frac, ren); err != nil {
				return err
			}
		}
	}
	ren.Exposure = exposure
	ren.Factor = factor
	ren.F = f
	fmt.Println(ren)
	// fmt.Println(histo.Max(frac.R), histo.Max(frac.G), histo.Max(frac.B))
	// if histo.Max(frac.R)+histo.Max(frac.G)+histo.Max(frac.B) == 0 {
	// 	out += "-black"
	// 	return nil
	// }
	// fmt.Println("Longest orbit:", mandel.Max)

	// Plot and render to file.
	fmt.Println(ren.Factor)
	plot.Plot(ren, frac)
	plot.Render(ren.Image, filePng, fileJpg, out)
	sum := 0
	for _, k := range frac.Method.Keys {
		sum += k
	}
	for _, k := range frac.Method.Keys {
		fmt.Printf("%.5f\t", float64(k)/float64(sum))
	}
	fmt.Println()
	return nil
}

// fillHistograms creates a number of workers which finds orbits and stores
// their points in a histogram.
func fillHistograms(frac *fractal.Fractal, workers int) {
	bar, _ := barcli.New(int(tries * float64(width*height)))
	go func(bar *barcli.Bar) {
		for {
			if bar.Done() {
				return
			}
			bar.Print()
			time.Sleep(1000 * time.Millisecond)
		}
	}(bar)
	wg := new(sync.WaitGroup)
	wg.Add(workers)
	share := int(tries*float64(width*height)) / workers

	testPlot := make(chan OrbitPoint)
	go func(testPlot chan OrbitPoint, frac *fractal.Fractal) {
		// img := image.NewRGBA(image.Rect(0, 0, width, height))
		// grad := coloring.NewPrettyGradient(float64(frac.Iterations))
		// for op := range testPlot {
		// 	fmt.Println(op)
		// 	p := ptoc(op.P, frac)
		// 	// c := grad.DivergenceToColor(int(op.Len))

		// 	img.SetRGBA(p.X, p.Y, color.RGBA{0, 0, 0, 0})
		// }
	}(testPlot, frac)
	for n := 0; n < workers; n++ {
		// Our worker channel to send our orbits on!
		rng := rand7i.NewComplexRNG(int64(n+1) + seed)
		go arbitrary(testPlot, frac, &rng, share, wg, bar)
		// go iterative(frac, &rng, share, wg, bar)
	}
	wg.Wait()
	bar.SetMax()
	bar.Print()
}

// func toRGBA(c color.
type OrbitPoint struct {
	Len int64
	P   complex128
}

func ptoc(c complex128, frac *fractal.Fractal) (p image.Point) {
	r, i := real(c), imag(c)

	p.X = int((float64(frac.Width)/2.5)*frac.Zoom*(r+frac.OffsetReal) + float64(frac.Width)/2.0)
	p.Y = int((float64(frac.Height)/2.5)*frac.Zoom*(i+frac.OffsetImag) + float64(frac.Height)/2.0)

	return p
}

// arbitrary will try to find orbits in the complex function by choosing a
// random point in it's domain and iterating it a number of times to see if it
// converges or diverges.
func arbitrary(testPlot chan OrbitPoint, frac *fractal.Fractal, rng *rand7i.ComplexRNG, share int, wg *sync.WaitGroup, bar *barcli.Bar) {
	var potentials = make([]complex128, iterations)
	z := complex(0, 0)
	for i := 0; i < share; i++ {
		// Increase progress bar.
		bar.Inc()
		// Our random point which, hopefully, will create an orbit!

		// z = rng.Complex128Go()
		// z = complex(real(z), 0)
		// z = complex(0, imag(z))
		c := rng.Complex128Go()
		mandel.FieldLines(z, c, potentials, frac)

		// testPlot <- OrbitPoint{Len: brot(z, c, potentials, frac), P: c}
	}
	wg.Done()
}

func iterative(frac *fractal.Fractal, rng *rand7i.ComplexRNG, share int, wg *sync.WaitGroup, bar *barcli.Bar) {
	var potentials = make([]complex128, iterations)
	z := complex(0, 0)
	c := complex(0, 0)
	// h := (5000 / float64(share))
	h := 4 / math.Sqrt(float64(share))
	nudge := rand.Float64() * h
	for 100*nudge > h {
		nudge /= 10
	}
	h = h + nudge
	fmt.Println("test", h, nudge, h+nudge)
	var x, y float64
	var i int
	for y = -2; y <= 2; y += h {
		for x = -2; x <= 2; x += h {
			bar.Inc()
			// fmt.Println(c)
			c = complex(x, y)

			z = rng.Complex128Go()
			brot(z, c, potentials, frac)
			i++
		}
		// fmt.Println(y, float64(i)/float64(share))
	}
	// fmt.Println(i, share, h)
	wg.Done()
}

// getFunctionName returns the name of a function as string.
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func saveArt(frac *fractal.Fractal, ren *render.Render) (err error) {
	file, err := os.Create("r-g-b.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	gob.Register(colorful.Color{})
	err = enc.Encode(frac)
	if err != nil {
		return err
	}
	err = enc.Encode(ren)
	if err != nil {
		return err
	}
	return nil
}

func loadArt() (frac *fractal.Fractal, ren *render.Render, err error) {
	file, err := os.Open("r-g-b.gob")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	gob.Register(colorful.Color{})
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&frac); err != nil {
		return nil, nil, err
	}
	if err := dec.Decode(&ren); err != nil {
		return nil, nil, err
	}

	// Work around for function pointers and gobbing.
	switch ren.FName {
	case "github.com/karlek/wabisabi/plot.Log":
		ren.F = plot.Log
	case "github.com/karlek/wabisabi/plot.Exp":
		ren.F = plot.Exp
	case "github.com/karlek/wabisabi/plot.Lin":
		ren.F = plot.Lin
	case "github.com/karlek/wabisabi/plot.Sqrt":
		ren.F = plot.Sqrt
	}

	return frac, ren, nil
}
