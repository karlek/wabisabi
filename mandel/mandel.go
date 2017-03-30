package mandel

import (
	"image"

	"github.com/karlek/wabisabi/fractal"
	"github.com/karlek/wabisabi/histo"
)

// Credits: https://github.com/morcmarc/buddhabrot/blob/master/buddhabrot.go
func isInBulb(c complex128) bool {
	Cr, Ci := real(c), imag(c)
	// Main cardioid
	if !(((Cr-0.25)*(Cr-0.25)+(Ci*Ci))*(((Cr-0.25)*(Cr-0.25)+(Ci*Ci))+(Cr-0.25)) < 0.25*Ci*Ci) {
		// 2nd order period bulb
		if !((Cr+1.0)*(Cr+1.0)+(Ci*Ci) < 0.0625) {
			// smaller bulb left of the period-2 bulb
			if !((((Cr + 1.309) * (Cr + 1.309)) + Ci*Ci) < 0.00345) {
				// smaller bulb bottom of the main cardioid
				if !((((Cr + 0.125) * (Cr + 0.125)) + (Ci-0.744)*(Ci-0.744)) < 0.0088) {
					// smaller bulb top of the main cardioid
					if !((((Cr + 0.125) * (Cr + 0.125)) + (Ci+0.744)*(Ci+0.744)) < 0.0088) {
						return false
					}
				}
			}
		}
	}
	return true
}

// escaped returns all points in the domain of the complex function before
// diverging.
// func Escaped(plane func(complex128, complex128) complex128, c, coefficient complex128, points []complex128, iterations int64, bailout float64, width, height int, r, g, b histo.Histo) {
func Escaped(z, c complex128, points []complex128, frac *fractal.Fractal) {
	track(z, c, points, registerOrbits, frac)
}

func CalculationPath(z, c complex128, points []complex128, frac *fractal.Fractal) {
	track(z, c, points, registerPaths, frac)
}

func abs(c complex128) complex128 {
	// return complex(real(c), -imag(c))
	// return complex(math.Abs(real(c)), math.Abs(imag(c)))
	return complex(real(c)/imag(c), real(c))
	// return complex(real(c)*imag(c), -imag(c))
	// return complex(-imag(c), -real(c))
	// return complex(imag(c), real(c))
	// return complex(imag(c), real(c))
	// return complex(imag(c), imag(c))
	// return complex(real(c), real(c))
	// return complex(math.Abs(real(c)), math.Abs(imag(c)))
}

func track(z, c complex128, points []complex128, f func(int64, []complex128, *fractal.Fractal), frac *fractal.Fractal) {
	// We ignore all values that we know are in the bulb, and will therefore
	// converge.
	if isInBulb(c) {
		return
	}

	// Saved value for cycle-detection.
	var bfract complex128

	// Number of points that we will return.
	var num int

	// See if the complex function diverges before we reach our iteration count.
	var i int64
	for i = 0; i < frac.Iterations; i++ {
		// z = frac.Coef*complex(real(z), imag(z))*complex(real(z), imag(z)) + frac.Coef*complex(real(c), imag(c))
		// z = complex(math.Sin(real(z))*imag(z), imag(z))*complex(math.Sin(real(z))*imag(z), imag(z)) + c
		z = z*z + c
		// z1 := complex(math.Sin(real(z))*imag(z), imag(z))
		// z2 := complex(imag(z), math.Cos(math.Exp(real(z))))
		// z = frac.Coef*z1*z2 + c
		// z = z + z*c + c
		// z = abs(z*z) + c
		// z = z + c
		// z = cmplx.Pow(c, z) + c
		// z = cmplx.Sqrt(z) + cmplx.Pow(c, c)
		// z = z*c
		// z = cmplx.Sqrt(z) + cmplx.Pow(c, c)
		// z = z*c + cmplx.Sin(z+c)
		// z = z*c + cmplx.Cos(z+c)
		// z = z*z/c + 1
		// z = (z*z)/c + c*c
		// z = z*z + c
		// z = z*c + (z+1)/c
		// z = z * c
		// z = z * cmplx.Log(c)
		// z = cmplx.Sin(z) + 1 + z*cmplx.Log(c)
		// z = abs(z*z) + c
		// complex(-math.Abs(real(c)), imag(c))
		// complex(math.Abs(real(c)), imag(c))
		// complex(imag(c)-real(c), real(c)*imag(c))

		// Cycle-detection (See algorithmic explanation in README.md).
		if (i-1)&i == 0 && i > 1 {
			bfract = z
		} else if z == bfract {
			return
		}
		// This point diverges, so we all the preceeding points are interesting
		// and will be registered.
		if x, y := real(z), imag(z); x*x+y*y >= frac.Bailout {
			f(i, points, frac)
			return
		}

		points[num] = frac.Plane(z, c)
		num++
	}
	// This point converges; assumed under the number of iterations.
	return
}

