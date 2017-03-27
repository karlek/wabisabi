package render

import "image"

type Render struct {
	Image    *image.RGBA
	Factor   float64
	Exposure float64
	Points   int
	F        func(float64, float64) float64
}

// New returns a new render for fractals.
func New(width, height int, f func(float64, float64) float64, factor, exposure float64) *Render {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return &Render{Image: img, Factor: factor, Exposure: exposure}
}
