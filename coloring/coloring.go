// Package coloring contains utility functions useable when drawing any fractal.
package coloring

import (
	"image"
	"image/color"
	"math/rand"
	"time"
)

// Fractal is a general fractal representation. Have information about:
// iterations, zoom value and centering point.
type Fractal struct {
	Src    *image.RGBA // Image to write fractal on.
	Iter   float64     // Number of iterations to perform.
	Center complex128  // Point to center the fractal on.
	Zoom   float64     // Zoom value.
}

// Gradient is a list of colors.
type Gradient []color.Color

var (
	// PedagogicalGradient have a fixed transformation between colors for easier
	// visualization of divergence.s
	PedagogicalGradient = Gradient{
		color.RGBA{0, 0, 0, 0xff},       // Black.
		color.RGBA{0xff, 0xf0, 0, 0xff}, // Yellow.
		color.RGBA{0, 0, 0xff, 0xff},    // Blue.
		color.RGBA{0, 0xff, 0, 0xff},    // Green.
		color.RGBA{0xff, 0, 0, 0xff},    // Red.
	}
)

// NewRandomGradient creates a gradient of colors proportional to the number of
// iterations.
func NewRandomGradient(iterations float64) Gradient {
	seed := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	grad := make(Gradient, int64(iterations))
	for n := range grad {
		grad[n] = randomColor(seed)
	}
	return grad
}

// randomColor returns a random RGB color from a random seed.
func randomColor(seed *rand.Rand) color.RGBA {
	return color.RGBA{
		uint8(seed.Intn(255)),
		uint8(seed.Intn(255)),
		uint8(seed.Intn(255)),
		0xff} // No alpha.
}

// NewPrettyGradient creates a gradient of colors fading between purple and
// white. The smoothness is proportional to the number of iterations
func NewPrettyGradient(iterations float64) Gradient {
	grad := make(Gradient, int64(iterations))
	var col color.Color
	for n := range grad {
		// val ranges from [0..255]
		val := uint8(float64(n) / float64(iterations) * 255)
		if int64(n) < int64(iterations/2) {
			col = color.RGBA{val * 2, 0x00, val * 2, 0xff} // Shade of purple.
		} else {
			col = color.RGBA{val, val, val, 0xff} // Shade of white.
		}
		grad[n] = col
	}
	return grad
}

// DivergenceToColor returns a color depending on the number of iterations it
// took for the fractal to escape the fractal set.
func (g Gradient) DivergenceToColor(escapedIn int) color.Color {
	return g[escapedIn%len(g)]
}

type Type int

const (
	Modulo Type = iota
	IterationCount
)

type Coloring struct {
	Grad   Gradient
	mode   Type
	Keys   []int
	Ranges []float64
}

// function pointer which makes get unneccessary
func NewColoring(mode Type, grad Gradient, ranges []float64) *Coloring {
	if len(grad) != len(ranges) {
		panic("number of colors and ranges mismatch")
	}
	keys := make([]int, len(ranges)+1)
	return &Coloring{Grad: grad, mode: mode, Ranges: ranges, Keys: keys}
}

func (c *Coloring) Get(i int64, it int64) (float64, float64, float64) {
	switch c.mode {
	case Modulo:
		return c.modulo(i)
	case IterationCount:
		return c.iteration(i, it)
	default:
		return c.modulo(i)
	}
}

func (c *Coloring) modulo(i int64) (float64, float64, float64) {
	if i >= int64(len(c.Grad)) {
		i %= int64(len(c.Grad))
	}
	r, g, b, _ := (c.Grad)[i].RGBA()
	return float64(r>>8) / 256, float64(g>>8) / 256, float64(b>>8) / 256
}

func (c *Coloring) iteration(i int64, it int64) (float64, float64, float64) {
	key := -1
	for rID := len(c.Ranges) - 1; rID >= 0; rID-- {
		if float64(i)/float64(it) >= c.Ranges[rID] {
			// fmt.Printf("%.7f - %d\t\t%d\n", float64(i)/float64(it), i, it)
			key = rID
			break
		}
	}
	c.Keys[key+1]++
	if key == -1 {
		return 0, 0, 0
	}
	r, g, b, _ := (c.Grad)[key].RGBA()
	return float64(r>>8) / 256, float64(g>>8) / 256, float64(b>>8) / 256
}

func (g *Gradient) AddColor(c color.Color) {
	(*g) = append((*g), c)
}
