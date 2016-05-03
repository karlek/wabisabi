// buddha renders buddhabrot fractals.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"time"

	rand7i "github.com/7i/rand"
	"github.com/Sirupsen/logrus"
	"github.com/karlek/profile"
	"github.com/karlek/progress/barcli"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	// Width of the image.
	width = 4000
	// Height of the image.
	height = 4000
	// Number of iterations to compute before we assume that we converge.
	iterations = 100000
	wRatio     = float64(width) / float64(height)
	hRatio     = float64(height) / float64(width)
)

// stealPalette uses the dominant colors of an image as our gradient.
func stealPalette(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		logrus.Fatalln(err)
	}
	pal := createPalette(img, colors)
	for _, c := range pal {
		grad.AddColor(c.Col)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	// Seed our random seed.
	rand.Seed(seed)
	logrus.Println(seed)

	// Parse our non-trivial function flags.
	parseAdvancedFlags()

	// If we haven't specified a value for _factor_, scale it according to the
	// number of tries.
	if fmt.Sprint(flag.Lookup("factor").Value) == fmt.Sprint(flag.Lookup("factor").DefValue) {
		factor = 1.0 / (tries * 10)
	}
	// Handle interrupts as fails, so we can chain with an image viewer.
	inter := make(chan os.Signal, 1)
	signal.Notify(inter, os.Interrupt)
	go func(inter chan os.Signal) {
		<-inter
		os.Exit(1)
	}(inter)

	// Start our profiling.
	defer profile.Start(profile.All).Stop()

	// Render the (anti-)buddhabrot.
	if err := buddha(); err != nil {
		logrus.Fatalln(err)
	}
}

// Initialize allocates memory for our image and histograms.
func initialize() (img *image.RGBA, r, g, b *Histo) {
	// Output image with black background.
	return image.NewRGBA(image.Rect(0, 0, width, height)), &Histo{}, &Histo{}, &Histo{}
}

var bar *barcli.Bar

func buddha() (err error) {
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
	}
	logrus.WithFields(settings).Println("Config")

	logrus.Println("[.] Initializing.")
	// Initializing our image and arrays in memory.
	img, r, g, b := initialize()

	// Load previous histograms and render the image with, maybe, new options.
	if load {
		logrus.Println("[-] Loading visits.")
		r, g, b, err = loadHisto()
		if err != nil {
			return err
		}
		plot(img, r, g, b)
		return render(img, out)
	}

	// Default case uses a random gradient to plot the brot.
	//
	// If a path to an image file is supplied, the dominant colors from that
	// image are extracted.
	if palettePath == "" {
		logrus.Println("[.] Using random gradient.")
		for i := 0; i < colors; i++ {
			grad.AddColor(colorful.Color{
				rand.Float64(),
				rand.Float64(),
				rand.Float64()})
		}
	} else {
		logrus.Println("[.] Stealing palette from image.")
		stealPalette(palettePath)
	}

	logrus.Println("[-] Calculating visited points.")
	fillHistograms(r, g, b, runtime.NumCPU())

	// Saving the histograms for future re-rendering.
	if save {
		logrus.Println("[i] Saving r, g, b channels")
		if err := gobHisto(r, g, b); err != nil {
			return err
		}
	}

	// Color scale the histograms and plot to the brot.
	logrus.Println("[/] Creating image.")
	if trapPath == "" {
		plot(img, r, g, b)
	} else {
		// TODO(_): Write orbit trap functionality.
		trap(img, trapPath, r, g, b)
	}
	// Render the image to a file, can be either jpg or png.
	return render(img, out)
}

// fillHistograms creates a number of workers which finds orbits and stores
// their points in a histogram.
func fillHistograms(r, g, b *Histo, workers int) {
	bar, _ = barcli.New(int(tries * width * height))
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
	share := int((tries)*width*height) / workers
	for n := 0; n < workers; n++ {
		// Our worker channel to send our orbits on!
		rng := rand7i.NewComplexRNG(int64(n + 1))
		go arbitrary(r, g, b, &rng, share, wg)
	}
	wg.Wait()
	bar.SetMax()
	bar.Print()
}