// Converged returns all points in the domain of the complex function before
// diverging.
func Converged(z, c complex128, points []complex128, frac *fractal.Fractal) {
	converged(z, c, points, registerOrbits, frac)
}
func converged(z, c complex128, points []complex128, f func(int64, []complex128, *fractal.Fractal), frac *fractal.Fractal) {
	// Saved value for cycle-detection.
	var bfract complex128

	// Number of points that we will return.
	var num int

	// See if the complex function diverges before we reach our iteration count.
	var i int64
	for i = 0; i < frac.Iterations; i++ {
		z = frac.Coef*complex(real(z), imag(z))*complex(real(z), imag(z)) + frac.Coef*complex(real(c), imag(c))

		// Cycle-detection (See algorithmic explanation in README.md).
		if (i-1)&i == 0 && i > 1 {
			bfract = z
		} else if z == bfract {
			f(i, points, frac)
			return
		}
		// This point diverges. Since it's the anti-buddhabrot, we are not
		// interested in these points.
		if x, y := real(z), imag(z); x*x+y*y >= frac.Bailout {
			return
		}

		points[num] = frac.Plane(z, c)
		num++
	}
	// This point converges; assumed under the number of iterations. Since it's
	// the anti-buddhabrot we register the orbit.
	// registerOrbits(points, width, height, num, iterations, r, g, b)
	return
}

// Primitive returns all points in the domain of the complex function
// diverging.
func Primitive(z, c complex128, points []complex128, frac *fractal.Fractal) {
	primitive(z, c, points, registerOrbits, frac)
}
func primitive(z, c complex128, points []complex128, f func(int64, []complex128, *fractal.Fractal), frac *fractal.Fractal) {
	// Saved value for cycle-detection.
	var bfract complex128

	// Number of points that we will return.
	var num int

	// See if the complex function diverges before we reach our iteration count.
	var i int64
	for i = 0; i < frac.Iterations; i++ {
		z = frac.Coef*complex(real(z), imag(z))*complex(real(z), imag(z)) + frac.Coef*complex(real(c), imag(c))

		// Cycle-detection (See algorithmic explanation in README.md).
		if (i-1)&i == 0 && i > 1 {
			bfract = z
		} else if z == bfract {
			f(i, points, frac)
			return
		}
		// This point diverges. Since it's the primitive brot we register the
		// orbit.
		if x, y := real(z), imag(z); x*x+y*y >= frac.Bailout {
			f(i, points, frac)
			return
		}
		// Save the point.
		points[num] = frac.Plane(z, c)
		num++
	}
	// This point converges; assumed under the number of iterations.
	// Since it's the primitive brot we register the orbit.
	f(i, points, frac)
	return
}

