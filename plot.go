package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
)

func trap(img *image.RGBA, trapPath string, r, g, b *Histo) {

}

func plotImp() (err error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	impMax := max(&importance)
	for x, col := range importance {
		for y, v := range col {
			if importance[x][y] == 0 {
				continue
			}
			c := uint8(value(v, impMax))
			img.SetRGBA(x, y, color.RGBA{c, c, c, 255})
		}
	}
	return render(img, "importance")
}

// plot visualizes the histograms values as an image. It equalizes the
// histograms with a color scaling function to emphazise hidden features.
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

// plotCol plots a column of pixels. The RGB-value of the pixel is based on the
// frequency in the histogram. Higher value equals brighter color.
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

// exp is an exponential color scaling function.
func exp(x float64) float64 {
	return (1 - math.Exp(-factor*x))
}

// log is an logaritmic color scaling function.
func log(x float64) float64 {
	return math.Log1p(factor * x)
}

// sqrt is a square root color scaling function.
func sqrt(x float64) float64 {
	return math.Sqrt(factor * x)
}

// lin is a linear color scaling function.
func lin(x float64) float64 {
	return x
}

// value calculates the color value of the pixel.
func value(v, max float64) float64 {
	return math.Min(f(v)*scale(max), 255)
}

// scale equalizes the histogram distribution for each value.
func scale(max float64) float64 {
	return (255 * exposure) / f(max)
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