// arbitrary will try to find orbits in the complex function by choosing a
// random point in it's domain and iterating it a number of times to see if it
// converges or diverges.
func arbitrary(r, g, b *Histo, rng *rand7i.ComplexRNG, random io.Reader, tries int, wg *sync.WaitGroup) {
	var potentials [iterations]complex128
	for i := 0; i < tries; i++ {
		// Increase progress bar.
		bar.Inc()

		// Our random point which, hopefully, will create an orbit!
		c := randomPoint(rng, random)
		brot(c, &potentials, r, g, b)
	}
	wg.Done()
}

func randomPoint(rng *rand7i.ComplexRNG, random io.Reader) complex128 {
	// return complex(4*randfloat(random)-2, 4*randfloat(random)-2)
	return rng.Complex128Go()
}

// getFunctionName returns the name of a function as string.
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func trap(img *image.RGBA, trapPath string, r, g, b *Histo) {

}

func plot(img *image.RGBA, r, g, b *Histo) {
	// The highest number orbits passing through a point.
	rMax, gMax, bMax := max(r), max(g), max(b)
	logrus.Println("[i] Histo:", rMax, gMax, bMax)
	logrus.Printf("[i] Function: %s, factor: %.2f, exposure: %.2f", getFunctionName(f), factor, exposure)
	// We iterate over every point in our histogram to color scale and plot
	// them.
	wg := new(sync.WaitGroup)
	wg.Add(len(r))
	for x, col := range r {
		go plotCol(wg, x, &col, img, r, g, b, rMax, bMax, gMax)
	}
	wg.Wait()
}

func plotCol(wg *sync.WaitGroup, x int, col *[height]float64, img *image.RGBA, r, g, b *Histo, rMax, bMax, gMax float64) {
	for y := range col {
		// We skip to plot the black points for faster rendering. A side
		// effect is that rendering png images will have a transparent
		// background.
		if r[x][y] == 0 &&
			g[x][y] == 0 &&
			b[x][y] == 0 {
			continue
		}
		c := color.RGBA{
			uint8(value(r[x][y], rMax)),
			uint8(value(g[x][y], gMax)),
			uint8(value(b[x][y], bMax)),
			255}
		// We flip x <=> y to rotate the image to an upright position.
		img.Set(y, x, c)
	}
	wg.Done()
}

func exp(x float64) float64 {
	return (1 - math.Exp(-factor*x))
}
func log(x float64) float64 {
	return math.Log1p(factor * x)
}
func sqrt(x float64) float64 {
	return math.Sqrt(factor * x)
}
func lin(x float64) float64 {
	return x
}
func value(v, max float64) float64 {
	return math.Min(f(v)*scale(max), 255)
}
func scale(max float64) float64 {
	return (255 * exposure) / f(max)
}

// ptoc converts a point from the complex function to a pixel coordinate.
//
// Stands for point to coordinate, which is actually a really shitty name
// because of it's ambiguous character haha.
func ptoc(c complex128) (p image.Point) {
	r, i := real(c), imag(c)

	p.X = int((zoom*float64(width)/2.8)*(r+real(offset)) + float64(width)/2.0)
	p.Y = int((zoom*float64(height)/2.8)*(i+imag(offset)) + float64(height)/2.0)

	return p
}

// render creates an output image file.
func render(img image.Image, filename string) (err error) {
	enc := func(img image.Image, filename string) (err error) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if filePng {
			return png.Encode(file, img)
		}
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 76})
	}

	if filePng {
		filename += ".png"
	} else if fileJpg {
		filename += ".jpg"
	}
	logrus.Println("[!] Encoding:", filename)
	defer logrus.Println("[!] Done :D")
	return enc(img, filename)
}
