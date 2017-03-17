package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
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
	"github.com/karlek/wabisabi/histo"
	"github.com/karlek/wabisabi/mandel"
	"github.com/karlek/wabisabi/plot"
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
	// grad.AddColor(colorful.Color{0.11, 0.07, 0.08})
	// grad.AddColor(colorful.Color{0.1, 0.1, 0.1})
	// grad.AddColor(colorful.Color{0.3, 0.3, 0.3})
	grad.AddColor(colorful.Color{1, 0, 0})
	grad.AddColor(colorful.Color{0, 1, 0})
	grad.AddColor(colorful.Color{0, 0, 1})
	grad.AddColor(colorful.Color{1, 1, 1})
	ranges := []float64{
		50.0 / float64(iterations),
		100.0 / float64(iterations),
		500.0 / float64(iterations),
		2000.0 / float64(iterations),
		// 0.005,
		// 0.01,
		// 0.02,
		// 0.5,
		// 0.1,
		// 0.2,
		// 0.3,
		// 0.4,
	}
	method := coloring.NewColoring(coloring.IterationCount, grad, ranges)

	logrus.Println("[.] Initializing.")
	frac := fractal.New(width, height, iterations, method, coefficient, bailout, plane, zoom, offsetReal, offsetImag)
	img := image.NewRGBA(image.Rect(0, 0, frac.Width, frac.Height))
	// Load previous histograms and render the image with, maybe, new options.
	if load {
		logrus.Println("[-] Loading visits.")
		frac, err := loadRen()
		if err != nil {
			return err
		}
		plot.Plot(img, factor, exposure, frac)
		plot.Render(img, false, true, "/tmp/a")
		fmt.Println(histo.Max(frac.R), histo.Max(frac.G), histo.Max(frac.B))
		return nil
	}

	// Fill our histogram bins of the orbits.
	fillHistograms(frac, img, runtime.NumCPU())

	if save {
		logrus.Println("[i] Saving r, g, b channels")
		if err := saveRen(frac); err != nil {
			return err
		}
	}
	fmt.Println(histo.Max(frac.R), histo.Max(frac.G), histo.Max(frac.B))
	fmt.Println("Longest orbit:", mandel.Max)
	if histo.Max(frac.R)+histo.Max(frac.G)+histo.Max(frac.B) == 0 {
		out += "-black"
		return nil
	}
	// Plot and render to file.
	plot.Plot(img, factor, exposure, frac)
	plot.Render(img, false, true, out)
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
func fillHistograms(frac *fractal.Fractal, img *image.RGBA, workers int) {
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
	for n := 0; n < workers; n++ {
		// Our worker channel to send our orbits on!
		rng := rand7i.NewComplexRNG(int64(n+1) + seed)
		go arbitrary(frac, &rng, share, wg, img, bar)
	}
	wg.Wait()
	bar.SetMax()
	bar.Print()
}

// arbitrary will try to find orbits in the complex function by choosing a
// random point in it's domain and iterating it a number of times to see if it
// converges or diverges.
func arbitrary(frac *fractal.Fractal, rng *rand7i.ComplexRNG, share int, wg *sync.WaitGroup, img *image.RGBA, bar *barcli.Bar) {
	var potentials = make([]complex128, iterations)
	for i := 0; i < share; i++ {
		// Increase progress bar.
		bar.Inc()
		// Our random point which, hopefully, will create an orbit!
		c := rng.Complex128Go()
		mandel.Escaped(c, potentials, frac)
	}
	wg.Done()
}

// getFunctionName returns the name of a function as string.
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func saveRen(ren *fractal.Fractal) (err error) {
	file, err := os.Create("r-g-b.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	gob.Register(colorful.Color{})
	err = enc.Encode(ren)
	if err != nil {
		return err
	}
	return nil
}

func loadRen() (ren *fractal.Fractal, err error) {
	file, err := os.Open("r-g-b.gob")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	gob.Register(colorful.Color{})
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&ren); err != nil {
		return nil, err
	}
	return ren, nil
}
