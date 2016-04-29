# Wabisabi

Wabisabi is a renderer of buddhabrot and it's family members. It shares it's name with a Japanese asthethic called [Wabi-sabi](https://en.wikipedia.org/wiki/Wabi-sabi). Referencing the impossibility of creating the real buddhabrot and learning to accept the beauty in reality and it's flaws. 

[!Anti-buddhabrot](https://github.com/karlek/wabisabi/raw/master/img/anti.jpg)
[!Original buddhabrot](https://github.com/karlek/wabisabi/raw/master/img/original.jpg)

__Features__

* Saving and loading of histograms to re-render with different exposures.
* Calculating the original, anti- and primitive- buddhabrot.
* Exploring the different planes of Zr, Zi, Cr and Ci.
* Different histogram equalization functions (think color scaling).
* Using the color palette of an image to color the brot.
* Change the co-efficient of the complex function i.e __a__\*z\*z+__a__\*c
* Zooming.
* Multiple CPU support. 
* Hand optimized assembly(!) for generating random complex points. Thank you [7i](https://github.com/7i)

[!Benchmark](https://github.com/karlek/wabisabi/raw/master/img/benchmark.png)

__Future features__

* Metropolis-hastings algorithm for faster zooming.
* Orbit trapping; would be amazing!

## Install

```fish
$ go get github.com/karlek/wabisabi
```

## Run

```fish
# Be sure to limit the memory usage beforehand; wabisabi is greedy little devil.
$ ulimit -Sv 4000000 # Where the number is the memory in kB.
$ wabisabi
```

## Complex functions

```go
z = |z*z| + c
complex(real(c), -imag(c))
complex(-math.Abs(real(c)), imag(c))
complex(math.Abs(real(c)), imag(c))
```

## Z<sub>0</sub>

```go
z := randomPoint(random)
z := complex(math.Sin(real(c)), math.Sin(imag(c)))
```

## Future

* Only allow a certain type of orbits. 
    - Convex hull to check roundness?
    - Constant increment on certain axis indicates spirals?
    - Iteration length is not connected to orbits type?
    - How many orbit types are there?
    - Find more ways to discern different kinds of orbits. 
* Downsizing
    - Currently this feature is not supported in wabisabi, but _imagemagicks_ `convert` command supports resizing: `convert a.jpg -resize 25600x25600 b.jpg` 
* Super sampling
    - Actually not sure how this differs from rendering a larger buddhabrot and just downsizing it?
        + Probably is just skipping the render and resizing step and calculating the values in the histograms directly.
* Since the orbits reminds me of a circle; it could be possible to unravel the circle and convert them into sine-waves to create tones :D
    - Outer convex hull -> Radius (max, min (amplitude)) 
* Test slices instead of fixed size arrays for runtime allocation of iterations and width/height.
* More than 3 histograms?
    - Only makes sense with color spaces with more than 3 values such as cmyk?

## Fun stuff

Interesting bug:

```go
p.X = int((zoom*float64(width)/2.8)*(r+real(offset))) + width/2
p.Y = int((zoom*float64(height)/2.8)*(i+imag(offset))) + height/2
```

Fix
```go
p.X = int((zoom*float64(width)/2.8)*(r+real(offset)) + width/2)
p.Y = int((zoom*float64(height)/2.8)*(i+imag(offset)) + height/2)
```

Created crosses by rounding coordinates numbers.


