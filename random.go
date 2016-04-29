package main

import (
	"io"
)

var p = make([]byte, 4)

func randfloat(random io.Reader) float64 {
	random.Read(p)
	b0, b1, b2, b3 := float64(p[0]), float64(p[1]), float64(p[2]), float64(p[3])
	return (1 / 256.0) * (b0 + (1/256.0)*(b1+(1/256.0)*(b2+(1/256.0)*b3)))
}
