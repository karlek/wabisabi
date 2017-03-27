package render

import (
	"image"
	"reflect"
	"runtime"
)

type Render struct {
	Image    *image.RGBA
	Factor   float64
	Exposure float64
	Points   int
	F        func(float64, float64) float64
	FName    string
}

// New returns a new render for fractals.
func New(width, height int, f func(float64, float64) float64, factor, exposure float64) *Render {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return &Render{Image: img, F: f, FName: getFunctionName(f), Factor: factor, Exposure: exposure}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
