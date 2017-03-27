package fractal

import (
	"github.com/karlek/wabisabi/coloring"
	"github.com/karlek/wabisabi/histo"
	"github.com/lucasb-eyer/go-colorful"
)

type Fractal struct {
	Width, Height int
	R, G, B       histo.Histo
	Method        *coloring.Coloring
	// Fractal specific.
	Iterations int64
	Plane      func(complex128, complex128) complex128
	Coef       complex128
	Bailout    float64
	Zoom       float64
	OffsetReal float64
	OffsetImag float64
	Seed       int64
}

// New returns a new render for fractals.
func New(width, height int, iterations int64, method *coloring.Coloring, coef complex128, bailout float64, plane func(complex128, complex128) complex128, zoom, offsetReal, offsetImag float64, seed int64) *Fractal {
	r, g, b := histo.New(width, height), histo.New(width, height), histo.New(width, height)
	return &Fractal{Width: width, Height: height, Iterations: iterations, R: r, G: g, B: b, Method: method, Coef: coef, Bailout: bailout, Plane: plane, Zoom: zoom, OffsetReal: offsetReal, OffsetImag: offsetImag, Seed: seed}
}

func NewStd() *Fractal {
	var grad coloring.Gradient
	grad.AddColor(colorful.Color{1, 0, 0})
	grad.AddColor(colorful.Color{0, 1, 0})
	grad.AddColor(colorful.Color{0, 0, 1})
	method := coloring.NewColoring(coloring.IterationCount, grad, []float64{100 / 1e6, 0.2, 0.5})
	return New(4096, 4096, 1e6, method, complex(1, 0), 4, zrzi, 1, 0.4, 0, 1)
}
func zrzi(z, c complex128) complex128 { return complex(real(z), imag(z)) }
func zrcr(z, c complex128) complex128 { return complex(real(z), real(c)) }
func zrci(z, c complex128) complex128 { return complex(real(z), imag(c)) }
func crci(z, c complex128) complex128 { return complex(real(c), imag(c)) }
func crzi(z, c complex128) complex128 { return complex(real(c), imag(z)) }
func zici(z, c complex128) complex128 { return complex(imag(z), imag(c)) }
