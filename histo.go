package main

import (
	"encoding/gob"
	"os"
)

// Histo is a histogram of buddhabrot divergent orbits.
type Histo [width][height]float64

// max finds the highest value in the histogram. Used for color scaling
// algorithms.
func max(v *Histo) (max float64) {
	max = -1
	for _, row := range v {
		for _, v := range row {
			if v > max {
				max = v
			}
		}
	}
	return max
}

// gobHisto saves histograms to a gob file for future re-rendering.
func gobHisto(vs ...*Histo) (err error) {
	file, err := os.Create("r-g-b.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	for _, v := range vs {
		err = enc.Encode(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadHisto loads a previously calculated histogram file for re-rendering.
func loadHisto() (r, g, b *Histo, err error) {
	file, err := os.Open("r-g-b.gob")
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&r); err != nil {
		return nil, nil, nil, err
	}
	if err := dec.Decode(&g); err != nil {
		return nil, nil, nil, err
	}
	if err := dec.Decode(&b); err != nil {
		return nil, nil, nil, err
	}
	return r, g, b, nil
}
