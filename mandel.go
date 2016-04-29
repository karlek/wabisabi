package main

import (
	"io"
	"math"
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
func escaped(c complex128, points *[iterations]complex128, r, g, b *Histo) {
	// We ignore all values that we know are in the bulb, and will therefore
	// converge.
	if isInBulb(c) {
		return
	}

	// Saved value for cycle-detection.
	var brent complex128

	// Number of points that we will return.
	var num int

	// z is the point of the function.
	z := complex(0, 0)

	// See if the complex function diverges before we reach our iteration count.
	for i := 0; i < iterations; i++ {
		z = z*z + c
		/// z = complex(coefficient, 0)*complex(real(z), imag(z))*complex(real(z), imag(z)) + complex(coefficient, 0)*complex(real(c), imag(c))

		// Cycle-detection (See algorithmic explanation in README.md).
		if (i-1)&i == 0 && i > 1 {
			brent = z
		} else if z == brent {
			return
		}
		// This point diverges, so we all the preceeding points are interesting
		// and will be registered.
		if x, y := real(z), imag(z); x*x+y*y >= bailout {
			// Ignore lower iteration orbits to reduce noise.
			if num < 1000 {
				return
			}
			registerOrbits(points, num, r, g, b)
			return
		}

		/// points[num] = plane(z, c)
		points[num] = z
		num++
	}
	// This point converges; assumed under the number of iterations.
	return
}

// converged returns all points in the domain of the complex function before
// diverging.
func converged(c complex128, points *[iterations]complex128, r, g, b *Histo) {
	// Saved value for cycle-detection.
	var brent complex128

	// Number of points that we will return.
	var num int

	// z is the point of the function.
	z := complex(0, 0)

	// See if the complex function diverges before we reach our iteration count.
	for i := 0; i < iterations; i++ {
		z = z*z + c
		// Cycle-detection (See algorithmic explanation in README.md).
		if (i-1)&i == 0 && i > 1 {
			brent = z
		} else if z == brent {
			registerOrbits(points, num, r, g, b)
			return
		}
		// This point diverges, so we all the preceeding points are interesting
		// and will be registered.
		if x, y := real(z), imag(z); x*x+y*y >= bailout {
			return
		}

		points[num] = z
		num++
	}
	// This point converges; assumed under the number of iterations.
	registerOrbits(points, num, r, g, b)
	return
}

// registerOrbits register the points in an orbit in r,g,b channels depending on
// it's iteration count. Orbits with low iteration count will be ignored to
// reduce noise.
func registerOrbits(points *[iterations]complex128, it int, r, g, b *Histo) {
	red, green, blue := grad.Get(it % len(grad))
	for _, z := range points[:it] {
		p := ptoc(z)
		// Ignore points outside image.
		if p.X >= width || p.Y >= height || p.X < 0 || p.Y < 0 {
			continue
		}
		// Get color from gradient based on iteration count of the orbit.
		r[p.X][p.Y] += float64(red)
		g[p.X][p.Y] += float64(green)
		b[p.X][p.Y] += float64(blue)
	}
}

// primitive returns all points in the domain of the complex function
// diverging.
func primitive(random io.Reader, points *[iterations]complex128, c complex128) []complex128 {
	// Saved value for cycle-detection.
	var brent complex128

	// Number of points that we will return.
	var num int

	// z is the point of the function.
	z := complex(0, 0)

	// See if the complex function diverges before we reach our iteration count.
	for i := 0; i < iterations; i++ {
		z = z*z + c
		// Cycle-detection (See algorithmic explanation in README.md).
		if (i-1)&i == 0 && i > 1 {
			brent = z
		} else if z == brent {
			return points[:num]
		}
		// This point diverges, so we all the preceeding points are interesting
		// and will be returned.
		if x, y := real(z), imag(z); x*x+y*y >= bailout {
			return points[:num]
		}
		// Save the point.
		points[num] = plane(z, c)
		num++
	}
	// This point converges; assumed under the number of iterations.
	return points[:num]
}

func abs(c complex128) complex128 {
	return complex(math.Abs(real(c)), math.Abs(imag(c)))
}