func Bresenham(start, end image.Point, points []image.Point) []image.Point {
	// Bresenham's
	var cx int = start.X
	var cy int = start.Y

	var dx int = end.X - cx
	var dy int = end.Y - cy
	if dx < 0 {
		dx = 0 - dx
	}
	if dy < 0 {
		dy = 0 - dy
	}

	var sx int
	var sy int
	if cx < end.X {
		sx = 1
	} else {
		sx = -1
	}
	if cy < end.Y {
		sy = 1
	} else {
		sy = -1
	}
	var err int = dx - dy

	var n int
	for n = 0; n < cap(points); n++ {
		points = append(points, image.Point{cx, cy})
		if cx == end.X && cy == end.Y {
			return points
		}
		var e2 int = 2 * err
		if e2 > (0 - dy) {
			err = err - dy
			cx = cx + sx
		}
		if e2 < dx {
			err = err + dx
			cy = cy + sy
		}
	}
	return points
}

// ptoc converts a point from the complex function to a pixel coordinate.
//
// Stands for point to coordinate, which is actually a really shitty name
// because of it's ambiguous character haha.
func ptoc(c complex128, frac *fractal.Fractal) (p image.Point) {
	r, i := real(c), imag(c)

	p.X = int((float64(frac.Width)/2.5)*frac.Zoom*(r+frac.OffsetReal) + float64(frac.Width)/2.0)
	p.Y = int((float64(frac.Height)/2.5)*frac.Zoom*(i+frac.OffsetImag) + float64(frac.Height)/2.0)

	return p
}

var importance = histo.Histo{}

var Max int64

// registerOrbits register the points in an orbit in r, g, b channels depending
// on it's iteration count.
func registerOrbits(it int64, points []complex128, frac *fractal.Fractal) {
	if it > Max {
		Max = it
	}

	// Get color from gradient based on iteration count of the orbit.
	red, green, blue := frac.Method.Get(it, frac.Iterations)
	for _, p := range points[:it] {
		registerPoint(p, it, frac, red, green, blue)
	}

	// if p, ok := pointImp(points[0], frac.Width, frac.Height); ok && imp != 0 {
	// 	importance[p.X][p.Y] += imp / float64(frac.Iterations)
	// }
}

func registerPoint(z complex128, it int64, frac *fractal.Fractal, red, green, blue float64) {
	if p, ok := point(z, frac); ok {
		increase(p, it, frac, red, green, blue)
	}
}

func point(z complex128, frac *fractal.Fractal) (image.Point, bool) {
	// Convert the complex point to a pixel coordinate.
	p := ptoc(z, frac)

	// Ignore points outside image.
	if p.X >= frac.Width || p.Y >= frac.Height || p.X < 0 || p.Y < 0 {
		return p, false
	}
	return p, true
}

func pointImp(z complex128, width, height int) (image.Point, bool) {
	var p image.Point
	// Convert the complex point to a pixel coordinate.
	r, i := real(z), imag(z)

	p.X = int((float64(width)/2.5)*(r+0.4) + float64(width)/2.0)
	p.Y = int((float64(height)/2.5)*i + float64(height)/2.0)

	// Ignore points outside image.
	if p.X >= width || p.Y >= height || p.X < 0 || p.Y < 0 {
		return p, false
	}
	return p, true
}

func increase(p image.Point, it int64, frac *fractal.Fractal, red, green, blue float64) {
	frac.R[p.X][p.Y] += red
	frac.G[p.X][p.Y] += green
	frac.B[p.X][p.Y] += blue
}

func registerPaths(it int64, points []complex128, frac *fractal.Fractal) {
	// Get color from gradient based on iteration count of the orbit.
	red, green, blue := frac.Method.Get(it, frac.Iterations)
	first := true
	var last image.Point
	bresPoints := make([]image.Point, 0, frac.Points)
	for _, z := range points[:it] {
		// Convert the complex point to a pixel coordinate.
		p, ok := point(z, frac)
		if !ok {
			continue
		}
		if first {
			first = false
			last = p
			continue
		}
		for _, prim := range Bresenham(last, p, bresPoints) {
			frac.R[prim.X][prim.Y] += red
			frac.G[prim.X][prim.Y] += green
			frac.B[prim.X][prim.Y] += blue
		}
		last = p
	}
}
