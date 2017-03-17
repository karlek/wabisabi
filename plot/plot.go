package plot

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"reflect"
	"runtime"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/karlek/wabisabi/fractal"
	"github.com/karlek/wabisabi/histo"
)

var (
	f          = Exp
	importance histo.Histo
)

func Trap(img *image.RGBA, trapPath string, r, g, b histo.Histo) {

}

func PlotImp(width, height int, filePng, fileJpg bool) (err error) {
	fscale := func(v, max float64) float64 {
		return math.Min(((1 - math.Exp(-1*v)) / (1 - math.Exp(-1*max)) * 255), 255)
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	impMax := histo.Max(importance)
	for x, col := range importance {
		for y, v := range col {
			if importance[x][y] == 0 {
				continue
			}
			c := uint8(fscale(v, impMax))
			img.SetRGBA(y, x, color.RGBA{c, c, c, 255})
		}
	}
	return Render(img, filePng, fileJpg, "importance")
}

// Plot visualizes the histograms values as an image. It equalizes the
// histograms with a color scaling function to emphazise hidden features.
func Plot(img *image.RGBA, factor, exposure float64, frac *fractal.Fractal) {
	// The highest number orbits passing through a point.
	rMax, gMax, bMax := histo.Max(frac.R), histo.Max(frac.G), histo.Max(frac.B)
	// We iterate over every point in our histogram to color scale and plot
	// them.
	wg := new(sync.WaitGroup)
	wg.Add(len(frac.R))
	for x, col := range frac.R {
		go plotCol(wg, x, col, img, frac, factor, exposure, rMax, bMax, gMax)
	}
	wg.Wait()
}

// plotCol plots a column of pixels. The RGB-value of the pixel is based on the
// frequency in the histogram. Higher value equals brighter color.
func plotCol(wg *sync.WaitGroup, x int, col []float64, img *image.RGBA, frac *fractal.Fractal, factor, exposure, rMax, bMax, gMax float64) {
	for y := range col {
		// We skip to plot the black points for faster fracdering. A side
		// effect is that rendering png images will have a transpafract
		// background.
		if frac.R[x][y] == 0 &&
			frac.G[x][y] == 0 &&
			frac.B[x][y] == 0 {
			continue
		}
		c := color.RGBA{
			uint8(value(frac.R[x][y], rMax, factor, exposure)),
			uint8(value(frac.G[x][y], gMax, factor, exposure)),
			uint8(value(frac.B[x][y], bMax, factor, exposure)),
			255}
		// We flip x <=> y to rotate the image to an upright position.
		img.SetRGBA(y, x, c)
	}
	wg.Done()
}

// Exp is an exponential color scaling function.
func Exp(x, factor float64) float64 {
	return (1 - math.Exp(-factor*x))
}

// Log is an logaritmic color scaling function.
func Log(x, factor float64) float64 {
	return math.Log1p(factor * x)
}

// Sqrt is a square root color scaling function.
func Sqrt(x, factor float64) float64 {
	return math.Sqrt(factor * x)
}

// Lin is a linear color scaling function.
func Lin(x, factor float64) float64 {
	return x
}

// value calculates the color value of the pixel.
func value(v, max, factor, exposure float64) float64 {
	return math.Min(f(v, factor)*scale(max, factor, exposure), 255)
}

// scale equalizes the histogram distribution for each value.
func scale(max, factor, exposure float64) float64 {
	return (255 * exposure) / f(max, factor)
}

// Render creates an output image file.
func Render(img image.Image, filePng, fileJpg bool, filename string) (err error) {
	enc := func(img image.Image, filename string) (err error) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if filePng {
			return png.Encode(file, img)
		}
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
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

// getFunctionName returns the name of a function as string.
func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}